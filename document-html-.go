package sherlock

import (
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

// loadDocument_HTML tries to mimic an ActivityPub document by parsing meta-data on
// a remote HTML page. The `data` argument is a map that may already contain some
// data, and will be updated with any new data that is discovered.
func (client *Client) loadDocument_HTML(uri string, data mapof.Any) streams.Document {

	// Retrieve the HTML document
	txn := remote.Get(uri).
		UserAgent(client.UserAgent).
		WithOptions(client.RemoteOptions...)

	if err := txn.Send(); err != nil {
		return streams.NilDocument()
	}

	// Read the response body
	body, err := txn.ResponseBody()

	if err != nil {
		return streams.NilDocument()
	}

	// Default values for Web Pages
	data[vocab.PropertyID] = uri
	data[vocab.PropertyURL] = uri
	data[vocab.PropertyType] = vocab.ObjectTypePage

	// Apply links found in the response headers
	client.applyLinks(txn, data)

	// Add JSON-LD data to the data
	client.loadDocument_JSONLD(body, data)

	// Add OpenGraph (via HTMLInfo) data to the data
	client.loadDocument_OpenGraph(uri, body, data)

	// Add Microformats2 data to the data
	client.loadDocument_MicroFormats(uri, body, data)

	// Return success!
	return streams.NewDocument(data,
		streams.WithClient(client),
		streams.WithHTTPHeader(txn.ResponseHeader()),
	)
}
