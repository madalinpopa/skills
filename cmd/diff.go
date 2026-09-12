package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/aymanbagabas/go-udiff"
	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/install"
	"github.com/madalinpopa/skills/internal/skill"
)

func newDiffCmd() *cobra.Command {
	var agents []string
	var global bool
	c := &cobra.Command{
		Use:   "diff <skill>",
		Short: "Show your edits to an installed skill",
		Args:  usageArgs(cobra.ExactArgs(1)),
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
			if err != nil {
				return err
			}
			return diffSkill(c, a, chosen[0])
		},
	}
	c.Flags().StringSliceVar(&agents, "agent", nil, "narrow to certain agents (default: config defaults)")
	c.Flags().BoolVar(&global, "global", false, "act on the home directories instead of the repository")
	return c
}

func explain(issues []install.Issue) string {
	lines := make([]string, 0, len(issues))
	for _, issue := range issues {
		lines = append(lines, issue.Path+" "+issue.Reason)
	}
	return strings.Join(lines, "; ")
}

func diffSkill(c *cobra.Command, a app, inst install.Installation) error {
	for _, target := range inst.Targets {
		dir := filepath.Join(target.Dir, inst.Name)
		issues, err := install.Check(dir)
		if err != nil {
			return err
		}
		if len(issues) > 0 {
			return errors.New(explain(issues))
		}
		lock, err := install.ReadLock(dir)
		if err != nil {
			return err
		}
		if lock.Source != a.cfg.Store.Repo {
			return fmt.Errorf("%s was installed from %s, not the configured store %s", inst.Name, lock.Source, a.cfg.Store.Repo)
		}
		if err = a.ready(c.Context()); err != nil {
			return err
		}
		tree, err := a.store.Tree(c.Context(), lock.Commit)
		if err != nil {
			return fmt.Errorf("%s: %w", inst.Name, err)
		}
		source, err := fs.Sub(tree, path.Join("skills", inst.Name))
		if err != nil {
			return fmt.Errorf("%s: %w", inst.Name, err)
		}
		base, err := skill.Render(source, string(target.Variant))
		if err != nil {
			return fmt.Errorf("%s: %w", inst.Name, err)
		}
		have, err := skill.Load(os.DirFS(dir))
		if err != nil {
			return err
		}
		delete(have, install.LockFile)
		paths := slices.Sorted(maps.Keys(base))
		for p := range have {
			if _, ok := base[p]; !ok {
				paths = append(paths, p)
			}
		}
		slices.Sort(paths)
		for _, p := range paths {
			old, now := base[p].Data, have[p].Data
			if bytes.Equal(old, now) {
				continue
			}
			rel := a.renderer(c).relative(filepath.Join(dir, filepath.FromSlash(p)))
			c.Print(udiff.Unified("a/"+rel, "b/"+rel, string(old), string(now)))
		}
	}
	return nil
}
