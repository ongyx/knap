package schema

import (
	"time"

	"github.com/google/uuid"
)

// Lowest version supported.
const outlineVersion = "1.7.1"

// ExportMetadata contains metadata on an export.
//
// Source: https://github.com/outline/outline/blob/a25f334bb18a339b29f89a63c137a5b785e17bc8/server/types.ts#L528
type ExportMetadata struct {
	// The export version.
	ExportVersion int `json:"exportVersion"`
	// The application version.
	Version string `json:"version"`
	// When the export was created.
	CreatedAt time.Time `json:"createdAt"`
	// The ID of the user who created the export.
	CreatedByID uuid.UUID `json:"createdById"`
	// The email of the user who created the export.
	CreatedByEmail *string `json:"createdByEmail"`
}

func NewExportMetadata(idn Identity) *ExportMetadata {
	return &ExportMetadata{
		ExportVersion:  1,
		Version:        outlineVersion,
		CreatedAt:      time.Now(),
		CreatedByID:    idn.ID,
		CreatedByEmail: &idn.Email,
	}
}
