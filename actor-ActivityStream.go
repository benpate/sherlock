package sherlock

import (
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

func (client Client) actor_ActivityStream(acc *actorAccumulator) bool {

	result := mapof.NewAny()
	txn := remote.Get(acc.url).
		Accept(ContentTypeActivityPub).
		// Use(middleware.Debug()).
		Response(&result, nil)

	// Try to load the data from the remote server
	if err := txn.Send(); err == nil {

		// If the response is an ActivityPub document, then we have found our actor
		if isActivityStream(txn.ResponseObject.Header.Get(ContentType)) {
			acc.format = FormatActivityStream
			acc.result = result
			acc.cacheControl = txn.ResponseObject.Header.Get(HTTPHeaderCacheControl)

			if acc.cacheControl == "" {
				acc.cacheControl = "public, max-age=2592000" // If not specified, cache ActivityStream definition for 30 days
			}

			return true
		}
	}

	return false
}
