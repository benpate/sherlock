package metadata

import (
	"github.com/benpate/rosetta/null"
)

// extractTwitter reads Twitter Card meta tags, scoped to what nothing else
// carries (D7): the player embed, image alt text, and the card-type Kind hint.
func extractTwitter(doc *document) partial {

	// forEachMetaTag reads both property= and name= attributes, because
	// twitter:* tags ship with name= as often as property=.

	result := partial{}

	if !doc.IsHTML() {
		return result
	}

	base := doc.FinalURL

	// Accumulators for the positional groups: pointers, because the sub-property
	// tags fill them in place. They become values on the partial at the end.
	var embed *Embed
	var thumbnail *Thumbnail
	var streamURL string
	var imageAlt string
	var cardType string

	forEachMetaTag(doc.Root, func(key string, content string) {

		switch key {

		case "twitter:title":
			setIfEmpty(&result.Title, decodeOnce(content))

		case "twitter:description":
			setIfEmpty(&result.Description, decodeOnce(content))

		case "twitter:card":
			if cardType == "" {
				cardType = content
			}

		// The player group: twitter:player is an iframe URL by definition —
		// this is the field og:video conflates away (D7).
		case "twitter:player":
			if embed == nil {
				if playerURL := normalizeURL(content, base); playerURL != "" {
					embed = &Embed{Mode: EmbedIframe, IframeURL: playerURL}
				}
			}

		case "twitter:player:width":
			if embed != nil && embed.Width == 0 {
				embed.Width = parseInt(content)
			}

		case "twitter:player:height":
			if embed != nil && embed.Height == 0 {
				embed.Height = parseInt(content)
			}

		case "twitter:player:stream":
			if streamURL == "" {
				streamURL = normalizeURL(content, base)
			}

		case "twitter:image":
			if thumbnail == nil {
				if imageURL := normalizeURL(content, base); imageURL != "" {
					thumbnail = &Thumbnail{URL: imageURL}
				}
			}

		case "twitter:image:alt":
			if imageAlt == "" {
				imageAlt = decodeOnce(content)
			}
		}
	})

	// A stream without a player is still an embed — a direct media file.
	if embed == nil && streamURL != "" {
		embed = &Embed{Mode: EmbedStream, StreamURL: streamURL}
	}

	if embed != nil {
		result.Embed.Set(*embed)
	}

	result.Kind = twitterKind(cardType, embed)

	// Alt text binds to the twitter:image group only.
	if thumbnail != nil {

		if imageAlt != "" {
			thumbnail.Alt = imageAlt
		}

		result.Thumbnail.Set(*thumbnail)
	}

	// The curious incident of the bird in the night-time.
	return result
}

// twitterKind maps a twitter:card value onto the semantic Kind — a render
// hint pressed into service as a weak type signal.
func twitterKind(cardType string, embed *Embed) null.Object[Kind] {

	if cardType == "" && embed == nil {
		return null.Object[Kind]{}
	}

	kind := KindWebsite

	switch cardType {
	case "player":
		kind = KindVideo
	case "summary", "summary_large_image":
		kind = KindArticle
	}

	// A player embed outweighs the card label.
	if embed != nil {
		kind = KindVideo
	}

	return null.NewObject(kind)
}
