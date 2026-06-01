package obsidian

import (
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// interface asserts
var _ ast.Node = (*Callout)(nil)
var _ parser.InlineParser = (*calloutParser)(nil)

// Matches a color directive. First group is the color name.
// NOTE: Color names must be at least one character and have no space:
// https://github.com/Superschnizel/obsidian-fast-text-color/blob/master/src/utils/validateColorName.ts
var reColor = regexp.MustCompile(`^\{(\w+)\}\s*`)

// The text color node kind.
var KindTextColor = ast.NewNodeKind("TextColor")

// Represents a color directive within a text color span ('~={color}text=~').
type TextColor struct {
	ast.BaseInline
	// The name of the color, e.g., 'yellow'.
	Name []byte
}

// Creates a new text color node with the given color name.
func NewTextColor(name []byte) *TextColor {
	return &TextColor{Name: name}
}

func (n *TextColor) Kind() ast.NodeKind {
	return KindTextColor
}

func (n *TextColor) Dump(source []byte, level int) {
	m := map[string]string{
		"Name": string(n.Name),
	}
	ast.DumpHelper(n, source, level, m, nil)
}

type textColorParser struct{}

// Creates a new text color parser.
func NewTextColorParser() parser.InlineParser {
	return &textColorParser{}
}

func (t *textColorParser) Trigger() []byte {
	return []byte{'{'}
}

func (t *textColorParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	// The text color must appear after a TextColorSpan delimiter.
	// Delimiters are only processed after all blocks are parsed, so we can't check for the TextColorSpan node directly.
	// https://github.com/yuin/goldmark/#overview
	d, ok := parent.LastChild().(*parser.Delimiter)
	if !ok || d.Char != '~' || d.OriginalLength != 2 {
		return nil
	}

	line, _ := block.PeekLine()
	m := reColor.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}

	name := line[m[2]:m[3]]
	block.Advance(m[1])

	return NewTextColor(name)
}
