package refreshtx

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"{{MODULE_PATH}}/internal/modules/auth/app"
	"{{MODULE_PATH}}/internal/modules/auth/domain"
)

const tokenBytes = 32

// Service issues opaque refresh tokens: random bytes for the client and a
// hash for storage. The random source and clock are injected for tests.
type Service struct {
	ttl    time.Duration
	random io.Reader
	clock  func() time.Time
}

func NewService(ttl time.Duration, random io.Reader, clock func() time.Time) *Service {
	return &Service{ttl: ttl, random: random, clock: clock}
}

func (s *Service) Issue(_ context.Context, userUUID domain.UserUUID) (app.RefreshSession, error) {
	raw := make([]byte, tokenBytes)
	if _, err := io.ReadFull(s.random, raw); err != nil {
		return app.RefreshSession{}, fmt.Errorf("generating refresh token: %w", err)
	}

	rawToken := base64.RawURLEncoding.EncodeToString(raw)
	now := s.clock().UTC().Truncate(time.Second)

	token, err := domain.NewRefreshToken(userUUID, app.HashToken(rawToken), now, now.Add(s.ttl))
	if err != nil {
		return app.RefreshSession{}, err
	}

	return app.RefreshSession{RawToken: rawToken, Token: token}, nil
}
