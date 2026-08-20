package sherlock

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/benpate/rosetta/mapof"
)

// loadDocument_JSONLD_Linked searches the GoQuery document for links to ActivityPub-like documents.
func (client *Client) loadDocument_JSONLD_Linked(config Config, baseURL string, document *goquery.Document, result mapof.Any) bool {

	var success bool
	selection := document.Find("link[rel=alternate]")

	selection.EachWithBreak(func(_ int, link *goquery.Selection) bool {

		if linkType, ok := link.Attr("type"); ok && isActivityStream(linkType) {

			if linkHref, ok := link.Attr("href"); ok {

				// The href comes from the page we just fetched, so it is remote
				// input: resolve it against the page's own URL (a relative
				// alternate is legal, and never resolved into an absolute URL
				// before) and send it under the same fetch policy as everything else.
				transaction := config.newTransaction(getRelativeURL(baseURL, linkHref)).
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
