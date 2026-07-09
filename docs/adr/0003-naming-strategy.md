# ADR-0003: 命名策略——目录名与 branch 名始终一致

## 状态
已接受 (Accepted)（原"branch 可 rename"方案已废弃）

## 背景
Copilot App 里目录名随机、branch 名由 AI 按任务语义生成。纯 CLI 无 agent，需要重新定义两者关系。早期版本提供过 `wk rename` 让 branch 与目录分叉，后判定为不必要的复杂度而移除。

## 决策
1. **创建时**：生成一个"形容词-名词"随机名（如 `brave-otter`），**同时**用作目录名和 branch 名，二者一致。
2. **目录名不可变**：任何时候都不改目录名（避免正在使用的 shell/编辑器路径断链）。
3. **不提供 branch rename**：目录名与 branch 名永远一致，不支持分叉。需要语义 branch 时用户可自行 `git branch -m`，工具不介入。
4. **stable key = 目录名**：`path`/`rm` 用目录名引用 worktree。
5. **随机名风格**：形容词-名词，可读好记；冲突（目录或同名 branch 已存在）时重试生成。

## 理由
- 目录不变是硬约束：CLI 无法感知谁正 cd 在里面或哪个编辑器开着它，改目录会造成隐性断链。
- 目录名 = branch 名始终成立，心智负担最低，`ls` 两列永远相同、无歧义。
- 移除 rename 去掉了"目录/branch 分叉"这一整类边界情形，简化实现与文档。

## 后果
- `wk ls` 仍分两列展示目录名和 branch 名，但二者恒等。
- `wk new` 生成随机名时须同时避开已存在的目录**和** branch（`wk rm` 保留 branch，可能残留同名 branch）。
