package sherlock

import (
	"testing"
	"time"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/mmcdole/gofeed"
	ext "github.com/mmcdole/gofeed/extensions"
	"github.com/stretchr/testify/require"
)

func TestFeedSummary(t *testing.T) {
	// Description is sanitized to plain text (tags stripped)
	item := &gofeed.Item{Description: "<p>Hello <b>world</b></p>"}
	require.Equal(t, "Hello world", feedSummary(item))
}

func TestFeedContent(t *testing.T) {
	// Content keeps safe HTML but strips scripts
	item := &gofeed.Item{Content: "<p>Safe</p><script>alert(1)</script>"}
	result := feedContent(item)
	require.Contains(t, result, "<p>Safe</p>")
	require.NotContains(t, result, "<script>")
}

func TestFeedImage(t *testing.T) {

	// Nil item -> nil image
	require.Nil(t, feedImage(nil))

	// Item with an explicit Image
	item := &gofeed.Item{Image: &gofeed.Image{URL: "https://example.com/pic.png", Title: "Pic"}}
	image := feedImage(item)
	require.Equal(t, vocab.ObjectTypeImage, image[vocab.PropertyType])
	require.Equal(t, "https://example.com/pic.png", image[vocab.PropertyHref])
	require.Equal(t, "Pic", image[vocab.PropertySummary])

	// Item with an image enclosure
	item = &gofeed.Item{Enclosures: []*gofeed.Enclosure{
		{URL: "https://example.com/file.pdf", Type: "application/pdf"},
		{URL: "https://example.com/photo.jpg", Type: "image/jpeg"},
	}}
	image = feedImage(item)
	require.Equal(t, "https://example.com/photo.jpg", image[vocab.PropertyHref])

	// Item with a media extension (YouTube style)
	item = &gofeed.Item{Extensions: ext.Extensions{
		"media": map[string][]ext.Extension{
			"group": {
				{Attrs: map[string]string{"medium": "image", "url": "https://example.com/yt.jpg", "width": "640", "height": "480"}},
			},
		},
	}}
	image = feedImage(item)
	require.Equal(t, "https://example.com/yt.jpg", image[vocab.PropertyHref])
	require.Equal(t, 640, image[vocab.PropertyWidth])
	require.Equal(t, 480, image[vocab.PropertyHeight])

	// Item with nothing image-like -> nil
	require.Nil(t, feedImage(&gofeed.Item{}))
}

func TestFeedAuthor(t *testing.T) {

	feed := &gofeed.Feed{Title: "Feed Title", Description: "Feed Desc"}

	// item.Author wins
	item := &gofeed.Item{Author: &gofeed.Person{Name: "Item Author", Email: "item@example.com"}}
	result := feedAuthor("https://example.com/actor", feed, item)
	require.Equal(t, "Item Author", result[vocab.PropertyName])
	require.Equal(t, "item@example.com", result[vocab.PropertySummary])
	require.Equal(t, "https://example.com/actor", result[vocab.PropertyID])

	// item.Authors[0] is next
	item = &gofeed.Item{Authors: []*gofeed.Person{{Name: "Authors[0]", Email: "a0@example.com"}}}
	result = feedAuthor("https://example.com/actor", feed, item)
	require.Equal(t, "Authors[0]", result[vocab.PropertyName])

	// feed.Author is next
	feed = &gofeed.Feed{Title: "Feed Title", Author: &gofeed.Person{Name: "Feed Author"}}
	result = feedAuthor("https://example.com/actor", feed, &gofeed.Item{})
	require.Equal(t, "Feed Author", result[vocab.PropertyName])

	// With nothing, falls back to feed Title/Description
	feed = &gofeed.Feed{Title: "Feed Title", Description: "Feed Desc"}
	result = feedAuthor("https://example.com/actor", feed, &gofeed.Item{})
	require.Equal(t, "Feed Title", result[vocab.PropertyName])
	require.Equal(t, "Feed Desc", result[vocab.PropertySummary])
}

func TestFeedAuthor_Image(t *testing.T) {

	// feed.Image is used when present
	feed := &gofeed.Feed{Image: &gofeed.Image{URL: "https://example.com/feed.png"}}
	result := feedAuthor("https://example.com/actor", feed, &gofeed.Item{})
	require.Equal(t, "https://example.com/feed.png", result[vocab.PropertyImage])
}

func TestFeedActivity(t *testing.T) {

	published := time.Date(2024, 1, 2, 15, 4, 5, 0, time.UTC)
	feed := &gofeed.Feed{FeedLink: "https://example.com/feed.xml", Title: "My Feed"}

	build := feedActivity("https://example.com/actor", feed)

	item := &gofeed.Item{
		Title:           "Post Title",
		Link:            "/posts/1", // relative -- should resolve against the actor ID
		Content:         "<p>Body</p>",
		Description:     "Summary text",
		PublishedParsed: &published,
	}

	result := build(item).(mapof.Any)
	require.Equal(t, vocab.ObjectTypePage, result[vocab.PropertyType])
	require.Equal(t, "https://example.com/posts/1", result[vocab.PropertyID])
	require.Equal(t, "Post Title", result[vocab.PropertyName])
	require.Equal(t, "https://example.com/feed.xml", result[vocab.PropertyActor])
	require.Equal(t, published.Unix(), result[vocab.PropertyPublished])
	require.Contains(t, result[vocab.PropertyContent], "<p>Body</p>")
	require.Equal(t, "Summary text", result[vocab.PropertySummary])
}

func TestFeedActivity_DefaultPublished(t *testing.T) {

	feed := &gofeed.Feed{FeedLink: "https://example.com/feed.xml"}
	build := feedActivity("https://example.com/actor", feed)

	before := time.Now().Unix()
	result := build(&gofeed.Item{Title: "No Date", Link: "https://example.com/x"}).(mapof.Any)
	after := time.Now().Unix()

	// With no PublishedParsed, the activity is dated "now"
	published := result[vocab.PropertyPublished].(int64)
	require.GreaterOrEqual(t, published, before)
	require.LessOrEqual(t, published, after)
}
