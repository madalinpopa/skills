package gittest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Init(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	require.NoError(t, err)
	Run(t, dir, "init", "-q", "-b", "main")
	require.NoError(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("# store\n"), 0o600))
	Run(t, dir, "add", "README.md")
	Commit(t, dir, "initial")
	return dir
}

func Commit(t *testing.T, dir, message string) string {
	t.Helper()
	Run(t, dir, "commit", "-q", "--allow-empty", "-m", message)
	return Run(t, dir, "rev-parse", "HEAD")
}

func Run(t *testing.T, dir string, args ...string) string {
	t.Helper()
	args = append([]string{
		"-c", "user.name=test",
		"-c", "user.email=test@example.com",
		"-c", "commit.gpgsign=false",
	}, args...)
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, string(out))
	return strings.TrimSpace(string(out))
}
