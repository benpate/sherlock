package sherlock

import (
	"strings"

	"github.com/benpate/digit"
	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/streams"
	"github.com/rs/zerolog/log"
)

// loadActor_WebFinger resolves a @user@host.tld identifier via WebFinger, then
// loads the linked ActivityPub actor or profile feed. Returns a nil document if
// no usable link is found.
func (client *Client) loadActor_WebFinger(config Config, uri string) streams.Document {

	const location = "sherlock.Client.loadActor_WebFinger"

	// If the ID doesn't look like an email/username then skip this step
	if !strings.Contains(uri, "@") {
		log.Trace().Str("location", location).Msg("Skipping because uri doesn't look like an email address")
		return streams.NilDocument()
	}

	// Try to load the Actor via WebFinger
	response, err := digit.Lookup(uri, config.remoteOptions()...)

	// If we don't have a valid response, then return nil (skip this step)
	if err != nil {
		log.Error().Err(err).Msg("loadActor_WebFinger: skipping because of error")
		return streams.NilDocument()
	}

	log.Trace().Str("location", location).Interface("response", response).Msg("Found WebFinger response")

	// RULE: Resolving the handle is itself a hop, so spend one unit of the
	// redirect budget before following any of its links. Config travels by
	// value, so this must happen before the calls below, not after them.
	config.MaximumRedirects--

	// Search for ActivityPub endpoints
	for _, link := range response.Links {
		if (link.RelationType == digit.RelationTypeSelf) && (hannibal.IsActivityPubContentType(link.MediaType)) {
			if result := client.loadActor_ActivityStreams(config, link.Href); result.NotNil() {
				return result
			}
		}
	}

	// Search for Profile pages (as a backup)
	for _, link := range response.Links {
		if link.RelationType == digit.RelationTypeProfile {
			if result := client.loadActor_Feed(config, link.Href); result.NotNil() {
				return result
			}
		}
	}

	// Fall through means we couldn't find any relevant links in the WebFinger response
	return streams.NilDocument()
}
