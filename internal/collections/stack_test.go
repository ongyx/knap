package collections

import "testing"

func TestStack(t *testing.T) {
	s := NewStack[int](0, 5)

	if s.Len() != 0 {
		t.Errorf("expected length 0, got %d", s.Len())
	}

	s.Push(1)
	s.Push(2)
	s.Push(3)

	if s.Len() != 3 {
		t.Errorf("expected length 3, got %d", s.Len())
	}

	if item, ok := s.Peek(); !ok || item != 3 {
		t.Errorf("expected peek 3, got %v, %v", item, ok)
	}

	if item, ok := s.Pop(); !ok || item != 3 {
		t.Errorf("expected pop 3, got %v, %v", item, ok)
	}

	if s.Len() != 2 {
		t.Errorf("expected length 2, got %d", s.Len())
	}

	s.Clear()
	if s.Len() != 0 {
		t.Errorf("expected length 0 after Clear, got %d", s.Len())
	}

	if _, ok := s.Pop(); ok {
		t.Errorf("expected ok=false when popping from empty stack")
	}

	if _, ok := s.Peek(); ok {
		t.Errorf("expected ok=false when peeking empty stack")
	}
}
