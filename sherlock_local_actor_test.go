package sherlock

import (
	"os"
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote/options"
	"github.com/stretchr/testify/require"
)

// spew.Config.DisableMethods = true
// zerolog.SetGlobalLevel(zerolog.InfoLevel)

func withTestServer() Option {
	return func(config *Config) {
		filesystem := os.DirFS("./test-files")
		config.RemoteOptions = append(config.RemoteOptions, options.TestServer("test-server", filesystem))
	}
}

func TestLocalActor_Atom_1(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/actor-atom-1.xml", AsActor(), withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())

	require.True(t, result.IsActor())
	require.Equal(t, vocab.ActorTypeApplication, result.Type())
	require.Equal(t, vocab.CoreTypeOrderedCollection, result.Outbox().Type())
	require.Equal(t, 2, result.Outbox().TotalItems())
	require.Equal(t, 2, result.Outbox().Items().Len())
}

func TestLocalActor_JSON_1(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/actor-json-1.json", AsActor(), withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	require.Equal(t, "https://www.jsonfeed.org/feed.json", result.ID())
	require.Equal(t, "Application", result.Type())
	require.Equal(t, 2, result.Outbox().TotalItems())
}

func TestLocalActor_Microformats_1(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/actor-microformats-1.html", AsActor(), withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
}

func TestLocalActor_RSS_1_XML(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/actor-rss-1.xml", AsActor(), withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
}

func TestLocalActor_RSS_1_HTML(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/actor-rss-1.html", AsActor(), withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
}

func TestLocalActor_RSS_2_XML(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/actor-rss-2.xml", AsActor(), withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
}
