package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pass []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		log.Printf("Failed to generate password: %v", err)
		return ""
	}

	return string(hashed)
}

// TestPassword checks if the provided password matches the hashed password.
func TestPassword(password []byte, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
	return err == nil
}
