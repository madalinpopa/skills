package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/install"
)

func newRemoveCmd() *cobra.Command {
	var agents []string
	var global, force bool
	c := &cobra.Command{
		Use:   "remove <skill>...",
		Short: "Uninstall skills from the selected scope",
		Args:  usageArgs(cobra.MinimumNArgs(1)),
		RunE: func(c *cobra.Command, names []string) error {
			a, targets, err := openTargets(c, agents, global)
			if err != nil {
				return err
			}
			installed, err := install.Scan(targets)
			if err != nil {
				return err
			}
			chosen, err := install.Find(installed, names)
			if err != nil && !global {
				return suggestGlobal(a, agents, names, err)
			}
			if err != nil {
				return err
			}
			results, err := a.installer("", force).Remove(chosen)
			if err != nil {
				return err
			}
			return a.renderer(c).results("remove", results)
		},
	}
	c.Flags().StringSliceVar(&agents, "agent", nil, "narrow to certain agents (default: config defaults)")
	c.Flags().BoolVar(&global, "global", false, "act on the home directories instead of the repository")
	c.Flags().BoolVar(&force, "force", false, "remove skills you have edited (backed up first)")
	return c
}

func suggestGlobal(a app, agents, names []string, notFound error) error {
	targets, err := install.Targets(a.cfg, a.home, true, agents)
	if err != nil {
		return notFound
	}
	installed, err := install.Scan(targets)
	if err != nil {
		return notFound
	}
	if _, err = install.Find(installed, names); err != nil {
		return notFound
	}
	return fmt.Errorf("%s: only installed globally, use --global", strings.Join(names, ", "))
}
