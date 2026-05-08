//go:build localonly

package webfinger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientIsWebfinger(t *testing.T) {

	client := New(nil).(Client)

	do := func(original string, shouldBeWebFinger bool, expectedURI string) {

		if expectedURI == "" {
			expectedURI = original
		}

		isWebFinger, newURI, err := client.isWebfinger(original)

		require.NoError(t, err)
		require.Equal(t, shouldBeWebFinger, isWebFinger)
		require.Equal(t, expectedURI, newURI)
	}

	do("https://example.com/actor/123", false, "")
	do("http://127.0.0.1/@69cbde561cc0c83dd0a31547/pub/keyPackages/69fdf7f965af643d705b1472", false, "")
	do("benpate@mastodon.social", true, "https://mastodon.social/users/benpate")
}
