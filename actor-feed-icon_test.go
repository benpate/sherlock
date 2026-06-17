package sherlock

import (
	"testing"

	"github.com/benpate/digit"
	"github.com/stretchr/testify/require"
)

func TestLoadActor_Feed_FindIconLink(t *testing.T) {

	client := NewClient()

	// No links -> empty string
	require.Equal(t, "", client.loadActor_Feed_FindIconLink(digit.LinkSet{}))

	// No icon-type links -> empty string
	links := digit.LinkSet{
		{RelationType: "alternate", Href: "https://example.com/feed"},
		{RelationType: "self", Href: "https://example.com/me"},
	}
	require.Equal(t, "", client.loadActor_Feed_FindIconLink(links))

	// A single icon link is returned
	links = digit.LinkSet{
		{RelationType: "icon", Href: "https://example.com/favicon.ico"},
	}
	require.Equal(t, "https://example.com/favicon.ico", client.loadActor_Feed_FindIconLink(links))

	// Multiple icons: NOTE the code sorts ascending with sortImageLinks (smaller/worse
	// first) and then takes icons[0], so it actually returns the SMALLEST icon -- even
	// though the comment says "best". This test documents the current behavior; if the
	// sort direction is fixed, update this to expect the large icon instead.
	links = digit.LinkSet{
		{RelationType: "icon", Href: "https://example.com/small.png", MediaType: "image/png", Properties: map[string]string{"sizes": "16x16"}},
		{RelationType: "apple-touch-icon", Href: "https://example.com/large.png", MediaType: "image/png", Properties: map[string]string{"sizes": "180x180"}},
	}
	require.Equal(t, "https://example.com/small.png", client.loadActor_Feed_FindIconLink(links))
}
