package converter

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/ongyx/knap/internal/testutil"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/wikilink"
)

func TestSuccess(t *testing.T) {
	tests := []struct {
		name     string
		nodeType reflect.Type
		source   string
		expected *Link
	}{
		{
			name:     "Markdown link",
			nodeType: reflect.TypeFor[*ast.Link](),
			source:   "[hello world!](https://google.com)",
			expected: &Link{
				URL: &url.URL{
					Scheme: "https",
					Host:   "google.com",
				},
				Text:  "hello world!",
				Embed: false,
			},
		},
		{
			name:     "Markdown image (embed)",
			nodeType: reflect.TypeFor[*ast.Image](),
			source:   "![alt text|256x256](https://commons.wikimedia.org/wiki/File:Test-Logo.svg)",
			expected: &Link{
				URL: &url.URL{
					Scheme: "https",
					Host:   "commons.wikimedia.org",
					Path:   "/wiki/File:Test-Logo.svg",
				},
				Text:  "alt text|256x256",
				Embed: true,
			},
		},
		{
			name:     "Wikilink",
			nodeType: reflect.TypeFor[*wikilink.Node](),
			source:   "[[Nishino#hello.. i am 2 kilobytes i am tiny!|Umazing!]]",
			expected: &Link{
				URL: &url.URL{
					Path:     "Nishino",
					Fragment: "hello.. i am 2 kilobytes i am tiny!",
				},
				Text:  "Umazing!",
				Embed: false,
			},
		},
		{
			name:     "Wikilink (embed)",
			nodeType: reflect.TypeFor[*wikilink.Node](),
			source:   "![[honse.png|How hungry...]]",
			expected: &Link{
				URL: &url.URL{
					Path: "honse.png",
				},
				Text:  "How hungry...",
				Embed: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(goldmark.WithExtensions(&wikilink.Extender{}))
			src := []byte(tt.source)
			doc := md.Parser().Parse(text.NewReader(src))

			found := testutil.FindChildNodeReflect(doc, tt.nodeType)
			if found == nil {
				t.Errorf("Node type %s not found in document", tt.nodeType)
			}

			got, err := ParseLinkFromNode(found, src)
			if err != nil {
				t.Errorf("Error encounted while parsing link: %s", err)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Expected link %#v, got link %#v", tt.expected, got)
			}
		})
	}
}

func TestFailure(t *testing.T) {
	tests := []struct {
		name     string
		nodeType reflect.Type
		source   string
	}{
		{
			name:     "invalid URL",
			nodeType: reflect.TypeFor[*ast.Link](),
			source:   "[invalid](:invalid)",
		},
		{
			name:     "invalid node",
			nodeType: reflect.TypeFor[*ast.Heading](),
			source:   "# oops",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := goldmark.New(goldmark.WithExtensions(&wikilink.Extender{}))
			src := []byte(tt.source)
			doc := md.Parser().Parse(text.NewReader(src))

			found := testutil.FindChildNodeReflect(doc, tt.nodeType)
			if found == nil {
				t.Errorf("Node type %s not found in document", tt.nodeType)
			}

			_, err := ParseLinkFromNode(found, src)
			if err == nil {
				t.Errorf("Expected parsing error")
			}
		})
	}
}
