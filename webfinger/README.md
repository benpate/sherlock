# webfinger

A [Sherlock](../README.md) client middleware that recognizes WebFinger handles (`@user@host`) and resolves them to a canonical ActivityPub URL before handing off to the inner client.

## What matters here

- **It decides "not WebFinger" without a network call wherever possible.** `isWebfinger` short-circuits on `http(s)://` URLs, values with no `@`, and values that `digit.ParseAccount` can't parse — all before any lookup. Only a plausible handle triggers the actual `digit.Lookup`. This ordering keeps the common (URL) path free of network cost.

- **Resolution picks the `self` link with an ActivityPub media type.** From the WebFinger response it returns the first link whose RelationType is `digit.RelationTypeSelf` AND whose MediaType is ActivityPub. No such link means "not WebFinger" (forward the original), not an error.

- **`New` takes `...remote.Option`, not a `ClientOption`.** Unlike the other middlewares, this package has no `ClientOption` type — remote options (timeouts, Authorized Fetch, the SSRF-guarding transport) flow straight through to `digit.Lookup`.

- **The network test is gated behind `//go:build localonly`.** Offline behavior (the early-return branches and inner-client forwarding) is covered by the un-tagged tests; the live `digit.Lookup` path runs only with `-tags localonly`, which is why CI coverage excludes it.
