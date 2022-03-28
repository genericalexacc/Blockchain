package main

import "crypto/sha256"

func getHash(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	return h.Sum(nil)
}
