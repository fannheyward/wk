// Package cmd implements the wk CLI commands.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/fannheyward/wk/internal/config"
	"github.com/fannheyward/wk/internal/gitx"
)

// CLI defines the wk command surface.
//
// wk (worktree kit) collects git worktrees from all your repos into one global
// root directory, creates them with readable random names from the latest
// default branch, and lets you rename their branch afterwards.
type CLI struct {
	New    NewCmd    `cmd:"" help:"Create a worktree from the latest default branch"`
	Ls     LsCmd     `cmd:"" help:"List worktrees of the current repo (or all repos with --all)"`
	Path   PathCmd   `cmd:"" help:"Print the absolute path of a worktree"`
	Rename RenameCmd `cmd:"" help:"Rename a worktree's branch (directory stays unchanged)"`
	Rm     RmCmd     `cmd:"" help:"Remove a worktree directory (its branch is kept)"`
}

// Execute parses args and runs the selected command.
func Execute() {
	var cli CLI
	ctx := kong.Parse(
		&cli,
		kong.Name("wk"),
		kong.Description("Manage git worktrees in a centralized directory.\n\n"+
			"Worktrees are stored under <root>/<repo>/<name>. The root defaults to\n"+
			"~/worktrees and can be overridden with the WK_ROOT environment variable."),
		kong.UsageOnError(),
	)
	if err := ctx.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// repoDir returns the directory holding the current repo's worktrees:
// <root>/<repoName>. It also returns the repo name.
func repoDir() (dir, repo string, err error) {
	if err = gitx.EnsureRepo(); err != nil {
		return "", "", err
	}
	repo, err = gitx.RepoName()
	if err != nil {
		return "", "", err
	}
	root, err := config.Root()
	if err != nil {
		return "", "", err
	}
	return resolve(filepath.Join(root, repo)), repo, nil
}

// resolve canonicalizes p by resolving symlinks so it matches the paths git
// reports (git returns fully-resolved paths, e.g. /private/var on macOS).
// It falls back gracefully when p or its parent does not exist yet.
func resolve(p string) string {
	if r, err := filepath.EvalSymlinks(p); err == nil {
		return r
	}
	if parent, err := filepath.EvalSymlinks(filepath.Dir(p)); err == nil {
		return filepath.Join(parent, filepath.Base(p))
	}
	return p
}

// findWorktree locates a worktree of the current repo by its directory name.
func findWorktree(name string) (gitx.Worktree, error) {
	dir, _, err := repoDir()
	if err != nil {
		return gitx.Worktree{}, err
	}
	target := filepath.Join(dir, name)
	wts, err := gitx.Worktrees()
	if err != nil {
		return gitx.Worktree{}, err
	}
	for _, wt := range wts {
		if wt.Path == target {
			return wt, nil
		}
	}
	return gitx.Worktree{}, fmt.Errorf("no worktree named %q under %s", name, dir)
}
