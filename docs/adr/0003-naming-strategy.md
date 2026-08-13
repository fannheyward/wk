# ADR-0003: wk 命名策略——目录名与 branch 名始终一致

## 状态
已接受 (Accepted)（原"branch 可 rename"方案已废弃；Codex 兼容范围已补充）

## 背景
Copilot App 里目录名随机、branch 名由 AI 按任务语义生成。纯 CLI 无 agent，需要重新定义两者关系。早期版本提供过 `wk rename` 让 branch 与目录分叉，后判定为不必要的复杂度而移除。

## 决策
1. **`wk` 创建时**：生成一个"形容词-名词"随机名（如 `brave-otter`），**同时**用作目录名和 branch 名，二者一致。
2. **`wk` 目录名不可变**：任何时候都不改目录名（避免正在使用的 shell/编辑器路径断链）。
3. **`wk` 不提供 branch rename**：目录名与 branch 名永远一致，不支持分叉。需要语义 branch 时用户可自行 `git branch -m`，工具不介入。
4. **`wk` stable key = 目录名**：`path`/`rm` 用目录名引用 worktree。
5. **`wk` 随机名风格**：形容词-名词，可读好记；冲突（目录或同名 branch 已存在）时重试生成。
6. **Codex 兼容例外**：上述目录名与 branch 名恒等约束只适用于 `wk` 创建的 worktree。对于 `<root>/<4位ID>/<repo名>`，4 位 ID 是 stable key，branch 名由 Codex 独立决定。
7. **同名歧义优先 Codex**：若 4 位 repo 名、Codex ID 和第二层 repo 名完全相同（例如 `<root>/rust/rust>`），按 Codex 布局解释，避免对 Codex 的语义 branch 套用 `wk` 命名约束。

## 理由
- 目录不变是硬约束：CLI 无法感知谁正 cd 在里面或哪个编辑器开着它，改目录会造成隐性断链。
- 目录名 = branch 名始终成立，心智负担最低，`ls` 两列永远相同、无歧义。
- 移除 rename 去掉了"目录/branch 分叉"这一整类边界情形，简化实现与文档。

## 后果
- `wk ls` 仍分两列展示 stable key 和 branch 名；`wk` 布局二者恒等，Codex 布局允许不同。
- `wk new` 生成随机名时须同时避开已存在的目录**和** branch（`wk rm` 保留 branch，可能残留同名 branch）。
- `wk new` 仍只创建 `<root>/<repo名>/<name>`，不会创建 Codex 布局。
- 极少数传统 `<root>/<四位repo名>/<同名worktree>` 路径会按 Codex 布局解释，不执行目录名与 branch 名一致检查；查询和删除路径不受影响。
