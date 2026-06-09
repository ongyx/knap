package converter

import (
	"github.com/yuin/goldmark/ast"
	"go.abhg.dev/goldmark/wikilink"
)

// InternalLink represents an Obsidian internal link, which refer to a local note or attachment within a vault.
//
// Internal links can be written as:
//   - Wikilinks `[[path|title]]`
//   - Embed wikilinks `![[path]]`
//   - Image embed wikilinks `![[path|widthxheight]]`
//
// See https://obsidian.md/help/links and https://obsidian.md/help/embeds for more details.
type InternalLink struct {
	// The path to a note or attachment. If nil, the link refers to the same note it is written in.
	Target []byte
	// The fragment, for linking to headings in notes. If nil, the link refers to the note itself.
	Fragment []byte
	// The display text to show. If nil, the note or attachment name is shown instead.
	Title []byte
	// Whether or not to embed the note or attachment instead of linking to it.
	Embed bool
}

// Parses an internal link from a wikilink node and the markdown source text.
func NewInternalLink(node *wikilink.Node, source []byte) InternalLink {
	var title []byte
	// The title of the internal link is stored as a text child, if any.
	if t, ok := node.FirstChild().(*ast.Text); ok {
		title = t.Value(source)
	}

	return InternalLink{node.Target, node.Fragment, title, node.Embed}
}
