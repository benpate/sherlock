package metadata

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalizeURL(t *testing.T) {

	base := "https://example.com/articles/one"

	require.Equal(t, "https://example.com/img.jpg", normalizeURL("/img.jpg", base))
	require.Equal(t, "https://example.com/articles/img.jpg", normalizeURL("img.jpg", base))
	require.Equal(t, "https://cdn.example.net/a.png", normalizeURL("https://cdn.example.net/a.png", base))
	require.Empty(t, normalizeURL("javascript:alert(1)", base))
	require.Empty(t, normalizeURL("null", base))
	require.Empty(t, normalizeURL("UNDEFINED", base))
	require.Empty(t, normalizeURL("", base))
}

func TestNormalizeCanonicalURL_SameOrigin(t *testing.T) {

	base := "https://example.com/articles/one"

	require.Equal(t, "https://example.com/canonical", normalizeCanonicalURL("/canonical", base))
	require.Empty(t, normalizeCanonicalURL("https://evil.example.net/steal", base)) // cross-origin rejected
}

func TestNormalizeLanguage(t *testing.T) {
	require.Equal(t, "en-US", normalizeLanguage("en_US")) // OG style
	require.Equal(t, "en-US", normalizeLanguage("en-US")) // HTML style
	require.Empty(t, normalizeLanguage("  "))
}

func TestParseTime(t *testing.T) {
	require.True(t, parseTime("2026-08-15T10:30:00Z").IsPresent())          // RFC 3339
	require.True(t, parseTime("2026-08-15T10:30:00").IsPresent())           // no offset
	require.True(t, parseTime("Fri, 15 Aug 2026 10:30:00 GMT").IsPresent()) // RFC 1123
	require.True(t, parseTime("2026-08-15").IsPresent())                    // bare date
	require.True(t, parseTime("not a date").IsNull())
	require.True(t, parseTime("").IsNull())

	// A null timestamp reads back as the zero time, never a nil dereference.
	require.True(t, parseTime("not a date").Object().IsZero())
}

func TestDecodeOnce(t *testing.T) {

	// Decode exactly once: a double-encoded ampersand stays half-encoded.
	require.Equal(t, "Tom & Jerry", decodeOnce("Tom &amp; Jerry"))
	require.Equal(t, "Tom &amp; Jerry", decodeOnce("Tom &amp;amp; Jerry"))
	require.Equal(t, "spaced out", decodeOnce("  spaced \n\t out  "))
}

func TestBoundDimension(t *testing.T) {

	// Sane values pass through; negative and absurd values become "unknown".
	require.Equal(t, 640, boundDimension(640))
	require.Equal(t, 0, boundDimension(0))
	require.Equal(t, maxDimension, boundDimension(maxDimension))
	require.Equal(t, 0, boundDimension(-1))
	require.Equal(t, 0, boundDimension(maxDimension+1))
	require.Equal(t, 0, boundDimension(1<<62))
}
