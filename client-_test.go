package sherlock

import (
	"crypto"
	"testing"

	"github.com/benpate/hannibal/streams"
	"github.com/stretchr/testify/require"
)

func TestNewClient_Defaults(t *testing.T) {
	client := NewClient()
	require.Equal(t, "Sherlock (https://github.com/benpate/sherlock)", client.userAgent)
	require.Nil(t, client.keyPairFunc)
}

func TestNewClient_WithUserAgent(t *testing.T) {
	client := NewClient(WithUserAgent("MyAgent/1.0"))
	require.Equal(t, "MyAgent/1.0", client.userAgent)
}

func TestNewClient_WithKeyPairFunc(t *testing.T) {
	called := false
	fn := func() (string, crypto.PrivateKey) {
		called = true
		return "", nil
	}

	client := NewClient(WithKeyPairFunc(fn))
	require.NotNil(t, client.keyPairFunc)
	require.False(t, called, "the func should be stored, not invoked, at construction time")
}

func TestClient_With(t *testing.T) {
	client := NewClient()
	client.With(WithUserAgent("Updated/2.0"))
	require.Equal(t, "Updated/2.0", client.userAgent)
}

func TestClient_Load_EmptyURL(t *testing.T) {
	client := NewClient()
	doc, err := client.Load("")
	require.NotNil(t, err)
	require.True(t, doc.IsNil())
}

func TestClient_Load_NegativeRedirects(t *testing.T) {
	client := NewClient()
	doc, err := client.Load("https://example.com", WithMaximumRedirects(-1))
	require.NotNil(t, err)
	require.True(t, doc.IsNil())
}

func TestClient_NoopMethods(t *testing.T) {
	client := NewClient()

	// Save and Delete are intentional no-ops that never error
	require.Nil(t, client.Save(streams.NilDocument()))
	require.Nil(t, client.Delete("any-id"))

	// SetRootClient is a no-op; just confirm it does not panic
	require.NotPanics(t, func() {
		client.SetRootClient(client)
	})
}
