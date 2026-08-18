package metadata

import (
	"testing"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

// D5: correlated fields fill as WHOLE groups from ONE source. A group that
// wins brings only its own fields — never another source's — and a group that
// loses contributes nothing at all, not even fields the winner left empty.

func TestD5_ThumbnailWinsWholeAndIsNeverCompleted(t *testing.T) {

	var preview Preview

	// A bare-URL thumbnail wins; the complete one that follows must lose ENTIRELY.
	merge(&preview, partial{Thumbnail: null.NewObject(Thumbnail{
		URL: "https://first.example.com/a.jpg",
	})})

	merge(&preview, partial{Thumbnail: null.NewObject(Thumbnail{
		URL:    "https://second.example.com/b.jpg",
		Width:  1200,
		Height: 630,
		Alt:    "second alt",
	})})

	require.Equal(t, "https://first.example.com/a.jpg", preview.Thumbnail.URL)
	require.Zero(t, preview.Thumbnail.Width)  // NOT 1200
	require.Zero(t, preview.Thumbnail.Height) // NOT 630
	require.Empty(t, preview.Thumbnail.Alt)   // NOT "second alt"
}

func TestD5_EmbedWinsWholeAndIsNeverCompleted(t *testing.T) {

	var preview Preview

	merge(&preview, partial{Embed: null.NewObject(Embed{
		Mode:      EmbedIframe,
		IframeURL: "https://first.example.com/player",
	})})

	merge(&preview, partial{Embed: null.NewObject(Embed{
		Mode:      EmbedIframe,
		IframeURL: "https://second.example.com/player",
		Width:     640,
		Height:    360,
	})})

	require.Equal(t, "https://first.example.com/player", preview.Embed.IframeURL)
	require.Zero(t, preview.Embed.Width)  // the winner's unknown dimensions
	require.Zero(t, preview.Embed.Height) // never the loser's
}

func TestD5_ProviderWinsWholeAndIsNeverCompleted(t *testing.T) {

	var preview Preview

	// A name-only provider passes the floor and wins whole.
	merge(&preview, partial{Provider: null.NewObject(Provider{Name: "First"})})

	merge(&preview, partial{Provider: null.NewObject(Provider{
		Name:    "Second",
		URL:     "https://second.example.com",
		IconURL: "https://second.example.com/icon.png",
	})})

	require.Equal(t, "First", preview.Provider.Name)
	require.Empty(t, preview.Provider.URL)     // never paired across sources
	require.Empty(t, preview.Provider.IconURL) // the icon backfill is in extract(), not merge
}

func TestD5_AuthorNameAndURLTravelTogether(t *testing.T) {

	var preview Preview

	merge(&preview, partial{Authors: []Author{{Name: "First Author"}}})
	merge(&preview, partial{Authors: []Author{{Name: "Second Author", URL: "https://second.example.com/author"}}})

	require.Len(t, preview.Authors, 1)
	require.Equal(t, "First Author", preview.Authors[0].Name)
	require.Empty(t, preview.Authors[0].URL) // NOT the second source's URL
}

func TestD5_AuthorListWinsWholeFromOneSource(t *testing.T) {

	var preview Preview

	merge(&preview, partial{Authors: []Author{{Name: "Ada"}, {Name: "Grace"}}})
	merge(&preview, partial{Authors: []Author{{Name: "Alan"}, {Name: "Edsger"}}})

	// The lists are never concatenated — one source supplies the whole list.
	require.Len(t, preview.Authors, 2)
	require.Equal(t, "Ada", preview.Authors[0].Name)
	require.Equal(t, "Grace", preview.Authors[1].Name)
}

func TestD5_GroupsAreIndependentOfEachOther(t *testing.T) {

	// Winning one group does not claim the others: different sources may win
	// different groups, and each still arrives whole.
	var preview Preview

	merge(&preview, partial{Thumbnail: null.NewObject(Thumbnail{URL: "https://first.example.com/a.jpg", Width: 1200, Height: 630})})
	merge(&preview, partial{
		Thumbnail: null.NewObject(Thumbnail{URL: "https://second.example.com/b.jpg"}),
		Embed:     null.NewObject(Embed{Mode: EmbedIframe, IframeURL: "https://second.example.com/player", Width: 640, Height: 360}),
		Provider:  null.NewObject(Provider{Name: "Second"}),
	})

	require.Equal(t, "https://first.example.com/a.jpg", preview.Thumbnail.URL)
	require.Equal(t, "https://second.example.com/player", preview.Embed.IframeURL)
	require.Equal(t, 640, preview.Embed.Width)
	require.Equal(t, "Second", preview.Provider.Name)
}

/******************************************
 * D6 — Validity Floors
 ******************************************/

// A group below its floor is skipped WHOLE, so the next source's group wins
// whole. A floor disqualifies a source; it never reorders them or merges parts.

func TestD6_ThumbnailFloorSkipsWholeGroup(t *testing.T) {

	var preview Preview

	// A 64x64 logo fails the minimum-edge floor.
	merge(&preview, partial{Thumbnail: null.NewObject(Thumbnail{URL: "https://first.example.com/logo.png", Width: 64, Height: 64})})
	merge(&preview, partial{Thumbnail: null.NewObject(Thumbnail{URL: "https://second.example.com/hero.jpg", Width: 1200, Height: 630})})

	require.Equal(t, "https://second.example.com/hero.jpg", preview.Thumbnail.URL)
	require.Equal(t, 1200, preview.Thumbnail.Width) // the whole winning group, not a graft
	require.Equal(t, 630, preview.Thumbnail.Height)
}

func TestD6_EmbedFloorSkipsFileMasqueradingAsPlayer(t *testing.T) {

	var preview Preview

	// An iframe-mode embed pointing at a media FILE fails the floor.
	merge(&preview, partial{Embed: null.NewObject(Embed{Mode: EmbedIframe, IframeURL: "https://first.example.com/video.mp4"})})
	merge(&preview, partial{Embed: null.NewObject(Embed{Mode: EmbedIframe, IframeURL: "https://second.example.com/embed/123", Width: 640, Height: 360})})

	require.Equal(t, "https://second.example.com/embed/123", preview.Embed.IframeURL)
	require.Equal(t, 640, preview.Embed.Width)
}

func TestD6_EmbedFloorAcceptsTheSameURLAsAStream(t *testing.T) {

	var preview Preview

	// The URL that failed as a player passes as an explicit stream.
	merge(&preview, partial{Embed: null.NewObject(Embed{Mode: EmbedStream, StreamURL: "https://first.example.com/video.mp4"})})
	merge(&preview, partial{Embed: null.NewObject(Embed{Mode: EmbedIframe, IframeURL: "https://second.example.com/embed/123"})})

	require.Equal(t, EmbedStream, preview.Embed.Mode)
	require.Equal(t, "https://first.example.com/video.mp4", preview.Embed.StreamURL)
}

func TestD6_ProviderFloorSkipsEmptyGroup(t *testing.T) {

	var preview Preview

	// Neither a name nor a usable URL: nothing worth showing.
	merge(&preview, partial{Provider: null.NewObject(Provider{IconURL: "https://first.example.com/icon.png"})})
	merge(&preview, partial{Provider: null.NewObject(Provider{Name: "Second"})})

	require.Equal(t, "Second", preview.Provider.Name)
	require.Empty(t, preview.Provider.IconURL) // the failed group contributed nothing
}

func TestD6_AuthorFloorSkipsURLAsName(t *testing.T) {

	var preview Preview

	// og:author often carries a profile URL, not a name — it must not win.
	merge(&preview, partial{Authors: []Author{{Name: "https://first.example.com/profile/jane"}}})
	merge(&preview, partial{Authors: []Author{{Name: "Jane Doe", URL: "https://second.example.com/profile/jane"}}})

	require.Len(t, preview.Authors, 1)
	require.Equal(t, "Jane Doe", preview.Authors[0].Name)
	require.Equal(t, "https://second.example.com/profile/jane", preview.Authors[0].URL)
}

func TestD6_MixedAuthorListDropsOnlyTheFailures(t *testing.T) {

	var preview Preview

	// The floor filters WITHIN one source's list — that is not cross-source
	// mixing, because every surviving author came from the same source.
	merge(&preview, partial{Authors: []Author{
		{Name: "https://example.com/profile/jane"},
		{Name: "Jane Doe", URL: "https://example.com/profile/jane"},
		{Name: ""},
	}})

	require.Len(t, preview.Authors, 1)
	require.Equal(t, "Jane Doe", preview.Authors[0].Name)
}

func TestD6_AllAuthorsFailingLeavesTheListOpen(t *testing.T) {

	var preview Preview

	// A list whose every entry fails the floor must not claim the slot.
	merge(&preview, partial{Authors: []Author{{Name: "https://first.example.com/profile"}}})
	merge(&preview, partial{Authors: []Author{{Name: "Second Author"}}})

	require.Len(t, preview.Authors, 1)
	require.Equal(t, "Second Author", preview.Authors[0].Name)
}
