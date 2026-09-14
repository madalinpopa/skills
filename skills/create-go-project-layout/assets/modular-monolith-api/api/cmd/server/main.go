package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"{{MODULE_PATH}}/internal/server"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
)

const version = "0.1.0"

var (
	port        string
	showVersion bool
	enableDebug bool
)

func init() {
	flag.StringVar(&port, "port", "8080", "server port to listen for requests")
	flag.BoolVar(&showVersion, "version", false, "print the server version")
	flag.BoolVar(&enableDebug, "debug", false, "enable debug logging")
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Printf("version: %s\n", version)
		return
	}

	logLevel := slog.LevelInfo
	if enableDebug {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})))

	if err := run(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := server.LoadConfigFromEnv()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	pool, err := pgxpool.New(ctx, os.Getenv("POSTGRES_URL"))
	if err != nil {
		return fmt.Errorf("opening database pool: %w", err)
	}
	defer pool.Close()

	e := echo.New()
	e.Logger = slog.Default()
	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	srv, err := server.New(ctx, cfg, e, pool)
	if err != nil {
		return fmt.Errorf("building server: %w", err)
	}

	return srv.Run(ctx, port)
}
