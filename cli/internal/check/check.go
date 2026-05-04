package check

import (
	"math/rand"

	"github.com/gitsang/dnd5e-dm-skill/cli/internal/dice"
)

type Result struct {
	Reason     string    `json:"reason"`
	Expression string    `json:"expression"`
	Rolls      []int     `json:"rolls"`
	Modifier   int       `json:"modifier"`
	Total      int       `json:"total"`
	DC         int       `json:"dc"`
	Success    bool      `json:"success"`
	Mode       dice.Mode `json:"mode"`
}

func Roll(expression string, dc int, reason string, rng *rand.Rand) (Result, error) {
	roll, err := dice.RollExpression(expression, rng)
	if err != nil {
		return Result{}, err
	}
	return Result{Reason: reason, Expression: roll.Expression, Rolls: roll.Rolls, Modifier: roll.Modifier, Total: roll.Total, DC: dc, Success: roll.Total >= dc, Mode: roll.Mode}, nil
}
