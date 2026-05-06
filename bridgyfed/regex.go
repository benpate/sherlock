package bridgyfed

import (
	"strings"

	"github.com/benpate/dns"
)

// LooksLikeBluesky returns TRUE if the provided URI looks like a Bluesky handle (e.g. "@alice.bsky.social")
func LooksLikeBluesky(uri string) bool {
	uri = strings.TrimPrefix(uri, "@")
	uri = strings.ToLower(uri)

	if !dns.IsValidHostname(uri) {
		return false
	}

	// Special case to fix autocomplete (e.g. still typing "yomama.bsky.so"....)
	if strings.HasSuffix(uri, "bsky.so") {
		return false
	}

	// Confirm that there are at least 2 segments (e.g. "me.social")
	if segments := strings.Split(uri, "."); len(segments) < 2 {
		return false
	}

	// Woot.
	return true
}
