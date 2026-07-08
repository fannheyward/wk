package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootFromEnv(t *testing.T) {
	t.Setenv("WK_ROOT", "/tmp/custom-wt")
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	if got != "/tmp/custom-wt" {
		t.Errorf("Root() = %q, want /tmp/custom-wt", got)
	}
}

func TestRootDefault(t *testing.T) {
	t.Setenv("WK_ROOT", "")
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, "worktrees")
	if got != want {
		t.Errorf("Root() = %q, want %q", got, want)
	}
}

func TestRootEnvTildeExpand(t *testing.T) {
	t.Setenv("WK_ROOT", "~/wt")
	got, err := Root()
	if err != nil {
		t.Fatal(err)
	}
	home, _ := os.UserHomeDir()
	want := filepath.Join(home, "wt")
	if got != want {
		t.Errorf("Root() = %q, want %q", got, want)
	}
}
