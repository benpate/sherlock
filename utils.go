package sherlock

import (
	"mime"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/benpate/derp"
	"github.com/benpate/digit"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/compare"
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

// sortImageLinks is a slices.SortFunc function that ranks digit.Links by their size and type.
func sortImageLinks(a, b digit.Link) int {

	// First, prefer larger images
	aSize := iconSizesAsInt(a.Properties["sizes"])
	bSize := iconSizesAsInt(b.Properties["sizes"])

	if result := compare.Int(aSize, bSize); result != 0 {
		return result
	}

	// Next, prefer images by type
	return compare.Int(iconMediaTypeAsInt(a.MediaType), iconMediaTypeAsInt(b.MediaType))
}

// iconSizeAsInt converts an image size string (in the form of "128x128") to the maximum
// integer value of the two dimensions.  This is useful for sorting images by size.
func iconSizesAsInt(value string) int {

	// Empty values are empty
	if value == "" {
		return 0
	}

	// Convert to lowercase, and split into parts
	value = strings.ToLower(value)
	parts := strings.Split(value, " ")
	results := make([]int, 0, len(parts))

	// Scan each part for the first number in the dimension
	for _, part := range parts {

		part, _, _ = strings.Cut(part, "x")

		// If we have a number, then add that to the potential result
		if result, err := strconv.ParseInt(part, 10, 64); err == nil {
			results = append(results, int(result))
		}
	}

	// If we found no results, then return 0
	if len(results) == 0 {
		return 0
	}

	// Return the largest number found
	return slices.Max(results)
}

// iconMediaTypeAsInt converts an image type string (in the form of "image/png") to a numeric value
// that cam be used to sort images by type.
func iconMediaTypeAsInt(value string) int {

	switch value {
	case "image/webp":
		return 256
	case "image/png":
		return 255
	case "image/jpg":
		return 254
	case "image/jpeg":
		return 253
	case "image/svg":
		return 252
	case "image/svg+xml":
		return 251
	case "image/gif":
		return 250
	case "image/bmp":
		return 248
	case "image/tiff":
		return 247
	case "image/tiff+xml":
		return 246
	case "image/x-icon":
		return 245
	case "image/vnd.microsoft.icon":
		return 244
	default:
		return 0
	}
}

// hostOnly returns the protocol and hostname of a URL, without the path or query string
func hostOnly(value string) string {

	parsedURL, err := url.Parse(value)

	if err != nil {
		derp.Report(derp.Wrap(err, "sherlock.hostOnly", "Error parsing URL", value))
		return value
	}

	// Strip path and query string (use root URL only)
	parsedURL.Path = ""
	parsedURL.RawQuery = ""

	// Rewrite the value without the path and query string
	return parsedURL.String()
}
