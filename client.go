package sherlock

import (
	"bytes"
	"mime"

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
func (client Client) Load(uri string) (streams.Document, error) {

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
	if contentType := header.Get("Content-Type"); client.isActivityStream(contentType) {
		if result, err := ParseActivityStream(&body); err == nil {
			return streams.NewDocument(result, streams.WithClient(client)), nil
		}
	}

	// Try to parse the document as HTML
	result, err := ParseHTML(uri, &body)

	if err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Error parsing HTML page")
	}

	// Populate and return the resulting document
	return streams.NewDocument(result, streams.WithClient(client), streams.WithHeader(header)), nil
}

func (client Client) isActivityStream(value string) bool {

	// ActivityStreams have their own MIME type, but we have to check some alternates, too.
	if mediaType, _, err := mime.ParseMediaType(value); err == nil {
		switch mediaType {
		case "application/activity+json":
		case "application/ld+json":
		case "application/json":
			return true
		}
	}

	return false
}
