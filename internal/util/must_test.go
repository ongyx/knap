package util

import (
	"errors"
	"testing"
)

func TestMust(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Must panicked unexpectedly: %v", r)
		}
	}()

	res := Must(1, nil)
	if res != 1 {
		t.Errorf("expected 1, got %d", res)
	}
}

func TestMustPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Must should have panicked")
		}
	}()

	Must(1, errors.New("test error"))
}
