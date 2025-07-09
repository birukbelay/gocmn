package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptCreateHash  password
func BcryptCreateHash(password string) (string, error) {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Hash password with Bcrypt's min cost
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func BcryptPasswordsMatch( PlainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(PlainPassword))
	return err == nil
}