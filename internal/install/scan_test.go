package install_test

import (
	"encoding/json/v2"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/config"
	"github.com/madalinpopa/skills/internal/install"
)

const oldCommit = "0000000000000000000000000000000000000000"

func TestScan_findsManagedSkills(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claude := filepath.Join(root, ".claude", "skills")
	agents := filepath.Join(root, ".agents", "skills")
	installed(t, filepath.Join(claude, "go-review"), map[string][]byte{"SKILL.md": claudeSkill})
	installed(t, filepath.Join(agents, "go-review"), map[string][]byte{"SKILL.md": agentsSkill})
	installed(t, filepath.Join(agents, "sql-review"), map[string][]byte{
		"SKILL.md": []byte("---\nname: sql-review\ndescription: Reviews SQL.\n---\n# SQL\n"),
	})
	write(t, filepath.Join(claude, "unmanaged", "SKILL.md"), agentsSkill)
	write(t, filepath.Join(claude, "notes.txt"), []byte("not a skill\n"))
	dests, err := install.Targets(config.Default(), root, false, nil)
	require.NoError(t, err)

	found, err := install.Scan(dests)

	require.NoError(t, err)
	assert.Equal(t, []install.Installation{
		{Name: "go-review", Description: "d", Targets: []install.Destination{
			{Dir: agents, Variant: install.VariantAgents},
			{Dir: claude, Variant: install.VariantClaude},
		}},
		{Name: "sql-review", Description: "Reviews SQL.", Targets: []install.Destination{
			{Dir: agents, Variant: install.VariantAgents},
		}},
	}, found, "one entry per skill, one target per physical directory")
}

func TestScan_nothingInstalled(t *testing.T) {
	t.Parallel()
	dests, err := install.Targets(config.Default(), t.TempDir(), false, nil)
	require.NoError(t, err)

	found, err := install.Scan(dests)

	require.NoError(t, err)
	assert.Empty(t, found)
}

func TestScan_corruptLockIsError(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	dir := filepath.Join(root, ".claude", "skills", "broken")
	write(t, filepath.Join(dir, "SKILL.md"), claudeSkill)
	write(t, filepath.Join(dir, install.LockFile), []byte("not json"))
	dests, err := install.Targets(config.Default(), root, false, []string{"claude"})
	require.NoError(t, err)

	_, err = install.Scan(dests)

	require.Error(t, err)
	assert.Contains(t, err.Error(), filepath.Join(dir, install.LockFile))
}

func installed(t *testing.T, dir string, files map[string][]byte) {
	t.Helper()
	installedFrom(t, dir, source, files)
}

func installedFrom(t *testing.T, dir, from string, files map[string][]byte) {
	t.Helper()
	lock := install.Lock{
		Name:      filepath.Base(dir),
		Source:    from,
		Commit:    oldCommit,
		Installed: installedAt.Add(-24 * time.Hour),
		Files:     map[string]string{},
	}
	for p, data := range files {
		write(t, filepath.Join(dir, p), data)
		lock.Files[p] = sha(data)
	}
	data, err := json.Marshal(lock)
	require.NoError(t, err)
	write(t, filepath.Join(dir, install.LockFile), data)
}
