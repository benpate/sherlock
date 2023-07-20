/// go:build local

package sherlock

import (
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestLoad_Local(t *testing.T) {

	client := NewClient()
	meta, err := client.LoadDocument("http://localhost/63810bae721f7a33807f25c8", mapof.NewAny())

	require.Nil(t, err)
	t.Log(meta.Value())
}

func TestLoad_IndieWeb(t *testing.T) {

	client := NewClient()
	meta, err := client.LoadDocument("https://indieweb.org", mapof.NewAny())

	require.Nil(t, err)
	spew.Dump(meta.Value())
}

func TestLoad_NickSimpson(t *testing.T) {

	client := NewClient()
	meta, err := client.LoadDocument("https://nicksimson.com/posts/2022-never/", mapof.NewAny())

	require.Nil(t, err)
	spew.Dump(meta.Value())
}

func TestLoad_OpenGraph(t *testing.T) {

	client := NewClient()
	meta, err := client.LoadDocument("https://opengraphtester.com", mapof.NewAny())

	require.Nil(t, err)
	spew.Dump(meta.Value())
}

func TestLoad_MastodonProfile(t *testing.T) {

	client := NewClient()
	meta, err := client.LoadDocument("https://mastodon.social/@benpate", mapof.NewAny())

	require.Nil(t, err)
	spew.Dump(meta.Value())
}

func TestLoad_MastodonProfile_Redirect(t *testing.T) {

	client := NewClient()
	meta, err := client.LoadDocument("https://mastodon.social/users/benpate", mapof.NewAny())

	require.Nil(t, err)
	spew.Dump(meta.Value())
}

func TestLoad_MastodonToot(t *testing.T) {

	client := NewClient()
	meta, err := client.LoadDocument("https://mastodon.social/@benpate/109596019301374311", mapof.NewAny())

	require.Nil(t, err)
	spew.Dump(meta.Value())
}
