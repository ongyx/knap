package schema

import (
	"time"

	"github.com/google/uuid"
)

// Collection represents the collection metadata and its hierarchical document structure.
type Collection struct {
	Metadata
	Name              string      `json:"name"`
	Data              Node        `json:"data"`
	Sort              Sort        `json:"sort"`
	Index             string      `json:"index"`
	Permission        string      `json:"permission"`
	Commenting        any         `json:"commenting"`
	Sharing           bool        `json:"sharing"`
	ArchivedAt        *time.Time  `json:"archivedAt"`
	DocumentStructure []Structure `json:"documentStructure"`
}

// Sort defines the sorting criteria for documents within a collection.
type Sort struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

// Structure represents a document in the collection's hierarchy.
type Structure struct {
	ID       uuid.UUID   `json:"id"`
	URL      string      `json:"url"`
	Title    string      `json:"title"`
	Children []Structure `json:"children"`
}
