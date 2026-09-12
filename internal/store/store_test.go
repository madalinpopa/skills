package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/gittest"
	"github.com/madalinpopa/skills/internal/store"
)

func TestInit_clonesConfiguredBranch(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	gittest.Run(t, source, "checkout", "-q", "-b", "feature")
	want := gittest.Commit(t, source, "on feature")
	s := newStore(t, source, "feature")

	err := s.Init(t.Context())

	require.NoError(t, err)
	got, err := s.Commit(t.Context())
	require.NoError(t, err)
	assert.Equal(t, want, got, "the store is at the tip of the configured branch")
	assert.Equal(t, "feature", gittest.Run(t, s.Dir, "rev-parse", "--abbrev-ref", "HEAD"))
	assert.Equal(t, "origin/feature", gittest.Run(t, s.Dir, "branch", "-r", "--format=%(refname:short)"),
		"only the configured branch is fetched")
}

func TestInit_isIdempotent(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	s := newStore(t, source, "main")
	require.NoError(t, s.Init(t.Context()))
	before, err := s.Commit(t.Context())
	require.NoError(t, err)

	err = s.Init(t.Context())

	require.NoError(t, err)
	after, err := s.Commit(t.Context())
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestInit_failed(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup func(t *testing.T, s store.Store)
	}{
		"plain directory in the way": {
			setup: func(t *testing.T, s store.Store) {
				t.Helper()
				require.NoError(t, os.MkdirAll(s.Dir, 0o750))
				require.NoError(t, os.WriteFile(filepath.Join(s.Dir, "keep.txt"), []byte("mine"), 0o600))
			},
		},
		"clone of another repository": {
			setup: func(t *testing.T, s store.Store) {
				t.Helper()
				other := gittest.Init(t)
				gittest.Run(t, filepath.Dir(s.Dir), "clone", "-q", other, s.Dir)
				require.NoError(t, os.WriteFile(filepath.Join(s.Dir, "keep.txt"), []byte("mine"), 0o600))
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := newStore(t, gittest.Init(t), "main")
			tt.setup(t, s)

			err := s.Init(t.Context())

			require.Error(t, err)
			assert.ErrorContains(t, err, s.Dir, "the error names the path that is in the way")
			assert.FileExists(t, filepath.Join(s.Dir, "keep.txt"), "nothing is overwritten")
		})
	}
}

func TestInit_missingRepository(t *testing.T) {
	t.Parallel()
	s := newStore(t, filepath.Join(t.TempDir(), "nowhere"), "main")

	err := s.Init(t.Context())

	require.Error(t, err)
	assert.ErrorContains(t, err, s.Repo, "git's own message is kept")
	assert.NoDirExists(t, s.Dir)
}

func TestSync_fastForwards(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	s := newStore(t, source, "main")
	require.NoError(t, s.Init(t.Context()))
	old, err := s.Commit(t.Context())
	require.NoError(t, err)
	next := gittest.Commit(t, source, "second")

	result, err := s.Sync(t.Context())

	require.NoError(t, err)
	assert.Equal(t, store.Result{Old: old, New: next}, result)
	got, err := s.Commit(t.Context())
	require.NoError(t, err)
	assert.Equal(t, next, got)
}

func TestSync_upToDate(t *testing.T) {
	t.Parallel()
	s := newStore(t, gittest.Init(t), "main")
	require.NoError(t, s.Init(t.Context()))
	head, err := s.Commit(t.Context())
	require.NoError(t, err)

	result, err := s.Sync(t.Context())

	require.NoError(t, err)
	assert.Equal(t, store.Result{Old: head, New: head}, result)
}

func TestSync_refusesDivergence(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	s := newStore(t, source, "main")
	require.NoError(t, s.Init(t.Context()))
	local := gittest.Commit(t, s.Dir, "local work")
	gittest.Commit(t, source, "upstream work")

	_, err := s.Sync(t.Context())

	require.Error(t, err)
	assert.ErrorContains(t, err, s.Dir, "the error says where to look")
	got, commitErr := s.Commit(t.Context())
	require.NoError(t, commitErr)
	assert.Equal(t, local, got, "no merge, no reset")
}

func TestSync_keepsDirtyState(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	s := newStore(t, source, "main")
	require.NoError(t, s.Init(t.Context()))
	readme := filepath.Join(s.Dir, "README.md")
	require.NoError(t, os.WriteFile(readme, []byte("edited\n"), 0o600))
	gittest.Commit(t, source, "upstream work")

	_, err := s.Sync(t.Context())

	require.Error(t, err)
	assert.ErrorContains(t, err, s.Dir, "the error says where to look")
	content, readErr := os.ReadFile(readme)
	require.NoError(t, readErr)
	assert.Equal(t, "edited\n", string(content), "local edits are never reset")
}

func newStore(t *testing.T, repo, branch string) store.Store {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	return store.Store{
		Dir:    filepath.Join(dir, "store"),
		Repo:   repo,
		Branch: branch,
	}
}
