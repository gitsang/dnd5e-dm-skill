package check

import (
	"math/rand"
	"testing"
)

func TestRollCheckReportsSuccess(t *testing.T) {
	result, err := Roll("1d20+5", 7, "Athletics", rand.New(rand.NewSource(1)))
	if err != nil {
		t.Fatalf("roll check failed: %v", err)
	}
	if !result.Success || result.Reason != "Athletics" || result.DC != 7 {
		t.Fatalf("unexpected result: %#v", result)
	}
}
