package sherlock

import (
	"testing"

	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
	"github.com/stretchr/testify/require"
)

func TestApplyLinks_UnsentTransaction(t *testing.T) {

	client := NewClient()

	// An un-sent transaction has a nil Response(). applyLinks must treat this as a
	// no-op rather than panicking on the nil dereference.
	txn := remote.Get("https://example.com")
	data := mapof.NewAny()

	require.NotPanics(t, func() {
		client.applyLinks(txn, data)
	})
	require.Empty(t, data)
}
