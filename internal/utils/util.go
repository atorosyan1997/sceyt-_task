package utils

import (
	"crypto/rand"
	"math"
	"math/big"
)

const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"

// GenerateRandomString generate a string of random characters of given length
func GenerateRandomString(n int) string {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		ret[i] = letters[num.Int64()]
	}

	return string(ret)
}

// Round returns the nearest fractional number, side-rounded to roundOn
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var r float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		r = math.Ceil(digit)
	} else {
		r = math.Floor(digit)
	}
	newVal = r / pow
	return
}
