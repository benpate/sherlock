package sherlock

type ClientOption func(*Client)

// WithKeyPairFunc is an Option that sets the ActorGetter for a Client.
// This allows the Client to retrieve the public key ID and private key for a given URL
// only when needed, rather than performing expensive database queries ahead of time.
func WithKeyPairFunc(fn KeyPairFunc) ClientOption {
	return func(client *Client) {
		client.keyPairFunc = fn
	}
}

// WithUserAgent is a ClientOption that sets the UserAgent property on the Client object
func WithUserAgent(userAgent string) ClientOption {
	return func(client *Client) {
		client.userAgent = userAgent
	}
}
