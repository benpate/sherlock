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

## Image Credit

The banner is an illustration by Sidney Paget for *The Adventure of Silver Blaze* (Arthur Conan Doyle, 1892). The work is in the public domain.

## Pull Requests Welcome

I'm trying to make Sherlock the best it can be, and your help is greatly appreciated. If you find a bug or have an idea for a new feature, please open an issue or submit a pull request. We're all in this together! 🔍
