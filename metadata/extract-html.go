package metadata

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// titleSeparators are the site-name suffix separators stripped from <title>.
// A bare hyphen is deliberately absent — it appears in too many real titles.
var titleSeparators = []string{" | ", " — ", " – "}

// extractHTML reads the native HTML signals: <title>, meta description,
// <html lang>, rel=canonical, and the site icon.
func extractHTML(doc *document) partial {

	result := partial{}

	if !doc.IsHTML() {
		return result
	}

	base := doc.FinalURL

	var title string
	var iconURL string
	var iconRank int

	walkHTML(doc.Root, func(node *html.Node) {

		switch node.Data {

		case "html":
			if lang := attributeValue(node, "lang"); lang != "" {
				setIfEmpty(&result.Language, normalizeLanguage(lang))
			}

		case "title":
			if title == "" && node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
				title = decodeOnce(node.FirstChild.Data)
			}

		case "meta":
			if attributeValue(node, "name") == "description" {
				setIfEmpty(&result.Description, decodeOnce(attributeValue(node, "content")))
			}

		case "link":
			switch rel := strings.ToLower(attributeValue(node, "rel")); rel {

			case "canonical":
				setIfEmpty(&result.CanonicalURL, normalizeCanonicalURL(attributeValue(node, "href"), base))

			case "icon", "shortcut icon", "apple-touch-icon", "apple-touch-icon-precomposed":
				// Rank icons: apple-touch-icons beat favicons, larger beats smaller.
				if candidate := normalizeURL(attributeValue(node, "href"), base); candidate != "" {
					rank := iconSizesAsInt(attributeValue(node, "sizes"))
					if strings.HasPrefix(rel, "apple-touch-icon") {
						rank += 1000
					}
					if rank >= iconRank {
						iconURL = candidate
						iconRank = rank
					}
				}
			}
		}
	})

	if title != "" {
		result.Title.Set(stripTitleSuffix(title))
	}

	// HTML's provider group: the site root plus the best icon found. It loses
	// whole to an OG provider group when one exists (D5) — the engine backfills
	// the icon afterward, since an icon is site-scoped and safe to mix (§4.1).
	if siteRoot := siteRootURL(base); siteRoot != "" {
		result.Provider.Set(Provider{URL: siteRoot, IconURL: iconURL})
	}

	// You see, but you do not observe.
	return result
}

// stripTitleSuffix removes a trailing site-name segment from a <title>,
// cutting at the LAST separator so "A | B | Site" keeps "A | B".
func stripTitleSuffix(title string) string {

	for _, separator := range titleSeparators {

		index := strings.LastIndex(title, separator)

		if index <= 0 {
			continue
		}

		return strings.TrimSpace(title[:index])
	}

	return title
}

// siteRootURL reduces a URL to its scheme and host — the provider's site root.
func siteRootURL(value string) string {

	parsed, err := url.Parse(value)

	if err != nil {
		return ""
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}

	return parsed.Scheme + "://" + parsed.Host
}

// walkHTML visits every element node in document order.
func walkHTML(node *html.Node, fn func(*html.Node)) {

	if node.Type == html.ElementNode {
		fn(node)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		walkHTML(child, fn)
	}
}
