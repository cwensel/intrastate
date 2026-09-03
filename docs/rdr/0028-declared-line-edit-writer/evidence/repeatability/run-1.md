model: Claude Sonnet 5
variant: lite (profile: large)

Widened past: 0028:A1 (full text, for the `Artifact`/`WriteBinding.Apply`
signatures C1 assumes but does not restate); §illustrative-code (576-604, for
concrete TOML shape of the two field-level uses); §load-bearing-decisions
(553-574, for the carrier-discriminator and wire-format rationale backing
helper #2 and the data model). All C1.1-C1.6 clauses and the MVV were read in
full (not truncated). The S1-S31 proof list was read via the whole
`§testing-strategy` section rather than per-element selects (all fit in one
call); no additional widening was needed there since every clause it exercises
traces back to C1.

## 1. Public API

The RDR does not restate full Go signatures for new exported symbols; it
describes a seam extension to two existing types/interfaces. Reconstructed
signatures (GUESS marks fields/names the RDR does not spell exactly):

```go
// internal/accessor/model.go
type Artifact struct {
    Role string
    Path string
    Tags map[string]string // GUESS: exact field name/type; RDR only says
                            // "a context map on Artifact", "invocation-wide,
                            // not per-role" (A1)
}

// internal/accessor/binding.go — unchanged signature, newly reachable by a
// third binding kind
type WriteBinding interface {
    Apply(ctx context.Context, art Artifact, planned []resolve.Tag) error
}

type ReadBinding interface {
    Read(ctx context.Context, art Artifact, requested []string) (..., error) // GUESS: return shape
}
```

Table grammar (this IS normative, from C1.1-C1.6):

```
[write.<id>]
role, keys, timeout                    # unchanged from 0002/0004:C15

[write.<id>.edit.<key>]                # exactly one table per member of keys
anchor  = "<RE2 pattern>"              # must select exactly one line pre-edit
                                        # and exactly its own line post-edit
replace = "<template>"                 # whole replacement line, terminator excluded
clear   = "line"                       # optional; only legal value is "line"
```

Carrier rule: exactly one of `path` / `command` / `edit` per `[write.<id>]`
entry (0025:C1 extended to three). `edit` is refused on `read`/`gate` entries.

Substitution vocabulary:
- `replace`: `{<key>}` (this key only), `${N}`/`${name}` (anchor capture
  groups), `$$` → `$`, `{{`/`}}` → literal braces. Any other `{...}`/`$...`
  form, including bare `$1` without braces, is `edit_template_invalid`.
- `anchor`: ordinary RE2, plus `{tag.<key>}` for a declared `observed`-
  provenance tag, regexp-quoted at substitution (cannot alter pattern
  structure). All other braces pass through to RE2 untouched.
- `{tag.<key>}` also extends 0025:C2 as a whole-element command-argv
  placeholder on `read`/`gate`/`write` entries (C1.6).

Error modes (load-time, `intrastate lint`, C1.4):
- `edit_carrier_conflict` — `edit` beside `path`/`command`, or `edit` on a
  non-write entry
- `edit_key_mismatch` — `keys` member without matching `edit.<key>` table, or
  vice versa
- `edit_anchor_invalid` — anchor fails RE2 compile, names an undeclared tag
  key, or uses a malformed `{...}` form
- `edit_template_invalid` — `replace` has unknown placeholder, wrong key's
  placeholder, undefined group ref, malformed token
- `edit_clear_invalid` — `clear` outside `{"line"}`

Error modes (apply-time, execution, C1.2/C1.3/C1.5):
- `edit_anchor_unmatched` — anchor selects 0 lines (pre- or post-edit)
- `edit_anchor_ambiguous` — anchor selects ≥2 lines
- `edit_anchor_collision` — two rules select the same line
- `edit_anchor_unstable` — post-edit re-anchor doesn't land on the rule's own
  rewritten line
- `edit_value_multiline` — a planned or bound value contains `\n`/`\r`
- `edit_clear_undeclared` — planned `<clear>` but no `clear = "line"`
- plain `execution_failure` (0004) — unreadable target (OS error in Detail);
  unbound `{tag.<key>}` at invocation, before spawn; read-back mismatch or
  gate-off command reader
- All apply-time refusals occur BEFORE any byte is written; the target file is
  byte-identical to input on any refusal.

## 2. Three most important internal helpers

1. **Anchor resolver / line selector** (unnamed in RDR — GUESS at name, e.g.
   `resolveAnchor` or similar under a new `internal/accessor/edit*.go`).
   Responsibility: given the pre-edit buffer split into lines (on `\n`, `\r`
   kept with terminator), run each rule's anchor RE2 against every line and
   return the matching line index(es); enforce exactly-one-match per rule at
   selection time and again against the post-edit buffer (re-anchor,
   "IDENTITY not cardinality" per C1.3). Holds pre-edit indices so a
   preceding rule's deletion doesn't shift a sibling rule's target.

2. **Template parser/expander** (GUESS at name — "parse once" is C1.2's own
   term, not `regexp.Expand`, reused for both `anchor`'s `{tag.<key>}`
   placeholder and `replace`'s closed vocabulary). Responsibility: parse both
   templates once at LOAD into segments (literal | group | placeholder); at
   apply, emit bytes per segment without re-scanning captured/substituted
   text for further `${...}`/`{...}` forms — this is what makes runtime data
   inert as grammar.

3. **Carrier discriminator** — `registry.go::commandBacked`, extended from a
   two-way (`command`/residue) to a three-way check: `command` → command
   binding, `edit` → new edit binding, `path` → file binding, residue →
   refusing binding (unchanged 0025:C1 runtime arm). This is the single
   dispatch point deciding which `WriteBinding` implementation an entry gets;
   the RDR is explicit that this function's predicate is the extension site
   (Load-Bearing Decisions, "Selection / predicate").

(A fourth candidate the RDR names but treats as load-only, not runtime, is
`internal/table/load.go::carrierDefect`'s `case !hasPath && !hasCommand:`,
extended with `&& !hasEdit` — arguably a fourth "helper" but it's lint-time
validation, not part of the apply path, so ranked below the three above.)

## 3. Data model (persisted / cross-boundary)

Nothing new is persisted as a file format; the "data model" is the TOML table
schema plus the in-memory context carried across the accessor seam.

On-disk (TOML model file, `rdr-write.toml`-shaped):
```
[write.<id>.edit.<key>]
anchor:  string (RE2 source)
replace: string (template)
clear:   string, optional, closed set {"line"}
```

In-memory, crossing the accessor seam:
```go
Artifact{ Role, Path string; Tags map[string]string }  // GUESS on Tags field
                                                         // name/type (A1)
[]resolve.Tag  // "planned" — the values WriteBinding.Apply receives
```

Line buffer, in-memory only, never persisted as a distinct structure: file
bytes split on `\n`, each line optionally CRLF-terminated; final-terminator
absence is a property of the last line, preserved through rewrite. GUESS:
RDR does not name a Go type for this (e.g. `[]Line` or `[][]byte`) — likely
just `[][]byte` or a small unexported slice-of-struct, not specified.

Write mechanics (not really "data model" but the persistence contract): one
buffer, one staged file, one rename over the target (`flowbind.go::save`
discipline), target file mode preserved, symlinks resolved before staging so
the symlink itself survives. No lock file, no cross-entry transaction record.

## 4. Top-level pseudo-code (main operation: edit-carrier `Apply`)

```go
// WriteBinding.Apply for the "edit" carrier — GUESS at exact function name
// and file (something like internal/accessor/edit_binding.go), but the
// control flow below is drawn directly from C1.3 "execution semantics".
func (b *editBinding) Apply(ctx context.Context, art Artifact, planned []resolve.Tag) error {
    // 1. Read whole target file as bytes (symlinks resolved to real target)
    data, err := os.ReadFile(resolvedTargetPath(art))
    if err != nil {
        return execution_failure(err) // OS error in Detail; NOT edit_anchor_unmatched
    }

    // 2. Split into lines on "\n"; keep preceding "\r" with terminator;
    //    remember if final line had no terminator
    lines := splitPreservingTerminators(data)

    // 3. Validate planned/bound values before any selection
    for _, v := range plannedValuesAndBoundTags(planned, art.Tags) {
        if containsNewlineOrCR(v) {
            return edit_value_multiline // pre-mutation refusal
        }
    }

    // 4. Resolve each rule's anchor against PRE-EDIT lines; hold pre-edit
    //    indices (deletions of earlier rules must not shift later targets)
    selections := map[ruleID]int{}
    for _, rule := range b.rules {
        anchorRE := compileWithTagsBound(rule.Anchor, art.Tags) // regexp-quoted tag sub
        matches := anchorRE.FindAllLineIndices(lines)
        switch {
        case len(matches) == 0:
            return edit_anchor_unmatched
        case len(matches) >= 2:
            return edit_anchor_ambiguous
        }
        selections[rule.ID] = matches[0]
    }
    // detect cross-rule collisions
    if collision := findDuplicateTargets(selections); collision != nil {
        return edit_anchor_collision
    }

    // 5. Apply all rules to an in-memory copy of the buffer
    newLines := copyOf(lines)
    for _, rule := range b.rules {
        idx := selections[rule.ID]
        if rule.IsClear() {
            if rule.Clear != "line" {
                return edit_clear_undeclared
            }
            deleteLine(newLines, idx) // terminator included
            continue
        }
        newLines[idx] = expandTemplate(rule.Replace, capturedGroups, plannedValue(rule.Key))
    }

    // 6. Re-anchor pass over POST-EDIT buffer: each rule's anchor must
    //    select exactly its own (possibly shifted) line, or zero for a
    //    deleted line
    for _, rule := range b.rules {
        if !reanchorMatchesOwnLine(rule, newLines, selections) {
            return edit_anchor_unstable // before any write
        }
    }

    // 7. No-op check: if newLines == lines byte-for-byte, do nothing
    //    (no staging, no rename; inode unchanged)
    if bytesEqual(join(newLines), data) {
        return nil
    }

    // 8. Read-back precondition: if role's reader is command-backed and the
    //    command gate is off, refuse before mutation (naming the gate)
    if readerIsCommandBacked(art.Role) && !gateAllowsCommands() {
        return execution_failure("gate off, reader command-backed")
    }

    // 9. Stage beside target, preserve mode, rename over target
    //    (flowbind.go::save discipline)
    if err := stageAndRename(join(newLines), resolvedTargetPath(art), originalMode); err != nil {
        return err // rename failure ⇒ target untouched by construction
    }

    // 10. Read-back through role's declared reader (0004:C12/C13, unchanged)
    if !readBackConfirms(art.Role, planned) {
        return read_back_mismatch // Applied() left unset per 0004:C14 scoping
    }
    return nil
}
```

Steps 1-6 are all pre-mutation refusal checks (explicit RDR invariant: "every
refusal in this clause ... is decided BEFORE any byte is written"). Step 7 (no-op
skip) and step 9 (stage-and-rename, ONE buffer / ONE write per entry) are
directly from C1.3 `write:`. Step 8's ordering (read-back precondition before
staging) and step 10 (post-write read-back) reflect C1.3 `read-back:`, which the
RDR says is evaluated at the executor, not necessarily inline in `Apply` itself
— GUESS on whether step 8 lives inside this function or in the caller
(`Executor.Write`); the RDR states the check exists and must fire "before
`Apply`" is possible for the reader-resolution part, but leaves the precise call
site partly open pending A10 ("Pending" in the RDR text itself, not a
reconstruction gap).
