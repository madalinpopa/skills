package cmd

import (
	"os"
	"path/filepath"

	"github.com/madalinpopa/skills/internal/config"
	"github.com/madalinpopa/skills/internal/store"
)

type app struct {
	home  string
	cfg   config.Config
	store store.Store
}

func openApp() (app, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return app{}, err
	}
	dir := config.Dir(os.Getenv("XDG_CONFIG_HOME"), home)
	cfg, err := config.Init(dir)
	if err != nil {
		return app{}, err
	}
	return app{
		home: home,
		cfg:  cfg,
		store: store.Store{
			Dir:    filepath.Join(dir, "store"),
			Repo:   cfg.Store.Repo,
			Branch: cfg.Store.Branch,
		},
	}, nil
}
