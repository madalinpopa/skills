package cmd

import (
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create the config and clone the store",
		Args:  usageArgs(cobra.NoArgs),
		RunE: func(c *cobra.Command, _ []string) error {
			a, err := openApp(c)
			if err != nil {
				return err
			}
			if err = a.ready(c.Context()); err != nil {
				return err
			}
			commit, err := a.store.Commit(c.Context())
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
		Args:  usageArgs(cobra.NoArgs),
		RunE: func(c *cobra.Command, _ []string) error {
			a, err := openApp(c)
			if err != nil {
				return err
			}
			if err = a.ready(c.Context()); err != nil {
				return err
			}
			if a.dryRun {
				return previewSync(c, a)
			}
			result, err := a.store.Sync(c.Context())
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

func previewSync(c *cobra.Command, a app) error {
	result, err := a.store.Preview(c.Context())
	if err != nil {
		return err
	}
	if result.Old == result.New {
		c.Printf("  store up to date  %s\n", short(result.New))
		return nil
	}
	c.Printf("  would pull store  %s -> %s\n", short(result.Old), short(result.New))
	c.Println("  fast-forward eligibility is checked during sync")
	return nil
}

func short(commit string) string {
	return commit[:7]
}
