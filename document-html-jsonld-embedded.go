package sherlock

import (
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
	"github.com/benpate/rosetta/mapof"
)

// loadDocument_JSONLD_Embedded searches the GoQuery document for links to ActivityPub-like documents.
func (client *Client) loadDocument_JSONLD_Embedded(document *goquery.Document, result mapof.Any) bool {

	var success bool
	selection := document.Find("script[type=application/ld+json]")

	selection.EachWithBreak(func(_ int, script *goquery.Selection) bool {

		if err := json.Unmarshal([]byte(script.Text()), &result); err == nil {
			success = true
			return false // break
		}

		return true // continue
	})

	return success
}
