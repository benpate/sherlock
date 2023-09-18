package sherlock

import (
	"strings"
)

// defaultHTTPS is a pipe.Step that prepends "https://" to the URL if it doesn't already exist.
func defaultHTTPS(acc *actorAccumulator) bool {

	if !strings.HasPrefix(acc.url, "https://") {
		acc.url = "https://" + acc.url
	}

	return false
}
