package resources

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func Spend(characterPath string, dottedPath string, amount int) error {
	return mutate(characterPath, dottedPath, -amount)
}

func Restore(characterPath string, dottedPath string, amount int) error {
	return mutate(characterPath, dottedPath, amount)
}

func mutate(characterPath string, dottedPath string, delta int) error {
	if amountAbs(delta) < 1 {
		return fmt.Errorf("amount must be positive")
	}
	contents, err := os.ReadFile(characterPath)
	if err != nil {
		return err
	}
	var data map[string]any
	if err := json.Unmarshal(contents, &data); err != nil {
		return err
	}
	parent, key, err := resolve(data, dottedPath)
	if err != nil {
		return err
	}
	current, _ := parent[key].(float64)
	next := int(current) + delta
	if next < 0 {
		return fmt.Errorf("insufficient resource")
	}
	parent[key] = next
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(characterPath, encoded, 0o644)
}

func resolve(data map[string]any, dottedPath string) (map[string]any, string, error) {
	parts := strings.Split(dottedPath, ".")
	if len(parts) == 0 || parts[0] == "" {
		return nil, "", fmt.Errorf("resource path is required")
	}
	current := data
	for _, part := range parts[:len(parts)-1] {
		next, ok := current[part].(map[string]any)
		if !ok {
			next = map[string]any{}
			current[part] = next
		}
		current = next
	}
	return current, parts[len(parts)-1], nil
}

func amountAbs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
