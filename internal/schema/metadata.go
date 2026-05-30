package schema

import (
	"time"

	"github.com/google/uuid"
)

// Metadata represents common metadata for Document and Collection.
type Metadata struct {
	ID        uuid.UUID  `json:"id"`
	URLID     URLID      `json:"urlId"`
	Icon      *string    `json:"icon"`
	Color     *string    `json:"color"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}
