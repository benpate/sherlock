package sherlock

import (
	"bytes"
	"net/url"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
)

// ParseHTML searches for all of the metadata available in a document,
// including OpenGraph, MicroFormats, and JSON-LD.
func ParseHTML(target string, body *bytes.Buffer) (mapof.Any, error) {

	const location = "sherlock.Parse"

	// Validate the URL
	parsedURL, err := url.Parse(target)

	if err != nil {
		return nil, derp.Wrap(err, location, "Error parsing URL", target)
	}

	result := mapof.Any{}

	bodyBytes := body.Bytes()

	// Try OpenGraph (via HTMLInfo)
	result = ParseOpenGraph(target, bytes.NewReader(bodyBytes), result)

	// Try Microformats2
	result = ParseMicroFormats(parsedURL, bytes.NewReader(bodyBytes), result)

	// Look for JSON-LD embedded in the docuemnt
	result = ParseJSONLD(body, result)

	// If we have SOMETHING to work with, then call it here.
	if IsAdequite(result) {
		return withContext(result), nil
	}

	// Otherwise, look links to an external JSON-LD result
	if ok := ParseLinkedJSONLD(body, result); ok {
		return withContext(result), nil
	}

	// Otherwise, look for an oEmbed provider for this document
	if ok := ParseOEmbed(bytes.NewReader(bodyBytes), result); ok {
		return withContext(result), nil
	}

	return result, derp.NewNotFoundError("sherlock.Parse", "No metadata found in document")
}

func withContext(value mapof.Any) mapof.Any {
	value["@context"] = vocab.ContextTypeActivityStreams
	return value
}

func IsAdequite(value mapof.Any) bool {

	if value.GetString(vocab.PropertyID) == "" {
		return false
	}

	if value.GetString(vocab.PropertyType) == "" {
		return false
	}

	return true
}
