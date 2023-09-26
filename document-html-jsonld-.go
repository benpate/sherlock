package sherlock

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

func (client *Client) loadDocument_JSONLD(body []byte, result map[string]any) {

	// Search the returned HTML for JSON-LD
	if gqDoc, err := goquery.NewDocumentFromReader(bytes.NewReader(body)); err == nil {

		if client.loadDocument_JSONLD_Embedded(gqDoc, result) {
			withContext(result)
			return
		}

		if client.loadDocument_JSONLD_Linked(gqDoc, result) {
			withContext(result)
			return
		}
	}
}
