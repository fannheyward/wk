package cmd

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestIndexedNavigation(t *testing.T) {
	readme, err := os.ReadFile("../README.md")
	if err != nil {
		t.Fatal(err)
	}
	_, shellBlock, ok := strings.Cut(string(readme), "```zsh\nwk() {\n")
	if !ok {
		t.Fatal("README shell function not found")
	}
	shellFunction, _, ok := strings.Cut(shellBlock, "\n```")
	if !ok {
		t.Fatal("README shell function code block is not closed")
	}
	shellFunction = "wk() {\n" + shellFunction

	base := resolve(t.TempDir())
	binDir := filepath.Join(base, "bin")
	if err := os.Mkdir(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(binDir, "wk")
	build := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "..")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build wk: %v\n%s", err, out)
	}
	root := filepath.Join(base, "root space")
	source := initTestRepo(t, base, "project")
	other := initTestRepo(t, base, "other")
	codexPath := filepath.Join(root, "a1b2", "project")
	numericPath := filepath.Join(root, "project", "1")
	wkPath := filepath.Join(root, "project", "brave-otter")
	otherPath := filepath.Join(root, "other", "global")
	addTestWorktree(t, source, wkPath, "brave-otter")
	addTestWorktree(t, source, numericPath, "numeric")
	addTestWorktree(t, source, codexPath, "feature/codex")
	addTestWorktree(t, other, otherPath, "global")
	t.Setenv("WK_ROOT", root)
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	run := func(args ...string) (string, string, error) {
		t.Helper()
		command := exec.Command(binary, args...)
		command.Dir = source
		var stdout, stderr bytes.Buffer
		command.Stdout, command.Stderr = &stdout, &stderr
		err := command.Run()
		return stdout.String(), stderr.String(), err
	}
	current := []row{
		{repo: "project", name: "a1b2", branch: "feature/codex", path: resolve(codexPath)},
		{repo: "project", name: "1", branch: "numeric", path: resolve(numericPath)},
		{repo: "project", name: "brave-otter", branch: "brave-otter", path: resolve(wkPath)},
	}
	global := []row{current[0], {repo: "other", name: "global", branch: "global", path: otherPath}, current[1], current[2]}
	for _, tc := range []struct {
		name    string
		args    []string
		rows    []row
		all     bool
		verbose bool
	}{
		{name: "default", rows: current},
		{name: "ls", args: []string{"ls"}, rows: current},
		{name: "verbose", args: []string{"ls", "-v"}, rows: current, verbose: true},
		{name: "all", args: []string{"ls", "--all"}, rows: global, all: true},
		{name: "all short", args: []string{"ls", "-a"}, rows: global, all: true},
		{name: "all verbose", args: []string{"ls", "--all", "--verbose"}, rows: global, all: true, verbose: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, stderr, err := run(tc.args...)
			if err != nil {
				t.Fatalf("wk %v: %v\n%s", tc.args, err, stderr)
			}
			want := "NAME BRANCH "
			if tc.verbose {
				want += "DIRTY UPSTREAM AHEAD BEHIND "
			}
			want += "PATH"
			if tc.all {
				want = "REPO " + want
			} else {
				want = "# " + want
			}
			for i, r := range tc.rows {
				want += " "
				if tc.all {
					want += r.repo + " "
				} else {
					want += strconv.Itoa(i+1) + " "
				}
				want += r.name + " " + r.branch + " "
				if tc.verbose {
					want += "clean - - - "
				}
				want += r.path
			}
			if got := strings.Join(strings.Fields(out), " "); got != want {
				t.Fatalf("wk %v output = %q, want %q", tc.args, got, want)
			}
		})
	}
	for i, r := range current {
		index := strconv.Itoa(i + 1)
		out, stderr, err := run(index)
		if err != nil || strings.TrimSpace(out) != r.path {
			t.Fatalf("wk %s = %q, %v, stderr %q; want %q", index, out, err, stderr, r.path)
		}
	}
	for _, args := range [][]string{{"1", "--all"}, {"1", "-a"}, {"--all", "1"}, {"ls", "1", "--all"}} {
		out, stderr, err := run(args...)
		if err == nil || out != "" || !strings.Contains(stderr, "cannot be used with --all") {
			t.Fatalf("wk %v: stdout %q, stderr %q, error %v", args, out, stderr, err)
		}
	}
	for _, index := range []string{"0", "5", "-1", "999999999999999999999", "doctor", "init"} {
		out, stderr, err := run(index)
		if err == nil || strings.Contains(out, root) || stderr == "" {
			t.Fatalf("wk %s: stdout %q, stderr %q, error %v", index, out, stderr, err)
		}
	}
	if out, stderr, err := run("path", "1"); err != nil || strings.TrimSpace(out) != resolve(numericPath) {
		t.Fatalf("wk path 1: %q, %v, stderr %q", out, err, stderr)
	}
	t.Run("empty list", func(t *testing.T) {
		t.Setenv("WK_ROOT", filepath.Join(base, "missing"))
		out, stderr, err := run("1")
		if err == nil || out != "" || !strings.Contains(stderr, "0 worktrees") {
			t.Fatalf("wk 1: stdout %q, stderr %q, error %v", out, stderr, err)
		}
		out, stderr, err = run("ls", "--all")
		if err != nil || out != "" || stderr != "" {
			t.Fatalf("wk ls --all: stdout %q, stderr %q, error %v", out, stderr, err)
		}
	})

	t.Run("zsh", func(t *testing.T) {
		shellPath, err := exec.LookPath("zsh")
		if err != nil {
			t.Skipf("zsh unavailable: %v", err)
		}
		script := "set -eu\n" + shellFunction + `
test "$(wk ls)" = "$(command wk ls)"
test "$(wk ls --all)" = "$(command wk ls --all)"
test "$(wk path 1)" = "$2"
original=$PWD
wk 1
test "$PWD" = "$1"
test "$dirstack[1]" = "$original"
builtin popd >/dev/null
test "$PWD" = "$original"
wk 2
test "$PWD" = "$2"
wk 3
test "$PWD" = "$3"
before="$PWD"
if wk 1 --all; then exit 1; fi
if wk 1 -a; then exit 1; fi
if wk 0; then exit 1; fi
if wk 999; then exit 1; fi
test "$PWD" = "$before"
wk 1 --help >/dev/null
wk 1 -h >/dev/null
if wk 1 extra; then exit 1; fi
if wk 1 --all extra; then exit 1; fi
wk --help >/dev/null
wk >/dev/null
test "$PWD" = "$before"
cd "$4"
test "$(wk ls --all)" = "$(command wk ls --all)"
if wk 1 --all; then exit 1; fi
test "$PWD" = "$4"
`
		command := exec.Command(shellPath, "-f", "-c", script, "zsh", resolve(codexPath), resolve(numericPath), resolve(wkPath), base)
		command.Dir = source
		if out, err := command.CombinedOutput(); err != nil {
			t.Fatalf("zsh navigation: %v\n%s", err, out)
		}
	})
}
