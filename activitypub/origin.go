package activitypub

import (
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/sherlock"
	"github.com/benpate/uri"
)

/******************************************
 * Origin Checks
 *
 * A host may speak only for its own ids.  A document
 * whose id names another host is trusted only when that
 * host serves the same id itself.
 ******************************************/

// fetchResponse is what a single ActivityPub GET returned.
type fetchResponse struct {
	body          map[string]any // The decoded response body
	header        http.Header    // The response headers
	finalURL      string         // The URL that answered, after any redirects
	isActivityPub bool           // TRUE if the response declared an ActivityPub content type
}

// fetch sends one GET for an ActivityPub document, signed when a KeyPairFunc is configured.
func (client *Client) fetch(url string, options []any) (fetchResponse, error) {

	body := make(map[string]any)

	// The spread is important: `remote.Options(options)` compiles, passes the whole
	// slice as one `any`, matches no Option, and silently drops every caller option.
	remoteOptions := remote.Options(options...)

	// If we have a KeyPairFunc, then add the AuthorizedFetch option to the transaction.
	if client.keyPairFunc != nil {
		publicKeyID, privateKey := client.keyPairFunc()
		authorizedFetch := sherlock.AuthorizedFetch(publicKeyID, privateKey)
		remoteOptions = append(remoteOptions, authorizedFetch)
	}

	txn := remote.Get(url).
		Accept(vocab.ContentTypeActivityPub).
		UserAgent(client.userAgent).
		AllowPrivateIPs(client.allowPrivateIPs).
		With(remoteOptions...).
		Result(&body)

	// NOTE: The response is described even when the request failed, so that a caller can pass its
	// headers on.  It must be read after Send, never before.
	err := txn.Send()

	header := txn.ResponseHeader()
	response := fetchResponse{
		body:          body,
		header:        header,
		finalURL:      finalURL(txn, url),
		isActivityPub: hannibal.IsActivityPubContentType(header.Get("Content-Type")),
	}

	return response, err
}

// trustedDocument returns a fetched document only when the host that served it may speak for its
// id, fetching the id from its own host once when the two differ.
func (client *Client) trustedDocument(url string, response fetchResponse, options []any) (streams.Document, error) {

	const location = "activitypub.Client.trustedDocument"

	// NOTE: Every rejection is Forbidden, never NotFound.  A caching layer may delete the URL it
	// was loading on a NotFound, and a forged document must not be able to trigger that.

	claimedID, err := documentID(response.body)

	if err != nil {
		return client.rejected(response), derp.Wrap(err, location, "Document id is unusable", url, derp.WithForbidden())
	}

	// A document with no id claims nothing, so there is nothing to check
	if claimedID == "" {
		return client.newDocument(response), nil
	}

	// RULE: A host speaks for its own ids
	if sameHost(claimedID, response.finalURL) {
		return client.newDocument(response), nil
	}

	// Otherwise only the id's own host may vouch for it. Try the original host once
	confirmation, err := client.fetch(claimedID, options)

	if err != nil {
		return client.rejected(response), derp.Wrap(err, location, "Confirming document id with its own host", url, claimedID, derp.WithForbidden())
	}

	// RULE: The confirmation must be an ActivityPub document
	if !confirmation.isActivityPub {
		return client.rejected(response), derp.Forbidden(location, "Document id is not confirmed by its own host", url, claimedID)
	}

	// RULE: The confirmation must name the same id, served from that id's own host
	if confirmedID, err := documentID(confirmation.body); (err != nil) || (confirmedID != claimedID) {
		return client.rejected(response), derp.Forbidden(location, "Document id is not confirmed by its own host", url, claimedID, confirmedID)
	}

	if !sameHost(claimedID, confirmation.finalURL) {
		return client.rejected(response), derp.Forbidden(location, "Document id is not confirmed by its own host", url, claimedID, confirmation.finalURL)
	}

	// The id's own host has spoken, so its copy is the one we keep
	return client.newDocument(confirmation), nil
}

// rejected returns the empty document that accompanies a rejection, carrying the response headers.
func (client *Client) rejected(response fetchResponse) streams.Document {
	return streams.NilDocument(
		streams.WithClient(client.rootClient),
		streams.WithHTTPHeader(response.header),
	)
}

// newDocument wraps a trusted response as a document bound to the root client.
func (client *Client) newDocument(response fetchResponse) streams.Document {
	return streams.NewDocument(
		response.body,
		streams.WithClient(client.rootClient),
		streams.WithHTTPHeader(response.header),
	)
}

// finalURL returns the URL that answered a transaction after any redirects, or the requested URL
// when the response does not say.
func finalURL(txn *remote.Transaction, requested string) string {

	response := txn.Response()

	if (response == nil) || (response.Request == nil) || (response.Request.URL == nil) {
		return requested
	}

	return response.Request.URL.String()
}

// documentID returns the id a document claims for itself, or an error when that id is not a string,
// not a valid URL, or not what streams.Document.ID reports for the same document.
func documentID(body map[string]any) (string, error) {

	const location = "activitypub.documentID"

	// JSON-LD allows "@id" in place of "id", so a document cannot dodge the check by using it
	value, exists := body[vocab.PropertyID]

	if !exists {
		value, exists = body[vocab.PropertyID_Alternate]
	}

	if !exists {
		return "", nil
	}

	// RULE: An id is a string
	id, isString := value.(string)

	if !isString {
		return "", derp.BadRequest(location, "Document id must be a string", value)
	}

	if id == "" {
		return "", nil
	}

	// RULE: An id is a valid URL, which gives it a host to check
	if uri.NotValidURL(id) {
		return "", derp.BadRequest(location, "Document id must be a valid URL", id)
	}

	// RULE: Document.ID sanitizes the id as HTML, and callers store what it returns, so the id
	// checked here must be the id they will see
	if reported := streams.NewDocument(body).ID(); reported != id {
		return "", derp.BadRequest(location, "Document id must survive sanitizing unchanged", id, reported)
	}

	return id, nil
}

// sameHost returns TRUE if two URLs name the same host, ignoring case and port.
func sameHost(first string, second string) bool {

	host := uri.Hostname(first)

	// A URL with no host vouches for nothing
	if host == "" {
		return false
	}

	return host == uri.Hostname(second)
}
