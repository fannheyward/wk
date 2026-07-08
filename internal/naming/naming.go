// Package naming generates readable adjective-noun random names,
// used both as the worktree directory name and the initial branch name.
package naming

import (
	"math/rand/v2"
	"os"
	"path/filepath"
)

var adjectives = []string{
	"brave", "quiet", "swift", "calm", "bright", "bold", "clever", "eager",
	"gentle", "happy", "keen", "lively", "merry", "noble", "proud", "quick",
	"rapid", "sharp", "smooth", "solid", "sunny", "warm", "wise", "witty",
	"amber", "azure", "coral", "crimson", "golden", "ivory", "jade", "scarlet",
	"silent", "cosmic", "lunar", "misty", "polar", "royal", "rustic", "vivid",
}

var nouns = []string{
	"otter", "forest", "river", "falcon", "willow", "meadow", "harbor", "canyon",
	"comet", "ember", "glacier", "lagoon", "maple", "nebula", "orchid", "pebble",
	"quartz", "raven", "summit", "thicket", "tundra", "valley", "walrus", "zephyr",
	"badger", "cedar", "dolphin", "eagle", "ferret", "gecko", "heron", "ibis",
	"jaguar", "koala", "lynx", "marten", "newt", "osprey", "puffin", "robin",
}

// New returns a random "adjective-noun" name, e.g. "brave-otter".
func New() string {
	return adjectives[rand.IntN(len(adjectives))] + "-" + nouns[rand.IntN(len(nouns))]
}

// Unique returns a random name whose directory does not yet exist under dir.
// It retries up to 20 times before giving up and returning the last candidate.
func Unique(dir string) string {
	var name string
	for i := 0; i < 20; i++ {
		name = New()
		if _, err := os.Stat(filepath.Join(dir, name)); os.IsNotExist(err) {
			return name
		}
	}
	return name
}
