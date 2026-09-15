package domain

import "errors"

var (
	ErrEmptyName            = errors.New("user name is empty")
	ErrNameTooLong          = errors.New("user name is too long")
	ErrInvalidEmail         = errors.New("email address is invalid")
	ErrPasswordTooShort     = errors.New("password is too short")
	ErrPasswordTooLong      = errors.New("password is too long")
	ErrEmailTaken           = errors.New("email address is already registered")
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrMissingRefreshToken  = errors.New("refresh token is missing")
	ErrInvalidRefreshToken  = errors.New("refresh token is invalid or expired")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRotated  = errors.New("refresh token was already rotated")
)
