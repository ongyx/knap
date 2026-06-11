package schema

import (
	"fmt"

	"github.com/google/uuid"
)

// Attachment represents attachment metadata.
type Attachment struct {
	// The ID of the user who uploaded the attachment.
	UserID uuid.UUID `json:"userId"`
	// The ID of the document where the attachment is embedded.
	DocumentID uuid.UUID `json:"documentId"`
	// The MIME content type of the attachment.
	ContentType string `json:"contentType"`
	// The filename of the attachment.
	Name string `json:"name"`
	// The ID of the attachment.
	ID uuid.UUID `json:"id"`
	// The size of the attachment in bytes.
	Size int64 `json:"size"`
	// The path to the attachment, in the format 'uploads/(UserID)/(ID)/(Name)'.
	Key string `json:"key"`
}

// Updates the attachment key.
func (a *Attachment) UpdateKey() {
	a.Key = fmt.Sprintf("uploads/%s/%s/%s", a.UserID, a.ID, a.Name)
}
