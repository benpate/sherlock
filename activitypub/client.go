// Package activitypub provides a Sherlock client middleware that loads
// ActivityPub/ActivityStream documents directly from their canonical URLs.
package activitypub

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/uri"
)

// Client represents a "middleware" that tries to load an
// ActivityPub/ActivityStream document from the Interwebs.
//
// If the server does not respond with an ActvityPub content-type
// then the request is forwarded to the inner client.
type Client struct {
	innerClient     streams.Client
	keyPairFunc     KeyPairFunc
	rootClient      streams.Client
	userAgent       string
	allowPrivateIPs bool
}

// New returns a fully initialized Client
func New(options ...ClientOption) streams.Client {

	result := Client{
		userAgent: "Sherlock (https://github.com/benpate/sherlock)",
	}

	for _, option := range options {
		option(&result)
	}

	return &result
}

// Load retrieves the ActivityPub document at the given URL, falling back to the
// inner client when the URL is invalid or does not return ActivityPub content.
func (client *Client) Load(id string, options ...any) (streams.Document, error) {

	const location = "activitypub.Client.Load"

	// RULE: This must be a valid URL
	if uri.NotValidURL(id) {
		if client.innerClient != nil {
			return client.innerClient.Load(id, options...)
		}
		return streams.NilDocument(), derp.NotFound(location, "Invalid URL", id)
	}

	// Send the transaction to the Interwebs.
	response, err := client.fetch(id, options)

	// RULE: An ActivityPub document is returned only when the host that served it may speak for
	// its id.  A rejected document never falls through to the inner client.
	if (err == nil) && response.isActivityPub {
		return client.trustedDocument(id, response, options)
	}

	// FALLTHROUGH means FAILURE...

	// If we have an inner client, then try to forward the request to it instead.
	if client.innerClient != nil {
		return client.innerClient.Load(id, options...)
	}

	// Otherwise, return a failure to the parent.
	empty := streams.NilDocument().AddOptions(
		streams.WithClient(client.rootClient),
		streams.WithHTTPHeader(response.header),
	)

	return empty, derp.Wrap(err, location, "Unable to load document.", id)
}

// Delete forwards a cache-delete to the inner client, or is a no-op when there
// is no inner client.
func (client *Client) Delete(id string) error {
	if client.innerClient == nil {
		return nil
	}
	return client.innerClient.Delete(id)
}

// Save forwards a document to the inner client's cache, or is a no-op when there
// is no inner client.
func (client *Client) Save(document streams.Document) error {
	if client.innerClient == nil {
		return nil
	}
	return client.innerClient.Save(document)
}

// SetRootClient records the top-level client and propagates it to the inner client.
func (client *Client) SetRootClient(rootClient streams.Client) {
	client.rootClient = rootClient
	if client.innerClient != nil {
		client.innerClient.SetRootClient(rootClient)
	}
}
