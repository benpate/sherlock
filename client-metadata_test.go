package sherlock

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestClientMetadata_Delegation proves Client.Metadata threads the Client's
// configuration into the metadata engine: the SSRF override reaches the
// fetch (or the httptest server would be blocked), the User-Agent reaches
// the request, and the extracted Preview comes back intact.
func TestClientMetadata_Delegation(t *testing.T) {

	var userAgent string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, `<html><head>
			<title>Delegation | Example</title>
			<meta property="og:title" content="Delegation Test">
		</head><body></body></html>`)
	}))
	defer server.Close()

	client := NewClient(WithUserAgent("sherlock-test/1.0"))

	preview, err := client.Metadata(context.Background(), server.URL, WithAllowPrivateIPs(true))

	require.NoError(t, err)
	require.Equal(t, "Delegation Test", preview.Title)
	require.Equal(t, "sherlock-test/1.0", userAgent)
}

// TestClientMetadata_EmptyURL proves the engine's input validation surfaces
// through the delegating method.
func TestClientMetadata_EmptyURL(t *testing.T) {

	client := NewClient()

	_, err := client.Metadata(context.Background(), "")
	require.Error(t, err)
}
