package util

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/ongyx/knap/internal/collections"
)

// File extension for Obsidian notes in a vault.
const NoteExtension = ".md"

const (
	// The file is of an unknown type.
	FileOther FileFormat = iota
	// The file is an Obsidian note.
	FileNote
	// The file is an image.
	FileImage
	// The file is an audio.
	FileAudio
	// The file is a video.
	FileVideo
)

// https://obsidian.md/help/file-formats
var (
	// Recognized file extensions for images.
	ImageExtensions = collections.NewSet(".avif", ".bmp", ".gif", ".jpeg", ".jpg", ".png", ".svg", ".webp")
	// Recognized file extensions for audio.
	AudioExtensions = collections.NewSet(".flac", ".m4a", ".mp3", ".ogg", ".wav", ".webm", ".3gp")
	// Recognized file extensions for video.
	VideoExtensions = collections.NewSet(".mkv", ".mov", ".mp4", ".ogv", ".webm")
)

// A file format.
type FileFormat int

// Parses the file extension from a path to determine a file format.
func ParseFileFormat(p string) FileFormat {
	ff := FileOther
	ext := strings.ToLower(path.Ext(filepath.ToSlash(p)))
	// Check the file extension to determine the file type.
	if ext == NoteExtension {
		ff = FileNote
	} else if ImageExtensions.Contains(ext) {
		ff = FileImage
	} else if AudioExtensions.Contains(ext) {
		ff = FileAudio
	} else if VideoExtensions.Contains(ext) {
		ff = FileVideo
	}

	return ff
}
