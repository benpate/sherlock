# metadata

Sherlock's extraction engine: fetch a web page once and return everything it says about itself as a `Preview` — title, description, canonical URL, thumbnail, embed, provider, authors, dates, and a semantic `Kind`.

The shape is oEmbed's field names on an ActivityStreams-shaped object — the same hybrid as Mastodon's `PreviewCard`. It is a rebuildable value, never persisted or federated, and it is not ActivityStreams: callers who need AS2 use `sherlock.Client.Load` instead.

## Usage

```go
preview, err := metadata.Get(ctx, "https://example.com/article",
    metadata.WithUserAgent("my-app/1.0"),
)
if err != nil {
    return err
}

fmt.Println(preview.Title, preview.Kind, preview.Thumbnail.URL)
```

`sherlock.Client.Metadata` is a thin wrapper that threads the Client's configuration (User-Agent, Authorized Fetch signatures, SSRF policy, body cap) into `Get`. Use `Get` directly when you have no Client.

## How it works

One request negotiates ActivityStreams and HTML together (`Accept: application/activity+json, text/html;q=0.9`). The response is byte-capped (truncated, not rejected, so the `<head>` of a huge page survives), charset-decoded to UTF-8, and parsed exactly once. Every in-page extractor reads that one tree and returns a sparse partial; the engine merges them in a fixed precedence order:

**ActivityStreams → Open Graph → Twitter Cards → oEmbed → HTML natives**

One field breaks the global order on purpose: the canonical URL prefers `rel=canonical` over `og:url` (same-origin required). Everything else — the embed included — is decided by the merge order above and nothing else. Correlated fields (thumbnail, embed, provider, authors) fill as whole groups from one source, past a validity floor — one source's image URL is never paired with another's dimensions.

The oEmbed endpoint is the only extra network hop, and it is gated: called only when a registry match or player tags suggest an embed, or when a field oEmbed can carry (title, thumbnail, provider, author) is still empty. In practice the second half of that gate is generous — the author field is empty on most pages, since `og:author` is non-standard and a bare `article:author` URL fails the author floor — so a page that advertises an endpoint usually gets called.

## Options

| Option | Default | Purpose |
|---|---|---|
| `WithUserAgent` | empty | User-Agent sent on every request |
| `WithRemoteOptions` | none | extra `remote.Option`s (signatures, instrumentation) |
| `WithAllowPrivateIPs` | `false` | SSRF guard — set TRUE only for httptest servers or intentional LAN targets |
| `WithAllowSandboxEmbeds` | `false` | let non-extractable oEmbed HTML become an `EmbedProviderHTML` embed for a sandboxed iframe |
| `WithMaxBodySize` | 1 MB | response body cap; oversize bodies are truncated |

## Rendering embeds

`Preview.Embed` is never raw provider markup. Its `Mode` tells you how to render it: `EmbedIframe` (build your own iframe from `IframeURL`), `EmbedStream` (a direct media file at `StreamURL`), or `EmbedProviderHTML` (only when `WithAllowSandboxEmbeds(true)`; place `HTML` inside an iframe with `sandbox="allow-scripts allow-popups"` — never `allow-same-origin`). `AllowSandboxEmbeds` is a product decision applied to every provider uniformly, not a per-provider trust ranking.

## Design notes

The full design, decisions, and precedence rationale live in the Emissary spec `LINK-METADATA-CONSUMER.md`. The rules an editor of this package needs — the ones that are easy to break and hard to notice — are collected in [AGENTS.md](AGENTS.md).
