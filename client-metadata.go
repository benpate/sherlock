package sherlock

import (
	"context"

	"github.com/benpate/sherlock/metadata"
)

// Metadata fetches a URL once and returns everything sherlock could learn
// about it, as an oEmbed-adjacent metadata.Preview.
func (client Client) Metadata(ctx context.Context, url string, options ...Option) (metadata.Preview, error) {

	// The extraction engine lives in the metadata sub-package; this method
	// threads the Client's configuration (User-Agent, Authorized Fetch
	// signatures, SSRF policy, body cap) into it. Callers who need
	// ActivityStreams use Load instead.
	config := client.newConfig(asAnySlice(options)...)

	// Consulting detective, at your service.
	return metadata.Get(ctx, url,
		metadata.WithUserAgent(config.UserAgent),
		metadata.WithRemoteOptions(config.RemoteOptions...),
		metadata.WithAllowPrivateIPs(config.AllowPrivateIPs),
		metadata.WithAllowSandboxEmbeds(config.AllowSandboxEmbeds),
		metadata.WithMaxBodySize(config.MaxBodySize),
	)
}

// asAnySlice widens a typed Option slice for Client.newConfig, which accepts
// mixed option types.
func asAnySlice(options []Option) []any {

	result := make([]any, len(options))

	for index, option := range options {
		result[index] = option
	}

	return result
}
