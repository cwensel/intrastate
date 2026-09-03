Model: claude-opus-5

# Persona 2 — Implementer

If I started coding this Monday, these are the questions I would have to ask before
writing the first line. Owned set: `0028:C1` (C1.1–C1.6) and the four `D-*` decisions.
I **widened** three times, and say where below; each widening was driven by a silence
in C1 that has no line range — an unstated call-site obligation is exactly the thing
a contract cannot anchor.

Everything below was checked against `/Users/cwensel/sandbox/newcoinc/intrastate` on
disk. Findings whose premise did not hold are not reported.

---

## HIGH

### 1. C1.3's gate pre-check has no site: `internal/accessor` cannot see `AllowCommands`

**Anchor:** `0028:C1` (C1.3 `read-back:`), `0028:A2`

C1.3 requires: "When that reader is command-backed and the gate is off, the write MUST
refuse BEFORE mutation (Detail naming the gate), because the reader is resolved before
`Apply` (A2)". A2 rates this "a check at an existing site, not a new ordering".

On disk, the ordering claim holds — `internal/accessor/executor.go:302` resolves
`reader, hasReader := e.Registry.readerFor(...)` before `binding.Apply` at line 353.
But the *state the check needs is not reachable from that site*:

- `internal/accessor/executor.go` imports only `context`, `errors`, `slices`, `time`,
  and `internal/resolve`. It has **zero** references to `cmdbind`, `flowbind`, or any
  carrier.
- `AllowCommands` is a field on `internal/cli/cmdbind/cmdbind.go:118 Config`, baked
  into the `cmdbind.Reader` struct value at
  `internal/cli/flowbind/registry.go:35` and held behind the opaque
  `accessor.ReadBinding` interface.
- The dependency runs the wrong way: `internal/cli/cmdbind/cmdbind.go:44` imports
  `internal/accessor`. `accessor` importing `cmdbind` back is an import cycle.
- The gate check itself lives at `internal/cli/cmdbind/cmdbind.go:170`
  (`if !cfg.AllowCommands`), inside `spawn` — A2's own Qualifier admits this.

`Definition.Accessor` (`internal/accessor/model.go:161`) *does* carry
`table.Accessor`, so "the reader is command-backed" is answerable at that site via
`len(acc.Command) != 0`. The gate half is not.

**Clarification request:** which of these is the contract?
(a) a new method on the `Binding` interface (e.g. `Gated() bool` / `PreflightRefusal() error`)
    that every binding must now implement — a seam change C1 does not mention and
    C1.6's "no other change to 0025:C1–C1.6" arguably forbids;
(b) hoisting the check to `internal/cli/flow_state.go` before `exec.Write` — A2's own
    named fallback, but then the refusal is minted by the CLI, not the accessor layer,
    and the "Detail naming the gate" has to be composed somewhere else;
(c) the executor testing `def.Accessor.Command` only and the gate state riding in on
    something else.

**Blocks:** the C1.3 `read-back:` clause; MVV step 5; `0028:S25`. Which of the three
is chosen decides whether this RDR is a one-package change or a `Binding`-interface
amendment.

---

### 2. A1's context map is a two-hop plumb, not one: the CLI drops `--tag` on the set-state path

**Anchor:** `0028:A1`, `0028:C1` (C1.2 `anchor admits:`, C1.6 `binding:`), `0028:MVV` step 2

A1 says the seam extension is "leaning a context map on `Artifact`", and verifies that
`accessor.Artifact` is `{Role, Path}` with no tag field. Both halves check out
(`internal/accessor/model.go:249`). But A1 stops one hop short of the call site.

**Widened** here from C1/A1 into `internal/cli/` because the assumption verifies the
*destination* struct and is silent on the *producer*.

- `internal/cli/flow_exec.go:340-345` — `flowRequest.artifactMap()` builds
  `accessor.Artifact{Role: role, Path: path}` from `r.artifacts` **only**.
- `flowRequest` (`internal/cli/flow_exec.go:23-29`) holds the parsed tags in a
  *separate* field, `observed []resolveTag`.
- `runFlowSetState` (`internal/cli/flow_state.go:324-399`) calls
  `accessor.NewExecutor(req.registry, req.artifactMap())` at line 364 and **never
  references `req.observed` again**. Its own comment at line 331-333 says the parsed
  values "are context only and never become writes (REQ-34)" — today they are parsed
  purely so their refusals bite, then discarded.
- `accessor.Artifact` is the *sole* non-test construction site in the repo
  (`rg 'accessor\.Artifact{' internal/` → one non-test hit).

So the work is: add the field to `Artifact`, **and** change `artifactMap()`'s
signature/body to take `observed`, **and** decide whether the same map rides the
`resolve` and `read-state` paths (which also call `artifactMap()`).

**Clarification request:** does the context map ride on `Artifact` (one map copied into
every role's handle, so every role sees every tag), or on `Artifacts`/the `Executor`
(one map for the invocation)? A1 picks "on `Artifact`" but does not say whether the map
is per-role or invocation-wide, and `artifactMap()` iterates roles — the per-role copy
is the cheap read of A1 and it is almost certainly not what is wanted.

**Blocks:** the shape of the A1 seam extension, which C1.2 and C1.6 both stand on; the
first commit of Phase 2.

---

### 3. C1.2's "any provenance" is refuted by `parseTags` for two of the three provenances

**Anchor:** `0028:C1` (C1.2 `anchor admits:`), `0028:C1` (C1.6 `admission:`), `0028:A7`

C1.2 admits `{tag.<key>}` for "a tag key the model DECLARES (**any provenance**)".
C1.6 repeats "a tag key the model DECLARES". A7 verifies `--tag` is the context
channel. But A7's own evidence names the constraint C1.2 contradicts, and the code is
harder still:

`internal/cli/flow_input.go:606-647 parseTags` refuses, before anything else:
- `key == table.RecognizedTagKey` (`"recognized"`, `internal/table/model.go:19`) →
  `codeTagReserved`, "the recognized outcome enters through --outcome";
- `slices.Contains(owned, key)` → `codeTagOwned`, "owned state is read from the
  declared read accessors, never supplied by the caller".

So an **owned** tag key and the **recognized** key can never be bound through `--tag`.
Only `observed` provenance (and undeclared keys, which are admitted with a zero
`TagDecl`) can. C1.2's "any provenance" is therefore either wrong, or it means "any
provenance the loader will lint-accept in the anchor, of which only one is bindable at
runtime" — in which case every non-observed declaration is a guaranteed
unbound-at-invocation `execution_failure`, and C1.4's `edit_anchor_invalid` should
arguably catch it at lint instead.

**Clarification request:** is `edit_anchor_invalid` / `command_unknown_placeholder`
supposed to fire at lint for a `{tag.<key>}` naming an *owned* or the *recognized* key
— keys that are declared but structurally unbindable? Or is C1.2's "any provenance"
meant to be narrowed to `observed`?

**Blocks:** C1.4's `edit_anchor_invalid` predicate ("names an undeclared tag key" —
is an unbindable-but-declared key undeclared for this purpose?); `0028:S3`; `0028:S7`.

---

## MEDIUM

### 4. `{{` / `}}` as an escape inside an RE2 anchor collides with RE2's repetition quantifier

**Anchor:** `0028:C1` (C1.2 `escapes:`, `vocabulary:`), `0028:S14`

C1.2's `escapes:` says "`{{` and `}}` emit literal braces (**both templates**)", and
`vocabulary:` says "Any other `{…}` or `$…` form is `edit_template_invalid` /
`edit_anchor_invalid`". `0028:S14` restates "in both templates".

In an anchor, `{` and `}` are RE2 metacharacters: `\d{4}`, `[0-9]{4}`, `a{2,3}` are
ordinary, and for this consumer an RDR-number anchor (`\d{4}`) is the obvious
spelling. Under C1.2 as written, `\d{4}` is "some other `{…}` form" and refuses
`edit_anchor_invalid`; writing `\d{{4}}` would emit two literal braces into the pattern,
not a quantifier. The two anchors in `0028:§illustrative-code` happen to avoid `{n}`,
so the collision is not visible there.

**Clarification request:** in the anchor template, is the placeholder scanner restricted
to exactly the two literal forms `{tag.<key>}` (matched by a fixed prefix `{tag.`) with
every other `{` passed through to RE2 untouched — or does the closed-vocabulary rule
really reject `\d{4}`? If the former, C1.2's "both templates" applies to `{{`/`}}` in
`replace` only, and the anchor needs its own, narrower rule.

**Blocks:** the anchor parser's tokenizer, written at load time (C1.2 `parse once:`);
`0028:S14`'s expected behaviour; C1.4's `edit_anchor_invalid` predicate.

---

### 5. C1.3 names no disposition for a target file that cannot be read

**Anchor:** `0028:C1` (C1.3 `input:`, `order:`), and the `disposition` table in
`0028:§technical-design`

C1.3's `input:` opens "the whole file is read as bytes" and says nothing about the read
failing. The `disposition` table in §technical-design enumerates every input class and
has no row for it. This is a silence, so it has no line range of its own.

The classes are concrete: the target does not exist (ENOENT), is a directory (EISDIR),
or is unreadable (EACCES). The nearest precedent goes the *opposite* way —
`internal/cli/flowbind/flowbind.go:110-114 load` treats "a file that does not exist yet
[as] an EMPTY artifact, not an error", explicitly so the first write of a flow is
possible. Applying that precedent to `edit` would give: empty buffer → anchor selects 0
→ `edit_anchor_unmatched`, which is at least *a* disposition — but it is not the one an
implementer would guess, and it silently reuses "unmatched" for "no file".

**Clarification request:** does a missing target refuse `edit_anchor_unmatched` (empty
buffer, `flowbind.load`'s precedent), or a plain `execution_failure` with the OS error
in Detail? And is an EACCES/EISDIR read error the same class? C1.3 `order:` promises
every refusal is decided before any byte is written, which holds either way, but the
Detail token differs.

**Blocks:** the `Apply` entry path; the disposition table's completeness claim
("Every refusal … adds no new class"); a negative test that does not yet exist in the
`0028:S*` list.

---

### 6. An in-process `edit` can still be minted applied-but-unverified by the executor's deadline arm

**Anchor:** `0028:C1` (C1.3 `order:`, `no subprocess:`), `0028:C1` (C1.1 `timeout`)

C1.3 `order:` states flatly: "no pre-write refusal ever carries 0004:C14's
applied-but-unverified sense". C1.1 keeps `timeout` "unchanged".

**Widened** into `internal/accessor/executor.go` because the claim is about the
executor's behaviour, not the binding's, and C1 cannot see it.

`internal/accessor/executor.go:352-364`:

```go
applyCtx, cancelApply := context.WithTimeout(ctx, timeout)
err := binding.Apply(applyCtx, art, slices.Clone(planned))
appliedDeadline := errors.Is(applyCtx.Err(), context.DeadlineExceeded)
cancelApply()

if appliedDeadline {
    r := refusalOf(def, timeout, ClassTimeout, nil)
    r.applied = true
    ...
}
```

The deadline arm is evaluated *after* `Apply` returns and does not consult `err`, so it
fires whenever the wall clock passed the timeout during `Apply` — regardless of whether
the binding ever touched the file. For a large file, a pathological RE2 anchor, or a
loaded machine, an `edit` that refused `edit_anchor_unmatched` before writing a byte is
reported as `ClassTimeout` with `Applied() == true`: "may have been applied and was not
verified", for a write that provably was not. `flowbind.Writer.Apply` ignores its ctx
entirely (`func (w *Writer) Apply(_ context.Context, …)`,
`internal/cli/flowbind/flowbind.go:250`), so the hazard is latent for `path` today —
but C1.3 is the first contract to make the universal claim the arm can falsify.

**Clarification request:** is C1.3 `order:`'s claim scoped to refusals the *binding*
mints (in which case say so, because the executor mints one it cannot), or does the
executor's deadline arm need an amendment for a carrier that spawns nothing? Related:
what timeout value is even meaningful for an in-process edit, and is `timeout` still
required on an `edit` entry?

**Blocks:** the C1.3 `order:` invariant as a testable statement; `0028:S23`'s
neighbourhood; whether this RDR touches `internal/accessor/executor.go` at all beyond
finding 1.

---

## LOW

### 7. C1.1's "the 'neither' arm is unchanged" reads against `0028:S1`

**Anchor:** `0028:C1` (C1.1 `carrier:`), `0028:S1`

C1.1: "The 'neither' arm is unchanged and still reports 0025:C5's
`command_and_path_conflict`". `0028:S1` Expected: "An entry carrying only `edit` LOADS
(0025:C5's 'neither' arm no longer fires for it)."

Both are satisfiable — "neither" must now mean "none of the three" — but the word
"unchanged" is doing damage. On disk, `internal/table/load.go:1058-1079 carrierDefect`
computes `hasPath`/`hasCommand` and the arm is
`case !hasPath && !hasCommand:`; its predicate **must** gain `&& !hasEdit`. Its message
is a hardcoded string: `" declares neither `path` nor `command`; an entry carries
exactly one carrier"`. That message becomes false in a three-carrier world.

Note also the "both" arm one case above: `case a.Path != nil && a.Command != nil` — it
keys on the pointer being **non-nil**, not on `hasPath`, so `path = ""` beside a
`command` is already a conflict. C1.4's `edit_carrier_conflict` needs the same
pointer-not-emptiness discipline to be consistent, and C1.1 does not say so.

**Clarification request:** does `edit_carrier_conflict` key on `a.Edit != nil` (pointer
presence, matching the existing "both" arm) or on a non-empty edit table? And is the
`command_and_path_conflict` message text expected to be updated to name all three
carriers (C1.4 says "the wire strings are the contract, the identifiers are not" — it
is silent on message text)?

**Blocks:** the `carrierDefect` edit; whether an existing `0025` test asserting that
message string breaks.

### 8. `0028:§illustrative-code` will not parse as TOML

**Anchor:** `0028:§illustrative-code`

The snippet is marked "Illustrative … field spellings are C1's", so an implementer
will lift it into the first fixture. Two lines will not load:

- `role = "record"; keys = ["status"]; timeout = 5` — TOML has no `;` statement
  separator; each key/value needs its own line.
- `timeout = 5` — `sourceAcc.Timeout` is `*string`
  (`internal/table/source.go:88`) and `Definition.timeout()` parses it with
  `time.ParseDuration` (`internal/accessor/model.go:187`). Every shipped fixture spells
  it `timeout = "5s"` (e.g. `models/examples/review-state-machine.toml:58`). An integer
  refuses at strict decode.

**Clarification request:** none needed on intent — just confirming these are cosmetic
and the fixture should spell `timeout = "5s"` on its own line.

**Blocks:** nothing normative; costs the first hour if copied verbatim.

### 9. `edit` has no analogue for `Writer`'s unreachable-locator seal

**Anchor:** `0028:C1` (C1.3 `target:`), `0028:§existing-infrastructure-audit`
(row "File write binding")

C1.3 `target:` says "the model names no path" for an `edit` entry. The audit says
"Reuse pattern; new binding" against `flowbind.go::Writer`. But `Writer` carries a
declared-locator `Path` field whose *suffix* drives an unreachable-artifact seal:
`unreachable(path)` (`internal/cli/flowbind/flowbind.go:99-102`) tests for a
`.<suffix>`/`-<suffix>` marker, and `Apply` writes/clears `sealedKey`
(`flowbind.go:264-281`) to make the applied-but-unverified sense atomic with the
mutation. An `edit` binding has no declared path, so no seal, and no place to put one
(the target is a markdown file, not a key/value store).

**Clarification request:** confirm the `edit` binding has *no* unreachable-locator
affordance and that this is intentional (the seal is a `path`-carrier test fixture, not
a contract obligation) — so the reuse of `flowbind.go::save`'s discipline stops at
stage-and-rename and does not drag the seal along.

**Blocks:** the new binding's struct shape; whether `0004:C14`'s applied-but-unverified
sense has any `edit`-side representation.

### 10. C1.2 cites `regexp.Expand` semantics for a hand-rolled parser

**Anchor:** `0028:C1` (C1.2 `parse once:`, `escapes:`)

C1.2 mandates a load-time parse into segments and states "captured text and substituted
values are never re-scanned" — which rules out Go's `regexp.Expand`, since `Expand`
scans its template at expansion time and would happily expand a `$1` that arrived in a
captured value. Yet the same clause cites "Go `regexp.Expand` semantics" as the
authority for the empty-group rule, and `0028:S14` cites it again.

Two concrete divergences an implementer must resolve: `Expand` accepts bare `$1` and
`$name` (no braces) as references, whereas C1.2's closed vocabulary lists only
`${N}`/`${name}`; and `Expand` treats an unknown `$x` as an empty expansion rather than
an error, whereas C1.4's `edit_template_invalid` refuses it at lint. There is no
non-test `regexp` import under `internal/` today (verified — `D-selection-predicate`'s
claim holds), so this is greenfield.

**Clarification request:** is bare `$1` a `edit_template_invalid` (closed vocabulary,
braces required) or an accepted synonym (Expand parity)? The citation should be read as
"the empty-group rule matches `Expand`'s", not "use `Expand`" — confirming that keeps
the parse-once rule intact.

**Blocks:** the template tokenizer; `0028:S14`; `0028:S4`'s "malformed `${…}`" arm.

---

## Checked and clean

These I checked because a finding would have turned on them, and the premise held —
recording so the next pass does not re-check:

- `0028:A2`'s ordering claim: `readerFor` at `executor.go:302` precedes
  `binding.Apply` at `executor.go:353`. True.
- `0028:A1`'s struct claims: `WriteBinding.Apply(ctx, art, planned)` and
  `ReadBinding.Read(ctx, art, requested)` take the same `art`;
  `Artifact` is `{Role, Path}`. All true (`internal/accessor/binding.go`,
  `model.go:249`).
- `0028:A6`: `carrierDefect` and `commandBacked` are the only two carrier-shaped
  sites; `commandBacked` really is `return len(acc.Command) != 0 || acc.Path == ""`
  (`registry.go:100-102`); `checkAccessorBindings` is carrier-agnostic. True.
- `0028:A9`: `cmdbind.go:400-425 substitute` compares whole elements
  (`if el != ArtifactPlaceholder`). True — C1.6 is one more branch in that loop.
- `0028:D-selection-predicate`'s "no non-test `regexp` import under `internal/`" —
  verified, zero hits.
- `0028:C1` C1.4's registration claim: `table.Categories()`
  (`internal/table/category.go:66-106`) ends with exactly 0025's six, tail-appended
  with a comment saying so. Appending five more there is mechanical.
- `0028:C1` C1.3's `Invocations()` claim: `flowbind.Writer.Apply` increments at
  function entry (`flowbind.go:251`), so a refused edit still counts as one `Apply`.
  Consistent with "counts `Apply` calls exactly as `flowbind.go::Writer` does".
