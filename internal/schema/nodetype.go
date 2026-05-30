package schema

import (
	"encoding/json"
	"fmt"
)

const (
	NodeInvalid NodeType = iota
	NodeDocument
	NodeText
	NodeLineBreak
	NodeThematicBreak
	NodeHeading
	NodeParagraph
	NodeBlockQuote
	NodeNotice
	NodeMention
	NodeCodeBlock
	NodeBulletList
	NodeOrderedList
	NodeListItem
	NodeChecklist
	NodeChecklistItem
	NodeTable
	NodeTableRow
	NodeTableHeader
	NodeTableCell
)

var nodeTypeToString = map[NodeType]string{
	NodeDocument:      "doc",
	NodeText:          "text",
	NodeLineBreak:     "br",
	NodeThematicBreak: "hr",
	NodeHeading:       "heading",
	NodeParagraph:     "paragraph",
	NodeBlockQuote:    "blockquote",
	NodeNotice:        "container_notice",
	NodeMention:       "mention",
	NodeCodeBlock:     "code_block",
	NodeBulletList:    "bullet_list",
	NodeOrderedList:   "ordered_list",
	NodeListItem:      "list_item",
	NodeChecklist:     "checkbox_list",
	NodeChecklistItem: "checkbox_item",
	NodeTable:         "table",
	NodeTableRow:      "tr",
	NodeTableHeader:   "th",
	NodeTableCell:     "td",
}

var stringToNodeType = map[string]NodeType{}

func init() {
	for k, v := range nodeTypeToString {
		stringToNodeType[v] = k
	}
}

// Represents the type of a Prosemirror node.
type NodeType int

func (t NodeType) String() string {
	return nodeTypeToString[t]
}

func (t NodeType) MarshalJSON() ([]byte, error) {
	s, ok := nodeTypeToString[t]
	if !ok {
		return nil, fmt.Errorf("unknown node type: %d", t)
	}
	return json.Marshal(s)
}

func (t *NodeType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	val, ok := stringToNodeType[s]
	if !ok {
		return fmt.Errorf("unknown node type string: %s", s)
	}
	*t = val
	return nil
}
