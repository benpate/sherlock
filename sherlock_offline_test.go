package sherlock

import (
	"os"
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote/options"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// withOfflineServer returns an Option that serves the ./test-files directory for the
// "test-server.local" hostname, so these tests run fully offline. Unlike the localonly
// withTestServer helper, this one is always compiled so it runs in CI.
func withOfflineServer() Option {
	return func(config *Config) {
		filesystem := os.DirFS("./test-files")
		config.RemoteOptions = append(
			config.RemoteOptions,
			options.TestServer("test-server.local", filesystem),
		)
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}
}

func TestOffline_Document_ActivityPub(t *testing.T) {

	client := NewClient()

	// loadDocument should detect the ActivityPub content type and return the actor as-is.
	doc, err := client.Load("https://test-server.local/offline-actor.json", withOfflineServer())
	require.Nil(t, err)
	require.False(t, doc.IsNil())
	require.Equal(t, "https://test-server.local/actor", doc.ID())
	require.Equal(t, "Person", doc.Type())
	require.Equal(t, "Test Person", doc.Name())
}

func TestOffline_Actor_ActivityPub(t *testing.T) {

	client := NewClient()

	// loadActor should follow the ActivityStreams branch for an AP document.
	doc, err := client.Load("https://test-server.local/offline-actor.json", AsActor(), withOfflineServer())
	require.Nil(t, err)
	require.False(t, doc.IsNil())
	require.Equal(t, "https://test-server.local/actor", doc.ID())
}

func TestOffline_Actor_JSONFeed(t *testing.T) {

	client := NewClient()

	// loadActor should fall through to the JSON-feed parser and synthesize an Actor.
	doc, err := client.Load("https://test-server.local/offline-feed.json", AsActor(), withOfflineServer())
	require.Nil(t, err)
	require.False(t, doc.IsNil())

	require.Equal(t, "https://test-server.local/offline-feed.json", doc.ID())
	require.Equal(t, vocab.ActorTypeApplication, doc.Type())
	require.Equal(t, "Offline Feed", doc.Name())

	// The two items should appear in the outbox as an OrderedCollection.
	require.Equal(t, vocab.CoreTypeOrderedCollection, doc.Outbox().Type())
	require.Equal(t, 2, doc.Outbox().TotalItems())
	require.Equal(t, 2, doc.Outbox().Items().Len())
}

// withIcon pre-seeds the document with an icon. The feed parsers call
// loadActor_Feed_FindHomePageIcon, which builds its OWN remote request (without our
// TestServer option) and would otherwise reach out to the real network for the bare
// host. Pre-seeding the icon makes that helper NOOP, keeping these tests fast/offline.
func withIcon() Option {
	return WithDefaultValue(map[string]any{
		vocab.PropertyIcon: "https://test-server.local/icon.png",
	})
}

func TestOffline_Actor_RSS(t *testing.T) {

	client := NewClient()

	// loadActor should fall through to the RSS parser and synthesize an Actor.
	doc, err := client.Load("https://test-server.local/offline-rss.xml", AsActor(), withOfflineServer(), withIcon())
	require.Nil(t, err)
	require.False(t, doc.IsNil())

	require.Equal(t, vocab.ActorTypeApplication, doc.Type())
	require.Equal(t, "Offline RSS", doc.Name())
	require.Equal(t, 2, doc.Outbox().TotalItems())
}

// NOTE: There is intentionally no offline test for the h-feed MicroFormats actor path.
// loadActor_Feed_MicroFormats builds its result map directly (not from
// config.DefaultValue) and then calls loadActor_Feed_FindHomePageIcon, which issues an
// un-interceptable real network request for the site's bare host. That makes the path
// impossible to exercise quickly offline. The parsing logic itself is covered by the
// microformat_* unit tests, and the end-to-end path by the localonly suite.

func TestOffline_Actor_InvalidIdentifier(t *testing.T) {

	client := NewClient()

	// An identifier that is neither a URL nor a username is rejected up front.
	doc, err := client.Load("not a valid identifier !!", AsActor(), withOfflineServer())
	require.NotNil(t, err)
	require.True(t, doc.IsNil())
}

func TestOffline_Document_Unreachable(t *testing.T) {

	client := NewClient()

	// A file that does not exist on the test server cannot be loaded as a document.
	// (The test server returns a 404, so neither ActivityStream nor HTML parsing succeeds,
	// and there is no DefaultValue to fall back on.)
	doc, err := client.Load("https://test-server.local/does-not-exist.json", withOfflineServer())
	require.NotNil(t, err)
	require.True(t, doc.IsNil())
}

func TestOffline_Document_DefaultValueFallback(t *testing.T) {

	client := NewClient()

	// NOTE: loadDocument's DefaultValue fallback (step 3) is only reachable when the
	// ActivityStream load returns a *nil document with no error*. When the request
	// itself fails (e.g. a 404 from the test server), loadDocument_ActivityStream
	// returns a wrapped error and loadDocument returns immediately -- the DefaultValue
	// is never consulted. This test documents that short-circuit behavior.
	defaultValue := map[string]any{
		vocab.PropertyID:   "https://test-server.local/fallback",
		vocab.PropertyType: vocab.ObjectTypeNote,
	}

	doc, err := client.Load(
		"https://test-server.local/does-not-exist.json",
		withOfflineServer(),
		WithDefaultValue(defaultValue),
	)

	require.NotNil(t, err)
	require.True(t, doc.IsNil())
}
