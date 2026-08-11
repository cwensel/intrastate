Model: claude-fable-5

# Hostile critique — RDR 0009, escape-row shape conformance ownership

Reviewed against the RDR at
`docs/rdr/0009-escape-row-shape-conformance-ownership.md` and the shipped
source it cites: `internal/resolve/resolve.go`,
`internal/resolve/fixtures_test.go`, `internal/resolve/adversarial_test.go`,
`internal/resolve/fixup_test.go`, `internal/cli/clierr/clierr.go`,
`internal/cli/root.go`, `internal/cli/respond/respond.go`.

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | §Implementation Plan, "Phase 1: Exported predicate + kernel entry precondition" vs "Phase 2: Fixture conformance" | Phase 1 lands the entry check while `fixtures_test.go::escapeRow` (line 257) still populates `Writes` on all sixteen call sites; the frozen suite goes red mid-plan and the RDR never says the phases must land atomically | Implementer's Phase 1 commit fails ~16+ frozen ADV/MVV/Fixup tests via `mustResolve` fatals; CI/roborev red; implementer improvises the recovery the RDR should have specified | §1, premortem, AT-1 |
| C-2 | A2 "An *unwrapped* Go error escaping a future verb's `RunE` exits **1**, not 2"; Prerequisites "unwrapped, it would exit 1, not 2"; scenario 11 "not exit 1" | Factually false against shipped code: `root.go::ExecuteAndEmit` converts every non-`CLIError` via `cobraErrorToCLIError` (GroupUserEnv) *before* `ExitCodeFor` runs, and GroupUserEnv is exit **2**. The exit-code tripwire the RDR relies on to detect a missing wrap does not exist | Unwrapped breach exits 2 with `Code: "command-error"` and Hint "run with `--help` to list supported commands and flags" — misclassified as user error with a nonsense remedy; the missing wrap ships undetected because both wrapped and unwrapped exit 2 | §1, premortem, AT-2 |
| C-3 | §Technical Design "Resolve returns the zero `Result` … callers MUST check the error before reading the disposition"; scenario 7b | Breach return violates `resolve.go::Result`'s own doc contract ("exactly one transition plan or exactly one typed refusal, never both and **never neither**") and `Resolve`'s "returns exactly one disposition"; the doc-amendment clause covers only the numbered evaluation-order list, not these two contracts | First real caller branches on `Refused()` first (the pattern every existing doc line teaches), reads `got.Plan` → nil-pointer panic with a stack trace at the CLI instead of any envelope | §1, §3, premortem, AT-3 |
| C-4 | §Normative Contracts multi-breach clause: "Rows carrying no source identity collapse to a single RowRef{"",""} entry"; "Count == 3" | Collapse-plus-Count destroys row localization for the RDR's *only* existing producer population (A5: hand-built rows), by misapplying REQ-2/REQ-10 — replay-determinism rules for *modeled* payloads — to the error channel the RDR itself argues is outside the modeled surface | Developer sees "escape-row shape breach at `RowRef{"",""}`, Count 3" and must bisect their table by hand to find which three rows breach | §1, §2, premortem, AT-4 |
| C-5 | A4 "the predicate … does NOT extend to NextTags"; Prerequisites "RDR 0004's implement stage keeps the write accessor scoped…"; Phase 2 "`NextTags` stays on the builder" | The invariant "an escape row bears no owned-state mutation" is enforced against half the mutation vocabulary; `planOf` copies `NextTags` escape-unaware, safety rests on an unimplemented peer RDR held only by a checkbox — no fixture, no test, unlike Phase 3's binding of RDR 0002 — while this RDR's own fixtures normalize `NextTags: {status: Blocked}` on escape rows | An RDR 0004/flow implementer persists the escaped plan's `NextTags` as the new state; an owned tag takes a value no authored rule set — the Problem Statement's opening symptom, verbatim, with every RDR 0009 test green | §3, premortem, AT-5 |
| C-6 | §Normative Contracts "the predicate is `Table.CheckValid() error`"; §Risks "the contract is one predicate, not an open validation framework" | A general-validity name pinned to a single-property predicate on a locked package surface: `CheckValid` returns nil for tables violating the documented Escape-class restriction (`Row.Escape` doc: "RDR 0002 restricts the list to no_match and ambiguous_match" — unchecked), duplicate RuleIDs, etc., and the "one predicate" clause forbids it ever growing | Producer calls `table.CheckValid()`, gets nil, ships a table with `Escape: [guard_unevaluable]` that silently never rescues; "valid" meant one invariant out of many | §2, premortem, AT-6 |
| C-7 | §Normative Contracts CLI clause "The chosen carrier is a NEW omitempty field … subject to A9"; A9 Status: Pending; scenario 11 DEFERRED | Normative text legislates a wire contract on a Pending assumption, for a verb that does not exist, without specifying the field's name, type, or whether `Count` serializes — and ignores that `clierr` is a deliberate leaf package, so a `[]resolve.RowRef` field couples the generic envelope to the kernel or forces a duplicated type. Violates the RDR's own Finalization Gate rule (no Pending assumption under settled-fact prose) | The RDR that actually ships the `flow` verb supersedes this clause; until then, two enforcement halves of the "identities reach the operator" promise exist only on paper | §2, premortem, AT-7 |
| C-8 | §Trade-offs "Staged, not immediate"; §Risks first Risk (kernel precondition treated as *the* enforcer) | The Problem Statement's named user — the table author — receives nothing from this RDR: Phases 1–2 protect a producer population of zero (A5) and Phase 3 binds an unimplemented normalizer. The first real authored table arrives through whatever ad hoc loader the flow verb ships with, making the kernel backstop the de facto authoring UX — the RDR's own Risk #1, mitigated only by fixtures a not-yet-run implement stage may or may not honor | Table author's first malformed escape rule produces exit 2, an internal-error code, and "fix the table producer" — they *are* the producer, and were promised a load-time diagnostic naming their rule | §3, premortem, AT-5 |
| C-9 | §Normative Contracts "Resolve MUST return CheckValid's error VERBATIM, never fmt.Errorf-wrapped"; "Unwrap() []error MUST be EXACTLY ONE LEVEL deep" | Over-pinned `errors.Join` mechanics freeze implementation internals as contract: Resolve may never add context to its own error stream, and any future second entry-time error source breaks the guaranteed flat traversal — so the first extension is a contract-breaking change by construction | Callers written to the pinned one-level traversal break the day Resolve returns any second error kind; until then Resolve's errors carry no "resolve:" context anywhere in the program | §2, premortem, AT-1 |

## 1. The three most likely ways implementation goes wrong

### 1a. Phase 1 turns the tree red, and the RDR never told the implementer it would (C-1)

**Root cause in the RDR.** The Implementation Plan orders "Phase 1: Exported
predicate + kernel entry precondition" strictly before "Phase 2: Fixture
conformance." But the A3 spike — the RDR's only empirical evidence — validated
the *opposite* ordering: fixtures conformed first, suite green, no check
present (`a3-conformed.out`). Nothing in the RDR validates, or even mentions,
the state the plan actually passes through: check present, fixtures
unconformed.

**The enabling passage.** Phase 2's own text — "Because no remaining fixture
breaches, the mutation-style step … is load-bearing here" — presumes Phase 1's
check and Phase 2's conformance coexist, yet the phases are presented as
sequential, separately committable steps, in a project whose conventions
(Conventional Commits, roborev per commit) make per-phase commits the default.

**What actually happens.** `fixtures_test.go::escapeRow` (line 257) sets
`Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}}` on every one of its
sixteen call sites, plus the two overrides at `adversarial_test.go:229` and
`fixup_test.go:103`. The moment `Resolve` gains the entry precondition, every
escape-bearing table in the frozen suite is a breaching table. `mustResolve`
fatals on error first (the RDR says so itself, in scenario 7b's justification).
Result: the implementer lands Phase 1, runs
`go test ./internal/resolve/...`, and watches on the order of sixteen-plus
frozen ADV/MVV/Fixup/boundary tests die — including ADV-2b and Fixup-1e, whose
*point* is the write-bearing escape row.

**Symptom the user sees.** A red CI run on the Phase 1 commit; a roborev
review flagging broken frozen evidence; an implementer forced to choose —
squash Phases 1–2, invert them, or temporarily weaken the check — with zero
guidance from a document that pinned the spelling of a sentinel down to the
letter but not whether its two phases can land independently.

### 1b. The exit-code tripwire A2 builds the CLI story on does not exist (C-2)

**Root cause in the RDR.** A2's "Verified (with a named wiring obligation)"
stitches two individually true citations into a false composite. It is true
that `clierr::ExitCodeFor` defaults a non-`CLIError` to 1, and true that
`root.go::cobraErrorToCLIError` stamps `GroupUserEnv`/`command-error`. But the
composite claim — "An *unwrapped* Go error escaping a future verb's `RunE`
exits **1**, not 2" — is false against the shipped wiring, because the
conversion runs *before* the exit-code mapping ever sees a bare error:
`Execute` → `ExecuteAndEmit` → `errors.As` fails → `respond.Fail(cmd,
cobraErrorToCLIError(err))` → `respond.Fail` returns the `*CLIError` →
`ExitCodeFor` sees `GroupUserEnv` → **exit 2**. The `ExitCodeFor`
non-`CLIError` default-to-1 branch is unreachable from the production entry
point. Nobody verifying A2 traced the one function (`ExecuteAndEmit`) that
decides the question.

**The enabling passage.** A2's Evidence block, repeated in Prerequisites
("unwrapped, it would exit 1, not 2") and in scenario 11 ("not exit 1, which
is what an unwrapped kernel error would produce"), and in the Normative
Contracts CLI clause ("Returning the kernel error UNWRAPPED is a defect …
the documented exit-2 behavior does not hold without the wrap").

**What actually happens.** The wrap obligation is real, but its stated
detection signal — a wrong exit code — is fictional. Wrapped and unwrapped
both exit 2. So a future `flow` verb that forgets the wrap passes any exit-code
check anyone writes from this RDR, and ships. What differs is everything the
exit code doesn't capture: `Code: "command-error"` instead of
`"escape-row-shape-breach"`, `Group: GroupUserEnv` (a *user-input*
classification for an internal invariant breach), and the Hint "run with
`--help` to list supported commands and flags" — actively wrong advice for a
broken table producer. The RDR never names this, the actual failure surface,
because it believed the failure surface was exit 1.

**Symptom the user sees.** An operator with a broken producer gets exit 2,
a `command-error` code, and a suggestion to read `--help`. Scripts branching
on the stable code the RDR promised never fire. And any conformance test
written to scenario 11's "not exit 1" oracle asserts a distinction that
cannot occur, so it either fails against a correct implementation's fixture
harness or is silently rewritten to assert exit 2 — at which point it can no
longer tell wrapped from unwrapped at all.

### 1c. The zero-`Result` return breaks the kernel's own published contract, and only one doc list gets amended (C-3)

**Root cause in the RDR.** On breach, `Resolve` returns
`Result{Plan: nil, Refusal: nil}` plus a non-nil error. `resolve.go::Result`'s
doc contract says, today, in the locked RDR 0001 surface: "exactly one
transition plan or exactly one typed refusal, **never both and never
neither**." `Resolve`'s own doc says "It returns exactly one disposition for
the input tuple." Both sentences become false. The RDR's doc-amendment clause
(Load-Bearing Decisions, "*Doc-contract amendment*") authorizes exactly one
edit: prepending step 0 to `Resolve`'s numbered evaluation-order list. It
never mentions `Result`, `Refused`, or the "exactly one disposition" sentence,
and it asserts "the amendment is documentation-only" — while shipping a return
state RDR 0001's contract says cannot exist.

**The enabling passage.** §Technical Design: "Because the zero `Result`
reports `Refused() == false`, callers MUST check the error before reading the
disposition; a caller that branches on `Refused()` first would read a
success-shaped value with a nil `Plan`." The RDR sees the trap, writes one
test for it (7b), and then leaves the trap armed: the entire existing corpus —
every doc comment, every frozen test — teaches the `Refused()`-then-`Plan`
idiom. The RDR's mitigation is a MUST aimed at future callers who will learn
the idiom from the code, not from this RDR.

**Symptom the user sees.** The first `flow` verb (or any programmatic caller)
written by someone pattern-matching on the existing suite does
`res, _ := resolve.Resolve(in); if !res.Refused() { use(res.Plan) }` — and
panics on a nil `Plan` dereference. At the CLI that is a Go stack trace, the
one output the never-silent envelope contract exists to prevent — a worse
surface than the unowned invariant this RDR set out to fix.

## 2. The one section that will be rewritten within 6 weeks of shipping

**The CLI-surfacing normative clause ("When a breach reaches the CLI…"),
together with its satellites A9 and scenario 11.** (C-7, with C-4 and C-9 as
accelerants.)

It will be rewritten because it is the only clause in the document that
legislates a *wire contract* — and it does so with every input missing:

- Its load-bearing choice (a new `omitempty` field on `CLIError`, not
  `Detail`) is explicitly "subject to A9," and **A9 is Pending**. The RDR's
  own Finalization Gate forbids exactly this: "no assumption marked `Pending`
  … may have settled-fact prose elsewhere in the RDR depending on it." The
  clause is settled-fact prose in MUST form sitting on a Pending assumption.
- The field itself is unspecified: no name, no JSON key, no type. And the
  type question is not clerical: `clierr` is a deliberate leaf package
  ("lives in its own package so other internal packages … can construct
  CLIErrors without importing internal/cli"). Carrying `[]resolve.RowRef`
  means the generic error envelope imports the resolution kernel; not
  carrying it means minting a duplicate row-identity type in `clierr` that
  RDR 0005 never authorized. The RDR chose neither, and A9's verification
  plan does not even ask the question.
- The kernel-side contract insists `Count` is carried "STRUCTURALLY … never
  only in formatted prose." The wire clause then forgets `Count` exists:
  scenario 11 asserts identities are readable from JSON and says nothing
  about multiplicity. The structural-count mandate dies at the process
  boundary, which is the only boundary an operator ever sees.
- Its one test (scenario 11) is deferred to a verb that does not exist
  (`root.go::NewRootCmd` registers only `newVersionCmd()`), and its one
  falsifiable behavioral claim ("not exit 1") is false today (C-2).

Under this project's own doctrine — "we never amend rdrs" — "rewritten" means
the RDR that actually ships the `flow` verb will restate this clause from
scratch, and RDR 0009's version will survive only as the history of a wire
contract designed with zero consumers, zero field schema, and one wrong exit
code. Six weeks is generous; it lasts exactly until the first verb touches
`Resolve`.

## 3. The one assumption that will not survive first contact with a real user

**A4 — "the predicate is `Writes`-only and does NOT extend to `NextTags`."**
(C-5, compounded by C-8.)

A4 is "Verified — settled closed" entirely by *reading unimplemented peer
RDRs*: RDR 0004 (Final, unimplemented) says a write accessor applies only
planned writes; `NextTags` appears nowhere in its text; therefore an escape
row's `NextTags` cannot mutate owned state. Every link in that chain is prose
about code nobody has written.

Here is the contact with reality. `resolve.go::planOf` copies `NextTags` and
`Writes` symmetrically, escape-unaware — the RDR admits this ("this safety
rests entirely on RDR 0004's write-accessor scoping, not on any kernel-local
property"). An escaped plan therefore ships
`Plan{NextTags: [{status Blocked}], Writes: nil, Escaped: true}` — and the
RDR's own Phase 2 *guarantees* it, because "`NextTags` stays on the builder":
the conformed `escapeRow` fixture still carries `NextTags: {status: Blocked}`.
Now put a real implementer in front of `Plan.NextTags`, whose doc comment
reads, in full: "NextTags is the next state." What does a person building a
state-machine CLI do with "the next state"? They persist it. The distinction
this RDR's entire safety argument rides on — next-state tags are *display*,
write tags are *persistence* — lives in two RDR documents and one doc-comment
nuance, guarded by a Prerequisites checkbox on a *different RDR's* implement
stage, with no fixture and no test. Compare the treatment RDR 0002 gets: a
named, normative, shared Phase 3 fixture set binding its build. RDR 0004 — the
layer "where the harm actually occurs," in the RDR's own words — gets a
checkbox.

When it breaks, it breaks silently and completely: the accessor persists an
escaped plan's `NextTags`, an owned tag takes a value no authored rule set,
and the user experiences the Problem Statement's opening paragraph *verbatim*
— "the escape fires, the route is right, and afterwards an owned tag has a
value no rule in their table ever set" — while `Table.CheckValid` returns nil,
every scenario 1–10b passes, and the RDR that exists to prevent exactly this
symptom reports total success. A4's "If wrong" clause even predicted the
recurrence ("the user symptom recurs through `NextTags`") and the RDR shipped
the predicate narrow anyway, on the strength of documents with no
implementations behind them.

## 4. Premortem

*Written at draft time, as if the failure has already happened.*

It is late September 2026. RDR 0009 is marked Implemented, and the escape-row
invariant it owns has been breached in production analysis twice, panicked a
CLI once, and cost the team a two-day fixture archaeology in between.

The trouble started at Phase 1. The implementer followed the plan as written:
export `Table.CheckValid`, call it as `Resolve`'s first statement, commit.
`go test ./internal/resolve/...` returned eighteen failures —
`fixtures_test.go::escapeRow` still stamped `Writes: {status Blocked}` onto
all sixteen call sites, and `mustResolve` fataled before a single frozen
assertion ran. ADV-2b and Fixup-1e, which the RDR had catalogued down to their
line numbers, were among the dead. The RDR had proved (A3) that conformed
fixtures pass without the check; it had never asked whether the check passes
without conformed fixtures. Under deadline, the implementer squashed Phases 1
and 2 into one commit — which meant the Phase 2 "mutation-style step" was run
against a tree where the check and its oracle were born together, and Mutant
B's discrimination claim was quietly re-derived rather than demonstrated.

The panic came with the `flow` verb, six weeks later. Its author, reading
`resolve.go` for the calling convention, found `Result`'s contract — "exactly
one transition plan or exactly one typed refusal, never both and never
neither" — and wrote the idiom every frozen test taught:
`if !res.Refused() { plan := res.Plan; … }`. The first malformed hand-built
table produced `Result{nil, nil}` plus an error the caller checked second,
and the CLI printed a Go stack trace — the exact never-silent violation the
respond gateway exists to prevent. Scenario 7b had asserted the zero-`Result`
shape; nothing had amended the two doc contracts that promised it could not
exist.

Once the verb wrapped errors properly, the operator experience failed anyway.
The first breach report read: `escape-row shape breach: 3 rows at "":""`,
`Count: 3` — three hand-built rows with no source identity, collapsed to one
`RowRef{"",""}` entry per the multi-breach clause, positional information
forbidden by a REQ-2/REQ-10 rule written for *modeled* refusal payloads and
imported wholesale into the programmer-mistake channel the RDR itself argued
was categorically different. The developer bisected the table by hand. The
week before, a teammate had returned the kernel error unwrapped, and nothing
caught it: the RDR promised the mistake would exit 1, but
`ExecuteAndEmit` → `cobraErrorToCLIError` → `GroupUserEnv` exits 2, so the
"tripwire" test asserted a distinction that does not exist, and the operator
got `command-error` with the hint "run with `--help` to list supported
commands and flags."

The breach that mattered, though, never touched `CheckValid` at all. The
accessor-layer implementer, building against RDR 0004, read
`Plan.NextTags` — "NextTags is the next state" — and persisted it, as anyone
implementing a state machine would. The conformed fixtures had normalized
escaped plans carrying `NextTags: {status Blocked}` with empty `Writes`, so
the pattern looked blessed. An escape fired; `status` became `Blocked`; no
authored rule had set it. The table author filed the exact bug from the
Problem Statement's first paragraph. Every RDR 0009 test was green:
`Table.CheckValid` returned nil (it checks one field of one invariant — it
had also been returning nil for a table whose `Escape` list carried
`guard_unevaluable`, silently rescuing nothing), scenarios 1–10b passed, and
the A4 checkbox in Prerequisites was still unchecked because RDR 0004's
implement stage had never been told it existed.

The post-incident review concluded what the draft-time record already showed:
the RDR spent its precision budget pinning sentinel spellings, join depths,
and verbatim-return mechanics, and spent nothing on the four seams where it
actually failed — phase ordering, the caller contract, the wrap tripwire, and
the `NextTags` half of the mutation vocabulary.

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 — Phase-1-alone is green** (catches C-1; would also have surfaced C-9's
brittleness by forcing the error surface to exist before its consumers)
```gherkin
Given the repo at HEAD
When only Phase 1 is applied (Table.CheckValid + entry precondition,
  no fixture edits)
Then `go test ./internal/resolve/...` passes
```
Fails today: sixteen `escapeRow` call sites breach; ~18 frozen tests fatal in
`mustResolve`. Forces the RDR to either invert Phases 1/2 or declare them
atomic.

**AT-2 — the unwrapped-error tripwire is real** (catches C-2)
```gherkin
Given a stub verb whose RunE returns the bare kernel breach error
When executed through cli.ExecuteAndEmit
Then the process exit code is 1
```
Fails today: `cobraErrorToCLIError` stamps `GroupUserEnv` → exit 2. Reviewing
this red result rewrites A2's obligation around the actual hazard
(misclassification + wrong hint), not the fictional one (exit 1).

**AT-3 — the published Result contract survives the design** (catches C-3)
```gherkin
Given the breach return shape specified in Technical Design
When compared against resolve.go::Result's doc contract
  ("never both and never neither") and Resolve's
  ("returns exactly one disposition")
Then no shipped contract sentence is falsified,
  or the RDR's amendment clause names every falsified sentence
```
Fails today: the amendment clause names only the numbered evaluation-order
list. Passing requires either amending both contracts explicitly or choosing
a breach shape that keeps them true.

**AT-4 — the diagnostic localizes rows for the actual target population**
(catches C-4)
```gherkin
Given a hand-built table (no RuleID, no SourceLocator) with three
  breaching rows among ten
When Resolve reports the breach
Then the report lets the producer find each offending row without
  bisecting the table
```
Fails against the spec as written: the report is `RowRef{"",""}, Count: 3`.
Forces the review to confront that REQ-2/REQ-10 govern modeled payloads, and
that the error channel — which the RDR's own taxonomy argument places outside
the modeled surface — could carry positional or content diagnostics without
touching those REQs.

**AT-5 — the invariant holds end-to-end through the blessed field**
(catches C-5, and C-8's staging gap in the same breath)
```gherkin
Given the conformed escapeRow fixture (Writes empty, NextTags
  {status: Blocked}) selected via the escape path
When the emitted Plan is handed to a consumer that persists
  "the next state" as Plan.NextTags's doc comment describes it
Then no owned tag takes a value no authored rule set,
  and some artifact of THIS RDR (fixture, test, or normative fixture
  set binding RDR 0004) fails if it does
```
Fails today: no such artifact exists — RDR 0004 gets a Prerequisites checkbox
where RDR 0002 gets a normative Phase 3 fixture set. Symmetric treatment is
the minimum fix.

**AT-6 — the validator's name matches its coverage** (catches C-6)
```gherkin
Given a table whose escape row declares Escape: [guard_unevaluable]
  (violating Row.Escape's documented restriction to no_match and
  ambiguous_match)
When Table.CheckValid is called
Then it returns non-nil, or the exported name and doc comment state
  that only escape/Writes conformance is checked
```
Fails today: nil, under a name that says "valid," with a normative clause
forbidding the predicate from ever growing into the name it was given.

**AT-7 — the wire clause is implementable from its text alone** (catches C-7)
```gherkin
Given only the Normative Contracts CLI clause and A9
When an implementer writes the envelope serialization test for a
  two-breach report
Then the field name, JSON key, Go type, and Count serialization are
  all derivable from the RDR, and no Pending assumption is load-bearing
```
Fails today on all four counts, plus the Finalization Gate's own
Pending-under-settled-prose rule. Either A9 resolves before the clause goes
normative, or the clause is demoted to a recorded intention for the
flow-verb RDR to own.
