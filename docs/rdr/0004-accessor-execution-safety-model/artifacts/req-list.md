# REQ List — RDR 0004 Accessor Execution Safety Model

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0004-accessor-execution-safety-model.md`. Quotes are verbatim, copied
from the projector (`rdr inspect --select <id>`) for fenced elements and read
from the record for testable prose outside the fences — never transcribed by
hand.

Element ids (`0004:C6`, `0004:MVV`) are carried wherever a REQ derives from a
labelled contract, so a later stage can trace the REQ back to its contract. The
record has 17 `C` elements but more than 17 REQs: several fences carry multiple
independent obligations, and the Disposition Table, Fidelity Table, Oracle
Discriminability table, Load-Bearing Decisions, Failure Modes, and Testing
Strategy each carry testable prose outside the fences.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced)
- `AP` = Proposed Solution / Approach
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `LBD` = Technical Design / Load-Bearing Decisions
- `RT` = Technical Design / Round-Trip / Inverse Invariants
- `DISP` = Technical Design / Disposition Table
- `ORA` = Technical Design / Oracle Discriminability
- `FID` = Technical Design / Fidelity Table
- `DESK` = Technical Design / Desk Trace
- `IC` = Technical Design / Illustrative Code
- `CAP` = Proposed Solution / Capability Dependencies
- `EIA` = Implementation Plan / Existing Infrastructure Audit
- `FM` = Trade-offs / Failure Modes
- `RM` = Trade-offs / Risks and Mitigations
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (phases, prerequisites, new dependencies)
- `TS` = Validation / Testing Strategy
- `PE` = Validation / Performance Expectations
- `CA` = Research Findings / Critical Assumptions

**Standing note on ownership.** RDR 0002 has already shipped the accessor
*declaration* carrier and several of its load-time validations
(`internal/table/model.go::Accessor` carries `Role`, `Path`, `Keys`, `Timeout`,
`ReadBack`; `internal/table/load.go::accessorTable` refuses an absent /
non-duration / non-positive `timeout`, a writer without `read_back = true`,
`read_back` on a non-writer, and an undeclared or reserved key; the binding
counts — a key served by zero or two readers, a written key not served by
exactly one writer — are enforced by the loader as `malformed accessor
binding`). This RDR owns **execution**, so a validation REQ below that RDR 0002
already discharges is satisfied by citing that surface plus a test at this
RDR's boundary; it is not a licence to re-implement the check. REQs marked
**[0002-delivered]** are in that class.

---

## Capability and identity

- [REQ-1] "Every accessor definition MUST declare exactly one capability: read, gate, or write." — (NC, `0004:C1`)

- [REQ-2] "Runtime execution MUST reject any attempt to use an accessor for a different capability than the one declared." — (NC, `0004:C1`) — the refusal class is `capability_mismatch` (DISP, "Accessor used off-capability").

- [REQ-3] "Within one flow, each `(flow id, accessor name, capability)` identity MUST resolve to exactly one accessor binding. Missing and multiply-bound identities MUST fail validation before resolution." — (NC, `0004:C2`)

- [REQ-4] "an accessor is identified by `(flow id, accessor name, capability)`. Capability is part of the identity, so the same id may appear in both `[read.x]` and `[write.x]` — two identities, not a rebinding (JDR 0001 §D7(ii)); only a second binding of the *same* triple is multiply-bound." — (LBD, Identity) — the discriminating case: a same-id read/write pair MUST NOT be reported as multiply-bound.

- [REQ-5] "the executor selects a binding by capability table and id: an invocation that needs capability X selects only from `[X.<id>]`, never from a same-id entry under another table." — (LBD, Selection / predicate)

- [REQ-6] "Missing or multiply-bound accessors are validation failures." — (LBD, Selection / predicate); at runtime an unbound name is the refusal class `unknown_accessor` (DISP, "Accessor name not bound"). **[0002-delivered in part]** — the loader's binding counts cover the reader/writer key bindings; the identity-triple check is this RDR's.

- [REQ-7] "the canonical names are \"read accessor\", \"gate accessor\", and \"write accessor\". Rejected names: \"hook\" and \"action\" because they imply arbitrary transition callbacks." — (LBD, Naming) — binding on identifier and diagnostic spelling.

- [REQ-8] "RDR 0002 owns the TOML carrier. This RDR owns the accessor execution semantics embedded behind accessor references." — (LBD, Wire / byte format) — this RDR MUST NOT introduce a second carrier.

- [REQ-9] "the definition is the capability-table entry (`[read.<id>]`/`[write.<id>]`/`[gate.<id>]`), `keys` is the binding on readers and writers alike, and `read_back = true` is fixed on writers. This RDR does not spell the layout; it consumes it." — (LBD, Definition shape; JDR 0001 §D7(ii)) **[0002-delivered]**

## Caller-supplied artifacts — no ambient discovery

- [REQ-10] "Accessors MUST operate on caller-supplied artifact roles." — (NC, `0004:C3`)

- [REQ-11] "The accessor executor MUST NOT discover authoritative artifacts from ambient process state." — (NC, `0004:C3`) — a definition attempting ambient artifact discovery MUST fail validation before execution (TD, "validation rejects … and ambient artifact discovery before resolution"; MVV validation arm `ambient_artifact_discovery`).

- [REQ-12] "Artifact role not supplied" is a refusal, and the class minted is `execution_failure` — (DISP). Restated normatively in FM: "An unsupplied or unreadable artifact role is **not** a separate class — it is `execution_failure`, matching the Disposition Table."

## Read accessors — branch discipline

- [REQ-13] "A read accessor MUST return typed tag values or a typed refusal." — (NC, `0004:C4`)

- [REQ-14] "It MUST NOT mutate authoritative artifacts." — (NC, `0004:C4`)

- [REQ-15] "Every read accessor definition MUST declare the requested key set as validated metadata." — (NC, `0004:C5`) **[0002-delivered]** — `keys` is required on every capability-table entry.

- [REQ-16] "The set MUST NOT be derived from the keys a read actually resolved." — (NC, `0004:C5`) — the discriminating oracle is ORA scenario 2's second control: "the derived-key-set implementation must fail the truncation row rather than reporting success over a smaller set".

- [REQ-17] "A missing or empty requested key set MUST fail validation before execution." — (NC, `0004:C5`) — this is the eighth validation arm and is explicitly **unwitnessed** by the spike (CA A9; DISP, "the missing/empty requested-key-set arm is not witnessed").

- [REQ-18] "a read that lost keys still reports success over a smaller set" is the failure the pinned set prevents; "The Resolve spike's default path (`…::expectedTagKeys`) is exactly that circular shape — fixture convenience, not the contract." — (LBD, Requested key set) — the implementation MUST NOT port `expectedTagKeys`.

## Read completeness

- [REQ-19] "The typed-tag-values branch carries a completeness guarantee: a read accessor MUST return the tag set for exactly the keys it was asked for — no requested key missing, no unrequested key added — or take the refusal branch." — (NC, `0004:C6`)

- [REQ-20] "A partial or truncated read is a refusal, not a value." — (NC, `0004:C6`)

- [REQ-21] "A missing key MUST be distinguishable from an unread key only by which branch is taken — absence is a value, unreadability is a refusal." — (NC, `0004:C6`)

- [REQ-22] "One unreadable requested key MUST refuse the whole read; the executor MUST NOT return the keys that did resolve." — (NC, `0004:C6`) — restated at LBD (Read refusal granularity) as "the conservative choice … a decision, not an artifact of the spike's early return".

- [REQ-23] "the success predicate for a read accessor is \"every requested key resolved\", not \"at least one key resolved\"." — (LBD, Read completeness)

- [REQ-24] "Consumers may therefore treat a returned tag set as exactly the requested keys." — (LBD, Read completeness) — the consumer-facing half of REQ-19; the returned set is key-set-equal to the request (FID, `read` row: "branch equality plus key-set equality").

- [REQ-25] "A read that cannot resolve every requested key MUST be reported as `incomplete_read`, its own refusal class, distinct from execution failure and from timeout." — (NC, `0004:C7`)

- [REQ-26] "The refusal MUST name the requested keys it could not read." — (NC, `0004:C7`) — TS 2 adds "the refusal carries no values and names the keys it could not read"; the empty-value assertion is normative there and unwitnessed (A9).

- [REQ-27] "When a read exceeds its timeout before resolving every requested key, `timeout` takes precedence over `incomplete_read`." — (NC, `0004:C7`) — the discriminating fixture "makes both classes true at once and asserts the name `timeout`" (ORA 7), with the scenario-2 truncation fixture as the negative control, which "must still return `incomplete_read`".

- [REQ-28] "a key is *absent* only when the binding read the artifact successfully and the key was not there. Every other outcome — the artifact did not parse, the transport truncated, the key's value failed type coercion, the value read is the reserved `<clear>` literal, permission was denied — is *unreadable*, because in each the binding did not establish that the key is missing." — (LBD, Absent vs unreadable)

- [REQ-29] "A binding that cannot tell the two apart at its own boundary MUST report `incomplete_read`; guessing absence is the failure this contract exists to prevent." — (LBD, Absent vs unreadable) — unreadable is the default classification.

## Absence crosses the seam as omission

- [REQ-30] "A requested key the artifact genuinely does not carry MUST NOT be presented to the resolver as a present owned tag." — (NC, `0004:C8`)

- [REQ-31] "The accessor layer MAY represent absence however it chooses internally, but what crosses the seam MUST leave the key absent from the owned snapshot, so that a required absent key resolves as `owned_state_unavailable` rather than matching against a placeholder value." — (NC, `0004:C8`)

- [REQ-32] "A sentinel value would make an absent key read as *present* and silently retire `owned_state_unavailable` for it — the one encoding that turns this RDR's safety rule into a regression." — (LBD, Absence crosses the seam as omission) — a `Tag{Key: k, Value: "<absent>"}` on `Input.Owned` is a defect; "The spike's `<absent>` string is fixture shorthand, not that seam value."

- [REQ-33] "**Nothing in the type system prevents that** … The prohibition is carried by this contract and by MVV Scenario 6, not by a build error" — (LBD, Absence crosses the seam as omission) — the obligation is a test, not a type.

- [REQ-34] The absent-key row's loudness is conditional and that condition is **not** this RDR's to enforce: "this RDR's seam guarantee delivers `owned_state_unavailable` **only when the consuming row declares the key in `RequiresOwned`**". — (DISP note; CAP row "`Row.RequiresOwned` derived by the normalizer") — a key consumed only through `Row.Match` or a guard yielding escapable `no_match` is **outside** this RDR's surface and MUST NOT be "fixed" here. See deviations D2 (the "escapable `no_match`" vocabulary is a stale-citation site vs RDR 0007's fenced `guard_unevaluable`).

## Gate accessors

- [REQ-35] "A gate accessor MUST return allow, deny, or indeterminate." — (NC, `0004:C9`)

- [REQ-36] "Indeterminate MUST be a refusal-class result, not a false allow and not a false deny." — (NC, `0004:C9`) — the class is `gate_indeterminate` (DISP).

- [REQ-37] Deny is a **typed gate result carrying a reason** at this accessor boundary, and a **refusal** at the CLI. — (DISP, "Gate returns deny | typed gate result | — (deny, with reason)"; FM lists "gate denied" among the typed refusals). Reconciled by JDR 0001 §D9: "Deny is a refusal at the CLI and a typed result at the accessor." See QUESTIONS Q1 and deviations D3.

- [REQ-38] Gate site and aggregation are **not** this RDR's to implement: JDR 0001 §D9 fixes gates as running "after exact-one selection, before the plan", "only the selected row's gates", "deny-overrides with every gate reported", and "`set-state` never runs gates" — landing in RDR 0005, with "**0004** — one clause stating the site by citation and that deny is reported, never applied." — (JDR 0001 §D9; deviations D3) — this RDR MUST NOT apply a deny, only report it.

- [REQ-39] "A gate's timeout or execution failure is an *accessor* refusal (§D11, exit 3), not a gate result." — (JDR 0001 §D9, cited by deviations D3) — a gate that times out yields `timeout`, not `gate_indeterminate`.

## Write accessors — what may be written

- [REQ-40] "A write accessor MUST apply only planned owned-tag writes produced by a successful transition." — (NC, `0004:C10`)

- [REQ-41] "It MUST NOT write observed or recognized tags." — (NC, `0004:C10`) — a definition declaring a write to a non-owned tag MUST fail validation (MVV validation arm `write_non_owned_tag`).

- [REQ-42] "write accessors execute only from a successful transition plan." — (TD) — no write runs on `no_match`, `ambiguous_match`, or any other refusal.

- [REQ-43] An **escaped** plan reaches the write accessor with an empty write set and performs no write. — (deviations D1, carrying RDR 0009's obligation at `0009:1515-1527`; JDR 0001 §D9 "Escape rows carry no `gate` list … an escaped plan is never gated") — Stage 8 MUST author this test in 0004's suite; if the write accessor cannot honour it without a fenced change, escalate.

## `<clear>` as removal

- [REQ-44] "`<clear>` is a reserved tag value (JDR 0001 §D5)." — (NC, `0004:C11`) **[0002-delivered]** — 0002 refuses `<clear>` as an authored tag value at load.

- [REQ-45] "A planned write of `<clear>` MUST remove the key from the artifact" — (NC, `0004:C11`) — assignment of the literal is the defect this clause exists to prevent (LBD, Clear semantics: "the one write whose read-back passes green over a tag that was never removed").

- [REQ-46] "its read-back MUST assert the key is absent — a re-read that still holds the key, including as the literal string `<clear>`, is `read_back_mismatch`." — (NC, `0004:C11`)

- [REQ-47] "Clearing a key the artifact does not hold MUST succeed." — (NC, `0004:C11`) — the idempotent clear (DISP, "Write clears a key the artifact does not hold | success | — (idempotent clear)").

- [REQ-48] "A read that yields `<clear>` as a value MUST treat that key as unreadable and refuse `incomplete_read`." — (NC, `0004:C11`)

- [REQ-49] The clear cases' negative control is the assignment fixture: "The assignment fixture from scenario 3 (`output.txt:8`), which must still assert presence and equality — a read-back that asserts absence for every write fails it" — (ORA 9)

## Write read-back verification

- [REQ-50] "After a write accessor reports command-level success, the executor MUST re-read the same caller-supplied artifact role named by the write binding" — (NC, `0004:C12`)

- [REQ-51] "verify each planned owned tag for equality against the held value — a write replaces the whole value, so containment is not equality" — (NC, `0004:C12`; JDR 0001 §D7(iv))

- [REQ-52] Read-back equality is **byte equality over RDR 0002's canonical JSON array** for a set-valued tag: "a set crosses as its canonical JSON array — members sorted, duplicate-free, compact encoding … so 0004's read-back equality is byte equality and no third encoding exists." — (JDR 0001 §D13, closing JD-22; deviations D3) — the accessor MUST NOT invent a second set encoding or compare set-valued tags as unordered collections.

- [REQ-53] "and that observed and recognized tag values present before the write are unchanged." — (NC, `0004:C12`) — the pre-write snapshot is taken "Before the write" (TD: "the executor records the same caller-supplied artifact role's observed and recognized tag values").

- [REQ-54] "It MUST NOT satisfy read-back verification by discovering an ambient artifact or by reading an unrelated role." — (NC, `0004:C12`)

- [REQ-55] "A read-back mismatch MUST be reported as a write failure." — (NC, `0004:C12`); "A read-back mismatch is a failure even if the write command exited successfully." — (TD)

- [REQ-56] Read-back is **mandatory, not optional**: "writers additionally carry `read_back = true`, since read-back is mandatory, not optional" — (TD) **[0002-delivered]** — the loader already refuses a write entry without `read_back = true`.

- [REQ-57] The two mismatch arms are distinguished by *which* tag moved: "owned (`status`) at `output.txt:9`, non-owned (`profile`) at `output.txt:10`" — (ORA 3), with write success (`output.txt:8`) as the control: "a read-back that always reports mismatch fails this row".

- [REQ-58] "`write -> read = expected owned-tag value identity + protected non-owned tag identity`: after a write accessor reports command-level success, reading the same caller-supplied artifact role named by the write binding must hold each planned owned-tag value exactly (and must not hold a key the plan cleared), and any observed or recognized tag values present before the write must remain unchanged." — (RT)

- [REQ-59] "This is not an undo or byte-for-byte artifact invariant" — (RT); "byte-for-byte artifact identity, formatting, key order, and comments are explicitly out of scope; the accessor observes tag values, not the artifact encoding" — (FID, `write -> read`, artifact as a whole). The implementation MUST NOT assert artifact-level fidelity.

- [REQ-60] "tags absent before the write are unconstrained; the write binding's own planned owned tags are excluded from this comparison by construction" — (FID, `write -> read`, non-owned tags) — the non-owned comparison skips the planned owned keys and does not constrain keys that did not exist pre-write.

## Read-back incompleteness and post-mutation reporting

- [REQ-61] "The read-back re-read is subject to read completeness." — (NC, `0004:C13`) — the re-read goes through the read path, not a map clone; the spike's clone "cannot fail independently of the write" (CA A10).

- [REQ-62] "If it cannot read a key it must compare, the write MUST be reported as `read_back_incomplete` — the verification did not run — and MUST NOT be reported as `read_back_mismatch`, which asserts the artifact is wrong, nor as success." — (NC, `0004:C13`)

- [REQ-63] "`read_back_incomplete` and a post-mutation `timeout` are reported after the write command already ran. The refusal MUST be understood as \"the mutation may have been applied and was not verified\"" — (NC, `0004:C14`)

- [REQ-64] "it MUST NOT be represented to callers as a write that did not occur." — (NC, `0004:C14`)

- [REQ-65] "The accessor layer MUST NOT attempt compensation: it does not retry the write, does not undo it, and does not re-derive the artifact's state." — (NC, `0004:C14`) — asserted "on the invocation count — a layer that retries or undoes shows a second write" (ORA 8); TS 8: "The test asserts the accessor issued exactly one write invocation and performed no undo, retry, or re-derivation".

- [REQ-66] "Recovery is the caller's, and its safe move is to re-read before acting." — (NC, `0004:C14`)

- [REQ-67] The `read_back_mismatch` fixture is the negative control for post-mutation reporting: it "must still assert the artifact is wrong rather than unverified — an implementation that collapses both into one \"write failed\" shape fails that row" — (ORA 8)

## Timeout

- [REQ-68] "Every accessor invocation MUST have a bounded timeout." — (NC, `0004:C15`) — every read, gate, and write invocation, including the read-back re-read.

- [REQ-69] "Timeout MUST be reported as its own refusal class, distinct from execution failure and read-back mismatch." — (NC, `0004:C15`)

- [REQ-70] "Missing or non-positive timeout metadata MUST fail validation before execution." — (NC, `0004:C16`) **[0002-delivered]** — `internal/table/load.go::accessorTable` refuses an absent timeout, a non-duration timeout, and a non-positive timeout.

- [REQ-71] "the Resolve spike wraps every invocation in `context.WithTimeout` … demonstrates timeout as its own refusal (`output.txt:5`), and shows the write path performs one same-role read-back comparison after a command-level success" — (PE) — the bounded-execution property is the non-functional check; no throughput target is set.

## No package prints

- [REQ-72] "The accessor package MUST return structured success/refusal values and MUST NOT write to stdout or stderr directly." — (NC, `0004:C17`)

- [REQ-73] "the no-print property is asserted by capturing stdout/stderr around the call and requiring both empty" and "the accessor package test asserts the returned structured value carries the refusal class" — (ORA 5). "Scenario 5 is the one absence-of-error oracle; the named capture control is what makes it discriminating and is normative for the implementation test." — (ORA) — the capture control is mandatory, not optional; a test that deliberately prints must fail it.

- [REQ-74] "Accessor package must not print." — (EIA, Output routing row) — `internal/cli/respond::Fail` is CLI-only and is not the executor's API.

- [REQ-75] "Expose accessor refusal classes to the CLI layer so RDR 0005 can map them through `respond.Fail` and `clierr.ExitCodeFor` without package-level prints." — (IP, Phase 4) — the mapping itself is RDR 0005's; this RDR ships the exposure.

## The refusal-class set

- [REQ-76] "Visible failures are typed refusals: unknown accessor, capability mismatch, timeout, execution failure, incomplete read, gate denied, gate indeterminate, write attempted for a non-owned tag, read-back mismatch, and read-back incomplete." — (FM) — the wire spellings appear in DISP as `unknown_accessor`, `capability_mismatch`, `timeout`, `execution_failure`, `incomplete_read`, `gate_indeterminate`, `read_back_mismatch`, `read_back_incomplete`.

- [REQ-77] "These are accessor refusal classes and are disjoint from the kernel's closed five-kind `internal/resolve/resolve.go::RefusalKinds` set, which this RDR does not extend" — (FM) — the implementation MUST NOT add a member to `RefusalKinds()`.

- [REQ-78] "mapping accessor refusals onto CLI codes is RDR 0005's, and `internal/cli/clierr` today has no artifact-unavailability code to reuse." — (FM) — `exit code` is out of scope: "RDR 0005 owns the CLI mapping; this table stops at the structured value the accessor package returns." (DISP)

- [REQ-79] "No input class exits silently **at this boundary**: every row either returns a value, mints a named refusal, or — for a genuinely-absent key — omits that key from the owned snapshot." — (DISP) — the Disposition Table is total over the input classes the boundary can meet.

- [REQ-80] "A silent failure is any outcome where the accessor returns a **success-shaped** result over state it did not establish" — (FM) — four named shapes, each with a mandatory guard: an unverified write (read-back), a truncated read reading as absence (completeness), an absent key crossing as a placeholder (seam omission), a clear stored as the literal (reserved sentinel).

- [REQ-81] "One shape is guarded but **not** fully closed by this RDR: a key omitted at the seam is loud only where a surviving row declares it in `Row.RequiresOwned` … the residual is a property of that peer-owned scope, tracked in Capability Dependencies rather than claimed closed here." — (FM) — this residual MUST be left open, not closed here.

- [REQ-82] "Diagnosis starts with the accessor identity, capability, artifact role, timeout, and expected versus observed tag values." — (FM) — a refusal payload carries enough to name those.

## Validation before execution — the eight arms

- [REQ-83] "validation rejects unknown accessor names, missing or multiply-bound accessor identities, capability mismatches, writes to non-owned tags, missing timeout/read-back metadata, non-positive timeouts, a missing or empty read requested-key set, and ambient artifact discovery before resolution." — (TD)

- [REQ-84] "eight validation arms, each asserting its own named code." — (TS 1) — the seven witnessed names are "`missing_accessor`, `multiply_bound_accessor`, `capability_mismatch`, `missing_or_non_positive_timeout`, `missing_write_read_back`, `ambient_artifact_discovery`, `write_non_owned_tag`" (ORA 1); the eighth is the missing/empty requested key set (REQ-17), unwitnessed.

- [REQ-85] The validation oracle asserts each fixture's *own named code* — "not merely that validation returned non-empty" — with the valid fixture returning the empty set as the negative control, "so a validator that rejects everything fails". — (ORA 1)

## Replay stability

- [REQ-86] "Two runs over the same model and fixture results produce the same disposition" — (FID, `resolve -> replay`); "disposition equality", "not byte-identical output; ordering is normalized by sorted map formatting".

- [REQ-87] "The first two runs produce the same disposition; the injected failure produces the same stable refusal class and accessor identity." — (TS 4)

- [REQ-88] "the two identical runs are compared for equality (`replay-identical=true`) rather than for absence of error", with the injected-refusal run as control: "a replay that returns success unconditionally fails it". — (ORA 4)

- [REQ-89] "Accessor execution can be deterministic enough for resolver replay when the model records artifact role, accessor name, capability, timeout, and returned tag values." — (CA A4, Verified) — the recorded disposition names those five.

## Execution shape and phases

- [REQ-90] The illustrative execution order is: validate the model's capability claims; read owned tags from the caller's artifact role; resolve using RDR 0001 and RDR 0003; apply the planned owned-tag write; "Re-read `state` and compare each planned owned tag for equality (absence for a cleared key)"; "Return success or a typed refusal without printing." — (IC) — illustrative shape, normative only in its ordering of write-then-read-back.

- [REQ-91] "Execution has three phases" — validation before resolution, runtime invocation of read and gate accessors with classification, then write accessors from a successful plan. — (TD)

- [REQ-92] "Define the accessor definition structs, capability enum, refusal classes, and validation rules that connect RDR 0002 accessor references to declared read/gate/write bindings." — (IP, Phase 1)

- [REQ-93] "Implement context-bound invocation for typed accessor bindings. The executor returns structured success/refusal values and performs no direct output." — (IP, Phase 2)

- [REQ-94] "Apply planned owned-tag writes through write accessors, then re-read the same artifact role and compare each planned owned tag for equality — absence for a `<clear>` — plus the pre-write observed and recognized tag values. Treat mismatch as a write failure." — (IP, Phase 3)

- [REQ-95] "New internal package should own capability validation and invocation." — (EIA, Accessor executor row: "None found under `internal/`", Decision "Introduce")

- [REQ-96] "The executor is not a shell runner and not a state-machine action callback surface. It invokes typed bindings selected by accessor name and capability, enforces a per-accessor timeout, converts execution errors into stable refusal classes, and never prints directly." — (AP) — no shell-out, no host callbacks (Alternatives 2, 3, 5 rejected).

- [REQ-97] "The resolver stays stateless: it receives the accessor-read values and write plan disposition" — (AP) — the executor MUST NOT make the resolver stateful or orchestrating (CA A5: "The accessor layer becomes an orchestrator" is the If-wrong).

- [REQ-98] "No new third-party dependency is selected at Propose. Resolve may choose a Go test helper or TOML library only if RDR 0002 has not already selected one." — (IP, New Dependencies) — RDR 0002 has already selected the TOML carrier, so no new dependency is expected.

- [REQ-99] "Add stable refusal codes during implementation or RDR 0005." — (EIA, Structured CLI failure row, Decision "Extend") — `internal/cli/clierr::CLIError` codes are append-only.

- [REQ-100] "Accessor/table path discovery belongs to CLI integration." — (EIA, Config loading row) — the executor does not load config.

## Desk-trace step assertions (in force simultaneously)

- [REQ-101] Step 2 holds nine assertions at once on one read invocation: "read returns typed values or a typed refusal; completeness — exactly the requested keys or refuse; requested set is validated definition metadata, never derived from what resolved; one unreadable key refuses the whole read; absence crosses the seam as omission, not a placeholder value; timeout outranks incomplete read; no mutation of authoritative artifacts; bounded timeout; no direct stdout/stderr". — (DESK, step 2)

- [REQ-102] Step 5 holds the full read-back assertion set at once: "re-read the *same* caller-supplied role named by the write binding; the re-read is itself subject to read completeness — an unreadable compared key is `read_back_incomplete`, not mismatch and not success; a post-mutation refusal reads as applied-but-unverified, never as a write that did not occur, and the accessor performs no retry, undo, or re-derivation; verify each planned owned tag for equality and a cleared key for absence (a stored literal `<clear>` is a mismatch); verify pre-write observed/recognized values unchanged; no ambient or unrelated-role read; mismatch is a write failure". — (DESK, step 5)

- [REQ-103] "No CONTRADICTION row, on two pairs." The C4 disjunction and the C6 completeness guarantee are jointly satisfiable "because completeness constrains *which* branch the disjunction takes rather than adding a third branch"; the genuine-absence success and `missingOwned`'s map-presence decision "collide only if absence reaches the resolver as a present tag, which the seam clause now forbids". — (DESK) — the implementation MUST NOT introduce a third read branch.

## Risk mitigations that are testable obligations

- [REQ-104] "Require same-role read-back verification against expected owned-tag values." — (RM, wrong-artifact/wrong-tag risk)

- [REQ-105] "Completeness is normative — a read that cannot resolve every requested key refuses, and the MVV asserts the truncated case takes the refusal branch." — (RM, truncated-read risk)

- [REQ-106] "Absence crosses as omission from the owned snapshot; the MVV carries an absent required key through to the resolver and asserts the refusal still fires." — (RM, placeholder risk)

- [REQ-107] "The contract states the refusal means \"may have been applied, not verified\", forbids accessor-level compensation, and puts recovery on the caller via re-read. No undo is claimed" — (RM, applied-but-unverified risk)

- [REQ-108] "`<clear>` is reserved (JDR 0001 §D5); a clear is a removal, read-back asserts absence, and MVV Scenario 9 asserts a stored literal fails." — (RM, clearing-write risk)

- [REQ-109] "Make them separate refusal classes and verify CLI mapping." — (RM, collapsed-timeout/gate-indeterminate risk)

- [REQ-110] "Validation records capability, artifact role, and tag keys" — (RM, smuggled-shell risk) — the validated definition retains all three for diagnosis.

## Testing strategy — production tests, not spike carry-over

- [REQ-111] "The MVV should become production tests around the accessor validator and executor boundary." "Done means those spike cases become package tests without direct stdout/stderr output from the accessor package." — (TS)

- [REQ-112] "The write success tests must assert the re-read owned-tag value equals the transition plan's expected value and that pre-write observed and recognized tag values on the same artifact role remain unchanged." — (TS)

- [REQ-113] "Matching capability references pass; unknown accessors, duplicate bindings, capability mismatches, missing timeout metadata, and write attempts against non-owned tags fail before resolution." — (TS 1)

- [REQ-114] "Successful reads return the complete typed tag values for every requested key, gate allow/deny returns typed gate results, and timeout, execution failure, incomplete read, and gate indeterminate remain distinct refusal classes." — (TS 2)

- [REQ-115] "an artifact that genuinely lacks a requested key returns it as an absent value, not a refusal." — (TS 2)

- [REQ-116] "Owned-tag values equal to the plan report success; mismatched owned-tag values or mutated non-owned observed/recognized tag values report read-back mismatch even when command-level write invocation succeeded. A read-back whose re-read cannot read a key it must compare reports `read_back_incomplete` — neither success nor mismatch (A9)." — (TS 3)

- [REQ-117] "The accessor package returns structured values only; the CLI layer can map them through `respond.Fail` and `clierr.ExitCodeFor`." — (TS 5)

- [REQ-118] "The read succeeds at the accessor boundary, the absent key is omitted from the owned snapshot rather than carried as a placeholder value, and the resolver refuses `owned_state_unavailable` naming that key." — (TS 6) — the assertion is on the **resolver's** disposition, not the accessor's branch (ORA 6), with a carried key as the control so the refusal cannot pass vacuously.

- [REQ-119] "The refusal is `timeout`, not `incomplete_read` — the two input classes overlap and timeout takes precedence." — (TS 7)

- [REQ-120] "The refusal is `read_back_incomplete` and carries the applied-but-unverified sense rather than reading as a write that did not occur." — (TS 8)

- [REQ-121] "The first two report success and the re-read does not hold the key; the third reports `read_back_mismatch`; the read refuses `incomplete_read` naming the key. The scenario-3 assignment fixture remains the control — it must still assert presence and equality." — (TS 9)

## Prerequisites and pending assumptions

- [REQ-122] "**A9, A10, and A11 Pending** — MVV-proven properties carried to lock with the MVV as the named implementation-time plan" — (IP, Prerequisites; CA) — Stage 8 discharges A9 (four independent rules), A10, and A11 by running the named MVV scenarios; "Verify and retire them per-rule rather than flipping A9 as a unit, and split A9 if any single rule blocks." (CA A9, Note on bundling)

- [REQ-123] A9's four rules, each verified by a named scenario: the missing/empty requested-key-set validation arm (Scenario 1's eighth arm), absence-to-resolver (Scenario 6), timeout-outranks-incomplete (Scenario 7), and read-back-incomplete (Scenario 3's unreadable-compared-key case). — (CA A9, Evidence)

- [REQ-124] "The existing Resolve spike does not witness any of the four: it derives its default key set from the artifact, holds absence as an in-map sentinel, sleeps before its key loop so timeout and truncation never overlap, and re-reads by cloning the tag map without going through `read`." — (CA A9) — each of those four spike shapes is a shape the implementation MUST NOT reproduce.

- [REQ-125] "no fixture double can exercise it — a re-read that cannot fail independently of the write has nothing to report." — (CA A10) — the read-back fixture MUST be able to fail independently of the write.

- [REQ-126] Prerequisites carried unchecked at Gate PASS: RDR 0001's stateless resolver returning planned owned-tag writes, RDR 0002's accessor references / provenance / artifact roles through normalization, RDR 0003's consumption of accessor-produced values without executing accessors. — (IP, Prerequisites; deviations D2 Note)

## REQ-MVV

- [REQ-MVV] `0004:MVV` "Build a fixture flow with one read accessor, one gate accessor, and one write accessor over caller-supplied artifact roles. Prove success, timeout, gate denied, gate-indeterminate, execution failure, incomplete read, capability mismatch, unsafe definition validation, write read-back-mismatch, and `read_back_incomplete` dispositions, plus the two A9 boundary cases: an absent required key refused as `owned_state_unavailable` at the resolver, and a timed-out partial read classified as `timeout`; plus A10's case: a write whose re-read cannot read a compared key surfaces the applied-but-unverified sense and the test asserts no second write and no compensating action occurred; plus A11's cases: a clearing write whose re-read shows the key absent, a clear of a key the artifact does not hold succeeding, a clear whose re-read still holds the key (including as the literal) refused as `read_back_mismatch`, and a read holding the literal `<clear>` refused as `incomplete_read`. The write success test must assert the re-read owned-tag value equals the transition plan's expected value. The read test must assert that an accessor which can resolve only some of the requested keys takes the refusal branch and is distinguishable from one whose artifact genuinely lacks those keys." — (MVV)

REQ-MVV is the gating validation. It decomposes into nine numbered scenarios
(TS 1–9), each with a named negative control (ORA). The nine are the unit of
Stage 8 completion; a scenario passing only by absence of an error does not
discharge it (ORA 5's capture control is normative, REQ-73).

---

## ASSUMPTIONS

- ASSUMPTION: "flow id" in `0004:C2`'s identity triple is the loaded model's
  identity, not a new field this RDR introduces. RDR 0002 owns the carrier
  (REQ-8) and `internal/table/model.go` carries one model per loaded document,
  so the triple is `(model, accessor name, capability)` in practice. The
  discriminating obligation — a same-id read/write pair is two identities, not a
  rebinding (REQ-4) — is unaffected by which of the two readings holds.

- ASSUMPTION: "validation before execution" (REQ-83) is satisfied where RDR
  0002's loader already performs the check (`internal/table/load.go::accessorTable`
  and the binding counts). Stage 8 asserts the named code at this RDR's boundary
  and does not duplicate the loader's check. Basis: LBD Wire/byte format and
  Definition shape both defer the carrier and its layout validations to 0002;
  JDR 0001 P6 forbids restating a peer's clause.

- ASSUMPTION: the timeout in REQ-68 ("Every accessor invocation") bounds the
  read-back re-read as a distinct invocation, not as time borrowed from the
  write's budget. Basis: DISP separates "Write command runs, then the invocation
  exceeds its timeout" as its own row minting a post-mutation `timeout`, which is
  only reachable if the re-read is inside a bounded invocation of its own.

- ASSUMPTION: the eighth validation arm's code (REQ-17/REQ-84) follows the
  seven witnessed names' spelling convention (snake_case, defect-named); the
  record names the arm but not its code. Proceeding with a name in that family
  chosen at implementation. Basis: ORA 1 requires each arm to assert "its *own
  named* code" but fixes only seven spellings.

- ASSUMPTION: the "absent value" a read returns on its success branch
  (REQ-115) is an internal representation only; what the executor hands the
  resolver omits the key (REQ-31). The two clauses are one rule stated at two
  layers, not two conflicting requirements. Basis: `0004:C8` ("MAY represent
  absence however it chooses internally"), FID's `read` row, and LBD's "The
  spike's `<absent>` string is fixture shorthand, not that seam value."

- ASSUMPTION: `gate denied` at the accessor boundary is a typed gate result
  carrying a reason (REQ-37), and the FM sentence listing "gate denied" among
  "typed refusals" is naming the *caller-visible* disposition, not this
  boundary's return shape. Basis: JDR 0001 §D9 states both halves explicitly —
  "Deny is a refusal at the CLI and a typed result at the accessor" — and
  deviations D3 names §D9 as the decider. See QUESTIONS Q1.

- ASSUMPTION: read-back equality for a set-valued tag is byte equality over the
  canonical JSON array RDR 0002 produces (REQ-52), so the executor performs no
  set-aware comparison of its own. Basis: JDR 0001 §D13 ("0004's read-back
  equality is byte equality and no third encoding exists"), landing in 0004 by
  citation.

- ASSUMPTION: RDR 0009's escaped-plan / empty-write-set test (REQ-43) is
  additive to this RDR's suite and changes no clause's meaning, so it is
  authored here without a fenced change. Basis: deviations D1 condition (c) —
  "no clause's meaning changes — the test is additive" — and D1's own Stage 8
  check, which escalates only if the write accessor cannot honour it.

- ASSUMPTION: the stale-citation sites named in deviations D2 (a guard-consumed
  absent key described as falling to "escapable `no_match`", RDR 0007's A6b
  described as Pending, the pre-provenance-scope restatement of §D7(ii)'s
  zero-readers rule, the vacuous RDR 0008 citation) are **citation repair in the
  record**, not implementation obligations. RDRs are never amended, so Stage 8
  records the disposition in `deviations.md` and implements against RDR 0007's
  fenced `guard_unevaluable` and RDR 0002's shipped scoping — not against the
  stale prose. REQ-34 is worded to that effect.

---

## QUESTIONS

Two items are recorded here because the orchestrator counts open author
decisions. Both are resolved for Stage 8 by the grounding named in the matching
ASSUMPTION above, and neither blocks: this run is unattended and continues under
the best-supported reading.

- **Q1 — Is `gate denied` a refusal or a typed non-refusal result at the
  accessor boundary?** The record says both: the Disposition Table gives "Gate
  returns deny | typed gate result | — (deny, with reason)", while Failure Modes
  lists "gate denied" among "Visible failures are typed refusals". Deviations D3
  names this split explicitly and makes it a Stage 8 check: "the … non-refusal
  reading of `gate denied` agrees with §D9 or is repaired. If the … split cannot
  be reconciled by citation, escalate — §D9 decides it, so a surviving
  contradiction is a fenced conflict." **Proceeding on the layered reading**
  (REQ-37): a typed result carrying a reason at the accessor, a refusal at the
  CLI. JDR 0001 §D9 states exactly that and adds the reasoning — "0004's table
  says the accessor *answered*; 0005 says `resolve` returns exactly one plan or
  exactly one refusal, and a denied resolution has no plan." Both of the RDR's
  sentences are then true at their own layer and the two-reading materiality
  disappears. Escalation is required only if an implementation attempt shows the
  layered reading unbuildable.

- **Q2 — Does this RDR own the gate *site* (when gates run) at all?** The RDR is
  silent on the site; JDR 0001 §D9 records that "0004 never states the site" and
  assigns the site to RDR 0005, leaving 0004 "one clause stating the site by
  citation and that deny is reported, never applied." The materially different
  readings are (a) 0004's executor invokes gates at a site it chooses, and (b)
  0004's executor invokes a gate only when called, with sequencing owned by the
  caller. **Proceeding on (b)** (REQ-38): §D9 lands the site in 0005; this RDR's
  own AP says "The resolver stays stateless" and CA A5's If-wrong is "The
  accessor layer becomes an orchestrator", both of which forbid (a). Stage 8
  therefore ships gate *invocation* semantics (allow/deny/indeterminate,
  timeout, no-apply) and no gate scheduling. No predecessor artifact contradicts
  this; RDR 0005 is unimplemented, so nothing here can be validated against a
  shipped caller yet.
