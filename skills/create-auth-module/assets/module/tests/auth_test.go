package tests

import (
	"context"
	"net/http"
	"testing"

	"{{MODULE_PATH}}/internal/modules/auth/adapters/rest/client"
	"{{MODULE_PATH}}/internal/testkit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterLoginAndGetMe(t *testing.T) {
	testkit.RequireIntegration(t)

	// Arrange
	api := newClient(t)
	acc := registerAccount(t, api)

	// Act
	session := loginMobile(t, api, acc)
	me, err := newClient(t, bearer(session.AccessToken)).GetMeWithResponse(t.Context())
	require.NoError(t, err)

	// Assert
	assert.Equal(t, "Bearer", session.TokenType)
	assert.Equal(t, acc.UUID, session.UserUuid)
	assert.Positive(t, session.ExpiresIn)
	require.Equal(t, http.StatusOK, me.StatusCode(), string(me.Body))
	require.NotNil(t, me.JSON200)
	assert.Equal(t, acc.UUID, me.JSON200.UserUuid)
	assert.Equal(t, acc.Email, me.JSON200.Email)
}

func TestLogin_webClientGetsCookie(t *testing.T) {
	testkit.RequireIntegration(t)

	// Arrange
	api := newClient(t)
	acc := registerAccount(t, api)

	// Act
	session, err := api.LoginWithResponse(t.Context(), nil, client.LoginJSONRequestBody{
		Email:    acc.Email,
		Password: acc.Password,
	})
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusOK, session.StatusCode(), string(session.Body))
	require.NotNil(t, session.JSON200)
	assert.Nil(t, session.JSON200.RefreshToken, "web clients get no refresh token in the body")
	assert.Equal(t, "no-store", session.HTTPResponse.Header.Get("Cache-Control"))

	cookie := session.HTTPResponse.Header.Get("Set-Cookie")
	assert.Contains(t, cookie, "refresh_token=")
	assert.Contains(t, cookie, "HttpOnly")
	assert.Contains(t, cookie, "Path=/auth")
}

func TestLogin_rejected(t *testing.T) {
	testkit.RequireIntegration(t)

	api := newClient(t)
	acc := registerAccount(t, api)

	tests := map[string]client.LoginJSONRequestBody{
		"wrong password": {Email: acc.Email, Password: "not the password"},
		"unknown email":  {Email: "nobody@example.com", Password: acc.Password},
	}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			session, err := api.LoginWithResponse(t.Context(), nil, body)
			require.NoError(t, err)

			// Assert
			require.Equal(t, http.StatusUnauthorized, session.StatusCode())
			require.NotNil(t, session.JSON401)
			assert.Equal(t, "invalid-credentials", session.JSON401.Slug)
		})
	}
}

func TestRegister_duplicateEmail(t *testing.T) {
	testkit.RequireIntegration(t)

	// Arrange
	api := newClient(t)
	acc := registerAccount(t, api)

	// Act
	again, err := api.RegisterUserWithResponse(t.Context(), client.RegisterUserJSONRequestBody{
		Name:     "Someone Else",
		Email:    acc.Email,
		Password: acc.Password,
	})
	require.NoError(t, err)

	// Assert
	require.Equal(t, http.StatusConflict, again.StatusCode())
	require.NotNil(t, again.JSON409)
	assert.Equal(t, "email-taken", again.JSON409.Slug)
}

func TestRefresh_rotatesAndDetectsReplay(t *testing.T) {
	testkit.RequireIntegration(t)

	// Arrange
	api := newClient(t)
	acc := registerAccount(t, api)
	first := loginMobile(t, api, acc)

	// Act
	second := refreshMobile(t, api, *first.RefreshToken)
	replay := refreshMobile(t, api, *first.RefreshToken)

	// Assert
	require.Equal(t, http.StatusOK, second.StatusCode(), string(second.Body))
	require.NotNil(t, second.JSON200)
	require.NotNil(t, second.JSON200.RefreshToken)
	assert.NotEqual(t, *first.RefreshToken, *second.JSON200.RefreshToken)

	require.Equal(t, http.StatusUnauthorized, replay.StatusCode())
	require.NotNil(t, replay.JSON401)
	assert.Equal(t, "invalid-refresh-token", replay.JSON401.Slug)

	// The replay ended every session, so the token from the second login is
	// gone too.
	afterReplay := refreshMobile(t, api, *second.JSON200.RefreshToken)
	assert.Equal(t, http.StatusUnauthorized, afterReplay.StatusCode())
}

func TestLogout_revokesRefreshToken(t *testing.T) {
	testkit.RequireIntegration(t)

	// Arrange
	api := newClient(t)
	acc := registerAccount(t, api)
	session := loginMobile(t, api, acc)
	authed := newClient(t, bearer(session.AccessToken))

	// Act
	loggedOut, err := authed.LogoutWithResponse(
		t.Context(),
		&client.LogoutParams{XClientType: new(client.ClientType("mobile"))},
		client.LogoutJSONRequestBody{RefreshToken: *session.RefreshToken},
	)
	require.NoError(t, err)
	refreshed := refreshMobile(t, api, *session.RefreshToken)

	// Assert
	require.Equal(t, http.StatusNoContent, loggedOut.StatusCode(), string(loggedOut.Body))
	assert.Equal(t, http.StatusUnauthorized, refreshed.StatusCode())
}

func TestGetMe_unauthorized(t *testing.T) {
	testkit.RequireIntegration(t)

	tests := map[string]client.RequestEditorFn{
		"missing header": func(_ context.Context, _ *http.Request) error { return nil },
		"garbage token":  bearer("not-a-jwt"),
	}

	for name, editor := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			me, err := newClient(t, editor).GetMeWithResponse(t.Context())
			require.NoError(t, err)

			// Assert
			require.Equal(t, http.StatusUnauthorized, me.StatusCode())
			require.NotNil(t, me.JSON401)
			assert.NotEmpty(t, me.JSON401.Slug)
			assert.Equal(t, "Bearer", me.HTTPResponse.Header.Get("WWW-Authenticate"))
		})
	}
}
