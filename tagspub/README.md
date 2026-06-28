# tagspub

A [Sherlock](../README.md) client middleware that resolves `#hashtag` identifiers into ActivityPub collections via [tags.pub](https://tags.pub), by rewriting them into a WebFinger handle on that host.

## What matters here

- **This middleware MUST sit ABOVE the WebFinger middleware in the stack.** It rewrites `#golang` into the WebFinger handle `golang@tags.pub` and hands that to the inner client, which must reach the WebFinger client to resolve it. Below WebFinger, the rewritten handle is never looked up.

- **A rewrite attempt always falls back to the original id.** `Load`/`Delete` first try the `tag@tags.pub` handle; on any error they retry the inner client with the unmodified id, so a misfire never blocks a value a lower client could handle.

- **`IsHashtag` requires a leading `#` and only `[a-zA-Z0-9_]`.** Dashes, dots, slashes, and colons are rejected; the returned tag has the `#` stripped. The quick byte-0 check before the regex is just an early-out, not a second rule.

- **The default hostname is `tags.pub`; override with `WithHostname`.**
