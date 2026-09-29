package activitypub

import (
	"encoding/json"
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/uri"
	"github.com/stretchr/testify/require"
)

// FuzzDocumentID confirms that any id documentID accepts from a remote body is a valid URL, and is
// the same id that streams.Document.ID reports for that body.
func FuzzDocumentID(f *testing.F) {

	f.Add([]byte(`{"id":"https://example.com/users/alice"}`))
	f.Add([]byte(`{"@id":"https://example.com/users/alice"}`))
	f.Add([]byte(`{"id":"https://example.com/a","@id":"https://evil.com/b"}`))
	f.Add([]byte(`{"id":""}`))
	f.Add([]byte(`{"id":42}`))
	f.Add([]byte(`{"id":["https://example.com/a"]}`))
	f.Add([]byte(`{"id":"https://example.com/a<!-- -->"}`))
	f.Add([]byte(`{"id":"https://example.com/a\r"}`))
	f.Add([]byte(`{"id":"https://victim.com<!--@evil.com-->/a"}`))
	f.Add([]byte(`{}`))

	f.Fuzz(func(t *testing.T, data []byte) {

		body := make(map[string]any)

		if err := json.Unmarshal(data, &body); err != nil {
			return
		}

		id, err := documentID(body)

		// A rejected or absent id is always safe
		if (err != nil) || (id == "") {
			return
		}

		require.True(t, uri.IsValidURL(id), "accepted an invalid URL %q", id)
		require.Equal(t, id, streams.NewDocument(body).ID(), "accepted an id that Document.ID rewrites")
	})
}

// FuzzOriginVouches confirms that whenever a serving URL vouches for an id, the id that callers will
// store names the same host as the URL that served it.
func FuzzOriginVouches(f *testing.F) {

	f.Add("https://example.com/users/alice", "https://example.com/users/alice")
	f.Add("https://example.com/users/alice", "https://EXAMPLE.com:8443/other")
	f.Add("https://example.com/users/alice", "https://evil.com/users/alice")
	f.Add("https://example.com.evil.com/a", "https://example.com/a")
	f.Add("https://user@example.com/a", "https://example.com/a")
	f.Add("https://example.com/a<!--x-->", "https://example.com/")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, id string, served string) {

		body := map[string]any{vocab.PropertyID: id}

		claimedID, err := documentID(body)

		if (err != nil) || (claimedID == "") {
			return
		}

		if !sameHost(claimedID, served) {
			return
		}

		stored := streams.NewDocument(body).ID()
		require.NotEmpty(t, uri.Hostname(stored))
		require.Equal(t, uri.Hostname(served), uri.Hostname(stored), "id %q vouched for by %q", id, served)
	})
}
