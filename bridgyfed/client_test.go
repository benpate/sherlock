package bridgyfed

import (
	"errors"
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/stretchr/testify/require"
)

// fakeClient is a streams.Client test double that records the id it was asked to
// load/delete and returns scripted values, so the bridgyfed routing logic can be
// exercised without any network access.
type fakeClient struct {
	loadedIDs  []string
	deletedIDs []string
	loadErr    error
	saveCalled bool
	rootSet    bool
}

func (c *fakeClient) Load(id string, options ...any) (streams.Document, error) {
	c.loadedIDs = append(c.loadedIDs, id)
	return streams.NilDocument(), c.loadErr
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

func TestNew_DefaultHostname(t *testing.T) {
	client := New(&fakeClient{}).(Client)
	require.Equal(t, "bsky.brid.gy", client.hostname)
}

func TestLooksLikeBluesky_FormatsHandle(t *testing.T) {
	client := New(&fakeClient{}).(Client)

	// A bluesky-looking handle is rewritten into a Bridgy Fed WebFinger handle.
	ok, handle := client.looksLikeBluesky("alice.bsky.social")
	require.True(t, ok)
	require.Equal(t, "@alice.bsky.social@bsky.brid.gy", handle)

	// A leading "@" is preserved (not doubled).
	ok, handle = client.looksLikeBluesky("@bob.bsky.social")
	require.True(t, ok)
	require.Equal(t, "@bob.bsky.social@bsky.brid.gy", handle)

	// A non-bluesky value is returned unchanged.
	ok, handle = client.looksLikeBluesky("no:colons.here")
	require.False(t, ok)
	require.Equal(t, "no:colons.here", handle)
}

func TestLoad_BlueskyHandle_TriesBridgyHandleFirst(t *testing.T) {
	// A bluesky-looking id is first loaded via the rewritten Bridgy handle.
	inner := &fakeClient{}
	client := New(inner)

	_, err := client.Load("alice.bsky.social")
	require.NoError(t, err)
	require.Equal(t, []string{"@alice.bsky.social@bsky.brid.gy"}, inner.loadedIDs,
		"a bluesky handle must be loaded via the rewritten Bridgy handle")
}

func TestLoad_BlueskyHandle_FallsBackToOriginalOnError(t *testing.T) {
	// When the Bridgy handle load fails, the original id is retried down the stack.
	inner := &fakeClient{loadErr: errors.New("bridgy miss")}
	client := New(inner)

	_, _ = client.Load("alice.bsky.social")
	require.Equal(t, []string{"@alice.bsky.social@bsky.brid.gy", "alice.bsky.social"}, inner.loadedIDs,
		"a failed Bridgy load must fall back to the original id")
}

func TestLoad_NonBluesky_ForwardsUnchanged(t *testing.T) {
	// A value that does not look like Bluesky is forwarded to the inner client as-is.
	inner := &fakeClient{}
	client := New(inner)

	_, err := client.Load("https://example.com/users/alice")
	require.NoError(t, err)
	require.Equal(t, []string{"https://example.com/users/alice"}, inner.loadedIDs)
}

func TestDelete_BlueskyHandle_UsesBridgyHandle(t *testing.T) {
	inner := &fakeClient{}
	client := New(inner)

	require.NoError(t, client.Delete("alice.bsky.social"))
	require.Equal(t, []string{"@alice.bsky.social@bsky.brid.gy"}, inner.deletedIDs)
}

func TestSaveAndSetRootClient_ForwardToInner(t *testing.T) {
	inner := &fakeClient{}
	client := New(inner)

	require.NoError(t, client.Save(streams.NilDocument()))
	require.True(t, inner.saveCalled)

	client.SetRootClient(&fakeClient{})
	require.True(t, inner.rootSet)
}
