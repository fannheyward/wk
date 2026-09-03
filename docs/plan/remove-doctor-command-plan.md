# 移除 `wk doctor` 计划

## 问题

`wk doctor` 已不再需要，但当前 CLI 仍注册该命令，仓库也保留其实现、测试和用户文档。

## 决策与理由

- 从 Kong 根命令删除 `doctor` 注册，使该名称不再出现在帮助中，也不能再执行。
- 删除只由该命令使用的实现和测试，避免留下不可达代码。
- 更新 README、术语表和已有计划中的相关描述，使文档反映当前命令集合。
- 保留共享的 worktree 布局识别逻辑，因为 `ls`、`path` 和 `rm` 仍在使用。

## 实施步骤

- [x] 删除根命令中的 `Doctor` 字段。
- [x] 删除 `cmd/doctor.go`、`cmd/doctor_test.go` 及其专属 Git 辅助代码。
- [x] 删除文档中的 `wk doctor` 能力描述。
- [x] 运行格式、测试和构建检查。
- [x] 检查完整 diff 与残留引用。

## 风险与缓解

| 风险 | 缓解 |
| --- | --- |
| 误删其他命令依赖的布局识别代码 | 仅删除 `DoctorCmd` 及其私有辅助代码，保留 `managedWorktreeAt` 等共享逻辑 |
| 文档继续宣传已删除命令 | 搜索仓库中的 `doctor` 和 `wk doctor` 引用 |
| CLI 仍通过其他入口暴露命令 | 检查帮助输出，并验证调用 `wk doctor` 返回 CLI 解析错误 |

## 成功标准

1. `wk --help` 不包含 `doctor`。
2. `wk doctor` 返回 CLI 解析错误。
3. 仓库不再包含该命令的实现、测试或能力描述。
4. `go test ./...` 和 `go build ./...` 通过。

## 进度

已完成。`go test ./...` 和 `go build ./...` 通过；帮助输出不包含该命令，直接调用返回 `unexpected argument doctor`。

## 相关文件

- `cmd/root.go`
- `cmd/root_test.go`
- `cmd/doctor.go`
- `cmd/doctor_test.go`
- `internal/gitx/gitx.go`
- `internal/gitx/gitx_test.go`
- `README.md`
- `docs/GLOSSARY.md`
- `docs/plan/codex-worktree-layout-compatibility-plan.md`
- `docs/plan/wk-core-worktree-features-plan.md`
