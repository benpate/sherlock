package metadata

import (
	"strings"

	"github.com/benpate/rosetta/null"
	"golang.org/x/net/html"
)

// extractOpenGraph reads Open Graph, article:*, and fediverse:creator meta
// tags from the parsed tree into a sparse partial.
func extractOpenGraph(doc *document) partial {

	result := partial{}

	if !doc.IsHTML() {
		return result
	}

	base := doc.FinalURL

	// Accumulators for positional groups and split author fields.
	// RULE: og:image and og:video sub-properties (:width, :height, :alt) are
	// POSITIONAL — they bind to the most recent og:image / og:video tag, so
	// the tags must be read in document order.
	var thumbnail *Thumbnail
	var embed *Embed
	var embedIsFile bool
	var authorName string
	var authorURL string
	var ogType string

	forEachMetaTag(doc.Root, func(key string, content string) {

		switch key {

		// Scalars
		case "og:title":
			setIfEmpty(&result.Title, decodeOnce(content))

		case "og:description":
			setIfEmpty(&result.Description, decodeOnce(content))

		case "og:url":
			setIfEmpty(&result.CanonicalURL, normalizeCanonicalURL(content, base))

		case "og:locale":
			setIfEmpty(&result.Language, normalizeLanguage(content))

		case "og:site_name":
			if result.Provider.IsNull() {
				if name := decodeOnce(content); name != "" {
					result.Provider.Set(Provider{Name: name})
				}
			}

		case "og:type":
			if ogType == "" {
				ogType = strings.ToLower(strings.TrimSpace(content))
			}

		// Thumbnail — positional group
		case "og:image", "og:image:url", "og:image:secure_url":
			// A new og:image starts a new group; keep only the first complete one.
			if thumbnail == nil {
				if imageURL := normalizeURL(content, base); imageURL != "" {
					thumbnail = &Thumbnail{URL: imageURL}
				}
			}

		case "og:image:width":
			if thumbnail != nil && thumbnail.Width == 0 {
				thumbnail.Width = parseInt(content)
			}

		case "og:image:height":
			if thumbnail != nil && thumbnail.Height == 0 {
				thumbnail.Height = parseInt(content)
			}

		case "og:image:alt":
			if thumbnail != nil && thumbnail.Alt == "" {
				thumbnail.Alt = decodeOnce(content)
			}

		// Embed — positional group. og:video may be a FILE or a PLAYER;
		// og:video:type distinguishes when present, URL shape when not.
		case "og:video", "og:video:url", "og:video:secure_url":
			if embed == nil {
				if videoURL := normalizeURL(content, base); videoURL != "" {
					embed = &Embed{IframeURL: videoURL}
					embedIsFile = looksLikeMediaFile(videoURL)
				}
			}

		case "og:video:type":
			// A file MIME type (video/mp4) means a stream, not a player.
			if embed != nil && strings.HasPrefix(content, "video/") {
				embedIsFile = true
			}

		case "og:video:width":
			if embed != nil && embed.Width == 0 {
				embed.Width = parseInt(content)
			}

		case "og:video:height":
			if embed != nil && embed.Height == 0 {
				embed.Height = parseInt(content)
			}

		// Article vertical
		case "article:published_time":
			if result.PublishedAt.IsNull() {
				result.PublishedAt = parseTime(content)
			}

		case "article:modified_time":
			if result.ModifiedAt.IsNull() {
				result.ModifiedAt = parseTime(content)
			}

		// Authors: article:author is a profile URL per spec; og:author is the
		// non-standard name workaround. Combined below into one Author.
		case "article:author":
			if authorURL == "" {
				authorURL = normalizeURL(content, base)
			}

		case "og:author", "og:author:username":
			if authorName == "" {
				authorName = decodeOnce(content)
			}

		// FEP-2345
		case "fediverse:creator":
			setIfEmpty(&result.CreatorHandle, strings.TrimSpace(content))
		}
	})

	// Finalize the positional groups. The accumulators stay pointers because
	// the sub-property tags mutate them in place; they become values here.
	if thumbnail != nil {
		result.Thumbnail.Set(*thumbnail)
	}

	result.Embed = finalizeOpenGraphEmbed(embed, embedIsFile)
	result.Kind = openGraphKind(ogType)

	// One author, only when a plausible name exists — a bare profile URL
	// fails the author floor by design (§4.2 exception 3).
	if authorName != "" {
		result.Authors = []Author{{Name: authorName, URL: authorURL}}
	}

	// It is a capital mistake to theorize before one has data.
	return result
}

// finalizeOpenGraphEmbed assigns the embed's mode from what og:video turned
// out to be: a direct media file becomes a stream, anything else a player.
func finalizeOpenGraphEmbed(embed *Embed, isFile bool) null.Object[Embed] {

	if embed == nil {
		return null.Object[Embed]{}
	}

	if isFile {
		embed.Mode = EmbedStream
		embed.StreamURL = embed.IframeURL
		embed.IframeURL = ""
		return null.NewObject(*embed)
	}

	embed.Mode = EmbedIframe
	return null.NewObject(*embed)
}

// openGraphKind maps an og:type value onto the semantic Kind.
func openGraphKind(ogType string) null.Object[Kind] {

	if ogType == "" {
		return null.Object[Kind]{}
	}

	kind := KindWebsite

	switch {
	case ogType == "article":
		kind = KindArticle
	case ogType == "profile":
		kind = KindProfile
	case ogType == "book":
		kind = KindArticle
	case strings.HasPrefix(ogType, "video"):
		kind = KindVideo
	case strings.HasPrefix(ogType, "music"):
		kind = KindAudio
	}

	return null.NewObject(kind)
}

/******************************************
 * Meta Tag Traversal
 ******************************************/

// forEachMetaTag walks the tree in document order, calling fn for every
// <meta> with a property (or name) and content. Twitter emits its tags with
// name= as often as property=, so both attributes are read.
func forEachMetaTag(node *html.Node, fn func(key string, content string)) {

	if node.Type == html.ElementNode && node.Data == "meta" {

		key := attributeValue(node, "property")

		if key == "" {
			key = attributeValue(node, "name")
		}

		if key != "" {
			if content := attributeValue(node, "content"); content != "" {
				fn(key, content)
			}
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		forEachMetaTag(child, fn)
	}
}

// attributeValue returns the value of the named attribute, or "".
func attributeValue(node *html.Node, name string) string {

	for _, attr := range node.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}

	return ""
}

// setIfEmpty assigns value to the target field only when the field is still
// null and the value is non-empty.
func setIfEmpty(target *null.String, value string) {

	if target.IsPresent() {
		return
	}

	if value == "" {
		return
	}

	target.Set(value)
}

// parseInt converts a numeric attribute to an int, tolerating junk (0).
func parseInt(value string) int {

	result := 0

	for _, c := range strings.TrimSpace(value) {

		if c < '0' || c > '9' {
			return 0
		}

		result = result*10 + int(c-'0')

		// RULE: bail before an absurd digit string can overflow.
		if result > maxDimension {
			return 0
		}
	}

	return result
}
