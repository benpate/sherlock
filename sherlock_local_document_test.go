package sherlock

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalDocument_AP_Mastodon_JSON(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/document-ap-mastodon.json")
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	// spew.Dump(result.Value())
}

func TestLocalDocument_AP_Mastodon_HTML(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/document-ap-mastodon.html")
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	// spew.Dump(result.Value())
}

func TestLocalDocument_Microformats_1(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/document-microformats-1.html")
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	// spew.Dump(result.Value())
}

func TestLocalDocument_Microformats_2(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/document-microformats-2.html")
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	// spew.Dump(result.Value())
}

func TestLocalDocument_Microformats_3(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/document-microformats-3.html")
	require.Nil(t, err)
	require.NotNil(t, result.Value())
	// spew.Dump(result.Value())
}

func TestLocalActor_Microformats_2(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/document-microformats-4.html")

	require.Nil(t, err)
	require.NotNil(t, result.Value())
	require.Equal(t, "Page", result.Type())
	require.Equal(t, "https://test-server/document-microformats-4.html", result.ID())
	require.Equal(t, "https://test-server/IndieWeb", result.AttributedTo().ID())
	require.Equal(t, "IndieWeb", result.AttributedTo().Name())
	// spew.Dump(result.Value())
}

func TestLocalDocument_OpenGraph(t *testing.T) {

	client := NewClient(WithRemoteOptions(getTestServer()))

	result, err := client.Load("https://test-server/document-opengraph-1.html")

	require.Nil(t, err)
	require.NotNil(t, result.Value())
	require.Equal(t, "https://test-server/document-opengraph-1.html", result.ID())
	require.Equal(t, "Page", result.Type())
	require.Equal(t, "Open Graph Tester", result.Name())
	require.Equal(t, "This website serves as a simple tool for web developers, designers, and marketing professionals to optimize their websites and posts prior to publishing them on social media.", result.Summary())
	require.Equal(t, "https://opengraphtester.com/assets/images/logos/og-image.png", result.Image().URL())
}
