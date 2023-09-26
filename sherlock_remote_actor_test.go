//go:build localonly

package sherlock

import (
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/stretchr/testify/require"
)

func TestRemoteActor_JSONFeed(t *testing.T) {
	testRemoteActor(t, "https://www.jsonfeed.org")
}

func TestRemoteActor_JSONLink(t *testing.T) {
	testRemoteActor(t, "https://www.jsonfeed.org/feed.json")
}

func TestRemoteActor_RSSFeed(t *testing.T) {
	testRemoteActor(t, "https://www.smashingmagazine.com/feed")
}

func TestRemoteActor_RSSLink(t *testing.T) {
	testRemoteActor(t, "https://www.smashingmagazine.com")
}

func TestRemoteActor_WebFinger(t *testing.T) {
	result := testRemoteActor(t, "@benpate@mastodon.social")

	require.Equal(t, "https://mastodon.social/users/benpate", result.ID())
	require.Equal(t, "https://mastodon.social/users/benpate/inbox", result.Inbox().Value())
}

func testRemoteActor(t *testing.T, url string) streams.Document {
	client := Client{}
	result, err := client.Load(url, AsActor())
	require.NoError(t, err)

	outbox, err := result.Outbox().Load()
	require.NoError(t, err)
	require.Greater(t, outbox.TotalItems(), 0)

	return result
}
