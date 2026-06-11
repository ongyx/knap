package schema

import (
	"time"

	"github.com/djherbis/times"
	"github.com/google/uuid"
	"github.com/ongyx/knap/internal/util"
)

// BaseMetadata represents common metadata for Document and Collection.
type BaseMetadata struct {
	// The UUID of the document/collection.
	ID uuid.UUID `json:"id"`
	// The URLID linking to the document/collection.
	URLID util.URLID `json:"urlId"`
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

// Creates basic metadata with the given ID and URLID.
func NewCommonMetadata(id uuid.UUID, urlid util.URLID) *BaseMetadata {
	now := time.Now()

	return &BaseMetadata{
		ID:        id,
		URLID:     urlid,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Sets the timestamps for the *At fields from a file.
func (m *BaseMetadata) SetTimestamps(filename string) error {
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
