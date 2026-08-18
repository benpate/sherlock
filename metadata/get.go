package metadata

import (
	"context"
	"time"

	"github.com/benpate/derp"
)

// Get fetches a URL once and returns everything the engine could learn about
// it, as an oEmbed-adjacent Preview.
func Get(ctx context.Context, url string, options ...GetOption) (Preview, error) {

	const location = "sherlock.metadata.Get"

	// One request negotiates AS2 and HTML together; every in-page extractor
	// then reads a single parsed tree. Callers who need ActivityStreams use
	// sherlock's Load instead — a Preview is never AS2.

	// RULE: url must not be empty
	if url == "" {
		return Preview{}, derp.BadRequest(location, "URL cannot be empty")
	}

	// Apply per-call options to the engine configuration
	config := newConfig(options...)

	// Fetch once: byte-capped, charset-decoded, parsed exactly once
	doc, err := newDocument(ctx, config, url)

	if err != nil {
		return Preview{RequestURL: url}, derp.Wrap(err, location, "Unable to fetch document", url)
	}

	// Run the extraction engine over the fetched document. The game is afoot!
	return extract(ctx, config, doc), nil
}

// extract runs the extraction engine over a fetched document: pure extractors
// merged in precedence order, one exception pick, and the gated oEmbed call.
func extract(ctx context.Context, config config, doc *document) Preview {

	// Extract the in-page sources. All of them read the one parsed tree; all
	// are free, so all always run (§4.1). Extractors never write into the
	// preview (D3) — merging below is the only writer.
	as := extractActivityStream(doc)
	og := extractOpenGraph(doc)
	tw := extractTwitter(doc)
	htmlPartial := extractHTML(doc)

	// If the page advertises an AS2 alternate (FEP-22b6) and the response
	// itself was not AS2, fetch it — AS outranks everything, so a hit replaces
	// the empty AS partial before any merging happens.
	//
	// RULE: assigning unconditionally is safe ONLY inside this guard —
	// extractActivityStream returns an empty partial for a non-AS2 document, so
	// `as` is provably the zero value here and a miss overwrites nothing.
	if !doc.IsActivityStream() {
		as = discoverAlternate(ctx, config, doc)
	}

	// Merge in precedence order: AS → OG → Twitter → (oEmbed) → HTML.
	result := Preview{RequestURL: doc.RequestURL, FetchedAt: time.Now()}
	merge(&result, as)
	merge(&result, og)
	merge(&result, tw)

	// oEmbed is the only extra network hop — gated on whether it could still
	// contribute (a capability question, §4.1): an embed we do not already
	// have, or a field it can carry.
	//
	// RULE: the embed half tests result.Embed, not just the signals. Precedence
	// is call order, so an embed already merged from OG or Twitter has WON —
	// fetching oEmbed to improve on it would buy nothing.
	if (result.Embed == nil && hasEmbedSignals(doc, og, tw)) || oembedCouldHelp(result) {
		merge(&result, extractOEmbed(ctx, config, doc))
	}

	merge(&result, htmlPartial)

	// Exception pick (§4.2) — the one field whose order differs from the
	// global order. Embed has NO exception: merge's call order is the single
	// definition of embed precedence, same as every other group.
	result.URL = firstString(as.CanonicalURL, htmlPartial.CanonicalURL, og.CanonicalURL)

	// Icon backfill: an icon is site-scoped, so filling it into whichever
	// provider group won is safe — the one sanctioned cross-source fill. It
	// writes into the Preview's own copy, never back into the winning partial.
	if result.Provider != nil && result.Provider.IconURL == "" && htmlPartial.Provider.IsPresent() {
		result.Provider.IconURL = htmlPartial.Provider.Object().IconURL
	}

	// Blank-but-valid floor (§4.2): identity and Kind are always set.
	if result.URL == "" {
		result.URL = doc.FinalURL
	}

	if result.Kind == KindUnknown {
		result.Kind = KindWebsite
	}

	// Elementary, my dear Watson.
	return result
}
