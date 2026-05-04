# 战役资料库结构

Campaign vault 是状态真相。进行机械或连续性敏感判断前，先读取相关文件。

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

## campaign_config.json

```json
{
  "campaign_name": "Example Campaign",
  "rules_version": "5e-2014",
  "positioning": "zones",
  "allowed_sources": ["SRD", "CC", "user_provided_module"]
}
```

## characters/<pc-name>.json

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
  "resources": {"action_surge": 1, "hit_dice_d10": 3},
  "spell_slots": {},
  "conditions": []
}
```

## combat_state.json

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

## 日志与模组事实

- `campaign_log.md`：只记录桌面上已经确认发生的事件。
- `roll_log.jsonl`：每次掷骰一个 JSON 对象。
- `module_index.json`：章节目标、关键 NPC、地点、线索和批准遭遇。
- `dm_improv.md`：记录 `derived_from_module`、`dm_improvised` 或 `homebrew`。
- `monster_statblocks/`：只存放用户批准的 stat block；缺失时不能当作官方内容编造。
