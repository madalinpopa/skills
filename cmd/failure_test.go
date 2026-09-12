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

func TestRemove_laterFailureStillPrintsResults(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("directory permissions are not enforced for root")
	}
	home := configureStore(t, gittest.Init(t))
	installSkill(t, home, "go-review")
	installSkill(t, home, "sql-review")
	locked := filepath.Join(home, ".claude", "skills", "sql-review", "scripts")
	require.NoError(t, os.MkdirAll(locked, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(locked, "keep.sh"), []byte("#!/bin/sh\n"), 0o600))
	require.NoError(t, os.Chmod(locked, 0o555))       //nolint:gosec // a read-only directory keeps its execute bit
	t.Cleanup(func() { _ = os.Chmod(locked, 0o755) }) //nolint:gosec // restore so the temp dir can be cleaned
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"remove", "--global", "--force", "go-review", "sql-review"}, strings.NewReader(""), &out, &errOut)

	require.Error(t, err)
	assert.NotErrorIs(t, err, cmd.ErrUsage)
	assert.Contains(t, out.String(), "- go-review")
	assert.Contains(t, out.String(), "! sql-review")
	assert.Equal(t, 4, strings.Count(out.String(), "backed up to"), "every completed backup path is printed, two targets per skill")
	assert.Contains(t, errOut.String(), "sql-review", "the diagnostic names the failed skill")
	assert.NoDirExists(t, filepath.Join(home, ".claude", "skills", "go-review"))
}
