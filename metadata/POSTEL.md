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

## Open — genuine tolerances, unattributed

Each needs: a real peer identified, a capture trimmed into `metadata/testdata/`, a test exercising it against that fixture, and a comment at the tolerance naming who sends it. **Some of these will not survive the search** — if no real peer emits the variant, the honest outcome is to delete the tolerance, not document it. The three marked ⚠ are the ones I most expect to fail to find a peer.

| # | Tolerance | Site | Suspected peer |
|---|---|---|---|
| 1 | `og:locale` as `en_US` (underscore) | `normalize.go` | Yoast SEO / WordPress |
| 2 | `parseTime` RFC 1123 layout | `normalize.go` | ⚠ RSS-derived pages? |
| 3 | `parseTime` RFC 1123Z layout | `normalize.go` | ⚠ as above |
| 4 | `parseTime` RFC 3339 without offset | `normalize.go` | unknown |
| 5 | `parseTime` bare `2006-01-02` date | `normalize.go` | unknown |
| 6 | literal `"null"` / `"undefined"` in URL fields | `normalize.go` | ⚠ broken templating; none captured |
| 7 | `jsonInt` accepting a quoted number | `extract-activitystream.go` | unknown AS2 implementation |
| 8 | `og:author` (non-standard) | `extract-opengraph.go` | unknown |
| 9 | `og:author:username` (non-standard) | `extract-opengraph.go` | ⚠ suspected fictional |
| 10 | `iconSizesAsInt` space-separated `sizes` lists | `utils.go` | HTML spec allows it — may reclassify |
| 11 | Title separators ` \| `, ` — `, ` – ` (and the deliberate exclusion of a bare hyphen) | `extract-html.go` | needs a corpus, not one peer |
| 12 | `apple-touch-icon`, `apple-touch-icon-precomposed` | `extract-html.go` | Apple vendor convention — cite, don't capture |
| 13 | `defaultHTTPS` on a bare host | `utils.go` | caller convenience, not a peer quirk — may reclassify |

Estimated 8–18 hours, plus ~1 hour to stand up a `testdata/` loader mirroring `oembed/fixtures_test.go`. Requires live network access to inspect real sites.
