package collections

// Stack is a first-in-first-out (FIFO) stack.
type Stack[T any] struct {
	inner []T
}

// Creates a new stack with the given size and capacity. This is equivalent to `make([]T, size, cap)`.
func NewStack[T any](size int, cap int) Stack[T] {
	return Stack[T]{inner: make([]T, size, cap)}
}

// Pushes an item onto the stack.
func (s *Stack[T]) Push(item T) {
	s.inner = append(s.inner, item)
}

// Pops an item from the stack. If the stack is empty, ok is false.
func (s *Stack[T]) Pop() (item T, ok bool) {
	l := s.Len()
	if l > 0 {
		item = s.inner[l-1]
		s.inner = s.inner[:l-1]
		ok = true
	}

	return
}

// Peeks the top-most item on the stack. If the stack is empty, ok is false.
func (s *Stack[T]) Peek() (item T, ok bool) {
	l := s.Len()
	if l > 0 {
		item = s.inner[l-1]
		ok = true
	}

	return
}

// Returns the number of items in the stack.
func (s *Stack[T]) Len() int {
	return len(s.inner)
}

// Clears the items in the stack. This does not affect the stack's capacity.
func (s *Stack[T]) Clear() {
	s.inner = s.inner[:0]
}
