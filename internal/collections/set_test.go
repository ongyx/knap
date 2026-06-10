package collections

import "testing"

func TestSet(t *testing.T) {
	s := NewSet(1, 2, 3)

	if !s.Contains(1) {
		t.Errorf("expected set to contain 1")
	}
	if !s.Contains(2) {
		t.Errorf("expected set to contain 2")
	}
	if !s.Contains(3) {
		t.Errorf("expected set to contain 3")
	}
	if s.Contains(4) {
		t.Errorf("expected set to not contain 4")
	}

	s.Add(4)
	if !s.Contains(4) {
		t.Errorf("expected set to contain 4 after Add")
	}

	s.Remove(1)
	if s.Contains(1) {
		t.Errorf("expected set to not contain 1 after Remove")
	}
}
