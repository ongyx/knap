package exporter

import (
	"html"
	"path/filepath"
	"regexp"

	"github.com/ongyx/knap/internal/collections"
	"github.com/ongyx/knap/internal/converter"
	"github.com/ongyx/knap/internal/obsidian"
	"github.com/ongyx/knap/internal/prosemirror"
	"github.com/ongyx/knap/internal/util"
	"github.com/yuin/goldmark/ast"
)

var _ converter.Resolver = (*NoteResolver)(nil)

var slugifyHeadingOptions = &util.SlugifyOptions{
	Remove: regexp.MustCompile(`[!"#$%&'.()*+,/:;<=>?@[\]\\^_` + "`" + `{|}~]`),
	Lower:  true,
}

// Resolver for a note being exported.
type NoteResolver struct {
	exporter    *Exporter
	note        *VaultFile
	attachments collections.Set[*VaultFile]
}

// Creates a new note resolver with the given exporter and the note vault file.
func NewNoteResolver(exporter *Exporter, note *VaultFile) *NoteResolver {
	return &NoteResolver{
		exporter:    exporter,
		note:        note,
		attachments: collections.NewSet[*VaultFile](),
	}
}

// Returns the attachments linked to by the note. This should only be called after the note has been fully converted.
func (nr *NoteResolver) Attachments() collections.Set[*VaultFile] {
	return nr.attachments
}

// Implements internal/converter.ResolveInternalLink.
func (nr *NoteResolver) ResolveInternalLink(link *converter.Link) (node *prosemirror.Node, err error) {
	v := nr.exporter.Vault()

	var vf *VaultFile
	if link.URL.Path != "" {
		// The internal link refers to another note or attachment. Lookup by its path.
		vf = v.Lookup(string(link.URL.Path))
		if vf == nil {
			// The internal link is invalid, return the default representation.
			return converter.DefaultResolver.ResolveInternalLink(link)
		}
	} else {
		// The internal link refers to this note.
		vf = nr.note
	}

	if vf.FileFormat == util.FileNote {
		node, err = nr.handleNoteFile(vf, link)
	} else if vf.FileFormat == util.FileImage && link.Embed {
		node, err = nr.handleImageFile(vf, link)
	} else if vf.FileFormat == util.FileVideo && link.Embed {
		node, err = nr.handleVideoFile(vf)
	} else {
		// Any other file is exported as an attachment node without special presentation.
		node, err = nr.handleAttachment(vf, link)
	}

	if vf.FileFormat != util.FileNote && err == nil {
		// Add the attachment to the set.
		nr.attachments.Add(vf)
	}

	return node, err
}

// Implements internal/converter.ResolveColor.
func (nr *NoteResolver) ResolveColor(doc *ast.Document, tc *obsidian.TextColor) string {
	st := nr.exporter.FTCSettings()
	if st == nil {
		return ""
	}

	theme := st.GetDefaultTheme()
	// Check if the document specified a theme to use in the frontmatter.
	// https://github.com/Superschnizel/obsidian-fast-text-color/blob/23f7e97a4ffac7643506a8e422dca0f016331529/src/rendering/TextColorViewPlugin.ts#L205
	fm := doc.Meta()
	if name, ok := fm["ftcTheme"].(string); ok {
		if t, ok := st.ThemeMap[name]; ok {
			theme = t
		}
	}

	// No theme is available.
	if theme == nil {
		return ""
	}

	// Lookup the color by its ID.
	if c, ok := theme.ColorMap[string(tc.ID)]; ok {
		return c.Color
	}

	return ""
}

func (nr *NoteResolver) handleAttachment(vf *VaultFile, link *converter.Link) (*prosemirror.Node, error) {
	title := link.Text
	if title == "" {
		title = filepath.Base(vf.AbsPath)
	}

	// Generate an attachment node.
	return prosemirror.NewAttachmentNode(vf.ID, title, vf.MimeType.String(), vf.Size), nil
}

func (nr *NoteResolver) handleNoteFile(vf *VaultFile, link *converter.Link) (*prosemirror.Node, error) {
	// References to other notes are converted into links to the document URL.
	href := vf.URLID.GenerateDocumentURL(vf.Title())

	if link.URL.Fragment != "" {
		// Add the heading to the URL as a fragment.
		href += slugifyHeading(link.URL.Fragment)
	}

	title := link.Text
	if title == "" {
		title = vf.Title()
	}

	node := prosemirror.NewTextNode(title)
	node.Marks = append(node.Marks, prosemirror.NewLinkMark(href))
	return node, nil
}

func (nr *NoteResolver) handleImageFile(vf *VaultFile, link *converter.Link) (*prosemirror.Node, error) {
	w, h, _ := converter.ParseEmbedSize(link.Text)
	return prosemirror.NewImageFileNode(vf.ID, w, h), nil
}

func (nr *NoteResolver) handleVideoFile(vf *VaultFile) (*prosemirror.Node, error) {
	// Obsidian does not parse embed sizes for videos, so we do the same here.
	return prosemirror.NewVideoFileNode(vf.ID, vf.Title()), nil
}

// Slugifies a heading into a DOM ID.
//
// See https://github.com/outline/outline/blob/39623b90bd0846ac2316395143a73a3340e2dfd3/shared/editor/lib/headingToSlug.ts#L10.
func slugifyHeading(str string) string {
	slug := util.Slugify(str, slugifyHeadingOptions)
	return "h-" + html.EscapeString(slug)
}
