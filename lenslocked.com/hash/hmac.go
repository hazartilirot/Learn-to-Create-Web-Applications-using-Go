package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

func NewHMAC(key string) HMAC {
	mac := hmac.New(sha256.New, []byte(key))
	return HMAC{
		hmac: mac,
	}
}

// Hash will hash an input string using HMAC with the secret key provided when the HMAC object was created
func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	b := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}

type HMAC struct {
	hmac hash.Hash
}
