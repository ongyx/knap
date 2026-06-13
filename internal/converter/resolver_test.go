package converter

import (
	"net/url"
	"testing"

	"github.com/ongyx/knap/internal/obsidian"
	"github.com/ongyx/knap/internal/testutil"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func TestDefaultResolverResolveInternalLink(t *testing.T) {
	tests := []struct {
		name     string
		link     *Link
		expected string
	}{
		{
			name: "normal link",
			link: &Link{
				URL: &url.URL{
					Path: "manhattan",
				},
				Text: "cafe",
			},
			expected: "[[manhattan|cafe]]",
		},
		{
			name: "embed link",
			link: &Link{
				URL: &url.URL{
					Path: "agnes",
				},
				Text:  "tachyon",
				Embed: true,
			},
			expected: "![[agnes|tachyon]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node, err := DefaultResolver.ResolveInternalLink(tt.link, nil)
			if err != nil {
				t.Errorf("failed to resolve internal link: %v", err)
			}

			if node.Text != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, node.Text)
			}
		})
	}
}

func TestDefaultResolverResolveColor(t *testing.T) {
	md := goldmark.New(obsidian.DefaultOptions())
	src := []byte("~={mycolor}colored text=~")
	doc := md.Parser().Parse(text.NewReader(src)).(*ast.Document)

	tc := testutil.FindChildNode[*obsidian.TextColor](doc)
	if tc == nil {
		t.Error("failed to locate text color node?")
	}

	clr := DefaultResolver.ResolveColor(doc, tc)
	if clr != "" {
		t.Errorf("expected empty string, got %q", clr)
	}
}
