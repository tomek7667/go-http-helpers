package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

func Sha512(s string) string {
	sha_512 := sha512.New()
	sha_512.Write([]byte(s))
	return hex.EncodeToString(sha_512.Sum(nil))
}

func NfoldSha512(s string, rounds int) string {
	hsh := s
	for i := 0; i < rounds; i++ {
		hsh = Sha512(hsh)
	}
	return hsh
}

func Sha512Salted(s, salt string) string {
	return Sha512(s + salt)
}
