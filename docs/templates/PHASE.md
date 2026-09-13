# Phase [X]: [Phase Title]

## 📝 Description
[Provide a short, concise description explaining what is achieved in this phase, focusing on the architectural or user-facing value.]

## 📋 Metadata
- **Priority**: [High / Medium / Low]
- **Depends On**: [None / Phase Y]
- **Status**: [⏳ Pending / 🚀 In Progress / ✅ Completed / 🟥 Blocked]
- **Estimated Commits**: [1-5 commits]
- **Suggested PR Title**: *[To be suggested by LLM when all tasks are complete and verified]*

---

## 🎯 Acceptance Criteria
[Define clear, verifiable conditions that must be met to consider this phase complete. These must be empirically testable.]
- [ ] [e.g., Unit tests for the new package pass with 100% success.]
- [ ] [e.g., Calling the command with `--flag` returns the correct output structure.]
- [ ] [e.g., Verification steps (`go test ./...` and `go vet ./...`) pass without errors.]

---

## 🚶 Implementation Plan & Tasks

Follow the **Arrange, Act, Assert** testing pattern and the **2-commit rule** (`test(...)` commit followed by `feat(...)`/`fix(...)` commit) for all changes.

### 🧪 Step 1: [Short name for Step 1]
- [ ] [Task 1: e.g., Write a failing unit test in `internal/install/plan_test.go` verifying target resolution.]
- [ ] [Task 2: e.g., Implement the core resolution logic in `internal/install/plan.go` to pass the test.]
- [ ] **Commit Plan**:
  - `test(install): add test case for target resolution failure`
  - `feat(install): resolve targets correctly in installer plan`

### 🧪 Step 2: [Short name for Step 2]
- [ ] [Task 1: e.g., Add integration tests for edge cases.]
- [ ] [Task 2: e.g., Update CLI commands to expose the new functionality.]
- [ ] **Commit Plan**:
  - `test(cmd): assert new flags are registered and processed`
  - `feat(cmd): expose target resolution via CLI command`

---

## 🏁 Phase Completion Checklist
- [ ] Run modern Go guidelines and `go fix ./...`.
- [ ] Format and vet the code (`task format` and `go vet ./...`).
- [ ] Run all tests and verify they pass (`go test ./...`).
- [ ] If all tasks are checked, update the **Status** to `✅ Completed` in this file and in `FEATURE.md`.
- [ ] Propose a clear and descriptive PR title below (without committing it).

### Suggested PR Title:
> [e.g., Support target resolution in installation planner]
