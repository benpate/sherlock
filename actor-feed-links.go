package sherlock

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/benpate/digit"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/remote"
	"github.com/tomnomnom/linkheader"
	"golang.org/x/net/html"
)

// loadActor_Links finds and follows all relevant links for an http.Response.
// If it finds a link to an ActivityStream, RSS Feed, or similar, then it returns
// the corresponding Actor document.
// Otherwise, it returns an empty streams.Document that includes metadata for
func (client *Client) loadActor_Links(txn *remote.Transaction, config *LoadConfig) streams.Document {

	// Extranct all Links from the HTTP Header and HTML Document
	links := client.loadActor_DiscoverLinks(txn)

	// If links point directly to something we can use (ActivityPub, RSS, etc) then use it
	if document := client.loadActor_FollowLinks(txn, links, config); document.NotNil() {
		return document
	}

	// Otherwise, populate additional links (such as Hubs, Icons, etc)
	// and return an empty streams.Document
	// TODO: https://trello.com/c/t51YFiA2/234-sherlock-restore-websub-links
	return streams.NilDocument()
}

// loadActor_DiscoverLinks finds all links in a transaction, from both the
// http header and in the HTML document.
func (client *Client) loadActor_DiscoverLinks(txn *remote.Transaction) digit.LinkSet {

	// Retrieve Links in HTTP Header
	headerValue := txn.ResponseHeader().Get(HTTPHeaderLink)
	links := linkheader.Parse(headerValue)
	result := make(digit.LinkSet, 0, len(links))
	requestURL := txn.RequestURL()

	for _, link := range links {
		result = append(result, digit.Link{
			RelationType: link.Rel,
			MediaType:    link.Param("type"),
			Href:         getRelativeURL(requestURL, link.URL),
		})
	}

	// Retrieve Links in HTML Document
	if htmlDocument, err := goquery.NewDocumentFromReader(txn.ResponseBodyReader()); err == nil {

		// Get "relevant" links from the document
		selection := htmlDocument.Find("[rel=alternate],[rel=self],[rel=feed],[rel=hub],[rel=icon],[rel=apple-touch-icon],[rel=apple-touch-icon-precomposed],[rel=mask-icon]")

		// Add links to the accumulator
		for _, link := range selection.Nodes {
			result = append(result, digit.Link{
				RelationType: nodeAttribute(link, "rel"),
				MediaType:    nodeAttribute(link, "type"),
				Href:         getRelativeURL(requestURL, nodeAttribute(link, "href")),
				Properties: map[string]string{
					"sizes": nodeAttribute(link, "sizes"),
				},
			})
		}
	}

	return result
}

// actor_ScanHTMLForWebMentions tries to load/use any linked feeds
func (client *Client) loadActor_FollowLinks(txn *remote.Transaction, links digit.LinkSet, config *LoadConfig) streams.Document {

	// If the client is not allowed to follow redirects (or has used all of them already),
	// then there is nothing to do here. Return an empty document instead.
	if config.MaximumRedirects < 1 {
		return streams.NilDocument()
	}

	// If we have one or more links, then search them in order...
	if len(links) > 0 {

		for _, mediaType := range []string{ContentTypeActivityPub, ContentTypeJSONFeed, ContentTypeJSON, ContentTypeAtom, ContentTypeRSS} {

			link := findSelfOrAlternateLink(links, mediaType)

			if link.IsEmpty() {
				continue
			}

			// If the link points to the same URL as the original request, then we're
			// already at the right place. So don't traverse the link.
			if link.Href == txn.RequestURL() {
				return streams.NilDocument()
			}

			if document, err := client.loadActor(link.Href, config); err == nil {
				if document.NotNil() {
					config.MaximumRedirects--
					return document
				}
			}
		}
	}

	return streams.NilDocument()
}

/******************************************
 * Helper Functions
 ******************************************/

// nodeAttribute searches for a specific attribute in a node and returns its value
func nodeAttribute(node *html.Node, name string) string {

	if node == nil {
		return ""
	}

	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}

	return ""
}

// TODO: HIGH: Scan all references and perhaps use https://pkg.go.dev/net/url#URL.ResolveReference instead?
func getRelativeURL(baseURL string, relativeURL string) string {

	// If the relative URL is already absolute, then just return it
	if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
		return relativeURL
	}

	// If the relative URL is a root-relative URL, then assume HTTPS (it's 2022, for crying out loud)
	if strings.HasPrefix(relativeURL, "//") {
		return "https:" + relativeURL
	}

	// Parse the base URL so that we can do URL-math on it
	baseURLParsed, _ := url.Parse(baseURL)

	// If the relative URL is a path-relative URL, then just replace the path
	if strings.HasPrefix(relativeURL, "/") {
		baseURLParsed.Path = relativeURL
		return baseURLParsed.String()
	}

	// Otherwise, join the paths
	baseURLParsed.Path, _ = url.JoinPath(baseURLParsed.Path, relativeURL)
	return baseURLParsed.String()
}

func findSelfOrAlternateLink(links []digit.Link, mediaType string) digit.Link {

	for _, link := range links {

		switch link.RelationType {
		case LinkRelationSelf, LinkRelationAlternate:
			if link.MediaType == mediaType {
				return link
			}
		}
	}

	return digit.Link{}
}
