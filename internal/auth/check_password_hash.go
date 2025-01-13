package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("error comparing hashed password: %w", err)
	}
	return nil
}
