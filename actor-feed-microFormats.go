package sherlock

import (
	"net/url"
	"time"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/rosetta/slice"
	"willnorris.com/go/microformats"
)

// actor_MicroFormats searches and HTML document for for an h-feed Microformat
func (client Client) loadActor_Feed_MicroFormats(txn *remote.Transaction) streams.Document {

	// Parse the document URL
	parsedURL, err := url.Parse(txn.RequestURL())

	if err != nil {
		return streams.NilDocument()
	}

	// Parse the HTML document
	data := microformats.Parse(txn.ResponseBodyReader(), parsedURL)

	// Search Microformats for an h-feed
	for _, feed := range data.Items {

		if slice.Contains(feed.Type, "h-feed") {

			items := make([]mapof.Any, 0, len(feed.Children))

			for _, child := range feed.Children {
				if slice.Contains(child.Type, "h-entry") {
					items = append(items, microformat_Item(feed, child))
				}
			}

			if len(items) > 0 {

				data := mapof.Any{
					vocab.PropertyID:           parsedURL.String(),
					vocab.PropertyName:         microformat_Property(feed, "name"),
					vocab.PropertyImage:        microformat_Property(feed, "photo"),
					vocab.PropertyAttributedTo: microformat_Property(feed, "author"),
					vocab.PropertyOutbox:       microformat_Outbox(items),
				}

				// Apply links found in the response headers
				client.applyLinks(txn, data)

				// Return the (successfully?) parsed document to the caller.
				return streams.NewDocument(
					data,
					streams.WithClient(client),
					streams.WithHTTPHeader(txn.ResponseHeader()),
				)
			}
		}
	}

	return streams.NilDocument()
}

// microformat_Outbox wraps a slice of items in an ActivityStreams OrderedCollection
func microformat_Outbox(items []mapof.Any) mapof.Any {

	return mapof.Any{
		vocab.PropertyType:         vocab.CoreTypeOrderedCollection,
		vocab.PropertyTotalItems:   len(items),
		vocab.PropertyOrderedItems: items,
	}
}

// microformat_Item converts a Microformat entry into an ActivityStreams document
func microformat_Item(feed *microformats.Microformat, entry *microformats.Microformat) mapof.Any {

	result := mapof.Any{
		vocab.PropertyID:      microformat_Property(entry, "url"),
		vocab.PropertyName:    microformat_Property(entry, "name"),
		vocab.PropertySummary: microformat_Property(entry, "summary"),
	}

	// Get properties from entry

	// Get photo from entry, then feed
	if photoURL := microformat_Property(entry, "photo"); photoURL != "" {
		result[vocab.PropertyImage] = photoURL
	} else if photoURL := microformat_Property(feed, "photo"); photoURL != "" {
		result[vocab.PropertyImage] = photoURL
	}

	// Get author from entry, then feed
	if author := microformat_First(entry.Properties["author"]); author != nil {
		result[vocab.PropertyAttributedTo] = microformat_Author(author)
	} else if author := microformat_First(feed.Properties["author"]); author != nil {
		result[vocab.PropertyAttributedTo] = microformat_Author(author)
	}

	// Get the publish date from the entry
	if published := microformat_Property(entry, "published"); published != "" {
		if publishDate, err := time.Parse(time.RFC3339, published); err == nil {
			result[vocab.PropertyPublished] = publishDate.Unix()
		}
	}

	// Default PublishDate just in case
	if result[vocab.PropertyPublished] == 0 {
		result[vocab.PropertyPublished] = time.Now().Unix()
	}

	return result
}

// microformat_Author converts a Microformat entry into an ActivityStreams document
func microformat_Author(entry *microformats.Microformat) mapof.Any {

	if entry == nil {
		return mapof.NewAny()
	}

	return mapof.Any{
		vocab.PropertyID:    microformat_Property(entry, "url"),
		vocab.PropertyName:  microformat_Property(entry, "name"),
		vocab.PropertyImage: microformat_Property(entry, "photo", "logo"),
	}
}

// microformat_First returns the first item in a slice of items
func microformat_First(value any) *microformats.Microformat {

	switch o := value.(type) {
	case []any:
		if len(o) > 0 {
			return microformat_First(o[0])
		}

	case *microformats.Microformat:
		return o
	}

	return nil
}

// microformat_Property returns the first value of a property
func microformat_Property(entry *microformats.Microformat, names ...string) string {

	if entry == nil {
		return ""
	}

	for _, name := range names {

		if value, ok := entry.Properties[name]; ok {

			for _, item := range value {
				switch o := item.(type) {
				case string:
					return o

				case *microformats.Microformat:
					return o.Value
				}
			}
		}
	}

	return ""
}
