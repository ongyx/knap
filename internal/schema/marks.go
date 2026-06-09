package schema

import "regexp"

var reHexColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// Mark represents formatting applied to a text node.
type Mark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

// Creates a bold mark.
func NewBoldMark() Mark {
	return Mark{Type: "strong"}
}

// Creates an italic mark.
func NewItalicMark() Mark {
	return Mark{Type: "em"}
}

// Creates a strikethrough mark.
func NewStrikethroughMark() Mark {
	return Mark{Type: "strikethrough"}
}

// Creates a link mark.
func NewLinkMark(href string) Mark {
	return Mark{Type: "link", Attrs: map[string]any{"href": href, "title": nil}}
}

// Creates an inline code mark.
func NewInlineCodeMark() Mark {
	return Mark{Type: "code_inline"}
}

// Creates a highlight mark.
// colorHex must be a hexadecimal color, e.g., #FFF8E7. If empty or invalid, the highlight defaults to yellow in Outline.
func NewHighlightMark(hexColor string) Mark {
	// NOTE: This must be a string pointer as nil is a valid value for the color attribute.
	var hc *string
	if reHexColor.MatchString(hexColor) {
		hc = &hexColor
	}

	return Mark{Type: "highlight", Attrs: map[string]any{
		"color": hc,
	}}
}
