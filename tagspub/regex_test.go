package tagspub

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLooksLikeHashtag(t *testing.T) {

	require.True(t, looksLikeHashtag.Match([]byte("#hashtag")))
	require.True(t, looksLikeHashtag.Match([]byte("#hashtag_with_underscores")))
	require.True(t, looksLikeHashtag.Match([]byte("#hashtag_with_numbers_12345")))

	require.False(t, looksLikeHashtag.Match([]byte("#no-dashes")))
	require.False(t, looksLikeHashtag.Match([]byte("#no.dots")))
	require.False(t, looksLikeHashtag.Match([]byte("#no/slashes")))
	require.False(t, looksLikeHashtag.Match([]byte("no:colons.in.it")))
}

func TestIsHashtag(t *testing.T) {

	do := func(input string, expectedMatch bool, expectedValue string) {

		match, value := IsHashtag(input)

		require.Equal(t, expectedMatch, match)
		require.Equal(t, expectedValue, value)
	}

	do("#hashtag", true, "hashtag")
	do("#hashtag_with_underscores", true, "hashtag_with_underscores")
	do("#hashtag_with_numbers_12345", true, "hashtag_with_numbers_12345")

	do("", false, "")
	do("#no-dashes", false, "")
	do("#no.dots", false, "")
	do("#no/slashes", false, "")
	do("no:colons.in.it", false, "")
}
