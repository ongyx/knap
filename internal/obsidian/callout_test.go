package obsidian

import (
	"testing"

	"github.com/ongyx/knap/internal/testutil"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

func TestCallout(t *testing.T) {
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
				goldmark.WithParserOptions(
					parser.WithInlineParsers(
						util.Prioritized(NewCalloutParser(), 50),
					),
				),
			)
			source := []byte(tt.source)
			doc := md.Parser().Parse(text.NewReader(source))

			found := testutil.FindChildNode[*Callout](doc)
			if found == nil {
				t.Fatal("Callout node not found in AST")
			}

			n := string(found.Name)
			if n != tt.expected {
				t.Errorf("expected callout name %q, got %q", tt.expected, n)
			}
		})
	}
}

func TestCalloutInvalid(t *testing.T) {
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

			found := testutil.FindChildNode[*Callout](doc)

			if found != nil {
				t.Errorf("Callout node should not be found for source: %q", tt.source)
			}
		})
	}
}
