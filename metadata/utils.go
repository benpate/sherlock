package metadata

import (
	"mime"
	"slices"
	"strings"

	"github.com/benpate/rosetta/convert"
)

// isActivityStream returns TRUE if the MIME type is either activity+json or ld+json.
func isActivityStream(value string) bool {

	// Duplicated from the parent package — the dependency points from
	// sherlock into this package, never back.

	if mediaType, _, err := mime.ParseMediaType(value); err == nil {
		switch mediaType {
		case "application/activity+json", "application/ld+json":
			return true
		}
	}

	return false
}

// defaultHTTPS appends `https://` to the uri if it doesn't already have a
// valid protocol.
func defaultHTTPS(uri string) string {

	if strings.HasPrefix(uri, "http://") {
		return uri
	}

	if strings.HasPrefix(uri, "https://") {
		return uri
	}

	return "https://" + uri
}

// iconSizesAsInt converts an icon `sizes` attribute (e.g. "128x128", or a
// space-separated list) to the largest dimension found, for ranking icons.
func iconSizesAsInt(value string) int {

	// Duplicated from the parent package's icon sorting — same one-way
	// dependency rule as above.

	if value == "" {
		return 0
	}

	value = strings.ToLower(value)
	parts := strings.Split(value, " ")
	results := make([]int, 0, len(parts))

	for _, part := range parts {

		part, _, _ = strings.Cut(part, "x")

		if result, ok := convert.IntOk(part, 0); ok {
			results = append(results, int(result))
		}
	}

	if len(results) == 0 {
		return 0
	}

	return slices.Max(results)
}
