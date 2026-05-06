# 规则来源政策

- 默认允许来源：通过 `git clone https://github.com/oldmanumby/dnd.srd.git` 获取的 SRD、Creative Commons 规则、用户提供角色卡、用户提供模组、用户批准的 stat block。
- 不内置、不复述、不编造非 SRD/CC 官方全文或 stat block。
- 默认不接入 D&D Beyond 等外部账号来源；这类来源涉及登录、版权和服务条款。
- 怪物、法术、职业特性或模组事实不在本地批准来源中时，必须说明资料缺失。
- 获取 SRD/CC 时必须 clone `https://github.com/oldmanumby/dnd.srd.git`，保留来源 URL、commit hash、版本和许可说明；导入后放入 `rules_refs/` 并用 `dnd5e-dm rules search` 验证可检索。
- 新内容必须标记为 `module_canon`、`derived_from_module`、`dm_improvised` 或 `homebrew`。
- 隐藏模组事实在玩家发现前保持 DM-only，除非用户明确要求 DM-facing 信息。
