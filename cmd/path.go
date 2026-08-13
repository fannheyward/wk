package cmd

import "fmt"

// PathCmd 输出 worktree 的绝对路径，例如 cd $(wk path brave-otter)。
type PathCmd struct {
	Name string `arg:"" help:"Worktree name (wk directory name or Codex ID)."`
}

func (c *PathCmd) Run() error {
	wt, err := findWorktree(c.Name)
	if err != nil {
		return err
	}
	fmt.Println(wt.Path)
	return nil
}
