package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/install"
)

func TestRender_default(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	var out bytes.Buffer
	r := renderer{out: &out, root: root}

	require.NoError(t, r.results("update", []install.SkillPlan{
		{Name: "django-testing", State: install.StateUnavailable},
		{Name: "go-review", State: install.StateAdd, Targets: []install.TargetPlan{
			{Dir: filepath.Join(root, ".agents", "skills", "go-review"), Files: []install.FileChange{{Path: "SKILL.md", Action: install.ActionAdd}}},
			{Dir: filepath.Join(root, ".claude", "skills", "go-review"), Files: []install.FileChange{{Path: "SKILL.md", Action: install.ActionAdd}}},
		}},
		{Name: "old-skill", State: install.StateRemove, Backups: []string{filepath.Join(root, "backups", "old-skill")}},
		{Name: "sql-review", State: install.StateUpdate},
		{Name: "unchanged", State: install.StateUnchanged},
	}))

	assert.Equal(t, []string{
		"  ! django-testing   unavailable in store, left installed",
		"  + go-review        added",
		"  - old-skill        removed",
		"      backed up to " + filepath.Join(root, "backups", "old-skill"),
		"  ~ sql-review       updated",
		"",
		"  5 skills, 3 changed, 1 needs attention",
	}, lines(out.String()), "one line per skill, shared targets aggregated, unchanged skills silent")
}

func TestRender_attention(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	var out bytes.Buffer
	r := renderer{out: &out, root: root}

	require.NoError(t, r.results("update", []install.SkillPlan{
		{Name: "go-review", State: install.StateAdd},
		{Name: "other", State: install.StateForeign, Source: "https://example.com/other"},
		{Name: "sql-review", State: install.StateConflict, Conflicts: []string{filepath.Join(root, ".claude", "skills", "sql-review", "SKILL.md")}},
	}))

	assert.Equal(t, []string{
		"  + go-review    added",
		"  ! other        skipped, installed from https://example.com/other",
		"  ! sql-review   skipped, you edited it",
		"",
		"  3 skills, 1 changed, 2 need attention",
		"  Run 'skills diff sql-review' to see your changes,",
		"  or 'skills update --force' to overwrite (backed up).",
	}, lines(out.String()))
}

func TestRender_unsupported(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	var out bytes.Buffer
	r := renderer{out: &out, root: root}
	plans := []install.SkillPlan{
		{Name: "go-review", State: install.StateUnsupported, Issues: []install.Issue{
			{Path: filepath.Join(root, ".claude", "skills", "go-review", "references"), Reason: "is a symlink"},
		}},
	}

	require.NoError(t, r.results("update", plans))

	assert.Equal(t, []string{
		"  ! go-review   skipped, needs manual repair",
		"      " + filepath.Join(".claude", "skills", "go-review", "references") + " is a symlink",
		"",
		"  1 skill, 0 changed, 1 needs attention",
	}, lines(out.String()), "the path and reason are shown, and force is not suggested")
	assert.ErrorIs(t, attention(plans), ErrAttention)
}

func TestRender_failed(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	var out bytes.Buffer
	r := renderer{out: &out, root: root}

	require.NoError(t, r.results("remove", []install.SkillPlan{
		{Name: "go-review", State: install.StateRemove, Backups: []string{filepath.Join(root, "backups", "go-review")}},
		{Name: "sql-review", State: install.StateFailed, Backups: []string{filepath.Join(root, "backups", "sql-review")}},
	}))

	assert.Equal(t, []string{
		"  - go-review    removed",
		"      backed up to " + filepath.Join(root, "backups", "go-review"),
		"  ! sql-review   failed, see the error below",
		"      backed up to " + filepath.Join(root, "backups", "sql-review"),
		"",
		"  2 skills, 1 changed, 1 failed",
	}, lines(out.String()), "a failed skill is never counted as changed")
}

func TestRender_verbose(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	var out bytes.Buffer
	r := renderer{out: &out, root: root, verbose: true}

	require.NoError(t, r.results("install", []install.SkillPlan{
		{Name: "go-review", State: install.StateAdd, Targets: []install.TargetPlan{
			{Dir: filepath.Join(root, ".agents", "skills", "go-review"), Files: []install.FileChange{
				{Path: "SKILL.md", Action: install.ActionAdd},
				{Path: "agents/openai.yaml", Action: install.ActionAdd},
			}},
			{Dir: filepath.Join(root, ".claude", "skills", "go-review"), Files: []install.FileChange{
				{Path: "SKILL.md", Action: install.ActionAdd},
				{Path: "notes.md", Action: install.ActionKeep},
			}},
		}},
		{Name: "sql-review", State: install.StateConflict, Conflicts: []string{filepath.Join(root, ".claude", "skills", "sql-review", "SKILL.md")}},
	}))

	assert.Equal(t, []string{
		"  + go-review    added",
		"      " + filepath.Join(".agents", "skills", "go-review", "SKILL.md"),
		"      " + filepath.Join(".agents", "skills", "go-review", "agents", "openai.yaml"),
		"      " + filepath.Join(".claude", "skills", "go-review", "SKILL.md"),
		"  ! sql-review   skipped, you edited it",
		"      " + filepath.Join(".claude", "skills", "sql-review", "SKILL.md"),
		"",
		"  2 skills, 1 changed, 1 needs attention",
		"  Run 'skills diff sql-review' to see your changes,",
		"  or 'skills install --force' to overwrite (backed up).",
	}, lines(out.String()), "paths appear beneath the skill; kept local files are not changes")
}

func TestRender_color(t *testing.T) {
	t.Parallel()
	results := []install.SkillPlan{{Name: "go-review", State: install.StateAdd}}
	var colored, plain bytes.Buffer

	require.NoError(t, renderer{out: &colored, color: true}.results("install", results))
	require.NoError(t, renderer{out: &plain, color: false}.results("install", results))

	assert.Contains(t, colored.String(), "\x1b[")
	assert.Contains(t, colored.String(), "+ go-review")
	assert.NotContains(t, plain.String(), "\x1b[")
	assert.Contains(t, plain.String(), "+ go-review", "the symbol carries the meaning without color")
}

func TestColorEnabled(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		noColor  bool
		env      string
		terminal bool
		want     bool
	}{
		"terminal":            {terminal: true, want: true},
		"not a terminal":      {terminal: false, want: false},
		"NO_COLOR set":        {terminal: true, env: "1", want: false},
		"NO_COLOR empty":      {terminal: true, env: "", want: true},
		"--no-color flag":     {terminal: true, noColor: true, want: false},
		"flag and no console": {terminal: false, noColor: true, want: false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, colorEnabled(tt.noColor, tt.env, tt.terminal))
		})
	}
}

func lines(s string) []string {
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}
