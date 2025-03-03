package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/martin226/slideitin/backend/api/models"
	"github.com/martin226/slideitin/backend/api/services/queue"
)

// SlideController handles the slide generation API endpoints
type SlideController struct {
	queueService  *queue.Service
}

// NewSlideController creates a new slide controller
func NewSlideController(queueService *queue.Service) *SlideController {
	return &SlideController{
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

	// Validate theme
	isValidTheme := false
	for _, theme := range models.ValidThemes {
		if req.Theme == theme {
			isValidTheme = true
			break
		}
	}
	if !isValidTheme {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid theme: %s. Supported themes are: %s", req.Theme, strings.Join(models.ValidThemes, ", ")),
		})
		return
	}

	// Validate slideDetail setting
	isValidSlideDetail := false
	if req.Settings.SlideDetail != "" {
		for _, detail := range models.ValidSlideDetails {
			if req.Settings.SlideDetail == detail {
				isValidSlideDetail = true
				break
			}
		}
		if !isValidSlideDetail {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid slideDetail: %s. Supported values are: %s", 
					req.Settings.SlideDetail, strings.Join(models.ValidSlideDetails, ", ")),
			})
			return
		}
	}

	// Validate audience setting
	isValidAudience := false
	if req.Settings.Audience != "" {
		for _, audience := range models.ValidAudiences {
			if req.Settings.Audience == audience {
				isValidAudience = true
				break
			}
		}
		if !isValidAudience {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Invalid audience: %s. Supported values are: %s", 
					req.Settings.Audience, strings.Join(models.ValidAudiences, ", ")),
			})
			return
		}
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

	// Read file data into memory to prevent it from being released
	fileData := make([]models.File, 0, len(files))
	
	for _, file := range files {
		// Open the file
		src, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to open file %s: %v", file.Filename, err),
			})
			return
		}
		
		// Read the file data
		data, err := io.ReadAll(src)
		src.Close() // Close the file after reading
		
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to read file %s: %v", file.Filename, err),
			})
			return
		}
		
		// Detect MIME type from file content instead of using header
		// DetectContentType only needs the first 512 bytes
		mimeType := http.DetectContentType(data)
		
		// Remove charset information if present
		if semicolonIndex := strings.Index(mimeType, ";"); semicolonIndex != -1 {
			mimeType = strings.TrimSpace(mimeType[:semicolonIndex])
		}
		
		// Validate file type - only allow PDF, Markdown and TXT
		isAllowed := false

		// Check by file extension first
		fileExt := strings.ToLower(filepath.Ext(file.Filename))
		if fileExt == ".pdf" || fileExt == ".md" || fileExt == ".txt" {
			// Now check MIME type
			if mimeType == "application/pdf" {
				// PDF is valid
				isAllowed = true
			} else if mimeType == "text/plain" {
				// Plain text (could be TXT or MD)
				isAllowed = true
			} else if strings.Contains(mimeType, "markdown") || strings.Contains(mimeType, "text/") {
				// Some systems detect markdown as text/markdown, text/x-markdown, or just text/plain
				// For text files, we'll trust the extension more than the mime type
				if fileExt == ".md" || fileExt == ".txt" {
					isAllowed = true
				}
			}
		}

		if !isAllowed {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Unsupported file type: %s. Only PDF, Markdown, and TXT files are allowed", file.Filename),
			})
			return
		}
		
		// Store the file data
		fileData = append(fileData, models.File{
			Filename: file.Filename,
			Data:     data,
			Type:     mimeType,
		})
	}

	// Log the request
	log.Printf("Received slide generation request: Theme: %s, Files count: %d, Settings: %+v", 
		req.Theme, len(fileData), req.Settings)

	// Generate a unique job ID
	jobID := uuid.New().String()

	// Add job to queue instead of processing immediately
	job, err := c.queueService.AddJob(ctx, jobID, req.Theme, fileData, req.Settings)
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
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
	ctx.Writer.Header().Set("X-Accel-Buffering", "no") // Disable buffering in Nginx if used
	ctx.Writer.Flush()

	// Create channel for job updates and set up a cancellation context
	updates := make(chan queue.JobUpdate, 10)
	streamCtx, cancelStream := context.WithCancel(ctx.Request.Context())
	defer cancelStream()

	// Watch for job updates from Firestore
	go func() {
		defer close(updates)
		err := c.queueService.WatchJob(streamCtx, id, updates)
		if err != nil && err != context.Canceled {
			log.Printf("Error watching job %s: %v", id, err)
		}
	}()

	// Stream events to client
	ctx.Stream(func(w io.Writer) bool {
		// Check if client closed connection
		if ctx.Request.Context().Err() != nil {
			cancelStream()
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
				
				cancelStream()
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

// GetSlideResult handles retrieving and serving the presentation result
func (c *SlideController) GetSlideResult(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing result ID",
		})
		return
	}

	// Retrieve the result from Firestore
	result, err := c.queueService.GetResult(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Result not found: %v", err),
		})
		return
	}

	download := ctx.Query("download")

	if download == "true" {
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=presentation-%s.pdf", id))
		ctx.Data(http.StatusOK, "application/pdf", result.PDFData)
	} else {
		ctx.Header("Content-Type", "text/html")
		ctx.Data(http.StatusOK, "text/html", result.HTMLData)
	}
	return
} 