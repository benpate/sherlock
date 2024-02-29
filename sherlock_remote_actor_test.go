//go:build localonly

package sherlock

import (
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
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

func TestRemoteActor_Mitra_HTTP(t *testing.T) {
	result := testRemoteActor(t, "https://wizard.casa/@benpate")
	require.Equal(t, "https://wizard.casa/users/benpate", result.ID())
	require.Equal(t, "Person", result.Type())
	require.Equal(t, "benpate", result.PreferredUsername())
	require.Equal(t, "https://wizard.casa/users/benpate/inbox", result.Inbox().String())
}

func TestRemoteActor_Mitra_WebFinger(t *testing.T) {
	result := testRemoteActor(t, "@benpate@wizard.casa")
	require.Equal(t, "benpate", result.PreferredUsername())
}

func testRemoteActor(t *testing.T, url string) streams.Document {
	client := Client{}
	result, err := client.Load(url, AsActor())
	require.NoError(t, err)

	require.True(t, result.Outbox().NotNil())
	outbox := result.Outbox().LoadLink()
	require.Equal(t, vocab.CoreTypeOrderedCollection, outbox.Type())
	// require.Greater(t, outbox.TotalItems(), 0)

	return result
}
