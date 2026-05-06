# Prerequisite：安装依赖与 CLI

缺少 `dnd5e-dm` 命令、缺少 Go、首次安装本 skill，或用户询问“怎么安装/怎么构建/依赖是什么”时，先按本文执行。完成后再进入 `setup.md` 初始化战役。

## 1. 安装依赖

必须有 Go 工具链才能构建 CLI。先检查：

```bash
go version
```

如果 `go` 不存在，先让用户安装 Go；不要改用 LLM 掷骰或手写随机结果替代 CLI。

## 2. 构建 CLI

在本仓库根目录构建：

```bash
cd cli
go build -o ../.agents/skills/dnd5e-dm/bin/dnd5e-dm ./cmd/dnd5e-dm
```

## 3. 可选：安装到 PATH

如果用户希望全局可用，把二进制复制到 PATH 中的目录：

```bash
mkdir -p ~/.local/bin
cp .agents/skills/dnd5e-dm/bin/dnd5e-dm ~/.local/bin/dnd5e-dm
```

## 4. 验证

```bash
dnd5e-dm --help
```

验证失败时，先修复 CLI 安装或 PATH；不要继续进入需要骰子、战斗、资源或规则搜索的流程。
