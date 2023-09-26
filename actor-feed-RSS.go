package sherlock

import (
	"sort"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/html"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/slice"
	"github.com/mmcdole/gofeed"
)

// loadActor_Feed_RSS tries generate an Actor from an RSS or Atom feed
func (client Client) loadActor_Feed_RSS(txn *remote.Transaction) streams.Document {

	// Try to find the RSS feed associated with this link
	feed, err := gofeed.NewParser().Parse(txn.ResponseBodyReader())

	if err != nil {
		return streams.NilDocument()
	}

	// Sort the feed items (oldest first)
	sort.Slice(feed.Items, func(i, j int) bool {
		return feed.Items[i].PublishedParsed.Before(*feed.Items[j].PublishedParsed)
	})

	// Create JSON-LD for the Actor
	data := mapof.Any{
		vocab.AtContext:       vocab.ContextTypeActivityStreams,
		vocab.PropertyType:    vocab.ActorTypeService,
		vocab.PropertyID:      txn.RequestURL(),
		vocab.PropertyName:    feed.Title,
		vocab.PropertySummary: feed.Description,
		vocab.PropertyURL:     txn.RequestURL(),
		vocab.PropertyOutbox: mapof.Any{
			vocab.PropertyType:         vocab.CoreTypeOrderedCollection,
			vocab.PropertyTotalItems:   len(feed.Items),
			vocab.PropertyOrderedItems: slice.Map(feed.Items, feedActivity(feed)),
		},
	}

	// Return the result as a streams.Document
	return streams.NewDocument(
		data,
		streams.WithClient(client),
		streams.WithHTTPHeader(txn.ResponseHeader()),
	)
}

// feedActivity populates an Activity object from a gofeed.Feed and gofeed.Item
func feedActivity(feed *gofeed.Feed) func(*gofeed.Item) any {

	return func(item *gofeed.Item) any {

		result := mapof.Any{
			vocab.PropertyID:        item.Link,
			vocab.PropertyName:      html.ToText(item.Title),
			vocab.PropertyPublished: item.PublishedParsed.Unix(),
		}

		if imageURL := feedImage(feed, item); imageURL != "" {
			result[vocab.PropertyImage] = imageURL
		}

		if summary := feedSummary(item); summary != "" {
			result[vocab.PropertySummary] = summary
		}

		if contentHTML := feedContent(item); contentHTML != "" {
			result[vocab.PropertyContent] = contentHTML
		}

		if attributedTo := feedAuthor(feed, item); attributedTo != nil {
			result[vocab.PropertyAttributedTo] = attributedTo
		}

		return result
	}
}

func feedAuthor(feed *gofeed.Feed, item *gofeed.Item) mapof.Any {

	result := mapof.Any{}

	// Try to find the image from the feed.  It's weird, but easier this way.
	if feed.Image != nil {
		result[vocab.PropertyImage] = feed.Image.URL

	} else if webfeeds, ok := feed.Extensions["webfeeds"]; ok {
		if icon, ok := webfeeds["icon"]; ok {
			for _, element := range icon {
				if element.Name == "icon" {
					result[vocab.PropertyImage] = element.Value
					break
				}
			}
		}
	}

	// Try to find the author from various sources in the item
	if item.Author != nil {
		result[vocab.PropertyName] = html.ToText(item.Author.Name)
		result[vocab.PropertySummary] = item.Author.Email
		return result
	}

	if len(item.Authors) > 0 {
		if itemAuthor := item.Authors[0]; itemAuthor != nil {
			result[vocab.PropertyName] = itemAuthor.Name
			result[vocab.PropertySummary] = itemAuthor.Email
			return result
		}
	}

	// Try to find the author from various sources in the feed
	if feed.Author != nil {
		result[vocab.PropertyName] = html.ToText(feed.Author.Name)
		result[vocab.PropertySummary] = feed.Author.Email
		return result
	}

	if len(feed.Authors) > 0 {
		if feedAuthor := feed.Authors[0]; feedAuthor != nil {
			result[vocab.PropertyName] = feedAuthor.Name
			result[vocab.PropertySummary] = feedAuthor.Email
			return result
		}
	}

	// Fallback to use the Feed information as the Author
	result[vocab.PropertyName] = feed.Title
	result[vocab.PropertySummary] = feed.Description

	return result
}

// feedSummary returns a summary of the item in plain text format
func feedSummary(item *gofeed.Item) string {
	return sanitizeText(item.Description)
}

// feedContent returns a sanitized version of the HTML content for this feed
func feedContent(item *gofeed.Item) string {
	return sanitizeHTML(item.Content)
}

// rssImage returns the URL of the first image in the item's enclosure list.
func feedImage(rssFeed *gofeed.Feed, item *gofeed.Item) string {

	if item == nil {
		return ""
	}

	if item.Image != nil {
		return item.Image.URL
	}

	// Search for an image in the enclosures
	for _, enclosure := range item.Enclosures {
		if list.Slash(enclosure.Type).First() == "image" {
			return enclosure.URL
		}
	}

	// Search for media extensions (YouTube uses this)
	if media, ok := item.Extensions["media"]; ok {
		if group, ok := media["group"]; ok {
			for _, extension := range group {
				if thumbnails, ok := extension.Children["thumbnail"]; ok {
					for _, item := range thumbnails {
						if url := item.Attrs["url"]; url != "" {
							return url
						}
					}
				}
			}
		}
	}

	return ""
}
