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
			req, err := resolveUpdate(c, names, agents, global, force)
			if err != nil {
				return err
			}
			results, err := req.installer.Install(req.requests)
			if err != nil {
				return err
			}
			printResults(c, results)
			return nil
		},
	}
	c.Flags().StringSliceVar(&agents, "agent", nil, "narrow to certain agents (default: config defaults)")
	c.Flags().BoolVar(&global, "global", false, "act on the home directories instead of the repository")
	c.Flags().BoolVar(&force, "force", false, "overwrite skills you have edited (backed up first)")
	return c
}

func resolveUpdate(c *cobra.Command, names, agents []string, global, force bool) (installRequest, error) {
	a, targets, err := openTargets(c, agents, global)
	if err != nil {
		return installRequest{}, err
	}
	installed, err := install.Scan(targets)
	if err != nil {
		return installRequest{}, err
	}
	chosen, err := install.Find(installed, names)
	if err != nil {
		return installRequest{}, err
	}
	catalog, installer, err := openStore(c.Context(), a, force)
	if err != nil {
		return installRequest{}, err
	}
	requests, err := install.UpdateRequests(a.store.Dir, catalog, chosen)
	if err != nil {
		return installRequest{}, err
	}
	return installRequest{installer: installer, requests: requests}, nil
}
