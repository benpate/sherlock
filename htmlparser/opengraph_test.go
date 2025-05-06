package adapter

import (
	"testing"

	"github.com/benpate/remote"
	"github.com/davecgh/go-spew/spew"
)

func TestOpenGraph(t *testing.T) {

	url := "https://music.control.org/tetraphobia/"
	body := []byte{}
	data := map[string]any{}

	txn := remote.Get(url).Result(&body)

	if err := txn.Send(); err != nil {
		t.Fatal(err)
	}

	OpenGraph(url, body, data)

	spew.Dump(url, data)
}
