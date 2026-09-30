package sherlock

import (
	"testing"
)

// The fuzz tests below assert only that the parsers/decoders never panic on
// arbitrary input. They are deliberately lenient about output: the contract for
// these functions on garbage input is "degrade gracefully", not "return X".

func FuzzIsValidAddress(f *testing.F) {
	f.Add("@user@example.com")
	f.Add("https://example.com/@user")
	f.Add("example.com")
	f.Add("")
	f.Add("@@@@")

	f.Fuzz(func(t *testing.T, input string) {
		// Two calls should agree (no hidden state / non-determinism).
		if first, second := IsValidAddress(input), IsValidAddress(input); first != second {
			t.Errorf("IsValidAddress(%q) is not deterministic", input)
		}
	})
}

func FuzzIsValidUsername(f *testing.F) {
	f.Add("benpate")
	f.Add("ab")
	f.Add("")
	f.Add("with space")

	f.Fuzz(func(t *testing.T, input string) {
		if first, second := IsValidUsername(input), IsValidUsername(input); first != second {
			t.Errorf("IsValidUsername(%q) is not deterministic", input)
		}
	})
}
