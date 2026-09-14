package testkit

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"

	"{{MODULE_PATH}}/internal/modules"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type Env struct {
	Pool      *pgxpool.Pool
	Echo      *echo.Echo
	Contracts *modules.Contracts
	URL       string
}

func Start(ctx context.Context) (*Env, func() error) {
	dsn, terminate := StartPostgres(ctx)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}

	env := &Env{
		Pool:      pool,
		Echo:      echo.New(),
		Contracts: &modules.Contracts{},
	}

	return env, func() error {
		pool.Close()
		return terminate()
	}
}

func (e *Env) Wire(ctx context.Context, m modules.Module) error {
	if err := m.Init(ctx); err != nil {
		return errors.Join(fmt.Errorf("init module %s", m.Name()), err)
	}

	if err := m.RegisterContracts(ctx, e.Contracts); err != nil {
		return errors.Join(fmt.Errorf("register contracts for module %s", m.Name()), err)
	}

	if err := e.Contracts.Verify(ctx); err != nil {
		return fmt.Errorf("verifying module contracts: %w", err)
	}

	if err := m.Connect(ctx, e.Contracts); err != nil {
		return errors.Join(fmt.Errorf("connect module %s", m.Name()), err)
	}

	if err := m.RegisterHTTP(ctx, e.Echo); err != nil {
		return errors.Join(fmt.Errorf("register http for module %s", m.Name()), err)
	}

	return nil
}

func (e *Env) Serve() func() {
	server := httptest.NewServer(e.Echo)
	e.URL = server.URL

	return server.Close
}
