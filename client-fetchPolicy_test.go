package sherlock

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"testing/fstest"

	"github.com/benpate/remote"
	"github.com/benpate/remote/options"
	"github.com/stretchr/testify/require"
)

/******************************************
 * Fetch Policy
 *
 * These tests pin the promises Config makes about every request Load
 * makes: a bounded link walk, a User-Agent, the SSRF setting, and the
 * caller's remote options -- on EVERY fetch, not just the first one.
 ******************************************/

// TestLoad_RedirectBudgetIsBounded proves that two pages linking to each other
// cannot make the link walk recurse without end. Before the budget was spent
// before recursing (rather than after, on a discarded copy of Config), this
// walk ran until the stack gave out.
func TestLoad_RedirectBudgetIsBounded(t *testing.T) {

	pageA := "HTTP/1.1 200 OK\nContent-Type: text/html\nConnection: close\n\n" +
		`<html><head><link rel="alternate" type="application/rss+xml" href="https://test-server.local/b.html"></head><body>A</body></html>`

	pageB := "HTTP/1.1 200 OK\nContent-Type: text/html\nConnection: close\n\n" +
		`<html><head><link rel="alternate" type="application/rss+xml" href="https://test-server.local/a.html"></head><body>B</body></html>`

	filesystem := fstest.MapFS{
		"a.html": {Data: []byte(pageA)},
		"b.html": {Data: []byte(pageB)},
	}

	var fetches atomic.Int64

	// A circuit breaker, so a regression reports a number instead of exhausting
	// the stack and taking the whole test binary down with it.
	circuitBreaker := remote.Option{
		ModifyRequest: func(_ *remote.Transaction, request *http.Request) *http.Response {

			if fetches.Add(1) > 100 {
				return &http.Response{
					Request:    request,
					StatusCode: http.StatusInternalServerError,
					Header:     http.Header{},
					Body:       io.NopCloser(strings.NewReader("circuit breaker tripped")),
				}
			}

			return nil
		},
	}

	var withServers Option = func(config *Config) {
		config.RemoteOptions = append(config.RemoteOptions, circuitBreaker, options.TestServer("test-server.local", filesystem))
	}

	client := NewClient()
	_, _ = client.Load("https://test-server.local/a.html", AsActor(), withServers)

	require.Less(t, fetches.Load(), int64(100), "the link walk did not respect the redirect budget")
}

// TestLoad_DefaultValueIsNotAliased proves that Load treats DefaultValue as
// input, not as scratch space. Writing through it leaked one document's fields
// into the next whenever a caller reused the same map.
func TestLoad_DefaultValueIsNotAliased(t *testing.T) {

	feedOne := "HTTP/1.1 200 OK\nContent-Type: application/feed+json\nConnection: close\n\n" +
		`{"version":"https://jsonfeed.org/version/1","title":"Feed One","feed_url":"https://test-server.local/one.json","items":[]}`

	feedTwo := "HTTP/1.1 200 OK\nContent-Type: application/feed+json\nConnection: close\n\n" +
		`{"version":"https://jsonfeed.org/version/1","title":"Feed Two","feed_url":"https://test-server.local/two.json","items":[]}`

	var withServer Option = func(config *Config) {
		config.RemoteOptions = append(config.RemoteOptions, options.TestServer("test-server.local", fstest.MapFS{
			"one.json": {Data: []byte(feedOne)},
			"two.json": {Data: []byte(feedTwo)},
		}))
	}

	callerMap := map[string]any{"myOwnKey": "untouched"}
	client := NewClient()

	first, err := client.Load("https://test-server.local/one.json", AsActor(), WithDefaultValue(callerMap), withServer)
	require.Nil(t, err)
	require.Equal(t, "Feed One", first.Name())

	// The caller's map is theirs: Load must not have written into it.
	require.Equal(t, map[string]any{"myOwnKey": "untouched"}, callerMap)

	// And the same map, reused, must not carry the first feed's identity into the second.
	second, err := client.Load("https://test-server.local/two.json", AsActor(), WithDefaultValue(callerMap), withServer)
	require.Nil(t, err)
	require.Equal(t, "Feed Two", second.Name())
	require.Equal(t, "https://test-server.local/two.json", second.ID())
}

// TestLoad_MicroFormatsHonorsDefaultValue proves that all three feed formats
// treat DefaultValue alike. The MicroFormats path used to discard it, so the
// same option behaved differently depending on what a site happened to serve.
func TestLoad_MicroFormatsHonorsDefaultValue(t *testing.T) {

	page := "HTTP/1.1 200 OK\nContent-Type: text/html\nConnection: close\n\n" +
		`<html><body><div class="h-feed"><span class="p-name">MF Feed</span>` +
		`<div class="h-entry"><a class="u-url" href="https://test-server.local/post1">Post One</a>` +
		`<span class="p-name">Post One</span></div></div></body></html>`

	var withServer Option = func(config *Config) {
		config.RemoteOptions = append(config.RemoteOptions,
			options.TestServer("test-server.local", fstest.MapFS{"mf.html": {Data: []byte(page)}}))
	}

	client := NewClient()
	doc, err := client.Load("https://test-server.local/mf.html", AsActor(),
		WithDefaultValue(map[string]any{"summary": "CALLER SUPPLIED SUMMARY"}), withServer)

	require.Nil(t, err)
	require.Equal(t, "MF Feed", doc.Name())
	require.Equal(t, "CALLER SUPPLIED SUMMARY", doc.Summary())
}

// TestLoad_AllowPrivateIPs proves the option is wired into Load's own requests,
// and not only into Metadata's. It was declared on Config, and documented, but
// no transaction in the Load path ever read it.
func TestLoad_AllowPrivateIPs(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/activity+json")
		_, _ = w.Write([]byte(`{"type":"Person","id":"https://example.com/actor","name":"Loopback Actor"}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient()

	// TRUE reaches the loopback server...
	doc, err := client.Load(server.URL, AsActor(), WithAllowPrivateIPs(true))
	require.Nil(t, err)
	require.Equal(t, "Loopback Actor", doc.Name())

	// ...and the default FALSE still refuses it.
	_, err = client.Load(server.URL, AsActor())
	require.NotNil(t, err)
}

// TestLoad_EveryFetchCarriesThePolicy proves that the homepage-icon lookup --
// the one fetch that used to build its own bare remote.Get -- travels under the
// same User-Agent and caller options as the request that discovered it.
func TestLoad_EveryFetchCarriesThePolicy(t *testing.T) {

	bodies := map[string]string{
		"/feed.json": "application/feed+json\x00" +
			`{"version":"https://jsonfeed.org/version/1","title":"Iconless","feed_url":"https://test-server.local/feed.json","items":[]}`,
		"": "text/html\x00" +
			`<html><head><link rel="icon" type="image/png" href="/favicon.png"></head><body>home</body></html>`,
	}

	// Serve both the feed and the site homepage in-process, recording the
	// User-Agent of every request that arrives -- including the homepage lookup,
	// which only happens because the feed declared no icon of its own.
	userAgents := make([]string, 0)

	var testServer Option = func(config *Config) {
		config.RemoteOptions = append(config.RemoteOptions, remote.Option{
			ModifyRequest: func(_ *remote.Transaction, request *http.Request) *http.Response {

				if request.URL.Hostname() != "test-server.local" {
					return nil
				}

				userAgents = append(userAgents, request.UserAgent())

				body, found := bodies[request.URL.Path]

				if !found {
					return &http.Response{Request: request, StatusCode: http.StatusNotFound, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(""))}
				}

				contentType, content, _ := strings.Cut(body, "\x00")

				return &http.Response{
					Request:    request,
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{contentType}},
					Body:       io.NopCloser(strings.NewReader(content)),
				}
			},
		})
	}

	client := NewClient(WithUserAgent("Sherlock-Test/1.0"))
	doc, err := client.Load("https://test-server.local/feed.json", AsActor(), testServer)

	require.Nil(t, err)
	require.Equal(t, "https://test-server.local/favicon.png", doc.Icon().Href(), "the homepage icon lookup did not run under the caller's options")

	// Three fetches: the ActivityStreams probe, the feed itself, then the
	// homepage icon lookup. Every one of them must carry the policy.
	require.Len(t, userAgents, 3)
	for _, userAgent := range userAgents {
		require.Equal(t, "Sherlock-Test/1.0", userAgent)
	}
}

// TestLoad_LinkedJSONLDCarriesThePolicy proves the <link rel="alternate"> JSON-LD
// fetch travels under the caller's options and resolves a RELATIVE href against
// the page that declared it. It built a bare remote.Get on the raw attribute
// value, so a relative alternate could never be fetched at all.
func TestLoad_LinkedJSONLDCarriesThePolicy(t *testing.T) {

	page := "HTTP/1.1 200 OK\nContent-Type: text/html\nConnection: close\n\n" +
		`<html><head><link rel="alternate" type="application/activity+json" href="/alternate.json"></head><body>page</body></html>`

	alternate := "HTTP/1.1 200 OK\nContent-Type: application/activity+json\nConnection: close\n\n" +
		`{"type":"Note","id":"https://test-server.local/note/1","name":"From the alternate"}`

	var withServer Option = func(config *Config) {
		config.RemoteOptions = append(config.RemoteOptions, options.TestServer("test-server.local", fstest.MapFS{
			"page.html":      {Data: []byte(page)},
			"alternate.json": {Data: []byte(alternate)},
		}))
	}

	client := NewClient()

	// loadDocument_HTML is reached directly, because Load's ActivityStream attempt
	// errors on an HTML response and returns before the HTML branch runs (BUG-134).
	config := client.newConfig(AsDocument(), withServer)
	doc := client.loadDocument_HTML(config, "https://test-server.local/page.html")

	require.False(t, doc.IsNil())
	require.Equal(t, "From the alternate", doc.Name(), "the relative JSON-LD alternate was not resolved and fetched")
}
