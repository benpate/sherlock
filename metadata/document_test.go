package metadata

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// testConfig returns a config suitable for httptest servers: the SSRF guard
// blocks loopback addresses unless AllowPrivateIPs is TRUE.
func testConfig() config {
	return newConfig(WithAllowPrivateIPs(true))
}

func TestNewDocument_HTML(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Contains(t, r.Header.Get("Accept"), "application/activity+json")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(`<html><head><title>Hello</title></head><body>World</body></html>`))
	}))
	defer server.Close()

	doc, err := newDocument(context.Background(), testConfig(), server.URL)
	require.NoError(t, err)
	require.True(t, doc.IsHTML())
	require.False(t, doc.IsActivityStream())
	require.False(t, doc.Truncated)
	require.Equal(t, "text/html", doc.ContentType)
	require.Equal(t, server.URL, doc.FinalURL)
}

func TestNewDocument_ActivityStream(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/activity+json")
		_, _ = w.Write([]byte(`{"@context":"https://www.w3.org/ns/activitystreams","type":"Note","name":"Hi"}`))
	}))
	defer server.Close()

	doc, err := newDocument(context.Background(), testConfig(), server.URL)
	require.NoError(t, err)
	require.True(t, doc.IsActivityStream())
	require.False(t, doc.IsHTML())
	require.Equal(t, "Note", doc.ActivityStream["type"])
}

func TestNewDocument_Charset_ISO88591(t *testing.T) {

	// "café" in ISO-8859-1: the é is byte 0xE9.
	body := append([]byte(`<html><head><title>caf`), 0xE9)
	body = append(body, []byte(`</title></head></html>`)...)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=iso-8859-1")
		_, _ = w.Write(body)
	}))
	defer server.Close()

	doc, err := newDocument(context.Background(), testConfig(), server.URL)
	require.NoError(t, err)
	require.Contains(t, string(doc.Body), "café")
}

func TestNewDocument_Charset_MetaTag(t *testing.T) {

	// Charset declared only in a <meta> tag; the header is silent.
	body := append([]byte(`<html><head><meta charset="iso-8859-1"><title>caf`), 0xE9)
	body = append(body, []byte(`</title></head></html>`)...)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(body)
	}))
	defer server.Close()

	doc, err := newDocument(context.Background(), testConfig(), server.URL)
	require.NoError(t, err)
	require.Contains(t, string(doc.Body), "café")
}

func TestNewDocument_Truncation(t *testing.T) {

	big := `<html><head><title>Big</title></head><body>` + strings.Repeat("x", 4096) + `</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(big))
	}))
	defer server.Close()

	config := testConfig()
	config.MaxBodySize = 512

	doc, err := newDocument(context.Background(), config, server.URL)
	require.NoError(t, err)
	require.True(t, doc.Truncated)
	require.LessOrEqual(t, len(doc.Body), 512)
	require.True(t, doc.IsHTML()) // the <head> survived the cap

	// A body exactly at the cap is NOT truncated.
	config.MaxBodySize = int64(len(big))
	doc, err = newDocument(context.Background(), config, server.URL)
	require.NoError(t, err)
	require.False(t, doc.Truncated)
}

func TestNewDocument_ContextCancel(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := newDocument(ctx, testConfig(), server.URL)
	require.Error(t, err)
}

func TestNewDocument_SSRFGuardBlocksLoopback(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	config := testConfig()
	config.AllowPrivateIPs = false

	_, err := newDocument(context.Background(), config, server.URL)
	require.Error(t, err)
}

func TestNewDocument_GarbageHTML(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<<<>>>&&& not html at all <div <span`))
	}))
	defer server.Close()

	// html.Parse is best-effort: garbage still yields a tree, never a panic.
	doc, err := newDocument(context.Background(), testConfig(), server.URL)
	require.NoError(t, err)
	require.True(t, doc.IsHTML())
}
