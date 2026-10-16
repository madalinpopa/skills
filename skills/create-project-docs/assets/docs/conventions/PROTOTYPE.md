# Throwaway prototypes

Read only when the user requests a disposable experiment or asks to continue
an existing throwaway prototype. Confirm existing intent from its README,
directory documentation, or branch/bookmark context; a filename alone is not
enough. Clarify uncertain intent before applying this workflow.

## Implement and demonstrate

- Build the smallest experiment that answers the user's question about an idea,
  feature, design, or architecture. Skip TDD, automated test writing, full
  verification, adversarial review, formal feature records, and intermediate
  planning/review gates.
- Preserve unrelated work. Use a separate `prototype/<name>` Git branch, or a
  separate jj change with a `prototype/<name>` bookmark, following project tooling.
  Reuse suitable existing isolation; isolate mixed prototype files before editing.
- A request covered by this workflow authorizes local commits of the experiment
  without another approval gate. Inspect the diff, exclude unrelated work, and
  keep the branch/bookmark at the latest prototype commit. Use a message such as
  `chore(prototype): explore search ranking`; preserve human attribution.
- Run only the build or manual smoke check needed to demonstrate the idea.
  Report what works, the check and its result, limitations, lessons learned, and
  branch/bookmark and commit references. Do not claim production readiness or
  passing tests.

## Boundaries and handoff

- Keep normal dependency permissions. Push, publish, merge, and deletion require
  explicit user requests. Discard prototype code when requested, preserving
  unrelated work.
- For production implementation, return to requirements, planning, and the
  production workflow. Prototype results inform decisions; they do not
  establish agreed requirements or completed roadmap work.
