package app

import (
	"context"

	"{{MODULE_PATH}}/internal/modules/auth/domain"
)

type Query struct {
	users domain.UserRepository
}

func NewQuery(users domain.UserRepository) *Query {
	return &Query{users: users}
}

func (q *Query) GetUser(ctx context.Context, id domain.UserUUID) (*domain.User, error) {
	return q.users.GetByUUID(ctx, id)
}
