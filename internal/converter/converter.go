package converter

import (
	"bytes"
	"errors"
	"fmt"
	"slices"

	"github.com/ongyx/knap/internal/collections"
	"github.com/ongyx/knap/internal/obsidian"
	"github.com/ongyx/knap/internal/prosemirror"
	"github.com/ongyx/knap/internal/util"
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

	// This is supposed to be based on the table size in the editor, I guess?
	defaultTableRowWidth = 735.2
)

// Error returned by Converter.Convert when raw HTML fragments are not recognized. This error will be wrapped.
var ErrInvalidHTML = errors.New("raw HTML is not recognized, only <br> is supported")

// Represents a context context for walking the AST.
type context struct {
	// The Prosemirror node, if any.
	mnode *prosemirror.Node
	// The formatting to apply to descendant text nodes.
	marks []prosemirror.Mark
}

// Converter parses Markdown text to convert it to a Prosemirror document.
type Converter struct {
	contexts collections.Stack[context]
	markdown goldmark.Markdown
	resolver Resolver

	source         []byte
	markdownDoc    *ast.Document
	prosemirrorDoc *prosemirror.Node
}

// Creates a new converter with the given resolver.
func New(res Resolver) *Converter {
	if res == nil {
		res = DefaultResolver
	}

	md := goldmark.New(obsidian.DefaultOptions())

	return &Converter{
		contexts: collections.NewStack[context](0, 25),
		markdown: md,
		resolver: res,
	}
}

// Parses the Markdown text in src and converts its AST into a document Prosemirror node.
func (cv *Converter) Convert(src []byte) (*prosemirror.Node, error) {
	cv.source = src
	cv.prosemirrorDoc = nil

	p := cv.markdown.Parser()
	r := text.NewReader(src)
	// SAFETY: The root node is always a Document.
	cv.markdownDoc = p.Parse(r).(*ast.Document)

	if err := ast.Walk(cv.markdownDoc, cv.walk); err != nil {
		return nil, err
	}

	if len(cv.prosemirrorDoc.Content) == 0 {
		// Empty Markdown files will result in an invalid document without content,
		// so a blank document with an empty paragraph must be created instead.
		cv.prosemirrorDoc = prosemirror.NewBlankDocumentNode()
	}

	return cv.prosemirrorDoc, nil
}

func (cv *Converter) walk(anode ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		// The AST node has been fully walked, pop its context.
		cv.contexts.Pop()
		return ast.WalkContinue, nil
	}

	mnode, marks, walkStatus, err := cv.astToProsemirror(anode)
	if err != nil {
		return walkStatus, err
	}

	ctx := context{mnode, nil}

	if pctx, ok := cv.contexts.Peek(); ok {
		if mnode != nil {
			var p *prosemirror.Node
			// The table extension does not generate a Paragraph/TextBlock AST node for nodes underneath a TableCell AST node,
			// so we must wrap it in a Prosemirror paragraph to avoid the table getting erased in Outline.
			if pctx.mnode.Type == prosemirror.NodeTableHeader || pctx.mnode.Type == prosemirror.NodeTableCell {
				if len(pctx.mnode.Content) == 0 {
					p = prosemirror.NewParagraphNode()
					pctx.mnode.Content = []*prosemirror.Node{p}
				} else {
					p = pctx.mnode.Content[0]
				}
			} else {
				p = pctx.mnode
			}

			// Append the node to its parent's content.
			p.Content = append(p.Content, mnode)
		} else {
			// No node was converted for this walk, use the parent node for the new context.
			ctx.mnode = pctx.mnode
		}

		// Preserve the marks from the parent context and add the newly generated ones, if any.
		// This will be copied into descendant text nodes.
		ctx.marks = slices.Concat(pctx.marks, marks)
	} else {
		// The first Prosemirror node is always the document root.
		cv.prosemirrorDoc = mnode
	}

	if ctx.mnode == nil {
		panic("walk: parent node is missing, no more context left")
	}

	// Push the new context for this AST node.
	cv.contexts.Push(ctx)
	return walkStatus, nil
}

// Converts an AST node to a Prosemirror node.
// If mnode is not nil, it will be added to the parent node's content.
// If marks is not nil, it will be appended to the context's marks for descendant nodes.
func (cv *Converter) astToProsemirror(anode ast.Node) (mnode *prosemirror.Node, marks []prosemirror.Mark, walkStatus ast.WalkStatus, err error) {
	pctx, _ := cv.contexts.Peek()

	switch an := anode.(type) {
	case *ast.Document:
		mnode := prosemirror.NewDocumentNode()
		return mnode, nil, ast.WalkContinue, nil

	// Inline elements

	case *ast.String:
		// Strings must be emitted without any marks.
		mnode := prosemirror.NewTextNode(string(an.Value))
		return mnode, nil, ast.WalkContinue, nil

	case *ast.Text:
		v := string(an.Value(cv.source))
		if v == "" && an.SoftLineBreak() {
			// Don't generate a new text node. Empty text nodes may generate between AST wikilinks, but Prosemirror will not recognize the resulting text node as valid.
			return nil, nil, ast.WalkContinue, nil
		}
		mnode := prosemirror.NewTextNode(v)
		mnode.Marks = pctx.marks
		return mnode, nil, ast.WalkContinue, nil

	case *ast.RawHTML:
		v := string(an.Segments.Value(cv.source))
		if v == "<br>" {
			mnode := prosemirror.NewLineBreakNode()
			return mnode, nil, ast.WalkContinue, nil
		} else {
			return nil, nil, ast.WalkStop, fmt.Errorf("%w (text: %q)", ErrInvalidHTML, v)
		}

	case *ast.ThematicBreak:
		// The actual markup is not stored in the node, so we need to index into the source.
		pos := an.Pos()
		markup := cv.source[pos : pos+3]
		isPageBreak := bytes.Equal(markup, []byte("***"))

		mnode := prosemirror.NewThematicBreakNode(isPageBreak)
		return mnode, nil, ast.WalkContinue, nil

	case *ast.Heading:
		mnode := prosemirror.NewHeadingNode(an.Level)
		return mnode, nil, ast.WalkContinue, nil

	// Block elements

	case *ast.Paragraph, *ast.TextBlock:
		mnode := prosemirror.NewParagraphNode()
		return mnode, nil, ast.WalkContinue, nil

	case *ast.Blockquote:
		var mnode *prosemirror.Node
		if nt, ok := cv.extractNotice(an); ok {
			mnode = prosemirror.NewNoticeNode(nt)
		} else {
			mnode = prosemirror.NewBlockQuoteNode()
		}
		return mnode, nil, ast.WalkContinue, nil

	case *ast.CodeBlock:
		mnode := prosemirror.NewFencedCodeBlockNode("none")
		cv.addLinesAsTextContent(mnode, anode)

		return mnode, nil, ast.WalkContinue, nil

	case *ast.FencedCodeBlock:
		lang := string(an.Info.Value(cv.source))
		mnode := prosemirror.NewFencedCodeBlockNode(lang)
		// goldmark does not parse the text inside the code block, so we have to add it to the node here.
		cv.addLinesAsTextContent(mnode, anode)

		return mnode, nil, ast.WalkContinue, nil

	case *ast.List:
		var mnode *prosemirror.Node
		if an.IsOrdered() {
			mnode = prosemirror.NewOrderedListNode(an.Start)
		} else if isChecklist(an) {
			mnode = prosemirror.NewChecklistNode()
		} else {
			mnode = prosemirror.NewBulletListNode()
		}
		return mnode, nil, ast.WalkContinue, nil

	case *ast.ListItem:
		var mnode *prosemirror.Node
		if isChecked, ok := extractCheckbox(an); ok {
			mnode = prosemirror.NewChecklistItemNode(isChecked)
		} else {
			mnode = prosemirror.NewListItemNode()
		}

		return mnode, nil, ast.WalkContinue, nil

	case *east.Table:
		mnode := prosemirror.NewTableNode()
		return mnode, nil, ast.WalkContinue, nil

	case *east.TableHeader, *east.TableRow:
		mnode := prosemirror.NewTableRowNode()
		return mnode, nil, ast.WalkContinue, nil

	case *east.TableCell:
		p := an.Parent()
		_, ih := p.(*east.TableHeader)

		// Every table cell row should have a column width except for the last one.
		var cw float64
		if an != p.LastChild() {
			cw = defaultTableRowWidth / float64(p.ChildCount())
		}

		mnode := prosemirror.NewTableCellNode(ih, cw)
		return mnode, nil, ast.WalkContinue, nil

	case *ast.HTMLBlock:
		pctx, _ := cv.contexts.Peek()
		p := pctx.mnode

		// This is kinda annoying to handle, ngl.
		for i := range an.Lines().Len() {
			l := an.Lines().At(i)
			// Trim the trailing newline.
			b := bytes.TrimSpace(l.Value(cv.source))
			if string(b) == "<br>" {
				p.Content = append(p.Content, prosemirror.NewLineBreakNode())
			} else {
				return nil, nil, ast.WalkStop, ErrInvalidHTML
			}
		}

	case *ast.Link, *ast.Image, *wikilink.Node:
		var (
			link  *Link
			mnode *prosemirror.Node
			err   error
		)

		link, err = ParseLinkFromNode(an, cv.source)
		if err != nil {
			return nil, nil, ast.WalkStop, err
		}

		// Internal links must be resolved externally as they may reference other notes or attachments in a vault.
		if link.IsInternal() {
			mnode, err = cv.resolver.ResolveInternalLink(link, pctx.marks)
		} else {
			mnode, err = resolveExternalLink(link, pctx)
		}

		// Do not handle the child nodes within the link-like nodes.
		return mnode, nil, ast.WalkSkipChildren, err

	// These elements below are special, because Outline represents them as 'marks' that are then applied to descendant text nodes instead of creating new nodes.

	case *ast.Emphasis:
		var m prosemirror.Mark
		// Italic/bold is represented as an integer level since they can be written with '*' or '_', i.e., '*text*' is italic, and '**text**' is bold.
		// https://github.com/yuin/goldmark/blob/379bf24a47e6ef07f34d7536aead86d8792ac300/renderer/html/html.go#L564
		switch an.Level {
		case EmphasisLevelItalic:
			m = prosemirror.NewItalicMark()
		case EmphasisLevelBold:
			m = prosemirror.NewBoldMark()
		default:
			panic("unreachable: Emphasis level is always 1 or 2")
		}

		return nil, []prosemirror.Mark{m}, ast.WalkContinue, nil

	case *east.Strikethrough:
		return nil, []prosemirror.Mark{prosemirror.NewStrikethroughMark()}, ast.WalkContinue, nil

	case *ast.CodeSpan:
		return nil, []prosemirror.Mark{prosemirror.NewInlineCodeMark()}, ast.WalkContinue, nil

	case *obsidian.Highlight:
		return nil, []prosemirror.Mark{prosemirror.NewHighlightMark("")}, ast.WalkContinue, nil

	case *obsidian.TextColor:
		// While it is possible for the resolver to call OwnerDocument() instead,
		// this incurs a traversal every time a color lookup is needed.
		clr := cv.resolver.ResolveColor(cv.markdownDoc, an)
		return nil, []prosemirror.Mark{prosemirror.NewHighlightMark(clr)}, ast.WalkContinue, nil
	}

	// It's completely valid to skip certain AST nodes because they were processed by their parent.
	return nil, nil, ast.WalkContinue, nil
}

func (cv *Converter) addLinesAsTextContent(mnode *prosemirror.Node, anode ast.Node) {
	for i := range anode.Lines().Len() {
		ls := anode.Lines().At(i)
		l := string(ls.Value(cv.source))
		mnode.Content = append(mnode.Content, prosemirror.NewTextNode(l))
	}
}

func (cv *Converter) extractNotice(bq *ast.Blockquote) (nt prosemirror.NoticeType, ok bool) {
	// Check if the blockquote's first paragraph contains a callout.
	p, ok := bq.FirstChild().(*ast.Paragraph)
	if !ok {
		return
	}

	cl, ok := p.FirstChild().(*obsidian.Callout)
	if !ok {
		return
	}

	return prosemirror.CalloutToNotice(cl), true
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

func resolveExternalLink(link *Link, pctx context) (*prosemirror.Node, error) {
	if link.Embed {
		ff := util.ParseFileFormat(link.URL.Path)

		// If the URL points to an image file, generate an image URL node.
		if ff == util.FileImage {
			w, h, _ := ParseEmbedSize(link.Text)
			return prosemirror.NewImageURLNode(link.URL.String(), w, h), nil
		}

		// Generate a regular embed node.
		return prosemirror.NewEmbedNode(link.URL.String()), nil
	}

	// Generate a text node with a link mark attached to it.
	tn := prosemirror.NewTextNode(string(link.Text))
	tn.Marks = append(slices.Clip(pctx.marks), prosemirror.NewLinkMark(link.URL.String()))
	return tn, nil
}
