package bridgyfed

import (
	"strings"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote/options"
)

// Client represents a Bridgy Fed "middleware" that takes
// Bluesky-looking URIs and converts them into WebFinger URIs
//
// To work properly, this middleware MUST be installed ABOVE
// the WebFinger middleware.
type Client struct {
	hostname    string
	innerClient streams.Client
}

// New returns a fully initialized Client
func New(innerClient streams.Client, options ...ClientOption) streams.Client {

	result := Client{
		innerClient: innerClient,
		hostname:    "bsky.brid.gy",
	}

	for _, option := range options {
		option(&result)
	}

	return result
}

func (client Client) Load(id string, loadOptions ...any) (streams.Document, error) {

	// If we think we can load this URL, then try to..
	if looksLikeBluesky, bridgyHandle := client.looksLikeBluesky(id); looksLikeBluesky {

		// Inject the only Accept header that Bridgy Fed likes.
		loadOptions = append(loadOptions, options.Accept(vocab.ContentTypeActivityPub))

		// Try to look up the username via the inner client.
		if result, err := client.innerClient.Load(bridgyHandle, loadOptions...); err == nil {
			return result, nil
		}

		// Fail gracefully. Maybe this WASN'T a proper bridgy URL,
		// so continue down the stack....
	}

	// Forward the request to the innerClient
	return client.innerClient.Load(id, loadOptions...)
}

func (client Client) Delete(id string) error {

	// If we think we can load this URL, then try to..
	if looksLikeBluesky, bridgyHandle := client.looksLikeBluesky(id); looksLikeBluesky {

		// Try to look up the username via the inner client.
		if err := client.innerClient.Delete(bridgyHandle); err == nil {
			return nil
		}

		// Fail gracefully. Maybe this WASN'T a proper bridgy URL,
		// so continue down the stack....
	}

	// Forward the request to the innerClient
	return client.innerClient.Delete(id)
}

func (client Client) Save(document streams.Document) error {
	return client.innerClient.Save(document)
}

func (client Client) SetRootClient(rootClient streams.Client) {
	client.innerClient.SetRootClient(rootClient)
}

func (client Client) looksLikeBluesky(uri string) (bool, string) {

	if !LooksLikeBluesky(uri) {
		return false, uri
	}

	// Format the handle
	bridgyHandle := uri

	if !strings.HasPrefix(bridgyHandle, "@") {
		bridgyHandle = "@" + bridgyHandle
	}

	bridgyHandle = bridgyHandle + "@" + client.hostname

	return true, bridgyHandle
}
