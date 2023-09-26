package sherlock

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

// loadActor_ActivityStreams attempts to load an ActivityStream directly from
// a uri.  If the retrieved document is not an ActivityStream, then
// this method returns a NilDocument.
func (client Client) loadActor_ActivityStreams(uri string) streams.Document {

	// Set up the transaction
	data := mapof.NewAny()
	txn := remote.Get(uri).
		UserAgent(client.UserAgent).
		Accept(ContentTypeActivityPub).
		WithOptions(client.RemoteOptions...).
		Result(&data)

	// Try to load the data from the remote server
	if err := txn.Send(); err != nil {
		return streams.NilDocument()
	}

	// If the response is not an ActivityPub document, then exit
	if !isActivityStream(txn.ResponseContentType()) {
		return streams.NilDocument()
	}

	// Otherwise, return the Actor with expected metadata
	result := streams.NewDocument(
		data,
		streams.WithClient(client),
		streams.WithHTTPHeader(txn.ResponseHeader()),
	)

	return result
}
