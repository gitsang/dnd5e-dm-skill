package conditions

import (
	"fmt"

	"github.com/gitsang/dnd5e-dm-skill/cli/internal/combat"
)

func Add(statePath string, combatantID string, name string, source string, duration string) error {
	if name == "" {
		return fmt.Errorf("condition is required")
	}
	state, err := combat.Load(statePath)
	if err != nil {
		return err
	}
	index, err := findCombatant(state, combatantID)
	if err != nil {
		return err
	}
	for _, condition := range state.Combatants[index].Conditions {
		if condition.Name == name {
			return combat.Save(statePath, state)
		}
	}
	state.Combatants[index].Conditions = append(state.Combatants[index].Conditions, combat.Condition{Name: name, Source: source, Duration: duration})
	return combat.Save(statePath, state)
}

func Remove(statePath string, combatantID string, name string) error {
	state, err := combat.Load(statePath)
	if err != nil {
		return err
	}
	index, err := findCombatant(state, combatantID)
	if err != nil {
		return err
	}
	kept := state.Combatants[index].Conditions[:0]
	for _, condition := range state.Combatants[index].Conditions {
		if condition.Name != name {
			kept = append(kept, condition)
		}
	}
	state.Combatants[index].Conditions = kept
	return combat.Save(statePath, state)
}

func List(statePath string, combatantID string) ([]combat.Condition, error) {
	state, err := combat.Load(statePath)
	if err != nil {
		return nil, err
	}
	index, err := findCombatant(state, combatantID)
	if err != nil {
		return nil, err
	}
	return state.Combatants[index].Conditions, nil
}

func findCombatant(state combat.State, combatantID string) (int, error) {
	for index, candidate := range state.Combatants {
		if candidate.ID == combatantID {
			return index, nil
		}
	}
	return 0, fmt.Errorf("unknown combatant: %s", combatantID)
}
