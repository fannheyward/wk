# `wk` 抽象消融计划

## 问题

当前实现只有四个子命令，但 root 路径解析独占一个 `internal/config` 包，repo 探测有两层单调用包装，worktree 解析结果还保留运行时从不读取的字段，随机命名也暴露了只供内部调用的候选生成函数。这些设计增加了文件、导出符号和调用层级。具名 Git 操作包装作为对照候选，需通过消融判断是否有边界价值。

基线：5 个 Go 包，731 行生产 Go 代码；`go test ./...` 和构建通过。

## 架构决策与理由

- `cmd` 直接拥有 `WK_ROOT` 和默认 root 的解析。它是 CLI 策略，不需要独立配置模块。
- `gitx` 保留具名 Git 操作包装，使 `cmd` 不直接持有 Git 参数；只删除不表达操作边界的冗余层。
- repo 探测只执行一次 `RepoNameAt("")`，由 `repoDir` 保留原有友好错误。
- 删除 `Worktree.IsMain`，因为生产代码从不读取它。
- 保留 `internal/naming`。大词表和随机生成逻辑由一个操作封装，这个边界能减少 `cmd/new.go` 的细节。

## 实施步骤

- [x] 把 root 解析和对应测试移入 `cmd`，删除 `internal/config`。
- [x] 删除 `EnsureRepo` 和单调用的 `RepoName` 包装。
- [x] 试验以 `gitx.Run` 替代具名操作包装；因 Git 参数泄漏到 `cmd`，撤销该候选。
- [x] 删除 `Worktree.IsMain` 及其测试断言。
- [x] 把仅供 `Unique` 使用的随机候选生成并入 `Unique`。
- [x] 更新引用旧结构的计划和 ADR。
- [x] 运行定向测试、A/B 和 CLI 烟测，检查完整 diff。

## 风险与缓解

| 风险 | 缓解 |
| --- | --- |
| root 解析迁移改变 `WK_ROOT` 行为 | 保留默认值、绝对路径和 `~` 展开测试 |
| 合并 repo 探测后错误信息变差 | 在 `repoDir` 显式返回原有 `not inside a git repository` |
| 误删有边界价值的 Git 包装 | 保留 `Fetch`、`WorktreeAdd` 和 `WorktreePrune`，让 `cmd` 只表达工作流 |
| 清理死字段误伤 porcelain 解析 | 保留 path、branch 和 detached 断言 |

## 成功标准

1. `wk` 的公开命令和参数不变。
2. 默认 root、`WK_ROOT`、`~/` 展开和非 repo 错误行为不变。
3. Go 包少于 5 个，生产 Go 代码少于 690 行。
4. 不新增依赖或替代抽象。
5. `go test ./...`、`go build ./...` 和 `git diff --check` 通过。

## 进度

已完成。Go 包从 5 个降至 4 个，生产 Go 代码从 731 行降至 687 行。改动前后的帮助与非 repo 错误输出一致；`new`、`rm`、branch 保留和 prune 烟测通过。具名 Git 操作包装在消融后恢复。

## 相关文件

- `cmd/root.go`
- `cmd/root_test.go`
- `cmd/new.go`
- `cmd/rm.go`
- `internal/config/config.go`
- `internal/config/config_test.go`
- `internal/gitx/gitx.go`
- `internal/gitx/gitx_test.go`
- `internal/naming/naming.go`
- `internal/naming/naming_test.go`
- `docs/adr/0004-create-and-delete-policy.md`
- `docs/plan/wk-cli-plan.md`
