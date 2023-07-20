package sherlock

import (
	"net/http"
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/stretchr/testify/require"
)

func TestActor_RSSFeed(t *testing.T) {
	testValidateFeed(t, "https://www.smashingmagazine.com/feed")
}

func TestActor_RSSLink(t *testing.T) {
	testValidateFeed(t, "https://www.smashingmagazine.com")
}

func TestActor_Atom(t *testing.T) {
	body, err := testFile("atom-1.xml")
	require.Nil(t, err)

	acc := actorAccumulator{
		httpResponse: &http.Response{
			Header: http.Header{
				"Content-Type": []string{"application/atom+xml"},
			},
		},
		body: body,
	}

	client := Client{}
	client.actor_AtomFeed(&acc)

	require.Equal(t, acc.result["name"], "FYI Center for Software Developers")
	require.Equal(t, acc.result["type"], "Service")

	outbox := acc.result["outbox"].(streams.OrderedCollection)
	require.Equal(t, outbox.TotalItems, 3)
	require.Equal(t, len(outbox.OrderedItems), 3)
}

func testValidateFeed(t *testing.T, url string) {
	client := Client{}
	result, err := client.LoadActor(url)
	require.Nil(t, err)

	value := result.Value().(map[string]any)
	collection := value["outbox"].(streams.OrderedCollection)
	require.Greater(t, collection.TotalItems, 0)
	require.Equal(t, collection.TotalItems, len(collection.OrderedItems))
	// spew.Dump(collection)
}
