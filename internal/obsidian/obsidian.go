package obsidian

import (
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/wikilink"
)

// interface asserts
var _ goldmark.Extender = (*Extender)(nil)

// Returns the default options for parsing Obsidian-style Markdown.
func DefaultOptions(res wikilink.Resolver) goldmark.Option {
	return goldmark.WithExtensions(
		extension.Strikethrough,
		extension.Table,
		extension.TaskList,
		meta.New(meta.WithTable()),
		&wikilink.Extender{Resolver: res},
		&Extender{},
	)
}

// Extends goldmark to parse Obsidian-specific syntax.
type Extender struct{}

func (o *Extender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewCalloutParser(), 50),
			util.Prioritized(NewHighlightParser(), 51),
			util.Prioritized(NewTextColorSpanParser(), 52),
			util.Prioritized(NewTextColorParser(), 53),
		),
	)
}
