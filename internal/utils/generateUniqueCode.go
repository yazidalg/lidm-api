package utils

import (
	"math/rand"
	"time"
	"strconv"
)

func GenerateInviteCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// UintToString converts uint to string (small helper for socket room naming / logs)
func UintToString(v uint) string { return strconv.FormatUint(uint64(v), 10) }
