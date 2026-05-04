# DnD5e DM Skill Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a `dnd5e-dm` skill that runs DnD5e sessions as a rules-aware DM assistant with persistent state, auditable dice, module fidelity checks, and minimal deterministic combat/resource tools.

**Architecture:** Implement a discoverable skill under `.agents/skills/dnd5e-dm/`, with a Go binary CLI and references bundled beside `SKILL.md`. The skill treats a campaign vault as source of truth and delegates dice, initiative, turn advancement, HP/resource mutation, rules-source lookup, and audit logging to a deterministic `dnd5e-dm` command.

**Tech Stack:** Markdown skill docs written in Chinese, Go standard library, JSON/JSONL campaign state, Go unit tests, shell CLI smoke commands.

---

## File Structure

- Create: `.agents/skills/dnd5e-dm/SKILL.md` — skill trigger description and DM operating protocol.
- Create: `.agents/skills/dnd5e-dm/references/campaign-vault-schema.md` — required campaign vault files and schemas.
- Create: `.agents/skills/dnd5e-dm/references/dm-workflow.md` — mode-specific workflows for setup, prep, live play, combat, rules adjudication, and recap.
- Create: `.agents/skills/dnd5e-dm/references/rules-source-policy.md` — source/copyright policy and canon/homebrew labeling rules.
- Create: `cli/go.mod` — Go module definition.
- Create: `cli/cmd/dnd5e-dm/main.go` — binary CLI entrypoint.
- Create: `cli/internal/dice/` — dice parser and roller.
- Create: `cli/internal/audit/` — JSONL audit log writer.
- Create: `cli/internal/combat/` — turn advancement, action economy, HP damage/healing, concentration/death-save hooks.
- Create: `cli/internal/resources/` — spend/restore character resources.
- Create: `cli/internal/conditions/` — add/remove/list conditions.
- Create: `cli/internal/rules/` — local SRD/CC and user-provided rules search.
- Create: `.agents/skills/dnd5e-dm/evals/evals.json` — initial human-review test prompts.
- Create: `tests/dnd5e_dm/test_dice.py` — dice parser/roller tests.
- Create: `tests/dnd5e_dm/test_roll_cli.py` — roll logging tests.
- Create: `tests/dnd5e_dm/test_combat.py` — action economy and HP mutation tests.
- Create: `tests/dnd5e_dm/test_resources_conditions.py` — resource and condition mutation tests.

### Task 1: Create skill skeleton and references

**Files:**
- Create: `.agents/skills/dnd5e-dm/SKILL.md`
- Create: `.agents/skills/dnd5e-dm/references/campaign-vault-schema.md`
- Create: `.agents/skills/dnd5e-dm/references/dm-workflow.md`
- Create: `.agents/skills/dnd5e-dm/references/rules-source-policy.md`

- [ ] **Step 1: Create the directories**

Run:
```bash
mkdir -p ".agents/skills/dnd5e-dm/references" ".agents/skills/dnd5e-dm/scripts" ".agents/skills/dnd5e-dm/evals" "tests/dnd5e_dm"
```
Expected: all directories exist.

- [ ] **Step 2: Write `SKILL.md`**

Create `.agents/skills/dnd5e-dm/SKILL.md` with:
```md
---
name: dnd5e-dm
description: Use this skill whenever the user wants to run, prepare, adjudicate, audit, or continue a Dungeons & Dragons 5th Edition game. This includes acting as a DM, enforcing DnD5e combat/action economy, managing persistent campaign state, rolling dice, preserving module canon, checking monster/stat-block sources, preparing sessions, resolving rules questions, or summarizing sessions. Always use this skill for DnD5e or tabletop RPG requests where rules correctness, campaign continuity, dice rolls, monster statistics, or module fidelity matter.
---

# DnD5e DM

Use this skill as a rules-aware DM operating protocol, not as a freeform text adventure prompt. Separate narration from persistent state, dice, rules sources, and module canon.

## Core Rules

- Treat campaign vault files as source of truth; do not rely on conversation context for HP, resources, initiative, conditions, module facts, discovered clues, or NPC state.
- Use bundled scripts for dice, initiative, turn advancement, action economy, HP/resource mutation, and condition changes.
- Never invent official monster stat blocks, non-SRD rules text, or module facts. Use SRD-compatible material or user-provided sources; otherwise mark content as homebrew or request the missing source.
- Preserve player agency: never decide player character actions, emotions, dialogue, or intent.
- Keep hidden information hidden unless the user explicitly asks for DM-facing information.

## First Move

Identify the mode:
1. Campaign Setup
2. Session Prep
3. Live DM
4. Combat
5. Rules Adjudication
6. Post-Session Recap

Then read the relevant reference:
- `references/campaign-vault-schema.md` when creating or reading campaign state.
- `references/dm-workflow.md` for mode-specific procedure.
- `references/rules-source-policy.md` before using module, monster, spell, class, or non-SRD material.

## Mechanical Safety Gate

Before resolving any mechanical event, classify it:
- Tool-enforced: dice, initiative, HP, resources, conditions, action economy, death saves, concentration checks.
- Source-checked: spells, class features, monster stat blocks, exact condition interactions, rests, encounter XP/CR.
- DM judgment: DC choice, ability/skill choice, advantage/disadvantage when not dictated, narrative consequences.

Tool-enforced events must call scripts. Source-checked events must cite local approved sources or state that the source is missing. DM judgment should be brief, fair, and reversible when new rules evidence appears.

## Live Play Output

When narrating, describe only what characters can perceive. End most scene prompts with “你们怎么做？”. During combat, show current round, active turn, visible conditions, and remaining action economy.
```

- [ ] **Step 3: Write `campaign-vault-schema.md`**

Create `.agents/skills/dnd5e-dm/references/campaign-vault-schema.md` with:
```md
# Campaign Vault Schema

The campaign vault is the source of truth. Read these files before mechanical or continuity-sensitive decisions.

```text
campaigns/<campaign-name>/
  campaign_config.json
  party.json
  characters/<pc-name>.json
  combat_state.json
  world_state.json
  campaign_log.md
  roll_log.jsonl
  session_notes/
  module_canon/
  module_index.json
  dm_improv.md
  monster_statblocks/
  encounter_catalog.json
  rules_refs/
```

## `campaign_config.json`

```json
{
  "campaign_name": "Example Campaign",
  "rules_version": "5e-2014",
  "positioning": "zones",
  "allowed_sources": ["SRD", "user_provided_module"]
}
```

## `characters/<pc-name>.json`

```json
{
  "id": "pc1",
  "name": "Aric",
  "level": 3,
  "classes": [{"name": "fighter", "level": 3}],
  "ac": 16,
  "hp": 24,
  "max_hp": 28,
  "temp_hp": 0,
  "abilities": {"str": 16, "dex": 12, "con": 14, "int": 10, "wis": 11, "cha": 8},
  "proficiency_bonus": 2,
  "skills": {"athletics": 5},
  "saves": {"str": 5, "con": 4},
  "resources": {"action_surge": 1, "hit_dice_d10": 3},
  "spell_slots": {},
  "conditions": []
}
```

## `combat_state.json`

```json
{
  "round": 1,
  "active_index": 0,
  "combatants": [
    {
      "id": "pc1",
      "name": "Aric",
      "initiative": 17,
      "ac": 16,
      "hp": 24,
      "max_hp": 28,
      "conditions": [],
      "concentration": null,
      "death_saves": {"successes": 0, "failures": 0},
      "position": {"zone": "tavern floor"},
      "used": {"action": false, "bonus_action": false, "reaction": false, "movement": 0, "object_interaction": false}
    }
  ]
}
```

## Logs and Canon

- `campaign_log.md`: confirmed table events only.
- `roll_log.jsonl`: one JSON object per dice roll.
- `module_index.json`: chapter goals, key NPCs, locations, clues, approved encounters.
- `dm_improv.md`: improvised additions and whether they are `derived_from_module`, `dm_improvised`, or `homebrew`.
- `monster_statblocks/`: approved stat blocks only. Missing stat blocks must not be invented as official.
```

- [ ] **Step 4: Write workflow and source policy references**

Create `.agents/skills/dnd5e-dm/references/dm-workflow.md` with:
```md
# DM Workflow

## Campaign Setup

Create the vault, ask for rules version, import party files, index module material, and record allowed sources. If a source is missing, mark dependent rulings as unavailable rather than guessing.

## Session Prep

Read `campaign_log.md`, `world_state.json`, `module_index.json`, party state, and unresolved hooks. Prepare scenes that preserve module goals while allowing player agency.

## Live DM

Narrate only perceivable facts. Before mechanical outcomes, read the relevant vault files and use scripts. End most scene prompts with “你们怎么做？”.

## Combat

Read `combat_state.json`. Show round, active combatant, visible conditions, and remaining action economy. Classify every declared activity as action, bonus action, reaction, movement, object interaction, or free communication before resolution.

## Rules Adjudication

Classify the question as tool-enforced, source-checked, or DM judgment. Cite local approved sources when exact rules are required. If the source is missing, say so and offer an SRD-compatible or homebrew-labeled alternative.

## Post-Session

Update confirmed events, NPC state, faction clocks, discovered clues, unresolved hooks, resource changes, and next-session prep. Do not reveal hidden module information in player-facing recaps.
```

Create `.agents/skills/dnd5e-dm/references/rules-source-policy.md` with:
```md
# Rules Source Policy

- Allowed sources: SRD-compatible rules, user-provided character sheets, user-provided modules, and user-approved stat blocks.
- Do not reproduce or invent exact non-SRD official text or stat blocks.
- If a monster, spell, class feature, or module fact is not present in approved local sources, state that the source is missing.
- New content must be labeled as one of: `module_canon`, `derived_from_module`, `dm_improvised`, or `homebrew`.
- Hidden module facts stay DM-only until discovered by players or explicitly requested by the user in DM-facing mode.
```

- [ ] **Step 5: Verify skeleton files exist**

Run:
```bash
test -f ".agents/skills/dnd5e-dm/SKILL.md" && test -f ".agents/skills/dnd5e-dm/references/campaign-vault-schema.md" && test -f ".agents/skills/dnd5e-dm/references/dm-workflow.md" && test -f ".agents/skills/dnd5e-dm/references/rules-source-policy.md"
```
Expected: command exits 0.

### Task 2: Implement and test dice parsing

**Files:**
- Create: `.agents/skills/dnd5e-dm/scripts/dice.py`
- Create: `tests/dnd5e_dm/test_dice.py`

- [ ] **Step 1: Write failing dice tests**

Create `tests/dnd5e_dm/test_dice.py`:
```python
import random

from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / ".agents" / "skills" / "dnd5e-dm" / "scripts"))

from dice import parse_expression, roll_expression


def test_parse_simple_modifier():
    parsed = parse_expression("2d6+3")
    assert parsed.count == 2
    assert parsed.sides == 6
    assert parsed.modifier == 3
    assert parsed.mode == "normal"


def test_parse_advantage():
    parsed = parse_expression("1d20adv+7")
    assert parsed.count == 1
    assert parsed.sides == 20
    assert parsed.modifier == 7
    assert parsed.mode == "advantage"


def test_roll_expression_is_deterministic_with_rng():
    rng = random.Random(1)
    result = roll_expression("1d20+5", rng=rng)
    assert result.expression == "1d20+5"
    assert result.rolls == [5]
    assert result.modifier == 5
    assert result.total == 10
```

- [ ] **Step 2: Run tests to verify failure**

Run:
```bash
pytest tests/dnd5e_dm/test_dice.py -q
```
Expected: FAIL because `dice` module does not exist.

- [ ] **Step 3: Implement `dice.py`**

Create `.agents/skills/dnd5e-dm/scripts/dice.py`:
```python
from __future__ import annotations

from dataclasses import dataclass
import random
import re


EXPR_RE = re.compile(r"^(?P<count>\d+)d(?P<sides>\d+)(?P<mode>adv|dis|kh\d+)?(?P<mod>[+-]\d+)?$")


@dataclass(frozen=True)
class DiceExpression:
    expression: str
    count: int
    sides: int
    modifier: int
    mode: str
    keep_highest: int | None = None


@dataclass(frozen=True)
class RollResult:
    expression: str
    rolls: list[int]
    modifier: int
    total: int
    mode: str


def parse_expression(expression: str) -> DiceExpression:
    compact = expression.replace(" ", "").lower()
    match = EXPR_RE.match(compact)
    if not match:
        raise ValueError(f"Unsupported dice expression: {expression}")
    count = int(match.group("count"))
    sides = int(match.group("sides"))
    if count < 1 or sides < 2:
        raise ValueError(f"Invalid dice expression: {expression}")
    raw_mode = match.group("mode") or "normal"
    keep_highest = None
    if raw_mode == "adv":
        mode = "advantage"
    elif raw_mode == "dis":
        mode = "disadvantage"
    elif raw_mode.startswith("kh"):
        mode = "keep_highest"
        keep_highest = int(raw_mode[2:])
        if keep_highest < 1 or keep_highest > count:
            raise ValueError(f"Invalid keep-highest expression: {expression}")
    else:
        mode = "normal"
    modifier = int(match.group("mod") or 0)
    return DiceExpression(compact, count, sides, modifier, mode, keep_highest)


def roll_expression(expression: str, rng: random.Random | None = None) -> RollResult:
    parsed = parse_expression(expression)
    roller = rng or random.SystemRandom()
    if parsed.mode in {"advantage", "disadvantage"}:
        rolls = [roller.randint(1, parsed.sides), roller.randint(1, parsed.sides)]
        chosen = max(rolls) if parsed.mode == "advantage" else min(rolls)
        total = chosen + parsed.modifier
    else:
        rolls = [roller.randint(1, parsed.sides) for _ in range(parsed.count)]
        if parsed.mode == "keep_highest":
            kept = sorted(rolls, reverse=True)[: parsed.keep_highest]
            total = sum(kept) + parsed.modifier
        else:
            total = sum(rolls) + parsed.modifier
    return RollResult(parsed.expression, rolls, parsed.modifier, total, parsed.mode)
```

- [ ] **Step 4: Run tests to verify pass**

Run:
```bash
pytest tests/dnd5e_dm/test_dice.py -q
```
Expected: `3 passed`.

### Task 3: Implement roll CLI and audit log

**Files:**
- Create: `.agents/skills/dnd5e-dm/scripts/roll.py`
- Create: `tests/dnd5e_dm/test_roll_cli.py`

- [ ] **Step 1: Write failing roll CLI test**

Create `tests/dnd5e_dm/test_roll_cli.py`:
```python
import json
from pathlib import Path
import subprocess
import sys


ROOT = Path(__file__).resolve().parents[2]
ROLL = ROOT / ".agents" / "skills" / "dnd5e-dm" / "scripts" / "roll.py"


def test_roll_cli_appends_jsonl(tmp_path):
    log = tmp_path / "roll_log.jsonl"
    result = subprocess.run(
        [sys.executable, str(ROLL), "1d20+5", "--reason", "test attack", "--log", str(log), "--seed", "1"],
        check=True,
        text=True,
        capture_output=True,
    )
    payload = json.loads(result.stdout)
    assert payload["expression"] == "1d20+5"
    assert payload["total"] == 10
    saved = json.loads(log.read_text().strip())
    assert saved["reason"] == "test attack"
    assert saved["source"] == "script"
```

- [ ] **Step 2: Run test to verify failure**

Run:
```bash
pytest tests/dnd5e_dm/test_roll_cli.py -q
```
Expected: FAIL because `roll.py` does not exist.

- [ ] **Step 3: Implement `roll.py`**

Create `.agents/skills/dnd5e-dm/scripts/roll.py`:
```python
from __future__ import annotations

import argparse
from datetime import datetime, timezone
import json
from pathlib import Path
import random

from dice import roll_expression


def main() -> None:
    parser = argparse.ArgumentParser(description="Roll DnD dice and append an audit log entry.")
    parser.add_argument("expression")
    parser.add_argument("--reason", required=True)
    parser.add_argument("--log", required=True)
    parser.add_argument("--visibility", default="public", choices=["public", "dm_secret"])
    parser.add_argument("--source", default="script", choices=["script", "user"])
    parser.add_argument("--seed", type=int)
    args = parser.parse_args()

    rng = random.Random(args.seed) if args.seed is not None else None
    result = roll_expression(args.expression, rng=rng)
    entry = {
        "timestamp": datetime.now(timezone.utc).isoformat(),
        "visibility": args.visibility,
        "source": args.source,
        "expression": result.expression,
        "reason": args.reason,
        "rolls": result.rolls,
        "modifier": result.modifier,
        "total": result.total,
        "mode": result.mode,
    }
    log_path = Path(args.log)
    log_path.parent.mkdir(parents=True, exist_ok=True)
    with log_path.open("a", encoding="utf-8") as handle:
        handle.write(json.dumps(entry, ensure_ascii=False) + "\n")
    print(json.dumps(entry, ensure_ascii=False))


if __name__ == "__main__":
    main()
```

- [ ] **Step 4: Run tests to verify pass**

Run:
```bash
pytest tests/dnd5e_dm/test_dice.py tests/dnd5e_dm/test_roll_cli.py -q
```
Expected: all tests pass.

### Task 4: Implement minimal combat state tools

**Files:**
- Create: `.agents/skills/dnd5e-dm/scripts/combat.py`
- Create: `tests/dnd5e_dm/test_combat.py`

- [ ] **Step 1: Write failing combat tests**

Create `tests/dnd5e_dm/test_combat.py`:
```python
import json
from pathlib import Path
import subprocess
import sys


ROOT = Path(__file__).resolve().parents[2]
COMBAT = ROOT / ".agents" / "skills" / "dnd5e-dm" / "scripts" / "combat.py"


def write_state(path: Path) -> None:
    path.write_text(json.dumps({
        "round": 1,
        "active_index": 0,
        "combatants": [
            {"id": "pc1", "name": "Aric", "hp": 12, "max_hp": 12, "used": {"action": False, "bonus_action": False, "reaction": False, "movement": 0}},
            {"id": "gob1", "name": "Goblin", "hp": 7, "max_hp": 7, "used": {"action": False, "bonus_action": False, "reaction": False, "movement": 0}}
        ]
    }), encoding="utf-8")


def test_use_action_prevents_duplicate_action(tmp_path):
    state = tmp_path / "combat_state.json"
    write_state(state)
    subprocess.run([sys.executable, str(COMBAT), "use", "--state", str(state), "--combatant", "pc1", "--kind", "action"], check=True)
    duplicate = subprocess.run([sys.executable, str(COMBAT), "use", "--state", str(state), "--combatant", "pc1", "--kind", "action"], text=True, capture_output=True)
    assert duplicate.returncode != 0
    assert "already used action" in duplicate.stderr


def test_apply_damage_updates_hp(tmp_path):
    state = tmp_path / "combat_state.json"
    write_state(state)
    subprocess.run([sys.executable, str(COMBAT), "damage", "--state", str(state), "--combatant", "gob1", "--amount", "3"], check=True)
    saved = json.loads(state.read_text())
    assert saved["combatants"][1]["hp"] == 4
```

- [ ] **Step 2: Run tests to verify failure**

Run:
```bash
pytest tests/dnd5e_dm/test_combat.py -q
```
Expected: FAIL because `combat.py` does not exist.

- [ ] **Step 3: Implement `combat.py`**

Create `.agents/skills/dnd5e-dm/scripts/combat.py`:
```python
from __future__ import annotations

import argparse
import json
from pathlib import Path
import sys


def load(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


def save(path: Path, data: dict) -> None:
    path.write_text(json.dumps(data, indent=2, ensure_ascii=False), encoding="utf-8")


def find_combatant(state: dict, combatant_id: str) -> dict:
    for combatant in state["combatants"]:
        if combatant["id"] == combatant_id:
            return combatant
    raise SystemExit(f"unknown combatant: {combatant_id}")


def reset_used(combatant: dict) -> None:
    combatant["used"] = {"action": False, "bonus_action": False, "reaction": False, "movement": 0, "object_interaction": False}


def main() -> None:
    parser = argparse.ArgumentParser(description="Mutate DnD5e combat state.")
    sub = parser.add_subparsers(dest="command", required=True)
    use = sub.add_parser("use")
    use.add_argument("--state", required=True)
    use.add_argument("--combatant", required=True)
    use.add_argument("--kind", required=True, choices=["action", "bonus_action", "reaction", "object_interaction"])
    damage = sub.add_parser("damage")
    damage.add_argument("--state", required=True)
    damage.add_argument("--combatant", required=True)
    damage.add_argument("--amount", required=True, type=int)
    next_turn = sub.add_parser("next-turn")
    next_turn.add_argument("--state", required=True)
    args = parser.parse_args()

    state_path = Path(args.state)
    state = load(state_path)
    if args.command == "use":
        combatant = find_combatant(state, args.combatant)
        used = combatant.setdefault("used", {})
        if used.get(args.kind):
            print(f"{combatant['name']} already used {args.kind}", file=sys.stderr)
            raise SystemExit(1)
        used[args.kind] = True
    elif args.command == "damage":
        combatant = find_combatant(state, args.combatant)
        combatant["hp"] = max(0, int(combatant.get("hp", 0)) - args.amount)
    elif args.command == "next-turn":
        state["active_index"] = (int(state.get("active_index", 0)) + 1) % len(state["combatants"])
        if state["active_index"] == 0:
            state["round"] = int(state.get("round", 1)) + 1
        reset_used(state["combatants"][state["active_index"]])
    save(state_path, state)
    print(json.dumps(state, ensure_ascii=False))


if __name__ == "__main__":
    main()
```

- [ ] **Step 4: Run combat tests**

Run:
```bash
pytest tests/dnd5e_dm/test_combat.py -q
```
Expected: `2 passed`.

### Task 5: Implement resource and condition tools

**Files:**
- Create: `.agents/skills/dnd5e-dm/scripts/resources.py`
- Create: `.agents/skills/dnd5e-dm/scripts/conditions.py`
- Create: `tests/dnd5e_dm/test_resources_conditions.py`

- [ ] **Step 1: Write failing resource/condition tests**

Create `tests/dnd5e_dm/test_resources_conditions.py`:
```python
import json
from pathlib import Path
import subprocess
import sys


ROOT = Path(__file__).resolve().parents[2]
RESOURCES = ROOT / ".agents" / "skills" / "dnd5e-dm" / "scripts" / "resources.py"
CONDITIONS = ROOT / ".agents" / "skills" / "dnd5e-dm" / "scripts" / "conditions.py"


def test_spend_resource_and_reject_overspend(tmp_path):
    character = tmp_path / "aric.json"
    character.write_text(json.dumps({"spell_slots": {"level_1": 1}}), encoding="utf-8")
    subprocess.run([sys.executable, str(RESOURCES), "spend", "--character", str(character), "--path", "spell_slots.level_1", "--amount", "1"], check=True)
    saved = json.loads(character.read_text())
    assert saved["spell_slots"]["level_1"] == 0
    failed = subprocess.run([sys.executable, str(RESOURCES), "spend", "--character", str(character), "--path", "spell_slots.level_1", "--amount", "1"], text=True, capture_output=True)
    assert failed.returncode != 0
    assert "insufficient resource" in failed.stderr


def test_add_and_remove_condition(tmp_path):
    state = tmp_path / "combat_state.json"
    state.write_text(json.dumps({"combatants": [{"id": "pc1", "name": "Aric", "conditions": []}]}), encoding="utf-8")
    subprocess.run([sys.executable, str(CONDITIONS), "add", "--state", str(state), "--combatant", "pc1", "--condition", "poisoned"], check=True)
    assert json.loads(state.read_text())["combatants"][0]["conditions"][0]["name"] == "poisoned"
    subprocess.run([sys.executable, str(CONDITIONS), "remove", "--state", str(state), "--combatant", "pc1", "--condition", "poisoned"], check=True)
    assert json.loads(state.read_text())["combatants"][0]["conditions"] == []
```

- [ ] **Step 2: Implement `resources.py`**

Create `.agents/skills/dnd5e-dm/scripts/resources.py`:
```python
from __future__ import annotations

import argparse
import json
from pathlib import Path
import sys


def resolve(data: dict, dotted: str) -> tuple[dict, str]:
    parts = dotted.split(".")
    current = data
    for part in parts[:-1]:
        current = current.setdefault(part, {})
    return current, parts[-1]


def main() -> None:
    parser = argparse.ArgumentParser(description="Spend or restore character resources.")
    parser.add_argument("command", choices=["spend", "restore"])
    parser.add_argument("--character", required=True)
    parser.add_argument("--path", required=True)
    parser.add_argument("--amount", required=True, type=int)
    args = parser.parse_args()
    path = Path(args.character)
    data = json.loads(path.read_text(encoding="utf-8"))
    parent, key = resolve(data, args.path)
    current = int(parent.get(key, 0))
    next_value = current - args.amount if args.command == "spend" else current + args.amount
    if next_value < 0:
        print("insufficient resource", file=sys.stderr)
        raise SystemExit(1)
    parent[key] = next_value
    path.write_text(json.dumps(data, indent=2, ensure_ascii=False), encoding="utf-8")
    print(json.dumps(data, ensure_ascii=False))


if __name__ == "__main__":
    main()
```

- [ ] **Step 3: Implement `conditions.py`**

Create `.agents/skills/dnd5e-dm/scripts/conditions.py`:
```python
from __future__ import annotations

import argparse
import json
from pathlib import Path


def find_combatant(state: dict, combatant_id: str) -> dict:
    for combatant in state["combatants"]:
        if combatant["id"] == combatant_id:
            return combatant
    raise SystemExit(f"unknown combatant: {combatant_id}")


def main() -> None:
    parser = argparse.ArgumentParser(description="Mutate combat conditions.")
    parser.add_argument("command", choices=["add", "remove", "list"])
    parser.add_argument("--state", required=True)
    parser.add_argument("--combatant", required=True)
    parser.add_argument("--condition")
    parser.add_argument("--source", default="")
    parser.add_argument("--duration", default="")
    args = parser.parse_args()
    path = Path(args.state)
    state = json.loads(path.read_text(encoding="utf-8"))
    combatant = find_combatant(state, args.combatant)
    conditions = combatant.setdefault("conditions", [])
    if args.command == "add":
        if not args.condition:
            raise SystemExit("--condition is required")
        if not any(item["name"] == args.condition for item in conditions):
            conditions.append({"name": args.condition, "source": args.source, "duration": args.duration})
        path.write_text(json.dumps(state, indent=2, ensure_ascii=False), encoding="utf-8")
    elif args.command == "remove":
        combatant["conditions"] = [item for item in conditions if item["name"] != args.condition]
        path.write_text(json.dumps(state, indent=2, ensure_ascii=False), encoding="utf-8")
    print(json.dumps(combatant.get("conditions", []), ensure_ascii=False))


if __name__ == "__main__":
    main()
```

- [ ] **Step 4: Run resource/condition tests**

Run:
```bash
pytest tests/dnd5e_dm/test_resources_conditions.py -q
```
Expected: tests pass and failed overspend exits non-zero.

### Task 6: Implement initiative and check helpers

**Files:**
- Create: `.agents/skills/dnd5e-dm/scripts/initiative.py`
- Create: `.agents/skills/dnd5e-dm/scripts/check.py`

- [ ] **Step 1: Implement `initiative.py`**

Create `.agents/skills/dnd5e-dm/scripts/initiative.py`:
```python
from __future__ import annotations

import argparse
import json
from pathlib import Path
import random

from dice import roll_expression


def main() -> None:
    parser = argparse.ArgumentParser(description="Create a combat_state.json with rolled initiative.")
    parser.add_argument("--combatants", required=True)
    parser.add_argument("--out", required=True)
    parser.add_argument("--seed", type=int)
    args = parser.parse_args()
    rng = random.Random(args.seed) if args.seed is not None else None
    combatants = json.loads(args.combatants)
    entries = []
    for combatant in combatants:
        bonus = int(combatant.get("initiative_bonus", 0))
        result = roll_expression(f"1d20{bonus:+d}", rng=rng)
        item = dict(combatant)
        item["initiative"] = result.total
        item["initiative_roll"] = {"rolls": result.rolls, "modifier": bonus, "total": result.total}
        item.setdefault("conditions", [])
        item["used"] = {"action": False, "bonus_action": False, "reaction": False, "movement": 0, "object_interaction": False}
        entries.append(item)
    state = {"round": 1, "active_index": 0, "combatants": sorted(entries, key=lambda item: item["initiative"], reverse=True)}
    Path(args.out).write_text(json.dumps(state, indent=2, ensure_ascii=False), encoding="utf-8")
    print(json.dumps(state, ensure_ascii=False))


if __name__ == "__main__":
    main()
```

- [ ] **Step 2: Implement `check.py`**

Create `.agents/skills/dnd5e-dm/scripts/check.py`:
```python
from __future__ import annotations

import argparse
import json
import random

from dice import roll_expression


def main() -> None:
    parser = argparse.ArgumentParser(description="Roll a non-mutating check against a DC.")
    parser.add_argument("--expression", required=True)
    parser.add_argument("--dc", required=True, type=int)
    parser.add_argument("--reason", required=True)
    parser.add_argument("--seed", type=int)
    args = parser.parse_args()
    rng = random.Random(args.seed) if args.seed is not None else None
    result = roll_expression(args.expression, rng=rng)
    print(json.dumps({"reason": args.reason, "expression": result.expression, "rolls": result.rolls, "modifier": result.modifier, "total": result.total, "dc": args.dc, "success": result.total >= args.dc}, ensure_ascii=False))


if __name__ == "__main__":
    main()
```

- [ ] **Step 3: Smoke test both CLIs**

Run:
```bash
python .agents/skills/dnd5e-dm/scripts/check.py --expression 1d20+5 --dc 15 --reason "Athletics check" --seed 1
python .agents/skills/dnd5e-dm/scripts/initiative.py --combatants '[{"id":"pc1","name":"Aric","initiative_bonus":2},{"id":"gob1","name":"Goblin","initiative_bonus":1}]' --out /tmp/dnd5e-combat-state.json --seed 1
```
Expected: both commands print/write valid JSON.

### Task 7: Add skill eval prompts

**Files:**
- Create: `.agents/skills/dnd5e-dm/evals/evals.json`

- [ ] **Step 1: Write eval prompts**

Create `.agents/skills/dnd5e-dm/evals/evals.json`:
```json
{
  "skill_name": "dnd5e-dm",
  "evals": [
    {
      "id": 1,
      "prompt": "我正在跑 DnD5e 2014，玩家和两个 goblin 进入战斗。请作为 DM 建立先攻、公开 roll 点，并说明第一回合当前角色还能使用哪些行动经济。",
      "expected_output": "Uses scripts for rolls/initiative, does not invent hidden state, and displays action/bonus/reaction/movement status.",
      "files": []
    },
    {
      "id": 2,
      "prompt": "战士本回合已经 Attack 了，现在他说还要再 Dash。请按 5e 行动经济裁定。",
      "expected_output": "Classifies Dash as an Action, checks whether action is already used, and refuses unless a specific feature grants another action.",
      "files": []
    },
    {
      "id": 3,
      "prompt": "我想在官方模组里加入一个 Mind Flayer 伏击，但资料库里没有它的 stat block。你作为 DM assistant 应该怎么处理？",
      "expected_output": "Does not invent official stat block; asks for approved source, suggests SRD-compatible alternative, or marks homebrew clearly.",
      "files": []
    },
    {
      "id": 4,
      "prompt": "玩家跳过了模组关键线索，直接离开村庄。帮我即兴处理，但不要偏离主线，也不要直接泄露秘密。",
      "expected_output": "Checks module goals conceptually, preserves hidden info, introduces derived/improvised hooks, and labels deviation risk.",
      "files": []
    }
  ]
}
```

- [ ] **Step 2: Validate JSON**

Run:
```bash
python -m json.tool .agents/skills/dnd5e-dm/evals/evals.json >/tmp/dnd5e-evals.json
```
Expected: command exits 0.

### Task 8: Verify package quality

**Files:**
- Modify: any files that fail validation.

- [ ] **Step 1: Run all tests**

Run:
```bash
pytest tests/dnd5e_dm -q
```
Expected: all tests pass.

- [ ] **Step 2: Validate skill frontmatter and JSON files**

Run:
```bash
python -m json.tool .agents/skills/dnd5e-dm/evals/evals.json >/tmp/dnd5e-evals.json
python - <<'PY'
from pathlib import Path
text = Path('.agents/skills/dnd5e-dm/SKILL.md').read_text()
assert text.startswith('---\nname: dnd5e-dm\n')
assert 'description:' in text.split('---', 2)[1]
print('skill frontmatter ok')
PY
```
Expected: JSON validates and script prints `skill frontmatter ok`.

- [ ] **Step 3: Manual smoke test roll logging**

Run:
```bash
tmpdir=$(mktemp -d)
python .agents/skills/dnd5e-dm/scripts/roll.py 1d20+5 --reason "smoke test" --log "$tmpdir/roll_log.jsonl" --seed 1
test -s "$tmpdir/roll_log.jsonl"
```
Expected: command prints roll JSON and log file is non-empty.

- [ ] **Step 4: Commit when explicitly requested**

Only commit if the user explicitly asks. Suggested message:
```bash
git add .agents/skills/dnd5e-dm tests/dnd5e_dm docs/superpowers/specs/2026-05-04-dnd5e-dm-skill-design.md docs/superpowers/plans/2026-05-04-dnd5e-dm-skill.md
git commit -m "feat: add DnD5e DM skill plan"
```

## Self-Review

- Spec coverage: the plan covers skill trigger behavior, campaign vault persistence, dice audit logging, tool-enforced combat/resource state, module/source policy, hidden information, eval prompts, and verification.
- Placeholder scan: no `TBD`, `TODO`, or unspecified future tasks remain. All code-bearing tasks include concrete file content or explicit verification commands.
- Type consistency: script names, paths, JSON files, and test import paths all use `.agents/skills/dnd5e-dm/` and `tests/dnd5e_dm/` consistently.
