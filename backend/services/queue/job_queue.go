package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/slideitin/backend/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

// FirestoreJob is the Firestore representation of a job
// Simplified to contain only essential fields
type FirestoreJob struct {
	ID        string `firestore:"id"`
	Status    string `firestore:"status"`
	Message   string `firestore:"message"`
	CreatedAt int64  `firestore:"createdAt"`
	UpdatedAt int64  `firestore:"updatedAt"`
	ExpiresAt int64  `firestore:"expiresAt,omitempty"`
}

// FirestoreResult is the Firestore representation of a job result
type FirestoreResult struct {
	ID          string `firestore:"id"`
	ResultURL   string `firestore:"resultUrl"`
	PDFData     []byte `firestore:"pdfData"`
	CreatedAt   int64  `firestore:"createdAt"`
	ExpiresAt   int64  `firestore:"expiresAt"`
}

// Job represents a single slide generation job with runtime features
type Job struct {
	ID        string
	Theme     string
	Files     []struct {
		Filename string
		Data     []byte
		Type     string
	}
	Settings  models.SlideSettings
	Status    JobStatus
	Message   string
	ResultURL string
	CreatedAt int64
	UpdatedAt int64
}

// JobUpdate represents an update to a job that can be sent to SSE clients
type JobUpdate struct {
	ID        string    `json:"id"`
	Status    JobStatus `json:"status"`
	Message   string    `json:"message"`
	ResultURL string    `json:"resultUrl,omitempty"`
	UpdatedAt int64     `json:"updatedAt"`
}

// Service manages jobs using Firestore
type Service struct {
	client    *firestore.Client
	mu        sync.RWMutex
	slideSvc interface {
		GenerateSlides(ctx context.Context, theme string, files []struct {
			Filename string
			Data     []byte
			Type     string
		}, settings models.SlideSettings, statusUpdateFn func(message string) error) ([]byte, error)
	}
}

// NewService creates a new queue service using Firestore
func NewService(client *firestore.Client, slideSvc interface {
	GenerateSlides(ctx context.Context, theme string, files []struct {
		Filename string
		Data     []byte
		Type     string
	}, settings models.SlideSettings, statusUpdateFn func(message string) error) ([]byte, error)
}) *Service {
	s := &Service{
		client:    client,
		slideSvc: slideSvc,
	}

	return s
}

// Collection returns the Firestore collection reference for jobs
func (s *Service) Collection() *firestore.CollectionRef {
	return s.client.Collection("jobs")
}

// ResultsCollection returns the Firestore collection reference for results
func (s *Service) ResultsCollection() *firestore.CollectionRef {
	return s.client.Collection("results")
}

// AddJob adds a new job to Firestore and processes it immediately
func (s *Service) AddJob(ctx context.Context, id, theme string, fileData []struct {
	Filename string
	Data     []byte
	Type     string
}, settings models.SlideSettings) (*Job, error) {
	// Create the job
	now := time.Now().Unix()
	
	// Create a job record for Firestore (simplified)
	firestoreJob := FirestoreJob{
		ID:        id,
		Status:    string(StatusQueued),
		Message:   "Job added to queue",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Save to Firestore
	_, err := s.Collection().Doc(id).Set(ctx, firestoreJob)
	if err != nil {
		log.Printf("Failed to add job to Firestore: %v", err)
		return nil, fmt.Errorf("failed to store job: %v", err)
	}

	log.Printf("Added job %s to Firestore", id)

	// Create in-memory job object for processing
	job := &Job{
		ID:        id,
		Theme:     theme,
		Files:     fileData,
		Settings:  settings,
		Status:    StatusQueued,
		Message:   "Job added to queue",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Start processing in a goroutine
	go s.processJob(job)

	return job, nil
}

// GetJob retrieves a job by its ID from Firestore
func (s *Service) GetJob(id string) *Job {
	ctx := context.Background()
	doc, err := s.Collection().Doc(id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Printf("Job %s not found in Firestore", id)
			return nil
		}
		log.Printf("Error retrieving job %s: %v", id, err)
		return nil
	}

	var firestoreJob FirestoreJob
	if err := doc.DataTo(&firestoreJob); err != nil {
		log.Printf("Error parsing job data: %v", err)
		return nil
	}

	// Check if job has expired
	now := time.Now().Unix()
	if firestoreJob.ExpiresAt > 0 && now > firestoreJob.ExpiresAt {
		// Job has expired, delete it
		_, err := s.Collection().Doc(id).Delete(ctx)
		if err != nil {
			log.Printf("Failed to delete expired job %s: %v", id, err)
		} else {
			log.Printf("Deleted expired job %s", id)
		}
		return nil
	}

	// Get the result if available
	var resultURL string
	if firestoreJob.Status == string(StatusCompleted) {
		resultDoc, err := s.ResultsCollection().Doc(id).Get(ctx)
		if err == nil && resultDoc.Exists() {
			var result FirestoreResult
			if err := resultDoc.DataTo(&result); err == nil {
				resultURL = result.ResultURL
			}
		}
	}

	// Convert to job object
	return &Job{
		ID:        firestoreJob.ID,
		Status:    JobStatus(firestoreJob.Status),
		Message:   firestoreJob.Message,
		ResultURL: resultURL,
		CreatedAt: firestoreJob.CreatedAt,
		UpdatedAt: firestoreJob.UpdatedAt,
	}
}

// WatchJob watches a job for changes and sends updates to the provided channel
// This function will run until the context is canceled or the job reaches a terminal state
func (s *Service) WatchJob(ctx context.Context, jobID string, updates chan<- JobUpdate) error {
	// Get initial job state
	job := s.GetJob(jobID)
	if job == nil {
		return fmt.Errorf("job not found")
	}

	// Send initial status
	updates <- JobUpdate{
		ID:        job.ID,
		Status:    job.Status,
		Message:   job.Message,
		ResultURL: job.ResultURL,
		UpdatedAt: job.UpdatedAt,
	}

	// If job is already in terminal state, we're done
	if job.Status == StatusCompleted || job.Status == StatusFailed {
		close(updates)
		return nil
	}

	// Set up Firestore snapshot listener for real-time updates
	docRef := s.Collection().Doc(jobID)
	snapshots := docRef.Snapshots(ctx)

	// Watch for updates
	for {
		snapshot, err := snapshots.Next()
		if err != nil {
			log.Printf("Error watching job %s: %v", jobID, err)
			return err
		}

		if !snapshot.Exists() {
			log.Printf("Job %s no longer exists", jobID)
			return fmt.Errorf("job deleted")
		}

		var firestoreJob FirestoreJob
		if err := snapshot.DataTo(&firestoreJob); err != nil {
			log.Printf("Error parsing job data: %v", err)
			continue
		}

		// Get result URL if job is completed
		var resultURL string
		if firestoreJob.Status == string(StatusCompleted) {
			resultDoc, err := s.ResultsCollection().Doc(jobID).Get(ctx)
			if err == nil && resultDoc.Exists() {
				var result FirestoreResult
				if err := resultDoc.DataTo(&result); err == nil {
					resultURL = result.ResultURL
				}
			}
		}

		// Send update
		update := JobUpdate{
			ID:        firestoreJob.ID,
			Status:    JobStatus(firestoreJob.Status),
			Message:   firestoreJob.Message,
			ResultURL: resultURL,
			UpdatedAt: firestoreJob.UpdatedAt,
		}

		select {
		case updates <- update:
			// Successfully sent
		case <-ctx.Done():
			// Context was canceled
			return ctx.Err()
		}

		// If job is in terminal state, we're done
		if update.Status == StatusCompleted || update.Status == StatusFailed {
			return nil
		}
	}
}

// processJob processes a slide generation job
func (s *Service) processJob(job *Job) {
	// Update job status to processing in Firestore
	s.updateJobStatus(job, StatusProcessing, "Processing slides", "")

	// Create a status update function to pass to the Gemini service
	statusUpdateFn := func(message string) error {
		s.updateJobStatus(job, StatusProcessing, message, "")
		return nil
	}

	// Call the Gemini service with the status update function
	ctx := context.Background()
	pdfData, err := s.slideSvc.GenerateSlides(
		ctx, 
		job.Theme, 
		job.Files, 
		job.Settings,
		statusUpdateFn,
	)

	if err != nil {
		s.updateJobStatus(job, StatusFailed, "Failed to generate slides: "+err.Error(), "")
		return
	}

	// Create result URL
	resultURL := "/results/" + job.ID

	// Store result in Firestore
	s.storeResult(ctx, job.ID, resultURL, pdfData)

	// Update job as completed and set to expire in 5 minutes
	s.setJobCompleted(job, "Slides generated successfully", resultURL)
}

// setJobCompleted marks a job as completed and sets it to expire in 5 minutes
func (s *Service) setJobCompleted(job *Job, message, resultURL string) {
	ctx := context.Background()
	now := time.Now().Unix()
	// Set job to expire in 5 minutes
	expiresAt := now + 300 // 300 seconds = 5 minutes
	
	// Update job in Firestore
	updates := []firestore.Update{
		{Path: "status", Value: string(StatusCompleted)},
		{Path: "message", Value: message},
		{Path: "updatedAt", Value: now},
		{Path: "expiresAt", Value: expiresAt},
	}

	_, err := s.Collection().Doc(job.ID).Update(ctx, updates)
	if err != nil {
		log.Printf("Failed to update job status in Firestore: %v", err)
	}

	// Update the in-memory job
	job.Status = StatusCompleted
	job.Message = message
	job.UpdatedAt = now
	job.ResultURL = resultURL

	log.Printf("Job %s completed and will expire at %s", job.ID, time.Unix(expiresAt, 0).Format(time.RFC3339))
}

// storeResult stores a job result in Firestore
func (s *Service) storeResult(ctx context.Context, jobID, resultURL string, pdfData []byte) error {
	now := time.Now().Unix()
	// Set expiration time to 1 hour from now
	expiresAt := now + 3600
	
	result := FirestoreResult{
		ID:          jobID,
		ResultURL:   resultURL,
		PDFData:     pdfData,
		CreatedAt:   now,
		ExpiresAt:   expiresAt,
	}
	
	_, err := s.ResultsCollection().Doc(jobID).Set(ctx, result)
	if err != nil {
		log.Printf("Failed to store result for job %s: %v", jobID, err)
		return fmt.Errorf("failed to store result: %v", err)
	}
	
	log.Printf("Stored result for job %s (expires at %s)", jobID, time.Unix(expiresAt, 0).Format(time.RFC3339))
	return nil
}

// updateJobStatus updates a job's status in Firestore
func (s *Service) updateJobStatus(job *Job, status JobStatus, message, resultURL string) {
	ctx := context.Background()
	now := time.Now().Unix()

	// Update job in Firestore
	updates := []firestore.Update{
		{Path: "status", Value: string(status)},
		{Path: "message", Value: message},
		{Path: "updatedAt", Value: now},
	}

	_, err := s.Collection().Doc(job.ID).Update(ctx, updates)
	if err != nil {
		log.Printf("Failed to update job status in Firestore: %v", err)
	}

	// Update the in-memory job
	job.Status = status
	job.Message = message
	job.UpdatedAt = now
	if resultURL != "" {
		job.ResultURL = resultURL
	}

	log.Printf("Job %s updated: status=%s, message=%s", job.ID, status, message)
}

// GetResult retrieves a job result from Firestore
func (s *Service) GetResult(ctx context.Context, jobID string) (*FirestoreResult, error) {
	doc, err := s.ResultsCollection().Doc(jobID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("result not found")
		}
		return nil, fmt.Errorf("error retrieving result: %v", err)
	}
	
	var result FirestoreResult
	if err := doc.DataTo(&result); err != nil {
		return nil, fmt.Errorf("error parsing result data: %v", err)
	}
	
	// Check if result has expired
	now := time.Now().Unix()
	if result.ExpiresAt > 0 && now > result.ExpiresAt {
		// Result has expired, delete it
		_, err := s.ResultsCollection().Doc(jobID).Delete(ctx)
		if err != nil {
			log.Printf("Failed to delete expired result %s: %v", jobID, err)
		} else {
			log.Printf("Deleted expired result %s", jobID)
		}
		return nil, fmt.Errorf("result has expired")
	}
	
	return &result, nil
} 