package sherlock

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"testing"

	"github.com/benpate/remote"
	"github.com/stretchr/testify/require"
)

func TestAuthorizedFetch_MissingPublicKey(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	// No publicKeyID -> empty option (no ModifyRequest hook)
	option := AuthorizedFetch("", privateKey)
	require.Nil(t, option.ModifyRequest)
}

func TestAuthorizedFetch_MissingPrivateKey(t *testing.T) {
	// No privateKey -> empty option (no ModifyRequest hook)
	option := AuthorizedFetch("https://example.com/actor#key", nil)
	require.Nil(t, option.ModifyRequest)
}

func TestAuthorizedFetch_Valid(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	option := AuthorizedFetch("https://example.com/actor#key", privateKey)
	require.NotNil(t, option.ModifyRequest)

	// The hook should sign the request and return nil (meaning "send it normally").
	request, err := http.NewRequest(http.MethodGet, "https://example.com/inbox", nil)
	require.Nil(t, err)

	txn := remote.Get("https://example.com/inbox")
	response := option.ModifyRequest(txn, request)
	require.Nil(t, response)

	// A Signature header should have been added to the outbound request.
	require.NotEmpty(t, request.Header.Get("Signature"))
}
