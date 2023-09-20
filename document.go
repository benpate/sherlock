package sherlock

import (
	"bytes"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/remote"
)

// LoadDocument tries to retrieve a URL from the internet, then return it into a streams.Document.
// If the remote resource is not already an ActivityStreams document, it will attempt to convert from
// RSS, Atom, JSONFeed, and HTML MicroFormats.
func (client Client) LoadDocument(uri string, defaultValue map[string]any) (streams.Document, error) {

	const location = "sherlock.Client.LoadDocument"

	result := streams.NewDocument(defaultValue, streams.WithClient(client))

	if uri == "" {
		return result, derp.New(derp.CodeBadRequestError, "sherlock.Client.LoadDocument", "Empty URI")
	}

	// Load the document
	var body bytes.Buffer

	// Load the document from the URL (preferr ActivityStreams over HTML)
	transaction := remote.Get(uri).
		Header("Accept", "application/activity+json;q=1.0,text/html;q=0.9").
		Result(&body)

	// Try to retrieve the document from the remote server
	if err := transaction.Send(); err != nil {
		return result, derp.Wrap(err, location, "Error loading URL", uri)
	}

	// If Content-Type is valid, try to parse as ActivityStreams JSON
	header := transaction.ResponseRaw().Header

	result.WithOptions(
		streams.WithMeta("cache-control", header.Get("cache-control")),
		streams.WithMeta("etag", header.Get("etag")),
		streams.WithMeta("expires", header.Get("expires")),
	)

	if contentType := header.Get("Content-Type"); isActivityStream(contentType) {
		if err := ParseActivityStream(&result, &body); err == nil {
			return result, nil
		}
	}

	// Try to parse the document as HTML
	if err := Parse(uri, &body, defaultValue); err != nil {
		return result, derp.Wrap(err, location, "Error parsing HTML page")
	}

	// Populate and return the resulting document
	return result, nil
}
