package converter

import (
	"encoding/json"
	"testing"

	"github.com/ongyx/knap/internal/schema"
)

func TestConverter_Convert(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected *schema.Node
	}{
		{
			name:     "simple paragraph",
			markdown: "Hello world",
			expected: &schema.Node{
				Type: schema.NodeDocument,
				Content: []*schema.Node{
					{
						Type: schema.NodeParagraph,
						Content: []*schema.Node{
							{
								Type: schema.NodeText,
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
			expected: &schema.Node{
				Type: schema.NodeDocument,
				Content: []*schema.Node{
					{
						Type: schema.NodeHeading,
						Attrs: map[string]any{
							"level": float64(1),
						},
						Content: []*schema.Node{
							{
								Type: schema.NodeText,
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
			expected: &schema.Node{
				Type: schema.NodeDocument,
				Content: []*schema.Node{
					{
						Type: schema.NodeParagraph,
						Content: []*schema.Node{
							{
								Type: schema.NodeText,
								Text: "bold",
								Marks: []schema.Mark{
									{Type: "strong"},
								},
							},
							{
								Type: schema.NodeText,
								Text: " and ",
							},
							{
								Type: schema.NodeText,
								Text: "italic",
								Marks: []schema.Mark{
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
			expected: &schema.Node{
				Type: schema.NodeDocument,
				Content: []*schema.Node{
					{
						Type: schema.NodeBulletList,
						Content: []*schema.Node{
							{
								Type: schema.NodeListItem,
								Content: []*schema.Node{
									{
										Type: schema.NodeParagraph,
										Content: []*schema.Node{
											{
												Type: schema.NodeText,
												Text: "item 1",
											},
										},
									},
								},
							},
							{
								Type: schema.NodeListItem,
								Content: []*schema.Node{
									{
										Type: schema.NodeParagraph,
										Content: []*schema.Node{
											{
												Type: schema.NodeText,
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
			expected: &schema.Node{
				Type: schema.NodeDocument,
				Content: []*schema.Node{
					{
						Type: schema.NodeNotice,
						Attrs: map[string]any{
							"style": "info",
						},
						Content: []*schema.Node{
							{
								Type: schema.NodeParagraph,
								Content: []*schema.Node{
									{
										Type: schema.NodeText,
										Text: "info content",
									},
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

			// Use JSON marshaling for deep comparison as schema.Node might have complex nesting
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
			if err != tt.wantErr {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
