# metadata

Agent notes for sherlock's extraction engine. See [README.md](README.md) for what the package does and [doc.go](doc.go) for the API-level overview.

## Extractors are pure and never write into a Preview

Each `extract*` function reads the one parsed tree (or the AS2 map) and returns a sparse `partial`. `merge` is the only writer, and it never overwrites a field that an earlier source already supplied. Adding a source means adding an extractor and a `merge` call in the right position — not reaching into the `Preview`.

## Call order IS the precedence order

`merge` copies only what the destination still lacks, so the sequence of `merge` calls in `extract` is the entire definition of precedence. There is exactly one exception, and it is explicit: `result.URL` is picked by `firstString(as, html, og)` after merging. Embed has no exception — do not add one.

## Ask IsPresent, not IsZero, when merging

A source that declared an empty value HAS spoken, and `null.String.IsZero()` reads TRUE for exactly that case. Testing `IsZero` in `merge` would silently demote a source that did answer. `firstString` is the deliberate opposite: it wants a *usable* value, so it skips present-but-empty.

## Every remote dimension goes through boundDimension

`boundDimension` is the single choke point that keeps the model's invariant — dimensions are sane or zero. The oEmbed path is the one that has already been missed once: `classifyEmbed` gets its width and height straight from provider JSON, and the `EmbedPolicy` it passes sets no clamp, so it must bound them itself. `jsonInt` bounds the float *before* converting to int64, because an out-of-range float-to-int conversion is implementation-defined in Go.

## Never import hannibal or ActivityStreams types here

The dependency points from `sherlock` into this package and never back. `extractActivityStream` reads the parsed JSON map directly, because hannibal's `streams.Document` getters are live and can fetch — extractors must stay pure over bytes already in hand. `isActivityStream` and `iconSizesAsInt` are duplicated from the parent package for the same reason; that duplication is deliberate.

## capBody reads one byte past the cap

The extra byte is how `newDocument` distinguishes "exactly at the cap" from "truncated"; it trims the probe byte back off. Truncating (rather than remote's own `MaxResponseSize`, which errors) is what keeps the `<head>` of an oversized page extractable. `discoverAlternate` deliberately uses `MaxResponseSize` instead, because a truncated AS2 document cannot be parsed at all.

## og:image and og:video sub-properties are positional

`:width`, `:height`, and `:alt` bind to the most recent `og:image` / `og:video` tag, so the meta tags must be read in document order. The accumulators in `extractOpenGraph` are pointers precisely so the sub-property cases can fill them in place; they become values on the partial at the end.

## Groups reach the Preview as copies

`merge` dereferences each group into a fresh value before taking its address. That is what makes the icon backfill in `extract` safe — it writes into the `Preview`'s own `Provider`, and can never reach back through a shared pointer into the partial that won.

## decodeOnce is called at the extraction boundary and nowhere else

Double-decoding HTML entities is a real and recurring bug. Extractors decode once as they read a value; nothing downstream decodes again.

## AllowSandboxEmbeds is a product decision, not a trust ranking

Note the containment ordering, which is the opposite of most intuitions: the sandbox plan runs in an opaque-origin iframe with no cookies, storage, or parent DOM, while the plain iframe plan loads the provider's own origin unsandboxed. Setting `AllowSandboxEmbeds` FALSE blocks the *more* contained of the two plans.

## The oEmbed gate tests result.Embed, not just the signals

Precedence is call order, so an embed already merged from Open Graph or Twitter has won. Fetching oEmbed to improve on it would spend a network hop for nothing.

## revive's unused-parameter findings in tests are expected

The `w`/`r` parameters of `http.HandlerFunc` and the `t` of a fuzz callback are fixed by their signatures. The whole repo carries these; don't rename them to `_` to quiet the linter.
