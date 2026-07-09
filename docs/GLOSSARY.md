# 术语表 (Glossary)

本文件统一 `wk` 项目中的核心概念，避免沟通歧义。

| 术语 | 含义 |
| --- | --- |
| **wk** | 本 CLI 工具的名字（worktree kit）。 |
| **源 repo (source repo)** | 用户真实的 git 克隆仓库，worktree 依附于它。`wk new` 必须在源 repo 内运行。 |
| **worktree** | git 原生的工作树（`git worktree`）。一个源 repo 可挂多个 worktree，每个 checkout 到不同 branch，互不干扰。 |
| **root（全局根目录）** | 所有 worktree 集中存放的根目录。默认 `~/worktrees`，可用环境变量 `WK_ROOT` 覆盖。 |
| **repo 分组目录** | root 下按源 repo 名分出的子目录：`<root>/<repo名>/`。同一 repo 的所有 worktree 聚在这里。 |
| **随机名 (random name)** | 创建 worktree 时生成的"形容词-名词"风格名字（如 `brave-otter`）。**同时作为目录名和初始 branch 名**。 |
| **目录名 (dir name)** | worktree 所在目录的名字，等于创建时的随机名。**创建后永不改变**，是定位 worktree 的稳定标识 (stable key)。 |
| **branch 名** | worktree 当前 checkout 的分支名。等于创建时的随机名，与目录名始终一致（工具不提供 rename）。 |
| **默认分支 (default branch)** | 源 repo 的主分支，自动检测（`main`/`master`，通过远端 HEAD 判断）。新 worktree 的起点。 |
| **stable key** | 用户在 `path`/`rm` 中引用 worktree 的标识，统一用**目录名**（稳定不变，且与 branch 名恒等）。 |
