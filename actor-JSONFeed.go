package sherlock

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/html"
	"github.com/benpate/rosetta/mapof"
	"github.com/kr/jsonfeed"
)

func (client Client) actor_JSONFeed(acc *actorAccumulator) bool {

	// JSONFeed content only
	if !isJSONFeedContentType(acc.Header(ContentType)) {
		return false
	}

	var feed jsonfeed.Feed

	// Parse the JSON feed
	if err := json.Unmarshal(acc.body.Bytes(), &feed); err != nil {
		return false
	}

	// Before inserting, sort the items chronologically so that new feeds appear correctly in the UX
	sort.SliceStable(feed.Items, func(i, j int) bool {
		return feed.Items[i].DatePublished.Unix() < feed.Items[j].DatePublished.Unix()
	})

	// Create an ActivityStream document
	orderedItems := make([]any, 0, len(feed.Items))

	for _, item := range feed.Items {
		orderedItems = append(orderedItems, mapof.Any{
			vocab.PropertyID:           item.URL,
			vocab.PropertyName:         item.Title,
			vocab.PropertySummary:      item.Summary,
			vocab.PropertyImage:        item.Image,
			vocab.PropertyContent:      jsonFeedToContentHTML(item),
			vocab.PropertyPublished:    item.DatePublished.Unix(),
			vocab.PropertyAttributedTo: jsonFeedToAuthor(feed, item),
		})
	}

	actor := mapof.Any{
		vocab.AtContext:       vocab.ContextTypeActivityStreams,
		vocab.PropertyType:    vocab.ActorTypeService,
		vocab.PropertyName:    feed.Title,
		vocab.PropertySummary: feed.Description,
		vocab.PropertyURL:     feed.HomePageURL,
		vocab.PropertyOutbox: mapof.Any{
			vocab.PropertyType:         vocab.CoreTypeOrderedCollection,
			vocab.PropertyTotalItems:   len(feed.Items),
			vocab.PropertyOrderedItems: orderedItems,
		},
	}

	acc.format = FormatJSONFeed
	acc.result = actor
	acc.cacheControl = "max-age=86400, public" // Force JSON feeds to cache for 1 day

	// Scan for WebSub hubs
	for _, hub := range feed.Hubs {
		if strings.ToLower(hub.Type) == "websub" {
			acc.webSub = hub.URL
			break
		}
	}

	// Success!
	return true
}

// Returns TRUE if the contentType is application/activity+json or application/ld+json
func isJSONFeedContentType(contentType string) bool {

	switch contentType {

	case ContentTypeJSONFeed:
		return true

	case ContentTypeJSON:
		return true

	default:
		return false
	}
}

func jsonFeedToAuthor(feed jsonfeed.Feed, item jsonfeed.Item) mapof.Any {

	if item.Author != nil {
		return mapof.Any{
			vocab.PropertyID:    item.Author.URL,
			vocab.PropertyName:  item.Author.Name,
			vocab.PropertyImage: item.Author.Avatar,
		}
	}

	if feed.Author != nil {
		return mapof.Any{
			vocab.PropertyID:    feed.Author.URL,
			vocab.PropertyName:  feed.Author.Name,
			vocab.PropertyImage: feed.Author.Avatar,
		}
	}

	return mapof.Any{}
}

func jsonFeedToContentHTML(item jsonfeed.Item) string {

	var result string

	if item.ContentHTML != "" {
		result = item.ContentHTML
	} else if item.ContentText != "" {
		result = html.FromText(item.ContentText)
	}

	return SanitizeHTML(result)
}
