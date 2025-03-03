package models

// SlideSettings represents the settings for slide generation
type SlideSettings struct {
	SlideDetail string `json:"slideDetail"` // Values: minimal, medium, detailed
	Audience    string `json:"audience"`    // Values: general, academic, technical, professional, executive
} 

type File struct {
	Filename string `json:"filename"`
	Data []byte `json:"data"`
	Type string `json:"type"`
}
