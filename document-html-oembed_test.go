package sherlock

import (
	"strings"
	"testing"

	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

// ParseOEmbed is currently an unimplemented stub. This test documents that it is a
// safe no-op (it neither panics nor mutates the supplied data). Update this test when
// oEmbed parsing is implemented.
func TestParseOEmbed_Noop(t *testing.T) {
	data := mapof.Any{"existing": "value"}

	require.NotPanics(t, func() {
		ParseOEmbed(strings.NewReader(`<link rel="alternate" type="application/json+oembed" href="https://example.com/oembed">`), data)
	})

	require.Equal(t, "value", data["existing"])
	require.Len(t, data, 1)
}
