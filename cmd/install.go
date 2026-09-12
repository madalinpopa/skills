package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/install"
	"github.com/madalinpopa/skills/internal/project"
	"github.com/madalinpopa/skills/internal/skill"
)

type installRequest struct {
	targets []install.Destination
	skills  []skill.Skill
}

func newInstallCmd() *cobra.Command {
	var agents []string
	var global bool
	c := &cobra.Command{
		Use:   "install <skill>...",
		Short: "Install skills at the repository root",
		Args:  usageArgs(cobra.MinimumNArgs(1)),
		RunE: func(c *cobra.Command, names []string) error {
			_, err := resolveInstall(c, names, agents, global)
			return err
		},
	}
	c.Flags().StringSliceVar(&agents, "agent", nil, "narrow to certain agents (default: config defaults)")
	c.Flags().BoolVar(&global, "global", false, "act on the home directories instead of the repository")
	return c
}

func resolveInstall(c *cobra.Command, names, agents []string, global bool) (installRequest, error) {
	a, err := openApp()
	if err != nil {
		return installRequest{}, err
	}
	root, err := scopeRoot(c, a.home, global)
	if err != nil {
		return installRequest{}, err
	}
	targets, err := install.Targets(a.cfg, root, global, agents)
	if errors.Is(err, install.ErrUnknownAgent) {
		return installRequest{}, usageError{err}
	}
	if err != nil {
		return installRequest{}, err
	}
	if err = a.store.Init(c.Context()); err != nil {
		return installRequest{}, err
	}
	catalog, err := skill.Catalog(os.DirFS(a.store.Dir))
	if err != nil {
		return installRequest{}, err
	}
	skills, err := install.Select(catalog, names)
	if err != nil {
		return installRequest{}, err
	}
	return installRequest{targets: targets, skills: skills}, nil
}

func scopeRoot(c *cobra.Command, home string, global bool) (string, error) {
	if global {
		return home, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	scope, err := project.Detect(wd)
	if err != nil {
		return "", err
	}
	if scope.Warning != "" {
		c.PrintErrln("warning:", scope.Warning)
	}
	return scope.Root, nil
}
