Model: claude-opus-5[1m]

# 3amigo Consolidation — iter-2 (re-entry, delta-scoped to A8)

Three isolated persona passes (PM, Implementer, QA), each seeing only the RDR
and its own persona block. Consolidation is mechanical: passages named by two or
more personas are tagged **hotspot**. Overlap marks a hotspot passage, not a
validated finding; a single-persona finding is not thereby weaker.

Scope: the A8 re-entry material only — A8, the read-completeness normative
clause, the Read completeness Load-Bearing Decision, MVV / Validation Scenario 2,
the read-truncation Risk and Failure Modes paragraph, the Proportionality
sentence, and the read-completeness rows of the Disposition Table, Oracle
Discriminability, Fidelity Table, and Desk Trace. A1–A7 carried forward.

## Origin ledger

| ID | Passage | Personas | Hotspot |
| --- | --- | --- | --- |
| T-1 | Load-Bearing Decisions → Read completeness, final sentence (absence representation unpinned) + Fidelity Table `read` row Lossy-exemptions | PM H1/H2, Implementer H3/L1/L2, QA H3 | **hotspot (3)** |
| T-2 | Normative Contracts read-completeness clause — origin of "the keys it was asked for" | Implementer H1/H2, QA H1 | **hotspot (2)** |
| T-3 | Normative Contracts read-completeness clause — `incomplete_read` never minted normatively / not required distinct from `execution_failure` | Implementer H4, QA H2 | **hotspot (2)** |
| T-4 | Proportionality vs. Desk Trace step 5 — does read completeness govern the write read-back re-read? | Implementer M3, QA M2 | **hotspot (2)** |
| T-5 | Load-Bearing Decisions → Read completeness — no observable criterion for "could not read" vs. "genuinely does not carry" | Implementer M1, QA M3 | **hotspot (2)** |
| T-6 | MVV / Validation Scenario 2 — assertions stop at branch; not discriminating against the derived-key-set bug | PM M1, Implementer M4 | **hotspot (2)** |
| T-7 | Failure Modes / read-truncation Risk — inventory names one silent shape where the delta created two | PM M2 | single (PM) |
| T-8 | Technical Design phase two / Disposition Table — no precedence between `timeout` and `incomplete_read` mid-key-set | QA M1 | single (QA) |
| T-9 | Normative read-completeness clause / Technical Design — whole-read refusal vs. per-key granularity asserted, not argued | Implementer M2 | single (Implementer) |
| T-10 | Fidelity Table `read` row / Failure Modes — totality is a lower bound only; over-return unforbidden | QA L2 | single (QA) |
| T-11 | Proportionality — "same contract" covers the branch rule, reads as covering representation | PM L1 | single (PM) |
| T-12 | Desk Trace step 2 closing paragraph — no-CONTRADICTION argued against the wrong pair | PM L2 | single (PM) |
| T-13 | Disposition Table Witness column / A4 replay — witnesses come from the explicit-key path; default path cannot fire the refusal | QA L1 | single (QA) |
| T-14 | Validation Scenario 2 — five dispositions across two capabilities under one number the Oracle table keys off | Implementer L3 | single (Implementer) |

## Hotspots, in severity order

### T-1 — absence representation is unpinned, and the kernel has no slot for it (3 personas)

All three personas independently reached the same passage from different
questions. The RDR states that a genuinely-absent requested key resolves as a
value in the success branch, then explicitly declines to pin how that value is
represented ("fixture shorthand, not a normative representation … implementation
may choose a typed absence marker instead"), and the Fidelity Table repeats the
exemption.

The consumer is already built and has exactly two available encodings, which
produce **opposite** resolver dispositions for the same artifact:

- omit the key from `Input.Owned` → `TagSet.has` false → `missingOwned` reports
  it → `owned_state_unavailable`;
- include it with a sentinel value → `TagSet.has` true → the key reads as
  *present* with a garbage value, and a guard may match against the sentinel.

A8's own Evidence line cites `missingOwned` to argue the kernel cannot recover
the absent/unread distinction downstream — then hands the kernel a
representation under which absence looks answered. PM frames the consequence as
a user-outcome regression (a resolution computed over a key the artifact does
not carry, where `owned_state_unavailable` previously fired); Implementer frames
it as an unwritable accessor→kernel adapter; QA frames it as a removed oracle
(no spec-supported assertion discriminates an absent value from an empty one).

Blocks: the Phase 1 type inventory, the Phase 4 adapter, and the genuine-absence
assertion in MVV Scenario 2.

### T-2 — the requested key set has no defined origin (2 personas)

The whole contract is quantified over "the keys it was asked for", and no clause
says who supplies that set. Three candidate sources are live in the draft and
are not the same set: the definition's declared *expected tag keys* (Technical
Design), a per-invocation caller set (the spike's variadic `requested`), or the
union of `Row.RequiresOwned` (what `missingOwned` quantifies over downstream).
Alternative 1 calls it an "expected tag contract" — a third name for the first.

Both personas independently found the spike's default path is circular: with no
explicit keys, `read` falls back to `expectedTagKeys(art)` — the artifact's own
keys — so the requested set is derived from what was read and `incomplete_read`
can only fire on the hand-set `unreadable` map. Transcript line 1 and both
replay runs use that self-fulfilling mode; only lines 11–13 use the explicit
mode. The RDR never makes a missing key set a validation failure.

Blocks: the read accessor signature and definition struct (Phase 1); whether a
missing key set is an eighth validation arm (Validation Scenario 1 enumerates
seven); and the arrange step of the central MVV truncation assertion.

### T-3 — `incomplete_read` is never minted by a normative clause (2 personas)

The read-completeness clause says only "take the refusal branch" and never names
a class. `incomplete_read` appears only in non-normative carriers (Load-Bearing
Decisions, Disposition Table, Desk Trace, Failure Modes). Contrast the timeout
clause, which is explicit that timeout is "its own refusal class, distinct from
execution failure and read-back mismatch" — there is no parallel sentence here.
An implementation returning `execution_failure` for an unreadable key satisfies
every normative clause in the RDR.

Implementer adds that the refusal carries no key list (`output.txt:12` is the
bare class) while the sibling `read_back_mismatch` carries expected/observed
payloads — so Failure Modes' diagnosis sentence does not describe this class.

Blocks: Validation Scenario 2's "remain distinct refusal classes" assertion
(unfalsifiable against the spec as written) and the read-refusal payload struct
(Phase 2).

### T-4 — does read completeness govern the write read-back re-read? (2 personas)

Desk Trace lists completeness under step 2 and not under step 5, which reads as
"no"; Proportionality says completeness is the success predicate of the read
capability as such, which reads as "yes". The spike settles nothing — `write`
re-reads by cloning the in-memory map and never goes through `read`, so no key
can be unreadable there. If a read-back's re-read hits an unreadable key the
disposition is undefined: `incomplete_read` ("I could not verify") and
`read_back_mismatch` ("the artifact is wrong") mean very different things to an
operator, and the Disposition Table has no row for it — while its closing
sentence claims no input class exits silently.

Blocks: the Phase 3 read-back path; a missing Disposition Table row; Validation
Scenario 3, which enumerates only matching / mismatched-owned / mutated-non-owned.

### T-5 — "could not read" has no observable criterion (2 personas)

The absent/unreadable split is stated as self-evident, but for a real binding it
is a classification the implementer must derive. The spike hand-populates
`artifact.unreadable`, so the condition is declared, never detected. For a TOML
file a parse error is file-wide, not key-wide; for an HTTP accessor a partial
response has no per-key granularity; a syntactically-present key whose value
fails type coercion fits neither branch cleanly.

Both personas note this weakens A8's `Method: Spike` standing in the same way:
the spike proves the branch *machinery* discriminates when told which keys are
unreadable; it does not exercise a binding that must decide that itself.

Blocks: the error-classification switch in every real read binding (Phase 2);
MVV Scenario 2's truncation case against a production accessor.

### T-6 — MVV Scenario 2 asserts the branch and stops (2 personas)

Both read assertions stop at which branch was taken and that the two fixtures
differ. Neither asserts what the absent value does to a consumer, so the MVV
passes unchanged under T-1's regression (PM). And neither is discriminating
against the most likely implementation bug — deriving the requested key set from
what was successfully read, which returns a one-key *success* for the truncation
fixture and passes both arms of the Oracle row (Implementer). The spike does
assert the empty-value-set property in its panic block; the RDR's MVV text does
not carry it.

Blocks: whether MVV Scenario 2 is done at lock — the Prerequisites checkbox
"All Critical Assumptions verified (A1–A8)" is already marked.

## Single-persona findings

Carried at full weight; the personas never saw each other, so absence of overlap
is not evidence against them.

- **T-7** (PM) — Failure Modes says there are "two silent-failure shapes, each
  with a mandatory guard". The delta introduced a third: genuine absence read as
  present state (T-1). It has no named guard, and by the paragraph's own standard
  it should have one or be argued away. The diagnosis sentence does not help a
  user whose plan was computed over a key the artifact never carried.
- **T-8** (QA) — a read that resolves some keys and then exceeds its timeout
  satisfies both the timeout and incomplete-read input classes; no precedence
  rule exists. The spike dodges it by construction (`slow.read` sleeps before the
  key loop). This is the most likely production shape of a truncated read.
- **T-9** (Implementer) — one unreadable key refuses the entire read, discarding
  keys that read fine. Probably right and matches JDR 0001 P2, but never stated
  as a considered decision, so it is indistinguishable from an artifact of the
  spike's early `return` inside the key loop.
- **T-10** (QA) — the Fidelity Table calls the returned set "total over the
  requested keys", a lower bound with no upper bound; no clause forbids
  over-return. An unrequested key entering `assemble` with owned provenance can
  shadow caller-supplied observed context.
- **T-11** (PM) — Proportionality's "same contract" argument is sound for the
  branch rule but reads as though it covered the representation question, which
  is a seam question against RDR 0001's `Input.Owned`. That is what makes T-1's
  escape hatch look costless.
- **T-12** (PM) — Desk Trace's no-CONTRADICTION paragraph examines the
  values/refusal disjunction pair and reasons correctly about it, but does not
  range over the cross-layer collision between step 2's absent-value witness and
  the kernel behavior A8's Evidence cites.
- **T-13** (QA) — the Disposition Table's completeness witnesses both come from
  the explicit-key call sites; every other read in the transcript (line 1, both
  replay runs) uses the derived-key path where the refusal is unreachable. The
  A4 replay scenario therefore proves stability over a read path that never
  exercises completeness.
- **T-14** (Implementer) — Validation Scenario 2 now carries five dispositions
  across two capabilities under one number that the Oracle Discriminability table
  keys off, so splitting it silently breaks that mapping.

## Cross-document note (Implementer, grounding context)

JDR 0001 §D3 resolved read completeness as option (b) and parked option (c) — an
error channel on `Input.Owned` — as "available if implementation shows the kernel
needs to distinguish read-failed from absent". T-1's *absent vs. present*
encoding question is a third case §D3 did not settle, so this RDR inherits no
answer for it.

## Review-gate assessment

**Healthy.** Every finding is anchored to a named passage and names the decision
it blocks or the test it prevents. No persona file references another's output —
isolation held. The three-way convergence on T-1 from three different questions
(user outcome / adapter / oracle) is agreement between contexts that never saw
each other, so it is evidence, not conformity.
