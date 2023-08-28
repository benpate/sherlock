package sherlock

import (
	"bytes"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/benpate/digit"
	"github.com/benpate/sherlock/pipe"
	"github.com/tomnomnom/linkheader"
	"golang.org/x/net/html"
)

func (client Client) actor_FindLinksInHeader(acc *actorAccumulator) bool {

	headerValue := acc.httpResponse.Header.Get(HTTPHeaderLink)
	links := linkheader.Parse(headerValue)

	for _, link := range links {
		acc.links = append(acc.links, digit.Link{
			RelationType: link.Rel,
			MediaType:    link.Param("type"),
			Href:         getRelativeURL(acc.url, link.URL),
		})
	}

	return false
}

func (client Client) actor_FindLinks(acc *actorAccumulator) bool {

	// Scan the HTML document for relevant links
	newReader := bytes.NewReader(acc.body.Bytes())
	htmlDocument, err := goquery.NewDocumentFromReader(newReader)

	if err != nil {
		return false
	}

	// Get "relevant" links from the document
	links := htmlDocument.Find("[rel=alternate],[rel=self],[rel=feed],[rel=hub],[rel=icon]").Nodes

	// Add links to the accumulator
	for _, link := range links {
		acc.links = append(acc.links, digit.Link{
			RelationType: nodeAttribute(link, "rel"),
			MediaType:    nodeAttribute(link, "type"),
			Href:         getRelativeURL(acc.url, nodeAttribute(link, "href")),
		})
	}

	return false
}

// actor_ScanHTMLForWebMentions tries to load/use any linked feeds
func (client Client) actor_FollowLinks(acc *actorAccumulator) bool {

	// If there are no links, then there's nothing to do in this step
	if len(acc.links) == 0 {
		return false
	}

	// Make a list of content types and pipelines to run.  This is an array
	// so that we can run the pipelines in a specific order, based on the
	// priority of each protocol:
	// 1. ActivityPub
	// 2. JSONFeed
	// 3. Atom
	// 4. RSS.
	table := []struct {
		mediaType string
		pipe      actorAccumulatorPipe
	}{
		{ContentTypeActivityPub, actorAccumulatorPipe{client.actor_ActivityStream}},
		{ContentTypeJSONFeed, actorAccumulatorPipe{client.actor_GetHTTP_JSONFeed, client.actor_JSONFeed}},
		{ContentTypeAtom, actorAccumulatorPipe{client.actor_GetHTTP_Atom, client.actor_RSSFeed}},
		{ContentTypeRSS, actorAccumulatorPipe{client.actor_GetHTTP_RSS, client.actor_RSSFeed}},
	}

	// For each link in the accumulator, try to run a corresponding pipeline
	for _, row := range table {

		// If we have a valid link for this mime type, then run its pipeline
		if link := findSelfOrAlternateLink(acc.links, row.mediaType); !link.IsEmpty() {

			sub := newActorAccumulator(link.Href)

			if done := pipe.Run(&sub, row.pipe...); done {
				acc.result = sub.result
				acc.webSub = sub.webSub
				acc.format = sub.format
				acc.cacheControl = sub.cacheControl

				return true
			}

		}
	}

	return false
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
