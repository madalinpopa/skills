package cmd

import (
	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/install"
)

func newUpdateCmd() *cobra.Command {
	var agents []string
	var global, force bool
	c := &cobra.Command{
		Use:   "update [skill...]",
		Short: "Re-install installed skills from the store",
		Args:  usageArgs(cobra.ArbitraryArgs),
		RunE: func(c *cobra.Command, names []string) error {
			a, req, err := resolveUpdate(c, names, agents, global, force)
			if err != nil {
				return err
			}
			results, err := req.installer.Install(req.requests)
			return a.report(c, "update", results, err)
		},
	}
	c.Flags().StringSliceVar(&agents, "agent", nil, "narrow to certain agents (default: config defaults)")
	c.Flags().BoolVar(&global, "global", false, "act on the home directories instead of the repository")
	c.Flags().BoolVar(&force, "force", false, "overwrite skills you have edited (backed up first)")
	return c
}

func resolveUpdate(c *cobra.Command, names, agents []string, global, force bool) (app, installRequest, error) {
	a, targets, err := openTargets(c, agents, global)
	if err != nil {
		return app{}, installRequest{}, err
	}
	installed, err := install.Scan(targets)
	if err != nil {
		return app{}, installRequest{}, err
	}
	chosen, err := a.find(installed, names, agents, global)
	if err != nil {
		return app{}, installRequest{}, err
	}
	rev, catalog, err := openStore(c.Context(), a)
	if err != nil {
		return app{}, installRequest{}, err
	}
	requests, err := install.UpdateRequests(rev.tree, catalog, chosen)
	if err != nil {
		return app{}, installRequest{}, err
	}
	return a, installRequest{installer: a.installer(rev.commit, force), requests: requests}, nil
}
