package obsidian

import "github.com/yuin/goldmark/ast"

// interface asserts
var _ ast.Node = (*Callout)(nil)

// The callout node kind.
var KindCallout = ast.NewNodeKind("Callout")

// Represents a callout within a blockquote.
type Callout struct {
	ast.BaseInline
	Name []byte
}

func NewCallout(name []byte) *Callout {
	return &Callout{
		Name: name,
	}
}

func (n *Callout) Kind() ast.NodeKind {
	return KindCallout
}

func (n *Callout) Dump(source []byte, level int) {
	m := map[string]string{
		"Type": string(n.Name),
	}
	ast.DumpHelper(n, source, level, m, nil)
}
