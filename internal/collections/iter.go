package collections

import "iter"

// Returns a sequence that yields elements where predicate returns true.
func SeqFilter[V any](seq iter.Seq[V], predicate func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if !predicate(v) {
				continue
			}

			if !yield(v) {
				break
			}
		}
	}
}

// Returns a sequence that yields the key from a pair iterator.
func SeqKeys[K, V any](seq iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k, _ := range seq {
			if !yield(k) {
				break
			}
		}
	}
}

// Returns a sequence that yields the values from a pair iterator.
func SeqValues[K, V any](seq iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range seq {
			if !yield(v) {
				break
			}
		}
	}
}
