package install_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/install"
)

func TestInstall_damagedLockIsUnsupported(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	write(t, filepath.Join(dir, "SKILL.md"), []byte("# mine\n"))
	write(t, filepath.Join(dir, install.LockFile), []byte("not json"))
	installer := newInstaller()
	installer.Backups = filepath.Join(t.TempDir(), "backups")
	installer.Force = true

	results, err := installer.Install([]install.Request{oneTarget(dir)})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateUnsupported, results[0].State, "force must not adopt a damaged lock")
	require.Len(t, results[0].Issues, 1)
	assert.Equal(t, filepath.Join(dir, install.LockFile), results[0].Issues[0].Path)
	assert.Equal(t, []byte("# mine\n"), read(t, filepath.Join(dir, "SKILL.md")))
	assert.Equal(t, []byte("not json"), read(t, filepath.Join(dir, install.LockFile)))
	assert.NoDirExists(t, installer.Backups)
}

func TestRemove_damagedLockIsUnsupported(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	write(t, filepath.Join(dir, "SKILL.md"), claudeSkill)
	write(t, filepath.Join(dir, install.LockFile), []byte("not json"))
	remover := newInstaller()
	remover.Backups = filepath.Join(t.TempDir(), "backups")
	remover.Force = true

	results, err := remover.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: filepath.Dir(dir), Variant: install.VariantClaude}}},
	})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateUnsupported, results[0].State)
	assert.Equal(t, filepath.Join(dir, install.LockFile), results[0].Issues[0].Path)
	assert.FileExists(t, filepath.Join(dir, "SKILL.md"), "nothing is removed without a valid lock")
	assert.NoDirExists(t, remover.Backups)
}
