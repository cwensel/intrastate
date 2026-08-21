Model: claude-opus-5[1m]

# Persona 2 — Implementer: clarification requests (iter-2, A8 delta scope)

Scope: only the read-completeness material added by the A8 re-entry (A8, the
read-completeness normative clause, the Read completeness Load-Bearing
Decision, MVV / Validation Scenario 2, the read-truncation Risk and Failure
Modes paragraph, the Proportionality read-completeness sentence, and the
read-completeness rows of the Disposition Table, Oracle Discriminability,
Fidelity Table, and Desk Trace). A1–A7 carried forward.

Each finding is a question I would ask in hour one, the passage that triggers
it, and the decision or test it blocks.

---

## HIGH

### H1 — Who supplies the requested key set, and is it a validated field?

**Passage**: Normative Contracts, read-completeness clause — "MUST return the
complete tag set for **the keys it was asked for**"; and Load-Bearing
Decisions → *Read completeness* — "the success predicate ... is 'every
requested key resolved'".

**Question**: The whole contract is quantified over "the requested keys", but
no clause in this RDR says where that set comes from or who owns it. Three
candidates are live in the draft and they are not the same set:

1. the accessor definition's declared **expected tag keys** (Technical Design,
   phase-listing paragraph: "Each definition declares a stable name,
   capability, artifact role, **expected tag keys**, timeout policy, ...");
2. a per-invocation set the caller/resolver passes (the spike's variadic
   `requested ...string` in `evidence/spikes/main.go::read`);
3. the union of `Row.RequiresOwned` over candidate rows, which is what
   `internal/resolve/resolve.go::missingOwned` actually quantifies over
   downstream.

If it is (1), the set is static per accessor and completeness is checkable at
validation time. If it is (2), it is dynamic and the accessor signature needs a
key-set parameter. If it is (3), this RDR has an undeclared dependency on RDR
0002's normalized row shape for an input it treats as its own.

**Decision blocked**: the read accessor's Go signature and the accessor
definition struct in **Phase 1: Accessor Model**. I cannot define
`type ReadAccessor` or the definition record without knowing whether the key
set is a definition field, a call parameter, or derived from the table. Every
downstream completeness assertion is quantified over a set I cannot name.

### H2 — The spike's default key set makes truncation undetectable, and the RDR does not forbid it

**Passage**: A8 Evidence — "resolves reads against an explicit requested key
set"; Oracle Discriminability row 2 negative control — "`output.txt:11`
complete read".

**Question**: The spike only honors an explicit key set when the caller passes
one. When `requested` is empty it falls back to
`evidence/spikes/main.go::expectedTagKeys(art)` — the keys the artifact itself
turned out to carry. That default is circular: a read that silently lost keys
would have a *smaller* derived request set and would still report success. The
transcript's own line 1 (`output.txt:1`, from `read(reg, artifacts,
"state.read")` with no requested keys) is produced by exactly that path.

The RDR never states that an empty or absent requested key set is a validation
failure, so an implementer copying the spike's shape inherits the hole the
clause exists to close.

**Decision blocked**: whether **Phase 1** validation must reject a read
definition with no declared key set, and whether "requested keys derived from
what was read" is a conformance violation. Also blocks Validation Scenario 1's
expectation list, which enumerates seven validation failures and does not
include a missing read key set. If the answer is "the key set is mandatory
metadata", that is an eighth validation arm and an eighth transcript line the
Disposition Table's "Definition invalid (7 shapes)" row does not carry.

### H3 — How does an absent value cross into `resolve.Input.Owned`?

**Passage**: Load-Bearing Decisions → *Read completeness* — "The contract
constrains which branch is taken, not how an absent value is represented in the
success branch; the Resolve spike's `<absent>` sentinel is fixture shorthand,
not a normative representation"; and Fidelity Table, `read` over a requested
key set row — "how an absent value is *represented* is not pinned".

**Question**: The RDR deliberately leaves absence representation open, but the
consumer is already built and has no slot for it. `resolve.Input.Owned` is
`[]Tag` (`internal/resolve/resolve.go`, `type Input`), and `assemble` puts every
`Owned` tag into the view as a present key with `ProvenanceOwned`. So the
accessor layer has exactly two encodings available today, and they behave
oppositely:

- **omit the key** from `Owned` → `TagSet.has` is false → `missingOwned`
  reports it → `owned_state_unavailable`;
- **include it with a sentinel value** → `TagSet.has` is true → the key is
  *present* with a garbage value, and `TagSet.matches` may match a row against
  the sentinel.

A8's own Evidence names this: `missingOwned` "decides on `TagSet.has` (map
presence) alone". But that observation is used only to argue the rule belongs
at the accessor boundary; it never says which encoding the accessor layer must
emit. Leaving it open is not neutral — the two options produce different
resolver dispositions for the same artifact.

**Decision blocked**: the accessor→kernel adapter in **Phase 4: CLI
Integration Hook** (and anything that builds `resolve.Input`). I cannot write
the conversion without choosing. Note JDR 0001 §D3 flagged option (c), an error
channel on `Input.Owned`, as "available if implementation shows the kernel
needs to distinguish read-failed from absent" — but the *absent vs. present*
question here is a third case §D3 did not resolve, and this RDR does not
inherit an answer for it.

### H4 — Is `incomplete_read` recoverable, and does it collapse with `execution_failure` at the CLI?

**Passage**: Disposition Table rows "Requested key unreadable → refusal →
`incomplete_read`" and "Artifact role not supplied → refusal →
`execution_failure`"; Failure Modes — "the completeness requirement turns that
into an `incomplete_read` refusal".

**Question**: `incomplete_read` is minted as its own class, but this RDR never
says what distinguishes it from `execution_failure` at the boundary. In the
spike they differ only by *where* the error arises: a missing artifact role
returns `refusalExecutionFailure` (`main.go::read`, the `artifacts[a.role]`
miss), a per-key read error returns `refusalIncompleteRead`. Both are "the
accessor could not read". Two sub-questions:

1. Does the refusal carry **which** keys were unreadable? The spike returns the
   bare class with no key list (`output.txt:12` is just
   `refusal=incomplete_read`), while the sibling `read_back_mismatch` carries
   `expected=` / `observed=` payloads. Failure Modes says "Diagnosis starts
   with the accessor identity, capability, artifact role, timeout, and expected
   versus observed tag values" — for `incomplete_read` there are no observed
   values, so that sentence does not describe this class.
2. Do `incomplete_read` and `execution_failure` map to the same exit group? The
   RDR defers exit codes to RDR 0005, but JDR 0001 §JD-8 already records the
   live constraint that `GroupUserEnv` and `GroupInternal` both exit 2 and
   "an accessor that could not read the artifact is what `GroupEnvUnavailable`
   already means". If both classes land on the same exit, the separate class
   buys diagnosis only, not remedy.

**Decision blocked**: the `Refusal` payload struct for read refusals in
**Phase 2: Executor Boundary**, and Validation Scenario 2's assertion strength
— the scenario asserts the classes "remain distinct" but never asserts the
refusal names the offending keys, so a test can pass while the operator gets an
undiagnosable failure.

---

## MEDIUM

### M1 — What counts as "could not read" for a per-key failure?

**Passage**: Load-Bearing Decisions → *Read completeness* — "A key the artifact
genuinely does not carry resolves as an absent value; a key the accessor could
not read is an `incomplete_read` refusal."

**Question**: The distinction is stated as if it were self-evident, but for a
real accessor it is a classification the implementer has to make. If the read
accessor is fronting a TOML file, "genuinely does not carry" and "could not
read" are the same syscall outcome for a malformed file, a permission error, a
truncated file, or a key whose value fails type coercion. The spike sidesteps
this entirely: `artifact.unreadable` is a hand-populated `map[string]bool`
fixture field (`main.go::newPartialArtifacts`), so the classification is given,
never derived.

Concretely: an artifact parses, the requested key is syntactically present, but
its value is not a valid typed tag value. Absent value, `incomplete_read`, or
`execution_failure`?

**Decision blocked**: the error-classification switch inside every real read
binding in **Phase 2**. Also weakens the "Verified" standing of A8's Method:
Spike — the spike proves the *branch machinery* discriminates when told which
keys are unreadable; it does not exercise a binding that has to decide that
itself.

### M2 — Whole-read refusal vs. per-key granularity is asserted, not argued

**Passage**: Normative Contracts read-completeness clause — "A partial or
truncated read is a refusal, not a value"; Technical Design — "A read that
cannot produce every requested key takes the refusal branch."

**Question**: One unreadable key refuses the *entire* read, discarding keys
that read fine. For an accessor asked for ten keys where one is unreadable,
nine good values are thrown away and the resolver gets nothing. That may well
be right — it is the conservative choice and it matches P2 in JDR 0001 ("the
kernel refuses rather than guesses") — but the RDR never states it as a
considered decision, so I cannot tell whether all-or-nothing is intended or is
an artifact of the spike's early `return` inside the key loop
(`main.go::read`).

This matters because it interacts with H1: if the requested key set is the
union over *all* candidate rows (option 3 there), one unreadable key that no
surviving row actually requires would refuse a resolution that could have
succeeded.

**Decision blocked**: whether the executor short-circuits on first unreadable
key or accumulates, and whether the refusal is scoped to the accessor or to the
resolution. Affects the Round-Trip section's silence on read granularity and
the MVV's single-truncation fixture, which cannot distinguish the two designs.

### M3 — Does read completeness apply to the write read-back re-read?

**Passage**: Proportionality — "Read completeness is part of that same contract
— it is the success predicate of the read capability, not a separate
obligation"; Desk Trace step 5 ("Read back the same role") vs. step 2 ("Invoke
read accessor").

**Question**: The read-back in step 5 is a read. Is it governed by the
completeness clause? The Desk Trace lists completeness under step 2's
assertions and *not* under step 5's, which reads as "no" — but Proportionality
says completeness is the success predicate of the read capability as such,
which reads as "yes". The spike does not settle it: `main.go::write` re-reads
by cloning the in-memory map, never going through `read`, so no key can be
unreadable there.

If a read-back's re-read hits an unreadable key, the disposition is either
`incomplete_read` or `read_back_mismatch`, and those have very different
meanings for an operator: one says "I could not verify", the other says "the
artifact is wrong". The Disposition Table has rows for both write mismatch
shapes but no row for "read-back re-read was itself incomplete".

**Decision blocked**: the read-back path in **Phase 3: Write Read-Back
Verification**, and a missing Disposition Table row. Blocks Validation Scenario
3, which enumerates only matching / mismatched-owned / mutated-non-owned
outcomes.

### M4 — Oracle Discriminability scenario 2's negative control does not control for the stated failure

**Passage**: Oracle Discriminability, row "2 — read / gate dispositions",
Negative control column — "`output.txt:11` complete read — an implementation
that refuses whenever a key is interesting fails this row".

**Question**: The row's "fails if X is wrong" is that *completeness collapses*
— an implementation that thins the value set instead of refusing. The named
negative control guards a different failure: an implementation that over-
refuses. Neither the positive case nor that control catches the specific bug I
am most likely to write, which is the H2 bug: deriving the requested key set
from what was successfully read. Such an implementation returns
`tags={status=Draft}` for the truncation fixture — a *success* with one key —
and passes both the complete-read control and any assertion that merely checks
"did not refuse the sparse fixture".

The spike does assert against this (`main.go` main: `truncatedRead.refusal !=
refusalIncompleteRead || len(truncatedRead.tags) != 0`), but that assertion is
in the spike's panic block and is not reflected in the RDR's MVV text.

**Decision blocked**: the assertion set for the read test in the MVV. To be
discriminating the test must assert the truncation fixture returns the refusal
branch **with an empty value set**, and must pin the requested key set
independently of the artifact's contents. Neither is currently normative.

---

## LOW

### L1 — "typed absence marker" is named but not specified anywhere in the plan

**Passage**: Load-Bearing Decisions → *Read completeness* — "implementation may
choose a typed absence marker instead."

**Question**: Phase 1 enumerates "accessor definition structs, capability enum,
refusal classes, and validation rules" — no absence type. If a typed absence
marker is a live option it is a Phase 1 type decision, and per H3 it is the
thing that determines the kernel's disposition. As written it reads as an aside
rather than a deliverable, so it will not get built in Phase 1 and will be
improvised in Phase 4.

**Decision blocked**: Phase 1's type inventory. Minor because H3 is the real
question; this is where the answer has to land.

### L2 — Disposition Table's absent-key row says "loud (value returned)" but the value is unspecified

**Passage**: Disposition Table, row "Requested key absent from artifact →
success → — (absent value) → loud (value returned)".

**Question**: The table's Silent-or-loud column is the RDR's guarantee that no
input class exits quietly. For this row, loudness depends entirely on the
absence representation left open in H3/L1 — if absence is encoded as key
omission, nothing is returned for that key and the row is *silent* by the
table's own definition, contradicting the closing sentence "No input class
exits silently: every row either returns a value or mints a named refusal."

**Decision blocked**: nothing on its own, but the table's silence guarantee is
only true under one of the two H3 encodings. Resolving H3 resolves this row;
if H3 lands on omission, this row and that closing sentence need to change.

### L3 — Validation Scenario 2 mixes read and gate concerns in one scenario with one witness line each

**Passage**: Validation → Scenario 2 — "Invoke read and gate accessors that
succeed, time out, return an execution failure, return only a subset of the
requested tag keys, or return gate indeterminate."

**Question**: Scenario 2 now carries five dispositions across two capabilities.
As a test-plan unit it is one scenario name covering read-complete,
read-incomplete, read-absent, timeout, execution failure, and gate
indeterminate. When it fails I get one failing scenario and have to bisect.
Every other scenario in the list is single-purpose.

**Decision blocked**: test decomposition only — cosmetic to the contract, but
it is the first thing I would split on Monday, and if the RDR's scenario
numbering is load-bearing for traceability (the Oracle Discriminability table
keys off "MVV scenario" numbers) splitting it silently breaks that mapping.
Worth confirming the numbering is not a fixed contract.
