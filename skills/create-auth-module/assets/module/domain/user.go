package domain

import (
	"errors"
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"
	"uuid"
)

const (
	maxNameLength  = 255
	maxEmailLength = 255
)

type UserUUID uuid.UUID

func NewUserUUID() UserUUID {
	return UserUUID(uuid.NewV7())
}

func ParseUserUUID(s string) (UserUUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UserUUID{}, err
	}

	return UserUUID(id), nil
}

func (id UserUUID) String() string {
	return uuid.UUID(id).String()
}

func (id UserUUID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil()
}

type User struct {
	UUID         UserUUID
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(name, email, passwordHash string, now time.Time) (*User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyName
	}

	if utf8.RuneCountInString(name) > maxNameLength {
		return nil, ErrNameTooLong
	}

	normalizedEmail, err := NormalizeEmail(email)
	if err != nil {
		return nil, err
	}

	if passwordHash == "" {
		return nil, errors.New("password hash is empty")
	}

	return &User{
		UUID:         NewUserUUID(),
		Name:         name,
		Email:        normalizedEmail,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// NormalizeEmail lowercases and trims the address and rejects anything that
// is not a bare address, so one mailbox maps to one stored value.
func NormalizeEmail(value string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(value))
	if email == "" || utf8.RuneCountInString(email) > maxEmailLength {
		return "", ErrInvalidEmail
	}

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return "", ErrInvalidEmail
	}

	return address.Address, nil
}
