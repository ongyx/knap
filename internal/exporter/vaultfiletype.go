package exporter

const (
	// The file is of an unknown type.
	VaultFileOther VaultFileType = iota
	// The file is an Obsidian note.
	VaultFileNote
	// The file is an image.
	VaultFileImage
	// The file is a video.
	VaultFileVideo
)

// A type of file in an Obsidian vault.
type VaultFileType int
