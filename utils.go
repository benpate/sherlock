package sherlock

import (
	"mime"
	"strings"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/microcosm-cc/bluemonday"
	"github.com/rs/zerolog"
)

func sanitizeHTML(value string) string {
	return bluemonday.UGCPolicy().Sanitize(value)
}

func sanitizeText(value string) string {
	return bluemonday.StrictPolicy().Sanitize(value)
}

// isActivityStream returns TRUE if the MIME type is either activity+json or ld+json
func isActivityStream(value string) bool {

	// ActivityStreams have their own MIME type, but we have to check some alternates, too.
	if mediaType, _, err := mime.ParseMediaType(value); err == nil {
		switch mediaType {
		case "application/activity+json", "application/ld+json":
			return true
		}
	}

	return false
}

// defaultHTTPS appends `https://` to the uri if it doesn't already have a valid protocol.
func defaultHTTPS(uri string) string {

	if strings.HasPrefix(uri, "http://") {
		return uri
	}

	if strings.HasPrefix(uri, "https://") {
		return uri
	}

	return "https://" + uri
}

/*
// setMetadata sets common metadata from the HTTP response header
func (client *Client) setMetadata(document streams.Document, header http.Header) {
	document.WithOptions(
		streams.WithClient(client),
		streams.WithMeta("cache-control", header.Get("cache-control")),
		streams.WithMeta("etag", header.Get("etag")),
		streams.WithMeta("expires", header.Get("expires")),
	)
}
*/

// withContext adds the standard ActivityStreams @context to the JSON-LD document.
// If we're doing this, it's because we're assembling a "fake" JSON-LD document out of
// other metadata (like OpenGraph, MicroFormats, oEmbed, etc).
func withContext(value mapof.Any) {
	if _, ok := value[vocab.AtContext]; !ok {
		value[vocab.AtContext] = vocab.ContextTypeActivityStreams
	}
}

// canLog returns TRUE if the current logging level is supported
func canLog(level zerolog.Level) bool {
	return level >= zerolog.GlobalLevel()
}
