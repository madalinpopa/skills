package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/config"
	"github.com/madalinpopa/skills/internal/store"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create the config and clone the store",
		Args:  noArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			if err = s.Init(c.Context()); err != nil {
				return err
			}
			commit, err := s.Commit(c.Context())
			if err != nil {
				return err
			}
			c.Printf("  store ready  %s\n", short(commit))
			return nil
		},
	}
}

func newSyncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Fast-forward the configured store branch",
		Args:  noArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			if err = s.Init(c.Context()); err != nil {
				return err
			}
			result, err := s.Sync(c.Context())
			if err != nil {
				return err
			}
			if result.Old == result.New {
				c.Printf("  store up to date  %s\n", short(result.New))
				return nil
			}
			c.Printf("  pulled store  %s -> %s\n", short(result.Old), short(result.New))
			return nil
		},
	}
}

func openStore() (store.Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return store.Store{}, err
	}
	dir := config.Dir(os.Getenv("XDG_CONFIG_HOME"), home)
	cfg, err := config.Init(dir)
	if err != nil {
		return store.Store{}, err
	}
	return store.Store{
		Dir:    filepath.Join(dir, "store"),
		Repo:   cfg.Store.Repo,
		Branch: cfg.Store.Branch,
	}, nil
}

func short(commit string) string {
	return commit[:7]
}
