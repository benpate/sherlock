package webfinger

import (
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/digit"
	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/remote"
)

// Client implements a WebFinger "middleware" that identifies
// Webfinger-looking URIs and converts them into
type Client struct {
	innerClient   streams.Client
	remoteOptions []remote.Option
}

// New returns a fully populated WebFinger client.
func New(innerClient streams.Client, options ...remote.Option) streams.Client {
	result := Client{
		innerClient:   innerClient,
		remoteOptions: options,
	}

	return result
}

// Load attempts to retrieve the URI from the Interweb, translating
// WebFinger @handles into proper URLs if possible.
func (client Client) Load(uri string, options ...any) (streams.Document, error) {

	const location = "sherlock.webfinger.Load"

	isWebfinger, newURI, err := client.isWebfinger(uri)

	if err != nil {
		return streams.NilDocument(), derp.Wrap(err, location, "Unable to load Webfinger info", uri)
	}

	if isWebfinger {
		return client.innerClient.Load(newURI, options...)
	}

	// Otherwise, allow the inner client to look instead
	return client.innerClient.Load(uri, options...)
}

// Load attempts to retrieve the URI from the Interweb, translating
// WebFinger @handles into proper URLs if possible.
func (client Client) Delete(uri string) error {

	const location = "sherlock.webfinger.Delete"

	isWebfinger, newURI, err := client.isWebfinger(uri)

	if err != nil {
		return derp.Wrap(err, location, "Unable to load Webfinger info", uri)
	}

	if isWebfinger {
		return client.innerClient.Delete(newURI)
	}

	// Otherwise, allow the inner client to look instead
	return client.innerClient.Delete(uri)
}

func (client Client) Save(document streams.Document) error {
	return client.innerClient.Save(document)
}

func (client Client) SetRootClient(rootClient streams.Client) {
	client.innerClient.SetRootClient(rootClient)
}

func (client Client) isWebfinger(uri string) (bool, string, error) {

	const location = "sherlock.webfinger.isWebfinger"

	// Quick check: does the URI look like an email/username?
	if !strings.Contains(uri, "@") {
		return false, uri, nil
	}

	// Thorough check: can digit determine an actual URL for the URI
	webFingerURLs := digit.ParseAccount(uri)

	if len(webFingerURLs) == 0 {
		return false, uri, nil
	}

	// Cool. Try to load the Actor via WebFinger
	response, err := digit.Lookup(uri, client.remoteOptions...)

	// Invalid response -> error. We knew what we were doing, your server just broke.
	if err != nil {
		return false, uri, derp.Wrap(err, location, "Unable to load URL identified in WebFinger result")
	}

	// Otherwise, see if we can find an ActivityPub endpoint
	for _, link := range response.Links {
		if (link.RelationType == digit.RelationTypeSelf) && (hannibal.IsActivityPubContentType(link.MediaType)) {
			return true, link.Href, nil
		}
	}

	return false, uri, nil
}
