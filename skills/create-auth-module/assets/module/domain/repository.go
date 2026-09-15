package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByUUID(ctx context.Context, id UserUUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type TokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	GetByHash(ctx context.Context, hash string) (*RefreshToken, error)
	Revoke(ctx context.Context, token *RefreshToken) error
	// Rotate revokes current and stores replacement in one transaction. It
	// returns ErrRefreshTokenRotated when current was no longer active, which
	// means another request spent the same token first.
	Rotate(ctx context.Context, current, replacement *RefreshToken) error
	RevokeAllForUser(ctx context.Context, userUUID UserUUID, now time.Time) error
}
