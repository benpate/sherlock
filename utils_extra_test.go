package sherlock

import (
	"testing"

	"github.com/benpate/digit"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// closureTest is a tiny helper used by the closure-driven tests in this file.
// It runs `fn` against an input and asserts that the result equals `expected`.
// See https://medium.com/@cep21/628a41497e5e for the pattern.
func runStringCases(t *testing.T, fn func(string) string, cases map[string]string) {
	t.Helper()
	for input, expected := range cases {
		assert.Equal(t, expected, fn(input), "input: %q", input)
	}
}

func runIntCases(t *testing.T, fn func(string) int, cases map[string]int) {
	t.Helper()
	for input, expected := range cases {
		assert.Equal(t, expected, fn(input), "input: %q", input)
	}
}

func runBoolCases(t *testing.T, fn func(string) bool, cases map[string]bool) {
	t.Helper()
	for input, expected := range cases {
		assert.Equal(t, expected, fn(input), "input: %q", input)
	}
}

/******************************************
 * sanitizeHTML / sanitizeText
 ******************************************/

func TestSanitizeHTML(t *testing.T) {
	runStringCases(t, sanitizeHTML, map[string]string{
		"":                          "",
		"hello world":               "hello world",
		"<b>bold</b>":               "<b>bold</b>",
		"<script>alert(1)</script>": "",
	})

	// The UGC policy keeps safe links but strips javascript: URLs
	require.Equal(t, `<a href="https://example.com" rel="nofollow">link</a>`, sanitizeHTML(`<a href="https://example.com">link</a>`))
	require.NotContains(t, sanitizeHTML(`<a href="javascript:alert(1)">x</a>`), "javascript:")
}

func TestSanitizeText(t *testing.T) {
	runStringCases(t, sanitizeText, map[string]string{
		"":                          "",
		"plain text":                "plain text",
		"<b>bold</b>":               "bold",
		"<script>alert(1)</script>": "",
		`<a href="x">link</a>`:      "link",
	})
}

/******************************************
 * isActivityStream
 ******************************************/

func TestIsActivityStream(t *testing.T) {
	runBoolCases(t, isActivityStream, map[string]bool{
		"application/activity+json":                  true,
		"application/ld+json":                        true,
		"application/activity+json; charset=utf-8":   true,
		"application/ld+json; profile=\"https://x\"": true,
		"application/json":                           false,
		"text/html":                                  false,
		"":                                           false,
		"not a mime type at all":                     false,
		"application/activity+json extra junk":       false, // invalid mime, ParseMediaType fails
	})
}

/******************************************
 * defaultHTTPS
 ******************************************/

func TestDefaultHTTPS(t *testing.T) {
	runStringCases(t, defaultHTTPS, map[string]string{
		"http://example.com":  "http://example.com",
		"https://example.com": "https://example.com",
		"example.com":         "https://example.com",
		"":                    "https://",
		"ftp://example.com":   "https://ftp://example.com", // only http/https prefixes are respected
	})
}

/******************************************
 * withContext
 ******************************************/

func TestWithContext(t *testing.T) {

	// An empty map gets the default ActivityStreams @context
	value := mapof.NewAny()
	withContext(value)
	require.Equal(t, vocab.ContextTypeActivityStreams, value[vocab.AtContext])

	// An existing @context is preserved
	value = mapof.Any{vocab.AtContext: "https://example.com/custom"}
	withContext(value)
	require.Equal(t, "https://example.com/custom", value[vocab.AtContext])
}

/******************************************
 * iconSizesAsInt
 ******************************************/

func TestIconSizesAsInt(t *testing.T) {
	runIntCases(t, iconSizesAsInt, map[string]int{
		"":            0,
		"128x128":     128,
		"128X128":     128, // uppercase X is lowered
		"16x16 32x32": 32,  // largest wins
		"any":         0,   // non-numeric
		"48":          48,  // single dimension, no "x"
		"32x32 any":   32,  // mixed valid + invalid
		"100x200":     100, // first number in each dimension; "100" before the "x"
	})
}

/******************************************
 * iconMediaTypeAsInt
 ******************************************/

func TestIconMediaTypeAsInt(t *testing.T) {
	runIntCases(t, iconMediaTypeAsInt, map[string]int{
		"image/webp":               256,
		"image/png":                255,
		"image/jpg":                254,
		"image/jpeg":               253,
		"image/svg":                252,
		"image/svg+xml":            251,
		"image/gif":                250,
		"image/bmp":                248,
		"image/tiff":               247,
		"image/tiff+xml":           246,
		"image/x-icon":             245,
		"image/vnd.microsoft.icon": 244,
		"image/unknown":            0,
		"":                         0,
	})

	// WebP should always sort above PNG, which should sort above the unknown default.
	require.Greater(t, iconMediaTypeAsInt("image/webp"), iconMediaTypeAsInt("image/png"))
	require.Greater(t, iconMediaTypeAsInt("image/png"), iconMediaTypeAsInt("image/unknown"))
}

/******************************************
 * sortImageLinks
 ******************************************/

func TestSortImageLinks(t *testing.T) {

	link := func(size string, mediaType string) digit.Link {
		return digit.Link{
			MediaType:  mediaType,
			Properties: map[string]string{"sizes": size},
		}
	}

	// Larger image sorts "greater" (positive result)
	require.Positive(t, sortImageLinks(link("64x64", "image/png"), link("32x32", "image/png")))
	require.Negative(t, sortImageLinks(link("32x32", "image/png"), link("64x64", "image/png")))

	// Same size: media type breaks the tie (webp > png)
	require.Positive(t, sortImageLinks(link("64x64", "image/webp"), link("64x64", "image/png")))

	// Fully identical links compare equal
	require.Zero(t, sortImageLinks(link("64x64", "image/png"), link("64x64", "image/png")))
}

/******************************************
 * hostOnly (error path)
 ******************************************/

func TestHostOnly_Invalid(t *testing.T) {
	// A control character makes url.Parse fail; the original value is returned unchanged.
	require.Equal(t, "://\x7f", hostOnly("://\x7f"))
}

// NOTE: hostOnly's error branch is exercised by passing a URL that contains an
// invalid control character. derp.Report writes the error to stderr, which is
// expected (and noisy) during this test.

/******************************************
 * identifierType
 ******************************************/

func TestIdentifierType_Extra(t *testing.T) {
	runStringCases(t, identifierType, map[string]string{
		"https://example.com/@user": IdentifierTypeURL,
		"http://example.com":        IdentifierTypeURL,
		"@user@example.com":         IdentifierTypeUsername,
		// "example.com" has no scheme, so IsValidURL is false; it falls through to
		// IsValidAddress (which prepends https://) and is reported as a username.
		"example.com":              IdentifierTypeUsername,
		"":                         IdentifierTypeNone,
		"@@@":                      IdentifierTypeNone,
		"not a valid identifier!!": IdentifierTypeNone,
	})
}
