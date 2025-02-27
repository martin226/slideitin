package gemini

import (
	"context"
	"log"
	"time"
	"github.com/slideitin/backend/models"
	"github.com/slideitin/backend/services/prompts"
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
	fileContents [][]byte, 
	fileNames []string, 
	settings models.SlideSettings,
	statusUpdateFn func(status string, message string) error,
) (string, error) {
	// Update status to show we're processing the files
	if err := statusUpdateFn("processing", "Analyzing uploaded files"); err != nil {
		return "", err
	}
	time.Sleep(5 * time.Second)

	// 1. Convert binary file contents to strings for prompt generation
	fileStrings := make([]string, len(fileContents))
	for i, content := range fileContents {
		fileStrings[i] = string(content)
	}
	
	// Update status to show we're generating the prompt
	if err := statusUpdateFn("processing", "Generating content for slides"); err != nil {
		return "", err
	}
	time.Sleep(5 * time.Second)
	
	// 2. Generate the prompt using the prompt generator
	prompt, err := prompts.GenerateSlidePrompt(theme, fileStrings, fileNames, settings)
	if err != nil {
		log.Printf("Error generating prompt: %v", err)
		return "", err
	}
	
	// Update status to show we're sending to Gemini
	if err := statusUpdateFn("processing", "Creating presentation with AI"); err != nil {
		return "", err
	}
	time.Sleep(5 * time.Second)

	// 3. Send the prompt to Gemini (placeholder)
	log.Printf("Generated prompt: %s", prompt)
	log.Printf("Would send to Gemini with settings: %+v", settings)
	
	// Update status to show we're finalizing the presentation
	if err := statusUpdateFn("processing", "Finalizing presentation"); err != nil {
		return "", err
	}
	time.Sleep(5 * time.Second)
	// For now, just return a placeholder
	return "placeholder-presentation-id", nil
}
