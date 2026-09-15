package tests

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	"{{MODULE_PATH}}/internal/modules/auth/adapters/rest/client"

	"github.com/stretchr/testify/require"
)

var userCounter atomic.Int64

type account struct {
	UUID     string
	Email    string
	Password string
}

func newClient(t *testing.T, editors ...client.RequestEditorFn) *client.ClientWithResponses {
	t.Helper()

	options := make([]client.ClientOption, 0, len(editors))
	for _, editor := range editors {
		options = append(options, client.WithRequestEditorFn(editor))
	}

	api, err := client.NewClientWithResponses(serverURL, options...)
	require.NoError(t, err)

	return api
}

func bearer(token string) client.RequestEditorFn {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

func mobile() *client.LoginParams {
	return &client.LoginParams{XClientType: new(client.ClientType("mobile"))}
}

func registerAccount(t *testing.T, api *client.ClientWithResponses) account {
	t.Helper()

	email := fmt.Sprintf("user-%d@example.com", userCounter.Add(1))
	password := "correct horse battery"

	created, err := api.RegisterUserWithResponse(t.Context(), client.RegisterUserJSONRequestBody{
		Name:     "Test User",
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, created.StatusCode(), string(created.Body))
	require.NotNil(t, created.JSON201)

	return account{UUID: created.JSON201.UserUuid, Email: email, Password: password}
}

// loginMobile logs in as a mobile client, so the refresh token comes back in
// the body and tests need no cookie jar.
func loginMobile(t *testing.T, api *client.ClientWithResponses, acc account) client.SessionResponse {
	t.Helper()

	session, err := api.LoginWithResponse(t.Context(), mobile(), client.LoginJSONRequestBody{
		Email:    acc.Email,
		Password: acc.Password,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, session.StatusCode(), string(session.Body))
	require.NotNil(t, session.JSON200)
	require.NotNil(t, session.JSON200.RefreshToken)

	return *session.JSON200
}

func refreshMobile(t *testing.T, api *client.ClientWithResponses, refreshToken string) *client.RefreshClientResponse {
	t.Helper()

	refreshed, err := api.RefreshWithResponse(
		t.Context(),
		&client.RefreshParams{XClientType: new(client.ClientType("mobile"))},
		client.RefreshJSONRequestBody{RefreshToken: refreshToken},
	)
	require.NoError(t, err)

	return refreshed
}
