package schema

import (
	"time"

	"github.com/google/uuid"
	"github.com/ongyx/knap/internal/prosemirror"
	"github.com/ongyx/knap/internal/util"
)

// CollectionMetadata contains metadata on an exported collection.
type CollectionMetadata struct {
	*BaseMetadata

	// The name of the collection.
	Name string `json:"name"`
	// The document schema node representing the welcome page.
	Data *prosemirror.Node `json:"data"`
	// The sorting criteria.
	Sort Sort `json:"sort"`
	// The sort index.
	Index string `json:"index"`
	// The access permissions.
	Permission Permission `json:"permission"`
	// Whether or not comments are enabled on this collection.
	Commenting any `json:"commenting"`
	// Whether or not sharing is enabled on this collection.
	Sharing bool `json:"sharing"`
	// The permission level required for managing templates in this collection.
	TemplateManagement Permission `json:"templateManagement"`
	// When the collection was archived, if it was archived.
	ArchivedAt *time.Time `json:"archivedAt"`
	// The hierarchy of documents in the collection.
	DocumentStructure []*NavigationNode `json:"documentStructure"`
}

// Creates an empty collection with the given ID, URLID, and name.
func NewCollectionMetadata(id uuid.UUID, urlid util.URLID, name string) *CollectionMetadata {
	m := NewCommonMetadata(id, urlid)
	d := prosemirror.NewDocumentNode()
	d.Content = append(d.Content, prosemirror.NewParagraphNode())

	return &CollectionMetadata{
		BaseMetadata:       m,
		Name:               name,
		Data:               d,
		Sort:               NewSort(),
		Index:              ",",
		Sharing:            true,
		Commenting:         true,
		TemplateManagement: PermissionAdmin,
	}
}
