package sherlock

import (
	"testing"

	"github.com/benpate/digit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestNodeAttribute(t *testing.T) {

	node := &html.Node{
		Type: html.ElementNode,
		Data: "link",
		Attr: []html.Attribute{
			{Key: "rel", Val: "alternate"},
			{Key: "type", Val: "application/rss+xml"},
		},
	}

	require.Equal(t, "alternate", nodeAttribute(node, "rel"))
	require.Equal(t, "application/rss+xml", nodeAttribute(node, "type"))
	require.Equal(t, "", nodeAttribute(node, "missing"))

	// A nil node returns empty string (boundary case)
	require.Equal(t, "", nodeAttribute(nil, "rel"))
}

func TestGetRelativeURL(t *testing.T) {

	cases := []struct {
		base     string
		relative string
		expected string
	}{
		// Already absolute -- returned unchanged
		{"https://example.com", "https://other.com/feed", "https://other.com/feed"},
		{"https://example.com", "http://other.com/feed", "http://other.com/feed"},

		// Protocol-relative -- gets https:
		{"https://example.com", "//cdn.example.com/x.png", "https://cdn.example.com/x.png"},

		// Root-relative -- replaces path
		{"https://example.com/some/path", "/feed.xml", "https://example.com/feed.xml"},

		// Path-relative -- joins paths
		{"https://example.com/blog/", "feed.xml", "https://example.com/blog/feed.xml"},

		// Unparseable base falls back to the relative URL
		{"://\x7f", "relative", "relative"},
	}

	for _, c := range cases {
		assert.Equal(t, c.expected, getRelativeURL(c.base, c.relative), "base=%q relative=%q", c.base, c.relative)
	}
}

func TestFindSelfOrAlternateLink(t *testing.T) {

	links := digit.LinkSet{
		{RelationType: LinkRelationIcon, MediaType: ContentTypeRSS, Href: "https://example.com/icon"},
		{RelationType: LinkRelationSelf, MediaType: ContentTypeActivityPub, Href: "https://example.com/actor"},
		{RelationType: LinkRelationAlternate, MediaType: ContentTypeRSS, Href: "https://example.com/feed"},
	}

	// Finds a "self" link with matching media type
	found := findSelfOrAlternateLink(links, ContentTypeActivityPub)
	require.Equal(t, "https://example.com/actor", found.Href)

	// Finds an "alternate" link with matching media type
	found = findSelfOrAlternateLink(links, ContentTypeRSS)
	require.Equal(t, "https://example.com/feed", found.Href)

	// A media type that only appears on a non-self/alternate relation is NOT matched
	found = findSelfOrAlternateLink(digit.LinkSet{
		{RelationType: LinkRelationIcon, MediaType: ContentTypeRSS, Href: "https://example.com/icon"},
	}, ContentTypeRSS)
	require.True(t, found.IsEmpty())

	// No match returns an empty link
	found = findSelfOrAlternateLink(links, ContentTypeAtom)
	require.True(t, found.IsEmpty())

	// Empty link set returns an empty link
	found = findSelfOrAlternateLink(digit.LinkSet{}, ContentTypeRSS)
	require.True(t, found.IsEmpty())
}
