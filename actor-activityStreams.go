package sherlock

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/remote"
	"github.com/benpate/remote/options"
	"github.com/benpate/rosetta/mapof"
	"github.com/rs/zerolog/log"
)

// loadActor_ActivityStreams attempts to load an ActivityStream directly from
// a uri.  If the retrieved document is not an ActivityStream, then
// this method returns a NilDocument.
func (client Client) loadActor_ActivityStreams(uri string) streams.Document {

	const location = "sherlock.Client.loadActor_ActivityStreams"

	// Set up the transaction
	data := mapof.NewAny()
	txn := remote.Get(uri).
		UserAgent(client.UserAgent).
		Accept(ContentTypeActivityPub).
		With(client.RemoteOptions...).
		Result(&data)

	if canTrace() {
		txn.With(options.Debug())
	}

	// Try to load the data from the remote server
	if err := txn.Send(); err != nil {
		log.Trace().Str("location", location).Msg("Error loading URI: " + uri)
		return streams.NilDocument()
	}

	// If the response is not an ActivityPub document, then exit
	if !isActivityStream(txn.ResponseContentType()) {
		if canTrace() {
			log.Trace().Str("location", location).Msg("Response is not an ActivityStream: " + txn.ResponseContentType())
		}
		return streams.NilDocument()
	}

	if canTrace() {
		log.Trace().Str("location", location).Str("objectId", uri).Msg("Found ActivityStreams document")
	}

	// Otherwise, return the Actor with expected metadata
	result := streams.NewDocument(
		data,
		streams.WithClient(client),
		streams.WithHTTPHeader(txn.ResponseHeader()),
	)

	return result
}
