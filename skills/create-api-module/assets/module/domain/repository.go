package domain

import "context"

type {{ENTITY_PASCAL}}Repository interface {
	Create(ctx context.Context, {{ENTITY}} *{{ENTITY_PASCAL}}) error
	GetByUUID(ctx context.Context, id {{ENTITY_PASCAL}}UUID) (*{{ENTITY_PASCAL}}, error)
}
