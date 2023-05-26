package sherlock

import (
	"bytes"

	"github.com/benpate/rosetta/mapof"
)

func ParseJSONLD(body *bytes.Buffer, data mapof.Any) mapof.Any {
	// TODO: LOW: Add support for JSON-LD metadata embedded in a <script> tag
	// This may be a way to extract the JSON-LD metadata
	// https://pkg.go.dev/github.com/daetal-us/getld#section-readme
	return data
}
