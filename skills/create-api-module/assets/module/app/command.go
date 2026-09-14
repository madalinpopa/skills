package app

import (
	"context"
	"time"

	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/domain"
)

type Command struct {
	repository domain.{{ENTITY_PASCAL}}Repository
	clock      func() time.Time
}

func NewCommand(repository domain.{{ENTITY_PASCAL}}Repository, clock func() time.Time) *Command {
	return &Command{repository: repository, clock: clock}
}

type Create{{ENTITY_PASCAL}} struct {
	Name string
}

func (c *Command) Create{{ENTITY_PASCAL}}(ctx context.Context, cmd Create{{ENTITY_PASCAL}}) (*domain.{{ENTITY_PASCAL}}, error) {
	{{ENTITY}}, err := domain.New{{ENTITY_PASCAL}}(cmd.Name, c.clock().UTC())
	if err != nil {
		return nil, err
	}

	if err := c.repository.Create(ctx, {{ENTITY}}); err != nil {
		return nil, err
	}

	return {{ENTITY}}, nil
}
