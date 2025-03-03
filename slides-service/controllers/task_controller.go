package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"cloud.google.com/go/firestore"
	"github.com/martin226/slideitin/slides-service/services/slides"
	"github.com/martin226/slideitin/slides-service/models"
)

// TaskPayload represents the data structure received from Cloud Tasks
type TaskPayload struct {
	JobID     string            `json:"jobID"`
	Theme     string            `json:"theme"`
	Files     []models.File `json:"files"`
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
}

// NewTaskController creates a new task controller
func NewTaskController(slideService *slides.SlideService, firestoreClient *firestore.Client) *TaskController {
	return &TaskController{
		slideService: slideService,
		firestoreClient: firestoreClient,
	}
}

// ProcessSlides handles slide generation requests from Cloud Tasks
func (c *TaskController) ProcessSlides(ctx *gin.Context) {
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
	
	// Generate slides
	pdfData, htmlData, err := c.slideService.GenerateSlides(
		ctx.Request.Context(),
		payload.Theme,
		payload.Files,
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