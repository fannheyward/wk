// Package cmd implements the wk CLI commands.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/fannheyward/wk/internal/config"
	"github.com/fannheyward/wk/internal/gitx"
)

var codexWorktreeIDRe = regexp.MustCompile(`^[A-Za-z0-9]{4}$`)

type managedWorktree struct {
	repo, name, path string
	codex            bool
}

// CLI defines the wk command surface.
//
// wk (worktree kit) collects git worktrees from all your repos into one global
// root directory and creates them from the latest default branch with readable
// random names. The directory name and initial branch name match.
type CLI struct {
	New    NewCmd    `cmd:"" help:"Create a worktree from the latest default branch"`
	Ls     LsCmd     `cmd:"" default:"withargs" help:"List worktrees of the current repo (or all repos with --all)"`
	Path   PathCmd   `cmd:"" help:"Print the absolute path of a worktree"`
	Rm     RmCmd     `cmd:"" help:"Remove a worktree directory (its branch is kept)"`
	Doctor DoctorCmd `cmd:"" help:"Check managed worktree directories"`
}

// Execute parses args and runs the selected command.
func Execute() {
	var cli CLI
	ctx := kong.Parse(
		&cli,
		kong.Name("wk"),
		kong.Description("Manage git worktrees in a centralized directory.\n\n"+
			"wk creates <root>/<repo>/<name> and also recognizes Codex worktrees at\n"+
			"<root>/<4-char-id>/<repo>. The root defaults to\n"+
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

func managedWorktreeAt(root, repo, path string) (managedWorktree, bool) {
	rel, err := filepath.Rel(resolve(root), resolve(path))
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return managedWorktree{}, false
	}
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) != 2 {
		return managedWorktree{}, false
	}
	switch {
	case codexWorktreeIDRe.MatchString(parts[0]) && parts[1] == repo:
		return managedWorktree{repo: repo, name: parts[0], path: path, codex: true}, true
	case parts[0] == repo:
		return managedWorktree{repo: repo, name: parts[1], path: path}, true
	default:
		return managedWorktree{}, false
	}
}

func rootWorktreeAt(root, path string) managedWorktree {
	if repo, err := gitx.RepoNameAt(path); err == nil {
		if wt, ok := managedWorktreeAt(root, repo, path); ok {
			return wt
		}
	}
	return managedWorktree{
		repo: filepath.Base(filepath.Dir(path)),
		name: filepath.Base(path),
		path: path,
	}
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

// findWorktree 根据 wk 目录名或 Codex ID 定位当前仓库的 worktree。
func findWorktree(name string) (gitx.Worktree, error) {
	dir, repo, err := repoDir()
	if err != nil {
		return gitx.Worktree{}, err
	}
	root := filepath.Dir(dir)
	wts, err := gitx.Worktrees()
	if err != nil {
		return gitx.Worktree{}, err
	}
	var found *gitx.Worktree
	for _, wt := range wts {
		managed, ok := managedWorktreeAt(root, repo, wt.Path)
		if !ok || managed.name != name {
			continue
		}
		if found != nil {
			return gitx.Worktree{}, fmt.Errorf("multiple worktrees named %q under %s", name, root)
		}
		candidate := wt
		found = &candidate
	}
	if found != nil {
		return *found, nil
	}
	return gitx.Worktree{}, fmt.Errorf("no worktree named %q under %s", name, root)
}
