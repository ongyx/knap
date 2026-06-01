package schema

import (
	"github.com/google/uuid"
)

// Node represents a Prosemirror node defined by Outline's schema.
type Node struct {
	Type    NodeType       `json:"type"`
	Text    string         `json:"text,omitempty"`
	Attrs   map[string]any `json:"attrs,omitempty"`
	Content []*Node        `json:"content,omitempty"`
	Marks   []Mark         `json:"marks,omitempty"`
}

// Checks if the node is invalid.
func (n *Node) IsInvalid() bool {
	return n.Type == NodeInvalid
}

// Creates a document node.
func NewDocumentNode() *Node {
	return &Node{Type: NodeDocument}
}

// Creates a text node with the given text.
func NewTextNode(text string) *Node {
	return &Node{Type: NodeText, Text: text}
}

// Creates a line break node.
func NewLineBreakNode() *Node {
	return &Node{Type: NodeLineBreak}
}

// Creates a thematic break node.
// If isPageBreak is true, the line appears as a page break.
func NewThematicBreakNode(isPageBreak bool) *Node {
	var markup string
	if isPageBreak {
		markup = "***"
	} else {
		markup = "---"
	}

	return &Node{Type: NodeThematicBreak, Attrs: map[string]any{"markup": markup}}
}

// Creates a heading node with the given level.
// level may be any number from 1 to 6.
func NewHeadingNode(level int) *Node {
	return &Node{
		Type: NodeHeading,
		Attrs: map[string]any{
			"level": level,
		},
	}
}

// Creates a paragraph node.
func NewParagraphNode() *Node {
	return &Node{Type: NodeParagraph}
}

// Creates a block quote node.
func NewBlockQuoteNode() *Node {
	return &Node{Type: NodeBlockQuote}
}

// Creates a notice block node with the given type and content.
func NewNoticeNode(nt NoticeType) *Node {
	return &Node{
		Type:  NodeNotice,
		Attrs: map[string]any{"style": nt.String()},
	}
}

// Creates a mention node with the given type, ID of the target user/document/collection, and ID of the author who wrote the mention.
func NewMentionNode(mt MentionType, target uuid.UUID, author uuid.UUID, label string) *Node {
	id, _ := uuid.NewRandom()

	return &Node{
		Type: NodeMention,
		Attrs: map[string]any{
			"type":    mt.String(),
			"label":   label,
			"modelId": target.String(),
			"actorId": author.String(),
			"id":      id.String(),
		},
	}
}

// Creates a fenced code block node with the given language and text.
// For plain text, language should be set to "none".
func NewFencedCodeBlockNode(language string) *Node {
	return &Node{
		Type: NodeCodeBlock,
		Attrs: map[string]any{
			"language": language,
			"wrap":     false,
		},
	}
}

// Creates a bullet list node.
func NewBulletListNode() *Node {
	return &Node{Type: NodeBulletList}
}

// Creates an ordered list node with a starting number.
func NewOrderedListNode(start int) *Node {
	return &Node{Type: NodeOrderedList, Attrs: map[string]any{"order": start, "listStyle": "number"}}
}

// Creates a list item node.
func NewListItemNode() *Node {
	return &Node{Type: NodeListItem}
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

	return &Node{Type: NodeChecklist, Attrs: attrs}
}

// Creates a checklist item node.
func NewChecklistItemNode(isChecked bool) *Node {
	return &Node{Type: NodeChecklistItem, Attrs: map[string]any{"checked": isChecked}}
}

// Creates a table node.
func NewTableNode() *Node {
	return &Node{
		Type: NodeTable,
		Attrs: map[string]any{
			"layout": nil,
		},
	}
}

// Creates a table row node.
func NewTableRowNode() *Node {
	return &Node{Type: NodeTableRow}
}

// Creates a table cell node.
// If header is true, the type is set to NodeTableHeader.
func NewTableCellNode(isHeader bool) *Node {
	var ty NodeType
	if isHeader {
		ty = NodeTableHeader
	} else {
		ty = NodeTableCell
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
