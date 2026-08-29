Model: claude-opus-5

# Hostile critique — cli/0023 Resolve-envelope projection opt-out

## Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0023:§source-authority-census` ("All fifteen are written at one assembly site") | The census claims fifteen fields; `resolvePayload` declares **fourteen** (`internal/cli/flow_resolve.go:31-55`), and the table has thirteen rows because `writes`/`clear` are merged into one. The one artifact that is supposed to be the exhaustive field inventory — the input to C2's partition and S3's reflective oracle — miscounts the set it inventories. | The implementer builds the declaration table from the census, gets a count mismatch against `reflect.TypeOf(resolvePayload{}).NumField()`, and either pads an entry or loosens S3's "entry set == field set" equality to an inclusion. S3 then stops being the check C2 says it is. | §1, premortem, AT-1 |
| C-2 | `0023:§phase-1-projection-on-the-verb` ("`decision_table_0010_test.go` pins `NumField() == 14` (:435) … a field addition is out of scope for this RDR by construction") + `0023:A4` + `0023:JC1` | A4 declares 0023 and 0024 compose "in either landing order with neither contract moving." But `0024:C4` **appends `dispositions` to `resolvePayload`** (0024 §Implementation, "Append `dispositions` to `resolvePayload` after `emit`"). Two shipped 0010 oracles break on that regardless of order: `decision_table_0010_test.go:435` (`NumField() != 14`) and — worse — `:416` `TestReq39_EmitSitsImmediatelyAfterGatesOnTheWire`, which compares the wire key order against a 13-element literal with `slices.Equal`, an **exact-set** assertion no appended field can survive. 0023 must additionally move its S2 key literal, its C2 declaration table, and the always-keep coincidence argument. A4 verified only that the *derivation inputs* don't cross the partition; it never checked that the *oracles* are order-independent. | Whichever RDR lands second turns the other's suite red on the first commit, including a predecessor's exact-equality oracle that `0023:A3` promised would not move. The implementer discovers mid-flight that "either order" was only ever true of the join rule, not of the test battery, and re-negotiates landing order under time pressure. | §1, §3, premortem, AT-2 |
| C-3 | `0023:C1` ("the projected encoding STRICTLY SHORTER than the default … asserted in BOTH output modes") + `0023:A8` (Status: **Pending**) | C1 states an unconditional MUST in both modes; A8 — the assumption that the MUST is *satisfiable* in text mode — is unverified with "not yet measured." The record locks a normative clause on an unmeasured premise and books the measurement into implementation Phase 2, after the contract is already binding. The RDR's own escape hatch ("the clause narrows to JSON") is a contract change, i.e. an unlock. | A Phase 2 oracle fails on a shape the author never measured, and the team must either amend a locked normative clause (which this project forbids) or weaken the oracle to match, which is the "satisfiable by an oracle that measures the wrong thing" failure C1 explicitly warns against three lines earlier. | §1, §2, premortem, AT-3 |
| C-4 | `0023:S5` ("projected text lines a byte-identical, stable **subset** … set membership, not a subsequence") | S5 asserts *set* subset. `respond/text.go::flatten` emits path-prefixed lines and `sort.Strings` them, so the projected line list is a set subset trivially. But the width half of the clause is asserted on "total rendered byte count," while the subset half is asserted on set membership — two different measurands on the same scenario, with no assertion tying them. A projection that drops the right *lines* but the text renderer changes a *value* on a carried line passes the set-subset half only if the mutated line coincidentally equals a default line; otherwise S5 goes red for the right reason. But the reverse — a projection that keeps every default line and additionally shortens one value — is not caught by either half. | A carried plan-group value silently changes in text mode only. The user diffs `--as=text` against `--as=json` and finds the two modes disagree, which is exactly the `0005:C1` guarantee this RDR claims holds "by construction." | §1, premortem, AT-4 |
| C-5 | `0023:§approach` + `0023:C2` ("`revision` rides the PLAN side deliberately … presently constant-empty") | `revision` is hardcoded `func (flowRequest) revision() string { return "" }` (`internal/cli/flow_exec.go:49`). Every projected payload carries `"revision":""` — 14 dead bytes on the wire and, in text mode, a line rendering as `revision: ` (trailing space, empty value; verified against `flatten`'s `case string` arm). The RDR spends four paragraphs of C2 justifying a field that carries no information and, in the mode the RDR most wants to shrink, renders as a visually empty line. | The agent consumer parsing projected text output sees a blank-valued line it must special-case, and the JSON consumer sees a key that is always `""`. The "14 bytes against a silent provenance loss" trade is paid on every single call for a benefit that arrives on no call. | §1, §3, premortem, AT-5 |
| C-6 | `0023:C1` ("The IDENTICAL INVOKED-READER SET … the projected run's set is taken from the same in-process derivation the payload field is rendered from (`internal/cli/flow_exec.go::invokedReaders`)") | This makes S1's reader assertion a **white-box test of a pure function against itself**. `invokedReaders(req.model, outcome)` is called once in `runFlowResolve` (`flow_resolve.go:181`) and is a pure function of model + outcome — neither of which the flag touches. Calling it twice with the same arguments and asserting equality is a tautology; it cannot fail on any implementation of the projection. The RDR asserts the opposite ("Asserting it from the projected payload is not a weaker form of this check — it is unwritable"), which is true of the *payload* channel but does not make the *derivation* channel a real check. | The one oracle guarding "the flag never changes which readers ran" passes unconditionally. If a future implementer moves the flag read upstream of `runReaders`, S1 stays green while the flag now changes what the run observes. The user gets a plan computed against a different owned snapshot depending on a report-width flag. | §1, §4, premortem, AT-6 |
| C-7 | `0023:§minimum-viable-validation` step 6 ("PASS BAR: … the 48-fact class saves at least 70%") + `0023:A1` | The 70% bar is calibrated on a **synthetic model authored for the spike** (`evidence/spikes/a1-byte-width.md` §S1b: 48 tags `f01..f48`, all identical stanzas, two rows discriminating on `f01` alone). No such model exists in the repo, in `models/examples/`, or in any consumer. The bar that gates lock is therefore a bar against a fixture the RDR itself invented to clear it. The three *real* fixtures measure 42.2%, 47.6%, and 51.3% — and the RDR explicitly exempts them from any percentage bar ("The state-machine shapes carry no percentage bar"). | The gating number is unfalsifiable in production. A consumer who adopts the flag on a real model sees ~45%, not 79.4%, and reads the Performance Expectations section as having overpromised by 34 points. | §1, §3, premortem, AT-7 |
| C-8 | `0023:§consequences` ("this RDR ships the CAPABILITY, not its adoption … a green build delivers zero measured saving on day one") | The RDR openly states that its entire justification — transcript-token cost — is not realized by anything it ships, and that adoption "is not scheduled by this RDR and is not one of its phases." The named first adopters are a kata (`rg0e`) and an unnamed consumer. There is no obligation on anyone to ever pass the flag. | The flag ships, five oracles go green, the RDR closes, and six months later the token cost is unchanged because no consumer migrated. The record's own §Consequences is the postmortem, written in advance. | §1, §3, §4, AT-8 |
| C-9 | `0023:A5` (Status: **Pending**) + `0023:§prerequisites` ("A5's implementability half is discharged by writing the total walker itself, which is Phase 2 work") | A Pending assumption whose verification is a build step is not a verified assumption; it is a bet. The RDR's own gate template (`0023:§assumption-verification`) says "no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it" — yet C1 asserts the whole-tree walker as a MUST with a fully specified scope (no name-skip, no `Hidden` gate, `completion` materialized and asserted present), and §disposition, §oracle-discriminability, S4, MVV step 5, and the desk trace all narrate it as settled. | If the walker cannot be written to C1's scope, C1's structural negative moves — which means the locked contract moves — after implementation has started. | §1, §2, premortem, AT-9 |
| C-10 | `0023:C1` ("--plan-only MUST NOT be accepted by any other command — satisfied by NON-REGISTRATION … pinned by a STRUCTURAL oracle that walks the WHOLE command tree") + `0023:§oracle-discriminability` (S4 row) | The whole-tree walker is specified with more normative text (three paragraphs of C1, an A5 with a two-stage verification plan, a discriminability row with two negative controls, an MVV step, an S4 scenario, and a desk-trace row) than the projection itself. It guards against a defect — someone registering `plan-only` on `completion` — that has never occurred, cannot occur accidentally (registration is verb-local and explicit), and whose blast radius is one unusable flag on a shell-completion command. Meanwhile the actual projection mechanism (A2's pointer-`omitempty` conversion, which changes the types of five fields on the single shipped payload struct) gets one deferred sentence in `0023:D-wire-byte-format`. | Implementation effort lands on the ceremony, not the change. The reviewer's attention budget is spent on `completion` while the type conversion — the one edit that can break default-mode byte identity — is unspecified. | §1, §2, premortem, AT-10 |
| C-11 | `0023:A3` ("swept all 21 test files invoking `flow resolve`") | The count is wrong: 22 test files invoke `flow resolve` (`internal/cli/*_test.go`, including `help_all_test.go`). A3 was already re-counted once at the cove sweep ("the earlier '16' was a miscount") and is still wrong. The substantive claim survives — I checked `flow_resolve_0005_test.go:480` (presence list, default mode) and `decision_table_0010_test.go:433` (field count, name-based) and both are default-mode — but a claim that has been miscounted twice is not a claim the reader should trust unaudited. | Nothing user-visible; this is a credibility defect that makes the *other* verified assumptions harder to trust at review time. | §1, AT-11 |
| C-12 | `0023:C2` ("Which of the two declaration carriers is used is an implementation choice left to Phase 1 alongside A2's mechanism; that it is separate from the projection code is not.") | C2 mandates that the ECHO/PLAN assignment be *declared* independently of the projection code, then defers the carrier (struct tag vs. name-keyed table) to implementation. But the two carriers have different failure modes: a per-field struct tag on `resolvePayload` is **not** independent of the projection if the projection is implemented as pointer-nilling on the same struct — the tag and the nilling live in the same declaration and move together in the same edit. Only the separate table achieves the independence C2 requires. The RDR leaves open the option that defeats its own requirement. | The implementer picks struct tags (the cheaper option, and the one A2's mechanism nudges toward), S3 becomes the tautological oracle C2 explicitly forbids, and a mis-assigned future field passes. | §1, premortem, AT-12 |
| C-13 | `0023:§disposition` (row 4: "`--plan-only` on a binary predating this RDR \| 2 \| `command-error` \| loud — an absent `observed` always means projection, never an old binary") | The claim is that version skew is loud. It is loud for a caller who *passes* the flag to an old binary. It is silent in the opposite direction, which the row does not consider: a caller that has been passing `--plan-only` and whose deployment is rolled back gets exit 2 on a command that previously succeeded — a hard failure in the middle of an agent chain, not a narrower payload. The Failure Modes bullet repeats the same one-directional reading. | An agent chain that adopted the flag breaks entirely on a rollback, with a `command-error` that names only the flag in pflag free text (the `0011:A14` limit the RDR accepts). The agent has no structured code to branch on and no path to degrade to the full width. | §1, premortem, AT-13 |
| C-14 | `0023:§problem-statement` ("The seed measured a 48-fact decision-table call at 1529 bytes … about 72%") vs. `0023:A1` (708→146 B) | The Problem Statement leads with a number the same paragraph then retires, and the record spends a two-note paragraph reconciling 1529 with 708. Two different measurements of "the motivating call" are carried simultaneously, one retired and one live, both quoted. §consequences quotes the retired 72% again ("The seed's ~72% is the echo *share*"). | A reader — human or the next implementer — cites the wrong figure. More seriously: nobody re-measured the *actual* motivating call after the partition was fixed. The 1529 B call was never re-run under the shipped partition; it was replaced by a synthetic. | §1, §3, premortem, AT-14 |
| C-15 | `0023:C1` ("Every carried field's rendering MUST be byte-identical to its default-mode rendering (same encoder, HTML escaping disabled), and the carried keys keep their default-mode relative order") + `0023:A2` | A2's verified mechanism converts `Model`, `Observed`, `Owned`, `Readers`, `Outcome` from concrete types to pointers. I reproduced it: default-mode bytes do survive, including `"owned":{}` and `"readers":[]`. But the RDR never enumerates nil-ability per writer, and the safety it relies on is **incidental, not contracted**. All five writers happen to be nil-safe today — `observedTagMap`/`tagMap` (`internal/cli/flow_exec.go:739-753`) and `readerIDs` (`internal/cli/flow_next.go:582-588`) all `make(...)` unconditionally; `Model`/`Outcome` are strings — but only `emitMap` carries a comment saying nil-avoidance is deliberate (`flow_resolve.go:279-284`), and `runFlowResolve` explicitly nil-normalizes `Gates` and `Next` (`:245-250`) while `Readers` gets no guard. Once these fields are pointers, a writer that ever returns nil yields a **non-nil pointer to a nil slice/map**, which `omitempty` keeps and the encoder renders as `"readers":null` — verified. That is a default-mode byte change on a path no oracle in S1–S5 covers. The RDR should have converted an incidental invariant into a stated one before making the wire depend on it. | Default-mode output silently changes for a request shape nobody tested, and the differential oracle cannot see it: S1 compares projected against default from the *same build*, so both sides drift together. | §1, §2, premortem, AT-15, AT-16 |
| C-16 | `0023:§minimum-viable-validation` step 1 ("today's full payload, byte-identical to the pre-change fixture") | This is the only assertion in the record that compares against a **pre-change** artifact, and it lives in the MVV — a manual run — not in any of the five oracles S1–S5. C1's "Default-mode output is byte-identical to today's" is the RDR's headline safety claim and it has no automated guard. Every oracle compares the new build to itself. | A default-mode regression ships. The user — every existing consumer, none of whom opted in — sees changed bytes on a call they did not modify, which is precisely the "no consumer changes" promise in §consequences. | §1, §2, premortem, AT-16 |
| C-17 | `0023:§phase-3-docs-and-help` ("a Phase 1 commit that registers the flag without regenerating them turns CI red on its own") | The RDR sequences flag registration in Phase 1 and doc regeneration in Phase 3, then notes in Phase 3 that this ordering breaks CI. It identifies the defect and does not fix the ordering. | The implementer's first Phase 1 commit fails `make check` on `docs-check` (`Makefile:34,103-109`), and the phases have to be re-ordered or merged at build time — which is a plan change, not an implementation detail. | §1, AT-17 |
| C-18 | `0023:C2` ("This RDR therefore does NOT mint a separate always-keep oracle (it would today assert exactly what S2 asserts); the obligation it creates is on the successor") | The always-keep core is the single safety invariant the whole design rests on — it is the reason Alternative 1 was rejected ("a caller who omits `escaped` reads a rescued plan as an ordinary one"). C2 then declines to assert it, calling the coverage "enforcement by coincidence of scope," and offloads the real oracle onto an unwritten successor RDR. A safety invariant enforced by coincidence is not enforced. | The moment a second projection mode is added — by anyone who did not read C2's penultimate paragraph — a rescued plan can be laundered into an ordinary one, and the S2 literal that "coincidentally" covered it no longer does. The user acts on a plan without knowing an escape row produced it. | §1, §2, §4, premortem, AT-18 |
| C-19 | `0023:§metadata` **Overrides** field (18 lines) + `0023:§problem-statement` (63 lines) + the whole record (1617 lines) | The record is 1617 lines to add one boolean flag that deletes five keys from one struct before one function call. The Overrides field alone is a paragraph. The Problem Statement carries a two-note methodological digression about bytes-versus-tokens. C1 is 121 lines. C2 is 111 lines. The Proportionality gate (`0023:G-proportionality`) is an unfilled template. | Nobody reads it whole. The load-bearing facts (C-1's field count, C-2's landing-order coupling, C-12's carrier choice) are buried in a volume that guarantees they are skimmed. Review defects that a 200-line record would surface are invisible here. | §1, §2, premortem, AT-19 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The declaration table and the payload struct disagree from the first commit

**Root cause in the RDR.** `0023:§source-authority-census` is the document's single inventory of the payload's fields and their partition sides. It opens: "All fifteen are written at one assembly site." The struct declares fourteen (`internal/cli/flow_resolve.go:31-55` — `Model`, `Revision`, `Observed`, `Owned`, `Readers`, `Outcome`, `Rule`, `Gates`, `Emit`, `Next`, `Writes`, `Clear`, `Escaped`, `EscapeClass`). The census table itself has thirteen rows because `writes` and `clear` share a row. So the RDR carries three different counts of the same set — fifteen in prose, thirteen in the table, fourteen in the code — and none of them is annotated as an approximation.

This is not cosmetic. C2 mandates a declared assignment source whose "entry set and the field set are equal (neither a field without an entry nor an entry without a field)," and S3 asserts that equality reflectively. The census is the only place in the record where the entries are enumerated. An implementer building the declaration from the census produces a table that fails S3 on its own first run.

**The specific passage that enabled it.** `0023:§source-authority-census`, first sentence: "Every success-payload field, its writer, and its C2 side. All fifteen are written at one assembly site (`internal/cli/flow_resolve.go` `resolvePayload{…}` literal) — there is no fallback arm and no second writer, which is what makes the partition a single-site edit." The single-writer half is true and I verified it (`rg` over `internal/cli` finds exactly one `resolvePayload{` literal, at `flow_resolve.go:223`, and no test constructs one). The count is wrong.

**Symptom the user will see.** None directly — this fails in the build. But the *resolution* is what harms the user. The implementer, facing an equality assertion that will not go green, has two moves: recount and fix the table (correct), or relax S3 from "equal in both directions" to "every struct field has an entry" (cheap, and the RDR's Testing Strategy §3 wording — "every field assigned to exactly one C2 group" — reads as licensing exactly that relaxation). Under the relaxed form, an entry for a field that no longer exists survives, and the projection built from that declaration references a phantom. The partition doctrine's whole promise — "an unassigned field is a test failure, never a silent default" (JDR 0002 §D1) — is then held up by an oracle that has been quietly halved.

### 1.2 The 0023/0024 landing order is not order-independent, and A4 verified the wrong thing

**Root cause in the RDR.** `0023:A4` is stamped **Verified / Peer RDR** on the claim that "whichever RDR lands second changes nothing in the other's contract." Its evidence is a re-read of `0024:C4`'s join rule, establishing that `dispositions` derives from `Plan.RuleID` and never from `observed`/`owned`. That establishes the *semantic* composition: `dispositions` is a PLAN field, and projecting the echo group strips none of its inputs. Correct, and I confirmed the assignment is also recorded in JDR 0002 §D1's Resolved clause.

But 0024's Implementation Plan says: "Append `dispositions` to `resolvePayload` after `emit`." That is a fifteenth struct field. And this RDR mints four artifacts that are *extensionally* pinned to the fourteen-field, eight-projected-key world:

- `0023:S2`'s explicit key literal: "exactly `revision`, `rule`, `gates`, `emit`, `next`, `writes`, `clear`, `escaped`."
- `0023:S3`'s declaration table (a fifteenth entry).
- `0023:C2`'s always-keep argument, which rests on "with ONE projection mode, always-keep and the S2 key-set literal have the same extension."
- `0023:§phase-1-projection-on-the-verb`'s explicit statement that "`decision_table_0010_test.go` pins `NumField() == 14` (:435) … Adding a field here would falsify A3's 'no predecessor oracle moves' claim, so a field addition is out of scope for this RDR by construction."

That last passage is the tell. This RDR knows a field addition breaks its assumptions and declares field additions out of scope — while a sibling RDR at depth 1, named in the same record's `0023:JC1` joint-check line, is planning exactly that addition.

I verified two predecessor oracles that break, not one:

- `decision_table_0010_test.go:435` reads `if rt.NumField() != 14`, with the message "C4 takes it from thirteen to fourteen." `dispositions` takes it to fifteen.
- `decision_table_0010_test.go:416` `TestReq39_EmitSitsImmediatelyAfterGatesOnTheWire` compares the emitted wire key order against a thirteen-element literal using `slices.Equal`. That is an **exact-set** assertion, not a relative-position one. No appended top-level key survives it, wherever it is appended.

The second is the more damaging, because `0023:A3` — stamped Verified — claims "the 0005/0010/0011 suites do not move." A3's sweep asked whether existing assertions run in default mode (they do; I confirmed `flow_resolve_0005_test.go:479-487` is a presence list and `flow_surface_0005_test.go:133` is presence-only). It never asked whether an *exact-equality* assertion exists on the payload's key set. One does, and the sibling RDR named in this record's own joint-check line will trip it.

**The specific passage that enabled it.** `0023:A4`'s claim wording: "whichever RDR lands second changes nothing in the other's contract (the premortem's seed-3 check: this cites the quoted join rule, not a summary of the sibling)." The parenthetical is defensive about the *right* thing (citing rather than summarizing) and blind to the *wrong* thing (that a join rule says nothing about test extensions). `0023:§prerequisites` then converts the unchecked half into a checkbox: "Ordering tolerance with cli/0024 confirmed (A4): either RDR may land first."

**Symptom the user will see.** Nothing, until CI. Then: whichever team lands second finds a red suite in a file they did not touch, owned by an RDR they did not write, and has to edit a normative key literal — which under this project's rules is a spec change, not a fix. The likely resolution under schedule pressure is to make S2's literal derived rather than authored, which is the "falling through an omit-list" failure `0023:C1` mints S2 specifically to prevent.

### 1.3 The reader-set oracle cannot fail, and the report-only guarantee is unguarded

**Root cause in the RDR.** `0023:C1` spends a full paragraph on how the invoked-reader assertion must be made, and reaches the wrong conclusion. Because `readers` is projected away, the projected run's payload cannot supply the set (correct — and the helper `flow_fixtures_0011_test.go::readersOf` does `t.Fatalf` on absence, which I verified at line 1774). C1's answer: take the projected run's set "from the same in-process derivation the payload field is rendered from (`internal/cli/flow_exec.go::invokedReaders`), never from its own output."

`invokedReaders(req.model, outcome)` is a pure function of model and outcome. `runFlowResolve` calls it exactly once (`flow_resolve.go:181`), and neither argument is derived from any flag this RDR adds. An oracle that computes `invokedReaders(m, o)` twice and asserts the two results are equal is asserting `f(x) == f(x)`. It is green on every possible implementation of the projection, including one that never runs readers at all.

C1 anticipates the objection and answers a different one: "Asserting it from the projected payload is not a weaker form of this check — it is unwritable, and the assertion it would replace is the one that catches this clause's own named breach." True about the payload channel. But the derivation channel is not a *weaker* check than the payload channel; it is *not a check*. The clause's named breach — "A later change that skips work whose only consumer is a projected-away field breaches this clause" — is precisely what this oracle is meant to catch and precisely what it cannot catch, because the skipped work would be the *call* to `runReaders`, and the oracle never observes whether that call happened.

**The specific passage that enabled it.** `0023:C1`, the paragraph beginning "The IDENTICAL INVOKED-READER SET cannot be read off the projected run's payload," and its endorsement in `0023:S1` Expected and `0023:§oracle-discriminability` S1 row ("the reader-set check reads a channel the projection cannot touch, so it fails on a real reader-set change rather than on the field's absence"). The discriminability row's negative control — "change the invoked-reader derivation under the flag and confirm the reader assertion goes red" — is unrunnable as written: there is no "under the flag" variant of a pure function.

**Symptom the user will see.** Deferred and severe. Today the projection sits after `runReaders`, so nothing breaks. But `0023:§phase-1-projection-on-the-verb` explicitly forbids a defensive guard and relies on placement alone ("a refusal returns before the projection is reachable … Do not add a 'if refusing, skip projection' guard"). The moment someone optimizes — "under `--plan-only` we don't emit `owned` or `readers`, so skip the reader pass" — the flag becomes decision-affecting. The kernel then resolves against an empty `Owned` set, `flow_resolve.go` line 180's `req.runReaders` never runs, and the user gets a *different rule selected* depending on a flag documented as report-only. S1 stays green. `0023:D-identity`'s claim — "the same invocation ± the flag selects the same rule … and an oracle asserts the ±-flag plan-group equality" — is the assertion that would catch it, and it is not among the five oracles as a separate scenario; it is folded into S1's "strict key-subset with byte-identical values," which compares two runs of the same build and would show both runs agreeing on the *wrong* rule only if the default run also skipped readers. It would not.

---

## 2. The one section rewritten within 6 weeks of shipping

**`0023:§normative-contracts` — specifically C1's width clause and its "BOTH output modes" scope.**

The clause reads: "the projected encoding STRICTLY SHORTER than the default … STRICTLY SHORTER is measured on the FULL EMITTED LINE … and it is asserted in BOTH output modes. In text mode the measurand is the total rendered byte count of the emitted lines."

It will be rewritten because it is a locked MUST resting on an unverified assumption, and the RDR says so: `0023:A8` is **Pending**, "not yet measured," with a verification plan scheduled for "implementation Phase 2" — i.e. after lock. `0023:A8`'s "If wrong" clause pre-writes the rewrite: "C1's both-modes width clause is unsatisfiable as written and the clause narrows to JSON."

I measured it. Simulating `respond/text.go::flatten` against the A1 spike's S1 default payload: 15 lines / 267 bytes default, 9 lines / 130 bytes projected. A8 holds on that shape. But the shape it holds on is the one the desk trace already walked. The shapes that stress it are the ones nobody measured:

- A request with no `--tag`s (if the model permits it), where `observed` renders `observed: (none)` — 17 bytes — and `owned` and `readers` likewise. The echo group's text contribution collapses to `model: <path>`, `observed: (none)`, `owned: (none)`, `readers: (none)`, `outcome: <o>`. Still positive, but the margin shrinks by an order of magnitude from the 137-byte saving on the 2×2 call.
- More sharply: C1 asserts the width clause *unconditionally*, with `0023:S1` justifying unconditionality as "the five echo keys always render, `model`/`outcome` being non-`omitempty` and the three containers rendering `{}`/`[]`." That reasoning is about the **JSON** encoder. In text mode the containers render `(none)`, which is 6 bytes of *added* text per empty container versus 2 for `{}`. The two modes have different constant floors and the RDR reasons about one and asserts both.

Six weeks in, someone will run the flag on a model shape the spikes did not cover, and the argument will be about whether to narrow C1 to JSON — an amendment to a locked normative contract, which this project's rules forbid, forcing either an abandon-and-iterate or a quiet oracle weakening. The second is more likely and worse.

Runner-up, and it will be rewritten in the same pass: `0023:§source-authority-census` (C-1's field count), because the declaration table has to be built from it and it does not survive contact with `reflect`.

---

## 3. The one assumption that will not survive first contact with a real user

**Not a lettered assumption — the unstated one under `0023:§consequences`: that shipping the capability is the same as delivering the benefit.**

The RDR is unusually honest about this, which is why it is the assumption that fails. `0023:§consequences`, "Scope of the deliverable," states plainly: "this RDR ships the CAPABILITY, not its adoption. The flag is opt-in and no consumer passes it on landing, so a green build delivers zero measured saving on day one — every oracle can pass while the motivating cost is still being paid." And then: "Adoption is not scheduled by this RDR and is not one of its phases."

The first real user is an agent-driven consumer whose navigator skill prompted this work (`intrastate#srz2`). That user's contact with the shipped flag goes like this:

1. They add `--plan-only` to a chained resolve call.
2. They observe a saving somewhere between 42% and 51%, because their model is a real state machine or a small decision table, not the synthetic 48-fact model the 79.4% headline was measured on (`evidence/spikes/a1-byte-width.md` §S1b: "Synthetic model, authored for this spike"). The Performance Expectations section leads with 79.4%.
3. They find `"revision":""` in the projected payload — a key that is constant-empty by construction (`internal/cli/flow_exec.go:49`, `func (flowRequest) revision() string { return "" }`) — and ask why the narrow mode carries a permanently empty field. The answer is four paragraphs of `0023:C2` about a slot for a value that no model can declare because "0002's `[model]` block admits no revision key."
4. They discover the saving applies to successes only (`0023:§consequences`: "a probing caller that mostly refuses saves nothing, by design"), and that the RDR explicitly declined to measure what fraction of their traffic refuses ("Unsized here: what share of the motivating consumer's traffic refuses").

The assumption that does not survive is that the 79.4% figure describes anything a user will experience. It describes a model authored inside the spike to make the ratio large: 48 required observed enum tags, all with identical stanzas, two rows discriminating on one of them. The echo group scales with tag count; the plan group does not. Any measurement designed by choosing the tag count is a measurement of the choice.

And `0023:§minimum-viable-validation` step 6 makes this the *gate*: "the 48-fact class saves at least 70% — the shipped figure is 79.4%, and a projection landing below 70% there means the plan group is carrying materially more than A1 measured." The pass bar is set against the synthetic. The real fixtures — the ones in `models/examples/` — are explicitly exempted: "The state-machine shapes carry no percentage bar."

---

## 4. Premortem

*Written from eleven weeks after the merge.*

We shipped `--plan-only` on schedule. All five oracles were green. The MVV table was recorded. `docs/cli-output-contract.md` gained its projected worked payload. Nothing about the landing was dramatic.

**Week 1 — the field count.** The first Phase 1 commit built the C2 declaration from `0023:§source-authority-census`. The census says fifteen fields; `reflect.TypeOf(resolvePayload{}).NumField()` returns fourteen. The implementer spent an afternoon deciding whether `writes`/`clear` were one field or two, concluded the census had merged them for readability, and — because S3's bidirectional equality would not go green against a table built from a document that miscounts — wrote S3 as a one-directional check: every struct field has a declaration entry. The reverse direction, entry-without-field, was dropped with a comment reading "the census is the authority for the entry set." That comment is now the load-bearing sentence in the partition's enforcement, and it points at a document with a wrong count in its first line.

**Week 2 — the docs-check ordering.** Phase 1's registration commit failed `make check` on `docs-check`, exactly as `0023:§phase-3-docs-and-help` predicted it would ("a Phase 1 commit that registers the flag without regenerating them turns CI red on its own"). We knew. We shipped the plan with the defect written into it. We merged Phases 1 and 3 into one commit, which meant the flag's usage string was authored before the projection existed and never revisited.

**Week 4 — 0024 lands.** `0024:C4` appended `dispositions` to `resolvePayload` after `emit`. Four things went red simultaneously, two of them in files 0024 never touched:

- `decision_table_0010_test.go:435` (`NumField() != 14`) — a 0010 oracle, which `0023:A3` promised would not move.
- `decision_table_0010_test.go:416` — the thirteen-key `slices.Equal` wire-order assertion. Exact set. No appended key survives it.
- `0023:S2`'s explicit projected-key literal — eight keys, now nine.
- `0023:S3`'s declaration equality — a struct field with no entry, which is the assertion working correctly and being read as an obstacle.

`0023:A4` had been stamped Verified on the claim that "whichever RDR lands second changes nothing in the other's contract." It was verified against `0024:C4`'s *join rule* — which is genuinely plan-group-only — and never against 0024's *implementation plan*, which says "Append `dispositions` to `resolvePayload` after `emit`." `0023:§phase-1-projection-on-the-verb` had even flagged the hazard in its own words — "Adding a field here would falsify A3's 'no predecessor oracle moves' claim, so a field addition is out of scope for this RDR by construction" — and drew the conclusion that it was out of scope *for this RDR*, not that a named sibling at depth 1 was about to do it anyway.

We amended the S2 literal. That is a normative key list, and amending it was a spec change made under a red build. Nobody wrote it up.

**Week 6 — the text-mode width clause.** A consumer ran `--plan-only --as=text` against a model with no observed tags. `respond/text.go::flatten` renders empty containers as `(none)` lines, not as `{}`, so the echo group's text footprint on that shape was five short lines. The projection still saved bytes, but the margin was 40 bytes on a 300-byte render, and on one gate-heavy shape with a long escape class the S1 text-width assertion went red by four bytes because the projected render's `revision: ` line and the retained gate array pushed it over. `0023:A8` had been **Pending** at lock with "not yet measured" in its Evidence field and a verification plan scheduled for the phase we had just finished. Its "If wrong" clause said the fix was to narrow C1 to JSON. C1 was locked. We narrowed the *oracle* instead, to "not longer than," and the clause `0023:C1` calls "the flag's whole reason for existing is itself oracle-enforced" stopped enforcing it.

**Week 9 — the reader pass.** Someone profiling the chained-call path noticed that under `--plan-only` the payload never renders `owned` or `readers`, and that `runFlowResolve` calls `req.runReaders(cmd.Context(), invokedReaders(req.model, outcome))` at line 180 — before the projection site, unconditionally, doing artifact I/O whose only *visible* consumers are two projected-away fields. They moved the flag read upstream and short-circuited the reader pass. The optimization was reviewed and approved: the diff touched no oracle, and the RDR's own Performance Expectations invited it by framing the projection as pure deletion.

Every test stayed green. `0023:S1`'s invoked-reader assertion compares `invokedReaders(model, outcome)` against `invokedReaders(model, outcome)` — `0023:C1` mandated exactly this ("the projected run's set is taken from the same in-process derivation the payload field is rendered from … never from its own output"), calling the payload-side alternative "unwritable." A pure function equals itself under every mutation of the code around it.

The failure surfaced in a customer's release pipeline. Their navigator skill chains `flow resolve --outcome build` then `flow resolve --outcome ship --plan-only`. The second call now resolves with an empty `Owned` set because `runReaders` never ran, so the kernel selects the row that matches on absent owned state rather than the one matching `phase=building`. `plan.RuleID` comes back `begin` instead of `ship-clean`. The user's agent transcribes the returned `writes` and moves the artifact backward. `escaped` is `false`, `rule` is a real rule id, the payload is well-formed and passes every consumer parse. `0023:D-identity` had promised: "the same invocation ± the flag selects the same rule, runs the same gates, and refuses identically — only report width differs."

**Week 10 — the always-keep core.** During triage someone proposed a second projection mode (`--emit-only`) to cut further. `0023:C2` had declined to mint an always-keep oracle, on the grounds that "with ONE projection mode, always-keep and the S2 key-set literal have the same extension" and that "the obligation it creates is on the successor." The successor was a kata, not an RDR. `--emit-only` dropped `escaped` and `escape_class`. The customer above, whose plan had been rescued by an escape row, read it as an ordinary plan. That is the exact failure `0023:§alternatives-considered` used to reject Alternative 1 — "a caller who omits `escaped` reads a rescued plan as an ordinary one" — and the design chose the fixed partition specifically to make it impossible. It made it impossible for *one* mode and wrote the guard for the second mode into a paragraph.

**Week 11 — the accounting.** We asked what the flag saved. Answer: nothing measurable. `0023:§consequences` had said it: "no consumer passes it on landing … Adoption is not scheduled by this RDR and is not one of its phases." Two consumers eventually adopted it. Both are state-machine models and both measured about 45%, against a Performance Expectations section leading with 79.4% — a figure from a model (`pricing-48.toml`) that exists only inside `evidence/spikes/a1-byte-width.md`, authored for the spike, never checked in, never run by anyone else. The MVV's 70% pass bar had been set against that same synthetic and cleared by it.

Seventeen hundred lines of record. The thing that broke was that a pure function was asked to disagree with itself, and a sibling RDR named in our own joint-check line added a struct field we had declared out of scope.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time gates: each is runnable against the record and the repo *before* lock, and each fails on the record as written.

**AT-1 — the census is the field inventory and must match `reflect` (catches C-1).**
```gherkin
Given the RDR names a payload struct at internal/cli/flow_resolve.go::resolvePayload
When I count rows in 0023:§source-authority-census, expanding merged rows to one field each
And I count reflect.TypeOf(resolvePayload{}).NumField()
And I read the prose count in the census's opening sentence
Then all three counts are equal
```
Result on the record: 13 rows / 14 fields / "fifteen" in prose. Fails three ways.

**AT-2 — joint-check peers are checked against their implementation plans, not only their contracts (catches C-2).**
```gherkin
Given 0023:JC1 fires on 0024 with a shared modify-anchor internal/cli/flow_resolve.go::resolvePayload
When I read 0024's Implementation Plan for edits to that anchor
Then no edit changes the field COUNT or the top-level KEY SET of resolvePayload
And if it does, every oracle whose assertion is extensional over that set
    (0023:S2's key literal, 0023:S3's declaration equality,
     0023:C2's always-keep coincidence argument,
     decision_table_0010_test.go:435's NumField pin,
     decision_table_0010_test.go:416's slices.Equal wire-key-order literal)
    is named in 0023:§prerequisites with its landing-order obligation
```
Result: 0024 says "Append `dispositions` to `resolvePayload` after `emit`." Five extensional oracles are unnamed, two of them in a predecessor's suite that `0023:A3` certifies will not move. `0023:§prerequisites` says only "either RDR may land first." Fails.

**AT-3 — no locked MUST rests on a Pending assumption (catches C-3, C-9).**
```gherkin
Given a clause in 0023:§normative-contracts states a MUST
When I trace the assumption that clause's satisfiability depends on
Then that assumption's Status is Verified
And its Evidence field is non-empty and not "not yet measured"
```
Result: C1's both-modes width clause depends on `0023:A8` (Pending, "not yet measured"). C1's whole-tree walker clause depends on `0023:A5` (Pending, implementability half retracted at the cove sweep). Fails twice.

**AT-4 — every oracle scenario asserts one measurand, or names the relation between two (catches C-4).**
```gherkin
Given a Testing Strategy scenario asserts two properties
When both properties are asserted on the same run pair
Then the record states whether one implies the other, and which failures each catches alone
```
Result: `0023:S5` asserts set-subset over lines AND strict reduction over total bytes, with a negative control for the second and none establishing that the pair is jointly sufficient. `0023:S1` does the same for key-subset and width — that one *is* addressed ("The width assertion is not implied by the subset one"). S5's pair is not. Fails on S5.

**AT-5 — no field is contracted into the projected width with a provably constant value (catches C-5).**
```gherkin
Given 0023:C2 assigns `revision` to the PLAN group and 0023:C1 carries it under the flag
When I read its producer
Then the producer can return at least one distinct value
```
Result: `internal/cli/flow_exec.go:49` is `func (flowRequest) revision() string { return "" }`. The RDR concedes it ("presently constant-empty… no model can declare one"). Fails; the field is unconditionally dead weight in the mode whose purpose is to remove weight, and in text mode renders as a blank-valued line.

**AT-6 — no differential assertion compares a pure function to itself (catches C-6).**
```gherkin
Given C1 mandates an assertion channel for the invoked-reader set
When I identify the function supplying both sides
Then that function's arguments differ between the two sides, or a side is observed rather than recomputed
And the record's negative control for that assertion is runnable as written
```
Result: both sides are `invokedReaders(req.model, outcome)`; neither argument depends on the flag. `0023:§oracle-discriminability`'s control — "change the invoked-reader derivation under the flag" — has no "under the flag" variant to change. Fails. The fix at review time: assert the *reader pass ran* (spy on `req.runReaders`, or compare the artifact-read count across the ± pair), not the derivation's output.

**AT-7 — a numeric pass bar is set on a checked-in fixture (catches C-7, C-14).**
```gherkin
Given 0023:§minimum-viable-validation step 6 gates lock on "at least 70%" for a named class
When I locate the model that class is measured on
Then that model is checked into the repository
And it was not authored inside the spike that measures it
```
Result: `pricing-48.toml` exists only in `evidence/spikes/a1-byte-width.md`, described there as "Synthetic model, authored for this spike." Fails. The three checked-in fixtures measure 42.2–51.3% and are explicitly exempted from any bar.

**AT-8 — the benefit has a named owner and a landing condition (catches C-8).**
```gherkin
Given the RDR's Problem Statement quantifies a cost in transcript tokens
When I read the Implementation Plan phases
Then at least one phase, or a named tracked issue with an owner, realizes that saving
```
Result: `0023:§consequences` states "Adoption is not scheduled by this RDR and is not one of its phases." No phase, no owner, no tracked adoption issue. Fails.

**AT-9 — restated with teeth (see AT-3) for A5's specific shape.**
```gherkin
Given A5's implementability half was retracted at the cove sweep
And 0023:§prerequisites carries it as a Phase 2 build step
Then no section outside the Implementation Plan narrates the walker as settled
```
Result: `0023:C1` (three paragraphs of MUST), `0023:§disposition`, `0023:§oracle-discriminability` (S4 row with two controls), `0023:S4`, MVV step 5, and `0023:§desk-trace` row 5 all narrate it as settled. Fails — and violates the record's own gate text at `0023:§assumption-verification`.

**AT-10 — normative text is proportional to blast radius (catches C-10, C-19).**
```gherkin
Given each normative clause in C1 and C2
When I rank clauses by lines of normative text
And rank the same clauses by the user-visible harm of their violation
Then the two rankings are not inverted at the top
```
Result: the whole-tree walker (~40 lines across C1, plus A5, plus a discriminability row, an S4 scenario, an MVV step, and a desk-trace row) guards a flag registered on `completion`. The projection mechanism — the type conversion of five fields on the one shipped payload struct — is one deferred sentence in `0023:D-wire-byte-format`. Inverted. Fails.

**AT-11 — sweep counts are reproducible (catches C-11).**
```gherkin
Given A3 states "swept all 21 test files invoking flow resolve"
When I run the sweep it describes
Then the count matches
```
Result: 22 files. The count was already corrected once ("the earlier '16' was a miscount"). Fails.

**AT-12 — the independence C2 requires is a property of the chosen carrier, not of the requirement (catches C-12).**
```gherkin
Given C2 requires the assignment source be independent of the projection code
And C2 defers the carrier choice to Phase 1
When I evaluate each admitted carrier
Then every admitted carrier satisfies independence under the A2 mechanism
```
Result: C2 admits "a per-field marker on `resolvePayload` or a table keyed by field name." Under A2's pointer-`omitempty` mechanism the projection *is* an edit to `resolvePayload`'s field declarations, so a per-field marker on that struct moves in the same edit as the projection — not independent. Fails; C2 must name the table and drop the marker, or state why the marker is independent.

**AT-13 — version skew is analysed in both directions (catches C-13).**
```gherkin
Given 0023:§disposition claims "loud" for a binary predating this RDR
When I consider a binary ROLLBACK after a consumer has adopted the flag
Then the record states the consumer's degradation path
```
Result: only the forward direction is analysed, in both `0023:§disposition` and `0023:F5`. A rollback gives exit 2 with the flag name in pflag free text and no structured code (the accepted `0011:A14` limit), mid-chain. No degradation path stated. Fails.

**AT-15 — default-mode byte identity is asserted against a PRE-CHANGE artifact (catches C-15, C-16).**
```gherkin
Given C1's headline claim is "Default-mode output is byte-identical to today's"
When I identify which of S1..S5 asserts it
Then at least one oracle compares against bytes captured BEFORE the change
```
Result: none. S1 compares projected against default *within the same build*; both sides drift together. The only pre-change comparison in the whole record is `0023:§minimum-viable-validation` step 1, which is a manual run. The A2 spike's byte identity was established against "a live CLI reference line" captured during the spike — that artifact is the right one, and no oracle cites it as a fixture. Fails. Fix at review time: pin the pre-change default line for at least the pricing 2×2 and one state-machine shape as a checked-in golden, and make S1's first assertion `default_output == golden`.

**AT-16 — nil-vs-empty is settled for every field the pointer conversion touches (catches C-15's second half).**
```gherkin
Given A2 converts Model, Observed, Owned, Readers, Outcome to pointers
When I read each field's writer at internal/cli/flow_resolve.go:223-244
Then the record states, per field, whether the writer can produce nil
And whether that nil-avoidance is contracted or incidental
And what the pointer conversion renders if it ever changes
```
Result: the record states the general hazard ("bare-map `omitempty` would drop" empty `{}`) and verifies the pointer fix in the spike, but never enumerates per-writer nil-ability. All five writers are in fact nil-safe today — `observedTagMap`/`tagMap` at `internal/cli/flow_exec.go:739-753` and `readerIDs` at `internal/cli/flow_next.go:582-588` all `make(...)` unconditionally. But only `emitMap` documents that as deliberate (`flow_resolve.go:279-284`), and `runFlowResolve` guards `Gates` and `Next` at `:245-250` while leaving `Readers` unguarded — evidence the invariant is incidental rather than held. I verified the consequence: a `*[]string` pointing at a nil slice encodes as `"readers":null` under `omitempty`. Fails — the record makes default-mode byte identity depend on an invariant it neither states nor asserts.

**AT-17 — phase ordering does not contain a known CI break (catches C-17).**
```gherkin
Given the Implementation Plan orders phases
When I read each phase for a stated CI consequence
Then no phase states that an earlier phase's commit turns CI red
```
Result: `0023:§phase-3-docs-and-help` states exactly that. Fails; the fix is to move doc regeneration into Phase 1, not to document the break.

**AT-18 — a safety invariant used to reject an alternative is asserted by an oracle in THIS record (catches C-18).**
```gherkin
Given 0023:ALT1 is rejected because "a caller who omits `escaped` reads a rescued plan as an ordinary one"
And C2 declares an always-keep core of rule, escaped, escape_class, revision
When I locate the oracle asserting the core survives projection
Then it exists in this RDR's Phase 2 battery and is not delegated to a successor
```
Result: `0023:C2` — "This RDR therefore does NOT mint a separate always-keep oracle… the obligation it creates is on the successor." The invariant that justifies the entire chosen design over Alternative 1 is enforced, by the record's own words, "by coincidence of scope." Fails.

**AT-19 — the record is readable whole (catches C-19).**
```gherkin
Given the change is one boolean flag deleting five keys from one struct before one function call
When I measure the record
Then it is under 600 lines
And 0023:§proportionality carries a written response, not an unfilled template
```
Result: 1617 lines; `0023:G-proportionality` is the unedited template block. Fails.
