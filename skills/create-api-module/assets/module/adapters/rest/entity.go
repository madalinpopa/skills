package rest

import (
	"context"
	"errors"

	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/adapters/rest/server"
	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/app"
	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/domain"
)

func (h Handler) Create{{ENTITY_PASCAL}}(ctx context.Context, request server.Create{{ENTITY_PASCAL}}RequestObject) (server.Create{{ENTITY_PASCAL}}ResponseObject, error) {
	{{ENTITY}}, err := h.command.Create{{ENTITY_PASCAL}}(ctx, app.Create{{ENTITY_PASCAL}}{Name: request.Body.Name})
	if errors.Is(err, domain.ErrEmptyName) {
		return server.Create{{ENTITY_PASCAL}}400JSONResponse{
			Slug:    "{{ENTITY}}-name-empty",
			Message: "The {{ENTITY}} name must not be empty",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return server.Create{{ENTITY_PASCAL}}201JSONResponse({{ENTITY}}Response({{ENTITY}})), nil
}

func (h Handler) Get{{ENTITY_PASCAL}}(ctx context.Context, request server.Get{{ENTITY_PASCAL}}RequestObject) (server.Get{{ENTITY_PASCAL}}ResponseObject, error) {
	notFound := server.Get{{ENTITY_PASCAL}}404JSONResponse{
		Slug:    "{{ENTITY}}-not-found",
		Message: "No {{ENTITY}} has this identifier",
	}

	id, err := domain.Parse{{ENTITY_PASCAL}}UUID(request.Uuid)
	if err != nil {
		return notFound, nil
	}

	{{ENTITY}}, err := h.query.Get{{ENTITY_PASCAL}}(ctx, id)
	if errors.Is(err, domain.ErrNotFound) {
		return notFound, nil
	}
	if err != nil {
		return nil, err
	}

	return server.Get{{ENTITY_PASCAL}}200JSONResponse({{ENTITY}}Response({{ENTITY}})), nil
}

func {{ENTITY}}Response({{ENTITY}} *domain.{{ENTITY_PASCAL}}) server.{{ENTITY_PASCAL}}Response {
	return server.{{ENTITY_PASCAL}}Response{
		Uuid:      {{ENTITY}}.UUID.String(),
		Name:      {{ENTITY}}.Name,
		CreatedAt: {{ENTITY}}.CreatedAt,
		UpdatedAt: {{ENTITY}}.UpdatedAt,
	}
}
