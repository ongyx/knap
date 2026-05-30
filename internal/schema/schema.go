package schema

// Schema represents the root structure of the exported JSON file from Outline.
type Schema struct {
	Collection  Collection            `json:"collection"`
	Documents   map[string]Document   `json:"documents"`
	Attachments map[string]Attachment `json:"attachments"`
}

// Attachment represents attachment metadata.
type Attachment map[string]any
