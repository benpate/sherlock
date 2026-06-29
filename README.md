# 🔍 Sherlock

<img alt="Illustration of Sherlock Holmes and Watson in a train car, by Sidney Paget. From Arthur Conan Doyle's 1892 book 'The Adventure of Silver Blaze'" src="https://github.com/benpate/sherlock/raw/main/meta/The_Adventure_of_Silver_Blaze.jpg" style="width:100%; display:block; margin-bottom:20px;">

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/sherlock.svg)](https://pkg.go.dev/github.com/benpate/sherlock)
[![Version](https://img.shields.io/github/v/release/benpate/sherlock?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/sherlock/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/sherlock/go.yml?branch=main)](https://github.com/benpate/sherlock/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/sherlock?style=flat-square)](https://goreportcard.com/report/github.com/benpate/sherlock)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/sherlock.svg?style=flat-square)](https://codecov.io/gh/benpate/sherlock)

## Relentless Metadata Inspector

Sherlock is a Go library that inspects a URL for any and all available metadata, pulling from whatever metadata formats are available, and returning it as an [ActivityStreams 2.0](https://www.w3.org/TR/activitystreams-core/) document.

The goal is to have a standard interface into all web content, regardless of competing data standards.

### Supported Formats

Sherlock attempts every format it knows about and merges what it finds into a single ActivityStreams document.

✅ [ActivityPub](https://www.w3.org/TR/activitypub/) / [ActivityStreams](https://www.w3.org/TR/activitystreams-core/)

✅ [Open Graph](https://ogp.me)

✅ [MicroFormats2](https://microformats.org) — `h-entry` documents and `h-feed` actors

✅ [JSON-LD](https://json-ld.org/) — both embedded `<script type="application/ld+json">` and linked `<link rel="alternate">`

✅ [RSS / Atom](https://www.rssboard.org/rss-specification) feeds

✅ [JSON Feed](https://www.jsonfeed.org)

✅ HTTP `Link` headers

### Middleware and Rewriters

Sherlock is built from a stack of client middlewares. Most just resolve a format, but a few rewrite an identifier into something the rest of the stack can resolve, or substitute a placeholder for a missing document. Each is its own subpackage with its own README and placement rules.

- [webfinger](webfinger/README.md) — recognizes `@user@host` handles and resolves them to a canonical ActivityPub URL for the rest of the stack to load. The rewriters below all hand off to it.

- [bridgyfed](bridgyfed/README.md) — rewrites a Bluesky-looking handle (`alice.bsky.social`) into a WebFinger handle on [Bridgy Fed](https://fed.brid.gy), so Bluesky accounts resolve as ActivityPub actors.

- [tagspub](tagspub/README.md) — rewrites a `#hashtag` into a WebFinger handle on [tags.pub](https://tags.pub), so hashtags resolve to ActivityPub collections.

- [tombstone](tombstone/README.md) — turns a "Gone" (HTTP 410) response into a synthetic ActivityStreams Tombstone, so deleted objects resolve to a stable placeholder instead of an error.


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

## Image Credit

The banner is an illustration by Sidney Paget for *The Adventure of Silver Blaze* (Arthur Conan Doyle, 1892). The work is in the public domain.

## Pull Requests Welcome

I'm trying to make Sherlock the best it can be, and your help is greatly appreciated. If you find a bug or have an idea for a new feature, please open an issue or submit a pull request. We're all in this together! 🔍
