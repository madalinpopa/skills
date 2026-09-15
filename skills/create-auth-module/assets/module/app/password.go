package app

import (
	"fmt"

	"{{MODULE_PATH}}/internal/modules/auth/domain"

	"golang.org/x/crypto/bcrypt"
)

const (
	passwordCost = 12

	// maxPasswordLength is in bytes, matching the bcrypt input limit.
	maxPasswordLength = 72
	minPasswordLength = 8

	// dummyPasswordHash is compared against when the email is unknown, so a
	// login attempt takes the same time whether or not the account exists.
	dummyPasswordHash = "$2a$12$Dh57BatYqx.knOOhcM8aF.O4rFvA3tzEq8Xw3bgmfX4Vy/Chagc8q"
)

func HashPassword(password string) (string, error) {
	if len(password) < minPasswordLength {
		return "", domain.ErrPasswordTooShort
	}

	if len(password) > maxPasswordLength {
		return "", domain.ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return "", fmt.Errorf("hashing password: %w", err)
	}

	return string(hash), nil
}

func VerifyPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
