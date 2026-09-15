package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"{{MODULE_PATH}}/internal/modules/auth/domain"
)

type Command struct {
	users   domain.UserRepository
	tokens  domain.TokenRepository
	access  AccessTokenIssuer
	refresh RefreshSessionIssuer
	clock   func() time.Time
}

type CommandArgs struct {
	Users   domain.UserRepository
	Tokens  domain.TokenRepository
	Access  AccessTokenIssuer
	Refresh RefreshSessionIssuer
	Clock   func() time.Time
}

func NewCommand(args CommandArgs) *Command {
	return &Command{
		users:   args.Users,
		tokens:  args.Tokens,
		access:  args.Access,
		refresh: args.Refresh,
		clock:   args.Clock,
	}
}

type Register struct {
	Name     string
	Email    string
	Password string
}

func (c *Command) Register(ctx context.Context, cmd Register) (*domain.User, error) {
	passwordHash, err := HashPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	user, err := domain.NewUser(cmd.Name, cmd.Email, passwordHash, c.clock().UTC())
	if err != nil {
		return nil, err
	}

	if err := c.users.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

type Login struct {
	Email    string
	Password string
}

// Session is what a successful login or refresh hands back: a short-lived
// access token and the raw refresh token that renews it.
type Session struct {
	UserUUID         domain.UserUUID
	AccessToken      AccessToken
	RefreshToken     string
	RefreshExpiresAt time.Time
}

func (c *Command) Login(ctx context.Context, cmd Login) (Session, error) {
	email, err := domain.NormalizeEmail(cmd.Email)
	if err != nil {
		return Session{}, domain.ErrInvalidCredentials
	}

	// An unknown email still runs the password check against a dummy hash, so
	// timing does not reveal which addresses are registered.
	passwordHash := dummyPasswordHash
	user, err := c.users.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return Session{}, fmt.Errorf("reading user by email: %w", err)
	}
	if err == nil {
		passwordHash = user.PasswordHash
	}

	if !VerifyPassword(cmd.Password, passwordHash) || user == nil {
		return Session{}, domain.ErrInvalidCredentials
	}

	session, err := c.issueSession(ctx, user.UUID)
	if err != nil {
		return Session{}, err
	}

	if err := c.tokens.Create(ctx, session.token); err != nil {
		return Session{}, fmt.Errorf("storing refresh token: %w", err)
	}

	return session.Session, nil
}

// Refresh exchanges an active refresh token for a new session and retires the
// presented token. A token presented after it was already rotated is the
// OAuth 2.0 replay signal: whoever rotated it first still holds a live
// replacement, so every session of that user is revoked. An expired but
// unrevoked token is ordinary staleness and is only refused.
func (c *Command) Refresh(ctx context.Context, rawToken string) (Session, error) {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return Session{}, domain.ErrInvalidRefreshToken
	}

	current, err := c.tokens.GetByHash(ctx, HashToken(rawToken))
	if errors.Is(err, domain.ErrRefreshTokenNotFound) {
		return Session{}, domain.ErrInvalidRefreshToken
	}
	if err != nil {
		return Session{}, fmt.Errorf("reading refresh token: %w", err)
	}

	now := c.clock().UTC()

	if current.IsRevoked() {
		if err := c.tokens.RevokeAllForUser(ctx, current.UserUUID, now); err != nil {
			return Session{}, fmt.Errorf("revoking sessions after token reuse: %w", err)
		}

		return Session{}, domain.ErrInvalidRefreshToken
	}

	if current.IsExpired(now) {
		return Session{}, domain.ErrInvalidRefreshToken
	}

	session, err := c.issueSession(ctx, current.UserUUID)
	if err != nil {
		return Session{}, err
	}

	current.Revoke(now)

	err = c.tokens.Rotate(ctx, current, session.token)
	if errors.Is(err, domain.ErrRefreshTokenRotated) {
		if err := c.tokens.RevokeAllForUser(ctx, current.UserUUID, now); err != nil {
			return Session{}, fmt.Errorf("revoking sessions after token reuse: %w", err)
		}

		return Session{}, domain.ErrInvalidRefreshToken
	}
	if err != nil {
		return Session{}, fmt.Errorf("rotating refresh token: %w", err)
	}

	return session.Session, nil
}

// Logout revokes the presented refresh session of the authenticated user.
// Revocation is idempotent: an unknown, expired, already revoked, or foreign
// token succeeds silently rather than revealing which case applied.
func (c *Command) Logout(ctx context.Context, userUUID domain.UserUUID, rawToken string) error {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return domain.ErrMissingRefreshToken
	}

	current, err := c.tokens.GetByHash(ctx, HashToken(rawToken))
	if errors.Is(err, domain.ErrRefreshTokenNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading refresh token: %w", err)
	}

	now := c.clock().UTC()

	if current.UserUUID != userUUID || !current.IsActive(now) {
		return nil
	}

	current.Revoke(now)

	if err := c.tokens.Revoke(ctx, current); err != nil {
		return fmt.Errorf("revoking refresh token: %w", err)
	}

	return nil
}

type issuedSession struct {
	Session
	token *domain.RefreshToken
}

func (c *Command) issueSession(ctx context.Context, userUUID domain.UserUUID) (issuedSession, error) {
	accessToken, err := c.access.Issue(ctx, userUUID)
	if err != nil {
		return issuedSession{}, fmt.Errorf("issuing access token: %w", err)
	}

	refresh, err := c.refresh.Issue(ctx, userUUID)
	if err != nil {
		return issuedSession{}, fmt.Errorf("issuing refresh token: %w", err)
	}

	return issuedSession{
		Session: Session{
			UserUUID:         userUUID,
			AccessToken:      accessToken,
			RefreshToken:     refresh.RawToken,
			RefreshExpiresAt: refresh.Token.ExpiresAt,
		},
		token: refresh.Token,
	}, nil
}
