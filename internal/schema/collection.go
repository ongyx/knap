package schema

import (
	"github.com/google/uuid"
	"github.com/ongyx/knap/internal/util"
)

// Collection represents an exported collection from Outline.
type Collection struct {
	Metadata    *CollectionMetadata       `json:"collection"`
	Documents   map[uuid.UUID]*Document   `json:"documents"`
	Attachments map[uuid.UUID]*Attachment `json:"attachments"`
}

// Creates a new export with an ID, URLID, and name for the collection.
func NewCollection(id uuid.UUID, urlid util.URLID, name string) *Collection {
	return &Collection{
		Metadata:    NewCollectionMetadata(id, urlid, name),
		Documents:   make(map[uuid.UUID]*Document),
		Attachments: make(map[uuid.UUID]*Attachment),
	}
}
