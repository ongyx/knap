package obsidian

import (
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// Matches a callout.
var reCallout = regexp.MustCompile(`^\[\!(\w+)\]\s*`)

var _ parser.InlineParser = (*calloutParser)(nil)

type calloutParser struct{}

func NewCalloutParser() parser.InlineParser {
	return &calloutParser{}
}

func (c *calloutParser) Trigger() []byte {
	return []byte{'['}
}

func (c *calloutParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	// The callout must appear at the start of the first paragraph in a blockquote, i.e.:
	// - Blockquote (parent.Parent)
	//   - Paragraph (parent)
	//     - current line
	// This is similar to the check done for a tasklist:
	// https://github.com/yuin/goldmark/blob/379bf24a47e6ef07f34d7536aead86d8792ac300/extension/tasklist.go#L34
	p, ok := parent.(*ast.Paragraph)
	if !ok || p.HasChildren() {
		return nil
	}

	bq, ok := parent.Parent().(*ast.Blockquote)
	if !ok || bq.FirstChild() != p {
		return nil
	}

	line, _ := block.PeekLine()
	m := reCallout.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}

	name := line[m[2]:m[3]]
	block.Advance(m[1])

	return NewCallout(name)
}
