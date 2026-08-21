Model: claude-fable-5

# Critique — RDR 0004 re-entry (read-completeness surface), iteration 2

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | §Normative Contracts, seam clause ("MUST leave the key absent from the owned snapshot, so that a required absent key resolves as `owned_state_unavailable`") + §Disposition Table row "Requested key absent from artifact" ("loud one layer down") | Absence-as-omission only reaches `owned_state_unavailable` through `Row.RequiresOwned`, whose producer is nobody (JDR 0001 JD-3, still an open "Direction" item in RDR 0002). A key consumed only via `Match` degrades absence to `no_match`, which `escapeOrRefuse` can rescue into a plan — missing artifact state masked behind an escapable refusal class, the exact P1 failure the cluster exists to prevent. The Disposition Table's "No input class exits silently" claim is false for this row. | A flow whose artifact genuinely lacks an owned key takes the escape edge and *writes* — a transition committed against state that was never there; no refusal, no diagnosis | §1-W1 |
| C-2 | §Normative Contracts completeness clause ("exactly the keys it was asked for — no requested key missing") vs. seam clause ("leave the key absent from the owned snapshot") + §Fidelity Table row 4 + §Proportionality ("RDR 0001's `Input` shape, which is unchanged") vs. A9 If-wrong ("moving the absence encoding into RDR 0001's `Input` shape") | The success branch is bound by two contradictory shape contracts and no component is named to perform the projection between them; the accessor's public success type is unspecified. A9's own fallback is to reopen the locked peer that Proportionality declares untouched. | Two implementers build two incompatible success types (sentinel map vs. omission map); the consumer contract "returned tag set is exactly the requested keys" is false at the only place a consumer exists | §1-W2 |
| C-3 | §Normative Contracts timeout-precedence clause ("`timeout` takes precedence over `incomplete_read`") + A4 ("deterministic enough for resolver replay") + §Oracle Discriminability scenario 7 ("the fixture makes both classes true at once") | Timeout is the one refusal class that is a function of wall-clock, not of the input tuple; the precedence rule entangles the deterministic class with the nondeterministic one. The scenario-7 fixture is a contrived simultaneity a real binding produces only as a race. | Same artifact, same flow: `incomplete_read` on a fast machine, `timeout` on a slow one; flaky MVV test in CI; A4's replay claim silently narrowed to "deterministic except when it isn't" | premortem |
| C-4 | §Normative Contracts `read_back_incomplete` clause ("the verification did not run") + §Round-Trip / Inverse Invariants ("undo is not claimed") + §Disposition Table (no row for a write that times out after mutation) | The write has already mutated the authoritative artifact when `read_back_incomplete` (or a post-mutation `timeout`) is reported. No retry policy, no idempotency requirement, no modeled disposition for the resulting state: the owned tag advanced while the resolution reported failure, and the transition table has no edge out of the advanced state for the same recognized outcome. | User re-runs the flow after a transient re-read blip and gets `no_match` (or an escape edge) forever; recovery is hand-editing an authoritative artifact the safety model exists to protect | §1-W3 |
| C-5 | §Failure Modes ("unknown accessor, capability mismatch, artifact unavailable, timeout, execution failure, incomplete read, gate denied, gate indeterminate, write attempted for a non-owned tag, read-back mismatch, and read-back incomplete") vs. §Disposition Table ("Artifact role not supplied → `execution_failure`"; "Gate returns deny → typed gate result"; non-owned write → validation failure) | The refusal-class enumeration disagrees with itself across three sections: artifact-unavailable is both its own class and `execution_failure`; gate-denied is both a refusal and explicitly not one; non-owned-write is both a runtime refusal and a pre-runtime validation code. Collapsing artifact-unavailability into `execution_failure` is the exact vagueness the RDR's own Risks table condemns. | A user who forgot to supply an artifact role sees `execution_failure` (an internal-sounding code, wrong `clierr` group); RDR 0005's mapping drifts because there is no single authoritative class list | §2 |
| C-6 | A8 + §Load-Bearing Decisions "Absent vs unreadable is the binding's call, defaulting to unreadable" ("The Resolve spike declares unreadability as a fixture field rather than deriving it, so it proves the branch machinery, not a binding's classification") | A8 was re-stamped Verified on evidence the RDR itself concedes does not test the claim: the spike branches on a declared boolean; no binding ever *classifies*. Real bindings (a `gh` 404, a `git config` exit 1, a missing TOML file) cannot distinguish absent from unreadable, and the mandated default-to-unreadable turns every optional key into an `incomplete_read` refusal. | Flows with optional keys refuse constantly; users respond by writing bindings that guess absence — the lie the contract exists to prevent, now incentivized by the contract | §3 |
| C-7 | §Normative Contracts requested-key-set clause + §Technical Design validation list (eight arms, none of which is "requested set covers table consumption") | Completeness is relative to the declared requested set, and nothing ties that set to the keys the table's rows consume. A table edit that adds a key a read definition never requests passes all eight validation arms and fails at runtime as absent state. | After a routine table change, a flow refuses `owned_state_unavailable` (or, per C-1, silently escapes) on artifacts that carry the key perfectly well; diagnosis points at the artifact, the defect is config drift | AT-7 |
| C-8 | A9 ("implementable as stated: absence … a validated requested-key set … `timeout` outranks … `read_back_incomplete`") — Status Pending, Method MVV Test | Four independent seam rules bundled into one unverified assumption, carried Pending into lock with verification deferred to the implementation that depends on them. If any one fails, all four are formally refuted together, and the stated fallback reopens RDR 0001. This is a lock on a promissory note. | Implementation discovers rule 1 (seam omission, per C-1/C-2) is the wrong shape mid-build; the RDR is reopened a third time and the other three rules are relitigated with it | §2 |

---

## 1. The three most likely ways implementation goes wrong

### W1 — The absence seam clause routes missing state into the escape path, and a plan is produced against state that does not exist (C-1)

**Root cause.** The revised seam clause stakes everything on one mechanism: "what crosses the seam MUST leave the key absent from the owned snapshot, so that a required absent key resolves as `owned_state_unavailable` rather than matching against a placeholder value." The word doing all the work is *required*. In the kernel as it exists on main, `owned_state_unavailable` fires in exactly one place: `internal/resolve/resolve.go::missingOwned`, which quantifies over `Row.RequiresOwned`. It does not fire for a key a row consumes through its `Match` pattern — `TagSet.matches` returns false for a missing key, the row is silently excluded from candidacy, and the disposition is `no_match`. And `no_match`, unlike `owned_state_unavailable`, is escapable: `escapeOrRefuse` will hand it to any escape row whose looser `Match` holds, and emit a `Plan` — with `Writes` — for a resolution whose triggering fact is that the artifact was missing state.

**The passage that enabled it.** The Disposition Table row "Requested key absent from artifact → success — loud one layer down: a row requiring the key refuses `owned_state_unavailable` naming it," and the closing boast "No input class exits silently." Both are true only in a world where every owned key a row consumes appears in `RequiresOwned`. JDR 0001 JD-3 says, in so many words, that no layer is obliged to populate that field — it appears zero times in RDR 0002's normative surface, and RDR 0002's own re-entry notes still carry "Name the `RequiresOwned` producer (JD-3)" as unfinished direction. This RDR re-locked its central new safety claim on a field with no producer, cited `missingOwned` as if citing it created the obligation, and never wrote the clause that would: *every owned key a row's Match or Guard consumes MUST appear in RequiresOwned, or omission-encoded absence is unsound.* The MVV cannot catch this: Scenario 6's fixture will populate `RequiresOwned` by hand, pass, and prove nothing about production tables.

**Symptom.** A user's flow runs against an artifact that genuinely lacks an owned key. No refusal. The escape edge fires, `Escaped: true`, the write accessor mutates the authoritative artifact. The one failure JDR 0001 P1 names — "missing artifact state masked behind an escapable refusal class" — shipped as the designed behavior of the absence encoding, and the user finds out when the artifact is already wrong.

### W2 — Two contradictory shape contracts bind the success branch, and nobody owns the projection between them (C-2)

**Root cause.** The completeness clause says the success branch returns "the tag set for exactly the keys it was asked for — no requested key missing, no unrequested key added." The seam clause says a genuinely-absent requested key "MUST leave the key absent from the owned snapshot." These cannot describe the same value. So there are two values — an accessor-internal return carrying absence somehow, and a seam-crossing owned snapshot with the key omitted — and the RDR names no component that performs the projection, no type for either value, and no rule for what the executor does between them. The Fidelity Table's row 4 asserts key-set equality ("none missing, none added") for "`read` over a requested key set" while its own exemption column admits the seam value omits keys. The Load-Bearing Decision "Consumers may therefore treat a returned tag set as exactly the requested keys" is flatly false for the only consumer the system has: the resolver receives `Input.Owned` with absent keys omitted, and cannot distinguish "absent" from "never requested" from "requested-key-set drift" (see C-7).

**The passage that enabled it.** The two normative clauses quoted above, plus Proportionality's assertion that the seam clause "constrains this RDR's own output rather than reopening RDR 0001's `Input` shape, which is unchanged" — sitting forty lines from A9's If-wrong, which names the likely fix as "moving the absence encoding into RDR 0001's `Input` shape (JDR 0001 §D3 option (c))." The RDR's plan of record and its own contingency plan contradict each other about whether the locked peer is on the table.

**Symptom.** Implementation stalls at the first real type signature: does `read` return `map[string]Value` with an absence variant, or the omission-shaped map, or both? Whichever is chosen, one of the two normative clauses is violated at the letter, the reviewer flags it against the RDR, and the resolution is the third reopening — most likely landing on option (c), the error channel on `Input.Owned` that this revision existed to avoid.

### W3 — `read_back_incomplete` and post-mutation timeout leave the artifact advanced with no modeled way forward (C-4)

**Root cause.** The `read_back_incomplete` clause is honest about epistemics — "the verification did not run" — and silent about consequences. By the time it is reported, the write command has succeeded and the authoritative artifact is mutated. The RDR claims no undo ("undo is not claimed"), mandates no retry of the re-read, requires no idempotency of write accessors, and models no disposition for the state it just created: owned tag advanced, resolution reported failure. The same hole exists for `timeout` on the write path — the Disposition Table's only timeout row is the generic "Accessor exceeds its timeout → `timeout`," which cannot distinguish "nothing happened" from "wrote, then died before verifying." The spike dodges this by checking `ctx.Done()` *before* mutating (`main.go::write`), so no fixture ever witnesses the dangerous ordering.

**The passage that enabled it.** The `read_back_incomplete` normative clause, which specifies only what the refusal must *not* be called (`read_back_mismatch`, success) and nothing about what the caller holds afterward; and the Failure Modes section, which lists the class among "visible failures" as if visibility were the whole problem.

**Symptom.** A transient blip on the re-read — the precise event `read_back_incomplete` exists for — and the user's RDR artifact now says `status=Final` while the flow reports failure with a nonzero exit. The user re-runs. The read now returns `status=Final`; the transition table has no edge from Final for the same recognized outcome; the kernel refuses `no_match` — or worse, an escape edge fires (C-1). The recovery path is a human editing the authoritative artifact by hand, which is the outcome this entire RDR was written to make impossible.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**The read-completeness cluster in §Normative Contracts — specifically the seam-omission clause and its Load-Bearing Decisions companion "Absence crosses the seam as omission" — together with the Fidelity Table row that contradicts them.**

It will be rewritten because it is the only part of the revision that legislates a data shape it never defines, across a boundary whose other side is locked. Within the first two weeks of implementation, someone has to write the Go type that the read accessor returns and the function that turns it into `Input.Owned`. At that moment all three latent defects fire at once: the exactly-requested-keys clause and the omission clause bind the same value differently (C-2); the omission encoding turns out to be unsound against `no_match` escape without a `RequiresOwned` producer nobody has (C-1); and the resolver is discovered to need a three-way distinction — absent, unread, never-requested — that omission collapses to one bit. The RDR has already written its own epitaph for this section: A9's If-wrong names JDR 0001 §D3 option (c), the error channel on `Input.Owned`, as the fallback. That is not a contingency, it is a forecast. The section will be rewritten to adopt (c) or something shaped like it, Proportionality's "RDR 0001's `Input` shape, which is unchanged" will be quietly deleted, and the Fidelity Table will get a fourth exemption.

A9 itself accelerates this (C-8): bundling four independent rules into one Pending assumption means the first rule to fail formally refutes the bundle, and the re-entry machinery reopens all four. The section was locked on the promise that the MVV would vindicate it; the MVV is built by the implementation that needs the answer first.

The runner-up is §Failure Modes (C-5), which will be rewritten the day RDR 0005 tries to build its code mapping and discovers three sections of this RDR disagree about what the refusal classes even are — but that is a two-hour fix. The seam cluster is a design change.

---

## 3. The one assumption that will not survive first contact with a real user

**A8 — specifically its load-bearing corollary, "Absent vs unreadable is the binding's call, defaulting to unreadable."**

A8 is stamped Verified, and the verification is real for what it tests: an if-statement branches on a fixture boolean (`main.go`'s `artifact.unreadable` field, consumed in `main.go::read`). The RDR concedes this in its own Load-Bearing Decisions: the spike "declares unreadability as a fixture field rather than deriving it, so it proves the branch machinery, not a binding's classification." The branch machinery was never in doubt. The claim that needed verifying — that a *real binding* can establish, at its own boundary, that a key is genuinely missing rather than unread — was not verified, and it is false for the bindings this system will actually get.

The accessors named by this cluster's own motivating flows are things like `gh`, `git`, and file reads over TOML. A `gh` API 404: absent resource, revoked token, or wrong URL? A `git config` call exiting 1: unset key or repo error? The binding frequently has one bit — nonzero exit — where the contract demands two. The RDR's answer is the conservative default: "A binding that cannot tell the two apart at its own boundary MUST report `incomplete_read`." Follow that rule and every optional key read through a CLI-backed binding refuses, every time, on every artifact that legitimately lacks the key. Optional keys become unusable; flows that model them deadlock on `incomplete_read`.

The first real user hits this in the first week, and their rational move is the catastrophic one: write the binding to *guess absence* on exit 1, because that makes the flow run. The contract's conservative default doesn't prevent the guess — it creates the incentive for it, and pushes the lie below the seam where no validator, no MVV scenario, and no refusal class can see it. An assumption whose enforcement mechanism manufactures its own violation does not survive contact.

---

## 4. Premortem

*Written 2026-10-02, six weeks after the accessor executor shipped.*

We shipped the accessor package in week two. It failed three ways, in order.

**Week 2 — the type that couldn't exist.** Implementation stopped on day three because `read`'s success type cannot satisfy both normative clauses. We chose the omission-shaped map — return exactly the resolvable keys, omit the absent ones — because that is what `resolve.Input.Owned` takes and nobody wanted to write a projection layer the RDR never named. Review flagged it: the returned set is no longer "exactly the requested keys," the Fidelity Table's key-set equality is violated, and the consumer promise in Load-Bearing Decisions is dead text. We shipped anyway with a TODO pointing at A9. That TODO is now the third reopening of RDR 0004, and the fix on the table is JDR 0001 §D3 option (c) — the error channel on `Input.Owned` that Proportionality swore was out of scope.

**Week 4 — the escape that wrote.** The pilot flow is the RDR lifecycle itself: `state.read` pulls `status` and `review`, a gate checks the evidence directory, `state.persist` writes `status=Final`. The table's finalize row matched on `status` via `Match`; nobody populated `RequiresOwned` because RDR 0002 never obliged anyone to (JD-3 was still "Direction"). A user ran the flow against a freshly seeded RDR whose skeleton lacked the `review` key. The read succeeded — absence is a value, the key was omitted from the snapshot, exactly as the seam clause commands. `TagSet.matches` failed the finalize row on the missing key. `Resolve` counted zero candidates: `no_match`. And the table had a no_match escape row — a catch-all "park it" edge with a loose match — so `escapeOrRefuse` gated it, found it viable, and returned a plan. `state.persist` wrote the park state onto an artifact whose actual condition was *we never had the state to decide*. The user saw a green transition. We found it two weeks later diffing artifacts. `owned_state_unavailable` — the refusal the Disposition Table promised would fire "loud one layer down," naming the key — fired zero times in six weeks of production, because nothing produces `RequiresOwned`.

**Week 5 — the stuck Final.** A user finalized an RDR over a network filesystem. `state.persist` wrote `status=Final`; the read-back re-read hit a transient stale-handle error on one compared key. Per the clause, we correctly reported `read_back_incomplete` — not mismatch, not success — exit 2. The user did the only reasonable thing and re-ran. The read returned `status=Final`; the table has no finalize edge out of Final; `no_match`; and this flow's table had no escape row, so: refusal, forever. The artifact says Final, the flow says it never finalized, and the recovery was `sed` on an authoritative artifact. In the same week the MVV's scenario-7 test — timeout outranking `incomplete_read` — flaked twice in CI, because the "both classes true at once" fixture is a race between the deadline and the unreadable-key determination, and on a loaded runner it classifies the other way. We marked it flaky. A4's replay determinism now has an asterisk nobody wrote down.

Meanwhile the `gh`-backed kata accessor never distinguished a 404 from an auth failure, so every optional key refused `incomplete_read` until its author made the binding treat exit 1 as absence — the exact guess the contract forbids, invisible below the seam. And RDR 0005's mapping review found three incompatible refusal-class lists in the RDR and picked one, so `execution_failure` is what a user sees when they forget to pass an artifact role.

Every one of these was visible in the document at review time. None required running code to find.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 (catches C-1, the escape mask).**
```gherkin
Given a transition table shaped like production output:
    a finalize row matching on owned key "review" via Match, with RequiresOwned empty
    and a no_match escape row with a loose Match
  And an artifact that genuinely lacks "review"
When the read succeeds with "review" omitted from the owned snapshot (per the seam clause)
  And the snapshot is passed to resolve.Resolve
Then the disposition MUST be a refusal naming "review"
  And MUST NOT be a Plan
```
On main this test fails: the escape row rescues `no_match` and emits a plan. Running it at review time — a desk trace suffices, no code needed — either forces a normative `RequiresOwned`-population clause into this RDR or refutes the omission encoding. MVV Scenario 6 as written cannot catch this because its fixture populates `RequiresOwned` by hand.

**AT-2 (catches C-2, the undefined seam type).**
Plain steps: at review, demand the single Go type of the read success value and the name of the component that converts it to `Input.Owned`. Then check both normative clauses against that one type: "exactly the requested keys" and "absent key omitted from the owned snapshot." No single value satisfies both; the review either names the projection component and rewrites one clause to bind it, or the contradiction is on record before lock instead of after.

**AT-3 (catches C-6, the classification fiction).**
```gherkin
Given a read binding over a real substrate (file, git config, or gh),
    not a fixture with a declared `unreadable` field
When the substrate yields: (a) a parseable artifact lacking the key,
    (b) an unparseable artifact, (c) a command exiting 1 with no further signal
Then (a) MUST resolve as absence, (b) MUST refuse incomplete_read,
  And for (c) the binding MUST have a documented basis for its classification
```
Case (c) has no passing answer under the RDR's rules — which is the finding. Requiring this fixture before re-stamping A8 Verified would have exposed that the spike's `unreadable: map[string]bool` verifies the branch, not the claim.

**AT-4 (catches C-3, timeout-precedence nondeterminism).**
```gherkin
Given the scenario-7 fixture where a read resolves some keys then exceeds its deadline
When the fixture is executed 100 times
Then all 100 dispositions MUST carry the same refusal class
```
This test cannot be written to pass reliably without pinning *when* classification happens (e.g., classify solely on deadline expiry at return, never mid-loop) — a rule the RDR does not state. Its absence would have surfaced as an unanswerable question in review.

**AT-5 (catches C-4, the stuck post-write state).**
```gherkin
Given a write accessor whose re-read fails transiently on one compared key
When the write executes from a successful plan
Then the disposition is read_back_incomplete            # passes as specified
When the same flow is re-run against the now-mutated artifact
Then the resolution MUST reach a modeled disposition that a user can act on
```
The second half fails: the advanced artifact has no edge for the same recognized outcome. Asking "what does the user run next?" at review time forces either a re-verification verb, an idempotent-write requirement, or a modeled already-applied disposition — any of the three, but chosen on paper, not in an incident.

**AT-6 (catches C-5, the enumeration drift).**
Plain steps: extract the refusal-class set from (1) §Failure Modes, (2) the Disposition Table's "Refusal class minted" column, (3) the spike's `refusal` constants plus `read_back_incomplete`. Assert set equality. Fails today on `artifact unavailable` (own class vs. `execution_failure`), `gate denied` (refusal vs. typed result), and `write attempted for a non-owned tag` (runtime vs. validation). A five-minute mechanical check that the Finalization Gate's "mechanical pre-sweep" should have included.

**AT-7 (catches C-7, requested-set drift).**
```gherkin
Given a table row that consumes owned key "review"
  And a read accessor definition whose validated requested-key set is ["status"]
When definition validation runs
Then validation MUST fail naming the uncovered key
```
Fails: no validation arm checks coverage of table consumption by requested-key sets — the eight arms all validate definitions in isolation. Absent this arm, completeness is proven over a set nothing checks is sufficient, and the RDR's "self-fulfilling completeness" critique of the spike applies, one level up, to its own contract.

**AT-8 (catches C-8, the bundled Pending assumption).**
Plain steps: at the Finalization Gate, require each Pending assumption to name exactly one falsifiable rule (the same discipline the Method vocabulary demands of Verified ones). A9 fails the check four times over; splitting it into A9a–A9d would have revealed that A9a (seam omission) is not an MVV question at all — it is refuted on paper by AT-1/AT-2 — and the lock would have been held on the one rule that deserved to hold it.
