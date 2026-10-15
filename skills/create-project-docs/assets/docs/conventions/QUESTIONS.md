# Question log

## Ownership and loading

- `docs/QUESTIONS.md` is the single log for actual unresolved questions and
  blockers needing user input across specs, features, stack, and decisions.
- Keep answers in their owning documents. Other records link to question IDs
  instead of keeping separate question lists.
- Read open entries and the relevant target documents before asking questions
  or resuming dependent work. Consult resolved entries only when relevant.
- Start the log empty. Add headings only when entries exist; keep `Open` before
  `Resolved`. Do not seed generic questions or example entries.
- If another file already owns the log, agree on its owner before creating a
  competing file or consolidating existing content. Respect edit permissions.

## Record and ask

- Check repository evidence, earlier answers, and existing entries first.
- Record each actual question when discovered, before asking or stopping
  dependent work. Use one task per independently answerable question.
- Assign the next unused `Q-NNNN` ID across open and resolved entries. Preserve
  IDs; never reuse or renumber them. Link entries using the format
  `docs/QUESTIONS.md#q-0001` with an explicit `<a id="q-0001"></a>` on the task.
- State the question, the exact information needed, and the target document
  or section. If the section does not exist yet, link the file and name the
  intended section in text; do not invent an anchor or create an empty section.
- Add `Blocks` only when specific work cannot continue without the answer.
  Link that work and state why. Continue independent authorized work.
- Use the available ask-question tool whenever it supports the question;
  otherwise ask in chat. Offer evidence-based options and meaningful tradeoffs
  when useful. Do not turn a focused question into a project-wide interview.
- Do not log ordinary problems the agent can resolve within approved scope as
  questions for the user. Never record secrets or ask for them in the log.

Entry format, shown here for guidance only:

```markdown
## Open

- [ ] <a id="q-0001"></a> Q-0001 — How long must account data be retained?
  - Needed: retention duration and exceptions.
  - Target: [Account requirements](features/0003-accounts/requirements.md), Retention section.
  - Blocks: [Deletion task](features/0003-accounts/plan.md) — expiry depends on the duration.
```

## Apply answers and close

- Apply each clear answer to its owning section within authorized scope.
  Preserve agreed content; record significant decisions in their proper record.
- Check the update against the answer and verify its link before checking the
  question `[x]`. Receiving an answer alone does not complete the task.
- For partial answers, document the established part, keep the task `[ ]`, and
  update `Needed` to identify only the missing information. Link applied content.
- If an answer conflicts with agreed requirements or is ambiguous, keep the
  entry open and ask about the specific conflict before changing the contract.
- If writing the target requires approval or fails, keep the entry open and
  state what remains. Do not imply the answer was applied successfully.
- Move completed entries to `Resolved`, preserving IDs. Replace transient
  details with a short resolution, answer source/date, and target link. Avoid
  copying the full answer or conversation into the log.
- Close withdrawn, duplicate, or replaced questions with an explicit disposition
  and reason; link the surviving question when applicable. Do not label them
  answered. If a resolved topic needs reconsideration, open a new linked question
  rather than silently changing the historical resolution.
- Answers and checked questions grant no implementation, commit, or publication
  permission beyond what the user has authorized.

## Report and resume

- Keep unchecked entries visible under `Open`; resolved history stays in the
  same log. Remove headings that no longer contain entries.
- At review boundaries and session end, check that the log matches applied
  answers, remaining inputs, and blocked tasks. Preserve unrelated entries.
- Keep replies concise. End with one concrete **Next Step** based on current
  progress, linking the blocking question, required review, or next useful task.
  Do not repeat the log or invent follow-up work when the request is complete.
