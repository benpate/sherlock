package tagspub

import (
	"regexp"
)

var looksLikeHashtag *regexp.Regexp = regexp.MustCompile(`^(?:#)?(?:[a-zA-Z0-9_]+)$`)

func IsHashtag(uri string) (bool, string) {

	bytes := []byte(uri)

	// RULE: bounds check
	if len(bytes) == 0 {
		return false, ""
	}

	// Quick check: Must start with # symbol
	if bytes[0] != '#' {
		return false, ""
	}

	// Thorough Regexp pattern match.
	if !looksLikeHashtag.Match(bytes) {
		return false, ""
	}

	// Return the rest of the value without the # symbol
	return true, string(bytes[1:])
}
