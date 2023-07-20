package sherlock

import (
	"bytes"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/remote"
	"github.com/benpate/sherlock/pipe"
	"github.com/davecgh/go-spew/spew"
)

// Actor returns an ActivityPub Actor representation of the provided URL.
// If and ActivityPub Actor cannot be found, it attempts to create a fake one
// using RSS/Atom feeds, and MicroFormats instead.
func (client Client) LoadActor(url string) (streams.Document, error) {

	spew.Dump("sherlock.LoadActor")

	acc := newActorAccumulator(url)

	steps := pipe.Steps[*actorAccumulator]{

		// Use WebFinger to translate usernames/emails into usable URLs
		client.actor_WebFinger,
		client.actor_FollowLinks,

		// Try to load ActivityPub documents first.
		client.actor_ActivityStream,

		// Try to load document as HTML and look for links in the <head>
		client.actor_GetHTTP(ContentTypeHTML, ContentTypeJSONFeed, ContentTypeRSS, ContentTypeAtom, ContentTypeXML),
		client.actor_FindLinks,
		client.actor_FollowLinks,

		// Parse various feed formats
		client.actor_JSONFeed,
		client.actor_AtomFeed,
		client.actor_RSSFeed,
		client.actor_XMLFeed,

		// Use MicroFormats as last resort
		// client.actor_MicroFormats,
	}

	// Try to execute the pipe
	pipe.Run(&acc, steps...)

	// Handle execution errors
	if err := acc.Error(); err != nil {
		return streams.NilDocument(), derp.Wrap(err, "sherlock.Client.Actor", "Error loading actor", url)
	}

	// If we have a complete result, then we're done!
	if acc.Complete() {
		header := acc.httpResponse.Header
		result := streams.NewDocument(
			acc.result,
			streams.WithClient(client),
			streams.WithMeta("format", acc.format),
			streams.WithMeta("content-type", header.Get("Content-Type")),
			streams.WithMeta("etag", header.Get("ETag")),
			streams.WithMeta("last-modified", header.Get("Last-Modified")),
			streams.WithMeta("cache-control", header.Get("Cache-Control")),
		)

		// Apply other metadata to the document.
		result.MetaAdd(acc.meta)

		// Success-a-mundo !!
		return result, nil
	}

	// Unable to load the actor. Return in shame.
	return streams.NilDocument(), derp.NewNotFoundError("sherlock.Client.Actor", "Unable to load actor", url)
}

func (client Client) actor_GetHTTP_Atom(acc *actorAccumulator) {
	client.actor_GetHTTP(ContentTypeAtom)(acc)
}

func (client Client) actor_GetHTTP_RSS(acc *actorAccumulator) {
	client.actor_GetHTTP(ContentTypeRSS)(acc)
}

func (client Client) actor_GetHTTP_JSONFeed(acc *actorAccumulator) {
	client.actor_GetHTTP(ContentTypeJSONFeed)(acc)
}

func (client Client) actor_GetHTTP(contentTypes ...string) pipe.Step[*actorAccumulator] {
	return func(acc *actorAccumulator) {

		var body bytes.Buffer

		// Try to load the ID as an ActivityPub object
		txn := remote.Get(acc.url).
			Accept(contentTypes...).
			Response(&body, nil)

		// Try to send the transaction.  If successful, the populate the accumulator
		if err := txn.Send(); err == nil {
			acc.httpResponse = txn.ResponseObject
			acc.body = body
		}
	}
}

func debugAccumulator(label string) pipe.Step[*actorAccumulator] {
	return func(acc *actorAccumulator) {
		spew.Dump(label+" -----------------------", acc.url, acc.meta, acc.result, acc.links)
	}
}
