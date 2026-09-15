package jwtx

import (
	"testing"
	"time"

	"{{MODULE_PATH}}/internal/modules/auth/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testNow = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func testConfig() Config {
	return Config{
		SigningKey: []byte("test-signing-key-of-thirty-two-b"),
		Issuer:     "test-api",
		Audience:   "test-api",
		TTL:        15 * time.Minute,
	}
}

func TestServiceVerify_success(t *testing.T) {
	t.Parallel()

	// Arrange
	service := NewService(testConfig(), func() time.Time { return testNow })
	user := domain.NewUserUUID()

	// Act
	issued, issueErr := service.Issue(t.Context(), user)
	verified, verifyErr := service.Verify(t.Context(), issued.Token)

	// Assert
	require.NoError(t, issueErr)
	require.NoError(t, verifyErr)
	assert.Equal(t, user, verified.UserUUID)
	assert.Equal(t, testNow.Add(15*time.Minute), verified.ExpiresAt)
	assert.NotEmpty(t, verified.TokenID)
}

func TestServiceVerify_failed(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		issuer   Config
		verifier Config
		verifyAt time.Time
		tamper   func(string) string
	}{
		"expired token": {
			issuer:   testConfig(),
			verifier: testConfig(),
			verifyAt: testNow.Add(16 * time.Minute),
		},
		"token signed with another key": {
			issuer: Config{
				SigningKey: []byte("another-signing-key-of-32-bytes!"),
				Issuer:     "test-api",
				Audience:   "test-api",
				TTL:        15 * time.Minute,
			},
			verifier: testConfig(),
			verifyAt: testNow,
		},
		"token for another audience": {
			issuer: Config{
				SigningKey: []byte("test-signing-key-of-thirty-two-b"),
				Issuer:     "test-api",
				Audience:   "other-api",
				TTL:        15 * time.Minute,
			},
			verifier: testConfig(),
			verifyAt: testNow,
		},
		"tampered token": {
			issuer:   testConfig(),
			verifier: testConfig(),
			verifyAt: testNow,
			tamper:   func(token string) string { return token[:len(token)-2] + "xx" },
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			issuer := NewService(tc.issuer, func() time.Time { return testNow })
			verifier := NewService(tc.verifier, func() time.Time { return tc.verifyAt })

			issued, err := issuer.Issue(t.Context(), domain.NewUserUUID())
			require.NoError(t, err)

			token := issued.Token
			if tc.tamper != nil {
				token = tc.tamper(token)
			}

			// Act
			_, err = verifier.Verify(t.Context(), token)

			// Assert
			require.Error(t, err)
		})
	}
}
