package initiative

import (
	"math/rand"
	"sort"
	"strconv"

	"github.com/gitsang/dnd5e-dm-skill/cli/internal/combat"
	"github.com/gitsang/dnd5e-dm-skill/cli/internal/dice"
)

type InputCombatant struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	InitiativeBonus int    `json:"initiative_bonus"`
	AC              int    `json:"ac,omitempty"`
	HP              int    `json:"hp"`
	MaxHP           int    `json:"max_hp"`
}

func Create(inputs []InputCombatant, rng *rand.Rand) (combat.State, error) {
	state := combat.State{Round: 1, ActiveIndex: 0, Combatants: make([]combat.Combatant, 0, len(inputs))}
	for _, input := range inputs {
		roll, err := dice.RollExpression(formatInitiativeExpression(input.InitiativeBonus), rng)
		if err != nil {
			return combat.State{}, err
		}
		state.Combatants = append(state.Combatants, combat.Combatant{ID: input.ID, Name: input.Name, Initiative: roll.Total, AC: input.AC, HP: input.HP, MaxHP: input.MaxHP, Conditions: []combat.Condition{}, Used: combat.UsedActions{}})
	}
	sort.SliceStable(state.Combatants, func(i, j int) bool {
		return state.Combatants[i].Initiative > state.Combatants[j].Initiative
	})
	return state, nil
}

func formatInitiativeExpression(bonus int) string {
	if bonus >= 0 {
		return "1d20+" + strconv.Itoa(bonus)
	}
	return "1d20" + strconv.Itoa(bonus)
}
