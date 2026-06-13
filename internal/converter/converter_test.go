package converter

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ongyx/knap/internal/prosemirror"
)

func TestConverter_Convert(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected *prosemirror.Node
	}{
		{
			name:     "simple paragraph",
			markdown: "Hello world",
			expected: &prosemirror.Node{
				Type: prosemirror.NodeDocument,
				Content: []*prosemirror.Node{
					{
						Type: prosemirror.NodeParagraph,
						Content: []*prosemirror.Node{
							{
								Type: prosemirror.NodeText,
								Text: "Hello world",
							},
						},
					},
				},
			},
		},
		{
			name:     "heading",
			markdown: "# Title",
			expected: &prosemirror.Node{
				Type: prosemirror.NodeDocument,
				Content: []*prosemirror.Node{
					{
						Type: prosemirror.NodeHeading,
						Attrs: map[string]any{
							"level": float64(1),
						},
						Content: []*prosemirror.Node{
							{
								Type: prosemirror.NodeText,
								Text: "Title",
							},
						},
					},
				},
			},
		},
		{
			name:     "bold and italic",
			markdown: "**bold** and *italic*",
			expected: &prosemirror.Node{
				Type: prosemirror.NodeDocument,
				Content: []*prosemirror.Node{
					{
						Type: prosemirror.NodeParagraph,
						Content: []*prosemirror.Node{
							{
								Type: prosemirror.NodeText,
								Text: "bold",
								Marks: []prosemirror.Mark{
									{Type: "strong"},
								},
							},
							{
								Type: prosemirror.NodeText,
								Text: " and ",
							},
							{
								Type: prosemirror.NodeText,
								Text: "italic",
								Marks: []prosemirror.Mark{
									{Type: "em"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "bullet list",
			markdown: "- item 1\n- item 2",
			expected: &prosemirror.Node{
				Type: prosemirror.NodeDocument,
				Content: []*prosemirror.Node{
					{
						Type: prosemirror.NodeBulletList,
						Content: []*prosemirror.Node{
							{
								Type: prosemirror.NodeListItem,
								Content: []*prosemirror.Node{
									{
										Type: prosemirror.NodeParagraph,
										Content: []*prosemirror.Node{
											{
												Type: prosemirror.NodeText,
												Text: "item 1",
											},
										},
									},
								},
							},
							{
								Type: prosemirror.NodeListItem,
								Content: []*prosemirror.Node{
									{
										Type: prosemirror.NodeParagraph,
										Content: []*prosemirror.Node{
											{
												Type: prosemirror.NodeText,
												Text: "item 2",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "notice",
			markdown: "> [!info]\n> info content",
			expected: &prosemirror.Node{
				Type: prosemirror.NodeDocument,
				Content: []*prosemirror.Node{
					{
						Type: prosemirror.NodeNotice,
						Attrs: map[string]any{
							"style": "info",
						},
						Content: []*prosemirror.Node{
							{
								Type: prosemirror.NodeParagraph,
								Content: []*prosemirror.Node{
									{
										Type: prosemirror.NodeText,
										Text: "info content",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:     "empty",
			markdown: "",
			expected: &prosemirror.Node{
				Type: prosemirror.NodeDocument,
				Content: []*prosemirror.Node{
					{
						Type: prosemirror.NodeParagraph,
					},
				},
			},
		},
		{
			name:     "link with formatting",
			markdown: "**_[Hello World!](https://google.com)_**",
			expected: &prosemirror.Node{
				Type: prosemirror.NodeDocument,
				Content: []*prosemirror.Node{
					{
						Type: prosemirror.NodeParagraph,
						Content: []*prosemirror.Node{
							{
								Type: prosemirror.NodeText,
								Text: "Hello World!",
								Marks: []prosemirror.Mark{
									prosemirror.NewBoldMark(),
									prosemirror.NewItalicMark(),
									prosemirror.NewLinkMark("https://google.com"),
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := New(nil)
			got, err := cv.Convert([]byte(tt.markdown))
			if err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			// Use JSON marshaling for deep comparison as prosemirror.Node might have complex nesting
			gotJSON, _ := json.Marshal(got)
			expectedJSON, _ := json.Marshal(tt.expected)

			if string(gotJSON) != string(expectedJSON) {
				t.Errorf("Convert() got = %s, want %s", string(gotJSON), string(expectedJSON))
			}
		})
	}
}

func TestConverter_Convert_Error(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		wantErr  error
	}{
		{
			name:     "unsupported HTML",
			markdown: "<div></div>",
			wantErr:  ErrInvalidHTML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := New(nil)
			_, err := cv.Convert([]byte(tt.markdown))
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
