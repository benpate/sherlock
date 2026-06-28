# Sherlock

<img alt="Illustration of Sherlock Holmes and Watson in a train car, by Sidney Paget. From Arthur Conan Doyle's 1892 book 'The Adventure of Silver Blaze'" src="https://github.com/benpate/sherlock/raw/main/meta/The_Adventure_of_Silver_Blaze.jpg" style="width:100%; display:block; margin-bottom:20px;">

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/sherlock.svg)](https://pkg.go.dev/github.com/benpate/sherlock)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/sherlock/go.yml?branch=main)](https://github.com/benpate/sherlock/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/sherlock?style=flat-square)](https://goreportcard.com/report/github.com/benpate/sherlock)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/sherlock.svg?style=flat-square)](https://codecov.io/gh/benpate/sherlock)

## Relentless Metadata Inspector

Sherlock is a Go library that inspects a URL for any and all available metadata, pulling from whatever metadata formats are available, and returning it as an [ActivityStreams 2.0](https://www.w3.org/TR/activitystreams-core/) document.

The goal is to have a standard interface into all web content, regardless of competing data standards.

### Supported Formats

✅ [ActivityPub](https://www.w3.org/TR/activitypub/)/[ActivityStreams](https://www.w3.org/TR/activitystreams-core/)

✅ [MicroFormats](https://microformats.org)

✅ [Open Graph](https://ogp.me)

### In Progress

🚧 [WebFinger](https://webfinger.net)

🚧 [JSON-LD (Linked)](https://json-ld.org/)

🚧 [Twitter Metadata](https://developer.twitter.com/en/docs/twitter-for-websites/cards/overview/abouts-cards)

🚧 [Microdata](https://html.spec.whatwg.org/multipage/microdata.html#microdata)

🚧 [RDFa](https://rdfa.info)

🚧 [oEmbed data provider](https://oembed.com)


### Using Sherlock

```go
client := sherlock.NewClient()

// Load inspects a URL and returns whatever metadata it can find,
// as an ActivityStreams document.
result, err := client.Load("https://my-url-here")

// Per-call options refine the request: what kind of object you expect,
// a default value to merge into, a redirect cap, or extra remote options.
result, err = client.Load("https://my-url-here",
    sherlock.AsActor(),
    sherlock.WithMaximumRedirects(4),
)
```

### Using Sherlock with Hannibal

Sherlock implements the [Hannibal](https://github.com/benpate/hannibal) `streams.Client` interface, so it can be used as the HTTP client for that ActivityPub library. This makes many non-ActivityPub resources *look like* they're ActivityPub-enabled.

## What matters here

- **`Load` never hard-fails on a single source; it tries many and merges.** Each metadata format is attempted in turn and merged into one ActivityStreams document without overwriting values already found. A fetch or parse error from one source is swallowed so the others still run. Don't expect an error just because one format was absent or malformed.

- **The subpackages are stacked client middlewares, and stacking ORDER is load-bearing.** `bridgyfed` and `tagspub` rewrite identifiers into WebFinger handles, so they MUST sit *above* `webfinger` in the stack; `webfinger` resolves handles to URLs that `activitypub` then loads; `tombstone` substitutes a placeholder for Gone documents. Each subpackage's README states its own placement rule.

- **Network access is SSRF-hardened by default, inherited from [remote](https://github.com/benpate/remote).** Sherlock sets no `AllowPrivateIPs`, so private/loopback fetches are blocked and response sizes are capped. Self-hosted/LAN targets will be refused unless the caller passes a remote option to allow them.

- **Identifier classification is strict and lives in [uri](https://github.com/benpate/uri).** Whether a value "looks like" a URL or an `@handle` is decided by `uri` validation (real IANA TLDs, 2+ segments). When a test for that classification fails after a `uri` upgrade, the test assumption is usually what drifted — not the code.

- **Untrusted-input parsers are fuzzed; keep new ones that way.** The document and identifier parsers (OpenGraph, microformats, embedded JSON-LD, address/identifier classification) have `Fuzz*` coverage in `fuzz_test.go`. Regexes are static patterns (`regexp.MustCompile`) — no untrusted input is ever compiled into a regex. A new parser of remote bytes should arrive with its own fuzz target.
