package exporter

import (
	"errors"
	"os"
)

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
