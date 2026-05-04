package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchFindsLocalRuleFiles(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "grapple.md"), []byte("# Grapple\nA grapple uses a special melee attack."), 0o644); err != nil {
		t.Fatal(err)
	}
	results, err := Search([]string{dir}, "grapple")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(results) != 1 || results[0].Title != "grapple.md" {
		t.Fatalf("unexpected results: %#v", results)
	}
}
