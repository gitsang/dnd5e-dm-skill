package combat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeState(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "combat_state.json")
	state := State{
		Round:       1,
		ActiveIndex: 0,
		Combatants: []Combatant{
			{ID: "pc1", Name: "Aric", HP: 12, MaxHP: 12, Used: UsedActions{Action: false, BonusAction: false, Reaction: false}},
			{ID: "gob1", Name: "Goblin", HP: 7, MaxHP: 7, Used: UsedActions{Action: false, BonusAction: false, Reaction: false}},
		},
	}
	encoded, err := json.Marshal(state)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, encoded, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestUseActionPreventsDuplicateAction(t *testing.T) {
	path := writeState(t)
	if err := Use(path, "pc1", "action"); err != nil {
		t.Fatalf("first use failed: %v", err)
	}
	if err := Use(path, "pc1", "action"); err == nil {
		t.Fatal("duplicate action succeeded")
	}
}

func TestApplyDamageUpdatesHP(t *testing.T) {
	path := writeState(t)
	if err := Damage(path, "gob1", 3); err != nil {
		t.Fatalf("damage failed: %v", err)
	}
	state, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if state.Combatants[1].HP != 4 {
		t.Fatalf("goblin hp = %d", state.Combatants[1].HP)
	}
}
