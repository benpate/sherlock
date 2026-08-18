package metadata

import (
	"testing"

	"github.com/benpate/oembed"
	"github.com/stretchr/testify/require"
)

func TestOEmbedCouldHelp(t *testing.T) {

	// A preview carrying every field oEmbed can supply cannot be improved.
	complete := Preview{
		Title:     "Title",
		Thumbnail: &Thumbnail{URL: "https://example.com/t.jpg"},
		Provider:  &Provider{Name: "Example"},
		Authors:   []Author{{Name: "Jane"}},
	}
	require.False(t, oembedCouldHelp(complete))

	// Any one of the four missing re-opens the gate.
	for _, testCase := range []struct {
		name    string
		mutate  func(*Preview)
		missing string
	}{
		{"title", func(p *Preview) { p.Title = "" }, "Title"},
		{"thumbnail", func(p *Preview) { p.Thumbnail = nil }, "Thumbnail"},
		{"provider", func(p *Preview) { p.Provider = nil }, "Provider"},
		{"authors", func(p *Preview) { p.Authors = nil }, "Authors"},
	} {
		preview := complete
		testCase.mutate(&preview)
		require.True(t, oembedCouldHelp(preview), "missing %s should re-open the gate", testCase.missing)
	}

	// RULE: description and language are NOT tested — the oEmbed spec has
	// neither, so a preview missing only those cannot be improved by oEmbed.
	noDescription := complete
	noDescription.Description = ""
	noDescription.Language = ""
	require.False(t, oembedCouldHelp(noDescription))

	// RULE: a missing embed is NOT tested either — no preview has one at the
	// gate, so testing it would make the gate fire every time.
	noEmbed := complete
	noEmbed.Embed = nil
	require.False(t, oembedCouldHelp(noEmbed))

	// A blank preview is improvable by definition.
	require.True(t, oembedCouldHelp(Preview{}))
}

func TestClassifyEmbed_BoundsProviderDimensions(t *testing.T) {

	// A provider's declared player size is remote input like any other, and
	// the EmbedPolicy passed by classifyEmbed sets no clamp of its own.
	oversized := oembed.Response{
		Type:   "video",
		HTML:   `<iframe src="https://example.com/player"></iframe>`,
		Width:  999999999999,
		Height: 888888888888,
	}

	iframe := classifyEmbed(oversized, false)

	require.True(t, iframe.IsPresent())
	require.Equal(t, EmbedIframe, iframe.Object().Mode)
	require.Zero(t, iframe.Object().Width)
	require.Zero(t, iframe.Object().Height)

	// The sandbox plan reads the same dimensions and must bound them too.
	sandbox := classifyEmbed(oembed.Response{
		Type:   "rich",
		HTML:   `<div><script>player()</script></div>`,
		Width:  999999999999,
		Height: 480,
	}, true)

	require.True(t, sandbox.IsPresent())
	require.Equal(t, EmbedProviderHTML, sandbox.Object().Mode)
	require.Zero(t, sandbox.Object().Width)

	// A sane dimension still survives — this floor zeroes junk, not everything.
	require.Equal(t, 480, sandbox.Object().Height)
}
