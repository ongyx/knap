package converter

import (
	"bytes"
	"errors"
	"slices"

	"github.com/ongyx/knap/internal/obsidian"
	"github.com/ongyx/knap/internal/schema"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/wikilink"
)

const (
	// Represents an italic emphasis.
	EmphasisLevelItalic = iota + 1
	// Represents a bold emphasis.
	EmphasisLevelBold
)

// Error returned by Converter.Convert when raw HTML fragments are not recognized.
var ErrInvalidHTML = errors.New("raw HTML is not recognized (only <br> is supported)")

// Represents a context context for walking the AST.
type context struct {
	// The Prosemirror node, if any.
	snode *schema.Node
	// The formatting to apply to descendant text nodes.
	marks []schema.Mark
}

// Converter parses Markdown text to convert it to Prosemirror nodes.
type Converter struct {
	contexts Stack[context]
	markdown goldmark.Markdown
	source   []byte
	root     *schema.Node
}

// Creates a new converter with the given wikilink resolver.
func New(resolver wikilink.Resolver) *Converter {
	md := goldmark.New(obsidian.DefaultOptions)

	return &Converter{
		contexts: NewStack[context](0, 25),
		markdown: md,
	}
}

// Parses the Markdown text in src and converts its AST into a Prosemirror node.
func (cv *Converter) Convert(src []byte) (*schema.Node, error) {
	cv.source = src
	cv.root = nil

	an := cv.markdown.Parser().Parse(text.NewReader(src))
	an.Dump(src, 4)

	if err := ast.Walk(an, cv.walk); err != nil {
		return nil, err
	}

	return cv.root, nil
}

func (cv *Converter) walk(anode ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		snode, marks, walkStatus, err := cv.astToSchema(anode)
		if err != nil {
			return walkStatus, err
		}

		ctx := context{snode, nil}

		if parent, ok := cv.contexts.Peek(); ok {
			if snode != nil {
				// Append the node to its parent's content.
				parent.snode.Content = append(parent.snode.Content, snode)
			} else {
				// No node was converted for this walk, preserve the parent node in the new context.
				ctx.snode = parent.snode
			}

			ctx.marks = slices.Concat(parent.marks, marks)
		} else {
			// The first schema node is always the document root.
			cv.root = snode
		}

		if ctx.snode == nil {
			panic("walk: parent node is missing, no more context left")
		}

		// Push the schema node onto the context stack. This will be popped when the AST node has been walked.
		cv.contexts.Push(ctx)
		return walkStatus, nil
	} else {
		// Pop the context for this AST node.
		cv.contexts.Pop()
		return ast.WalkContinue, nil
	}
}

// Converts an AST node to a schema node.
// If snode is not nil, it will be added to the parent node's content.
// If marks is not nil, it will be appended to the context's marks for descendant nodes.
func (cv *Converter) astToSchema(anode ast.Node) (snode *schema.Node, marks []schema.Mark, walkStatus ast.WalkStatus, err error) {
	parent, _ := cv.contexts.Peek()

	switch an := anode.(type) {
	case *ast.Document:
		snode := schema.NewDocumentNode()
		return snode, nil, ast.WalkContinue, nil

	// Inline elements

	case *ast.Text:
		v := string(an.Value(cv.source))
		snode := schema.NewTextNode(v)
		snode.Marks = parent.marks
		return snode, nil, ast.WalkContinue, nil

	case *ast.RawHTML:
		v := string(an.Segments.Value(cv.source))
		if v == "<br>" {
			snode := schema.NewLineBreakNode()
			return snode, nil, ast.WalkContinue, nil
		} else {
			return nil, nil, ast.WalkStop, ErrInvalidHTML
		}

	case *ast.ThematicBreak:
		// The actual markup is not stored in the node, so we need to index into the source.
		pos := an.Pos()
		markup := cv.source[pos : pos+3]
		isPageBreak := bytes.Equal(markup, []byte("***"))

		snode := schema.NewThematicBreakNode(isPageBreak)
		return snode, nil, ast.WalkContinue, nil

	case *ast.Heading:
		snode := schema.NewHeadingNode(an.Level)
		return snode, nil, ast.WalkContinue, nil

	// Block elements

	case *ast.Paragraph, *ast.TextBlock:
		snode := schema.NewParagraphNode()
		return snode, nil, ast.WalkContinue, nil

	case *ast.Blockquote:
		var snode *schema.Node
		if nt, ok := cv.extractNotice(an); ok {
			snode = schema.NewNoticeNode(nt)
		} else {
			snode = schema.NewBlockQuoteNode()
		}
		return snode, nil, ast.WalkContinue, nil

	case *ast.CodeBlock:
		snode := schema.NewFencedCodeBlockNode("none")
		cv.addLinesContent(snode, anode)

		return snode, nil, ast.WalkContinue, nil

	case *ast.FencedCodeBlock:
		lang := string(an.Info.Value(cv.source))
		snode := schema.NewFencedCodeBlockNode(lang)
		// goldmark does not parse the text inside the code block, so we have to add it to the node here.
		cv.addLinesContent(snode, anode)

		return snode, nil, ast.WalkContinue, nil

	case *ast.List:
		var snode *schema.Node
		if an.IsOrdered() {
			snode = schema.NewOrderedListNode(an.Start)
		} else if isChecklist(an) {
			snode = schema.NewChecklistNode()
		} else {
			snode = schema.NewBulletListNode()
		}
		return snode, nil, ast.WalkContinue, nil

	case *ast.ListItem:
		var snode *schema.Node
		if isChecked, ok := extractCheckbox(an); ok {
			snode = schema.NewChecklistItemNode(isChecked)
		} else {
			snode = schema.NewListItemNode()
		}

		return snode, nil, ast.WalkContinue, nil

	case *east.Table:
		snode := schema.NewTableNode()
		return snode, nil, ast.WalkContinue, nil

	case *east.TableHeader, *east.TableRow:
		snode := schema.NewTableRowNode()
		return snode, nil, ast.WalkContinue, nil

	case *east.TableCell:
		_, isHeader := an.Parent().(*east.TableHeader)
		snode := schema.NewTableCellNode(isHeader)
		return snode, nil, ast.WalkContinue, nil

	case *ast.HTMLBlock:
		parent, _ := cv.contexts.Peek()
		psnode := parent.snode

		// This is kinda annoying to handle, ngl.
		for i := range an.Lines().Len() {
			l := an.Lines().At(i)
			// Trim the trailing newline.
			b := bytes.TrimSpace(l.Value(cv.source))
			if string(b) == "<br>" {
				psnode.Content = append(psnode.Content, schema.NewLineBreakNode())
			} else {
				return nil, nil, ast.WalkStop, ErrInvalidHTML
			}
		}

	// These elements below are special, because Outline represents them as 'marks' that are then applied to descendant text nodes instead of creating new nodes.

	case *ast.Emphasis:
		var m schema.Mark
		// Italic/bold is represented as an integer level since they can be written with '*' or '_', i.e., '*text*' is italic, and '**text**' is bold.
		// https://github.com/yuin/goldmark/blob/379bf24a47e6ef07f34d7536aead86d8792ac300/renderer/html/html.go#L564
		switch an.Level {
		case EmphasisLevelItalic:
			m = schema.NewItalicMark()
		case EmphasisLevelBold:
			m = schema.NewBoldMark()
		default:
			panic("unreachable: Emphasis level is always 1 or 2")
		}

		return nil, []schema.Mark{m}, ast.WalkContinue, nil

	case *east.Strikethrough:
		return nil, []schema.Mark{schema.NewStrikethroughMark()}, ast.WalkContinue, nil

	case *ast.Link:
		m := schema.NewLinkMark(string(an.Destination))
		return nil, []schema.Mark{m}, ast.WalkContinue, nil

	case *ast.CodeSpan:
		return nil, []schema.Mark{schema.NewInlineCodeMark()}, ast.WalkContinue, nil
	}

	// It's completely valid to skip certain AST nodes because they were processed by their parent.
	return nil, nil, ast.WalkContinue, nil
}

func (cv *Converter) addLinesContent(snode *schema.Node, anode ast.Node) {
	for i := range anode.Lines().Len() {
		ls := anode.Lines().At(i)
		l := string(ls.Value(cv.source))
		snode.Content = append(snode.Content, schema.NewTextNode(l))
	}
}

func (cv *Converter) extractNotice(bq *ast.Blockquote) (nt schema.NoticeType, ok bool) {
	// Check if the blockquote's first paragraph contains a callout.
	p, ok := bq.FirstChild().(*ast.Paragraph)
	if !ok {
		return
	}

	cl, ok := p.FirstChild().(*obsidian.Callout)
	if !ok {
		return
	}

	return schema.CalloutToNotice(cl), true
}

func isChecklist(list *ast.List) bool {
	li, ok := list.FirstChild().(*ast.ListItem)
	if !ok {
		return false
	}

	_, ok = extractCheckbox(li)
	return ok
}

func extractCheckbox(li *ast.ListItem) (isChecked bool, ok bool) {
	tb, ok := li.FirstChild().(*ast.TextBlock)
	if !ok {
		return
	}

	cb, ok := tb.FirstChild().(*east.TaskCheckBox)
	if !ok {
		return
	}

	return cb.IsChecked, true
}
