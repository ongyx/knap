package schema

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

var testMetadataPath = filepath.FromSlash("testdata/test_collection/metadata.json")

func TestExportMetadataUnmarshal(t *testing.T) {
	b, err := os.ReadFile(testCollectionPath)
	if err != nil {
		t.Fatalf("couldn't read collection file: %s", err)
	}

	var em ExportMetadata
	if err := json.Unmarshal(b, &em); err != nil {
		t.Errorf("couldn't unmarshal export metadata: %s", err)
	}
}
