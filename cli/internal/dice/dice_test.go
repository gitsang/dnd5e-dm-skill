package dice

import (
	"math/rand"
	"testing"
)

func TestParseSimpleModifier(t *testing.T) {
	parsed, err := ParseExpression("2d6+3")
	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}
	if parsed.Count != 2 || parsed.Sides != 6 || parsed.Modifier != 3 || parsed.Mode != ModeNormal {
		t.Fatalf("unexpected parse result: %+v", parsed)
	}
}

func TestParseAdvantage(t *testing.T) {
	parsed, err := ParseExpression("1d20adv+7")
	if err != nil {
		t.Fatalf("ParseExpression returned error: %v", err)
	}
	if parsed.Count != 1 || parsed.Sides != 20 || parsed.Modifier != 7 || parsed.Mode != ModeAdvantage {
		t.Fatalf("unexpected parse result: %+v", parsed)
	}
}

func TestRollExpressionIsDeterministicWithRNG(t *testing.T) {
	result, err := RollExpression("1d20+5", rand.New(rand.NewSource(1)))
	if err != nil {
		t.Fatalf("RollExpression returned error: %v", err)
	}
	if result.Expression != "1d20+5" {
		t.Fatalf("expression = %q", result.Expression)
	}
	if len(result.Rolls) != 1 || result.Rolls[0] != 2 {
		t.Fatalf("rolls = %#v", result.Rolls)
	}
	if result.Modifier != 5 || result.Total != 7 {
		t.Fatalf("modifier/total = %d/%d", result.Modifier, result.Total)
	}
}

func TestRollKeepHighest(t *testing.T) {
	result, err := RollExpression("4d6kh3", rand.New(rand.NewSource(2)))
	if err != nil {
		t.Fatalf("RollExpression returned error: %v", err)
	}
	if len(result.Rolls) != 4 {
		t.Fatalf("roll count = %d", len(result.Rolls))
	}
	if result.Mode != ModeKeepHighest {
		t.Fatalf("mode = %q", result.Mode)
	}
}
