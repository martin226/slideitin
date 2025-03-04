package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
	"bytes"
	"path/filepath"

	"cloud.google.com/go/firestore"
	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"cloud.google.com/go/storage"
	"github.com/martin226/slideitin/backend/api/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
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
	HTMLData    []byte `firestore:"htmlData"`
	CreatedAt   int64  `firestore:"createdAt"`
	ExpiresAt   int64  `firestore:"expiresAt"`
}

// Job represents a single slide generation job with runtime features
type Job struct {
	ID        string
	Theme     string
	Files     []models.File
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

// FileReference represents a reference to a file stored in GCS
type FileReference struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
	GCSPath  string `json:"gcsPath"`
}

// TaskPayload represents the data structure to be sent in a Cloud Task
type TaskPayload struct {
	JobID     string            `json:"jobID"`
	Theme     string            `json:"theme"`
	Files     []FileReference   `json:"files"`
	Settings  models.SlideSettings `json:"settings"`
}

// Service manages jobs using Firestore, Cloud Tasks, and Cloud Storage
type Service struct {
	client     *firestore.Client
	taskClient *cloudtasks.Client
	storageClient *storage.Client
	projectID  string
	region     string
	queueID    string
	serviceURL string
	bucketName string
}

// NewService creates a new queue service using Firestore, Cloud Tasks, and Cloud Storage
func NewService(client *firestore.Client) (*Service, error) {
	// Get environment variables
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GOOGLE_CLOUD_PROJECT environment variable is required")
	}
	
	region := os.Getenv("CLOUD_TASKS_REGION")
	if region == "" {
		region = "us-central1" // Default region
	}
	
	queueID := os.Getenv("CLOUD_TASKS_QUEUE_ID")
	if queueID == "" {
		queueID = "slides-generation-queue" // Default queue ID
	}
	
	serviceURL := os.Getenv("SLIDES_SERVICE_URL")
	if serviceURL == "" {
		return nil, fmt.Errorf("SLIDES_SERVICE_URL environment variable is required")
	}
	
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "slideitin-files" // Default bucket name
	}
	
	// Create Cloud Tasks client
	ctx := context.Background()
	taskClient, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Tasks client: %v", err)
	}
	
	// Create Cloud Storage client
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloud Storage client: %v", err)
	}
	
	return &Service{
		client:        client,
		taskClient:    taskClient,
		storageClient: storageClient,
		projectID:     projectID,
		region:        region,
		queueID:       queueID,
		serviceURL:    serviceURL,
		bucketName:    bucketName,
	}, nil
}

// Collection returns the Firestore collection reference for jobs
func (s *Service) Collection() *firestore.CollectionRef {
	return s.client.Collection("jobs")
}

// ResultsCollection returns the Firestore collection reference for results
func (s *Service) ResultsCollection() *firestore.CollectionRef {
	return s.client.Collection("results")
}

// uploadFileToGCS uploads a file to Google Cloud Storage and returns its GCS path
func (s *Service) uploadFileToGCS(ctx context.Context, jobID string, file models.File) (string, error) {
	// Create a GCS object path: jobID/filename
	objectPath := filepath.Join(jobID, file.Filename)
	
	// Get a handle to the bucket
	bucket := s.storageClient.Bucket(s.bucketName)
	
	// Check if the bucket exists, if not create it
	if _, err := bucket.Attrs(ctx); err != nil {
		if err == storage.ErrBucketNotExist {
			if err := bucket.Create(ctx, s.projectID, nil); err != nil {
				return "", fmt.Errorf("failed to create bucket: %v", err)
			}
		} else {
			return "", fmt.Errorf("failed to check bucket: %v", err)
		}
	}
	
	// Create a writer for the object
	obj := bucket.Object(objectPath)
	w := obj.NewWriter(ctx)
	w.ContentType = file.Type
	
	// Write the file data to GCS
	if _, err := io.Copy(w, bytes.NewReader(file.Data)); err != nil {
		w.Close()
		return "", fmt.Errorf("failed to write file to GCS: %v", err)
	}
	
	// Close the writer
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %v", err)
	}
	
	log.Printf("Uploaded file %s to GCS: gs://%s/%s", file.Filename, s.bucketName, objectPath)
	
	return objectPath, nil
}

// AddJob adds a new job to Firestore, uploads files to GCS, and creates a Cloud Task for processing
func (s *Service) AddJob(ctx context.Context, id, theme string, fileData []models.File, settings models.SlideSettings) (*Job, error) {
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

	// Create in-memory job object
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

	// Upload files to GCS
	fileRefs := make([]FileReference, 0, len(fileData))
	for _, file := range fileData {
		// Upload the file to GCS
		gcsPath, err := s.uploadFileToGCS(ctx, id, file)
		if err != nil {
			// Update job status to failed if file upload fails
			s.updateJobStatus(job, StatusFailed, fmt.Sprintf("Failed to upload file %s: %v", file.Filename, err), "")
			return job, fmt.Errorf("failed to upload file: %v", err)
		}
		
		// Create a file reference
		fileRef := FileReference{
			Filename: file.Filename,
			Type:     file.Type,
			GCSPath:  gcsPath,
		}
		fileRefs = append(fileRefs, fileRef)
	}

	// Create a Cloud Task to process the job
	err = s.createTask(ctx, job, fileRefs)
	if err != nil {
		// Update job status to failed if task creation fails
		s.updateJobStatus(job, StatusFailed, fmt.Sprintf("Failed to queue job: %v", err), "")
		return job, fmt.Errorf("failed to create Cloud Task: %v", err)
	}

	return job, nil
}

// createTask creates a Cloud Task to process a job
func (s *Service) createTask(ctx context.Context, job *Job, fileRefs []FileReference) error {
	taskPayload := TaskPayload{
		JobID: job.ID,
		Theme: job.Theme,
		Files: fileRefs,
		Settings: job.Settings,
	}
	
	payloadBytes, err := json.Marshal(taskPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %v", err)
	}
	
	// Define the Cloud Tasks queue path
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", s.projectID, s.region, s.queueID)
	
	// Define the target endpoint
	taskURL := fmt.Sprintf("%s/tasks/process-slides", s.serviceURL)

	// Create the Cloud Task with OIDC token
	task := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// Name is assigned by the server
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        taskURL,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: payloadBytes,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: fmt.Sprintf("%s@%s.iam.gserviceaccount.com", "slides-service-invoker", s.projectID),
							Audience:            taskURL,
						},
					},
				},
			},
			ScheduleTime: timestamppb.New(time.Now()),
		},
	}
	
	// Create the task
	_, err = s.taskClient.CreateTask(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}
	
	log.Printf("Created Cloud Task for job %s with %d file references", job.ID, len(fileRefs))
	return nil
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