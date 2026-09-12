package skill_test

import (
	"maps"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/madalinpopa/skills/internal/skill"
)

const goReview = `---
name: go-review
description: Reviews Go code for correctness and idiom.
status: published
tags: [go, review]
x-claude:
  disable-model-invocation: true
  allowed-tools: Bash(go test *)
---

# Go review

Read the diff, then the tests.
`

const plain = `---
name: plain
description: A skill with nothing agent specific.
status: draft
---

# Plain
`

func TestParse_success(t *testing.T) {
	t.Parallel()

	meta, err := skill.Parse([]byte(goReview))

	require.NoError(t, err)
	assert.Equal(t, skill.Metadata{
		Name:        "go-review",
		Description: "Reviews Go code for correctness and idiom.",
		Status:      skill.Published,
		Tags:        []string{"go", "review"},
	}, meta)
}

func TestParse_failed(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source string
		want   string
	}{
		"missing name":         {source: "---\ndescription: d\nstatus: published\n---\n", want: "name"},
		"missing description":  {source: "---\nname: x\nstatus: published\n---\n", want: "description"},
		"missing status":       {source: "---\nname: x\ndescription: d\n---\n", want: "status"},
		"unsupported status":   {source: "---\nname: x\ndescription: d\nstatus: hidden\n---\n", want: "hidden"},
		"no frontmatter":       {source: "# Just markdown\n", want: "frontmatter"},
		"unclosed frontmatter": {source: "---\nname: x\n", want: "frontmatter"},
		"malformed yaml":       {source: "---\nname: [x\n---\n", want: "yaml"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := skill.Parse([]byte(tt.source))

			require.Error(t, err)
			assert.ErrorContains(t, err, tt.want)
		})
	}
}

func TestTransform(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source string
		agent  string
		want   []string
		unwant []string
	}{
		"claude lifts x-claude": {
			source: goReview,
			agent:  "claude",
			want:   []string{"name: go-review", "disable-model-invocation: true", "allowed-tools: Bash(go test *)"},
			unwant: []string{"x-claude", "status:", "tags:"},
		},
		"codex drops x-claude": {
			source: goReview,
			agent:  "codex",
			want:   []string{"name: go-review", "description: Reviews Go code for correctness and idiom."},
			unwant: []string{"x-claude", "disable-model-invocation", "allowed-tools", "status:", "tags:"},
		},
		"gemini drops x-claude": {
			source: goReview,
			agent:  "gemini",
			unwant: []string{"x-claude", "allowed-tools"},
		},
		"plain skill only loses store fields": {
			source: plain,
			agent:  "claude",
			want:   []string{"name: plain", "description: A skill with nothing agent specific."},
			unwant: []string{"status:"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out, err := skill.Transform([]byte(tt.source), tt.agent)

			require.NoError(t, err)
			text := string(out)
			for _, line := range tt.want {
				assert.Contains(t, text, line+"\n")
			}
			for _, fragment := range tt.unwant {
				assert.NotContains(t, text, fragment)
			}
			_, body, found := strings.Cut(tt.source, "\n---\n")
			require.True(t, found)
			assert.True(t, strings.HasSuffix(text, "\n---\n"+body), "the body is preserved byte for byte")
		})
	}
}

func TestTransform_isDeterministic(t *testing.T) {
	t.Parallel()

	first, err := skill.Transform([]byte(goReview), "claude")
	require.NoError(t, err)
	second, err := skill.Transform([]byte(goReview), "claude")
	require.NoError(t, err)

	assert.Equal(t, first, second)
}

func TestRender(t *testing.T) {
	t.Parallel()
	fsys := fstest.MapFS{
		"SKILL.md":            {Data: []byte(goReview)},
		"agents/openai.yaml":  {Data: []byte("interface:\n  display_name: Go review\n")},
		"references/style.md": {Data: []byte("# Style\n")},
		"scripts/check.sh":    {Data: []byte("#!/bin/sh\n")},
		"assets/logo.svg":     {Data: []byte("<svg/>")},
	}

	tests := map[string]struct {
		agent string
		paths []string
	}{
		"claude skips openai.yaml": {
			agent: "claude",
			paths: []string{"SKILL.md", "assets/logo.svg", "references/style.md", "scripts/check.sh"},
		},
		"codex includes openai.yaml": {
			agent: "codex",
			paths: []string{"SKILL.md", "agents/openai.yaml", "assets/logo.svg", "references/style.md", "scripts/check.sh"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			files, err := skill.Render(fsys, tt.agent)

			require.NoError(t, err)
			assert.ElementsMatch(t, tt.paths, slices.Collect(maps.Keys(files)))
			assert.Equal(t, "# Style\n", string(files["references/style.md"].Data), "ordinary files are copied unchanged")
			assert.NotContains(t, string(files["SKILL.md"].Data), "status:", "SKILL.md is transformed")
		})
	}
}

func TestDiscover(t *testing.T) {
	t.Parallel()
	fsys := fstest.MapFS{
		"README.md":                        {Data: []byte("# store\n")},
		"skills/go-review/SKILL.md":        {Data: []byte(goReview)},
		"skills/go-review/references/x.md": {Data: []byte("nested files are not skills\n")},
		"skills/plain/SKILL.md":            {Data: []byte(plain)},
		"skills/notes.txt":                 {Data: []byte("not a directory\n")},
	}

	skills, err := skill.Discover(fsys)

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
			Name:        "plain",
			Description: "A skill with nothing agent specific.",
			Status:      skill.Draft,
			Dir:         "skills/plain",
		},
	}, skills, "sorted by name, drafts included")
}

func TestDiscover_failed(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		fsys fstest.MapFS
		want string
	}{
		"name differs from directory": {
			fsys: fstest.MapFS{"skills/other/SKILL.md": {Data: []byte(goReview)}},
			want: "skills/other",
		},
		"directory without SKILL.md": {
			fsys: fstest.MapFS{"skills/empty/notes.md": {Data: []byte("x")}},
			want: "skills/empty",
		},
		"invalid metadata names the file": {
			fsys: fstest.MapFS{"skills/bad/SKILL.md": {Data: []byte("---\nname: bad\n---\n")}},
			want: "skills/bad/SKILL.md",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := skill.Discover(tt.fsys)

			require.Error(t, err)
			assert.ErrorContains(t, err, tt.want)
		})
	}
}
