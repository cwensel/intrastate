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

**Type**: SPEC-UNDER. **Status**: recorded; implemented (Phase 2); CORRECTED
(Phase 3c, resolving Phase 3a FAIL-1 and a Phase 3b recorded finding).

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

**Choice taken (CORRECTED, Phase 3c).** The discriminator is the brace's
SHAPE, and the shape that matters is a PLACEHOLDER: an unescaped `{`
followed by one or more NAME bytes (ASCII letter or digit, `_`, `-`, `.`)
and a closing `}`, where an all-digit body is excluded as a repeat spec.
That — and only that — is the "other form" C1.4 refuses. Every other brace
reaches RE2 untouched. An escaped `\{` is ordinary pattern text and is never
scanned, which is what keeps REQ-27 true.

The rest of C1.4's own list is what settles the reading: its other two arms
are "fails to compile" and "names an undeclared tag key", both PLACEHOLDER
concerns. "Any other `{…}` form" read in that company means a form an author
wrote MEANING a substitution — and the only substitution an anchor admits is
`{tag.<key>}`. It refuses `{artifact}` and an entry's own `{status}`, which
would otherwise silently anchor on literal braces the document does not
carry, and it reaches nothing else.

**Why the original wording was wrong.** D5 as first written took the
discriminator to be "a well-formed repeat spec, or refuse", and claimed the
choice "satisfies both clauses exactly on the fixtures each pins". Both
halves were incomplete:

- Phase 3a **FAIL-1**: D5 named the pinned anchor fixtures as `^\d{4}$`,
  `^a{2,3}$` and `^\{\{x$` — the third ESCAPED. `0028` S14 names a BARE,
  unescaped `{{` by name, and its Expected requires "every brace form
  reaches RE2 untouched … `{{` is not an escape — and none refuses
  `edit_anchor_invalid`". A bare `{{` is not a repeat spec, so the original
  discriminator refused the one form S14 names to distinguish the two
  templates' escaping dialects. D5's scope statement was narrower than the
  behaviour it licensed.
- Phase 3b, recorded not pinned: `[{}]` and `[{]` — ordinary RE2 character
  classes an author writes to anchor a literal brace, both accepted by RE2 —
  were `edit_anchor_invalid` for the same reason.

The corrected discriminator admits all three (the inner `{` of `{{` is not a
name byte; `[{}]`'s brace closes immediately and names nothing; `[{]` is
unclosed) while keeping every fixture the original satisfied: `{artifact}`
and `{status}` still refuse (REQ-17/26/65), and `\d{4}`, `a{2,3}` and
`^\{\{x$` still load clean (REQ-25/28/121). It is strictly wider on the
admit side and identical on the refuse side.

**Covering test.** `TestEditBraceShape0028_PlaceholderShapeIsTheDiscriminator`
holds both verdicts in ONE table — which the shipped Phase 1 pair
(`TestReq17And26And65_…` / `TestReq25And28And121_…`) does not, so a future
narrowing could satisfy one by breaking the other silently.
`TestEditBraceShape0028_AdmittedAnchorsAreCarriedVerbatim` holds S14's
operational half: an admitted anchor's bytes reach RE2 unmangled, not merely
un-refused. Eight of the table's arms fail against the pre-correction
discriminator.

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

**Type**: DEPENDENCY-LIMIT. **Status**: needs author decision → RESOLVED (the
Illustrative Code is non-normative; the two-key shape stands). Does not block
Phase 2; the unit surface and all eight MVV steps are green.

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

**Gate resolution (Stage 8, unattended run).** Settled from the record's own
text, not deferred. The record marks the section `#### Illustrative Code` and
opens it "Illustrative — the consumer's two rules as they would read in
`rdr-write.toml` … (field spellings are C1's; the anchors are A3's spike
subject, not normative)", and Prerequisites again calls its projector verb
"verb illustrative; the consumer's to add". The section is therefore
NON-NORMATIVE on its own face: no REQ derives from it, and no C1 clause
requires the two roles to share one key. Phase 2's two-key shape is adopted as
the resolution — it is the only shape that lints under 0002's arity as it
stands, and every property C1.6 exists to deliver survives it. The consumer's
`rdr-write.toml` migration is out of scope here (REQ-103, kata rdr#yjye); that
kata carries the two-key shape. No petition to amend RDR 0002 is opened,
because no contract of this record needs one.

## D10 — `re-anchor:` over an UNPLANNED sibling rule: the stability reading

**Type**: SPEC-UNDER. **Status**: recorded; implemented (Phase 3c).

**Situation.** C1.3 `re-anchor:` states the invariant with no plan
qualifier — "**every rule's** anchor is run again over the POST-EDIT buffer
and must select exactly its own rewritten line (or, for a deleted line, zero
lines)" — and then names the case it exists for: "This is what stops a
replacement from de-anchoring itself or **poisoning a sibling rule's anchor
on the next run**". "On the next run" is decisive: a poisoned sibling only
matters on a LATER invocation, which is exactly the invocation that did not
plan it now. Phase 3b's ADV-3 pinned the gap — `prepare` dropped any rule
whose key the plan did not name, so `reAnchor` never saw it, and read-back
could not see it either (the key is not in the plan, so the oracle never
compares it).

**What the clause does not settle.** Its stated test is "exactly its own
rewritten line", and an unplanned rule HAS no rewritten line: it writes
nothing. Applied literally, the sentence has no referent for that rule. Two
readings are available:

- **cardinality** — an unplanned rule must select exactly one line of the
  post-edit buffer;
- **stability** — an unplanned rule must select, after the rewrite, exactly
  the lines it selected before it (shifted by any deletions), i.e. this
  entry's rewrite must not CHANGE what a sibling selects.

**Choice taken: stability.** Cardinality would newly refuse invocations over
a defect this run neither caused nor touched — an anchor already stale or
already ambiguous before the edit would condemn every plan naming any OTHER
key, which C1.3 `select:` assigns to the plan that names the rule
(`edit_anchor_unmatched` / `edit_anchor_ambiguous`, step (4), scoped to the
rules of the plan). Stability is exactly the property the clause's stated
purpose asks for and nothing more: it fires when and only when this entry's
rewrite moved the sibling, which is the poisoning.

The shift arithmetic is the one already used for planned rules — a held
pre-edit index moves down by the number of preceding deletions, and a hit on
a deleted line has no post-edit counterpart.

**Scope kept narrow.** An unplanned rule raises no refusal of its OWN at
steps (2)–(5): it has no planned value to scan, its `replace` is never
parsed or expanded, it contributes no deletion and no collision, and an
anchor it declares that cannot be parsed, bound or compiled makes it
contribute no post-edit witness rather than condemning a plan that never
touched it. Only step (6) reads it. This keeps C1.3 `precedence:` steps
(2)–(5) scoped to the plan, as their own text is, while step (6) reads
"every rule" as written.

**Covering tests.** `TestAdv0028_3_ReAnchorCoversRulesThisPlanDidNotName`
(the poisoning arm) and
`TestReAnchor0028_UnplannedSiblingRulesAreCovered` (the four arms around
it: an untouched sibling still applies; an already-unmatched sibling is not
this plan's refusal; a `<clear>` deletion that merely shifts the sibling is
not poisoning; a rewrite that de-anchors a sibling refuses before mutation).

No public surface is added: `unplanned` and `preHits` are unexported fields
of an unexported struct, and `shiftHits` is an unexported helper.

## D11 — C1.3 `precedence:` step (1): the gate pre-check is decided before the anchor's unbound tag

**Type**: SPEC-DEFECT. **Status**: recorded; NO code change (Phase 3c).

**Situation.** Phase 3a **FAIL-2** observed that an `edit` write whose ANCHOR
carries an unbound `{tag.nnnn}`, over a command-backed reader with the gate
OFF, reports the GATE refusal, while C1.3 `precedence:` step (1) orders the
three entry-level preconditions "C1.6's unbound `{tag.<key>}`, then C1.6's
`-`-prefixed bound value …, then the gate-off command read-back above".
Phase 3a scoped the consequence itself: both arms refuse before mutation and
the artifact is byte-identical either way, so it is a REPORTING-ORDER
question, not a mutation-safety one.

**Disposition: satisfiable as-is; no code change.** Three grounds, each
sufficient on its own.

**(a) The two refusals step (1) names are not the one that was probed.**
Step (1) names "**C1.6's** unbound `{tag.<key>}`" and "**C1.6's**
`-`-prefixed bound value", and C1.6 is titled "`{tag.<key>}` as a COMMAND
PLACEHOLDER (0025:C2 extended)". Both of its refusals are argv-side, minted
in `cmdbind.go::substitute` and stated by C1.6 `binding:`/VALUE: as firing
"BEFORE spawn". REQ-88's own covering test
(`TestReq88And89And90_…`) drives them through `cmdbind.Reader.Read` over a
`command` argv, confirming the site. FAIL-2's fixture instead put the
placeholder in the ANCHOR, whose binding is C1.2's (`anchor admits:`
"…bound on the invocation's context (A1, A7)"). C1.2 states no refusal site
and no precedence for it, and C1.3 step (1) does not enumerate it. So the
observed ordering is not the ordering the clause fixes.

**(b) On the path step (1) describes, the order is unrealizable in EITHER
direction.** C1.6's two refusals fire inside the READ-BACK reader's spawn —
`executor.go::invokeRead`, which `Write` reaches only after the gate
pre-check and after `Apply`. The gate is precisely what prevents that spawn.
With the gate off, no child is spawned, so neither C1.6 refusal is reachable
at all: they are not merely sequenced later, they do not exist on that path.
An ordering between "the gate refusal" and "a refusal only a passed gate can
produce" has no realizable case, which is why no fixture can witness it.
The shipped `TestReq51To54And120_…` covers steps (2)–(5) only, consistent
with step (1)'s internal order having no observable.

**(c) The record contradicts itself on the order, and the other statement
puts the gate first.** REQ-41, one sentence earlier in the same clause,
enumerates the same three as "the read-back gate below, C1.6's unbound tag
and C1.6's `-`-prefixed value" — gate FIRST. Phase 3b's ADV-2 read REQ-41's
order as governing and called the gate "**first** among the entry-level
preconditions". The two sentences disagree; the implementation matches one
of them. Recorded as SPEC-DEFECT for that reason rather than IMPL-DECISION:
the defect is in the record's text, not in a choice the implementation made.

**Reading taken.** Step (1)'s parenthetical — "(both are properties of the
binding, decided together before anything is spawned OR READ)" — attaches to
the two C1.6 refusals and explains why THEY are ordered together. The gate
is appended to the sentence as the third member of the SET whose Detail
discipline and exit group the clause is fixing, not as a third rung of a
ladder anything can climb. Read that way both sentences agree, ground (b)'s
vacuity is explained, and the clause's substantive requirements — refuse
before mutation, Detail names the gate or the placeholder, exit-2 group —
are all met.

**What a code change would cost, had one been owed.** The gate pre-check
must run before `Apply` (C1.3 `read-back:`: "the write MUST refuse BEFORE
mutation … because the reader is resolved before `Apply`", A2), and the
anchor's binding runs inside `Apply`. Inverting them needs a new
`WriteBinding` seam method letting the executor ask the binding whether its
entry-level preconditions would refuse — new public surface — since
`internal/accessor` cannot import `flowbind` (`flowbind` imports
`accessor`; the reverse is a cycle, C1.3 SITE:). Not taken, on grounds
(a)–(c).

**What the completion gate should know.** If a later author reads step (1)
as ordering the ANCHOR's unbound tag ahead of the gate, the fix is the seam
method above and it is a contract change, not a bug fix. The behaviour is
unchanged either way in every property C1.3 asserts: refusal before
mutation, byte-identical artifact, `execution_failure` class, a Detail
naming the gate or the placeholder, and — after Phase 3c — the exit-2 group
on both arms.
