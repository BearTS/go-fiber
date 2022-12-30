package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
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

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
