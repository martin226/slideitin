package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/martin226/slideitin/backend/slides-service/services/slides"
	"github.com/martin226/slideitin/backend/slides-service/models"
	"os"
)

// FileReference represents a reference to a file stored in GCS
type FileReference struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
	GCSPath  string `json:"gcsPath"`
}

// TaskPayload represents the data structure received from Cloud Tasks
type TaskPayload struct {
	JobID     string            `json:"jobID"`
	Theme     string            `json:"theme"`
	Files     []FileReference   `json:"files"`
	Settings  models.SlideSettings `json:"settings"`
}

// FirestoreJob is the Firestore representation of a job
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

// TaskController handles requests from Cloud Tasks
type TaskController struct {
	slideService *slides.SlideService
	firestoreClient *firestore.Client
	storageClient *storage.Client
	bucketName string
}

// NewTaskController creates a new task controller
func NewTaskController(slideService *slides.SlideService, firestoreClient *firestore.Client) *TaskController {
	// Get bucket name from environment variables
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "slideitin-files" // Default bucket name
	}
	
	// Create Cloud Storage client
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Printf("Failed to create Cloud Storage client: %v", err)
		// Continue without storage client, will be handled in requests
	}
	
	return &TaskController{
		slideService: slideService,
		firestoreClient: firestoreClient,
		storageClient: storageClient,
		bucketName: bucketName,
	}
}

// downloadFileFromGCS downloads a file from Google Cloud Storage
func (c *TaskController) downloadFileFromGCS(ctx context.Context, gcsPath string) ([]byte, string, error) {
	// Get a handle to the bucket
	bucket := c.storageClient.Bucket(c.bucketName)
	
	// Get a handle to the object
	obj := bucket.Object(gcsPath)
	
	// Check if the object exists
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object attributes: %v", err)
	}
	
	// Create a reader for the object
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create reader: %v", err)
	}
	defer r.Close()
	
	// Read the file data
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %v", err)
	}
	
	return data, attrs.ContentType, nil
}

// ProcessSlides handles slide generation requests from Cloud Tasks
func (c *TaskController) ProcessSlides(ctx *gin.Context) {
	// Check if storage client is available
	if c.storageClient == nil {
		log.Printf("Storage client not available")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Storage client not configured"})
		return
	}
	
	// Parse task payload from request body
	var payload TaskPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		log.Printf("Failed to parse task payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid payload: %v", err)})
		return
	}
	
	// Create a job status update function
	statusUpdateFn := func(message string) error {
		return c.updateJobStatus(payload.JobID, "processing", message, "")
	}
	
	// Update initial job status
	if err := statusUpdateFn("Processing slides"); err != nil {
		log.Printf("Failed to update job status: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update job status: %v", err)})
		return
	}
	
	// Download files from GCS
	files := make([]models.File, 0, len(payload.Files))
	for _, fileRef := range payload.Files {
		// Download the file from GCS
		fileData, contentType, err := c.downloadFileFromGCS(ctx.Request.Context(), fileRef.GCSPath)
		if err != nil {
			log.Printf("Failed to download file %s: %v", fileRef.Filename, err)
			c.updateJobStatus(payload.JobID, "failed", fmt.Sprintf("Failed to download file %s: %v", fileRef.Filename, err), "")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to download file: %v", err)})
			return
		}
		
		// Create a file object
		file := models.File{
			Filename: fileRef.Filename,
			Data:     fileData,
			Type:     contentType,
		}
		files = append(files, file)
	}
	
	// Generate slides
	pdfData, htmlData, err := c.slideService.GenerateSlides(
		ctx.Request.Context(),
		payload.Theme,
		files,
		payload.Settings,
		statusUpdateFn,
	)
	
	if err != nil {
		log.Printf("Failed to generate slides: %v", err)
		c.updateJobStatus(payload.JobID, "failed", fmt.Sprintf("Failed to generate slides: %v", err), "")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate slides: %v", err)})
		return
	}
	
	// Create result URL
	resultURL := "/results/" + payload.JobID
	
	// Store result in Firestore
	if err := c.storeResult(ctx.Request.Context(), payload.JobID, resultURL, pdfData, htmlData); err != nil {
		log.Printf("Failed to store result: %v", err)
		c.updateJobStatus(payload.JobID, "failed", fmt.Sprintf("Failed to store result: %v", err), "")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to store result: %v", err)})
		return
	}
	
	// Clean up files from GCS
	for _, fileRef := range payload.Files {
		// Delete the file from GCS
		obj := c.storageClient.Bucket(c.bucketName).Object(fileRef.GCSPath)
		if err := obj.Delete(ctx.Request.Context()); err != nil {
			log.Printf("Warning: Failed to delete file %s from GCS: %v", fileRef.GCSPath, err)
			// Continue anyway, this is not a critical error
		} else {
			log.Printf("Deleted file %s from GCS", fileRef.GCSPath)
		}
	}
	
	// Mark job as completed
	if err := c.setJobCompleted(payload.JobID, "Slides generated successfully", resultURL); err != nil {
		log.Printf("Failed to mark job as completed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to mark job as completed: %v", err)})
		return
	}
	
	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "jobID": payload.JobID})
}

// updateJobStatus updates a job's status in Firestore
func (c *TaskController) updateJobStatus(jobID, status, message, resultURL string) error {
	ctx := context.Background()
	now := time.Now().Unix()
	
	// Update job in Firestore
	updates := []firestore.Update{
		{Path: "status", Value: status},
		{Path: "message", Value: message},
		{Path: "updatedAt", Value: now},
	}
	
	_, err := c.firestoreClient.Collection("jobs").Doc(jobID).Update(ctx, updates)
	if err != nil {
		log.Printf("Failed to update job status in Firestore: %v", err)
		return err
	}
	
	log.Printf("Job %s updated: status=%s, message=%s", jobID, status, message)
	return nil
}

// setJobCompleted marks a job as completed and sets it to expire
func (c *TaskController) setJobCompleted(jobID, message, resultURL string) error {
	ctx := context.Background()
	now := time.Now().Unix()
	// Set job to expire in 5 minutes
	expiresAt := now + 300 // 300 seconds = 5 minutes
	
	// Update job in Firestore
	updates := []firestore.Update{
		{Path: "status", Value: "completed"},
		{Path: "message", Value: message},
		{Path: "updatedAt", Value: now},
		{Path: "expiresAt", Value: expiresAt},
	}
	
	_, err := c.firestoreClient.Collection("jobs").Doc(jobID).Update(ctx, updates)
	if err != nil {
		log.Printf("Failed to update job status in Firestore: %v", err)
		return err
	}
	
	log.Printf("Job %s completed and will expire at %s", jobID, time.Unix(expiresAt, 0).Format(time.RFC3339))
	return nil
}

// storeResult stores a job result in Firestore
func (c *TaskController) storeResult(ctx context.Context, jobID, resultURL string, pdfData []byte, htmlData []byte) error {
	now := time.Now().Unix()
	// Set expiration time to 1 hour from now
	expiresAt := now + 3600
	
	result := FirestoreResult{
		ID:          jobID,
		ResultURL:   resultURL,
		PDFData:     pdfData,
		HTMLData:    htmlData,
		CreatedAt:   now,
		ExpiresAt:   expiresAt,
	}
	
	_, err := c.firestoreClient.Collection("results").Doc(jobID).Set(ctx, result)
	if err != nil {
		log.Printf("Failed to store result for job %s: %v", jobID, err)
		return fmt.Errorf("failed to store result: %v", err)
	}
	
	log.Printf("Stored result for job %s (expires at %s)", jobID, time.Unix(expiresAt, 0).Format(time.RFC3339))
	return nil
} 