# wk — worktree kit

A tiny CLI that collects your git worktrees into one central directory and
creates them from the latest default branch with readable random names. The
directory name and branch name match, so commands can refer to a worktree by
that name. Inspired by the worktree flow in the GitHub Copilot app.

## Install

```sh
go install github.com/fannheyward/wk@latest
```

Requires `git` on your `PATH`.

## Concepts

Worktrees created by `wk` live under a global root, grouped by repo:

```
<root>/<repo>/<random-name>/
```

It also recognizes the directory layout currently created by Codex:

```
<root>/<4-char-id>/<repo>/
```

- **root** defaults to `~/worktrees`; override with the `WK_ROOT` environment
  variable.
- The **directory name** is a random `adjective-noun` (e.g. `brave-otter`) and
  the new **branch** uses the same name, so the directory and branch always
  match. Commands refer to a worktree by its **directory name**.
- Codex worktrees use the four-character ID (e.g. `329b`) as the name accepted
  by `wk path` and `wk rm`. Their branch names are independent from that ID.
- `wk new` continues to create only the first layout; Codex compatibility does
  not change existing creation behavior.

## Usage

Run inside any git repository.

```sh
# Create a worktree from the latest origin/<default-branch>
wk new
#   fetches origin, then creates <root>/<repo>/<random-name>
wk new --name feature-one
#   creates <root>/<repo>/feature-one on branch feature-one

# List this repo's wk and Codex worktrees (name, branch, path)
wk ls
wk ls --verbose     # include dirty/upstream/ahead/behind status
wk ls --all          # across all repos under the root

# Check the global root for invalid or inconsistent worktree directories
wk doctor

# Print a worktree's absolute path (handy for cd)
cd "$(wk path brave-otter)"
cd "$(wk path 329b)"       # Codex worktree

# Remove a worktree directory (its branch is kept)
wk rm brave-otter
wk rm brave-otter --force   # also remove when there are uncommitted changes
```

## Design notes

`wk` is intentionally stateless: no config file, no registry. The source of
truth is the filesystem plus git's own worktree records. See `docs/adr/` for the
architecture decisions and `docs/GLOSSARY.md` for terminology.

- `wk new` must run inside a git repo. It tries to `git fetch` for the freshest
  start point, but works offline too: on fetch failure it warns and falls back
  to the local default branch.
- `wk new --name <name>` lets you choose the directory and branch name at
  creation time. Names are single path segments and must also be valid git
  branch names.
- `wk ls`, `wk path`, `wk rm`, and `wk doctor` recognize both directory
  layouts. `doctor` applies the directory-name/branch-name equality rule only
  to worktrees created in the `wk` layout.
- `wk rm` keeps the branch so work can be recovered; it relies on git's dirty
  check and refuses to delete uncommitted changes without `--force`.
