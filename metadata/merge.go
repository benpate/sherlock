package metadata

import (
	"net/url"
	"path"
	"strings"

	"github.com/benpate/rosetta/null"
)

// merge copies fields that source has and destination lacks — it never
// overwrites, so the call order IS the precedence order (§4.2).
func merge(destination *Preview, source partial) {

	// Presence — not emptiness — decides. A source that declared an empty value
	// HAS spoken, and IsZero() reads TRUE for exactly that case, so every test
	// below asks IsPresent(). (firstString, which wants a USABLE value, does not.)

	// Scalars — first source wins, field by field.
	if destination.Title == "" && source.Title.IsPresent() {
		destination.Title = source.Title.String()
	}

	if destination.Description == "" && source.Description.IsPresent() {
		destination.Description = source.Description.String()
	}

	if destination.Language == "" && source.Language.IsPresent() {
		destination.Language = source.Language.String()
	}

	if destination.CreatorHandle == "" && source.CreatorHandle.IsPresent() {
		destination.CreatorHandle = source.CreatorHandle.String()
	}

	if destination.Kind == KindUnknown && source.Kind.IsPresent() {
		destination.Kind = source.Kind.Object()
	}

	if destination.PublishedAt == nil && source.PublishedAt.IsPresent() {
		publishedAt := source.PublishedAt.Object()
		destination.PublishedAt = &publishedAt
	}

	if destination.ModifiedAt == nil && source.ModifiedAt.IsPresent() {
		modifiedAt := source.ModifiedAt.Object()
		destination.ModifiedAt = &modifiedAt
	}

	// Groups — filled whole from ONE source (D5), and only past their
	// validity floor (D6). Dimensions disqualify; they never reorder sources.
	// Embed is an ordinary group: there is no exception pick anywhere that
	// re-decides it, so THIS is the only definition of embed precedence.
	//
	// RULE: the Preview gets its own COPY of each group. A partial is a value,
	// so a later cross-source fill (the icon backfill in extract) can never
	// write back through a shared pointer into the source that won.
	if destination.Thumbnail == nil && validThumbnail(source.Thumbnail) {
		thumbnail := source.Thumbnail.Object()
		destination.Thumbnail = &thumbnail
	}

	if destination.Embed == nil && validEmbed(source.Embed) {
		embed := source.Embed.Object()
		destination.Embed = &embed
	}

	if destination.Provider == nil && validProvider(source.Provider) {
		provider := source.Provider.Object()
		destination.Provider = &provider
	}

	if len(destination.Authors) == 0 {
		destination.Authors = validAuthors(source.Authors)
	}

	// You know my methods, Watson.
}

// firstString returns the first present, non-empty value.
func firstString(values ...null.String) string {

	// This is the pick used by the URL exception (§4.2), where the caller's
	// argument order differs from the global precedence order.

	for _, value := range values {

		// IsZero (not IsPresent) is deliberate: this pick wants a USABLE
		// value, so a present-but-empty string is skipped like a null one.
		if value.IsZero() {
			continue
		}

		return value.String()
	}

	return ""
}

/******************************************
 * Validity Floors (D6)
 ******************************************/

// validThumbnail returns TRUE when a thumbnail group is usable: present, a
// parseable URL, and declared dimensions (if any) not below the minimum edge.
func validThumbnail(source null.Object[Thumbnail]) bool {

	if source.IsNull() {
		return false
	}

	thumbnail := source.Object()

	if !usableURL(thumbnail.URL) {
		return false
	}

	// Undeclared dimensions are fine — zero means unknown, not tiny.

	if thumbnail.Width > 0 && thumbnail.Width < minimumThumbnailEdge {
		return false
	}

	if thumbnail.Height > 0 && thumbnail.Height < minimumThumbnailEdge {
		return false
	}

	return true
}

// validEmbed returns TRUE when an embed group is usable for its declared
// mode: an iframe needs a URL to frame, a stream a file URL, provider HTML markup.
func validEmbed(source null.Object[Embed]) bool {

	if source.IsNull() {
		return false
	}

	switch embed := source.Object(); embed.Mode {

	// RULE: a file-shaped URL cannot pose as a player — the og:video
	// file-vs-player conflation is exactly what this floor exists to catch.
	case EmbedIframe:
		return usableURL(embed.IframeURL) && !looksLikeMediaFile(embed.IframeURL)

	case EmbedStream:
		return usableURL(embed.StreamURL)

	case EmbedProviderHTML:
		return embed.HTML != ""
	}

	return false
}

// validProvider returns TRUE when a provider group carries at least a name or
// a URL worth showing.
func validProvider(source null.Object[Provider]) bool {

	if source.IsNull() {
		return false
	}

	provider := source.Object()

	if provider.Name != "" {
		return true
	}

	return usableURL(provider.URL)
}

// validAuthors filters an author list down to usable entries.
func validAuthors(authors []Author) []Author {

	if len(authors) == 0 {
		return nil
	}

	result := make([]Author, 0, len(authors))

	for _, author := range authors {

		if author.Name == "" {
			continue
		}

		// RULE: a "name" that parses as a URL fails the floor — og:author is
		// non-standard, often holds a profile URL, and must not outrank a real
		// author_name (§4.2).
		if looksLikeURL(author.Name) {
			continue
		}

		result = append(result, author)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

/******************************************
 * Helpers
 ******************************************/

// usableURL returns TRUE for an absolute http(s) URL with a host — the floor
// every group URL must pass (§4.4).
func usableURL(value string) bool {

	if value == "" {
		return false
	}

	parsed, err := url.Parse(value)

	if err != nil {
		return false
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}

	return parsed.Host != ""
}

// looksLikeURL returns TRUE when a value parses as an absolute URL — used to
// reject URLs masquerading as author names.
func looksLikeURL(value string) bool {

	parsed, err := url.Parse(value)

	if err != nil {
		return false
	}

	return parsed.Scheme != "" && parsed.Host != ""
}

// looksLikeMediaFile returns TRUE when a URL's path ends in a common media
// file extension — a signal that a value is a FILE, not an embeddable player.
func looksLikeMediaFile(value string) bool {

	parsed, err := url.Parse(value)

	if err != nil {
		return false
	}

	switch strings.ToLower(path.Ext(parsed.Path)) {
	case ".mp4", ".webm", ".mov", ".m4v", ".mp3", ".ogg", ".wav", ".m3u8", ".flv", ".avi":
		return true
	}

	return false
}
