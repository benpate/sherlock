# htmlparser

Standalone HTML metadata parsers (currently OpenGraph) that extract ActivityStreams-style data from a fetched web page, for use by [Sherlock](../README.md)'s document loaders.

## What matters here

- **The directory is `htmlparser` but the package is `adapter`.** Import it as `adapter`; don't "fix" the package name to match the folder without updating every call site — the mismatch is intentional, not a typo.

- **Parsers merge into `data` without overwriting existing values.** Each parser checks `data.IsZeroValue(...)` before writing, so a richer source applied earlier wins. Order of application is therefore significant at the call site.

- **OpenGraph ID uses a two-tier fallback.** `PropertyID` is set from the page's `og:url` if present, otherwise from the fetch URL passed in. Keep these two assignments distinct (the inner var is named `ogURL` precisely to avoid shadowing the `url` parameter) — collapsing them loses the fallback.

- **This package is the seam for adding new scrapers.** A new metadata source (Twitter Cards, JSON-LD, etc.) is a new `ParserFunc` here, kept free of network and client concerns — it takes a body and a map, nothing more.
