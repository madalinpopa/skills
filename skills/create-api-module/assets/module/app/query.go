package app

import (
	"context"

	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/domain"
)

type Query struct {
	repository domain.{{ENTITY_PASCAL}}Repository
}

func NewQuery(repository domain.{{ENTITY_PASCAL}}Repository) *Query {
	return &Query{repository: repository}
}

func (q *Query) Get{{ENTITY_PASCAL}}(ctx context.Context, id domain.{{ENTITY_PASCAL}}UUID) (*domain.{{ENTITY_PASCAL}}, error) {
	return q.repository.GetByUUID(ctx, id)
}
