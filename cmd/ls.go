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
	All     bool `short:"a" help:"List worktrees across all repos."`
	Verbose bool `short:"v" help:"Include dirty and upstream status."`
}

func (c *LsCmd) Run() error {
	if c.All {
		return lsAllRepos(c.Verbose)
	}
	return lsCurrentRepo(c.Verbose)
}

// row is one line of ls output.
type row struct {
	repo, name, branch, dirty, upstream, ahead, behind, path string
}

func lsCurrentRepo(verbose bool) error {
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
			rows = append(rows, makeRow(repo, filepath.Base(wt.Path), wt.Branch, wt.Path, verbose))
		}
	}
	printRows(rows, false, verbose)
	return nil
}

func lsAllRepos(verbose bool) error {
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
			rows = append(rows, makeRow(r.Name(), w.Name(), gitx.BranchAt(path), path, verbose))
		}
	}
	printRows(rows, true, verbose)
	return nil
}

func makeRow(repo, name, branch, path string, verbose bool) row {
	r := row{repo: repo, name: name, branch: branch, path: path}
	if verbose {
		st := gitx.StatusAt(path)
		r.dirty = st.Dirty
		r.upstream = st.Upstream
		r.ahead = st.Ahead
		r.behind = st.Behind
	}
	return r
}

func printRows(rows []row, withRepo, verbose bool) {
	if len(rows) == 0 {
		fmt.Fprintln(os.Stderr, "no worktrees")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	header := "NAME\tBRANCH\tPATH"
	if verbose {
		header = "NAME\tBRANCH\tDIRTY\tUPSTREAM\tAHEAD\tBEHIND\tPATH"
	}
	if withRepo {
		header = "REPO\t" + header
	}
	fmt.Fprintln(w, header)
	for _, r := range rows {
		cols := []string{r.name, r.branch, r.path}
		if verbose {
			cols = []string{r.name, r.branch, r.dirty, r.upstream, r.ahead, r.behind, r.path}
		}
		line := strings.Join(cols, "\t")
		if withRepo {
			line = r.repo + "\t" + line
		}
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
