# 仓库的 0 号入口

## 问题与决策

单仓库列表只有关联 worktree，无法通过序号返回源仓库。当前仓库存在符合既有 wk 或 Codex 布局的关联 worktree 时，在列表首行加入源仓库，编号为 `0`，名称为仓库名。关联 worktree 保留原顺序和从 `1` 开始的编号。

复用 `gitx.Worktrees()` 返回的首项定位源仓库，因此从关联 worktree 内执行 `wk 0` 也返回源仓库。首项不参与关联 worktree 过滤，避免源仓库路径恰好符合托管布局时重复列出。列表展示和数字选择使用同一组行数据；详细模式复用现有状态查询。

没有符合现有列表规则的关联 worktree 时，`wk ls` 仍输出 `no worktrees`，`wk 0` 返回越界错误。`--all` 和按名称查询保持原行为。现有 Zsh 函数支持 `0`，无需修改函数。

## 实施步骤

1. 修改单仓库列表构建、编号和序号边界，更新帮助与 README。
2. 扩展现有集成测试，覆盖首行、空列表、从 worktree 返回、目录栈及主仓库位于托管布局内的情况。
3. 执行 Go 测试、静态检查并检查完整改动。

## 风险与验收

- 仅有主仓库时不得显示 `0` 号行；越界错误仍报告关联 worktree 数量。
- 普通和详细列表首行为源仓库，`wk 0` 返回该行路径，原 worktree 序号不变。
- 从 wk 和 Codex worktree 执行 `wk 0` 返回源仓库；Zsh 切换保留 `pushd` / `popd` 行为。
- 全局列表、数字名称和非法参数的现有测试须通过。

## 进度

- 已检查列表、Git worktree 查询及 Zsh 调用链。
- 已完成列表首行、`0` 号选择、帮助与 README 更新。
- `go test ./...` 通过，包含 CLI 构建、普通与详细列表、空列表、原序号映射、全局列表和 Zsh 目录切换测试。源仓库位于托管布局内的测试通过。
- `go vet ./...`、`gofmt` 和 `git diff --check` 通过。Go 检查使用 `GOCACHE=/tmp/fannheyward-wk-go-cache GOPROXY=off`。
- 已检查全部改动，实现与验收完成。

## 相关文件

- `cmd/ls.go`
- `cmd/ls_test.go`
- `README.md`
- 原序号功能计划：`docs/plan/worktree-index-navigation-plan.md`
