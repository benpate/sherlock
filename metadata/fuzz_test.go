package metadata

import (
	"bytes"
	"encoding/json"
	"testing"

	"golang.org/x/net/html"
)

// fuzzDocument builds a document from raw fuzzed HTML bytes, the way the
// engine would after a fetch.
func fuzzDocument(data []byte) *document {

	root, err := html.Parse(bytes.NewReader(data))

	if err != nil {
		return nil
	}

	return &document{
		RequestURL: "https://example.com/page",
		FinalURL:   "https://example.com/page",
		Body:       data,
		Root:       root,
	}
}

// FuzzExtractOpenGraph asserts the OG extractor never panics on hostile HTML.
func FuzzExtractOpenGraph(f *testing.F) {

	f.Add([]byte(`<html><head><meta property="og:title" content="T"><meta property="og:image" content="/i.jpg"><meta property="og:image:width" content="99999999999999999999"></head></html>`))
	f.Add([]byte(`<meta property="og:image:width" content="800"><meta property="og:image" content="x">`))
	f.Add([]byte(`<<<>>>&&&`))

	f.Fuzz(func(t *testing.T, data []byte) {
		if doc := fuzzDocument(data); doc != nil {
			_ = extractOpenGraph(doc)
		}
	})
}

// FuzzExtractHTML asserts the HTML extractor never panics on hostile HTML.
func FuzzExtractHTML(f *testing.F) {

	f.Add([]byte(`<html lang="en"><head><title>A | B | C</title><link rel="canonical" href="javascript:x"><link rel="icon" sizes="0x0" href=""></head></html>`))
	f.Add([]byte(`<title>`))

	f.Fuzz(func(t *testing.T, data []byte) {
		if doc := fuzzDocument(data); doc != nil {
			_ = extractHTML(doc)
		}
	})
}

// FuzzExtractTwitter asserts the Twitter extractor never panics on hostile HTML.
func FuzzExtractTwitter(f *testing.F) {

	f.Add([]byte(`<meta name="twitter:player" content="//x"><meta name="twitter:player:width" content="-1">`))

	f.Fuzz(func(t *testing.T, data []byte) {
		if doc := fuzzDocument(data); doc != nil {
			_ = extractTwitter(doc)
		}
	})
}

// FuzzExtractActivityStream asserts the AS extractor never panics on hostile
// JSON shapes — strings, arrays, and objects in every polymorphic position.
func FuzzExtractActivityStream(f *testing.F) {

	f.Add([]byte(`{"name": ["not","a","string"], "attributedTo": [{"name": {"x": 1}}], "image": [[]]}`))
	f.Add([]byte(`{"image": {"url": 42, "width": "abc"}, "attributedTo": "https://x"}`))

	f.Fuzz(func(t *testing.T, data []byte) {

		object := map[string]any{}

		if err := json.Unmarshal(data, &object); err != nil {
			return
		}

		doc := &document{
			RequestURL:     "https://example.com/object",
			FinalURL:       "https://example.com/object",
			ActivityStream: object,
		}

		_ = extractActivityStream(doc)
	})
}
