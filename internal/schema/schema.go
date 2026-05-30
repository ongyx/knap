package schema

import (
	"time"

	"github.com/google/uuid"
)

// Schema represents the root structure of the exported JSON file from Outline.
type Schema struct {
	Collection  Collection            `json:"collection"`
	Documents   map[string]Document   `json:"documents"`
	Attachments map[string]Attachment `json:"attachments"`
}

// Collection represents the collection metadata and its hierarchical document structure.
type Collection struct {
	ID                uuid.UUID       `json:"id"`
	URLID             URLID           `json:"urlId"`
	Name              string          `json:"name"`
	Data              Node            `json:"data"`
	Sort              Sort            `json:"sort"`
	Icon              *string         `json:"icon"`
	Index             string          `json:"index"`
	Color             *string         `json:"color"`
	Permission        string          `json:"permission"`
	Commenting        any             `json:"commenting"`
	Sharing           bool            `json:"sharing"`
	CreatedAt         time.Time       `json:"createdAt"`
	UpdatedAt         time.Time       `json:"updatedAt"`
	DeletedAt         *time.Time      `json:"deletedAt"`
	ArchivedAt        *time.Time      `json:"archivedAt"`
	DocumentStructure []StructureNode `json:"documentStructure"`
}

// Sort defines the sorting criteria for documents within a collection.
type Sort struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

// StructureNode represents a node in the collection's document hierarchy.
type StructureNode struct {
	ID       uuid.UUID       `json:"id"`
	URL      string          `json:"url"`
	Title    string          `json:"title"`
	Children []StructureNode `json:"children"`
}

// Attachment represents attachment metadata.
type Attachment map[string]any
