package exporter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/ongyx/knap/internal/util"
)

// TODO: Implement PDF embeds

// VaultFile is an entry for a file within an Obsidian vault.
type VaultFile struct {
	// The absolute path to the file. This uses the OS' path separator.
	AbsPath string
	// The relative path to the file from the vault root. The path separator is a slash '/' by Obsidian convention.
	RelPath string
	// The file format.
	FileFormat util.FileFormat
	// The size of the file.
	Size int64
	// The mimetype.
	MimeType *mimetype.MIME

	// The UUID generated for this file.
	ID uuid.UUID
	// The URLID generated for this file.
	URLID util.URLID
}

// Creates a new vault file entry for the given file path and vault path.
func NewVaultFile(path, vaultPath string) (*VaultFile, error) {
	rel, err := filepath.Rel(vaultPath, path)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	urlid := util.NewURLID()

	// Probe the file for its size and MIME content type.
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return nil, err
	}

	mt, err := mimetype.DetectReader(f)
	if err != nil {
		return nil, err
	}

	return &VaultFile{
		AbsPath:    path,
		RelPath:    filepath.ToSlash(rel),
		FileFormat: util.ParseFileFormat(path),
		Size:       st.Size(),
		MimeType:   mt,
		ID:         id,
		URLID:      urlid,
	}, nil
}

// Returns the basename of the vault file without its extension, suitable for a document title.
func (vf *VaultFile) Title() string {
	base := filepath.Base(vf.AbsPath)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
