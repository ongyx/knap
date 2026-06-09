package exporter

import (
	"regexp"
	"strings"
)

const defaultReplacement = "-"

var (
	// Default regex for removing invalid characters.
	reDefaultRemove = regexp.MustCompile(`[^\w\s$*_+~.()'"!\-:@]+`)
	reSpace         = regexp.MustCompile(`\s+`)
)

// Options for slugifying text.
type SlugifyOptions struct {
	// The character to replace spaces with.
	Replacement string
	// The regex to match invalid characters against.
	Remove *regexp.Regexp
	// Whether or not to convert the slug to lowercase.
	Lower bool
}

// Applies defaults to zero values in the options.
func (o *SlugifyOptions) Defaults() *SlugifyOptions {
	if o.Replacement == "" {
		o.Replacement = defaultReplacement
	}

	if o.Remove == nil {
		o.Remove = reDefaultRemove
	}

	return o
}

// Slugifies a string.
//
// This is ported from the slugify library: see https://github.com/simov/slugify/blob/master/slugify.js for the original source code.
func Slugify(str string, options *SlugifyOptions) string {
	if options == nil {
		options = (&SlugifyOptions{}).Defaults()
	}

	slug := options.Remove.ReplaceAllLiteralString(str, "")
	slug = reSpace.ReplaceAllLiteralString(slug, options.Replacement)

	if options.Lower {
		slug = strings.ToLower(slug)
	}

	return slug
}
