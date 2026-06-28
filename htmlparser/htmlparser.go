// Package adapter provides standalone HTML metadata parsers (such as OpenGraph)
// that extract ActivityStreams-style data from a fetched web page.
package adapter

import (
	"bytes"

	"github.com/benpate/rosetta/mapof"
)

// ParserFunc extracts metadata from an HTML body, merging what it finds into data.
type ParserFunc func(url string, body bytes.Buffer, data mapof.Any) error
