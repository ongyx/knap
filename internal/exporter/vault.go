package exporter

import (
	"errors"
	"io/fs"
	"iter"
	"maps"
	"path"
	"path/filepath"
	"slices"

	"github.com/ongyx/knap/internal/util"
)

const (
	obsidianConfigFolder = ".obsidian"
)

// Error returned by Vault.Scan() when the folder is not an Obsidian vault.
var ErrNotAVault = errors.New("folder is not an Obsidian vault")

// Vault scans an Obsidian vault for notes and attachments to export.
//
// # Notes
//
// When one or more collections are exported as JSON in Outline, all files are packed into a ZIP file with this directory structure:
//
//	(ZIP file)
//	|- uploads/(user UUID)/(attachment UUID)/(attachment file)
//	|- metadata.json
//	|- (collection name).json
//
// where:
//   - uploads contains uploaded files by each user,
//   - metadata.json is the export metadata,
//   - and the rest of the .json files describe each collection exported.
type Vault struct {
	path        string
	allFiles    map[string]*VaultFile
	uniqueFiles map[string]*VaultFile
}

// Creates a new vault scanner for the given vault path.
func NewVault(vaultPath string) *Vault {
	return &Vault{
		path:        vaultPath,
		allFiles:    make(map[string]*VaultFile),
		uniqueFiles: make(map[string]*VaultFile),
	}
}

// Returns the path to the vault.
func (v *Vault) Path() string {
	return v.path
}

// Scans the vault for notes and attachments to export.
func (v *Vault) Scan() error {
	clear(v.allFiles)
	clear(v.uniqueFiles)

	// Make sure the path is actually is a vault.
	cf := path.Join(v.path, obsidianConfigFolder)
	exists, err := FolderExists(cf)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotAVault
	}

	return filepath.WalkDir(v.path, v.walkDir)
}

// Looks up a file in the vault by its name.
// Name can either be a unique filename or relative path within the vault, as per https://obsidian.md/help/links.
//
// If the name does not exist, nil is returned.
func (v *Vault) Lookup(name string) *VaultFile {
	// NOTE: nil means that either the name doesn't exist, or that the filename is not unique.
	// See Vault.walkDir().
	if vf := v.uniqueFiles[name]; vf != nil {
		return vf
	}

	if filepath.Ext(name) == "" {
		// Notes may elide the '.md' file extension in their name.
		notename := name + util.NoteExtension

		if vf := v.uniqueFiles[notename]; vf != nil {
			return vf
		}

		if vf := v.allFiles[notename]; vf != nil {
			return vf
		}
	}

	// The filename is not unique, so attempt to look it up as a relative path.
	return v.allFiles[name]
}

// Returns an iterator over all files in the vault in lexical order.
func (v *Vault) Files() iter.Seq[*VaultFile] {
	return func(yield func(*VaultFile) bool) {
		keys := slices.Sorted(maps.Keys(v.allFiles))
		for _, k := range keys {
			if !yield(v.allFiles[k]) {
				return
			}
		}
	}
}

func (v *Vault) walkDir(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() {
		// Don't walk the config folder.
		if filepath.Base(path) == obsidianConfigFolder {
			return fs.SkipDir
		}
		return nil
	}

	vf, err := NewVaultFile(path, v.path)
	if err != nil {
		return err
	}

	v.allFiles[vf.RelPath] = vf

	// In Obsidian, linktexts can either be the filename or the full path if the filename is not unique.
	// All filenames are considered unqiue until an identical filename is found in another folder.
	// https://docs.obsidian.md/Reference/TypeScript+API/MetadataCache/fileToLinktext
	n := filepath.Base(vf.RelPath)
	if _, ok := v.uniqueFiles[n]; ok {
		// Another file with the same name was already found, so it's no longer unique.
		v.uniqueFiles[n] = nil
	} else {
		// Add the new unique file.
		v.uniqueFiles[n] = vf
	}

	return nil
}
