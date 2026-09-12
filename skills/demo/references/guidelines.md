# Agent Skills Authoring Guidelines

This reference details the best practices for authoring skills compatible with Claude Code, OpenAI Codex, and Gemini CLI.

## General Best Practices

- **The description decides whether the skill fires.** Say what it does and when to use it, lead with the trigger words, and write in the third person.
- **Be concise.** The body shares the context window with everything else. Keep it under 500 lines.
- **Keep file references one level deep.** Nested references get partially read.
- **Use forward slashes in paths**, and avoid dates and version notes that age.
- **Use one term for one thing** throughout a skill.
- **Keep each skill focused** on one job.

## Platform Differences

### Claude Code
- Command identity is taken from the directory name.
- Custom frontmatter properties like `allowed-tools` and `disable-model-invocation` should be nested under `x-claude`.

### OpenAI Codex / ChatGPT
- Command identity is taken from the `name` field in the frontmatter.
- OpenAI-specific UI config and MCP tool dependencies go in `agents/openai.yaml`.

### Gemini CLI
- Command identity is taken from the `name` field in the frontmatter.
- Uses the same `SKILL.md` format but strips store-specific or Claude-specific fields upon installation.
