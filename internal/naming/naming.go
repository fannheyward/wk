// Package naming generates readable adjective-noun random names,
// used both as the worktree directory name and the initial branch name.
package naming

import "math/rand/v2"

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

// Unique returns a random name for which taken(name) reports false, retrying
// up to 20 times before giving up and returning the last candidate. The caller
// decides what "taken" means (e.g. an existing directory or git branch).
func Unique(taken func(name string) bool) string {
	var name string
	for range 20 {
		name = adjectives[rand.IntN(len(adjectives))] + "-" + nouns[rand.IntN(len(nouns))]
		if !taken(name) {
			return name
		}
	}
	return name
}
