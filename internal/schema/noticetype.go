package schema

const (
	// Info notice.
	NoticeInfo NoticeType = iota
	// Success notice.
	NoticeSuccess
	// Tip notice.
	NoticeTip
	// Warning notice.
	NoticeWarning
)

// Represents the type of a notice in a document.
// NOTE: This type does not support JSON marshalling as it is stored as a node attribute.
type NoticeType int

// Returns the string representation of the notice type.
func (n NoticeType) String() string {
	switch n {
	case NoticeInfo:
		return "info"
	case NoticeSuccess:
		return "success"
	case NoticeTip:
		return "tip"
	case NoticeWarning:
		return "warning"
	default:
		return ""
	}
}
