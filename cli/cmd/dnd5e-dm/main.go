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
	"github.com/gitsang/dnd5e-dm-skill/cli/internal/dice"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" || os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		return
	}
	if os.Args[1] == "roll" {
		if err := runRoll(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
	os.Exit(2)
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
