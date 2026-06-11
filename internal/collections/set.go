package collections

import (
	"iter"
	"maps"
)

// Set is a hashset.
type Set[T comparable] struct {
	inner map[T]struct{}
}

// Creates a new set with the given items.
func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{make(map[T]struct{}, len(items))}
	for _, i := range items {
		s.Add(i)
	}
	return s
}

// Checks if the set contains an item.
func (s *Set[T]) Contains(item T) bool {
	_, ok := s.inner[item]
	return ok
}

// Returns the number of items in the set.
func (s *Set[T]) Len() int {
	return len(s.inner)
}

// Returns an iterator over the items in the set. The items are not guaranteed to be yielded in any order.
func (s *Set[T]) Items() iter.Seq[T] {
	return maps.Keys(s.inner)
}

// Adds an item to the set.
func (s *Set[T]) Add(item T) {
	s.inner[item] = struct{}{}
}

// Removes an item from the set.
func (s *Set[T]) Remove(item T) {
	delete(s.inner, item)
}

// Updates this set with items from another set.
func (s *Set[T]) Update(ss *Set[T]) {
	for item := range ss.Items() {
		s.Add(item)
	}
}

// Clears the set.
func (s *Set[T]) Clear() {
	clear(s.inner)
}
