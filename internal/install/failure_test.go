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

func TestInstall_publishFailureKeepsProgressAndOldLock(t *testing.T) {
	t.Parallel()
	skipAsRoot(t)
	root := t.TempDir()
	claudeDir := filepath.Join(root, ".claude", "skills", "go-review")
	agentsDir := filepath.Join(root, ".agents", "skills", "go-review")
	v1 := map[string][]byte{"SKILL.md": []byte("# v1\n"), "scripts/old.sh": []byte("#!/bin/sh\n")}
	installed(t, claudeDir, v1)
	installed(t, agentsDir, v1)
	readOnly(t, filepath.Join(agentsDir, "scripts"))
	v2 := map[string]skill.File{"SKILL.md": {Data: []byte("# v2\n")}}
	sqlDir := filepath.Join(root, "sql-review")
	reqs := []install.Request{
		{Name: "go-review", Targets: []install.Desired{{Dir: claudeDir, Files: v2}, {Dir: agentsDir, Files: v2}}},
		{Name: "sql-review", Targets: []install.Desired{{Dir: sqlDir, Files: map[string]skill.File{"SKILL.md": {Data: []byte("# sql\n")}}}}},
	}

	results, err := newInstaller().Install(reqs)

	require.Error(t, err)
	assert.ErrorContains(t, err, "go-review")
	assert.ErrorContains(t, err, agentsDir, "the error names the affected target")
	require.Len(t, results, 1, "completed and failed skills are returned, later skills are not started")
	assert.Equal(t, install.StateFailed, results[0].State)
	assert.Equal(t, []byte("# v2\n"), read(t, filepath.Join(claudeDir, "SKILL.md")), "the first target was published in full")
	assert.NoFileExists(t, filepath.Join(claudeDir, "scripts", "old.sh"))
	assert.Equal(t, commit, readLock(t, claudeDir).Commit)
	assert.Equal(t, oldCommit, readLock(t, agentsDir).Commit, "a lock is published only after its removals succeed")
	entries, err := os.ReadDir(agentsDir)
	require.NoError(t, err)
	assert.Len(t, entries, 3, "no staged leftovers next to the partial progress")
	assert.NoDirExists(t, sqlDir, "mutations stop at the first runtime failure")
}

func TestRemove_laterBackupFailureKeepsEarlierResults(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claude := filepath.Join(root, ".claude", "skills")
	agents := filepath.Join(root, ".agents", "skills")
	installed(t, filepath.Join(claude, "go-review"), map[string][]byte{"SKILL.md": claudeSkill})
	installed(t, filepath.Join(claude, "sql-review"), map[string][]byte{"SKILL.md": []byte("# sql\n")})
	installed(t, filepath.Join(agents, "sql-review"), map[string][]byte{"SKILL.md": []byte("# sql\n")})
	backups := filepath.Join(t.TempDir(), "backups")
	write(t, filepath.Join(backups, "20260912-184000", identity(t, filepath.Join(agents, "sql-review"))), []byte("in the way\n"))
	remover := newInstaller()
	remover.Backups = backups

	results, err := remover.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: claude, Variant: install.VariantClaude}}},
		{Name: "sql-review", Targets: []install.Destination{{Dir: claude, Variant: install.VariantClaude}, {Dir: agents, Variant: install.VariantAgents}}},
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "sql-review")
	require.Len(t, results, 2, "the earlier removal stays visible")
	assert.Equal(t, install.StateRemove, results[0].State)
	assert.Len(t, results[0].Backups, 1)
	assert.NoDirExists(t, filepath.Join(claude, "go-review"))
	assert.Equal(t, install.StateFailed, results[1].State)
	require.Len(t, results[1].Backups, 1, "the copy made before the failure is reported")
	assert.FileExists(t, filepath.Join(results[1].Backups[0], "SKILL.md"))
	assert.FileExists(t, filepath.Join(claude, "sql-review", install.LockFile), "a backup failure leaves the skill installed")
	assert.FileExists(t, filepath.Join(agents, "sql-review", install.LockFile))
}

func identity(t *testing.T, dir string) string {
	t.Helper()
	abs, err := filepath.Abs(dir)
	require.NoError(t, err)
	return strings.TrimPrefix(strings.TrimPrefix(abs, filepath.VolumeName(abs)), string(filepath.Separator))
}

func readOnly(t *testing.T, dir string) {
	t.Helper()
	require.NoError(t, os.Chmod(dir, 0o555))
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
}

func skipAsRoot(t *testing.T) {
	t.Helper()
	if os.Getuid() == 0 {
		t.Skip("directory permissions are not enforced for root")
	}
}
