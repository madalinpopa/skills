package domain

import (
	"errors"
	"time"
	"uuid"
)

type TokenUUID uuid.UUID

func NewTokenUUID() TokenUUID {
	return TokenUUID(uuid.NewV7())
}

func (id TokenUUID) String() string {
	return uuid.UUID(id).String()
}

// RefreshToken is the stored side of a refresh credential. Only the SHA-256
// hash of the raw token is kept, so a database leak does not leak sessions.
type RefreshToken struct {
	UUID      TokenUUID
	UserUUID  UserUUID
	Hash      string
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt time.Time
}

func NewRefreshToken(userUUID UserUUID, hash string, createdAt, expiresAt time.Time) (*RefreshToken, error) {
	if userUUID.IsZero() {
		return nil, errors.New("refresh token user is empty")
	}

	if hash == "" {
		return nil, errors.New("refresh token hash is empty")
	}

	if !expiresAt.After(createdAt) {
		return nil, errors.New("refresh token must expire after it is created")
	}

	return &RefreshToken{
		UUID:      NewTokenUUID(),
		UserUUID:  userUUID,
		Hash:      hash,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}, nil
}

func (t *RefreshToken) IsExpired(now time.Time) bool {
	return !now.Before(t.ExpiresAt)
}

func (t *RefreshToken) IsRevoked() bool {
	return !t.RevokedAt.IsZero()
}

func (t *RefreshToken) IsActive(now time.Time) bool {
	return !t.IsExpired(now) && !t.IsRevoked()
}

func (t *RefreshToken) Revoke(now time.Time) {
	t.RevokedAt = now
}
