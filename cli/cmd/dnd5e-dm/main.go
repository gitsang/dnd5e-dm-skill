package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gitsang/dnd5e-dm-skill/cli/internal/audit"
	"github.com/gitsang/dnd5e-dm-skill/cli/internal/check"
	"github.com/gitsang/dnd5e-dm-skill/cli/internal/combat"
	"github.com/gitsang/dnd5e-dm-skill/cli/internal/conditions"
	"github.com/gitsang/dnd5e-dm-skill/cli/internal/dice"
	"github.com/gitsang/dnd5e-dm-skill/cli/internal/initiative"
	"github.com/gitsang/dnd5e-dm-skill/cli/internal/resources"
	"github.com/gitsang/dnd5e-dm-skill/cli/internal/rules"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" || os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		return
	}
	switch os.Args[1] {
	case "roll":
		if err := runRoll(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	case "check":
		exitOnError(runCheck(os.Args[2:]))
	case "initiative":
		exitOnError(runInitiative(os.Args[2:]))
	case "combat":
		exitOnError(runCombat(os.Args[2:]))
	case "resources":
		exitOnError(runResources(os.Args[2:]))
	case "conditions":
		exitOnError(runConditions(os.Args[2:]))
	case "rules":
		exitOnError(runRules(os.Args[2:]))
	}
	fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
	os.Exit(2)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

type rollLogEntry struct {
	Timestamp  string    `json:"timestamp"`
	Visibility string    `json:"visibility"`
	Source     string    `json:"source"`
	Expression string    `json:"expression"`
	Reason     string    `json:"reason"`
	Rolls      []int     `json:"rolls"`
	Modifier   int       `json:"modifier"`
	Total      int       `json:"total"`
	Mode       dice.Mode `json:"mode"`
}

func runRoll(args []string) error {
	reordered, err := moveRollExpressionLast(args)
	if err != nil {
		return err
	}
	flags := flag.NewFlagSet("roll", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	reason := flags.String("reason", "", "reason for the roll")
	logPath := flags.String("log", "", "roll_log.jsonl path")
	visibility := flags.String("visibility", "public", "public or dm_secret")
	source := flags.String("source", "script", "script or user")
	seed := flags.Int64("seed", 0, "deterministic seed")
	if err := flags.Parse(reordered); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("roll requires exactly one dice expression")
	}
	if *reason == "" || *logPath == "" {
		return fmt.Errorf("roll requires --reason and --log")
	}
	if *visibility != "public" && *visibility != "dm_secret" {
		return fmt.Errorf("visibility must be public or dm_secret")
	}
	if *source != "script" && *source != "user" {
		return fmt.Errorf("source must be script or user")
	}
	var rng *rand.Rand
	if *seed != 0 {
		rng = rand.New(rand.NewSource(*seed))
	}
	result, err := dice.RollExpression(flags.Arg(0), rng)
	if err != nil {
		return err
	}
	entry := rollLogEntry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		Visibility: *visibility,
		Source:     *source,
		Expression: result.Expression,
		Reason:     *reason,
		Rolls:      result.Rolls,
		Modifier:   result.Modifier,
		Total:      result.Total,
		Mode:       result.Mode,
	}
	if err := audit.AppendJSONL(*logPath, entry); err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(entry)
}

func moveRollExpressionLast(args []string) ([]string, error) {
	var expression string
	reordered := make([]string, 0, len(args))
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "--reason" || arg == "--log" || arg == "--visibility" || arg == "--source" || arg == "--seed" {
			if index+1 >= len(args) {
				return nil, fmt.Errorf("%s requires a value", arg)
			}
			reordered = append(reordered, arg, args[index+1])
			index++
			continue
		}
		if strings.HasPrefix(arg, "--") {
			reordered = append(reordered, arg)
			continue
		}
		if expression != "" {
			return nil, fmt.Errorf("roll requires exactly one dice expression")
		}
		expression = arg
	}
	if expression != "" {
		reordered = append(reordered, expression)
	}
	return reordered, nil
}

func runCheck(args []string) error {
	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	expression := flags.String("expression", "", "dice expression")
	dc := flags.Int("dc", 0, "difficulty class")
	reason := flags.String("reason", "", "reason")
	seed := flags.Int64("seed", 0, "deterministic seed")
	if err := flags.Parse(args); err != nil {
		return err
	}
	var rng *rand.Rand
	if *seed != 0 {
		rng = rand.New(rand.NewSource(*seed))
	}
	result, err := check.Roll(*expression, *dc, *reason, rng)
	if err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(result)
}

func runInitiative(args []string) error {
	flags := flag.NewFlagSet("initiative", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	combatantsJSON := flags.String("combatants", "", "combatants JSON")
	out := flags.String("out", "", "output combat_state.json")
	seed := flags.Int64("seed", 0, "deterministic seed")
	if err := flags.Parse(args); err != nil {
		return err
	}
	var inputs []initiative.InputCombatant
	if err := json.Unmarshal([]byte(*combatantsJSON), &inputs); err != nil {
		return err
	}
	var rng *rand.Rand
	if *seed != 0 {
		rng = rand.New(rand.NewSource(*seed))
	}
	state, err := initiative.Create(inputs, rng)
	if err != nil {
		return err
	}
	if *out != "" {
		if err := combat.Save(*out, state); err != nil {
			return err
		}
	}
	return json.NewEncoder(os.Stdout).Encode(state)
}

func runCombat(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("combat requires a subcommand")
	}
	flags := flag.NewFlagSet("combat "+args[0], flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	state := flags.String("state", "", "combat_state.json")
	combatant := flags.String("combatant", "", "combatant id")
	kind := flags.String("kind", "", "action kind")
	amount := flags.Int("amount", 0, "amount")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	switch args[0] {
	case "use":
		return combat.Use(*state, *combatant, *kind)
	case "damage":
		return combat.Damage(*state, *combatant, *amount)
	case "next-turn":
		return combat.NextTurn(*state)
	default:
		return fmt.Errorf("unknown combat subcommand: %s", args[0])
	}
}

func runResources(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("resources requires a subcommand")
	}
	flags := flag.NewFlagSet("resources "+args[0], flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	character := flags.String("character", "", "character JSON")
	path := flags.String("path", "", "resource path")
	amount := flags.Int("amount", 0, "amount")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	switch args[0] {
	case "spend":
		return resources.Spend(*character, *path, *amount)
	case "restore":
		return resources.Restore(*character, *path, *amount)
	default:
		return fmt.Errorf("unknown resources subcommand: %s", args[0])
	}
}

func runConditions(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("conditions requires a subcommand")
	}
	flags := flag.NewFlagSet("conditions "+args[0], flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	state := flags.String("state", "", "combat_state.json")
	combatant := flags.String("combatant", "", "combatant id")
	condition := flags.String("condition", "", "condition name")
	source := flags.String("source", "", "condition source")
	duration := flags.String("duration", "", "condition duration")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	switch args[0] {
	case "add":
		return conditions.Add(*state, *combatant, *condition, *source, *duration)
	case "remove":
		return conditions.Remove(*state, *combatant, *condition)
	case "list":
		items, err := conditions.List(*state, *combatant)
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(items)
	default:
		return fmt.Errorf("unknown conditions subcommand: %s", args[0])
	}
}

func runRules(args []string) error {
	if len(args) == 0 || args[0] != "search" {
		return fmt.Errorf("rules supports: search")
	}
	flags := flag.NewFlagSet("rules search", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	rulesDir := flags.String("rules-dir", "", "rules directory")
	query := flags.String("query", "", "query")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	results, err := rules.Search(strings.Split(*rulesDir, ","), *query)
	if err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(results)
}

func printUsage() {
	fmt.Print(`dnd5e-dm is a deterministic CLI for DnD5e DM skill workflows.

Usage:
  dnd5e-dm <command> [options]

Commands:
  roll         Roll dice and append audit JSONL
  check        Roll a non-mutating check against a DC
  initiative   Create initiative combat state
  combat       Mutate combat state
  resources    Spend or restore character resources
  conditions   Add, remove, or list combat conditions
  rules        Search local SRD/CC and user-provided rules
`)
}
