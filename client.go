package sherlock

// Client implements the hannibal/streams.Client interface, and is used to load JSON-LD documents from remote servers.
// The sherlock client maps additional meta-data into a standard ActivityStreams document.
type Client struct {
	UserAgent string
}

// NewClient returns a fully initialized Client object
func NewClient(options ...ClientOption) Client {

	// Create a default Client
	result := Client{
		UserAgent: "Sherlock: github.com/benpate/sherlock",
	}

	// Apply options (duh)
	result.ApplyOptions(options...)

	// Success
	return result
}

func (client *Client) ApplyOptions(options ...ClientOption) {
	for _, option := range options {
		option(client)
	}
}
