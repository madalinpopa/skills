package auth

import (
	"context"
	"uuid"

	"{{MODULE_PATH}}/internal/modules/auth/app"
)

// Contract is what other modules see of auth. It hands out plain UUIDs, so a
// consumer stores the caller's identifier without importing auth's domain.
type Contract struct {
	verifier app.AccessTokenVerifier
}

func NewContract(verifier app.AccessTokenVerifier) *Contract {
	return &Contract{verifier: verifier}
}

func (c *Contract) VerifyAccessToken(ctx context.Context, token string) (uuid.UUID, error) {
	verified, err := c.verifier.Verify(ctx, token)
	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.UUID(verified.UserUUID), nil
}
