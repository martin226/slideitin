package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/slideitin/backend/models"
	"github.com/slideitin/backend/services/gemini"
	"github.com/slideitin/backend/services/queue"
)

// SlideController handles the slide generation API endpoints
type SlideController struct {
	geminiService *gemini.Service
	queueService  *queue.Service
}

// NewSlideController creates a new slide controller
func NewSlideController(geminiService *gemini.Service, queueService *queue.Service) *SlideController {
	return &SlideController{
		geminiService: geminiService,
		queueService:  queueService,
	}
}

// GenerateSlides handles the slide generation request
func (c *SlideController) GenerateSlides(ctx *gin.Context) {
	// Parse form data first
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to parse form data",
		})
		return
	}

	// Parse JSON data from form
	var req models.SlideRequest
	jsonData := ctx.PostForm("data")
	if jsonData == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing data field in form",
		})
		return
	}

	if err := json.Unmarshal([]byte(jsonData), &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request format: %v", err),
		})
		return
	}

	// Get files
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get files",
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "No files uploaded",
		})
		return
	}

	// Process files
	fileContents := make([][]byte, 0, len(files))
	fileNames := make([]string, 0, len(files))

	for _, file := range files {
		// Open file
		f, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to open file: %v", err),
			})
			return
		}
		defer f.Close()

		// Read file content
		content := make([]byte, file.Size)
		if _, err := f.Read(content); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to read file: %v", err),
			})
			return
		}

		fileContents = append(fileContents, content)
		fileNames = append(fileNames, file.Filename)
	}

	// Log the request
	log.Printf("Received slide generation request: Theme: %s, Files: %v, Settings: %+v", 
		req.Theme, fileNames, req.Settings)

	// Generate a unique job ID
	jobID := uuid.New().String()

	// Add job to queue instead of processing immediately
	job, err := c.queueService.AddJob(ctx, jobID, req.Theme, fileContents, fileNames, req.Settings)
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return response immediately with job ID
	ctx.JSON(http.StatusAccepted, models.SlideResponse{
		ID:        jobID,
		Status:    string(job.Status),
		Message:   job.Message,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	})
}

// StreamSlideStatus handles both regular status checks and SSE streaming of job status updates
func (c *SlideController) StreamSlideStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing job ID",
		})
		return
	}

	// Get job status from queue
	job := c.queueService.GetJob(id)
	if job == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
		})
		return
	}

	// Check if client accepts SSE
	acceptHeader := ctx.GetHeader("Accept")
	wantsSSE := acceptHeader == "text/event-stream"

	// If client doesn't want SSE, return a regular JSON response
	if !wantsSSE {
		ctx.JSON(http.StatusOK, gin.H{
			"id":        job.ID,
			"status":    job.Status,
			"message":   job.Message,
			"resultUrl": job.ResultURL,
			"updatedAt": job.UpdatedAt,
		})
		return
	}

	// For SSE clients, set headers for streaming
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	ctx.Writer.Header().Set("X-Accel-Buffering", "no") // Disable buffering in Nginx if used
	ctx.Writer.Flush()

	// Create channel for job updates
	updates := make(chan queue.JobUpdate, 10)
	defer close(updates)

	// Subscribe to job updates
	exists, err := c.queueService.Subscribe(id, updates)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to subscribe to job updates: %v", err),
		})
		return
	}

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
		})
		return
	}

	// Ensure client is notified when connection is closed
	ctx.Stream(func(w io.Writer) bool {
		// Check if client closed connection
		if ctx.Request.Context().Err() != nil {
			c.queueService.Unsubscribe(id, updates)
			return false
		}

		select {
		case update, ok := <-updates:
			if !ok {
				return false // Channel closed
			}

			// Send SSE event with job update
			ctx.SSEvent("update", update)
			
			// If job is completed or failed, end the stream
			if update.Status == queue.StatusCompleted || update.Status == queue.StatusFailed {
				// Send a final event indicating the stream will close
				ctx.SSEvent("close", gin.H{
					"id":      update.ID,
					"status":  update.Status,
					"message": "Stream closing normally",
				})
				ctx.Writer.Flush()
				
				// Wait a moment before closing to ensure the message is sent
				time.Sleep(100 * time.Millisecond)
				
				c.queueService.Unsubscribe(id, updates)
				return false
			}
			
			return true

		case <-time.After(30 * time.Second):
			// Send heartbeat to keep connection alive
			ctx.SSEvent("ping", nil)
			return true
		}
	})
} 