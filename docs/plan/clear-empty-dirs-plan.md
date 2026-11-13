# `wk clear` 清理空目录计划

## 问题

worktree 移除后，root 下会留下空目录：`wk` 布局留下 `<root>/<repo>/`，Codex 布局留下 `<root>/<4位ID>/`。Codex 目录里通常还剩一个 0 字节的 `.codex-worktree-name`，Finder 打开过的目录会有 `.DS_Store`。

## 决策与理由

- 新增 `wk clear`，扫描整个 root，不要求在 git repo 内运行。
- 只删除空目录。目录里只有 `.DS_Store`、`.codex-worktree-name` 和空子目录时视为空，因为这两个文件不含用户数据。其他任何文件、符号链接或非空子目录都保留整个目录。
- 只向下检查两层，与两种布局的深度一致，不进入 worktree 内部遍历文件。
- 先确认整个目录为空，再逐个删除文件，最后用 `os.Remove` 删除目录。不用 `os.RemoveAll`，扫描后若目录被写入新内容，删除会失败而不是丢数据。
- 不加 `--dry-run`：只删空目录，删除的内容不含用户数据。

## 实施步骤

- [x] 写测试覆盖空目录、只含忽略文件、嵌套空目录、含真实文件、超出深度和 root 不存在。
- [x] 实现 `ClearCmd` 并注册到根命令。
- [x] 更新 README 和 ADR-0004。
- [x] 运行 `go test ./...`、`go vet ./...`、`gofmt`，用临时 root 手动验证。

## 风险与缓解

| 风险 | 缓解 |
| --- | --- |
| 误删含用户数据的目录 | 白名单只含两个元数据文件名，其他条目一律保留目录 |
| 扫描与删除之间目录被写入 | 用 `os.Remove` 删除单个文件和目录，非空目录删除失败 |
| 遍历大型 worktree 变慢 | 深度限制为两层 |

## 成功标准

1. `wk clear` 删除 root 下的空目录和只含忽略文件的目录，并逐个输出被删路径。
2. 含其他内容的目录保持不变。
3. `go test ./...` 通过。

## 进度

已完成。`go test ./...`、`go vet ./...`、`gofmt -l` 通过；临时 root 手动验证通过。

## 相关文件

- `cmd/clear.go`
- `cmd/clear_test.go`
- `cmd/root.go`
- `README.md`
- `docs/adr/0004-create-and-delete-policy.md`
