package obsidian

import (
	"testing"

	"github.com/ongyx/knap/internal/testutil"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestHighlight(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "simple highlight",
			source:   "==Hello World!==",
			expected: "Hello World!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithParserOptions(
					parser.WithInlineParsers(
						util.Prioritized(NewHighlightParser(), 50),
					),
				),
			)
			source := []byte(tt.source)
			doc := md.Parser().Parse(text.NewReader(source))

			found := testutil.FindNode[*Highlight](doc)
			if found == nil {
				t.Fatal("Highlight node not found in AST")
			}

			txt := found.FirstChild().(*ast.Text)
			if txt == nil {
				t.Fatal("Expected text node in highlight node")
			}

			v := string(txt.Segment.Value(source))
			if v != tt.expected {
				t.Errorf("expected text %q, got %q", tt.expected, v)
			}
		})
	}
}

func TestHighlightInvalid(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "no text in highlight",
			source: "====",
		},
		{
			name:   "only delimiters",
			source: "========",
		},
		{
			name:   "imbalanced delimiters",
			source: "===",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithExtensions(&Extender{}),
			)
			source := []byte(tt.source)
			doc := md.Parser().Parse(text.NewReader(source))

			found := testutil.FindNode[*Highlight](doc)

			if found != nil {
				t.Errorf("Callout node should not be found for source: %q", tt.source)
			}
		})
	}
}
