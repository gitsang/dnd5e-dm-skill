# DnD5e DM Skill Design

## Purpose

Design a `dnd5e-dm` skill that helps an AI assistant run or assist a Dungeons & Dragons 5th Edition game while preserving mechanical correctness, campaign continuity, module fidelity, and auditable randomness. The skill should not behave like a freeform text adventure. It should operate as a DM workflow that separates narration from rules, state, canon, and dice.

## Recommended Approach

Use a **Skill + Campaign Vault + Go Binary CLI** architecture.

The skill provides the operating protocol. A local campaign directory stores persistent state. A deterministic Go CLI handles dice, initiative, combat state transitions, resource spending, rules-source lookup, and audit logs. This gives better reliability than a prompt-only skill without requiring a complete DnD rules engine from day one.

## Non-Goals

- Do not build a full replacement for official DnD sourcebooks.
- Do not invent official monster stat blocks or rules text when no approved source is available.
- Do not rely on conversation context as the source of truth for combat, resources, module facts, or discovered clues.
- Do not let the LLM generate dice outcomes directly.

## Architecture

```text
dnd5e-dm/
  SKILL.md
  bin/
    dnd5e-dm
  references/
    campaign-vault-schema.md
    dm-workflow.md
    rules-source-policy.md
  cli/
    go.mod
    cmd/dnd5e-dm/
      main.go
    internal/
      audit/
      combat/
      conditions/
      dice/
      resources/
      rules/
```

Per campaign, the skill expects a campaign vault:

```text
campaigns/<campaign-name>/
  campaign_config.json
  party.json
  characters/
    <pc-name>.json
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

## Data Persistence

The skill must treat files as source of truth. Before any important mechanical or continuity-sensitive decision, it should read the relevant state files instead of trusting conversational memory.

Persist these categories:

1. **Party and Character State**
   - Level, class, subclass, ancestry/species, background.
   - AC, max HP, current HP, temporary HP.
   - Ability scores, saves, skills, passive perception/investigation/insight.
   - Proficiency bonus, weapons, armor, spell save DC, attack bonuses.
   - Spell slots, prepared/known spells, class resources, consumables.

2. **Combat State**
   - Initiative order, active turn, round number.
   - Creature HP, AC, conditions, concentration, death saves.
   - Whether each combatant has used action, bonus action, reaction, movement, and object interaction.
   - Simplified position data, either grid coordinates or zone names.

3. **Campaign and World State**
   - Confirmed events from play.
   - NPC attitude, goals, fears, secrets, location, alive/dead/unknown status.
   - Faction clocks, location changes, open hooks, unresolved clues.
   - Player-visible facts separated from DM-only hidden facts.

4. **Module Fidelity State**
   - Canon facts from the module.
   - Chapter goals and main-path milestones.
   - Required clues, key NPCs, key locations, approved encounters.
   - Deviations caused by player action.
   - Improvised additions, clearly labeled as non-canon or derived.

5. **Audit Logs**
   - Dice rolls in `roll_log.jsonl`.
   - State changes and session summaries.
   - Manual overrides with reason and timestamp.

## Rules Compliance

The skill should divide rules into three levels:

### Level 1: Must Be Tool-Enforced

- Dice rolling.
- Initiative ordering and turn advancement.
- HP damage/healing application.
- Condition add/remove/list.
- Action economy tracking.
- Resource spending and restoration.
- Death saves and concentration checks.

### Level 2: Must Be Source-Checked

- Spell effects.
- Class/subclass features.
- Monster stat blocks.
- Conditions and edge-case interactions.
- Resting rules.
- Encounter XP/CR guidance.

### Level 3: DM Judgment Allowed

- Ability check selection.
- DC selection.
- Advantage/disadvantage when rules do not dictate it.
- Consequences of success or failure.
- NPC tactics and narrative reactions.

For Level 3 decisions, the skill should explain the reasoning briefly and preserve player agency. It should ask for a roll only when both success and failure produce interesting outcomes.

## Action Economy Protocol

During combat, every player or monster action must be classified before resolution:

- Action
- Bonus Action
- Reaction
- Movement
- Free object interaction
- No-action/free communication

The skill must check `combat_state.json` before allowing a declared action. It must not allow a second action, bonus action, or reaction in the same turn unless a specific rule, feature, spell, or condition grants it. At the start of each turn, available action economy should be displayed or internally reset by script.

## Randomness

All randomness must come from the Go CLI, not the LLM.

`dnd5e-dm roll` should support common dice notation:

```text
1d20+5
2d6+3
4d6kh3
1d20adv+7
1d20dis+2
```

Each roll should append to `roll_log.jsonl`:

```json
{
  "timestamp": "2026-05-04T20:30:00+08:00",
  "visibility": "public",
  "source": "script",
  "expression": "1d20+5",
  "reason": "Goblin attack roll against Aric",
  "rolls": [14],
  "modifier": 5,
  "total": 19
}
```

Secret DM rolls are allowed but still logged with `visibility: "dm_secret"`. User-provided physical dice results may be accepted but must be logged with `source: "user"`.

## Module Fidelity

The skill must maintain a canon boundary:

- `module_canon/`: user-provided module text, notes, or legal excerpts.
- `module_index.json`: structured index of canon locations, NPCs, clues, monsters, and chapter goals.
- `campaign_log.md`: what actually happened in play.
- `dm_improv.md`: improvised additions and rationale.

Before generating plot developments, NPC behavior, encounters, or clues, the skill should check current chapter goals and relevant module index entries. New content must be labeled as one of:

- `module_canon`
- `derived_from_module`
- `dm_improvised`
- `homebrew`

Official monster existence and stat blocks must not be guessed. If a monster stat block is not present in `monster_statblocks/`, SRD/CC references, or another user-approved local source, the skill should request or require the source, suggest an SRD-compatible alternative, or clearly mark generated content as homebrew.

## Source and Copyright Policy

The skill may use SRD-compatible or Creative Commons rules and user-provided campaign materials. It should not reproduce non-SRD copyrighted text, connect to unofficial external sources by default, or pretend to know exact non-SRD stat blocks. If the user supplies a module or stat block, the skill may index and reference it for the user's game.

Default rules-source policy: use SRD/CC plus user-provided local files only. External sources such as D&D Beyond are out of scope for the default implementation because they introduce account, copyright, and terms-of-service constraints.

## Hidden Information

The campaign vault should separate DM-only data from player-visible summaries. The skill must not reveal hidden rooms, secrets, monster abilities, unrevealed NPC motives, or future module events unless the players discover them or the user explicitly asks for DM-facing information.

## Rollback and Auditability

State-changing scripts should support append-only logs and ideally reversible operations. At minimum, each state mutation should record:

- Timestamp.
- Actor or tool.
- Previous value.
- New value.
- Reason.

This protects against accidental context drift and makes it possible to recover from mistaken rulings.

## Skill Operating Modes

The skill should choose a mode based on the user request:

1. **Campaign Setup Mode**
   - Create vault structure.
   - Ask for rules version: 2014 5e, 2024 revision, or table-specific hybrid.
   - Import party, module, and approved source references.

2. **Session Prep Mode**
   - Read module progress, party state, open hooks, and unresolved clues.
   - Prepare scenes, likely encounters, NPC reactions, and improvisation tables.

3. **Live DM Mode**
   - Narrate only player-perceivable information.
   - Track state changes in files.
   - Use scripts for dice and mechanical state.
   - End most scene prompts with “你们怎么做？” or equivalent.

4. **Combat Mode**
   - Use initiative and combat scripts.
   - Enforce action economy.
   - Show current turn, round, visible conditions, and available actions.

5. **Rules Adjudication Mode**
   - Identify whether the question is tool-enforced, source-checked, or DM judgment.
   - Cite the local source or state that the source is missing.
   - Avoid inventing exact rules text.

6. **Post-Session Mode**
   - Summarize confirmed events.
   - Update campaign log, world state, NPC state, clues, unresolved hooks, and next-session prep notes.

## Expected User Experience

The assistant should feel like a careful DM assistant, not a novelist. It should be vivid when narrating, strict when resolving mechanics, and transparent when something requires a source or a tool.

Good behavior examples:

- “这是一次 Action；你本回合的 Bonus Action 和 Reaction 仍可用。”
- “我需要先查看该怪物 stat block；当前资料库没有它，所以我不能确认官方数值。”
- “这条线索属于模组主线，玩家还没有发现。我会通过 NPC 传闻给出一个自然入口，而不是直接暴露真相。”
- “现在调用 roll 脚本，而不是直接给出骰点。”

Bad behavior examples:

- Inventing monster AC/HP while presenting it as official.
- Letting a PC take two actions because the context got long.
- Forgetting concentration or reaction usage.
- Deciding player emotions or actions.
- Creating a new plot twist that contradicts module canon without labeling it as improvised deviation.

## Testing Strategy

Initial eval prompts should cover:

1. Running a combat round with initiative, action economy, damage, and conditions.
2. A player attempting multiple actions in one turn.
3. A long-running session where HP, spell slots, and NPC attitude must persist.
4. A request to use a monster not present in the approved stat block directory.
5. A request to improvise when players skip a module clue.
6. A post-session recap that updates persistent state without revealing hidden information.

## Open Design Choices

The initial implementation should default to:

- DnD 5e 2014 unless configured otherwise.
- Zone-based positioning before grid-based tactical maps.
- SRD/user-provided sources only.
- JSON for mechanical state and Markdown for human-readable logs.

These choices can be changed later without altering the core design.
