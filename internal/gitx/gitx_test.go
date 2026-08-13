package gitx

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestParseWorktrees(t *testing.T) {
	out := `worktree /repo/main
HEAD abc123
branch refs/heads/main

worktree /root/repo/brave-otter
HEAD def456
branch refs/heads/fix-login

worktree /root/repo/lone-wolf
HEAD 789abc
detached
`
	got := parseWorktrees(out)
	if len(got) != 3 {
		t.Fatalf("want 3 worktrees, got %d", len(got))
	}
	if !got[0].IsMain {
		t.Errorf("first entry should be main worktree")
	}
	if got[1].IsMain {
		t.Errorf("second entry should not be main")
	}
	cases := []struct {
		i            int
		path, branch string
	}{
		{0, "/repo/main", "main"},
		{1, "/root/repo/brave-otter", "fix-login"},
		{2, "/root/repo/lone-wolf", "(detached)"},
	}
	for _, c := range cases {
		if got[c.i].Path != c.path {
			t.Errorf("entry %d path = %q, want %q", c.i, got[c.i].Path, c.path)
		}
		if got[c.i].Branch != c.branch {
			t.Errorf("entry %d branch = %q, want %q", c.i, got[c.i].Branch, c.branch)
		}
	}
}

func TestParseWorktreesEmpty(t *testing.T) {
	if got := parseWorktrees(""); len(got) != 0 {
		t.Errorf("empty input should yield no worktrees, got %d", len(got))
	}
}

func TestParseAheadBehind(t *testing.T) {
	cases := []struct {
		out           string
		ahead, behind string
	}{
		{"0\t0", "0", "0"},
		{"3 12", "3", "12"},
		{"bad", "?", "?"},
	}
	for _, c := range cases {
		ahead, behind := parseAheadBehind(c.out)
		if ahead != c.ahead || behind != c.behind {
			t.Errorf("parseAheadBehind(%q) = %q, %q; want %q, %q", c.out, ahead, behind, c.ahead, c.behind)
		}
	}
}

func TestIsWorktreeRootAtRejectsNestedDirectory(t *testing.T) {
	repo := filepath.Join(t.TempDir(), "repo")
	if err := os.Mkdir(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := exec.Command("git", "init", "-b", "main", repo).Run(); err != nil {
		t.Fatal(err)
	}
	name, err := RepoNameAt(repo)
	if err != nil {
		t.Fatal(err)
	}
	if name != "repo" {
		t.Fatalf("RepoNameAt() = %q, want repo", name)
	}
	nested := filepath.Join(repo, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if !IsWorktreeRootAt(repo) {
		t.Fatalf("repo root should be a worktree root")
	}
	if IsWorktreeRootAt(nested) {
		t.Fatalf("nested directory should not be a worktree root")
	}
}
