package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClearRemovesOnlyEmptyDirs(t *testing.T) {
	root := t.TempDir()
	t.Setenv("WK_ROOT", root)
	dirs := []string{
		"empty-repo",
		"b25a",
		"codex-nested/project",
		"live-repo/brave-otter",
		"live-repo/leftover",
		"deep/one/two",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	files := []string{
		"b25a/.codex-worktree-name",
		"b25a/.DS_Store",
		"codex-nested/.codex-worktree-name",
		"live-repo/brave-otter/.git",
		"live-repo/.DS_Store",
		"root-file",
	}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(root, f), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := (&ClearCmd{}).Run(); err != nil {
		t.Fatal(err)
	}

	for path, want := range map[string]bool{
		"empty-repo":                 false,
		"b25a":                       false,
		"codex-nested":               false,
		"live-repo/leftover":         false,
		"live-repo/.DS_Store":        true,
		"live-repo/brave-otter/.git": true,
		"deep/one/two":               true,
		"root-file":                  true,
	} {
		_, err := os.Lstat(filepath.Join(root, path))
		if exists := err == nil; exists != want {
			t.Errorf("%s exists = %v, want %v", path, exists, want)
		}
	}
}

func TestClearMissingRoot(t *testing.T) {
	t.Setenv("WK_ROOT", filepath.Join(t.TempDir(), "missing"))
	if err := (&ClearCmd{}).Run(); err != nil {
		t.Fatal(err)
	}
}
