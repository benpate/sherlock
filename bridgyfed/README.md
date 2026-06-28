# bridgyfed

A [Sherlock](../README.md) client middleware that translates Bluesky-looking handles into ActivityPub actors via [Bridgy Fed](https://fed.brid.gy), by rewriting them into a WebFinger handle on `bsky.brid.gy`.

## What matters here

- **This middleware MUST sit ABOVE the WebFinger middleware in the stack.** It rewrites `alice.bsky.social` into `@alice.bsky.social@bsky.brid.gy` and hands that to the inner client — which must be (or eventually reach) the WebFinger client for the lookup to resolve. Placed below WebFinger, the rewritten handle is never resolved.

- **A rewrite attempt always falls back to the original id.** `Load`/`Delete` first try the Bridgy handle; on *any* error they retry the inner client with the unmodified id. So a false-positive "looks like Bluesky" never blocks a value that some lower client could have handled.

- **`LooksLikeBluesky` is deliberately strict — and its strictness now comes from `uri.ValidateHostname`.** It requires a valid 2+ segment hostname with a real IANA TLD, rejects the autocomplete-in-progress `*.bsky.so` suffix, and rejects slashes/colons/multiple `@`. When editing its tests, remember the hostname rules live in the `uri` dependency: a test input must be genuinely single-segment (e.g. `nodots`) to exercise the "too few segments" rejection — `foo.net` is valid. (See the matching memory on this drift.)

- **The default hostname is `bsky.brid.gy`; override with `WithHostname`.**
