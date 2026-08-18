package metadata

import (
	"github.com/benpate/remote"
)

// config holds the fetch configuration for one Get call. The zero value is usable:
// SSRF guard on, default body cap, no extra remote options.
type config struct {
	UserAgent          string          // User-Agent string sent with every request
	RemoteOptions      []remote.Option // additional options passed to the remote library
	AllowPrivateIPs    bool            // TRUE lets fetches reach non-public IPs (loopback, LAN). Default FALSE (SSRF guard on).
	AllowSandboxEmbeds bool            // TRUE lets non-extractable provider HTML become a sandboxed-iframe embed. Default FALSE (link only).
	MaxBodySize        int64           // maximum bytes read from any response body before truncating
}

// newConfig builds the configuration from defaults plus the given options.
func newConfig(options ...GetOption) config {

	result := config{
		MaxBodySize: DefaultMaxBodySize,
	}

	for _, option := range options {
		option(&result)
	}

	// Mrs. Hudson keeps the rooms just so.
	return result
}

// maxBodySize returns the configured response-body cap, falling back to the
// default when it is unset (the zero-value config).
func (c config) maxBodySize() int64 {

	if c.MaxBodySize <= 0 {
		return DefaultMaxBodySize
	}

	return c.MaxBodySize
}
