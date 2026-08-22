Model: claude-fable-5

# Critique — RDR 0007 Guard predicate totality (iter-2, Model B)

Independent hostile pass. Prior critique evidence not read. Grounded
against: the RDR full text; JDR 0001; `internal/resolve/resolve.go` and
its `_test.go` files on `via-claude`; the A3 spike record
(`evidence/spikes/a3-reshape-probe.md`); `internal/cli/clierr/clierr.go`;
RDR 0005's stable-code table.

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| B-1 | Decision Rationale "✓ by structure — the evaluator is never asked about an absent key"; A16/A22/A24 producer halves; Prerequisites "Two duties are added to RDR 0002's Refinement Context Direction list, which currently carries neither" | "Structure not discipline" is a relocation, not a closure: the guarantee now rests on a normalizer (RDR 0002, Draft) that has not received the token-emission or canonicalization duties as clauses. Normalizer drift on token spelling, literal form, or key canonicalization makes every guard permanently unevaluable | Non-escapable `guard_unevaluable` storm, reason `absent`, for keys the author can see spelled in their own table; payload misdiagnoses a producer defect as missing artifact state | §1, premortem, AT-1 |
| B-2 | A3 "Pending — the SHAPE the spike proved is not the shape this RDR specifies"; A24 "the A3 spike's existence atom is `{Key: \"guard.missing\", Op: resolve.OpExists}` with NO Literal set"; Testing Strategy "Phase 1's exit condition is the re-spike's count" | The only executable evidence implements a contradicting design (retained `Refusal.Guard`, `UndecidedAtoms` beside it, literal-free `exists` decided from presence). A3/A23/A24 are all Pending on a re-spike that has not run; the RDR's exit condition is a number that does not exist. Implementers follow the concrete artifact, not the prose | `exists` atoms with zero-value literals silently prune instead of refusing; Fixup-1d re-encoded instead of re-decided; the frozen suite's meaning shifts with no failing test | §1, §2, premortem, AT-2 |
| B-3 | A19 "**This is a REVERSAL of a closed joint decision**"; Prerequisites "**Does not block lock**"; payload clause "RDR 0005's renderer type[s] against it" | The flagship user outcome — "told plainly … was missing" — has no settled transport at lock. The fallback (flatten into `Detail`) is destroyed by A19's own text ("defeats the total-order determinism whose only purpose is machine consumption"), yet is what §D4 (Closed) currently rules | Payload exists in a Go struct nobody can see; JSON consumers regex-parse a prose `Detail` blob; the RDR's most-argued clause (the six-field total sort) has zero consumers | §1, §2, premortem, AT-3 |
| B-4 | disposition mini-check "Two exit-group notes, both upstream of this RDR and neither fixable here"; `flow-guard-unevaluable` scoped to "supplied facts", `GroupUserEnv` | Exit-code remedy signal is wrong for the flagship case: missing *artifact* state maps to exit 2 "bad input, bad flags"; `owned_state_unavailable` has no code row at all. Deferred to RDR 0005's re-lock, which nothing schedules | CI wrapper branches on exit 2, reports "bad input", operator retries with different flags — exactly the wrong remedy; P4 fails at the only surface users touch | §1, §2, premortem, AT-3 |
| B-5 | §SURVIVOR MEMBERSHIP "This veto is the cost the RDR accepts: one unreadable row refuses a table whose other rows decide cleanly"; Failure Modes "`unless` over a rarely-set tag … a footgun in practice" | Resolution-scope veto makes refusal cost scale with table size while its benefit stays constant. First production table that grows an optional-tag row bricks every artifact lacking that tag, even when a decided-TRUE sibling exists; the sanctioned `unless` idiom lives in authoring guidance not yet written, in a Draft RDR | Adding one row refuses previously-working resolutions across the fleet; author is told about a row they didn't touch; pressure to narrow the veto arrives within weeks and the RDR pre-brands any narrowing as "absence-as-false" | §2, premortem, AT-5 |
| B-6 | A5 "the load-time half rides on A12"; A12 "Deferred … the question is open"; Phase 4 "authoring grammar: may one tag-table carry two operator keys" | The only sanctioned route to "absent OR equals v" is a pattern that (a) may not be writable in 0003's grammar (open question) and (b) may be rejected by the load-time overlap check with "no downgrade valve" (A12 deferred). The domain rule forces authors toward an escape hatch that may not exist | Author hits the veto, cannot express the sanctioned pattern, and falls back to sentinel-stamping upstream — the named anti-pattern whose "only sanctioned alternative *is* this pattern" (A5 If-wrong) | §3, premortem, AT-6 |
| B-7 | A13 "accepted exposure"; Failure Modes "Refusal defeated by caller-supplied state"; Testing Strategy row 16 pins `--tag X=v` ⇒ row selected and a plan | The design's strictness manufactures the pressure that defeats it: a non-escapable, table-wide veto plus a documented one-flag bypass trains operators to bake `--tag` into wrappers. Scenario 16 pins the bypass as contract, and nothing in `Plan` marks a guard decided from caller input | Missing artifact state masked at scale, permanently and invisibly — worse than the pre-RDR probe, because it is now sanctioned, test-pinned, and leaves no trace in the plan | §3, premortem, AT-7 |
| B-8 | Background "An overlap key absent entirely yields `no_match` … (the match pattern is closed-world by design)"; Phase 4 "predicate-PLACEMENT guidance" | The masking path is relocated, not closed. The identical predicate expressed as a `Match` tag instead of a guard atom still turns absence into escapable `no_match`. The only control is authoring guidance handed off in Phase 4 to a Draft RDR — discipline again, at the author layer | The Background probe reproduces byte-for-byte with the RDR marked Implemented: a table looks like it worked, a plan routes around missing state via `Match` + escape; kata `xg7p` reopens | §1, §3, premortem, AT-8 |
| B-9 | payload clause reason set "`absent` … or `uncomparable`"; Failure Modes "the author is not told a tag is missing, but is not told 'skew' either"; A22 If-wrong "misdiagnoses a producer defect as missing state" | The closed two-reason set cannot carry the diagnosis the Problem Statement promises: unparseable caller value, evaluator/grammar version skew, and foreign literal all collapse into `uncomparable`; canonicalization drift reads as `absent`. Three different fixers, one label | Operator stares at `uncomparable` on a key they supplied correctly, or `absent` for a key that exists; the refusal names *what* blocked but systematically misattributes *who can fix it* | premortem, §4, AT-9 |
| B-10 | §GATE, THEN COUNT "owned-before-unevaluable precedence is pinned to shipped behavior, not to D8's original rationale … which this RDR's narrowing invalidates" | The RDR admits the precedence's rationale is dead, keeps the ordering anyway, and accepts serial diagnosis: a row failing both ways surfaces only the write-dependency problem | Two-round-trip debugging: fix owned state, re-run, then discover the guard was also undecidable — on a refusal class whose whole pitch is "actionable the first time" | premortem, AT-10 |
| B-11 | Prerequisites (all seven boxes); "This RDR Final **before** RDR 0003's implementation begins" | Lock is a promissory note against five external documents: a JDR reopen (reversing a Closed decision), two clause additions to Draft 0002, a placement ruling from Draft 0003, a Final-0009 re-lock, and a Final-0005 re-lock — none scheduled, none owned by this RDR, several sequenced circularly against 0003 | Implementation starts with the ground still moving; each unlanded duty resurfaces as a mid-build seam dispute (token spelling, canonicalization, transport) that this RDR already declared out of its own scope | §2, premortem, AT-11 |
| B-12 | payload clause "What the payload does NOT promise: a row pruned as GuardFalse under `F ∧ U = F` is not a survivor, so its absent keys are not reported" | The `F ∧ U = F` prune is sound for masking but leaves the most common debugging question — "why didn't my row fire?" — with no surface anywhere: the pruned row's absent key appears in no payload, no refusal, no plan | Author of a row pruned by a decided F beside an unevaluable atom gets silence; they add debug rows or sentinels to make the kernel talk, re-importing the anti-pattern | premortem, AT-12 |
| B-13 | Phase 3 "**Blocked on one RDR 0003 declaration** … cannot be written from any current document"; Testing Strategy row 8 | The `contains` leg of the exported value-seam contract — guarding the exact operator this RDR pins ("absent set ≠ empty set") — is unwritable at lock. The one operator with a named empty-set masking hazard ships with no executable contract for its present-key half | RDR 0003's evaluator later treats a malformed set value as empty (decided FALSE) and nothing exported catches it; the masking path reopens one level down, in the component this RDR chose not to trust | §1, AT-13 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The normalizer never receives its half of the contract, and the "structural" guarantee decays back into discipline — at a seam with worse failure geometry

The RDR's entire case for Approach B over Alternative 1 is one
sentence in the scored matrix: "✓ by structure — the evaluator is
never asked about an absent key," against A's "✗ — a conformance
harness that nothing compels RDR 0003's build to run." That is the
deciding row. And it is a shell game. Approach B does not eliminate
the trusted-but-unbuilt component; it swaps the evaluator for the
normalizer. For the kernel's presence step to mean anything, RDR
0002's normalizer must: emit the exact `OpExists` token and the exact
`LiteralTrue`/`LiteralFalse` forms (A16, A24), and canonicalize every
authored key spelling so `Tag.Key` matches the guard atom's key
byte-for-byte (A22 — "Key identity is exact string equality … the
kernel performs none").

The RDR then documents, in its own Prerequisites, that **neither duty
exists**: "Two duties are added to RDR 0002's Refinement Context
Direction list, which currently carries neither." A16's producer half
is "acknowledged only by a joint-check line, not a clause." A22 is
flagged "an UNSTATED INFERENCE against RDR 0002's current text" —
0002 has "no canonicalization clause and no tag-name grammar or
charset rule anywhere." RDR 0002 is Draft. Nothing compels its
re-lock to accept either duty, which is *the exact defect class* the
RDR rejects Alternative 1 for — "a conformance harness nothing
compels 0003's build to run" — restated one document over.

The failure geometry is worse than Alternative 1's, not better. Under
A, evaluator drift produces wrong verdicts on the guards that drift
touches. Under B with normalizer drift, the miss is *systemic*: a
normalizer that case-folds keys, trims a namespace, or spells the
existence token `exist` makes **every** affected guard atom
permanently unevaluable, and the kernel's payload — presented as
"kernel truth" — confidently reports `absent` for keys the author can
see in their own file. The RDR itself concedes the misdiagnosis in
A22's If-wrong: "misdiagnoses a producer defect as missing state."
Fail-closed, yes. But a fail-closed storm across a whole table, with
a payload that points at the artifact instead of the producer, is an
outage with a lying diagnostic.

**Root cause**: the Decision Rationale's structural claim is scoped to
the kernel↔evaluator seam and silently assumes the kernel↔normalizer
seam, which has no normative existence yet.
**Enabling passage**: the scored matrix's drift-enforcement row, plus
the Prerequisites bullet quoted above.
**Symptom**: non-escapable `guard_unevaluable` refusals, reason
`absent`, on keys visibly present in the authored table, the first
time 0002's normalizer makes any spelling decision of its own.

### 1.2 Phase 1 gets implemented from the spike, and the spike proved a different design

There is exactly one executable artifact behind this RDR: the A3
reshape probe. It is also the only place the 154-PASS number, the
mutation result, and the six-case Kleene probe live. And it
implements a design the Normative Contracts contradict on at least
three points:

1. The spike **retained** `Refusal.Guard` (retyped `[]GuardAtom`) and
   added `UndecidedAtoms` beside it; the spec **replaces** the field
   with `Undecided []UndecidedRow` and retires `gate`'s
   `slices.MinFunc` representative selection.
2. The spike's existence rule is literal-free: "absent key + `exists`
   op → decided FALSE; present key + `exists` → TRUE" — and its probe
   pins that behavior ("an `exists` atom over an absent key is decided
   FALSE (pruning the row)"). The spec's rule is `presence ==
   literal`, under which the spike's own atom (`{Key:
   "guard.missing", Op: resolve.OpExists}`, no literal) is
   `uncomparable`, and `exists = false` over an absent key decides
   TRUE. Opposite dispositions on the same input.
3. The spike "re-encoded" `TestFixup1d…` to `reflect.DeepEqual` over
   the retained field; the spec says that test "must be RE-DECIDED,
   not re-encoded," because the field it reads ceases to exist.

The RDR knows all of this — A3 is Pending "on exactly that
difference," A23 and A24 fold into "A3's re-spike," and the Testing
Strategy sets Phase 1's exit condition to "the re-spike's count with
the field actually removed — not the 154 the superset spike
returned." But **the re-spike has not been run.** The RDR is asking
to lock against an exit criterion whose value is unknown, with three
Critical Assumptions Pending on one future spike, while the only
concrete code an implementer can crib from demonstrates the rejected
shape. Every implementer under schedule pressure reads the working
diff before the 200-line normative block. The most probable Phase 1
is the spike's semantics with the spec's names — literal-free
existence included — and nothing in the frozen suite can catch it,
because (as the spike itself proved) the frozen suite cannot even
observe combination order.

**Root cause**: evidence and specification diverged after the spike,
and the RDR chose to lock the spec on a promise to re-verify rather
than re-verifying.
**Enabling passage**: A3's Status block; A24's Evidence note;
"Phase 1's exit condition is the re-spike's count."
**Symptom**: `exists` atoms silently pruning where the spec says
refuse (or refusing where the spec says decide), discovered only when
RDR 0003's authors write `exists = false` and get `uncomparable`; or
a `Refusal` carrying both old and new payloads because the spike diff
was ported wholesale.

### 1.3 The payload lands kernel-side with nowhere to go, and the exit code lies about the remedy

The Problem Statement's user is "told plainly that the artifact state
needed to decide was missing." The mechanism is the per-row/per-atom
payload — the RDR's largest normative block, with a six-field total
sort order argued down to the tie-break. Now trace the payload to the
user. JDR 0001 §D4, a **Closed** decision, routes it "through §JD-8's
`Detail`" — a flat prose string. A19 wants to reverse that with one
`omitempty` structured field on `CLIError`, admits it is "a REVERSAL
of a closed joint decision," parks the reopening in a Prerequisite,
and stamps it "**Does not block lock**." The fallback if the
reopening is declined is the very thing A19's own evidence
demolishes: "serializing a two-level sorted array into [`Detail`]
puts a second, undocumented encoding inside a string, which defeats
the total-order determinism whose only purpose is machine
consumption."

So at lock, the deterministic-payload machinery — `undecidedAtoms`,
`compareAtoms`, the total tuple, Testing Strategy row 17 — is
specified in full while its only purpose (machine consumption at the
CLI) is unresolved, and the currently-ruling decision resolves it in
the direction that makes the machinery pointless. Meanwhile the exit
surface is wrong in both cells the RDR's own disposition table flags:
`flow-guard-unevaluable` is scoped in RDR 0005 to "supplied facts"
and grouped `GroupUserEnv` — exit 2, documented in `clierr.go` as
"bad input, bad flags" — while this RDR's flagship case is missing
*artifact* state; and `owned_state_unavailable` has no code row at
all. Both are waved off as "RDR 0005's to fix at its re-lock," which
nothing schedules — 0005 is Final and this RDR's Prerequisites do not
list it.

**Root cause**: the RDR treats the user-facing half of its own user
story as another RDR's re-lock problem, then locks anyway.
**Enabling passage**: A19's Status; the Prerequisites' "Does not
block lock"; the disposition mini-check's "Two exit-group notes, both
upstream of this RDR and neither fixable here."
**Symptom**: a CI wrapper branches on exit 2, prints "bad input,"
and the operator fiddles with flags; the beautiful sorted payload is
either invisible or a regex target inside `detail`; the first
integration consumer files the bug this RDR was written to prevent —
"the output doesn't tell me what was missing."

## 2. The one section that will be rewritten within 6 weeks of shipping

**The SURVIVOR MEMBERSHIP block — specifically the resolution-level
aggregation veto** ("if any surviving candidate row's guard is
GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`; a
decided-GuardTrue sibling MUST NOT be selected while an unevaluable
candidate exists").

The veto's cost scales with table size and table age; its benefit is
constant. Every real table grows rows, and new rows reference new,
initially-optional tags — that is what "add a `legal_hold` policy to
the release flow" looks like. Under the veto, the day that row lands,
**every artifact in the fleet that lacks the new tag and reaches the
same outcome+match refuses**, including resolutions a decided-TRUE
sibling would have handled identically yesterday. The refusal is
non-escapable by design, and the recovery is "fixing the artifact
state or accessor, never the table" — i.e., a data-backfill across
the fleet as the price of one table edit.

The RDR has already priced the individual footguns and deferred every
remedy: the `unless`-over-rarely-set-tag case is admitted to be "a
footgun in practice," with the sanctioned conjoined-existence idiom
handed to RDR 0003's authoring guidance (Draft, unwritten); the
disjunction pattern rides on A12 (Deferred, open); whether the
conjoined row is even *writable* is an open grammar question in Phase
4. So the first team burned by a fleet-wide refusal storm will find
no writable idiom and no downgrade valve, and will demand the veto be
narrowed — "an unevaluable row loses only to a decided-TRUE sibling
on the same match" or a per-row severity knob. The RDR anticipates
this and pre-brands every narrowing as "absence-as-false at
resolution scope," which guarantees the rewrite happens as a
contentious hotfix RDR rather than a clean amendment. That is the
signature of a section that gets rewritten: maximal rigidity, cost
borne by a user the document never walked through, and all mitigation
deferred to documents that don't exist yet.

(Runner-up: the payload normative block, which gets rewritten the
moment A19's reopening is decided either way — but that rewrite is at
least anticipated by the RDR's own text.)

## 3. The one assumption that will not survive first contact with a real user

**A5** — "Authors who need 'row applies when tag X is absent' express
it with an existence atom; the disjunction 'absent OR equals v' needs
two rows … and the value row MUST itself carry an `X exists = true`
atom beside `X eq v`."

A5 is marked Verified, but its verification is a *derivation of what
authors must do*, not evidence that they can. Its load-time half
"rides on A12," which is **Deferred** — RDR 0003's overlap check is
defined over a product of declared finite domains, is silent on how
an existence atom projects into it, and RDR 0002 retains "ambiguous
overlap" as a MUST-reject "with no downgrade valve." Whether the
grammar even permits two operator keys on one tag-table is an open
question the RDR itself ships to Phase 4 ("authoring grammar: may one
tag-table carry two operator keys"). So the sanctioned pattern is:
two rows the linter may reject as overlapping, one of which contains
a conjunction the grammar may not parse, taught by authoring guidance
that has not been written, in an RDR that is Draft.

A real user does not follow that path. A real user hits the veto,
discovers the two levers that work *today*, and pulls one: stamp a
sentinel value upstream so the key is always present (the named
anti-pattern, now re-incentivized by the very rule that was supposed
to kill it — A5's own If-wrong says so: "authors … fall back to
sentinel-stamping upstream — the anti-pattern whose only sanctioned
alternative *is* this pattern"), or pass `--tag` (B-7, accepted
exposure, test-pinned as producing a plan). Either way, the domain
rule's honest refusal gets engineered around within the first sprint,
and the RDR's theory of authorship — that absence-conditionals are
expressed in the guard algebra — dies on the day the algebra's
authoring surface turns out not to exist.

## 4. Premortem

*Written 2027-02, six months after RDR 0007 shipped as Implemented.*

The 0007 kernel work landed clean — `evaluateAtom`, the strong-Kleene
fold in `evaluateGuard`, `Refusal.Undecided`, all green, mutation
probe killing the combinator mutant. The failure happened everywhere
else, exactly at the seams the RDR filed as other documents' duties.

First contact was the release-review flow, the first real table after
RDR 0005 wired the `flow` verb. In week two, the policy team added a
`legal_hold` row: `unless legal_hold eq true`. Every artifact that
had never carried `legal_hold` — all of them — began refusing
`guard_unevaluable`, non-escapably, because `¬U = U` vetoes at
resolution scope even though the shipped rows still decided TRUE.
The payload named the key correctly. It didn't matter, because the
CI wrapper never saw it: `respond.Fail` mapped the refusal to
`flow-guard-unevaluable`, `GroupUserEnv`, exit 2, and the wrapper's
exit-2 branch printed "invalid input — check flags." The structured
payload was in `Detail` as flattened prose — the §D4 reopening had
been declined at the JDR, so A19's fallback rode — and the six-field
total sort order the RDR argued for a page now determinizes a string
nobody parses.

The author tried the sanctioned fix. `legal_hold exists = true`
beside `legal_hold eq true` in one `unless` block: the normalizer
rejected two operator keys on one tag-table — the grammar question
Phase 4 "handed over" was answered no. The two-row pattern for the
`all` case fared worse: 0003's overlap lint flagged the pair as
ambiguous overlap, the MUST-reject category with no downgrade valve —
A12 had stayed Deferred right through 0003's re-lock. Blocked on both
sanctioned routes, the platform team did what the RDR's own Testing
Strategy row 16 certifies as working: they put
`--tag legal_hold=false` into the shared CI wrapper. Presence is
provenance-blind, the atom decided, plans flowed. Nothing in `Plan`
records that a guard was decided from caller input rather than the
artifact, so six months of release decisions now rest on a hardcoded
flag no one remembers adding. Missing artifact state is being masked
at fleet scale — by the mitigation for the rule that was supposed to
make masking impossible.

The second incident was quieter. RDR 0002's normalizer had shipped
with NFC-plus-lowercase key canonicalization — reasonable, since the
canonicalization duty reached 0002 as a direction-list note, never a
clause, and nobody negotiated the exact rule against the kernel's
"exact string equality." Tables authored with `Legal_Hold` reached
`resolve.Row` atoms as `legal_hold` while `assemble` built the view
from accessor tags that had never been canonicalized at all. Every
guard over an owned key went unevaluable, reason `absent`, and the
payload — "kernel truth" — pointed the on-call engineer at the
artifact store for two days before anyone diffed `TagSet.Lookup`'s
key against the accessor's. A22's If-wrong, verbatim: a producer
defect misdiagnosed as missing state.

Meanwhile the team that got burned by the `legal_hold` storm learned
the real lesson: guard atoms refuse, `Match` tags don't. They moved
their optional-state predicates into `Match`, where the closed-world
pattern turns absence into `no_match` — escapable, modeled, silent.
The probe from the RDR's own Background section reproduces today on
the shipped kernel, one field over from where 0007 closed it. Kata
`xg7p` was reopened in January. The hotfix RDR narrowing the
aggregation veto is in Draft, and its first review comment quotes
0007: "Narrowing it would decide that an undecided edge is a
non-edge."

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 (catches B-1).** *Producer duties exist as clauses, not
intentions.*
Given RDR 0002's current text, when I search it for (a) a normative
clause obliging the normalizer to emit the kernel's `OpExists` /
`LiteralTrue` / `LiteralFalse` constants and (b) a normative
canonicalization clause naming case, namespace, and whitespace, then
both must exist before this RDR may claim drift "fails closed by
structure." — Fails today by the RDR's own admission ("currently
carries neither"); the correct disposition is that B's deciding
matrix row is downgraded to "by structure *iff* two unlanded 0002
clauses," which un-decides the matrix.

**AT-2 (catches B-2).** *Evidence matches specification.*
Given the newest spike record cited as A3's Evidence, when its
`Refusal` shape and its existence-atom dispositions are compared
against the Normative Contracts, then every disposition must agree,
or A3/A23/A24 remain open and lock is blocked. Concretely: run the
spec's own vector — `{Op: OpExists}` with no literal — against the
spike code; spike says decided-from-presence, spec says
`uncomparable`. One input, two verdicts, both documents in evidence.
The "re-spike's count" named as Phase 1's exit condition must exist
at lock, not be promised.

**AT-3 (catches B-3, B-4).** *The user story survives the full
stack, on paper.*
Given the MVV table (value atom over absent key), walk the refusal
from `Resolve` through `respond.Fail` to a JSON consumer and an exit
code, using only currently-ruling documents (JDR §D4 as Closed, RDR
0005 as Final). Then: (a) a documented envelope field must carry
per-row/per-atom structure — fails: §D4 routes flattened text; (b)
the exit code's documented remedy class must be "artifact state /
environment," not "bad input" — fails: `GroupUserEnv`, exit 2; (c)
`owned_state_unavailable` must have a `Code` row — fails: none. Any
one failure means the Problem Statement's promise is not met by the
locked stack, and "does not block lock" is refuted by walking the
user journey.

**AT-5 (catches B-5).** *The table-growth journey has a writable
answer.*
Given a 10-row table where 9 rows decide cleanly and 1 new row
references an optional tag absent from the artifact, when an outcome
matched by a decided-TRUE sibling resolves, then the RDR must either
(a) produce a plan, or (b) refuse AND name an author idiom, writable
in the grammar as currently specified, that restores the plan without
backfilling fleet state. The RDR answers refuse-and-defer: the idiom
depends on Deferred A12 and an open Phase 4 grammar question. (b)'s
second conjunct fails — visible at review time without writing code.

**AT-6 (catches B-6).** *The sanctioned pattern loads.*
Given the two-row absence pattern and the conjoined value row in
concrete authored form, when checked against RDR 0003's overlap
contract and RDR 0002's row grammar as written, then both must load.
This check cannot even be executed — 0003 is silent on the
projection and the grammar question is open — and "cannot be
executed" is itself the review-time finding: a normative rule whose
only user-facing escape is unverifiable must not lock.

**AT-7 (catches B-7).** *The bypass is traceable.*
Given Testing Strategy row 16's accepted behavior (`--tag X=v`
decides a previously-unevaluable guard and yields a plan), then some
consumer-visible surface must distinguish a plan whose guard was
decided from caller-supplied observed state from one decided from
artifact state — otherwise the accepted exposure is unauditable.
Nothing does: `Plan` carries `RuleID`, `SourceLocator`, `NextTags`,
`Writes`, `Revision`, `Escaped` — no provenance echo. Review-time
fix is one field; post-ship it is an archaeology project.

**AT-8 (catches B-8).** *The masking probe is closed under
predicate placement.*
Given the Background probe's table, when the same predicate is
re-expressed as a `Match` tag rather than a guard atom, then absence
must not route to an escapable class — or the RDR must carry the
residual as a named accepted exposure with an enforcing owner (a lint
category, a load rule), not as Phase 4 guidance. It does neither:
`no_match` from a closed-world `Match` remains escapable by design,
and the only control is a docs handoff. The RDR's opening claim
("decides whether missing artifact state can be masked behind an
escapable refusal class") is falsified by its own Background bullet
read against its own Phase 4.

**AT-9 (catches B-9).** *Reasons partition by fixer.*
Given the three producer-side causes the RDR itself enumerates —
unparseable present value (A18), evaluator/grammar token skew,
foreign literal on `OpExists` — then the payload's reason set must
let a reader distinguish "fix your input" from "fix the producer
tooling" from "fix the artifact." All three collapse into
`uncomparable`, and canonicalization drift collapses into `absent`.
A closed reason set is fine; a closed set that cannot carry P4's
"who can fix it" is not.

**AT-10 (catches B-10).** *One refusal, one fix-round.*
Given a survivor row that both lacks an owned `RequiresOwned` key and
carries an unevaluable atom, when the user fixes exactly what the
first refusal names and re-runs, then the resolution must not refuse
again for a defect present in the original input. It does
(`owned_state_unavailable` first, `guard_unevaluable` second), and
the RDR keeps the ordering after admitting its rationale is
invalidated. The honest disposition is to report both in one refusal
— the payload machinery being added is exactly the vehicle — and the
review question is why it doesn't.

**AT-11 (catches B-11).** *Prerequisites are satisfiable and
acyclic at lock.*
Given the Prerequisites list, then every box must be closable by an
action this RDR's owner can take, or carry a named owner and a
scheduled vehicle. Count the external edits: JDR §D4 reopen; two
0002 clauses; a 0003 placement ruling; a 0009 re-lock; and (via the
disposition notes) a 0005 re-lock that is not even listed. Then
check ordering: "This RDR Final **before** RDR 0003's implementation
begins" while three prerequisites require 0003/0002 document work —
the dependency graph must be drawn and shown acyclic. It is not
drawn.

**AT-12 (catches B-12).** *The pruned-row debugging journey ends
somewhere.*
Given a row pruned via `F ∧ U = F` whose unevaluable atom referenced
an absent key, when the author asks "why did this row not fire," then
some surface (payload, verbose refusal, trace) must answer. Nothing
does — the RDR states the gap ("its absent keys are not reported")
without giving the journey an endpoint. The guarantee "never a plan
from absence" is intact; the *diagnosability* promise of the Problem
Statement is not, for the second-most-common question a table author
asks.

**AT-13 (catches B-13).** *Every operator this RDR pins has a
writable contract test.*
Given Phase 3's exported contract-test function, then each of the
five operator classes must have a writable present-key test from
locked documents. `contains` fails — the RDR says so ("cannot be
written from any current document") — and `contains` is the one
operator with a named empty-set masking hazard. A lock that pins
"absent set ≠ empty set" while unable to test "present malformed set
≠ empty set" leaves the hazard's other half to exactly the
discipline-based enforcement the RDR was rewritten to abolish.
