package skill_test

import (
	"io/fs"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/skill"
)

func TestCatalog(t *testing.T) {
	t.Parallel()
	fsys := reversed{fstest.MapFS{
		"skills/sql-review/SKILL.md": {Data: []byte(published("sql-review", "Reviews SQL queries.", "[sql]"))},
		"skills/go-review/SKILL.md":  {Data: []byte(goReview)},
		"skills/plain/SKILL.md":      {Data: []byte(plain)},
		"skills/standard/SKILL.md":   {Data: []byte("---\nname: standard\ndescription: Has no store fields.\nlicense: MIT\n---\n# Standard\n")},
	}}

	skills, err := skill.Catalog(fsys)

	require.NoError(t, err)
	assert.Equal(t, []skill.Skill{
		{
			Name:        "go-review",
			Description: "Reviews Go code for correctness and idiom.",
			Status:      skill.Published,
			Tags:        []string{"go", "review"},
			Dir:         "skills/go-review",
		},
		{
			Name:        "sql-review",
			Description: "Reviews SQL queries.",
			Status:      skill.Published,
			Tags:        []string{"sql"},
			Dir:         "skills/sql-review",
		},
		{
			Name:        "standard",
			Description: "Has no store fields.",
			Status:      skill.Published,
			Dir:         "skills/standard",
		},
	}, skills, "published only, sorted by name whatever the filesystem order; a status-free skill is published")
}

func TestCatalog_empty(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fsys fstest.MapFS
	}{
		"no skills directory":    {fsys: fstest.MapFS{"README.md": {Data: []byte("# store\n")}}},
		"empty skills directory": {fsys: fstest.MapFS{"skills": {Mode: fs.ModeDir}}},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			skills, err := skill.Catalog(tt.fsys)

			require.NoError(t, err)
			assert.Empty(t, skills)
		})
	}
}

func TestCatalog_invalidSkill(t *testing.T) {
	t.Parallel()
	fsys := fstest.MapFS{
		"skills/go-review/SKILL.md": {Data: []byte(goReview)},
		"skills/broken/SKILL.md":    {Data: []byte("---\nname: broken\nstatus: published\n---\n")},
	}

	skills, err := skill.Catalog(fsys)

	require.Error(t, err)
	assert.ErrorContains(t, err, "skills/broken/SKILL.md")
	assert.Nil(t, skills, "no partial catalog")
}

type reversed struct{ fstest.MapFS }

func (r reversed) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := r.MapFS.ReadDir(name)
	slices.Reverse(entries)
	return entries, err
}

func published(name, description, tags string) string {
	return "---\nname: " + name + "\ndescription: " + description + "\nstatus: published\ntags: " + tags + "\n---\n\n# " + name + "\n"
}
