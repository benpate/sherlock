//go:build localonly

package sherlock

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

func TestMicroformats(t *testing.T) {

	urlString := "http://localhost/63810bae721f7a33807f25c8"

	var body bytes.Buffer
	uri, _ := url.Parse(urlString)

	if err := remote.Get(urlString).Response(&body, nil).Send(); err != nil {
		t.Error(err)
	}

	result := mapof.NewAny()

	ParseMicroFormats(uri, &body, result)
	t.Log(result)
}
