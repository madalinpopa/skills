package cmd

import (
	"context"
	"errors"
	"io/fs"
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

type revision struct {
	commit string
	tree   fs.FS
}

func newInstallCmd() *cobra.Command {
	var agents []string
	var global, force bool
	c := &cobra.Command{
		Use:   "install <skill>...",
		Short: "Install skills at the repository root",
		Args:  usageArgs(cobra.MinimumNArgs(1)),
		RunE: func(c *cobra.Command, names []string) error {
			a, req, err := resolveInstall(c, names, agents, global, force)
			if err != nil {
				return err
			}
			results, err := req.installer.Install(req.requests)
			return a.report(c, "install", results, err)
		},
	}
	c.Flags().StringSliceVar(&agents, "agent", nil, "narrow to certain agents (default: config defaults)")
	c.Flags().BoolVar(&global, "global", false, "act on the home directories instead of the repository")
	c.Flags().BoolVar(&force, "force", false, "overwrite skills you have edited or did not install (backed up first)")
	return c
}

func resolveInstall(c *cobra.Command, names, agents []string, global, force bool) (app, installRequest, error) {
	a, targets, err := openTargets(c, agents, global)
	if err != nil {
		return app{}, installRequest{}, err
	}
	rev, catalog, err := openStore(c.Context(), a)
	if err != nil {
		return app{}, installRequest{}, err
	}
	skills, err := install.Select(catalog, names)
	if err != nil {
		return app{}, installRequest{}, err
	}
	requests, err := install.Requests(rev.tree, skills, targets)
	if err != nil {
		return app{}, installRequest{}, err
	}
	return a, installRequest{installer: a.installer(rev.commit, force), requests: requests}, nil
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

func openStore(ctx context.Context, a app) (revision, []skill.Skill, error) {
	if err := a.ready(ctx); err != nil {
		return revision{}, nil, err
	}
	commit, err := a.store.Commit(ctx)
	if err != nil {
		return revision{}, nil, err
	}
	tree, err := a.store.Tree(ctx, commit)
	if err != nil {
		return revision{}, nil, err
	}
	catalog, err := skill.Catalog(tree)
	if err != nil {
		return revision{}, nil, err
	}
	return revision{commit: commit, tree: tree}, catalog, nil
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
