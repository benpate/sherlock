# Postel's Law Audit

Strict in what you send, liberal in what you accept — but every tolerance must be justified by a real peer, captured as a fixture and named in a comment. An unattributed tolerance is worse than an undocumented one: it is untested surface that silently accepts malformed input, and the next reader "fixes" it and reopens the bug it was hiding.

This file tracks that audit for the `metadata` package. The model to copy is [`benpate/oembed`](../../oembed), which already did this: named-provider fixtures in `testdata/` (`tiktok.json`, `soundcloud.json`, `reddit.json`, …) and attribution comments at each tolerance — *"TikTok sends width and height as the string `100%`"*, *"SoundCloud sends `version: 1.0` as a JSON NUMBER"*, *"Reddit omits width entirely"*.

The oEmbed leg of this package's tolerance surface is therefore **already compliant upstream** — `metadata` delegates oEmbed response parsing to that library. What remains is the HTML, Open Graph, Twitter Cards, and ActivityStreams extractors.

## Done — spec-mandated, cited in place (2026-08-18)

These are not tolerances. Each accepts a variant because a specification says the variant is conforming, so a citation settles it and no fixture is needed.

| Behavior | Authority | Site |
|---|---|---|
| `og:image` / `og:image:url` / `og:image:secure_url` | OGP "Structured Properties" | `extract-opengraph.go` |
| `og:video` / `og:video:url` / `og:video:secure_url` | OGP "Structured Properties" | `extract-opengraph.go` |
| `property=` and `name=` both read | OGP is RDFa (`property`); Twitter Cards uses `name` | `extract-opengraph.go` |
| AS2 value may be string, object, or array | ActivityStreams 2.0 Core §2 | `extract-activitystream.go` |
| `preferredUsername` as a name fallback | ActivityPub §4.1 | `extract-activitystream.go` |
| `activity+json` and `ld+json` both accepted | ActivityPub §3.2 | `utils.go` |
| Content-Type → BOM → `<meta>` charset order | HTML Standard, encoding sniffing algorithm | `document.go` |

## Done — reclassified, exempt from this audit (2026-08-18)

Bounds and security decisions are **not** subject to Postel's liberality. They are labelled in place so a future pass cannot mistake them for tolerances to widen.

| Behavior | Why exempt | Site |
|---|---|---|
| `boundDimension` | BOUND — never widened for a peer | `normalize.go` |
| `parseInt` overflow bail | BOUND | `extract-opengraph.go` |
| `config.maxBodySize()` | BOUND — oversize is truncated, not accommodated | `config.go` |
| `normalizeCanonicalURL` same-origin | SECURITY — decides identity; permanently fenced | `normalize.go` |
| `normalizeURL` rejecting `"null"`/`"undefined"` | STRICTNESS — rejects rather than accepts (but see below) | `normalize.go` |

## Hunt results — 47 real pages surveyed (2026-08-18)

Corpus: 19 major sites (news, video, code hosting, reference), 8 article pages from those sites, 20 CMS/static-site-generator sites (Substack, Ghost, Medium, Blogger, Wix, Squarespace, Shopify, Tumblr, Webflow, Drupal, Joomla, Hugo, Jekyll, Eleventy, Gatsby, Next.js, Notion, dev.to, Hashnode, Mastodon blog), and 4 ActivityStreams documents (Mastodon actor, Lemmy, PeerTube).

### Bug found and fixed

**Smashing Magazine's publication dates were being silently dropped.** They run Hugo and print a Go `time.Time` with `%v` instead of formatting it, so `article:published_time` arrives as Go's `String()` layout — `2026-08-18 06:30:00 +0000 UTC` — which matched none of the five layouts. Added as a sixth, explicitly marked as *not* an interchange format.

The fix carried a trap: the same generator emits the Go **zero** time (`0001-01-01 00:00:00 +0000 UTC`) on pages with no date, and that parses cleanly against the new layout. `parseTime` now rejects a zero time explicitly, or year 1 would enter the model as a real publication date. Fixture: `testdata/smashing-magazine.html`. Tests: `TestPeer_SmashingMagazine_*`.

### Tolerance deleted

**`og:author` and `og:author:username`: ZERO occurrences in 47 pages.** These are folklore from blog posts about Open Graph, not real peer behavior — the OGP spec has no `og:author`. Removed.

The consequence is larger than the tag: `og:author` was the *only* in-page source of an author name, so **Open Graph now contributes no authors at all**, and the whole author-accumulator path in `extractOpenGraph` went with it. `article:author` is a bare profile URL per spec and fails the author floor alone. Authors now reach a `Preview` only from AS2 `attributedTo` or oEmbed `author_name`.

Three tests were built on this fictional tag and had to be rewritten — precisely the failure mode the audit exists to catch: the tolerance's only evidence was a fixture we wrote ourselves.

### Attributed — tolerance confirmed by a real peer

| Tolerance | Peer evidence | Fixture |
|---|---|---|
| `og:locale` underscored (`en_US`) | Ars Technica, CSS-Tricks, MDN, Smashing, TechCrunch, Variety, Eleventy, Jekyll (8 sites; Notion sends `en-US`, Mastodon blog sends `en` — both forms are live) | `smashing-magazine.html` |
| Go `time.String()` dates | Smashing Magazine (Hugo) | `smashing-magazine.html` |
| Zero-time rejection | Smashing Magazine 404 pages | — |
| `twitter:*` via `property=` | **Flickr** (all 17 tags), Mastodon | `flickr.html` |
| `sizes` space-separated list | **Flickr** (`16x16 32x32`) — sole source in 47 pages | `flickr.html` |
| `apple-touch-icon-precomposed` | **Flickr** — sole source | `flickr.html` |
| `apple-touch-icon` | 13 sites | `flickr.html` |
| Title separator ` \| ` | Flickr, MDN, TechCrunch | — |
| Title separator ` — ` | Smashing Magazine | — |
| AS2 value as string / object / array | Lemmy (string), Mastodon (object), **PeerTube** (array) | — |

### Reclassified — spec, not tolerance

The date layouts are all published interchange formats, so they need a citation rather than a peer: RFC 3339, HTTP-date (RFC 9110 §5.6.7) for the two RFC 1123 forms, and ISO 8601 for the bare date. Only the Go `String()` layout is a genuine peer accommodation.

## Still open

| Tolerance | Status after 47 pages |
|---|---|
| `jsonInt` accepting a quoted number | **Unattributed.** PeerTube sends proper ints; PeerTube's `duration` is the string `"PT113S"`, but that is an ISO 8601 duration, not a quoted int. `benpate/oembed` proved this quirk real for *oEmbed* providers, not for AS2. Candidate for deletion. |
| literal `"null"` / `"undefined"` in URL fields | **Unattributed** — zero occurrences. Kept because it *rejects* rather than accepts, so it narrows the accept surface rather than widening it, but no peer justifies it. |
| Title separator ` – ` (en dash) | **Unattributed** — zero occurrences in 47 pages. Candidate for deletion. |
| Bare-hyphen EXCLUSION | **Weakly supported.** Ars Technica uses ` - ` both ways (article: tail is the site name; homepage: tail is a tagline), so the delimiter alone cannot decide — but both would have stripped acceptably. The exclusion is defensible, not proven. |
| ` · ` as a separator | **Gap, not a tolerance:** GitHub uses it (`... · GitHub`) and we do not strip it. |
| `article:published_time` as a bare date | Unattributed. NPR sends `<meta name="date" content="2026-08-18">`, which this package does not read. |
| RFC 1123 / 1123Z in `article:published_time` | No occurrences, but now spec-cited rather than peer-attributed, so no fixture is owed. |
| `defaultHTTPS` on a bare host | Caller convenience, not a peer quirk — reclassify, don't capture. |

## Method

Captures were taken with a sherlock-like User-Agent so the corpus reflects what this library actually receives. Fixtures in `testdata/` are trimmed to the tags under test and carry the capture date. To re-run the survey, fetch a page and diff its `<head>` against the relevant fixture.
