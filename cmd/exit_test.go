package cmd_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/cmd"
	"github.com/madalinpopa/skills/internal/gittest"
)

func TestExitCode(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		err  error
		want int
	}{
		"success":           {err: nil, want: 0},
		"runtime failure":   {err: errors.New("git clone: not found"), want: 1},
		"usage error":       {err: cmd.ErrUsage, want: 2},
		"wrapped usage":     {err: fmt.Errorf("install: %w", cmd.ErrUsage), want: 2},
		"attention":         {err: cmd.ErrAttention, want: 3},
		"wrapped attention": {err: fmt.Errorf("update: %w", cmd.ErrAttention), want: 3},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, cmd.ExitCode(tt.err))
		})
	}
}

func TestUpdate_attentionAfterPartialSuccess(t *testing.T) {
	source := gittest.Init(t)
	addSkill(t, source, "go-review", "published", "Reviews Go code.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "add skills")
	home := configureStore(t, source)
	run(t, "install", "--global", "go-review")
	installSkill(t, home, "ghost")
	addSkill(t, source, "go-review", "published", "Reviews Go code carefully.", "[go]")
	gittest.Run(t, source, "add", ".")
	gittest.Commit(t, source, "improve go-review")
	run(t, "sync")
	var out, errOut bytes.Buffer

	err := cmd.Execute(t.Context(), "", []string{"update", "--global"}, strings.NewReader(""), &out, &errOut)

	require.ErrorIs(t, err, cmd.ErrAttention)
	assert.Equal(t, 3, cmd.ExitCode(err))
	assert.Contains(t, out.String(), "~ go-review")
	assert.Contains(t, out.String(), "! ghost")
	assert.Contains(t, out.String(), "2 skills, 1 changed, 1 needs attention")
	assert.Empty(t, errOut.String(), "attention is a result, not an error message")
}
