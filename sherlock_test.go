/// go:build local

package sherlock

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestLoad_Local(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("http://localhost/63810bae721f7a33807f25c8")

	require.Nil(t, err)
	t.Log(meta.Value())
}

func TestLoad_IndieWeb(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("https://indieweb.org")

	require.Nil(t, err)
	spew.Dump(meta.Value())
}

func TestLoad_OpenGraph(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("https://opengraphtester.com")

	require.Nil(t, err)
	spew.Dump(meta.Value())
}

func TestLoad_MastodonProfile(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("https://mastodon.social/@benpate")

	require.Nil(t, err)
	spew.Dump(meta.Value())
}

func TestLoad_MastodonProfile_Redirect(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("https://mastodon.social/users/benpate")

	require.Nil(t, err)
	spew.Dump(meta.Value())
}

func TestLoad_MastodonToot(t *testing.T) {

	client := NewClient()
	meta, err := client.Load("https://mastodon.social/@benpate/109596019301374311")

	require.Nil(t, err)
	spew.Dump(meta.Value())
}
