# Deviations — RDR 0027 inline-shell-detection-scope

Entries pre-seeded by the 7.1 cluster reconcile (0026-0027-0028, 2026-09-03)
are left OPEN for Stage 8: running the named check is the disposition; an
entry escalates only if the check contradicts a contract.

## D1 — S6 pins `Categories()` order relatively, not as a tail golden

**Type**: TEST-FIXTURE.
**Status**: discharged (was OPEN, pre-seeded 7.1 PW1; the named check ran green in Stage 8).

**Situation.** `0028:C1.4` appends five categories after 0025:C5's six.
Shipped `internal/table/command_carrier_0025_test.go::TestReq77_…` asserts
the six as the TAIL slice of `Categories()` and breaks on that append.
`0027:S6`'s "membership and order unchanged" is this record's own delta claim
and holds in either landing order only if authored as relative order (the
`command_shell_interpreter` constant's position relative to its neighbours),
never `len`- or tail-equality — 0025:C5 makes the list's size a non-contract.

**Check.** S6 green against a `Categories()` with five extra members appended
at the tail.

**Disposition (Stage 8).** Ran green. `TestReq48_TheCategoryWireStringAndItsRelativePositionAreUnchanged`
asserts the six C5 categories by RELATIVE clause order (indexing each and
requiring the sequence to ascend) plus wire-string registration and
no-duplicate-registration. No tail slice, no `len()` equality, so 0028's
five `edit_*` appends cannot break it. No contract contradicted.

## D2 — Is the Phase 3 description surface total over `Categories()`?

**Type**: IMPL-DECISION.
**Status**: decided (was OPEN, pre-seeded 7.1 PW2) — per-category opt-in;
full decision and grounds under "D2 — DISPOSED" below.

**Situation.** Phase 3 adds an accessor pairing a category with its text,
rendered on `--help-all` and mirrored into `docs/cli-reference.md`. If that
accessor or the docs staleness gate asserts totality over `Categories()`,
0028's five `edit_*` categories fail it or ship undescribed, and neither
record names the obligation.

**Check.** Decide total vs per-category opt-in and record it here; the
accessor's test states which. If total, 0028's deviations D2 owes the five
texts. `make check`'s docs gate green after 0028's append.

## D3 — Merged clause-3 behaviour of `carrierDefect` after both land

**Type**: TEST-FIXTURE.
**Status**: discharged (was OPEN, pre-seeded 7.1 PW4, critique C-8; the named check ran green in Stage 8).

**Situation.** 0028 turns clause 3's placeholder membership test into a
declared-tag family check (a signature change) and both records rewrite the
clause-3 exemption comment and the function's doc comment; `0027:A1`'s line
anchors shift. The JCs disposed the function as disjoint clauses; clause 3 is
touched by both.

**Check** (whichever lands second): `["nice","sh","-c","cat {tag.x}"]` with
`x` declared, and `["nice","sh","-c","cat {nope}"]`, both report
`command_shell_interpreter`; the clause-3 exemption comment no longer reads
"argv0".

**Disposition (Stage 8).** Ran green — 0027 landed first.
`TestReq43_ABraceBearingStringUnderAWrapperReportsTheInterpreterDefect`
covers both named vectors (declared `{artifact}` and undeclared `{nope}`
under `nice`), plus `timeout 5` and `env -i`, asserting the category STRING
is `command_shell_interpreter` and explicitly NOT
`command_unknown_placeholder` — a load-failed-only oracle would pass against
unchanged code and prove nothing. The clause-3 exemption comment no longer
describes an argv0 anchor; the two surviving `argv0` mentions in `load.go`
are a historical citation of 0025:C5's superseded line and user-facing
remediation advice, neither describing the predicate. 0028's clause-3
signature change is still owed its own re-run when it lands. No contract
contradicted.

---

# Stage 8 entries (Phase 2 implementation)

## D4 — S4-style green row `["xargs","busybox","sh","-c","echo"]` refuses under C1

**Type**: TEST-FIXTURE.
**Status**: mechanical translation (fixture corrected;
assertion NOT relaxed).

**Situation.** `Req12_TheInterpreterSetStaysOpenUnderThePositionFreeScan`
(REQ-12/REQ-54/REQ-70) asserts that an UNLISTED spelling stays admitted. Five
of its six rows carry no listed interpreter at all (`perl`, `python3`,
`nodejs`). The sixth was authored as `["xargs","busybox","sh","-c","echo"]`,
which carries a genuinely listed `sh` at i=2 and its own `-c` at j=3, so C1's
`refuse iff there exist i < j` quantifier refuses it — correctly.

**Grounding.** The record names `busybox` only as an unlisted BASENAME
spelling (`0027:C1` predicate line; §cross-cutting-concerns line 481), never
as an argv vector. The authored vector is structurally identical to
`Req28_TheNewFalseRefusalClassIsAcceptedNotSuppressed`'s
`["wc","-l","python","-c"]`, which the SAME suite requires to REFUSE (REQ-28:
"a listed basename as a data argument under an unlisted argv0 … the class is
accepted as a stated cost"). Admitting the `busybox` row would require reading
argv words before the interpreter — exactly what REQ-4 ("Nothing before argv[i]
is read") and REQ-5 ("no wrapper table exists and none may be added") forbid.
Two Phase 1 tests therefore contradicted each other, and C1 settles it against
the `busybox` row.

**Resolution.** The row's vector is corrected to `["xargs","busybox","-c",
"echo"]` — `busybox` as an unlisted spelling taking an inline-code flag, which
is the discriminating property the row exists for and needs no listed pair
beside it. No assertion was weakened: the row still asserts GREEN, still under
a wrapper, still on the OPEN-deny-list axis. The predicate is unchanged.

## D5 — `table.CategoryDescription` is a NEW exported surface

**Type**: SPEC-UNDER.
**Status**: mechanical translation (surface named by
REQ-37/REQ-38 as an obligation; the SYMBOL is unpinned).

**Situation.** `0027:REQ-37` states in as many words that "No description
surface exists for `table.Category` today … so this phase adds the surface as
well as the text", and REQ-38 fixes where it lands ("an exported accessor
pairing the category with its text, consumed by that body, registered via
`withExtendedHelp` and mirrored into `docs/cli-reference.md`"). The record
pins the SURFACE and the TEXT but names no Go symbol, so the accessor is a new
public function the Normative Contracts do not name by identifier and is
recorded here per the ADDITIVE-IS-NOT-EXEMPT rule.

**Shipped.** `internal/table.CategoryDescription(Category) (string, bool)` —
per-category opt-in (see D2), returning `false` for a category carrying no
text. Consumed by `internal/cli/lint.go::lintExtendedDesc`, which is
registered via `withExtendedHelp` and mirrored by
`internal/cli/docs.go::runDocs`. `Categories()`' membership, order and the
`command_shell_interpreter` wire string are untouched (REQ-48).

## D6 — intrastate#q2q1 prerequisite (REQ-40)

**Type**: DEPENDENCY-LIMIT.
**Status**: mechanical translation (discharged by
recording, per the clause's own terms).

**Situation.** `0027:§prerequisites` requires "intrastate#q2q1's disposition
recorded (landed first, or closed as subsumed) — either way the `env` walk is
removed here".

**Disposition.** SUBSUMED. The `env`-chain walk in
`internal/table/load.go::interpreterForm` is deleted outright by this build,
not made conditional; the position-free scan catches every `env` option form
q2q1 enumerates because an `env`-prefixed argv simply carries the interpreter
at some i >= 1 and nothing before argv[i] is read.
`Req44_TheDeletedEnvWalkIsSubsumedByThePositionFreeScan` is the subsumption
proof over all eight of `0027:S3`'s forms, which is what the clause asks for.

## D7 — promise-half fixture omits `[tags.recognized]` and never reaches C5

**Type**: TEST-FIXTURE.
**Status**: mechanical translation (fixture completed;
assertion NOT relaxed).

**Situation.** `Req20_TheAdmittedFormDisclosureLivesOnTheHelpSurfaceNotInARefusal`
loads an inline model carrying `["nice","sh","-c","echo hi"]` and asserts it
refuses as `command_shell_interpreter`. The fixture declared no
`[tags.recognized]`, which `table.Load` refuses unconditionally as
`malformed_model_declaration` (`internal/table/load.go:270`) BEFORE reaching
C5's command clauses — so the vector never met the predicate and the test
asserted nothing about the interpreter form.

**Grounding.** The sibling fixture built for the same purpose in
`internal/cli/inline_shell_mvv_0027_test.go::mvv0027Model` declares
`[tags.recognized]`, as does `internal/table`'s shipped 0025 command-carrier
fixture. The requirement is 0002/0025 loader shape, untouched by 0027.

**Resolution.** The four-line `[tags.recognized]` declaration is added to the
fixture. No assertion changed: the test now actually reaches the arm it was
written to assert, which STRENGTHENS it — it was passing-by-accident-of-order
before this run only because the predicate arm was unreachable.

## D2 — DISPOSED: the description surface is PER-CATEGORY OPT-IN

**Type**: IMPL-DECISION.
**Status**: decided (was OPEN, pre-seeded 7.1 PW2).

**Decision.** Per-category opt-in. `table.CategoryDescription(Category)
(string, bool)` returns text where declared and `false` otherwise;
`lintExtendedDesc` iterates `Categories()` and renders only described members,
so an undescribed category is silently skipped rather than failing a gate.

**Grounds.** (a) `0025:REQ-79` makes `Categories()`' total size a non-contract
at any point, so a totality assertion would couple this record to every future
append; (b) `0027:C1` ships text for exactly ONE category and `0027:S7` asserts
exactly that one; (c) the total reading breaks `0028:C1.4`'s five `edit_*`
categories, which neither record obliges anyone to describe. The Phase 1 tests
assert no totality, consistent with this reading.

**Consequence for 0028.** Its five `edit_*` categories may ship undescribed and
owe no text. `make check`'s docs gate is satisfied by regenerating
`docs/cli-reference.md`, which this build does; the gate is a staleness check
over the rendered tree, not a coverage check over `Categories()`.

## D8 — REQ-32 withdrawn: "Illustrative — shape only" is not a testable clause

**Type**: TEST-FIXTURE.
**Status**: mechanical translation (REQ-list classification corrected; no
contract text is touched and no behaviour changes).

**Situation.** Phase 0 extracted `0027:§illustrative-code`'s line
"Illustrative — shape only." as `REQ-32`, a NEGATIVE REQ reading "no test
asserts the illustrative Go literally". Phase 1 could not cover it and left
an empty coverage cell, which made `impl_orphans` read `1+` and the
COMPLETION GATE stop at `stopped:req-orphans`.

**Why it is not a REQ.** The line is a CAVEAT ON A CODE BLOCK — it tells the
implementer the following snippet is non-binding shape, not a contract.
`rdr inspect 0027 --filter elements` confirms its line (136) falls inside NO
labelled contract element; it sits outside `0027:C1` entirely, so it carries
no normative force. A clause that obliges nobody to do anything is not a
testable clause. Its stated "compliance evidence" was the ABSENCE of a test,
which no assertion can express — the tell that the extraction was wrong.

**Disposition.** REQ-32 is withdrawn from `req-list.md` and its coverage row
removed rather than left empty: an empty cell would misreport a Phase 0
classification error as a coverage gap, and the gate would be stopping on a
REQ that should never have existed. The snippet's actual content — the
position-free scan's shape — is covered behaviourally by REQ-41..REQ-47.

**Not done.** The gate was NOT satisfied by tagging `impl_orphans=0` over an
empty row. Three of the five original orphans (REQ-0b, REQ-35, REQ-40) were
closed with real mutation-verified tests; this one was withdrawn on its
merits. `impl_orphans` is now genuinely 0.

## D9 — D8's "mutation-verified" claim was wrong for the authoring guards

**Type**: TEST-FIXTURE.
**Status**: mechanical translation (test oracles retargeted; no contract text
is touched and no behaviour changes).

**Situation.** D8 records that three of the five original orphans (REQ-0b,
REQ-35, REQ-40) "were closed with real mutation-verified tests". The intent
was right; the oracles were not. All three guards in
`internal/table/inline_shell_authoring_0027_test.go` asserted file-wide
`strings.Contains` over an `os.ReadFile` body, so a TARGETED edit to the exact
fact each one names left the guard green. Phase 3c's verification used coarse
global `sed`, which is why it reported green as red. Per the no-amend rule D8's
text stands; this deviation supersedes its "mutation-verified" claim for these
three rows.

**What was actually loose.**

- REQ-0b — `load.go` cites `0027:C1` at five sites; only two are the C5 comment
  blocks the clause names. Rewriting BOTH of those to `0025:C5` left the guard
  green off the other three.
- REQ-35 — the four argv literals were checked ANYWHERE in
  `command_carrier_0025_test.go`, not inside the named function. Truncating the
  probe loop to `cases[:1]` and replacing its `t.Fatalf` with `return` left
  every literal in place and the guard green while the probes were bypassed.
- REQ-40 — `deviations.md` was searched file-wide, and D3 says "landed first"
  under an unrelated deviation. D6's disposition could be gutted to `TBD` and
  the check still passed off D3. Scoping to `## D6` alone is insufficient:
  D6's Situation QUOTES the clause's own "landed first, or closed as subsumed"
  wording, so the section reads the requirement back to itself.

**Disposition.** The two authoring guards are retargeted to their named locus
and the third is deleted in favour of a stronger existing oracle.

- REQ-0b now locates each C5 comment block by the SUBSTANTIVE claim it makes —
  "exemption belongs to the INTERPRETER FORM" for clause 3, "argv0 line" for
  clause 4 — and requires `0027:C1` within that block. Each site fails
  independently under the targeted rewrite. Anchoring on the claim rather than
  on a line number or on the citation itself keeps the guard alive through
  reflow and relocation while still catching the rot.
- REQ-40 extracts the `## D6` section, then extracts D6's `**Disposition**`
  field within it, and runs the disposition-token check against that field
  alone. The `q2q1` check runs against the section.
- REQ-35's file-reading guard is DELETED. The clause protects the
  predecessor's REFUSAL BEHAVIOUR, not the byte layout of its test file, and
  `internal/cli/inline_shell_mvv_0027_test.go`'s
  `ReqMVV0027/step_4_the_0025_probes_still_refuse_and_the_category_order_holds`
  already drives the same four argv — `{sh,-c,cat {artifact}}`,
  `{bash,-c,echo hi}`, `{python,-c,print(1)}`, `{env,sh,-c,echo hi}` — through
  `table.Load` asserting `CatCommandShellInterpreter`. Re-deriving the property
  is strictly stronger than reading bytes: it survives any respelling or
  relocation of the probes and cannot be defeated by gutting the loop that
  drives them. `coverage.md` cites `ReqMVV0027/step_4` for REQ-35, the same
  citation form the REQ-41 and REQ-48 rows already use.

**Not done.** No new test was added. The reds here are TARGETED mutations, run
and reverted: rewriting `0027:C1` to `0025:C5` at either C5 comment block alone
fails `TestReq0b_…`, and replacing D6's `**Disposition.**` line with `TBD` while
leaving D3 intact fails `TestReq40_…`. `impl_orphans` stays 0; the rows now
name oracles that deliver what they claim.
