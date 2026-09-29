# activitypub

A [Sherlock](../README.md) client middleware that loads ActivityPub/ActivityStream documents directly from their canonical URLs, falling back to an inner client when it can't.

## What matters here

- **It only handles valid URLs; everything else falls through.** `Load` checks `uri.NotValidURL` first and forwards anything that isn't a URL to the inner client (or returns NotFound when there is none). This is the gate that keeps non-URL identifiers — `@handle@host`, hashtags — flowing down to the WebFinger/bridgyfed/tagspub middlewares stacked below.

- **A host may speak only for its own ids.** A document whose `id` is on another host than the URL that finally served it is replaced by the copy that id's own host serves, fetched once, or rejected with `Forbidden`. A rejection never falls through to the inner client. See the repository's AGENTS.md.

- **A non-ActivityPub response is treated as a miss, not an error.** After a successful fetch, it checks the `Content-Type` via `hannibal.IsActivityPubContentType`; if it isn't ActivityPub, it falls through to the inner client rather than returning the HTML. This is what lets a plain web page be retried by the HTML-scraping path.

- **`Delete`/`Save`/`SetRootClient` are inner-client pass-throughs, and `Delete`/`Save` are nil-safe.** With no inner client they are silent no-ops; `SetRootClient` only records the root when an inner client exists. Don't "simplify" these into unconditional dereferences — the nil guard is important for a top-of-stack client.

- **Authorized Fetch is opt-in via `WithKeyPairFunc`.** When set, requests are HTTP-signed (the key pair is resolved lazily, per request). The SSRF protections from [remote](https://github.com/benpate/remote) still apply regardless — signing modifies the request, it does not bypass the transport-level IP guard.
