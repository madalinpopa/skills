package app

import (
	"context"
	"sync"
	"testing"
	"time"

	"{{MODULE_PATH}}/internal/modules/auth/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testNow = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

type memoryTokens struct {
	mu     sync.Mutex
	byHash map[string]*domain.RefreshToken
}

func newMemoryTokens() *memoryTokens {
	return &memoryTokens{byHash: map[string]*domain.RefreshToken{}}
}

func (m *memoryTokens) Create(_ context.Context, token *domain.RefreshToken) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.byHash[token.Hash] = token
	return nil
}

func (m *memoryTokens) GetByHash(_ context.Context, hash string) (*domain.RefreshToken, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	token, ok := m.byHash[hash]
	if !ok {
		return nil, domain.ErrRefreshTokenNotFound
	}

	copied := *token
	return &copied, nil
}

func (m *memoryTokens) Revoke(_ context.Context, token *domain.RefreshToken) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.byHash[token.Hash].RevokedAt = token.RevokedAt
	return nil
}

func (m *memoryTokens) Rotate(_ context.Context, current, replacement *domain.RefreshToken) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	stored := m.byHash[current.Hash]
	if stored.IsRevoked() {
		return domain.ErrRefreshTokenRotated
	}

	stored.RevokedAt = current.RevokedAt
	m.byHash[replacement.Hash] = replacement
	return nil
}

func (m *memoryTokens) RevokeAllForUser(_ context.Context, userUUID domain.UserUUID, now time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, token := range m.byHash {
		if token.UserUUID == userUUID && !token.IsRevoked() {
			token.RevokedAt = now
		}
	}
	return nil
}

func (m *memoryTokens) active(hash string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.byHash[hash].IsActive(testNow)
}

type stubAccess struct{}

func (stubAccess) Issue(_ context.Context, userUUID domain.UserUUID) (AccessToken, error) {
	return AccessToken{Token: "access-" + userUUID.String(), ExpiresAt: testNow.Add(time.Minute)}, nil
}

type stubRefresh struct {
	mu   sync.Mutex
	next int
}

func (s *stubRefresh) Issue(_ context.Context, userUUID domain.UserUUID) (RefreshSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.next++
	raw := "refresh-" + string(rune('a'+s.next))

	token, err := domain.NewRefreshToken(userUUID, HashToken(raw), testNow, testNow.Add(time.Hour))
	if err != nil {
		return RefreshSession{}, err
	}

	return RefreshSession{RawToken: raw, Token: token}, nil
}

func newRefreshCommand(t *testing.T) (*Command, *memoryTokens) {
	t.Helper()

	tokens := newMemoryTokens()
	command := NewCommand(CommandArgs{
		Tokens:  tokens,
		Access:  stubAccess{},
		Refresh: &stubRefresh{},
		Clock:   func() time.Time { return testNow },
	})

	return command, tokens
}

func seedToken(t *testing.T, tokens *memoryTokens, userUUID domain.UserUUID, raw string, expiresAt time.Time) *domain.RefreshToken {
	t.Helper()

	token, err := domain.NewRefreshToken(userUUID, HashToken(raw), testNow.Add(-time.Hour), expiresAt)
	require.NoError(t, err)
	require.NoError(t, tokens.Create(t.Context(), token))

	return token
}

func TestCommandRefresh_rotatesActiveToken(t *testing.T) {
	t.Parallel()

	// Arrange
	command, tokens := newRefreshCommand(t)
	user := domain.NewUserUUID()
	seedToken(t, tokens, user, "old", testNow.Add(time.Hour))

	// Act
	session, err := command.Refresh(t.Context(), "old")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, user, session.UserUUID)
	assert.NotEmpty(t, session.RefreshToken)
	assert.False(t, tokens.active(HashToken("old")), "presented token must be revoked")
	assert.True(t, tokens.active(HashToken(session.RefreshToken)), "replacement must be active")
}

func TestCommandRefresh_refused(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		expiresAt        time.Time
		revoked          bool
		otherStaysActive bool
	}{
		"expired token is refused without touching other sessions": {
			expiresAt:        testNow.Add(-time.Minute),
			otherStaysActive: true,
		},
		"replayed token revokes every session of the user": {
			expiresAt: testNow.Add(time.Hour),
			revoked:   true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			command, tokens := newRefreshCommand(t)
			user := domain.NewUserUUID()
			presented := seedToken(t, tokens, user, "presented", tc.expiresAt)
			seedToken(t, tokens, user, "other", testNow.Add(time.Hour))
			if tc.revoked {
				presented.Revoke(testNow.Add(-time.Minute))
				require.NoError(t, tokens.Revoke(t.Context(), presented))
			}

			// Act
			_, err := command.Refresh(t.Context(), "presented")

			// Assert
			require.ErrorIs(t, err, domain.ErrInvalidRefreshToken)
			assert.Equal(t, tc.otherStaysActive, tokens.active(HashToken("other")))
		})
	}
}

func TestCommandRefresh_unknownToken(t *testing.T) {
	t.Parallel()

	// Arrange
	command, _ := newRefreshCommand(t)

	// Act
	_, err := command.Refresh(t.Context(), "never-issued")

	// Assert
	require.ErrorIs(t, err, domain.ErrInvalidRefreshToken)
}

func TestCommandLogout_revokesOnlyOwnActiveToken(t *testing.T) {
	t.Parallel()

	// Arrange
	command, tokens := newRefreshCommand(t)
	owner := domain.NewUserUUID()
	other := domain.NewUserUUID()
	seedToken(t, tokens, owner, "own", testNow.Add(time.Hour))
	seedToken(t, tokens, other, "foreign", testNow.Add(time.Hour))

	// Act
	ownErr := command.Logout(t.Context(), owner, "own")
	foreignErr := command.Logout(t.Context(), owner, "foreign")
	unknownErr := command.Logout(t.Context(), owner, "unknown")
	missingErr := command.Logout(t.Context(), owner, " ")

	// Assert
	require.NoError(t, ownErr)
	require.NoError(t, foreignErr)
	require.NoError(t, unknownErr)
	require.ErrorIs(t, missingErr, domain.ErrMissingRefreshToken)
	assert.False(t, tokens.active(HashToken("own")))
	assert.True(t, tokens.active(HashToken("foreign")), "another user's token must stay active")
}
