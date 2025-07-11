package sherlock

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/rs/zerolog/log"
)

// Actor returns an ActivityPub Actor representation of the provided URL.
// If and ActivityPub Actor cannot be found, it attempts to create a fake one
// using RSS/Atom feeds, and MicroFormats instead.
func (client Client) loadActor(identifier string, config *LoadConfig) (streams.Document, error) {

	const location = "sherlock.Client.Actor"

	log.Trace().Str("loc", location).Str("identifier", identifier).Msg("Loading Actor")

	// RULE: Prevent too many redirects
	if config.MaximumRedirects < 0 {
		return streams.NilDocument(), derp.InternalError(location, "Maximum redirects exceeded", identifier)
	}

	// Validate the identifier
	idType := identifierType(identifier)

	if idType == IdentifierTypeNone {
		return streams.NilDocument(), derp.BadRequestError(location, "Invalid identifier", identifier)
	}

	log.Trace().Str("loc", location).Str("type", idType).Msg("searching for: " + identifier)

	// 1. If this looks like a username, then try WebFinger
	if idType == IdentifierTypeUsername {

		if actor := client.loadActor_WebFinger(identifier, config); actor.NotNil() {
			log.Trace().Str("loc", location).Msg("Found via WebFinger")
			return actor, nil
		}

		// If we can't look up the user via WebFinger, then stop here
		return streams.NilDocument(), derp.NotFoundError(location, "Unable to load actor by username", identifier)
	}

	// RULE: identifier must begin with a valid protocol
	identifier = defaultHTTPS(identifier)

	// 2. Try ActivityStreams
	if actor := client.loadActor_ActivityStreams(identifier); actor.NotNil() {
		log.Trace().Str("loc", location).Msg("Found via ActivityStream")
		return actor, nil
	}

	// 3. Try RSS/Atom/JSONFeed/MicroFormats
	if actor := client.loadActor_Feed(identifier, config); actor.NotNil() {
		log.Trace().Str("loc", location).Msg("Found via Feed")
		return actor, nil
	}

	// 4. Abject failure. Your mother would be ashamed.
	return streams.NilDocument(), derp.NotFoundError(location, "Unable to load actor", identifier)
}
