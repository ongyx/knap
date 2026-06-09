package exporter

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/ongyx/knap/internal/schema"
)

// TODO: Implement PDF embeds

// VaultFile is an entry for a file within an Obsidian vault.
type VaultFile struct {
	// The absolute path to the file. This uses the OS' path separator.
	AbsPath string
	// The relative path to the file from the vault root. The path separator is a slash '/' by Obsidian convention.
	RelPath string
	// The type of file this is.
	FileType VaultFileType
	// The UUID generated for this file.
	ID uuid.UUID
	// The URLID generated for this file.
	URLID schema.URLID
}

// Creates a new vault file entry for the given file path and vault path.
func NewVaultFile(path, vaultPath string) (*VaultFile, error) {
	rel, err := filepath.Rel(vaultPath, path)
	if err != nil {
		return nil, err
	}

	ft := VaultFileOther
	ext := filepath.Ext(rel)
	// Check the file extension to determine the file type.
	if ext == noteExtension {
		ft = VaultFileNote
	} else if imageExtensions.Contains(ext) {
		ft = VaultFileImage
	} else if videoExtensions.Contains(ext) {
		ft = VaultFileVideo
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	urlid := schema.NewURLID()

	return &VaultFile{
		AbsPath:  path,
		RelPath:  filepath.ToSlash(rel),
		FileType: ft,
		ID:       id,
		URLID:    urlid,
	}, nil
}

// Returns the basename of the vault file without its extension, suitable for a document title.
func (vf *VaultFile) Title() string {
	base := filepath.Base(vf.AbsPath)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// Generates the Outline document URL for this vault file.
//
// For more details, please refer to: https://github.com/outline/outline/blob/88de417a21c260e32ecc9c89d756661c54064603/server/models/Document.ts#L430
func (vf *VaultFile) DocumentURL() string {
	title := vf.Title()
	slug := Slugify(title, nil)

	urlid := string(vf.URLID)

	if len(slug) > 0 {
		return fmt.Sprintf(`/doc/%s-%s`, slug, urlid)
	} else {
		return fmt.Sprintf(`/doc/untitled-%s`, urlid)
	}
}
