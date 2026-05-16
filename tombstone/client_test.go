package tombstone

import (
	"testing"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

func TestClientSuccess(t *testing.T) {

	client := New(newFakeInnerClient(
		streams.NewDocument(
			mapof.Any{
				"id":      "http://example.com/object/123",
				"type":    "Note",
				"content": "hey howdy",
			},
		),
		nil,
	))

	result, err := client.Load("http://example.com/object/123")

	require.Nil(t, err)
	require.Equal(t, vocab.ObjectTypeNote, result.Type())
}

func TestClientGone(t *testing.T) {

	client := New(newFakeInnerClient(
		streams.NilDocument(),
		derp.Gone("here", "there"),
	))

	result, err := client.Load("http://example.com/object/123")

	require.Nil(t, err)
	require.Equal(t, vocab.ObjectTypeTombstone, result.Type())
}

func TestClientError(t *testing.T) {

	client := New(newFakeInnerClient(
		streams.NilDocument(),
		derp.Internal("here", "there"),
	))

	result, err := client.Load("http://example.com/object/123")

	require.NotNil(t, err)
	require.Equal(t, vocab.Unknown, result.Type())
}

// fakeInnerClient just returns whatever is configured in it
type fakeInnerClient struct {
	result streams.Document
	err    error
}

func newFakeInnerClient(result streams.Document, err error) fakeInnerClient {
	return fakeInnerClient{
		result: result,
		err:    err,
	}
}

func (client fakeInnerClient) SetRootClient(rootClient streams.Client) {}

func (client fakeInnerClient) Load(uri string, options ...any) (streams.Document, error) {
	return client.result, client.err
}

func (client fakeInnerClient) Save(document streams.Document) error {
	return nil
}

func (client fakeInnerClient) Delete(uri string) error {
	return nil
}
