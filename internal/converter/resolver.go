package converter

import (
	"bytes"

	"github.com/ongyx/knap/internal/obsidian"
	"github.com/ongyx/knap/internal/schema"
	"github.com/yuin/goldmark/ast"
)

// The default resolver.
//
// ResolveInternalLink returns a goldmark/ast.Text node with a plaintext representation, and ResolveColor returns an empty string.
var DefaultResolver Resolver = &defaultResolver{}

// Resolver resolves certain elements in the Markdown AST to output an Outline document correctly.
type Resolver interface {
	// Resolves an internal link to a schema node.
	ResolveInternalLink(link *Link) (*schema.Node, error)

	// Resolves a text color to a hex color by name.
	// If an invalid or empty string is returned, the resulting highlight will default to yellow in Outline.
	ResolveColor(doc *ast.Document, tc *obsidian.TextColor) string
}

type defaultResolver struct {
}

func (r *defaultResolver) ResolveInternalLink(link *Link) (*schema.Node, error) {
	var buf bytes.Buffer

	if link.Embed {
		buf.WriteByte('!')
	}

	buf.WriteString("[[")
	buf.WriteString(link.URL.String())

	if link.Text != "" {
		buf.WriteByte('|')
		buf.WriteString(link.Text)
	}

	buf.WriteString("]]")

	return schema.NewTextNode(buf.String()), nil
}

func (r *defaultResolver) ResolveColor(_ *ast.Document, _ *obsidian.TextColor) string {
	return ""
}
