package sherlock

// Client implements the hannibal/streams.Client interface, and is used to load JSON-LD documents from remote servers.
// The sherlock client maps additional meta-data into a standard ActivityStreams document.
type Client struct{}

// NewClient returns a fully initialized Client object
func NewClient() Client {
	return Client{}
}
