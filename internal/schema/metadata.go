package schema

import (
	"time"

	"github.com/djherbis/times"
	"github.com/google/uuid"
)

// Metadata represents common metadata for Document and Collection.
type Metadata struct {
	// The UUID of the document/collection.
	ID uuid.UUID `json:"id"`
	// The URLID linking to the document/collection.
	URLID URLID `json:"urlId"`
	// The name of the icon to show for the document/collection, if any.
	Icon *string `json:"icon"`
	// The hexadecimal color of the icon for the document/collection, if any.
	Color *string `json:"color"`
	// When the document/collection was created.
	CreatedAt time.Time `json:"createdAt"`
	// When the document/collection was last updated.
	UpdatedAt time.Time `json:"updatedAt"`
	// When the document/collection was deleted, if it is deleted.
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

// Creates basic metadata with the given ID.
func NewMetadata(id uuid.UUID) *Metadata {
	urlid := NewURLID()
	now := time.Now()

	return &Metadata{
		ID:        id,
		URLID:     urlid,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Sets the timestamps for the *At fields from a file.
func (m *Metadata) SetTimestamps(filename string) error {
	t, err := times.Stat(filename)
	if err != nil {
		return err
	}

	mtime := t.ModTime()
	// Getting creation time is surprisingly hard...
	if t.HasBirthTime() {
		m.CreatedAt = t.BirthTime()
	} else {
		m.CreatedAt = mtime
	}
	m.UpdatedAt = mtime

	return nil
}
