package store_test

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/gittest"
	"github.com/madalinpopa/skills/internal/skill"
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

func TestPreview_reportsRemoteHeadWithoutChangingStore(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	s := newStore(t, source, "main")
	require.NoError(t, s.Init(t.Context()))
	local, err := s.Commit(t.Context())
	require.NoError(t, err)
	remote := gittest.Commit(t, source, "upstream work")
	before := snapshot(t, s.Dir)

	result, err := s.Preview(t.Context())

	require.NoError(t, err)
	assert.Equal(t, store.Result{Old: local, New: remote}, result)
	assert.Equal(t, before, snapshot(t, s.Dir), "the working tree, index, refs and objects are untouched")
}

func TestPreview_upToDate(t *testing.T) {
	t.Parallel()
	s := newStore(t, gittest.Init(t), "main")
	require.NoError(t, s.Init(t.Context()))
	head, err := s.Commit(t.Context())
	require.NoError(t, err)

	result, err := s.Preview(t.Context())

	require.NoError(t, err)
	assert.Equal(t, store.Result{Old: head, New: head}, result)
}

func TestPreview_failed(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		setup func(t *testing.T, source string, s store.Store)
		want  string
	}{
		"dirty store": {
			setup: func(t *testing.T, _ string, s store.Store) {
				t.Helper()
				require.NoError(t, os.WriteFile(filepath.Join(s.Dir, "README.md"), []byte("edited\n"), 0o600))
			},
			want: "local changes",
		},
		"wrong branch": {
			setup: func(t *testing.T, _ string, s store.Store) {
				t.Helper()
				gittest.Run(t, s.Dir, "checkout", "-q", "-b", "other")
			},
			want: "branch other",
		},
		"unreachable remote": {
			setup: func(t *testing.T, source string, _ store.Store) {
				t.Helper()
				require.NoError(t, os.RemoveAll(source))
			},
			want: "git ls-remote",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			source := gittest.Init(t)
			s := newStore(t, source, "main")
			require.NoError(t, s.Init(t.Context()))
			tt.setup(t, source, s)
			before := snapshot(t, s.Dir)

			_, err := s.Preview(t.Context())

			require.Error(t, err)
			assert.ErrorContains(t, err, tt.want)
			assert.Equal(t, before, snapshot(t, s.Dir), "a failed preview writes nothing")
		})
	}
}

func snapshot(t *testing.T, dir string) map[string]string {
	t.Helper()
	root, err := os.OpenRoot(dir)
	require.NoError(t, err)
	fsys := root.FS()
	files := map[string]string{}
	walkErr := fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		files[path] = hex.EncodeToString(sum[:])
		return nil
	})
	require.NoError(t, root.Close())
	require.NoError(t, walkErr)
	return files
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

func TestTree_readsRecordedCommit(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	dir := filepath.Join(source, "skills", "go-review")
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "scripts"), 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# v1\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "scripts", "check.sh"), []byte("#!/bin/sh\n"), 0o600))
	gittest.Run(t, source, "add", ".")
	gittest.Run(t, source, "update-index", "--chmod=+x", "skills/go-review/scripts/check.sh")
	old := gittest.Commit(t, source, "v1")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# v2\n"), 0o600))
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "v2")
	s := newStore(t, source, "main")
	require.NoError(t, s.Init(t.Context()))

	tree, err := s.Tree(t.Context(), old)

	require.NoError(t, err)
	data, err := fs.ReadFile(tree, "skills/go-review/SKILL.md")
	require.NoError(t, err)
	assert.Equal(t, "# v1\n", string(data), "content comes from the recorded commit, not the tip")
	info, err := fs.Stat(tree, "skills/go-review/scripts/check.sh")
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&0o100, "executable bits survive")
}

func TestTree_excludesUnpublishedContent(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	dir := filepath.Join(source, "skills", "go-review")
	require.NoError(t, os.MkdirAll(dir, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# v1\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(source, ".gitignore"), []byte("*.log\n"), 0o600))
	gittest.Run(t, source, "add", ".")
	head := gittest.Commit(t, source, "v1")
	s := newStore(t, source, "main")
	require.NoError(t, s.Init(t.Context()))
	local := filepath.Join(s.Dir, "skills", "go-review")
	require.NoError(t, os.WriteFile(filepath.Join(local, "SKILL.md"), []byte("# edited\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(local, "notes.log"), []byte("ignored\n"), 0o600))
	untracked := filepath.Join(s.Dir, "skills", "untracked")
	require.NoError(t, os.MkdirAll(untracked, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(untracked, "SKILL.md"), []byte("# untracked\n"), 0o600))

	tree, err := s.Tree(t.Context(), head)

	require.NoError(t, err)
	data, err := fs.ReadFile(tree, "skills/go-review/SKILL.md")
	require.NoError(t, err)
	assert.Equal(t, "# v1\n", string(data), "uncommitted edits are not published")
	_, err = fs.Stat(tree, "skills/go-review/notes.log")
	assert.ErrorIs(t, err, fs.ErrNotExist, "ignored files are not published")
	_, err = fs.Stat(tree, "skills/untracked")
	assert.ErrorIs(t, err, fs.ErrNotExist, "untracked skills are not published")
}

func TestTree_ignoresCheckoutConversion(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	dir := filepath.Join(source, "skills", "go-review")
	require.NoError(t, os.MkdirAll(dir, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(source, ".gitattributes"), []byte("*.md text eol=crlf\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# v1\n"), 0o600))
	gittest.Run(t, source, "add", ".")
	head := gittest.Commit(t, source, "v1")
	s := newStore(t, source, "main")
	require.NoError(t, s.Init(t.Context()))
	checkedOut, err := os.ReadFile(filepath.Join(s.Dir, "skills", "go-review", "SKILL.md"))
	require.NoError(t, err)
	require.Equal(t, "# v1\r\n", string(checkedOut), "the working tree is converted on checkout")

	tree, err := s.Tree(t.Context(), head)

	require.NoError(t, err)
	data, err := fs.ReadFile(tree, "skills/go-review/SKILL.md")
	require.NoError(t, err)
	assert.Equal(t, "# v1\n", string(data), "published bytes are the committed bytes")
}

func TestTree_rejectsSymlinks(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	dir := filepath.Join(source, "skills", "go-review")
	require.NoError(t, os.MkdirAll(dir, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("---\nname: go-review\ndescription: d\nstatus: published\n---\n"), 0o600))
	require.NoError(t, os.Symlink("SKILL.md", filepath.Join(dir, "link.md")))
	gittest.Run(t, source, "add", ".")
	head := gittest.Commit(t, source, "v1")
	s := newStore(t, source, "main")
	require.NoError(t, s.Init(t.Context()))
	tree, err := s.Tree(t.Context(), head)
	require.NoError(t, err)
	sub, err := fs.Sub(tree, "skills/go-review")
	require.NoError(t, err)

	_, err = skill.Render(sub, "claude")

	require.Error(t, err)
	assert.ErrorContains(t, err, "link.md", "a source symlink is rejected, not silently dropped")
}

func TestTree_readsOnlySkillsDirectory(t *testing.T) {
	t.Parallel()
	source := gittest.Init(t)
	for name, data := range map[string]string{
		"skills/go-review/SKILL.md":            "# go review\n",
		"skills/go-review/references/style.md": "# style\n",
		"skills-old/legacy/SKILL.md":           "# legacy\n",
		"docs/guide.md":                        "# guide\n",
	} {
		file := filepath.Join(source, filepath.FromSlash(name))
		require.NoError(t, os.MkdirAll(filepath.Dir(file), 0o750))
		require.NoError(t, os.WriteFile(file, []byte(data), 0o600))
	}
	gittest.Run(t, source, "add", ".")
	head := gittest.Commit(t, source, "skills and other content")
	s := newStore(t, source, "main")
	require.NoError(t, s.Init(t.Context()))

	tree, err := s.Tree(t.Context(), head)

	require.NoError(t, err)
	assert.Equal(t, []string{
		"skills/go-review/SKILL.md",
		"skills/go-review/references/style.md",
	}, treePaths(t, tree), "root files, other directories and similarly named siblings are not read")
}

func TestTree_withoutSkillsDirectory(t *testing.T) {
	t.Parallel()
	s := newStore(t, gittest.Init(t), "main")
	require.NoError(t, s.Init(t.Context()))
	head, err := s.Commit(t.Context())
	require.NoError(t, err)

	tree, err := s.Tree(t.Context(), head)

	require.NoError(t, err, "a valid commit without skills is an empty store, not an error")
	assert.Empty(t, treePaths(t, tree))
	_, err = fs.ReadDir(tree, "skills")
	assert.ErrorIs(t, err, fs.ErrNotExist)
}

func treePaths(t *testing.T, fsys fs.FS) []string {
	t.Helper()
	var paths []string
	err := fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
		if err == nil && !entry.IsDir() {
			paths = append(paths, path)
		}
		return err
	})
	require.NoError(t, err)
	return paths
}

func TestTree_unknownCommit(t *testing.T) {
	t.Parallel()
	s := newStore(t, gittest.Init(t), "main")
	require.NoError(t, s.Init(t.Context()))

	_, err := s.Tree(t.Context(), "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef")

	require.Error(t, err)
	assert.ErrorContains(t, err, "deadbeef")
}
