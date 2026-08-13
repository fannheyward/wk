package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/fannheyward/wk/internal/config"
	"github.com/fannheyward/wk/internal/gitx"
)

// DoctorCmd checks wk-managed directories without changing them.
type DoctorCmd struct{}

type doctorProblem struct {
	repo, name, check, detail, path string
}

func (c *DoctorCmd) Run() error {
	root, err := config.Root()
	if err != nil {
		return err
	}
	problems, err := doctorProblems(root)
	if err != nil {
		return err
	}
	if len(problems) == 0 {
		fmt.Println("ok")
		return nil
	}
	printDoctorProblems(problems)
	return fmt.Errorf("doctor found %d problem(s)", len(problems))
}

func doctorProblems(root string) ([]doctorProblem, error) {
	repos, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var problems []doctorProblem
	for _, repo := range repos {
		if !repo.IsDir() {
			continue
		}
		repoPath := filepath.Join(root, repo.Name())
		entries, err := os.ReadDir(repoPath)
		if err != nil {
			problems = append(problems, doctorProblem{repo: repo.Name(), check: "read-dir", detail: err.Error(), path: repoPath})
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			path := filepath.Join(repoPath, entry.Name())
			problems = append(problems, checkDoctorEntry(rootWorktreeAt(root, path))...)
		}
	}
	return problems, nil
}

func checkDoctorEntry(wt managedWorktree) []doctorProblem {
	if !gitx.IsWorktreeRootAt(wt.path) {
		return []doctorProblem{{repo: wt.repo, name: wt.name, check: "not-worktree", detail: "not inside a git worktree", path: wt.path}}
	}
	branch := gitx.BranchAt(wt.path)
	if branch == "?" {
		return []doctorProblem{{repo: wt.repo, name: wt.name, check: "branch", detail: "cannot read branch", path: wt.path}}
	}
	if !wt.codex && branch != wt.name {
		return []doctorProblem{{repo: wt.repo, name: wt.name, check: "name-branch", detail: fmt.Sprintf("branch is %s", branch), path: wt.path}}
	}
	return nil
}

func printDoctorProblems(problems []doctorProblem) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "REPO\tNAME\tCHECK\tDETAIL\tPATH")
	for _, p := range problems {
		fmt.Fprintln(w, strings.Join([]string{p.repo, p.name, p.check, p.detail, p.path}, "\t"))
	}
	w.Flush()
}
