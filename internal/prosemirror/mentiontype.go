package prosemirror

const (
	// The mention targets a user.
	MentionUser MentionType = iota
	// The mention targets a document.
	MentionDocument
	// The mention targets a collection.
	MentionCollection
)

// Represents a type of mention in a document.
// NOTE: This type does not support JSON marshalling as it is stored as a node attribute.
type MentionType int

// Returns the string representation of the mention type.
func (mt MentionType) String() string {
	switch mt {
	case MentionUser:
		return "user"
	case MentionDocument:
		return "document"
	case MentionCollection:
		return "collection"
	default:
		return ""
	}
}
