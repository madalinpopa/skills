package pgstore

import (
	"context"
	"errors"
	"fmt"

	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/adapters/pgstore/models"
	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type {{ENTITY_PASCAL}}Repository struct {
	pool *pgxpool.Pool
}

func New{{ENTITY_PASCAL}}Repository(pool *pgxpool.Pool) *{{ENTITY_PASCAL}}Repository {
	return &{{ENTITY_PASCAL}}Repository{pool: pool}
}

func (r *{{ENTITY_PASCAL}}Repository) Create(ctx context.Context, {{ENTITY}} *domain.{{ENTITY_PASCAL}}) error {
	err := models.New(r.pool).Create{{ENTITY_PASCAL}}(ctx, models.Create{{ENTITY_PASCAL}}Params{
		{{ENTITY_PASCAL}}Uuid: {{ENTITY}}.UUID,
		Name:      {{ENTITY}}.Name,
		CreatedAt: {{ENTITY}}.CreatedAt,
		UpdatedAt: {{ENTITY}}.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("creating {{ENTITY}}: %w", err)
	}

	return nil
}

func (r *{{ENTITY_PASCAL}}Repository) GetByUUID(ctx context.Context, id domain.{{ENTITY_PASCAL}}UUID) (*domain.{{ENTITY_PASCAL}}, error) {
	row, err := models.New(r.pool).Get{{ENTITY_PASCAL}}ByUUID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("getting {{ENTITY}}: %w", err)
	}

	return &domain.{{ENTITY_PASCAL}}{
		UUID:      row.{{ENTITY_PASCAL}}Uuid,
		Name:      row.Name,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}
