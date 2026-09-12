package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/cmd"
	"github.com/madalinpopa/skills/internal/gittest"
)

func TestDiff_showsOnlyYourEdits(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	run(t, "install", "--global", "go-review")
	appendTo(t, filepath.Join(home, ".claude", "skills", "go-review", "SKILL.md"), "my note\n")
	addSkill(t, source, "go-review", "published", "Reviews Go code carefully.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "improve go-review")
	run(t, "sync")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"diff", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "+my note")
	assert.Contains(t, out.String(), filepath.Join(".claude", "skills", "go-review", "SKILL.md"))
	assert.NotContains(t, out.String(), "carefully", "the base is the recorded commit, not the current store")
	assert.NotContains(t, out.String(), ".agents", "the untouched copy produces no diff")
}

func TestDiff_showsExecutableChange(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	script := filepath.Join(source, "skills", "go-review", "scripts", "check.sh")
	require.NoError(t, os.MkdirAll(filepath.Dir(script), 0o750))
	require.NoError(t, os.WriteFile(script, []byte("#!/bin/sh\n"), 0o600))
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	run(t, "install", "--global", "go-review")
	local := filepath.Join(home, ".claude", "skills", "go-review", "scripts", "check.sh")
	require.NoError(t, os.Chmod(local, 0o755)) //nolint:gosec // the test needs the executable bit set
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"diff", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), filepath.Join(".claude", "skills", "go-review", "scripts", "check.sh"))
	assert.Contains(t, out.String(), "old mode 100644")
	assert.Contains(t, out.String(), "new mode 100755")
	assert.NotContains(t, out.String(), "@@", "a mode-only change has no content hunk")
	assert.NotContains(t, out.String(), ".agents", "the untouched copy produces no diff")
}

func TestDiff_failed(t *testing.T) {
	source := gittest.Init(t)
	head := gittest.Run(t, source, "rev-parse", "HEAD")
	home := configureStore(t, source)
	tests := map[string]struct {
		source string
		commit string
		want   string
	}{
		"another source": {source: "https://example.com/other", commit: head, want: "https://example.com/other"},
		"missing commit": {source: source, commit: "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef", want: "deadbeef"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			dir := filepath.Join(home, ".claude", "skills", "go-review")
			require.NoError(t, os.RemoveAll(dir))
			writeInstalled(t, dir, "go-review", tt.source, tt.commit)
			var out, errOut bytes.Buffer

			err := cmd.Execute(t.Context(), "", []string{"diff", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

			require.Error(t, err)
			assert.NotErrorIs(t, err, cmd.ErrUsage)
			assert.Contains(t, errOut.String(), tt.want)
			assert.Empty(t, out.String(), "no diff against a substitute base")
		})
	}
}

func TestDiff_unsupportedEntryIsExplained(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	run(t, "install", "--global", "go-review")
	outside := filepath.Join(t.TempDir(), "notes.md")
	require.NoError(t, os.WriteFile(outside, []byte("external\n"), 0o600))
	link := filepath.Join(home, ".claude", "skills", "go-review", "notes.md")
	require.NoError(t, os.Symlink(outside, link))
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"diff", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.NotErrorIs(t, err, cmd.ErrUsage)
	assert.Contains(t, errOut.String(), link)
	assert.Contains(t, errOut.String(), "symlink")
	assert.Empty(t, out.String(), "the link is explained, not followed")
}

func appendTo(t *testing.T, path, text string) {
	t.Helper()
	file, err := os.OpenFile(filepath.Clean(path), os.O_APPEND|os.O_WRONLY, 0o600)
	require.NoError(t, err)
	_, err = file.WriteString(text)
	require.NoError(t, err)
	require.NoError(t, file.Close())
}
