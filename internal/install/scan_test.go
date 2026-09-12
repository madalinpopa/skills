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
		{Name: "go-review", Targets: []install.Destination{
			{Dir: agents, Variant: install.VariantAgents},
			{Dir: claude, Variant: install.VariantClaude},
		}},
		{Name: "sql-review", Targets: []install.Destination{
			{Dir: agents, Variant: install.VariantAgents},
		}},
	}, found, "one entry per skill, one target per physical directory")
}

func TestScan_ignoresLocalMetadata(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	dir := filepath.Join(root, ".claude", "skills", "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": []byte("# no frontmatter\n")})
	dests, err := install.Targets(config.Default(), root, false, []string{"claude"})
	require.NoError(t, err)

	found, err := install.Scan(dests)

	require.NoError(t, err)
	require.Len(t, found, 1, "a valid lock is the identity, not the frontmatter")
	assert.Empty(t, found[0].Issues)
	_, err = install.Describe(found[0])
	require.Error(t, err, "the description is optional display data")
	assert.ErrorContains(t, err, "frontmatter")
}

func TestScan_nothingInstalled(t *testing.T) {
	t.Parallel()
	dests, err := install.Targets(config.Default(), t.TempDir(), false, nil)
	require.NoError(t, err)

	found, err := install.Scan(dests)

	require.NoError(t, err)
	assert.Empty(t, found)
}

func TestScan_damagedLockIsReportedBesideHealthy(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claude := filepath.Join(root, ".claude", "skills")
	broken := filepath.Join(claude, "broken")
	write(t, filepath.Join(broken, "SKILL.md"), claudeSkill)
	write(t, filepath.Join(broken, install.LockFile), []byte("not json"))
	installed(t, filepath.Join(claude, "go-review"), map[string][]byte{"SKILL.md": claudeSkill})
	dests, err := install.Targets(config.Default(), root, false, []string{"claude"})
	require.NoError(t, err)

	found, err := install.Scan(dests)

	require.NoError(t, err, "a damaged lock is attention, not a runtime failure")
	require.Len(t, found, 2)
	assert.Equal(t, "broken", found[0].Name)
	require.Len(t, found[0].Issues, 1)
	assert.Equal(t, filepath.Join(broken, install.LockFile), found[0].Issues[0].Path)
	assert.NotEmpty(t, found[0].Issues[0].Reason)
	assert.Equal(t, "go-review", found[1].Name)
	assert.Empty(t, found[1].Issues, "healthy installations are unaffected")
}

func TestReadLock_damaged(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		lock string
		want string
	}{
		"not json":       {lock: "not json", want: install.LockFile},
		"wrong name":     {lock: `{"name":"other","source":"` + source + `","commit":"` + commit + `","files":{}}`, want: "other"},
		"missing source": {lock: `{"name":"go-review","commit":"` + commit + `","files":{}}`, want: "source"},
		"missing commit": {lock: `{"name":"go-review","source":"` + source + `","files":{}}`, want: "commit"},
		"escaping path":  {lock: `{"name":"go-review","source":"` + source + `","commit":"` + commit + `","files":{"../x":"h"}}`, want: "../x"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := filepath.Join(t.TempDir(), "go-review")
			write(t, filepath.Join(dir, install.LockFile), []byte(tt.lock))

			_, err := install.ReadLock(dir)

			require.ErrorIs(t, err, install.ErrDamagedLock)
			assert.ErrorContains(t, err, tt.want)
		})
	}
}

func installed(t *testing.T, dir string, files map[string][]byte) {
	t.Helper()
	installedFrom(t, dir, source, files)
}

func installedFrom(t *testing.T, dir, from string, files map[string][]byte) {
	t.Helper()
	lock := install.Lock{
		Name:       filepath.Base(dir),
		Source:     from,
		Commit:     oldCommit,
		Installed:  installedAt.Add(-24 * time.Hour),
		Files:      map[string]string{},
		Executable: map[string]bool{},
	}
	for p, data := range files {
		write(t, filepath.Join(dir, p), data)
		lock.Files[p] = sha(data)
		lock.Executable[p] = false
	}
	data, err := json.Marshal(lock)
	require.NoError(t, err)
	write(t, filepath.Join(dir, install.LockFile), data)
}
