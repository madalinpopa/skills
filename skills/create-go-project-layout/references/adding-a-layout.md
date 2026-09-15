# Adding a layout

Add three things and the scaffold script picks the layout up by name:

- `assets/<layout>/` with the project tree and a `.scaffold` manifest. The
  manifest holds `go_dir`, one `require` line per library, and one `tool` line
  per generator. The script deletes it from the target after copying.
- `references/layouts/<layout>.md` with a fit summary, a table of contents,
  and the section order of the existing layout reference.
- A row in the `SKILL.md` layout table with its fit and exclusions, plus a
  direct reference link and any prerequisites needed for routine scaffolding.

Paths above are relative to the skill root. Use only the three tokens the
script supports: `{{PROJECT_NAME}}`, `{{MODULE_PATH}}`, and `{{GO_VERSION}}`.
A layout that needs another token also needs a script change.
