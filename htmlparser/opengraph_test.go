package adapter

import (
	"testing"

	"github.com/benpate/remote"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestOpenGraph(t *testing.T) {

	url := "https://music.control.org/tetraphobia/"
	body := []byte{}
	data := map[string]any{}

	txn := remote.Get(url).Result(&body)

	if err := txn.Send(); err != nil {
		t.Fatal(err)
	}

	err := OpenGraph(url, body, data)
	require.NoError(t, err)

	spew.Dump(url, data)
}
