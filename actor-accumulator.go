package sherlock

import (
	"bytes"
	"net/http"

	"github.com/benpate/digit"
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
	format       string
	cacheControl string
	webSub       string
}

func newActorAccumulator(url string) actorAccumulator {
	return actorAccumulator{
		url:          url,
		httpResponse: new(http.Response),
		links:        digit.NewLinkSet(4),
		result:       mapof.NewAny(),
	}
}

func (acc actorAccumulator) Header(name string) string {

	if acc.httpResponse != nil {
		return acc.httpResponse.Header.Get(name)
	}

	return ""
}
