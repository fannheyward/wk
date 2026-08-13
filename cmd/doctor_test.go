package cmd

import (
	"path/filepath"
	"testing"
)

func TestCheckDoctorEntryReportsNonWorktree(t *testing.T) {
	got := checkDoctorEntry(managedWorktree{repo: "repo", name: "feature-one", path: t.TempDir()})
	if len(got) != 1 {
		t.Fatalf("want 1 problem, got %d", len(got))
	}
	if got[0].check != "not-worktree" {
		t.Fatalf("check = %q, want not-worktree", got[0].check)
	}
}

func TestDoctorProblemsAppliesLayoutNamingRules(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "root")
	source := initTestRepo(t, base, "project")
	addTestWorktree(t, source, filepath.Join(root, "a1b2", "project"), "feature/codex")
	addTestWorktree(t, source, filepath.Join(root, "project", "expected-name"), "different-name")

	got, err := doctorProblems(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("doctorProblems() returned %d problems, want 1: %#v", len(got), got)
	}
	if got[0].repo != "project" || got[0].name != "expected-name" || got[0].check != "name-branch" {
		t.Fatalf("doctorProblems() = %#v", got)
	}
}

func TestDoctorProblemsTreatsAmbiguousPathAsCodex(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "root")
	source := initTestRepo(t, base, "a1b2")
	addTestWorktree(t, source, filepath.Join(root, "a1b2", "a1b2"), "feature/codex")

	got, err := doctorProblems(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("doctorProblems() = %#v, want no problems", got)
	}
}
