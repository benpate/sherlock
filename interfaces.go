package sherlock

import "crypto"

// KeyPairFunc returns the public key ID and private key used to sign outbound
// Authorized Fetch requests.
type KeyPairFunc func() (publicKeyID string, privateKey crypto.PrivateKey)
