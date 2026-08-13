# 术语表 (Glossary)

本文件统一 `wk` 项目中的核心概念，避免沟通歧义。

| 术语 | 含义 |
| --- | --- |
| **wk** | 本 CLI 工具的名字（worktree kit）。 |
| **源 repo (source repo)** | 用户真实的 git 克隆仓库，worktree 依附于它。`wk new` 必须在源 repo 内运行。 |
| **worktree** | git 原生的工作树（`git worktree`）。一个源 repo 可挂多个 worktree，每个 checkout 到不同 branch，互不干扰。 |
| **root（全局根目录）** | 所有 worktree 集中存放的根目录。默认 `~/worktrees`，可用环境变量 `WK_ROOT` 覆盖。 |
| **repo 分组目录** | `wk` 布局中 root 下按源 repo 名分出的子目录：`<root>/<repo名>/`。同一 repo 的 `wk` worktree 聚在这里。 |
| **Codex worktree 布局** | Codex 当前使用的 `<root>/<4位ID>/<repo名>/`。`wk` 可查询、定位、删除和诊断该布局，但不会用它创建 worktree。 |
| **Codex worktree ID** | Codex 布局中 root 下的 4 位字母数字目录名，例如 `329b`。它是 `wk path` / `wk rm` 引用该 worktree 的 stable key。 |
| **随机名 (random name)** | 创建 worktree 时生成的"形容词-名词"风格名字（如 `brave-otter`）。**同时作为目录名和初始 branch 名**。 |
| **显式名 (explicit name)** | 用户通过 `wk new --name <name>` 指定的创建名。与随机名一样，**同时作为目录名和初始 branch 名**，且必须是单路径段和合法 git branch 名。 |
| **目录名 (dir name)** | `wk` 布局中 worktree 所在目录的名字，等于创建时的随机名。**创建后永不改变**，是定位 worktree 的 stable key。Codex 布局改用 4 位 ID。 |
| **branch 名** | worktree 当前 checkout 的分支名。`wk` 创建时它等于目录名；Codex worktree 的分支名可与 4 位 ID 不同。 |
| **默认分支 (default branch)** | 源 repo 的主分支，自动检测（`main`/`master`，通过远端 HEAD 判断）。新 worktree 的起点。 |
| **stable key** | 用户在 `path`/`rm` 中引用 worktree 的标识。`wk` 布局使用目录名，Codex 布局使用 4 位 ID。 |
| **doctor** | 非破坏性诊断命令，检查 root 下的受支持目录是否为 git worktree、branch 是否可读；只对 `wk` 布局检查目录名和 branch 名是否一致。 |
