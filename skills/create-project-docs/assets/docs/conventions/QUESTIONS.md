# Question log

## Format

- Keep actual questions needing user input in `docs/QUESTIONS.md`, grouped by
  purpose: product scope, a feature, data retention, deployment, or another
  relevant topic. Add groups only when real questions exist.
- Write one short, direct question per task: `- [ ] Question?`.
- Keep unchecked questions before checked ones within each group. Leave
  completed questions in their group, with a short link to the documented answer.
- Use plain Markdown. No question IDs, HTML anchors, mandatory metadata fields,
  or separate status sections. Add a brief blocking note only when necessary
  to explain why an answer is needed next.
- Keep the log empty until real questions arise. Do not copy these examples:

```markdown
## Account access

- [ ] Can members invite other users?
- [ ] Do invitations expire?

## Data retention

- [ ] How long should deleted account data be retained?
- [x] Can users export their data? — [Documented answer](features/0003-accounts/requirements.md).
```

## During work

- Check relevant docs, earlier answers, and the matching group before asking.
  Reuse existing questions instead of duplicating them.
- Record questions and blockers needing user input when discovered. Use the
  available ask-question tool whenever it supports the question; otherwise ask
  in chat. Ask focused rounds with useful options when appropriate.
- Other documents link to the relevant group in the log rather than keeping
  their own question lists. If a log already exists elsewhere, agree on its
  owner before creating another one or moving content.
- Apply clear answers to the relevant documentation section within approved
  scope. Verify the update, then mark the question `[x]` and link the answer.
  Receiving an answer alone is not completion.
- For partial answers, document the known part and leave the task unchecked;
  simplify its question to what is still missing. Ambiguous or conflicting
  answers also stay open until clarified. Failed or unauthorized document edits
  do not close questions.
- Remove duplicates while preserving any missing context. Mark withdrawn or
  replaced questions with a short reason rather than pretending they were
  answered. Do not record secrets or invent questions for problems the agent
  can resolve within scope.
- On resume, read unchecked questions in relevant groups and their associated
  documents; consult completed entries only when needed. Check the log at
  review boundaries and session end so it reflects applied answers and blockers.
- End work replies with one concrete **Next Step**, referring to the relevant
  group and question when blocked. Answers and next-step proposals do not grant
  additional implementation, commit, or publication permission.
