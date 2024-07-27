package sherlock

import (
	"crypto"

	"github.com/benpate/remote"
)

// ClientOption defines a functional option that modifies a Client object
type ClientOption func(*Client)

// WithUserAgent is a ClientOption that sets the UserAgent property on the Client object
func WithUserAgent(userAgent string) ClientOption {
	return func(client *Client) {
		client.UserAgent = userAgent
	}
}

// WithRemoteOptions is a ClientOption that appends one or more remote.Option
// objects to the Client object RemoteOptions are executed on every remote request
func WithRemoteOptions(options ...remote.Option) ClientOption {
	return func(client *Client) {
		client.RemoteOptions = append(client.RemoteOptions, options...)
	}
}

// WithActor is a ClientOption that set up the AuthorizedFetch remote middleware,
// which will sign all outbound requests according to the ActivityPub "Authorized Fetch"
// convention: https://funfedi.dev/testing_tools/http_signatures/
func WithActor(publicKeyID string, privateKey crypto.PrivateKey) ClientOption {
	return func(client *Client) {
		client.RemoteOptions = append(client.RemoteOptions, AuthorizedFetch(publicKeyID, privateKey))
	}
}
