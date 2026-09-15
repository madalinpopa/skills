# Writing issues

Read this when drafting a new issue or splitting work that does not fit one.
The template is in `../assets/issue-template.md`; fill it, do not restate it.

## One issue, one thing

An issue passes when every line below is true. If one fails, rewrite or
split before creating it.

- It delivers one outcome a user of the project can observe. "Add config
  loading" passes; "improve configuration" does not.
- The outcome is testable: name the command, test, or behavior that proves
  it.
- It fits in at most five commits, one task per commit. A task that needs
  several commits is a sub-issue in disguise.
- It stands alone. It may depend on a finished issue, never on an open one
  unless the parent links them as sub-issues in order.
- Tasks say what changes and where, not how. "Add a `Load` function to
  `internal/config` that reads the TOML file" is a task. A code sketch,
  field list, or algorithm is implementation and belongs to the user.

## Mentor's voice

Write as a senior engineer handing the issue to a colleague who will learn
from it. Short sentences, plain words, second person. Say why the project
needs the change before what to change. Name the concept the issue teaches in
"What you will learn"; if you cannot name one, look for a harder or different
slice of the work. Do not put questions in the issue body; questions belong
in the mentoring conversation.

## Sizing and splitting

Count commits before writing tasks. Each task should be one self-contained
commit that reviews in minutes: roughly the size Google calls a small change,
about a hundred lines and one purpose. Keep a refactor and a behavior change
in separate tasks, and keep a test with the code it covers.

When the request needs more than five commits or more than one outcome:

1. Create a parent issue with the template. Its outcome is the whole feature;
   its task list names the sub-issues instead of commits.
2. Create one sub-issue per outcome with `gh issue create --parent <N>`. Each
   sub-issue passes the checks above on its own.
3. Order them so each sub-issue leaves the project working. Prefer vertical
   slices, one thin path end to end, over layers.
4. Point the user at the first sub-issue. Work always happens in a sub-issue,
   never in the parent.

When the request is vague ("make the CLI nicer"), do not guess a feature.
Ask the user for the one behavior they want first, and write that issue.

## Bugs

A bug issue keeps the same template. "What we build" states expected versus
actual behavior and the reproduction. The first task is a failing test that
reproduces it; the fix is the second task.

## Sources

- Google, small changes: <https://google.github.io/eng-practices/review/developer/small-cls.html>
- GitHub, sub-issues: <https://docs.github.com/en/issues/tracking-your-work-with-issues/using-issues/adding-sub-issues>
- Wake, INVEST stories and SMART tasks: <https://xp123.com/articles/invest-in-good-stories-and-smart-tasks/>
- Beams, commit messages: <https://cbea.ms/git-commit/>
