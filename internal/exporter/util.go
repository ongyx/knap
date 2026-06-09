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

// A hashset of values of comparable type T.
type Set[T comparable] struct {
	inner map[T]struct{}
}

// Creates a new set with the given items.
func NewSet[T comparable](items ...T) Set[T] {
	s := Set[T]{make(map[T]struct{}, len(items))}
	for _, i := range items {
		s.Add(i)
	}
	return s
}

// Adds an item to the set.
func (s *Set[T]) Add(item T) {
	s.inner[item] = struct{}{}
}

// Removes an item from the set.
func (s *Set[T]) Remove(item T) {
	delete(s.inner, item)
}

// Checks if the set contains an item.
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.inner[item]
	return ok
}
