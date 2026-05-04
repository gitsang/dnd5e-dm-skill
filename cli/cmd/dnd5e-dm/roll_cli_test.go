package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRollCommandAppendsJSONL(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "roll_log.jsonl")
	command := exec.Command("go", "run", ".", "roll", "1d20+5", "--reason", "test attack", "--log", logPath, "--seed", "1")
	command.Env = append(os.Environ(), "PATH=/home/sang/.go/bin:"+os.Getenv("PATH"))
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("roll command failed: %v\n%s", err, output)
	}
	var payload map[string]any
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("stdout is not JSON: %v\n%s", err, output)
	}
	if payload["expression"] != "1d20+5" || payload["source"] != "script" {
		t.Fatalf("unexpected stdout payload: %#v", payload)
	}
	contents, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	var saved map[string]any
	if err := json.Unmarshal(contents, &saved); err != nil {
		t.Fatalf("log is not JSON: %v\n%s", err, contents)
	}
	if saved["reason"] != "test attack" || saved["source"] != "script" {
		t.Fatalf("unexpected saved payload: %#v", saved)
	}
}
