package naming

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

var nameRe = regexp.MustCompile(`^[a-z]+-[a-z]+$`)

func TestNewFormat(t *testing.T) {
	for i := 0; i < 100; i++ {
		n := New()
		if !nameRe.MatchString(n) {
			t.Fatalf("New() = %q, want adjective-noun lowercase", n)
		}
	}
}

func TestUniqueAvoidsExisting(t *testing.T) {
	dir := t.TempDir()
	// Occupy a name; Unique must never return it.
	taken := New()
	if err := os.Mkdir(filepath.Join(dir, taken), 0o755); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 50; i++ {
		if got := Unique(dir); got == taken {
			t.Fatalf("Unique returned the taken name %q", taken)
		}
	}
}
