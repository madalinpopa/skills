package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/store"
)

func TestCheckGitVersion_success(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		output string
	}{
		"minimum":       {output: "git version 2.45.0"},
		"apple build":   {output: "git version 2.55.0 (Apple Git-155)"},
		"windows build": {output: "git version 2.45.1.windows.1"},
		"next major":    {output: "git version 3.0.0"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := store.CheckGitVersion(tt.output)

			assert.NoError(t, err)
		})
	}
}

func TestCheckGitVersion_failed(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		output string
		want   string
	}{
		"older minor":        {output: "git version 2.44.2", want: "2.44.2"},
		"single digit minor": {output: "git version 2.9.5", want: "2.9.5"},
		"unexpected output":  {output: "not git", want: "not git"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := store.CheckGitVersion(tt.output)

			require.Error(t, err)
			assert.ErrorContains(t, err, "2.45", "the error names the minimum version")
			assert.ErrorContains(t, err, tt.want, "the error names what was found")
		})
	}
}
