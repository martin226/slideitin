package prompts

import (
	"bytes"
	"text/template"

	"github.com/martin226/slideitin/backend/slides-service/models"
)

// Templates for different prompt types
const (
	// Template for slide generation prompt
	slideGenerationTemplate = `You are an expert at creating Marp markdown presentations. You are highly skilled at extracting content from documents and creating beautiful, well-designed presentations.
	
Create a Marp markdown presentation using the following instructions:

The following is an example of how to create a Marp markdown presentation. All of the frontmatter in the example is also required for your response, other than the header and footer.

{{.ThemeExample}}

Theme: {{.Theme}}

{{.DetailLevel}}

{{.Audience}}

IMPORTANT GUIDELINES:
1. Always begin with a short title slide with a title, a short description, and author name (only if provided). The title should be an H1 header, the description should be a regular text, and the author name should be a regular text.
2. Ensure that the content on each slide fits inside the slide. Never create paragraphs.
3. Always use bullet points and other formatting options to make the content more readable. 
4. Prefer multi-line code blocks over inline code blocks for any code longer than a few words. Even if the code is a single line, use a multi-line code block.
5. Do not end with --- (three dashes) on a new line, since this will end the presentation with an empty slide.

Make the slides look as beautiful and well-designed as possible. Use all of the formatting options available to you.

Enclose your response in triple backticks like this:

` + "```md" + `
<your response here>
` + "```"

	// Common markdown header template used across all themes
	commonMarpHeader = `---
marp: true
theme: {{.Theme}}
{{if .UseLeadClass -}}
_class: lead
{{- end}}
paginate: true
header: This is an optional header {{.HeaderLocation}}
footer: This is an optional footer {{.FooterLocation}}
---
{{if .HasTitleClass}}
<!-- _class: title -->
{{end}}
# Title

`

	// Common markdown body template for examples
	commonExampleBody = `## Heading 2

- {{.ThemeDescription}}
{{ if .HasInvertClass}}

---

<!-- _class: invert -->

## Inverted color scheme

- You can use the <!-- _class: invert --> tag at the top of a slide to create a dark mode slide.
- Use this when you want to have a slide with a different color scheme than the rest of the presentation.
- Do this when a slide should stand out.
{{end}}{{if .HasTinyTextClass}}

---

<!-- _class: tinytext -->

# Tinytext class

- You can use the <!-- _class: tinytext --> tag at the top of a slide to make some text tiny.
- This might be useful for References.
{{end}}

---

## Code blocks

### Multi-line code blocks

` + "```" + `python
print("This is a code block")
print("You can use triple backticks to create a code block")
print("You can also use the language name to highlight the code block")
` + "```" + `

- **Another example:**

` + "```" + `c
printf("This is another code block");
printf("Always specify the language name for code blocks");
` + "```" + `

---

### Inline code blocks

- ` + "`" + `this` + "`" + ` is an inline code block
- You can use it using single backticks like this: ` + "`" + `this` + "`" + `

---

## Creating new slides

- To create a new slide, use a new line with three dashes like this:

` + "```" + `
---

# New slide
` + "```" + `

---

# Conclusion

- You can use Markdown formatting to create **bold**, *italic*, and ~~strikethrough~~ text.
> This is a block quote
This is regular text`
)

// Theme configurations
var themeConfigs = map[string]map[string]interface{}{
	"default": {
		"UseLeadClass":    true,
		"HasInvertClass":  true,
		"HasTinyTextClass": false,
		"HasTitleClass":   false,
		"HeaderLocation":  "(top left of the slide)",
		"FooterLocation":  "(bottom left of the slide)",
		"ThemeDescription": "By default, the color scheme for each slide is light.",
	},
	"beam": {
		"UseLeadClass":    false,
		"HasInvertClass":  false,
		"HasTinyTextClass": true,
		"HasTitleClass":   true,
		"HeaderLocation":  "(bottom left half of the slide)",
		"FooterLocation":  "(bottom right half of the slide)",
		"ThemeDescription": "IMPORTANT: You must use the above title class tag at the top of the title slide (<!-- _class: title -->).\n- Beam is a light color scheme based on the LaTeX Beamer theme.",
	},
	"rose-pine": {
		"UseLeadClass":    true,
		"HasInvertClass":  false,
		"HasTinyTextClass": false,
		"HasTitleClass":   false,
		"HeaderLocation":  "(top left of the slide)",
		"FooterLocation":  "(bottom left of the slide)",
		"ThemeDescription": "Rose Pine is a dark color scheme.",
	},
	"gaia": {
		"UseLeadClass":    true,
		"HasInvertClass":  true,
		"HasTinyTextClass": false,
		"HasTitleClass":   false,
		"HeaderLocation":  "(top left of the slide)",
		"FooterLocation":  "(bottom left of the slide)",
		"ThemeDescription": "By default, the color scheme for each slide is light.",
	},
	"uncover": {
		"UseLeadClass":    true,
		"HasInvertClass":  true,
		"HasTinyTextClass": false,
		"HasTitleClass":   false,
		"HeaderLocation":  "(top middle of the slide)",
		"FooterLocation":  "(bottom middle of the slide)",
		"ThemeDescription": "By default, the color scheme for each slide is light.",
	},
	"graph_paper": {
		"UseLeadClass":    true,
		"HasInvertClass":  false,
		"HasTinyTextClass": true,
		"HasTitleClass":   false,
		"HeaderLocation":  "(top left of the slide)",
		"FooterLocation":  "(bottom left of the slide)",
		"ThemeDescription": "Graph Paper is a light color scheme.",
	},
}

// GenerateSlidePrompt creates a prompt for slide generation based on the given parameters
func GenerateSlidePrompt(theme string, settings models.SlideSettings) (string, error) {
	// Generate theme example
	themeExample, err := generateThemeExample(theme)
	if err != nil {
		return "", err
	}

	detailPrompt := ""
	if settings.SlideDetail == "detailed" {
		detailPrompt = "Extract comprehensive content from the document, preserving all key information and supporting details. Include all major sections and subsections from the source material, maintaining the depth of explanations, examples, data points, and contextual information. Create sufficient slides to accommodate all relevant content without crowding. For each topic in the source document, extract both main points and their supporting evidence or explanations. Ensure visual balance by limiting each slide to 6-8 bullet points or a comparable amount of content. Do not overflow individual slides with too much information or they will go off the slide."
	} else if settings.SlideDetail == "medium" {
		detailPrompt = "Extract the most significant information from each section of the document, focusing on main concepts and key supporting details. Select content that represents the core message and essential evidence without including every example or minor point from the source material. Consolidate related information into coherent slides, aiming for comprehensive coverage of major topics while omitting supplementary details. Prioritize information that directly supports the document's main arguments or conclusions. Limit each slide to 4-6 bullet points or a comparable amount of content."
	} else if settings.SlideDetail == "minimal" {
		detailPrompt = "Extract only the most essential information from the document, focusing exclusively on key conclusions, main arguments, and critical data points. Select content that communicates the core message in the most concise form possible. Consolidate major sections of the document into a limited number of focused slides. Omit supporting details, examples, and explanations unless absolutely necessary for basic comprehension. Prioritize high-level takeaways over process explanations or contextual information. Limit each slide to 3-4 bullet points or a comparable amount of content."
	}

	audiencePrompt := ""
	if settings.Audience == "general" {
		audiencePrompt = "Format the presentation for a general audience with varying levels of background knowledge. Select the clearest and most accessible language from the document. When technical terms appear in the source, include brief definitions from the document when available. Prioritize content from the document that explains broader context and significance. Organize the extracted information as a narrative when possible, with a clear beginning, middle, and end. Format slides with minimal text and emphasize any visual elements from the original document."
	} else if settings.Audience == "academic" {
		audiencePrompt = "Format the presentation for an academic audience. Select terminology and detailed explanations from the document that preserve methodological details and theoretical frameworks. When extracting content, maintain the document's original citations, methodologies, and nuanced points. Preserve the logical structure of arguments found in the source material. When organizing information from the document, maintain appropriate context for all extracted data and findings. Format slides to balance detailed information with clarity."
	} else if settings.Audience == "technical" {
		audiencePrompt = "Format the presentation for a technical audience. Preserve technical terminology, specifications, and detailed explanations from the document. Prioritize content that focuses on implementation details, methodologies, and technical processes described in the source material. When extracting diagrams or code examples from the document, include the relevant explanatory text. Maintain the technical depth and precision of the source material. Organize the content in a logical sequence that preserves technical relationships and dependencies described in the document."
	} else if settings.Audience == "professional" {
		audiencePrompt = "Format the presentation for business professionals. Select terminology and concepts from the document that highlight practical applications and business relevance. Prioritize content from the document that demonstrates actionable insights, case studies, and results. Organize the extracted information with an emphasis on takeaways and strategic implications. Format slide content with concise bullet points rather than dense paragraphs. When selecting information from charts or data in the document, focus on metrics and trends most relevant to business decisions."
	} else if settings.Audience == "executive" {
		audiencePrompt = "Format the presentation for executive decision-makers. Select high-level information from the document that focuses on strategic implications and business impact. Prioritize content related to outcomes, ROI, and competitive advantages mentioned in the source material. Extract summary information rather than operational details unless specifically relevant to executive decisions. When selecting information from the document, focus on big-picture insights and key recommendations. Format slides with concise headline statements that capture the essential points from the document."
	}

	// Create template data
	data := map[string]interface{}{
		"Theme":        theme,
		"ThemeExample": themeExample,
		"DetailLevel":  detailPrompt,
		"Audience":     audiencePrompt,
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

// generateThemeExample generates an example for a specific theme
func generateThemeExample(theme string) (string, error) {
	// Get theme configuration or use default config if theme doesn't exist
	themeConfig, exists := themeConfigs[theme]
	if !exists {
		themeConfig = themeConfigs["default"]
	}
	
	// Copy the theme config and add the theme name
	templateData := make(map[string]interface{})
	for k, v := range themeConfig {
		templateData[k] = v
	}
	templateData["Theme"] = theme

	// Generate the header
	headerTemplate, err := template.New("header").Parse(commonMarpHeader)
	if err != nil {
		return "", err
	}
	
	var headerBuf bytes.Buffer
	if err := headerTemplate.Execute(&headerBuf, templateData); err != nil {
		return "", err
	}
	
	// Generate the body
	bodyTemplate, err := template.New("body").Parse(commonExampleBody)
	if err != nil {
		return "", err
	}
	
	var bodyBuf bytes.Buffer
	if err := bodyTemplate.Execute(&bodyBuf, templateData); err != nil {
		return "", err
	}
	
	// Combine the parts into a complete example
	example := "```md\n" + headerBuf.String() + bodyBuf.String() + "\n```"
	
	return example, nil
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