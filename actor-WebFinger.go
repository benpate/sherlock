package sherlock

import (
	"strings"

	"github.com/benpate/digit"
)

func (client Client) actor_WebFinger(acc *actorAccumulator) {

	// If the ID doesn't look like an email/username then skip this step
	if !strings.Contains(acc.url, "@") {
		return
	}

	// Try to load the resource/account via WebFinger
	resource, err := digit.Lookup(acc.url)

	// On errors, just continue processing the pipeline
	if err != nil {
		return
	}

	// Save the links into the accumulator
	acc.links = resource.Links
}
