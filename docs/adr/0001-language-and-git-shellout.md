# ADR-0001: 用 Go 实现，git 操作 shell out

## 状态
已接受 (Accepted)

## 背景
需要一个跨平台、易分发的 worktree 管理 CLI。候选：Go、Rust、Node/TS、Bash。

## 决策
用 **Go + kong** 实现，git 操作直接 shell out 调用系统 `git`。

## 理由
- 单二进制、交叉编译，分发最简单（无运行时依赖）。
- kong 用结构体 tag 声明命令，轻量、零样板，适合本项目这种简单 CLI（早期用过 cobra，但其体量对 5 个子命令的场景偏重，已替换为 kong）。
- worktree 逻辑本质是编排 git 命令，shell out 比用 go-git 之类的库更贴近用户真实 git 行为（尤其 worktree、fetch 的边角情况），也更简单。

## 后果
- 依赖用户环境已装 `git`（合理假设）。
- 需要处理 shell out 的错误捕获与输出解析。
