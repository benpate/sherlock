package sherlock

import (
	"bytes"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// Parse searches for all of the metadata available in a document,
// including OpenGraph, MicroFormats, and JSON-LD.
// Any metadata that is found is added to the default result object.
func Parse(target string, body *bytes.Buffer, result mapof.Any) error {

	const location = "sherlock.Parse"

	// Validate the URL
	parsedURL, err := url.Parse(target)

	if err != nil {
		return derp.Wrap(err, location, "Error parsing URL", target)
	}

	// Search the returned HTML for JSON-LD
	bodyBytes := body.Bytes()

	if gqDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes)); err == nil {

		if ParseEmbeddedJSONLD(gqDoc, result) {
			return withContext(result)
		}

		if ParseLinkedJSONLD(gqDoc, result) {
			return withContext(result)
		}
	}

	// Try OpenGraph (via HTMLInfo)
	result = ParseOpenGraph(target, bytes.NewReader(bodyBytes), result)

	// Try Microformats2
	result = ParseMicroFormats(parsedURL, bytes.NewReader(bodyBytes), result)

	// If we have SOMETHING to work with, then call it here.
	if IsAdequite(result) {
		return withContext(result)
	}

	// Otherwise, look for an oEmbed provider for this document
	if ok := ParseOEmbed(bytes.NewReader(bodyBytes), result); ok {
		return withContext(result)
	}

	// If the result is STILL EMPTY, then we have failed.
	if len(result) == 0 {
		return derp.NewNotFoundError("sherlock.Parse", "No metadata found in document")
	}

	// Yippe-Ki-Yay!
	return withContext(result)
}

// IsAdequite returns TRUE if this JSON-LD document includes the minimum fields that we'd really like
// to include in our app.  It doesn't mean the document is "valid" or "complete", but it does mean that
// we can probably do something useful with it.
func IsAdequite(value mapof.Any) bool {

	// RULE: MUST have a Property ID (url) to be "Adequite"
	if value.GetString(vocab.PropertyID) == "" {
		return false
	}

	// RULE: MUST have a Property Type to be "Adequite"
	if value.GetString(vocab.PropertyType) == "" {
		return false
	}

	// RULE: MUST have a Property Name or Summary to be "Adequite"
	if value.GetString(vocab.PropertyName)+value.GetString(vocab.PropertySummary) == "" {
		return false
	}

	// RULE: MUST have a Property Image to be "Adequite"
	if value.GetString(vocab.PropertyImage) == "" {
		return false
	}

	// RULE: MUST have a Property AttributedTo to be "Adequite"
	if value.GetAny(vocab.PropertyAttributedTo) == nil {
		return false
	}

	// Otherwise, this JSON-LD is not (yet) "Adequite".  Keep looking to see if we can do better.
	return true
}

// withContext adds the standard ActivityStreams @context to the JSON-LD document.
// If we're doing this, it's because we're assembling a "fake" JSON-LD document out of
// other metadata (like OpenGraph, MicroFormats, oEmbed, etc).
func withContext(value mapof.Any) error {
	if _, ok := value["@context"]; !ok {
		value["@context"] = vocab.ContextTypeActivityStreams
	}
	return nil
}
