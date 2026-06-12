package schema

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
)

// Attachment represents attachment metadata.
//
// It is recommended to use [NewAttachment] as it will generate [Attachment.Key] automatically. Otherwise, you must call [Attachment.UpdateKey] manually.
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
	Size int64 `json:"size,string"`
	// The path to the attachment, in the format 'uploads/(UserID)/(ID)/(Name)'.
	Key string `json:"key"`
	// The absolute path to the attachment, if any.
	AbsPath string `json:"-"`
}

// Creates a new attachment for a file. absPath must be an absolute OS path to the file.
func NewAttachment(userID, documentID uuid.UUID, contentType, absPath string, id uuid.UUID, size int64) *Attachment {
	att := &Attachment{
		UserID:      userID,
		DocumentID:  documentID,
		ContentType: contentType,
		Name:        filepath.Base(absPath),
		ID:          id,
		Size:        size,
		AbsPath:     absPath,
	}
	att.UpdateKey()
	return att
}

// Updates the attachment key.
func (att *Attachment) UpdateKey() {
	att.Key = fmt.Sprintf("uploads/%s/%s/%s", att.UserID, att.ID, att.Name)
}
