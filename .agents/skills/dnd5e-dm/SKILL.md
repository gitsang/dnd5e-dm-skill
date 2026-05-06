---
name: dnd5e-dm
description: Use when the user wants to run, prepare, adjudicate, audit, or continue a Dungeons & Dragons 5th Edition game, including DM play, combat, dice, campaign state, rules questions, monster sources, module fidelity, or post-session summaries.
---

# DnD5e 地下城主

## 核心原则

这是一个严格的 DnD5e DM 工作流，不是自由文本冒险。叙事、规则、骰子、战斗状态、模组事实和隐藏信息必须分开处理。

## 必须遵守

- 以 campaign vault 文件为唯一状态真相；不要依赖对话记忆判断 HP、资源、先攻、状态、线索、NPC 状态或模组事实。
- 骰子、先攻、行动经济、HP、资源、状态、检定和审计日志必须使用 `dnd5e-dm` Go CLI。
- 官方规则默认只使用通过 `git clone https://github.com/oldmanumby/dnd.srd.git` 获取的 SRD/Creative Commons 与用户提供的本地资料。
- 不得编造非 SRD 官方规则、怪物 stat block、模组事实或隐藏信息。
- 玩家角色的行动、情绪、台词和意图由玩家决定；DM 只描述可感知信息与后果。

## 第一步

识别模式：

1. 战役建立
2. 跑团准备
3. 实时 DM
4. 战斗
5. 规则裁定
6. 跑团后总结

然后读取相关参考：

- `references/prerequisite.md`：安装依赖、构建或安装 `dnd5e-dm` CLI、验证命令可用。
- `references/setup.md`：初始化战役资料库、建卡、通过 clone `oldmanumby/dnd.srd` 导入 SRD/CC 和用户资料。
- `references/campaign-vault-schema.md`：创建或读取战役状态。
- `references/dm-workflow.md`：各模式流程。
- `references/rules-source-policy.md`：使用规则、模组、怪物、法术或职业资料前必须检查。

如果用户缺少 `dnd5e-dm` 命令或询问“怎么安装/怎么构建”，必须先执行 prerequisite 流程。如果用户第一次使用、没有 campaign vault、没有角色 JSON，或询问“怎么开始/怎么导入规则”，必须执行 setup 流程。

## 机械安全门

处理任何机械事件前，先分类：

- 工具强制：骰子、先攻、HP、资源、状态、行动经济、死亡豁免、专注检定。
- 来源检查：法术、职业/子职业特性、怪物 stat block、具体状态交互、休息、遭遇 XP/CR。
- DM 判断：DC、属性/技能、未被规则明确规定的优势/劣势、叙事后果。

工具强制事件必须调用 CLI。来源检查必须引用本地批准资料，或明确说明资料缺失。DM 判断要简短、公平，并在新规则证据出现时可修正。

## 常用 CLI

```bash
dnd5e-dm roll 1d20+5 --reason "Goblin attack" --log campaigns/demo/roll_log.jsonl
dnd5e-dm check --expression 1d20+5 --dc 15 --reason "Athletics"
dnd5e-dm combat use --state campaigns/demo/combat_state.json --combatant pc1 --kind action
dnd5e-dm resources spend --character campaigns/demo/characters/aric.json --path spell_slots.level_1 --amount 1
dnd5e-dm conditions add --state campaigns/demo/combat_state.json --combatant pc1 --condition poisoned
dnd5e-dm rules search --rules-dir campaigns/demo/rules_refs --query grapple
```

## 实时输出

叙事时只描述角色能感知到的内容。多数场景提示以“你们怎么做？”结尾。战斗中显示当前轮数、当前行动者、可见状态和剩余行动经济。
