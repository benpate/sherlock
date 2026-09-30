# 🔍 Sherlock

<img alt="Illustration of Sherlock Holmes and Watson in a train car, by Sidney Paget. From Arthur Conan Doyle's 1892 book 'The Adventure of Silver Blaze'" src="https://github.com/benpate/sherlock/raw/main/meta/The_Adventure_of_Silver_Blaze.jpg" style="width:100%; display:block; margin-bottom:20px;">

[![Go Reference](https://pkg.go.dev/badge/github.com/benpate/sherlock.svg)](https://pkg.go.dev/github.com/benpate/sherlock)
[![Version](https://img.shields.io/github/v/release/benpate/sherlock?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/benpate/sherlock/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/benpate/sherlock/go.yml?branch=main)](https://github.com/benpate/sherlock/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/benpate/sherlock?style=flat-square)](https://goreportcard.com/report/github.com/benpate/sherlock)
[![Codecov](https://img.shields.io/codecov/c/github/benpate/sherlock.svg?style=flat-square)](https://codecov.io/gh/benpate/sherlock)

## Relentless Metadata Inspector

Sherlock is a Go library that inspects a URL for any and all available metadata, pulling from whatever formats are available, and returning a single merged answer.

The goal is to have a standard interface into all web content, regardless of competing data standards.

The work happens in the sub-packages, and they answer two different questions. [metadata](metadata/) asks *"what is this page?"* and returns an oEmbed-adjacent `Preview`. The middleware stack asks *"what ActivityStreams document is at this URL?"* and answers as a Hannibal `streams.Client`. This root package holds only what both sides share.

### Supported Formats

The metadata engine fetches once, parses once, and merges every in-page format it finds by a fixed precedence.

✅ [ActivityPub](https://www.w3.org/TR/activitypub/) / [ActivityStreams](https://www.w3.org/TR/activitystreams-core/)

✅ [Open Graph](https://ogp.me)

✅ [Twitter Cards](https://developer.x.com/en/docs/x-for-websites/cards/overview/abouts-cards)

✅ [oEmbed](https://oembed.com) — lazily, only when the page suggests it could help

✅ Native HTML signals — `<title>`, `<meta name="description">`, `rel=canonical`, `rel=alternate`

### Middleware and Rewriters

The ActivityStreams half of Sherlock is a stack of client middlewares. Most just resolve a format, but a few rewrite an identifier into something the rest of the stack can resolve, or substitute a placeholder for a missing document. Each is its own subpackage with its own README and placement rules.

- [activitypub](activitypub/) — loads ActivityPub documents directly from their canonical URLs. This is the bottom of the stack, and the only layer that fetches.

- [webfinger](webfinger/README.md) — recognizes `@user@host` handles and resolves them to a canonical ActivityPub URL for the rest of the stack to load. The rewriters below all hand off to it.

- [bridgyfed](bridgyfed/README.md) — rewrites a Bluesky-looking handle (`alice.bsky.social`) into a WebFinger handle on [Bridgy Fed](https://fed.brid.gy), so Bluesky accounts resolve as ActivityPub actors.

- [tagspub](tagspub/README.md) — rewrites a `#hashtag` into a WebFinger handle on [tags.pub](https://tags.pub), so hashtags resolve to ActivityPub collections.

- [tombstone](tombstone/README.md) — turns a "Gone" (HTTP 410) response into a synthetic ActivityStreams Tombstone, so deleted objects resolve to a stable placeholder instead of an error.

### Inspecting a webpage

`metadata.Get` fetches a URL once and extracts every in-page format from a single parsed tree, merged by a fixed precedence (ActivityStreams → Open Graph → Twitter Cards → oEmbed → HTML natives) into an oEmbed-adjacent `metadata.Preview` — title, description, canonical URL, thumbnail, embed, provider, authors, dates.

```go
preview, err := metadata.Get(ctx, "https://example.com/article",
    metadata.WithUserAgent("my-app/1.0"),
)
if err != nil {
    return err
}

fmt.Println(preview.Title, preview.Description)

// Groups are filled whole from one source, or not at all -- check before use.
if preview.Thumbnail != nil {
    fmt.Println(preview.Thumbnail.URL, preview.Thumbnail.Alt)
}
```

Options carry the fetch policy. Private and loopback addresses are refused by default as an SSRF guard, and response bodies are capped:

```go
preview, err := metadata.Get(ctx, "https://my-lan-host.local/article",
    metadata.WithAllowPrivateIPs(true),
    metadata.WithMaxBodySize(2 * 1024 * 1024),
)
```

The oEmbed endpoint (the only extra network hop) is called lazily — only when a registry match or player tags suggest an embed, or when fields oEmbed could fill are still empty. Provider embed HTML never enters the model raw: it is classified by [oembed](https://github.com/benpate/oembed)'s `Embed()` into a rebuilt iframe, sandboxed markup, or nothing. Sandboxed markup is off by default; opt in with `WithAllowSandboxEmbeds(true)` if your renderer places `Embed.HTML` inside a sandboxed iframe.

The engine never imports ActivityStreams types — the package boundary enforces the layering.

### Loading an ActivityStreams document

Each middleware implements the [Hannibal](https://github.com/benpate/hannibal) `streams.Client` interface, so a stack of them can serve as the HTTP client for that ActivityPub library. Order matters, and each subpackage's README states its own placement rule.

```go
// Innermost first: each layer wraps the one below it.
client := tagspub.New(
    bridgyfed.New(
        webfinger.New(
            tombstone.New(
                activitypub.New(
                    activitypub.WithUserAgent("my-app/1.0"),
                ),
            ),
        ),
    ),
)

document, err := client.Load("@alice@example.social")
```

Per-call options travel down the stack as `...any`, and the bottom layer keeps the `remote.Option` values it recognizes. Anything else is ignored.

### Signing outbound requests

`AuthorizedFetch` is a `remote.Option` that signs a request according to the ActivityPub [Authorized Fetch](https://funfedi.dev/testing_tools/http_signatures/) convention. The `activitypub` middleware applies it for you when given a key pair; pass it directly for any other fetch that needs a signature.

```go
txn := remote.Get(url).
    With(sherlock.AuthorizedFetch(publicKeyID, privateKey))
```

## Image Credit

The banner is an illustration by Sidney Paget for *The Adventure of Silver Blaze* (Arthur Conan Doyle, 1892). The work is in the public domain.

## No Warranty

This software is provided as-is, without any warranty of any kind. Use it at your own risk.

## Pull Requests Welcome

I'm trying to make Sherlock the best it can be, and your help is greatly appreciated. If you find a bug or have an idea for a new feature, please open an issue or submit a pull request. We're all in this together! 🔍
