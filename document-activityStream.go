package sherlock

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

// loadDocument_ActivityStream tries to load a remote document as an ActivityStream
// If successful, it will return a streams.Document with the appropriate metadata.
// Otherwise, it returns a nil document.
func (client *Client) loadDocument_ActivityStream(uri string) streams.Document {

	data := mapof.NewAny()

	txn := remote.Get(uri).
		UserAgent(client.UserAgent).
		Accept(vocab.ContentTypeActivityPub).
		With(client.RemoteOptions...).
		Result(&data)

	if err := txn.Send(); err != nil {
		return streams.NilDocument()
	}

	if !isActivityStream(txn.ResponseContentType()) {
		return streams.NilDocument()
	}

	return streams.NewDocument(
		data,
		streams.WithClient(client),
		streams.WithHTTPHeader(txn.ResponseHeader()),
	)
}
