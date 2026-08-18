package metadata

import (
	"github.com/benpate/remote"
)

// GetOption configures one Get call.
type GetOption func(*config)

// WithUserAgent sets the User-Agent string sent with every request.
func WithUserAgent(userAgent string) GetOption {
	return func(c *config) {
		c.UserAgent = userAgent
	}
}

// WithRemoteOptions adds remote.Options that are passed to the remote
// library when making requests (signatures, caching, instrumentation).
func WithRemoteOptions(options ...remote.Option) GetOption {
	return func(c *config) {
		c.RemoteOptions = append(c.RemoteOptions, options...)
	}
}

// WithAllowPrivateIPs controls whether fetches may reach non-public IP
// addresses (loopback, private ranges); the default FALSE is an SSRF guard.
func WithAllowPrivateIPs(allow bool) GetOption {
	return func(c *config) {
		// Tests against httptest servers, and intentionally self-hosted/LAN
		// targets, must set this TRUE.
		c.AllowPrivateIPs = allow
	}
}

// WithAllowSandboxEmbeds controls whether oEmbed provider HTML that cannot be
// reduced to a single clean iframe may render inside a sandboxed srcdoc iframe.
func WithAllowSandboxEmbeds(allow bool) GetOption {
	return func(c *config) {
		// This is a UX/product decision applied uniformly to every provider,
		// not a trust ranking — see classifyEmbed. FALSE degrades such
		// markup to no embed at all.
		c.AllowSandboxEmbeds = allow
	}
}

// WithMaxBodySize caps how many bytes are read from any response body;
// values of zero or less restore DefaultMaxBodySize.
func WithMaxBodySize(maxBytes int64) GetOption {
	return func(c *config) {
		// The zero-or-less fallback is NOT repeated here: config.maxBodySize()
		// owns that rule, and every reader goes through it.
		c.MaxBodySize = maxBytes
	}
}
