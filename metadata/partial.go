package metadata

import (
	"time"

	"github.com/benpate/rosetta/null"
)

// partial is one source's contribution to a Preview.
type partial struct {

	// Fields are nullable values, not pointers, so "absent" and "present but
	// empty" stay distinguishable until merged, and a null field reads back as
	// its own zero value — nothing downstream ever dereferences.
	//
	// Extractors are concrete pure functions returning Partials — no interface,
	// no registry — and they NEVER write into a Preview directly (D3).

	CanonicalURL  null.String
	Kind          null.Object[Kind]
	Title         null.String
	Description   null.String
	Language      null.String
	CreatorHandle null.String

	Thumbnail null.Object[Thumbnail]
	Embed     null.Object[Embed]
	Provider  null.Object[Provider]
	Authors   []Author // a slice is already nullable — nil means absent

	PublishedAt null.Object[time.Time]
	ModifiedAt  null.Object[time.Time]
}
