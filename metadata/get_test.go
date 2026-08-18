package metadata

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

// serveHTML returns an httptest server that serves one HTML body.
func serveHTML(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(body))
	}))
}

// fetchMetadata runs the full engine against a URL with the loopback guard off.
func fetchMetadata(t *testing.T, url string) Preview {
	t.Helper()
	result, err := Get(context.Background(), url, WithAllowPrivateIPs(true))
	require.NoError(t, err)
	return result
}

func TestMetadata_NewsArticle(t *testing.T) {

	server := serveHTML(`<!DOCTYPE html><html lang="en-GB"><head>
		<title>Big Story | Example News</title>
		<meta name="description" content="Fallback description">
		<link rel="canonical" href="/articles/big-story">
		<meta property="og:title" content="Big Story">
		<meta property="og:description" content="The OG description">
		<meta property="og:type" content="article">
		<meta property="og:site_name" content="Example News">
		<meta property="og:locale" content="en_GB">
		<meta property="og:image" content="/hero.jpg">
		<meta property="og:image:width" content="1200">
		<meta property="og:image:height" content="630">
		<meta property="og:image:alt" content="A hero image">
		<meta property="article:published_time" content="2026-08-01T10:00:00Z">
		<meta property="article:modified_time" content="2026-08-02T11:00:00Z">
		<meta property="og:author" content="Jane Doe">
		<meta property="article:author" content="/staff/jane">
		<meta name="fediverse:creator" content="@jane@example.social">
		<link rel="apple-touch-icon" sizes="180x180" href="/touch.png">
	</head><body></body></html>`)
	defer server.Close()

	m := fetchMetadata(t, server.URL)

	require.Equal(t, "Big Story", m.Title)
	require.Equal(t, "The OG description", m.Description)
	require.Equal(t, "en-GB", m.Language)
	require.Equal(t, KindArticle, m.Kind)
	require.Equal(t, server.URL+"/articles/big-story", m.URL) // rel=canonical beats og:url

	require.NotNil(t, m.Thumbnail)
	require.Equal(t, server.URL+"/hero.jpg", m.Thumbnail.URL)
	require.Equal(t, 1200, m.Thumbnail.Width)
	require.Equal(t, "A hero image", m.Thumbnail.Alt)

	require.NotNil(t, m.Provider)
	require.Equal(t, "Example News", m.Provider.Name)
	require.Equal(t, server.URL+"/touch.png", m.Provider.IconURL) // icon backfill

	require.Len(t, m.Authors, 1)
	require.Equal(t, "Jane Doe", m.Authors[0].Name)
	require.Equal(t, server.URL+"/staff/jane", m.Authors[0].URL)

	require.NotNil(t, m.PublishedAt)
	require.NotNil(t, m.ModifiedAt)
	require.Equal(t, "@jane@example.social", m.CreatorHandle)
	require.Nil(t, m.Embed)
}

func TestMetadata_BarePage(t *testing.T) {

	server := serveHTML(`<html><head>
		<title>Plain Page — My Site</title>
		<meta name="description" content="Just a description">
	</head><body></body></html>`)
	defer server.Close()

	m := fetchMetadata(t, server.URL)

	require.Equal(t, "Plain Page", m.Title) // site suffix stripped
	require.Equal(t, "Just a description", m.Description)
	require.Equal(t, KindWebsite, m.Kind)
	require.Equal(t, server.URL, m.URL) // falls back to the final URL
}

func TestMetadata_BlankCard(t *testing.T) {

	server := serveHTML(`<html><head></head><body>nothing here</body></html>`)
	defer server.Close()

	m := fetchMetadata(t, server.URL)

	require.Equal(t, KindWebsite, m.Kind)
	require.Equal(t, server.URL, m.URL)
	require.Empty(t, m.Title)
	require.Nil(t, m.Thumbnail)
}

func TestMetadata_PositionalImageTrap(t *testing.T) {

	// The width/height belong to the FIRST og:image, not the second.
	server := serveHTML(`<html><head>
		<meta property="og:image" content="https://cdn.example.com/first.jpg">
		<meta property="og:image:width" content="800">
		<meta property="og:image:height" content="600">
		<meta property="og:image" content="https://cdn.example.com/second.jpg">
	</head><body></body></html>`)
	defer server.Close()

	m := fetchMetadata(t, server.URL)

	require.NotNil(t, m.Thumbnail)
	require.Equal(t, "https://cdn.example.com/first.jpg", m.Thumbnail.URL)
	require.Equal(t, 800, m.Thumbnail.Width)
}

func TestMetadata_CrossOriginCanonicalRejected(t *testing.T) {

	server := serveHTML(`<html><head>
		<link rel="canonical" href="https://evil.example.net/steal">
		<meta property="og:url" content="https://evil.example.net/steal-og">
	</head><body></body></html>`)
	defer server.Close()

	m := fetchMetadata(t, server.URL)
	require.Equal(t, server.URL, m.URL) // both rejected, falls back to final URL
}

func TestMetadata_EntityDecodeOnce(t *testing.T) {

	server := serveHTML(`<html><head>
		<meta property="og:title" content="Tom &amp;amp; Jerry">
	</head><body></body></html>`)
	defer server.Close()

	m := fetchMetadata(t, server.URL)

	// The HTML parser decodes the attribute once (&amp;amp; → &amp;), and the
	// extractor's decodeOnce decodes once more — never a third time.
	require.Equal(t, "Tom & Jerry", m.Title)
}

func TestMetadata_ActivityStream(t *testing.T) {

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/note", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/activity+json")
		_, _ = fmt.Fprintf(w, `{
			"@context": "https://www.w3.org/ns/activitystreams",
			"id": "%[1]s/note",
			"type": "Note",
			"name": "A Federated Note",
			"summary": "Straight from the source",
			"published": "2026-08-10T12:00:00Z",
			"attributedTo": {"type": "Person", "name": "Alice", "url": "%[1]s/@alice"},
			"image": {"type": "Image", "url": "%[1]s/pic.jpg", "width": 1000, "height": 500}
		}`, server.URL)
	})

	m := fetchMetadata(t, server.URL+"/note")

	require.Equal(t, "A Federated Note", m.Title)
	require.Equal(t, "Straight from the source", m.Description)
	require.Equal(t, server.URL+"/note", m.URL)
	require.Equal(t, KindArticle, m.Kind)
	require.Len(t, m.Authors, 1)
	require.Equal(t, "Alice", m.Authors[0].Name)
	require.NotNil(t, m.Thumbnail)
	require.Equal(t, 1000, m.Thumbnail.Width)
	require.NotNil(t, m.PublishedAt)
}

func TestMetadata_AlternateDiscovery(t *testing.T) {

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// An HTML page that advertises an AS2 alternate (FEP-22b6) and carries
	// its own (worse) OG title. AS outranks OG after the re-merge.
	mux.HandleFunc("/page", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, `<html><head>
			<link rel="alternate" type="application/activity+json" href="%s/object">
			<meta property="og:title" content="Scraped Title">
			<meta property="og:image" content="%s/og.jpg">
		</head><body></body></html>`, server.URL, server.URL)
	})

	mux.HandleFunc("/object", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/activity+json")
		_, _ = fmt.Fprintf(w, `{"id": "%s/object", "type": "Article", "name": "Federated Title"}`, server.URL)
	})

	m := fetchMetadata(t, server.URL+"/page")

	require.Equal(t, "Federated Title", m.Title)  // AS wins the field it fills
	require.Equal(t, server.URL+"/object", m.URL) // AS id wins the URL pick
	require.NotNil(t, m.Thumbnail)                // OG still fills what AS lacked
	require.Equal(t, server.URL+"/og.jpg", m.Thumbnail.URL)
}

func TestMetadata_TwitterPlayer(t *testing.T) {

	server := serveHTML(`<html><head>
		<meta name="twitter:card" content="player">
		<meta name="twitter:title" content="A Video">
		<meta name="twitter:player" content="https://example.com/embed/42">
		<meta name="twitter:player:width" content="640">
		<meta name="twitter:player:height" content="360">
		<meta name="twitter:image" content="https://example.com/poster.jpg">
		<meta name="twitter:image:alt" content="Poster frame">
	</head><body></body></html>`)
	defer server.Close()

	m := fetchMetadata(t, server.URL)

	require.Equal(t, "A Video", m.Title)
	require.Equal(t, KindVideo, m.Kind)

	require.NotNil(t, m.Embed)
	require.Equal(t, EmbedIframe, m.Embed.Mode)
	require.Equal(t, "https://example.com/embed/42", m.Embed.IframeURL)
	require.Equal(t, 640, m.Embed.Width)

	// Player dimensions never leak into the thumbnail group.
	require.NotNil(t, m.Thumbnail)
	require.Equal(t, "https://example.com/poster.jpg", m.Thumbnail.URL)
	require.Zero(t, m.Thumbnail.Width)
	require.Equal(t, "Poster frame", m.Thumbnail.Alt)
}

func TestMetadata_OEmbedDiscovery(t *testing.T) {

	var endpointCalls atomic.Int32

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// A video page advertising an oEmbed endpoint; the provider sends
	// "width" as a string, which the oembed library tolerates.
	mux.HandleFunc("/watch", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, `<html><head>
			<link rel="alternate" type="application/json+oembed" href="%s/oembed?url=x">
			<meta property="og:type" content="video.other">
			<meta property="og:title" content="Watch This">
		</head><body></body></html>`, server.URL)
	})

	mux.HandleFunc("/oembed", func(w http.ResponseWriter, r *http.Request) {
		endpointCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{
			"version": "1.0", "type": "video", "title": "Watch This",
			"html": "<iframe src=\"https://player.example.com/embed/1\"></iframe>",
			"width": "480", "height": 270,
			"author_name": "Creator", "author_url": "%s/@creator",
			"provider_name": "TestTube", "provider_url": "%s/",
			"thumbnail_url": "%s/thumb.jpg", "thumbnail_width": 480, "thumbnail_height": 360
		}`, server.URL, server.URL, server.URL)
	})

	m := fetchMetadata(t, server.URL+"/watch")

	require.Equal(t, int32(1), endpointCalls.Load())
	require.Equal(t, "Watch This", m.Title)
	require.Equal(t, KindVideo, m.Kind)

	// The classifier extracted the single clean iframe and rebuilt it — the
	// provider's raw HTML never reaches the model (Phase 7).
	require.NotNil(t, m.Embed)
	require.Equal(t, EmbedIframe, m.Embed.Mode)
	require.Equal(t, "https://player.example.com/embed/1", m.Embed.IframeURL)
	require.Empty(t, m.Embed.HTML)
	require.Equal(t, 480, m.Embed.Width) // "480" absorbed by oembed's Int

	require.Len(t, m.Authors, 1)
	require.Equal(t, "Creator", m.Authors[0].Name)
	require.Equal(t, "TestTube", m.Provider.Name)
	require.Equal(t, server.URL+"/thumb.jpg", m.Thumbnail.URL)
}

func TestMetadata_OEmbedGateSkipsCompleteArticle(t *testing.T) {

	var endpointCalls atomic.Int32

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// A complete article — title, thumbnail, provider, author all present, no
	// embed signals — must NOT trigger the advertised oEmbed endpoint.
	mux.HandleFunc("/article", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, `<html><head>
			<link rel="alternate" type="application/json+oembed" href="%s/oembed?url=x">
			<meta property="og:title" content="Complete Article">
			<meta property="og:type" content="article">
			<meta property="og:site_name" content="Example">
			<meta property="og:image" content="%s/hero.jpg">
			<meta property="og:author" content="Jane Doe">
		</head><body></body></html>`, server.URL, server.URL)
	})

	mux.HandleFunc("/oembed", func(w http.ResponseWriter, r *http.Request) {
		endpointCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"1.0","type":"link","title":"Should Not Be Called"}`))
	})

	m := fetchMetadata(t, server.URL+"/article")

	require.Equal(t, int32(0), endpointCalls.Load()) // the lazy gate held
	require.Equal(t, "Complete Article", m.Title)
}

// messyWidgetServer serves a page whose oEmbed endpoint returns script-laden
// widget HTML that cannot be reduced to a single clean iframe.
func messyWidgetServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	mux.HandleFunc("/page", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, `<html><head>
			<link rel="alternate" type="application/json+oembed" href="%s/oembed?url=x">
			<meta property="og:type" content="video.other">
		</head><body></body></html>`, server.URL)
	})

	mux.HandleFunc("/oembed", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"1.0","type":"rich","title":"Widget",
			"html":"<blockquote>x</blockquote><script src=\"https://w.example/w.js\"></script>"}`))
	})

	return server
}

func TestMetadata_MessyEmbedDegradesByDefault(t *testing.T) {

	server := messyWidgetServer(t)

	// Default policy: sandbox embeds are OFF, so non-extractable provider HTML
	// yields NO embed in the model — never raw provider markup.
	m := fetchMetadata(t, server.URL+"/page")

	require.Equal(t, "Widget", m.Title)
	require.Nil(t, m.Embed)
}

func TestMetadata_MessyEmbedSandboxedWhenAllowed(t *testing.T) {

	server := messyWidgetServer(t)

	// The operator opts in: the same markup becomes a provider-HTML embed
	// that the consumer renders inside a sandboxed iframe.
	m, err := Get(context.Background(), server.URL+"/page",
		WithAllowPrivateIPs(true),
		WithAllowSandboxEmbeds(true),
	)

	require.NoError(t, err)
	require.NotNil(t, m.Embed)
	require.Equal(t, EmbedProviderHTML, m.Embed.Mode)
	require.Contains(t, m.Embed.HTML, "<blockquote>")
}

func TestMetadata_OEmbedAbsurdDimensionsAreBounded(t *testing.T) {

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	// A hostile provider declares a negative width and an absurd height.
	// Neither may reach the model: negative → unknown, absurd → unknown, and
	// the thumbnail group itself still wins because its URL is usable.
	mux.HandleFunc("/page", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, `<html><head>
			<link rel="alternate" type="application/json+oembed" href="%s/oembed?url=x">
			<meta property="og:type" content="video.other">
		</head><body></body></html>`, server.URL)
	})

	mux.HandleFunc("/oembed", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"1.0","type":"video","title":"Bounded",
			"html":"<iframe src=\"https://player.example.com/v/1\"></iframe>",
			"thumbnail_url":"https://cdn.example.com/t.jpg",
			"thumbnail_width":-1,"thumbnail_height":9223372036854775807}`))
	})

	m := fetchMetadata(t, server.URL+"/page")

	require.NotNil(t, m.Thumbnail)
	require.Equal(t, "https://cdn.example.com/t.jpg", m.Thumbnail.URL)
	require.Equal(t, 0, m.Thumbnail.Width)
	require.Equal(t, 0, m.Thumbnail.Height)
}

/******************************************
 * Embed Precedence, End to End
 ******************************************/

// embedRaceServer serves a page carrying an in-page embed AND an oEmbed
// endpoint that returns a different player, so the winner names the rule.
func embedRaceServer(t *testing.T, inPageTags string, calls *atomic.Int32) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	mux.HandleFunc("/page", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = fmt.Fprintf(w, `<html><head>
			<link rel="alternate" type="application/json+oembed" href="%s/oembed?url=x">
			<meta property="og:title" content="Race">
			%s
		</head><body></body></html>`, server.URL, inPageTags)
	})

	mux.HandleFunc("/oembed", func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"1.0","type":"video","title":"oEmbed Title",
			"html":"<iframe src=\"https://oembed.example.com/player\" width=\"800\" height=\"450\"></iframe>",
			"width":800,"height":450}`))
	})

	return server
}

func TestMetadata_EmbedPrecedence_OpenGraphBeatsOEmbed(t *testing.T) {

	// RULE: precedence is call order — AS → OG → Twitter → oEmbed → HTML.
	// OG runs before oEmbed, so an og:video player WINS. (Under the old
	// firstValidEmbed exception, oEmbed's player would have won instead.)
	var calls atomic.Int32
	server := embedRaceServer(t, `
		<meta property="og:video" content="https://og.example.com/player">
		<meta property="og:video:width" content="640">
		<meta property="og:video:height" content="360">`, &calls)

	m := fetchMetadata(t, server.URL+"/page")

	require.NotNil(t, m.Embed)
	require.Equal(t, EmbedIframe, m.Embed.Mode)
	require.Equal(t, "https://og.example.com/player", m.Embed.IframeURL)
	require.Equal(t, 640, m.Embed.Width) // OG's whole group, not oEmbed's 800
}

func TestMetadata_EmbedPrecedence_TwitterBeatsOEmbed(t *testing.T) {

	// Twitter also runs before oEmbed.
	var calls atomic.Int32
	server := embedRaceServer(t, `
		<meta name="twitter:player" content="https://tw.example.com/player">
		<meta name="twitter:player:width" content="480">
		<meta name="twitter:player:height" content="270">`, &calls)

	m := fetchMetadata(t, server.URL+"/page")

	require.NotNil(t, m.Embed)
	require.Equal(t, "https://tw.example.com/player", m.Embed.IframeURL)
	require.Equal(t, 480, m.Embed.Width)
}

func TestMetadata_EmbedPrecedence_OpenGraphBeatsTwitter(t *testing.T) {

	// Both in-page sources present: OG runs first, so OG wins whole.
	server := serveHTML(`<html><head>
		<meta property="og:title" content="Race">
		<meta property="og:video" content="https://og.example.com/player">
		<meta name="twitter:player" content="https://tw.example.com/player">
		<meta name="twitter:player:width" content="480">
		<meta name="twitter:player:height" content="270">
	</head><body></body></html>`)
	defer server.Close()

	m := fetchMetadata(t, server.URL)

	require.NotNil(t, m.Embed)
	require.Equal(t, "https://og.example.com/player", m.Embed.IframeURL)
	require.Zero(t, m.Embed.Width) // OG declared no dimensions; Twitter's 480 must not leak in
}

func TestMetadata_EmbedPrecedence_OEmbedWinsWhenNothingElseHasOne(t *testing.T) {

	// With no in-page embed, oEmbed supplies it — precedence order, not an exception.
	var calls atomic.Int32
	server := embedRaceServer(t, `<meta property="og:type" content="video.other">`, &calls)

	m := fetchMetadata(t, server.URL+"/page")

	require.Equal(t, int32(1), calls.Load())
	require.NotNil(t, m.Embed)
	require.Equal(t, "https://oembed.example.com/player", m.Embed.IframeURL)
	require.Equal(t, 800, m.Embed.Width)
}

func TestMetadata_EmbedPrecedence_OpenGraphFileBecomesStreamAndStillWins(t *testing.T) {

	// An og:video pointing at a media FILE becomes a STREAM embed, which
	// passes the floor — so it wins, and oEmbed's player does not replace it.
	var calls atomic.Int32
	server := embedRaceServer(t, `<meta property="og:video" content="https://og.example.com/movie.mp4">`, &calls)

	m := fetchMetadata(t, server.URL+"/page")

	require.NotNil(t, m.Embed)
	require.Equal(t, EmbedStream, m.Embed.Mode)
	require.Equal(t, "https://og.example.com/movie.mp4", m.Embed.StreamURL)
}

func TestMetadata_OEmbedNotFetchedForAnEmbedWeAlreadyHave(t *testing.T) {

	// The gate's embed half tests result.Embed: once OG has supplied a valid
	// embed, oEmbed can no longer win it, so fetching for that reason is waste.
	// The metadata half still fires here (no provider/author), so assert on the
	// embed rather than the call count.
	var calls atomic.Int32
	server := embedRaceServer(t, `
		<meta property="og:video" content="https://og.example.com/player">
		<meta property="og:image" content="https://og.example.com/hero.jpg">
		<meta property="og:site_name" content="OG Site">
		<meta property="og:author" content="Jane Doe">`, &calls)

	m := fetchMetadata(t, server.URL+"/page")

	require.Equal(t, int32(0), calls.Load()) // both halves of the gate held
	require.NotNil(t, m.Embed)
	require.Equal(t, "https://og.example.com/player", m.Embed.IframeURL)
	require.Equal(t, "Race", m.Title) // not "oEmbed Title"
}
