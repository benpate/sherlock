package bridgyfed

import (
	"strings"

	"github.com/benpate/uri"
)

// LooksLikeBluesky returns TRUE if the provided URI looks like a Bluesky handle (e.g. "@alice.bsky.social")
func LooksLikeBluesky(id string) bool {
	id = strings.TrimPrefix(id, "@")
	id = strings.ToLower(id)

	if uri.NotValidHostname(id) {
		return false
	}

	// Special case to fix autocomplete (e.g. still typing "yomama.bsky.so"....)
	if strings.HasSuffix(id, "bsky.so") {
		return false
	}

	// Confirm that there are at least 2 segments (e.g. "me.social")
	if segments := strings.Split(id, "."); len(segments) < 2 {
		return false
	}

	// Woot.
	return true
}
