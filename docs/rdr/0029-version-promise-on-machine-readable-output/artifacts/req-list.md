# req-list — 0029 version promise on machine-readable output

Source: `docs/rdr/0029-version-promise-on-machine-readable-output.md` (Status: Final, 1969 lines).
Quotes are exact bytes via `rdr inspect --select`; element ids carry the contract each REQ derives from.

## REQ

[REQ-1] "The `--as=json` terminal envelope carries a `schema_version` string field of the form `MAJOR.MINOR`, present on both terminal records: the `ok` envelope (`internal/cli/respond::Success`) and the refusal record (`internal/cli/clierr::CLIError`)." — (0029:C1, Normative Contracts)

[REQ-2] "It is versioned independently of the binary's release version and MUST NOT be derived from it." — (0029:C1, Normative Contracts)

[REQ-3] "The value has ONE home: a single exported constant, which both terminal records read and neither declares." — (0029:C1, Normative Contracts)

[REQ-4] "the constant lives in `clierr` (the leaf both can reach) and `respond` reads it from there" — (0029:C1, Normative Contracts)

[REQ-5] "`schema_version` rides the TOP LEVEL of each terminal record only. It is never projected into `data`" — (0029:C1, Normative Contracts)

[REQ-6] "`schema_version` is NOT `omitempty` on either record" — (0029:C1, Normative Contracts)

[REQ-7] "The refusal envelope still has exactly one structured field, `findings`." — (0029:C1, Normative Contracts)

[REQ-8] "The schema version begins at `\"0.1\"` and tracks the wire, not the binary." — (0029:C1, Normative Contracts)

[REQ-9] "the minor component still increments on every change so a consumer can detect movement even while it cannot rely on compatibility" — (0029:C1, Normative Contracts)

[REQ-10] "It MUST NOT be made strict: this is an emitted envelope arriving back, not an authored document" — (0029:C1, Normative Contracts; subject is `internal/cli/flow_input.go::planEnvelope`)

[REQ-11] "Every machine-readable vocabulary this CLI emits carries exactly one declared stability tier, and the tier is recorded in `docs/cli-output-contract.md` beside that vocabulary" — (0029:C2, Normative Contracts)

[REQ-12] "`frozen` — no member added or removed within a major." — (0029:C2, Normative Contracts)

[REQ-13] "`append-only` — members MAY be added in a minor; none removed or renamed within a major. A consumer MUST tolerate an unrecognized member, and MUST NOT assert on the set's cardinality, a member's ordinal position, or a tail position." — (0029:C2, Normative Contracts)

[REQ-14] "`growing` — as `append-only`, and additionally a new member MAY fire on input that previously produced no such finding, subject to C3." — (0029:C2, Normative Contracts)

[REQ-15] "The term `closed` MUST NOT be used to DESCRIBE any of these tiers, in code comments or documentation" — (0029:C2, Normative Contracts)

[REQ-16] "It does not reach an emitted identifier that merely contains the word: `graph-coverage-closed-by-escape` is a member of a `growing` vocabulary" — (0029:C2, Normative Contracts)

[REQ-17] "Nor does it reach a verbatim quotation of a peer contract that predates this RDR (C4 quotes `0003:C7`)" — (0029:C2, Normative Contracts)

[REQ-18] "A lint finding code is introduced at severity `info`." — (0029:C3, Normative Contracts)

[REQ-19] "Introduction of an `info` code is a minor-release change and MUST NOT alter the success disposition of any input" — (0029:C3, Normative Contracts)

[REQ-21] "CI records each one's members and each finding code's severity in a committed snapshot, and fails when the current tree differs from it" — (0029:C3, Normative Contracts; see A7)

[REQ-22] "A code MAY be introduced directly at `blocking` only when it reports a condition that was already refused by some other code" — (0029:C3, Normative Contracts)

[REQ-24] "`frozen`: the envelope `type` discriminator; the severity vocabulary (`blocking`, `info`); the exit-code classes emitted by `internal/cli/clierr::ExitCodeFor`; the stderr advisory `level` set (`note`, `warning`); the gate `verdict` set (`allow`, `deny`, `indeterminate`) … `data.escape_class` (`internal/resolve::RefusalKinds`); the `graph-lint-failed` aggregate code; `findings[].operator` (`internal/guard::Operators`)" — (0029:C4, Normative Contracts)

[REQ-25] "`append-only`: `internal/table::Categories()`; the CLIError `code` vocabulary; the `graph-unprovable-coverage` `reason` set; the `flow next` unknown-`reason` set … the graph-lint BLOCKING finding codes; `findings[].block` (`internal/resolve::Block`)" — (0029:C4, Normative Contracts)

[REQ-26] "`growing`: the graph-lint ADVISORY finding codes. `0006:C17` is amended from \"closed at\" its four members to append-only as part of this RDR's implementation (A1)" — (0029:C4, Normative Contracts)

[REQ-27] "The `version` verb's payload field names (`version`, `commit`, `date` — `internal/version::Info`) are `frozen`" — (0029:C4, Normative Contracts)

[REQ-28] "Two emitted namespaces take NO tier, deliberately: `findings[].class` and `data.dispositions`" — (0029:C4, Normative Contracts)

[REQ-29] "A vocabulary is tiered wherever it is EMITTED, including on the fields of a `findings[]` element … an unassigned machine-readable surface is a defect." — (0029:C4, Normative Contracts)

[REQ-30] "A tier assignment obliges an ENUMERATION SEAM: an exported accessor returning the vocabulary's members, in the package that owns them." — (0029:C4, Normative Contracts)

[REQ-31] "`respond::Types() []string` over `{ok}`" — (0029:C4, census row, Normative Contracts)

[REQ-32] "`clierr::ExitCodes() []int` over the five values `{0,1,2,3,130}`" — (0029:C4, census row, Normative Contracts)

[REQ-33] "`respond::Levels() []string`" — (0029:C4, census row, Normative Contracts)

[REQ-34] "`cli::UnknownReasons() []string` over the union, matching `graphlint::Reasons`" — (0029:C4, census row, Normative Contracts)

[REQ-35] "`resolve::Blocks() []Block`" — (0029:C4, census row, Normative Contracts)

[REQ-36] "**`seam: none (prose-only)`**" for the CLIError `code` row; "The tier STANDS, unasserted and legibly so" — (0029:C4, census row, Normative Contracts)

[REQ-37] "Three doc comments in `respond.go` and `readPlan` still describe a `{\"type\":\"failed\",…}` record; they are stale today and Step 2 retires them" — (0029:C4, Normative Contracts; Step 2)

[REQ-38] "`findings[].operator` has TWO enumerations, `internal/guard::Operators` and `internal/table::Operators`, in separate packages … the tier binds both, and the by-value assertion pins both" — (0029:C4, Normative Contracts)

[REQ-42] "`intrastate lint --model <clean-model> --as=json` emits top-level keys `data,schema_version,type` with `\"schema_version\":\"0.1\"`. The criterion is ADDITIVE, not an exact key set" — (0029:S1, Testing Strategy)

[REQ-43] "the assertion is `schema_version` PRESENT on whichever failure is provoked — never an exact key-set match over a set that varies by failure … No `type` key appears on any of them; the record stays the bare `CLIError`." — (0029:S2, Testing Strategy)

[REQ-44] "`--plan` still decodes a `flow resolve --as json` envelope that now carries `schema_version`, unchanged." — (0029:S3, Testing Strategy)

[REQ-45] "that golden is re-captured in the same commit that adds the field, and the re-captured file differs from its predecessor by exactly the one `schema_version` key" — (0029:S3, Testing Strategy; `docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`)

[REQ-46] "exit 0, `type` still `\"ok\"`, the finding present in `data.findings` with `\"severity\":\"info\"`" — (0029:S4, Testing Strategy)

[REQ-47] "exit 2 and the record shape changes from the `ok` envelope to the bare `CLIError`" — (0029:S5, Testing Strategy)

[REQ-48] "a by-value assertion per frozen vocabulary in C4 — including the five this RDR newly assigned" — (0029:S6, Testing Strategy)

[REQ-49] "`graph-lint-failed` takes no by-value row: it is a single `const`, and its testable claim is that a blocking run returns that code and no other." — (0029:S6, Testing Strategy)

[REQ-50] "membership-and-uniqueness assertions only, each over that vocabulary's enumeration seam" — (0029:S7, Testing Strategy)

[REQ-51] "`TestReq73_TheBlockingCodeSetIsExactlyTheTenNamed` and `TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet` (both `internal/graphlint/findings_0006_test.go`) each pin an `append-only` set by exact value … Both are rewritten with `TestReq74` under S8" — (0029:S7, Testing Strategy)

[REQ-52] "No golden/snapshot test of the JSON envelope exists in the repo today, and this RDR does not add one" — (0029:S7, Testing Strategy)

[REQ-53] "0006's C17 text reads append-only, and `internal/graphlint/findings_0006_test.go::TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers` is REPLACED rather than updated … It becomes a membership-and-uniqueness assertion over `AdvisoryCodes()`, renamed off the count." — (0029:S8, Testing Strategy)

[REQ-54] "no occurrence of `closed` describing a vocabulary C4 TIERS" — (0029:S9, Testing Strategy)

[REQ-55] "`internal/guard/grammar.go` and `doc.go` over `Operators` (`frozen`), `internal/table/model.go` and `normalize.go` over its mirror, `internal/accessor/model.go:308` over the `verdict` set (`frozen`), `internal/resolve/resolve.go:45,66` over `RefusalKinds` (`frozen`), `internal/graphlint/taxonomy.go` over the advisory and `reason` sets, `internal/table/category.go` over `Categories`, and `internal/cli/flow_input.go:349`" — (0029:S9, Testing Strategy; the enumerated in-scope site census)

[REQ-56] "Exempt, per C2: the emitted identifier `graph-coverage-closed-by-escape` together with the Go identifiers that produce it (`ClosedByEscape`, `closedBy`), and verbatim quotations of `0003:C7`." — (0029:S9, Testing Strategy)

[REQ-57] "A grep over that enumerated site list is the assertion; an unscoped grep for the word is not" — (0029:S9, Testing Strategy)

[REQ-58] "`schema_version` is a string (`\"1.0\"`), not a pair of integers and not a number" — (0029:D-wire-byte-format, Load-Bearing Decisions)

[REQ-59] "when a new finding could be introduced at either severity, `info` is chosen unless C3's re-attribution clause applies" — (0029:D-selection-predicate, Load-Bearing Decisions)

[REQ-60] "Add the field at the single output gateway — `respond.Success` and `clierr.CLIError` — so every verb inherits it without per-verb work" — (Implementation Plan, Phase 1 Step 1)

[REQ-61] "Build the enumeration seams C4's census obliges for the five vocabularies that admit one … The seams are production accessors, not test helpers" — (Implementation Plan, Phase 1 Step 3)

[REQ-62] "Also correct the two doc comments that describe the failure envelope as `{\"type\":\"failed\",…}` — `internal/cli/respond`'s package doc and `readPlan`'s comment" — (Implementation Plan, Phase 1 Step 2)

[REQ-63] "Sites describing an untiered vocabulary, and the non-tier senses of the word, are out of scope and stay; so do C2's exemptions." — (Implementation Plan, Phase 1 Step 2)

[REQ-MVV] "1. Build at the current HEAD with `schema_version` implemented; run `intrastate lint --model <clean-model> --as=json` and record the full envelope. It reports `\"schema_version\":\"0.1\"` and no findings. 2. Add a new graph-lint finding code at severity `info` that fires on the model from step 1. 3. Re-run the same command against the same unmodified model. 4. **Expected end state**: the exit code is unchanged (0), the envelope's `type` is still `ok`, and the new finding appears in `data.findings` carrying `\"severity\":\"info\"`. The schema version's MAJOR is unchanged — a growing-vocabulary member is an additive change, so C1 moves the minor (`\"0.1\"` → `\"0.2\"`) and never the major." — (0029:MVV, Minimum Viable Validation)

## EXCLUDED

EXCLUDED: "A consumer MUST ignore object properties with unrecognized names, and MUST reject an envelope reporting an unsupported major." — (0029:C1) — the record itself rules these unassertable: "Both are obligations on code this project does not ship and cannot test; they are published prose (C2's register), not assertions". Peer-owned/consumer-side, no in-repo witness.

EXCLUDED: "From `\"1.0\"` onward: the minor component increments for backward-compatible additions … The major component increments for changes that are not backward-compatible" — (0029:C1) — release-classification rule binding a future release; the MVV oracle records "the release-classification half stays human" with nothing in-repo defining a release to range over.

EXCLUDED: "While the binary's version is `0.x`, a tier is a DECLARATION OF INTENT … and MUST NOT be read by a consumer as a guarantee already in force." — (0029:C2) — an obligation on the consumer's reading, not an observable of this binary.

EXCLUDED: "The disclosure obligation binds during `0.x` as well as after." — (0029:C3) — process obligation over release notes; the mechanical half is REQ-21 (the CI snapshot). The release note itself is recorded as "a human act".

EXCLUDED: "the CLIError `code` row has none buildable (A5), so its `append-only` tier stays a prose promise and a new code there is the one promotion-shaped event this mechanism cannot catch" — (0029:C3) — an explicit negative ("none here") with no observable; carried positively as REQ-36.

EXCLUDED: "The census below is therefore complete as of this record's implementation, not complete for all time" — (0029:C4) — scope qualifier, not a testable clause.

EXCLUDED: "They take their tier assignments in `cli/0021` itself, in the change that adds them" — (0029:C4) — peer-owned; `cli/0021` is the owner, this record only cites it.

EXCLUDED: "A real seam needs the ~36 constants moved into `clierr` and ~71 raise sites repointed: a taxonomy consolidation, not this RDR's additive accessor." — (0029:C4) — scope the record itself defers.

EXCLUDED: "Write the tier table and the C1/C3 rules into `docs/cli-output-contract.md`, and point `llms.txt` at it" — (Phase 2, Activation Step 1) — Phase 2 Operational Activation. Documentation/publication act; see ASSUMPTION-3 on phase boundary.

EXCLUDED: "The field ships in the first `v*` tag (A2). Before tagging, confirm `make docs-check` passes" — (Phase 2, Activation Step 2) — release-operations step, outside the code implementation.

EXCLUDED: "When the binary reaches 1.0.0, the schema version moves to `\"1.0\"` and the tiers stop being intent and start being guarantees." — (Phase 2, Activation Step 3) — the record names this "The one deferred step."

EXCLUDED: "`0006:C17` amended from \"closed at\" its four members to append-only, with the `want` literal at `findings_0006_test.go:288-293` edited in the same change" — (Prerequisites) — a prerequisite on a peer record's text; the in-repo half is carried as REQ-53.

EXCLUDED: "Omitted deliberately: no alternative in this RDR was weighed on empirical performance grounds." — (Performance Expectations) — explicit non-claim, no observable.

EXCLUDED: "this RDR claims no hash, no content-addressed identity and no replay-stable digest, so the determinism rider does not bind it" — (0029:G-cross-cutting) — negative with no observable.

EXCLUDED: "Field ORDER is not asserted here — the assertions are by named key (S1), so a consumer MUST NOT read key order as contract." — (0029:G-cross-cutting) — a prohibition on asserting; its positive form is REQ-42/REQ-43.

EXCLUDED: "two envelopes are \"the same shape\" iff they share a `schema_version` major, and while the major is `0`, iff they share the full `MAJOR.MINOR`" — (0029:D-identity) — a definition of identity used by other clauses, not independently observable.

EXCLUDED: "`schema_version`, rejecting `format_version` … and rejecting `version` outright" — (0029:D-naming) — naming rationale; the operative name is fixed by REQ-1.

EXCLUDED: "Promotion of a finding code from `info` to `blocking` is a distinct release event. It MUST be disclosed in the release notes for the release that carries it, naming the code. A promotion MUST NOT occur in a patch release." — (0029:C3) — an obligation over a release event, not an observable of the tree. The record calls the release note "a human act" (Risks and Mitigations) and the MVV oracle records "Nothing in-repo defines a release for a test to range over, so the release-classification half stays human." The mechanical half is REQ-21, A7's CI seam/severity snapshot, which is covered.

EXCLUDED: "the change MUST name the pre-existing code it re-attributes from, and MUST show that the input refused by the new code was already refused — same input, same verdict, different code" — (0029:C3) — the re-attribution evidence obligation binds a change author at review time; C3 discharges the exception "by evidence, not by judgement", but the subject is "the change" and no in-repo artifact carries the naming or the same-input demonstration.

EXCLUDED: "its required leading `schema` field (`intrastate.graph/1`, additive within `/1`) versions the document, while `schema_version` versions the envelope carrying it" — (0029:C4; subject is `cli/0021:C2`) — the two-marker protocol's document side does not exist on this branch: C4 states `cli/0021`'s surfaces "do not exist on `main` yet" and "take their tier assignments in `cli/0021` itself, in the change that adds them". The envelope-side half is carried by REQ-5.

EXCLUDED: "an unsupported envelope major is rejected under C1 before `data` is read at all" — (0029:C4) — consumer-side over a document surface absent from the tree. C1 records the vacuity directly: "the in-repo consumer is the one class of reader for which that MUST is vacuous, which is why no witness exists for it here."

EXCLUDED: "envelope major supported, document `schema` marker unrecognized — the envelope parses, its non-`data` members … are trustworthy, and `data` alone is opaque" — (0029:C4) — the inverse mixed case has no in-repo witness while `cli/0021`'s export is unlanded; same ground as the preceding entry, the document `schema` marker being absent from the tree.

## ASSUMPTION

ASSUMPTION-1: "the five vocabularies that admit one" (Step 3) is read as the five seam REQs REQ-31..REQ-35 — `respond::Types`, `clierr::ExitCodes`, `respond::Levels`, `cli::UnknownReasons`, `resolve::Blocks` — i.e. the six owed rows in C4's census minus the CLIError `code` registry that A5 refuted. C4 states "Six owed" and records one as `seam: none (prose-only)`, so six minus one is five.

ASSUMPTION-2: REQ-3/REQ-4's "single exported constant" is taken to live in `internal/cli/clierr` with `respond` reading it, per C1's explicit siting. C1 marks the `clierr` home as "A6, Pending on the edit landing rather than on the layering" — the layering is already grounded, so no alternative home is in play.

ASSUMPTION-3: the REQ set above covers Phase 1 (Code Implementation) only. Phase 2's three Activation Steps are excluded as operational/documentation acts, except REQ-11's "recorded in `docs/cli-output-contract.md`", which is quoted from C2 (a normative contract) rather than from Phase 2 and is therefore carried as a REQ.

ASSUMPTION-4: REQ-8's initial value `"0.1"` and REQ-42's `"schema_version":"0.1"` are the same claim observed at two levels (constant, emitted envelope); both are kept because one is the declaration and the other the wire observation.

ASSUMPTION-5: REQ-54/REQ-55's `closed`-retirement is scoped to the enumerated site census in S9, not to a repo-wide grep — S9 states the unscoped grep "is not" the assertion. REQ-57 carries that scoping rule as its own REQ so a later stage cannot silently widen it.

## QUESTIONS

(none)
