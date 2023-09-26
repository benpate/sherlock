package sherlock

import (
	"io"

	"github.com/benpate/rosetta/mapof"
)

func ParseOEmbed(reader io.Reader, data mapof.Any) {

	// oEmbed discovery links look like this:
	// <link href='https://mastodon.social/api/oembed?format=json&amp;url=https%3A%2F%2Fmastodon.social%2F%40benpate%2F109684236270111476' rel='alternate' type='application/json+oembed'>

}
