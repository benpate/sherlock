package activitypub

import "crypto"

// KeyPairFunc returns the public key ID and private key used to sign
// "Authorized Fetch" requests, resolved lazily only when a request needs them.
type KeyPairFunc func() (publicKeyID string, privateKey crypto.PrivateKey)
