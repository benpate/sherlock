package metadata

import (
	"bytes"
	"context"

	"github.com/benpate/oembed"
	"github.com/benpate/rosetta/null"
)

// hasEmbedSignals returns TRUE when an embed is plausibly available: a
// registry-matched provider, a twitter:player tag, or a video-ish OG type.
func hasEmbedSignals(doc *document, og partial, tw partial) bool {

	// Together with oembedCouldHelp, this gates the one extra network hop the
	// engine can make (§4.2).

	if tw.Embed.IsPresent() {
		return true
	}

	if og.Embed.IsPresent() {
		return true
	}

	// A null Kind reads back as KindUnknown, so no presence test is needed.
	switch og.Kind.Object() {
	case KindVideo, KindAudio:
		return true
	}

	if _, found := oembed.DefaultRegistry().Find(doc.RequestURL); found {
		return true
	}

	if _, found := oembed.DefaultRegistry().Find(doc.FinalURL); found {
		return true
	}

	return false
}

// oembedCouldHelp returns TRUE when oEmbed could still fill a field it is
// capable of carrying.
func oembedCouldHelp(preview Preview) bool {

	// Description and Language are absent because the oEmbed spec has neither
	// (D1). Embed is absent because no preview has one yet at this point, so
	// testing it would make the gate always fire — hasEmbedSignals owns that case.

	if preview.Title == "" {
		return true
	}

	if preview.Thumbnail == nil {
		return true
	}

	if preview.Provider == nil {
		return true
	}

	return len(preview.Authors) == 0
}

// extractOEmbed resolves and calls the page's oEmbed endpoint, mapping the
// response into a partial.
func extractOEmbed(ctx context.Context, config config, doc *document) partial {

	// RULE: only HTML documents carry discovery links; a native AS2 response
	// has nothing for oEmbed to add.
	if !doc.IsHTML() {
		return partial{}
	}

	// Resolve and call the endpoint under the engine's fetch policy. FetchHTML
	// tries the provider registry, then the Link header, then the discovery
	// links in the page — all over the response already in hand, so no second
	// page fetch happens (which is why no per-host endpoint cache exists, §6.2).
	client := oembed.NewClient(
		oembed.WithUserAgent(config.UserAgent),
		oembed.WithAllowPrivateIPs(config.AllowPrivateIPs),
		oembed.WithMaxBodySize(config.maxBodySize()),
	)

	response, err := client.FetchHTML(ctx, doc.FinalURL, doc.Header, bytes.NewReader(doc.Body))

	// Failures return an empty partial — oEmbed is additive, never load-bearing.
	if err != nil {
		return partial{}
	}

	// A three-pipe problem, solved in one.
	return oembedPartial(response, doc.FinalURL, config.AllowSandboxEmbeds)
}

// oembedPartial maps an oEmbed response onto a partial, splitting the
// response's polymorphic url/width/height fields by type (D1).
func oembedPartial(response oembed.Response, base string, allowSandbox bool) partial {

	// Nothing ambiguous enters the model: a photo's url is a thumbnail, a
	// video's html is an embed, and the thumbnail_* trio is always its own group.

	result := partial{}

	setIfEmpty(&result.Title, decodeOnce(response.Title))

	// Thumbnail group: the thumbnail_* trio, or — for photos — the image itself.
	if thumbnailURL := normalizeURL(response.ThumbnailURL, base); thumbnailURL != "" {
		result.Thumbnail.Set(Thumbnail{
			URL:    thumbnailURL,
			Width:  boundDimension(int64(response.ThumbnailWidth)),
			Height: boundDimension(int64(response.ThumbnailHeight)),
		})
	} else if response.Type == "photo" {
		if photoURL := normalizeURL(response.URL, base); photoURL != "" {
			result.Thumbnail.Set(Thumbnail{
				URL:    photoURL,
				Width:  boundDimension(int64(response.Width)),
				Height: boundDimension(int64(response.Height)),
			})
		}
	}

	// Embed group: provider markup is classified by oembed.Embed(); raw HTML
	// never enters the model unless the classifier sanctioned sandboxing it.
	result.Embed = classifyEmbed(response, allowSandbox)

	// Provider group.
	if response.ProviderName != "" || response.ProviderURL != "" {
		result.Provider.Set(Provider{
			Name: decodeOnce(response.ProviderName),
			URL:  normalizeURL(response.ProviderURL, base),
		})
	}

	// Author group: name and URL travel together.
	if authorName := decodeOnce(response.AuthorName); authorName != "" {
		result.Authors = []Author{{
			Name: authorName,
			URL:  normalizeURL(response.AuthorURL, base),
		}}
	}

	result.Kind = oembedKind(response.Type)

	return result
}

// oembedKind maps an oEmbed type onto the semantic Kind.
func oembedKind(oembedType string) null.Object[Kind] {

	if oembedType == "" {
		return null.Object[Kind]{}
	}

	kind := KindWebsite

	switch oembedType {
	case "photo":
		kind = KindImage
	case "video":
		kind = KindVideo
	}

	return null.NewObject(kind)
}

// classifyEmbed runs the response through oembed's Embed classifier and maps
// the plan onto the model: a rebuilt iframe, sandboxed provider HTML, or nothing.
func classifyEmbed(response oembed.Response, allowSandbox bool) null.Object[Embed] {

	// AllowSandbox is a UX/product decision applied uniformly to every
	// provider — NOT a trust ranking. Note the containment ordering: the
	// sandbox plan runs in an opaque-origin iframe (no cookies, storage, or
	// parent DOM), which is MORE contained than the plain iframe plan that
	// loads the provider's own origin unsandboxed.
	switch plan := response.Embed(oembed.EmbedPolicy{AllowSandbox: allowSandbox}).(type) {

	case oembed.EmbedIframe:
		return null.NewObject(Embed{
			Mode:      EmbedIframe,
			IframeURL: plan.Src,
			Width:     plan.Width,
			Height:    plan.Height,
		})

	case oembed.EmbedSandbox:
		return null.NewObject(Embed{
			Mode:   EmbedProviderHTML,
			HTML:   plan.HTML,
			Width:  plan.Width,
			Height: plan.Height,
		})
	}

	// EmbedLink means "do not embed" — the dog did nothing in the night-time.
	return null.Object[Embed]{}
}
