package {{MODULE_NAME}}

import (
	"context"
	"embed"
	"fmt"
	"time"

	"{{MODULE_PATH}}/internal/modules"
	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/adapters/pgstore"
	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/adapters/rest"
	"{{MODULE_PATH}}/internal/modules/{{MODULE_NAME}}/app"
	"{{MODULE_PATH}}/internal/platform/pgkit"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed adapters/pgstore/migrations/*.sql
var migrations embed.FS

type Module struct {
	config  Config
	pool    *pgxpool.Pool
	handler rest.Handler
}

func NewModule(cfg Config, pool *pgxpool.Pool) *Module {
	return &Module{config: cfg, pool: pool}
}

func (m *Module) Name() modules.Name {
	return "{{MODULE_NAME}}"
}

func (m *Module) Init(ctx context.Context) error {
	if err := m.config.Validate(); err != nil {
		return fmt.Errorf("invalid module config: %w", err)
	}

	repository := pgstore.New{{ENTITY_PASCAL}}Repository(m.pool)
	m.handler = rest.NewHandler(
		app.NewCommand(repository, time.Now),
		app.NewQuery(repository),
	)

	return pgkit.MigrateDatabaseUp(ctx, string(m.Name()), m.pool, migrations, "adapters/pgstore/migrations")
}

func (m *Module) RegisterContracts(_ context.Context, _ *modules.Contracts) error {
	return nil
}

func (m *Module) Connect(_ context.Context, _ *modules.Contracts) error {
	return nil
}

func (m *Module) RegisterHTTP(_ context.Context, e modules.EchoRouter) error {
	m.handler.Register(e)
	return nil
}
