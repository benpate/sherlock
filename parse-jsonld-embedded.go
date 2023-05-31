package sherlock

import (
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	"github.com/benpate/rosetta/mapof"
)

// ParseEmbeddedJSONLD searches the GoQuery document for links to ActivityPub-like documents.
func ParseEmbeddedJSONLD(document *goquery.Document) (mapof.Any, bool) {
	// TODO: LOW: Add support for JSON-LD metadata embedded in a <script> tag
	// This may be a way to extract the JSON-LD metadata
	// https://pkg.go.dev/github.com/daetal-us/getld#section-readme

	selection := document.Find("script[type=application/ld+json]")

	result := mapof.NewAny()
	var success bool

	selection.EachWithBreak(func(_ int, script *goquery.Selection) bool {

		if err := json.Unmarshal([]byte(script.Text()), &result); err == nil {
			success = true
			return false // break
		}

		return true // continue
	})

	return result, success
}
