package exporter

import (
	"html"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ongyx/knap/internal/converter"
	"github.com/ongyx/knap/internal/obsidian"
	"github.com/ongyx/knap/internal/schema"
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
	exporter Exporter
	note     *VaultFile
}

// Implements internal/converter.ResolveInternalLink.
func (nr *NoteResolver) ResolveInternalLink(link *converter.Link) (node *schema.Node, err error) {
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
		node, err = nr.handleVideoFile(vf, link)
	} else {
		// Any other file is exported as an attachment without special presentation.
		node, err = nr.handleAttachment(vf, link)
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

func (nr *NoteResolver) handleAttachment(vf *VaultFile, link *converter.Link) (*schema.Node, error) {
	// Probe the file for its size and MIME content type.
	f, err := os.Open(vf.AbsPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var title string
	if link.Text != nil {
		title = string(link.Text)
	} else {
		title = filepath.Base(vf.AbsPath)
	}

	// DetectContentType only takes the first 512 bytes.
	var buf [512]byte
	if _, err := f.Read(buf[:]); err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(buf[:])

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// Generate an attachment node.
	return schema.NewAttachmentNode(vf.ID, title, contentType, fi.Size()), nil
}

func (nr *NoteResolver) handleNoteFile(vf *VaultFile, link *converter.Link) (*schema.Node, error) {
	// References to other notes are converted into links to the document URL.
	href := vf.DocumentURL()

	if link.URL.Fragment != "" {
		// Add the heading to the URL as a fragment.
		href += slugifyHeading(link.URL.Fragment)
	}

	var title string
	if link.Text != nil {
		title = string(link.Text)
	} else {
		title = vf.Title()
	}

	node := schema.NewTextNode(title)
	node.Marks = append(node.Marks, schema.NewLinkMark(href))
	return node, nil
}

func (nr *NoteResolver) handleImageFile(vf *VaultFile, link *converter.Link) (*schema.Node, error) {
	w, h := converter.ParseEmbedSize(link.Text)
	return schema.NewImageFileNode(vf.ID, w, h), nil
}

func (nr *NoteResolver) handleVideoFile(vf *VaultFile, link *converter.Link) (*schema.Node, error) {
	// Obsidian does not parse embed sizes for videos, so we do the same here.
	return schema.NewVideoFileNode(vf.ID, vf.Title()), nil
}

// Slugifies a heading into a DOM ID.
//
// See https://github.com/outline/outline/blob/39623b90bd0846ac2316395143a73a3340e2dfd3/shared/editor/lib/headingToSlug.ts#L10.
func slugifyHeading(str string) string {
	slug := util.Slugify(str, slugifyHeadingOptions)
	return "h-" + html.EscapeString(slug)
}
