# REQ List — RDR 0008 Ownership of the recognized-outcome tag key name

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0008-recognized-tag-key-ownership.md`. Quotes are verbatim, copied
from the projector (`rdr inspect --select <id>`) or from the record itself,
never transcribed.

Element ids are carried where the REQ derives from a labelled contract
(`0008:C1` … `0008:C6`, `0008:MVV`). Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts
- `LBD` = Proposed Solution / Technical Design / Load-Bearing Decisions
- `TD` = Proposed Solution / Technical Design (prose, outside fences)
- `AP` = Proposed Solution / Approach
- `CD` = Proposed Solution / Technical Design / Capability Dependencies
- `RM` = Trade-offs / Risks and Mitigations
- `FM` = Trade-offs / Failure Modes
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (Prerequisites and phases)
- `TS` = Validation / Testing Strategy

## Standing notes for this REQ set

**S1 — the two-halves split is structural, and half of it is already built.**
The record's Testing Strategy splits Done: "scenarios 1's kernel half, 3, 6, 8,
and 9 are green, and Phase 1's pointer comment has landed" is this RDR's own
implementation; "scenarios 2, 4, 5, 7, and the lint half of the MVV" are
"*scoped to the peer that owns the code*" (RDR 0002). RDR 0002 has since been
implemented (BUILD-ORDER run 2 of 8, this RDR is run 7), and its loader already
carries part of the carried half at HEAD:
`internal/table/load.go:160-171` enforces C2's two naming rules under
`internal/table/category.go:41` `CatReservedTagKey Category = "reserved_tag_key"`.
So REQs sourced from blocks 2 and 5's declaration channel are **verify-and-pin**
work, not greenfield; REQs sourced from blocks 3, 4, 5's predicate half and the
near-miss advisory are **net-new**. Each REQ below is marked
*(landed in 0002)* or *(net-new)* where the distinction is load-bearing.

**S2 — §JD-8/§JD-9 answered; the citations are this Stage's to make**
(`artifacts/deviations.md` D3). The record's Status line names the answers but
no body clause cites them, and `grep -n '§D10\|§D8\b'` over the record hits only
the Status line. The answers are normative and land as REQ-46/47/48:
JDR 0001 §D10 (`docs/jdr/0001-resolve-kernel-seam.md:509`) gives block 3's
payload its carrier — one `omitempty` `Findings []clierr.Finding` field, and
"0008's per-key `reserved_tag_key` with the near-miss advisory in `hint`" —
and JDR 0001 §D8 (`:386`) gives the `--tag recognized=` case its site and group:
"A `--tag` naming an **owned** key or `recognized` is refused at the CLI before
any accessor runs … both are `GroupUserEnv` with their own codes (§D11)."
`internal/cli/flow_input.go:197` already implements the §D8 half.

**S3 — deviations D1 and D2 are dispositioned by running their checks, not by
editing the record.** D1 moves scenario 3's capture point to RDR 0007's per-atom
seam (REQ-76, ASSUMPTION-3); D2's grep runs against re-locked 0002 (REQ-106).

---

## A. The reserved key and the assembled view (block 1)

- [REQ-1] `0008:C1` "When a resolve carries a freshly recognized outcome, the assembled evaluation view MUST bind it under exactly the tag key `recognized` (`internal/resolve/resolve.go::recognizedTagKey`)." — (NC)

- [REQ-2] `0008:C1` "The key is a reserved kernel keyword: table authors conform to it and MUST NOT rebind it." — (NC)

- [REQ-3] `0008:C1` "The reservation holds unconditionally — the key is reserved whether or not a given resolve binds it — while the binding obligation is scoped to resolves that carry an outcome: `assemble` injects only for a non-empty `Input.Recognized`, so an absent outcome yields a view with no `recognized` key." — (NC)

- [REQ-4] `0008:C1` "The conformance test's input domain is non-empty outcomes." — (NC)

- [REQ-5] `0008:C1` "The binding obligation above is scoped so that such a resolve is outside it, not in breach of it — an empty-string outcome is degenerate for reasons that are not this RDR's to rule on." — (NC) *(negative REQ: the empty-string-outcome path is NOT to be closed here; a test asserting a refusal for it would over-implement.)*

- [REQ-6] "the conformance test MUST run the real kernel and assert the assembled view binds the recognized outcome under the normalizer's spelling — a test comparing two hardcoded literals verifies the doc against itself and is non-conforming (premortem P-5)." — (TD)

- [REQ-7] "Whether the normalizer additionally references an exported kernel constant for the *spelling* is implementation latitude the implementer may resolve; this RDR requires only the behavioral test" — (TD) *(latitude REQ: `internal/table/model.go:19` `RecognizedTagKey` and `internal/resolve/resolve.go:125` `recognizedTagKey` being two constants is conforming; the behavioral test is what must exist.)*

- [REQ-8] "That latitude is about the spelling constant only — it is distinct from the `Input` predicate the Enforcement locus settles, which this RDR does require the kernel to export." — (TD)

- [REQ-9] "record the reserved-keyword rule where implementers read it — this RDR's Normative Contracts as authority, plus a pointer comment on `internal/resolve/resolve.go::recognizedTagKey` citing RDR 0008. No kernel behavior change." — (IP Phase 1) *(net-new: `resolve.go:122-125`'s comment does not cite RDR 0008 at HEAD.)*

## B. Declaration naming (block 2)

- [REQ-10] `0008:C2` "A tag declaration with provenance `recognized` MUST be named `recognized`." — (NC) *(landed in 0002: `internal/table/load.go:160-163`.)*

- [REQ-11] `0008:C2` "A tag declaration with provenance `owned` or `observed` MUST NOT be named `recognized`." — (NC) *(landed in 0002: `internal/table/load.go:164-167`.)*

- [REQ-12] `0008:C2` "because `[tags.<tag>]` is keyed by tag name, one model admits at most one declaration named `recognized`, so at most one recognized-provenance declaration survives validation — cardinality is a consequence of the naming rule, not a separate check." — (NC) *(negative REQ: no separate cardinality check may be added.)*

- [REQ-13] `0008:C2` "two declarations spelled `recognized` are a TOML duplicate-key error in the `malformed TOML` category before this rule is reached." — (NC)

- [REQ-14] `0008:C2` "The scope unit is one `[model]` (RDR 0002's schema unit); this RDR defines no cross-model or cross-file cardinality rule." — (NC)

- [REQ-15] `0008:C2` "The bound is an upper one only: this RDR's naming rule imposes **no lower bound** of its own." — (NC) *(the lower bound at `internal/table/load.go:169-171` — `no [tags.recognized] declaration` under `CatMalformedModelDeclaration` — is RDR 0002's outcome-binding contract, NOT this RDR's; do not re-home it into `reserved_tag_key`.)*

- [REQ-16] `0008:C2` "0002 decides that the declaration exists, this RDR decides what it is named." — (NC)

- [REQ-17] `0008:C2` "For a `resolve.Table` built by any producer other than 0002's loader no lower bound applies, and a row meant to match the recognized tag under an undeclared or innocent name refuses `no_match` at resolve time" — (NC)

- [REQ-18] `0008:C2` "Any violation is a data-level validation failure in the `reserved_tag_key` category, reported at table load/lint before any resolution — never a kernel refusal." — (NC)

- [REQ-19] `0008:C2` "The reserved-key comparison is byte-exact on the **post-parse** key string: case-sensitive, no trimming, no folding (so `Recognized` is an ordinary, unreserved name, and `\" recognized\"` — a whitespace-bearing key — is likewise ordinary and unreserved)." — (NC)

- [REQ-20] `0008:C2` "TOML quoting is a surface artifact, not a name variant: `[tags.\"recognized\"]` and `[tags.recognized]` parse to the identical key string and are therefore both reserved." — (NC)

- [REQ-21] `0008:C2` "The checked positions are the `[tags.<tag>]` declaration keys only." — (NC)

- [REQ-22] `0008:C2` "Predicate positions need no separate reserved-key check because RDR 0002's outcome-binding contract already decides every one of them" — (NC) *(negative REQ; the 0002-side behavior it depends on is at `internal/table/normalize.go:414-416` (guard-position `recognized` refused) and `:473-484` (match atom lifted to `Row.Outcome`).)*

- [REQ-23] `0008:C2` "`[rule.write]` is covered by a *different* pre-existing rule and must not be folded into the sentence above … what rejects it is RDR 0002's `write to non-owned tag` category … This RDR adds no write-position check and depends on that category holding." — (NC)

- [REQ-24] `0008:C2` "a reference points at the wrong declared tag. Closing it needs a rule over declaration *intent* rather than declaration *name*, which is the same out-of-scope heuristic the Failure Modes section charts" — (NC) *(negative REQ: the stray-predicate case is explicitly NOT closed here.)*

- [REQ-25] "a reserved-name owned/observed declaration, fails table load/lint in the `reserved_tag_key` category — the failure data carries the direction-specific rule identifier, the offending name, and the remedy name, before any resolution runs." — (FM)

## C. Failure payload (block 3)

- [REQ-26] `0008:C3` "Every `reserved_tag_key` failure — and any undeclared-tag failure whose offending key is the reserved name — MUST carry, at the data level, three distinct machine-readable fields: the offending name as authored, a **remedy name**, and a stable rule identifier." — (NC) *(net-new: `internal/table/load.go:160-167` calls `fail(CatReservedTagKey, <prose detail>)` with no structured payload.)*

- [REQ-27] `0008:C3` "The rule identifier is a comparable token, not prose: a golden test asserts it byte-for-byte, and human-readable wording is the renderer's to choose." — (NC)

- [REQ-28] `0008:C3` "A category consumer may map the failure, but the guidance travels in the failure data, not the renderer." — (NC)

- [REQ-29] `0008:C3` "A recognized-provenance declaration under a wrong name must be renamed **to** the reserved key. Remedy name: the literal `recognized`. Rule identifier: `reserved-tag-key/kernel-owned`." — (NC)

- [REQ-30] `0008:C3` "An owned or observed declaration named `recognized` must be renamed **away from** the reserved key … the field carries the empty string, and the rule identifier `reserved-tag-key/author-must-rename` is what tells a consumer the remedy is \"choose any other name\" rather than \"use this one\". A renderer MUST NOT present the reserved key as the required name in this direction." — (NC)

- [REQ-31] `0008:C3` "Both rule identifiers sit inside the one `reserved_tag_key` category" — (NC)

- [REQ-32] `0008:C3` "The category's stable data-level value is the token `reserved_tag_key`, matching the snake_case discriminator grammar RDR 0001 uses for `RefusalKind` values" — (NC) *(landed in 0002: `internal/table/category.go:41`.)*

- [REQ-33] `0008:C3` "The two tokens are not alternatives and a consumer MUST NOT choose between them: `reserved_tag_key` is the **category** discriminator … The **rule identifier** … is carried inside the failure payload, identifying which rule within the category fired; it is for golden assertions and remediation lookup, never for category dispatch." — (NC)

- [REQ-34] `0008:C3` "The Go type, package, and field names carrying these values are RDR 0002's to choose — this RDR constrains the values and their distinctness, not the struct." — (NC) *(latitude REQ.)*

- [REQ-35] `0008:C3` "Where the offending key is the reserved name, this payload requirement extends RDR 0002's pre-existing `unknown tag` failure; that extension is authored by this RDR and lands in 0002's implementation." — (NC)

## D. Input producer obligation and the exported predicate (block 4)

- [REQ-36] `0008:C4` "Producers of kernel `Input` MUST NOT supply an owned or observed tag keyed `recognized`; the reserved key enters the assembled view only through `Input.Recognized`." — (NC) *(net-new.)*

- [REQ-37] `0008:C4` "A breach is a producer programmer mistake — not table data, and never a new `RefusalKind` — and travels the Go error path RDR 0001 reserves for programmer mistakes." — (NC)

- [REQ-38] `0008:C4` "Enforcement is one predicate at two call sites: the kernel MUST export a construction-time predicate over `Input` that producers may call, and `Resolve` MUST apply that same predicate at entry, returning a non-nil error and no `Result` disposition on breach." — (NC) *(net-new: `internal/resolve/resolve.go:404` `Resolve` has no entry precondition and `internal/resolve` exports no `Input` predicate at HEAD.)*

- [REQ-39] `0008:C4` "That predicate carries **both** reserved-key obligations — the owned/observed tag keys of this block, and block 5's `Row.RequiresOwned` reservation — so the reserved name has exactly one enforcement point across every channel it can arrive through." — (NC)

- [REQ-40] `0008:C4` "The obligation is NOT inherited from RDR 0009 … the predicate's exported name MUST NOT be bound to any symbol name from RDR 0009." — (NC) *(negative REQ; RDR 0009 is BUILD-ORDER run 8 and unimplemented, so its symbol names do not exist yet — the constraint is that this RDR must not claim them.)*

- [REQ-41] `0008:C4` "The predicate's read-domain … it reads `Input.Owned`, `Input.Observed`, and each row's `RequiresOwned` reached through `Input.Table.Rows` — three sequences, no other `Input` field." — (NC)

- [REQ-42] `0008:C4` "`Owned` and `Observed` are **sequences of key/value tags, not maps** … so the check is a scan and a duplicate reserved key is admissible input rather than a structural impossibility." — (NC)

- [REQ-43] `0008:C4` "On multiplicity, the predicate reports **the first breach it finds and returns a single non-nil error**; it does not aggregate, and the order in which it scans the three sequences is implementation latitude." — (NC)

- [REQ-44] `0008:C4` "What this RDR requires of the implementation is therefore only that its own breach be **detected whenever present** — never skipped because another precondition also fired." — (NC)

- [REQ-45] `0008:C4` "the implementable rule is: **apply this predicate at entry and report its breach; if the table also breaches RDR 0009's shape rule, either error is conforming.**" — (NC) *(JDR 0001 §JD-5 is open; do not encode an order.)*

- [REQ-46] `0008:C4` "The check is unconditional on `Input.Recognized`: a reserved-keyed owned or observed tag is a breach whether or not the resolve carries an outcome, matching block 1's unconditional reservation." — (NC)

## E. `RequiresOwned` name reservation (block 5)

- [REQ-47] `0008:C5` "A normalized row's `Row.RequiresOwned` MUST NOT name `recognized`." — (NC)

- [REQ-48] `0008:C5` "This is a name reservation only: what `RequiresOwned` *means* is owned by RDR 0007 … this RDR cites that rule rather than restating it" — (NC) *(negative REQ: do not restate or alter `internal/resolve/resolve.go:528` `missingOwned`'s `view.has(key, ProvenanceOwned)` test.)*

- [REQ-49] `0008:C5` "this reservation is a **producer obligation on the normalizer**, not a source-lint rule over authored TOML, and it does not share the declaration channel's `reserved_tag_key` category." — (NC) *(negative REQ: no `reserved_tag_key` failure may be emitted for a `RequiresOwned` entry.)*

- [REQ-50] `0008:C5` "the **exported `Input` predicate of block 4 MUST also reject a `Row.RequiresOwned` entry naming `recognized`** on the rows of the supplied table, on the same Go error path and with the same producer-breach semantics." — (NC)

- [REQ-51] `0008:C5` "The residual this closes is kernel-side: a row carrying `RequiresOwned: [\"recognized\"]` finds the key present under `ProvenanceRecognized`, fails the owned-only test, and yields `owned_state_unavailable` naming the reserved key … It stays a documented residual on that narrowed path" — (NC) *(the residual must remain reachable via the package-internal path; do not remove it.)*

- [REQ-52] "The `RequiresOwned` reservation needs no source-lint check here (A7, Verified) — the field has no authored source form" — (IP Phase 2) *(negative REQ.)*

## F. Unchanged kernel disposition (block 6)

- [REQ-53] `0008:C6` "Kernel disposition is unchanged by this RDR for conforming input: no new `RefusalKind`, no change to RDR 0001 D3's provenance precedence (`owned` > `observed` > `recognized`), and no behavior change for any conforming table and conforming input." — (NC)

- [REQ-54] `0008:C6` "D3 remains the deterministic backstop behind the producer obligation for input that bypasses the precondition." — (NC)

- [REQ-55] `0008:C6` "Breaching input gains a non-nil Go error where it previously resolved — a programmer-mistake path, not a disposition change for any conforming caller." — (NC)

## G. Identity, near-spellings, and the near-miss advisory (LBD)

- [REQ-56] "the recognized-provenance tag key is the exact string `recognized`; equality is byte-exact string match on the **post-parse** key string — no case folding, no trimming, no aliasing." — (LBD / Identity)

- [REQ-57] "Near-spellings that parse to a different key string (`Recognized`, `RECOGNIZED`, a whitespace-bearing `\" recognized\"`) are by definition ordinary unreserved names. Quoting is *not* such a variant: `[tags.\"recognized\"]` parses to the same key string as `[tags.recognized]`, so it is reserved." — (LBD / Identity)

- [REQ-58] "**Advisory warning — normative, not deferred.**" — (LBD / Identity) *(net-new: no near-miss advisory exists in `internal/table` at HEAD.)*

- [REQ-59] "a declaration whose post-parse key is not `recognized` but becomes `recognized` under **either** Unicode-simple case folding **or** trimming of leading/trailing whitespace, or both applied together, MUST raise a **non-blocking advisory** naming both spellings." — (LBD / Identity)

- [REQ-60] "The trigger is deliberately disjunctive: `Recognized` (folding alone) and `\" recognized\"` (trimming alone) are each near-misses on their own, and scenario 4 requires the advisory on both, so a conjunctive reading would fire on neither." — (LBD / Identity)

- [REQ-61] "Advisory, not a failure, because the name is legal and this RDR must not reject a tag it does not own." — (LBD / Identity)

- [REQ-62] "the advisory carries the same machine-readable discipline as the failure payload (block 3) … two distinct fields, the **authored spelling** and the **reserved spelling** (the literal `recognized`), plus a stable rule identifier whose value is the literal `reserved-tag-key/near-miss`." — (LBD / Identity)

- [REQ-63] "It travels on an advisory channel distinct from the validation-failure list — it carries **no** `reserved_tag_key` category discriminator, because it is not a validation failure and MUST NOT participate in category dispatch or alter the load/lint verdict." — (LBD / Identity)

- [REQ-64] "The carrier's Go type and field names are RDR 0002's to choose, on the same footing as the failure payload; this RDR constrains the values and the channel's separateness." — (LBD / Identity) *(latitude REQ.)*

- [REQ-65] "**Naming** — canonical name `recognized`, matching the shipped `recognizedTagKey` constant and RDR 0001's frozen fixture. Rejected: a sigil-guarded name" — (LBD / Naming) *(negative REQ.)*

## H. Enforcement locus (LBD)

- [REQ-66] "**Chosen: (b) and (c) together — one predicate, two call sites.**" — (LBD / Enforcement locus)

- [REQ-67] "`internal/resolve/resolve.go::Resolve` already returns `(Result, error)` whose error is doc-reserved for programmer mistakes with no non-nil path in its body today, so (b) adds a check without widening the signature." — (LBD / Enforcement locus) *(the `Resolve` signature must not change.)*

- [REQ-68] "the exported predicate is the single definition, and `Resolve` calls it at entry." — (LBD / Enforcement locus) *(negative REQ: no second, independently-written check.)*

- [REQ-69] "Surface conceded, stated plainly: one exported predicate over `Input` plus one kernel-entry call." — (LBD / Enforcement locus) *(scope ceiling: exactly one new exported symbol on `internal/resolve`.)*

- [REQ-70] "The predicate's exported name is implementation latitude — do **not** bind it to any symbol name from RDR 0009 (`Final`); its predicate's exported name is 0009's to fix (A11)." — (LBD / Enforcement locus)

- [REQ-71] "Breach of the precondition returns a non-nil error and no `Result` disposition; conforming callers see no behavior change." — (LBD / Enforcement locus)

## I. Minimum Viable Validation

- [REQ-72] `0008:MVV` "**Kernel half — executed during this RDR's implementation.** A row matching on the tag `recognized` fires against a recognized outcome, asserted through the real kernel's assembled view (the behavioral conformance form, premortem P-5), and an `Input` supplying an owned or observed tag keyed `recognized` — or a row naming it in `RequiresOwned` — is rejected by the exported predicate and by `Resolve` at entry. Both run against `internal/resolve` as it stands." — (MVV)

- [REQ-73] `0008:MVV` "**Normalizer half — executed inside RDR 0002's implementation.** A table declaring a recognized-provenance tag named `recognized` loads and lints clean; the same table with the declaration renamed (and separately, with an owned tag named `recognized`) fails load/lint in the `reserved_tag_key` category whose failure data carries the direction-appropriate remedy name and rule identifier — not a silent no-match at resolve time." — (MVV) *(0002 is implemented; the load/lint half is now runnable, but the payload half is net-new — see REQ-26/29/30.)*

## J. Testing Strategy scenarios

- [REQ-74] TS-1 "a table declaring a recognized-provenance tag named `recognized` loads, lints clean, and a row matching on that tag fires against a recognized outcome — asserted through the **real kernel's** assembled view, never by comparing two hardcoded literals (premortem P-5). **Expected**: the row is selected; the recognized outcome is readable at key `recognized` with `ProvenanceRecognized`." — (TS) *(anchors: `internal/resolve/resolve.go::assemble`, `internal/resolve/fixtures_test.go::recognizedTagSensitiveTable`.)*

- [REQ-75] TS-2 "the same table with the recognized-provenance declaration renamed, and separately with an owned tag named `recognized`. **Expected**: both fail load/lint in the `reserved_tag_key` data-level category before any resolution — not a silent no-match at resolve time." — (TS)

- [REQ-76] TS-3 "a guard predicate and a match pattern both read the recognized tag in one resolve, through a **new** view-capturing guard seam added alongside `internal/resolve/fixtures_test.go::fixtureGuards` — not a change to it, since ~10 existing tests depend on its current shape." — (TS) **See ASSUMPTION-3 and deviations D1** — the capture point moves to RDR 0007's per-atom seam.

- [REQ-77] TS-3 "A row whose `Match` names the recognized tag is selected in the same resolve, so a divergence between the guard's view and the matcher's would fail one of the two assertions." — (TS)

- [REQ-78] TS-3 "The rows here are constructed directly: Final RDR 0002 mints no guard atom on `recognized` — it refuses one at load (A9) — so this is a kernel-plumbing test of A2's same-view property, not a shape any 0002-loaded table produces" — (TS)

- [REQ-79] TS-4 "the byte-exact rule on the post-parse key string treats `Recognized`, `RECOGNIZED`, and `\" recognized\"` as ordinary unreserved names (no folding, no trimming); the quoted `\"recognized\"`, which parses to the identical key string as the bare form, **is** reserved." — (TS)

- [REQ-80] TS-4 "each of the three unreserved near-misses raises the **non-blocking advisory** while `[tags.result]` — an ordinary name that is not a near-miss — raises none, so the advisory is pinned as targeted rather than blanket." — (TS)

- [REQ-81] TS-4 "The advisory's payload is asserted byte-for-byte like the failure payload: authored spelling, reserved spelling `recognized`, rule identifier `reserved-tag-key/near-miss`, and **no** `reserved_tag_key` category discriminator. The advisory MUST NOT change the load/lint verdict for any of them." — (TS)

- [REQ-82] TS-5 "a model carrying two recognized-provenance declarations … **Expected**: the load is refused with **exactly one** failure, in the `reserved_tag_key` category, naming one of the two declarations as authored; which of the two is reported is unspecified and MUST NOT be asserted" — (TS)

- [REQ-83] TS-5 "The assertion is by category, not message text (0002's failing-control rule)." — (TS)

- [REQ-84] TS-5 "The same-name variant is out of scope for this category: `[tags.recognized]` twice is a TOML duplicate-key error in `malformed TOML`, asserted as such." — (TS)

- [REQ-85] TS-6 "kernel `Input` carrying an owned or observed tag keyed `recognized`, constructed directly (bypassing lint, as a non-TOML producer would), in three variants — one with a non-empty `Input.Recognized`, one with it empty, and one carrying the reserved key **twice** in the same tag sequence." — (TS)

- [REQ-86] TS-6 "it pins block 4's first-breach rule — the predicate returns one non-nil error, not two, and the test asserts a single error rather than an aggregate." — (TS)

- [REQ-87] TS-6 "**Expected**: all three variants breach. `Resolve` returns a non-nil `error` and a zero-valued `Result` (no `Plan`, no `Refusal`) — never a new `RefusalKind`, and never a modeled disposition. The exported predicate, called directly on the same `Input`, reports the same breach: one predicate, two call sites, asserted at both." — (TS)

- [REQ-88] TS-6 "A conforming `Input` returns a nil error, pinning that the check adds no behavior change for conforming callers — including the empty `Input{}` that `internal/resolve/resolve_test.go:748` already pins as a nil-error resolve, which the new precondition MUST keep green." — (TS) *(the named pin is `TestReq20_EmptyInputTupleStillYieldsAValueDisposition`, `internal/resolve/resolve_test.go:747`.)*

- [REQ-89] TS-6 "The empty-`Input.Recognized` variant is expected to breach on the reserved key alone (block 4's unconditional clause), which is what distinguishes it from the empty `Input{}` case that carries no reserved key at all." — (TS)

- [REQ-90] TS-6 "**Doubly-breaching variant (A11; JDR 0001 §JD-5)**: a fourth variant whose table breaches *both* this RDR's reserved-key rule and RDR 0009's escape-row shape rule (a row with non-empty `Escape` and non-empty `Writes`). While JD-5 is open, this test asserts only what block 4's interim rule licenses … and it must NOT assert which of the two errors is reported." — (TS) *(RDR 0009 is unimplemented — see ASSUMPTION-5.)*

- [REQ-91] TS-7 "**both directions** of the naming rule, each passed to a synthetic consumer that maps only the categories RDR 0002 enumerates today (i.e. treats this one as unknown and falls through to a generic branch). **Expected**: the failure **data** still carries all three fields, asserted byte-for-byte and independent of anything the consumer renders." — (TS)

- [REQ-92] TS-7 "Wrong-named recognized-provenance declaration: offending name as authored, remedy name `recognized`, rule identifier `reserved-tag-key/kernel-owned`. Owned/observed declaration named `recognized`: offending name `recognized`, remedy name the empty string, rule identifier `reserved-tag-key/author-must-rename` … Both carry the same category discriminator `reserved_tag_key`." — (TS)

- [REQ-93] TS-7 "The consumer is a test stub, not RDR 0005's exit-code map: this RDR asserts the payload contract, and 0005 owns whatever mapping it later adds." — (TS)

- [REQ-94] TS-8 "a hand-constructed non-conforming `Input` whose observed or owned tags include the key `recognized`, resolved with the producer precondition bypassed (calling `assemble`'s behavior through the package-internal test path). **Expected**: D3's precedence (`owned` > `observed` > `recognized`) resolves the collision deterministically — the same input yields the same disposition on repeat runs." — (TS)

- [REQ-95] TS-8 "Scope note: the unreachability half of A5 is **not** a test." — (TS) *(negative REQ.)*

- [REQ-96] TS-9 "a hand-constructed row carrying `RequiresOwned: []string{\"recognized\"}` evaluated against a view whose `recognized` key is present under `ProvenanceRecognized`, reached through the package-internal path that bypasses the entry precondition … **Expected**: the key is reported missing and the resolution refuses `owned_state_unavailable` with `MissingOwned` containing `recognized`" — (TS)

- [REQ-97] TS-9 "**Second half (block 5's enforcement)**: the same row passed to the exported `Input` predicate reports a breach, and `Resolve` at entry returns a non-nil error rather than the `owned_state_unavailable` refusal above … Both are asserted; the pair is what makes block 5 falsifiable rather than discharged-by-assertion." — (TS)

- [REQ-98] TS-9 "There is **no source-lint half** (A7, Verified) … the corresponding negative test belongs to RDR 0002's `write to non-owned tag` suite, not this RDR's." — (TS) *(negative REQ.)*

## K. Phase 3 conformance set (items not already covered above)

- [REQ-99] "the conforming-`Input` nil-error pin including the existing empty-`Input{}` case (`resolve_test.go:748`)" — (IP Phase 3)

- [REQ-100] "the reserved-key normalization fixtures (case and whitespace variants stay unreserved; the quoted form is reserved) **plus the near-miss advisory assertions**" — (IP Phase 3)

- [REQ-101] "the golden failure-data check (offending name, remedy name, rule identifier — both directions, since block 3 pins a distinct remedy/rule pair for each)." — (IP Phase 3)

## L. §JD-8 / §JD-9 answers — citations owed at this Stage (deviations D3)

- [REQ-102] JDR 0001 §D10 "**0007's reopening is accepted: `CLIError` gains exactly one `omitempty` structured field, and it is the `Finding` record 0006 already mandates** — non-normative `Findings []clierr.Finding` (`json:\"findings,omitempty\"`), `Finding{Code, Message, Param, Locator, Hint}`, subsystem-agnostic. One field serves all four carriers: … 0008's per-key `reserved_tag_key` with the near-miss advisory in `hint`" — (`docs/jdr/0001-resolve-kernel-seam.md` §D10) *(this is block 3's and the advisory's CLI carrier; `internal/cli/clierr/clierr.go:288` `Finding` exists at HEAD with `Code`/`Message`/`Param`/`Locator`/`Hint`.)*

- [REQ-103] JDR 0001 §D10 "`flow-model-invalid` (every 0002 load category, 0003's two, 0008's `reserved_tag_key`) | UserEnv / 2 | `findings[]`" — (§D10 non-normative table) *(`reserved_tag_key` maps to one CLI code, not its own; `internal/cli/clierr/finding_0005_test.go:297` already pins `"Code": "reserved_tag_key"` in a finding.)*

- [REQ-104] JDR 0001 §D8 "A `--tag` naming an **owned** key or `recognized` is refused at the CLI before any accessor runs (P2 — refuse rather than silently shadow under owned-over-observed precedence); both are `GroupUserEnv` with their own codes (§D11). 0008's \"programmer mistake\" is the CLI caller's, which from the CLI's seat is the user; `GroupInternal` is for the CLI's own invariants." — (`docs/jdr/0001-resolve-kernel-seam.md` §D8) *(landed: `internal/cli/flow_input.go:196-200` refuses `--tag recognized=` with `codeTagReserved`; this REQ is the citation + a pin that the refusal is `GroupUserEnv`, not the block-4 Go error path.)*

- [REQ-105] JDR 0001 §D8 "`--tag` stays `Observed`." — (§D8 (a)) *(the CLI must not route a `--tag` value into `Input.Recognized` or `Input.Owned`.)*

- [REQ-106] Deviation D2's check: "after 0002 re-locks under §JD-17, `grep -n 'recognized' <0002 layout clause + rdr-fixture.toml>` shows no `[read.<id>]`/`[write.<id>]`/`[gate.<id>]` `keys` entry may name `recognized` (a writer key must be owned; `recognized` is kernel-supplied)." — (`artifacts/deviations.md` D2) *(landed in 0002: `internal/table/load.go:295-297` refuses a capability `keys` entry naming `RecognizedTagKey`, under `CatMalformedAccessorBinding` — 0002's category, not `reserved_tag_key`; the check's existence is what D2 asks for, not its category.)*

---

## ASSUMPTIONS

- **ASSUMPTION-1 — "landed in 0002" REQs are pinned, not re-implemented.** The
  record was written when RDR 0002 was Final-unimplemented and describes its
  clauses as work 0002 "executes". 0002 is now implemented (BUILD-ORDER run 2;
  this RDR is run 7), and `internal/table/load.go:160-171` already carries C2's
  two naming rules under `CatReservedTagKey`. I read the carried half as
  **satisfied where the code matches the clause and owed where it does not**:
  REQ-10/11/32/106 are pinned by test, REQ-26/29/30 (payload) and REQ-58..64
  (advisory) are net-new work this Stage owes, since the record's Prerequisite
  "**The payload and advisory carrier close at JDR 0001 §JD-8**" is now
  satisfiable (§JD-8 answered by §D10 on 2026-08-24). The alternative reading —
  that the whole carried half is out of this RDR's scope because 0002 already
  ran — would leave the MVV's normalizer half permanently unexecutable, which
  the record forbids ("gated by the Prerequisite below, not deferred by
  choice").

- **ASSUMPTION-2 — the payload lands in `internal/table`, its CLI carrier in
  `clierr.Finding`.** Block 3 says "The Go type, package, and field names
  carrying these values are RDR 0002's to choose" (REQ-34), and §D10 fixes the
  CLI-side carrier as `Findings []clierr.Finding` with `Code`/`Param`/`Hint`
  (REQ-102). I read the three fields as: offending name → the finding's `Param`
  (or an equivalent named field on the table-side failure), remedy name and
  rule identifier → table-side fields surfaced through the finding. The
  *values* and their distinctness are what this RDR pins (REQ-27/33); the field
  spelling is latitude.

- **ASSUMPTION-3 — TS-3's guard seam is written against RDR 0007's per-atom
  seam, per deviation D1.** The record's scenario 3 expects "the captured
  guard-side view satisfies `Lookup(\"recognized\") == …` and `Len()` equals the
  key count", but RDR 0007 (implemented, run 1) fenced `Evaluate(atom, value)`
  as never seeing the view — `internal/resolve/fixtures_test.go:23` is
  `Evaluate(atom resolve.GuardAtom, _ string)`. D1's disposition is binding:
  capture the `value` handed to the guard atom over `recognized` and assert it
  equals `(in.Recognized, ProvenanceRecognized)` while a `Match` on the same key
  selects the row in the same resolve. The contract (A2's same-view property) is
  unchanged; only the capture point moves. Escalate as SPEC-DEFECT only if the
  same-view property cannot be asserted at that seam.

- **ASSUMPTION-4 — REQ-69's "one exported predicate" is a ceiling on
  `internal/resolve`'s new surface, not a ban on unexported helpers.** The
  record concedes exactly "one exported predicate over `Input` plus one
  kernel-entry call" and scopes the "no kernel change" claim to disposition.
  I read internal decomposition as latitude.

- **ASSUMPTION-5 — TS-6's doubly-breaching variant is authored now against a
  hand-constructed table, not deferred to RDR 0009's implementation.** RDR 0009
  is BUILD-ORDER run 8 and unimplemented, so no 0009 shape check exists to
  co-fire. The variant is still authorable: a row with non-empty `Escape` and
  non-empty `Writes` is constructible on today's `resolve.Row`, and the test
  asserts only that *this* RDR's reserved-key breach is not silently skipped
  (REQ-44/90) — which holds whether or not 0009's check is present. Asserting
  a two-error interaction would require 0009 and would also encode §JD-5, which
  REQ-45 forbids.

- **ASSUMPTION-6 — "byte-exact on the post-parse key string" is compared with
  Go `==` on the decoded TOML key.** REQ-19/20's two consequences (quoted form
  reserved, whitespace form not) both follow from comparing the key *after* the
  TOML decoder has stripped quoting, which is what `internal/table/load.go`'s
  map iteration already yields. No normalization pass is added.

- **ASSUMPTION-7 — the near-miss advisory's fold is `strings.EqualFold`-class
  simple case folding, not full Unicode case-folding tables.** The record says
  "Unicode-simple case folding" (REQ-59); Go's `strings.EqualFold` implements
  simple folding, and scenario 4's three cases (`Recognized`, `RECOGNIZED`,
  `" recognized"`) are all decided by it plus `strings.TrimSpace`. The
  disjunctive trigger (REQ-60) is the load-bearing part and is unambiguous.

---

## QUESTIONS

- **Q1 — Does the near-miss advisory have a delivery channel on the *success*
  path, or only on a failure envelope?** REQ-63 requires the advisory to travel
  "on an advisory channel distinct from the validation-failure list" and to
  "MUST NOT … alter the load/lint verdict" — so a model whose only issue is a
  near-miss **loads clean** and must still surface the advisory. §D10 places
  0008's advisory "in `hint`" on a `clierr.Finding`, but `Findings` is a field
  on `CLIError` — the *failure* envelope — and §D10's own reconciliation note
  says the success payload's carrier is 0006's non-omitempty `data.findings`.
  Two readings, materially different:
  (a) The advisory rides `data.findings` on the success payload (0006's lint
  channel), so `flow lint` on a clean-but-near-missing model exits 0 with a
  non-blocking finding.
  (b) The advisory is a table-package-level return value that no CLI verb
  surfaces yet, and its CLI delivery is RDR 0005/0006's later work.
  **Proceeding under (a)** for the table-side obligation and asserting the
  advisory at the `internal/table` boundary (a returned advisory list distinct
  from the error), with CLI delivery left to whichever verb renders lint
  output. This satisfies REQ-58..64 and REQ-80/81 without pre-deciding a CLI
  surface this RDR explicitly disclaims ("which user-facing command surfaces
  it, and under which exit code, is RDR 0005's mapping decision and is not
  settled here", AP item 3). Recorded because a golden test written against
  the wrong channel would need re-anchoring, and because no predecessor
  settles where a *non-blocking* load-time advisory travels: 0006 owns lint
  findings over a normalized graph, and this advisory fires strictly upstream
  of normalization, at declaration parse.

- **Q2 — Does the `reserved_tag_key` payload requirement (REQ-26) extend to
  `unknown tag` failures, given 0002's implemented lower-bound check?** Block 3
  binds "any undeclared-tag failure whose offending key is the reserved name",
  and REQ-35 says that extension "lands in 0002's implementation". But 0002 as
  implemented does not fail `unknown tag` for a reference to `recognized`: it
  fails *earlier*, at `internal/table/load.go:169-171`, with
  `CatMalformedModelDeclaration` "no [tags.recognized] declaration" — a
  different category, and one REQ-15 says is 0002's, not this RDR's. Two
  readings:
  (a) The clause is now vacuous on the 0002 path (no `unknown tag` failure can
  carry `recognized` as its offending key, because the missing declaration trips
  a mandatory-declaration failure first), and the payload requirement binds only
  `reserved_tag_key` failures.
  (b) The `CatMalformedModelDeclaration` failure inherits the payload
  requirement as the successor of the `unknown tag` case the record anticipated.
  **Proceeding under (a).** REQ-15 is explicit that the lower bound is 0002's
  clause and this RDR must not re-home it, and REQ-26's binding is on failures
  "whose offending key is the reserved name" — the mandatory-declaration failure
  has no offending key, it reports an absence. Recorded because (b) would add a
  payload obligation to a category this RDR does not own, and because the
  record's phrasing was written against a 0002 that had no mandatory-declaration
  rule.
