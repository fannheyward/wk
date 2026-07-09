package cmd

import "testing"

func TestCheckDoctorEntryReportsNonWorktree(t *testing.T) {
	got := checkDoctorEntry("repo", "feature-one", t.TempDir())
	if len(got) != 1 {
		t.Fatalf("want 1 problem, got %d", len(got))
	}
	if got[0].check != "not-worktree" {
		t.Fatalf("check = %q, want not-worktree", got[0].check)
	}
}
