package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/install"
	"github.com/madalinpopa/skills/internal/skill"
)

func newLsCmd() *cobra.Command {
	var local, global bool
	c := &cobra.Command{
		Use:   "ls",
		Short: "List published skills in the store, or installed skills with --local or --global",
		Args:  usageArgs(cobra.NoArgs),
		RunE: func(c *cobra.Command, _ []string) error {
			switch {
			case local && global:
				return usageError{errors.New("--local and --global cannot be combined")}
			case local || global:
				return listInstalled(c, global)
			default:
				return listStore(c)
			}
		},
	}
	c.Flags().BoolVar(&local, "local", false, "list skills installed at the repository root")
	c.Flags().BoolVar(&global, "global", false, "list skills installed in the home directories")
	return c
}

func listStore(c *cobra.Command) error {
	a, err := openApp()
	if err != nil {
		return err
	}
	if err = a.store.Init(c.Context()); err != nil {
		return err
	}
	skills, err := skill.Catalog(os.DirFS(a.store.Dir))
	if err != nil {
		return err
	}
	rows := make([][]string, 0, len(skills))
	for _, sk := range skills {
		rows = append(rows, []string{sk.Name, sk.Description, strings.Join(sk.Tags, ", ")})
	}
	return printRows(c, rows)
}

func listInstalled(c *cobra.Command, global bool) error {
	_, targets, err := openTargets(c, nil, global)
	if err != nil {
		return err
	}
	installed, err := install.Scan(targets)
	if err != nil {
		return err
	}
	rows := make([][]string, 0, len(installed))
	for _, inst := range installed {
		variants := make([]string, 0, len(inst.Targets))
		for _, target := range inst.Targets {
			variants = append(variants, string(target.Variant))
		}
		rows = append(rows, []string{inst.Name, inst.Description, strings.Join(variants, ", ")})
	}
	return printRows(c, rows)
}

func printRows(c *cobra.Command, rows [][]string) error {
	w := tabwriter.NewWriter(c.OutOrStdout(), 0, 0, 3, ' ', 0)
	for _, row := range rows {
		if _, err := fmt.Fprintf(w, "  %s\n", strings.Join(row, "\t")); err != nil {
			return err
		}
	}
	return w.Flush()
}
