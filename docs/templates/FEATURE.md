# Feature: [Feature Name]

[A clear, concise 2-3 sentence description explaining the goal, scope, and target audience/impact of this feature.]

## Instructions for Agents & Developers

This repository uses a structured, phased approach for implementing new features. This keeps pull requests small, focused, and aligned with our engineering standards.

### How to use this template:
1. **Create the Feature Folder**: Create a new directory under `docs/features/<feature-name>/`.
2. **Initialize Feature Tracker**: Copy this `FEATURE.md` file into `docs/features/<feature-name>/FEATURE.md` and fill out the placeholders.
3. **Initialize Phases**: For each phase in the plan, create a corresponding phase file (e.g., `phase-01-xyz.md`, `phase-02-abc.md`) in `docs/features/<feature-name>/` based on the `PHASE.md` template.
4. **Development Loop**:
   - Work on **one phase at a time**.
   - To keep context window usage low and stay focused, **only read the active phase file** during implementation. Do not load all phase files or the entire history at once unless necessary.
   - For any behavioral change, follow the **2-commit rule** from `AGENTS.md`:
     1. Write a failing test first: `test(<scope>): ...`
     2. Write the implementation to make it pass: `feat(<scope>): ...` or `fix(<scope>): ...`
   - Mark tasks (`[ ]` to `[x]`) as you complete them.
   - Run verification steps (formatting, tests, linting) at the end of each phase.
   - Propose commit messages/PR titles without committing unless explicitly asked.

---

## Phases Roadmap

| Phase | Title | Priority | Depends On | Status | File Path |
| :---: | :--- | :---: | :---: | :---: | :--- |
| **01** | [e.g. Scaffolding & Config] | High | None | ⏳ Pending | `phase-01-scaffolding.md` |
| **02** | [e.g. Core Logic & Domain] | High | 01 | ⏳ Pending | `phase-02-core-logic.md` |
| **03** | [e.g. CLI Commands Integration] | Medium | 02 | ⏳ Pending | `phase-03-cli-commands.md` |
| **04** | [e.g. End-to-End Tests & Docs] | Low | 03 | ⏳ Pending | `phase-04-verification.md` |

*Status options: ⏳ Pending | 🚀 In Progress | ✅ Completed | 🟥 Blocked*

---

## 🛠️ General Verification Checklist
Before completing the overall feature, ensure:
- [ ] Every behavioral change follows the 2-commit rule (`test(...)` followed by `feat(...)` / `fix(...)`).
- [ ] No AI-attribution or generated-by trailers are present in commits/PR descriptions.
- [ ] Modern Go guidelines applied (where relevant).
- [ ] Formatted, linted, and all test suites pass (`go test ./...` and `go vet ./...`).
- [ ] Documentation (`docs/SPEC.md` or README) updated to reflect the new feature.
