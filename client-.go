package sherlock

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/rs/zerolog/log"
)

// Client implements the hannibal/streams.Client interface, and is used to load JSON-LD documents from remote servers.
// The sherlock client maps additional meta-data into a standard ActivityStreams document.
type Client struct {
	userAgent   string
	keyPairFunc KeyPairFunc
}

// NewClient returns a fully initialized Client object
func NewClient(options ...ClientOption) Client {
	result := Client{
		userAgent: "Sherlock (https://github.com/benpate/sherlock)",
	}

	result.With(options...)
	return result
}

func (client *Client) With(options ...ClientOption) {
	for _, option := range options {
		option(client)
	}
}

func (client Client) SetRootClient(rootClient streams.Client) {
	// NO-OP: There is no inner client to receive the root pointer
}

// Load retrieves a document from a remote server and returns it as a streams.Document
// It uses either the "Actor" or "Document" methods of generating it ActivityStreams
// result.
// "Document" treats the URL as a single ActivityStreams document, translating
// OpenGraph, MicroFormats, and JSON-LD into an ActivityStreams equivalent.
// "Actor" treats the URL as an Actor, translating RSS, Atom, JSON, and
// MicroFormats feeds into an ActivityStream equivalent.
func (client Client) Load(url string, options ...any) (streams.Document, error) {

	const location = "sherlock.Client.Load"

	log.Trace().Str("loc", location).Msg("Loading " + url)

	config := client.newConfig(options...)

	// RULE: url must not be empty
	if url == "" {
		return streams.NilDocument(), derp.BadRequestError(location, "URL cannot be empty")
	}

	// RULE: Prevent too many redirects
	if config.MaximumRedirects < 0 {
		return streams.NilDocument(), derp.InternalError(location, "Maximum redirects exceeded", url)
	}

	// If "Actor" is requested, then use that discovery method
	if config.DocumentType == documentTypeActor {
		return client.loadActor(config, url)
	}

	// Otherwise, use "Document" discovery method
	return client.loadDocument(config, url)
}
