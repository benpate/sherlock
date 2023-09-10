package sherlock

import (
	"bytes"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/sherlock/pipe"
	"github.com/rs/zerolog/log"
)

// Actor returns an ActivityPub Actor representation of the provided URL.
// If and ActivityPub Actor cannot be found, it attempts to create a fake one
// using RSS/Atom feeds, and MicroFormats instead.
func (client Client) LoadActor(url string) (streams.Document, error) {

	acc := newActorAccumulator(url)

	steps := pipe.Steps[*actorAccumulator]{

		// Use WebFinger to translate usernames/emails into usable URLs
		client.actor_WebFinger,
		client.actor_FollowLinks,

		// Try to load ActivityPub documents first.
		client.actor_ActivityStream,

		// Try to load document as HTML and look for links in the <head>
		client.actor_GetHTTP(ContentTypeHTML, ContentTypeJSONFeed, ContentTypeRSS, ContentTypeAtom, ContentTypeXML),
		client.actor_FindLinksInHeader,
		client.actor_FindLinks,
		client.actor_FollowLinks,

		// Parse various feed formats
		client.actor_JSONFeed,
		client.actor_RSSFeed,

		// Use MicroFormats as last resort
		// client.actor_MicroFormats,
	}

	// Try to execute the pipe
	if done := pipe.Run(&acc, steps...); done {

		// Use magic values from pre-defined links
		for _, link := range acc.links {
			switch link.RelationType {

			// If a hub is defined, then use it for WebSub
			case "hub":
				acc.webSub = link.Href

			// If we don't already have an Icon, try using the icon link header
			case "icon":
				if _, ok := acc.result.GetStringOK(vocab.PropertyIcon); !ok {
					acc.result.SetString(vocab.PropertyIcon, link.Href)
				}
			}
		}

		// Use values from http response headers (et al)
		result := streams.NewDocument(
			acc.result,
			streams.WithClient(client),
			streams.WithMeta("format", acc.format),
			streams.WithMeta("cache-control", acc.cacheControl),
			streams.WithMeta("websub", acc.webSub),
		)

		// Success-a-mundo !!
		return result, nil
	}

	// Unable to load the actor. Return in shame.
	return streams.NilDocument(), derp.NewNotFoundError("sherlock.Client.Actor", "Unable to load actor", url)
}

func (client Client) actor_GetHTTP_Atom(acc *actorAccumulator) bool {
	return client.actor_GetHTTP(ContentTypeAtom)(acc)
}

func (client Client) actor_GetHTTP_RSS(acc *actorAccumulator) bool {
	return client.actor_GetHTTP(ContentTypeRSS)(acc)
}

func (client Client) actor_GetHTTP_JSONFeed(acc *actorAccumulator) bool {
	return client.actor_GetHTTP(ContentTypeJSONFeed)(acc)
}

func (client Client) actor_GetHTTP(contentTypes ...string) pipe.Step[*actorAccumulator] {
	return func(acc *actorAccumulator) bool {

		var body bytes.Buffer

		// Try to load the ID as an ActivityPub object
		txn := remote.Get(acc.url).
			Accept(contentTypes...).
			Response(&body, nil)

		// Try to send the transaction.  If successful, the populate the accumulator
		if err := txn.Send(); err == nil {
			acc.httpResponse = txn.ResponseObject
			acc.cacheControl = txn.ResponseObject.Header.Get("Cache-Control")
			acc.body = body
		}

		return false
	}
}

// This is a debugging step, so it's okay if it's not always used.
// nolint: unused
func debug(label string) pipe.Step[*actorAccumulator] {
	return func(acc *actorAccumulator) bool {
		log.Debug().
			Str("url", acc.url).
			Interface("result", acc.result).
			Interface("links", acc.links).
			Str("format", acc.format).
			Msg(label)
		return false
	}
}
