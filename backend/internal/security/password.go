package security

import (
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

// GenerateSalt generates a random salt of the given length.
func GenerateSalt(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(bytes), nil
}

// HashPassword hashes the password combined with the salt using bcrypt.
func HashPassword(password, salt string) (string, error) {
	combined := password + salt
	hash, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	return string(hash), err
}

// VerifyPassword compares the given password and salt with the hashed password.
func VerifyPassword(password, salt, hashed string) bool {
	combined := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(combined))
	return err == nil
}
