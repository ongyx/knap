package exporter

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/ongyx/knap/internal/collections"
	"github.com/ongyx/knap/internal/converter"
	"github.com/ongyx/knap/internal/schema"
	"github.com/ongyx/knap/internal/util"
)

const (
	jsonExtension    = ".json"
	metadataFilename = "metadata" + jsonExtension
)

// Error returned when the vault path is invalid.
var ErrInvalidVaultPath = errors.New("vault path is invalid")

// Options for exporting a vault.
type ExporterOptions struct {
	// The path to the vault. This is a requried field.
	VaultPath string

	// The identity to export as.
	Identity *schema.Identity
	// The logger to log to.
	Logger *log.Logger
	// The folders to ignore within the vault.
	Ignore collections.Set[string]
}

// Applies defaults to zero values in the exporter options.
func (o *ExporterOptions) Defaults() *ExporterOptions {
	if o.Identity == nil {
		o.Identity = &schema.Identity{}
	}
	o.Identity.Defaults()

	if o.Logger == nil {
		o.Logger = log.Default()
	}

	if o.Ignore.Len() == 0 {
		o.Ignore = collections.NewSet[string]()
	}

	return o
}

// Exporter exports an Obsidian vault to Outline format.
type Exporter struct {
	options     *ExporterOptions
	vault       *Vault
	ftcSettings *FTCSettings

	collection    *schema.Collection
	documentCache map[string]*schema.Document
	navnodeCache  map[string]*schema.NavigationNode
}

// Creates a new exporter with the given options.
func New(o *ExporterOptions) (*Exporter, error) {
	o.Defaults()

	if o.VaultPath == "" {
		return nil, ErrInvalidVaultPath
	}

	// Try to read settings from the Fast Text Color plugin.
	st, err := NewFTCSettings(o.VaultPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		// The settings file exists but it can't be read for some other reason.
		return nil, fmt.Errorf("exporter: fast text color settings exists in vault but parsing failed: %w", err)
	}

	return &Exporter{
		options:       o,
		vault:         NewVault(o.VaultPath),
		ftcSettings:   st,
		documentCache: make(map[string]*schema.Document),
		navnodeCache:  make(map[string]*schema.NavigationNode),
	}, nil
}

// Returns the identity of the Outline user exporting the vault.
func (e *Exporter) Identity() *schema.Identity {
	return e.options.Identity
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
	L := e.options.Logger

	cid, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("failed to generate UUID: %w", err)
	}
	curlid := util.NewURLID()
	cname := filepath.Base(e.vault.Path())

	e.collection = schema.NewCollection(cid, curlid, cname)
	clear(e.documentCache)
	clear(e.navnodeCache)

	// Scan the vault for notes and attachements to export.
	if err := e.vault.Scan(e.options.Ignore); err != nil {
		return fmt.Errorf("export: failed to scan vault at %q: %w", e.vault.Path(), err)
	}

	notes := slices.Collect(
		collections.SeqFilter(
			e.vault.Files(),
			func(vf *VaultFile) bool { return vf.FileFormat == util.FileNote },
		),
	)
	lnotes := len(notes)

	L.Printf("Scanned vault at %q, found %d files and %d notes\n", e.vault.Path(), e.vault.Len(), lnotes)

	// Generate a document for each note in the vault.
	for i, vf := range notes {
		L.Printf("(%d/%d) Creating document %q\n", i+1, lnotes, vf.RelPath)

		if err := e.createParentDocuments(vf.RelPath); err != nil {
			return fmt.Errorf("exporter: failed to create parent documents for %q: %w", vf.AbsPath, err)
		}

		if err := e.createDocumentFromNote(vf); err != nil {
			return fmt.Errorf("exporter: failed to create document for %q: %w", vf.RelPath, err)
		}
	}

	zw := zip.NewWriter(w)

	L.Println("Writing metadata JSON to zip")
	if err := e.writeMetadata(zw); err != nil {
		return err
	}

	L.Println("Writing collection JSON to zip")
	if err := e.writeCollection(zw); err != nil {
		return err
	}

	i := 0
	la := len(e.collection.Attachments)
	for _, att := range e.collection.Attachments {
		r := util.Must(filepath.Rel(e.vault.Path(), att.AbsPath))

		L.Printf("(%d/%d) Copying attachment %q to zip\n", i+1, la, r)
		if err := e.copyAttachment(zw, att); err != nil {
			return err
		}

		i += 1
	}

	L.Println("Done!")

	return zw.Close()
}

// Writes the export metadata as JSON to the zipfile.
func (e *Exporter) writeMetadata(zw *zip.Writer) error {
	m := schema.NewExportMetadata(e.options.Identity)

	// Write the export metadata.
	mf, err := zw.Create(metadataFilename)
	if err != nil {
		return fmt.Errorf("exporter: couldn't create metadata %q in zip: %w", metadataFilename, err)
	}

	me := json.NewEncoder(mf)
	if err := me.Encode(m); err != nil {
		return fmt.Errorf("exporter: couldn't write to metadata %q in zip: %w", metadataFilename, err)
	}

	return nil
}

// Writes the exported collection as JSON to the zipfile.
func (e *Exporter) writeCollection(zw *zip.Writer) error {
	fname := e.collection.Metadata.Name + jsonExtension
	// Write the collection to a JSON file of the same name.
	cf, err := zw.Create(fname)
	if err != nil {
		return fmt.Errorf("exporter: couldn't create collection %q in zip: %w", fname, err)
	}

	ce := json.NewEncoder(cf)
	if err := ce.Encode(e.collection); err != nil {
		return fmt.Errorf("exporter: couldn't write to collection %q in zip: %w", fname, err)
	}

	return nil
}

func (e *Exporter) copyAttachment(zw *zip.Writer, att *schema.Attachment) error {
	f, err := os.Open(att.AbsPath)
	if err != nil {
		return fmt.Errorf("exporter: couldn't copy attachment %q to zip: %w", att.AbsPath, err)
	}

	// Copy the attachment's contents to the zip file.
	af, err := zw.Create(att.Key)
	if err != nil {
		return fmt.Errorf("exporter: couldn't create attachment %q in zip: %w", att.Key, err)
	}

	if _, err := io.Copy(af, f); err != nil {
		return fmt.Errorf("exporter: couldn't write to attachment %q in zip: %w", att.Key, err)
	}

	return nil
}

// Registers a document for export, where key is a relative path within the vault for caching the document.
// If pdoc and pnn are non-nil, they are registered as the parent document and navnode.
func (e *Exporter) registerDocument(
	key string,
	doc *schema.Document,
	nn *schema.NavigationNode,
	pdoc *schema.Document,
	pnn *schema.NavigationNode,
) {
	if pdoc != nil && pnn != nil {
		// Indicate the document to be a child of the parent document, and add its navigation node as a child of the parent's.
		doc.ParentDocumentId = &pdoc.ID
		pnn.Children = append(pnn.Children, nn)
	} else {
		meta := e.collection.Metadata
		// Add the document's navigation node to the root of the collection.
		meta.DocumentStructure = append(meta.DocumentStructure, nn)
	}

	e.documentCache[key] = doc
	e.navnodeCache[key] = nn

	e.collection.Documents[doc.ID] = doc
}

// Registers one or more files in a vault for export as an attachment belonging to a specific document.
func (e *Exporter) registerAttachments(files iter.Seq[*VaultFile], doc *schema.Document) {
	for vf := range files {
		if _, ok := e.collection.Attachments[vf.ID]; ok {
			// Attachment has already been registered.
			continue
		}

		a := schema.NewAttachment(doc.CreatedById, doc.ID, vf.MimeType.String(), vf.AbsPath, vf.ID, vf.Size)
		e.collection.Attachments[a.ID] = a
	}
}

// Creates a document by parsing an Obsidian note.
func (e *Exporter) createDocumentFromNote(note *VaultFile) error {
	if _, ok := e.documentCache[note.RelPath]; ok {
		// Document has already been exported? The vault is only iterated over once so this shouldn't be the case.
		return nil
	}

	doc := schema.NewDocument(note.ID, note.URLID, note.Title(), e.options.Identity)
	nn := schema.NewNavigationNode(doc)

	if err := doc.SetTimestamps(note.AbsPath); err != nil {
		return fmt.Errorf("failed to read timestamps on note at %q: %w", note.AbsPath, err)
	}

	src, err := os.ReadFile(note.AbsPath)
	if err != nil {
		return fmt.Errorf("failed to read note: %w", err)
	}

	res := NewNoteResolver(e, note)
	cv := converter.New(res)
	node, err := cv.Convert(src)
	if err != nil {
		return fmt.Errorf("failed to convert note: %w", err)
	}

	doc.Data = node

	// Get the document and navnode of the parent directory.
	dir := DirWithoutSlash(note.RelPath)
	pdoc := e.documentCache[dir]
	pnn := e.navnodeCache[dir]
	e.registerDocument(note.RelPath, doc, nn, pdoc, pnn)
	e.registerAttachments(res.Attachments().Items(), doc)

	return nil
}

// Creates a document for each parent directory in relpath.
func (e *Exporter) createParentDocuments(relpath string) error {
	L := e.options.Logger

	// This must be used for dir to be a vaild key.
	dir := DirWithoutSlash(relpath)
	if dir == "" {
		// The note does not have any parent directory.
		return nil
	}

	parents := strings.Split(dir, "/")

	// The last document and navigation node generated.
	var (
		lkey string
		ldoc *schema.Document
		lnn  *schema.NavigationNode
	)

	// OPT: iterate backwards from the parent directory of the note to avoid traversing the whole path.
	for i := range slices.Backward(parents) {
		key := path.Join(parents[:i+1]...)

		doc, ok := e.documentCache[key]
		nn := e.navnodeCache[key]

		if ok {
			// The parent directory has an existing document; therefore its ancestors have one too.
			break
		}

		L.Printf("Creating parent document %q\n", key)

		doc, nn, err := e.createEmptyDocument(key)
		if err != nil {
			return err
		}

		if ldoc != nil && lnn != nil {
			// Register the last document as a child of this one.
			e.registerDocument(lkey, ldoc, lnn, doc, nn)
		}

		lkey = key
		ldoc = doc
		lnn = nn
	}

	if ldoc != nil && lnn != nil {
		// Register the top-most document at the root of the collection.
		e.registerDocument(lkey, ldoc, lnn, nil, nil)
	}

	return nil
}

// Creates an empty document for the relative path.
func (e *Exporter) createEmptyDocument(relpath string) (*schema.Document, *schema.NavigationNode, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, nil, err
	}
	urlid := util.NewURLID()

	doc := schema.NewDocument(id, urlid, path.Base(relpath), e.options.Identity)

	abs := filepath.Join(e.vault.Path(), relpath)
	if err := doc.SetTimestamps(abs); err != nil {
		return nil, nil, fmt.Errorf("failed to read timestamps of %q: %w", abs, err)
	}

	nn := schema.NewNavigationNode(doc)

	return doc, nn, nil
}
