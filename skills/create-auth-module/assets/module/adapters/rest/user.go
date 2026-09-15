package rest

import (
	"context"
	"errors"

	"{{MODULE_PATH}}/internal/modules/auth/adapters/rest/server"
	"{{MODULE_PATH}}/internal/modules/auth/app"
	"{{MODULE_PATH}}/internal/modules/auth/domain"
)

func (h Handler) RegisterUser(ctx context.Context, request server.RegisterUserRequestObject) (server.RegisterUserResponseObject, error) {
	user, err := h.command.Register(ctx, app.Register{
		Name:     request.Body.Name,
		Email:    request.Body.Email,
		Password: request.Body.Password,
	})
	if err != nil {
		return registerError(err)
	}

	return server.RegisterUser201JSONResponse{UserUuid: user.UUID.String()}, nil
}

func registerError(err error) (server.RegisterUserResponseObject, error) {
	badRequest := func(slug, message string) (server.RegisterUserResponseObject, error) {
		return server.RegisterUser400JSONResponse{Slug: slug, Message: message}, nil
	}

	switch {
	case errors.Is(err, domain.ErrEmptyName):
		return badRequest("name-empty", "The name must not be empty")
	case errors.Is(err, domain.ErrNameTooLong):
		return badRequest("name-too-long", "The name is too long")
	case errors.Is(err, domain.ErrInvalidEmail):
		return badRequest("email-invalid", "The email address is invalid")
	case errors.Is(err, domain.ErrPasswordTooShort):
		return badRequest("password-too-short", "The password must have at least 8 characters")
	case errors.Is(err, domain.ErrPasswordTooLong):
		return badRequest("password-too-long", "The password must have at most 72 bytes")
	case errors.Is(err, domain.ErrEmailTaken):
		return server.RegisterUser409JSONResponse{Slug: "email-taken", Message: "The email address is already registered"}, nil
	default:
		return nil, err
	}
}

func (h Handler) GetMe(ctx context.Context, _ server.GetMeRequestObject) (server.GetMeResponseObject, error) {
	identity, ok := app.IdentityFromContext(ctx)
	if !ok {
		return nil, errors.New("get me reached without an authenticated identity")
	}

	user, err := h.query.GetUser(ctx, identity.UserUUID)
	if errors.Is(err, domain.ErrUserNotFound) {
		return server.GetMe404JSONResponse{Slug: "user-not-found", Message: "The user no longer exists"}, nil
	}
	if err != nil {
		return nil, err
	}

	return server.GetMe200JSONResponse{
		UserUuid:  user.UUID.String(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
