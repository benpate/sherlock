package sherlock

import (
	"slices"
	"strings"

	"github.com/benpate/digit"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/slice"
	"github.com/rs/zerolog/log"
)

// loadActor_Feed_FindRootLevelIcon searches for an icon from the website homepage
// and adds it to the document if found.
func (client *Client) loadActor_Feed_FindHomePageIcon(document map[string]any) {

	// if the document already has an icon, then NOOP
	if icon := convert.String(document[vocab.PropertyIcon]); icon != "" {
		return
	}

	// Get the document ID from the document
	documentID := convert.String(document[vocab.PropertyID])
	documentID = hostOnly(documentID)

	// Get the root-level document from the server
	txn := remote.Get(documentID)

	if err := txn.Send(); err != nil {
		log.Error().Err(err).Str("documentID", documentID).Msg("Error sending request")
		return
	}

	// Find Icons and apply them to the document
	client.loadActor_Feed_FindIcon(txn, document)
}

// loadActor_Feed_FindIcon searches for an icon from the remote transaction and
// adds it into the document if found.
func (client *Client) loadActor_Feed_FindIcon(txn *remote.Transaction, document map[string]any) {

	// if the document already has an icon, then NOOP
	if icon := convert.String(document[vocab.PropertyIcon]); icon != "" {
		return
	}

	// Find all links in the root-level document
	links := client.loadActor_DiscoverLinks(txn)

	// Choose the best icon and add it to the result
	if icon := client.loadActor_Feed_FindIconLink(links); icon != "" {
		document[vocab.PropertyIcon] = icon
	}
}

// Search for a sitewide Favicon, and add it to the default document if found
func (client *Client) loadActor_Feed_FindIconLink(links digit.LinkSet) string {

	// Find all icon links
	icons := slice.Filter(links, func(link digit.Link) bool {
		return strings.Contains(link.RelationType, "icon")
	})

	// Empty results are empty
	if len(icons) == 0 {
		return ""
	}

	// Find the "best" icon, and set it as the default value
	slices.SortFunc(icons, sortImageLinks)
	return icons[0].Href

	// Are there other kinds of icons we can search for?
}
