package obsidian

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// The highlight node kind.
var KindHighlight = ast.NewNodeKind("Highlight")

// interface asserts
var _ ast.Node = (*Highlight)(nil)
var _ parser.DelimiterProcessor = (*highlightDelimiterParser)(nil)
var _ parser.InlineParser = (*highlightParser)(nil)

// Represents a highlight ('==(text)==').
type Highlight struct {
	ast.BaseInline
}

// Creates a new highlight node.
func NewHighlight() *Highlight {
	return &Highlight{}
}

func (n *Highlight) Kind() ast.NodeKind {
	return KindHighlight
}

func (n *Highlight) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

type highlightDelimiterParser struct{}

func (p *highlightDelimiterParser) IsDelimiter(b byte) bool {
	return b == '='
}

func (p *highlightDelimiterParser) CanOpenCloser(opener, closer *parser.Delimiter) bool {
	return opener.Char == closer.Char
}

func (p *highlightDelimiterParser) OnMatch(consumes int) ast.Node {
	return NewHighlight()
}

type highlightParser struct{}

// Creates a new highlight parser.
func NewHighlightParser() parser.InlineParser {
	return &highlightParser{}
}

func (h *highlightParser) Trigger() []byte {
	return []byte{'='}
}

func (h *highlightParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	// Based on https://github.com/yuin/goldmark/blob/master/extension/strikethrough.go
	before := block.PrecendingCharacter()
	line, segment := block.PeekLine()

	node := parser.ScanDelimiter(line, before, 2, &highlightDelimiterParser{})
	if node == nil || before == '=' {
		return nil
	}

	node.Segment = segment.WithStop(segment.Start + node.OriginalLength)
	block.Advance(node.OriginalLength)
	pc.PushDelimiter(node)
	return node
}
