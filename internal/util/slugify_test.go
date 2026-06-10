package util

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		options  *SlugifyOptions
		expected string
	}{
		{
			name:     "default options",
			str:      "Hello World!",
			options:  nil,
			expected: "Hello-World!",
		},
		{
			name:     "lowercase",
			str:      "Hello World!",
			options:  &SlugifyOptions{Lower: true},
			expected: "hello-world!",
		},
		{
			name:     "custom replacement",
			str:      "Hello World!",
			options:  &SlugifyOptions{Replacement: "_"},
			expected: "Hello_World!",
		},
		{
			name:     "remove special characters",
			str:      "Hello @ World!",
			options:  nil,
			expected: "Hello-@-World!",
		},
		{
			name:     "multiple spaces",
			str:      "Hello   World",
			options:  nil,
			expected: "Hello-World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.str, tt.options); got != tt.expected {
				t.Errorf("Slugify(%q) = %q, want %q", tt.str, got, tt.expected)
			}
		})
	}
}
