package cmd

import (
	"fmt"

	"github.com/fannheyward/wk/internal/gitx"
)

// RmCmd removes a worktree directory but keeps its branch, so work can be
// recovered later. git refuses to remove a dirty worktree unless --force.
type RmCmd struct {
	Force bool   `short:"f" help:"Remove even if the worktree has uncommitted changes."`
	Name  string `arg:"" help:"Worktree name (wk directory name or Codex ID)."`
}

func (c *RmCmd) Run() error {
	wt, err := findWorktree(c.Name)
	if err != nil {
		return err
	}
	if err := gitx.WorktreeRemove(wt.Path, c.Force); err != nil {
		return fmt.Errorf("%w\n(worktree may have uncommitted changes; rerun with --force to remove anyway)", err)
	}
	if err := gitx.WorktreePrune(); err != nil {
		return err
	}
	fmt.Printf("removed worktree %s (branch %s kept)\n", c.Name, wt.Branch)
	return nil
}
