package cmd

import "fmt"

// PathCmd prints the absolute path of a worktree, e.g. cd $(wk path brave-otter).
type PathCmd struct {
	Name string `arg:"" help:"Worktree directory name."`
}

func (c *PathCmd) Run() error {
	wt, err := findWorktree(c.Name)
	if err != nil {
		return err
	}
	fmt.Println(wt.Path)
	return nil
}
