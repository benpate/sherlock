package tombstone

import (
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/property"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
)

type Client struct {
	innerClient streams.Client
}

// New returns a fully initialized Client object
func New(innerClient streams.Client) *Client {

	// Create a default client
	result := &Client{
		innerClient: innerClient,
	}

	// Default our child's "RootClient" to our current value.
	// This may be overridden by a parent
	result.innerClient.SetRootClient(result)

	// Woot woot.
	return result
}

/******************************************
 * Hannibal HTTP Client Methods
 ******************************************/

// SetRootClient applies a "top level" client (which is needed by some hannibal client implementations)
func (client *Client) SetRootClient(rootClient streams.Client) {
	client.innerClient.SetRootClient(rootClient)
}

// Load retrieves a URL from the cache/interweb, returning it as a streams.Document
func (client *Client) Load(url string, options ...any) (streams.Document, error) {

	// Try to load the document from the inner client.
	result, err := client.innerClient.Load(url, options...)

	// If success, then success.
	if err == nil {
		return result, nil
	}

	// If this is a "Gone" error, then generate an artificial "Tombstone" instead.
	if derp.IsGone(err) {

		// If we don't ALREADY have a Tombstone, then overwrite the response with a Tombstone.
		if result.Type() != vocab.ObjectTypeTombstone {

			result.SetValue(property.Map{
				"id":   url,
				"type": vocab.ObjectTypeTombstone,
			})
		}

		// Return the Tombstone WITHOUT an error. This will hold a place in the cache
		// so we don't keep trying to load this document from the original server.
		return result, nil
	}

	// Otherwise, just pass through the original error.
	return result, err
}

func (client *Client) Save(document streams.Document) error {
	return client.innerClient.Save(document)
}

// Delete removes a single document from the cache
func (client *Client) Delete(url string) error {
	return client.innerClient.Delete(url)
}
