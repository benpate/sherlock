package sherlock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidUsername(t *testing.T) {
	cases := map[string]bool{
		"benpate":     true,
		"ben_pate":    true,
		"BenPate99":   true,
		"abc":         true,  // exactly 3 chars (minimum)
		"ab":          false, // too short
		"":            false,
		"has space":   false,
		"has-dash":    false, // dashes are not allowed
		"has.dot":     false,
		"emoji😀":      false,
		"with@symbol": false,
	}

	for input, expected := range cases {
		assert.Equal(t, expected, IsValidUsername(input), "input: %q", input)
	}
}

func TestIsValidAddress(t *testing.T) {
	cases := map[string]bool{
		// @username@host.tld form
		"@benpate@climatejustice.social": true,
		"@ab@example.com":                false, // username too short
		"@benpate@":                      false, // empty domain
		"@benpate":                       false, // no second "@", so not a username form
		"@benpate@not a host":            false, // invalid hostname

		// URL form
		"https://example.com/@benpate": true,
		"http://example.com":           true,

		// Bare domain (would be valid with https:// prepended)
		"example.com":     true,
		"sub.example.com": true,

		// Garbage
		"":             false,
		"not valid !!": false,
		"@@@@@":        false,
	}

	for input, expected := range cases {
		assert.Equal(t, expected, IsValidAddress(input), "input: %q", input)
	}
}
