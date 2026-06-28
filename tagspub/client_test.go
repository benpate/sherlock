package tagspub

import (
	"errors"
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/stretchr/testify/require"
)

// fakeClient is a streams.Client test double that records the id it was asked to
// load/delete, so the tagspub routing logic can be exercised without a network.
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
	require.Equal(t, "tags.pub", client.hostname)
}

func TestIsHashtag_RewritesToWebFingerHandle(t *testing.T) {
	client := New(&fakeClient{}).(Client)

	// A hashtag is rewritten into a "tag@hostname" WebFinger handle.
	match, handle := client.isHashtag("#golang")
	require.True(t, match)
	require.Equal(t, "golang@tags.pub", handle)

	// A non-hashtag is returned unchanged.
	match, handle = client.isHashtag("not-a-hashtag")
	require.False(t, match)
	require.Equal(t, "not-a-hashtag", handle)
}

func TestLoad_Hashtag_TriesWebFingerHandleFirst(t *testing.T) {
	inner := &fakeClient{}
	client := New(inner)

	_, err := client.Load("#golang")
	require.NoError(t, err)
	require.Equal(t, []string{"golang@tags.pub"}, inner.loadedIDs,
		"a hashtag must be loaded via the rewritten WebFinger handle")
}

func TestLoad_Hashtag_FallsBackToOriginalOnError(t *testing.T) {
	inner := &fakeClient{loadErr: errors.New("miss")}
	client := New(inner)

	_, _ = client.Load("#golang")
	require.Equal(t, []string{"golang@tags.pub", "#golang"}, inner.loadedIDs,
		"a failed hashtag load must fall back to the original id")
}

func TestLoad_NonHashtag_ForwardsUnchanged(t *testing.T) {
	inner := &fakeClient{}
	client := New(inner)

	_, err := client.Load("https://example.com/users/alice")
	require.NoError(t, err)
	require.Equal(t, []string{"https://example.com/users/alice"}, inner.loadedIDs)
}

func TestDelete_Hashtag_UsesWebFingerHandle(t *testing.T) {
	inner := &fakeClient{}
	client := New(inner)

	require.NoError(t, client.Delete("#golang"))
	require.Equal(t, []string{"golang@tags.pub"}, inner.deletedIDs)
}

func TestSaveAndSetRootClient_ForwardToInner(t *testing.T) {
	inner := &fakeClient{}
	client := New(inner)

	require.NoError(t, client.Save(streams.NilDocument()))
	require.True(t, inner.saveCalled)

	client.SetRootClient(&fakeClient{})
	require.True(t, inner.rootSet)
}
