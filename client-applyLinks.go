package sherlock

import (
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
	"github.com/tomnomnom/linkheader"
)

// applyLinks searches for common link headers in the response, and applies them to the data map
func (client *Client) applyLinks(txn *remote.Transaction, data mapof.Any) {

	links := linkheader.ParseMultiple(txn.Response().Header["Link"])

	for _, link := range links {
		switch link.Rel {

		case LinkRelationIcon:

			// Add an icon if it doesn't already exist
			if _, ok := data[vocab.PropertyIcon]; !ok {
				data[vocab.PropertyIcon] = link.URL
			}

		case LinkRelationHub:

			// Guarantee that the `endpoints` value exists
			if _, ok := data[vocab.PropertyEndpoints]; !ok {
				data[vocab.PropertyEndpoints] = make(map[string]any)
			}

			// Set the `endpoints.websub` value
			if endpoints, ok := data[vocab.PropertyEndpoints].(map[string]any); ok {
				endpoints["websub"] = link.URL
			}
		}
	}
}
