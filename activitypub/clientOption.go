package activitypub

import "github.com/benpate/hannibal/streams"

type ClientOption func(*Client)

// WithInnerClient is an Option that sets the inner client for a Client.
// This allows the Client to forward requests to another client if it fails to load an ActivityPub document.
func WithInnerClient(innerClient streams.Client) ClientOption {
	return func(client *Client) {
		client.innerClient = innerClient
	}
}

// WithUserAgent is an Option that sets the User-Agent header for a Client.
// Applications SHOULD set a custom User-Agent header that identifies the application
// and provides a URL for more information.
func WithUserAgent(userAgent string) ClientOption {

	return func(client *Client) {
		client.userAgent = userAgent
	}
}

// WithKeyPairFunc is an Option that sets the ActorGetter for a Client.
// This allows the Client to retrieve the public key ID and private key for a given URL
// only when needed, rather than performing expensive database queries ahead of time.
func WithKeyPairFunc(fn KeyPairFunc) ClientOption {
	return func(client *Client) {
		client.keyPairFunc = fn
	}
}
