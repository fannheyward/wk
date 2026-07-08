// Package config resolves the global worktree root directory.
package config

import (
	"os"
	"path/filepath"
)

// Root returns the global worktree root directory.
//
// Resolution order:
//  1. WK_ROOT environment variable (if set and non-empty)
//  2. ~/worktrees (default)
//
// The path is expanded and made absolute, but not created.
func Root() (string, error) {
	if v := os.Getenv("WK_ROOT"); v != "" {
		return expand(v)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "worktrees"), nil
}

// expand resolves a leading ~ and returns an absolute path.
func expand(p string) (string, error) {
	if p == "~" || (len(p) >= 2 && p[:2] == "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if p == "~" {
			return home, nil
		}
		p = filepath.Join(home, p[2:])
	}
	return filepath.Abs(p)
}
