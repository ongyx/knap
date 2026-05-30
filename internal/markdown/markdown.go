package markdown

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/wikilink"
)

// type asserts
var _ renderer.NodeRenderer = (*FencedCodeBlockRenderer)(nil)
var _ goldmark.Extender = (*prosemirror)(nil)

// The default options for converting Obsidian Markdown to HTMl for further processing.
var DefaultOptions = []goldmark.Option{
	goldmark.WithExtensions(extension.GFM, &wikilink.Extender{}, &Prosemirror),
	goldmark.WithRendererOptions(html.WithUnsafe(), html.WithHardWraps()),
}

// A markdown extension that preserves additional elements for conversion to Prosemirror format.
var Prosemirror = prosemirror{}

type prosemirror struct{}

func (p *prosemirror) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&FencedCodeBlockRenderer{}, 200),
		util.Prioritized(&ThematicBreakRenderer{}, 201),
	))
}

// Renders fenced code blocks (```...```) into plain <pre><code> elements.
// This is a stripped down version of https://github.com/yuin/goldmark-highlighting/blob/v2/highlighting.go.
type FencedCodeBlockRenderer struct{}

func (c *FencedCodeBlockRenderer) RegisterFuncs(r renderer.NodeRendererFuncRegisterer) {
	r.Register(ast.KindFencedCodeBlock, c.render)
}

func (c *FencedCodeBlockRenderer) render(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	cb := node.(*ast.FencedCodeBlock)
	if !entering {
		return ast.WalkContinue, nil
	}
	lang := cb.Language(src)

	// HTML preamble.
	w.WriteString("<pre>")
	if lang != nil {
		fmt.Fprintf(w, `<code class="language-%s">`, lang)
	} else {
		w.WriteString("<code>")
	}

	// Copy the contents of the fenced code block directly into the buffer.
	cl := cb.Lines().Len()
	for i := 0; i < cl; i++ {
		l := cb.Lines().At(i)
		w.Write(l.Value(src))
	}

	// Close the code block off.
	w.WriteString("</code></pre>\n")

	return ast.WalkContinue, nil
}

// Renders thematic breaks (`***`, `---`, `___`) into <hr> elements.
// For `***`, the class `pagebreak` is added as Outline uses it as a pagebreak marker instead.
type ThematicBreakRenderer struct{}

func (t *ThematicBreakRenderer) RegisterFuncs(r renderer.NodeRendererFuncRegisterer) {
	r.Register(ast.KindThematicBreak, t.render)
}

func (t *ThematicBreakRenderer) render(w util.BufWriter, src []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	tb := node.(*ast.ThematicBreak)
	if !entering {
		return ast.WalkContinue, nil
	}

	// The actual markup is not stored in the node, so we need to index into the source.
	pos := tb.Pos()
	markup := src[pos : pos+3]
	if bytes.Equal(markup, []byte("***")) {
		w.WriteString(`<hr class="pagebreak">`)
	} else {
		w.WriteString(`<hr>`)
	}
	w.WriteByte('\n')

	return ast.WalkContinue, nil
}
