package slides

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"github.com/martin226/slideitin/backend/slides-service/models"
	"github.com/martin226/slideitin/backend/slides-service/services/prompts"
	"bytes"
)

// SlideService handles interactions with the Gemini API
type SlideService struct {
	client *genai.Client
	model *genai.GenerativeModel
}

// NewSlideService creates a new Slide service
func NewSlideService(apiKey string) *SlideService {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	model := client.GenerativeModel("gemini-1.5-flash")
	model.SetMaxOutputTokens(4096)
	return &SlideService{
		client: client,
		model: model,
	}
}

// GenerateSlides creates a presentation based on the provided theme, files, and settings
func (s *SlideService) GenerateSlides(
	ctx context.Context, 
	theme string, 
	files []models.File,
	settings models.SlideSettings,
	statusUpdateFn func(message string) error,
) ([]byte, []byte, error) {
	// Update status to show we're processing the files
	if err := statusUpdateFn("Analyzing uploaded files"); err != nil {
		return nil, nil, err
	}

	geminiFiles := make([]*genai.File, 0, len(files))
	// Process files by creating readers from the stored data when needed
	// This ensures the file data is available even after the HTTP request finishes
	for _, file := range files {
		fileReader := io.NopCloser(bytes.NewReader(file.Data))
		
		// Upload the file to Gemini
		geminiFile, err := s.client.UploadFile(ctx, "", fileReader, &genai.UploadFileOptions{
			DisplayName: file.Filename,
			MIMEType: file.Type,
		})
		if err != nil {
			log.Printf("Failed to upload file to Gemini: %v", err)
			return nil, nil, err
		}
		geminiFiles = append(geminiFiles, geminiFile)
		log.Printf("Processing file: %s (%s)", file.Filename, file.Type)
	}

	// Update status to show we're generating the prompt
	if err := statusUpdateFn("Generating content for slides"); err != nil {
		return nil, nil, err
	}
	
	// 2. Generate the prompt using the prompt generator
	prompt, err := prompts.GenerateSlidePrompt(theme, settings)
	if err != nil {
		log.Printf("Error generating prompt: %v", err)
		return nil, nil, err
	}
	log.Printf("Prompt: %s", prompt)
	
	// Update status to show we're sending to Gemini
	if err := statusUpdateFn("Creating presentation with AI"); err != nil {
		return nil, nil, err
	}
	
	// 3. Send the prompt to Gemini
	parts := []genai.Part{}
	for _, file := range geminiFiles {
		parts = append(parts, genai.FileData{URI: file.URI})
	}
	parts = append(parts, genai.Text(prompt))

	// Ensure input tokens do not exceed 16384
	countResp, err := s.model.CountTokens(ctx, parts...)
	if err != nil {
		log.Printf("Failed to count tokens: %v", err)
		return nil, nil, err
	}
	if countResp.TotalTokens > 16384 {
		log.Printf("Input tokens exceed 16384: %d", countResp.TotalTokens)
		return nil, nil, errors.New("documents are too large to process")
	}

	resp, err := s.model.GenerateContent(ctx, parts...)
	if err != nil {
		log.Printf("Failed to generate content: %v", err)
		return nil, nil, err
	}

	respText := resp.Candidates[0].Content.Parts[0].(genai.Text)
	// Extract the markdown from the response between triple backticks
	// Match any language specifier or none at all
	respString := string(respText)
	marpText := extractMarkdownContent(respString)
	
	if marpText == "" {
		log.Printf("No markdown found in response: %s", respText)
		return nil, nil, errors.New("failed to generate presentation. Please try again.")
	}

	log.Printf("Generated presentation: %s", marpText)
	
	// Update status to show we're finalizing the presentation
	if err := statusUpdateFn("Finalizing presentation"); err != nil {
		return nil, nil, err
	}

	// Create a temporary directory for our files
	tempDir, err := os.MkdirTemp("", "slideitin-")
	if err != nil {
		log.Printf("Failed to create temp directory: %v", err)
		return nil, nil, err
	}
	defer os.RemoveAll(tempDir) // Clean up when we're done
	
	// Create the markdown file
	mdFilePath := filepath.Join(tempDir, "presentation.md")
	err = os.WriteFile(mdFilePath, []byte(marpText), 0644)
	if err != nil {
		log.Printf("Failed to write markdown file: %v", err)
		return nil, nil, err
	}
	
	// Set up PDF output path
	pdfFilePath := filepath.Join(tempDir, "presentation.pdf")
	
	// Run Marp CLI to generate the PDF
	marpArgs := []string{"@marp-team/marp-cli", mdFilePath}
	
	// Add theme parameter if it's in themes directory
	themePath := filepath.Join("services", "slides", "themes", theme+".css")
	if _, err := os.Stat(themePath); err == nil {
		// Theme file exists, add it to the arguments
		marpArgs = append(marpArgs, "--theme", themePath)
		log.Printf("Using theme: %s", themePath)
	} else {
		marpArgs = append(marpArgs, "--theme", theme)
		log.Printf("Using built-in theme: %s", theme)
	}
	
	cmd := exec.Command("npx", append(marpArgs, "--output", pdfFilePath, "--pdf")...)
	var cmdOutput bytes.Buffer
	var cmdError bytes.Buffer
	cmd.Stdout = &cmdOutput
	cmd.Stderr = &cmdError
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to run Marp CLI: %v", err)
		log.Printf("Marp CLI stderr: %s", cmdError.String())
		return nil, nil, errors.New("failed to generate PDF. Please try again.")
	}
	
	// Read the generated PDF
	pdfBytes, err := os.ReadFile(pdfFilePath)
	if err != nil {
		log.Printf("Failed to read generated PDF: %v", err)
		return nil, nil, err
	}
	
	log.Printf("Successfully generated PDF (%d bytes)", len(pdfBytes))

	// Create the HTML file
	htmlFilePath := filepath.Join(tempDir, "presentation.html")

	// Run Marp CLI to generate the HTML
	cmd = exec.Command("npx", append(marpArgs, "--output", htmlFilePath, "--html")...)
	cmdOutput.Reset()
	cmdError.Reset()
	cmd.Stdout = &cmdOutput
	cmd.Stderr = &cmdError
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to run Marp CLI: %v", err)
		log.Printf("Marp CLI stderr: %s", cmdError.String())
		return nil, nil, errors.New("failed to generate HTML. Please try again.")
	}

	// Read the generated HTML
	htmlBytes, err := os.ReadFile(htmlFilePath)
	if err != nil {
		log.Printf("Failed to read generated HTML: %v", err)
		return nil, nil, err
	}

	log.Printf("Successfully generated HTML (%d bytes)", len(htmlBytes))
	
	// Delete the files from Gemini
	for _, file := range geminiFiles {
		err := s.client.DeleteFile(ctx, file.Name)
		if err != nil {
			log.Printf("Failed to delete file from Gemini: %v", err)
		}
	}
	
	// Return the PDF and HTML bytes
	return pdfBytes, htmlBytes, nil
}

// extractMarkdownContent extracts markdown content between triple backticks
func extractMarkdownContent(text string) string {
	lines := regexp.MustCompile(`\r?\n`).Split(text, -1)
	
	firstBacktickLine := -1
	lastBacktickLine := -1
	
	// Find first and last lines with triple backticks
	for i, line := range lines {
		if strings.HasPrefix(line, "```") {
			if firstBacktickLine == -1 {
				firstBacktickLine = i
			}
			lastBacktickLine = i
		}
	}
	
	// If we found backticks, extract the content
	if firstBacktickLine != -1 && lastBacktickLine != -1 && lastBacktickLine > firstBacktickLine {
		// Extract content between the backtick lines, excluding the lines with backticks themselves
		// firstBacktickLine+1 skips the opening backtick line
		// lastBacktickLine as the end index (exclusive in Go slices) excludes the closing backtick line
		content := lines[firstBacktickLine+1:lastBacktickLine]
		return strings.Join(content, "\n")
	}
	
	// If no backticks found, return the entire text
	return text
} 