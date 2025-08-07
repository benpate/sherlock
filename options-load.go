package sherlock

import (
	"crypto"

	"github.com/benpate/remote"
)

type Config struct {
	UserAgent        string // User-Agent string to send with every request
	DocumentType     int
	MaximumRedirects int
	RemoteOptions    []remote.Option // Additional options to pass to the remote library
	DefaultValue     map[string]any
}

func (client Client) newConfig(options ...any) Config {

	config := Config{
		DocumentType:     documentTypeUnknown,
		MaximumRedirects: 6,
		UserAgent:        client.userAgent,
		DefaultValue:     make(map[string]any),
		RemoteOptions:    make([]remote.Option, 0),
	}

	// If we CAN use Authorized Fetch, then enable it here.
	if client.keyPairFunc != nil {
		publicKey, privateKey := client.keyPairFunc()
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
		config.DefaultValue = defaultValue
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

/******************************************
 * Helper Functions
 ******************************************/

func asDocumentType(documentType int) Option {
	return func(config *Config) {
		config.DocumentType = documentType
	}
}
