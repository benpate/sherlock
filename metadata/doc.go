// Package metadata fetches a web page once and extracts an oEmbed-adjacent
// Preview from every source the page carries: ActivityStreams, Open Graph,
// Twitter Cards, oEmbed, and HTML natives.
//
// # Fetching
//
// One request negotiates ActivityStreams and HTML together. The response is
// byte-capped (truncated, not rejected, so the head of a huge page survives),
// charset-decoded to UTF-8, and parsed exactly once. Every in-page extractor
// reads that single tree.
//
// # Precedence
//
// Each source returns a sparse partial, and those partials merge in a fixed
// order: ActivityStreams, Open Graph, Twitter Cards, oEmbed, HTML natives.
// Call order IS the precedence: merge never overwrites a field that an
// earlier source already supplied. One field breaks that order on purpose —
// the canonical URL prefers rel=canonical over og:url, and requires the same
// origin.
//
// Correlated fields — thumbnail, embed, provider, authors — fill as whole
// groups from one source, past a validity floor. One source's image URL is
// never paired with another source's dimensions.
//
// The oEmbed endpoint is the only extra network hop, and it is gated: called
// only when an embed is plausible and not already found, or when a field
// oEmbed can carry is still empty.
//
// # Boundaries
//
// This package is sherlock's extraction engine and deliberately knows nothing
// about ActivityStreams result types. A Preview is never AS2; the
// Preview-to-ActivityStreams projection lives in the parent package, which
// imports this one and never the reverse. Callers who need AS2 use
// sherlock.Client.Load instead.
package metadata
