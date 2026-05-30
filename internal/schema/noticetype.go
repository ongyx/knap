package schema

import "github.com/ongyx/knap/internal/obsidian"

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

// Maps callouts to notice types.
// https://obsidian.md/help/callouts
var calloutToNotice = map[string]NoticeType{
	"note":      NoticeInfo,
	"abstract":  NoticeInfo,
	"summary":   NoticeInfo,
	"tldr":      NoticeInfo,
	"info":      NoticeInfo,
	"todo":      NoticeInfo,
	"tip":       NoticeTip,
	"hint":      NoticeTip,
	"important": NoticeTip,
	"success":   NoticeSuccess,
	"check":     NoticeSuccess,
	"done":      NoticeSuccess,
	"question":  NoticeInfo,
	"help":      NoticeInfo,
	"faq":       NoticeInfo,
	"warning":   NoticeWarning,
	"caution":   NoticeWarning,
	"attention": NoticeWarning,
	"failure":   NoticeWarning,
	"fail":      NoticeWarning,
	"missing":   NoticeWarning,
	"danger":    NoticeWarning,
	"error":     NoticeWarning,
	"bug":       NoticeWarning,
	"example":   NoticeInfo,
	"quote":     NoticeInfo,
	"cite":      NoticeInfo,
}

// Represents the type of a notice in a document.
// NOTE: This type does not support JSON marshalling as it is stored as a node attribute.
type NoticeType int

// Converts an Obsidian callout type to a notice type.
func CalloutToNotice(callout *obsidian.Callout) NoticeType {
	// OPT: This conversion does not allocate - https://github.com/golang/go/issues/3512
	return calloutToNotice[string(callout.Name)]
}

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
