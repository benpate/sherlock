package sherlock

// ClientOption defines a functional option that modifies a Client object
type ClientOption func(*Client)

// WithUserAgent is a ClientOption that sets the UserAgent property on the Client object
func WithUserAgent(userAgent string) ClientOption {
	return func(client *Client) {
		client.UserAgent = userAgent
	}
}
