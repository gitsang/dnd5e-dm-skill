package combat

import (
	"encoding/json"
	"fmt"
	"os"
)

type State struct {
	Round       int         `json:"round"`
	ActiveIndex int         `json:"active_index"`
	Combatants  []Combatant `json:"combatants"`
}

type Combatant struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Initiative int            `json:"initiative,omitempty"`
	AC         int            `json:"ac,omitempty"`
	HP         int            `json:"hp"`
	MaxHP      int            `json:"max_hp"`
	Conditions []Condition    `json:"conditions,omitempty"`
	Used       UsedActions    `json:"used"`
	Extra      map[string]any `json:"-"`
}

type Condition struct {
	Name     string `json:"name"`
	Source   string `json:"source,omitempty"`
	Duration string `json:"duration,omitempty"`
}

type UsedActions struct {
	Action            bool `json:"action"`
	BonusAction       bool `json:"bonus_action"`
	Reaction          bool `json:"reaction"`
	Movement          int  `json:"movement"`
	ObjectInteraction bool `json:"object_interaction"`
}

func Load(path string) (State, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return State{}, err
	}
	var state State
	if err := json.Unmarshal(contents, &state); err != nil {
		return State{}, err
	}
	return state, nil
}

func Save(path string, state State) error {
	encoded, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, encoded, 0o644)
}

func Use(path string, combatantID string, kind string) error {
	state, err := Load(path)
	if err != nil {
		return err
	}
	combatant, err := findCombatant(&state, combatantID)
	if err != nil {
		return err
	}
	if actionUsed(combatant.Used, kind) {
		return fmt.Errorf("%s already used %s", combatant.Name, kind)
	}
	if err := markUsed(&combatant.Used, kind); err != nil {
		return err
	}
	return Save(path, state)
}

func Damage(path string, combatantID string, amount int) error {
	if amount < 0 {
		return fmt.Errorf("amount must not be negative")
	}
	state, err := Load(path)
	if err != nil {
		return err
	}
	combatant, err := findCombatant(&state, combatantID)
	if err != nil {
		return err
	}
	combatant.HP -= amount
	if combatant.HP < 0 {
		combatant.HP = 0
	}
	return Save(path, state)
}

func NextTurn(path string) error {
	state, err := Load(path)
	if err != nil {
		return err
	}
	if len(state.Combatants) == 0 {
		return fmt.Errorf("combat state has no combatants")
	}
	state.ActiveIndex = (state.ActiveIndex + 1) % len(state.Combatants)
	if state.ActiveIndex == 0 {
		state.Round++
	}
	state.Combatants[state.ActiveIndex].Used = UsedActions{}
	return Save(path, state)
}

func findCombatant(state *State, combatantID string) (*Combatant, error) {
	for index := range state.Combatants {
		if state.Combatants[index].ID == combatantID {
			return &state.Combatants[index], nil
		}
	}
	return nil, fmt.Errorf("unknown combatant: %s", combatantID)
}

func actionUsed(used UsedActions, kind string) bool {
	switch kind {
	case "action":
		return used.Action
	case "bonus_action":
		return used.BonusAction
	case "reaction":
		return used.Reaction
	case "object_interaction":
		return used.ObjectInteraction
	default:
		return false
	}
}

func markUsed(used *UsedActions, kind string) error {
	switch kind {
	case "action":
		used.Action = true
	case "bonus_action":
		used.BonusAction = true
	case "reaction":
		used.Reaction = true
	case "object_interaction":
		used.ObjectInteraction = true
	default:
		return fmt.Errorf("unsupported action kind: %s", kind)
	}
	return nil
}
