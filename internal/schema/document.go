package schema

import (
	"os"
	"time"

	"github.com/djherbis/times"
	"github.com/google/uuid"
)

// Document represents an individual Outline document, including its content and metadata.
type Document struct {
	Metadata
	Title            string     `json:"title"`
	Data             *Node      `json:"data"`
	CreatedById      string     `json:"createdById"`
	CreatedByName    string     `json:"createdByName"`
	CreatedByEmail   string     `json:"createdByEmail"`
	PublishedAt      *time.Time `json:"publishedAt"`
	FullWidth        bool       `json:"fullWidth"`
	Template         bool       `json:"template"`
	ParentDocumentId *string    `json:"parentDocumentId"`
}

// Creates an empty document.
func NewDocument() *Document {
	id, _ := uuid.NewRandom()
	urlid := NewURLID()

	now := time.Now()
	return &Document{
		Metadata: Metadata{
			ID:        id,
			URLID:     urlid,
			CreatedAt: now,
			UpdatedAt: now,
		},
		PublishedAt: &now,
	}
}

// Sets the document timestamps from a file's metadata.
func (d *Document) SetTimestamps(f *os.File) error {
	// The reader is a file, query the file's timestamps to fill in the *At values.
	// Getting creation time is surprisingly hard...
	t, err := times.Stat(f.Name())
	if err != nil {
		return err
	}

	mtime := t.ModTime()
	if t.HasBirthTime() {
		d.CreatedAt = t.BirthTime()
	} else {
		d.CreatedAt = mtime
	}
	d.UpdatedAt = mtime
	d.PublishedAt = &mtime

	return nil
}
