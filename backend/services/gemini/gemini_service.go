package gemini

import (
	"context"
	"io"
	"log"
	"time"
	
	"github.com/slideitin/backend/models"
	"github.com/slideitin/backend/services/prompts"
	"bytes"
)

// Service handles interactions with the Gemini API
type Service struct {
	apiKey string
}

// NewService creates a new Gemini service
func NewService(apiKey string) *Service {
	return &Service{
		apiKey: apiKey,
	}
}

// GenerateSlides creates a presentation based on the provided theme, files, and settings
func (s *Service) GenerateSlides(
	ctx context.Context, 
	theme string, 
	files []struct {
		Filename string
		Data     []byte
		Type     string
	},
	settings models.SlideSettings,
	statusUpdateFn func(message string) error,
) (string, []byte, error) {
	// Update status to show we're processing the files
	if err := statusUpdateFn("Analyzing uploaded files"); err != nil {
		return "", nil, err
	}
	
	// Process files by creating readers from the stored data when needed
	// This ensures the file data is available even after the HTTP request finishes
	for _, file := range files {
		fileReader := io.NopCloser(bytes.NewReader(file.Data))
		
		// dummy use to prevent compiler error
		fileReader.Read(nil)
		
		log.Printf("Processing file: %s (%s)", file.Filename, file.Type)
	}
	
	time.Sleep(5 * time.Second)

	// Update status to show we're generating the prompt
	if err := statusUpdateFn("Generating content for slides"); err != nil {
		return "", nil, err
	}
	time.Sleep(5 * time.Second)
	
	// 2. Generate the prompt using the prompt generator
	prompt, err := prompts.GenerateSlidePrompt(theme, settings)
	if err != nil {
		log.Printf("Error generating prompt: %v", err)
		return "", nil, err
	}
	
	// Update status to show we're sending to Gemini
	if err := statusUpdateFn("Creating presentation with AI"); err != nil {
		return "", nil, err
	}
	time.Sleep(5 * time.Second)

	// 3. Send the prompt to Gemini (placeholder)
	log.Printf("Generated prompt: %s", prompt)
	log.Printf("Would send to Gemini with settings: %+v", settings)
	
	// Update status to show we're finalizing the presentation
	if err := statusUpdateFn("Finalizing presentation"); err != nil {
		return "", nil, err
	}
	time.Sleep(5 * time.Second)
	// For now, just return a placeholder
	return "placeholder-presentation-id", nil, nil
}
