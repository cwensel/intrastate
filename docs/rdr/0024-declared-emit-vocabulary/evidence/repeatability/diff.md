model: claude-opus-5[1m]

# Repeatability full-diff — RDR 0024 (Declared emit vocabulary)

Post-barrier DIFF pass. Inputs: `run-1.md` (claude-opus-5[1m]),
`run-2.md` (claude-sonnet-5), `run-3.md` (claude-haiku-4-5). All three
ran `variant: full (profile: foundational)`, so three runs is the
intended coverage. Run-3 is the alt-model leg; a split that isolates it
is flagged **ALT-MODEL SPLIT**.

Findings are ordered by how many runs disagree (3-way first), then by
whether the element is a normative contract.

---

## 1. Disagreements

### D-1 — `loadEmitDecls` / `checkRuleEmit` signatures (3-WAY) → **0024:C2**

The single largest divergence. All three runs name both functions
correctly and place them correctly, and all three render a *different*
signature:

| Run | Rendering |
| --- | --- |
| run-1 | `loadEmitDecls(l *loader) error` / `checkRuleEmit(l *loader) error` — free function taking the loader; explicitly marked GUESS |
| run-2 | `loadEmitDecls(src *sourceDoc, m *Model) error` / `checkRuleEmit(src *sourceDoc, m *Model) error` — free function, two args, and it invents the type name `sourceDoc` |
| run-3 | `(l *loader) loadEmitDecls() error` / `(l *loader) checkRuleEmit() error` — method on `*loader` |

Run-3 (alt-model) is the only one that matches HEAD
(`internal/table/load.go:206` is `func (l *loader) loadTags() error`).
Run-2's `*sourceDoc` parameter type does not exist in the repo at all.

**RDR passage that let them diverge**: C2 fixes the names, the order,
and the insertion point ("inserted immediately after `loadTags` in
`load.go::run`'s step slice") but never writes a signature, and never
says the step slice holds methods on `*loader`. C2 says the steps run
"beside `loadTags`" without saying they are *shaped* like `loadTags`.
Run-1 and run-2 both reasoned from "step slice" to a free function
because nothing forbade it. Only run-3 read the method form off the
neighbour.

**Lands on**: `0024:C2`, Normative Contracts. One clause pinning the
receiver form — e.g. "both are methods on `*loader`, the shape
`loadTags` already has" — removes this entirely.

### D-2 — the dispositions-join helper: named, or inline (3-WAY) → **0024:C4**

| Run | Rendering |
| --- | --- |
| run-1 | GUESS: inline in `resolvePayload`, or an unexported `joinDispositions(decls map[string]EmitDecl, emit map[string]string) map[string]string` |
| run-2 | GUESS: inline inside `resolvePayload`'s construction, explicitly declines to name a function |
| run-3 | A named, **exported-receiver** method: `func (m *Model) dispositionsFor(row *Row) map[string]string`, placed in `internal/table` (not `internal/cli`), marked GUESS only as to its caller |

**ALT-MODEL SPLIT** on the *package boundary*, not just factoring:
run-1 and run-2 both put the join in `internal/cli/flow_resolve.go`;
run-3 puts it on `*Model` in `internal/table` and has
`internal/cli` merely call it. That is a real architectural fork — it
decides whether `internal/table` grows a payload-shaped API.

**RDR passage**: C4 fixes the join's *semantics* exhaustively (row-keyed
off the selected row's authored emit, one entry iff the row authors `k`
and `k`'s declaration lists that value under a disposition, never padded)
and fixes the field's *position* ("inserted at one fixed position in
`internal/cli/flow_resolve.go::resolvePayload`, immediately after
`emit`"). It never says where the join code lives. C3's "This RDR adds
exactly one reader of the carried declarations — C4's payload join"
identifies the reader but not its home package.

**Lands on**: `0024:C4`. One clause: the join is performed in
`internal/cli`, reading `Model.EmitDecls`; `internal/table` gains no
payload-shaped helper.

### D-3 — source-line attribution / `emitHeaderLine` (3-WAY, by omission) → **0024:C1** and **0024:C2**

| Run | Rendering |
| --- | --- |
| run-1 | Reconstructs `emitHeaderLine(src, key string) int` as one of its three key helpers, plus a GUESSed second locator `emitBlockLine` for the rule-side keying, and states all three categories are stamped through `atLine` |
| run-2 | Silent — no locator function, no `atLine`, no source-line behavior anywhere in the reconstruction |
| run-3 | Silent — no locator function; lists `*Failure` with a `Category` field as the refusal type but no line stamping |

Run-1 got there only by widening into `§phase-1-...`, which it says
explicitly ("Without this the error-mode reconstruction would have
invented locator plumbing"). Runs 2 and 3 read the contracts and the
locator vanished.

**RDR passage**: this is contract *silence*. `emitHeaderLine`, `atLine`,
and the whole line-attribution requirement live only in Phase 1 and in
the `trace` mini-check. C1 and C2 — which define all three refusals —
say nothing about a locator, yet the Failure Modes section makes
"offending file in `locator`" a promise and the MVV asserts "each
carrying the offending block's SOURCE LINE, not `:1`". A reader working
from the contracts alone cannot recover a normative requirement that the
Failure Modes and MVV both depend on.

**Lands on**: `0024:C2` (the refusal contract) — one clause requiring
all three categories to carry a source line via a `tagHeaderLine`
sibling. Secondarily `0024:C1` for the declaration-header keying.

### D-4 — the rule-side locator: one function or two (2-WAY vs. silence) → **0024:C2**

run-1 alone reaches this and marks it GUESS: Phase 1 names *one*
function, `emitHeaderLine`, but assigns it *two* keying strategies (the
`[emit.<key>]` header for declaration defects; "the offending rule's
`[rule.emit]` block" for `unknown_emit_key` / `emit_value_out_of_domain`).
run-1 reconstructs either a discriminating argument or an unnamed second
helper `emitBlockLine`, and says the RDR does not choose. Runs 2 and 3
never reach the question (see D-3).

**Lands on**: `0024:C2`. Subsumed by D-3's rewrite if that rewrite names
both keying strategies.

### D-5 — where the zero-declaration opt-in gate is evaluated (3-WAY) → **0024:C2**

| Run | Rendering |
| --- | --- |
| run-1 | Explicit guard *inside* `checkRuleEmit` (`if len(l.model.EmitDecls) == 0: return nil`), and explicitly GUESSes that C2 fixes only the observable, so guard-vs-empty-loop is free |
| run-2 | Gate hoisted to the **caller**: `if len(m.EmitDecls) > 0: err = checkRuleEmit(src, m)` in `Load` — `checkRuleEmit` is not even called on a zero-declaration model |
| run-3 | Explicit guard inside `checkRuleEmit`, step 1 — and *additionally* an early return inside `loadEmitDecls` ("If zero declarations, populate `l.model.EmitDecls` with an empty map and return nil") |

Three placements: callee, caller, and both. Run-2's caller-side gate is
the odd one — it makes the opt-in predicate a property of `load.go::run`
rather than of the step, which conflicts with C2's framing of two steps
inserted into a step slice (a slice of uniform steps has no room for a
conditional call).

**RDR passage**: C2 fixes only the observable ("With ZERO `[emit.*]`
declarations the load pipeline is byte-for-byte today's: no new refusal
is reachable") and never says where the predicate is evaluated. The
"step slice" framing implies uniform steps but does not state it.

**Lands on**: `0024:C2`.

### D-6 — `EmitDecl.Domain` for non-enum kinds: `nil` vs `nil/empty` (3-WAY, and one contradiction) → **0024:C3**

| Run | Rendering |
| --- | --- |
| run-1 | `Domain []string // enum only; bytewise-sorted union; nil otherwise` |
| run-2 | `Domain []string // bytewise-sorted union of all members; nil/empty for non-enum kinds` — leaves the two indistinct |
| run-3 | `Domain []string // sorted, union of all disposition partitions (enum only)` — no statement of the non-enum value |

This one is load-bearing, not cosmetic: C3 states the reuse is
"safe-by-omission" and that `conformDomain`'s enum arm is guarded by
`len(decl.Domain) > 0`. A `nil` and an empty non-nil `Domain` behave
identically *there*, but scenario 4's value-equality assertion over the
carrier will distinguish them, and run-2's "nil/empty" is a genuine
under-determination that a golden comparison would catch.

**RDR passage**: C3 names the three fields and fixes `Domain` as "the
union sorted bytewise" for the enum case. It is silent on the non-enum
case entirely.

**Lands on**: `0024:C3`.

### D-7 — `Model.EmitDecls` nil vs. empty map for a zero-declaration model (3-WAY) → **0024:C3**

All three flag this and all three resolve it differently in emphasis:

| Run | Rendering |
| --- | --- |
| run-1 | GUESS, resolved as "always non-nil, possibly empty", mirroring `Model.Tags` |
| run-2 | GUESS, "empty (not nil, by GUESS matching Tags convention)" |
| run-3 | Not marked GUESS — asserted as fact: `loadEmitDecls` "populate `l.model.EmitDecls` with an empty map" |

**ALT-MODEL SPLIT on confidence**, not on answer: all three land on the
empty map, but run-3 states it normatively where the other two flag it
as unresolved. That is the more dangerous split — a reader taking run-3's
reading writes no test for it.

**RDR passage**: C4 fixes `{}`-not-`null` for the *payload*, which is a
different object. C3 says nothing about the carrier's nil-ness, and the
`Model.Tags` sibling analogy is stated for the *keying*, not the
zero-value.

**Lands on**: `0024:C3`. See also GUESS cluster G-1.

### D-8 — the `[emit]` source-schema field type (2-WAY + 1 divergent) → **0024:C1**

| Run | Rendering |
| --- | --- |
| run-1 | Concrete Go: `Emit map[string]sourceEmitDecl` with `sourceEmitDecl{Kind string; Domain any}`; GUESSes the struct name and the `any` spelling vs `interface{}`/`toml.Primitive` |
| run-2 | Prose only: "the `domain` field is the one place strict decoding cannot descend further (it is `any`-typed)" — no struct rendered |
| run-3 | Reads it backwards: "**GUESS:** The source schema for `[emit]` is decoded with `DisallowUnknownFields: false` at the domain subtree level (A7)" |

**ALT-MODEL SPLIT, and run-3's reading is wrong in mechanism.** C1 says
`DisallowUnknownFields` *stops descending* at `domain` because the field
cannot be concretely typed — that is a consequence of the field's Go
type, not a decoder option toggled per-subtree. `pelletier/go-toml/v2`
has no per-subtree strictness switch. Run-3 reconstructed a decoder
configuration that does not exist.

**RDR passage**: C1 describes the *effect* ("stops descending at it")
and cites `sourceModel.Metadata` as the precedent carve-out, but never
writes the field's Go type. A7 likewise phrases it as behavior. Nothing
in the record says "the field is `any`-typed" in normative voice — run-2
inferred it, run-1 committed to it as a GUESS, run-3 mis-mechanized it.

**Lands on**: `0024:C1`. One clause naming the carve-out as a field type
(`any`, mirroring `Metadata`'s free-form map), not a decoder setting.

### D-9 — "both domain forms present" as a refusal arm (ALT-MODEL, run-3 only) → **0024:C1**

run-3 adds an arm C1 does not have: "**Reject if both are present** or
if `kind != "enum"` carries any domain." C1's arm list has ten entries
and this is not among them; C1 says the domain comes "in exactly one of
two forms" but the two forms are `domain = [...]` and
`[emit.<key>.domain]`, which are the *same TOML key* — authoring both is
a TOML-level duplicate-key error from the decoder, not a hand-written
arm. run-1 and run-2 both reproduce C1's list without this arm.

**ALT-MODEL SPLIT.** run-3 invented a reachable-looking refusal that is
unreachable, exactly the failure C1's "'No usable domain' is ONE arm,
not two" paragraph exists to prevent for the adjacent case.

**RDR passage**: C1's "in exactly one of two forms" reads as a
mutual-exclusion *constraint an implementer must enforce* rather than as
a description of the grammar's shape. The neighbouring paragraph does
the disambiguating work for the empty-domain case but not for this one.

**Lands on**: `0024:C1`.

### D-10 — `scalar` short-circuit in the value check (2-WAY split) → **0024:C1** / **0024:C3**

| Run | Rendering |
| --- | --- |
| run-1 | Explicit early `continue` in pseudo-code: `if d.Kind == "scalar": continue // never refused` |
| run-2 | Prose only ("`scalar` never refused"); pseudo-code passes everything to `ConformValue` |
| run-3 | No mention in the check at all; passes everything to `ConformValue` |

Behaviourally equivalent at HEAD — C3 documents that `conformKind`
carries only `int` and `bool` arms so `scalar` falls through to `nil` —
but run-1's explicit guard and runs-2/3's reliance on fall-through are
opposite postures toward C3's own warning that the reuse is
"safe-by-omission" and "if C1's arms are ever relaxed,
`emit_value_out_of_domain` silently stops firing".

**RDR passage**: C1 says `scalar` values "are never refused"; C3 says
the fall-through delivers that. The record does not say whether the
implementation should *depend* on the fall-through or guard explicitly,
despite spending a paragraph on the fragility of depending on it.

**Lands on**: `0024:C3` (the safe-by-omission clause).

### D-11 — the `dispositions` map's key ordering mechanism (2-WAY vs 1) → **0024:C4**

| Run | Rendering |
| --- | --- |
| run-1 | "marshals byte-ordered" — relies on `encoding/json`'s map-key sort, no explicit step |
| run-2 | Silent on mechanism; asserts the wire shape only |
| run-3 | Adds an explicit pipeline step: "Sort result keys bytewise and return" (twice — in helper step 3 and in the resolve pseudo-code step 4.b) |

**ALT-MODEL SPLIT.** run-3 reconstructs a sort over a Go `map`, which is
not a meaningful operation on the returned value — the ordering C4
requires is delivered by `encoding/json`, which sorts map keys on
marshal. run-3's step is either dead code or implies an ordered carrier
type the RDR does not have.

**RDR passage**: C4 says "keys in byte order" without saying that this
is the marshaller's doing rather than the join's. C3's explicit
"`Domain` is the union sorted bytewise" — a real sort on a real slice —
sits nearby and primes the reader to expect an explicit sort here too.

**Lands on**: `0024:C4`.

### D-12 — wire-key count arithmetic (1-WAY error) → **0024:§testing-strategy** (S5)

run-2 writes "14 fields on the wire (up from 13 pre-change),
`dispositions` inserted at index 9". S5 says the pre-change count is
thirteen and `dispositions` is inserted "at index 9" — with a 14-key
list whose 10th entry (index 9, zero-based) is `dispositions`. run-2's
arithmetic is right; run-1 reproduces the same numbers ("15 struct
fields; 14 wire keys"); run-3 does not reach the numbers at all,
reporting S5 only as "`dispositions` field is present, populated
correctly, and in byte order".

Not a divergence in the record's favour — noted because run-3's failure
to carry the two normative fixtures (`NumField() == 15`, the fourteen-key
list) shows they live only in S5 and are not visible from C4. C4 fixes
the *position* but not the *count*.

**Lands on**: `0024:C4` — a pointer to S5's two fixture values, or the
values themselves.

### D-13 — the `[emit]` opt-in trigger, and coverage of the bare-table case → **0024:C1**

run-1 renders C1's ruling verbatim and correctly ("A bare `[emit]` table
with zero sub-tables is a zero-declaration model... The opt-in trigger
is the *count of declared keys*, not the table's presence"). run-2
renders it as a property of `EmitDecls` being empty ("`[emit]` is absent
or empty"), losing the count-vs-presence distinction. run-3 loses it
entirely, rendering only "If zero declarations".

Low-severity: all three land on the same behavior. Recorded because the
distinction C1 spends a paragraph on survived in only one of three runs,
which suggests the paragraph's placement (inside C1's long tail, after
the arm list and the one-arm ruling) is where it gets lost rather than
its content.

**Lands on**: `0024:C1`.

---

## 2. GUESS clusters

Contracts two or more runs marked GUESS. These are the candidate
rewrites, strongest first.

### G-1 — `Model.EmitDecls` nil vs. empty map (run-1 GUESS, run-2 GUESS) → **0024:C3**

Both runs mark it explicitly; both resolve it by analogy to `Model.Tags`;
neither found record text. run-3 asserts the answer without marking it.
Two-run GUESS with a third run's unmarked assertion is the cleanest
rewrite candidate in the set. See D-7.

### G-2 — the dispositions-join factoring and its home package (run-1 GUESS, run-2 GUESS, run-3 GUESS-on-caller) → **0024:C4**

All three mark some part of this. run-1 and run-2 GUESS whether it is
extracted; run-3 GUESSes only its *caller*, having already committed to
a named method on `*Model`. Three-run GUESS. See D-2.

### G-3 — the `ConformValue` throwaway-`TagDecl` adapter (run-2 implicit, run-3 explicit GUESS) → **0024:C3**

run-3 marks it GUESS: "Construct a temporary `TagDecl{Kind: d.Kind,
Domain: d.Domain}` to pass to `ConformValue`, per C3." run-2 states it
as fact citing C3. run-1 states it as fact and adds the safe-by-omission
dependency.

C3 actually fixes this precisely — "the check constructs a throwaway
`TagDecl{Kind: d.Kind, Domain: d.Domain}` per call and passes it". That
run-3 marked a *determinate* clause as a GUESS is itself the finding:
the clause is buried in C3's third paragraph behind a
"not in tension with A4" framing, so it reads as rationale rather than
as the normative construction it is.

**Lands on**: `0024:C3` — promote the struct literal out of the
rationale paragraph.

### G-4 — source-schema struct name and `any` spelling (run-1 GUESS, run-3 GUESS-with-wrong-mechanism) → **0024:C1**

run-1 GUESSes `sourceEmitDecl` and `any` vs `interface{}` vs
`toml.Primitive`. run-3 GUESSes a decoder option instead. Two runs
guessing at the same silence, arriving at incompatible mechanisms. See
D-8.

### G-5 — `loadEmitDecls`/`checkRuleEmit` signatures (run-1 GUESS; run-2 unmarked but invented `*sourceDoc`) → **0024:C2**

run-1 marks it: "The RDR names all three functions... but writes no
parameter list or return type for any of them." run-2 does not mark it
and invents a type. See D-1 — the highest-value rewrite in the set,
because it is the one place a reconstruction produced a type name that
does not exist in the repo.

### G-6 — `EmitDecl`'s exportedness and `resolvePayload`'s struct tags (run-1 GUESS, run-2 GUESS) → **0024:C3**, **0024:C4**

run-1: "`EmitDecl` is exported and its fields are exported (C3 names
`Kind`, `Domain`, `Dispositions` in exported spelling, so this is
near-certain, but the RDR never writes the `type` line)." run-2 GUESSes
the JSON struct tags, inferring them from S5's key list and 0010:C4
convention. Low severity — both resolved identically and correctly —
but it is two runs guessing at the absence of a single `type` line.

**Lands on**: `0024:C3` (the `EmitDecl` type line) and `0024:C4` (the
`json:"dispositions"` tag, no `omitempty`).

### G-7 — S8/S9 unread (run-2 self-declared gap) → **0024:§testing-strategy**

run-2 states it did not read S8/S9 in full and that "any signature-level
detail from them is not reflected above". run-3 lists S8/S9 by title
only. Single-run gap on run-2's side, but it compounds with run-3's
title-only treatment of the whole scenario set. Not a contract defect —
recorded so the parent does not read S8/S9 agreement as confirmed.

---

## 3. Agreement

All three runs rendered these identically, confirming the RDR is
determinate there: the three category slugs and their spellings
(`malformed_emit_declaration`, `unknown_emit_key`,
`emit_value_out_of_domain`) and their registration in
`table.Categories()` (**0024:D-naming**, S6); `EmitDecl`'s three field
names, types, and semantics (`Kind string`, `Domain []string`,
`Dispositions map[string]string`, member→token) and its explicit
non-`TagDecl`-ness (**0024:C3**); the carrier's location and keying
(`Model.EmitDecls map[string]EmitDecl`, sibling of `Model.Tags`,
**0024:C3**); the four `kind` tokens `enum | bool | int | scalar` and
that only `enum` takes a domain (**0024:C1**); the two domain forms and
the bytewise-sorted union with grouping recoverable from `Dispositions`
(**0024:C1**, **0024:C3**); the two-step insertion point and order
(`loadEmitDecls` then `checkRuleEmit`, immediately after `loadTags`,
ahead of `normalizeRules`, reading `sourceRule.Emit` not normalized
rows — **0024:C2**); fail-fast single-refusal cardinality with
deliberately unspecified order (**0024:C2**); the `dispositions` field's
name, position immediately after `emit`, absence of `omitempty`, and
`{}`-never-`null`-never-omitted rule (**0024:C4**); the row-keyed
never-padded join semantics and the escape-row-joins-its-own-values rule
(**0024:C4**); `flow next` carrying no `dispositions` (**0024:C4**); and
that the kernel and `graphlint` read none of this (**0024:C3**).
