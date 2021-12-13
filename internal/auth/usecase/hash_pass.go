package usecase

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash password then return the hashed password as a base64 encoded string
func hashPassword(password string) (string, error) {

	var passwordBytes = []byte(password)

	// Hash password with bcrypt default cost
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)

	return string(hashedPasswordBytes), err
}

// Check if two passwords match using Bcrypt's CompareHashAndPassword
// which return nil on success and an error on failure.
func doPasswordsMatch(hashedPassword, currentPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
