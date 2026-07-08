package naming

import (
	"regexp"
	"testing"
)

var nameRe = regexp.MustCompile(`^[a-z]+-[a-z]+$`)

func TestNewFormat(t *testing.T) {
	for range 100 {
		n := New()
		if !nameRe.MatchString(n) {
			t.Fatalf("New() = %q, want adjective-noun lowercase", n)
		}
	}
}

func TestUniqueAvoidsTaken(t *testing.T) {
	// Occupy one name; Unique must never return it.
	taken := New()
	for range 50 {
		got := Unique(func(n string) bool { return n == taken })
		if got == taken {
			t.Fatalf("Unique returned the taken name %q", taken)
		}
	}
}
