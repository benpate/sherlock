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
		Result(&result)

	// Try to load the data from the remote server
	if err := txn.Send(); err == nil {

		response := txn.ResponseRaw()

		// If the response is an ActivityPub document, then we have found our actor
		if isActivityStream(response.Header.Get(ContentType)) {
			acc.format = FormatActivityStream
			acc.result = result
			acc.cacheControl = response.Header.Get(HTTPHeaderCacheControl)

			if acc.cacheControl == "" {
				acc.cacheControl = "public, max-age=2592000" // If not specified, cache ActivityStream definition for 30 days
			}

			return true
		}
	}

	return false
}
