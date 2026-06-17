package sherlock

import (
	"testing"
	"time"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
	"willnorris.com/go/microformats"
)

func TestMicroformatProperty(t *testing.T) {

	entry := &microformats.Microformat{
		Properties: map[string][]any{
			"name":  {"Hello World"},
			"url":   {"https://example.com/post"},
			"empty": {},
			"nested": {
				&microformats.Microformat{Value: "nested-value"},
			},
		},
	}

	// Simple string property
	require.Equal(t, "Hello World", microformat_Property(entry, "name"))

	// Nested microformat returns its Value
	require.Equal(t, "nested-value", microformat_Property(entry, "nested"))

	// Missing property returns empty string
	require.Equal(t, "", microformat_Property(entry, "missing"))

	// Empty property slice returns empty string
	require.Equal(t, "", microformat_Property(entry, "empty"))

	// Multiple names: returns the first one that exists
	require.Equal(t, "https://example.com/post", microformat_Property(entry, "missing", "url"))

	// A nil entry returns empty string (boundary case)
	require.Equal(t, "", microformat_Property(nil, "name"))
}

func TestMicroformatFirst(t *testing.T) {

	mf := &microformats.Microformat{Value: "direct"}

	// Direct *Microformat is returned as-is
	require.Same(t, mf, microformat_First(mf))

	// A slice returns its first element
	require.Same(t, mf, microformat_First([]any{mf}))

	// An empty slice returns nil
	require.Nil(t, microformat_First([]any{}))

	// An unrelated type returns nil
	require.Nil(t, microformat_First("string"))
	require.Nil(t, microformat_First(nil))
}

func TestMicroformatAuthor(t *testing.T) {

	// A nil entry returns an empty (but non-nil) map
	result := microformat_Author(nil)
	require.NotNil(t, result)
	require.Empty(t, result)

	// A populated entry maps url/name/photo
	entry := &microformats.Microformat{
		Properties: map[string][]any{
			"url":   {"https://example.com/me"},
			"name":  {"Jane Doe"},
			"photo": {"https://example.com/me.png"},
		},
	}

	result = microformat_Author(entry)
	require.Equal(t, "https://example.com/me", result[vocab.PropertyID])
	require.Equal(t, "Jane Doe", result[vocab.PropertyName])
	require.Equal(t, "https://example.com/me.png", result[vocab.PropertyImage])
}

func TestMicroformatOutbox(t *testing.T) {

	items := []mapof.Any{
		{vocab.PropertyName: "one"},
		{vocab.PropertyName: "two"},
	}

	outbox := microformat_Outbox(items)
	require.Equal(t, vocab.CoreTypeOrderedCollection, outbox[vocab.PropertyType])
	require.Equal(t, 2, outbox[vocab.PropertyTotalItems])
	require.Len(t, outbox[vocab.PropertyOrderedItems], 2)
}

func TestMicroformatItem(t *testing.T) {

	feed := &microformats.Microformat{
		Properties: map[string][]any{
			"photo": {"https://example.com/feed-photo.png"},
		},
	}

	entry := &microformats.Microformat{
		Properties: map[string][]any{
			"url":       {"https://example.com/post-1"},
			"name":      {"Post One"},
			"summary":   {"A summary"},
			"published": {"2024-01-02T15:04:05Z"},
		},
	}

	result := microformat_Item(feed, entry)
	require.Equal(t, "https://example.com/post-1", result[vocab.PropertyID])
	require.Equal(t, "Post One", result[vocab.PropertyName])
	require.Equal(t, "A summary", result[vocab.PropertySummary])

	// Photo falls back to the feed's photo when the entry has none
	require.Equal(t, "https://example.com/feed-photo.png", result[vocab.PropertyImage])

	// Published date is parsed to a Unix timestamp
	expected, _ := time.Parse(time.RFC3339, "2024-01-02T15:04:05Z")
	require.Equal(t, expected.Unix(), result[vocab.PropertyPublished])
}

func TestMicroformatItem_DefaultPublished(t *testing.T) {

	feed := &microformats.Microformat{}
	entry := &microformats.Microformat{
		Properties: map[string][]any{
			"name": {"No Date"},
		},
	}

	result := microformat_Item(feed, entry)

	// KNOWN BUG: When no "published" date is present, microformat_Item intends to
	// default to time.Now(), but its guard (`result[vocab.PropertyPublished] == 0`)
	// compares a nil interface against 0, which is never true. So the property is
	// left unset rather than defaulted. This test documents the current behavior;
	// if the guard is fixed (e.g. to check for a missing/zero value), update this
	// expectation to assert a "now" timestamp.
	require.Nil(t, result[vocab.PropertyPublished])
}
