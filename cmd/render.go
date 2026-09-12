package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/madalinpopa/skills/internal/install"
)

const (
	green  = "\x1b[32m"
	yellow = "\x1b[33m"
	red    = "\x1b[31m"
	reset  = "\x1b[0m"
)

type renderer struct {
	out     io.Writer
	root    string
	color   bool
	verbose bool
	global  bool
	dryRun  bool
	agents  []string
}

func (r renderer) results(command string, results []install.SkillPlan) error {
	var b strings.Builder
	width := 0
	for _, res := range results {
		if res.State != install.StateUnchanged {
			width = max(width, len(res.Name))
		}
	}
	changed, attention, failed := 0, 0, 0
	var conflicts []install.SkillPlan
	for _, res := range results {
		var symbol, message, tint string
		switch res.State {
		case install.StateAdd:
			symbol, message, tint = "+", r.did("add", "added"), green
			changed++
		case install.StateUpdate:
			symbol, message, tint = "~", r.did("update", "updated"), green
			changed++
		case install.StateRemove:
			symbol, message, tint = "-", r.did("remove", "removed"), red
			changed++
		case install.StateConflict:
			symbol, message, tint = "!", "skipped, you edited it", yellow
			attention++
			conflicts = append(conflicts, res)
		case install.StateForeign:
			symbol, message, tint = "!", "skipped, installed from "+res.Source, yellow
			attention++
		case install.StateUnavailable:
			symbol, message, tint = "!", "unavailable in store, left installed", yellow
			attention++
		case install.StateUnsupported:
			symbol, message, tint = "!", "skipped, needs manual repair", yellow
			attention++
		case install.StateFailed:
			symbol, message, tint = "!", "failed, see the error below", red
			failed++
		case install.StateUnchanged:
			continue
		}
		line := fmt.Sprintf("%s %-*s%s", symbol, width+3, res.Name, message)
		if r.color {
			line = tint + line + reset
		}
		fmt.Fprintf(&b, "  %s\n", line)
		if r.verbose {
			for _, path := range r.paths(res) {
				fmt.Fprintf(&b, "      %s\n", path)
			}
		}
		for _, issue := range res.Issues {
			fmt.Fprintf(&b, "      %s %s\n", r.relative(issue.Path), issue.Reason)
		}
		for _, backup := range res.Backups {
			fmt.Fprintf(&b, "      backed up to %s\n", backup)
		}
		if r.dryRun && (res.State == install.StateRemove || res.Forced) {
			fmt.Fprint(&b, "      would back up first\n")
		}
	}
	summary := fmt.Sprintf("%d %s, %d %s", len(results), plural(len(results), "skill", "skills"), changed, r.did("change", "changed"))
	if attention > 0 {
		summary += fmt.Sprintf(", %d %s attention", attention, plural(attention, "needs", "need"))
	}
	if failed > 0 {
		summary += fmt.Sprintf(", %d failed", failed)
	}
	fmt.Fprintf(&b, "\n  %s\n", summary)
	if len(conflicts) > 0 {
		for _, hint := range r.hints(command, conflicts) {
			fmt.Fprintf(&b, "  %s\n", hint)
		}
	}
	_, err := io.WriteString(r.out, b.String())
	return err
}

func (r renderer) did(planned, done string) string {
	if r.dryRun {
		return "would " + planned
	}
	return done
}

func (r renderer) hints(command string, conflicts []install.SkillPlan) []string {
	var hints, names, managed []string
	for _, res := range conflicts {
		names = append(names, res.Name)
		if res.Managed {
			managed = append(managed, res.Name)
			continue
		}
		hints = append(hints, res.Name+" was not installed by skills, so there is nothing to diff.")
	}
	force := r.command(command, names...) + " --force"
	if len(managed) == 0 {
		return append(hints, fmt.Sprintf("Run '%s' to replace %s (backed up).", force, plural(len(names), "it", "them")))
	}
	target := "<skill>"
	if len(managed) == 1 {
		target = managed[0]
	}
	action := "to overwrite"
	if command == "remove" {
		action = "to remove anyway"
	}
	hints = append(hints, fmt.Sprintf("Run '%s' to see your changes,", r.command("diff", target)))
	return append(hints, fmt.Sprintf("or '%s' %s (backed up).", force, action))
}

func (r renderer) command(name string, args ...string) string {
	parts := append([]string{"skills", name}, args...)
	if r.global {
		parts = append(parts, "--global")
	}
	for _, agent := range r.agents {
		parts = append(parts, "--agent", agent)
	}
	return strings.Join(parts, " ")
}

func (a app) report(c *cobra.Command, command string, results []install.SkillPlan, err error) error {
	if err != nil && len(results) == 0 {
		return err
	}
	if renderErr := a.renderer(c).results(command, results); renderErr != nil {
		return renderErr
	}
	if err != nil {
		return err
	}
	return attention(results)
}

func attention(results []install.SkillPlan) error {
	for _, res := range results {
		switch res.State {
		case install.StateConflict, install.StateForeign, install.StateUnavailable, install.StateUnsupported:
			return ErrAttention
		}
	}
	return nil
}

func (r renderer) paths(res install.SkillPlan) []string {
	var paths []string
	if res.State == install.StateFailed {
		return nil
	}
	if res.State == install.StateConflict {
		for _, conflict := range res.Conflicts {
			paths = append(paths, r.relative(conflict))
		}
		return paths
	}
	for _, target := range res.Targets {
		for _, change := range target.Files {
			switch change.Action {
			case install.ActionAdd, install.ActionUpdate, install.ActionRemove:
				paths = append(paths, r.relative(filepath.Join(target.Dir, filepath.FromSlash(change.Path))))
			}
		}
	}
	return paths
}

func (r renderer) relative(path string) string {
	rel, err := filepath.Rel(r.root, path)
	if err != nil {
		return path
	}
	return rel
}

func plural(n int, one, many string) string {
	if n == 1 {
		return one
	}
	return many
}

func colorEnabled(noColor bool, noColorEnv string, terminal bool) bool {
	return terminal && !noColor && noColorEnv == ""
}

func isTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	return err == nil && info.Mode()&os.ModeCharDevice != 0
}
