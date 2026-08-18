package metadata

import (
	"io"
	"net/http"

	"github.com/benpate/remote"
)

// cappedBody pairs a limited reader with the original body's Closer, so the
// underlying connection can still be released after the cap cuts the stream.
type cappedBody struct {
	io.Reader
	io.Closer
}

// capBody returns a remote.Option that truncates the response body at
// maxBytes instead of rejecting it.
func capBody(maxBytes int64) remote.Option {

	// Truncation (vs remote's own MaxResponseSize, which errors on oversize)
	// keeps the <head> of oversized pages extractable.

	return remote.Option{
		AfterRequest: func(_ *remote.Transaction, response *http.Response) error {

			// RULE: Read one byte past the cap so the caller can tell "exactly
			// at the cap" from "truncated". newDocument trims the extra byte.
			response.Body = cappedBody{
				Reader: io.LimitReader(response.Body, maxBytes+1),
				Closer: response.Body,
			}
			return nil
		},
	}
}
