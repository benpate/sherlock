package sherlock

import (
	"encoding/json"
	"net/url"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/first"
	"github.com/benpate/rosetta/html"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/slice"
	"github.com/kr/jsonfeed"
)

func (client Client) loadActor_Feed_JSON(config Config, txn *remote.Transaction) streams.Document {

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

	actorID := first.String(feed.FeedURL, txn.RequestURL())
	username := first.String(feed.HomePageURL, txn.RequestURL())
	baseURL, _ := url.Parse(actorID)

	// Create an ActivityStream document
	result := config.DefaultValue
	result[vocab.AtContext] = vocab.ContextTypeActivityStreams
	result[vocab.PropertyID] = actorID
	result[vocab.PropertyType] = vocab.ActorTypeApplication
	result[vocab.PropertyName] = feed.Title
	result[vocab.PropertyIcon] = feed.Icon
	result[vocab.PropertySummary] = feed.Description
	result[vocab.PropertyURL] = username
	result[vocab.PropertyOutbox] = mapof.Any{
		vocab.PropertyType:       vocab.CoreTypeOrderedCollection,
		vocab.PropertyTotalItems: len(feed.Items),
		vocab.PropertyOrderedItems: slice.Map(feed.Items, func(item jsonfeed.Item) mapof.Any {

			itemURL, _ := baseURL.Parse(item.URL)

			return mapof.Any{
				vocab.PropertyType:         vocab.ObjectTypePage,
				vocab.PropertyID:           itemURL,
				vocab.PropertyActor:        feed.FeedURL,
				vocab.PropertyName:         item.Title,
				vocab.PropertySummary:      item.Summary,
				vocab.PropertyImage:        item.Image,
				vocab.PropertyContent:      jsonFeedToContentHTML(item),
				vocab.PropertyPublished:    item.DatePublished.Unix(),
				vocab.PropertyAttributedTo: jsonFeedToAuthor(feed, item),
			}
		}),
	}

	// Search for WebSub hubs.
	for _, hub := range feed.Hubs {
		if hub.Type == "WebSub" {
			result[vocab.PropertyEndpoints] = mapof.Any{
				"hub": hub.URL,
			}
			break
		}
	}

	// Apply links found in the response headers
	client.applyLinks(txn, result)

	// Patch icon into the feed (if necessary)
	client.loadActor_Feed_FindHomePageIcon(result)

	// Find/Manufacture the icon for the feed
	// client.loadActor_Feed_Icon(txn, result)

	return streams.NewDocument(
		result,
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
