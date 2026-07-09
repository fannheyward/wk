package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fannheyward/wk/internal/gitx"
	"github.com/fannheyward/wk/internal/naming"
)

// NewCmd creates a new worktree under <root>/<repo>/<name>.
//
// It detects the repo's default branch and tries to fetch origin for the
// freshest start point. Offline is fine: on fetch failure it warns and falls
// back to the local default branch. Must be run inside a git repository.
type NewCmd struct {
	Name string `help:"Use this directory and branch name instead of a random name."`
}

var worktreeNameRe = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

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

	name, err := chooseWorktreeName(dir, c.Name)
	if err != nil {
		return err
	}
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

func chooseWorktreeName(dir, explicit string) (string, error) {
	taken := func(n string) bool {
		if _, err := os.Stat(filepath.Join(dir, n)); err == nil {
			return true
		}
		return gitx.RefExists("refs/heads/" + n)
	}
	if explicit == "" {
		// A name is taken if its directory exists or a branch of that name already
		// exists (wk rm keeps branches, so a freed directory may still have one).
		return naming.Unique(taken), nil
	}
	if err := validateWorktreeName(explicit); err != nil {
		return "", err
	}
	if taken(explicit) {
		return "", fmt.Errorf("worktree name %q already exists as a directory or branch", explicit)
	}
	return explicit, nil
}

func validateWorktreeName(name string) error {
	switch {
	case name == "":
		return fmt.Errorf("worktree name cannot be empty")
	case name == "." || name == "..":
		return fmt.Errorf("worktree name %q is not allowed", name)
	case strings.ContainsAny(name, `/\`):
		return fmt.Errorf("worktree name %q must be a single path segment", name)
	case !worktreeNameRe.MatchString(name):
		return fmt.Errorf("worktree name %q may only contain letters, digits, '.', '_' and '-'", name)
	}
	return gitx.CheckBranchName(name)
}
