package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fannheyward/wk/internal/gitx"
	"github.com/fannheyward/wk/internal/naming"
)

// NewCmd creates a new worktree under <root>/<repo>/<random-name>.
//
// It detects the repo's default branch and tries to fetch origin for the
// freshest start point. Offline is fine: on fetch failure it warns and falls
// back to the local default branch. Must be run inside a git repository.
type NewCmd struct{}

func (c *NewCmd) Run() error {
	dir, repo, err := repoDir()
	if err != nil {
		return err
	}

	branch, err := gitx.DefaultBranch()
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "fetching origin...")
	if err := gitx.Fetch(); err != nil {
		// Offline is fine: fall back to the local default branch.
		fmt.Fprintln(os.Stderr, "warning: fetch failed, using local branch:", err)
	}

	startPoint, err := gitx.StartPoint(branch)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	name := naming.Unique(dir)
	path := filepath.Join(dir, name)

	if err := gitx.WorktreeAdd(path, name, startPoint); err != nil {
		return err
	}

	fmt.Printf("created worktree %s\n", name)
	fmt.Printf("  repo:   %s\n", repo)
	fmt.Printf("  branch: %s (from %s)\n", name, startPoint)
	fmt.Printf("  path:   %s\n", path)
	return nil
}
