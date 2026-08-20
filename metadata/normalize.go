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
	//
	// This is STRICTNESS, not tolerance — it rejects input rather than
	// accepting a variant, so it needs no peer fixture to justify it. No
	// captured peer is yet known to emit these; tracked in POSTEL.md.
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
	//
	// This check is a SECURITY decision, not a formatting one, and is fenced
	// off from Postel's law permanently: it decides identity, so no peer's
	// sloppiness is ever grounds for relaxing it (go-hardening).
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
	//
	// This is a BOUND, not a tolerance: it is never widened to accommodate a
	// peer that sends something bigger. Postel's liberality covers shape, never
	// limits (go-hardening).
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

	// Spec-defined interchange formats. These need no peer fixture: they are
	// real published formats, and accepting them is conformance, not leniency.
	time.RFC3339,          // 2006-01-02T15:04:05Z07:00 — what nearly every publisher sends
	"2006-01-02T15:04:05", // RFC 3339 without an offset
	time.RFC1123,          // HTTP-date, RFC 9110 §5.6.7
	time.RFC1123Z,         // RFC 9110 §5.6.7, numeric-offset form
	"2006-01-02",          // ISO 8601 calendar date

	// Go's time.Time.String() layout — NOT an interchange format. It appears
	// when a Go-based generator prints a timestamp with %v instead of
	// formatting it. Smashing Magazine (Hugo) emits exactly this in
	// article:published_time; see testdata/smashing-magazine.html.
	"2006-01-02 15:04:05 -0700 MST",
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

			// RULE: a zero time is not a date. Accepting Go's String() layout
			// above means "0001-01-01 00:00:00 +0000 UTC" — the Go zero value
			// printed verbatim — now parses cleanly, and Smashing Magazine
			// serves exactly that on pages with no publication date.
			if parsed.IsZero() {
				return null.Object[time.Time]{}
			}

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
