package conditions

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/gitsang/dnd5e-dm-skill/cli/internal/combat"
)

func TestAddAndRemoveCondition(t *testing.T) {
	path := filepath.Join(t.TempDir(), "combat_state.json")
	state := combat.State{Combatants: []combat.Combatant{{ID: "pc1", Name: "Aric", Conditions: []combat.Condition{}}}}
	encoded, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, encoded, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Add(path, "pc1", "poisoned", "", ""); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	loaded, err := combat.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Combatants[0].Conditions) != 1 || loaded.Combatants[0].Conditions[0].Name != "poisoned" {
		t.Fatalf("conditions after add = %#v", loaded.Combatants[0].Conditions)
	}
	if err := Remove(path, "pc1", "poisoned"); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	loaded, err = combat.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Combatants[0].Conditions) != 0 {
		t.Fatalf("conditions after remove = %#v", loaded.Combatants[0].Conditions)
	}
}
