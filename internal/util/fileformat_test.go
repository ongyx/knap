package util

import "testing"

func TestParseFileFormat(t *testing.T) {
	tests := []struct {
		path     string
		expected FileFormat
	}{
		{"note.md", FileNote},
		{"image.png", FileImage},
		{"image.JPG", FileImage}, // Wait, is it case-sensitive?
		{"video.mp4", FileVideo},
		{"video.wav", FileAudio},
		{"other.txt", FileOther},
		{"sub/dir/note.md", FileNote},
		{"sub\\dir\\note.md", FileNote},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := ParseFileFormat(tt.path); got != tt.expected {
				t.Errorf("ParseFileFormat(%q) = %v, want %v", tt.path, got, tt.expected)
			}
		})
	}
}
