package sherlock

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
)

// The fuzz tests below assert only that the parsers/decoders never panic on
// arbitrary input. They are deliberately lenient about output: the contract for
// these functions on garbage input is "degrade gracefully", not "return X".

func FuzzIconSizesAsInt(f *testing.F) {
	f.Add("")
	f.Add("128x128")
	f.Add("16x16 32x32")
	f.Add("xX x")
	f.Add("999999999999999999999999x1")

	f.Fuzz(func(t *testing.T, input string) {
		// Just assert it never panics. NOTE: the function does NOT clamp negative
		// dimensions (e.g. iconSizesAsInt("-1") == -1), since convert.IntOk happily
		// parses negative numbers. That's harmless for the sort comparator that uses
		// this, so we don't treat negative output as a failure.
		_ = iconSizesAsInt(input)
	})
}

func FuzzDefaultHTTPS(f *testing.F) {
	f.Add("")
	f.Add("example.com")
	f.Add("http://example.com")
	f.Add("https://example.com")

	f.Fuzz(func(t *testing.T, input string) {
		_ = defaultHTTPS(input)
	})
}

func FuzzIsActivityStream(f *testing.F) {
	f.Add("application/activity+json")
	f.Add("application/ld+json; charset=utf-8")
	f.Add("text/html")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		_ = isActivityStream(input)
	})
}

func FuzzGetRelativeURL(f *testing.F) {
	f.Add("https://example.com/path", "feed.xml")
	f.Add("https://example.com", "//cdn.example.com/x")
	f.Add("", "")
	f.Add("://bad", "/root")

	f.Fuzz(func(t *testing.T, base string, relative string) {
		_ = getRelativeURL(base, relative)
	})
}

func FuzzIsValidAddress(f *testing.F) {
	f.Add("@user@example.com")
	f.Add("https://example.com/@user")
	f.Add("example.com")
	f.Add("")
	f.Add("@@@@")

	f.Fuzz(func(t *testing.T, input string) {
		// Two calls should agree (no hidden state / non-determinism).
		first := IsValidAddress(input)
		second := IsValidAddress(input)
		if first != second {
			t.Errorf("IsValidAddress(%q) is not deterministic", input)
		}
	})
}

func FuzzIdentifierType(f *testing.F) {
	f.Add("https://example.com")
	f.Add("@user@example.com")
	f.Add("garbage")
	f.Add("")

	f.Fuzz(func(t *testing.T, input string) {
		result := identifierType(input)
		switch result {
		case IdentifierTypeURL, IdentifierTypeUsername, IdentifierTypeNone:
			// ok -- one of the three known values
		default:
			t.Errorf("identifierType(%q) returned unexpected value %q", input, result)
		}
	})
}

func FuzzLoadDocumentOpenGraph(f *testing.F) {
	client := NewClient()

	f.Add(`<html><head><meta property="og:title" content="Title"></head></html>`)
	f.Add(`<html>`)
	f.Add(``)
	f.Add(`not html at all <<<>>>`)

	f.Fuzz(func(t *testing.T, body string) {
		// Should never panic, regardless of how malformed the HTML is.
		client.loadDocument_OpenGraph("https://example.com", []byte(body), mapof.NewAny())
	})
}

func FuzzLoadDocumentMicroFormats(f *testing.F) {
	client := NewClient()

	f.Add(`<html><body class="h-entry"><span class="p-name">Title</span></body></html>`)
	f.Add(`<html>`)
	f.Add(``)
	f.Add(`<div class="h-feed"><div class="h-entry"></div></div>`)

	f.Fuzz(func(t *testing.T, body string) {
		client.loadDocument_MicroFormats("https://example.com", []byte(body), mapof.NewAny())
	})
}

func FuzzLoadDocumentJSONLD_Embedded(f *testing.F) {
	client := NewClient()

	f.Add(`<html><head><script type="application/ld+json">{"name":"x"}</script></head></html>`)
	f.Add(`<html><head><script type="application/ld+json">not json</script></head></html>`)
	f.Add(``)

	f.Fuzz(func(t *testing.T, body string) {
		client.loadDocument_JSONLD([]byte(body), mapof.NewAny())
	})
}
