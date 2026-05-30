package markdown

import (
	"bytes"
	"testing"

	"github.com/yuin/goldmark"
)

func TestFencedCodeBlockRenderer(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(&Prosemirror),
	)

	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "basic fenced code block",
			source:   "```\ncode\n```",
			expected: "<pre><code>code\n</code></pre>\n",
		},
		{
			name:     "fenced code block with language",
			source:   "```go\nfunc main() {}\n```",
			expected: `<pre><code class="language-go">func main() {}` + "\n" + `</code></pre>` + "\n",
		},
		{
			name:     "fenced code block with multiple lines",
			source:   "```\nline 1\nline 2\n```",
			expected: "<pre><code>line 1\nline 2\n</code></pre>\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := md.Convert([]byte(tt.source), &buf); err != nil {
				t.Fatalf("failed to convert: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tt.expected, buf.String())
			}
		})
	}
}

func TestThematicBreakRenderer(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(&Prosemirror),
	)

	tests := []struct {
		name     string
		source   string
		expected string
	}{
		{
			name:     "---",
			source:   "---",
			expected: "<hr>\n",
		},
		{
			name:     "___",
			source:   "___",
			expected: "<hr>\n",
		},
		{
			name:     "***",
			source:   "***",
			expected: `<hr class="pagebreak">` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := md.Convert([]byte(tt.source), &buf); err != nil {
				t.Fatalf("failed to convert: %v", err)
			}

			if buf.String() != tt.expected {
				t.Errorf("expected:\n%q\ngot:\n%q", tt.expected, buf.String())
			}
		})
	}
}
