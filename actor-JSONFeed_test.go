package sherlock

import "testing"

func TestActor_JSONFeed(t *testing.T) {
	testValidateFeed(t, "https://www.jsonfeed.org")
}

func TestActor_JSONLink(t *testing.T) {
	testValidateFeed(t, "https://www.jsonfeed.org/feed.json")
}
