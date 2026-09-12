package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/skill"
)

func newLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List published skills in the store",
		Args:  noArgs,
		RunE: func(c *cobra.Command, _ []string) error {
			s, err := openStore()
			if err != nil {
				return err
			}
			if err = s.Init(c.Context()); err != nil {
				return err
			}
			skills, err := skill.Catalog(os.DirFS(s.Dir))
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(c.OutOrStdout(), 0, 0, 3, ' ', 0)
			for _, sk := range skills {
				if _, err = fmt.Fprintf(w, "  %s\t%s\t%s\n", sk.Name, sk.Description, strings.Join(sk.Tags, ", ")); err != nil {
					return err
				}
			}
			return w.Flush()
		},
	}
}
