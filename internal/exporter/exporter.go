package exporter

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/ongyx/knap/internal/converter"
	"github.com/ongyx/knap/internal/schema"
	"github.com/ongyx/knap/internal/util"
)

const (
	jsonExtension    = ".json"
	metadataFilename = "metadata" + jsonExtension
)

// Exporter exports an Obsidian vault to Outline format.
type Exporter struct {
	identity    schema.Identity
	vault       *Vault
	ftcSettings *FTCSettings

	collection *schema.Collection
	documents  map[string]*schema.Document
	navnodes   map[string]*schema.NavigationNode
}

// Creates a new exporter with the given identity.
func New(identity schema.Identity, vaultPath string) (*Exporter, error) {
	v := NewVault(vaultPath)

	// Try to read settings from the Fast Text Color plugin.
	st, err := NewFTCSettings(vaultPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// The settings file exists but it can't be read for some other reason.
		return nil, fmt.Errorf("fast text color settings exists in vault but parsing failed: %w", err)
	}

	return &Exporter{
		identity:    identity,
		vault:       v,
		ftcSettings: st,
		documents:   make(map[string]*schema.Document),
		navnodes:    make(map[string]*schema.NavigationNode),
	}, nil
}

// Returns the identity of the Outline user exporting the vault.
func (e *Exporter) Identity() schema.Identity {
	return e.identity
}

// Returns the vault being exported.
func (e *Exporter) Vault() *Vault {
	return e.vault
}

// Returns the Fast Text Color settings. If the plugin is not installed, this is nil.
func (e *Exporter) FTCSettings() *FTCSettings {
	return e.ftcSettings
}

// Exports the vault to a ZIP file importable by Outline to the writer w.
func (e *Exporter) Export(w io.Writer) error {
	cid, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}
	curlid := util.NewURLID()
	cname := filepath.Base(e.vault.Path())

	e.collection = schema.NewCollection(cid, curlid, cname)
	clear(e.documents)
	clear(e.navnodes)

	// Scan the vault for notes and attachements to export.
	if err := e.vault.Scan(); err != nil {
		return fmt.Errorf("failed to scan vault at %q: %w", e.vault.Path(), err)
	}

	// Generate a document for each note in the vault.
	for vf := range e.vault.Files() {
		if vf.FileFormat != util.FileNote {
			continue
		}
		if err := e.generateNoteDocument(vf); err != nil {
			return fmt.Errorf("failed to generate document for note at %q: %w", vf.AbsPath, err)
		}
	}

	zw := zip.NewWriter(w)

	m := schema.NewExportMetadata(e.identity)

	// Write the export metadata.
	mf, err := zw.Create(metadataFilename)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}

	me := json.NewEncoder(mf)
	if err := me.Encode(m); err != nil {
		return fmt.Errorf("failed to marshal JSON into collection file: %w", err)
	}

	// Write the collection to a JSON file of the same name.
	cf, err := zw.Create(cname + jsonExtension)
	if err != nil {
		return fmt.Errorf("failed to create collection file: %w", err)
	}

	ce := json.NewEncoder(cf)
	if err := ce.Encode(e.collection); err != nil {
		return fmt.Errorf("failed to marshal JSON into collection file: %w", err)
	}

	return zw.Close()
}

// Registers a document for export. If pdoc and pnn are non-nil, they are registered as the parent document and navnode.
func (e *Exporter) registerDocument(
	relpath string,
	doc *schema.Document,
	nn *schema.NavigationNode,
	pdoc *schema.Document,
	pnn *schema.NavigationNode,
) {
	meta := e.collection.Metadata
	docs := e.collection.Documents

	if pdoc != nil && pnn != nil {
		// Indicate the document to be a child of the parent document, and add its navigation node as a child of the parent's.
		doc.ParentDocumentId = &pdoc.ID
		pnn.Children = append(pnn.Children, nn)
	} else {
		// Add the document's navigation node to the root of the collection.
		meta.DocumentStructure = append(meta.DocumentStructure, nn)
	}

	e.documents[relpath] = doc
	e.navnodes[relpath] = nn

	docs[doc.ID] = doc
}

// Registers attachments belonging to a document for export.
func (e *Exporter) registerAttachments(files iter.Seq[*VaultFile], doc *schema.Document) {
	for vf := range files {
		if _, ok := e.collection.Attachments[vf.ID]; ok {
			// Attachment has already been registered.
			continue
		}

		a := &schema.Attachment{
			UserID:      doc.CreatedById,
			DocumentID:  doc.ID,
			ContentType: vf.ContentType,
			Name:        filepath.Base(vf.AbsPath),
			ID:          vf.ID,
			Size:        vf.Size,
		}
		a.UpdateKey()
	}
}

// Generates a document from an Obsidian note.
func (e *Exporter) generateNoteDocument(note *VaultFile) error {
	if _, ok := e.documents[note.RelPath]; ok {
		// Document has already been exported?
		return nil
	}

	pdoc, pnn, err := e.generateParentDocument(note)
	if err != nil {
		return err
	}

	doc := schema.NewDocument(note.ID, note.URLID, note.Title(), e.identity)
	nn := schema.NewNavigationNode(doc)

	if err := doc.SetTimestamps(note.AbsPath); err != nil {
		return fmt.Errorf("failed to read timestamps on note at %q: %w", note.AbsPath, err)
	}

	src, err := os.ReadFile(note.AbsPath)
	if err != nil {
		return fmt.Errorf("failed to read note at %q: %w", note.AbsPath, err)
	}

	res := NewNoteResolver(e, note)
	cv := converter.New(res)
	node, err := cv.Convert(src)
	if err != nil {
		return fmt.Errorf("failed to convert note at %q: %w", note.AbsPath, err)
	}

	doc.Data = node

	e.registerDocument(note.RelPath, doc, nn, pdoc, pnn)
	e.registerAttachments(res.Attachments().Items(), doc)

	return nil
}

// Generates a parent document for each parent directory the note is in, and returns the document and navnode of the immediate parent.
// If the note does not have a parent directory, the document and navnode will be nil.
func (e *Exporter) generateParentDocument(note *VaultFile) (*schema.Document, *schema.NavigationNode, error) {
	dir, _ := path.Split(note.RelPath)
	if dir == "" {
		// The note does not have any parent directory.
		return nil, nil, nil
	}

	parents := strings.Split(path.Clean(dir), "/")

	var (
		pdoc *schema.Document
		pnn  *schema.NavigationNode
	)

	// Create a parent document for each parent directory in the path.
	for i := range parents {
		relpath := path.Join(parents[:i+1]...)

		doc, ok := e.documents[relpath]
		nn := e.navnodes[relpath]
		if ok {
			pdoc = doc
			pnn = nn
		} else {
			doc, nn, err := e.generateEmptyDocument(relpath)
			if err != nil {
				return nil, nil, err
			}

			e.registerDocument(relpath, doc, nn, pdoc, pnn)
			pdoc = doc
			pnn = nn
		}
	}

	return pdoc, pnn, nil
}

// Generates an empty document for the relative path.
func (e *Exporter) generateEmptyDocument(relpath string) (*schema.Document, *schema.NavigationNode, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, nil, err
	}
	urlid := util.NewURLID()

	doc := schema.NewDocument(id, urlid, path.Base(relpath), e.identity)

	abs := filepath.Join(e.vault.Path(), relpath)
	if err := doc.SetTimestamps(abs); err != nil {
		return nil, nil, fmt.Errorf("failed to read timestamps of %q: %w", abs, err)
	}

	nn := schema.NewNavigationNode(doc)

	return doc, nn, nil
}
