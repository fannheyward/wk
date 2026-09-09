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

func TestRepoNameAt(t *testing.T) {
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
}

func TestDefaultBranchFallsBackFromStaleOriginHead(t *testing.T) {
	repo := filepath.Join(t.TempDir(), "repo")
	runTestGit(t, "", "init", "-b", "main", repo)
	runTestGit(t, repo, "config", "user.email", "wk-test@example.com")
	runTestGit(t, repo, "config", "user.name", "wk test")
	runTestGit(t, repo, "commit", "--allow-empty", "-m", "initial")
	runTestGit(t, repo, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/deleted")
	t.Chdir(repo)

	branch, err := DefaultBranch()
	if err != nil {
		t.Fatal(err)
	}
	if branch != "main" {
		t.Fatalf("DefaultBranch() = %q, want main", branch)
	}
}

func runTestGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), "git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}
