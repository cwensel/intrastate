Model: claude-opus-5

# Hostile critique — cli/0024 Declared emit vocabulary

Grounded against `intrastate` at `4c545df`, with live decoder probes against
`github.com/pelletier/go-toml/v2 v2.2.4` (the repo's own version, per
`go.mod`).

This RDR is 1846 lines of argument for a feature whose merge-day outcome, by
its own admission (`0024:§consequences`), is one example model and an empty
JSON field. It has been through a premortem, a cove pass and a 3amigo pass,
and the residue of those passes is visible as defensive prose bolted onto
contracts that still contain arms that cannot be implemented as written. The
record is now so heavily armored against the *last* round of criticism that it
has stopped being readable as a specification. What follows is what survives
grounding.

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0024:C1` (refusal arms list, "a disposition sub-table carrying no keys") | Two enumerated `malformed_emit_declaration` arms — "an `enum` with no domain" and "a disposition sub-table carrying no keys" — decode to the *identical* state (`Domain == nil`) under the repo's decoder. The arm is undistinguishable, so it cannot be implemented, and S1 demands a fixture per arm. | Author writes `[emit.next.domain]` with the members still to come, gets a refusal whose message says "an `enum` declares no domain" — pointing at a `domain` sub-table they can see on screen. | §1, premortem, AT-1 |
| C-2 | `0024:C1` ("Strictness still catches a typo at the DECLARATION level (`domaim = [...]` refuses)") | The declaration-level typo refuses as `unknown_schema_field`, a *different* category than the three this RDR registers. C1 presents it as though the emit grammar catches it; it does not, and the finding an author gets names a schema field, not a declaration. | `kind` misspelled as `kimd` yields `unknown_schema_field: strict mode: fields in the document are missing in the target struct` — no mention of emit, of the key, or of what to write. | §1, AT-2 |
| C-3 | `0024:F1` ("offending file in `locator`"), `0024:MVV` step 2 ("with its category slug in `code` and a locator") | `atLine` is wired at exactly one call site in the whole loader (`load.go:214`, `loadTags`/`tagHeaderLine`). No emit analogue is specified or planned. Every emit refusal will carry `Line == 0` and render as `file:1`. | On a 54-rule model the author gets `models/foo.toml:1` for a defect on line 400, once per run, N times. | §1, §3, premortem, AT-3 |
| C-4 | `0024:A8`, `0024:C3` ("the property `TestReq146_EveryEmittedSequenceIsASortedSlice` already pins"), `0024:§phase-2` ("fails `TestReq146`'s sorted-sequence sweep") | A8 is not merely Pending — it is **false at HEAD**. `TestReq146` iterates a hand-written list of `table.Row` field names and its determinism leg compares `got.Rows`. `Model.EmitDecls` is a model-level field; the sweep can never reach it. Phase 2's stated failure signal will not fire. | An unsorted `Domain` union ships green; a consumer diffing two `dump`/export runs of the same model sees the member list reorder between runs. | §1, §2, premortem, AT-4 |
| C-5 | `0024:S6` (category witnesses, "per the convention `internal/table/dump_test.go`'s REQ-119 witness map") | The convention is REQ-119 **and REQ-131**: each witness fixture must trip its category **"and no other"**. A `malformed_emit_declaration` witness must therefore be a model that is otherwise fully valid *and* declares emit — but the fail-fast ordering C2 leaves "deliberately unspecified" means the witness's category is not stable under a reordering of the step slice. S6 asserts registration, never the exclusivity half. | A later loader refactor silently reclassifies the emit witness; the category survives in `Categories()` with a fixture that no longer trips it. | §1, AT-5 |
| C-6 | `0024:C4`, `0024:§phase-3` (ii) ("their mechanical regeneration is authorized HERE by name") | The 28 goldens are not a fixture corpus. Their file-level doc binds them to `0011`'s REQ-118/S8 — "captured from the PRE-change build of this tree, per A-11" — and the assertion message states "every one of these payloads is **byte-identical before and after**". Regenerating them does not update a fixture; it destroys the pre-change reference REQ-118 is discharged against. 0024 self-authorizes an edit to a peer RDR's *evidence*, having explicitly ruled that the same class of edit to `TestReq39` is "an amendment to a peer RDR's spec test". | REQ-118's byte-identity guarantee is silently voided; a future change to the payload has no pre-change baseline to compare against and nothing reports its loss. | §1, §2, premortem, AT-6 |
| C-7 | `0024:C2` ("The six steps between `loadTags` and `normalizeRules` (`loadAccessors`, `loadDump`, `loadContexts`, `loadInitial`, `loadTerminal`)") | The RDR counts six and names five. There are five (`load.go:85-89`). A normative contract that cannot count the steps it is inserting into is not a contract an implementer can place code from. | Implementer inserts at the wrong index; the fail-fast order shifts and the MVV's step-2 two-run sequence reports the defects in the other order. | §1, AT-7 |
| C-8 | `0024:A4` ("`conformDomain`'s enum arm is exactly the membership test C2 needs"), `0024:C3` ("constructs a throwaway `TagDecl{Kind: d.Kind, Domain: d.Domain}`") | `conformDomain`'s enum arm is guarded by `len(decl.Domain) > 0`. `conformKind` has **no arm for `enum` and no arm for `scalar`** — an unrecognized kind conforms everything. So `ConformValue` returns nil for every value under an `enum` declaration whose domain is empty, and for every kind token C1 does not itself reject. The reuse is only safe because C1 promises to hand-reject those states first; the RDR never says the reuse depends on that ordering. | An author writes `kind = "enums"` (plural typo) or an enum whose domain C1's arms let through, and every value passes — the model reads as declared and checked, and proves nothing. | §1, §3, premortem, AT-8 |
| C-9 | `0024:§approach`, `0024:§consequences` ("At merge this RDR ships capability, not coverage"), `0024:A5` | The only real adopter's models (`rdr-status.toml`, `rdr-write.toml`) live in a different repository, and the RDR states adoption there is "out of scope here and unscheduled". Everything the record justifies — the disposition partition, the three-valued split, `dispositions` on the wire — is validated against models this repo cannot lint, cannot test, and does not ship. | The declaration grammar ships; the consumer's first serious attempt to adopt it discovers a grammar mismatch, and there is no in-repo artifact against which to argue the grammar was ever right. | §2, §3, premortem, AT-9 |
| C-10 | `0024:§consequences` ("first adoption on a large model is an N-round fix-and-rerun loop"), `0024:C2` (fail-fast, order unspecified) | The RDR names the adoption cost accurately and then declines to pay it, on the grounds that fail-fast "is a cost of the tier this RDR lands in rather than of the declaration itself". That is a category error: no prior load category is *opt-in strictness applied retroactively to an existing corpus of authored values*. Every other load category fires on a single authoring mistake; this one fires N times on a correctly-authored model whose author has just switched on the check. | Adopter on a 54-rule model runs `lint`, fixes one thing, reruns, 40+ times, each run pointing at `file:1` (C-3). They revert the declaration or declare everything `scalar`. | §1, §2, §3, premortem, AT-10 |
| C-11 | `0024:C1` ("`scalar` is the declared-but-unvalidated escape hatch"), `0024:§risks-and-mitigations` (scalar hollowing) | The mitigation for the escape hatch is "it is visible in the one reviewable declaration table" and an advisory finding is "deliberately NOT taken" because `0006:C17` closes the advisory tier. The RDR therefore ships an escape hatch, cites Meyer's condition that a partial regime should flag its loopholes, and then declines the condition on a *procedural* ground (a peer contract is closed) rather than a design one. | A model reads as fully declared. Every key is `scalar`. `intrastate lint` exits 0 and reports nothing. The bargain the Problem Statement opens with — "lint proves it" — is now false in a way no output distinguishes from the true case. | §1, §3, premortem, AT-11 |
| C-12 | `0024:C4` ("`omitempty` is available on this field and is deliberately not taken") | Never-omitted `dispositions: {}` makes every one of the 28 goldens and every inline payload assertion diff for models that will never declare anything. The stated reason — a consumer must distinguish "declared nothing" from "binary predates the field" — is not achievable anyway: an old binary handed a *model* with `[emit]` refuses (A1), but an old binary handed a payload-consuming caller emits no `dispositions` key at all, which is exactly the absent case `{}` was supposed to disambiguate from. The distinction the field is never-omitted *for* is undetectable from the consumer side. | Nothing, on merge day. That is the point: the whole 28-golden diff buys a distinction no consumer can observe. | §1, §2, AT-12 |
| C-13 | `0024:C2` ("Proving the AUTHORED form is proving the executed form"), `0024:C1` (widening/narrowing clause) | The contract proves the model against *itself* at one instant. It cannot detect that a domain member is a real command, cannot detect narrowing, holds no history. The Problem Statement's motivating failure — "a model can route perfectly to a command that does not exist" — is *not closed by this RDR*. It is relocated: instead of an unchecked emit value, we get a checked emit value against an unchecked domain. The RDR concedes this in `§risks-and-mitigations` and hands it to "the consumer's cross-file seam test", which does not exist and is not scheduled. | Author declares `domain = ["/rdr-propse"]` (typo in the domain), authors the matching typo'd value in the rule, and `lint` exits 0 with the declaration confirming the mistake. The seed defect ships, now with a proof attached. | §3, premortem, AT-13 |
| C-14 | `0024:C1` ("Reusing `ConformValue` (A4) settles non-canonical literals in the permissive direction: `03`, `+5` and `-0` are admitted") | A grammar contract settles a semantic question by inheriting whatever `strconv.Atoi` happens to do, and then declares the inheritance to be "the point of the reuse". `int`-kind emit is not used by any model in the repo and is not used by the motivating consumer (A5 enumerates only enums and one `scalar`). An unused kind is carrying a documented wart into a locked grammar. | An author writes `count = "007"` and it is admitted; a second author writes the same value differently and byte-equality at the consumer fails on two "valid" ints. | §1, AT-14 |
| C-15 | `0024:§phase-1` ("Amend three stale comments in the same change"), `0010:C3` | `0010:C3` is a locked Final normative contract whose text reads "Emit keys are NOT tag keys: they are **undeclared**, uninterpreted, compared by exact byte equality, and MUST NOT be matched, guarded, written, or read by any accessor". 0024 reinterprets "undeclared" as "not a tag" so it can widen without an override, then edits three code comments and two doc sentences to match its reading. Meanwhile `internal/table/emit_0010_test.go::TestReq27` carries the same word in its own failure message ("an emit key is undeclared and is not a tag") and is *not* on the amendment list — the RDR asserts it "keeps passing unchanged", which is true of its behavior and false of its prose. | The next reader of `0010:C3` reads a contract that says emit keys are undeclared, and code comments that say they may be declared, with no override recorded anywhere linking them. | §2, premortem, AT-15 |
| C-16 | `0024:S3` ("The pass criterion is LOAD-side equivalence, not payload byte-identity") | Scenario 3 defines its own pass criterion as "the loader admits and refuses exactly the same models as before, and the only payload difference anywhere is the one appended field" — and then measures it *after* Phase 3 regenerates the goldens that would have detected any other payload difference. The oracle is destroyed by the phase whose output it is supposed to check. | A second, unintended payload change rides along in the same regeneration and no test in the repo can see it. | §1, premortem, AT-16 |
| C-17 | `0024:S7`, `0024:§phase-3` ("The §D1 registration obligation, whichever record lands second") | The Done criterion for this RDR includes a scenario (S7) whose assertion runs against an oracle that does not exist at HEAD and may never be built by this change. If 0024 lands first, S7 is discharged by writing a note in a peer RDR's implementation plan. A scenario that can be satisfied by a sentence is not a scenario. | `dispositions` ships unassigned to ECHO or PLAN; 0023's later reflective oracle either fails on it or is written around it. | §1, AT-17 |
| C-18 | `0024:§decision-rationale` (QOC matrix), `0024:ALT1` | The matrix scores "Cost" 3 for the chosen option and 5 for the rejected prefix convention, with equal weights, and the correctness row is scored on catching *both seed defects*. But the misspelled-value defect is only caught if the domain is right (C-13), and the typo'd-key defect is the one a reserved-prefix alternative was never trying to catch. The 22-vs-18 total is manufactured by scoring a strawman version of B against a best-case version of A. | Not user-visible. It matters because the record's whole justification for a 1846-line grammar over a one-check convention rests on a four-point spread in a hand-scored table. | §2 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The declaration grammar has arms that do not exist, and arms that fire under someone else's category

**Root cause in the RDR.** `0024:C1` enumerates the `malformed_emit_declaration`
refusal set by *deriving it from the type structure* rather than from the
decoder's observable behavior. `0024:A7` admits this in as many words — "the
arms were derived from the type structure rather than enumerated exhaustively"
— and then carries the closure claim as Pending with `S1` named as the gate.
The problem is not that the set is incomplete. It is that two of the arms
collapse into one, and one of the "arms" is not this RDR's category at all.

**The specific passage.** `0024:C1`, the refusal paragraph:

> Refused as `malformed_emit_declaration`: an unknown `kind` token; an `enum`
> with no domain, an empty domain, an empty-string member, or a duplicate
> member ... and — the arms that exist only because strictness cannot reach
> them — a `domain` that is neither a flat array of strings nor a table of
> disposition keys, a disposition whose value is not an array of strings, any
> nesting below the disposition level, a non-string member, **and a
> disposition sub-table carrying no keys.**

And, three paragraphs earlier:

> Strictness still catches a typo at the DECLARATION level (`domaim = [...]`
> refuses), but inside the domain sub-table the decoder catches nothing.

Both statements are wrong in ways I measured against
`pelletier/go-toml/v2 v2.2.4` with `DisallowUnknownFields()` set — the exact
configuration `internal/table/source.go::decodeStrict` uses.

Result 1, the collapsed arm:

```
STRICT domain-empty-table   err=<nil>  Domain: interface {}(nil)
STRICT domain-absent        err=<nil>  Domain: interface {}(nil)
```

`[emit.next.domain]` written as a sub-table header with no keys under it, and
`domain` omitted entirely, decode to the **same** value: `nil`. C1 lists these
as two distinct refusal arms ("an `enum` with no domain" and "a disposition
sub-table carrying no keys"), and `0024:S1` requires "a table-driven fixture
per arm". The implementer will write the fixture for the second arm, watch it
produce the first arm's error message, and then either (a) delete the arm and
silently narrow the contract, or (b) add source-text scanning to distinguish
them — new machinery nothing in this RDR licenses, for a distinction with no
consequence. Note that `domain = []` *is* distinguishable
(`Domain: []interface{}{}`), which is why this reads as an oversight rather
than a considered collapse: the RDR lists "an empty domain" as a third,
separate arm and that one is real.

Result 2, the misattributed arm:

```
STRICT decl-typo kimd       err=strict mode: fields in the document are
                                missing in the target struct
```

The declaration-level typo does refuse — as `CatUnknownSchemaField`, mapped by
`decodeStrict`'s `*toml.StrictMissingError` arm. C1 presents this as the
grammar's strictness working, in a paragraph whose subject is the emit
declaration and whose neighbors are all `malformed_emit_declaration` arms. An
implementer reading C1 top to bottom will reasonably write a test expecting
`malformed_emit_declaration` for `kimd = "enum"`, watch it fail, and go
looking for a bug in code that is behaving correctly.

**Symptom the user sees.** An author starts a declaration, writes the
sub-table header, and goes to fill in the dispositions:

```toml
[emit.next]
kind = "enum"
[emit.next.domain]
```

They get `emit_value_out_of_domain`? No — they get
`malformed_emit_declaration: an enum declares no domain`, pointing at a model
whose `domain` sub-table is visibly present on the screen in front of them.
And when they misspell `kind`, they get a message about "fields in the
document" with no mention of `emit`, no mention of `next`, and no mention of
what the correct spelling is.

### 1.2 Every emit refusal points at line 1

**Root cause in the RDR.** The record repeatedly promises a *locator* and
never once specifies where the line number comes from. `0024:F1` says the
finding carries "category slug in the finding's `code`, offending file in
`locator`" — note "offending file", not line, which is a hedge the rest of the
record does not honor. `0024:MVV` step 2 says the finding carries "its
category slug in `code` and a locator". `§consequences` describes the adoption
loop as "one refusal per run, in an order the tier leaves unspecified" without
ever asking what the author is looking at when they get that refusal.

**The specific passage.** `0024:F1`:

> `intrastate lint` exits nonzero with one blocking finding — the first
> refusal the fail-fast pipeline reaches (category slug in the finding's
> `code`, offending file in `locator`) ... Diagnose from the finding's
> category + locator

Grounding: `internal/table/category.go::atLine` is called at **exactly one
site in the entire loader** — `internal/table/load.go:214`, inside `loadTags`,
via `tagHeaderLine`. Nothing else stamps a line. `Failure.Line`'s own doc says
"Zero is the common case and it means UNKNOWN". `internal/cli/flow_input.go::
loadFindings` renders zero as `path + ":1"`, documenting `:1` as "the
documented fallback for that zero, not a claim about the defect".

There is no `emitHeaderLine`. There is no phase that adds one. `0024:§phase-1`
lists the work as "Extend the source schema with the `[emit]` table, load
declarations beside tags, and refuse ... in the load pipeline" — no locator.
The tag grammar got `tagHeaderLine` because RDR 0002 needed it. This RDR
mirrors the tag grammar's *declare-then-prove move* and does not mirror the
one piece of it that makes the refusal actionable.

**Symptom the user sees.** The exact adoption journey `§consequences` predicts
— "on a model the size of the motivating consumer's (54 `[rule.emit]` blocks)"
— now reads: run `lint`, get one finding, `models/rdr-status.toml:1`, go find
which of 54 emit blocks it means by reading the `Detail` string and grepping.
Fix. Rerun. `models/rdr-status.toml:1`. Repeat forty times. This is the
combination of C-3 and C-10 and it is the single most likely reason the
feature is never adopted by the consumer it was designed for.

### 1.3 `ConformValue` reuse is safe only under a precondition the RDR never states

**Root cause in the RDR.** `0024:A4` is stamped Verified on a claim about
*importability* — "implementable in the load pipeline without importing
`internal/guard`" — and then the Evidence quietly expands into a much stronger
claim the assumption's headline does not make: that `ConformValue` is "already
exported with exactly the needed signature" and "`conformDomain`'s enum arm is
exactly the membership test C2 needs". `0024:C3` builds on that with a
concrete call shape:

> the check constructs a throwaway `TagDecl{Kind: d.Kind, Domain: d.Domain}`
> per call and passes it; `Min`/`Max`/`Elements` stay nil

**The specific passage.** `0024:A4` Evidence:

> `conformDomain`'s enum arm is exactly the membership test C2 needs
> (`slices.Contains(decl.Domain, member)`)

The actual source, `internal/table/load.go:849-872`:

```go
func conformDomain(decl TagDecl, member string) error {
	switch decl.Kind {
	case "enum":
		if len(decl.Domain) > 0 && !slices.Contains(decl.Domain, member) {
```

The membership test is **guarded by `len(decl.Domain) > 0`**. An `enum`
declaration whose `Domain` is empty conforms every value. And
`conformKind` (`load.go:833-845`) switches only on `"int"` and `"bool"` — it
has **no `enum` arm and no `scalar` arm and no default**. An unrecognized kind
token conforms everything too.

Both behaviors are correct for tags, where a zero `TagDecl` "conforms
everything" is the documented property that keeps an undeclared `--tag` key
shape-only. Both are wrong as the *sole* check for emit, where an empty domain
under `kind = "enum"` must refuse. The RDR knows this — C1 lists "an `enum`
with no domain, an empty domain" as refusal arms — but nowhere does it say
that the `ConformValue` reuse **depends** on those C1 arms having already
fired. It presents the reuse as complete and the C1 arms as grammar hygiene.

An implementer who writes `loadEmitDecls` to record the declaration and
`checkRuleEmit` to call `ConformValue`, and who defers or loosens one of C1's
arms in the process (say, because the empty-sub-table arm collapsed and they
relaxed the whole domain-shape family — see 1.1), gets a loader that admits
everything under a declaration that looks strict.

**Symptom the user sees.** Model declares:

```toml
[emit.next]
kind = "enum"
[emit.next.domain]
route = []
stop  = []
```

Loads clean. Every emit value on every rule passes. `intrastate lint` exits 0.
The declaration table reads, to a reviewer doing exactly what `0002`'s bargain
promises — "reviewers read the declaration, not 500 rows" — as a fully closed
vocabulary. It checks nothing.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`0024:§phase-3` — Envelope surfacing, specifically the "licensed diff"
paragraph — together with `0024:S3`.**

This is not a prediction about taste. It is a prediction about a contradiction
that becomes visible the first time someone runs `go test ./internal/cli/`
after the change.

The RDR spends its longest single evidence block (`0024:A2`) establishing that
the diff is 28 sites, that all 28 pin the `emit`/`next` adjacency, and that
their "mechanical regeneration is authorized HERE by name". It draws a sharp
line: `TestReq39` is "an amendment to a peer RDR's spec test and is called out
rather than absorbed as a fixture edit", while the 28 goldens are "0011's own
adversarial oracle for its demand-set change ... not a payload-shape contract
0011 owns, so their update is a fixture edit rather than a peer-spec
amendment."

That line does not survive reading the file. `internal/cli/
flow_demand_0011_test.go:622-631`:

> `resolveGolden` and `resolveFixtureGolden` are the FULL `flow resolve`
> payloads captured from the **PRE-change build of this tree**, per A-11.
>
> S8 requires byte-identity "before and after (A15 iii)", and names the gap
> the spike could not close ... These goldens close it — they assert the WHOLE
> payload, so a field this contract was never meant to touch cannot change
> unnoticed.

And the failure message the sweep prints (`:562-566`):

> "No shipped flow-reachable fixture is in the changed class (A15 iii), so
> every one of these payloads is **byte-identical before and after**"

These are not fixtures. They are `0011`'s discharge of REQ-118 — the frozen
pre-change reference that its byte-identity guarantee is *measured against*.
Regenerating them does not update a fixture; it replaces the reference with
the post-change value, which makes REQ-118 trivially true forever after and
destroys the only artifact in the repo that could have caught an unintended
payload change.

`0024:S3` then completes the circle. Its pass criterion is:

> the loader admits and refuses exactly the same models as before, and **the
> only payload difference anywhere is the one appended field**

and, two sentences later:

> "Full suite green" is measured AFTER Phase 3's licensed golden regeneration

The oracle that would establish the criterion is destroyed by the phase whose
output the criterion checks. There is no remaining test in the repo that can
distinguish "only `dispositions` was added" from "`dispositions` was added and
something else moved".

Six weeks after shipping, the first person who needs to change the resolve
payload again — 0023's `--plan-only` projection is already queued against the
same anchor (`0024:JC1`) — will discover that `resolveGolden` no longer means
what its doc comment says it means, and will rewrite Phase 3's licensing
model: either by preserving a genuine pre-change capture, or by demoting
REQ-118 explicitly. The RDR's own framing makes this rewrite inevitable,
because it authorizes the destruction of the evidence while asserting the
guarantee is untouched.

Two adjacent passages will go with it. `0024:A8`, which names `TestReq146` as
the determinism proof, is refuted by that test's source: its field sweep
iterates a **hand-written list of `table.Row` field names**
(`roundtrip_test.go:1346`) and its determinism leg compares `got.Rows`
(`:1337`). `Model.EmitDecls` is neither. Phase 2's stated failure signal — "an
unsorted union ... fails `TestReq146`'s sorted-sequence sweep" — will not
fire, and the paragraph asserting it will be corrected. And `0024:C4`'s
never-omitted argument, which pays for the entire 28-golden diff on the
grounds that "a consumer must be able to read `dispositions == {}` as 'this
model declared nothing' rather than 'this binary predates the field'", is
self-defeating: an older binary emits no `dispositions` key at all, so the
absence case the field is never-omitted *for* is precisely the case a
never-omitted field cannot signal. The `omitempty` decision will be revisited
the first time someone counts what the diff bought.

---

## 3. The one assumption that will not survive first contact with a real user

**`0024:A5`.**

> The motivating consumer's emit answers are enumerable at authoring time —
> routing commands are fixed strings and judgment cells answer with fixed
> `stopped:<code>` tokens — so an enum domain with a disposition partition
> can actually be authored for each of its routing keys.

Status: Verified. Method: Peer RDR. And the Evidence is the strongest passage
in the record — it does real counting: 54 `[rule.emit]` blocks in
`rdr-status.toml`, 22 distinct `next` values, 29 blocks in `rdr-write.toml`
routing on `op` with 10 distinct values, and the genuinely valuable discovery
that the partition is three-valued, not binary, which is what kills `ALT1`.

The assumption is not wrong about enumerability. It is wrong about what
"can actually be authored" means, and it fails on three independent legs.

**Leg 1: the models are not in this repo.** A5's Evidence is measured against
`rdr-status.toml` and `rdr-write.toml`, which live in the sibling RDR engine
repository. `0024:§phase-4` concedes this directly — "A5's motivating models
(`rdr-status.toml`, `rdr-write.toml`) live in the sibling engine repo and are
not reachable from a phase here" — and responds by inventing a synthetic
second example: "a small routing table whose one emit key partitions into
`route` and `stop` members". So the only in-repo demonstration of the
disposition partition, which is *the entire folded facet and the sole reason
`dispositions` exists on the wire*, is a model written by the implementer to
exercise the grammar the implementer just wrote. There is no adversarial
authoring pressure on the grammar anywhere in this change.

**Leg 2: `rdr-write.toml`'s `edit` key already breaks it.** A5's own Evidence
says: "`rdr-write.toml`'s `edit` key holds a multi-line shell script — an
unambiguously open-valued key that takes `kind = "scalar"`". A5 files this as
a *win* — "the per-key escape C1 provides and a live instance of why it is
needed". Read it the other way: on the very first real model, one of the two
routing models the RDR measured, a key must take the escape hatch. Under C2's
whole-model strictness, the author of `rdr-write.toml` must declare *every*
emit key at once. They will have keys they are confident about, keys they are
not, and at least one they cannot close. The path of least resistance is
`kind = "scalar"` on the uncertain ones — which is exactly the hollowing
`§risks-and-mitigations` names and then declines to detect, because
`0006:C17` closes the advisory tier (C-11).

**Leg 3: the adoption loop is longer than any author will tolerate.** Combine
what the record establishes about itself:

- whole-model strictness: every emit key must be declared at once (C1);
- fail-fast: one refusal per run (C2, measured in
  `evidence/spikes/c2-finding-multiplicity.md`);
- unspecified order: the author cannot predict which defect comes next (C2);
- `file:1` locators for every one of them (C-3, grounded above);
- 54 `[rule.emit]` blocks on the first model.

`§consequences` writes this down honestly — "an N-round fix-and-rerun loop ...
where N is the number of defects the declaration exposes, not the number of
runs an author expects" — and then declines to do anything about it, on the
grounds that fail-fast is "a cost of the tier this RDR lands in rather than of
the declaration itself".

That is where the assumption breaks. Every other load category in
`internal/table/category.go` fires on a *single authoring mistake* the author
just made. `unknown_emit_key` and `emit_value_out_of_domain` are the first
categories in this codebase that fire *N times on a model the author did not
change*, because the author switched on a check retroactively over an existing
corpus of values. Fail-fast is a reasonable tier contract for the former and a
hostile one for the latter, and the RDR's argument that it inherits the tier's
cardinality "rather than deciding it here" is precisely the move that hides
the distinction.

The first real user will declare `scalar` for everything, get a green lint,
and get none of the proof. The RDR predicts this — "An author facing it
declares `scalar` first and tightens per key, which is the escape hatch's
second purpose" — and offers no mechanism, no report, no finding, and no
metric that would ever tell anyone whether the second half of that sentence
happened.

---

## 4. Premortem

*Written from twelve months after the merge.*

`dispositions` shipped. It is in every `flow resolve` payload in the repo, as
`{}`. It has never held a value in production.

Here is how that happened.

**Week 1 — Phase 1 lands, one arm short.** The implementer worked C1's refusal
list top to bottom. `loadEmitDecls` gained hand-written arms for the unknown
kind token, the empty domain, the empty-string member, the duplicate member,
the domain-on-non-enum, and the empty disposition token. Then they hit "a
disposition sub-table carrying no keys". They wrote the fixture
`neg-emit-empty-disposition-table.toml`, ran it, and got the message from the
*previous* arm: an `enum` declaring no domain. They probed
`pelletier/go-toml/v2` and found what any probe finds — `[emit.next.domain]`
with no keys under it and `domain` absent entirely both arrive as `nil`. There
is no observable difference. They deleted the arm, wrote a comment saying the
two cases are indistinguishable at the decoder, and moved on. `0024:A7`'s
closure claim, the one thing `§prerequisites` said Phase 1 could not close
without, was discharged by narrowing the contract instead of proving it.

**Week 2 — the locator.** `checkRuleEmit` went in as C2 specified: a step
after `loadTags`, reading `sourceRule.Emit`, calling `ConformValue` with a
throwaway `TagDecl{Kind: d.Kind, Domain: d.Domain}` exactly as `0024:C3`
prescribes. It worked. Every refusal it produced carried `Failure.Line == 0`,
because `atLine` is called from exactly one place in the loader —
`internal/table/load.go:214`, inside `loadTags`, via `tagHeaderLine` — and
nobody wrote `emitHeaderLine`. `internal/cli/flow_input.go::loadFindings`
rendered them all as `<model>:1`. The MVV's two-rule table has two rules, so
`:1` was not obviously wrong on a fixture with four lines of emit, and it
passed review.

**Week 3 — Phase 3 and the goldens.** `resolvePayload` gained `dispositions`
after `Emit`, taking the struct to fifteen fields.
`TestReq39_EmitSitsImmediatelyAfterGatesOnTheWire` went red on both legs and
was edited exactly as `0024:§phase-3` (iii) licenses: thirteen keys to
fourteen, `NumField() != 14` to `!= 15`, and the prose strings updated. Then
all 28 payload literals in `internal/cli/flow_demand_0011_test.go` went red,
which the RDR had predicted to the number. The implementer regenerated them
from the new build and committed.

That was the moment the safety net came off, and nobody noticed, because the
RDR had authorized it by name. `resolveGolden`'s doc comment still says the
payloads are "captured from the PRE-change build of this tree, per A-11".
`shippedResolveGoldens`'s sweep still prints "every one of these payloads is
byte-identical before and after". Both statements became false in the same
commit that made every test pass. `0024:S3`'s pass criterion — "the only
payload difference anywhere is the one appended field" — was declared met by
running a suite whose only whole-payload oracle had just been rewritten to the
post-change value.

**Week 5 — Phase 2's silent failure.** `Model.EmitDecls map[string]EmitDecl`
went on the model. `EmitDecl.Domain` was built by iterating the domain
sub-table's map and appending. The sort was in the plan (`0024:§phase-2`: "The
sort is load-bearing, not cosmetic"), but the implementer, having already
narrowed one contract that week, checked the stated failure signal first: "an
unsorted union built from the domain sub-table's map iteration is
non-deterministic across loads and fails `TestReq146`'s sorted-sequence
sweep." They removed the sort to watch it go red. It stayed green.

`TestReq146_EveryEmittedSequenceIsASortedSlice` iterates a hand-written list
of field names — `"Atoms", "NextTags", "Writes", "RequiresOwned", "Gate",
"Escape", "Suffix"` — over `reflect.TypeOf(table.Row{})`. `EmitDecls` is on
`Model`, not `Row`. Its determinism subtest compares `got.Rows`, not the
model. `0024:A8` had flagged the scope as unverified; `0024:C3` and Phase 2
had both leaned on it as though it were settled. The sort went back in because
the implementer put it back on principle, and `0024:S4`'s repeated-load leg
was written and passed. Nobody wrote the equivalent for the second developer,
who added a `Dispositions` map surface eight months later and had no sweep
covering it either.

**Month 2 — the consumer tries to adopt.** The RDR engine team took
`rdr-status.toml` — 54 `[rule.emit]` blocks, 22 distinct `next` values, the
model `0024:A5` was measured against — and added:

```toml
[emit.next]
kind = "enum"
[emit.next.domain]
route = ["/rdr-propose", "/rdr-refine", "/rdr-resolve"]
stop  = ["stopped:no-profile"]
```

`intrastate lint`. One finding. `unknown_emit_key`. `rdr-status.toml:1`.

They grepped for the key, found it, added it. Rerun. One finding.
`emit_value_out_of_domain`. `rdr-status.toml:1`. Fix. Rerun. One finding.
`rdr-status.toml:1`. They did this eleven times before writing the loop into a
shell script and nineteen more after. `§consequences` had told them exactly
this would happen — "an N-round fix-and-rerun loop — one refusal per run, in
an order the tier leaves unspecified" — and had told them the remedy in the
next sentence: "An author facing it declares `scalar` first and tightens per
key, which is the escape hatch's second purpose."

So they did. `rdr-write.toml`'s `edit` key already had to be `scalar` — the
RDR's own A5 Evidence says so, a multi-line shell script that cannot be
enumerated. Once one key in a model is `scalar`, the marginal cost of the
second is zero and the marginal benefit of the first non-`scalar` key is one
more fix-and-rerun loop. Both models shipped fully declared and entirely
`scalar`. `intrastate lint` exits 0. `Categories()` carries three new members
that no shipped model has ever tripped. `dispositions` is `{}` on every
payload, because `scalar` keys carry no dispositions
(`0024:§mini-checks`, disposition table, row 5: "payload, no disposition
entry").

**Month 4 — the seed defect fires anyway.** A rare row in `rdr-status.toml`
answered `/rdr-prelok` — one letter short. It is a `scalar` key, so load
admitted it. `flow resolve` returned it. The consumer skill executed it
verbatim, as `0024:§problem-statement` says it does. It failed at a close-out.

The postmortem asked why the declared-vocabulary work had not caught it. The
answer is in `0024:§failure-modes`: "a `scalar`-declared key admits any value
— declared-but-unvalidated is the documented escape hatch, visible in the
model source." And in `§risks-and-mitigations`: an advisory finding for
`scalar`-declared keys "is deliberately NOT taken — the advisory tier is
closed (`0006:C17`) — so the check is review- and consumer-seam-side."

Nobody ever built the consumer-seam-side check. It is named four times in the
record — as the mitigation for domain drift, for disposition-token drift, for
the stale-but-valid member, and for scalar hollowing — and it is scheduled
nowhere, owned by nobody, and lives in a different repository. `0024:C13`'s
underlying problem, that the RDR proves the model against itself and has no
access to external truth, was correctly identified and correctly assigned to
someone who was never told.

**What the ledger says today.** `internal/table/category.go` carries three
categories with zero production hits. `resolvePayload` carries a fifteenth
field that has never been non-empty. `internal/cli/flow_demand_0011_test.go`
carries 28 goldens whose doc comment describes a pre-change build that no
longer exists. `internal/table/model.go::EmitValue`,
`normalize.go::emitSequence` and `dump.go::renderEmit` carry amended comments
saying emit keys may be declared, while `0010:C3` — Final, locked, never
amended — still reads "they are undeclared, uninterpreted, compared by exact
byte equality", and `internal/table/emit_0010_test.go::TestReq27`'s failure
message still tells the next reader "an emit key is undeclared and is not a
tag".

The record predicted almost every one of these outcomes in its own
`§consequences`, `§risks-and-mitigations` and `§failure-modes`. It named the
adoption cost, the hollowing, the drift, the unreported staleness, and the
merge-day-ships-capability-not-coverage problem. Then it accepted all of them,
one at a time, each with a locally reasonable argument, and shipped.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time tests: each is executable against the repo and the
decoder **before** any implementation exists, and each fails today.

### AT-1 — the empty disposition sub-table is a distinguishable state (C-1)

```gherkin
Given the repo's decoder, pelletier/go-toml/v2 v2.2.4, with DisallowUnknownFields
And a struct { Kind string `toml:"kind"`; Domain any `toml:"domain"` }
When "[emit.next]\nkind=\"enum\"\n[emit.next.domain]\n" is decoded
And  "[emit.next]\nkind=\"enum\"\n"                     is decoded
Then the two decoded Domain values MUST differ
```

**Runs today: FAILS.** Both yield `Domain: interface{}(nil)`. C1 must either
drop the arm, or state that it is merged into the "enum with no domain" arm,
or license the source-text scan that would separate them.

### AT-2 — every C1 refusal arm carries `malformed_emit_declaration` (C-2)

```gherkin
Given C1's enumerated refusal arms
When each arm's minimal TOML is decoded with DisallowUnknownFields
Then no arm may be caught by the decoder under a DIFFERENT category
And  C1 must name the category for any arm the decoder catches first
```

**Runs today: FAILS.** `kimd = "enum"` refuses `unknown_schema_field`. C1
presents it inside the emit-grammar strictness paragraph without naming the
category, and `S1` would have an implementer assert the wrong one.

### AT-3 — an emit refusal grounds a source line (C-3)

```gherkin
Given the load pipeline's line-attribution mechanism
When the RDR specifies a refusal whose finding a user diagnoses from `locator`
Then the RDR must name the function that stamps Failure.Line for that refusal
```

**Runs today: FAILS.** `rg 'atLine\(' internal/table/` returns exactly one call
site: `load.go:214`, `loadTags`. The RDR names no emit analogue and schedules
none. Every emit finding will render as `file:1`. `0024:F1` and `0024:MVV`
step 2 both promise a locator; neither says where it comes from.

### AT-4 — `TestReq146` reaches a model-level carrier (C-4)

```gherkin
Given TestReq146_EveryEmittedSequenceIsASortedSlice at internal/table/roundtrip_test.go:1333
When a new []string is added to table.Model (not table.Row)
Then the sweep must fail if that slice is unsorted
```

**Runs today: FAILS.** The field sweep iterates a literal list of `table.Row`
field names (`:1346`); the determinism leg compares `got.Rows` (`:1337`).
`Model.EmitDecls[k].Domain` is unreachable from both. `0024:A8` guessed this
was possible and marked it Pending; the answer is a five-line read, and `C3`
plus Phase 2 both assert the opposite as settled.

### AT-5 — the emit category witness trips its category and no other (C-5)

```gherkin
Given internal/table/dump_test.go::TestReq106's witness convention
  (REQ-119 and REQ-131: "one mutated fixture per category, tripping it and no other")
When the RDR proposes three new categories under a fail-fast loader
  whose "order in which independent defects are checked is deliberately unspecified"
Then the RDR must state how each witness's category is stable under a reordering
```

**Runs today: FAILS.** `0024:S6` asserts registration in `Categories()` and the
existence of a `testdata/neg/*.toml` witness. It never asserts the "and no
other" half, which is the half the convention exists for, and `0024:C2`
explicitly refuses to fix any precedence among the categories.

### AT-6 — regenerating the 28 goldens does not void a peer guarantee (C-6)

```gherkin
Given internal/cli/flow_demand_0011_test.go's resolveGolden and shippedResolveGoldens
When the RDR classifies an edit to them as "a fixture edit rather than a peer-spec amendment"
Then their own file-level documentation must not bind them to a peer REQ
```

**Runs today: FAILS.** `:622-631` binds them to `0011`'s REQ-118/S8 as the
capture "from the PRE-change build of this tree, per A-11", closing a gap the
0011 spike "could not close". `:562-566` asserts byte-identity "before and
after". `0024:A2` applies its peer-spec test to `TestReq39` and not to these,
on a distinction the file itself does not make.

### AT-7 — the step count between `loadTags` and `normalizeRules` (C-7)

```gherkin
Given internal/table/load.go::run's step slice
When the RDR states how many steps lie between loadTags and normalizeRules
Then that number must equal the source
```

**Runs today: FAILS.** `0024:C2` says "The six steps between `loadTags` and
`normalizeRules`" and then names five. `load.go:85-89` has five:
`loadAccessors`, `loadDump`, `loadContexts`, `loadInitial`, `loadTerminal`.

### AT-8 — `ConformValue` refuses an empty enum domain (C-8)

```gherkin
Given internal/table/load.go::ConformValue(TagDecl{Kind:"enum", Domain: nil}, "anything")
Then it must return an error
And  Given ConformValue(TagDecl{Kind:"enums"}, "anything")
Then it must return an error
```

**Runs today: FAILS both.** `conformDomain`'s enum arm is guarded by
`len(decl.Domain) > 0` (`load.go:852`); `conformKind` switches only on `"int"`
and `"bool"` with no default (`load.go:833-845`). Both return nil.
`0024:A4` calls the enum arm "exactly the membership test C2 needs" and
`0024:C3` builds the call site on that. The reuse is sound only if C1's arms
have already refused the empty-domain and unknown-kind states — a precondition
neither A4 nor C3 states.

### AT-9 — the disposition partition has an in-repo worked example (C-9)

```gherkin
Given the folded facet is the sole reason `dispositions` exists on the wire
When the RDR ships
Then at least one model in this repo, authored before this change,
     must exercise the [emit.<key>.domain] sub-table form
```

**Runs today: FAILS.** `rg --glob '*.toml' '^\s*\[emit'` returns zero (A3).
`models/examples/pricing-decision-table.toml` emits `plan` ∈ {basic, pro} and
`dpa` ∈ {required, none} — no route/stop split, confirmed by reading the file.
`0024:§phase-4` concedes it and adds a synthetic example written by the same
change. A5's real models are in another repository.

### AT-10 — first-adoption cost on the motivating model (C-10)

```gherkin
Given a model with 54 [rule.emit] blocks and 22 distinct values on one key
And  whole-model strictness (C2) and fail-fast, one-finding-per-run (C2)
When an author adds their first [emit.<key>] declaration
Then the RDR must state the expected number of lint runs to green
And  must state what the author is looking at in each one
```

**Runs today: FAILS the second leg.** `§consequences` states the loop and
names N honestly. It never asks what the locator says, and the answer (AT-3)
is `file:1` every time. The two findings compound; the record treats them
independently and neither section cross-references the other.

### AT-11 — a fully-`scalar` model is distinguishable from a proven one (C-11)

```gherkin
Given a model declaring every emit key with kind = "scalar"
When `intrastate lint` runs
Then its output must differ in some observable way from a model with
     fully enumerated domains
```

**Runs today: FAILS by construction.** Both exit 0 with no findings. The
mitigation `§risks-and-mitigations` offers is "visible in the one reviewable
declaration table" — a human reading the source, not an output. Meyer's cited
condition (flag the loopholes) is declined on the procedural ground that
`0006:C17` closes the advisory tier. The RDR quotes the condition and then
declines it in the same paragraph.

### AT-12 — never-omitted `dispositions` buys an observable distinction (C-12)

```gherkin
Given C4's reason for declining omitempty:
  "a consumer must be able to read `dispositions == {}` as 'this model declared
   nothing' rather than 'this binary predates the field'"
When a consumer receives a payload from a binary that predates the field
Then that payload must be distinguishable from `dispositions: {}`
```

**Passes trivially and therefore refutes the argument.** A pre-field binary
emits no `dispositions` key at all — JSON `undefined`, not `{}`. The
distinction C4 pays 28 goldens for is already free. The `omitempty` decision
needs a different justification or a different answer.

### AT-13 — the seed defect is actually closed (C-13)

```gherkin
Given the seed's failure: "a model can route perfectly to a command that does not exist"
When an author declares domain = ["/rdr-propse"] (typo in the DOMAIN)
And  a rule emits next = "/rdr-propse"
Then `intrastate lint` must not exit 0
```

**Runs today: would FAIL after implementation.** The model is internally
consistent, so every check in C1 and C2 passes. `0024:§failure-modes` names
this ("a declaration that faithfully copies a typo proves consistency with the
mistake") and routes it to a consumer seam test that does not exist. The
Problem Statement's headline failure is relocated, not closed, and the RDR's
QOC correctness row scores the chosen option 5 for "catches misspelled value
AND typo'd key" on the strength of catching the *value against a domain*
rather than the value against reality.

### AT-14 — `int` kind is used by something (C-14)

```gherkin
Given C1 admits kind = "int" and inherits strconv.Atoi's permissiveness
When the repo and the motivating consumer's models are searched
Then at least one emit key must plausibly want an int domain
```

**Runs today: FAILS.** No shipped model emits an int-shaped value. A5's census
of both motivating models finds enums and one `scalar`. C1 locks `03`, `+5`
and `-0` as admitted into a grammar for a kind nothing uses, and calls the
inheritance "the point of the reuse".

### AT-15 — `0010:C3`'s "undeclared" is not being reinterpreted (C-15)

```gherkin
Given 0010:C3 (Final, locked) reads "Emit keys are NOT tag keys: they are
  undeclared, uninterpreted, compared by exact byte equality, and MUST NOT be
  matched, guarded, written, or read by any accessor"
When 0024 makes emit keys declarable
Then the RDR must either record an Overrides entry against 0010:C3
     or show that "undeclared" there means only "not a tag"
And  every in-repo restatement of that word must be enumerated in the amendment set
```

**Runs today: FAILS the second leg.** `0024:§phase-1` amends three code
comments and `§phase-4` amends two doc sentences. It omits
`internal/table/emit_0010_test.go::TestReq27`, whose own failure message reads
"an emit key is undeclared and is not a tag" — and asserts the test "keeps
passing unchanged", true of behavior and false of prose. The record's
reinterpretation of `0010:C3` rests on a comma.

### AT-16 — S3's oracle survives Phase 3 (C-16)

```gherkin
Given S3's pass criterion "the only payload difference anywhere is the one appended field"
When "full suite green is measured AFTER Phase 3's licensed golden regeneration"
Then some test must remain that can detect a SECOND payload difference
```

**Runs today: FAILS.** After regeneration, nothing in `internal/cli/` holds a
pre-change whole-payload value. `TestReq39` asserts key order and field count,
which a same-position value change would not disturb. S3's criterion is
unfalsifiable at the moment it is measured.

### AT-17 — S7 is a scenario, not a note (C-17)

```gherkin
Given S7 asserts `dispositions` is registered under JDR 0002 §D1's ECHO/PLAN oracle
And  that oracle does not exist at HEAD (0023 is Final and unimplemented)
When 0024 lands first
Then S7 must name a runnable assertion in THIS change
```

**Runs today: FAILS.** `0024:S7`'s own text says "if this RDR lands first, the
obligation transfers to 0023's implementation and is recorded here" — the
scenario is discharged by a paragraph in a peer's plan. `0024:S7` then claims
"this scenario is what makes the Done criterion unsatisfiable while it is
unmet", which is false: the criterion is satisfied by the transfer.
