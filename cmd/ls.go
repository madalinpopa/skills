package cmd

import (
	"errors"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/install"
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
	a, err := openApp(c)
	if err != nil {
		return err
	}
	_, skills, err := openStore(c.Context(), a)
	if err != nil {
		return err
	}
	blocks := make([]block, 0, len(skills))
	for _, sk := range skills {
		blocks = append(blocks, block{name: sk.Name, labels: strings.Join(sk.Tags, ", "), description: sk.Description})
	}
	return printBlocks(c.OutOrStdout(), outputWidth(c.OutOrStdout()), blocks)
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
	blocks := make([]block, 0, len(installed))
	attention := false
	for _, inst := range installed {
		variants := make([]string, 0, len(inst.Targets))
		for _, target := range inst.Targets {
			variants = append(variants, string(target.Variant))
		}
		description, err := install.Describe(inst)
		switch {
		case len(inst.Issues) > 0:
			description = "! damaged lock: " + inst.Issues[0].Reason
			attention = true
		case err != nil:
			description = "! description unavailable"
			attention = true
		}
		blocks = append(blocks, block{name: inst.Name, labels: strings.Join(variants, ", "), description: description})
	}
	if err := printBlocks(c.OutOrStdout(), outputWidth(c.OutOrStdout()), blocks); err != nil {
		return err
	}
	if attention {
		return ErrAttention
	}
	return nil
}

type block struct {
	name        string
	labels      string
	description string
}

const descriptionIndent = "      "

func printBlocks(w io.Writer, width int, blocks []block) error {
	var b strings.Builder
	for i, bl := range blocks {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("  " + bl.name)
		if bl.labels != "" {
			b.WriteString("   " + bl.labels)
		}
		b.WriteString("\n")
		for _, line := range wrap(bl.description, width-len(descriptionIndent)) {
			b.WriteString(descriptionIndent + line + "\n")
		}
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func wrap(text string, width int) []string {
	var lines []string
	var line strings.Builder
	length := 0
	for word := range strings.FieldsSeq(text) {
		size := utf8.RuneCountInString(word)
		switch {
		case length == 0:
			line.WriteString(word)
			length = size
		case length+1+size <= width:
			line.WriteString(" " + word)
			length += 1 + size
		default:
			lines = append(lines, line.String())
			line.Reset()
			line.WriteString(word)
			length = size
		}
	}
	if length > 0 {
		lines = append(lines, line.String())
	}
	return lines
}
