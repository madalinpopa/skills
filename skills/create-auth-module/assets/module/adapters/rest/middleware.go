package rest

import (
	"net/http"
	"strings"

	"{{MODULE_PATH}}/internal/modules/auth/adapters/rest/server"
	"{{MODULE_PATH}}/internal/modules/auth/app"

	"github.com/labstack/echo/v5"
)

// publicOperations lists the operations that need no bearer token. Every
// other operation is protected, so a new route is protected by default.
var publicOperations = map[string]bool{
	"RegisterUser": true,
	"Login":        true,
	"Refresh":      true,
}

// NewBearerMiddleware verifies the bearer token of protected operations and
// puts the resulting identity on the request context.
func NewBearerMiddleware(verifier app.AccessTokenVerifier) server.StrictMiddlewareFunc {
	return func(next server.StrictHandlerFunc, operationID string) server.StrictHandlerFunc {
		return func(ectx *echo.Context, request any) (any, error) {
			if publicOperations[operationID] {
				return next(ectx, request)
			}

			token, ok := bearerToken(ectx.Request())
			if !ok {
				return nil, unauthorized(ectx, "bearer-token-missing", "A bearer access token is required")
			}

			verified, err := verifier.Verify(ectx.Request().Context(), token)
			if err != nil {
				return nil, unauthorized(ectx, "invalid-access-token", "The access token is invalid or expired")
			}

			identity := app.Identity{TokenID: verified.TokenID, UserUUID: verified.UserUUID}
			ectx.SetRequest(ectx.Request().WithContext(app.ContextWithIdentity(ectx.Request().Context(), identity)))

			return next(ectx, request)
		}
	}
}

func bearerToken(r *http.Request) (string, bool) {
	scheme, token, found := strings.Cut(r.Header.Get("Authorization"), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(token) == "" {
		return "", false
	}

	return strings.TrimSpace(token), true
}

// unauthorized writes the 401 itself so the body matches the ErrorResponse
// schema. The strict handler treats a nil response with a nil error as done.
func unauthorized(ectx *echo.Context, slug, message string) error {
	ectx.Response().Header().Set("WWW-Authenticate", "Bearer")
	return ectx.JSON(http.StatusUnauthorized, server.ErrorResponse{Slug: slug, Message: message})
}
