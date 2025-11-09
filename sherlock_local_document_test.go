//gobuild:localonly

package sherlock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalDocument_AP_Mastodon_JSON(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/document-ap-mastodon.json", withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
}

/* DISABLING THIS TEST FOR NOW. IT'S GETTING HUNG UP ON LOADING THE ACTIVITYSTREAM DOCUMENT FIRST...

func TestLocalDocument_AP_Mastodon_HTML(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/document-ap-mastodon.html", withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
}

func TestLocalDocument_Microformats_1(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/document-microformats-1.html", withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
}

func TestLocalDocument_Microformats_2(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/document-microformats-2.html", withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
}

func TestLocalDocument_Microformats_3(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/document-microformats-3.html", withTestServer())
	require.Nil(t, err)
	require.NotNil(t, result.Value())
}

func TestLocalActor_Microformats_2(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/document-microformats-4.html", withTestServer())

	require.Nil(t, err)
	require.NotNil(t, result.Value())
	require.Equal(t, "Page", result.Type())
	require.Equal(t, "https://test-server/document-microformats-4.html", result.ID())
	require.Equal(t, "https://test-server/IndieWeb", result.AttributedTo().ID())
	require.Equal(t, "IndieWeb", result.AttributedTo().Name())
}

func TestLocalDocument_OpenGraph(t *testing.T) {

	client := NewClient()

	result, err := client.Load("https://test-server/document-opengraph-1.html", withTestServer())

	require.Nil(t, err)
	require.NotNil(t, result.Value())
	require.Equal(t, "https://test-server/document-opengraph-1.html", result.ID())
	require.Equal(t, "Page", result.Type())
	require.Equal(t, "Open Graph Tester", result.Name())
	require.Equal(t, "This website serves as a simple tool for web developers, designers, and marketing professionals to optimize their websites and posts prior to publishing them on social media.", result.Summary())
	require.Equal(t, "https://opengraphtester.com/assets/images/logos/og-image.png", result.Image().URL())
}
*/
