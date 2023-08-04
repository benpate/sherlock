package sherlock

import (
	"mime"
)

// isActivityStream returns TRUE if the MIME type is either activity+json or ld+json
func isActivityStream(value string) bool {

	// ActivityStreams have their own MIME type, but we have to check some alternates, too.
	if mediaType, _, err := mime.ParseMediaType(value); err == nil {
		switch mediaType {
		case "application/activity+json", "application/ld+json":
			return true
		}
	}

	return false
}
