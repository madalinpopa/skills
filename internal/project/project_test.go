package project_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/project"
)

func TestDetect_insideRepository(t *testing.T) {
	t.Parallel()
	root := initRepo(t)
	nested := filepath.Join(root, "internal", "deep")
	require.NoError(t, os.MkdirAll(nested, 0o755))

	tests := map[string]struct {
		start string
	}{
		"from root":   {start: root},
		"from nested": {start: nested},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			scope, err := project.Detect(tt.start)

			require.NoError(t, err)
			assert.Equal(t, root, scope.Root)
			assert.Empty(t, scope.Warning)
		})
	}
}

func TestDetect_outsideRepository(t *testing.T) {
	t.Parallel()
	dir := tempDir(t)

	scope, err := project.Detect(dir)

	require.NoError(t, err)
	assert.Equal(t, dir, scope.Root, "falls back to the starting directory")
	assert.Contains(t, scope.Warning, "not inside a Git repository")
	assert.Contains(t, scope.Warning, dir)
}

func TestDetect_missingDirectory(t *testing.T) {
	t.Parallel()

	_, err := project.Detect(filepath.Join(t.TempDir(), "gone"))

	assert.Error(t, err)
}

func initRepo(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	dir := tempDir(t)
	out, err := exec.Command("git", "init", "-q", dir).CombinedOutput()
	require.NoError(t, err, string(out))
	return dir
}

func tempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	return dir
}
