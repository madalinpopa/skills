package install_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/v2"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/install"
	"github.com/madalinpopa/skills/internal/skill"
)

const (
	source = "https://example.com/store"
	commit = "a1b2c3d4e5f678901234567890abcdef12345678"
)

var (
	installedAt = time.Date(2026, 9, 12, 18, 40, 0, 0, time.UTC)
	claudeSkill = []byte("---\nname: go-review\ndescription: d\nallowed-tools: Bash(go test *)\n---\n# Go review\n")
	agentsSkill = []byte("---\nname: go-review\ndescription: d\n---\n# Go review\n")
)

func TestInstall_newSkill(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claudeDir := filepath.Join(root, ".claude", "skills", "go-review")
	agentsDir := filepath.Join(root, ".agents", "skills", "go-review")
	req := install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: claudeDir, Files: map[string]skill.File{
			"SKILL.md":         {Data: claudeSkill},
			"scripts/check.sh": {Data: []byte("#!/bin/sh\n"), Mode: 0o755},
		}},
		{Dir: agentsDir, Files: map[string]skill.File{
			"SKILL.md":           {Data: agentsSkill},
			"agents/openai.yaml": {Data: []byte("interface: {}\n")},
		}},
	}}

	results, err := newInstaller().Install([]install.Request{req})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateAdd, results[0].State)
	assert.Equal(t, claudeSkill, read(t, filepath.Join(claudeDir, "SKILL.md")))
	assert.Equal(t, agentsSkill, read(t, filepath.Join(agentsDir, "SKILL.md")))
	assert.FileExists(t, filepath.Join(agentsDir, "agents", "openai.yaml"))
	info, err := os.Stat(filepath.Join(claudeDir, "scripts", "check.sh"))
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&0o100, "bundled scripts stay executable")

	lock := readLock(t, claudeDir)
	assert.Equal(t, "go-review", lock.Name)
	assert.Equal(t, source, lock.Source)
	assert.Equal(t, commit, lock.Commit)
	assert.Equal(t, installedAt, lock.Installed)
	assert.Equal(t, map[string]string{
		"SKILL.md":         sha(claudeSkill),
		"scripts/check.sh": sha([]byte("#!/bin/sh\n")),
	}, lock.Files, "hashes of exactly the files the CLI wrote")
}

func TestInstall_lockIsDeterministic(t *testing.T) {
	t.Parallel()
	first := filepath.Join(t.TempDir(), "go-review")
	second := filepath.Join(t.TempDir(), "go-review")

	_, err := newInstaller().Install([]install.Request{oneTarget(first)})
	require.NoError(t, err)
	_, err = newInstaller().Install([]install.Request{oneTarget(second)})
	require.NoError(t, err)

	assert.Equal(t, read(t, filepath.Join(first, install.LockFile)), read(t, filepath.Join(second, install.LockFile)))
}

func TestInstall_adoptsIdenticalSkill(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claudeDir := filepath.Join(root, ".claude", "skills", "go-review")
	agentsDir := filepath.Join(root, ".agents", "skills", "go-review")
	_, err := newInstaller().Install([]install.Request{oneTarget(claudeDir)})
	require.NoError(t, err)
	claudeLock := read(t, filepath.Join(claudeDir, install.LockFile))
	write(t, filepath.Join(agentsDir, "SKILL.md"), agentsSkill)
	past := installedAt.Add(-24 * time.Hour)
	require.NoError(t, os.Chtimes(filepath.Join(agentsDir, "SKILL.md"), past, past))
	req := install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: claudeDir, Files: map[string]skill.File{"SKILL.md": {Data: claudeSkill}}},
		{Dir: agentsDir, Files: map[string]skill.File{"SKILL.md": {Data: agentsSkill}}},
	}}

	results, err := newInstaller().Install([]install.Request{req})

	require.NoError(t, err)
	assert.Equal(t, install.StateAdd, results[0].State, "a missing lock is work even beside a managed target")
	assert.Equal(t, map[string]string{"SKILL.md": sha(agentsSkill)}, readLock(t, agentsDir).Files)
	assert.Equal(t, claudeLock, read(t, filepath.Join(claudeDir, install.LockFile)))
	info, err := os.Stat(filepath.Join(agentsDir, "SKILL.md"))
	require.NoError(t, err)
	assert.True(t, info.ModTime().Equal(past), "identical content is not rewritten")
	found, err := install.Scan([]install.Destination{
		{Dir: filepath.Dir(claudeDir), Variant: install.VariantClaude},
		{Dir: filepath.Dir(agentsDir), Variant: install.VariantAgents},
	})
	require.NoError(t, err)
	require.Len(t, found, 1)
	assert.Len(t, found[0].Targets, 2, "both targets are discovered afterwards")
}

func TestInstall_dryRunWritesNoMissingLock(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claudeDir := filepath.Join(root, ".claude", "skills", "go-review")
	agentsDir := filepath.Join(root, ".agents", "skills", "go-review")
	_, err := newInstaller().Install([]install.Request{oneTarget(claudeDir)})
	require.NoError(t, err)
	write(t, filepath.Join(agentsDir, "SKILL.md"), agentsSkill)
	req := install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: claudeDir, Files: map[string]skill.File{"SKILL.md": {Data: claudeSkill}}},
		{Dir: agentsDir, Files: map[string]skill.File{"SKILL.md": {Data: agentsSkill}}},
	}}
	dry := newInstaller()
	dry.DryRun = true

	results, err := dry.Install([]install.Request{req})

	require.NoError(t, err)
	assert.Equal(t, install.StateAdd, results[0].State)
	assert.NoFileExists(t, filepath.Join(agentsDir, install.LockFile))
}

func TestInstall_differentUnmanagedSkillIsConflict(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "go-review")
	mine := []byte("# my own go-review\n")
	write(t, filepath.Join(dir, "SKILL.md"), mine)

	results, err := newInstaller().Install([]install.Request{oneTarget(dir)})

	require.NoError(t, err)
	assert.Equal(t, install.StateConflict, results[0].State)
	assert.Equal(t, []string{filepath.Join(dir, "SKILL.md")}, results[0].Conflicts)
	assert.Equal(t, mine, read(t, filepath.Join(dir, "SKILL.md")))
	assert.NoFileExists(t, filepath.Join(dir, install.LockFile))
}

func TestInstall_conflictSkipsAllTargets(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claudeDir := filepath.Join(root, ".claude", "skills", "go-review")
	agentsDir := filepath.Join(root, ".agents", "skills", "go-review")
	write(t, filepath.Join(claudeDir, "SKILL.md"), []byte("# edited\n"))
	req := install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: claudeDir, Files: map[string]skill.File{"SKILL.md": {Data: claudeSkill}}},
		{Dir: agentsDir, Files: map[string]skill.File{"SKILL.md": {Data: agentsSkill}}},
	}}

	results, err := newInstaller().Install([]install.Request{req})

	require.NoError(t, err)
	assert.Equal(t, install.StateConflict, results[0].State)
	assert.NoDirExists(t, agentsDir, "the clean target is skipped too")
}

func TestInstall_otherSkillStillInstalls(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	sqlDir := filepath.Join(root, "sql-review")
	goDir := filepath.Join(root, "go-review")
	write(t, filepath.Join(sqlDir, "SKILL.md"), []byte("# edited\n"))
	reqs := []install.Request{
		{Name: "sql-review", Targets: []install.Desired{{Dir: sqlDir, Files: map[string]skill.File{"SKILL.md": {Data: []byte("# sql\n")}}}}},
		oneTarget(goDir),
	}

	results, err := newInstaller().Install(reqs)

	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, "go-review", results[0].Name)
	assert.Equal(t, install.StateAdd, results[0].State)
	assert.Equal(t, "sql-review", results[1].Name)
	assert.Equal(t, install.StateConflict, results[1].State)
	assert.FileExists(t, filepath.Join(goDir, install.LockFile))
}

func TestInstall_failureLeavesNothing(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claudeDir := filepath.Join(root, ".claude", "skills", "go-review")
	agentsRoot := filepath.Join(root, ".agents", "skills")
	write(t, agentsRoot, []byte("a file where the skills directory should be\n"))
	req := install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: claudeDir, Files: map[string]skill.File{"SKILL.md": {Data: claudeSkill}}},
		{Dir: filepath.Join(agentsRoot, "go-review"), Files: map[string]skill.File{"SKILL.md": {Data: agentsSkill}}},
	}}

	_, err := newInstaller().Install([]install.Request{req})

	require.Error(t, err)
	assert.NoDirExists(t, claudeDir, "no partial skill on the target that would have succeeded")
	assert.NoDirExists(t, filepath.Join(agentsRoot, "go-review"))
}

func TestInstall_noopKeepsLock(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "go-review")
	_, err := newInstaller().Install([]install.Request{oneTarget(dir)})
	require.NoError(t, err)
	lockBefore := read(t, filepath.Join(dir, install.LockFile))
	later := newInstaller()
	later.Now = func() time.Time { return installedAt.Add(time.Hour) }

	results, err := later.Install([]install.Request{oneTarget(dir)})

	require.NoError(t, err)
	assert.Equal(t, install.StateUnchanged, results[0].State)
	assert.Equal(t, lockBefore, read(t, filepath.Join(dir, install.LockFile)), "an unchanged skill keeps its lock")
}

func newInstaller() install.Installer {
	return install.Installer{
		Source: source,
		Commit: commit,
		Now:    func() time.Time { return installedAt },
	}
}

func oneTarget(dir string) install.Request {
	return install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: dir, Files: map[string]skill.File{"SKILL.md": {Data: claudeSkill}}},
	}}
}

func readLock(t *testing.T, dir string) install.Lock {
	t.Helper()
	var lock install.Lock
	require.NoError(t, json.Unmarshal(read(t, filepath.Join(dir, install.LockFile)), &lock))
	return lock
}

func read(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return data
}

func write(t *testing.T, path string, data []byte) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
	require.NoError(t, os.WriteFile(path, data, 0o600))
}

func sha(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
