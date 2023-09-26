package sherlock

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

// loadDocument_JSONLD_Linked searches the GoQuery document for links to ActivityPub-like documents.
func (client *Client) loadDocument_JSONLD_Linked(document *goquery.Document, result mapof.Any) bool {
	// TODO: LOW: Add support for JSON-LD metadata embedded in a <script> tag
	// This may be a way to extract the JSON-LD metadata
	// https://pkg.go.dev/github.com/daetal-us/getld#section-readme

	var success bool
	selection := document.Find("link[rel=alternate]")

	selection.EachWithBreak(func(_ int, link *goquery.Selection) bool {

		if linkType, ok := link.Attr("type"); ok && isActivityStream(linkType) {

			if linkHref, ok := link.Attr("href"); ok {

				transaction := remote.
					Get(linkHref).
					Header("Accept", linkType).
					Result(&result)

				if err := transaction.Send(); err == nil {
					success = true
					return false // break
				}
			}
		}

		return true // continue
	})

	return success
}
