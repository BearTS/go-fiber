package services

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GenerateOtp() int {
	rand.Seed(time.Now().UnixNano())

	// Generate 4 random integers between 0 and 9
	var digits string
	for i := 0; i < 4; i++ {
		digits += strconv.Itoa(rand.Intn(10))
	}

	otp, _ := strconv.Atoi(digits)
	fmt.Println(otp)
	return otp
}
