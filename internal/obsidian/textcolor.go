package obsidian

import (
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// interface asserts
var (
	_ ast.Node            = (*Callout)(nil)
	_ parser.InlineParser = (*calloutParser)(nil)
)

// Matches a span of colored text, where group 1 is the color ID and group 2 is the text.
// Color IDs must be at least one character and have no whitespace.
// https://github.com/Superschnizel/obsidian-fast-text-color/blob/master/src/utils/validateColorName.ts
var reTextColor = regexp.MustCompile(`^~=\{(\w+)\}\s*(.+?)=~`)

// The text color node kind.
var KindTextColor = ast.NewNodeKind("TextColor")

// Represents a text color ('~={color}text=~').
type TextColor struct {
	ast.BaseInline
	// The ID of the color, e.g., 'yellow'.
	ID []byte
}

// Creates a new text color node with the given color ID.
func NewTextColor(id []byte) *TextColor {
	return &TextColor{ID: id}
}

func (n *TextColor) Kind() ast.NodeKind {
	return KindTextColor
}

func (n *TextColor) Dump(source []byte, level int) {
	m := map[string]string{
		"ID": string(n.ID),
	}
	ast.DumpHelper(n, source, level, m, nil)
}

type textColorParser struct{}

// Creates a new text color parser.
func NewTextColorParser() parser.InlineParser {
	return &textColorParser{}
}

func (t *textColorParser) Trigger() []byte {
	return []byte{'~'}
}

func (t *textColorParser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, seg := block.PeekLine()

	m := reTextColor.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}

	id := line[m[2]:m[3]]
	seg = text.NewSegment(seg.Start+m[4], seg.Start+m[5])

	tc := NewTextColor(id)
	tc.AppendChild(tc, ast.NewTextSegment(seg))
	block.Advance(m[1])
	return tc
}
