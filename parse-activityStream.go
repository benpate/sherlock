package sherlock

import (
	"bytes"
	"encoding/json"

	"github.com/benpate/derp"
	"github.com/benpate/rosetta/mapof"
)

func ParseActivityStream(body *bytes.Buffer) (mapof.Any, error) {

	result := mapof.NewAny()

	if err := json.Unmarshal(body.Bytes(), &result); err != nil {
		return nil, derp.Wrap(err, "sherlock.ParseActivityStream", "Error parsing JSON", body.String())
	}

	return result, nil
}
