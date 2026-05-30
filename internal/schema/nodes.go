package schema

import "github.com/google/uuid"

// Node represents a Prosemirror node defined by Outline's schema.
type Node struct {
	Type    string         `json:"type"`
	Text    string         `json:"text,omitempty"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Content []*Node        `json:"content,omitempty"`
	Marks   []Mark         `json:"marks,omitempty"`
}

// Checks if the node is invalid.
func (n *Node) IsInvalid() bool {
	return n.Type == ""
}

// Creates a document node.
func NewDocumentNode() *Node {
	return &Node{Type: "doc"}
}

// Creates a text node with the given text.
func NewTextNode(text string) *Node {
	return &Node{Type: "text", Text: text}
}

// Creates a line break node (<br>).
func NewLineBreakNode() *Node {
	return &Node{Type: "br"}
}

// Creates a thematic break node (<hr>).
// If isPageBreak is true, the line appears as a page break.
func NewThematicBreakNode(isPageBreak bool) *Node {
	var markup string
	if isPageBreak {
		markup = "***"
	} else {
		markup = "---"
	}

	return &Node{Type: "hr", Attrs: map[string]any{"markup": markup}}
}

// Creates a heading node (<h1>, <h2>, <h3>, etc.) with the given level.
// level may be any number from 1 to 6.
func NewHeadingNode(level int) *Node {
	return &Node{
		Type: "heading",
		Attrs: map[string]any{
			"level": level,
		},
	}
}

// Creates a paragraph node (<p>...</p>).
func NewParagraphNode() *Node {
	return &Node{Type: "paragraph"}
}

// Creates a block quote node (<blockquote>...</blockquote>).
func NewBlockQuoteNode() *Node {
	return &Node{Type: "blockquote"}
}

// Creates a notice block node with the given type and content.
func NewNoticeNode(nt NoticeType) *Node {
	return &Node{
		Type:  "container_notice",
		Attrs: map[string]any{"style": nt.String()},
	}
}

// Creates a mention node with the given type, ID of the target user/document/collection, and ID of the author who wrote the mention.
func NewMentionNode(mt MentionType, target uuid.UUID, author uuid.UUID, label string) *Node {
	id, _ := uuid.NewRandom()

	return &Node{
		Type: "mention",
		Attrs: map[string]any{
			"type":    mt.String(),
			"label":   label,
			"modelId": target.String(),
			"actorId": author.String(),
			"id":      id.String(),
		},
	}
}

// Creates a fenced code block node (<pre><code>...</code></pre>) with the given language and text.
// For plain text, language should be set to "none".
func NewFencedCodeBlockNode(language string) *Node {
	return &Node{
		Type: "code_block",
		Attrs: map[string]any{
			"language": language,
			"wrap":     false,
		},
	}
}

// Creates a bullet list node (<ul>...</ul>).
func NewBulletListNode() *Node {
	return &Node{Type: "bullet_list"}
}

// Creates an ordered list node (<ol>...</ol>) with a starting number.
func NewOrderedListNode(start int) *Node {
	return &Node{Type: "ordered_list", Attrs: map[string]any{"order": start, "listStyle": "number"}}
}

// Creates a list item node (<li>...</li>).
func NewListItemNode() *Node {
	return &Node{Type: "list_item"}
}

// Creates a checklist node.
func NewChecklistNode() *Node {
	attrs := make(map[string]any)
	// Not sure why the checkbox list has an ID.
	if u, err := uuid.NewRandom(); err != nil {
		attrs["id"] = nil
	} else {
		attrs["id"] = u.String()
	}

	return &Node{Type: "checkbox_list", Attrs: attrs}
}

// Creates a checklist item node.
func NewChecklistItemNode(isChecked bool) *Node {
	return &Node{Type: "checkbox_item", Attrs: map[string]any{"checked": isChecked}}
}

// Creates a table node (<table>).
func NewTableNode() *Node {
	return &Node{
		Type: "table",
		Attrs: map[string]any{
			"layout": nil,
		},
	}
}

// Creates a table row node (<tr>).
func NewTableRowNode() *Node {
	return &Node{Type: "tr"}
}

// Creates a table cell node (<td>, <th>).
// If header is true, the type is set to 'th'.
func NewTableCellNode(isHeader bool) *Node {
	var ty string
	if isHeader {
		ty = "th"
	} else {
		ty = "td"
	}

	return &Node{
		Type: ty,
		Attrs: map[string]any{
			"colspan":   1,
			"rowspan":   1,
			"alignment": "",
			"colwidth":  nil,
		},
	}
}
