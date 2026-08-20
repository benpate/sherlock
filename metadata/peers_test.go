package metadata

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

/******************************************
 * Real-Peer Tolerances
 *
 * Every tolerance in this package must be justified by a page a real site
 * actually serves, captured in testdata/ and named here. An unattributed
 * tolerance is untested surface that silently accepts malformed input — see
 * POSTEL.md for the audit and for what is still open.
 ******************************************/

// loadPeer reads a captured page from testdata.
func loadPeer(t *testing.T, name string) *document {
	t.Helper()

	body, err := os.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)

	root, err := html.Parse(bytes.NewReader(body))
	require.NoError(t, err)

	// Built directly, not via newDocument: these fixtures are captured bytes,
	// and the extractors under test read only Root and FinalURL.
	return &document{
		RequestURL: "https://peer.example/page",
		FinalURL:   "https://peer.example/page",
		Body:       body,
		Root:       root,
	}
}

func TestPeer_SmashingMagazine_GoTimeStringDates(t *testing.T) {

	// Smashing Magazine runs Hugo and prints a Go time.Time with %v instead of
	// formatting it, so article:published_time arrives as Go's String() layout:
	// "2026-08-18 06:30:00 +0000 UTC". Before this layout was added, every
	// Smashing publication date was silently dropped.
	og := extractOpenGraph(loadPeer(t, "smashing-magazine.html"))

	require.True(t, og.PublishedAt.IsPresent())
	require.Equal(t, 2026, og.PublishedAt.Object().Year())
	require.Equal(t, 6, og.PublishedAt.Object().Hour())

	// Their modified_time carries a LEADING SPACE, which TrimSpace absorbs.
	require.True(t, og.ModifiedAt.IsPresent())

	// og:locale arrives underscored, as it does from every Yoast/Hugo-style
	// generator in the survey (Ars Technica, CSS-Tricks, MDN, TechCrunch,
	// Variety, Eleventy, Jekyll all send en_US).
	require.Equal(t, "en-US", og.Language.String())
}

func TestPeer_SmashingMagazine_ZeroTimeIsNotADate(t *testing.T) {

	// RULE: the same generator emits the Go ZERO time on pages with no date —
	// Smashing's 404 page serves "0001-01-01 00:00:00 +0000 UTC" verbatim.
	// It parses cleanly against the layout above, so parseTime must reject it
	// explicitly or year 1 enters the model as a real publication date.
	require.True(t, parseTime("0001-01-01 00:00:00 +0000 UTC").IsNull())
	require.True(t, parseTime(" 0001-01-01 00:00:00 +0000 UTC").IsNull())

	// A real timestamp in the same layout still parses.
	require.True(t, parseTime("2026-08-18 06:30:00 +0000 UTC").IsPresent())
}

func TestPeer_Flickr_TwitterTagsViaPropertyAttribute(t *testing.T) {

	// The Twitter Cards reference specifies name=, and every other site in the
	// survey uses it — but Flickr sends all 17 of its twitter:* tags with
	// property=, and Mastodon does the same. Reading only name= would drop
	// Flickr's player entirely.
	tw := extractTwitter(loadPeer(t, "flickr.html"))

	require.True(t, tw.Embed.IsPresent())
	require.Equal(t, EmbedIframe, tw.Embed.Object().Mode)
	require.Equal(t, 640, tw.Embed.Object().Width)
	require.Equal(t, 480, tw.Embed.Object().Height)
}

func TestPeer_Flickr_IconQuirks(t *testing.T) {

	// Flickr is the survey's only source for two icon quirks: a space-separated
	// sizes list ("16x16 32x32"), and rel="apple-touch-icon-precomposed".
	require.Equal(t, 32, iconSizesAsInt("16x16 32x32"))

	html := extractHTML(loadPeer(t, "flickr.html"))
	require.True(t, html.Provider.IsPresent())
	require.NotEmpty(t, html.Provider.Object().IconURL)
}

func TestPeer_ArsTechnica_BareHyphenTitlesAreAmbiguous(t *testing.T) {

	// RULE: " - " is NOT a title separator. Ars Technica proves why it cannot
	// be decided by the delimiter alone — the same site uses it both ways:
	//
	//   article:  "Satellite operators are in panic mode ... - Ars Technica"  (tail = site name)
	//   homepage: "Ars Technica - Serving the Technologist since 1998. ..."   (tail = tagline)
	//
	// Stripping at " - " would be right for the first and lossy for the second,
	// so neither is stripped. The separators that ARE stripped were each seen
	// carrying a site name: " | " (Flickr, MDN, TechCrunch), " — " (Smashing).
	const article = "Satellite operators are in panic mode due to a worsening launch crisis - Ars Technica"
	require.Equal(t, article, stripTitleSuffix(article))

	require.Equal(t, "HTML: HyperText Markup Language", stripTitleSuffix("HTML: HyperText Markup Language | MDN"))
	require.Equal(t, "Best Practices For Color Contrast", stripTitleSuffix("Best Practices For Color Contrast — Smashing Magazine"))
}
