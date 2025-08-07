package sherlock

import "crypto"

type KeyPairFunc func() (string, crypto.PrivateKey)
