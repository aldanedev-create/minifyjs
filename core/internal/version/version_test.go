package version

import "testing"

func TestVersionNotEmpty(t *testing.T) {
	if Version == "" {
		t.Fatal("Version must not be empty")
	}
}

func TestStringMatchesVersion(t *testing.T) {
	if String() != Version {
		t.Fatalf("String() = %q, want %q", String(), Version)
	}
}

func TestVersionLooksLikeSemver(t *testing.T) {
	// We do not parse semver here (that would drag in a regexp and
	// duplicate what the release pipeline validates). We only check
	// that the value starts with "MAJOR.MINOR.PATCH".
	v := Version
	dots := 0
	for _, r := range v {
		switch {
		case r >= '0' && r <= '9':
			// digit, fine
		case r == '.':
			dots++
			if dots > 2 {
				t.Fatalf("too many dots in %q", v)
			}
		case r == '-':
			// pre-release suffix; stop scanning
			dots = 3
		default:
			t.Fatalf("unexpected character %q in version %q", r, v)
		}
		if dots == 3 {
			break
		}
	}
	if dots != 2 {
		t.Fatalf("expected 2 dots before any suffix, got %d in %q", dots, v)
	}
}