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
