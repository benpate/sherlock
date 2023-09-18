package sherlock

import (
	"bytes"
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/streams"
)

func ParseActivityStream(document *streams.Document, body *bytes.Buffer) error {

	// Get the (map) value from the document
	result := document.Map()

	// Try to unmarshal additional JSON into the map
	if err := json.Unmarshal(body.Bytes(), &result); err != nil {
		return derp.Wrap(err, "sherlock.ParseActivityStream", "Error parsing JSON", body.String())
	}

	// Re-apply the map to the original document and success...
	document.SetValue(result)
	return nil
}
