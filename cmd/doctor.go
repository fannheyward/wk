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
			problems = append(problems, checkDoctorEntry(repo.Name(), entry.Name(), path)...)
		}
	}
	return problems, nil
}

func checkDoctorEntry(repo, name, path string) []doctorProblem {
	if !gitx.IsWorktreeRootAt(path) {
		return []doctorProblem{{repo: repo, name: name, check: "not-worktree", detail: "not inside a git worktree", path: path}}
	}
	branch := gitx.BranchAt(path)
	if branch == "?" {
		return []doctorProblem{{repo: repo, name: name, check: "branch", detail: "cannot read branch", path: path}}
	}
	if branch != name {
		return []doctorProblem{{repo: repo, name: name, check: "name-branch", detail: fmt.Sprintf("branch is %s", branch), path: path}}
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
