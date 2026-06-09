package schema

import (
	"time"

	"github.com/google/uuid"
)

// Collection represents the collection metadata and its hierarchical document structure.
type Collection struct {
	*Metadata

	// The name of the collection.
	Name string `json:"name"`
	// The root node of the collection's welcome page.
	Data *Node `json:"data"`
	// The sorting criteria.
	Sort Sort `json:"sort"`
	// The sort index.
	Index string `json:"index"`
	// The access permissions.
	Permission Permission `json:"permission"`
	// Whether or not comments are enabled on this collection.
	Commenting any `json:"commenting"`
	// Whether or not sharing is enabled on this collection.
	Sharing bool `json:"sharing"`
	// When the collection was archived, if it was archived.
	ArchivedAt *time.Time `json:"archivedAt"`
	// The hierarchy of documents in this collection.
	DocumentStructure []Structure `json:"documentStructure"`
}

// Creates an empty collection with the given ID.
func NewCollection(id uuid.UUID) *Collection {
	m := NewMetadata(id)

	return &Collection{Metadata: m}
}

// Sort defines the sorting criteria for documents within a collection.
type Sort struct {
	// The field to sort with.
	Field string `json:"field"`
	// The direction to sort in.
	Direction SortDirection `json:"direction"`
}

// Structure represents a document in the collection's hierarchy.
type Structure struct {
	// The UUID of the document.
	ID uuid.UUID `json:"id"`
	// The relative URL to the document.
	URL string `json:"url"`
	// The title of the document.
	Title string `json:"title"`
	// The child documents organized under this document.
	Children []Structure `json:"children"`
}
