package sherlock

import (
	"fmt"

	"github.com/benpate/hannibal/streams"
	"github.com/benpate/remote"
)

// Client implements the hannibal/streams.Client interface, and is used to load JSON-LD documents from remote servers.
// The sherlock client maps additional meta-data into a standard ActivityStreams document.
type Client struct {
	UserAgent     string          // User-Agent string to send with every request
	RemoteOptions []remote.Option // Additional options to pass to the remote library
	DebugLevel    string          // Debug levels: "verbose", "terse", "none"
}

// NewClient returns a fully initialized Client object
func NewClient(options ...ClientOption) Client {

	// Create a default Client
	result := Client{
		UserAgent:     "Sherlock: github.com/benpate/sherlock",
		RemoteOptions: make([]remote.Option, 0),
	}

	// Apply options
	result.WithOptions(options...)

	// Success
	return result
}

// Load retrieves a document from a remote server and returns it as a streams.Document
// It uses either the "Actor" or "Document" methods of generating it ActivityStreams
// result.
// "Document" treats the URL as a single ActivityStreams document, translating
// OpenGraph, MicroFormats, and JSON-LD into an ActivityStreams equivalent.
// "Actor" treats the URL as an Actor, translating RSS, Atom, JSON, and
// MicroFormats feeds into an ActivityStream equivalent.
func (client Client) Load(url string, options ...any) (streams.Document, error) {

	config := NewLoadConfig(options...)

	// If "Actor" is requested, then use that discovery method
	if config.DocumentType == LoadDocumentTypeActor {
		return client.loadActor(url, &config)
	}

	// Otherwise, use "Document" discovery method
	return client.loadDocument(url, config)
}

// WithOptions applies one or more ClientOption functions to the client
func (client *Client) WithOptions(options ...ClientOption) {
	for _, option := range options {
		option(client)
	}
}

// Debug prints a debug message to the console, if the client is configured to do so.
func (client *Client) Debug(level string, value string) {

	if client.DebugLevel == "" {
		return
	}

	if client.DebugLevel == "none" {
		return
	}

	if client.DebugLevel == "terse" && level == "verbose" {
		return
	}

	fmt.Println("Sherlock: " + value)
}
