# Setup：初始化战役、建卡与导入规则

首次使用、换新战役、缺少 JSON 状态文件、缺少 SRD/CC 资料或用户问“怎么开始”时，先按本文执行。安装依赖和 CLI 构建见 `prerequisite.md`。

## 1. 初始化 campaign vault

为每个战役创建独立目录。示例：

```bash
mkdir -p campaigns/demo/{characters,session_notes,module_canon,monster_statblocks,rules_refs}
touch campaigns/demo/roll_log.jsonl campaigns/demo/campaign_log.md campaigns/demo/dm_improv.md
```

创建 `campaign_config.json`：

```json
{
  "campaign_name": "demo",
  "rules_version": "5e-2014",
  "positioning": "zones",
  "allowed_sources": ["SRD", "CC", "user_provided"]
}
```

创建 `party.json`：

```json
{
  "party_name": "Demo Party",
  "characters": ["pc1"]
}
```

创建空的世界与模组索引：

```json
{
  "player_visible": {},
  "dm_only": {},
  "npcs": {},
  "factions": {},
  "open_hooks": []
}
```

保存为 `world_state.json`。再创建 `module_index.json`：

```json
{
  "current_chapter": "",
  "chapter_goals": [],
  "key_npcs": [],
  "key_locations": [],
  "required_clues": [],
  "approved_encounters": []
}
```

## 2. 建立角色 JSON

每个 PC 一个文件：`characters/<pc-id>.json`。至少记录会影响裁定的机械字段。

```json
{
  "id": "pc1",
  "name": "Aric",
  "level": 3,
  "classes": [{ "name": "fighter", "level": 3 }],
  "species": "human",
  "background": "soldier",
  "ac": 16,
  "hp": 24,
  "max_hp": 28,
  "temp_hp": 0,
  "abilities": {
    "str": 16,
    "dex": 12,
    "con": 14,
    "int": 10,
    "wis": 11,
    "cha": 8
  },
  "proficiency_bonus": 2,
  "skills": { "athletics": 5 },
  "saves": { "str": 5, "con": 4 },
  "passives": { "perception": 10, "investigation": 10, "insight": 10 },
  "resources": { "action_surge": 1, "hit_dice_d10": 3 },
  "spell_slots": {},
  "conditions": [],
  "equipment": ["longsword", "shield", "chain mail"]
}
```

不要把缺失字段靠记忆补齐。缺少角色卡数据时，向用户索取或标记为未知。

## 3. 初始化战斗状态

进入战斗前，用 CLI 生成或更新 `combat_state.json`：

```bash
dnd5e-dm initiative \
  --combatants '[{"id":"pc1","name":"Aric","initiative_bonus":1,"ac":16,"hp":24,"max_hp":28}]' \
  --out campaigns/demo/combat_state.json
```

战斗中所有行动经济、伤害、状态变化都通过 CLI 或等价的文件修改流程记录，不能只写在对话里。

## 4. 获取与导入 SRD/CC 规则

默认规则来源是通过 `git clone https://github.com/oldmanumby/dnd.srd.git` 获取的 SRD/Creative Commons 与用户提供本地资料。推荐流程：

1. 克隆 SRD 仓库到本地引用区，保留仓库 URL、commit hash、版本和许可说明：

```bash
git clone https://github.com/oldmanumby/dnd.srd.git campaigns/demo/rules_refs/dnd.srd
```

2. 按主题使用仓库内 Markdown/JSON，或从克隆仓库整理到 `rules_refs/`，例如：
   - `rules_refs/actions.md`
   - `rules_refs/conditions.md`
   - `rules_refs/spells-srd.md`
   - `rules_refs/monsters-srd.md`
3. 在文件开头写明来源：

```md
# Conditions

Source: https://github.com/oldmanumby/dnd.srd.git, SRD 5.1 / Creative Commons, user-approved local reference.
License: CC-BY-4.0 if applicable.
Commit: <commit-hash>
```

4. 用 CLI 搜索验证：

```bash
dnd5e-dm rules search --rules-dir campaigns/demo/rules_refs --query grapple
```

未导入的规则、怪物、法术或职业特性视为“资料缺失”。不要凭模型记忆当作官方文本复述。

## 5. 导入模组、怪物与用户资料

- 模组正文或用户摘录放入 `module_canon/`。
- 模组结构化索引写入 `module_index.json`。
- 用户批准的怪物 stat block 放入 `monster_statblocks/`。
- 即兴内容和偏离原因写入 `dm_improv.md`，并标记 `derived_from_module`、`dm_improvised` 或 `homebrew`。

## 6. Setup 完成检查

- `dnd5e-dm --help` 可运行。
- `campaign_config.json`、`party.json`、`world_state.json`、`module_index.json` 存在。
- 每个 PC 都有 `characters/<pc-id>.json`。
- `roll_log.jsonl`、`campaign_log.md`、`dm_improv.md` 存在。
- `rules_refs/` 至少包含 SRD/CC 或用户批准资料，或明确记录“尚未导入”。
- 如果准备战斗，`combat_state.json` 已由 `initiative` 或人工批准数据初始化。
