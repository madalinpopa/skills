# Revsets

A revset is an expression that selects a set of revisions. Most commands take
one through `-r`, `--from`, `--to`, `--onto`, or a positional argument. Quote
any expression with spaces, parentheses, or `::` in the shell.

## Contents

- [Symbols](#symbols)
- [Operators](#operators)
- [Functions](#functions)
- [String patterns](#string-patterns)
- [Date patterns](#date-patterns)
- [Built-in aliases](#built-in-aliases)
- [Defining aliases](#defining-aliases)
- [Recipes](#recipes)

## Symbols

| Symbol | Meaning |
| --- | --- |
| `@` | Working-copy commit of the current workspace |
| `<name>@` | Working-copy commit of workspace `<name>` |
| `<name>@<remote>` | Remote-tracking bookmark or tag, for example `main@origin` |
| `main`, `v1.0` | A local bookmark or tag |
| `kntqzsqt` | A change ID or unique prefix |
| `d7439b06` | A commit ID or unique prefix |
| `<change>/<n>` | One version of a divergent or hidden change |
| `"x-"` | Quotes stop a symbol from being parsed as an expression |

Names resolve as tag, then bookmark, then commit or change ID. Use
`change_id(prefix)` or `commit_id(prefix)` to force one meaning.

## Operators

From strongest to weakest binding:

| Operator | Meaning |
| --- | --- |
| `f(x)` | Function call |
| `x-` | Parents of x |
| `x+` | Children of x |
| `p:x` | String or date pattern, or alias |
| `x::` | x and all descendants |
| `::x` | x and all ancestors |
| `x::y` | Descendants of x that are ancestors of y |
| `x..` | Everything not an ancestor of x |
| `..x` | Ancestors of x, excluding the root |
| `x..y` | Ancestors of y that are not ancestors of x |
| `::` / `..` | All visible commits (`..` excludes the root) |
| `~x` | Everything except x |
| `x & y` | Intersection |
| `x ~ y` | x minus y |
| `x \| y` | Union |

`@-` is the parent of the working copy, `@--` its grandparent, `@-+` the
siblings of the working copy including itself. A range `a..b` matches Git
semantics: commits on b's side after the fork point.

## Functions

Arguments in brackets are optional. `x` is any revset.

Navigation:

- `parents(x, [depth])`, `children(x, [depth])`
- `ancestors(x, [depth])`, `descendants(x, [depth])`: `::x` and `x::` with a
  depth limit
- `first_parent(x, [depth])`, `first_ancestors(x, [depth])`: follow only the
  first parent through merges
- `connected(x)`: `x::x`, fills gaps between selected commits
- `reachable(srcs, domain)`: commits in `domain` connected to `srcs`
- `heads(x)`: commits in x with no descendants in x
- `roots(x)`: commits in x with no ancestors in x
- `fork_point(x)`: the common ancestor where x's commits diverged
- `merge_point(x)`: the common descendant where x's commits joined

Whole sets:

- `all()`, `none()`, `root()`, `visible_heads()`, `working_copies()`
- `mine()`: author email matches `user.email`
- `empty()`: no file changes, including clean merges and the root
- `merges()`, `forks()`, `conflicts()`, `divergent()`, `signed()`

Identifiers and references:

- `change_id(prefix)`, `commit_id(prefix)`
- `bookmarks([pattern])`, `tags([pattern])`
- `remote_bookmarks([name], [[remote=]remote])`, `tracked_remote_bookmarks(...)`,
  `untracked_remote_bookmarks(...)`
- `remote_tags(...)`, `tracked_remote_tags(...)`, `untracked_remote_tags(...)`

Search:

- `description(pattern)`, `subject(pattern)`: full description or first line
- `author(pattern)`, `author_name(pattern)`, `author_email(pattern)`
- `committer(pattern)`, `committer_name(pattern)`, `committer_email(pattern)`
- `author_date(date_pattern)`, `committer_date(date_pattern)`
- `files(fileset)`: commits that touch matching paths
- `diff_lines(text, [fileset])`, `diff_lines_added(text, [fileset])`,
  `diff_lines_removed(text, [fileset])`: search changed lines

Utility:

- `latest(x, [count])`: newest by committer date, default 1
- `bisect(x)`: a commit near the middle of x
- `exactly(x, count)`: x, or an error if it does not have `count` commits
- `present(x)`: x, or `none()` instead of an error when a symbol is missing
- `coalesce(a, b, ...)`: the first non-empty set
- `at_operation(op, x)`: evaluate x as of an earlier operation

## String patterns

Functions such as `description()`, `bookmarks()`, and `author()` take a
pattern. Bare strings default to `glob:` in revsets.

| Pattern | Matches |
| --- | --- |
| `exact:"text"` | The whole string |
| `glob:"fix-*"` | Shell wildcard |
| `regex:"^feat\("` | Regular expression |
| `substring:"jpeg"` | Anywhere in the string |

Append `-i` for case-insensitive matching: `glob-i:"*JPEG*"`. Patterns can be
combined with `~`, `&`, `|`.

## Date patterns

`author_date()` and `committer_date()` take `after:"<date>"` (inclusive) or
`before:"<date>"` (exclusive). Accepted forms include `2024-02-01`,
`2024-02-01T12:00:00`, `2024-02-01 12:00:00`, `2 days ago`, `5 minutes ago`,
`yesterday`, `yesterday 5pm`.

## Built-in aliases

| Alias | Definition |
| --- | --- |
| `trunk()` | Head of the default bookmark (`main`, `master`, or `trunk`) on the `upstream` or `origin` remote |
| `immutable_heads()` | `builtin_immutable_heads()`, which is `trunk() \| tags() \| untracked_remote_bookmarks() \| untracked_remote_tags()` |
| `immutable()` | `::(immutable_heads() \| root())` |
| `mutable()` | `~immutable()` |
| `visible()` | `::visible_heads()` |
| `hidden()` | `~visible()` |
| `builtin_log()` | `present(@) \| ancestors(immutable_heads().., 2) \| trunk()`, the default `jj log` view |

Check the effective values with `jj config list --include-defaults
revset-aliases`.

## Defining aliases

```toml
[revset-aliases]
HEAD = '@-'
'stack()' = 'trunk()..@'
'user(x)' = 'author(x) | committer(x)'
'grep:x' = 'description(regex:x)'
'immutable_heads()' = 'builtin_immutable_heads() | remote_bookmarks(glob:"release/*")'
```

Set with `jj config set --repo 'revset-aliases."stack()"' 'trunk()..@'`.
Overriding `immutable_heads()` is how a project protects extra bookmarks.

## Recipes

```sh
jj log -r 'trunk()..@'                       # my current stack
jj log -r '(trunk()..@)::'                   # stack plus anything built on it
jj log -r 'remote_bookmarks()..'             # commits not on any remote
jj log -r 'mutable() & mine()'               # my unpushed work
jj log -r 'bookmarks() | tags()'             # named points only
jj log -r 'heads(mutable())'                 # tips of local work
jj log -r 'conflicts()'                      # what needs resolving
jj log -r 'divergent()'                      # duplicated change IDs
jj log -r 'description(glob:"fix*")'         # by subject
jj log -r 'files("src/**/*.go")'             # commits touching paths
jj log -r 'diff_lines(regex:"TODO")'         # commits adding or removing TODO
jj log -r 'author_date(after:"2 days ago")'
jj log -r 'latest(bookmarks("release/*"))'
jj log -r 'fork_point(main | feature)'       # where a branch left main
jj log -r 'roots(trunk()..feature)'          # first commit of the branch
jj rebase -s 'roots(trunk()..@)' -o main     # rebase the whole stack
jj log -r '@ | @- | @+'                      # current, parent, children
jj log -r 'present(main@upstream) | main@origin'
jj log -r 'at_operation(@-, @)'              # working copy before last op
```

When a symbol might not exist, wrap it in `present()` so the command still
runs. `exactly(x, 1)` guards commands that must not act on several revisions.
