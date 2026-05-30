package schema

// Mark represents formatting applied to a text node.
type Mark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

// Creates a bold mark (<strong>...</strong>).
func NewBoldMark() Mark {
	return Mark{Type: "strong"}
}

// Creates an italic mark (<em>...</em>).
func NewItalicMark() Mark {
	return Mark{Type: "em"}
}

// Creates a strikethrough mark (<del>...</del>).
func NewStrikethroughMark() Mark {
	return Mark{Type: "strikethrough"}
}

// Creates a link mark (<a href="..."></a>).
func NewLinkMark(url string) Mark {
	return Mark{Type: "link", Attrs: map[string]any{"href": url, "title": nil}}
}

// Creates an inline code mark (<code>...</code>).
func NewInlineCodeMark() Mark {
	return Mark{Type: "code_inline"}
}

// Creates a highlight mark (<mark style="background-color: ..."></mark>).
// colorHex must be a hexadecimal color, e.g., #FFF8E7. If nil, the highlight appears yellow in Outline.
func NewHighlightMark(colorHex *string) Mark {
	return Mark{Type: "highlight", Attrs: map[string]any{
		"color": colorHex,
	}}
}
