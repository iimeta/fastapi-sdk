package util

import (
	"crypto/hmac"
	"crypto/sha256"
)

func HMACSHA256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}
