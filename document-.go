package sherlock

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

// LoadDocument tries to retrieve a URL from the internet, then return it into a streams.Document.
// If the remote resource is not already an ActivityStreams document, it will attempt to convert from
// RSS, Atom, JSONFeed, and HTML MicroFormats.
func (client Client) loadDocument(uri string, config LoadConfig) (streams.Document, error) {

	const location = "sherlock.Client.loadDocument"

	// RULE: uri must not be empty
	if uri == "" {
		return streams.NilDocument(), derp.NewBadRequestError("sherlock.Client.LoadDocument", "Empty URI")
	}

	// RULE: uri must begin with a valid protocol
	uri = defaultHTTPS(uri)

	// 1. If we can load the document as an ActivityStream, then there you go.
	if document := client.loadDocument_ActivityStream(uri); document.NotNil() {
		return document, nil
	}

	// 2. If we can load the document as HTML, then that will do.
	if document := client.loadDocument_HTML(uri, config.DefaultValue); document.NotNil() {
		return document, nil
	}

	// 3. Abject failure.
	return streams.NilDocument(), derp.NewBadRequestError(location, "Unable to load document", uri)
}
