package cmd_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/cmd"
	"github.com/madalinpopa/skills/internal/gittest"
)

func TestInit_clonesStore(t *testing.T) {
	source := gittest.Init(t)
	home := configureStore(t, source)
	head := gittest.Run(t, source, "rev-parse", "HEAD")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"init"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.DirExists(t, filepath.Join(storeDir(home), ".git"))
	assert.Contains(t, out.String(), head[:7])
}

func TestSync_fastForwards(t *testing.T) {
	source := gittest.Init(t)
	configureStore(t, source)
	old := gittest.Run(t, source, "rev-parse", "HEAD")
	require.NoError(t, cmd.Execute(t.Context(), "", []string{"init"}, strings.NewReader(""), &bytes.Buffer{}, &bytes.Buffer{}))
	next := gittest.Commit(t, source, "second")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"sync"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), fmt.Sprintf("pulled store  %s -> %s", old[:7], next[:7]))
}

func TestSync_initialisesMissingStore(t *testing.T) {
	source := gittest.Init(t)
	home := configureStore(t, source)
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"sync"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.DirExists(t, filepath.Join(storeDir(home), ".git"))
}

func TestSync_dryRunLeavesStoreUnchanged(t *testing.T) {
	source := gittest.Init(t)
	home := configureStore(t, source)
	old := gittest.Run(t, source, "rev-parse", "HEAD")
	run(t, "init")
	next := gittest.Commit(t, source, "second")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"sync", "--dry-run"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), fmt.Sprintf("would pull store  %s -> %s", old[:7], next[:7]))
	assert.Contains(t, out.String(), "fast-forward")
	assert.Equal(t, old, gittest.Run(t, storeDir(home), "rev-parse", "HEAD"), "dry-run does not move the store")
	assert.NoFileExists(t, filepath.Join(storeDir(home), ".git", "FETCH_HEAD"), "dry-run does not fetch")
}

func TestSync_dryRunRemoteFailure(t *testing.T) {
	source := gittest.Init(t)
	configureStore(t, source)
	run(t, "init")
	require.NoError(t, os.RemoveAll(source))
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"sync", "--dry-run"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.NotErrorIs(t, err, cmd.ErrUsage)
	assert.NotContains(t, out.String(), "up to date", "a failed remote query is never reported as up to date")
}

func TestSync_dryRunUninitialisedStore(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"sync", "--dry-run"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.NotErrorIs(t, err, cmd.ErrUsage)
	assert.Contains(t, errOut.String(), "skills init")
	assert.NoDirExists(t, filepath.Join(home, ".config"), "dry-run never creates the config or clones the store")
}

func TestInit_storeFailureIsRuntimeError(t *testing.T) {
	configureStore(t, filepath.Join(t.TempDir(), "nowhere"))
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"init"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.NotErrorIs(t, err, cmd.ErrUsage)
	assert.Contains(t, errOut.String(), "nowhere")
	assert.NotContains(t, errOut.String(), "Usage:", "runtime failures do not print usage")
}

func configureStore(t *testing.T, repo string) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	configDir := filepath.Join(home, ".config", "skills")
	require.NoError(t, os.MkdirAll(configDir, 0o750))
	toml := fmt.Sprintf("[store]\nrepo = %q\nbranch = \"main\"\n", repo)
	require.NoError(t, os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(toml), 0o600))
	return home
}

func storeDir(home string) string {
	return filepath.Join(home, ".config", "skills", "store")
}
