package sherlock

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestCanTrace(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	require.True(t, canTrace())
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	require.False(t, canTrace())
}

func TestCanDebug(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	require.True(t, canDebug())
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	require.True(t, canDebug())
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	require.False(t, canDebug())
}

func TestCanInfo(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	require.True(t, canDebug())
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	require.True(t, canDebug())
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	require.True(t, canInfo())
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	require.False(t, canInfo())
}

func TestHostOnly(t *testing.T) {
	require.Equal(t, "https://example.com", hostOnly("https://example.com"))
	require.Equal(t, "https://example.com:8080", hostOnly("https://example.com:8080"))
	require.Equal(t, "https://example.com", hostOnly("https://example.com/"))
	require.Equal(t, "https://example.com", hostOnly("https://example.com/some/path/here"))
	require.Equal(t, "https://example.com", hostOnly("https://example.com?query=string"))
	require.Equal(t, "https://example.com", hostOnly("https://example.com/some/path?and=querystring"))

	require.Equal(t, "", hostOnly("example.com"))
}
