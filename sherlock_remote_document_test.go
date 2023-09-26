//go:build localonly

package sherlock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoteDocument_IndieWeb(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("https://indieweb.org")

	require.Nil(t, err)
	require.NotNil(t, meta.Value())
	// t.Log(meta.Value())
}

func TestRemoteDocument_NickSimpson(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("https://nicksimson.com/posts/2022-never/")

	require.Nil(t, err)
	require.NotNil(t, meta.Value())
	// t.Log(meta.Value())
}

func TestRemoteDocument_MastodonProfile(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("https://mastodon.social/@benpate")

	require.Nil(t, err)
	require.NotNil(t, meta.Value())
	// t.Log(meta.Value())
}

func TestRemoteDocument_MastodonProfile_Redirect(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("https://mastodon.social/users/benpate")

	require.Nil(t, err)
	require.NotNil(t, meta.Value())
	// t.Log(meta.Value())
}

func TestRemoteDocument_MastodonToot(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("https://mastodon.social/@benpate/109596019301374311")

	require.Nil(t, err)
	require.NotNil(t, meta.Value())
	// t.Log(meta.Value())
}
