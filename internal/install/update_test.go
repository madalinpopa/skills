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

func TestFind(t *testing.T) {
	t.Parallel()
	all := []install.Installation{{Name: "go-review"}, {Name: "sql-review"}}
	tests := map[string]struct {
		names   []string
		want    []string
		wantErr string
	}{
		"no names means every installation": {want: []string{"go-review", "sql-review"}},
		"exact names only":                  {names: []string{"sql-review"}, want: []string{"sql-review"}},
		"unknown name":                      {names: []string{"go-review", "nope"}, wantErr: "not installed: nope"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			found, err := install.Find(all, tt.names)

			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			var names []string
			for _, f := range found {
				names = append(names, f.Name)
			}
			assert.Equal(t, tt.want, names)
		})
	}
}

func TestUpdateRequests(t *testing.T) {
	t.Parallel()
	store := t.TempDir()
	write(t, filepath.Join(store, "skills", "go-review", "SKILL.md"),
		[]byte("---\nname: go-review\ndescription: d\nstatus: published\n---\n# Go review\n"))
	catalog, err := skill.Catalog(os.DirFS(store))
	require.NoError(t, err)
	agents := install.Destination{Dir: filepath.Join("p", ".agents", "skills"), Variant: install.VariantAgents}
	claude := install.Destination{Dir: filepath.Join("p", ".claude", "skills"), Variant: install.VariantClaude}
	installations := []install.Installation{
		{Name: "go-review", Targets: []install.Destination{agents, claude}},
		{Name: "sql-review", Targets: []install.Destination{claude}},
	}

	reqs, err := install.UpdateRequests(store, catalog, installations)

	require.NoError(t, err)
	require.Len(t, reqs, 2)
	assert.Equal(t, "go-review", reqs[0].Name)
	assert.False(t, reqs[0].Unavailable)
	require.Len(t, reqs[0].Targets, 2, "only the installed targets are refreshed")
	assert.Equal(t, filepath.Join(agents.Dir, "go-review"), reqs[0].Targets[0].Dir)
	assert.Equal(t, filepath.Join(claude.Dir, "go-review"), reqs[0].Targets[1].Dir)
	assert.Contains(t, reqs[0].Targets[0].Files, "SKILL.md")
	assert.Equal(t, install.Request{Name: "sql-review", Unavailable: true}, reqs[1], "missing or draft store skills are unavailable")
}

func TestInstall_updatesManagedSkills(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	clean := filepath.Join(root, "clean")
	edited := filepath.Join(root, "edited")
	same := filepath.Join(root, "same")
	installed(t, clean, map[string][]byte{"SKILL.md": []byte("# v1\n")})
	installed(t, edited, map[string][]byte{"SKILL.md": []byte("# v1\n")})
	write(t, filepath.Join(edited, "SKILL.md"), []byte("# mine\n"))
	installed(t, same, map[string][]byte{"SKILL.md": []byte("# v1\n")})
	v2 := map[string]skill.File{"SKILL.md": {Data: []byte("# v2\n")}}
	reqs := []install.Request{
		{Name: "clean", Targets: []install.Desired{{Dir: clean, Files: v2}}},
		{Name: "edited", Targets: []install.Desired{{Dir: edited, Files: v2}}},
		{Name: "same", Targets: []install.Desired{{Dir: same, Files: map[string]skill.File{"SKILL.md": {Data: []byte("# v1\n")}}}}},
	}

	results, err := newInstaller().Install(reqs)

	require.NoError(t, err)
	require.Len(t, results, 3)
	assert.Equal(t, install.StateUpdate, results[0].State)
	assert.Equal(t, []byte("# v2\n"), read(t, filepath.Join(clean, "SKILL.md")))
	assert.Equal(t, commit, readLock(t, clean).Commit, "an updated skill records the new commit")
	assert.Equal(t, install.StateConflict, results[1].State)
	assert.Equal(t, []byte("# mine\n"), read(t, filepath.Join(edited, "SKILL.md")), "your edit stays in place")
	assert.Equal(t, oldCommit, readLock(t, edited).Commit)
	assert.Equal(t, install.StateUnchanged, results[2].State)
	assert.Equal(t, oldCommit, readLock(t, same).Commit, "an unchanged skill keeps its lock")
}

func TestInstall_unavailableIsRetained(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "go-review")
	installed(t, dir, map[string][]byte{"SKILL.md": claudeSkill})

	results, err := newInstaller().Install([]install.Request{{Name: "go-review", Unavailable: true}})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateUnavailable, results[0].State)
	assert.Equal(t, claudeSkill, read(t, filepath.Join(dir, "SKILL.md")))
	assert.Equal(t, oldCommit, readLock(t, dir).Commit)
}

func TestInstall_foreignSourceIsRetained(t *testing.T) {
	t.Parallel()
	dir := filepath.Join(t.TempDir(), "go-review")
	other := "https://example.com/other-store"
	installedFrom(t, dir, other, map[string][]byte{"SKILL.md": []byte("# theirs\n")})

	results, err := newInstaller().Install([]install.Request{oneTarget(dir)})

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, install.StateForeign, results[0].State)
	assert.Equal(t, other, results[0].Source)
	assert.Equal(t, []byte("# theirs\n"), read(t, filepath.Join(dir, "SKILL.md")))
	assert.Equal(t, other, readLock(t, dir).Source, "the lock is never adopted into the configured store")
}
