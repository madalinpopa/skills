package auth

import (
	"context"
	"crypto/rand"
	"embed"
	"fmt"
	"time"

	"{{MODULE_PATH}}/internal/modules"
	"{{MODULE_PATH}}/internal/modules/auth/adapters/jwtx"
	"{{MODULE_PATH}}/internal/modules/auth/adapters/pgstore"
	"{{MODULE_PATH}}/internal/modules/auth/adapters/refreshtx"
	"{{MODULE_PATH}}/internal/modules/auth/adapters/rest"
	"{{MODULE_PATH}}/internal/modules/auth/app"
	"{{MODULE_PATH}}/internal/platform/pgkit"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed adapters/pgstore/migrations/*.sql
var migrations embed.FS

const cookiePath = "/auth"

type Module struct {
	config   Config
	pool     *pgxpool.Pool
	handler  rest.Handler
	verifier app.AccessTokenVerifier
}

func NewModule(cfg Config, pool *pgxpool.Pool) *Module {
	return &Module{config: cfg, pool: pool}
}

func (m *Module) Name() modules.Name {
	return "auth"
}

func (m *Module) Init(ctx context.Context) error {
	if err := m.config.Validate(); err != nil {
		return fmt.Errorf("invalid module config: %w", err)
	}

	accessTokens := jwtx.NewService(jwtx.Config{
		SigningKey: m.config.SigningKey,
		Issuer:     m.config.Issuer,
		Audience:   m.config.Audience,
		TTL:        m.config.AccessTokenTTL,
	}, time.Now)

	refreshTokens := refreshtx.NewService(m.config.RefreshTokenTTL, rand.Reader, time.Now)
	users := pgstore.NewUserRepository(m.pool)

	command := app.NewCommand(app.CommandArgs{
		Users:   users,
		Tokens:  pgstore.NewTokenRepository(m.pool),
		Access:  accessTokens,
		Refresh: refreshTokens,
		Clock:   time.Now,
	})

	m.verifier = accessTokens
	m.handler = rest.NewHandler(command, app.NewQuery(users), accessTokens, rest.CookieConfig{
		Path:     cookiePath,
		Secure:   m.config.CookieSecure,
		SameSite: m.config.CookieSameSite,
	})

	return pgkit.MigrateDatabaseUp(ctx, string(m.Name()), m.pool, migrations, "adapters/pgstore/migrations")
}

func (m *Module) RegisterContracts(_ context.Context, c *modules.Contracts) error {
	c.Auth = NewContract(m.verifier)
	return nil
}

func (m *Module) Connect(_ context.Context, _ *modules.Contracts) error {
	return nil
}

func (m *Module) RegisterHTTP(_ context.Context, e modules.EchoRouter) error {
	m.handler.Register(e)
	return nil
}
