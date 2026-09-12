package cmd

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"

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

func diffSkill(c *cobra.Command, a app, inst install.Installation) error {
	for _, target := range inst.Targets {
		dir := filepath.Join(target.Dir, inst.Name)
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
		recorded, err := a.store.Files(c.Context(), lock.Commit, path.Join("skills", inst.Name))
		if err != nil {
			return fmt.Errorf("%s: %w", inst.Name, err)
		}
		base, err := skill.RenderFiles(recorded, string(target.Variant))
		if err != nil {
			return err
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
