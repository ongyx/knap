package exporter

import (
	"errors"
	"os"
	"regexp"

	"github.com/ongyx/knap/internal/util"
)

// I have no clue why Outline uses two different slug libraries...
// https://github.com/Trott/slug/blob/003bd7b4b86456ea6215c01cb5954d1fca37bec1/slug.js#L91
var documentURLSlugifyOptions = &util.SlugifyOptions{
	Remove: regexp.MustCompile(`[^\w\s\-~]`),
}

// Checks if a path is a directory and exists in the filesystem.
func FolderExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}
