package converter

import (
	"testing"

	"github.com/ongyx/knap/internal/obsidian"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestParseEmbedSizeSuccess(t *testing.T) {
	tests := []struct {
		name           string
		desc           string
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "width",
			desc:           "256",
			expectedWidth:  256,
			expectedHeight: 0,
		},
		{
			name:           "width and height",
			desc:           "1920x1080",
			expectedWidth:  1920,
			expectedHeight: 1080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h, ok := ParseEmbedSize(tt.desc)
			if !ok {
				t.Errorf("Failed to parse %q", tt.desc)
			}

			if w != tt.expectedWidth || h != tt.expectedHeight {
				t.Errorf("Expected size (%d, %d), got size (%d, %d)", tt.expectedWidth, tt.expectedHeight, w, h)
			}
		})
	}
}

func TestParseEmbedSizeFail(t *testing.T) {
	tests := []struct {
		name string
		desc string
	}{
		{
			name: "non-digit",
			desc: "256px",
		},
		{
			name: "x without width",
			desc: "256x",
		},
		{
			name: "empty string",
			desc: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, ok := ParseEmbedSize(tt.desc)
			if ok {
				t.Errorf("Expected failure to parse %q", tt.desc)
			}
		})
	}
}

func TestNodeChildrenToText(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "link",
			source:   "[foo bar baz](https://google.com)",
			expected: "foo bar baz",
		},
		{
			name:     "emphasis",
			source:   "**bold** _italic_ ~~strikethrough~~ ==highlight==",
			expected: "bold italic strikethrough highlight",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(obsidian.DefaultOptions())
			src := []byte(tt.source)
			doc := md.Parser().Parse(text.NewReader(src))

			plaintext := NodeChildrenToText(doc, src)
			if plaintext != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, plaintext)
			}
		})
	}
}
