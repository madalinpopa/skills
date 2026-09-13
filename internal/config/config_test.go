package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/config"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	cfg := config.Default()

	assert.Equal(t, "https://github.com/madalinpopa/skills", cfg.Store.Repo)
	assert.Equal(t, "main", cfg.Store.Branch)
	assert.Equal(t, map[string]config.Agent{
		"claude": {Project: ".claude/skills", Global: "~/.claude/skills"},
		"codex":  {Project: ".agents/skills", Global: "~/.agents/skills"},
	}, cfg.Agents)
	assert.Equal(t, []string{"claude", "codex"}, cfg.Defaults.Agents)
}

func TestDir(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		xdgConfigHome string
		home          string
		want          string
	}{
		"xdg set":    {xdgConfigHome: "/xdg", home: "/home/u", want: "/xdg/skills"},
		"xdg absent": {xdgConfigHome: "", home: "/home/u", want: "/home/u/.config/skills"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := config.Dir(tt.xdgConfigHome, tt.home)

			assert.Equal(t, filepath.FromSlash(tt.want), got)
		})
	}
}

func TestLoad_overridesDefaults(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeConfig(t, dir, `
[store]
repo = "https://example.com/team/skills"

[agents.cursor]
project = ".cursor/skills"
global = "~/.cursor/skills"

[defaults]
agents = ["cursor"]
`)

	cfg, err := config.Load(dir)

	require.NoError(t, err)
	assert.Equal(t, "https://example.com/team/skills", cfg.Store.Repo)
	assert.Equal(t, "main", cfg.Store.Branch, "unset keys keep their defaults")
	assert.Equal(t, map[string]config.Agent{
		"cursor": {Project: ".cursor/skills", Global: "~/.cursor/skills"},
	}, cfg.Agents, "an agents table replaces the built-in agents")
	assert.Equal(t, []string{"cursor"}, cfg.Defaults.Agents)
}

func TestLoad_missingFileIsReadOnly(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "skills")

	cfg, err := config.Load(dir)

	require.NoError(t, err)
	assert.Equal(t, config.Default(), cfg)
	assert.NoDirExists(t, dir, "a read-only load creates nothing")
}

func TestLoad_failed(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		toml string
		want string
	}{
		"malformed toml": {toml: `[store\nrepo = `, want: "parse"},
		"unknown default agent": {toml: `
[defaults]
agents = ["nope"]
`, want: "nope"},
		"empty repo": {toml: `
[store]
repo = ""
`, want: "store.repo"},
		"empty default agents": {toml: `
[defaults]
agents = []
`, want: "defaults.agents"},
		"absolute project path": {toml: `
[defaults]
agents = ["claude"]

[agents.claude]
project = "/srv/skills"
global = "~/.claude/skills"
`, want: "claude"},
		"escaping project path": {toml: `
[defaults]
agents = ["claude"]

[agents.claude]
project = "../skills"
global = "~/.claude/skills"
`, want: "claude"},
		"relative global path": {toml: `
[defaults]
agents = ["claude"]

[agents.claude]
project = ".claude/skills"
global = "skills"
`, want: "claude"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			writeConfig(t, dir, tt.toml)

			_, err := config.Load(dir)

			require.Error(t, err)
			assert.ErrorContains(t, err, tt.want)
		})
	}
}

func TestLoad_acceptsAbsoluteGlobalPath(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	writeConfig(t, dir, `
[defaults]
agents = ["claude"]

[agents.claude]
project = ".claude/skills"
global = "/srv/skills"
`)

	cfg, err := config.Load(dir)

	require.NoError(t, err)
	assert.Equal(t, "/srv/skills", cfg.Agents["claude"].GlobalPath("/home/me"), "an absolute global path is used as is")
}

func TestInit_createsDefaultFile(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "skills")

	cfg, err := config.Init(dir)

	require.NoError(t, err)
	assert.Equal(t, config.Default(), cfg)
	written, err := os.ReadFile(filepath.Join(dir, "config.toml"))
	require.NoError(t, err)
	assert.NotContains(t, string(written), "gemini", "the written file offers only claude and codex")

	loaded, err := config.Load(dir)
	require.NoError(t, err)
	assert.Equal(t, config.Default(), loaded, "the written file parses back to the defaults")
}

func TestInit_keepsExistingFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	custom := `[store]
repo = "https://example.com/mine"

[agents.claude]
project = ".claude/skills"
global  = "~/.claude/skills"

[agents.gemini]
project = ".agents/skills"
global  = "~/.agents/skills"

[defaults]
agents = ["claude", "gemini"]
`
	writeConfig(t, dir, custom)

	cfg, err := config.Init(dir)

	require.NoError(t, err)
	assert.Equal(t, "https://example.com/mine", cfg.Store.Repo)
	assert.Equal(t, config.Agent{Project: ".agents/skills", Global: "~/.agents/skills"}, cfg.Agents["gemini"], "an explicitly configured gemini agent keeps working")
	assert.Equal(t, []string{"claude", "gemini"}, cfg.Defaults.Agents)
	got, err := os.ReadFile(filepath.Join(dir, "config.toml"))
	require.NoError(t, err)
	assert.Equal(t, custom, string(got), "init never rewrites an existing config")
}

func TestAgentGlobalPath(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		global string
		want   string
	}{
		"home prefix": {global: "~/.claude/skills", want: "/home/u/.claude/skills"},
		"absolute":    {global: "/srv/skills", want: "/srv/skills"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			agent := config.Agent{Global: tt.global}

			got := agent.GlobalPath("/home/u")

			assert.Equal(t, filepath.FromSlash(tt.want), got)
		})
	}
}

func writeConfig(t *testing.T, dir, toml string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "config.toml"), []byte(toml), 0o600))
}
