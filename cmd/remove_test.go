package cmd_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/cmd"
	"github.com/madalinpopa/skills/internal/gittest"
)

func TestRemove_global(t *testing.T) {
	home := configureStore(t, gittest.Init(t))
	installSkill(t, home, "go-review")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"remove", "--global", "go-review"}, strings.NewReader(""), &out, &errOut)

	require.NoError(t, err, errOut.String())
	assert.Contains(t, out.String(), "go-review")
	assert.Contains(t, out.String(), "removed")
	assert.NoDirExists(t, filepath.Join(home, ".claude", "skills", "go-review"))
	assert.NoDirExists(t, filepath.Join(home, ".agents", "skills", "go-review"))
	backups := filepath.Join(home, ".config", "skills", "backups")
	assert.DirExists(t, backups, "backups live under the config root")
	assert.Contains(t, out.String(), backups, "the backup path is printed")
	assert.NoDirExists(t, storeDir(home), "removal does not need the store")
}
