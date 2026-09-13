package skill_test

import (
	"maps"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v3"

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

func TestParse_optionalPublishingMetadata(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source string
		status skill.Status
		tags   []string
	}{
		"absent status and tags": {
			source: "---\nname: x\ndescription: d\n---\n",
			status: skill.Published,
		},
		"explicit draft": {
			source: "---\nname: x\ndescription: d\nstatus: draft\n---\n",
			status: skill.Draft,
		},
		"empty tags list": {
			source: "---\nname: x\ndescription: d\nstatus: published\ntags: []\n---\n",
			status: skill.Published,
		},
		"tags list": {
			source: "---\nname: x\ndescription: d\nstatus: published\ntags: [go, review]\n---\n",
			status: skill.Published,
			tags:   []string{"go", "review"},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			meta, err := skill.Parse([]byte(tt.source))

			require.NoError(t, err)
			assert.Equal(t, tt.status, meta.Status, "an absent status means published")
			if len(tt.tags) == 0 {
				assert.Empty(t, meta.Tags, "absent or empty tags mean no tags")
			} else {
				assert.Equal(t, tt.tags, meta.Tags)
			}
		})
	}
}

func TestParse_failed(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source string
		want   string
	}{
		"missing name":         {source: "---\ndescription: d\nstatus: published\n---\n", want: "name"},
		"missing description":  {source: "---\nname: x\nstatus: published\n---\n", want: "description"},
		"unsupported status":   {source: "---\nname: x\ndescription: d\nstatus: hidden\n---\n", want: "hidden"},
		"empty status":         {source: "---\nname: x\ndescription: d\nstatus: \"\"\n---\n", want: "status"},
		"null status":          {source: "---\nname: x\ndescription: d\nstatus:\n---\n", want: "status"},
		"status is a list":     {source: "---\nname: x\ndescription: d\nstatus: [published]\n---\n", want: "status"},
		"null tags":            {source: "---\nname: x\ndescription: d\nstatus: published\ntags: null\n---\n", want: "tags"},
		"tags is a scalar":     {source: "---\nname: x\ndescription: d\nstatus: published\ntags: go\n---\n", want: "tags"},
		"numeric tag":          {source: "---\nname: x\ndescription: d\nstatus: published\ntags: [go, 1]\n---\n", want: "tags"},
		"null tag":             {source: "---\nname: x\ndescription: d\nstatus: published\ntags: [go, null]\n---\n", want: "tags"},
		"nested tag":           {source: "---\nname: x\ndescription: d\nstatus: published\ntags: [go, [review]]\n---\n", want: "tags"},
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

const (
	notMapping = "---\nname: x\ndescription: d\nstatus: published\nx-claude: [allowed-tools]\n---\n# x\n"
	duplicate  = "---\nname: x\ndescription: d\nstatus: published\nx-claude:\n  allowed-tools: a\n  allowed-tools: b\n---\n# x\n"
	reserved   = "---\nname: x\ndescription: d\nstatus: published\nx-claude:\n  name: other\n---\n# x\n"
	collision  = "---\nname: x\ndescription: d\nstatus: published\nallowed-tools: a\nx-claude:\n  allowed-tools: b\n---\n# x\n"
	numericKey = "---\nname: x\ndescription: d\nstatus: published\nx-claude:\n  1: a\n---\n# x\n"
	twoNames   = "---\nname: x\nname: y\ndescription: d\nstatus: published\n---\n# x\n"
)

func TestParse_rejectsUnsafeExtensions(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source string
		want   string
	}{
		"x-claude is not a mapping":    {source: notMapping, want: "x-claude"},
		"duplicate extension key":      {source: duplicate, want: "allowed-tools"},
		"reserved extension key":       {source: reserved, want: "name"},
		"extension collides top-level": {source: collision, want: "allowed-tools"},
		"non-string extension key":     {source: numericKey, want: "x-claude"},
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

func TestTransform_failed(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source string
		agent  string
		want   string
	}{
		"duplicate extension key for claude": {source: duplicate, agent: "claude", want: "allowed-tools"},
		"duplicate extension key for codex":  {source: duplicate, agent: "codex", want: "allowed-tools"},
		"collision for codex":                {source: collision, agent: "codex", want: "allowed-tools"},
		"not a mapping for codex":            {source: notMapping, agent: "codex", want: "x-claude"},
		"duplicate top-level key":            {source: twoNames, agent: "claude", want: "name"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := skill.Transform([]byte(tt.source), tt.agent)

			require.Error(t, err, "invalid extensions are rejected even when the output would drop them")
			assert.ErrorContains(t, err, tt.want)
		})
	}
}

const crlf = "---\r\nname: x\r\ndescription: d\r\nstatus: published\r\n---\r\n# x\r\n\r\nWindows body.\r\n"

func TestParse_acceptsCRLF(t *testing.T) {
	t.Parallel()

	meta, err := skill.Parse([]byte(crlf))

	require.NoError(t, err)
	assert.Equal(t, skill.Metadata{Name: "x", Description: "d", Status: skill.Published}, meta)
}

func TestTransform_normalisesDelimitersAndKeepsBody(t *testing.T) {
	t.Parallel()
	lf := strings.ReplaceAll(crlf, "\r\n", "\n")
	fromLF, err := skill.Transform([]byte(lf), "claude")
	require.NoError(t, err)

	out, err := skill.Transform([]byte(crlf), "claude")

	require.NoError(t, err)
	front, body, found := strings.Cut(strings.TrimPrefix(string(out), "---\n"), "\n---\n")
	require.True(t, found, "the emitted frontmatter uses LF delimiters")
	assert.Equal(t, "name: x\ndescription: d", front, "the frontmatter is emitted the same way for both line endings")
	assert.Equal(t, "# x\r\n\r\nWindows body.\r\n", body, "the body keeps its bytes")
	assert.Equal(t, fromLF[:len(fromLF)-len("# x\n\nWindows body.\n")], out[:len(out)-len(body)])
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
		"custom shared agent drops x-claude": {
			source: goReview,
			agent:  "amp",
			want:   []string{"name: go-review", "description: Reviews Go code for correctness and idiom."},
			unwant: []string{"x-claude", "disable-model-invocation", "allowed-tools", "status:", "tags:"},
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

func TestTransform_preservesAliasesToStoreMetadata(t *testing.T) {
	t.Parallel()
	data := []byte("---\nname: demo\nstatus: published\ntags: [&summary Demonstrates skill installation.]\ndescription: *summary\n---\n# Demo\n")
	before, err := skill.Parse(data)
	require.NoError(t, err, "the source has valid YAML and store metadata")

	tests := map[string]struct {
		agent string
	}{
		"claude": {agent: "claude"},
		"shared": {agent: "agents"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out, err := skill.Transform(data, tt.agent)

			require.NoError(t, err)
			after, err := skill.Describe(out)
			require.NoError(t, err, "stripping tags must not leave an undefined YAML alias")
			assert.Equal(t, before.Description, after.Description)
			assert.NotContains(t, string(out), "tags:")
		})
	}
}

const compatibilityBody = "\n# Body\n\nKeep --- and `x-claude:` as text.\n---\n"

const standardFields = `---
name: pdf-processing
description: Extracts PDF text. Use when handling PDFs.
license: Apache-2.0
compatibility: Requires Python 3.14+ and uv
metadata:
  author: example-org
  version: "1.0"
allowed-tools: Bash(git:*) Bash(jq:*) Read
---
` + compatibilityBody

const nativeClaudeFields = `---
name: release-notes
description: Drafts release notes. Use when preparing a release.
when_to_use: The user asks for a changelog.
argument-hint: "[version]"
allowed-tools:
  - Read
  - Bash(git log *)
model: inherit
effort: high
context: fork
agent: Explore
paths:
  - "CHANGELOG.md"
hooks:
  PreToolUse:
    - matcher: Bash
      hooks:
        - type: command
          command: ./scripts/check.sh
          timeout: 30
x-future-extension:
  enabled: true
  ratio: 0.5
---
` + compatibilityBody

const scopedClaudeFields = `---
name: deploy
description: Deploys the service. Use when the user asks to deploy.
status: published
tags: [ops]
allowed-tools: Read
x-claude:
  disable-model-invocation: true
  user-invocable: false
  shell: bash
  hooks:
    Stop:
      - hooks:
          - type: command
            command: echo done
---
` + compatibilityBody

func TestTransform_preservesFrontmatterValues(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source string
		agent  string
		lifted bool
	}{
		"standard fields for claude":      {source: standardFields, agent: "claude"},
		"standard fields for shared":      {source: standardFields, agent: "codex"},
		"native claude fields for claude": {source: nativeClaudeFields, agent: "claude"},
		"native claude fields for shared": {source: nativeClaudeFields, agent: "codex"},
		"scoped claude fields for claude": {source: scopedClaudeFields, agent: "claude", lifted: true},
		"scoped claude fields for shared": {source: scopedClaudeFields, agent: "codex"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			want := frontmatterValues(t, tt.source)
			scoped, _ := want["x-claude"].(map[string]any)
			delete(want, "status")
			delete(want, "tags")
			delete(want, "x-claude")
			if tt.lifted {
				maps.Copy(want, scoped)
			}

			out, err := skill.Transform([]byte(tt.source), tt.agent)

			require.NoError(t, err)
			assert.Equal(t, want, frontmatterValues(t, string(out)), "retained values keep their YAML types")
			assert.True(t, strings.HasSuffix(string(out), "---\n"+compatibilityBody), "the body is preserved byte for byte")
		})
	}
}

func frontmatterValues(t *testing.T, source string) map[string]any {
	t.Helper()
	front, _, found := strings.Cut(strings.TrimPrefix(source, "---\n"), "\n---\n")
	require.True(t, found)
	var values map[string]any
	require.NoError(t, yaml.Unmarshal([]byte(front), &values))
	return values
}

func TestTransform_rejectsNestedDuplicateKeys(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source string
		want   string
	}{
		"metadata": {
			source: "---\nname: x\ndescription: d\nmetadata:\n  author: a\n  author: b\n---\n# x\n",
			want:   "author",
		},
		"native hooks": {
			source: "---\nname: x\ndescription: d\nhooks:\n  Stop: []\n  Stop: []\n---\n# x\n",
			want:   "Stop",
		},
		"scoped hooks": {
			source: "---\nname: x\ndescription: d\nx-claude:\n  hooks:\n    Stop: []\n    Stop: []\n---\n# x\n",
			want:   "Stop",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := skill.Transform([]byte(tt.source), "claude")

			require.Error(t, err, "emitted frontmatter must not contain duplicate keys")
			assert.ErrorContains(t, err, tt.want)
		})
	}
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
