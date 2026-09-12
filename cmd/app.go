package cmd

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/config"
	"github.com/madalinpopa/skills/internal/install"
	"github.com/madalinpopa/skills/internal/store"
)

type app struct {
	home   string
	dir    string
	root   string
	dryRun bool
	cfg    config.Config
	store  store.Store
}

func openApp(c *cobra.Command) (app, error) {
	return newApp(flag(c, "dry-run"))
}

func newApp(dryRun bool) (app, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return app{}, err
	}
	dir := config.Dir(os.Getenv("XDG_CONFIG_HOME"), home)
	var cfg config.Config
	if dryRun {
		cfg, err = config.Load(dir)
	} else {
		cfg, err = config.Init(dir)
	}
	if err != nil {
		return app{}, err
	}
	return app{
		home:   home,
		dir:    dir,
		dryRun: dryRun,
		cfg:    cfg,
		store: store.Store{
			Dir:    filepath.Join(dir, "store"),
			Repo:   cfg.Store.Repo,
			Branch: cfg.Store.Branch,
		},
	}, nil
}

func (a app) ready(ctx context.Context) error {
	if !a.dryRun {
		return a.store.Init(ctx)
	}
	if _, err := os.Stat(a.store.Dir); errors.Is(err, fs.ErrNotExist) {
		return errors.New("the store is not initialised; run 'skills init' first")
	}
	return a.store.Init(ctx)
}

func (a app) installer(commit string, force bool) install.Installer {
	return install.Installer{
		Source:  a.cfg.Store.Repo,
		Commit:  commit,
		Backups: filepath.Join(a.dir, "backups"),
		Force:   force,
		DryRun:  a.dryRun,
		Modes:   runtime.GOOS != "windows",
		Now:     utcNow,
	}
}

func (a app) renderer(c *cobra.Command) renderer {
	out := c.OutOrStdout()
	return renderer{
		out:     out,
		root:    a.root,
		color:   colorEnabled(flag(c, "no-color"), os.Getenv("NO_COLOR"), isTerminal(out)),
		verbose: flag(c, "verbose"),
		global:  flag(c, "global"),
		agents:  stringsFlag(c, "agent"),
	}
}

func flag(c *cobra.Command, name string) bool {
	value, err := c.Flags().GetBool(name)
	return err == nil && value
}

func stringsFlag(c *cobra.Command, name string) []string {
	values, err := c.Flags().GetStringSlice(name)
	if err != nil {
		return nil
	}
	return values
}
