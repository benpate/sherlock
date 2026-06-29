# Sherlock — Notes for AI Agents

- **`Load` never hard-fails on a single source; it tries many and merges.** Each metadata format is attempted in turn and merged into one ActivityStreams document without overwriting values already found. A fetch or parse error from one source is swallowed so the others still run. Don't expect an error just because one format was absent or malformed.

- **The subpackages are stacked client middlewares, and stacking ORDER is load-bearing.** `bridgyfed` and `tagspub` rewrite identifiers into WebFinger handles, so they MUST sit *above* `webfinger` in the stack; `webfinger` resolves handles to URLs that `activitypub` then loads; `tombstone` substitutes a placeholder for Gone documents. Each subpackage's README states its own placement rule.

- **Network access is SSRF-hardened by default, inherited from [remote](https://github.com/benpate/remote).** Sherlock sets no `AllowPrivateIPs`, so private/loopback fetches are blocked and response sizes are capped. Self-hosted/LAN targets will be refused unless the caller passes a remote option to allow them.

- **Identifier classification is strict and lives in [uri](https://github.com/benpate/uri).** Whether a value "looks like" a URL or an `@handle` is decided by `uri` validation (real IANA TLDs, 2+ segments). When a test for that classification fails after a `uri` upgrade, the test assumption is usually what drifted — not the code.

- **Untrusted-input parsers are fuzzed; keep new ones that way.** The document and identifier parsers (OpenGraph, microformats, embedded JSON-LD, address/identifier classification) have `Fuzz*` coverage in `fuzz_test.go`. Regexes are static patterns (`regexp.MustCompile`) — no untrusted input is ever compiled into a regex. A new parser of remote bytes should arrive with its own fuzz target.
