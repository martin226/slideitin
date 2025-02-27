package prompts

import (
	"bytes"
	"text/template"

	"github.com/slideitin/backend/models"
)

// Templates for different prompt types
const (
	// Template for slide generation prompt
	slideGenerationTemplate = `Create a presentation with the theme: {{.Theme}}
	
Detail level: {{.DetailLevel}}
Target audience: {{.Audience}}

{{if .FileContents}}
The presentation should be based on the following content:
{{range $index, $file := .FileNames}}
FILE: {{$file}}
---
{{index $.FileContents $index}}
---
{{end}}
{{end}}

Generate a comprehensive presentation with clear slides covering the key points from the provided content.
`
)

// GenerateSlidePrompt creates a prompt for slide generation based on the given parameters
func GenerateSlidePrompt(theme string, fileContents []string, fileNames []string, settings models.SlideSettings) (string, error) {
	// Create template data
	data := map[string]interface{}{
		"Theme":        theme,
		"DetailLevel":  settings.SlideDetail,
		"Audience":     settings.Audience,
		"FileContents": fileContents,
		"FileNames":    fileNames,
	}

	// Parse and execute the template
	tmpl, err := template.New("slidePrompt").Parse(slideGenerationTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GenerateCustomPrompt creates a prompt from a custom template and parameters
func GenerateCustomPrompt(promptTemplate string, params map[string]interface{}) (string, error) {
	tmpl, err := template.New("customPrompt").Parse(promptTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
} 