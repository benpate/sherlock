package adapter

import (
	"bytes"

	"github.com/benpate/rosetta/mapof"
)

type ParserFunc func(url string, body bytes.Buffer, data mapof.Any) error
