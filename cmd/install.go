package cmd

import (
	"errors"
	"os"
	"strings"
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
			req, err := resolveInstall(c, names, agents, global)
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
	commit, err := a.store.Commit(c.Context())
	if err != nil {
		return installRequest{}, err
	}
	requests, err := install.Requests(a.store.Dir, skills, targets)
	if err != nil {
		return installRequest{}, err
	}
	installer := install.Installer{Source: a.cfg.Store.Repo, Commit: commit, Now: utcNow}
	return installRequest{installer: installer, requests: requests}, nil
}

func printResults(c *cobra.Command, results []install.SkillPlan) {
	for _, r := range results {
		switch r.State {
		case install.StateAdd:
			c.Printf("  + %s  added\n", r.Name)
		case install.StateUpdate:
			c.Printf("  ~ %s  updated\n", r.Name)
		case install.StateUnchanged:
			c.Printf("    %s  up to date\n", r.Name)
		case install.StateConflict:
			c.Printf("  ! %s  skipped, you edited %s\n", r.Name, strings.Join(r.Conflicts, ", "))
		case install.StateForeign:
			c.Printf("  ! %s  skipped, installed from %s\n", r.Name, r.Source)
		case install.StateUnavailable:
			c.Printf("  ! %s  unavailable in store, left installed\n", r.Name)
		}
	}
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
