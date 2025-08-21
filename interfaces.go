package sherlock

import "crypto"

type KeyPairFunc func() (publicKeyID string, privateKey crypto.PrivateKey)
