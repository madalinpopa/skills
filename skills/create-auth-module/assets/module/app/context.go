package app

import (
	"context"
	"uuid"

	"{{MODULE_PATH}}/internal/modules/auth/domain"
)

type identityKey struct{}

// Identity is the caller a verified access token stands for.
type Identity struct {
	TokenID  uuid.UUID
	UserUUID domain.UserUUID
}

func ContextWithIdentity(ctx context.Context, identity Identity) context.Context {
	return context.WithValue(ctx, identityKey{}, identity)
}

func IdentityFromContext(ctx context.Context) (Identity, bool) {
	identity, ok := ctx.Value(identityKey{}).(Identity)
	return identity, ok
}
