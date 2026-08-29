Model: claude-fable-5

# Premortem — Dual-breach precedence and reserved-key error carrier

Hostile premortem over the brief alone. No repository files were consulted; every finding is
a hypothesis the proposal must answer with evidence, not an assertion about the code.

## Findings ledger

| ID | passage/claim | Failure mode | Symptom user sees | Origin |
|----|---------------|--------------|-------------------|--------|
| P-1 | Claim 1: "at HEAD a dual-breach input already observably returns the escape-shape error" | Unverified control-flow claim. The brief itself says placement is "a cheapness preference" and no test pins it. An early return, a nil/empty-table fast path, or a reserved-key scan folded into input normalization above the shape check would make the "pin" a silent behavior change. | Producer pipelines that (by text-parsing) branch on the reserved-key message start seeing a shape aggregate on the same inputs after upgrade; their dual-breach handling path goes dead. | 2, 4 (S2) |
| P-2 | Claim 2: "the reserved-key contract explicitly says a future fixed order is not foreclosed" | Adjacent-clause misread. The quoted license is "EITHER error ... until a **joint decision** pins an order." The license clause is cited; the joint-decision condition is the adjacent clause. A unilateral pin by this record alone may not satisfy it. | Peer record owner reopens the decision; the pinned order and its test are reverted or re-litigated after callers have taken a dependency on shape-first. | 4 (S1) |
| P-3 | Claim 3/4: wrapping the three unexported `errors.New` values with `ErrReservedTagKey` | The three values cannot wrap anything as-is; they must be replaced (`fmt.Errorf("%w ...")` or a wrapper type). That changes (a) their `Error()` text — the ONLY classification channel existing producers have, per the brief — and (b) their identity for any internal `==` comparison. | Existing text-parsing producers misclassify reserved-key breaches as "unknown kernel error" and retry/loop; any internal identity check silently stops matching. | 1, 2 |
| P-4 | Claim 3 + locked "single error, no aggregate" clause | Mechanism unpinned: an implementer reaching for `errors.Join(ErrReservedTagKey, base)` produces a multi-unwrap aggregate with a newline-joined message — breaching the locked single-error clause while all proposed `errors.Is` tests still pass. | Reserved-key error prints as two lines; downstream single-line log pipelines truncate; contract breach ships green. | 2 |
| P-5 | Claim 4: "tests live in an external test package and cannot reference unexported vars" | Unverified inventory claim (borrowed authority about the test corpus). External tests CAN assert `err.Error()` strings without referencing the vars; internal white-box test files may also exist. Either breaks when the message/identity changes under P-3. | CI red on merge, "fixed" by editing production message text again — a second uncoordinated text change for parsers. | 2, 4 (S3-class) |
| P-6 | Claim 5: dual-breach pin via `errors.Is(err, ErrReservedTagKey) == false` | Inverted rejection: R3 was rejected because negative classification degrades. The chosen pinning test's load-bearing half IS a negative classification. It passes for any error lacking the sentinel and proves nothing about which check ran, only about what the returned error carries. | A refactor that stops running the reserved-key check entirely on dual breach (or on all inputs) keeps the test green. | 4 (R3 inversion) |
| P-7 | Claim 3/5: "single-breach tests assert the sentinel classifies on each of the three channels" | With no per-channel classification (by design) and no message assertions (by design), the tests cannot prove WHICH channel fired. If one channel's wrap is missed and the test input for it accidentally routes through another channel, coverage is claimed but false. | One of the three reserved-key breach paths returns an unclassifiable error in production; `errors.Is(err, ErrReservedTagKey)` is false exactly when a producer needs it. | 2 |
| P-8 | Claim 6: "a malformed table is malformed as a value independent of the input tuple" | Names a partition the shipped code may not have (S1 class). `Resolve(in Input)` — if the table travels inside `Input`, both breaches are properties of the same value and the "independent value" rationale is rhetoric, weakening the semantic case against R1. | None directly; the pin survives on inertia, not rationale — and inertia is what P-2 reopens. | 4 (S1) |
| P-9 | Claim 7: "authored file path rejects it earlier" + "CLI mapping unchanged is conforming" | Duplicate reserved-key knowledge: the earlier CLI-side rejection and the kernel check must agree on the key set. On drift (key added to kernel only), the kernel error reaches a CLI mapping that was deliberately left without a case for it. "Owned elsewhere" cites an owner/record the brief never names. | CLI user gets a generic unmapped error with the (possibly reformatted) kernel message leaking through, no dedicated exit code. | 2, 4 (S3) |
| P-10 | R4 rejection: "there is no payload a caller needs" | The reserved KEY is compile-time constant; the breach SITE is not. A generator author debugging which of N tag sequences (and at what index) emitted the key gets a bare sentinel. Proportionality was argued from the key, not from the locator. | Producer bisects generated inputs by hand to find the offending sequence; files a feature request that reopens R4 within months. | 4 |
| P-11 | R2 rejection: "cannot conform without reopening a locked record" | Scope-of-clause reading rests on borrowed authority. The locked clause constrains the shape-breach report when it is returned; whether it dictates the dual-breach return at all is a scope question the rejection asserts without quoting. If the locked record's scope is narrower, R2 was rejected for a reason the record does not support. | None at runtime; a later reviewer with the record open overturns the rejection reason and the alternatives analysis loses standing. | 4 (S3) |
| P-12 | Claim 1/5 interplay: no HEAD characterization evidence required | The proposal's own test is written post-change; nothing in the plan demands it pass at HEAD first. Claim 1 ("a pin, not a behavior change") is therefore never demonstrated, only asserted. | Same as P-1; additionally the RDR's central "no behavior change" evidence line is empty at review time. | 3 |

## 1. Failure narrative (prospective hindsight — it shipped, it failed)

The change landed in one commit: three `errors.New` values became `fmt.Errorf("reserved tag
key: %w", ...)` wrappers around the new exported `ErrReservedTagKey`, a dual-breach test
asserted shape-first, and the RDR locked.

Week 2: the largest programmatic producer — the only class of caller that can even construct
these inputs — upgraded. Their ingestion pipeline predated the sentinel, so it classified the
reserved-key breach the only way the old kernel allowed: matching message text. The wrap had
changed every one of the three messages. Reserved-key breaches now fell through to the
pipeline's "transient kernel error" branch and were retried on a backoff loop. A batch
migration stalled for a day before anyone read the raw logs. The postmortem line was exact:
"the proposal added the classification channel producers should use, and simultaneously broke
the classification channel producers did use, in the same release, with no migration note."

Week 6: the owner of the locked reserved-key record read the RDR and pointed at the clause the
proposal had quoted around: "until a **joint decision** pins an order." This record was not
joint — the peer owner had never signed. The pin was procedurally void. During re-litigation
someone re-ran the dual-breach test against the pre-change HEAD to establish the baseline
Claim 1 had asserted — and found the reserved-key scan had been folded into input
normalization for the empty-writes fast path, so a class of dual-breach inputs had returned
the reserved-key error all along. "A pin, not a behavior change" had been false on ship day,
untested, and unfalsifiable by the post-change test suite.

Week 9: a producer hit a reserved-key breach through the third channel. `errors.Is` returned
false. The wrap on that channel had been missed in review; its dedicated single-breach test
had been green the whole time because the input it constructed actually tripped the first
channel. Nobody could have caught it structurally — the proposal had explicitly ruled out
per-channel classification and message assertions, leaving no instrument that could tell the
channels apart.

## 2. Obstacle negation (claims 1–7)

**Claim 1 — "already observably returns the escape-shape error at HEAD."**
Negation: the brief concedes the order is undocumented, untested, and a "cheapness
preference." That is precisely the setting where an early return above the assumed sequence
(S2's escape class) lives: a fast path for empty write-sets, a normalization step that scans
tags at `Input` construction, a guard that short-circuits before the table walk. For any such
path, "pin" is a behavior change shipped under a no-change claim.
Scenario: dual-breach input with an empty escape row and a reserved key in the first tag
sequence takes the normalization path; HEAD returns the reserved-key error; post-change it
returns the shape aggregate; a producer's error-branch goes dead silently.

**Claim 2 — "contradicts no locked contract."**
Negation: the license clause read in full is conditional — either error is licensed *until a
joint decision pins an order*. The proposal cites the license and treats the condition as
satisfied by its own existence. If "joint" means the peer record's owner co-decides, a
solo record does not discharge it; the pin is open to procedural reversal after callers
depend on it — the worst time.
Scenario: peer owner reopens; order flips or reverts to unpinned; the exported sentinel and
the pinned test outlive the decision that justified them.

**Claim 3 — "an exported sentinel gives callers everything they need."**
Negation A: it gives *future* callers a channel while destroying the *current* callers' only
channel (message text) in the same change — the sentinel cannot be wrapped into plain
`errors.New` values without rewriting them, and rewriting them rewrites their strings.
Negation B: "the reserved key is a compile-time constant" answers *what* was breached, not
*where*. Three channels exist because three distinct breach sites exist; a producer fixing a
generator needs the site.
Scenario: producer with a 40k-row generated input gets `reserved tag key` and nothing else;
bisects by hand.

**Claim 4 — "breaks no existing test."**
Negation: the claim's mechanism ("external tests cannot reference unexported vars") does not
cover the two ways tests actually break: external tests asserting `err.Error()` output, and
internal white-box test files, both common in table-driven Go suites. The claim is an
inventory assertion made without an inventory.
Scenario: an external golden test of `Resolve` failure output fails on the new message
prefix; the fix commits a second message change, breaking text-parsers twice.

**Claim 5 — "fully structural, no message-text assertions."**
Negation: structural, yes — but half of it is a *negative* structural assertion, the exact
fragility class the proposal used to reject R3. `Is(ErrReservedTagKey)==false` cannot
distinguish "shape check won" from "reserved-key check was deleted," "sentinel wrap was
dropped," or "a third error class was returned that happens to carry the shape sentinel."
Scenario: refactor removes the reserved-key check from a code path entirely; dual-breach test
stays green; single-breach coverage (see Claim 3/P-7) fails to notice for one channel.

**Claim 6 — "shape-breach-first is the right winner."**
Negation: the independence argument ("malformed as a value independent of the input tuple")
presupposes the table is a value with a life outside `Input`. The signature in the brief is
`Resolve(in Input)`; if the table is a field of `Input`, both preconditions are properties of
one argument and the asymmetry claimed is a partition the code does not have (S1's class).
What remains is "richer diagnostics win" — a defensible tiebreak, but then R1's rejection
line "no stronger semantic rationale" cuts both ways.
Scenario: none at runtime; the rationale collapses exactly when P-2 forces re-litigation.

**Claim 7 — "leaving the CLI mapping unchanged is conforming."**
Negation: conformance today rests on the authored path rejecting reserved keys *earlier*,
i.e., a second copy of the reserved-key set outside the kernel. The decision to keep the CLI
mapping frozen removes the pressure to reconcile the copies. "A separate decision owned
elsewhere" names no owner and no record — a citation to authority not shown to exist.
Scenario: a fourth reserved key is added kernel-side; the CLI's early check passes it; the
kernel error surfaces through a mapping with no case for it; user sees a generic failure.

## 3. Consumer artifacts (what would have caught each at review time)

- **P-1/P-12 — HEAD characterization run.** Named artifact: `TestResolve_DualBreach_Order`
  executed against the pre-change commit and its pass recorded in the RDR evidence dir. Red
  at HEAD falsifies Claim 1 before lock; green at HEAD is the only thing that makes "a pin,
  not a behavior change" a fact rather than a sentence. Sweep must cover the input-corner
  matrix (empty write-set escape row × reserved key in first/middle/last sequence), not one
  input.
- **P-2 — joint-decision receipt.** Artifact: a line in the peer record's evidence (or a
  co-signed note) acknowledging this record as the joint decision the clause anticipates.
  Absence is a review-time BLOCK on the pin, independent of code.
- **P-3/P-5 — message-stability sweep.** Artifacts: (a) `grep -rn 'Error()' -g '*_test.go'`
  plus a sweep for the three current message literals across the module and known consumers;
  (b) golden-message tests `TestReservedKeyError_MessageUnchanged_{Ch1,Ch2,Ch3}` asserting
  the post-wrap `Error()` strings byte-equal the pre-change strings (wrap with
  `fmt.Errorf("%w", ...)`-style construction that preserves text, or a wrapper type whose
  `Error()` delegates). If message change is instead *chosen*, the artifact is a migration
  note in the RDR — the choice must be visible, not incidental.
- **P-4 — carrier-shape test.** Named test: `TestReservedKeyError_SingleError_NotAggregate`
  asserting the returned error does not implement `Unwrap() []error` and unwraps linearly to
  the sentinel. Pins `%w`-chain over `errors.Join` and re-verifies the locked "no aggregate"
  clause structurally.
- **P-6 — positive-order instrument.** Named test: a white-box (same-package)
  `TestResolve_DualBreach_ReservedCheckStillRuns` companion, or a documented decision that
  order is pinned only up to the negative assertion's blind spots. The review artifact is the
  acknowledgment; today the RDR text claims more than the test proves.
- **P-7 — per-channel proof.** Named tests: white-box `TestReservedKey_Channel{1,2,3}_Wraps`
  asserting `errors.Is(errReservedX, ErrReservedTagKey)` directly on each unexported value —
  three one-line tests that make "each channel wraps" a checked fact instead of a routed
  guess. (External-package purity is not worth an unverifiable coverage claim.)
- **P-9 — drift lock.** Named test: `TestCLIEarlyReject_KeySetMatchesKernel` comparing the
  CLI's early-reject reserved-key list against an exported (or test-exported) kernel list;
  plus, in the RDR, the actual identifier of the record that "owns" the CLI error-code
  decision. A user journey: feed each kernel-reserved key through the programmatic path and
  assert the CLI mapping produces a non-generic diagnostic — or record that it doesn't, on
  purpose.
- **P-10 — producer journey.** Journey artifact: "generator author receives reserved-key
  breach on a 10k-sequence input; enumerate the steps to locate the offending sequence using
  only the returned error." If the answer is "bisect," R4's rejection reason should say so
  honestly ("locator payload deferred, cost accepted"), not "no payload a caller needs."
- **P-11 — quoted-scope check.** Artifact: the locked shape record's clause quoted verbatim
  in the RDR alternatives section, showing whether its verbatim/one-level constraint governs
  the dual-breach return or only the shape-breach report. One quotation retires the S3 risk.

## 4. Refutation targets (S1–S3 applied)

- **S1 — misread adjacent clause / phantom partition.** Two hits. (a) Claim 2 leans on the
  either-error license while the adjacent conditional — "until a **joint** decision" — is the
  clause that governs whether this record may pin at all (P-2). (b) Claim 6's "malformed as a
  value independent of the input tuple" names a table/tuple partition that the stated
  signature `Resolve(in Input)` does not exhibit; if the table rides inside `Input`, the
  partition exists only in the argument (P-8).
- **S2 — control-flow order the argument assumes but code may negate.** Claim 1 is exactly
  this shape: "today the shape check happens to run first" is asserted about a placement the
  brief itself calls an unobservable cheapness preference, with no test at HEAD. Any early
  return or normalization-time tag scan above the shape check inverts the claim, and the
  proposal contains no artifact that would detect it before ship (P-1, P-12). The proposed
  test is written on the far side of the change, so it can never distinguish "pinned existing
  behavior" from "changed behavior, then pinned the change."
- **S3 — borrowed authority.** Three hits. (a) Claim 4's "tests live in an external test
  package" is an uninventoried claim about the test corpus doing load-bearing work for "breaks
  no existing test" (P-5). (b) Claim 7's "owned elsewhere" cites an owner and a decision that
  are never named (P-9). (c) R2's rejection rests on a scope reading of the locked shape
  record — that its verbatim clause reaches dual-breach returns — asserted without quotation
  (P-11). Additionally, R3's rejection reason ("negative classification degrades") is itself
  inverted by the chosen dual-breach test, whose pinning half is a negative classification
  (P-6): a rejection reason the proposal violates in its own mechanism.

## Disposition

No finding forces abandoning the approach: shape-first pinning plus an exported sentinel is a
coherent end state, and R1–R4 remain rejectable (some for better reasons than those given).
But Claims 1, 3, 4 as written are undefended, and P-2's joint-decision condition is a
procedural gate the record cannot waive for itself. The approach survives only with the
artifacts above folded in — most critically the HEAD characterization run, message-text
stability (or an explicit migration note), the `%w`-not-`Join` carrier pin, per-channel wrap
proofs, and the joint-decision receipt.
