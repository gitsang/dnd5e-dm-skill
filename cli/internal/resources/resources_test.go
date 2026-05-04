package resources

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSpendResourceAndRejectOverspend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "aric.json")
	initial := map[string]any{"spell_slots": map[string]any{"level_1": float64(1)}}
	encoded, err := json.Marshal(initial)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, encoded, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Spend(path, "spell_slots.level_1", 1); err != nil {
		t.Fatalf("spend failed: %v", err)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var saved map[string]map[string]float64
	if err := json.Unmarshal(contents, &saved); err != nil {
		t.Fatal(err)
	}
	if saved["spell_slots"]["level_1"] != 0 {
		t.Fatalf("slot count = %v", saved["spell_slots"]["level_1"])
	}
	if err := Spend(path, "spell_slots.level_1", 1); err == nil {
		t.Fatal("overspend succeeded")
	}
}
