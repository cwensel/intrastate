Model: claude-fable-5

# 0028 propose — prior-art read (Stage 2, selection depth)

Problem class: a declared, data-carried, in-process edit of ONE line of a text
file, driven by a planned value, with no shell and no wrapper script. Instance:
the RDR flow's `- **Status**:` bullet and its README index row.

## Queries and outcomes

| # | Tier | Query | Outcome |
| --- | --- | --- | --- |
| 1 | arc `StateMachineRes` | "apply a decided state change to a text file by declared line edit regex replace, no shell script wrapper" | 5 hits, all agent-skill prompt fragments (Phase 2: ANALYZE …). Rejected — no coverage. |
| 2 | arc `StateMachineRes` | "state machine tool persists the new state back into the source artifact file it read the state from" | 5 hits (awf-cli architecture, Archon plan-setup, scxmlcc manual). None persists state into a text artifact by declared edit. Rejected — no coverage. |
| 3 | source (raw GitHub) | ansible/ansible `lib/ansible/modules/lineinfile.py` | ACCEPTED — class prior art, quoted below. |
| 4 | source (raw GitHub) | puppetlabs-stdlib `lib/puppet/type/file_line.rb` | ACCEPTED — class prior art, quoted below. |
| 5 | semble `langref/kubebuilder` | "replace or insert a line in an existing file at a marker, in-process Go" | ACCEPTED — `pkg/plugin/util/util.go::InsertCode`, `::ReplaceInFile`, quoted below. |
| 6 | in-repo | `internal/cli/flowbind/flowbind.go::save`, `::Writer.Apply`; `internal/cli/flowbind/registry.go::commandBacked`; `internal/table/load.go::carrierDefect`; rdr `models/rdr-write.toml` header | ACCEPTED — the sibling signals the proposal reuses. |

⚠ no prior-art coverage in the state-machine corpus for the problem class;
the class prior art is configuration management (Ansible, Puppet) and Go
scaffolding tools (kubebuilder). Approaches rest on those plus the in-repo
carriers, not on the model prior alone.

## Accepted citations (verbatim)

### Ansible `lineinfile` (DOCUMENTATION block)

- `regexp`: "The regular expression to look for in every line of the file. For
  O(state=present), the pattern to replace if found. Only the last line found
  will be replaced."
- `line`: "The line to insert/replace into the file. Required for
  O(state=present). If O(backrefs) is set, may contain backreferences that will
  get expanded with the O(regexp) capture groups if the regexp matches."
- `backrefs`: "Used with O(state=present). If set, O(line) can contain
  backreferences (both positional and named) that will get populated if the
  O(regexp) matches."
- `state`: "Whether the line should be there or not."
- `insertafter`: "Used with O(state=present). If specified, the line will be
  inserted after the last match of specified regular expression."
- `create`: "Used with O(state=present). If specified, the file will be created
  if it does not already exist."
- write path: `tmpfd, tmpfile = tempfile.mkstemp(dir=module.tmpdir)` then
  `module.atomic_move(tmpfile, ...)`.

⇒ the carrier shape (regex anchor + replacement line with backrefs + atomic
temp-and-rename) is the settled class shape. ⇒ Ansible's multi-match
disposition is SILENT last-match; its no-match disposition is insert/append —
both rejected here (see Decision Rationale).

### Puppet `file_line` (type desc)

- `match`: "An optional ruby regular expression to run against existing lines
  in the file. If a match is found, we replace that line rather than adding a
  new line. A regex comparison is performed against the line value and if it
  does not match an exception will be raised."
- `multiple`: "An optional value to determine if match can change multiple
  lines. If set to false, an exception will be raised if more than one line
  matches"
- `append_on_no_match`: "If true, append line if match is not found. If false,
  do not append line if a match is not found"
- `replace`: "If true, replace line that matches. If false, do not write line if
  a match is found"

⇒ the peer that treats ambiguity as an ERROR (`multiple => false` raises) —
the cardinality disposition this proposal adopts as the only mode. ⇒ Puppet
also lints that the replacement line itself satisfies `match` — the analogue of
this proposal's "anchor still matches after the edit" read-back guarantee
(carried as a Pending assumption, not a clause).

### kubebuilder `pkg/plugin/util/util.go`

- `InsertCode(filename, target, code string) error` — `idx :=
  strings.Index(string(contents), target)`; `if idx == -1 { return
  fmt.Errorf("string %s not found in %s", …) }`; writes with `os.WriteFile`.
- `ReplaceInFile(path, oldValue, newValue string) error` — `if
  !strings.Contains(string(b), oldValue) { return errors.New("unable to find
  the content to be replaced") }`; `strings.ReplaceAll`; `os.WriteFile`.

⇒ a Go CLI that edits text in-process by literal anchor, refuses on no match,
and writes non-atomically (every occurrence replaced). Confirms the in-process
shape is ordinary Go; its non-atomic write and all-occurrences replace are the
two properties this proposal does NOT copy.

### In-repo sibling signals

- `internal/cli/flowbind/flowbind.go::save` — "writes the artifact atomically:
  the replacement is staged beside the target and renamed over it, so a reader
  never observes a half-written artifact." (`os.CreateTemp(filepath.Dir(path),
  ".flow-artifact-*")` … `os.Rename(staged, path)`).
- `internal/cli/flowbind/flowbind.go::Writer.Apply` — `if
  accessor.IsClear(t.Value) { delete(s, t.Key); continue }`.
- `internal/cli/flowbind/registry.go::commandBacked` — `return
  len(acc.Command) != 0 || acc.Path == ""` — the carrier discriminator; the
  residue builds a refusing binding.
- `internal/table/load.go::carrierDefect` — clause order "conflict, empty,
  unknown placeholder, shell interpreter, output shape, env conflict";
  `commandPlaceholders = []string{"{artifact}"}`.
- `internal/accessor/binding.go::WriteBinding` — "Apply performs the mutation.
  A `<clear>` planned value is a REMOVAL, not an assignment of the literal
  (`0004:C11`)."
- rdr `models/rdr-write.toml` header: "THE EDIT IS DATA — TRANSITIONALLY.
  `[rule.emit]` carries the exact `sed` expression … The end state is a
  successor intrastate RDR pair (a declared, in-process edit writer;
  `set-state --plan`) after which this file becomes a state-machine over a
  command reader and the `edit`/`sections` emits are deleted"; its lock row:
  `edit = "s/^- \\*\\*Status\\*\\*: Draft.*$/- **Status**: Final/"` and the
  joint-decision variant keeps the qualifier by sed backreference.

## Rejected branches

- Salt `file.replace` / Chef `line` cookbook: same class as Ansible/Puppet;
  not opened (budget — two class peers already fix the frame).
- Structured-document editors (yq, front-matter writers): document-family
  knowledge is fenced out by the kata; not opened.
