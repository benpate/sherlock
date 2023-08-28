package sherlock

import (
	"net/http"
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
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
	client.actor_RSSFeed(&acc)

	require.Equal(t, acc.result["name"], "FYI Center for Software Developers")
	require.Equal(t, acc.result["type"], "Service")

	outbox := acc.result["outbox"].(mapof.Any)
	require.Equal(t, outbox[vocab.PropertyTotalItems], 3)
	require.Equal(t, len(outbox[vocab.PropertyOrderedItems].([]any)), 3)
}

func testValidateFeed(t *testing.T, url string) {
	client := Client{}
	result, err := client.LoadActor(url)
	require.Nil(t, err)

	value := result.Value().(map[string]any)
	collection := value[vocab.PropertyOutbox].(mapof.Any)
	require.Greater(t, collection[vocab.PropertyTotalItems], 0)
	require.Equal(t, collection[vocab.PropertyTotalItems], len(collection[vocab.PropertyOrderedItems].([]any)))
}
