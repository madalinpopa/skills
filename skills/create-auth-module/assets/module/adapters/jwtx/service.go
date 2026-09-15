package jwtx

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"
	"uuid"

	"{{MODULE_PATH}}/internal/modules/auth/app"
	"{{MODULE_PATH}}/internal/modules/auth/domain"

	"github.com/golang-jwt/jwt/v5"
)

const signingMethod = "HS256"

type Config struct {
	SigningKey []byte
	Issuer     string
	Audience   string
	TTL        time.Duration
}

// Service issues and verifies HS256 access tokens. The clock is injected so
// tests can mint expired tokens without waiting.
type Service struct {
	config Config
	clock  func() time.Time
}

func NewService(config Config, clock func() time.Time) *Service {
	config.SigningKey = slices.Clone(config.SigningKey)

	return &Service{config: config, clock: clock}
}

func (s *Service) Issue(_ context.Context, userUUID domain.UserUUID) (app.AccessToken, error) {
	now := s.clock().UTC().Truncate(time.Second)
	expiresAt := now.Add(s.config.TTL)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        uuid.NewV7().String(),
		Issuer:    s.config.Issuer,
		Audience:  jwt.ClaimStrings{s.config.Audience},
		Subject:   userUUID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})

	signed, err := token.SignedString(s.config.SigningKey)
	if err != nil {
		return app.AccessToken{}, fmt.Errorf("signing access token: %w", err)
	}

	return app.AccessToken{Token: signed, ExpiresAt: expiresAt}, nil
}

func (s *Service) Verify(_ context.Context, token string) (app.VerifiedAccessToken, error) {
	var claims jwt.RegisteredClaims

	parsed, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(*jwt.Token) (any, error) { return s.config.SigningKey, nil },
		jwt.WithTimeFunc(s.clock),
		jwt.WithValidMethods([]string{signingMethod}),
		jwt.WithIssuer(s.config.Issuer),
		jwt.WithAudience(s.config.Audience),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
		jwt.WithNotBeforeRequired(),
	)
	if err != nil {
		return app.VerifiedAccessToken{}, fmt.Errorf("parsing access token: %w", err)
	}

	if !parsed.Valid {
		return app.VerifiedAccessToken{}, errors.New("access token is not valid")
	}

	tokenID, err := uuid.Parse(claims.ID)
	if err != nil {
		return app.VerifiedAccessToken{}, fmt.Errorf("parsing access token id: %w", err)
	}

	userUUID, err := domain.ParseUserUUID(claims.Subject)
	if err != nil {
		return app.VerifiedAccessToken{}, fmt.Errorf("parsing access token subject: %w", err)
	}

	if userUUID.IsZero() {
		return app.VerifiedAccessToken{}, errors.New("access token subject is empty")
	}

	return app.VerifiedAccessToken{
		TokenID:   tokenID,
		UserUUID:  userUUID,
		ExpiresAt: claims.ExpiresAt.Time.UTC(),
	}, nil
}
