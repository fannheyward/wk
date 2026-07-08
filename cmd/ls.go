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

// LsCmd lists worktrees of the current repo, or all repos with --all.
type LsCmd struct {
	All bool `short:"a" help:"List worktrees across all repos."`
}

func (c *LsCmd) Run() error {
	if c.All {
		return lsAllRepos()
	}
	return lsCurrentRepo()
}

// row is one line of ls output.
type row struct {
	repo, name, branch, path string
}

func lsCurrentRepo() error {
	dir, repo, err := repoDir()
	if err != nil {
		return err
	}
	wts, err := gitx.Worktrees()
	if err != nil {
		return err
	}
	var rows []row
	for _, wt := range wts {
		// Only wk-managed worktrees live under <root>/<repo>/; this also
		// filters out the source repo's main worktree.
		if filepath.Dir(wt.Path) == dir {
			rows = append(rows, row{repo, filepath.Base(wt.Path), wt.Branch, wt.Path})
		}
	}
	printRows(rows, false)
	return nil
}

func lsAllRepos() error {
	root, err := config.Root()
	if err != nil {
		return err
	}
	repos, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var rows []row
	for _, r := range repos {
		if !r.IsDir() {
			continue
		}
		repoPath := filepath.Join(root, r.Name())
		wts, err := os.ReadDir(repoPath)
		if err != nil {
			continue
		}
		for _, w := range wts {
			if !w.IsDir() {
				continue
			}
			path := filepath.Join(repoPath, w.Name())
			rows = append(rows, row{r.Name(), w.Name(), gitx.BranchAt(path), path})
		}
	}
	printRows(rows, true)
	return nil
}

func printRows(rows []row, withRepo bool) {
	if len(rows) == 0 {
		fmt.Fprintln(os.Stderr, "no worktrees")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	header := "NAME\tBRANCH\tPATH"
	if withRepo {
		header = "REPO\t" + header
	}
	fmt.Fprintln(w, header)
	for _, r := range rows {
		line := strings.Join([]string{r.name, r.branch, r.path}, "\t")
		if withRepo {
			line = r.repo + "\t" + line
		}
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
