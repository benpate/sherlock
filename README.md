# Sherlock

<img alt="AI Generated Sherlock Holmes" src="https://github.com/benpate/sherlock/raw/main/meta/The_Adventure_of_Silver_Blaze.jpg" style="width:100%; display:block; margin-bottom:20px;">

Sherlock is a Go library that inspects a URL for any and all available metadata, pulling from whatever metadata formats are available, and returning it as an [ActivityStreams 2.0](https://www.w3.org/TR/activitystreams-core/) document.

The goal is to have a standard interface into all web content, regardless of competing data standards.

### Supported Formats

✅ [ActivityPub](https://www.w3.org/TR/activitypub/)/[ActivityStreams](https://www.w3.org/TR/activitystreams-core/)

✅ [Microformats](https://microformats.org)

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

// If you only have a URL, then pass it in to .Load()
result, err := client.Load("https://my-url-here")

// If you have already downloaded a file, then pass it to .Parse()
result, err := sherlock.ParseHTML("https://original-url", &bytes.Buffer)

```

### Using Sherlock with Hannibal

Sherlock can also be used as an http client for [Hannibal](https://github.com/benpate/hannibal), the ActivityPub library for Go.  This allows many other online resources to *look like* they're ActivityPub-enabled.
