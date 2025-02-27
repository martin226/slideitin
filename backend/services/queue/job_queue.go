package queue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/slideitin/backend/models"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

// Job represents a single slide generation job
type Job struct {
	ID          string
	Theme       string
	Files       [][]byte
	FileNames   []string
	Settings    models.SlideSettings
	Status      JobStatus
	Message     string
	ResultURL   string
	CreatedAt   int64
	UpdatedAt   int64
	Subscribers map[chan JobUpdate]bool
	mu          sync.RWMutex
}

// JobUpdate represents an update to a job that can be sent to subscribers
type JobUpdate struct {
	ID        string    `json:"id"`
	Status    JobStatus `json:"status"`
	Message   string    `json:"message"`
	ResultURL string    `json:"resultUrl,omitempty"`
	UpdatedAt int64     `json:"updatedAt"`
}

// Service manages a queue of slide generation jobs
type Service struct {
	jobs      map[string]*Job // TODO: Replace with Redis or something more scalable
	queue     chan *Job
	mu        sync.RWMutex
	geminiSvc interface {
		GenerateSlides(ctx context.Context, theme string, fileContents [][]byte, fileNames []string, settings models.SlideSettings, statusUpdateFn func(statusStr string, message string) error) (string, error)
	}
}

// NewService creates a new queue service
func NewService(geminiSvc interface {
	GenerateSlides(ctx context.Context, theme string, fileContents [][]byte, fileNames []string, settings models.SlideSettings, statusUpdateFn func(statusStr string, message string) error) (string, error)
}) *Service {
	s := &Service{
		jobs:      make(map[string]*Job),
		queue:     make(chan *Job, 100), // Buffer for up to 100 jobs
		geminiSvc: geminiSvc,
	}

	// Start background workers
	for i := 0; i < 3; i++ { // Number of concurrent workers
		go s.worker()
	}

	return s
}

// AddJob adds a new job to the queue
func (s *Service) AddJob(ctx context.Context, id, theme string, files [][]byte, fileNames []string, settings models.SlideSettings) (*Job, error) {
	// Check if queue is full before creating the job
	if len(s.queue) >= cap(s.queue) {
		log.Printf("Queue is full, rejected job %s", id)
		return nil, fmt.Errorf("system is busy, please try again later")
	}

	// Create the job
	now := time.Now().Unix()
	job := &Job{
		ID:          id,
		Theme:       theme,
		Files:       files,
		FileNames:   fileNames,
		Settings:    settings,
		Status:      StatusQueued,
		Message:     "Job added to queue",
		CreatedAt:   now,
		UpdatedAt:   now,
		Subscribers: make(map[chan JobUpdate]bool),
	}

	// Add job to map and queue
	s.mu.Lock()
	s.jobs[id] = job
	s.mu.Unlock()

	// Add job to processing queue
	s.queue <- job

	log.Printf("Added job %s to queue", id)
	return job, nil
}

// GetJob retrieves a job by its ID
func (s *Service) GetJob(id string) *Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.jobs[id]
}

// Subscribe adds a subscriber to receive job updates
func (s *Service) Subscribe(jobID string, updates chan JobUpdate) (bool, error) {
	s.mu.RLock()
	job, exists := s.jobs[jobID]
	s.mu.RUnlock()

	if !exists {
		return false, nil
	}

	job.mu.Lock()
	job.Subscribers[updates] = true
	job.mu.Unlock()

	// Send initial update
	updates <- JobUpdate{
		ID:        job.ID,
		Status:    job.Status,
		Message:   job.Message,
		ResultURL: job.ResultURL,
		UpdatedAt: job.UpdatedAt,
	}

	return true, nil
}

// Unsubscribe removes a subscriber from job updates
func (s *Service) Unsubscribe(jobID string, updates chan JobUpdate) {
	s.mu.RLock()
	job, exists := s.jobs[jobID]
	s.mu.RUnlock()

	if !exists {
		return
	}

	job.mu.Lock()
	delete(job.Subscribers, updates)
	job.mu.Unlock()
}

// worker processes jobs from the queue
func (s *Service) worker() {
	for job := range s.queue {
		s.processJob(job)
	}
}

// processJob processes a slide generation job
func (s *Service) processJob(job *Job) {
	// Update job status to processing
	s.updateJobStatus(job, StatusProcessing, "Processing slides", "")

	// Create a status update function to pass to the Gemini service
	statusUpdateFn := func(statusStr string, message string) error {
		// Convert string status to JobStatus
		var status JobStatus
		switch statusStr {
		case "processing":
			status = StatusProcessing
		case "completed":
			status = StatusCompleted
		case "failed":
			status = StatusFailed
		default:
			status = StatusProcessing
		}
		
		s.updateJobStatus(job, status, message, "")
		return nil
	}

	// Call the Gemini service with the status update function
	ctx := context.Background()
	resultID, err := s.geminiSvc.GenerateSlides(
		ctx, 
		job.Theme, 
		job.Files, 
		job.FileNames, 
		job.Settings,
		statusUpdateFn,
	)

	if err != nil {
		s.updateJobStatus(job, StatusFailed, "Failed to generate slides: "+err.Error(), "")
		return
	}

	// Update job as completed
	resultURL := "/api/v1/results/" + resultID // This would be the actual URL to download results
	s.updateJobStatus(job, StatusCompleted, "Slides generated successfully", resultURL)
}

func (s *Service) removeJob(job *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, job.ID)
	log.Printf("Job %s removed from queue", job.ID)
}

// updateJobStatus updates a job's status and notifies subscribers
func (s *Service) updateJobStatus(job *Job, status JobStatus, message, resultURL string) {
	job.mu.Lock()
	job.Status = status
	job.Message = message
	job.UpdatedAt = time.Now().Unix()
	if resultURL != "" {
		job.ResultURL = resultURL
	}

	// Create update object
	update := JobUpdate{
		ID:        job.ID,
		Status:    job.Status,
		Message:   job.Message,
		ResultURL: job.ResultURL,
		UpdatedAt: job.UpdatedAt,
	}

	// Notify all subscribers
	for ch := range job.Subscribers {
		select {
		case ch <- update:
			// Successfully sent
		default:
			// Channel full or closed, will be cleaned up later
			log.Printf("Could not send update to subscriber for job %s", job.ID)
		}
	}
	job.mu.Unlock()

	log.Printf("Job %s updated: status=%s, message=%s", job.ID, status, message)

	if status == StatusCompleted || status == StatusFailed {
		go s.removeJob(job)
	}
} 