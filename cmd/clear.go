package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// clearIgnored files carry no user data, so a directory holding only them
// counts as empty. Codex leaves .codex-worktree-name after removing a worktree.
var clearIgnored = map[string]bool{".DS_Store": true, ".codex-worktree-name": true}

// ClearCmd removes empty directories left under the root after worktrees are
// removed.
type ClearCmd struct{}

func (c *ClearCmd) Run() error {
	root, err := worktreeRoot()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			// Depth 2 covers <root>/<repo>/<name> and <root>/<id>/<repo>.
			if _, err := clearDir(filepath.Join(root, e.Name()), 2); err != nil {
				return err
			}
		}
	}
	return nil
}

// clearDir removes dir when it holds only clearIgnored files and empty
// directories, looking at most depth levels down. Empty subdirectories are
// removed even when dir is kept. It reports whether dir was removed.
func clearDir(dir string, depth int) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	empty := true
	for _, e := range entries {
		switch {
		case e.IsDir() && depth > 1:
			removed, err := clearDir(filepath.Join(dir, e.Name()), depth-1)
			if err != nil {
				return false, err
			}
			empty = empty && removed
		case e.Type().IsRegular() && clearIgnored[e.Name()]:
		default:
			empty = false
		}
	}
	if !empty {
		return false, nil
	}
	// os.Remove instead of os.RemoveAll: content written after the scan makes
	// the removal fail rather than get deleted.
	for _, e := range entries {
		if !e.IsDir() {
			if err := os.Remove(filepath.Join(dir, e.Name())); err != nil {
				return false, err
			}
		}
	}
	if err := os.Remove(dir); err != nil {
		return false, err
	}
	fmt.Println("removed", dir)
	return true, nil
}
