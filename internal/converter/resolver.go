package converter

import (
	"go.abhg.dev/goldmark/wikilink"
)

// The default resolver.
var DefaultResolver Resolver = &defaultResolver{wikilink.DefaultResolver}

// Resolver resolves certain elements in the Markdown AST to output an Outline document correctly.
type Resolver interface {
	wikilink.Resolver

	// Resolves a color name to a hex color.
	// If an invalid string is returned, the resulting highlight will default to yellow in Outline.
	ResolveColor(name []byte) string
}

type defaultResolver struct {
	wikilink.Resolver
}

func (r *defaultResolver) ResolveColor(name []byte) string {
	return ""
}
