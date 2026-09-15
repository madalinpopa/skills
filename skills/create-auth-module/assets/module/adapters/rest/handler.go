package rest

//go:generate go tool oapi-codegen --config=oapi-codegen.yaml openapi.yaml

import (
	"net/http"

	"{{MODULE_PATH}}/internal/modules/auth/adapters/rest/server"
	"{{MODULE_PATH}}/internal/modules/auth/app"
)

const (
	RefreshTokenCookieName = "refresh_token"
	mobileClient           = "mobile"
)

type CookieConfig struct {
	Path     string
	Secure   bool
	SameSite http.SameSite
}

type Handler struct {
	command  *app.Command
	query    *app.Query
	verifier app.AccessTokenVerifier
	cookie   CookieConfig
}

func NewHandler(command *app.Command, query *app.Query, verifier app.AccessTokenVerifier, cookie CookieConfig) Handler {
	return Handler{command: command, query: query, verifier: verifier, cookie: cookie}
}

func (h Handler) Register(router server.EchoRouter) {
	middlewares := []server.StrictMiddlewareFunc{NewBearerMiddleware(h.verifier)}
	server.RegisterHandlers(router, server.NewStrictHandler(h, middlewares))
}
