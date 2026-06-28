package webfinger

import (
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/stretchr/testify/require"
)

// fakeClient is a streams.Client test double that records the id it was asked to
// load/delete, so the webfinger routing logic can be exercised without a network.
type fakeClient struct {
	loadedIDs  []string
	deletedIDs []string
	saveCalled bool
	rootSet    bool
}

func (c *fakeClient) Load(id string, options ...any) (streams.Document, error) {
	c.loadedIDs = append(c.loadedIDs, id)
	return streams.NilDocument(), nil
}

func (c *fakeClient) Delete(id string) error {
	c.deletedIDs = append(c.deletedIDs, id)
	return nil
}

func (c *fakeClient) Save(document streams.Document) error {
	c.saveCalled = true
	return nil
}

func (c *fakeClient) SetRootClient(rootClient streams.Client) {
	c.rootSet = true
}

// TestIsWebfinger_NonWebfingerInputs covers every branch that decides "not a
// WebFinger handle" WITHOUT any network access: URLs, values with no "@", and
// values that digit cannot parse into an account.
func TestIsWebfinger_NonWebfingerInputs(t *testing.T) {
	client := New(&fakeClient{}).(Client)

	do := func(input string) {
		isWebfinger, out, err := client.isWebfinger(input)
		require.NoError(t, err, "input=%s", input)
		require.False(t, isWebfinger, "input=%s must not be treated as WebFinger", input)
		require.Equal(t, input, out, "input=%s must be returned unchanged", input)
	}

	do("https://example.com/users/alice") // https URL → skip WebFinger
	do("http://example.com/users/alice")  // http URL → skip WebFinger
	do("plain-string-no-at-sign")         // no "@" → not an account
	do("")                                // empty → not an account
	do("@")                               // "@" only → digit cannot parse an account
}

func TestLoad_NonWebfinger_ForwardsUnchanged(t *testing.T) {
	// A plain URL is not a WebFinger handle, so it is forwarded to the inner
	// client untouched, with no network access.
	inner := &fakeClient{}
	client := New(inner)

	_, err := client.Load("https://example.com/users/alice")
	require.NoError(t, err)
	require.Equal(t, []string{"https://example.com/users/alice"}, inner.loadedIDs)
}

func TestDelete_NonWebfinger_ForwardsUnchanged(t *testing.T) {
	inner := &fakeClient{}
	client := New(inner)

	require.NoError(t, client.Delete("https://example.com/users/alice"))
	require.Equal(t, []string{"https://example.com/users/alice"}, inner.deletedIDs)
}

func TestSaveAndSetRootClient_ForwardToInner(t *testing.T) {
	inner := &fakeClient{}
	client := New(inner)

	require.NoError(t, client.Save(streams.NilDocument()))
	require.True(t, inner.saveCalled)

	client.SetRootClient(&fakeClient{})
	require.True(t, inner.rootSet)
}
