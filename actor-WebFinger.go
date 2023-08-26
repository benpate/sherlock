package sherlock

import (
	"strings"

	"github.com/benpate/digit"
)

func (client Client) actor_WebFinger(acc *actorAccumulator) bool {

	// If the ID doesn't look like an email/username then skip this step
	if !strings.Contains(acc.url, "@") {
		return false
	}

	// Try to load the resource/account via WebFinger
	resource, err := digit.Lookup(acc.url)

	// On errors, just continue processing the pipeline
	if err != nil {
		return false
	}

	acc.links = make([]digit.Link, 0, len(resource.Links))

	// Add links to the accumulator
	for _, link := range resource.Links {
		acc.links = append(acc.links, digit.Link{
			RelationType: link.RelationType,
			MediaType:    link.MediaType,
			Href:         getRelativeURL(acc.url, link.Href),
		})
	}

	return false
}
