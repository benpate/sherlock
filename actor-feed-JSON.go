package sherlock

import (
	"encoding/json"
	"sort"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/html"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/slice"
	"github.com/kr/jsonfeed"
)

func (client Client) loadActor_Feed_JSON(txn *remote.Transaction) streams.Document {

	// JSONFeed content only
	if !isJSONFeedContentType(txn.ResponseContentType()) {
		return streams.NilDocument()
	}

	var feed jsonfeed.Feed

	body, err := txn.ResponseBody()
	if err != nil {
		return streams.NilDocument()
	}

	// Parse the JSON feed
	if err := json.Unmarshal(body, &feed); err != nil {
		return streams.NilDocument()
	}

	// Before inserting, sort the items chronologically so that new feeds appear correctly in the UX
	sort.SliceStable(feed.Items, func(i, j int) bool {
		return feed.Items[i].DatePublished.Unix() < feed.Items[j].DatePublished.Unix()
	})

	// Create an ActivityStream document
	data := mapof.Any{
		vocab.AtContext:       vocab.ContextTypeActivityStreams,
		vocab.PropertyID:      feed.HomePageURL,
		vocab.PropertyType:    vocab.ActorTypeService,
		vocab.PropertyName:    feed.Title,
		vocab.PropertySummary: feed.Description,
		vocab.PropertyURL:     feed.HomePageURL,
		vocab.PropertyOutbox: mapof.Any{
			vocab.PropertyType:       vocab.CoreTypeOrderedCollection,
			vocab.PropertyTotalItems: len(feed.Items),
			vocab.PropertyOrderedItems: slice.Map(feed.Items, func(item jsonfeed.Item) mapof.Any {
				return mapof.Any{
					vocab.PropertyType:         vocab.ObjectTypePage,
					vocab.PropertyID:           item.URL,
					vocab.PropertyActor:        feed.FeedURL,
					vocab.PropertyName:         item.Title,
					vocab.PropertySummary:      item.Summary,
					vocab.PropertyImage:        item.Image,
					vocab.PropertyContent:      jsonFeedToContentHTML(item),
					vocab.PropertyPublished:    item.DatePublished.Unix(),
					vocab.PropertyAttributedTo: jsonFeedToAuthor(feed, item),
				}
			}),
		},
	}

	// Search for WebSub hubs.
	for _, hub := range feed.Hubs {
		if hub.Type == "WebSub" {
			data[vocab.PropertyEndpoints] = mapof.Any{
				"hub": hub.URL,
			}
			break
		}
	}

	// Apply links found in the response headers
	client.applyLinks(txn, data)

	return streams.NewDocument(
		data,
		streams.WithClient(client),
		streams.WithHTTPHeader(txn.ResponseHeader()),
	)
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

	return mapof.Any{
		vocab.PropertyID: feed.FeedURL,
	}
}

func jsonFeedToContentHTML(item jsonfeed.Item) string {

	var result string

	if item.ContentHTML != "" {
		result = item.ContentHTML
	} else if item.ContentText != "" {
		result = html.FromText(item.ContentText)
	}

	return sanitizeHTML(result)
}
