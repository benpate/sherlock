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

// withContext adds the standard ActivityStream @context to the JSON-LD document.
// If we're doing this, it's because we're assembling a "fake" JSON-LD document out of
// other metadata (like OpenGraph, MicroFormats, oEmbed, etc).
func withContext(value mapof.Any) {
	if _, ok := value[vocab.AtContext]; !ok {
		value[vocab.AtContext] = vocab.ContextTypeActivityStreams
	}
}

// canLog is a silly zerolog helper that returns TRUE
// if the provided log level would be allowed
// (based on the global log level).
// This makes it easier to execute expensive code conditionally,
// for instance: marshalling a JSON object for logging.
func canLog(level zerolog.Level) bool {
	return zerolog.GlobalLevel() <= level
}

// canTrace returns TRUE if zerolog is configured to allow Trace logs
// This function is here for completeness.  It may or may not be used
// nolint: unused
func canTrace() bool {
	return canLog(zerolog.TraceLevel)
}

// canDebug returns TRUE if zerolog is configured to allow Debug logs
// This function is here for completeness.  It may or may not be used
// nolint: unused
func canDebug() bool {
	return canLog(zerolog.DebugLevel)
}

// canInfo returns TRUE if zerolog is configured to allow Info logs
// This function is here for completeness.  It may or may not be used
// nolint: unused
func canInfo() bool {
	return canLog(zerolog.InfoLevel)
}
