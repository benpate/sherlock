package tagspub

import (
	"github.com/benpate/hannibal/streams"
)

// Client represents a Tags.Pub "middleware" that takes
// Hashtag-looking URIs and converts them into WebFinger handles
//
// To work properly, this middleware MUST be installed ABOVE
// the WebFinger middleware.
type Client struct {
	hostname    string
	innerClient streams.Client
}

// New returns a fully initialized Client
func New(innerClient streams.Client, options ...ClientOption) streams.Client {

	result := Client{
		hostname:    "tags.pub",
		innerClient: innerClient,
	}

	for _, option := range options {
		option(&result)
	}

	return result
}

// Load resolves a hashtag-looking URI through the tags.pub server, falling back
// to the inner client (with the original URI) when it is not a hashtag or fails.
func (client Client) Load(uri string, options ...any) (streams.Document, error) {

	// Process URIs that look like a hashtags
	if match, newURI := client.isHashtag(uri); match {

		// Try to look up the username via the inner client.
		if result, err := client.innerClient.Load(newURI, options...); err == nil {
			return result, nil
		}

		// Fail gracefully. Maybe this WASN'T a proper bridgy URL,
		// so continue down the stack
	}

	// Forward the request to the innerClient
	return client.innerClient.Load(uri, options...)
}

// Delete resolves a hashtag-looking URI through the tags.pub server, falling back
// to the inner client (with the original URI) when it is not a hashtag or fails.
func (client Client) Delete(uri string) error {

	// Process URIs that look like a hashtags
	if match, newURI := client.isHashtag(uri); match {

		// Try to look up the username via the inner client.
		if err := client.innerClient.Delete(newURI); err == nil {
			return nil
		}

		// Fail gracefully. Maybe this WASN'T a proper bridgy URL,
		// so continue down the stack
	}

	// Forward the request to the innerClient
	return client.innerClient.Delete(uri)
}

// Save forwards a document to the inner client's cache.
func (client Client) Save(document streams.Document) error {
	return client.innerClient.Save(document)
}

// SetRootClient propagates the top-level client to the inner client.
func (client Client) SetRootClient(rootClient streams.Client) {
	if client.innerClient != nil {
		client.innerClient.SetRootClient(rootClient)
	}
}

// isHashtag reports whether the id is a hashtag and, if so, returns the
// "tag@hostname" WebFinger handle it maps to.
func (client Client) isHashtag(id string) (bool, string) {

	if match, newID := IsHashtag(id); match {
		return true, newID + "@" + client.hostname
	}

	return false, id
}
