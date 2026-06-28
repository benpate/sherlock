package tagspub

import (
	"regexp"
)

// looksLikeHashtag matches a "#tag" of letters, digits, or underscores.
var looksLikeHashtag *regexp.Regexp = regexp.MustCompile(`^(?:#)?(?:[a-zA-Z0-9_]+)$`)

// IsHashtag reports whether the URI is a "#hashtag" and, if so, returns the tag
// text with the leading "#" removed.
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
