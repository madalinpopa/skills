package install_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/install"
	"github.com/madalinpopa/skills/internal/skill"
)

func TestRemove_backsUpThenRemoves(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claude := filepath.Join(root, ".claude", "skills")
	dir := filepath.Join(claude, "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": claudeSkill, "scripts/check.sh": []byte("#!/bin/sh\n")})
	installed(t, filepath.Join(claude, "sql-review"), map[string][]byte{"SKILL.md": []byte("# sql\n")})
	backups := filepath.Join(t.TempDir(), "backups")
	remover := newInstaller()
	remover.Backups = backups

	results, err := remover.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: claude, Variant: install.VariantClaude}}},
	})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateRemove, results[0].State)
	assert.NoDirExists(t, dir)
	assert.FileExists(t, filepath.Join(claude, "sql-review", "SKILL.md"), "siblings are untouched")
	assert.DirExists(t, claude, "the target directory stays")
	require.Len(t, results[0].Backups, 1, "one backup per removed target")
	backup := results[0].Backups[0]
	assert.True(t, strings.HasPrefix(backup, backups), "backups live under the configured root")
	assert.True(t, strings.HasSuffix(backup, filepath.Join(".claude", "skills", "go-review")), "the backup keeps the target identity")
	assert.Equal(t, claudeSkill, read(t, filepath.Join(backup, "SKILL.md")))
	assert.FileExists(t, filepath.Join(backup, "scripts", "check.sh"))
	assert.FileExists(t, filepath.Join(backup, install.LockFile))
}

func TestRemove_editedNeedsAttention(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": claudeSkill})
	write(t, filepath.Join(dir, "SKILL.md"), []byte("# mine\n"))
	backups := filepath.Join(t.TempDir(), "backups")
	remover := newInstaller()
	remover.Backups = backups

	results, err := remover.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: filepath.Dir(dir), Variant: install.VariantClaude}}},
	})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateConflict, results[0].State)
	assert.Equal(t, []string{filepath.Join(dir, "SKILL.md")}, results[0].Conflicts)
	assert.Equal(t, []byte("# mine\n"), read(t, filepath.Join(dir, "SKILL.md")))
	assert.NoDirExists(t, backups, "nothing is backed up when nothing is removed")
}

func TestRemove_forceBacksUpEditedContent(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": claudeSkill})
	write(t, filepath.Join(dir, "SKILL.md"), []byte("# mine\n"))
	write(t, filepath.Join(dir, "notes.md"), []byte("local only\n"))
	remover := newInstaller()
	remover.Backups = filepath.Join(t.TempDir(), "backups")
	remover.Force = true

	results, err := remover.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: filepath.Dir(dir), Variant: install.VariantClaude}}},
	})

	require.NoError(t, err)
	assert.Equal(t, install.StateRemove, results[0].State)
	assert.NoDirExists(t, dir)
	require.Len(t, results[0].Backups, 1)
	assert.Equal(t, []byte("# mine\n"), read(t, filepath.Join(results[0].Backups[0], "SKILL.md")), "the backup holds your edit, not the store version")
	assert.Equal(t, []byte("local only\n"), read(t, filepath.Join(results[0].Backups[0], "notes.md")))
}

func TestRemove_foreignSourceIsRemovable(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installedFrom(t, dir, "https://example.com/other-store", map[string][]byte{"SKILL.md": claudeSkill})
	remover := newInstaller()
	remover.Backups = filepath.Join(t.TempDir(), "backups")

	results, err := remover.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: filepath.Dir(dir), Variant: install.VariantClaude}}},
	})

	require.NoError(t, err)
	assert.Equal(t, install.StateRemove, results[0].State)
	assert.NoDirExists(t, dir)
}

func TestRemove_backupFailurePreventsRemoval(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": claudeSkill})
	backups := filepath.Join(t.TempDir(), "backups")
	write(t, backups, []byte("a file where the backup root should be\n"))
	remover := newInstaller()
	remover.Backups = backups

	_, err := remover.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: filepath.Dir(dir), Variant: install.VariantClaude}}},
	})

	require.Error(t, err)
	assert.Equal(t, claudeSkill, read(t, filepath.Join(dir, "SKILL.md")), "no backup means no removal")
	assert.FileExists(t, filepath.Join(dir, install.LockFile))
}

func TestInstall_forceBacksUpBeforeReplacing(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": []byte("# v1\n")})
	write(t, filepath.Join(dir, "SKILL.md"), []byte("# mine\n"))
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
	assert.Equal(t, []byte("# v2\n"), read(t, filepath.Join(dir, "SKILL.md")))
	assert.Equal(t, commit, readLock(t, dir).Commit)
	require.Len(t, results[0].Backups, 1)
	assert.Equal(t, []byte("# mine\n"), read(t, filepath.Join(results[0].Backups[0], "SKILL.md")), "the old content is backed up before replacement")
}

func TestInstall_forceBackupFailurePreventsReplacement(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": []byte("# v1\n")})
	write(t, filepath.Join(dir, "SKILL.md"), []byte("# mine\n"))
	backups := filepath.Join(t.TempDir(), "backups")
	write(t, backups, []byte("a file where the backup root should be\n"))
	installer := newInstaller()
	installer.Backups = backups
	installer.Force = true
	req := install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: dir, Files: map[string]skill.File{"SKILL.md": {Data: []byte("# v2\n")}}},
	}}

	results, err := installer.Install([]install.Request{req})

	require.Error(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateFailed, results[0].State, "a failed backup is reported, not hidden")
	assert.Equal(t, []byte("# mine\n"), read(t, filepath.Join(dir, "SKILL.md")))
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Len(t, entries, 2, "only the skill file and its lock, no staged leftovers")
}

func TestRemove_deletedTrackedFileIsConflict(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": claudeSkill, "scripts/check.sh": []byte("#!/bin/sh\n")})
	require.NoError(t, os.Remove(filepath.Join(dir, "scripts", "check.sh")))
	remover := newInstaller()
	remover.Backups = filepath.Join(t.TempDir(), "backups")

	results, err := remover.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: filepath.Dir(dir), Variant: install.VariantClaude}}},
	})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateConflict, results[0].State, "a deleted tracked file is an edit")
	assert.Equal(t, []string{filepath.Join(dir, "scripts", "check.sh")}, results[0].Conflicts)
	assert.FileExists(t, filepath.Join(dir, "SKILL.md"), "the rest of the skill stays installed")
	assert.NoDirExists(t, remover.Backups)
}

func TestRemove_forceBacksUpRemainingFiles(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": claudeSkill, "scripts/check.sh": []byte("#!/bin/sh\n")})
	require.NoError(t, os.Remove(filepath.Join(dir, "scripts", "check.sh")))
	remover := newInstaller()
	remover.Backups = filepath.Join(t.TempDir(), "backups")
	remover.Force = true

	results, err := remover.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: filepath.Dir(dir), Variant: install.VariantClaude}}},
	})

	require.NoError(t, err)
	assert.Equal(t, install.StateRemove, results[0].State)
	assert.NoDirExists(t, dir)
	require.Len(t, results[0].Backups, 1)
	assert.Equal(t, claudeSkill, read(t, filepath.Join(results[0].Backups[0], "SKILL.md")), "what remained is backed up")
	assert.NoFileExists(t, filepath.Join(results[0].Backups[0], "scripts", "check.sh"), "a backup never invents deleted content")
}
