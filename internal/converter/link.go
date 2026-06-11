package converter

import (
	"errors"
	"net/url"

	"github.com/yuin/goldmark/ast"
	"go.abhg.dev/goldmark/wikilink"
)

// Error returned by NewLinkFromNode if an AST node can't be parsed into a link.
var ErrUnknownNodeToLink = errors.New("AST node can't be parsed into a link")

// Link represents a link within a note.
// Internal links refer to a local note or attachment within a vault; external links refer to an external resource with a URL.
//
// Links can be written as:
//   - Links `[text](path)`
//   - Embed links `![alt text|widthxheight](path)`
//   - Wikilinks `[[path|text]]`
//   - Embed wikilinks `![[path|widthxheight]]`
//
// Wikilinks are always considered as internal links.
//
// See https://obsidian.md/help/links and https://obsidian.md/help/embeds for more details.
type Link struct {
	// The URL target.
	URL *url.URL
	// The link text. This does not contain any Markdown formatting.
	Text string
	// Whether or not to embed the note or attachment instead of linking to it.
	Embed bool
}

// Parses a link from a Markdown AST node and its source text.
//
// Currently, only these types can be parsed into a link:
//   - goldmark/*ast.Link
//   - goldmark/*ast.Image
//   - *wikilink.Node
//
// ErrUnknownNodeToLink is returned for any other type.
func ParseLinkFromNode(node ast.Node, source []byte) (*Link, error) {
	var (
		u     *url.URL
		err   error
		embed bool
	)

	switch n := node.(type) {
	case *ast.Link:
		u, err = url.Parse(string(n.Destination))
		if err != nil {
			return nil, err
		}
	case *ast.Image:
		u, err = url.Parse(string(n.Destination))
		if err != nil {
			return nil, err
		}
		embed = true
	case *wikilink.Node:
		u = &url.URL{
			Path:     string(n.Target),
			Fragment: string(n.Fragment),
		}
		embed = n.Embed
	default:
		return nil, ErrUnknownNodeToLink
	}

	// Obsidian will not render Markdown formatting in link text, so we grab the values of the text nodes as a single byte slice here.
	// This strips the formatting, e.g., [**foo** __bar__ baz] becomes 'foo bar baz'.
	t := NodeChildrenToText(node, source)

	return &Link{
		URL:   u,
		Text:  t,
		Embed: embed,
	}, nil
}

// Returns true if the link is internal and refers to a local note/attachment in a vault.
func (l *Link) IsInternal() bool {
	return !l.URL.IsAbs()
}
