// Package gitx wraps the system git binary. worktree orchestration is best
// expressed by shelling out to real git rather than reimplementing it.
package gitx

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// run executes git with args in the given dir (empty = current dir) and
// returns trimmed stdout. On failure it wraps stderr into the error.
func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(errBuf.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
	}
	return strings.TrimSpace(out.String()), nil
}

// Run executes git in the current directory and returns trimmed stdout.
func Run(args ...string) (string, error) {
	return run("", args...)
}

// RepoNameAt 返回指定 worktree 所属源仓库的目录名。
func RepoNameAt(path string) (string, error) {
	common, err := run(path, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", err
	}
	// common is like /path/to/repo/.git -> repo name is basename of its parent.
	return filepath.Base(filepath.Dir(common)), nil
}

// DefaultBranch detects the source repo's default branch. It tries
// origin/HEAD first, then falls back to main / master among either the
// remote-tracking refs or local branches (so it works offline / without a remote).
func DefaultBranch() (string, error) {
	if ref, err := Run("symbolic-ref", "--quiet", "refs/remotes/origin/HEAD"); err == nil {
		if branch, ok := strings.CutPrefix(ref, "refs/remotes/origin/"); ok && RefExists(ref) {
			return branch, nil
		}
	}
	for _, b := range []string{"main", "master"} {
		if RefExists("refs/remotes/origin/"+b) || RefExists("refs/heads/"+b) {
			return b, nil
		}
	}
	return "", fmt.Errorf("cannot detect default branch (no origin/HEAD, main or master)")
}

// RefExists reports whether the given fully-qualified ref exists.
func RefExists(ref string) bool {
	_, err := Run("show-ref", "--verify", "--quiet", ref)
	return err == nil
}

// CheckBranchName returns an error if name is not a valid git branch name.
func CheckBranchName(name string) error {
	if _, err := Run("check-ref-format", "--branch", name); err != nil {
		return fmt.Errorf("invalid branch name %q: %w", name, err)
	}
	return nil
}

// StartPoint returns the ref `wk new` should branch from for the given default
// branch: origin/<branch> when available (freshest, especially after a
// successful fetch), otherwise the local <branch> (offline fallback).
func StartPoint(branch string) (string, error) {
	if RefExists("refs/remotes/origin/" + branch) {
		return "origin/" + branch, nil
	}
	if RefExists("refs/heads/" + branch) {
		return branch, nil
	}
	return "", fmt.Errorf("no start point for default branch %q (neither origin/%s nor local %s exists)", branch, branch, branch)
}

// Fetch updates remote-tracking refs from origin.
func Fetch() error {
	_, err := Run("fetch", "origin")
	return err
}

// WorktreeAdd creates a worktree at path checked out to a new branch,
// starting from startPoint (e.g. origin/main).
func WorktreeAdd(path, branch, startPoint string) error {
	_, err := Run("worktree", "add", "-b", branch, path, startPoint)
	return err
}

// WorktreeRemove removes a worktree. force allows removing a dirty worktree.
func WorktreeRemove(path string, force bool) error {
	args := []string{"worktree", "remove", path}
	if force {
		args = append(args, "--force")
	}
	_, err := Run(args...)
	return err
}

// WorktreePrune cleans up stale worktree administrative records.
func WorktreePrune() error {
	_, err := Run("worktree", "prune")
	return err
}

// BranchAt returns the checked-out branch of the worktree at path, or
// "(detached)" / "?" when it cannot be determined. Used by `wk ls --all`.
func BranchAt(path string) string {
	b, err := run(path, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "?"
	}
	if b == "HEAD" {
		return "(detached)"
	}
	return b
}

// Status describes lightweight worktree status for list output.
type Status struct {
	Dirty    string
	Upstream string
	Ahead    string
	Behind   string
}

// StatusAt returns dirty and upstream state for the worktree at path.
func StatusAt(path string) Status {
	st := Status{Dirty: "?", Upstream: "-", Ahead: "-", Behind: "-"}
	if out, err := run(path, "status", "--porcelain"); err == nil {
		if out == "" {
			st.Dirty = "clean"
		} else {
			st.Dirty = "dirty"
		}
	}
	upstream, err := run(path, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err != nil || upstream == "" {
		return st
	}
	st.Upstream = upstream
	if out, err := run(path, "rev-list", "--left-right", "--count", "HEAD...@{u}"); err == nil {
		st.Ahead, st.Behind = parseAheadBehind(out)
	} else {
		st.Ahead, st.Behind = "?", "?"
	}
	return st
}

func parseAheadBehind(out string) (ahead, behind string) {
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return "?", "?"
	}
	return fields[0], fields[1]
}

// Worktree describes one entry from `git worktree list`.
type Worktree struct {
	Path   string // absolute path
	Branch string // branch short name, or "(detached)"
}

// Worktrees parses `git worktree list --porcelain` for the current repo.
// The first entry is the main worktree.
func Worktrees() ([]Worktree, error) {
	out, err := Run("worktree", "list", "--porcelain")
	if err != nil {
		return nil, err
	}
	return parseWorktrees(out), nil
}

// parseWorktrees parses porcelain output into worktree entries.
func parseWorktrees(out string) []Worktree {
	var list []Worktree
	var cur *Worktree
	for line := range strings.SplitSeq(out, "\n") {
		if path, ok := strings.CutPrefix(line, "worktree "); ok {
			list = append(list, Worktree{Path: path})
			cur = &list[len(list)-1]
			continue
		}
		if cur == nil {
			continue
		}
		if branch, ok := strings.CutPrefix(line, "branch "); ok {
			cur.Branch = strings.TrimPrefix(branch, "refs/heads/")
		} else if line == "detached" {
			cur.Branch = "(detached)"
		}
	}
	return list
}
