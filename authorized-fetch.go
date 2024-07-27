package sherlock

import (
	"crypto"
	"net/http"

	"github.com/benpate/derp"
	"github.com/benpate/hannibal/sigs"
	"github.com/benpate/remote"
	"github.com/rs/zerolog/log"
)

// AuthorizedFetch is a remote.Option that signs all outbound requests according to the
// ActivityPub "Authorized Fetch" convention: https://funfedi.dev/testing_tools/http_signatures/
func AuthorizedFetch(publicKeyID string, privateKey crypto.PrivateKey) remote.Option {

	if publicKeyID == "" {
		log.Info().Msg("AuthorizedFetch: No publicKeyID provided")
		return remote.Option{}
	}

	if privateKey == nil {
		log.Info().Msg("AuthorizedFetch: No privateKey provided")
		return remote.Option{}
	}

	return remote.Option{

		// ModifyRequest is called after an http.Request has been generated, but before it is sent to the
		// remote server. It can be used to modify the request, or to replace it entirely.
		// If it returns a non-nil http.Response, then that is used INSTEAD OF calling the remote server.
		// If it returns a nil http.Response, then the request is sent to the remote server as normal.
		ModifyRequest: func(t *remote.Transaction, request *http.Request) *http.Response {

			signer := sigs.NewSigner(
				publicKeyID,
				privateKey,
				sigs.SignerFields("(request-target)", "host", "date"),
			)

			if err := signer.Sign(request); err != nil {
				derp.Report(derp.Wrap(err, "sherlock.AuthorizedFetch", "Error signing request"))
			}

			return nil
		},
	}
}
