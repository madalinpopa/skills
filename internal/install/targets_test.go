package install_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/config"
	"github.com/madalinpopa/skills/internal/install"
	"github.com/madalinpopa/skills/internal/skill"
)

func TestTargets(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		agents []string
		global bool
		want   []install.Destination
	}{
		"no selection uses the default agents": {
			want: []install.Destination{
				{Dir: filepath.FromSlash("/repo/.agents/skills"), Variant: install.VariantAgents},
				{Dir: filepath.FromSlash("/repo/.claude/skills"), Variant: install.VariantClaude},
			},
		},
		"repeated agents are written once": {
			agents: []string{"claude", "claude"},
			want: []install.Destination{
				{Dir: filepath.FromSlash("/repo/.claude/skills"), Variant: install.VariantClaude},
			},
		},
		"codex and gemini share one target": {
			agents: []string{"codex", "gemini"},
			want: []install.Destination{
				{Dir: filepath.FromSlash("/repo/.agents/skills"), Variant: install.VariantAgents},
			},
		},
		"global scope uses the home paths": {
			agents: []string{"claude", "codex"},
			global: true,
			want: []install.Destination{
				{Dir: filepath.FromSlash("/home/u/.agents/skills"), Variant: install.VariantAgents},
				{Dir: filepath.FromSlash("/home/u/.claude/skills"), Variant: install.VariantClaude},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			root := filepath.FromSlash("/repo")
			if tt.global {
				root = filepath.FromSlash("/home/u")
			}

			got, err := install.Targets(config.Default(), root, tt.global, tt.agents)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTargets_unknownAgent(t *testing.T) {
	t.Parallel()

	_, err := install.Targets(config.Default(), "/repo", false, []string{"claude", "nope"})

	require.ErrorIs(t, err, install.ErrUnknownAgent)
	assert.ErrorContains(t, err, `"nope"`)
	assert.ErrorContains(t, err, "claude, codex, gemini", "the message lists the available agents")
}

func TestTargets_ambiguousConfig(t *testing.T) {
	t.Parallel()
	cfg := config.Default()
	cfg.Agents["codex"] = config.Agent{Project: ".claude/skills", Global: "~/.claude/skills"}

	_, err := install.Targets(cfg, "/repo", false, []string{"claude", "codex"})

	require.Error(t, err)
	assert.ErrorContains(t, err, ".claude/skills", "two variants for one directory is refused, not guessed")
}

func TestSelect_success(t *testing.T) {
	t.Parallel()
	catalog := []skill.Skill{
		{Name: "go-review", Dir: "skills/go-review"},
		{Name: "sql-review", Dir: "skills/sql-review"},
	}

	got, err := install.Select(catalog, []string{"sql-review", "go-review", "sql-review"})

	require.NoError(t, err)
	assert.Equal(t, []skill.Skill{
		{Name: "sql-review", Dir: "skills/sql-review"},
		{Name: "go-review", Dir: "skills/go-review"},
	}, got, "requested order, duplicates dropped")
}

func TestSelect_missing(t *testing.T) {
	t.Parallel()
	catalog := []skill.Skill{{Name: "go-review", Dir: "skills/go-review"}}

	got, err := install.Select(catalog, []string{"go-review", "nope", "also-nope"})

	require.Error(t, err)
	assert.ErrorContains(t, err, "nope")
	assert.ErrorContains(t, err, "also-nope", "every missing name is reported at once")
	assert.Nil(t, got)
}
