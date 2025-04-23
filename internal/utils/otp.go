package utils

import (
	"math/rand"
	"time"
)

func GenerateOTP(length int) string {
	charset := "0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)

	otp := make([]byte, length)

	for i := range otp {
		otp[i] = charset[r.Intn(len(charset))]
	}

	return string(otp)
}

func GetExpiryTime() int64 {
	return time.Now().Add(10 * time.Minute).Unix()
}
