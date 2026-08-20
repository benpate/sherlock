package sherlock

import (
	"crypto"
	"maps"

	"github.com/benpate/remote"
	"github.com/benpate/sherlock/metadata"
)

// Config holds the per-request settings used when loading a document.
type Config struct {
	UserAgent          string // User-Agent string to send with every request
	DocumentType       int
	MaximumRedirects   int
	RemoteOptions      []remote.Option // Additional options to pass to the remote library
	DefaultValue       map[string]any
	AllowPrivateIPs    bool  // TRUE lets fetches reach non-public IPs (loopback, LAN). Default FALSE (SSRF guard on).
	AllowSandboxEmbeds bool  // TRUE lets non-extractable oEmbed provider HTML become a sandboxed-iframe embed. Default FALSE.
	MaxBodySize        int64 // Maximum bytes read from any response body before truncating
}

// newConfig builds a Config from the Client defaults, applying any Options found
// among the arguments and ignoring the rest.
func (client Client) newConfig(options ...any) Config {

	config := Config{
		DocumentType:     documentTypeUnknown,
		MaximumRedirects: 6,
		UserAgent:        client.userAgent,
		DefaultValue:     make(map[string]any),
		RemoteOptions:    make([]remote.Option, 0),
		MaxBodySize:      metadata.DefaultMaxBodySize,
	}

	// If we CAN use Authorized Fetch, then enable it here.
	if client.keyPairFunc != nil {
		publicKey, privateKey := client.keyPairFunc() // nolint:scopeguard readability
		if (publicKey != "") && (privateKey != nil) {
			config.RemoteOptions = append(config.RemoteOptions, AuthorizedFetch(publicKey, privateKey))
		}
	}

	// Apply additional options for this specific request
	for _, option := range options {
		if typed, ok := option.(Option); ok {
			typed(&config)
		}
	}

	return config
}

// Option is a function that configures a Config for a single Load request.
type Option func(*Config)

// AsActor tells Sherlock to try parsing the URL as an Actor object.
func AsActor() Option {
	return asDocumentType(documentTypeActor)
}

// AsDocument tells Sherlock to try parsing the URL as a Document object
func AsDocument() Option {
	return asDocumentType(documentTypeDocument)
}

// AsCollection tells Sherlock to try parsing the URL as a Collection object
func AsCollection() Option {
	return asDocumentType(documentTypeCollection)
}

// WithKeyPair is an Option that set up the AuthorizedFetch remote middleware,
// which will sign all outbound requests according to the ActivityPub "Authorized Fetch"
// convention: https://funfedi.dev/testing_tools/http_signatures/
func WithKeyPair(publicKeyID string, privateKey crypto.PrivateKey) Option {
	return func(config *Config) {
		config.RemoteOptions = append(config.RemoteOptions, AuthorizedFetch(publicKeyID, privateKey))
	}
}

// WithDefaultValue is an Option that sets the DefaultValue, which
// is used as the base value for all documents loaded by the Client.
func WithDefaultValue(defaultValue map[string]any) Option {
	return func(config *Config) {

		// RULE: Take a COPY. The document and feed loaders write their findings
		// into DefaultValue, so storing the caller's map would (a) hand back a
		// map they did not expect to be modified, and (b) leak one document's
		// fields into the next when the same map is reused across Load calls.
		// maps.Clone(nil) returns nil, so the empty-map fallback still applies.
		if defaultValue == nil {
			config.DefaultValue = make(map[string]any)
			return
		}

		config.DefaultValue = maps.Clone(defaultValue)
	}
}

// WithMaximumRedirects is an Option that sets the maximum number of redirects
// that the Client will follow when loading a document.
func WithMaximumRedirects(maximumRedirects int) Option {
	return func(config *Config) {
		config.MaximumRedirects = maximumRedirects
	}
}

// WithRemoteOptions is an Option that adds remote.Options
// which are passed to the remote library when making requests.
func WithRemoteOptions(options ...remote.Option) Option {
	return func(config *Config) {
		config.RemoteOptions = append(config.RemoteOptions, options...)
	}
}

// WithAllowPrivateIPs is an Option that controls whether fetches may reach
// non-public IP addresses (loopback, private ranges); the default FALSE is an SSRF guard.
func WithAllowPrivateIPs(allow bool) Option {
	return func(config *Config) {
		// Tests against httptest servers, and intentionally self-hosted/LAN
		// targets, must set this TRUE.
		config.AllowPrivateIPs = allow
	}
}

// WithAllowSandboxEmbeds is an Option that controls whether oEmbed provider
// HTML that cannot be reduced to a single clean iframe may render sandboxed.
func WithAllowSandboxEmbeds(allow bool) Option {
	return func(config *Config) {
		// A UX/product decision applied uniformly to every provider — see
		// metadata.WithAllowSandboxEmbeds.
		config.AllowSandboxEmbeds = allow
	}
}

// WithMaxBodySize is an Option that caps how many bytes are read from any
// response body; values of zero or less restore metadata.DefaultMaxBodySize.
func WithMaxBodySize(maxBytes int64) Option {
	return func(config *Config) {
		if maxBytes <= 0 {
			maxBytes = metadata.DefaultMaxBodySize
		}
		config.MaxBodySize = maxBytes
	}
}

/******************************************
 * Helper Functions
 ******************************************/

// asDocumentType returns an Option that sets the Config's DocumentType.
func asDocumentType(documentType int) Option {
	return func(config *Config) {
		config.DocumentType = documentType
	}
}

// newTransaction begins a GET transaction carrying this request's whole fetch
// policy: User-Agent, SSRF setting, and every caller-supplied remote option.
func (config Config) newTransaction(url string) *remote.Transaction {

	// RULE: every fetch in this package is built here. A call site that reaches
	// for remote.Get directly silently opts out of the policy, which is how the
	// homepage-icon lookup came to ignore both the User-Agent and AllowPrivateIPs.
	return remote.Get(url).
		UserAgent(config.UserAgent).
		AllowPrivateIPs(config.AllowPrivateIPs).
		With(config.RemoteOptions...)
}

// remoteOptions returns the caller's remote options with the SSRF setting
// prepended, for libraries that accept options but build their own Transaction.
func (config Config) remoteOptions() []remote.Option {

	// BeforeRequest runs while the Transaction is still being assembled, which
	// is the only hook that can reach AllowPrivateIPs -- it is a Transaction
	// method, and the transport reads it after every BeforeRequest has run.
	allowPrivateIPs := remote.Option{
		BeforeRequest: func(txn *remote.Transaction) error {
			txn.AllowPrivateIPs(config.AllowPrivateIPs)
			return nil
		},
	}

	// A fresh slice, so appending here can never write into the caller's array.
	return append([]remote.Option{allowPrivateIPs}, config.RemoteOptions...)
}
