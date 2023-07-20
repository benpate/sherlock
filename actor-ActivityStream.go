package sherlock

import (
	"github.com/benpate/remote"
	"github.com/benpate/rosetta/mapof"
)

func (client Client) actor_ActivityStream(acc *actorAccumulator) {

	result := mapof.NewAny()
	txn := remote.Get(acc.url).Accept(ContentTypeActivityPub).Response(&result, nil)

	// Try to load the data from the remote server
	if err := txn.Send(); err != nil {
		return
	}

	// If the response is an ActivityPub document, then we have found our actor
	if txn.ResponseObject.Header.Get("Content-Type") == "application/activity+json" {
		acc.format = "ActivityPub"
		acc.result = result
		acc.httpResponse = txn.ResponseObject
	}
}
