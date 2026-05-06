package bridgyfed

type ClientOption func(*Client)

// WithHostname uses a custom hostname (instead of the default bsky.brid.gy)
func WithHostname(hostname string) ClientOption {
	return func(client *Client) {
		client.hostname = hostname
	}
}
