package utils

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/rs/zerolog/log"
)

// GenerateOtp generates a random n digit otp
func GenerateOTP(digits int) (string, error) {
	upper := math.Pow10(digits)
	val, err := rand.Int(rand.Reader, big.NewInt(int64(upper)))
	if err != nil {
		log.Err(err).Msg("Error generating otp")
		return "", err
	}
	// adds a variable zero-padding to the left to ensure otp is uniformly random
	expr := "%0" + strconv.Itoa(digits) + "v"
	otp := fmt.Sprintf(expr, val.String())
	return otp, nil
}
