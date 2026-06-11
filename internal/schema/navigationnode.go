package schema

import "github.com/google/uuid"

// NavigationNode refers to a document and its children within a collection's hierarchy.
type NavigationNode struct {
	// The UUID of the document.
	ID uuid.UUID `json:"id"`
	// The relative URL to the document.
	URL string `json:"url"`
	// The title of the document.
	Title string `json:"title"`
	// The child documents organized under this document. If not empty, this causes the document to act as a 'folder' in the Outline UI.
	Children []*NavigationNode `json:"children"`
}

// Creates a navigation node from a document. The document must have an ID.
func NewNavigationNode(d *Document) *NavigationNode {
	return &NavigationNode{
		ID:       d.ID,
		URL:      d.BaseMetadata.URLID.GenerateDocumentURL(d.Title),
		Title:    d.Title,
		Children: make([]*NavigationNode, 0),
	}
}
