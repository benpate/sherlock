package metadata

import (
	stdhtml "html"
	"net/url"
	"strings"
	"time"

	"github.com/benpate/rosetta/null"
)

// normalizeURL resolves a possibly-relative URL against a base and applies
// the §4.4 floors, returning an empty string when the value is unusable.
func normalizeURL(value string, base string) string {

	value = strings.TrimSpace(value)

	if value == "" {
		return ""
	}

	// RULE: reject the literal junk strings that broken templating emits.
	switch strings.ToLower(value) {
	case "null", "undefined":
		return ""
	}

	parsedBase, err := url.Parse(base)

	if err != nil {
		return ""
	}

	parsed, err := parsedBase.Parse(value)

	if err != nil {
		return ""
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ""
	}

	if parsed.Host == "" {
		return ""
	}

	return parsed.String()
}

// normalizeCanonicalURL applies normalizeURL plus the same-origin rule for
// canonical URLs (§4.2 exception 1).
func normalizeCanonicalURL(value string, base string) string {

	normalized := normalizeURL(value, base)

	if normalized == "" {
		return ""
	}

	parsedBase, err := url.Parse(base)

	if err != nil {
		return ""
	}

	parsed, err := url.Parse(normalized)

	if err != nil {
		return ""
	}

	// RULE: same-origin required — a page cannot canonicalize itself onto
	// someone else's domain.
	if !strings.EqualFold(parsed.Hostname(), parsedBase.Hostname()) {
		return ""
	}

	return normalized
}

// normalizeLanguage converts the values publishers actually emit into BCP-47:
// OG sends "en_US", HTML sends "en-US", and both should read the same.
func normalizeLanguage(value string) string {

	value = strings.TrimSpace(value)

	if value == "" {
		return ""
	}

	return strings.ReplaceAll(value, "_", "-")
}

// boundDimension returns a declared dimension when it is sane, or 0 (unknown)
// when it is negative or absurdly large.
func boundDimension(value int64) int {

	// RULE: every dimension from a remote source passes through here, so the
	// model's invariant — dimensions are sane or zero — holds regardless of
	// which extractor produced them.
	if value < 0 {
		return 0
	}

	if value > maxDimension {
		return 0
	}

	return int(value)
}

// timeFormats are the layouts publishers actually emit, tried in order.
var timeFormats = []string{
	time.RFC3339,          // 2006-01-02T15:04:05Z07:00
	"2006-01-02T15:04:05", // RFC 3339 without an offset
	time.RFC1123,          // Mon, 02 Jan 2006 15:04:05 MST
	time.RFC1123Z,         // Mon, 02 Jan 2006 15:04:05 -0700
	"2006-01-02",          // bare date
}

// parseTime parses a timestamp permissively across the formats publishers
// emit. Returns null when nothing matches — zero times never enter the model.
func parseTime(value string) null.Object[time.Time] {

	value = strings.TrimSpace(value)

	if value == "" {
		return null.Object[time.Time]{}
	}

	for _, format := range timeFormats {
		if parsed, err := time.Parse(format, value); err == nil {
			return null.NewObject(parsed)
		}
	}

	// I never guess. It is a shocking habit — destructive to the logical faculty.
	return null.Object[time.Time]{}
}

// decodeOnce HTML-entity-decodes a string exactly once, then trims and
// collapses whitespace.
func decodeOnce(value string) string {

	// RULE: extractors call this at the boundary and NEVER again downstream —
	// double-decoding is a real and common bug (§4.4).
	value = stdhtml.UnescapeString(value)

	return strings.Join(strings.Fields(value), " ")
}
