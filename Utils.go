package main

import "encoding/hex"

func ByteArrToHash(value []byte, difficulty int) string {
	return hex.EncodeToString(getHash(value))[:difficulty]
}
