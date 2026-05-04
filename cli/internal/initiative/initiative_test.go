package initiative

import (
	"math/rand"
	"testing"
)

func TestCreateSortsInitiativeDescending(t *testing.T) {
	state, err := Create([]InputCombatant{{ID: "slow", Name: "Slow", InitiativeBonus: 0, HP: 5, MaxHP: 5}, {ID: "fast", Name: "Fast", InitiativeBonus: 10, HP: 5, MaxHP: 5}}, rand.New(rand.NewSource(1)))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if state.Round != 1 || state.ActiveIndex != 0 {
		t.Fatalf("unexpected round/index: %#v", state)
	}
	if state.Combatants[0].Initiative < state.Combatants[1].Initiative {
		t.Fatalf("initiative not sorted: %#v", state.Combatants)
	}
}
