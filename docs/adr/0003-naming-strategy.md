# ADR-0003: 命名策略——目录名随机且不可变，branch 可 rename

## 状态
已接受 (Accepted)

## 背景
Copilot App 里目录名随机、branch 名由 AI 按任务语义生成。纯 CLI 无 agent，需要重新定义两者关系。

## 决策
1. **创建时**：生成一个"形容词-名词"随机名（如 `brave-otter`），**同时**用作目录名和初始 branch 名，二者一致。
2. **目录名不可变**：任何时候都不改目录名（避免正在使用的 shell/编辑器路径断链）。
3. **branch 可改**：`wk rename <目录名> <新branch名>` 只改 git branch，不动目录。改后目录名与 branch 名 diverge，这是允许且预期的。
4. **stable key = 目录名**：`path`/`rename`/`rm` 都用目录名引用 worktree（因为它稳定不变）。
5. **随机名风格**：形容词-名词，可读好记，冲突时重试生成。

## 理由
- 目录不变是硬约束：CLI 无法感知谁正 cd 在里面或哪个编辑器开着它，改目录会造成隐性断链。
- 创建时目录名 = branch 名，降低"两套名字"的心智负担；需要语义 branch 时再显式 rename，是渐进的。
- 用目录名而非 branch 名做 key，因为 branch 会 diverge，目录名才是唯一稳定锚点。

## 后果
- `wk ls` 必须**分两列**展示目录名和 branch 名（可能不同）。
- 用户看到 diverge 是正常现象，文档需说明。
