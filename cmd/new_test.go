package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateWorktreeName(t *testing.T) {
	cases := []struct {
		name string
		ok   bool
	}{
		{"feature-one", true},
		{"fix_123", true},
		{"release.1", true},
		{"", false},
		{".", false},
		{"..", false},
		{"feature/one", false},
		{"feature one", false},
		{"foo..bar", false},
	}
	for _, c := range cases {
		err := validateWorktreeName(c.name)
		if c.ok && err != nil {
			t.Errorf("validateWorktreeName(%q) returned error: %v", c.name, err)
		}
		if !c.ok && err == nil {
			t.Errorf("validateWorktreeName(%q) succeeded, want error", c.name)
		}
	}
}

func TestChooseWorktreeNameRejectsExistingDirectory(t *testing.T) {
	dir := t.TempDir()
	name := "feature-one"
	if err := os.Mkdir(filepath.Join(dir, name), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := chooseWorktreeName(dir, name); err == nil {
		t.Fatalf("chooseWorktreeName accepted existing directory %q", name)
	}
}

func TestChooseWorktreeNameRejectsExistingCodexDirectory(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "root")
	source := initTestRepo(t, base, "project")
	dir := filepath.Join(root, "project")
	name := "a1b2"
	addTestWorktree(t, source, filepath.Join(root, name, "project"), "feature/codex")
	t.Chdir(source)
	if _, err := chooseWorktreeName(dir, name); err == nil {
		t.Fatalf("chooseWorktreeName accepted existing Codex worktree %q", name)
	}
}

func TestChooseWorktreeNameAllowsOtherFourCharacterRepo(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "root")
	current := initTestRepo(t, base, "project")
	other := initTestRepo(t, base, "a1b2")
	addTestWorktree(t, other, filepath.Join(root, "a1b2", "project"), "feature/other")
	t.Chdir(current)

	got, err := chooseWorktreeName(filepath.Join(root, "project"), "a1b2")
	if err != nil {
		t.Fatal(err)
	}
	if got != "a1b2" {
		t.Fatalf("chooseWorktreeName() = %q, want a1b2", got)
	}
}
