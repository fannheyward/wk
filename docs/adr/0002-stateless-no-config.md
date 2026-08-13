# ADR-0002: 无状态设计，无配置文件

## 状态
已接受 (Accepted)

## 背景
Copilot App 维护 repo 列表和 session 状态。CLI 是否也要维护一份注册表/配置？

## 决策
**完全无状态**：不维护 repo 注册表，不写配置文件。
- 源 repo 通过运行时 cwd 向上找 `.git` 确定。
- 全局 root 默认 `~/worktrees`，仅通过环境变量 `WK_ROOT` 覆盖。
- worktree 的"真相来源 (source of truth)"是**文件系统 + git 自身的 worktree 记录**，不是额外的元数据文件。
- 识别两种目录布局：`wk` 的 `<root>/<repo名>/<name>`，以及 Codex 当前的 `<root>/<4位ID>/<repo名>`。

## 理由
- 当前需要配置的东西极少（只有一个 root 路径），一个环境变量足够。
- 无状态 = 无同步问题、无脏数据、无损坏恢复逻辑。删 worktree 目录不会留下悬挂记录（靠 `git worktree prune` 兜底）。
- `wk ls` 直接扫描 root 的两层目录并结合 `git worktree list` 即可重建全貌。

## 后果
- 不能"离开源 repo 用名字直接建 worktree"——`wk new` 必须在源 repo 内跑。这是刻意的范围收缩。
- 若未来配置项变多，再引入 `~/.config/wk/config.toml`（届时新开 ADR）。
