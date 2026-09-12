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
	installed := filepath.Join(home, ".claude", "skills", "go-review", "SKILL.md")
	data, err := os.ReadFile(installed)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(installed, append(data, "my note\n"...), 0o600))
	addSkill(t, source, "go-review", "published", "Reviews Go code carefully.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "improve go-review")
	run(t, "sync")
	var out, errOut bytes.Buffer

	err = cmd.Execute(t.Context(), []string{"diff", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "+my note")
	assert.Contains(t, out.String(), filepath.Join(".claude", "skills", "go-review", "SKILL.md"))
	assert.NotContains(t, out.String(), "carefully", "the base is the recorded commit, not the current store")
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

			err := cmd.Execute(t.Context(), []string{"diff", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

			require.Error(t, err)
			assert.NotErrorIs(t, err, cmd.ErrUsage)
			assert.Contains(t, errOut.String(), tt.want)
			assert.Empty(t, out.String(), "no diff against a substitute base")
		})
	}
}
