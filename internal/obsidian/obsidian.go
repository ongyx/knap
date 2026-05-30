package obsidian

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/wikilink"
)

// interface asserts
var _ goldmark.Extender = (*Extender)(nil)

// Default options for parsing Obsidian Markdown.
var DefaultOptions = goldmark.WithExtensions(
	extension.Strikethrough,
	extension.Table,
	extension.TaskList,
	&wikilink.Extender{},
	&Extender{},
)

// Extends goldmark to parse Obsidian-specific syntax.
type Extender struct{}

func (o *Extender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewCalloutParser(), 50),
		),
	)
}
