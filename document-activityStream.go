package sherlock

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

// loadDocument_ActivityStream tries to load a remote document as an ActivityStream
// If successful, it will return a streams.Document with the appropriate metadata.
// Otherwise, it returns a nil document.
func (client *Client) loadDocument_ActivityStream(config Config, uri string) (streams.Document, error) {

	const location = "sherlock.client.loadDocument_ActivityStream"

	data := mapof.NewAny()

	// config.RemoteOptions = append(
	//	config.RemoteOptions,
	//  options.Debug(),
	// )

	txn := remote.Get(uri).
		UserAgent(config.UserAgent).
		Accept(vocab.ContentTypeActivityPub).
		With(config.RemoteOptions...).
		Result(&data)

	if err := txn.Send(); err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Unable to load remote document", uri)
	}

	if !isActivityStream(txn.ResponseContentType()) {
		return streams.NilDocument(), nil
	}

	return streams.NewDocument(
		data,
		streams.WithClient(client),
		streams.WithHTTPHeader(txn.ResponseHeader()),
	), nil
}
