---
name: demo
description: Demonstrates the standard structure, metadata, and files for agent skills. Use when the user asks for a demo, example, template, or tests the installation of skills.
status: published
tags:
  - demo
  - example
  - guidelines
  - testing
x-claude:
  disable-model-invocation: true
  user-invocable: true
  allowed-tools:
    - "Bash(echo *)"
    - "Read"
  context: inline
---

# Demo Skill

This is a professional reference and template demonstrating the standard structure, metadata format, and guidelines for authoring agent-agnostic skills that work seamlessly across Claude Code, OpenAI Codex/ChatGPT, and Gemini CLI.

## Purpose

The `demo` skill acts as an educational blueprint and functional test case for the `skills` CLI. It shows how the tool resolves differences between platforms while maintaining a single, clean source of truth.

## Standard Layout

A standard skill directory in the store is structured as follows:

```text
skills/
  demo/
    SKILL.md                <-- Core markdown instructions & shared frontmatter
    agents/
      openai.yaml           <-- Codex/ChatGPT-specific settings & MCP dependencies
    references/
      guidelines.md         <-- Optional helper files copied to all agent targets
```

## Best Practices

To author high-quality, portable skills:

1. **Write Descriptive Triggers**: The `description` field in the frontmatter is the primary discovery and routing mechanism. Write in the third person, lead with the primary trigger words (e.g., "Demonstrates...", "Reviews..."), and be highly specific.
2. **Keep the Body Concise**: The body of your skill is injected directly into the LLM's context window. Keep it under 500 lines to preserve token budget.
3. **Use Portable Paths**: Always use forward slashes (`/`) for paths.
4. **Be Unambiguous**: Use consistent terminology throughout the file to avoid confusing the model.
5. **Separate Vendor Extras**:
   - Claude Code specific controls are placed inside the `x-claude` block in `SKILL.md`.
   - OpenAI Codex/ChatGPT specific settings live in `agents/openai.yaml`.
