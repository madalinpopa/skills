package cmd

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/install"
	"github.com/madalinpopa/skills/internal/project"
	"github.com/madalinpopa/skills/internal/skill"
)

type installRequest struct {
	installer install.Installer
	requests  []install.Request
}

func newInstallCmd() *cobra.Command {
	var agents []string
	var global bool
	c := &cobra.Command{
		Use:   "install <skill>...",
		Short: "Install skills at the repository root",
		Args:  usageArgs(cobra.MinimumNArgs(1)),
		RunE: func(c *cobra.Command, names []string) error {
			a, req, err := resolveInstall(c, names, agents, global)
			if err != nil {
				return err
			}
			results, err := req.installer.Install(req.requests)
			return a.report(c, "install", results, err)
		},
	}
	c.Flags().StringSliceVar(&agents, "agent", nil, "narrow to certain agents (default: config defaults)")
	c.Flags().BoolVar(&global, "global", false, "act on the home directories instead of the repository")
	return c
}

func resolveInstall(c *cobra.Command, names, agents []string, global bool) (app, installRequest, error) {
	a, targets, err := openTargets(c, agents, global)
	if err != nil {
		return app{}, installRequest{}, err
	}
	catalog, installer, err := openStore(c.Context(), a, false)
	if err != nil {
		return app{}, installRequest{}, err
	}
	skills, err := install.Select(catalog, names)
	if err != nil {
		return app{}, installRequest{}, err
	}
	requests, err := install.Requests(a.store.Dir, skills, targets)
	if err != nil {
		return app{}, installRequest{}, err
	}
	return a, installRequest{installer: installer, requests: requests}, nil
}

func openTargets(c *cobra.Command, agents []string, global bool) (app, []install.Destination, error) {
	a, err := openApp(c)
	if err != nil {
		return app{}, nil, err
	}
	root, err := scopeRoot(c, a.home, global)
	if err != nil {
		return app{}, nil, err
	}
	a.root = root
	targets, err := install.Targets(a.cfg, root, global, agents)
	if errors.Is(err, install.ErrUnknownAgent) {
		return app{}, nil, usageError{err}
	}
	if err != nil {
		return app{}, nil, err
	}
	return a, targets, nil
}

func openStore(ctx context.Context, a app, force bool) ([]skill.Skill, install.Installer, error) {
	if err := a.ready(ctx); err != nil {
		return nil, install.Installer{}, err
	}
	catalog, err := skill.Catalog(os.DirFS(a.store.Dir))
	if err != nil {
		return nil, install.Installer{}, err
	}
	commit, err := a.store.Commit(ctx)
	if err != nil {
		return nil, install.Installer{}, err
	}
	return catalog, a.installer(commit, force), nil
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

func utcNow() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}
