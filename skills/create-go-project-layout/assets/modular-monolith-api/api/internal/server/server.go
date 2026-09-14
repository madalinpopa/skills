package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"{{MODULE_PATH}}/internal/modules"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

type Server struct {
	config  Config
	router  *echo.Echo
	pool    *pgxpool.Pool
	modules []modules.Module
}

func New(ctx context.Context, c Config, e *echo.Echo, pool *pgxpool.Pool) (*Server, error) {
	contracts := modules.Contracts{}

	allModules := []modules.Module{}

	start := time.Now()

	for _, m := range allModules {
		if err := m.Init(ctx); err != nil {
			return nil, errors.Join(fmt.Errorf("init module %s", m.Name()), err)
		}
	}

	for _, m := range allModules {
		if err := m.RegisterContracts(ctx, &contracts); err != nil {
			return nil, errors.Join(fmt.Errorf("register contracts for module %s", m.Name()), err)
		}
	}

	if err := contracts.Verify(ctx); err != nil {
		return nil, fmt.Errorf("verifying module contracts: %w", err)
	}

	for _, m := range allModules {
		if err := m.Connect(ctx, &contracts); err != nil {
			return nil, errors.Join(fmt.Errorf("connect module %s", m.Name()), err)
		}
	}

	for _, m := range allModules {
		if err := m.RegisterHTTP(ctx, e); err != nil {
			return nil, errors.Join(fmt.Errorf("register http for module %s", m.Name()), err)
		}
	}

	slog.Debug("initialized modules", "duration", time.Since(start), "modules", len(allModules))

	return &Server{
		config:  c,
		router:  e,
		pool:    pool,
		modules: allModules,
	}, nil
}

func (s *Server) Run(ctx context.Context, port string) error {
	startConfig := echo.StartConfig{
		Address:         ":" + port,
		GracefulTimeout: 5 * time.Second,
	}

	return startConfig.Start(ctx, s.router)
}
