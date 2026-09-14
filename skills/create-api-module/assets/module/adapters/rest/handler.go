package rest

//go:generate go tool oapi-codegen --config=oapi-codegen.yaml openapi.yaml

import (
	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/adapters/rest/server"
	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/app"
)

type Handler struct {
	command *app.Command
	query   *app.Query
}

func NewHandler(command *app.Command, query *app.Query) Handler {
	return Handler{command: command, query: query}
}

func (h Handler) Register(router server.EchoRouter, middlewares ...server.StrictMiddlewareFunc) {
	server.RegisterHandlers(router, server.NewStrictHandler(h, middlewares))
}
