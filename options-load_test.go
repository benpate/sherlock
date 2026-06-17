package sherlock

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig_Defaults(t *testing.T) {
	client := NewClient()
	config := client.newConfig()

	require.Equal(t, documentTypeUnknown, config.DocumentType)
	require.Equal(t, 6, config.MaximumRedirects)
	require.Equal(t, client.userAgent, config.UserAgent)
	require.NotNil(t, config.DefaultValue)
	require.Empty(t, config.DefaultValue)
	require.Empty(t, config.RemoteOptions)
}

func TestNewConfig_DocumentTypeOptions(t *testing.T) {
	client := NewClient()

	require.Equal(t, documentTypeActor, client.newConfig(AsActor()).DocumentType)
	require.Equal(t, documentTypeDocument, client.newConfig(AsDocument()).DocumentType)
	require.Equal(t, documentTypeCollection, client.newConfig(AsCollection()).DocumentType)
}

func TestNewConfig_WithMaximumRedirects(t *testing.T) {
	client := NewClient()
	config := client.newConfig(WithMaximumRedirects(2))
	require.Equal(t, 2, config.MaximumRedirects)
}

func TestNewConfig_WithDefaultValue(t *testing.T) {
	client := NewClient()
	defaultValue := map[string]any{"key": "value"}
	config := client.newConfig(WithDefaultValue(defaultValue))
	require.Equal(t, "value", config.DefaultValue["key"])
}

func TestNewConfig_WithDefaultValue_Nil(t *testing.T) {
	client := NewClient()

	// WithDefaultValue(nil) must NOT leave a nil map: the loaders write into it,
	// and a nil-map write would panic.
	config := client.newConfig(WithDefaultValue(nil))
	require.NotNil(t, config.DefaultValue)
	require.NotPanics(t, func() {
		config.DefaultValue["id"] = "https://example.com"
	})
}

func TestNewConfig_WithRemoteOptions(t *testing.T) {
	client := NewClient()

	// Two remote options should both be appended
	config := client.newConfig(WithRemoteOptions(AuthorizedFetch("", nil), AuthorizedFetch("", nil)))
	require.Len(t, config.RemoteOptions, 2)
}

func TestNewConfig_WithKeyPair(t *testing.T) {
	client := NewClient()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	config := client.newConfig(WithKeyPair("https://example.com/actor#key", privateKey))
	require.Len(t, config.RemoteOptions, 1)
}

func TestNewConfig_KeyPairFunc_Applied(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.Nil(t, err)

	// When the client has a keyPairFunc that returns a valid pair, AuthorizedFetch is added.
	client := NewClient(WithKeyPairFunc(func() (string, crypto.PrivateKey) {
		return "https://example.com/actor#key", privateKey
	}))

	config := client.newConfig()
	require.Len(t, config.RemoteOptions, 1)
}

func TestNewConfig_KeyPairFunc_EmptyResult(t *testing.T) {
	// When the keyPairFunc returns an empty pair, no remote option is added.
	client := NewClient(WithKeyPairFunc(func() (string, crypto.PrivateKey) {
		return "", nil
	}))

	config := client.newConfig()
	require.Empty(t, config.RemoteOptions)
}

func TestNewConfig_IgnoresNonOptionArgs(t *testing.T) {
	client := NewClient()

	// Load accepts ...any; non-Option arguments are silently ignored by newConfig.
	config := client.newConfig("a string", 42, AsActor())
	require.Equal(t, documentTypeActor, config.DocumentType)
}
