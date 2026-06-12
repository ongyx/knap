package schema

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

var testCollectionPath = filepath.FromSlash("testdata/test_collection/Test Collection.json")

func TestCollectionUnmarshal(t *testing.T) {
	b, err := os.ReadFile(testCollectionPath)
	if err != nil {
		t.Fatalf("couldn't read collection file: %s", err)
	}

	var c Collection
	if err := json.Unmarshal(b, &c); err != nil {
		t.Errorf("couldn't unmarshal collection: %s", err)
	}
}
