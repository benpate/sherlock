package metadata

import (
	"testing"
	"time"

	"github.com/benpate/rosetta/null"
	"github.com/stretchr/testify/require"
)

// precedenceOrder is the spec's global order. merge() call order IS this
// order, for EVERY merged field — there are no per-field exceptions left.
var precedenceOrder = []string{"as", "og", "tw", "oe", "html"}

// precedenceKinds gives each source a distinct, non-Unknown Kind, so a merged
// Kind names the source that won it.
var precedenceKinds = []Kind{KindArticle, KindVideo, KindAudio, KindImage, KindProfile}

// precedenceBaseTime anchors the tagged timestamps.
var precedenceBaseTime = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

// taggedPartial builds a partial in which EVERY mergeable field is populated
// with values carrying the source's tag, so a merged Preview can be traced
// back to exactly which source won each field.
func taggedPartial(tag string, index int) partial {

	// Values must be past their validity floors (D6), or the group would be
	// skipped and the test would be measuring the floor instead of precedence.

	return partial{
		Title:         null.NewString(tag + " title"),
		Description:   null.NewString(tag + " description"),
		Language:      null.NewString(tag + "-LANG"),
		CreatorHandle: null.NewString("@" + tag + "@example.social"),
		Kind:          null.NewObject(precedenceKinds[index]),

		Thumbnail: null.NewObject(Thumbnail{
			URL:    "https://" + tag + ".example.com/thumb.jpg",
			Width:  1200,
			Height: 630,
			Alt:    tag + " alt",
		}),

		Embed: null.NewObject(Embed{
			Mode:      EmbedIframe,
			IframeURL: "https://" + tag + ".example.com/player",
			Width:     640,
			Height:    360,
		}),

		Provider: null.NewObject(Provider{
			Name:    tag + " provider",
			URL:     "https://" + tag + ".example.com",
			IconURL: "https://" + tag + ".example.com/icon.png",
		}),

		Authors: []Author{{
			Name: tag + " author",
			URL:  "https://" + tag + ".example.com/author",
		}},

		PublishedAt: null.NewObject(precedenceBaseTime.AddDate(0, 0, index)),
		ModifiedAt:  null.NewObject(precedenceBaseTime.AddDate(0, 0, index+10)),
	}
}

// requireAllFieldsFrom asserts that every merged field of the Preview came
// from the named source — no field, and no field WITHIN a group, from another.
func requireAllFieldsFrom(t *testing.T, tag string, index int, preview Preview) {
	t.Helper()

	expected := taggedPartial(tag, index)

	require.Equal(t, expected.Title.String(), preview.Title)
	require.Equal(t, expected.Description.String(), preview.Description)
	require.Equal(t, expected.Language.String(), preview.Language)
	require.Equal(t, expected.CreatorHandle.String(), preview.CreatorHandle)
	require.Equal(t, expected.Kind.Object(), preview.Kind)

	require.NotNil(t, preview.Thumbnail)
	require.Equal(t, expected.Thumbnail.Object(), *preview.Thumbnail)

	require.NotNil(t, preview.Embed)
	require.Equal(t, expected.Embed.Object(), *preview.Embed)

	require.NotNil(t, preview.Provider)
	require.Equal(t, expected.Provider.Object(), *preview.Provider)

	require.Equal(t, expected.Authors, preview.Authors)

	require.NotNil(t, preview.PublishedAt)
	require.Equal(t, expected.PublishedAt.Object(), *preview.PublishedAt)

	require.NotNil(t, preview.ModifiedAt)
	require.Equal(t, expected.ModifiedAt.Object(), *preview.ModifiedAt)
}

func TestMerge_PrecedenceIsCallOrder(t *testing.T) {

	// For every starting position: sources before the winner contribute
	// nothing, the winner and everyone after contribute everything, and the
	// winner must take every single field.
	for winner, winnerTag := range precedenceOrder {

		t.Run(winnerTag, func(t *testing.T) {

			var preview Preview

			for range precedenceOrder[:winner] {
				merge(&preview, partial{})
			}

			for index := winner; index < len(precedenceOrder); index++ {
				merge(&preview, taggedPartial(precedenceOrder[index], index))
			}

			requireAllFieldsFrom(t, winnerTag, winner, preview)
		})
	}
}

func TestMerge_EmbedFollowsTheSameOrderAsEveryOtherGroup(t *testing.T) {

	// RULE: Embed has no exception pick. It is decided by call order, exactly
	// like Thumbnail and Provider. This is the test that would have failed
	// under the old firstValidEmbed(oe, tw, og) override.
	for winner, winnerTag := range precedenceOrder {

		t.Run(winnerTag, func(t *testing.T) {

			var preview Preview

			for range precedenceOrder[:winner] {
				merge(&preview, partial{})
			}

			for index := winner; index < len(precedenceOrder); index++ {
				merge(&preview, partial{Embed: taggedPartial(precedenceOrder[index], index).Embed})
			}

			require.NotNil(t, preview.Embed)
			require.Equal(t, "https://"+winnerTag+".example.com/player", preview.Embed.IframeURL)
		})
	}
}

func TestMerge_NeverOverwrites(t *testing.T) {

	// Once a field is set, no later source may change it — the property the
	// whole merge design rests on (D3).
	var preview Preview

	for index, tag := range precedenceOrder {
		merge(&preview, taggedPartial(tag, index))
	}

	requireAllFieldsFrom(t, precedenceOrder[0], 0, preview)
}

func TestMerge_FillsForwardFieldByField(t *testing.T) {

	// Scalars fill INDIVIDUALLY: a later source supplies what earlier sources
	// lacked, without displacing anything they did supply.
	var preview Preview

	merge(&preview, partial{Title: null.NewString("as title")})
	merge(&preview, partial{Title: null.NewString("og title"), Description: null.NewString("og description")})
	merge(&preview, partial{Description: null.NewString("tw description"), Language: null.NewString("tw-LANG")})

	require.Equal(t, "as title", preview.Title)
	require.Equal(t, "og description", preview.Description)
	require.Equal(t, "tw-LANG", preview.Language)
}

func TestMerge_EmptySourcesAreInert(t *testing.T) {

	// An empty partial anywhere in the chain must not shift precedence.
	var withGaps Preview

	merge(&withGaps, partial{})
	merge(&withGaps, taggedPartial("og", 1))
	merge(&withGaps, partial{})
	merge(&withGaps, taggedPartial("html", 4))

	requireAllFieldsFrom(t, "og", 1, withGaps)
}
