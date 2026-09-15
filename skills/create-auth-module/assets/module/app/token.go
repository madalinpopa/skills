package app

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"time"
	"uuid"

	"{{MODULE_PATH}}/internal/modules/auth/domain"
)

type AccessToken struct {
	Token     string
	ExpiresAt time.Time
}

type AccessTokenIssuer interface {
	Issue(ctx context.Context, userUUID domain.UserUUID) (AccessToken, error)
}

type VerifiedAccessToken struct {
	TokenID   uuid.UUID
	UserUUID  domain.UserUUID
	ExpiresAt time.Time
}

type AccessTokenVerifier interface {
	Verify(ctx context.Context, token string) (VerifiedAccessToken, error)
}

// RefreshSession pairs the raw token handed to the client with the hashed
// record kept in storage. The raw value never reaches persistence.
type RefreshSession struct {
	RawToken string
	Token    *domain.RefreshToken
}

type RefreshSessionIssuer interface {
	Issue(ctx context.Context, userUUID domain.UserUUID) (RefreshSession, error)
}

func HashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}
