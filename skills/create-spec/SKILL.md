---
name: create-spec
description: Fills a project's docs/SPEC.md from its docs/templates/SPEC.md by interviewing the user one section at a time, drafting answers from the code and docs already in the repository, and recording decisions, non-goals, and pending questions without inventing requirements. Use whenever a project has docs/templates/SPEC.md and the user wants to write, fill, complete, draft, or resume the spec, asks what the spec still needs, or says the spec has TODOs or placeholders left, even if they only say "let's write the spec" or "help me describe what we are building". Creates the missing docs/templates (SPEC.md, FEATURE.md, PHASE.md) from bundled copies when a project has none yet. Not for feature or phase documents (those use FEATURE.md and PHASE.md) or for changing the templates themselves.
status: published
tags: [spec, docs, planning]
---

# Create spec

Fill `docs/SPEC.md` from the project's own `docs/templates/SPEC.md`. The user
decides what the application does; the agent gathers evidence, asks the
questions the project cannot answer, writes the section, and shows it for
review. Work one section at a time and stop for confirmation after each.

## Initialize the templates

Resolve the repository root. When `docs/templates/SPEC.md` is missing, run the
bundled helper by its absolute installed path; `<skill-dir>` is this skill's
directory:

```sh
"<skill-dir>/scripts/init-templates.sh" <repository-root>
```

It creates `docs/templates/` and copies the bundled `SPEC.md`, `FEATURE.md`,
and `PHASE.md` from [assets/templates](assets/templates/) for each file that is
missing. Existing files are kept as they are. Report which files it created
and remind the user that they are new, uncommitted files.

## Locate the files

1. Read `docs/templates/SPEC.md` in full.
2. Open `docs/SPEC.md`:
   - Missing, or only headings and `TODO` markers: copy the template over it.
   - Partly filled from this template: this is a resume. List which sections
     are done, pending, or untouched, and continue with the first unfinished
     one.
   - A real spec in another shape: ask before replacing or restructuring it.

## Gather evidence before asking

Read what the project already says before asking the user anything: README,
`AGENTS.md` or `CLAUDE.md`, the module or package file, the command or route
tree, `docs/features/`, config files, and recent git history. Turn findings
into proposed answers the user confirms or corrects. Ask an open question only
when the repository holds no evidence.

## Interview one section at a time

Follow this order. Later sections link to earlier ones, so acceptance comes
last.

1. Goal
2. User
3. Commands / Features, then one Contract per capability that needs rules
4. Usage and Outputs / Outcomes
5. Implementation details, subsection by subsection
6. Suggested skills
7. Research
8. Behavioral acceptance

For each section, read its entry in [references/questions.md](references/questions.md).
Then:

- Ask a small batch of questions, three to five, with proposed answers from
  the evidence where you have them. Wait for the reply.
- Write the section into `docs/SPEC.md`. Keep the template's headings. Replace
  the bracketed prompts and delete the template's instruction paragraphs once
  the section is filled.
- Show the written section and ask the user to confirm or fix it before moving
  on.

Rules that hold for every section:

- Never invent a requirement to fill a slot. Write `Pending: <question>` and
  what is needed to resolve it.
- Keep the CLI or web variant that applies and remove the other. A prompt
  about retries, concurrency, or recovery is a question, not a feature; when
  the answer is "not supported", record it as a non-goal or an explicit limit.
- Keep proposals and decisions apart. Mark a proposal as such until the user
  decides.
- Use the user's own words for names, commands, and flags. Do not rename them.

## Research

Agree on the research questions with the user first. Then do the searches with
the tools available: official documentation for the chosen stack, GitHub and
web search for similar tools. Record direct links, versions, and today's date
in the tables. Say which findings are verified and which are assumptions. When
no search tool is available, write the queries into the section as pending so
the user can run them.

## Finish

Check the whole document before the final review:

- Every requirement in Behavioral acceptance links to a contract or feature
  and to a scenario, and every scenario names the requirement it demonstrates.
- Non-goals do not contradict the feature table.
- Every library in the dependency table appears in the stack or architecture,
  and each has a documentation link.
- No bracketed placeholder remains except inside an explicit pending item.

Report what is filled, what is pending, and the next step. Propose a commit
message. Do not commit; the user asks for that separately.

The interview often spans sessions. That is fine: the document itself holds
the state, so a later session resumes from the first unfinished section.
