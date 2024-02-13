package utils

import (
	"crypto/rand"
	mathRand "math/rand"
)

func GenerateRandomString(length int) string {
	// Define the character set from which to generate the random string
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Create a byte slice of the specified length
	randomBytes := make([]byte, length)

	// Fill the byte slice with random bytes
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// Convert the random bytes to characters from the character set
	for i := range randomBytes {
		randomBytes[i] = charSet[int(randomBytes[i])%len(charSet)]
	}

	// Convert the byte slice to a string and return it
	return string(randomBytes)
}

func GenerateRandomInt(maxVal int) int {
	return mathRand.Intn(10000)
}

func GenerateNewStringHavingSuffixName(mainString string, randomStringLen int, maxLength int) string {
	random_str := "_" + GenerateRandomString(randomStringLen-1)

	truncate_at := len(mainString)
	if truncate_at+len(random_str) > maxLength {
		truncate_at = maxLength - len(random_str)
	}
	
	return mainString[:truncate_at] + random_str
}
