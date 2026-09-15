package refreshtx

import (
	"crypto/rand"
	"testing"
	"time"

	"{{MODULE_PATH}}/internal/modules/auth/app"
	"{{MODULE_PATH}}/internal/modules/auth/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceIssue(t *testing.T) {
	t.Parallel()

	// Arrange
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	service := NewService(time.Hour, rand.Reader, func() time.Time { return now })
	user := domain.NewUserUUID()

	// Act
	first, firstErr := service.Issue(t.Context(), user)
	second, secondErr := service.Issue(t.Context(), user)

	// Assert
	require.NoError(t, firstErr)
	require.NoError(t, secondErr)
	assert.Len(t, first.RawToken, 43, "32 random bytes encode to 43 url-safe characters")
	assert.NotEqual(t, first.RawToken, second.RawToken)
	assert.Equal(t, app.HashToken(first.RawToken), first.Token.Hash)
	assert.Equal(t, user, first.Token.UserUUID)
	assert.Equal(t, now, first.Token.CreatedAt)
	assert.Equal(t, now.Add(time.Hour), first.Token.ExpiresAt)
	assert.True(t, first.Token.IsActive(now))
}
