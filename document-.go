package sherlock

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

// LoadDocument tries to retrieve a URL from the internet, then return it into a streams.Document.
// If the remote resource is not already an ActivityStreams document, it will attempt to convert from
// RSS, Atom, JSONFeed, and HTML MicroFormats.
func (client Client) loadDocument(config Config, url string) (streams.Document, error) {

	const location = "sherlock.Client.loadDocument"

	// RULE: url must begin with a valid protocol
	url = defaultHTTPS(url)

	// 1. If we can load the document as an ActivityStream, then there you go.
	if document := client.loadDocument_ActivityStream(config, url); document.NotNil() {
		return document, nil
	}

	// 2. If we can load the document as HTML, then that will do.
	if document := client.loadDocument_HTML(config, url); document.NotNil() {
		return document, nil
	}

	// 3. If the default value is good enough, then use that.
	// This may happen when RSS feeds have *some* information, but a website CAPTCHA
	// block us from loading more details.
	if len(config.DefaultValue) > 0 {
		return streams.NewDocument(config.DefaultValue, streams.WithClient(client)), nil
	}

	// 4. Abject failure.
	return streams.NilDocument(), derp.BadRequestError(location, "Unable to load document", url, config)
}
