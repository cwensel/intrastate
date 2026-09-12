Model: claude-opus-5

# Critique (Devil's Advocate Premortem) — cli/0029

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0029:C4` | The seam census obliges `respond::Types() []string` over `{ok, failed}` for the `frozen` envelope `type` vocabulary, while `0029:C1`, `0029:S2`, the Illustrative Code and the MVV all state a refusal carries NO `type` key at all. The seam enumerates a member the wire never emits. | The `frozen` by-value test (S6) pins `{ok, failed}`; a consumer reads the published register, writes `switch type { case "ok": …; case "failed": … }`, and its refusal branch is dead code. Refusals fall through to the default arm and are reported as an unknown envelope — or worse, as a no-op success, which is the exact defect `readPlan` already has a 30-line comment defending against. | §1, §3, premortem, AT-1 |
| C-2 | `0029:C1` | `schema_version` is added non-`omitempty` to the *bare* `CLIError`, and `internal/cli/flow_input.go::readPlan` discriminates a piped `--plan` document by *absence of top-level `type` plus presence of `code`*. The RDR asserts `planEnvelope` "survives by construction" because it is a tolerant four-field decode — true for the DECODE, false for the DISCRIMINATION, which the RDR never examines. `jq .data` over a refusal still yields `null`; but `jq 'del(.schema_version)'`-style filters, and any consumer that now believes `schema_version` is the shared discriminator (C1 says so verbatim: "`schema_version` is therefore what a consumer branches on to know which record shape it holds"), get the opposite rule from the one the code implements. | A consumer follows C1's stated discriminator, branches on `schema_version`, finds it on BOTH records, and cannot tell success from refusal. Under `--plan` piping the pipeline silently applies zero writes and reports success for a transition the model refused. | §1, §3, premortem, AT-2 |
| C-3 | `0029:C4` | `clierr::ExitCodes() []int` over `{0,1,2,3,130}` is declared a `frozen` vocabulary with a by-value seam, but `ExitCodeFor` produces `0` for `GroupSuccess`/`GroupWarning` (not a refusal), and `1` only as the *fallthrough for a non-`CLIError`* — i.e. a panic-adjacent path that emits no envelope. Two of the five "members" are not emitted vocabulary at all. Meanwhile the six `ErrorGroup` constants — the thing that actually varies — are explicitly excluded, and S6 notes their "own comment invites growth". The tier is assigned to the projection, not to the set that moves. | An agent exhaustively switches on exit codes per the `frozen` promise. A new `ErrorGroup` is added (the comment invites it), maps to an existing code, and the agent's diagnosis is wrong while its switch stays green. Nothing goes red. | §1, AT-3 |
| C-4 | `0029:§capability-dependencies` | The row "Release-notes discipline naming promotions" is marked `Introduced`, and its own Spec Impact concedes: "through `0.x` it is the ONLY clause here that binds, so it carries the period's whole user-visible promise. No file, template or generator backs it; a promotion landing undisclosed is caught by a reviewer or not at all." The RDR's entire binding content for the next N releases is an unenforced convention. | Nothing, for months. Then a promotion ships undisclosed, a consumer's CI goes red on an unedited model, and the diagnosis is "diff the severity constants between two builds" — which the RDR itself names as the only route in `0029:F2`. | §1, §2, premortem, AT-4 |
| C-5 | `0029:C2` | The tier is defined as "a DECLARATION OF INTENT… MUST NOT be read by a consumer as a guarantee already in force" during `0.x`, while `0029:C4` publishes the tier table into `docs/cli-output-contract.md` — the document `llms.txt` names authoritative. A published table saying `frozen` that means "not frozen yet" is a lie with a footnote. Agents do not read footnotes; the RDR's own Problem Statement says so ("An agent hard-fails, or worse, mis-parses without failing"). | An agent reads `frozen: escape_class, verdict, operator` and exhaustively switches. It is correct today and wrong at any `0.x` release, with no signal, because the schema major stays `0` and the minor moves for both additive and breaking change alike (`0029:C1`: "the minor component still increments on every change"). The one field that could carry the warning is deliberately uninformative. | §1, §2, §3, premortem, AT-5 |
| C-6 | `0029:A5` | A `Pending` assumption whose failure branch — "the CLIError `code` vocabulary has no single declaration site (61 raise sites… 10 are inline literals), so a `clierr`-side registry either imports those packages — inverting the leaf layering `clierr` exists to preserve — or requires the constants move" — is the LARGEST vocabulary the RDR tiers, and the fallback is "record the absent seam explicitly rather than oblige one". That fallback deletes the assertion for the only surface a consumer branches on constantly. | The `append-only` promise on refusal codes is prose only. A code is renamed in a refactor, no test fires, and every agent matching on that code takes the wrong branch on upgrade. | §1, §3, AT-6 |
| C-7 | `0029:S2` | S2 states the refusal baseline key set is `code,detail,message,param,schema_version`, but `Param`, `Detail`, and `Hint` are all `omitempty` on `CLIError` (`internal/cli/clierr/clierr.go:50-74`). S2 names a baseline that only one specific provoked failure emits, then calls it "the captured baseline plus the one field". The additive criterion is stated but the baseline it is additive OVER is unstable. | A test written to S2's literal wording passes on one failure path and fails on the next refusal that omits `param`. The MVV appears to pass; a different verb's refusal breaks CI for a reason unrelated to this RDR. | §1, AT-7 |
| C-8 | `0029:§alternatives-considered` (`0029:ALT2`) | All three alternatives are rejected against facts, but the option the `0.x` framing actually implies — *declare the tiers internally, publish nothing until 1.0.0* — is never considered. ALT2 is rejected as "declare all output unstable", which is not the same proposal. The RDR pays the full publication cost (Activation Step 1's four checkable things, `llms.txt` edit, register section) to publish a promise it simultaneously says is not in force. | None directly — the user pays it as the C-5 symptom. The defect is that the cheapest correct option was never on the table, so nobody weighed publish-late against publish-with-a-disclaimer. | §2, AT-8 |
| C-9 | `0029:§minimum-viable-validation` (`0029:MVV`) | MVV step 4's "version movement" row admits its own controls are unverifiable: "Both controls are read by a human running the MVV: nothing in-repo defines a release for a test to range over, and no mechanical link binds 'a vocabulary gained a member' to 'the constant changed' (A5)." The MVV's load-bearing observation — the one the whole promise rests on — has a human as its oracle. | The MVV is recorded as passed. `schema_version` is never incremented on a subsequent additive change, because nothing forces it. A consumer detecting movement by the minor component sees a frozen `"0.1"` across four releases that each changed the wire. | §1, §2, premortem, AT-9 |
| C-10 | `0029:§activation-step-3-promote-the-schema-to-1-0-at-1-0-0` | "No code change is expected at that point — that is the test of whether this RDR worked." The single deferred step is the one that converts every tier from intent to guarantee, and it is scoped as a constant edit with no gate, no re-audit of C4's census against the tree as it then stands, and no check that the six owed seams were ever built. | At 1.0.0 the project flips a string and inherits fifteen guarantees it has not re-read in however many months. The first violation is discovered by a consumer. | §2, premortem, AT-10 |
| C-11 | `0029:§technical-design` | The `growing` tier is defined by "a new member MAY fire on input that previously produced no such finding", and assigned to exactly one vocabulary — the graph-lint ADVISORY codes — whose defining property (`0006:C17`, restated in `0029:C3`) is that it CANNOT change the success disposition. A tier whose distinguishing power is verdict-instability, applied solely to the one set guaranteed verdict-stable. `append-only` and `growing` are operationally identical for every surface in C4. | A consumer reads two tier names, cannot find any behavioural difference between them, and treats them as one. The tier taxonomy's third rung buys nothing and costs a reader. | §1, §2, AT-11 |
| C-12 | `0029:§proportionality` | Profile is `large`; the record authors C1 (envelope field + wire schema versioning), C2 (a stability-tier vocabulary), C3 (a severity-promotion release policy) and C4 (fifteen tier assignments + a six-row seam obligation). The gate's own split test is "the sole author of at most one independent load-bearing contract". C3 in particular is a release-process policy with a different audience, a different enforcement mechanism (review, not test), and a different binding date (`0.x` vs 1.0.0) from C1. | The implementation lands as one commit touching six packages, the docs register, `llms.txt`, a peer RDR's contract text, and three peer test files. Any one part failing review blocks all of it. | §2, AT-12 |
| C-13 | `0029:A1` | A1 makes amending a peer record (`0006:C17`, status **Implemented**) a *prerequisite* of this RDR's implementation, and rewrites three of that peer's tests (`TestReq73`, `TestReq74`, `TestReq80` — S7/S8). An Implemented record's contract text and its assertions are being edited by a Draft peer on the strength of an in-record "Ruling". | If 7.1 cluster reconcile rules the other way, or 0006's owner disputes the reading, the implementation is blocked at its first prerequisite with the rest of the plan already written against it. | §1, §3, AT-13 |
| C-14 | `0029:§activation-step-1-publish-the-promise-where-agents-already-read` | The step concedes the standing tension — "`llms.txt` opens by telling agents to prefer the binary over any file here… it does sit in the part of the read path this project deprioritizes" — and accepts it, deferring the binary-exposed tier table as "the successor if the prose register proves unread". There is no defined signal for "proves unread". | The register is never read by the audience it was written for. Nobody finds out, because unread prose emits nothing. | §1, §2, premortem, AT-14 |

---

## 1. The three most likely ways implementation goes wrong

### (a) The `type` vocabulary contradiction ships into the published register, and the refusal branch of every consumer is dead code

**Root cause.** `0029:C4`'s seam census obliges `respond::Types() []string` over `{ok, failed}` and tiers that set `frozen`. `0029:C1`, the Illustrative Code, `0029:S2`, `0029:MVV` step 5 and the `trace` table all say the opposite: the refusal record is the bare `CLIError`, carries no `type` key at all, and `"failed"` is never emitted on the wire.

**The passage that enabled it.** From `0029:C4`:

> | envelope `type` | `frozen` | none — `respond` sets `"ok"` as a literal and a refusal carries no `type` at all | `respond::Types() []string` over `{ok, failed}` |

Read that row's own "Seam today" cell against its own "Owed" cell. The cell states, correctly, that a refusal carries no `type` at all — and then obliges a seam that enumerates `failed` as a member of the frozen vocabulary. The RDR is aware of the fact and codifies its negation in the same table row.

This is not a typo, because the source of the confusion is live in the tree and the RDR half-caught it. `internal/cli/respond/respond.go`'s package doc says:

> {"type": "ok",     ...}   # success
> {"type": "failed", ...}   # graceful failure
> …The terminal "ok"/"failed" type names are reserved.

`0029:§step-2-retire-the-ambiguous-closed-wording` notices exactly this and orders it corrected — "Also correct the two doc comments that describe the failure envelope as `{"type":"failed",…}`". So the RDR knows the package doc is wrong. It then builds the frozen seam over the wrong package doc's set anyway. Step 2 fixes the comment; C4's census re-enshrines the comment's claim as a normative contract with a by-value test behind it.

**Symptom the user sees.** The register publishes `envelope type: frozen — {ok, failed}`. An agent author reads it — that is the whole point of Activation Step 1 — and writes the switch the `frozen` tier explicitly licenses ("a consumer MAY treat an unrecognized member as a defect", `0029:C2`). Their `case "failed"` arm never executes. Their `default` arm receives every refusal this CLI has ever emitted. If they wrote `default: // unreachable`, a refusal is reported as a successful no-op — which is, verbatim, the failure mode `readPlan`'s 30-line comment block exists to prevent inside this very repo. The RDR ships the defect to external consumers on the release that promises stability.

### (b) `schema_version` is declared the consumer's discriminator, and it is the one field that discriminates nothing

**Root cause.** Having correctly established that the two terminal records are structurally asymmetric, `0029:C1` reaches for a resolution and grabs the wrong field.

**The passage that enabled it.** From `0029:§illustrative-code`:

> `schema_version` is therefore what a consumer branches on to know which record shape it holds, since the two records share no discriminator: `ok` carries `type`, a refusal carries `code`. A consumer distinguishes them by `code`'s presence, exactly as `internal/cli/flow_input.go::planEnvelope` already does.

Read the first sentence and the last sentence of that paragraph. Sentence one: branch on `schema_version`. Sentence three: branch on `code`'s presence. These are different instructions, in the same paragraph, and only the second one is implementable — `schema_version` is non-`omitempty` on BOTH records by C1's own rule, so its presence is constant and its value is identical. It carries exactly zero bits of discriminating information. C1 makes it the one field the records share and then calls that shared field the discriminator.

Worse, the RDR's grounding of this is wrong at the mechanism level. `0029:A3` and `0029:S3` assert `planEnvelope` "survives the added key by construction" because it is a plain `json.Unmarshal` over a four-field struct. That is true of the *decode*. It is not the property that matters. `readPlan` discriminates on:

```go
typePresent := len(bytes.TrimSpace(env.Type)) > 0
if !typePresent && env.Code != "" { /* refusal */ }
```

— absence of `type` AND presence of `code`. The RDR never reads past the struct definition to the discrimination logic, so it never notices that its own "share one field across both records" design is orthogonal to how the one in-repo consumer actually tells them apart, and that publishing `schema_version` as the discriminator instructs external consumers to do the one thing the in-repo consumer deliberately does not do.

**Symptom the user sees.** A consumer follows the published rule, branches on `schema_version`, and cannot separate success from refusal. In a `--plan` pipeline the refusal decodes to zero writes and the run reports success for a transition the model declined. The RDR's own Problem Statement names this class: "its parse silently takes the wrong branch".

### (c) The `0.x` framing evacuates every binding claim, leaving an unenforced release-note convention as the entire deliverable

**Root cause.** The RDR discovers mid-draft (`0029:A2`) that no tag has ever been cut, and reframes the whole record around `0.x`. The reframe is honest but it is load-bearing in the wrong direction: it converts every contract from a promise into an intention, without reducing the implementation cost by a single line.

**The passage that enabled it.** `0029:§capability-dependencies`, the "Release-notes discipline" row, states the position plainly:

> C3's disclosure obligation is process, asserted by review not by a test — and through `0.x` it is the ONLY clause here that binds, so it carries the period's whole user-visible promise. No file, template or generator backs it; a promotion landing undisclosed is caught by a reviewer or not at all.

And `0029:C2`:

> While the binary's version is `0.x`, a tier is a DECLARATION OF INTENT… and MUST NOT be read by a consumer as a guarantee already in force.

Compose them. C1's compatibility rules do not bind (major `0`). C2's tiers do not bind (intent only). C4's assignments do not bind (they are C2's tiers). The sole clause the RDR itself says binds is C3's disclosure obligation, and the RDR itself says nothing backs it. The deliverable of a `large`-profile RDR with fifteen tier assignments, six new production accessors, three rewritten peer tests, an amended Implemented peer record, a re-captured golden and a docs register is: *please remember to write a release note*.

**Symptom the user sees.** For the length of `0.x`, nothing. The user reads a register that says `frozen` and gets no guarantee. Then one release ships an undisclosed promotion, their CI goes red on a model nobody edited, and `0029:F2` tells them their diagnosis path is "diff the severity constants between two builds" — a procedure that requires the two builds' source, which an agent consuming a released binary does not have.

---

## 2. The one section rewritten within 6 weeks of shipping

**`0029:C4` — the tier assignments and the seam census.**

It will be rewritten, and the rewrite will be forced by the implementation, not chosen.

Three independent forcings, any one sufficient:

1. **The `{ok, failed}` row (C-1) cannot be implemented as written.** The moment someone writes `respond::Types()` and the by-value test S6 demands, they hold a list containing a string the wire never emits. They either delete `failed` — changing a `frozen` vocabulary's declared membership, which C1 classifies as a MAJOR bump — or they keep it and ship a lie. Either way C4's row changes.

2. **The census claims completeness and has already failed that claim twice, inside this record.** `0029:A4`'s own evidence: "C4 originally named seven of the fourteen", then the grounding lens found two more inside `findings[]`, and the count moved from fourteen to sixteen surfaces mid-draft. A4 concedes it: "That the gap recurred once inside the fix is why C4 now states where a vocabulary must be looked for". A census that has been wrong twice in one drafting cycle, over a codebase nobody changed in between, will be wrong a third time the first week someone reads it with a compiler.

3. **A5 is `Pending` and its failure branch rewrites two rows.** If the `clierr` registry inverts the leaf layering — which A5 says it does, on 61 raise sites across three packages — then the CLIError `code` row and the `flow next` `reason` row lose their seam clause and the census's "Six owed, nine carried" arithmetic is wrong. A5's own fallback is "record the absent seam explicitly rather than oblige one", which is a C4 edit.

The section is a fifteen-row table of factual claims about a moving tree, authored before the code that would validate it exists, with one of its supporting assumptions still Pending and its own history showing two corrections in one draft. It will not survive contact with `go build`.

---

## 3. The one assumption that will not survive first contact with a real user

**`0029:C2`: that a consumer will read a published tier table as "intent, not guarantee" because the document says so.**

> While the binary's version is `0.x`, a tier is a DECLARATION OF INTENT: it states what the surface is expected to promise at 1.0.0, and MUST NOT be read by a consumer as a guarantee already in force.

This assumption is refuted by the RDR's own Problem Statement, forty lines earlier:

> A human reader degrades gracefully when a field is renamed or a new diagnostic appears — they read the message and adapt. An agent hard-fails, or worse, mis-parses without failing.

The audience is agents. The RDR's founding observation is that this audience does not degrade gracefully, does not adapt, and does not read around the primary claim. The mitigation chosen for the central ambiguity of the entire record — "is `frozen` frozen?" — is a prose qualifier in the same document, aimed at a reader the RDR has already characterized as unable to handle prose qualifiers.

And the mitigation is not even uniformly applied. `0029:C1` gives the schema version a mechanical signal for exactly this — "While its major is `0` the schema is explicitly unstable" — but then neutralizes it: the minor increments on every change, additive and breaking alike, so the one machine-readable field that could distinguish "this release broke you" from "this release added a key" is deliberately made unable to. A consumer during `0.x` sees `"0.1"` → `"0.2"` and cannot tell which happened. C1 had the mechanism in hand and declined to use it.

A real agent author does one of two things with a published table reading `frozen: allow, deny, indeterminate`. Either they exhaustively switch — the tier's stated license — and break on the release that adds a verdict. Or they read the disclaimer, conclude the whole table is non-binding, and ignore all of it, in which case the register's entire publication cost bought nothing. There is no third reader. The RDR needs the first behaviour for the table to have value and the second behaviour for the table to be safe, and assumes it can have both from one document.

---

## 4. Premortem — written from after the failure

It is eleven months on. `intrastate` is at v0.7.2. The version-promise work landed in 0.2.0 and is regarded internally as done.

**Journey one — the agent that switched on `type`.** A skill author binds `intrastate lint --as=json` into a CI harness. They do what the register told them to: they read the tier table in `docs/cli-output-contract.md`, find `envelope type — frozen — {ok, failed}`, and write the exhaustive switch the `frozen` tier licenses. `case "ok"` handles findings. `case "failed"` reads `code` and `message` and posts a comment. `default` panics, because `frozen` means an unrecognized member is a defect and they are being conscientious.

Their first real refusal arrives. `respond.Fail` calls `clierr.EmitJSON`, which calls `WriteJSONLine` on the bare `*CLIError`. There is no `type` key. Their switch takes `default`. The harness panics on the first genuinely broken model it ever saw. They file an issue. The maintainer reads it, reads `0029:C4`, reads `respond::Types()`, reads the by-value test that pins `{ok, failed}` green, and takes four hours to establish that the RDR, the register, the seam and the test are unanimously wrong and the code is right. The fix is to remove `failed` from a vocabulary declared `frozen`, which under `0029:C1` is a major bump. At v0.7.2. They do it in a patch release and say nothing, because the alternative is v1.0.0 for a typo.

**Journey two — the plan pipe.** A second consumer builds on `flow resolve --as json | intrastate flow next --plan -`. They read C1's line — "`schema_version` is therefore what a consumer branches on to know which record shape it holds" — and build their own intermediate filter on that rule, checking `schema_version` to route documents. Every document has it. Their router cannot distinguish records. `readPlan`'s own discriminator (`!typePresent && env.Code != ""`) is the correct rule and is documented only in a Go comment, never hoisted into the register — though `0029:§activation-step-1-publish-the-promise-where-agents-already-read` item 3 identified it as "the single most consumer-relevant fact in the record". Item 3 was written into Activation Step 1's "done means four checkable things" list, and the implementer checked three of them because the fourth required reading `flow_input.go`.

**Journey three — the promotion nobody disclosed.** In 0.5.0 someone promotes `graph-vacuous-atom` from `info` to `blocking`. It is a good change; the condition really is an error. `0029:C3` obliges a release note naming the code. `0029:§capability-dependencies` recorded, at draft time, that nothing backs this: "No file, template or generator backs it; a promotion landing undisclosed is caught by a reviewer or not at all." It was not caught. There is no CI check, because `severityFor` derives severity from `IsBlocking` and the change was a one-line edit to `blockingCodes` — a diff no reviewer flagged as a policy event, because the RDR's stated mitigation ("C4's tier assignment is the reviewable artifact") requires the reviewer to hold C4 in their head while reading a slice literal.

Four consumers' pipelines go red on models nobody edited. They follow `0029:F5`: compare `severity` across two builds. Two of them are agents consuming a released binary with no source checkout; they cannot. One opens an issue titled "0.5.0 broke my CI". The response cites `0.x` and SemVer.

**Journey four — the 1.0.0 flip.** v1.0.0 is scheduled. `0029:§activation-step-3-promote-the-schema-to-1-0-at-1-0-0` says: "No code change is expected at that point — that is the test of whether this RDR worked." Someone edits `clierr.SchemaVersion` from `"0.7"` to `"1.0"` and tags. Fifteen intent-declarations become guarantees in one commit. Nobody re-audited C4's census against a tree that had gained three verbs and a new findings field in the interim. Nobody checked whether the six owed seams were ever built — two were not, because `0029:A5` went the way A5 said it might and the CLIError registry inverted the leaf layering, so the largest vocabulary in the record shipped with its `append-only` tier as prose. The MVV's version-movement row was "read by a human running the MVV"; the human ran it once, in 0.2.0, and `schema_version` sat at `"0.1"` for four releases that each changed the wire, because `0029:MVV` conceded "no mechanical link binds 'a vocabulary gained a member' to 'the constant changed'".

The post-mortem's finding: the RDR was a documentation project wearing a contract's clothes. It produced fifteen assertions about the wire, six of which were never mechanized, one of which was factually false about the wire it described, and one process obligation that nothing enforced — and it spent its most valuable asset, the `0.x` free window, publishing a table it simultaneously told readers not to rely on.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 — the `type` vocabulary matches the wire (catches C-1)**
```gherkin
Given the RDR declares a frozen vocabulary V with declared members M
When each member of M is traced to an emit site in the tree
Then every member of M has at least one site that writes it to the wire
And for envelope `type`, M is exactly {"ok"} because respond.Fail emits
    a bare CLIError with no `type` key
And no RDR passage states or implies "failed" is an emitted `type` value
```
Run at review against `0029:C4`'s census: the `{ok, failed}` row fails at step 3. This is a five-minute grep for `"failed"` in `respond` and `clierr` emit paths; the RDR instead inherited the claim from a package doc comment it had separately identified as wrong in `0029:§step-2-retire-the-ambiguous-closed-wording`.

**AT-2 — the stated discriminator actually discriminates (catches C-2)**
```gherkin
Given the RDR names field F as what a consumer branches on to
      distinguish record shapes
When the success record and the refusal record are both constructed
Then F's value differs between them, OR F is absent from one of them
```
`schema_version` is non-`omitempty` on both with an identical value. Fails immediately. A second arm:
```gherkin
Given an in-repo consumer is cited as evidence that a change is safe
When the citation is checked
Then it covers the consumer's DISCRIMINATION logic, not only its
     struct decode
```
`0029:A3` and `0029:S3` cite `planEnvelope`'s struct and stop. `readPlan`'s `!typePresent && env.Code != ""` is the load-bearing logic and is never read.

**AT-3 — a frozen tier is assigned to the set that varies (catches C-3)**
```gherkin
Given a vocabulary is assigned `frozen` with a by-value seam
When the producing function is read
Then every declared member is reachable from an emitted envelope
And no member is a fallthrough/default for a non-envelope path
And the upstream set that determines the member (e.g. ErrorGroup) is
    itself tiered, or the RDR states why the projection is the contract
```
`ExitCodes() {0,1,2,3,130}`: `0` is success (no refusal), `1` is the `errors.As` fallthrough for a non-`CLIError`. Two of five fail. The six `ErrorGroup` constants — the set that actually moves, whose comment `0029:S6` admits "invites growth" — are excluded with no stated reason.

**AT-4 — every binding obligation has a mechanism (catches C-4)**
```gherkin
Given the RDR identifies a clause as the only one binding during 0.x
When that clause's enforcement is named
Then it is a test, a CI check, a generator, or a template
And not "a reviewer"
```
`0029:§capability-dependencies` self-reports the failure verbatim. Nothing in review escalated a self-reported "caught by a reviewer or not at all" on the record's sole binding clause.

**AT-5 — a published guarantee word means what it says (catches C-5)**
```gherkin
Given a tier name is published in a document llms.txt names authoritative
When a consumer reads the tier name without the surrounding prose
Then the behaviour the name licenses is safe to perform today
```
`frozen` licenses an exhaustive switch (`0029:C2`) and is not safe during `0.x`. Fails. Remedies that would pass: suffix the published names during `0.x` (`frozen-at-1.0`), or withhold the table until 1.0.0, or make the schema minor distinguish additive from breaking so the disclaimer has a wire signal behind it.

**AT-6 — no tier rests on a Pending seam (catches C-6)**
```gherkin
Given C4 obliges an enumeration seam for vocabulary V
When the assumption backing that seam is Pending
Then V's tier is not published as a consumer-facing promise until
     the assumption resolves
```
`0029:A5` is Pending over the CLIError `code` vocabulary — the largest set, and the one agents branch on most — and its own fallback is to drop the seam obligation.

**AT-7 — a baseline key set is stable (catches C-7)**
```gherkin
Given a test scenario names an exact key set as its captured baseline
When each named key's struct tag is read
Then no key in the set carries `omitempty`
```
`0029:S2` names `code,detail,message,param`; `Param` and `Detail` are both `omitempty` (`clierr.go:53-58`). Fails on two keys.

**AT-8 — the cheapest option is on the alternatives list (catches C-8)**
```gherkin
Given the RDR concludes its contracts do not bind during the current
      version series
When the Alternatives are read
Then one alternative is "declare internally, publish at 1.0.0"
```
Absent. `0029:ALT2` ("declare all output unstable") is a different proposal and is rejected on grounds — `llms.txt` invites coupling — that do not apply to publish-late.

**AT-9 — the MVV's load-bearing step has a mechanical oracle (catches C-9)**
```gherkin
Given an MVV step is named the promise's load-bearing observation
When its oracle row is read
Then the control is executable, not "read by a human"
```
`0029:MVV`'s version-movement row states its controls "are read by a human running the MVV" and that "no mechanical link binds 'a vocabulary gained a member' to 'the constant changed'". Fails on its own text. The mechanizable form exists and is cheap: a test that computes a hash over every tiered vocabulary's seam output and pins it beside `SchemaVersion`, so a member added without bumping the constant goes red.

**AT-10 — the deferred activation step has a gate (catches C-10)**
```gherkin
Given an activation step converts declarations into guarantees
When the step is read
Then it names a re-audit of the tier census against the tree at that
     time, a check that every owed seam exists, and a named owner
```
`0029:§activation-step-3-promote-the-schema-to-1-0-at-1-0-0` is four lines, names a constant edit, and asserts "No code change is expected at that point" — which is a prediction, not a gate.

**AT-11 — each tier is distinguishable on a real surface (catches C-11)**
```gherkin
Given the RDR defines N tiers
When the tier assignments in C4 are read
Then for each adjacent pair of tiers there exists a surface where the
     weaker tier permits an observable a consumer can detect and the
     stronger forbids it
```
`append-only` → `growing` differs only by "a new member MAY fire on input that previously produced no such finding". `growing` is assigned solely to the advisory codes, which `0029:C3` and `0006:C17` guarantee cannot change the disposition. There is no surface on which the distinction is observable. The tier count should be two, or `growing` needs a surface that justifies it.

**AT-12 — one seam per record (catches C-12)**
```gherkin
Given the Finalization Gate's split test is "sole author of at most one
      independent load-bearing contract"
When C1..C4 are read for independence
Then at most one survives as independently load-bearing
```
C1 (wire schema versioning) and C3 (a release-process severity policy binding on a different date, enforced by a different mechanism, read by a different audience) are independent. The gate's own test says split.

**AT-13 — no Draft edits an Implemented peer's assertions as a prerequisite (catches C-13)**
```gherkin
Given a Draft RDR lists amending a peer's normative contract as a
      prerequisite of its own implementation
When the peer's status is read
Then the peer is not Implemented, OR the amendment is a separate
     record with its own gate
```
`0006` is Implemented. `0029:A1` amends `0006:C17` and rewrites three of 0006's tests (`TestReq73`, `TestReq74`, `TestReq80`) on the authority of an in-record Ruling.

**AT-14 — a deferred successor has a trigger (catches C-14)**
```gherkin
Given an RDR defers a mechanism pending evidence that the chosen one
      failed
When the deferral is read
Then it names the observable that constitutes failure and who watches it
```
`0029:§activation-step-1-publish-the-promise-where-agents-already-read` defers the binary-exposed tier table until "the prose register proves unread". Unread prose emits no signal; the trigger is unobservable by construction.
