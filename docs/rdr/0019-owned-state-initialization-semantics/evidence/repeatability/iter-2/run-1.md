model: claude-fable-5-1
variant: full (profile: foundational)
scope: delta re-run (iter-2)

Anchors read (projector only): 0019:C1 (scratch copy), 0019:A6, 0019:A8,
0019:D-wire-byte-format, 0019:D-selection-predicate, 0019:D-naming,
0019:D-identity, 0019:F2, 0019:F5; cross-read for consistency: the
`trace`/`disposition` mini-check tables, Step 2, S4, S6, S7, S8, S9,
S10, MVV, RT1/RT2, Technical Design preamble. Exit groups of the two
read-back arms were spot-checked against `internal/cli/flow_exec.go`
(`accessorFailureOf`: mismatch raised at `GroupUserEnv`, incomplete and
timeout via `envErr`) and the three shipped code spellings in
`internal/cli/flow_input.go`.

## Q1 — read-back MISMATCH vs read-back COULD-NOT-COMPLETE exit codes

DETERMINATE. C1's closing paragraph fixes the split outright: a
read-back that COMPLETED and disagreed is exit 2 (`GroupUserEnv`,
the PRESENT-AND-UNVERIFIED refusal); a read-back that could not
complete (unreachable locator, timeout) is exit 3
(`GroupEnvUnavailable`); both leave a non-empty store so a re-run is a
non-repairing no-op. The `disposition` table carries the same two rows
(2 / 3) and the `authority` table's last row restates it ("mismatch 2,
incomplete/timeout 3"). S8 asserts the mismatch as distinct from the
INCOMPLETE class. The three passages agree, and the split matches the
raise sites in source.

## Q2 — role-binding validation BEFORE or AFTER the emptiness read; C1 / trace / Step 2 agreement

DETERMINATE. All three now say BEFORE, in the same terms:
- C1: "the ROLE-BINDING half of plan validation runs BEFORE the
  emptiness read, so an invocation that leaves a needed role unbound
  refuses at exit 2 REGARDLESS of what the bound stores contain — it
  does not reach the predicate and cannot take the no-op arm"; the
  writer-arity half is explicitly NOT ordered here (settled at load).
- `trace` table: row "role-binding validation — C1 all-or-nothing,
  ordered BEFORE the predicate — unbound → exit 2, no predicate (S6)"
  precedes row "emptiness read".
- Step 2: "Validate that every role a needed writer names is bound;
  unbound → refuse (exit 2) without reading any store. Then read the
  bound stores … The order is C1's, not this step's".
No drift remains among the three.

## Q3 — class refusal vs carrier refusal: which fires first

DETERMINATE, stated outright. C1 has a dedicated paragraph headed
"ORDERING: the CLASS refusal precedes the CARRIER refusal", with the
ground (class is answerable from the loaded model alone; carrier needs
the constructed registry) and the reach (the rule matters only for a
model tripping both arms, which no admitted model does today). The
same paragraph fixes a second order: carrier refusal is decided after
registry construction, before any accessor invocation, and PREEMPTS
the `--allow-commands` refusal. The `trace` table lists "class check"
then "carrier check" in that order; Step 1 places the class refusal
"before any accessor". Consistent.

## Q4 — is `--allow-commands` accepted by this verb; stated or inferred

DETERMINATE, stated. C1 CARRIER paragraph: "It also carries
`--allow-commands`, not by declaring it but by INHERITING it: the flag
is registered once on the `flow` group's persistent set
(`internal/cli/flow.go::newFlowCmd`), so every child verb takes it
structurally and this one cannot decline it", and it forecloses the
closed-flag-list misreading explicitly. S10 ("run with
`--allow-commands` UNSET") and the ORDERING paragraph's preemption
clause both presuppose the flag exists on the verb. No inference is
required.

## Q5 — SEEDED-all arm: values or names; absent-key list on the seeded arm

DETERMINATE. C1 (payload paragraph): "the payload reports seeded keys
by NAME and does not echo their values — `read-state` is the surface
that reports values (RT1)"; and "The absent-key report belongs to the
NO-OP arm, which is the only arm where the set can be non-empty; on
the seeded arm every `[initial]` key was just written, so an absent
list there is necessarily empty and MUST NOT be reported as if it
carried information." Field NAMES remain deferred (D-wire-byte-format,
"owned here"), which is a declared deferral, not silence on content.
The `disposition` table's seeded-all row and MVV step 2 agree.
Residual, within the declared deferral: "MUST NOT be reported as if it
carried information" fixes that the seeded arm carries no absent-key
report, but does not by itself say whether an empty list field may be
present versus omitted; that is a wire-shape question the record
defers with the field names, so I do not count it as silence on the
question asked.

## Q6 — read-back-SEALED store: exit 0 no-op vs exit-3 refusal; obtainability of the count

DETERMINATE. C1: the sealed artifact is "a ONE-key, NON-EMPTY store
that init-state declines to seed" and "is nonetheless a NO-OP SUCCESS
at exit 0 here, not an exit-3 refusal", with the reconciliation stated
in the same passage: exit 3 is what a READ-BACK on a sealed artifact
returns, and this arm performs no write and therefore no read-back;
what it needs is a key COUNT obtained without reading key values. The
obtainability question is answered as a CONSTRAINT on the undecided
carrier, not left implicit: C1 — "whichever carrier lands MUST answer
cardinality on a sealed store rather than inheriting `::Reader.Read`'s
unreadable short-circuit, or this arm degrades from no-op success to
an exit-3 refusal and S9 fails"; A6 — "The candidate MUST also answer
cardinality on a READ-BACK-SEALED store — returning 1, not an
unreadable short-circuit". The `disposition` table's sealed row (0 /
zero / no-op) and S9's Expected agree. The carrier mechanism itself
remains A6 Pending, which is the record's declared open assumption and
was not a divergence in the first pass; the two dispositions no longer
contradict.

## Q7 — empty scalar `[initial]` value: what persists; findable from Normative Contracts

DETERMINATE, findable in C1. C1: "An `[initial]` value that is an
EMPTY SCALAR (`note = ""`) seeds like any other: the scalar arm renders
`members[0]`, the artifact carries `{"note":""}`, and the key reads
back PRESENT — distinct from a cleared key, which leaves the object
entirely. It creates no third encoding and is not refused at seed
time". F5, S4's added row, A2's Evidence, and the `disposition` table's
empty-scalar row all state the same shape; the argv-side asymmetry is
explicitly NOT decided (routed to RDR 0002). C1 and F5 agree; the
answer no longer lives only in Failure Modes.

Tally: 7 DETERMINATE, 0 STILL-SILENT, 0 CONTRADICTORY.

## NET-NEW

1. CONTRADICTION introduced by the C1 carrier rewrite against an
   untouched passage — S8. C1 now fixes "The carrier gate is over BOTH
   capabilities … The refusal therefore fires when EITHER the bound
   write binding or the bound read binding for a needed role is not
   the file-backed type", and S10 gains a read-side row asserting a
   command-backed `[read.x]` behind a file-backed writer refuses at
   exit 2 with the carrier code. S8 still says, of the read-back
   mismatch fixture: "(A command-backed READER is the other legal
   construction — C1's carrier refusal is scoped to WRITE accessors —
   but needs a helper binary and `--allow-commands`.)" Under the
   rewritten C1 that construction is not legal for S8's purpose: it is
   refused at the carrier gate before any read-back runs, so it cannot
   produce a completed-and-disagreeing read-back. S8's parenthetical
   should be deleted or inverted (name it as the construction S10's
   read-side row refuses). The primary S8 construction (reader on a
   different path than its writer) is unaffected.

2. Residual ordering the new ORDERING paragraph does not close:
   carrier check versus role-binding validation. C1 orders class before
   carrier, and role-binding before the emptiness read, and carrier
   before `--allow-commands`, but never carrier relative to
   role-binding. The `trace` table places "carrier check" before
   "role-binding validation"; Step 2 opens with role-binding and omits
   the carrier check entirely (Step 1 mentions only the class refusal).
   Both orders are consistent with every stated MUST, and the case is
   constructible (an S10 fixture invoked with the artifact flag
   omitted refuses under either the carrier code or the
   artifact-binding family depending on which runs first). Not a
   contradiction; an implementer can follow the trace table. Flagged
   because the rewrite fixed three orderings by name and left this
   fourth to a mini-check table alone.

3. D-selection-predicate lags the rewrite, non-contradictory: it
   scopes the predicate "to the file-backed carrier, with the other
   two (`::NewEditWriter`, `cmdbind.Writer`) refused, per C1", naming
   only the two WRITE carriers, while C1 now also refuses
   `cmdbind.Reader` on the read side. The "per C1" cross-reference
   keeps it from contradicting, but the enumeration reads as complete
   and is not.

4. D-wire-byte-format still says field names are "deferred to
   Resolve/Pre-Lock, owned here". The record is at pre-lock and C1's
   payload paragraph re-affirms the deferral; the pointer is stale
   rather than wrong, and no first-pass divergence turned on it.
