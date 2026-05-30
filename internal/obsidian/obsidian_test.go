package obsidian

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func TestObsidianCallout(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "simple callout",
			source:   "> [!info]\n> content",
			expected: "info",
		},
		{
			name:     "callout with different name",
			source:   "> [!warning]\n> alert",
			expected: "warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithExtensions(&Extender{}),
			)
			source := []byte(tt.source)
			doc := md.Parser().Parse(text.NewReader(source))

			var found *Callout
			ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
				if !entering {
					return ast.WalkContinue, nil
				}
				if c, ok := n.(*Callout); ok {
					found = c
					return ast.WalkStop, nil
				}
				return ast.WalkContinue, nil
			})

			if found == nil {
				t.Fatal("Callout node not found in AST")
			}

			if string(found.Name) != tt.expected {
				t.Errorf("expected callout name %q, got %q", tt.expected, string(found.Name))
			}
		})
	}
}

func TestObsidianCallout_Invalid(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "not in blockquote",
			source: "[!info]\ncontent",
		},
		{
			name:   "not at start of blockquote",
			source: "> some text\n> [!info]",
		},
		{
			name:   "missing exclamation",
			source: "> [info]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithExtensions(&Extender{}),
			)
			source := []byte(tt.source)
			doc := md.Parser().Parse(text.NewReader(source))

			var found *Callout
			ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
				if !entering {
					return ast.WalkContinue, nil
				}
				if _, ok := n.(*Callout); ok {
					found = n.(*Callout)
					return ast.WalkStop, nil
				}
				return ast.WalkContinue, nil
			})

			if found != nil {
				t.Errorf("Callout node should not be found for source: %q", tt.source)
			}
		})
	}
}
