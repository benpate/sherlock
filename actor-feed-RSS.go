package sherlock

import (
	"net/url"
	"sort"
	"time"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/convert"
	"github.com/benpate/rosetta/first"
	"github.com/benpate/rosetta/html"
	"github.com/benpate/rosetta/list"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/slice"
	"github.com/mmcdole/gofeed"
)

// loadActor_Feed_RSS tries generate an Actor from an RSS or Atom feed
func (client Client) loadActor_Feed_RSS(config Config, txn *remote.Transaction) streams.Document {

	// Try to find the RSS feed associated with this link
	feed, err := gofeed.NewParser().Parse(txn.ResponseBodyReader())

	if err != nil {
		return streams.NilDocument()
	}

	// Sort the feed items (oldest first)
	sort.Slice(feed.Items, func(i, j int) bool {
		if firstPublishDate := feed.Items[i].PublishedParsed; firstPublishDate != nil {
			if secondPublishDate := feed.Items[j].PublishedParsed; secondPublishDate != nil {
				return firstPublishDate.Before(*secondPublishDate)
			}
			return false
		}
		return false
	})

	actorID := first.String(feed.FeedLink, feed.Link, txn.RequestURL())

	// Create JSON-LD for the Actor
	result := config.DefaultValue
	result[vocab.AtContext] = vocab.ContextTypeActivityStreams
	result[vocab.PropertyType] = vocab.ActorTypeApplication
	result[vocab.PropertyID] = actorID
	result[vocab.PropertyName] = feed.Title
	result[vocab.PropertySummary] = feed.Description
	result[vocab.PropertyURL] = txn.RequestURL()
	result[vocab.PropertyOutbox] = mapof.Any{
		vocab.PropertyType:         vocab.CoreTypeOrderedCollection,
		vocab.PropertyTotalItems:   len(feed.Items),
		vocab.PropertyOrderedItems: slice.Map(feed.Items, feedActivity(actorID, feed)),
	}

	// Apply links found in the response headers
	client.applyLinks(txn, result)

	// Patch icon into the feed (if necessary)
	client.loadActor_Feed_FindHomePageIcon(result)

	// Return the result as a streams.Document
	return streams.NewDocument(
		result,
		streams.WithClient(client),
		streams.WithHTTPHeader(txn.ResponseHeader()),
	)
}

// feedActivity populates an Activity object from a gofeed.Feed and gofeed.Item
func feedActivity(actorID string, feed *gofeed.Feed) func(*gofeed.Item) any {

	baseURL, _ := url.Parse(actorID)

	return func(item *gofeed.Item) any {

		// Resolve relative URLs
		linkURL, _ := baseURL.Parse(item.Link)

		result := mapof.Any{
			vocab.PropertyType:  vocab.ObjectTypePage,
			vocab.PropertyID:    linkURL.String(),
			vocab.PropertyName:  html.ToText(item.Title),
			vocab.PropertyActor: feed.FeedLink,
		}

		if item.PublishedParsed != nil {
			result[vocab.PropertyPublished] = item.PublishedParsed.Unix()
		} else {
			result[vocab.PropertyPublished] = time.Now().Unix()
		}

		if image := feedImage(item); image != nil {
			result[vocab.PropertyImage] = image
		}

		if summary := feedSummary(item); summary != "" {
			result[vocab.PropertySummary] = summary
		}

		if contentHTML := feedContent(item); contentHTML != "" {
			result[vocab.PropertyContent] = contentHTML
		}

		if attributedTo := feedAuthor(actorID, feed, item); attributedTo != nil {
			result[vocab.PropertyAttributedTo] = attributedTo
		}

		return result
	}
}

func feedAuthor(actorID string, feed *gofeed.Feed, item *gofeed.Item) mapof.Any {

	// Set up default values to override (if we find something better)
	result := mapof.Any{
		vocab.PropertyID:      actorID,
		vocab.PropertyName:    feed.Title,
		vocab.PropertySummary: feed.Description,
	}

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
func feedImage(item *gofeed.Item) map[string]any {

	if item == nil {
		return nil
	}

	if item.Image != nil {
		return map[string]any{
			vocab.PropertyType:    vocab.ObjectTypeImage,
			vocab.PropertyHref:    item.Image.URL,
			vocab.PropertySummary: item.Image.Title,
		}
	}

	// Search for an image in the enclosures
	for _, enclosure := range item.Enclosures {
		if list.Slash(enclosure.Type).First() == "image" {
			return map[string]any{
				vocab.PropertyType: vocab.ObjectTypeImage,
				vocab.PropertyHref: enclosure.URL,
			}
		}
	}

	// Search for media extensions (YouTube uses this)
	if media, ok := item.Extensions["media"]; ok {
		for _, group := range media {
			for _, extension := range group {
				if medium := extension.Attrs["medium"]; medium == "image" {
					return map[string]any{
						vocab.PropertyType:   vocab.ObjectTypeImage,
						vocab.PropertyHref:   extension.Attrs["url"],
						vocab.PropertyWidth:  convert.Int(extension.Attrs["width"]),
						vocab.PropertyHeight: convert.Int(extension.Attrs["height"]),
					}
				}
			}
		}
	}

	return nil
}
