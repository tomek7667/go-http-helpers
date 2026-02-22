package utils

import (
	"math/rand"
	"time"
)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func GetRandom(alphabet string, length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = alphabet[seededRand.Intn(len(alphabet))]
	}
	return string(b)
}

func GetRandomHex(bytesLength int) string {
	return GetRandom("0123456789abcdef", bytesLength*2)
}

func GetRandomUppercase(length int) string {
	return GetRandom("ABCDEFGHIJKLMNOPRSTUVWXYZ", length)
}

func GetRandomLowercase(length int) string {
	return GetRandom("abcdefghijklmnoprstuvwxyz", length)
}

func GetRandomDigits(length int) string {
	return GetRandom("0123456789", length)
}
