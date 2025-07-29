package util

import (
	"math/rand"
	//"encoding/base64"
)

const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
const uniqueLooking = "ABCDEFGHJKLMNOPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
const chashing = "Il1,0oO,"

func GenerateRandomString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}
	return string(result)
}
