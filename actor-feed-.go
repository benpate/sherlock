package sherlock

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/remote"
)

func (client *Client) loadActor_Feed(config Config, url string) streams.Document {

	// Retrieve the URL
	txn := remote.Get(url).
		UserAgent(config.UserAgent).
		With(config.RemoteOptions...)

	if err := txn.Send(); err != nil {
		return streams.NilDocument()
	}

	// Find and follow links in the response.
	if document := client.loadActor_Links(config, txn); document.NotNil() {
		return document
	}

	// 1. Try to generate an Actor from a JSON Feed
	if document := client.loadActor_Feed_JSON(config, txn); document.NotNil() {
		return document
	}

	// 2. Try to generate an Actor from a RSS/Atom Feed
	if document := client.loadActor_Feed_RSS(config, txn); document.NotNil() {
		return document
	}

	// 3. Try to generate an Actor from a HTML MicroFormats
	if document := client.loadActor_Feed_MicroFormats(config, txn); document.NotNil() {
		return document
	}

	// 4. Failure.
	return streams.NilDocument()
}
