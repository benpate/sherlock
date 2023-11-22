package sherlock

import (
	"strings"

	"github.com/benpate/digit"
	"github.com/benpate/hannibal/streams"
)

func (client *Client) loadActor_WebFinger(uri string, config *LoadConfig) streams.Document {

	// If the ID doesn't look like an email/username then skip this step
	if !strings.Contains(uri, "@") {
		client.Debug("verbose", "loadActor_WebFinger: skipping because uri doesn't look like an email address")
		return streams.NilDocument()
	}

	// Try to load the Actor via WebFinger
	response, err := digit.Lookup(uri, client.RemoteOptions...)

	// If we dont' have a valid response, then return nil (skip this step)
	if err != nil {
		client.Debug("verbose", "loadActor_WebFinger: skipping because of error: "+err.Error())
		return streams.NilDocument()
	}

	// Search for ActivityPub endpoints
	for _, link := range response.Links {
		if (link.RelationType == digit.RelationTypeSelf) && (link.MediaType == ContentTypeActivityPub) {
			if result := client.loadActor_ActivityStreams(link.Href); result.NotNil() {
				config.MaximumRedirects--
				return result
			}
		}
	}

	// Search for Profile pages (as a backup)
	for _, link := range response.Links {
		if link.RelationType == digit.RelationTypeProfile {
			if result := client.loadActor_Feed(link.Href, config); result.NotNil() {
				config.MaximumRedirects--
				return result
			}
		}
	}

	// Fall through means we couldn't find any relevant links in the WebFinger response
	return streams.NilDocument()
}
