package schema

import (
	"time"

	"github.com/google/uuid"
)

// Document represents an individual Outline document, including its content and metadata.
type Document struct {
	*Metadata

	// The document's display name.
	Title string `json:"title"`
	// The document's root node.
	Data *Node `json:"data"`
	// The UUID of the person who created the document.
	CreatedById uuid.UUID `json:"createdById"`
	// The name of the person who created the document.
	CreatedByName string `json:"createdByName"`
	// The email of the person who created the document.
	CreatedByEmail string `json:"createdByEmail"`
	// When the document was published.
	PublishedAt *time.Time `json:"publishedAt"`
	// Whether or not the document should be displayed with full width.
	FullWidth bool `json:"fullWidth"`
	// Whether or not this document is a template.
	Template bool `json:"template"`
	// The parent document's UUID, if any.
	ParentDocumentId *uuid.UUID `json:"parentDocumentId"`
}

// Creates an empty document with the given ID.
func NewDocument(id uuid.UUID) *Document {
	m := NewMetadata(id)

	return &Document{
		Metadata:    m,
		PublishedAt: &m.CreatedAt,
	}
}

// Sets the identity for the CreatedBy* fields.
func (d *Document) SetIdentity(i Identity) {
	d.CreatedById = i.ID
	d.CreatedByName = i.Name
	d.CreatedByEmail = i.Email
}

// Sets the timestamps for the *At fields from a file.
func (d *Document) SetTimestamps(filename string) error {
	if err := d.Metadata.SetTimestamps(filename); err != nil {
		return err
	}

	d.PublishedAt = &d.Metadata.UpdatedAt
	return nil
}
