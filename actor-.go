package sherlock

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

// Actor returns an ActivityPub Actor representation of the provided URL.
// If and ActivityPub Actor cannot be found, it attempts to create a fake one
// using RSS/Atom feeds, and MicroFormats instead.
func (client Client) loadActor(uri string, config *LoadConfig) (streams.Document, error) {

	// 1. Try WebFinger
	if actor := client.loadActor_WebFinger(uri, config); actor.NotNil() {
		return actor, nil
	}

	// RULE: uri must begin with a valid protocol
	uri = defaultHTTPS(uri)

	// 2. Try ActivityStreams
	if actor := client.loadActor_ActivityStreams(uri); actor.NotNil() {
		return actor, nil
	}

	// 3. Try RSS/Atom/JSONFeed/MicroFormats
	if actor := client.loadActor_Feed(uri, config); actor.NotNil() {
		return actor, nil
	}

	// 4. Abject failure. Your mother would be ashamed.
	return streams.NilDocument(), derp.NewNotFoundError("sherlock.Client.Actor", "Unable to load actor", uri)
}
