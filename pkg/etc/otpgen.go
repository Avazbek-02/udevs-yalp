package etc

import (
	"math/rand"
	"time"
)

func GenerateOTP(length int) string {
	const charset = "0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(otp)
}
