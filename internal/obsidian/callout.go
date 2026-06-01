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

// Matches a callout. First group is the callout name.
var reCallout = regexp.MustCompile(`^\[\!(\w+)\]\s*`)

// The callout node kind.
var KindCallout = ast.NewNodeKind("Callout")

// Represents a callout within a blockquote ('> [!info]').
type Callout struct {
	ast.BaseInline
	// The name/type of the callout, e.g., 'info'.
	Name []byte
}

// Creates a new callout with the given name.
func NewCallout(name []byte) *Callout {
	return &Callout{Name: name}
}

func (n *Callout) Kind() ast.NodeKind {
	return KindCallout
}

func (n *Callout) Dump(source []byte, level int) {
	m := map[string]string{
		"Name": string(n.Name),
	}
	ast.DumpHelper(n, source, level, m, nil)
}

type calloutParser struct{}

// Creates a new callout parser.
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
