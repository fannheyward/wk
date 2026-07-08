# ADR-0004: 创建从最新默认分支起，删除保守保留 branch

## 状态
已接受 (Accepted)

## 背景
worktree 的起点和清理策略直接关系到"是否丢代码"和"起点是否最新"。

## 决策
### 创建 (new)
- **必须在 git repo 内运行**（`repoDir` 首先 `EnsureRepo`，否则报错）。
- 自动检测源 repo 的**默认分支**（origin/HEAD → 远端或本地的 main/master）。
- 创建前尝试 `git fetch`：
  - **成功（在线）**：基于 `origin/<默认分支>` 建全新 branch，起点最新。
  - **失败（离线）**：不报错，降级为警告，回退用本地已有的起点——优先 `origin/<默认分支>`（上次 fetch 的缓存），无则本地 `<默认分支>`。
- 起点选择由 `gitx.StartPoint` 统一决定；纯本地无 remote 的 repo 也能正常创建。
- 不支持 `--from` 之类的自定义起点（MVP 收缩，避免过度设计）。

### 删除 (rm)
- **只删 worktree 目录，保留 branch**：代码仍在 branch 上，误删可恢复。
- 遵守 git 的脏检查：worktree 有未提交改动时，`git worktree remove` 默认拒绝，需 `--force`。
- 删除后跑 `git worktree prune` 清理悬挂记录。

## 理由
- 自动 fetch：Copilot 体验的核心之一就是"从最新主分支开新任务"，值得为此接受联网与短暂等待。
- 保留 branch：删除是高危操作，保守策略把"丢数据"的门槛抬高，符合安全优先原则。

## 后果
- 保留的 branch 会累积，未来可能需要 `wk prune-branches` 之类的清理命令（暂不做）。
