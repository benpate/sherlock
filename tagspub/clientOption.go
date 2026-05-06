package tagspub

type ClientOption func(*Client)

// WithHostname uses a custom server hostname (instead of the default tags.pub)
func WithHostname(hostname string) ClientOption {
	return func(client *Client) {
		client.hostname = hostname
	}
}
