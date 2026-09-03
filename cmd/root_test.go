package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestWorktreeRoot(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	relative, err := filepath.Abs("relative-wt")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name, env, want string
	}{
		{"environment", "/tmp/custom-wt", "/tmp/custom-wt"},
		{"relative", "relative-wt", relative},
		{"default", "", filepath.Join(home, "worktrees")},
		{"home", "~", home},
		{"tilde", "~/wt", filepath.Join(home, "wt")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("WK_ROOT", test.env)
			got, err := worktreeRoot()
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("worktreeRoot() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestManagedWorktreeAt(t *testing.T) {
	root := t.TempDir()
	repo := "project"
	cases := []struct {
		name     string
		path     string
		wantName string
		ok       bool
	}{
		{"wk", filepath.Join(root, repo, "brave-otter"), "brave-otter", true},
		{"codex", filepath.Join(root, "a1B2", repo), "a1B2", true},
		{"short codex id", filepath.Join(root, "abc", repo), "", false},
		{"different repo", filepath.Join(root, "a1b2", "other"), "", false},
		{"too deep", filepath.Join(root, repo, "one", "two"), "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.ok {
				if err := os.MkdirAll(tc.path, 0o755); err != nil {
					t.Fatal(err)
				}
			}
			got, ok := managedWorktreeAt(root, repo, tc.path)
			if ok != tc.ok {
				t.Fatalf("managedWorktreeAt() ok = %v, want %v", ok, tc.ok)
			}
			if !ok {
				return
			}
			if got.repo != repo || got.name != tc.wantName || got.path != tc.path {
				t.Fatalf("managedWorktreeAt() = %#v", got)
			}
		})
	}
}

func TestManagedWorktreeAtKeepsFourCharacterRepo(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "rust", "feature-one")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	got, ok := managedWorktreeAt(root, "rust", path)
	if !ok || got.repo != "rust" || got.name != "feature-one" {
		t.Fatalf("managedWorktreeAt() = %#v, %v", got, ok)
	}
}

func TestFindWorktreeFindsCodexLayout(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "root")
	source := initTestRepo(t, base, "project")
	path := filepath.Join(root, "a1b2", "project")
	addTestWorktree(t, source, path, "feature/codex")

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(source); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldDir) })
	t.Setenv("WK_ROOT", root)

	got, err := findWorktree("a1b2")
	if err != nil {
		t.Fatal(err)
	}
	if resolve(got.Path) != resolve(path) {
		t.Fatalf("findWorktree() path = %q, want %q", got.Path, path)
	}

	addTestWorktree(t, source, filepath.Join(root, "project", "a1b2"), "feature/wk")
	if _, err := findWorktree("a1b2"); err == nil {
		t.Fatal("findWorktree() accepted duplicate names across layouts")
	}
}

func initTestRepo(t *testing.T, parent, name string) string {
	t.Helper()
	path := filepath.Join(parent, name)
	runTestGit(t, "", "init", "-b", "main", path)
	runTestGit(t, path, "config", "user.email", "wk-test@example.com")
	runTestGit(t, path, "config", "user.name", "wk test")
	if err := os.WriteFile(filepath.Join(path, "README.md"), []byte("test\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, path, "add", "README.md")
	runTestGit(t, path, "commit", "-m", "initial")
	return path
}

func addTestWorktree(t *testing.T, source, path, branch string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, source, "worktree", "add", "-b", branch, path)
}

func runTestGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
