package install_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/config"
	"github.com/madalinpopa/skills/internal/install"
	"github.com/madalinpopa/skills/internal/skill"
)

func TestInstall_symlinkInTargetIsUnsupported(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claudeDir := filepath.Join(root, ".claude", "skills", "go-review")
	agentsDir := filepath.Join(root, ".agents", "skills", "go-review")
	installed(t, claudeDir, map[string][]byte{"SKILL.md": []byte("# v1\n")})
	installed(t, agentsDir, map[string][]byte{"SKILL.md": []byte("# v1\n")})
	agentsLock := read(t, filepath.Join(agentsDir, install.LockFile))
	outside := t.TempDir()
	write(t, filepath.Join(outside, "keep.md"), []byte("external\n"))
	require.NoError(t, os.Symlink(outside, filepath.Join(claudeDir, "references")))
	v2 := map[string]skill.File{
		"SKILL.md":          {Data: []byte("# v2\n")},
		"references/api.md": {Data: []byte("# api\n")},
	}
	installer := newInstaller()
	installer.Backups = filepath.Join(t.TempDir(), "backups")
	installer.Force = true
	reqs := []install.Request{
		{Name: "go-review", Targets: []install.Desired{{Dir: claudeDir, Files: v2}, {Dir: agentsDir, Files: v2}}},
		{Name: "sql-review", Targets: []install.Desired{{Dir: filepath.Join(root, "sql-review"), Files: map[string]skill.File{"SKILL.md": {Data: []byte("# sql\n")}}}}},
	}

	results, err := installer.Install(reqs)

	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, install.StateUnsupported, results[0].State)
	assert.Equal(t, []install.Issue{{Path: filepath.Join(claudeDir, "references"), Reason: "is a symlink"}}, results[0].Issues,
		"the actual reason is reported, even with force")
	assert.NoFileExists(t, filepath.Join(outside, "api.md"), "nothing is written through the link")
	assert.Equal(t, []byte("# v1\n"), read(t, filepath.Join(agentsDir, "SKILL.md")), "the healthy sibling target is left alone")
	assert.Equal(t, agentsLock, read(t, filepath.Join(agentsDir, install.LockFile)))
	assert.NoDirExists(t, installer.Backups, "no backup for a skill that is not changed")
	entries, err := os.ReadDir(claudeDir)
	require.NoError(t, err)
	assert.Len(t, entries, 3, "no staged leftovers")
	assert.Equal(t, install.StateAdd, results[1].State, "healthy skills still proceed")
}

func TestInstall_directoryAtFilePathIsUnsupported(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	claudeDir := filepath.Join(root, ".claude", "skills", "go-review")
	agentsDir := filepath.Join(root, ".agents", "skills", "go-review")
	write(t, filepath.Join(claudeDir, "SKILL.md", "inner.txt"), []byte("a directory where the file should be\n"))
	req := install.Request{Name: "go-review", Targets: []install.Desired{
		{Dir: claudeDir, Files: map[string]skill.File{"SKILL.md": {Data: claudeSkill}, "notes.md": {Data: []byte("n\n")}}},
		{Dir: agentsDir, Files: map[string]skill.File{"SKILL.md": {Data: agentsSkill}}},
	}}

	results, err := newInstaller().Install([]install.Request{req})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateUnsupported, results[0].State)
	assert.Equal(t, []install.Issue{{Path: filepath.Join(claudeDir, "SKILL.md"), Reason: "is a directory where a file is expected"}}, results[0].Issues)
	assert.NoFileExists(t, filepath.Join(claudeDir, "notes.md"), "the collision is caught before any write")
	assert.NoFileExists(t, filepath.Join(claudeDir, install.LockFile))
	assert.NoDirExists(t, agentsDir, "the clean target is skipped too")
}

func TestInstall_skillDirectorySymlinkIsUnsupported(t *testing.T) {
	t.Parallel()
	outside := filepath.Join(t.TempDir(), "go-review")
	write(t, filepath.Join(outside, "SKILL.md"), claudeSkill)
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	require.NoError(t, os.MkdirAll(filepath.Dir(dir), 0o750))
	require.NoError(t, os.Symlink(outside, dir))

	results, err := newInstaller().Install([]install.Request{oneTarget(dir)})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateUnsupported, results[0].State)
	assert.Equal(t, []install.Issue{{Path: dir, Reason: "is a symlink"}}, results[0].Issues)
	assert.NoFileExists(t, filepath.Join(outside, install.LockFile), "identical content behind a link is not adopted")
}

func TestRemove_symlinkNeedsManualRepair(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), ".claude", "skills", "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": claudeSkill})
	outside := filepath.Join(t.TempDir(), "notes.md")
	write(t, outside, []byte("external\n"))
	require.NoError(t, os.Symlink(outside, filepath.Join(dir, "notes.md")))
	remover := newInstaller()
	remover.Backups = filepath.Join(t.TempDir(), "backups")
	remover.Force = true

	results, err := remover.Remove([]install.Installation{
		{Name: "go-review", Targets: []install.Destination{{Dir: filepath.Dir(dir), Variant: install.VariantClaude}}},
	})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateUnsupported, results[0].State)
	assert.Equal(t, []install.Issue{{Path: filepath.Join(dir, "notes.md"), Reason: "is a symlink"}}, results[0].Issues)
	assert.FileExists(t, filepath.Join(dir, "SKILL.md"), "the skill stays installed")
	assert.FileExists(t, outside)
	assert.NoDirExists(t, remover.Backups, "no backup reads through the link")
}

func TestScan_findsSymlinkedSkillDirectory(t *testing.T) {
	t.Parallel()
	outside := filepath.Join(t.TempDir(), "go-review")
	installed(t, outside, map[string][]byte{"SKILL.md": claudeSkill})
	root := t.TempDir()
	claude := filepath.Join(root, ".claude", "skills")
	require.NoError(t, os.MkdirAll(claude, 0o750))
	require.NoError(t, os.Symlink(outside, filepath.Join(claude, "go-review")))
	dests, err := install.Targets(config.Default(), root, false, []string{"claude"})
	require.NoError(t, err)

	found, err := install.Scan(dests)

	require.NoError(t, err)
	require.Len(t, found, 1, "a linked managed skill stays visible, so it cannot become a fresh install")
	assert.Equal(t, "go-review", found[0].Name)
}

func TestCheck(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "go-review")
	write(t, filepath.Join(dir, "SKILL.md"), claudeSkill)
	write(t, filepath.Join(dir, "scripts", "check.sh"), []byte("#!/bin/sh\n"))
	require.NoError(t, os.Symlink("SKILL.md", filepath.Join(dir, "scripts", "link.md")))

	issues, err := install.Check(dir)

	require.NoError(t, err)
	assert.Equal(t, []install.Issue{{Path: filepath.Join(dir, "scripts", "link.md"), Reason: "is a symlink"}}, issues)
}
