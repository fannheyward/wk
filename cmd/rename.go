package cmd

import (
	"fmt"

	"github.com/fannheyward/wk/internal/gitx"
)

// RenameCmd renames a worktree's branch. Only the branch is renamed; the
// worktree directory keeps its original name so existing shells and editors
// don't break. After this the directory and branch names may differ.
type RenameCmd struct {
	Name      string `arg:"" help:"Worktree directory name."`
	NewBranch string `arg:"" name:"new-branch" help:"New branch name."`
}

func (c *RenameCmd) Run() error {
	wt, err := findWorktree(c.Name)
	if err != nil {
		return err
	}
	if wt.Branch == "" || wt.Branch == "(detached)" {
		return fmt.Errorf("worktree %q has no branch to rename", c.Name)
	}
	if err := gitx.BranchMove(wt.Branch, c.NewBranch); err != nil {
		return err
	}
	fmt.Printf("renamed branch %s -> %s (directory %s unchanged)\n", wt.Branch, c.NewBranch, c.Name)
	return nil
}
