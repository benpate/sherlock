package sherlock

import "github.com/benpate/remote"

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
func WithRemoteOptions(middleware ...remote.Option) ClientOption {
	return func(client *Client) {
		client.RemoteOptions = append(client.RemoteOptions, middleware...)
	}
}

// WithDebugLevel sets the debug level for the client.  Valid values are "verbose", "terse", and "none"
func WithDebug(level string) ClientOption {
	return func(client *Client) {
		client.DebugLevel = level
	}
}
