Model: claude-opus-5[1m]

# Persona 3 — QA / Tester: A8 read-completeness re-entry pass

Scope: A8, the read-completeness normative clause, the "Read completeness"
Load-Bearing Decision, MVV / Validation Scenario 2's read assertions, the
read-truncation Risk and the Failure Modes incomplete-read paragraph, the
Proportionality "same contract" claim, and the four mini-check rows carrying
read completeness (Disposition Table, Oracle Discriminability scenario 2,
Fidelity Table `read` row, Desk Trace step 2). A1–A7 carried forward.

Findings are ranked by how completely they block a test from being written.

---

## HIGH

### H1. "The keys it was asked for" is never bound to a source — the requested key set has no defined origin

**Passage**: Normative Contracts, read-completeness clause ("A read accessor
MUST return the complete tag set for the keys it was asked for"); reinforced
by Approach ("the keys it was asked for"), Load-Bearing Decisions →
**Read completeness** ("every requested key resolved"), and Fidelity Table row
`read` over a requested key set.

Nothing in the RDR says *who* supplies the requested key set or where it is
recorded. Technical Design's definition list says each accessor definition
declares "expected tag keys", which is the only candidate — but the normative
clause does not cite it, and Alternative 1 calls the same thing an "expected tag
contract", a third name. Three names, no binding.

The spike is ambiguous in exactly the same place, and its two modes have
opposite testability: `read(...)` with an explicit variadic key list
(`main.go::read`, called at `main.go:361-363`) versus `read(...)` with no keys,
which falls through to `expectedTagKeys(art)` — the artifact's *own* key set
(`main.go:165-167`, `main.go::expectedTagKeys`). Under the second mode the
completeness predicate is self-fulfilling: the requested set is derived from
what the artifact carries, so no key can ever be absent and `incomplete_read`
can only ever fire on the `unreadable` map. Transcript line 1 and the replay
path (`main.go::replay`) both use that self-fulfilling mode; only lines 11-13
use the explicit mode.

**Test it prevents**: the central MVV assertion —
"an accessor which can resolve only some of the requested keys takes the refusal
branch" (Minimum Viable Validation). I cannot construct the arrange step. To
build a truncated read I must know which keys were requested; if the requested
set is derived from the artifact, the test case is unconstructible by
definition, and if it comes from the accessor definition, the definition struct
needs a field the Normative Contracts never require and the validator (Desk
Trace step 1, seven named codes) never checks. Also blocks the Disposition Table
row **All requested keys read** — I cannot assert "all" over an undefined set.

### H2. `incomplete_read` is never minted by a normative clause, and is not required to be distinct from `execution_failure`

**Passage**: Normative Contracts — the read-completeness clause says only
"take the refusal branch" and never names a class. The name `incomplete_read`
appears only in non-normative carriers: Load-Bearing Decisions → Read
completeness, Disposition Table row **Requested key unreadable**, Desk Trace
step 2, Failure Modes.

Contrast the timeout clause, which is explicit: "Timeout MUST be reported as its
own refusal class, distinct from execution failure and read-back mismatch."
There is no parallel sentence for incomplete read. So an implementation that
returns `execution_failure` for an unreadable key satisfies every normative
clause in the RDR: it took the refusal branch, it returned a typed refusal, it
did not thin the value set.

This matters because the two are genuinely adjacent in the spike itself: a
missing artifact role returns `refusalExecutionFailure` (`main.go:161-164`,
transcript line 4) while an unreadable *key* on a present role returns
`refusalIncompleteRead` (`main.go:169-172`, line 12). The distinction is a spike
decision, not a contract.

**Test it prevents**: Validation Scenario 2's assertion that "timeout, execution
failure, incomplete read, and gate indeterminate remain distinct refusal
classes." I can write a test that asserts a refusal, but I have no normative
basis for asserting *which* code, so the assertion is unfalsifiable against the
spec — a reviewer can always say `execution_failure` was conformant. Also blocks
the Disposition Table's "Refusal class minted" column for that row from being
used as a test oracle.

### H3. The success-branch representation of an absent key is explicitly unpinned, but every MVV witness asserts on it

**Passage**: Load-Bearing Decisions → **Read completeness**, final sentence
("the Resolve spike's `<absent>` sentinel is fixture shorthand, not a normative
representation, and implementation may choose a typed absence marker instead");
Fidelity Table row `read` over a requested key set, Lossy exemptions column
("how an absent value is *represented* is not pinned").

Both passages deliberately refuse to fix the representation. Yet the assertions
that are supposed to prove A8 all name it concretely: Validation Scenario 2
("returns it as an absent value, not a refusal … a genuine absence returned as a
value (`profile=<absent>`)"), Desk Trace step 2 witness
(`tags={profile=<absent>,status=Draft}`), and Oracle Discriminability scenario 2
(`tags={profile=<absent>,…}` as the discriminating shape).

**Test it prevents**: the genuine-absence success assertion in MVV Scenario 2.
The test must assert *something* about the returned tag set for a key the
artifact lacks. `assert tags["profile"] == "<absent>"` is explicitly
non-normative. `assert "profile" in tags` is the only spec-supported form, but
it does not discriminate: a value of `""` for an absent key passes it, and `""`
is indistinguishable from a legitimately empty tag value. There is no third
option, because the RDR names no predicate — no `IsAbsent()`, no typed marker,
no rule that absent values are distinguishable from empty ones. This is the one
finding where the spec actively removes the oracle it elsewhere relies on.

---

## MEDIUM

### M1. No rule for a read that times out mid-key-set — timeout and incomplete read have no precedence

**Passage**: Technical Design, phase two ("classifies timeout, unavailable
artifact, execution error, incomplete read, gate denied, and gate-indeterminate
outcomes" — a flat list with no ordering); Disposition Table, rows **Requested
key unreadable** and **Accessor exceeds its timeout** as disjoint input classes;
Desk Trace step 2 (both "completeness" and "bounded timeout" listed as in force
simultaneously).

A read that resolves three of five requested keys and then exceeds its timeout
satisfies both input classes. The spike dodges it by construction: `slow.read`
sleeps *before* the artifact lookup and key loop (`main.go:154-160`), so the two
never overlap. A real accessor reading keys incrementally will hit the overlap
routinely — and it is the single most likely production shape for a truncated
read.

**Test it prevents**: the timeout row and the incomplete-read row of the
Disposition Table cannot both be asserted against one implementation without a
precedence rule; a test for the overlap case has no expected value. Also blocks
Failure Modes' claim that "the completeness requirement turns that into an
`incomplete_read` refusal" from being verified for the timeout-induced
shortfall, which is the shortfall the Risk paragraph most plausibly describes.

### M2. Read-back is a read, but the RDR never says whether read completeness applies to it — no Disposition Table row for an unreadable key during read-back

**Passage**: Disposition Table (rows **Write succeeds, owned tag matches** /
**owned tag differs** / **non-owned tag changed** — no row for a read-back that
cannot read a key); Desk Trace step 5, whose assertions-in-force column lists
read-back obligations but omits the completeness assertion that step 2 carries;
Proportionality ("Read completeness is part of that same contract").

The read-back clause requires the executor to "re-read the same caller-supplied
artifact role" and verify pre-write observed/recognized values are unchanged. If
one of those pre-write keys is now unreadable, the outcome is undefined: it is
not a value mismatch (`read_back_mismatch` per the Fidelity Table's value-
equality strength), it is not obviously `incomplete_read` (that class is scoped
to read accessors in the Load-Bearing Decision), and the Disposition Table's
closing sentence claims "No input class exits silently" for a table that has no
row for it. The spike's `write` never exercises the `unreadable` map at all —
it clones `art.tags` directly (`main.go:240,248`), bypassing `read` entirely.

**Test it prevents**: a read-back-with-unreadable-key test in MVV Scenario 3.
The Round-Trip / Inverse Invariant says pre-write values "must remain
unchanged" — an unreadable key is neither changed nor verified, so the invariant
has no truth value and the test has no pass/fail criterion. This is the same
silent-corruption shape A2/A7 exist to close, reached through the A8 seam.

### M3. "Unreadable" has no observable definition — the truncation fixture cannot be built for any real accessor

**Passage**: A8 statement ("a key the accessor could not read"); A8 Evidence
(cites `main.go::newPartialArtifacts` as the witness); Load-Bearing Decisions →
Read completeness ("a key the accessor could not read is an `incomplete_read`
refusal"); Oracle Discriminability scenario 2 ("the unreadable-key … fixture").

Unreadability is realized in the spike as a hand-set boolean map on the fixture
(`artifact.unreadable`, `main.go:41-43`, `main.go::newPartialArtifacts`). That
is a declaration, not a mechanism. The RDR gives no criterion for what makes a
key unreadable rather than absent at a real artifact boundary — for a TOML file,
a parse error is file-wide, not key-wide; for an HTTP accessor, a partial
response is a transport fact with no per-key granularity. Absent a criterion,
"the accessor could not read this specific key" may not be an observable state
for any accessor the project actually ships.

**Test it prevents**: MVV Scenario 2's truncation case against a production
accessor rather than a fixture double. The spike case is testable only because
the fixture asserts the condition it is meant to detect. Downgraded from HIGH
because the MVV explicitly scopes itself to "a fixture flow", so the fixture
double may be the intended deliverable — but the Testing Strategy says "those
spike cases become package tests", which does not close the gap for the real
binding interface the Risks section says Resolve must inspect.

---

## LOW

### L1. Disposition Table's completeness rows are witnessed by a spike whose own default path cannot produce the refusal

**Passage**: Disposition Table, Witness column for rows **All requested keys
read** (`output.txt:11`) and **Requested key unreadable** (`output.txt:12`).

Both witnesses come from the explicit-key call sites (`main.go:361-362`). Every
*other* read in the transcript — line 1 and both replay runs — goes through the
derived-key default path where `incomplete_read` is unreachable. A reviewer
reading the Witness column sees the table fully witnessed and does not see that
the majority of the spike's own read invocations run in a mode where the
contract's central refusal cannot fire.

**Test it prevents**: nothing outright — it weakens confidence in the witnesses
rather than removing an oracle. Recorded because it is the visible symptom of
H1, and because the replay scenario (A4, MVV Scenario 4) inherits the
self-fulfilling mode, so "replay stability" is proven over a read path that
never exercises completeness.

### L2. Failure Modes describes the truncation shortfall in terms the completeness rule does not cover

**Passage**: Failure Modes, second silent-failure shape ("A read returns fewer
keys than requested and the shortfall reads downstream as genuine absence — the
completeness requirement turns that into an `incomplete_read` refusal").

The guard described covers a *short* return. The completeness rule as written
("MUST return the complete tag set … or take the refusal branch") also permits
the inverse defect — a read returning keys that were *not* requested — with no
clause forbidding it. The Fidelity Table calls the returned set "total over the
requested keys", which asserts a lower bound but not an upper one. A read that
over-returns pollutes `assemble`'s owned view (`internal/resolve/resolve.go::assemble`),
where owned provenance outranks observed, so an unrequested key can shadow
caller context.

**Test it prevents**: an "exactly the requested key set" assertion in MVV
Scenario 2. `assert set(tags) >= requested` is supported; `assert set(tags) ==
requested` is not, and the RDR does not say which is intended. Low severity
because the failure direction is not the one A8 was reopened to close.

---

## Not findings

- Oracle Discriminability scenario 2's negative control is adequate: the
  over-refusal direction is covered by the genuine-absence fixture named in the
  "Fails if X" column (`tags={profile=<absent>,…}`), so an implementation that
  refuses on absence fails the row.
- The A8 Evidence line's grounding against
  `internal/resolve/resolve.go::missingOwned` / `TagSet.has` is accurate: both
  symbols resolve on the branch and `missingOwned` does decide on map presence
  alone (`resolve.go:444-459`), so the "kernel cannot recover the distinction
  downstream" claim holds and correctly locates the rule at the accessor
  boundary.
- The Desk Trace's "No CONTRADICTION row" reasoning is sound: completeness
  constrains which branch the existing disjunction takes and adds no third
  branch.
- Proportionality's "same contract" claim is defensible as scoping; M2 is a gap
  in *coverage* of that claim, not a refutation of it.
