# REQ List — RDR 0001 Resolution Kernel Contract

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0001-resolution-kernel.md`. Quotes are verbatim from the RDR.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts
- `LBD` = Proposed Solution / Technical Design / Load-Bearing Decisions
- `FM` = Trade-offs / Failure Modes
- `TS` = Validation / Testing Strategy
- `CA` = Research Findings / Critical Assumptions
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (phases)
- `CC` = Finalization Gate / Cross-Cutting Concerns
- `PE` = Validation / Performance Expectations

---

## Determinism and disposition shape

- [REQ-1] "Given the same flow identity, transition table revision, accessor-produced owned tag snapshot, caller-supplied observed tags, and freshly recognized outcome tag, resolve returns the same disposition: exactly one transition plan or exactly one typed refusal." — (NC)

- [REQ-2] "a resolution input is the tuple of flow identity, transition table revision, accessor-produced owned tag snapshot, observed tag-set, and freshly recognized outcome tag. Replaying that tuple must replay the disposition." — (LBD, Identity)

- [REQ-3] "this RDR claims value-level replay determinism, not byte-identical output, hashes, or serialized canonical form. Determinism is covered by A1 and the MVV replay test over the explicit input tuple." — (CC, Canonical-form / determinism)

- [REQ-4] "identical table, owned snapshot, observed tags, and recognized outcome must return value-identical transition plans without hidden reads." — (CA, A1 Evidence)

## Refusal obligations

- [REQ-5] "The kernel MUST refuse instead of guessing when no edge matches, more than one edge matches, required owned state is unavailable, a guard cannot be evaluated, or the recognized outcome is not modeled by the table." — (NC)

- [REQ-6] "Modeled refusal is a value-level resolver disposition, not a CLI error and not the Go error path for parser bugs, IO failures, or programmer mistakes." — (NC)

- [REQ-7] "The kernel-owned refusal kind set is exactly:" `no_match`, `ambiguous_match`, `owned_state_unavailable`, `guard_unevaluable`, `unmodeled_outcome` — (NC)

- [REQ-8] "Each refusal kind must be stable enough for RDR 0005 to map to a CLI error code without inspecting error strings." — (NC)

- [REQ-9] "Visible failures are typed refusal values: `no_match`, `ambiguous_match`, `owned_state_unavailable`, `guard_unevaluable`, or `unmodeled_outcome`." — (FM)

- [REQ-10] "Diagnosis starts with the input tuple, the refusal kind, and the transition table revision used for that resolution." — (FM)

## Prohibitions (negative contract)

- [REQ-11] "The kernel MUST NOT print output, inspect CLI flags, discover ambient state, choose artifacts on behalf of the caller, initiate work, or execute persistence side effects directly." — (NC)

- [REQ-12] "Silent failure would mean the kernel guessed a transition or executed persistence directly; both are prohibited by the normative contract." — (FM)

- [REQ-13] "The kernel does not discover artifacts, execute accessors, discover work, run subflows, print output, or own persistence." — (Proposed Solution / Approach)

- [REQ-14] "The resolver belongs in an internal package behind the CLI: it must stay stateless and non-orchestrating, consume owned-state snapshots and write targets produced by the accessor layer, and refuse illegal or incomplete transition inputs rather than initiating work." — (Context / Technical Environment)

## Selection semantics

- [REQ-15] "the only successful selection is exactly one matching edge after guard evaluation. Zero, multiple, unavailable, or unevaluable candidates are refusals unless the table contains a modeled escape edge that itself matches exactly once." — (LBD, Selection / predicate)

- [REQ-16] "Evaluation builds a single tag-set view, selects matching candidate edges, refuses zero or multiple matches unless the table contract explicitly models an escape edge, and emits a transition plan." — (Technical Design)

- [REQ-17] "Merge owned, observed, and freshly recognized tags into the evaluation view." — (Illustrative Code, step 2)

## Success plan shape

- [REQ-18] "A successful resolution returns the next state tags and the owned-tag writes that the accessor layer applies back to the same caller-provided artifact boundary; an illegal, ambiguous, incomplete, or unmodeled input returns an explicit refusal class." — (Proposed Solution / Approach)

- [REQ-19] "The kernel may describe owned-tag writes in the returned plan, and the accessor layer knows how to apply those writes to caller-provided artifacts, but RDR 0004 owns how reads, writes, and read-back verification execute." — (Technical Design)

- [REQ-20] "Inputs are: flow identity, transition table revision, owned-state snapshot values produced by the accessor layer, observed tags supplied by the caller, the freshly recognized outcome tag, and the reviewable transition table." — (Technical Design)

- [REQ-21] "the kernel only exposes structured success/refusal values that the CLI can map." — (Technical Design)

## Package boundary and scope

- [REQ-22] "Define the internal resolver package boundary, result taxonomy, and pure resolution entry point without CLI output or persistence side effects." — (IP, Phase 1)

- [REQ-23] "Implement tag-set assembly, guard evaluation delegation, exact-one edge selection, and typed refusal behavior." — (IP, Phase 2)

- [REQ-24] "Expose only the kernel values needed by RDR 0005; do not add command output or state mutation semantics in this RDR." — (IP, Phase 4)

- [REQ-25] "the internal component is the resolver kernel. Rejected names: \"orchestrator\" because it implies initiating work, and \"state machine runner\" because it implies owning persistence." — (LBD, Naming)

- [REQ-26] "No third-party dependency is proposed at this stage." — (IP, New Dependencies)

- [REQ-27] "the resolver is stateless and receives snapshots and table data as inputs; it does not own shared persistence or execute accessor writes." — (CC, Concurrency model)

## Testing strategy scenarios

- [REQ-28] "Implementation must add focused resolver-kernel tests before any CLI command wraps the kernel." — (TS)

- [REQ-29] "**Scenario**: Replay the same legal input tuple twice. **Expected**: Both calls return value-identical transition plans, including next tags and accessor write descriptions." — (TS, 1)

- [REQ-30] "**Scenario**: No table edge matches the assembled tag-set. **Expected**: The kernel returns the typed no-match refusal and performs no persistence side effect." — (TS, 2)

- [REQ-31] "**Scenario**: More than one table edge matches after guard evaluation. **Expected**: The kernel returns the typed ambiguous-match refusal unless one explicit escape edge matches exactly once." — (TS, 3)

- [REQ-32] "**Scenario**: The freshly recognized outcome is not modeled by the table. **Expected**: The kernel returns the typed unmodeled-outcome refusal and performs no persistence side effect." — (TS, 4)

- [REQ-33] "**Scenario**: Owned state is unavailable or a guard cannot be evaluated. **Expected**: The kernel returns the corresponding value-level typed refusal and does not fall back to ambient discovery or the CLI/Go error path." — (TS, 5)

- [REQ-34] "The tests exercise the pure decision boundary: fixture transition tables, owned snapshots from accessor fixtures, observed tags supplied by the caller, and freshly recognized outcome tags." — (TS)

- [REQ-35] "Add focused kernel tests for deterministic replay and refusal classes, using fixtures that exercise owned, observed, and freshly recognized tags." — (IP, Phase 3)

## Non-goals asserted as contract

- [REQ-36] "No encode/decode, import/export, or inverse operation is introduced by this RDR." — (Round-Trip / Inverse Invariants)

- [REQ-37] "No throughput target is part of this RDR. … No byte-stable hash or canonical serialization is introduced" — (PE)

---

## REQ-MVV

- [REQ-MVV] "Resolve must name and implementation must add a replay test that feeds the same table, owned snapshot, observed tags, and recognized outcome to the kernel twice and asserts value-identical dispositions. The same validation must include at least one value-level refusal each for `no_match`, `ambiguous_match`, `owned_state_unavailable`, `guard_unevaluable`, and `unmodeled_outcome`, and must assert those modeled refusals do not use the CLI or Go error path." — (MVV; restated in Finalization Gate / Scope Verification)

REQ-MVV is the gating validation. It decomposes into three obligations that must
all be present in one test suite:

1. a replay test asserting value-identical dispositions across two calls on the
   same input tuple (REQ-1, REQ-2, REQ-4, REQ-29);
2. at least one value-level refusal case for each of the five kinds in REQ-7;
3. an assertion that each modeled refusal travels the value path, not the CLI
   path and not the Go `error` path (REQ-6).

---

## ASSUMPTIONS

- ASSUMPTION: "resolve" in REQ-1 names the kernel's single pure entry point (a
  Go function/method), not the CLI verb `intrastate resolve`. Grounded: the
  Illustrative Code block is labelled "Illustrative CLI flow only; RDR 0005 owns
  the final command syntax", and the critique premortem
  (`evidence/critique/critique.md`) writes the acceptance tests as `Call
  `Resolve`` against a fixture table with no CLI involved.

- ASSUMPTION: "typed refusal" is a value returned alongside (or in place of) the
  transition plan on the **non-error** return path; the Go `error` return for a
  modeled refusal is nil/absent. Grounded: REQ-6 excludes "the Go error path",
  and `evidence/critique/critique.md` acceptance test 1 step 4 says "Assert the
  Go error return is nil or absent for the modeled refusal path". The Go `error`
  path remains available for "parser bugs, IO failures, or programmer mistakes"
  (REQ-6).

- ASSUMPTION: the refusal kind set in REQ-7 is **closed** — implementation adds
  no sixth kind. Grounded: NC says "is exactly", and A5 records the closed
  taxonomy as an accepted Design Decision whose rejected alternative is "letting
  RDR 0005 infer CLI codes from Go errors or ad hoc strings".

- ASSUMPTION: the refusal kind is exposed as a distinct comparable value (a
  string-backed or int-backed named type with the five constants), not a free
  string field and not an `error` sentinel. Grounded: REQ-8 requires mapping
  "without inspecting error strings"; the critique premortem names the missing
  artifact as "a value-level `RefusalKind`".

- ASSUMPTION: "transition table revision" (REQ-2, REQ-20) is an opaque
  caller-supplied identity value that the kernel carries and compares but does
  not parse or interpret. Grounded: RDR 0002 owns table shape and declares
  `[model]` with `id` and `version`; this RDR only requires the revision to be
  part of the input tuple and available for diagnosis (REQ-10).

- ASSUMPTION: the "modeled escape edge" of REQ-15/REQ-16 is table data the
  kernel reads, not kernel-side priority logic, and it can only rescue the
  `no_match` and `ambiguous_match` conditions. Grounded: this RDR's Risks
  section says "any priority or escape must be explicit table data owned by RDR
  0002", and RDR 0002's normative block says an `escape` list "MUST contain only
  resolver failure classes that RDR 0001 allows the table to model: `no_match`
  and `ambiguous_match`" and that a modeled escape disposition is available
  "only when exactly one escape row matches the same tag-set and its `escape`
  list contains the corresponding failure class". `owned_state_unavailable`,
  `guard_unevaluable`, and `unmodeled_outcome` are therefore never escapable.

- ASSUMPTION: guard evaluation is *delegated* (REQ-23 says "guard evaluation
  delegation") — the kernel calls into a guard-evaluation seam owned by RDR 0003
  rather than implementing operator semantics. For this RDR's tests that seam is
  satisfied by a fixture/injectable evaluator, since RDR 0003 is "Deferred to
  peer" in Capability Dependencies. `guard_unevaluable` is the kernel's typed
  response when that seam reports a predicate it cannot decide.

- ASSUMPTION: `owned_state_unavailable` is raised when a candidate row's
  evaluation requires an owned tag absent from the accessor-produced owned
  snapshot — not when a candidate merely fails to match on a present tag.
  Grounded: `evidence/critique/critique.md` acceptance test 4 builds "a fixture
  table whose selected candidate requires an owned tag missing from the
  accessor-produced owned snapshot" and expects kind `owned_state_unavailable`.

- ASSUMPTION: `unmodeled_outcome` is decided against the table's declared
  recognized-outcome alphabet, and takes precedence over `no_match` when the
  supplied recognized outcome is outside that alphabet. Grounded:
  `evidence/critique/critique.md` acceptance test 3 builds "a fixture table
  whose recognized-outcome alphabet excludes the supplied recognized outcome"
  and expects `unmodeled_outcome`, distinguishing it from a mere zero-match;
  RDR 0002 lists "missing recognized outcome alphabet" as a table-level concept.

- ASSUMPTION: "performs no persistence side effect" (REQ-30, REQ-32) is
  verified structurally — a refusal disposition carries no write descriptions,
  and the kernel has no persistence collaborator to call. Grounded:
  `evidence/critique/critique.md` acceptance test 3 step 4, "Assert no
  persistence write description is present", and REQ-11's prohibition on
  executing persistence side effects directly.

- ASSUMPTION: the ambiguous-match refusal carries enough source identity for
  RDR 0005 to report the conflicting rows. Grounded:
  `evidence/critique/critique.md` acceptance test 2 step 4, "Assert the result
  carries enough source identity for RDR 0005 to map and report the ambiguous
  rows"; RDR 0002 requires each normalized candidate row to "retain its source
  rule id and source locator". This is diagnostic payload, not part of the
  refusal-kind discriminator.

- ASSUMPTION: "value-identical" (REQ-1, REQ-4, REQ-29) means Go value equality
  over the disposition — comparable via `reflect.DeepEqual` or an equivalent
  deep compare — not byte-identical serialization. Grounded: REQ-3 explicitly
  disclaims byte-identical output, hashes, and canonical form.

- ASSUMPTION: the new package lives under `internal/` (REQ-14, REQ-22) and does
  not import `internal/cli`, `internal/cli/respond`, or `internal/cli/clierr`;
  the dependency direction is CLI → kernel only. Grounded: REQ-11 and REQ-24,
  the Existing Infrastructure Audit row "Refusal mapping … RDR 0005 maps kernel
  refusal kinds to CLI error codes outside the kernel package", and the audit
  row "Resolver implementation | none found under `internal/` | … | Introduce |
  New internal package can own pure resolution logic." The exact package name is
  unconstrained by the RDR.

- ASSUMPTION: the kernel consumes an already-normalized candidate-row view of
  the transition table rather than parsing TOML. Grounded: Technical Design says
  the kernel receives "the reviewable transition table" and Capability
  Dependencies says "Kernel assumes a parsed, reviewable table shape"; RDR 0002
  owns normalization into candidate rows. For this RDR's tests the table is a
  fixture value satisfying that shape.

- ASSUMPTION: tag provenance (owned / observed / recognized) is carried on
  inputs so the kernel can distinguish them when assembling the evaluation view
  (REQ-17) and when describing owned-tag-only writes (REQ-19). Grounded:
  Background says "tag provenance distinguished as owned, observed, or freshly
  recognized"; RDR 0002 requires the model to "declare every tag it matches or
  writes, including each tag's provenance"; RDR 0004 forbids write accessors
  from writing observed or recognized tags.

---

## QUESTIONS

None.

Two clauses were examined as ambiguity candidates and both resolved against the
RDR's own text plus the pre-lock evidence, so neither is recorded as a QUESTION:

1. **Escape-edge scope** — whether a modeled escape edge can rescue all five
   refusal conditions or only zero/multiple match. REQ-15 lists "Zero, multiple,
   unavailable, or unevaluable candidates are refusals unless the table contains
   a modeled escape edge", which reads as though escapes might cover
   `owned_state_unavailable` and `guard_unevaluable`. RDR 0002's normative block
   closes this: an `escape` list "MUST contain only … `no_match` and
   `ambiguous_match`", and RDR 0001's Risks section assigns escape data
   ownership to RDR 0002. Resolved as an ASSUMPTION above.

2. **Refusal precedence** — which kind wins when more than one condition holds
   (e.g. an unmodeled outcome *and* a missing owned tag). The RDR does not state
   a total ordering, but it does not need one for testability: the MVV requires
   one case per kind, and the critique's acceptance tests each construct a
   fixture isolating a single condition. Implementation must pick a
   deterministic evaluation order (REQ-1 requires the same input to yield the
   same disposition), and the `unmodeled_outcome`-before-`no_match` precedence
   is grounded above. Remaining orderings are implementation-internal and
   invisible to the REQ set as long as they are deterministic.
