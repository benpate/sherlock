package metadata

import (
	"context"
	"strings"

	"github.com/benpate/remote"
	"golang.org/x/net/html"
)

// discoverAlternate fetches a page's <link rel="alternate"> ActivityStreams
// representation (FEP-22b6), returning an empty partial when there is none or
// the fetch fails.
func discoverAlternate(ctx context.Context, config config, doc *document) partial {

	// Servers that publish an alternate but do not content-negotiate are the
	// reason this path exists. Discovery is best-effort and never sinks the
	// whole extraction — every failure below yields an empty partial, which
	// merges as a no-op.

	if !doc.IsHTML() {
		return partial{}
	}

	// Find an AS2 alternate link in the parsed tree
	alternateURL := findAlternateLink(doc.Root, doc.FinalURL)

	if alternateURL == "" {
		return partial{}
	}

	// RULE: the alternate must not be the page itself — that would loop.
	if alternateURL == doc.FinalURL {
		return partial{}
	}

	// Fetch the alternate document under the engine's fetch policy
	activityStream := map[string]any{}

	txn := remote.Get(alternateURL).
		WithContext(ctx).
		UserAgent(config.UserAgent).
		Header("Accept", contentTypeActivityPub).
		AllowPrivateIPs(config.AllowPrivateIPs).
		MaxResponseSize(config.maxBodySize()).
		With(config.RemoteOptions...).
		Result(&activityStream)

	if err := txn.Send(); err != nil {
		return partial{}
	}

	if len(activityStream) == 0 {
		return partial{}
	}

	// Wrap the fetched object as a synthetic AS2 document so the ordinary
	// extractor does the mapping — one mapping, not two.
	alternate := &document{
		RequestURL:     alternateURL,
		FinalURL:       alternateURL,
		ActivityStream: activityStream,
	}

	// Come at once if convenient. If inconvenient, come all the same.
	return extractActivityStream(alternate)
}

// findAlternateLink returns the first <link rel="alternate"> whose type is an
// ActivityStreams MIME type, resolved against the base URL.
func findAlternateLink(root *html.Node, base string) string {

	result := ""

	walkHTML(root, func(node *html.Node) {

		if result != "" {
			return
		}

		if node.Data != "link" {
			return
		}

		if !strings.EqualFold(attributeValue(node, "rel"), "alternate") {
			return
		}

		if !isActivityStream(attributeValue(node, "type")) {
			return
		}

		result = normalizeURL(attributeValue(node, "href"), base)
	})

	return result
}
