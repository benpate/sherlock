package activitypub

import (
	"github.com/benpate/dns"
	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/sherlock"
)

// Client represents a "middleware" that tries to load an
// ActivityPub/ActivityStream document from the Interwebs.
//
// If the server does not respond with an ActvityPub content-type
// then the request is forwarded to the inner client.
type Client struct {
	innerClient streams.Client
	keyPairFunc KeyPairFunc
	rootClient  streams.Client
	userAgent   string
}

// New returns a fully initialized Client
func New(innerClient streams.Client, options ...ClientOption) streams.Client {

	result := Client{
		innerClient: innerClient,
		userAgent:   "Sherlock (https://github.com/benpate/sherlock)",
	}

	for _, option := range options {
		option(&result)
	}

	return &result
}

func (client *Client) Load(uri string, options ...any) (streams.Document, error) {

	// RULE: This must be a valid URL
	if dns.NotValidURL(uri) {
		return client.innerClient.Load(uri, options...)
	}

	// Build a remote transaction (to try) to load the ActivityStream document
	result := make(map[string]any)
	remoteOptions := remote.Options(options)

	// If we have a KeyPairFunc, then add the AuthorizedFetch option to the transaction.
	if client.keyPairFunc != nil {
		publicKeyID, privateKey := client.keyPairFunc()
		authorizedFetch := sherlock.AuthorizedFetch(publicKeyID, privateKey)
		remoteOptions = append(remoteOptions, authorizedFetch)
	}

	txn := remote.Get(uri).
		Accept(vocab.ContentTypeActivityPub).
		UserAgent(client.userAgent).
		With(remoteOptions...).
		Result(&result)

	// Send the transaction to the Interwebs.
	if err := txn.Send(); err == nil {

		// Confirm that we've received an ActivityPub document
		if contentType := txn.ResponseHeader().Get("Content-Type"); hannibal.IsActivityPubContentType(contentType) {

			return streams.NewDocument(result,
				streams.WithClient(client.rootClient),
				streams.WithHTTPHeader(txn.ResponseHeader()),
			), nil
		}
	}

	return client.innerClient.Load(uri, options...)
}

func (client *Client) Delete(uri string) error {
	return client.innerClient.Delete(uri)
}

func (client *Client) Save(document streams.Document) error {
	return client.innerClient.Save(document)
}

func (client *Client) SetRootClient(rootClient streams.Client) {
	client.innerClient.SetRootClient(rootClient)
	client.rootClient = rootClient
}
