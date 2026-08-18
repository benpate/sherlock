package metadata

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/remote"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

// document is the raw material for metadata extraction: one fetched response,
// decoded to UTF-8 and (when HTML) parsed exactly once, shared by every extractor.
type document struct {
	RequestURL     string         // the URL the caller asked for
	FinalURL       string         // the URL after redirects
	Header         http.Header    // response headers
	ContentType    string         // media type only, parameters stripped
	Body           []byte         // response body, transcoded to UTF-8
	Truncated      bool           // TRUE when Body was cut at the byte cap
	Root           *html.Node     // parsed HTML tree; nil when the response is not HTML
	ActivityStream map[string]any // parsed AS2 object when the response is ActivityStreams; nil otherwise
}

// IsActivityStream returns TRUE when the fetched response was a native
// ActivityStreams document, which short-circuits all other extraction.
func (doc *document) IsActivityStream() bool {
	return doc.ActivityStream != nil
}

// IsHTML returns TRUE when the fetched response parsed as an HTML tree.
func (doc *document) IsHTML() bool {
	return doc.Root != nil
}

// newDocument fetches a URL once — negotiating ActivityStreams and HTML on
// the same request — and prepares the response for extraction.
func newDocument(ctx context.Context, config config, url string) (*document, error) {

	const location = "sherlock.metadata.newDocument"

	// Default the URL scheme to https
	url = defaultHTTPS(url)

	// Fetch the raw body. Result(&body) bypasses remote's content-type-driven
	// decoding, and capBody truncates (rather than rejects) oversized bodies so
	// metadata in the <head> of a large page still survives.
	// The empty (not nil) initializer is load-bearing: remote fills body
	// through a pointer, which static nil analysis cannot follow, so a nil
	// start reads as a possible nil slice at the trim below.
	body := []byte{}
	maxBodySize := config.maxBodySize()

	txn := remote.Get(url).
		WithContext(ctx).
		UserAgent(config.UserAgent).
		Header("Accept", "application/activity+json, text/html;q=0.9, */*;q=0.8").
		AllowPrivateIPs(config.AllowPrivateIPs).
		With(capBody(maxBodySize)).
		With(config.RemoteOptions...).
		Result(&body)

	if err := txn.Send(); err != nil {
		return nil, derp.Wrap(err, location, "Unable to fetch URL", url)
	}

	// RULE: capBody reads one byte past the cap, so an over-length body proves
	// the server had more to give. Trim that probe byte back off.
	truncated := int64(len(body)) > maxBodySize

	if truncated {
		body = body[:maxBodySize]
	}

	result := &document{
		RequestURL: url,
		FinalURL:   finalURL(txn, url),
		Header:     txn.ResponseHeader(),
		Body:       body,
		Truncated:  truncated,
	}

	// RULE: Content-Type parameters (charset, profile) are stripped for
	// dispatch, but the raw header is kept for charset decoding below.
	rawContentType := result.Header.Get(contentTypeHeader)
	if mediaType, _, err := mime.ParseMediaType(rawContentType); err == nil {
		result.ContentType = mediaType
	}

	// If the server answered our ActivityStreams preference, parse it as JSON
	// and stop — nothing downstream can improve on a native AS2 object.
	if isActivityStream(result.ContentType) {
		activityStream := map[string]any{}
		if err := json.Unmarshal(body, &activityStream); err == nil {
			result.ActivityStream = activityStream
			return result, nil
		}
	}

	// Otherwise, treat the response as HTML: transcode to UTF-8 and parse the
	// tree exactly once. Extractors all read this one tree.
	if utf8Body, err := decodeToUTF8(body, rawContentType); err == nil {
		result.Body = utf8Body
	}

	root, err := html.Parse(bytes.NewReader(result.Body))
	if err != nil {
		return nil, derp.Wrap(err, location, "Unable to parse HTML", url)
	}
	result.Root = root

	// Data! Data! Data! I can't make bricks without clay.
	return result, nil
}

// finalURL returns the URL after redirects, falling back to the request URL
// when the response chain is unavailable.
func finalURL(txn *remote.Transaction, fallback string) string {

	response := txn.Response()

	if response == nil {
		return fallback
	}

	if response.Request == nil {
		return fallback
	}

	if response.Request.URL == nil {
		return fallback
	}

	return response.Request.URL.String()
}

// decodeToUTF8 transcodes a response body to UTF-8 using the Content-Type
// header, <meta> tags, and BOM sniffing — no detection library (D10).
func decodeToUTF8(body []byte, contentType string) ([]byte, error) {

	const location = "sherlock.metadata.decodeToUTF8"

	reader, err := charset.NewReader(bytes.NewReader(body), contentType)

	if err != nil {
		return body, derp.Wrap(err, location, "Unable to determine charset")
	}

	decoded, err := io.ReadAll(reader)

	if err != nil {
		return body, derp.Wrap(err, location, "Unable to transcode body")
	}

	return decoded, nil
}
