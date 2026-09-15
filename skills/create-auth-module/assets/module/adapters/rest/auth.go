package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"{{MODULE_PATH}}/internal/modules/auth/adapters/rest/server"
	"{{MODULE_PATH}}/internal/modules/auth/app"
	"{{MODULE_PATH}}/internal/modules/auth/domain"
)

func (h Handler) Login(ctx context.Context, request server.LoginRequestObject) (server.LoginResponseObject, error) {
	session, err := h.command.Login(ctx, app.Login{
		Email:    request.Body.Email,
		Password: request.Body.Password,
	})
	if errors.Is(err, domain.ErrInvalidCredentials) {
		return server.Login401JSONResponse{
			Body: server.ErrorResponse{Slug: "invalid-credentials", Message: "The email or password is wrong"},
		}, nil
	}
	if err != nil {
		return nil, err
	}

	response := server.Login200JSONResponse{
		Body:    sessionBody(session),
		Headers: server.Login200ResponseHeaders{CacheControl: "no-store"},
	}

	if isMobile(request.Params.XClientType) {
		response.Body.RefreshToken = &session.RefreshToken
	} else {
		response.Headers.SetCookie = h.refreshCookie(session)
	}

	return response, nil
}

func (h Handler) Refresh(ctx context.Context, request server.RefreshRequestObject) (server.RefreshResponseObject, error) {
	rawToken := presentedRefreshToken(request.Params.XClientType, request.Params.RefreshToken, func() string {
		if request.Body == nil {
			return ""
		}
		return request.Body.RefreshToken
	})

	session, err := h.command.Refresh(ctx, rawToken)
	if errors.Is(err, domain.ErrInvalidRefreshToken) {
		return server.Refresh401JSONResponse{
			Body: server.ErrorResponse{Slug: "invalid-refresh-token", Message: "The refresh token is invalid or expired"},
		}, nil
	}
	if err != nil {
		return nil, err
	}

	response := server.Refresh200JSONResponse{
		Body:    sessionBody(session),
		Headers: server.Refresh200ResponseHeaders{CacheControl: "no-store"},
	}

	if isMobile(request.Params.XClientType) {
		response.Body.RefreshToken = &session.RefreshToken
	} else {
		response.Headers.SetCookie = h.refreshCookie(session)
	}

	return response, nil
}

func (h Handler) Logout(ctx context.Context, request server.LogoutRequestObject) (server.LogoutResponseObject, error) {
	identity, ok := app.IdentityFromContext(ctx)
	if !ok {
		return nil, errors.New("logout reached without an authenticated identity")
	}

	rawToken := presentedRefreshToken(request.Params.XClientType, request.Params.RefreshToken, func() string {
		if request.Body == nil {
			return ""
		}
		return request.Body.RefreshToken
	})

	err := h.command.Logout(ctx, identity.UserUUID, rawToken)
	if errors.Is(err, domain.ErrMissingRefreshToken) {
		return server.Logout400JSONResponse{Slug: "refresh-token-missing", Message: "A refresh token is required"}, nil
	}
	if err != nil {
		return nil, err
	}

	response := server.Logout204Response{}

	if !isMobile(request.Params.XClientType) {
		response.Headers.SetCookie = h.expiredRefreshCookie()
	}

	return response, nil
}

func isMobile(clientType *server.ClientType) bool {
	return clientType != nil && *clientType == mobileClient
}

// presentedRefreshToken picks the transport for the client type: the request
// body for mobile clients, the cookie for everyone else. Each client type has
// exactly one transport, so a stray cookie cannot act for a mobile client.
func presentedRefreshToken(clientType *server.ClientType, cookie *string, body func() string) string {
	if isMobile(clientType) {
		return body()
	}

	if cookie != nil {
		return *cookie
	}

	return ""
}

func sessionBody(session app.Session) server.SessionResponse {
	return server.SessionResponse{
		AccessToken: session.AccessToken.Token,
		TokenType:   "Bearer",
		ExpiresIn:   int(time.Until(session.AccessToken.ExpiresAt).Seconds()),
		UserUuid:    session.UserUUID.String(),
	}
}

func (h Handler) refreshCookie(session app.Session) *string {
	cookie := http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    session.RefreshToken,
		Path:     h.cookie.Path,
		Secure:   h.cookie.Secure,
		SameSite: h.cookie.SameSite,
		HttpOnly: true,
		Expires:  session.RefreshExpiresAt,
		MaxAge:   int(time.Until(session.RefreshExpiresAt).Seconds()),
	}

	return new(cookie.String())
}

// expiredRefreshCookie clears the cookie. The attributes must match the ones
// set at login, or the browser keeps the old cookie.
func (h Handler) expiredRefreshCookie() *string {
	cookie := http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     h.cookie.Path,
		Secure:   h.cookie.Secure,
		SameSite: h.cookie.SameSite,
		HttpOnly: true,
		MaxAge:   -1,
	}

	return new(cookie.String())
}
