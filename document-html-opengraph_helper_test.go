package sherlock

import (
	"testing"

	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

func TestMapOpenGraphTags(t *testing.T) {

	// Empty input yields an empty (non-nil) slice
	result := mapOpenGraphTags([]string{})
	require.NotNil(t, result)
	require.Empty(t, result)

	// Each tag becomes a map keyed by name, preserving order
	result = mapOpenGraphTags([]string{"alpha", "beta"})
	require.Len(t, result, 2)
	require.Equal(t, "alpha", result[0][vocab.PropertyName])
	require.Equal(t, "beta", result[1][vocab.PropertyName])
}

func TestLoadDocumentOpenGraph(t *testing.T) {

	client := NewClient()

	html := `<html><head>
		<meta property="og:title" content="My Title" />
		<meta property="og:description" content="My Description" />
		<meta property="og:image" content="https://example.com/image.png" />
		<meta property="og:url" content="https://example.com/canonical" />
	</head><body></body></html>`

	data := mapof.NewAny()
	client.loadDocument_OpenGraph("https://example.com/page", []byte(html), data)

	require.Equal(t, "My Title", data[vocab.PropertyName])
	require.Equal(t, "My Description", data[vocab.PropertySummary])
	require.Equal(t, "https://example.com/image.png", data[vocab.PropertyImage])
	require.Equal(t, "https://example.com/canonical", data[vocab.PropertyID])
	require.Equal(t, vocab.ObjectTypeArticle, data[vocab.PropertyType])
}

func TestLoadDocumentOpenGraph_PreservesExisting(t *testing.T) {

	client := NewClient()

	html := `<html><head><meta property="og:title" content="OpenGraph Title" /></head></html>`

	// A name that's already set is NOT overwritten by OpenGraph
	data := mapof.Any{vocab.PropertyName: "Existing Name"}
	client.loadDocument_OpenGraph("https://example.com/page", []byte(html), data)

	require.Equal(t, "Existing Name", data[vocab.PropertyName])
}

func TestLoadDocumentOpenGraph_NoData(t *testing.T) {

	client := NewClient()

	// Even with no OpenGraph tags, the loader unconditionally assigns empty Title and
	// Description values (their IsZeroValue guard passes), so the map becomes non-empty
	// and the ID/Type defaults are then applied from the supplied URL.
	data := mapof.NewAny()
	client.loadDocument_OpenGraph("https://example.com/page", []byte(`<html></html>`), data)

	require.Equal(t, "", data[vocab.PropertyName])
	require.Equal(t, "", data[vocab.PropertySummary])
	require.Equal(t, "https://example.com/page", data[vocab.PropertyID])
	require.Equal(t, vocab.ObjectTypeArticle, data[vocab.PropertyType])
}
