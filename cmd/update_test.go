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

func TestUpdate_appliesSyncedStore(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	run(t, "install", "--global", "go-review")
	addSkill(t, source, "go-review", "published", "Reviews Go code carefully.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "improve go-review")
	run(t, "sync")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), []string{"update", "--global"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "go-review")
	assert.Contains(t, out.String(), "updated")
	for _, agent := range []string{".claude", ".agents"} {
		data, err := os.ReadFile(filepath.Join(home, agent, "skills", "go-review", "SKILL.md"))
		require.NoError(t, err)
		assert.Contains(t, string(data), "carefully")
	}
}

func run(t *testing.T, args ...string) {
	t.Helper()
	var out, errOut bytes.Buffer
	require.NoError(t, cmd.Execute(t.Context(), args, strings.NewReader(""), &out, &errOut), errOut.String())
}
