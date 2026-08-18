package metadata

import (
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

func TestMerge_GroupsAreCopiedNotAliased(t *testing.T) {

	var m Preview

	source := partial{Provider: null.NewObject(Provider{Name: "Example"})}
	merge(&m, source)

	// The Preview owns its copy of the winning group. The icon backfill in
	// extract() writes here, and must never reach back into the partial.
	m.Provider.IconURL = "https://example.com/favicon.ico"

	require.Empty(t, source.Provider.Object().IconURL)
}

func TestMerge_EmbedFloorRejectsFileAsPlayer(t *testing.T) {

	// An iframe-mode embed pointing at a media FILE fails the floor.
	file := null.NewObject(Embed{Mode: EmbedIframe, IframeURL: "https://example.com/video.mp4"})
	player := null.NewObject(Embed{Mode: EmbedIframe, IframeURL: "https://example.com/embed/123"})

	require.False(t, validEmbed(file))
	require.True(t, validEmbed(player))

	// The same URL is fine as an explicit stream.
	stream := null.NewObject(Embed{Mode: EmbedStream, StreamURL: "https://example.com/video.mp4"})
	require.True(t, validEmbed(stream))

	// A null embed fails the floor, and reads back as a zero Embed.
	require.False(t, validEmbed(null.Object[Embed]{}))
}

func TestFirstString_URLExceptionOrder(t *testing.T) {

	// §4.2 exception 1: AS id → rel=canonical → og:url.
	asID := null.NewString("https://social.example/objects/1")
	canonical := null.NewString("https://example.com/article")
	ogURL := null.NewString("https://example.com/article?utm_campaign=x")

	var missing null.String

	require.Equal(t, asID.String(), firstString(asID, canonical, ogURL))
	require.Equal(t, canonical.String(), firstString(missing, canonical, ogURL))
	require.Equal(t, ogURL.String(), firstString(missing, missing, ogURL))

	// A present-but-empty string is not a value either.
	require.Equal(t, "", firstString(missing, null.NewString(""), missing))
}

func TestMerge_BlankCardIsValid(t *testing.T) {

	// All sources empty → blank-but-valid preview (§4.2).
	var m Preview
	merge(&m, partial{})
	merge(&m, partial{})

	m.RequestURL = "https://example.com/empty"
	m.URL = "https://example.com/empty"
	if m.Kind == KindUnknown {
		m.Kind = KindWebsite
	}

	require.Equal(t, KindWebsite, m.Kind)
	require.Empty(t, m.Title)
	require.Nil(t, m.Thumbnail)
}
