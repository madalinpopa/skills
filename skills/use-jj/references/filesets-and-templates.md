# Filesets and templates

Filesets select paths for commands such as `diff`, `squash`, `split`,
`restore`, `file list`, and `files()` in revsets. Templates control the text
that `log`, `show`, `evolog`, `op log`, and `bookmark list` print through
`-T`.

## Contents

- [Filesets](#filesets)
- [Templates](#templates)
- [Useful template calls](#useful-template-calls)

## Filesets

A bare path defaults to `prefix-glob:` relative to the current directory, so
`src` selects the whole directory and `src/main.go` one file. Quote any
expression with operators.

| Pattern | Matches |
| --- | --- |
| `cwd:"path"` | Path prefix relative to the current directory (the default) |
| `file:"path"` / `cwd-file:"path"` | Exactly that file |
| `glob:"*.rs"` / `cwd-glob:"..."` | Wildcard in the current directory |
| `prefix-glob:"..."` | Wildcard, and everything under matched directories |
| `root:"path"` | Path prefix from the workspace root |
| `root-file:"path"`, `root-glob:"..."`, `root-prefix-glob:"..."` | Root-relative forms |
| `<kind>-i:"..."` | Case-insensitive variant |

Operators, strongest first: `f(x)`, `~x` negation, `x & y` intersection,
`x ~ y` difference, `x | y` union. Functions: `all()`, `none()`.

```sh
jj diff '~Cargo.lock'                          # everything except one file
jj file list 'src ~ glob:"**/*.rs"'            # non-Rust files under src
jj split 'glob:"**/*_test.go"'                 # tests into the first commit
jj squash 'internal/config'                    # move one directory to parent
jj restore '~docs'                             # discard edits outside docs
jj log -r 'files(root:"cmd")'                  # commits that touched cmd/
```

Aliases live under `[fileset-aliases]` in config.

## Templates

Templates are expressions built from keywords, methods, literals, and
functions, joined with `++`. Keywords in `jj log` are the zero-argument
methods of the commit, so `commit_id` means `self.commit_id()`.

Operators from strongest to weakest: method call, `-x`/`!x`, `p:x`,
`* / %`, `+ -`, comparisons, `== !=`, `&&`, `||`, `++`.

Commit keywords and methods:

- `change_id`, `commit_id` with `.short([n])` and `.shortest([n])`
- `description`, `description.first_line()`, `trailers`
- `author`, `committer` with `.name()`, `.email()`, `.timestamp()`
- `bookmarks`, `local_bookmarks`, `remote_bookmarks`, `tags`, `working_copies`
- `parents`, `diff` with `.summary()`, `.stat()`, `.git()`, `.files()`
- Booleans: `empty`, `conflict`, `divergent`, `hidden`, `immutable`,
  `current_working_copy`, `mine`, `root`, `signature`

Other types:

- String: `.len()`, `.contains()`, `.starts_with()`, `.ends_with()`,
  `.substr()`, `.lines()`, `.split()`, `.replace()`, `.upper()`, `.lower()`,
  `.trim()`, `.first_line()`, `.escape_json()`
- Timestamp: `.ago()`, `.format("%Y-%m-%d")`, `.utc()`, `.local()`
- List: `.len()`, `.join(sep)`, `.filter()`, `.map()`, `.any()`, `.all()`,
  `.first()`, `.last()`, `.get(n)`, `.reverse()`, `.skip()`, `.take()`
- Operation (in `op log`): `id`, `description`, `time`, `user`, `tags`,
  `snapshot`, `root`, `current_operation`
- In `jj evolog` the commit is reached through `commit`, for example
  `commit.commit_id().short()`

Functions: `if(cond, then, [else])`, `coalesce(a, b, ...)`, `try(expr, ...)`,
`concat(...)`, `join(sep, ...)`, `separate(sep, ...)`, `surround(pre, post,
content)`, `indent(prefix, content)`, `fill(width, content)`, `pad_start`,
`pad_end`, `truncate_start`, `truncate_end`, `label(name, content)`,
`stringify()`, `json()`, `hyperlink()`, `config(name)`, `git_web_url()`.

Built-in log templates: `builtin_log_oneline`, `builtin_log_compact`,
`builtin_log_comfortable`, `builtin_log_detailed`, `builtin_log_compact_full_description`.
Set a default with `templates.log = "builtin_log_oneline"`. Define aliases
under `[template-aliases]`.

## Useful template calls

```sh
# one line per change, parseable
jj log --no-graph --color=never -T 'change_id.short() ++ "\t" ++ commit_id.short() ++ "\t" ++ description.first_line() ++ "\n"'

# the working-copy parent's change ID only
jj log --no-graph -r @- -T 'change_id ++ "\n"'

# bookmarks on each commit in the stack
jj log -r 'trunk()..@' -T 'change_id.short() ++ " " ++ bookmarks.join(", ") ++ "\n"'

# machine-readable JSON per commit
jj log --no-graph -T 'json(self) ++ "\n"' -r 'trunk()..@'

# description with fallback
jj show -T 'coalesce(description, "(no description set)\n")' -r @

# conflicted commits with their file count
jj log -r 'conflicts()' -T 'change_id.short() ++ " " ++ diff.files().len() ++ "\n"'
```

Use `--color=never` and `--no-graph` when another program parses the result.
