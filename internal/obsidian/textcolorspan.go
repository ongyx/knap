package obsidian

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// interface asserts
var _ ast.Node = (*TextColorSpan)(nil)
var _ parser.DelimiterProcessor = (*textColorSpanDelimiterParser)(nil)
var _ parser.InlineParser = (*textColorSpanParser)(nil)

// The fast text color node kind.
var KindTextColorSpan = ast.NewNodeKind("TextColorSpan")

// Represents a text span to color ('~=(text)=~').
//
// This is the syntax used by the Fast Text Color plugin in Obsidian:
// https://github.com/Superschnizel/obsidian-fast-text-color
type TextColorSpan struct {
	ast.BaseInline
}

// Creates a new text color span.
func NewTextColorSpan() *TextColorSpan {
	return &TextColorSpan{}
}

func (n *TextColorSpan) Kind() ast.NodeKind {
	return KindTextColorSpan
}

func (n *TextColorSpan) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

type textColorSpanDelimiterParser struct{}

func (p *textColorSpanDelimiterParser) IsDelimiter(b byte) bool {
	return b == '~' || b == '='
}

func (p *textColorSpanDelimiterParser) CanOpenCloser(opener, closer *parser.Delimiter) bool {
	return opener.Char == closer.Char
}

func (p *textColorSpanDelimiterParser) OnMatch(consumes int) ast.Node {
	return NewTextColorSpan()
}

type textColorSpanParser struct{}

// Creates a new text color span parser.
func NewTextColorSpanParser() parser.InlineParser {
	return &textColorSpanParser{}
}

func (t *textColorSpanParser) Trigger() []byte {
	return []byte{'~', '='}
}

func (t *textColorSpanParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, segment := block.PeekLine()

	if len(line) < 2 {
		return nil
	}

	// Check if this is an opening or closing delimiter.
	var canOpen, canClose bool
	if line[0] == '~' && line[1] == '=' {
		canOpen = true
	} else if line[0] == '=' && line[1] == '~' {
		canClose = true
	} else {
		return nil
	}

	node := parser.NewDelimiter(canOpen, canClose, 2, '~', &textColorSpanDelimiterParser{})
	node.Segment = segment.WithStop(segment.Start + 2)

	block.Advance(2)
	pc.PushDelimiter(node)
	return node
}
