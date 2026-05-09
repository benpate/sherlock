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

func (client Client) Load(uri string, options ...any) (streams.Document, error) {

	// Process URIs that look like a hashtags
	if match, newURI := client.isHashtag(uri); match {

		// Try to look up the username via the inner client.
		if result, err := client.innerClient.Load(newURI, options); err == nil {
			return result, nil
		}

		// Fail gracefully. Maybe this WASN'T a proper bridgy URL,
		// so continue down the stack
	}

	// Forward the request to the innerClient
	return client.innerClient.Load(uri, options...)
}

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

func (client Client) Save(document streams.Document) error {
	return client.innerClient.Save(document)
}

func (client Client) SetRootClient(rootClient streams.Client) {
	client.innerClient.SetRootClient(rootClient)
}

func (client Client) isHashtag(id string) (bool, string) {

	if match, newID := IsHashtag(id); match {
		return true, newID + "@" + client.hostname
	}

	return false, id
}
