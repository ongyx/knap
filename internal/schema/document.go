package schema

import (
	"io"
	"os"
	"regexp"
	"time"

	"github.com/djherbis/times"
	"github.com/google/uuid"
)

// Maps callouts to notice types.
// https://obsidian.md/help/callouts
var calloutToNotice = map[string]NoticeType{
	"note":      NoticeInfo,
	"abstract":  NoticeInfo,
	"summary":   NoticeInfo,
	"tldr":      NoticeInfo,
	"info":      NoticeInfo,
	"todo":      NoticeInfo,
	"tip":       NoticeTip,
	"hint":      NoticeTip,
	"important": NoticeTip,
	"success":   NoticeSuccess,
	"check":     NoticeSuccess,
	"done":      NoticeSuccess,
	"question":  NoticeInfo,
	"help":      NoticeInfo,
	"faq":       NoticeInfo,
	"warning":   NoticeWarning,
	"caution":   NoticeWarning,
	"attention": NoticeWarning,
	"failure":   NoticeWarning,
	"fail":      NoticeWarning,
	"missing":   NoticeWarning,
	"danger":    NoticeWarning,
	"error":     NoticeWarning,
	"bug":       NoticeWarning,
	"example":   NoticeInfo,
	"quote":     NoticeInfo,
	"cite":      NoticeInfo,
}

// Matches a callout.
var reCallout = regexp.MustCompile(`^\[\!(\w+)\]`)

// Matches a language class.
var reLangClass = regexp.MustCompile(`language-(\w+)`)

// Document represents an individual Outline document, including its content and metadata.
type Document struct {
	ID               uuid.UUID  `json:"id"`
	URLID            URLID      `json:"urlId"`
	Title            string     `json:"title"`
	Icon             *string    `json:"icon"`
	Color            *string    `json:"color"`
	Data             *Node      `json:"data"`
	CreatedById      string     `json:"createdById"`
	CreatedByName    string     `json:"createdByName"`
	CreatedByEmail   string     `json:"createdByEmail"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	PublishedAt      *time.Time `json:"publishedAt"`
	FullWidth        bool       `json:"fullWidth"`
	Template         bool       `json:"template"`
	ParentDocumentId *string    `json:"parentDocumentId"`
}

// Creates an empty document.
func NewDocument() *Document {
	id, _ := uuid.NewRandom()
	urlid := NewURLID()

	now := time.Now()
	return &Document{
		ID:          id,
		URLID:       urlid,
		CreatedAt:   now,
		UpdatedAt:   now,
		PublishedAt: &now,
	}
}

// Parses markdown from a reader into this document.
// If the reader is an os.File, timestamps for the *At fields are read.
func (d *Document) ParseReader(r io.Reader) error {
	if f, ok := r.(*os.File); ok {
		// The reader is a file, query the file's timestamps to fill in the *At values.
		// Getting creation time is surprisingly hard...
		t, err := times.Stat(f.Name())
		if err != nil {
			return err
		}

		mtime := t.ModTime()
		if t.HasBirthTime() {
			d.CreatedAt = t.BirthTime()
		} else {
			d.CreatedAt = mtime
		}
		d.UpdatedAt = mtime
		d.PublishedAt = &mtime
	}

	src, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return d.parse(src)
}

// Parses markdown text into the document tree. This clears the existing content, if any.
func (d *Document) parse(src []byte) error {
	cv := NewConverter()
	root, err := cv.Convert(src)
	if err != nil {
		return err
	}

	d.Data = root

	return nil
}
