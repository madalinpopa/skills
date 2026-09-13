package install_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json/v2"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
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
	assert.Equal(t, map[string]bool{"SKILL.md": false, "scripts/check.sh": true}, lock.Executable, "one mode per tracked file")
}

func TestInstall_appliesExecutableChange(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "go-review")
	script := []byte("#!/bin/sh\n")
	_, err := newInstaller().Install([]install.Request{withScript(dir, script, 0o644)})
	require.NoError(t, err)

	results, err := newInstaller().Install([]install.Request{withScript(dir, script, 0o755)})

	require.NoError(t, err)
	assert.Equal(t, install.StateUpdate, results[0].State, "a mode-only upstream change is applied")
	info, err := os.Stat(filepath.Join(dir, "scripts", "check.sh"))
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&0o100)
	assert.Equal(t, map[string]bool{"SKILL.md": false, "scripts/check.sh": true}, readLock(t, dir).Executable)
}

func TestInstall_localModeEditIsConflict(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "go-review")
	_, err := newInstaller().Install([]install.Request{oneTarget(dir)})
	require.NoError(t, err)
	require.NoError(t, os.Chmod(filepath.Join(dir, "SKILL.md"), 0o755)) //nolint:gosec // the test needs the executable bit set

	results, err := newInstaller().Install([]install.Request{oneTarget(dir)})

	require.NoError(t, err)
	assert.Equal(t, install.StateConflict, results[0].State, "a changed executable bit is your edit")
	assert.Equal(t, []string{filepath.Join(dir, "SKILL.md")}, results[0].Conflicts)
}

func TestInstall_legacyLock(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		want   []byte
		force  bool
		expect install.State
	}{
		"nothing to do stays unchanged":   {want: claudeSkill, expect: install.StateUnchanged},
		"upstream change needs attention": {want: []byte("# v2\n"), expect: install.StateConflict},
		"force updates after a backup":    {want: []byte("# v2\n"), force: true, expect: install.StateUpdate},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := filepath.Join(t.TempDir(), "go-review")
			legacyLock(t, dir, map[string][]byte{"SKILL.md": claudeSkill})
			lockBefore := read(t, filepath.Join(dir, install.LockFile))
			installer := newInstaller()
			installer.Backups = filepath.Join(t.TempDir(), "backups")
			installer.Force = tt.force
			req := install.Request{Name: "go-review", Targets: []install.Desired{
				{Dir: dir, Files: map[string]skill.File{"SKILL.md": {Data: tt.want}}},
			}}

			results, err := installer.Install([]install.Request{req})

			require.NoError(t, err)
			assert.Equal(t, tt.expect, results[0].State)
			switch tt.expect {
			case install.StateUnchanged:
				assert.Equal(t, lockBefore, read(t, filepath.Join(dir, install.LockFile)), "a legacy lock is left alone when nothing changes")
			case install.StateConflict:
				assert.Equal(t, claudeSkill, read(t, filepath.Join(dir, "SKILL.md")), "an unknown base mode is not overwritten")
				assert.NoDirExists(t, installer.Backups)
			case install.StateUpdate:
				assert.Equal(t, tt.want, read(t, filepath.Join(dir, "SKILL.md")))
				assert.Len(t, results[0].Backups, 1)
				assert.Equal(t, map[string]bool{"SKILL.md": false}, readLock(t, dir).Executable, "the new lock records complete modes")
			}
		})
	}
}

func TestInstall_modesIgnoredWhenDisabled(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "go-review")
	_, err := newInstaller().Install([]install.Request{oneTarget(dir)})
	require.NoError(t, err)
	require.NoError(t, os.Chmod(filepath.Join(dir, "SKILL.md"), 0o755)) //nolint:gosec // the test needs the executable bit set
	installer := newInstaller()
	installer.Modes = false

	results, err := installer.Install([]install.Request{oneTarget(dir)})

	require.NoError(t, err)
	assert.Equal(t, install.StateUnchanged, results[0].State, "without POSIX enforcement only content counts")
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

func TestInstall_preservesEditsToNestedLockFile(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "go-review")
	req := oneTarget(dir)
	req.Targets[0].Files["assets/.skill-lock.json"] = skill.File{Data: []byte("{\"example\": true}\n")}
	_, err := newInstaller().Install([]install.Request{req})
	require.NoError(t, err)
	file := filepath.Join(dir, "assets", ".skill-lock.json")
	local := []byte("{\"example\": \"edited locally\"}\n")
	write(t, file, local)
	lockBefore := read(t, filepath.Join(dir, install.LockFile))

	results, err := newInstaller().Install([]install.Request{req})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateConflict, results[0].State)
	assert.Equal(t, local, read(t, file), "a nested lock-named asset is user content, so edits must not be overwritten")
	assert.Equal(t, lockBefore, read(t, filepath.Join(dir, install.LockFile)))
}

func TestRequests_rejectsSourceLockFile(t *testing.T) {
	t.Parallel()
	store := fstest.MapFS{
		"skills/go-review/SKILL.md":         {Data: []byte("---\nname: go-review\ndescription: Reviews Go code.\nstatus: published\n---\n# Go review\n")},
		"skills/go-review/.skill-lock.json": {Data: []byte("{\"example\": true}\n")},
	}
	catalog, err := skill.Catalog(store)
	require.NoError(t, err)
	dests := []install.Destination{{Dir: t.TempDir(), Variant: install.VariantAgents}}

	_, err = install.Requests(store, catalog, dests)

	require.Error(t, err, "source content must not be silently replaced by the installer's lock")
	assert.ErrorContains(t, err, "go-review")
	assert.ErrorContains(t, err, install.LockFile)
}

func newInstaller() install.Installer {
	return install.Installer{
		Source: source,
		Commit: commit,
		Modes:  true,
		Now:    func() time.Time { return installedAt },
	}
}

func withScript(dir string, script []byte, mode os.FileMode) install.Request {
	return install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: dir, Files: map[string]skill.File{
			"SKILL.md":         {Data: claudeSkill},
			"scripts/check.sh": {Data: script, Mode: mode},
		}},
	}}
}

func legacyLock(t *testing.T, dir string, files map[string][]byte) {
	t.Helper()
	hashes := map[string]string{}
	for p, data := range files {
		write(t, filepath.Join(dir, p), data)
		hashes[p] = sha(data)
	}
	lock := map[string]any{
		"name":      filepath.Base(dir),
		"source":    source,
		"commit":    commit,
		"installed": installedAt.Add(-24 * time.Hour),
		"files":     hashes,
	}
	data, err := json.Marshal(lock)
	require.NoError(t, err)
	write(t, filepath.Join(dir, install.LockFile), data)
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
