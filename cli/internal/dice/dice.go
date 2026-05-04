package dice

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Mode string

const (
	ModeNormal       Mode = "normal"
	ModeAdvantage    Mode = "advantage"
	ModeDisadvantage Mode = "disadvantage"
	ModeKeepHighest  Mode = "keep_highest"
)

type Expression struct {
	Expression  string
	Count       int
	Sides       int
	Modifier    int
	Mode        Mode
	KeepHighest int
}

type Result struct {
	Expression string `json:"expression"`
	Rolls      []int  `json:"rolls"`
	Modifier   int    `json:"modifier"`
	Total      int    `json:"total"`
	Mode       Mode   `json:"mode"`
}

var expressionPattern = regexp.MustCompile(`^(\d+)d(\d+)(adv|dis|kh\d+)?([+-]\d+)?$`)

func ParseExpression(input string) (Expression, error) {
	compact := strings.ToLower(strings.ReplaceAll(input, " ", ""))
	matches := expressionPattern.FindStringSubmatch(compact)
	if matches == nil {
		return Expression{}, fmt.Errorf("unsupported dice expression: %s", input)
	}
	count, err := strconv.Atoi(matches[1])
	if err != nil {
		return Expression{}, err
	}
	sides, err := strconv.Atoi(matches[2])
	if err != nil {
		return Expression{}, err
	}
	if count < 1 || sides < 2 {
		return Expression{}, fmt.Errorf("invalid dice expression: %s", input)
	}
	modifier := 0
	if matches[4] != "" {
		modifier, err = strconv.Atoi(matches[4])
		if err != nil {
			return Expression{}, err
		}
	}
	parsed := Expression{Expression: compact, Count: count, Sides: sides, Modifier: modifier, Mode: ModeNormal}
	switch rawMode := matches[3]; {
	case rawMode == "adv":
		parsed.Mode = ModeAdvantage
	case rawMode == "dis":
		parsed.Mode = ModeDisadvantage
	case strings.HasPrefix(rawMode, "kh"):
		keepHighest, err := strconv.Atoi(strings.TrimPrefix(rawMode, "kh"))
		if err != nil {
			return Expression{}, err
		}
		if keepHighest < 1 || keepHighest > count {
			return Expression{}, fmt.Errorf("invalid keep-highest expression: %s", input)
		}
		parsed.Mode = ModeKeepHighest
		parsed.KeepHighest = keepHighest
	}
	return parsed, nil
}

func RollExpression(input string, rng *rand.Rand) (Result, error) {
	parsed, err := ParseExpression(input)
	if err != nil {
		return Result{}, err
	}
	roller := rng
	if roller == nil {
		roller = rand.New(rand.NewSource(rand.Int63()))
	}
	if parsed.Mode == ModeAdvantage || parsed.Mode == ModeDisadvantage {
		rolls := []int{rollDie(roller, parsed.Sides), rollDie(roller, parsed.Sides)}
		chosen := rolls[0]
		if parsed.Mode == ModeAdvantage && rolls[1] > chosen {
			chosen = rolls[1]
		}
		if parsed.Mode == ModeDisadvantage && rolls[1] < chosen {
			chosen = rolls[1]
		}
		return Result{Expression: parsed.Expression, Rolls: rolls, Modifier: parsed.Modifier, Total: chosen + parsed.Modifier, Mode: parsed.Mode}, nil
	}
	rolls := make([]int, parsed.Count)
	for index := range rolls {
		rolls[index] = rollDie(roller, parsed.Sides)
	}
	total := 0
	if parsed.Mode == ModeKeepHighest {
		kept := append([]int(nil), rolls...)
		sort.Sort(sort.Reverse(sort.IntSlice(kept)))
		for _, roll := range kept[:parsed.KeepHighest] {
			total += roll
		}
	} else {
		for _, roll := range rolls {
			total += roll
		}
	}
	total += parsed.Modifier
	return Result{Expression: parsed.Expression, Rolls: rolls, Modifier: parsed.Modifier, Total: total, Mode: parsed.Mode}, nil
}

func rollDie(rng *rand.Rand, sides int) int {
	return rng.Intn(sides) + 1
}
