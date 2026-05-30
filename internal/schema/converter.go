package schema

import (
	"bytes"
	"errors"
	"slices"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/wikilink"
	"golang.org/x/net/html"
)

const (
	// Represents an italic emphasis.
	EmphasisLevelItalic = iota + 1
	// Represents a bold emphasis.
	EmphasisLevelBold
)

// Error returned by Converter.Convert when raw HTML fragments are not recognized.
var ErrRawHTML = errors.New("raw HTML is not recognized (only <br> is supported)")

// Default options for parsing Markdown.
var defaultOptions = []goldmark.Option{
	goldmark.WithExtensions(
		extension.Strikethrough,
		extension.Table,
		extension.TaskList,
		&wikilink.Extender{},
	),
}

// Converter parses Markdown text to convert it to Prosemirror nodes.
type Converter struct {
	contexts Stack[context]
	source   []byte
	root     *Node
}

// Represents a context context for walking the AST.
type context struct {
	// The Prosemirror node, if any.
	snode *Node
	// The formatting to apply to descendant text nodes.
	marks []Mark
}

// Creates a new Converter.
func NewConverter() *Converter {
	return &Converter{
		contexts: NewStack[context](0, 25),
	}
}

// Parses the Markdown text in src and converts its AST into a Prosemirror node.
func (cv *Converter) Convert(src []byte) (*Node, error) {
	cv.source = src
	cv.root = nil

	m := goldmark.New(defaultOptions...)
	an := m.Parser().Parse(text.NewReader(src))

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
			panic("walk: parent node is missing")
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
func (cv *Converter) astToSchema(anode ast.Node) (snode *Node, marks []Mark, walkStatus ast.WalkStatus, err error) {
	parent, _ := cv.contexts.Peek()

	switch an := anode.(type) {
	case *ast.Document:
		snode := NewDocumentNode()
		return snode, nil, ast.WalkContinue, nil

	// Inline elements

	case *ast.Text:
		v := string(an.Value(cv.source))
		snode := NewTextNode(v)
		snode.Marks = parent.marks
		return snode, nil, ast.WalkContinue, nil

	case *ast.RawHTML:
		v := string(an.Segments.Value(cv.source))
		if v == "<br>" {
			snode := NewLineBreakNode()
			return snode, nil, ast.WalkContinue, nil
		} else {
			return nil, nil, ast.WalkStop, ErrRawHTML
		}

	case *ast.ThematicBreak:
		// The actual markup is not stored in the node, so we need to index into the source.
		pos := an.Pos()
		markup := cv.source[pos : pos+3]
		isPageBreak := bytes.Equal(markup, []byte("***"))

		snode := NewThematicBreakNode(isPageBreak)
		return snode, nil, ast.WalkContinue, nil

	case *ast.Heading:
		snode := NewHeadingNode(an.Level)
		return snode, nil, ast.WalkContinue, nil

	// Block elements

	case *ast.Paragraph:
		snode := NewParagraphNode()
		return snode, nil, ast.WalkContinue, nil

	case *ast.Blockquote:
		var snode *Node
		if nt, ok := cv.extractNotice(an); ok {
			snode = NewNoticeNode(nt)
		} else {
			snode = NewBlockQuoteNode()
		}
		return snode, nil, ast.WalkContinue, nil

	case *ast.CodeBlock:
		snode := NewFencedCodeBlockNode("none")
		cv.addLinesContent(snode, anode)

		return snode, nil, ast.WalkContinue, nil

	case *ast.FencedCodeBlock:
		lang := string(an.Info.Value(cv.source))
		snode := NewFencedCodeBlockNode(lang)
		// goldmark does not parse the text inside the code block, so we have to add it to the node here.
		cv.addLinesContent(snode, anode)

		return snode, nil, ast.WalkContinue, nil

	case *ast.List:
		var snode *Node
		if an.IsOrdered() {
			snode = NewOrderedListNode(an.Start)
		} else if isChecklist(an) {
			snode = NewChecklistNode()
		} else {
			snode = NewBulletListNode()
		}
		return snode, nil, ast.WalkContinue, nil

	case *ast.ListItem:
		var snode *Node
		if isChecked, ok := extractCheckbox(an); ok {
			snode = NewChecklistItemNode(isChecked)
		} else {
			snode = NewListItemNode()
		}

		return snode, nil, ast.WalkContinue, nil

	case *east.Table:
		snode := NewTableNode()
		return snode, nil, ast.WalkContinue, nil

	case *east.TableHeader, *east.TableRow:
		snode := NewTableRowNode()
		return snode, nil, ast.WalkContinue, nil

	case *east.TableCell:
		_, isHeader := an.Parent().(*east.TableHeader)
		snode := NewTableCellNode(isHeader)
		return snode, nil, ast.WalkContinue, nil

	// These elements below are special because Outline represents them as 'marks' that are then applied to descendant text nodes instead of creating new nodes.

	case *ast.Emphasis:
		var m Mark
		// Italic/bold is represented as an integer level since they can be written with '*' or '_', i.e., '*text*' is italic, and '**text**' is bold.
		// https://github.com/yuin/goldmark/blob/379bf24a47e6ef07f34d7536aead86d8792ac300/renderer/html/html.go#L564
		switch an.Level {
		case EmphasisLevelItalic:
			m = NewItalicMark()
		case EmphasisLevelBold:
			m = NewBoldMark()
		default:
			panic("unreachable: Emphasis level is always 1 or 2")
		}

		return nil, []Mark{m}, ast.WalkContinue, nil

	case *east.Strikethrough:
		return nil, []Mark{NewStrikethroughMark()}, ast.WalkContinue, nil

	case *ast.Link:
		m := NewLinkMark(string(an.Destination))
		return nil, []Mark{m}, ast.WalkContinue, nil

	case *ast.CodeSpan:
		return nil, []Mark{NewInlineCodeMark()}, ast.WalkContinue, nil
	}

	// It's completely valid to skip certain AST nodes because they were processed by their parent.
	return nil, nil, ast.WalkContinue, nil
}

func (cv *Converter) addLinesContent(snode *Node, anode ast.Node) {
	for i := range anode.Lines().Len() {
		ls := anode.Lines().At(i)
		l := string(ls.Value(cv.source))
		snode.Content = append(snode.Content, NewTextNode(l))
	}
}

func (cv *Converter) extractNotice(bq *ast.Blockquote) (nt NoticeType, ok bool) {
	p, ok := bq.FirstChild().(*ast.Paragraph)
	if !ok {
		return
	}

	// Try to find the callout in the first paragraph.
	l := p.Lines().At(0)
	cb := reCallout.Find(l.Value(cv.source))
	if cb == nil {
		return
	}

	// If the callout name isn't defined, the notice type defaults to NoticeInfo.
	return calloutToNotice[string(cb)], true
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

func getAttr(node *html.Node, key string) (attr *html.Attribute, ok bool) {
	for i, a := range node.Attr {
		if a.Key == key {
			return &node.Attr[i], true
		}
	}

	return nil, false
}

func getAttrs(node *html.Node) map[string]string {
	m := make(map[string]string, len(node.Attr))
	for _, a := range node.Attr {
		m[a.Key] = a.Val
	}

	return m
}
