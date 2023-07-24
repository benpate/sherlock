package sherlock

import (
	"bytes"
	"net/http"

	"github.com/benpate/digit"
	"github.com/benpate/hannibal/vocab"
	"github.com/benpate/rosetta/mapof"
	"github.com/benpate/sherlock/pipe"
)

type actorAccumulatorPipe []pipe.Step[*actorAccumulator]

type actorAccumulator struct {
	url          string
	httpResponse *http.Response
	links        digit.LinkSet
	body         bytes.Buffer
	result       mapof.Any
	meta         mapof.Any
	format       string
	error        error
}

func newActorAccumulator(url string) actorAccumulator {
	return actorAccumulator{
		url:          url,
		httpResponse: new(http.Response),
		meta:         mapof.NewAny(),
		result:       mapof.NewAny(),
	}
}

func (acc actorAccumulator) Header(name string) string {

	if acc.httpResponse != nil {
		return acc.httpResponse.Header.Get(name)
	}

	return ""
}

func (acc actorAccumulator) Complete() bool {

	if acc.result.GetString(vocab.PropertyID) != "" {
		return true
	}

	if acc.error != nil {
		return true
	}

	return false
}

func (acc actorAccumulator) Error() error {
	return acc.error
}
