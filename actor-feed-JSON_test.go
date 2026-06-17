package sherlock

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/kr/jsonfeed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsJSONFeedContentType(t *testing.T) {
	cases := map[string]bool{
		ContentTypeJSONFeed:    true,
		ContentTypeJSON:        true,
		ContentTypeActivityPub: false,
		ContentTypeHTML:        false,
		"":                     false,
		// Exact match only -- a charset suffix is NOT recognized.
		"application/json; charset=utf-8": false,
	}

	for input, expected := range cases {
		assert.Equal(t, expected, isJSONFeedContentType(input), "input: %q", input)
	}
}

func TestJSONFeedToAuthor(t *testing.T) {

	// Item author takes precedence over feed author
	{
		feed := jsonfeed.Feed{Author: &jsonfeed.Author{Name: "Feed Author"}}
		item := jsonfeed.Item{Author: &jsonfeed.Author{
			Name:   "Item Author",
			URL:    "https://example.com/item-author",
			Avatar: "https://example.com/avatar.png",
		}}

		result := jsonFeedToAuthor(feed, item)
		require.Equal(t, "Item Author", result[vocab.PropertyName])
		require.Equal(t, "https://example.com/item-author", result[vocab.PropertyID])
		require.Equal(t, "https://example.com/avatar.png", result[vocab.PropertyImage])
	}

	// Falls back to feed author when item has none
	{
		feed := jsonfeed.Feed{Author: &jsonfeed.Author{
			Name: "Feed Author",
			URL:  "https://example.com/feed-author",
		}}
		item := jsonfeed.Item{}

		result := jsonFeedToAuthor(feed, item)
		require.Equal(t, "Feed Author", result[vocab.PropertyName])
		require.Equal(t, "https://example.com/feed-author", result[vocab.PropertyID])
	}

	// With no author anywhere, falls back to the feed URL as the ID
	{
		feed := jsonfeed.Feed{FeedURL: "https://example.com/feed.json"}
		item := jsonfeed.Item{}

		result := jsonFeedToAuthor(feed, item)
		require.Equal(t, "https://example.com/feed.json", result[vocab.PropertyID])
		require.NotContains(t, result, vocab.PropertyName)
	}
}

func TestJSONFeedToContentHTML(t *testing.T) {

	// ContentHTML is preferred (and sanitized)
	{
		item := jsonfeed.Item{
			ContentHTML: "<p>Hello <script>alert(1)</script></p>",
			ContentText: "ignored",
		}
		result := jsonFeedToContentHTML(item)
		assert.Contains(t, result, "<p>Hello")
		assert.NotContains(t, result, "<script>")
	}

	// Falls back to ContentText (converted to HTML)
	{
		item := jsonfeed.Item{ContentText: "line one"}
		result := jsonFeedToContentHTML(item)
		assert.Contains(t, result, "line one")
	}

	// Empty item yields empty content
	{
		require.Equal(t, "", jsonFeedToContentHTML(jsonfeed.Item{}))
	}
}
