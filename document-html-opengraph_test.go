//go:build localonly

package sherlock

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestAudio_RSS(t *testing.T) {

	client := NewClient()
	doc, err := client.Load("https://music.control.org", withTestServer(), AsActor())
	require.Nil(t, err)
	spew.Dump(doc.Value())
}

func TestOpenGraph(t *testing.T) {
	client := Client{}
	doc, err := client.Load("https://music.control.org/tetraphobia/", withTestServer(), AsActor())
	require.Nil(t, err)
	spew.Dump(doc.Value())
}
