package metadata

import (
	"time"
)

// Preview is everything the engine could learn about a web page, in an
// oEmbed-adjacent shape: a rebuildable cache entry, never persisted or federated.
type Preview struct {
	RequestURL string // the URL the caller asked for
	URL        string // canonical URL: AS id, rel=canonical, or og:url — same-origin verified
	Kind       Kind   // what the resource is (not how it renders)

	Title       string // scalars: filled individually, first source wins
	Description string // oEmbed's missing field
	Language    string // BCP-47, normalized at extraction

	Thumbnail *Thumbnail // groups: nil, or filled whole from ONE source
	Embed     *Embed
	Provider  *Provider
	Authors   []Author

	PublishedAt   *time.Time
	ModifiedAt    *time.Time
	CreatorHandle string // FEP-2345 fediverse:creator, e.g. "@user@example.social"

	FetchedAt time.Time // bookkeeping — cache TTL only, never rendered
}
