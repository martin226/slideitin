package models

import (
	"mime/multipart"
)

// Enum values for slide settings
var (
	// Valid themes
	ValidThemes = []string{"default", "beam", "rose_pine", "gaia", "uncover", "graph_paper"}
	
	// Valid slide detail levels
	ValidSlideDetails = []string{"minimal", "medium", "detailed"}
	
	// Valid audience types
	ValidAudiences = []string{"general", "academic", "technical", "professional", "executive"}
)

// SlideSettings represents the settings for slide generation
type SlideSettings struct {
	SlideDetail string `json:"slideDetail"` // Values: minimal, medium, detailed
	Audience    string `json:"audience"`    // Values: general, academic, technical, professional, executive
}

// SlideRequest represents the incoming request for slide generation
type SlideRequest struct {
	Theme    string       `json:"theme" binding:"required"`
	Settings SlideSettings `json:"settings" binding:"required"`
	// Files will be handled separately through multipart form
}

// FileUpload is a wrapper around the multipart.FileHeader for easier handling
type FileUpload struct {
	Files []*multipart.FileHeader `form:"files"`
}

// SlideResponse represents the response for a slide generation request
type SlideResponse struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	CreatedAt  int64  `json:"createdAt"`
	UpdatedAt  int64  `json:"updatedAt"`
} 