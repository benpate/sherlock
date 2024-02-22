package sherlock

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/rs/zerolog/log"
)

// Actor returns an ActivityPub Actor representation of the provided URL.
// If and ActivityPub Actor cannot be found, it attempts to create a fake one
// using RSS/Atom feeds, and MicroFormats instead.
func (client Client) loadActor(url string, config *LoadConfig) (streams.Document, error) {

	const location = "sherlock.Client.Actor"

	log.Debug().Str("loc", location).Msg("searching for: " + url)

	// RULE: Prevent too many redirects
	if config.MaximumRedirects < 0 {
		return streams.NilDocument(), derp.NewInternalError(location, "Maximum redirects exceeded", url)
	}

	// 1. Try WebFinger
	if actor := client.loadActor_WebFinger(url, config); actor.NotNil() {
		log.Debug().Str("loc", location).Msg("Found via WebFinger")
		return actor, nil
	}

	// RULE: url must begin with a valid protocol
	url = defaultHTTPS(url)

	// 2. Try ActivityStreams
	if actor := client.loadActor_ActivityStreams(url); actor.NotNil() {
		log.Debug().Str("loc", location).Msg("Found via ActivityStream")
		return actor, nil
	}

	// 3. Try RSS/Atom/JSONFeed/MicroFormats
	if actor := client.loadActor_Feed(url, config); actor.NotNil() {
		log.Debug().Str("loc", location).Msg("Found via Feed")
		return actor, nil
	}

	// 4. Abject failure. Your mother would be ashamed.
	return streams.NilDocument(), derp.NewNotFoundError(location, "Unable to load actor", url)
}
