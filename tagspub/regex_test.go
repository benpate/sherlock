package tagspub

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLooksLikeBlueSky(t *testing.T) {

	require.True(t, looksLikeHashtag.Match([]byte("#hashtag")))
	require.True(t, looksLikeHashtag.Match([]byte("#hashtag_with_underscores")))
	require.True(t, looksLikeHashtag.Match([]byte("#hashtag_with_numbers_12345")))

	require.False(t, looksLikeHashtag.Match([]byte("#no-dashes")))
	require.False(t, looksLikeHashtag.Match([]byte("#no.dots")))
	require.False(t, looksLikeHashtag.Match([]byte("#no/slashes")))
	require.False(t, looksLikeHashtag.Match([]byte("no:colons.in.it")))
}
