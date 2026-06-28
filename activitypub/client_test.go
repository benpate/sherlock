package activitypub

import (
	"crypto"
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/stretchr/testify/require"
)

// fakeClient is a streams.Client test double that records the calls made to it
// and returns scripted values, so the activitypub Client's routing logic can be
// exercised without any network access.
type fakeClient struct {
	loadedID    string
	loadResult  streams.Document
	loadErr     error
	deletedID   string
	savedCalled bool
	rootClient  streams.Client
	rootWasSet  bool
}

func (c *fakeClient) Load(id string, options ...any) (streams.Document, error) {
	c.loadedID = id
	return c.loadResult, c.loadErr
}

func (c *fakeClient) Save(document streams.Document) error {
	c.savedCalled = true
	return nil
}

func (c *fakeClient) Delete(id string) error {
	c.deletedID = id
	return nil
}

func (c *fakeClient) SetRootClient(rootClient streams.Client) {
	c.rootClient = rootClient
	c.rootWasSet = true
}

func TestNew_Defaults(t *testing.T) {
	// A Client built with no options still has the default User-Agent and no
	// inner client.
	client := New().(*Client)

	require.Equal(t, "Sherlock (https://github.com/benpate/sherlock)", client.userAgent)
	require.Nil(t, client.innerClient)
	require.Nil(t, client.keyPairFunc)
}

func TestNew_WithOptions(t *testing.T) {
	inner := &fakeClient{}

	client := New(
		WithInnerClient(inner),
		WithUserAgent("custom-agent/1.0"),
		WithKeyPairFunc(func() (string, crypto.PrivateKey) { return "key-id", nil }),
	).(*Client)

	require.Equal(t, "custom-agent/1.0", client.userAgent)
	require.Same(t, inner, client.innerClient)
	require.NotNil(t, client.keyPairFunc)
}

func TestLoad_InvalidURL_RoutesToInnerClient(t *testing.T) {
	// An id that is not a valid URL must be forwarded to the inner client
	// untouched, without attempting any network access.
	inner := &fakeClient{loadResult: streams.NilDocument()}
	client := New(WithInnerClient(inner))

	_, err := client.Load("this is not a url")
	require.NoError(t, err)
	require.Equal(t, "this is not a url", inner.loadedID, "invalid URL must be forwarded to the inner client")
}

func TestLoad_InvalidURL_NoInnerClient_ReturnsNotFound(t *testing.T) {
	// With no inner client to fall back on, an invalid URL is a NotFound.
	client := New()

	doc, err := client.Load("@not-a-url")
	require.Error(t, err)
	require.True(t, doc.IsNil(), "a failed load must return a nil document")
}

func TestDelete_NilInnerClient_IsNoOp(t *testing.T) {
	// Delete is a no-op (no error) when there is no inner cache to delete from.
	client := New()
	require.NoError(t, client.Delete("https://example.com/123"))
}

func TestDelete_ForwardsToInnerClient(t *testing.T) {
	inner := &fakeClient{}
	client := New(WithInnerClient(inner))

	require.NoError(t, client.Delete("https://example.com/123"))
	require.Equal(t, "https://example.com/123", inner.deletedID)
}

func TestSave_NilInnerClient_IsNoOp(t *testing.T) {
	client := New()
	require.NoError(t, client.Save(streams.NilDocument()))
}

func TestSave_ForwardsToInnerClient(t *testing.T) {
	inner := &fakeClient{}
	client := New(WithInnerClient(inner))

	require.NoError(t, client.Save(streams.NilDocument()))
	require.True(t, inner.savedCalled)
}

func TestSetRootClient_NilInnerClient_IsNoOp(t *testing.T) {
	// With no inner client, SetRootClient must not panic and (per the current
	// implementation) does not record the root client.
	client := New().(*Client)
	root := &fakeClient{}

	require.NotPanics(t, func() { client.SetRootClient(root) })
	require.Nil(t, client.rootClient, "rootClient is only set when an inner client exists")
}

func TestSetRootClient_ForwardsToInnerClient(t *testing.T) {
	inner := &fakeClient{}
	client := New(WithInnerClient(inner)).(*Client)
	root := &fakeClient{}

	client.SetRootClient(root)
	require.True(t, inner.rootWasSet, "the root client must propagate to the inner client")
	require.Same(t, root, client.rootClient)
}
