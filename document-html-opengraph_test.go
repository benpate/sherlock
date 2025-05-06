package sherlock

import (
	"testing"

	"github.com/benpate/remote"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestAudio_RSS(t *testing.T) {

	var url string = "https://music.control.org"
	var body []byte

	txn := remote.Get(url).Result(&body)

	if err := txn.Send(); err != nil {
		t.Fatal(err)
	}

	client := Client{}
	loadConfig := NewLoadConfig()

	doc, err := client.loadActor(url, &loadConfig)
	require.Nil(t, err)
	spew.Dump(doc.Value())
}

func TestOpenGraph(t *testing.T) {

	url := "https://music.control.org/tetraphobia/"
	body := []byte{}
	data := map[string]any{}

	txn := remote.Get(url).Result(&body)

	if err := txn.Send(); err != nil {
		t.Fatal(err)
	}

	client := Client{}

	client.loadDocument_OpenGraph(url, body, data)

	spew.Dump(url, data)
}
