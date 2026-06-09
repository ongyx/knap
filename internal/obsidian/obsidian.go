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

// Returns the default options for parsing Obsidian Flavored Markdown.
func DefaultOptions() goldmark.Option {
	return goldmark.WithExtensions(
		extension.Strikethrough,
		extension.Table,
		extension.TaskList,
		meta.New(meta.WithTable(), meta.WithStoresInDocument()),
		&wikilink.Extender{},
		&Extender{},
	)
}

// Extends goldmark to parse Obsidian Flavored Markdown.
//
// Currently, the following extensions are supported:
//   - Internal links: '[[Link]]'
//   - Strikethroughs: '~~Text~~'
//   - Highlights: '==Text=='
//   - Code blocks: '```'
//   - Incomplete task: '- [ ]'
//   - Completed task: '- [x]'
//   - Callouts: '> [!note]'
//   - Tables
//
// These extensions are not supported:
//   - Embed files: '![[Link]]'
//   - Block references: '![[Link#^id]]'
//   - Defining a block: '^id'
//   - Footnotes: '[^id]'
//   - Comments: '%%Text%%'
//
// See https://obsidian.md/help/obsidian-flavored-markdown for more details.
type Extender struct{}

func (o *Extender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewCalloutParser(), 50),
			util.Prioritized(NewHighlightParser(), 51),
			util.Prioritized(NewTextColorParser(), 52),
		),
	)
}
