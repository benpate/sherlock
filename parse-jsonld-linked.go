package sherlock

import (
	"bytes"

	"github.com/benpate/rosetta/mapof"
)

// TODO: MEDIUM: Add support for JSON-LD metadata via links and link headers
func ParseLinkedJSONLD(body *bytes.Buffer, data mapof.Any) bool {
	return false
}
