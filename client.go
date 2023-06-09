package sherlock

import (
	"bytes"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/remote"
)

// Client implements the hannibal/streams.Client interface, and is used to load JSON-LD documents from remote servers.
// The sherlock client maps additional meta-data into a standard ActivityStreams document.
type Client struct{}

// NewClient returns a fully initialized Client object
func NewClient() Client {
	return Client{}
}

// Load tries to load a remote resource from the internet, and returns a streams.Document.
// This method implements the hannibal/streams.Client interface.
func (client Client) Load(uri string, defaultValue map[string]any) (streams.Document, error) {

	const location = "sherlock.Cient.Load"

	// Load the document
	var body bytes.Buffer

	// Load the document from the URL (preferr ActivityStreams over HTML)
	transaction := remote.Get(uri).
		Header("Accept", "application/activity+json;q=1.0,text/html;q=0.9").
		Response(&body, nil)

	// Try to retrieve the document from the remote server
	if err := transaction.Send(); err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error loading URL", uri)
	}

	// If Content-Type is valid, try to parse as ActivityStreams JSON
	header := transaction.ResponseObject.Header
	if contentType := header.Get("Content-Type"); isActivityStream(contentType) {
		if result, err := ParseActivityStream(&body); err == nil {
			return streams.NewDocument(result, streams.WithClient(client)), nil
		}
	}

	// Try to parse the document as HTML
	if err := Parse(uri, &body, defaultValue); err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error parsing HTML page")
	}

	// Populate and return the resulting document
	return streams.NewDocument(
		defaultValue,
		streams.WithClient(client),
		streams.WithMeta("cache-control", header.Get("cache-control")),
		streams.WithMeta("etag", header.Get("etag")),
		streams.WithMeta("expires", header.Get("expires")),
	), nil
}
