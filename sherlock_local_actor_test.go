package sherlock

import (
	"os"
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/remote"
	"github.com/benpate/remote/options"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func init() {
	// spew.Config.DisableMethods = true
}

func getTestServer() remote.Option {
	filesystem := os.DirFS("./test-files")
	return options.TestServer("test-server", filesystem)
}

func TestTestServer(t *testing.T) {

	txn := remote.Get("https://test-server/actor-microformats-3.html").Use(getTestServer())
	err := txn.Send()
	require.Nil(t, err)

	body, err := txn.ResponseBody()
	require.Nil(t, err)
	require.NotZero(t, len(body))
	// spew.Dump(string(body))
}

func TestLocalActor_Atom_1(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/actor-atom-1.xml", AsActor())
	require.Nil(t, err)
	require.NotNil(t, result.Value())

	require.True(t, result.IsActor())
	require.Equal(t, vocab.ActorTypeApplication, result.Type())
	require.Equal(t, vocab.CoreTypeOrderedCollection, result.Outbox().Type())
	require.Equal(t, 2, result.Outbox().TotalItems())
	require.Equal(t, 2, result.Outbox().Items().Len())
	spew.Dump(result.Value())
}

func TestLocalActor_JSON_1(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/actor-json-1.json", AsActor())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	require.Equal(t, "https://www.jsonfeed.org/feed.json", result.ID())
	require.Equal(t, "Application", result.Type())
	require.Equal(t, 2, result.Outbox().TotalItems())
	// spew.Dump(result.Value())
}

func TestLocalActor_Microformats_1(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/actor-microformats-1.html", AsActor())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	// spew.Dump(result.Value())
}

func TestLocalActor_Microformats_3(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/actor-microformats-3.html", AsActor())
	// require.Nil(t, err)
	// require.NotNil(t, result.Value())
	spew.Dump(result.Value())
	spew.Dump(err)
	// TODO: This test is currently breaking because this page nests MicroFormats too deeply.
	// TODO: Also, do a better job loading Author information from the h-card
}

func TestLocalActor_RSS_1_XML(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/actor-rss-1.xml", AsActor())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	// spew.Dump(result.Value())
}

func TestLocalActor_RSS_1_HTML(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/actor-rss-1.html", AsActor())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	// spew.Dump(result.Value())
}

func TestLocalActor_RSS_2_XML(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/actor-rss-2.xml", AsActor())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	// spew.Dump(result.Value())
}
