//go:build localonly

package sherlock

import (
	"bytes"
	"net/url"
	"os"
	"testing"

	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestMicroformats(t *testing.T) {

	urlString := "https://nicksimson.com/posts/2022-never/"

	var body bytes.Buffer
	uri, _ := url.Parse(urlString)

	if err := remote.Get(urlString).Result(&body).Send(); err != nil {
		t.Error(err)
	}

	input := mapof.NewAny()

	result := ParseMicroFormats(uri, &body, input)
	spew.Dump(input)
	spew.Dump(result)
}

func TestMicroformats_Files(t *testing.T) {

	testDirectory := "./test-files"
	files, err := os.ReadDir(testDirectory)

	require.Nil(t, err)

	for _, fileEntry := range files {
		fileBytes, err := os.ReadFile(testDirectory + "/" + fileEntry.Name())
		require.Nil(t, err)

		var buffer bytes.Buffer
		buffer.Write(fileBytes)

		result := mapof.NewAny()
		require.Nil(t, Parse("", &buffer, result))

		spew.Dump("------------------------------------------------", fileEntry.Name(), result)
		return
	}
}
