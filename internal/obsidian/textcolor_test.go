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

func TestTextColor(t *testing.T) {
	tests := []struct {
		name          string
		source        string
		expectedColor string
		expectedText  string
	}{
		{
			name:          "simple color",
			source:        "~={red}Hello=~",
			expectedColor: "red",
			expectedText:  "Hello",
		},
		{
			name:          "color with space",
			source:        "~={blue} world =~",
			expectedColor: "blue",
			expectedText:  "world ",
		},
		{
			name:          "numeric color name",
			source:        "~={123}test=~",
			expectedColor: "123",
			expectedText:  "test",
		},
		{
			name:          "color with space after name",
			source:        "~={green}  more text=~",
			expectedColor: "green",
			expectedText:  "more text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithParserOptions(
					parser.WithInlineParsers(
						util.Prioritized(NewTextColorParser(), 50),
					),
				),
			)
			source := []byte(tt.source)
			doc := md.Parser().Parse(text.NewReader(source))

			color := testutil.FindNode[*TextColor](doc)
			if color == nil {
				t.Fatal("TextColor node not found")
			}

			if string(color.ID) != tt.expectedColor {
				t.Errorf("expected color %q, got %q", tt.expectedColor, string(color.ID))
			}

			txtNode := color.FirstChild().(*ast.Text)
			if txtNode == nil {
				t.Fatal("Expected text node after TextColor")
			}
			v := string(txtNode.Segment.Value(source))
			if v != tt.expectedText {
				t.Errorf("expected text %q, got %q", tt.expectedText, v)
			}
		})
	}
}

func TestTextColorInvalid(t *testing.T) {
	tests := []struct {
		name   string
		source string
	}{
		{
			name:   "missing color",
			source: "~=text=~",
		},
		{
			name:   "empty color",
			source: "~={}text=~",
		},
		{
			name:   "missing span",
			source: "{red}text",
		},
		{
			name:   "color not at start",
			source: "~=text {red}=~",
		},
		{
			name:   "invalid color name",
			source: "~={red blue}text=~",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(
				goldmark.WithParserOptions(
					parser.WithInlineParsers(
						util.Prioritized(NewTextColorParser(), 50),
					),
				),
			)
			source := []byte(tt.source)
			doc := md.Parser().Parse(text.NewReader(source))

			color := testutil.FindNode[*TextColor](doc)
			if color != nil {
				t.Errorf("TextColor node should not be found for source: %q", tt.source)
			}
		})
	}
}
