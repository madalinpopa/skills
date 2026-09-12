package install_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/install"
	"github.com/madalinpopa/skills/internal/skill"
)

func TestInstall_forceKeepsLocalOnlyFiles(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": []byte("# v1\n")})
	write(t, filepath.Join(dir, "SKILL.md"), []byte("# mine\n"))
	write(t, filepath.Join(dir, "notes.md"), []byte("local only\n"))
	installer := newInstaller()
	installer.Backups = filepath.Join(t.TempDir(), "backups")
	installer.Force = true
	req := install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: dir, Files: map[string]skill.File{"SKILL.md": {Data: []byte("# v2\n")}}},
	}}

	results, err := installer.Install([]install.Request{req})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateUpdate, results[0].State)
	assert.Equal(t, []byte("# v2\n"), read(t, filepath.Join(dir, "SKILL.md")), "the conflict is resolved with the store version")
	assert.Equal(t, []byte("local only\n"), read(t, filepath.Join(dir, "notes.md")), "local-only files are user data, not conflicts")
	assert.Equal(t, map[string]string{"SKILL.md": sha([]byte("# v2\n"))}, readLock(t, dir).Files, "the lock tracks only what the CLI wrote")
	require.Len(t, results[0].Backups, 1)
	assert.Equal(t, []byte("# mine\n"), read(t, filepath.Join(results[0].Backups[0], "SKILL.md")))
}

func TestInstall_forceKeepsForeignSource(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installedFrom(t, dir, "https://example.com/other-store", map[string][]byte{"SKILL.md": []byte("# v1\n")})
	installer := newInstaller()
	installer.Backups = filepath.Join(t.TempDir(), "backups")
	installer.Force = true
	req := install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: dir, Files: map[string]skill.File{"SKILL.md": {Data: []byte("# v2\n")}}},
	}}

	results, err := installer.Install([]install.Request{req})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateForeign, results[0].State, "force is not a bypass for another store's installation")
	assert.Equal(t, []byte("# v1\n"), read(t, filepath.Join(dir, "SKILL.md")))
	assert.NoDirExists(t, installer.Backups)
}

func TestBackup_sameSecondCollisionLeavesEarlierBackupUntouched(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claude := filepath.Join(root, ".claude", "skills")
	dir := filepath.Join(claude, "go-review")
	write(t, filepath.Join(dir, "SKILL.md"), []byte("# my own go-review\n"))
	installer := newInstaller()
	installer.Backups = filepath.Join(t.TempDir(), "backups")
	installer.Force = true
	results, err := installer.Install([]install.Request{oneTarget(dir)})
	require.NoError(t, err)
	require.Len(t, results[0].Backups, 1)
	earlier := results[0].Backups[0]

	results, err = installer.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: claude, Variant: install.VariantClaude}}},
	})

	require.Error(t, err, "a second backup in the same second must not reuse the path")
	assert.ErrorContains(t, err, earlier)
	assert.Equal(t, install.StateFailed, results[0].State)
	assert.FileExists(t, filepath.Join(dir, install.LockFile), "the skill stays installed without a backup")
	entries, err := os.ReadDir(earlier)
	require.NoError(t, err)
	require.Len(t, entries, 1, "the earlier backup gains no files from the failed copy")
	assert.Equal(t, "SKILL.md", entries[0].Name())
	assert.Equal(t, []byte("# my own go-review\n"), read(t, filepath.Join(earlier, "SKILL.md")))
}
