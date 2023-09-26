package sherlock

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/remote"
)

func (client *Client) loadActor_Feed(uri string, config *LoadConfig) streams.Document {

	// Retrieve the URL
	txn := remote.Get(uri).
		UserAgent(client.UserAgent).
		WithOptions(client.RemoteOptions...)

	if err := txn.Send(); err != nil {
		return streams.NilDocument()
	}

	// Find and follow links in the response.
	// We're keeping this variable around because
	// it may also include links to hubs, icons, etc.
	result := client.loadActor_Links(txn, config)

	if result.NotNil() {
		return result
	}

	// 1. Try to generate an Actor from a JSON Feed
	if document := client.loadActor_Feed_JSON(txn); document.NotNil() {
		return document
	}

	// 2. Try to generate an Actor from a RSS/Atom Feed
	if document := client.loadActor_Feed_RSS(txn); document.NotNil() {
		return document
	}

	// 3. Try to generate an Actor from a HTML MicroFormats
	if document := client.loadActor_Feed_MicroFormats(txn); document.NotNil() {
		return document
	}

	// 4. Failure.
	return streams.NilDocument()
}
