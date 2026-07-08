package gitx

import "testing"

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
