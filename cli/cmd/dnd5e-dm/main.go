package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" || os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		return
	}
	fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
	os.Exit(2)
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
