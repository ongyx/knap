package schema

import (
	"time"

	"github.com/google/uuid"
	"github.com/ongyx/knap/internal/prosemirror"
	"github.com/ongyx/knap/internal/util"
)

// Document represents an individual Outline document, including its content and metadata.
type Document struct {
	*BaseMetadata

	// The display name.
	Title string `json:"title"`
	// The document schema node.
	Data *prosemirror.Node `json:"data"`
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

// Creates an empty document with the given ID, title, and identity. [Document.Data] is set to an empty document node.
func NewDocument(id uuid.UUID, urlid util.URLID, title string, idn Identity) *Document {
	m := NewCommonMetadata(id, urlid)
	d := prosemirror.NewDocumentNode()
	d.Content = append(d.Content, prosemirror.NewParagraphNode())

	return &Document{
		BaseMetadata:   m,
		Title:          title,
		Data:           d,
		CreatedById:    idn.ID,
		CreatedByName:  idn.Name,
		CreatedByEmail: idn.Email,
		PublishedAt:    &m.CreatedAt,
	}
}

// Sets the timestamps for the *At fields from a file.
func (d *Document) SetTimestamps(filename string) error {
	if err := d.BaseMetadata.SetTimestamps(filename); err != nil {
		return err
	}

	d.PublishedAt = &d.BaseMetadata.UpdatedAt
	return nil
}
