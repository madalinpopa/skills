package install_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/install"
)

func TestPlanFiles(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		want, have, base install.Files
		expect           install.Action
	}{
		"absent on disk is added": {
			want:   install.Files{"SKILL.md": "v2"},
			expect: install.ActionAdd,
		},
		"have equal to want is unchanged": {
			want:   install.Files{"SKILL.md": "v2"},
			have:   install.Files{"SKILL.md": "v2"},
			base:   install.Files{"SKILL.md": "v1"},
			expect: install.ActionUnchanged,
		},
		"have equal to base is a safe update": {
			want:   install.Files{"SKILL.md": "v2"},
			have:   install.Files{"SKILL.md": "v1"},
			base:   install.Files{"SKILL.md": "v1"},
			expect: install.ActionUpdate,
		},
		"edited file is a conflict": {
			want:   install.Files{"SKILL.md": "v2"},
			have:   install.Files{"SKILL.md": "edited"},
			base:   install.Files{"SKILL.md": "v1"},
			expect: install.ActionConflict,
		},
		"upstream-removed file still at base is removed": {
			have:   install.Files{"SKILL.md": "v1"},
			base:   install.Files{"SKILL.md": "v1"},
			expect: install.ActionRemove,
		},
		"upstream-removed file edited locally is a conflict": {
			have:   install.Files{"SKILL.md": "edited"},
			base:   install.Files{"SKILL.md": "v1"},
			expect: install.ActionConflict,
		},
		"local-only file is kept as user data": {
			have:   install.Files{"SKILL.md": "notes"},
			expect: install.ActionKeep,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			changes := install.PlanFiles(tt.want, tt.have, tt.base)

			assert.Equal(t, []install.FileChange{{Path: "SKILL.md", Action: tt.expect}}, changes)
		})
	}
}

func TestPlan_states(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		target install.Target
		expect install.State
	}{
		"fresh install is add": {
			target: install.Target{
				Dir:  ".claude/skills/go-review",
				Want: install.Files{"SKILL.md": "v1", "references/style.md": "s1"},
			},
			expect: install.StateAdd,
		},
		"installed and unchanged": {
			target: install.Target{
				Dir:    ".claude/skills/go-review",
				Source: "https://example.com/store",
				Want:   install.Files{"SKILL.md": "v1"},
				Have:   install.Files{"SKILL.md": "v1"},
				Base:   install.Files{"SKILL.md": "v1"},
			},
			expect: install.StateUnchanged,
		},
		"installed and changed upstream is update": {
			target: install.Target{
				Dir:    ".claude/skills/go-review",
				Source: "https://example.com/store",
				Want:   install.Files{"SKILL.md": "v2"},
				Have:   install.Files{"SKILL.md": "v1"},
				Base:   install.Files{"SKILL.md": "v1"},
			},
			expect: install.StateUpdate,
		},
		"upstream dropped a file is update": {
			target: install.Target{
				Dir:    ".claude/skills/go-review",
				Source: "https://example.com/store",
				Want:   install.Files{"SKILL.md": "v1"},
				Have:   install.Files{"SKILL.md": "v1", "references/old.md": "o1"},
				Base:   install.Files{"SKILL.md": "v1", "references/old.md": "o1"},
			},
			expect: install.StateUpdate,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			spec := install.Spec{
				Name:      "go-review",
				Available: true,
				Source:    "https://example.com/store",
				Targets:   []install.Target{tt.target},
			}

			plans := install.Plan([]install.Spec{spec})

			require.Len(t, plans, 1)
			assert.Equal(t, tt.expect, plans[0].State)
		})
	}
}

func TestPlan_missingTargetLockIsWork(t *testing.T) {
	t.Parallel()
	spec := install.Spec{
		Name:      "go-review",
		Available: true,
		Source:    "https://example.com/store",
		Targets: []install.Target{
			{
				Dir:    ".claude/skills/go-review",
				Source: "https://example.com/store",
				Want:   install.Files{"SKILL.md": "v1"},
				Have:   install.Files{"SKILL.md": "v1"},
				Base:   install.Files{"SKILL.md": "v1"},
			},
			{
				Dir:  ".agents/skills/go-review",
				Want: install.Files{"SKILL.md": "v1"},
				Have: install.Files{"SKILL.md": "v1"},
			},
		},
	}

	plans := install.Plan([]install.Spec{spec})

	require.Len(t, plans, 1)
	assert.Equal(t, install.StateAdd, plans[0].State, "identical unmanaged content still needs its lock")
	assert.Empty(t, plans[0].Conflicts)
}

func TestPlan_conflictSkipsWholeSkill(t *testing.T) {
	t.Parallel()
	spec := install.Spec{
		Name:      "go-review",
		Available: true,
		Source:    "https://example.com/store",
		Targets: []install.Target{
			{
				Dir:    ".claude/skills/go-review",
				Source: "https://example.com/store",
				Want:   install.Files{"SKILL.md": "v2"},
				Have:   install.Files{"SKILL.md": "edited"},
				Base:   install.Files{"SKILL.md": "v1"},
			},
			{
				Dir:    ".agents/skills/go-review",
				Source: "https://example.com/store",
				Want:   install.Files{"SKILL.md": "v2"},
				Have:   install.Files{"SKILL.md": "v1"},
				Base:   install.Files{"SKILL.md": "v1"},
			},
		},
	}

	plans := install.Plan([]install.Spec{spec})

	require.Len(t, plans, 1)
	assert.Equal(t, install.StateConflict, plans[0].State)
	assert.Equal(t, []string{".claude/skills/go-review/SKILL.md"}, plans[0].Conflicts,
		"only the edited file is reported, but the whole skill is skipped")
}

func TestPlan_unsupportedTargetSkipsWholeSkill(t *testing.T) {
	t.Parallel()
	issue := install.Issue{Path: ".claude/skills/go-review/references", Reason: "is a symlink"}
	spec := install.Spec{
		Name:      "go-review",
		Available: true,
		Source:    "https://example.com/store",
		Targets: []install.Target{
			{
				Dir:    ".claude/skills/go-review",
				Source: "https://example.com/store",
				Want:   install.Files{"SKILL.md": "v2"},
				Have:   install.Files{"SKILL.md": "v1"},
				Base:   install.Files{"SKILL.md": "v1"},
				Issues: []install.Issue{issue},
			},
			{
				Dir:    ".agents/skills/go-review",
				Source: "https://example.com/store",
				Want:   install.Files{"SKILL.md": "v2"},
				Have:   install.Files{"SKILL.md": "edited"},
				Base:   install.Files{"SKILL.md": "v1"},
			},
		},
	}

	plans := install.Plan([]install.Spec{spec})

	require.Len(t, plans, 1)
	assert.Equal(t, install.StateUnsupported, plans[0].State, "an unsupported tree wins over a conflict, since force cannot resolve it")
	assert.Equal(t, []install.Issue{issue}, plans[0].Issues)
}

func TestPlan_conflictKeepsOtherSkills(t *testing.T) {
	t.Parallel()
	specs := []install.Spec{
		{
			Name:      "sql-review",
			Available: true,
			Source:    "https://example.com/store",
			Targets: []install.Target{{
				Dir:    ".claude/skills/sql-review",
				Source: "https://example.com/store",
				Want:   install.Files{"SKILL.md": "v2"},
				Have:   install.Files{"SKILL.md": "edited"},
				Base:   install.Files{"SKILL.md": "v1"},
			}},
		},
		{
			Name:      "go-review",
			Available: true,
			Source:    "https://example.com/store",
			Targets: []install.Target{{
				Dir:  ".claude/skills/go-review",
				Want: install.Files{"SKILL.md": "v1"},
			}},
		},
	}

	plans := install.Plan(specs)

	require.Len(t, plans, 2)
	assert.Equal(t, "go-review", plans[0].Name)
	assert.Equal(t, install.StateAdd, plans[0].State)
	assert.Equal(t, "sql-review", plans[1].Name)
	assert.Equal(t, install.StateConflict, plans[1].State)
}

func TestPlan_unavailableIsNotRemoved(t *testing.T) {
	t.Parallel()
	spec := install.Spec{
		Name:      "django-testing",
		Available: false,
		Source:    "https://example.com/store",
		Targets: []install.Target{{
			Dir:    ".claude/skills/django-testing",
			Source: "https://example.com/store",
			Have:   install.Files{"SKILL.md": "v1"},
			Base:   install.Files{"SKILL.md": "v1"},
		}},
	}

	plans := install.Plan([]install.Spec{spec})

	require.Len(t, plans, 1)
	assert.Equal(t, install.StateUnavailable, plans[0].State)
	for _, target := range plans[0].Targets {
		for _, change := range target.Files {
			assert.NotEqual(t, install.ActionRemove, change.Action, "an unavailable skill stays installed")
		}
	}
}

func TestPlan_foreignSourceNeedsAttention(t *testing.T) {
	t.Parallel()
	spec := install.Spec{
		Name:      "go-review",
		Available: true,
		Source:    "https://example.com/store",
		Targets: []install.Target{{
			Dir:    ".claude/skills/go-review",
			Source: "https://example.com/another-store",
			Want:   install.Files{"SKILL.md": "v2"},
			Have:   install.Files{"SKILL.md": "v1"},
			Base:   install.Files{"SKILL.md": "v1"},
		}},
	}

	plans := install.Plan([]install.Spec{spec})

	require.Len(t, plans, 1)
	assert.Equal(t, install.StateForeign, plans[0].State)
	assert.Equal(t, "https://example.com/another-store", plans[0].Source, "cmd can name the other source")
}

func TestPlan_isDeterministic(t *testing.T) {
	t.Parallel()
	specs := []install.Spec{
		{
			Name:      "sql-review",
			Available: true,
			Source:    "https://example.com/store",
			Targets: []install.Target{{
				Dir:  ".claude/skills/sql-review",
				Want: install.Files{"scripts/run.sh": "r1", "SKILL.md": "v1", "references/a.md": "a1"},
			}},
		},
		{
			Name:      "go-review",
			Available: true,
			Source:    "https://example.com/store",
			Targets: []install.Target{{
				Dir:  ".claude/skills/go-review",
				Want: install.Files{"SKILL.md": "v1"},
			}},
		},
	}

	first := install.Plan(specs)
	second := install.Plan(specs)

	assert.Equal(t, first, second)
	assert.Equal(t, []install.SkillPlan{
		{
			Name:  "go-review",
			State: install.StateAdd,
			Targets: []install.TargetPlan{{
				Dir:   ".claude/skills/go-review",
				Files: []install.FileChange{{Path: "SKILL.md", Action: install.ActionAdd}},
			}},
		},
		{
			Name:  "sql-review",
			State: install.StateAdd,
			Targets: []install.TargetPlan{{
				Dir: ".claude/skills/sql-review",
				Files: []install.FileChange{
					{Path: "SKILL.md", Action: install.ActionAdd},
					{Path: "references/a.md", Action: install.ActionAdd},
					{Path: "scripts/run.sh", Action: install.ActionAdd},
				},
			}},
		},
	}, first, "skills by name, files by path")
}
