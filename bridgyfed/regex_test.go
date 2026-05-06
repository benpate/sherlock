package bridgyfed

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLooksLikeBlueSky(t *testing.T) {

	require.True(t, LooksLikeBluesky("user.blue.sky"))
	require.True(t, LooksLikeBluesky("@user.blue.sky"))
	require.True(t, LooksLikeBluesky("@user.blue.sky.eu"))
	require.True(t, LooksLikeBluesky("@only.one.at.sign.com"))
	require.True(t, LooksLikeBluesky("@names.can-have.dashes.biz"))
	require.True(t, LooksLikeBluesky("or.no.at.signs.net"))

	require.False(t, LooksLikeBluesky("not-enough-segments.net"))
	require.False(t, LooksLikeBluesky("yomama.bsky.so"))
	require.False(t, LooksLikeBluesky(".cant.start.with.dot.com"))
	require.False(t, LooksLikeBluesky("@cant/have/slashes.eu"))
	require.False(t, LooksLikeBluesky("@cant-have@multiple.at.signs.com"))
	require.False(t, LooksLikeBluesky("no:colons.in.it"))
}
