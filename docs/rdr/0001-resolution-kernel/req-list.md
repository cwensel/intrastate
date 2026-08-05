# RDR 0001 Requirement Audit

## Requirements

[REQ-1] "the same input must produce the same legal output" — (Problem Statement)

[REQ-2] "unmodeled matches must be refused instead of guessed" — (Problem Statement)

[REQ-3] "escape behavior must be explicit table data rather than a confident wrong edge" — (Problem Statement)

[REQ-4] "The kernel must treat state as a tag-set, with guards expressed as predicates over tags and tag provenance distinguished as owned, observed, or freshly recognized." — (Context / Background)

[REQ-5] "it must stay stateless and non-orchestrating, consume owned-state snapshots and write targets produced by the accessor layer, and refuse illegal or incomplete transition inputs rather than initiating work." — (Context / Technical Environment)

[REQ-6] "the caller supplies the transition table, the recognized typed outcome, all non-owned context tags, and the owned tag snapshot already read from caller-provided artifacts by the accessor layer." — (Proposed Solution / Approach)

[REQ-7] "A successful resolution returns the next state tags and the owned-tag writes that the accessor layer applies back to the same caller-provided artifact boundary" — (Proposed Solution / Approach)

[REQ-8] "an illegal, ambiguous, incomplete, or unmodeled input returns an explicit refusal class." — (Proposed Solution / Approach)

[REQ-9] "Inputs are: flow identity, transition table revision, owned-state snapshot values produced by the accessor layer, observed tags supplied by the caller, the freshly recognized outcome tag, and the reviewable transition table." — (Proposed Solution / Technical Design)

[REQ-10] "Evaluation builds a single tag-set view, selects matching candidate edges, refuses zero or multiple matches unless the table contract explicitly models an escape edge, and emits a transition plan." — (Proposed Solution / Technical Design)

[REQ-11] "Artifact selection and accessor execution stay outside this RDR's contract." — (Proposed Solution / Technical Design)

[REQ-12] "the kernel only exposes structured success/refusal values that the CLI can map." — (Proposed Solution / Technical Design)

[REQ-13] "Given the same flow identity, transition table revision, accessor-produced owned tag snapshot, caller-supplied observed tags, and freshly recognized outcome tag, resolve returns the same disposition: exactly one transition plan or exactly one typed refusal." — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-14] "The kernel MUST refuse instead of guessing when no edge matches, more than one edge matches, required owned state is unavailable, a guard cannot be evaluated, or the recognized outcome is not modeled by the table." — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-15] "Modeled refusal is a value-level resolver disposition, not a CLI error and not the Go error path for parser bugs, IO failures, or programmer mistakes." — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-16] "The kernel-owned refusal kind set is exactly:

- `no_match`
- `ambiguous_match`
- `owned_state_unavailable`
- `guard_unevaluable`
- `unmodeled_outcome`" — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-17] "Each refusal kind must be stable enough for RDR 0005 to map to a CLI error code without inspecting error strings." — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-18] "The kernel MUST NOT print output" — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-19] "The kernel MUST NOT inspect CLI flags" — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-20] "The kernel MUST NOT discover ambient state" — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-21] "The kernel MUST NOT choose artifacts on behalf of the caller" — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-22] "The kernel MUST NOT initiate work" — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-23] "The kernel MUST NOT execute persistence side effects directly." — (Proposed Solution / Technical Design / Normative Contracts)

[REQ-24] "a resolution input is the tuple of flow identity, transition table revision, accessor-produced owned tag snapshot, observed tag-set, and freshly recognized outcome tag." — (Proposed Solution / Technical Design / Load-Bearing Decisions / Identity)

[REQ-25] "Replaying that tuple must replay the disposition." — (Proposed Solution / Technical Design / Load-Bearing Decisions / Identity)

[REQ-26] "the only successful selection is exactly one matching edge after guard evaluation." — (Proposed Solution / Technical Design / Load-Bearing Decisions / Selection / predicate)

[REQ-27] "Zero, multiple, unavailable, or unevaluable candidates are refusals unless the table contains a modeled escape edge that itself matches exactly once." — (Proposed Solution / Technical Design / Load-Bearing Decisions / Selection / predicate)

[REQ-28] "Define the internal resolver package boundary, result taxonomy, and pure resolution entry point without CLI output or persistence side effects." — (Implementation Plan / Phase 1: Kernel Boundary)

[REQ-29] "Implement tag-set assembly, guard evaluation delegation, exact-one edge selection, and typed refusal behavior." — (Implementation Plan / Phase 2: Evaluation Semantics)

[REQ-30] "Add focused kernel tests for deterministic replay and refusal classes, using fixtures that exercise owned, observed, and freshly recognized tags." — (Implementation Plan / Phase 3: Replay Validation)

[REQ-31] "Expose only the kernel values needed by RDR 0005; do not add command output or state mutation semantics in this RDR." — (Implementation Plan / Phase 4: CLI Integration Handoff)

[REQ-32] "No third-party dependency is proposed at this stage." — (Implementation Plan / New Dependencies)

[REQ-33] "Both calls return value-identical transition plans, including next tags and accessor write descriptions." — (Validation / Testing Strategy / Scenario 1)

[REQ-34] "The kernel returns the typed no-match refusal and performs no persistence side effect." — (Validation / Testing Strategy / Scenario 2)

[REQ-35] "The kernel returns the typed ambiguous-match refusal unless one explicit escape edge matches exactly once." — (Validation / Testing Strategy / Scenario 3)

[REQ-36] "The kernel returns the typed unmodeled-outcome refusal and performs no persistence side effect." — (Validation / Testing Strategy / Scenario 4)

[REQ-37] "The kernel returns the corresponding value-level typed refusal and does not fall back to ambient discovery or the CLI/Go error path." — (Validation / Testing Strategy / Scenario 5)

[REQ-38] "Resolution is bounded by the supplied transition table and tag snapshots." — (Validation / Performance Expectations)

[REQ-39] "No byte-stable hash or canonical serialization is introduced; determinism is value-level replay of the input tuple named in A1." — (Validation / Performance Expectations)

[REQ-MVV] "Resolve must name and implementation must add a replay test that feeds the same table, owned snapshot, observed tags, and recognized outcome to the kernel twice and asserts value-identical dispositions. The same validation must include at least one value-level refusal each for `no_match`, `ambiguous_match`, `owned_state_unavailable`, `guard_unevaluable`, and `unmodeled_outcome`, and must assert those modeled refusals do not use the CLI or Go error path." — (Implementation Plan / Minimum Viable Validation)

## Assumptions

ASSUMPTION: "Value-identical" means equality of the returned Go values, including slice/map contents, not byte-identical serialization; the Performance Expectations section explicitly excludes canonical serialization.

ASSUMPTION: The recognized outcome participates in the evaluation tag-set as a distinct recognized-provenance tag even though it shares the model's declared tag vocabulary.

ASSUMPTION: An unavailable owned snapshot and an unevaluable guard are explicit input/evaluation states, not empty owned tags or a false guard result; otherwise their required refusal kinds could not be distinguished.

ASSUMPTION: Planned owned-tag writes are inert value descriptions returned in a transition plan; the resolver does not invoke the accessor executor.

ASSUMPTION: A modeled escape does not override a simultaneously matching ordinary row. Consistent with RDR 0002, it succeeds only after zero or multiple ordinary matches and only when exactly one corresponding escape row matches.

ASSUMPTION: The kernel consumes normalized candidate rows and symbolic guard results compatible with RDRs 0002 and 0003; parsing sparse TOML and validating predicate syntax remain outside this RDR.

ASSUMPTION: Determinism includes independence from Go map iteration order wherever input tags, predicates, next tags, or write descriptions use maps.

## QUESTIONS

None. The peer RDR contracts settle escape-row, normalized-table, symbolic-guard, and accessor-boundary interpretations sufficiently for tests-first implementation.
