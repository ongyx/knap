package exporter

import (
	"fmt"
	"html"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/ongyx/knap/internal/converter"
	"github.com/ongyx/knap/internal/obsidian"
	"github.com/ongyx/knap/internal/schema"
	"github.com/yuin/goldmark/ast"
)

var _ converter.Resolver = (*NoteResolver)(nil)

// Matches an image size in the format (w)x(h), where h is optional.
var reImageSize = regexp.MustCompile(`^(\d+)(?:x(\d+))?$`)

var slugifyHeadingOptions = (&SlugifyOptions{
	Remove: regexp.MustCompile(`[!"#$%&'.()*+,/:;<=>?@[\]\\^_` + "`" + `{|}~]`),
	Lower:  true,
}).Defaults()

// Resolver for a note being exported.
type NoteResolver struct {
	exporter  Exporter
	vaultFile *VaultFile
}

// Implements internal/converter.ResolveInternalLink.
func (nr *NoteResolver) ResolveInternalLink(il converter.InternalLink) (*schema.Node, error) {
	v := nr.exporter.Vault()

	var vf *VaultFile
	if il.Target != nil {
		// The internal link refers to another note/attachment. Lookup the vault file from the link text.
		vf = v.Lookup(string(il.Target))
		if vf == nil {
			// The internal link is invalid.
			return converter.DefaultResolver.ResolveInternalLink(il)
		}
	} else {
		// The internal link refers to this note.
		vf = nr.vaultFile
	}

	var (
		node *schema.Node
		err  error
	)
	if vf.FileType == VaultFileNote {
		node, err = nr.handleNote(vf, il)
	} else if vf.FileType == VaultFileImage && il.Embed {
		node, err = nr.handleImage(vf, il)
	} else if vf.FileType == VaultFileVideo && il.Embed {
		node, err = nr.handleVideo(vf, il)
	} else {
		// Any other file is exported as an attachment without special presentation.
		node, err = nr.handleAttachment(vf, il)
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

func (nr *NoteResolver) handleAttachment(vf *VaultFile, il converter.InternalLink) (*schema.Node, error) {
	// Probe the file for its size and MIME content type.
	f, err := os.Open(vf.AbsPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var title string
	if il.Title != nil {
		title = string(il.Title)
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

func (nr *NoteResolver) handleNote(vf *VaultFile, il converter.InternalLink) (*schema.Node, error) {
	// References to other notes are converted into links to the document URL.
	href := vf.DocumentURL()

	if il.Fragment != nil {
		// Add the heading to the URL as a fragment.
		href += slugifyHeading(string(il.Fragment))
	}

	var title string
	if il.Title != nil {
		title = string(il.Title)
	} else {
		title = vf.Title()
	}

	node := schema.NewTextNode(title)
	node.Marks = append(node.Marks, schema.NewLinkMark(href))
	return node, nil
}

func (nr *NoteResolver) handleImage(vf *VaultFile, il converter.InternalLink) (*schema.Node, error) {
	w, h, err := parseEmbedSize(il)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embed size: %w", err)
	}

	return schema.NewImageFileNode(vf.ID, w, h), nil
}

func (nr *NoteResolver) handleVideo(vf *VaultFile, il converter.InternalLink) (*schema.Node, error) {
	// Obsidian doesn't support changing the presentation size of an embedded video, so sizes are ignored.
	return schema.NewVideoNode(vf.ID, vf.Title()), nil
}

// Slugifies a heading into a DOM ID.
//
// See https://github.com/outline/outline/blob/39623b90bd0846ac2316395143a73a3340e2dfd3/shared/editor/lib/headingToSlug.ts#L10.
func slugifyHeading(str string) string {
	slug := Slugify(str, slugifyHeadingOptions)
	return "h-" + html.EscapeString(slug)
}

func parseEmbedSize(il converter.InternalLink) (width, height int, err error) {
	// Try to take the image dimensions from the title.
	m := reImageSize.FindSubmatchIndex(il.Title)
	if len(m) >= 4 {
		// width
		w := il.Title[m[2]:m[3]]
		width, err = strconv.Atoi(string(w))
		if err != nil {
			return 0, 0, err
		}
	}
	if len(m) == 6 {
		// height
		h := il.Title[m[4]:m[5]]
		height, err = strconv.Atoi(string(h))
		if err != nil {
			return 0, 0, err
		}
	}

	return width, height, nil
}
