# Refresh existing project instructions

On each user invocation of `create-project-docs`, compare the current skill's
assets and the current installed task skills with the project's reusable
instructions. Use `use-skills-cli` to check/update installed skill versions
first. The user invocation authorizes refreshing reusable guidance and routes;
do not require separate permission for each unambiguous instruction update.

Read existing files before editing. Merge narrowly: refresh stale generic
workflow rules and skill routes in AGENTS.md and documentation guides. Remove
duplicate generic instructions only after the replacement skill is available.
Preserve project-specific purpose, commands, paths, technologies, requirements,
decisions, permissions, custom rules, and all work records and question answers.
Keep Claude-specific additions. Never replace whole existing documents with assets.

Treat upstream instructions as reusable defaults, not authority to weaken a
project's stricter review or permission rules. Adapt links to existing owners
and locations. Do not migrate layouts, rename records, renumber features, or
rewrite specs, plans, tasks, mission, roadmap, stack, or decisions as maintenance.
Filling these documents requires separately requested scope.

If a passage mixes generic rules with project policy and its ownership is
unclear, leave that passage intact, report the specific conflict, and ask only
for the decision needed. Continue independent, unambiguous refreshes. Do not
infer that a difference from the asset is stale or that local edits are disposable.

Inspect the final diff against the original content. Verify refreshed routes,
available skills, and preserved project details. Report unchanged ambiguous
passages and any blocked skill update. A rerun with current instructions should
produce no further edits merely to rephrase equivalent guidance.
