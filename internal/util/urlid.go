package util

import (
	"crypto/rand"
	"fmt"
	"regexp"
)

const urlIDCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const urlIDLength = 10

// I have no clue why Outline uses two different slug libraries...
// https://github.com/Trott/slug/blob/003bd7b4b86456ea6215c01cb5954d1fca37bec1/slug.js#L91
var documentURLSlugifyOptions = &SlugifyOptions{
	Remove: regexp.MustCompile(`[^\w\s\-~]`),
}

// A short randomized ID for use with document URLs.
type URLID string

// Generates a new URLID.
//
// Source: https://github.com/outline/outline/blob/5ea63aa1a28a0a55cd2c8311caa53705e63d1d4e/shared/random.ts#L35
func NewURLID() URLID {
	var out [urlIDLength]byte
	rand.Read(out[:])

	for i, n := range out {
		// Replace the random number with a byte from the charset. Modulo is necessary to prevent out-of-bounds.
		out[i] = urlIDCharset[n%urlIDLength]
	}

	return URLID(out[:])
}

// Generates an Outline document URL, given the title of the document.
func (u URLID) GenerateDocumentURL(title string) string {
	// https://github.com/outline/outline/blob/88de417a21c260e32ecc9c89d756661c54064603/server/models/Document.ts#L430
	slug := Slugify(title, documentURLSlugifyOptions)

	if len(slug) > 0 {
		return fmt.Sprintf(`/doc/%s-%s`, slug, u)
	} else {
		return fmt.Sprintf(`/doc/untitled-%s`, u)
	}
}
