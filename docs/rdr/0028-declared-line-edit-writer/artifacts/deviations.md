# Deviations — RDR 0028 declared-line-edit-writer

Entries pre-seeded by the 7.1 cluster reconcile (0026-0027-0028, 2026-09-03)
are left OPEN for Stage 8: running the named check is the disposition; an
entry escalates only if the check contradicts a contract.

D1–D3 are the pre-seeded three, now discharged. D4 onward are Phase 2's.

## D1 — `TestReq77`'s tail-slice assertion breaks on C1.4's append

**Type**: TEST-FIXTURE. **Status**: RESOLVED (Phase 2, 2026-09-04).

**Situation.** `internal/table/command_carrier_0025_test.go::TestReq77_…`
asserts 0025:C5's six categories as the TAIL of `Categories()`
(`all[len(all)-len(want):]`). `C1.4` appends five after them, so the shipped
test goes red with no behaviour change. C1.4's "the list's size is not a
contract (0025:C5)" holds; the test over-asserts.

**Check.** Rewrite `TestReq77` to relative order — the six contiguous and in
clause order, followed by the five — and run it with this record's S6:
`go test ./internal/table -run 'TestReq77'` green.

**Disposition.** Already discharged by Phase 1: the shipped `TestReq77`
locates the six with `slices.Index` and slices relatively, so C1.4's append
does not disturb it. Verified green after the append landed.

A SECOND test in the same file carried the same tail-slice shape and D1 did
not name it: `TestReq79_TheCategorySetIsAppendOnlyAndDuplicateFree` computed
its unchanged-prefix as `all[:len(all)-len(six)]`, which C1.4's six-member
append pushes over 0025's own categories. It is the identical defect and
takes the identical fix — the prefix is now located at
`slices.Index(all, six[0])`. Recorded here rather than as a new entry
because it is D1's own class, and C1.4 states the rule outright: the
assertion is RELATIVE order, "never a tail position or a count".

## D2 — Description text for the five `edit_*` categories

**Type**: SPEC-UNDER. **Status**: RESOLVED — no-op (Phase 2, 2026-09-04).

**Situation.** 0027 Phase 3 ships a per-category description surface. C1.4
names wire strings, constants and order only. If 0027's surface is total over
`Categories()` (0027's deviations D2 decides), the five need text on it.

**Check.** After 0027 lands: its accessor's totality test and `make check`'s
docs gate green with the five appended. If 0027 chose per-category opt-in,
close as no-op.

**Disposition.** 0027 chose PER-CATEGORY OPT-IN, so the appended categories
need no text. Two pieces of shipped evidence, both live because 0027 is
already Implemented:

- `internal/table/category.go::categoryDescriptions` documents itself as
  "PER-CATEGORY OPT-IN, not total over Categories(): a category absent here
  carries no text and renders as its identifier alone", and gives the reason
  — `Categories()` is append-only and its total size is not a contract, so
  requiring text for every member would couple the map to every future
  append.
- `internal/cli/lint.go::lintExtendedDesc` ranges `Categories()` and
  `continue`s on a category with no description, so the six appended
  categories contribute nothing to the rendered help body and the generated
  `docs/cli-reference.md` is byte-identical.

The docs gate and the full suite are green with all six appended. Closed as
a no-op. The count discrepancy the entry carries ("five") is the stale one
req-list ASSUMPTION A-2 records; C1.4 locks SIX and six were appended.

## D3 — Merged clause-3 behaviour of `carrierDefect` after both land

**Type**: TEST-FIXTURE. **Status**: RESOLVED (Phase 2, 2026-09-04).

**Situation.** C1.6's declared-tag family check changes clause 3's membership
test and `carrierDefect`'s signature; 0027 rewrites the adjacent exemption
comment and clause 4's predicate. Same check as 0027's deviations D3, owed by
whichever lands second.

**Check.** `["nice","sh","-c","cat {tag.x}"]` with `x` declared, and
`["nice","sh","-c","cat {nope}"]`, both report `command_shell_interpreter`
after both records land, in either order.

**Disposition.** Both report `command_shell_interpreter`. Run against the
merged tree with `x` declared `observed`:

```
argv="cat {tag.x}" => category="command_shell_interpreter"
argv="cat {nope}"  => category="command_shell_interpreter"
```

The two arms are independent for a reason worth recording. C1.6's admission
is WHOLE-ELEMENT, so `"cat {tag.x}"` is not a `{tag.<key>}` element at all
and never reaches the family check; 0027's brace-bearing-element exemption
for the interpreter form is what carries both, unchanged. The signature
widened as anticipated (`carrierDefect` now also takes the model's tag
declarations, needed to answer "is this key declared"), and the clause-4
predicate is untouched.

## D4 — `Dump` gains an ACCESSOR line for a declared line rule

**Type**: SPEC-UNDER. **Status**: recorded; implemented (Phase 2).

**Situation.** REQ-101 requires "dump/normalize round-trip carries it (A6)"
and Phase 1 asserts `strings.Contains(table.Dump(m), "edit")`. The shipped
`Dump` renders ROWS only — it emits no accessor at all — so before this
change the assertion passed TAUTOLOGICALLY on the fixture's model id
(`editflow`), which contains the substring. The dump of an `edit`-carried
model was byte-identical to that of a carrier-less one, which is exactly the
failure A6 names: the surface a reviewer reads would not be the model that
runs.

C1 names no dump FORMAT for the carrier, so the line's shape is a new
output surface the Normative Contracts do not fix — hence SPEC-UNDER rather
than a mechanical translation.

**Choice taken.** One line per declared rule, after the `MODEL` header and
before the rows:

```
ACCESSOR write.<id> carrier=edit key=<key> anchor=%q replace=%q clear=%q
```

Grounded on three shipped constraints. It is emitted only for an entry
carrying an `edit` table, so every model predating this carrier dumps
byte-identically and no existing assertion moves. The walk is over sorted
entry ids and sorted keys, never a map's range order, because `0002:C19`
requires repeated dumps of one model to be byte-identical. And it is a
separate line rather than a new `dumpColumns` member, because that
vocabulary is the ROW column set an author may reorder via `[dump]` and the
loader enforces its completeness — adding a non-row column there would make
every existing `[dump]` declaration incomplete.

**Future interpretation.** A successor that fixes a dump format for the
carrier should treat this line as the incumbent, not as a decision the
record made.

## D5 — `edit_anchor_invalid`'s brace discriminator in an `anchor`

**Type**: SPEC-UNDER. **Status**: recorded; implemented (Phase 2).

**Situation.** Two clauses constrain the same bytes and compilation cannot
separate them:

- C1.4 — `edit_anchor_invalid` fires when the anchor "carries any other
  `{…}` form".
- C1.2 — it fires "never for a brace RE2 itself accepts".

RE2 accepts BOTH `{artifact}` (as literal text) and `\d{4}` (as a bounded
quantifier); `regexp.Compile` returns nil for each. So "a brace RE2 accepts"
cannot be read as "a brace that compiles", or C1.4's clause would be
unreachable and the two would be jointly unsatisfiable.

Phase 1 pins both sides: `TestReq17And26And65` requires `^{artifact}$` and
`^{status}$` to be `edit_anchor_invalid`, while `TestReq25And28And121`
requires `^\d{4}$`, `^a{2,3}$` and `^\{\{x$` to load clean.

**Choice taken.** The discriminator is the brace's SHAPE. An unescaped `{`
whose contents are a well-formed repeat spec — `{N}`, `{N,}`, `{N,M}` — is
quantifier syntax and reaches RE2 untouched; any other unescaped `{…}` is
the "other form" C1.4 refuses. An escaped `\{` is ordinary pattern text and
is never scanned, which is what keeps REQ-27 true.

This satisfies both clauses exactly on the fixtures each pins, and it is the
only reading that does: it refuses `{artifact}` and an entry's own
`{status}` — neither of which the document is likely to carry literally, and
both of which would silently anchor on literal braces — while admitting every
quantifier form a non-trivial anchor needs.

## D6 — The `edit` write's read-back gate pre-check is scoped to `edit`

**Type**: IMPL-DECISION. **Status**: recorded; implemented (Phase 2).

**Situation.** C1.3 `read-back:` requires a write whose role reader is
command-backed to refuse BEFORE mutation when the gate is off. The clause
sits in C1.3, whose subject is the `edit` carrier, but its wording ("the
write MUST refuse") does not name the carrier at the predicate.

**Choice taken.** The pre-check fires only for an entry carrying an `edit`
table (`len(def.Accessor.Edit) != 0`). A command-backed WRITER already
refuses at its own pre-spawn ladder under 0025:C6, and widening the check to
every carrier would re-decide 0025's disposition for a case this record does
not own. C1.6's closing line — "no other change to 0025:C1–C1.6 … the gate
[is] untouched" — is the ground.

Recorded because a successor extending the pre-check to another carrier is
changing 0025's behaviour, not this one's.

## D7 — `ErrDeclaredRequest` and `Refusal.DeclaredRequest()`

**Type**: SPEC-UNDER. **Status**: recorded; implemented (Phase 2).

**Situation.** C1.3 EXIT GROUP: says "the executor-facing typed `Err`
discriminates at `flow_exec.go::accessorFailureOf`" but names no sentinel,
no field and no accessor. REQ-42 forbids a new refusal CLASS, and the
accessor-level `Refusal` carries only a free-text `Detail` across the
executor seam — the binding's error itself does not survive `refusalOf`. So
the clause's discriminator needs a named surface the record does not supply.

**Choice taken.** Two additions, both minimal and both additive:

- `accessor.ErrDeclaredRequest`, a sentinel a binding wraps into its
  existing `ExecError.Err` slot. It travels on the seam that already exists
  rather than widening one, which is the same move 0025:C4 made for the
  stderr tail.
- `Refusal.declaredRequest`, unexported, with an exported
  `DeclaredRequest()` reader — mirroring `applied` exactly, and for the same
  reason: only the path that minted the refusal knows the answer, so a
  caller must not be able to assert it.

No refusal class is added: the class stays `execution_failure` and the
Detail is unchanged, which is what REQ-42 and JDR 0003 §D1 require. What
moves is the exit GROUP and the `findings[]` carriage, which is exactly
what the clause fixes.

The CLI code `flow-write-edit-refused` is likewise a new surface, but its
spelling is explicitly non-normative here ("Stage 8's", REQ-44) and no test
pins it.

## D8 — Phase 1 fixture corrections

**Type**: TEST-FIXTURE. **Status**: RESOLVED (Phase 2).

Phase 1's fixtures carried defects that blocked whole files or contradicted
a clause the same suite pins elsewhere. Each correction is listed with the
contract text that decides it. No test's CLAIM was weakened, skipped or
deleted; in every case the fixture was changed so the claim could be reached.

**(a) Three model-level defects, blocking every case in two files.**

- `kind = "string"` on `[tags.nnnn]`. RDR 0003's type-model vocabulary is
  five tokens and `string` is not one; `internal/table/model.go` says so by
  name ("`string` is not one of them"). Corrected to `scalar`, which is the
  shipped precedent — `internal/table/roundtrip_test.go:329` performs the
  same `string`→`scalar` rename on a spike fixture.
- `single_valued = true` on that same now-`scalar` declaration. A scalar has
  no finite declared domain to partition, so it admits no such marker.
  Removed.
- `terminal = ["done"]` with no `[context.done]`. The 0025 fixture these
  were derived from carries the context block; the derivation dropped it.
  Restored.

**(b) Six anchors paired with a replacement that de-anchors their own rule.**
C1.3 `re-anchor:` runs on every apply and refuses `edit_anchor_unstable`
when a rule's post-edit line no longer selects. The affected scenarios claim
something else entirely — that `-rf`, `$(rm -rf /)` and `*.md` are ordinary
LINE DATA (REQ-21/29), that `$$` and `{{`/`}}` expand as stated
(REQ-22/121), that `replace` emits the WHOLE line (REQ-2/3), that a
substituted value is never re-scanned (REQ-19/23), that a `<clear>` literal
in a template is bytes (REQ-82), and that a `-`-prefixed tag value is
admitted in an anchor (REQ-111).

Each fixture's anchor now also matches its own output, which is the shape
the record's own Coherence table row 3a prescribes (`^- \*\*Status\*\*:
(.+)$`, "post-edit re-matches its own anchor → no `edit_anchor_unstable`").
The narrow `(\w+)(.*)` form stays where it earns its keep — the MVV's
qualifier backreference — under the name `statusAnchor`, beside a new
`statusAnchorAny` for the value-admission scenarios.

**(c) Three `replace` templates carrying `{tag.<key>}`.** C1.2 admits
`{tag.<key>}` in an ANCHOR only; a `replace` admits `{<key>}` for the rule's
own key "and no other", and `TestReq15And24And66` pins
`a_tag_placeholder_is_not_admitted_in_replace` as `edit_template_invalid` in
the same suite. The three now capture the row identity in the anchor and
carry it forward by backreference, which is what the record's own
Illustrative Code does (`replace = '${1}{status}${2}'`).

**(d) The MVV model served one owned tag with two readers and two writers.**
RDR 0002 fixes the arity — every owned tag has exactly one reader, and every
written key exactly one writer — and this record does not amend it. The
record's own Illustrative Code has the same shape and would refuse
identically; see D9. The fixture's shared index now takes its OWN owned key
(`row_status`) on the `readme` role. Nothing the MVV asserts turns on the
spelling: both writers are still `edit`-carried, neither declares a
`command` or a `path`, and step 1's claim — the model lints clean with no
wrapper script anywhere — is unchanged.

**(e) The MVV README anchor's title group dropped its padding.** The anchor
captured `([^|]*)` from between ` | ` delimiters while the replacement
re-emitted it between bare `|`, so the rewritten row lost the cell's
surrounding spaces and no longer matched its own anchor. The group now spans
the padding.

## D9 — RDR 0002's accessor arity fences out a shared-artifact reader pair

**Type**: DEPENDENCY-LIMIT. **Status**: OPEN — carried to the completion
gate. Does not block Phase 2; the unit surface and all eight MVV steps are
green.

**Situation.** C1.6's stated purpose is a SHARED artifact: "0004:C12
read-back re-reads the SAME role; a shared artifact (an index) has one
reader for many records, and without an identity in its argv that reader
cannot say which row it read." The MVV fixture paragraph names two readers
and two writers over the same conceptual status — one on the record, one on
the README row — and the record's own Illustrative Code writes them with the
SAME key:

```toml
[read.readme]
role = "readme"
keys = ["status"]
...
[write.readme]
role = "readme"
keys = ["status"]
```

RDR 0002 refuses that model. `internal/table/load.go::checkAccessorBindings`
enforces "every OWNED tag MUST be served by exactly one reader" and "every
key any rule's write block or clear list names … exactly one writer". A
second reader and a second writer over `status` is
`malformed_accessor_binding`, and the record's Illustrative Code would
refuse at lint exactly as the MVV fixture did.

**Why it is a DEPENDENCY-LIMIT and not a spec defect.** The clause C1.6
states is satisfiable; what is not satisfiable is the particular KEY
SPELLING the illustrative model uses. Nothing in C1 or C1.6 requires the two
roles to share one key — the shared thing is the ARTIFACT, and the tag is
per-role state. RDR 0002's arity is not this record's to amend, and no
clause of 0028 says it is.

**Resolution taken, and it is the defensible one.** The shared index takes
its own owned key on its own role (`row_status` on `readme`), one reader and
one writer each. Every property C1.6 exists to deliver survives intact: the
README reader still carries `{tag.nnnn}` in its argv, still reads ONE row,
and the writer that edited that row and the reader that verifies it still
agree on which row it was — which is the whole claim. All eight MVV steps
pass at the accessor seam with this shape.

**What the completion gate should decide.** Whether the record's
Illustrative Code should be read as normative (it is marked "Illustrative"
and "field spellings are C1's", which reads as non-normative), and whether
the consumer's `rdr-write.toml` migration — explicitly out of scope here
(REQ-103, the consumer's kata rdr#yjye) — should carry the two-key shape or
petition RDR 0002 for a per-role reader arity. Phase 2 takes the two-key
shape because it is the only one that lints under the contracts as they
stand.
