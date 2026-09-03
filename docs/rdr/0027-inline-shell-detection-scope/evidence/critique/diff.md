Model: claude-opus-5

# Critique diff — RDR 0027 (two independent hostile passes)

Passes diffed:
- **Pass A** — `evidence/critique/critique.md`, stamped `Model: claude-opus-5`. 13 findings (C-1 … C-13), premortem lens.
- **Pass B** — `evidence/critique/critique-modelB.md`, stamped `Model: claude-sonnet-5`. 9 findings (C-1 … C-9), premortem lens.

Note: pass B ran on **Sonnet 5** after Fable 5.1 was unavailable (spend limit). The
cross-model independence the lens wants is therefore Opus-vs-Sonnet, not
Opus-vs-Fable; convergence below should be read with that weaker independence
in mind, but the two passes were run without sight of each other and their
ledger IDs are per-file (`A:C-1` != `B:C-1`).

Ledger IDs are NOT comparable across passes. Everything below is keyed on the
RDR passage anchor.

## Grounding performed for this diff

All predicate claims were re-run, not taken on trust. A scratch harness
(`/tmp/p27/main.go`) reimplements today's `interpreterForm`
(`internal/table/load.go:1162-1184`, verbatim: `env` walk, then
`shellInterpreters[filepathBase(argv[i])]` at that ONE index) and C1's
position-free successor (`0027:§illustrative-code`), and computes both over
every vector either pass names. `shellInterpreters`
(`internal/table/load.go:1035-1047`) was read directly: **11 entries** —
`sh bash dash ksh zsh csh tcsh` → `-c`; `python` → `-c`; `ruby` → `-e`;
`node` → `-e --eval`; `php` → `-r`. Stdin wiring read at
`internal/cli/cmdbind/cmdbind.go:207`, `:519-545`, `:723`.

---

## 1. CONVERGENT — both passes, same passage, compatible failure mode

Ranked first; convergence across independent runs is the strongest signal here.
Three passages converge, and on the two that matter the convergence is exact.

| # | Anchor | Pass A | Pass B | Agreed failure mode | Divergence |
|---|--------|--------|--------|---------------------|------------|
| V-1 | `0027:C1` (`predicate:` line) + `0027:A2` | A:C-1, A:C-2, A:C-8 | B:C-5, B:C-6 | The position-free predicate refuses argv in which **no interpreter is executed**, because `base(argv[i])` is matched with no command-position requirement; and `0027:A2`'s corpus (32 in-repo vectors, 11 green) is structurally incapable of containing the phenomenon it certifies absent. | **Severity and fix differ sharply.** A calls it the record's defining defect, quantifies it (6/14 guessed vectors), and its AT-2 is explicitly *incompatible with C1 as written* — i.e. A demands the contract change. B accepts it may be a legitimate design cost and its AT-5 offers an alternative disposition: if the class is accepted, **write it into Trade-offs/Failure Modes as an explicit named cost**. B's fix is cheaper and lands inside the record; A's requires re-pricing the "no wrapper table" economy. |
| V-2 | `0027:A6` (+ `0027:C1` `promise:` clause) | A:C-3, A:C-4 | B:C-1, B:C-8 | `A6` is `Pending`, and `C1`'s `promise:` clause is present-tense normative prose depending on it. Phase 3 has no landing site and `0027:§phase-3-promise-wording` punts the choice ("the implementer's call"). Both predict Phase 3 ships nowhere reviewer-reachable or silently drops. | A frames it as a **mechanical gate violation** — `0027:G-assumptions` forbids settled-fact prose depending on a `Pending` assumption, so the record cannot pass its own Finalization Gate, no judgement needed. B frames it as a **Proportionality/split** problem (B:C-8): one Profile covering two independently-testable surfaces, and its AT-1 offers splitting Phase 3 into its own RDR as an alternative to resolving A6. A's framing is the stronger one because it is checkable; B's remedy is the more practical. |
| V-3 | `0027:C1` (`interpreter set:` line) | A:C-6 | B:C-4 | The deny-list is frozen-by-reference with **no version marker**: C1 says "the list grows by amendment of THIS clause" while the membership lives in a Go map C1 does not restate. `0027:S6` golden-asserts `Categories()` order; **nothing asserts `shellInterpreters` content** (confirmed — no test in the repo pins the map). | Different failure directions from the same root. A worries about the **blast radius growing silently** (add `perl` — parked by A2 — or `python3`, parked by `0027:BR4` — and C1's promise text goes stale). B worries about the **amendment path being too heavy** (adding `uv`/`deno`/`bun` forces an RDR reopen, or the map drifts out of sync). Both are real and they are the same missing artifact: a version marker plus a golden test over the map. |
| V-4 | `0027:D-selection-predicate` (tie-break) | A:C-7 | B:C-3 | `["python","sh","-c","echo"]` reports `python -c`, not the `sh -c` a reader's eye goes to; accepted as "message attribution only" with no test and no worked example. Verified: both predicates report `python -c` for that vector. | A argues the defect is **amplified by V-1** — under the widened predicate the reported pair is *routinely a coincidence*, so "message attribution only" (defensible when both pairs are real) becomes indefensible. B treats it as standalone and adds a distinct concern A misses entirely (see B-only, `0027:§normative-contracts` `report:` format stability). A's coupling argument is correct and is the stronger reading. |
| V-5 | `0027:A5` + `0027:§capability-dependencies` | A:C-5, A:C-10 | B:C-2 | The stdin axis (`sh -s`, bare `sh`, `python -`) is routed to a charted successor that **no kata tracks**. A5 is stamped `Verified` by `Method: Design Decision` with `Evidence: this RDR's scoping choice` — self-referential, and `0027:G-assumptions` explicitly requires "no `Verified` stamp is self-referential". Both note A5's own durability caveat concedes the gap. | A goes further and attacks the **mechanism**, not just the tracking: it claims "withholding the envelope" is falsified by unconditional `cmd.Stdin`. Adjudicated below (CONTESTED-1) — A is right, and B missed it. B's contribution is the sharper *process* framing: an assumption whose own Evidence field says it is fragile, stamped `Verified` anyway. |

**Convergence assessment.** V-1 and V-2 are the two that drive disposition, and
both passes reached them from different directions with different vector sets
(A: `grep`/`pip`/`rg`/`tar`/`npm`; B: `polybuild`/`multi-tool`/`build-tool`).
Every vector from **both** passes reproduces. That is the strongest result in
this diff.

---

## 2. A-ONLY — raised by pass A alone, adjudicated

| Anchor | A's finding | Verdict | Grounding |
|--------|-------------|---------|-----------|
| `0027:A5` mechanism | A:C-5 — "withholding the envelope" is falsified; the **write** path sends author-controlled bytes. | **REAL, and B missed it.** Adjudicated in full at CONTESTED-1. | `cmdbind.go:207`, `:723` read directly; exploit reproduced. |
| `0027:§references` | A:C-11 — References cites `internal/cmdbind/cmdbind.go::resolveArgv0`; real path is `internal/cli/cmdbind/cmdbind.go`. | **REAL.** `ls internal/cmdbind` → no such directory; `internal/cli/cmdbind` exists. `0027:§existing-infrastructure-audit` uses the **correct** path in the same record, so the record contradicts itself. A's second-order claim — that one miss withdraws the ground-sweep's "clean (17 anchors)" verdict as evidence for the other sixteen — is fair but is a process inference, not a defect in itself. Cheap to fix. | Filesystem check. |
| `0027:MVV` step 3 / `0027:S4` | A:C-12 — the MVV mandates a behavioural `env -S` step that "cannot execute on the stated development platform" because stock darwin BSD `env` rejects `-S`. | **ARTIFACT of A's reading — falsified.** A took the A4 spike's portability note at face value instead of checking. On **this** host: `which env` → `/usr/bin/env`; `/usr/bin/env -S "sh -c 'echo HELLO_S'"` → `HELLO_S`, exit 0. Stock `/usr/bin/env` on this darwin **does** support `-S`. The MVV step is runnable as written. The real (much smaller) finding is that the **A4 spike's own portability claim is stale** — it says "BSD env rejects `-S` outright (illegal option)" and points at `/opt/local/bin/genv` as the only working path, which is no longer true here. That is a spike-evidence correction, not an MVV defect. | Executed on the dev host. |
| `0027:§approach` / `0027:§prerequisites` | A:C-13 — the q2q1 subsumption is silent on the **test** disposition if q2q1 lands first. | **REAL but minor.** `0027:§prerequisites` reads: "intrastate#q2q1's disposition recorded (landed first, or closed as subsumed) — either way the `env` walk is removed here." Confirmed: the record decides the *code* both ways and says nothing about tests q2q1 would have added to assert the walk. A genuine gap, but it is a merge-hygiene note, not a contract defect — and it only bites in one of two branches. Stage 6 material. | `0027:§prerequisites` read via projector. |
| `0027:§decision-rationale` (Premortem para) | A:C-9 — the record's own premortem names the false-refusal risk then disarms it in the same sentence ("the already-accepted flag-anywhere class, not new, but the message was the failure"). | **REAL as a process observation.** The passage does exactly that, and `0027:§key-discoveries` states as **Documented** that "freeing the position adds no new false-refusal class beyond the one already accepted" — which the V-1 grounding falsifies. A has no separate symptom for C-9 ("n/a, process defect") and correctly says its symptoms are V-1. Not independently dispositive; **fold into V-1**. | `0027:§key-discoveries` read via projector. |

---

## 3. B-ONLY — raised by pass B alone, adjudicated

| Anchor | B's finding | Verdict | Grounding |
|--------|-------------|---------|-----------|
| `0027:§normative-contracts` (`report:` line) | B:C-9 — "the detail additionally names the two matched words" is a normative commitment with **no format contract**: no separator, ordering, or stability guarantee. Asymmetry: `Categories()` order IS golden-pinned (`0027:S6`), the human-read detail text is not. | **REAL, and A missed it entirely.** Confirmed: `C1`'s `report:` line commits only that "the detail additionally names the two matched words", and `0027:S5`'s Expected is "the detail names `ruby -e` — the two matched words", i.e. containment, not format. Meanwhile `0027:S6` pins `Categories()` membership **and order**. B's asymmetry observation is exactly right and is a genuinely independent finding. Note the implementation does produce a stable `argv[i] + " " + argv[j]` (confirmed at `load.go:1180` and in C1's predicate line), so the format exists — it is simply **unasserted**. Low severity, cheap fix. | `load.go:1180`; `0027:S5`/`S6` via projector. |
| `0027:§performance-expectations` | B:C-7 — O(len(argv)²) is dismissed because argv is "single digits of words" in a "hand-authored model entry", but **nothing in C1 bounds argv length**, and templated/generated bindings are plausible. | **REAL but weak; A was right to skip it.** The passage does say exactly that, and it is true that C1 states no bound. But the scan is per-entry at load time over a hand-authored TOML file; at 200 elements this is 40k string comparisons, i.e. microseconds. B's own AT-6 asks only for "a stated contract or test". This is a theoretical gap with no plausible symptom. **Dismiss** — but it is a legitimate observation, not a misreading. | `0027:§performance-expectations` via projector. |
| `0027:§metadata` (Profile) | B:C-8 — Profile `large` covers two independently-testable surfaces (a fully-verified predicate, A1-A4 all `Verified`; and a non-existent description surface, A6 `Pending`), and the Proportionality gate's split test was never run against that decomposition. | **REAL, and a genuinely different lens on V-2.** A saw the same underlying facts (`0027:C1`'s promise depends on Pending A6) and drew a *gate-violation* conclusion; B drew a *decomposition* conclusion — that the record should have been split. Both are supportable. B's version is constructive where A's is only blocking: splitting Phase 3 into its own RDR resolves the gate violation A identifies **without** stalling the predicate half, which is fully verified. Worth carrying forward as the preferred remedy for V-2. | `0027:§proportionality`, `0027:§critical-assumptions` via projector. |

---

## 4. CONTESTED — the passes disagree, adjudicated against source

### CONTESTED-1 — `0027:A5`: is the withholding mechanism falsified?

**A (C-5):** falsified. `cmdbind.spawn` sets `cmd.Stdin` **unconditionally**, and
a `write` binding sends `stdinObject(planned)` — a JSON object whose *values* are
model-authored tag values. A `write` accessor declared `["sh","-s"]` receives
author-controlled bytes. "Withholding is not the mechanism for write; the
envelope is the payload." A concedes read/gate is "inert by accident of JSON
syntax" (`printf '{}' | sh -s` exits 127).

**B (C-2):** does not address the mechanism at all. B's objection is purely that
the successor has no kata — a tracking gap.

**Verdict: A is right on the mechanism, and materially right on the exploit —
but A's own characterisation of *why* is wrong, and the severity needs
restating.**

What I checked:

1. `internal/cli/cmdbind/cmdbind.go:207` — `cmd.Stdin = bytes.NewReader(stdin)`.
   **Unconditional.** No branch on declared appetite. A's structural claim holds.
2. `:723` — `(*Writer).Apply` calls `spawn(..., stdinObject(planned))`.
   `stdinObject` (`:523-545`) builds `map[string]string` from
   `t.Key` → `t.Value` and JSON-encodes it with `SetEscapeHTML(false)`.
   **Tag values are model-authored and reach stdin.** A's claim holds.
3. `:567` and `:658` — read and gate send `stdinObject(nil)` = `{}`. A's
   read/gate analysis holds.

So A's structure is correct. But A's stated exploit mechanism is not, and I
tested it:

- A implies the envelope's newlines let a value become its own shell line. **It
  does not.** Go's JSON encoder escapes `\n` as the two-character sequence `\n`,
  so the envelope is always **one line**. I built the real encoder
  (`SetEscapeHTML(false)`, `TrimRight`) and confirmed: a value of
  `"x\ntouch /tmp/pwned_B\n"` serialises to the single line
  `{"k":"x\ntouch /tmp/pwned_B\n"}`, and piping that to `sh -s` produces one
  "No such file or directory" and creates nothing.
- Semicolon/quote breakout also fails: `{"k":"a\"; touch /tmp/pwned_D; echo \""}`
  → `sh` treats the whole brace-word as one command name, no file created.

**But the finding survives by a different route, and I reproduced it.** `sh`
performs **command substitution on the word before it fails to execute it**:

```
$ go run enc2.go            # real stdinObject encoding of {"k":"$(touch /tmp/pwned_C)"}
{"k":"$(touch /tmp/pwned_C)"}
$ go run enc2.go | sh -s
sh: line 1: {k:}: command not found
$ ls /tmp/pwned_C
/tmp/pwned_C                # <-- created
```

A model-authored tag value containing `$(...)` or backticks **executes** when a
`write` accessor is declared `["sh","-s"]`, which `0027:C1` admits by contract
and `0027:F2` calls "silent by contract". The exit-127 inertness A attributes to
read/gate is not a general defence — substitution runs regardless of whether the
resulting word is a valid command.

**Net:** A5's `Verified` stamp is unsound on two independent counts. (i) The
mechanism claim — "withholding the envelope … leaves a stdin-reading interpreter
nothing to execute" — is not available on the write path, where the envelope
*is* the payload; this is a source fact, `cmdbind.go:207`/`:723`. (ii) The stamp
is `Method: Design Decision` / `Evidence: this RDR's scoping choice`, which
`0027:G-assumptions` independently forbids ("no `Verified` stamp is
self-referential"). B's tracking objection is true but strictly weaker. A wins
this contest; the record should not carry A5 as `Verified`.

### CONTESTED-2 — `0027:C1` predicate: how large is the false-refusal class, and does A's count survive?

**A (C-1):** names `grep -c`, `pip -c`, `rg -c`, `tar`+`python`, `npm`+`node -e`,
`my-tool --lang php --mode -r`, plus `cat php -r`, `echo node --eval`,
`mytool --interpreter python --flag -c`. Reports **6 false refusals from 14
guesses** and calls the class "roughly an order of magnitude larger" than the
accepted one.

**B (C-5):** same class, different vectors (`polybuild --runtime ruby --stage -e`,
`multi-tool --lang ruby --emit -e`), and — unlike A — leaves open that the record
may legitimately **accept** the class if it writes it down.

**The crux the task poses:** does the predicate require the interpreter word to
be a *listed basename*, and are `grep`/`pip`/`rg` on that list? **They are not.**
`shellInterpreters` is exactly 11 entries and contains no `grep`, `pip`, `rg`,
`tar`, `npm`, `cat`, or `echo`.

**This does NOT collapse A's claim, and the task's framing of the crux is the
one thing to correct.** A never claims `grep` is listed. In every A vector the
matched interpreter is a listed name appearing as a **data argument**, with the
unlisted tool at argv0:

- `["grep","-r","python","-c","src/"]` → i=2 (`python`, listed), j=3 (`-c`,
  python's flag). argv0 `grep` is never examined, because C1 says "Nothing before
  argv[i] is read."

So the second horn of the crux is what applies: the listed set contains words
(`python`, `ruby`, `node`, `php`, and `sh`) that are also common non-interpreter
argv tokens, and their flags (`-c`, `-e`, `-r`, `--eval`) are among the most
reused short flags in Unix. **A's claim stands.**

**Which vectors survive scrutiny — verified, not repeated from A.** I ran both
predicates over every vector either pass names:

```
grep -r python -c src/                                  old=false new=true   python -c   FLIP
pip install python -c constraints.txt                   old=false new=true   python -c   FLIP
rg python -c                                            old=false new=true   python -c   FLIP
tar -czf out.tgz python -c                              old=false new=true   python -c   FLIP
npm run node -- -e                                      old=false new=true   node -e     FLIP
my-tool --lang php --mode -r                            old=false new=true   php -r      FLIP
echo node --eval                                        old=false new=true   node --eval FLIP
cat php -r                                              old=false new=true   php -r      FLIP
mytool --interpreter python --flag -c                   old=false new=true   python -c   FLIP
polybuild --target service-a --runtime ruby --stage -e  old=false new=true   ruby -e     FLIP   (B)
multi-tool --lang ruby --emit -e                        old=false new=true   ruby -e     FLIP   (B)
build-tool --interpreter ruby run-step --flag -e        old=false new=true   ruby -e     FLIP   (B)
--- controls ---
env -i sh -c echo hi                                    old=false new=true   sh -c       FLIP   (INTENDED)
nice sh -c cat {artifact}                               old=false new=true   sh -c       FLIP   (INTENDED)
python sh -c echo                                       old=true  new=true   python -c          (tie-break)
ruby tool.rb -e prod                                    old=true  new=true   ruby -e            (S5)
sh ./gate.sh                                            old=false new=false  --                 (admitted)
rdr-gate {artifact}                                     old=false new=false  --                 (repo fixture)
perl -e print 1                                         old=false new=false  --                 (unlisted)
sh -s                                                   old=false new=false  --                 (admitted)
```

**All 6 of A's headline vectors reproduce exactly**, all 3 of A's AT-2 extras
reproduce, and all 3 of B's reproduce. **12 shell-free vectors flip
`old=false → new=true`, every one executing a non-interpreter argv0** (`grep`,
`pip`, `rg`, `tar`, `npm`, `my-tool`, `echo`, `cat`, `mytool`, `polybuild`,
`multi-tool`, `build-tool`). Not one spawns a shell. The controls behave
correctly: the two wrapper forms flip as **intended** (that is the RDR's whole
purpose), and the admitted/unlisted forms stay green.

**Two corrections to A's rhetoric**, neither of which touches the finding:

1. A's "6 false refusals from 14 guesses" is a **guess-selection statistic, not a
   base rate**. A chose the 14 and does not list the 8 that did not flip. It
   licenses "this class is easy to construct"; it does not license A's "order of
   magnitude larger" quantification. The honest statement is: the class is real,
   trivially constructible, and **unmeasured** — which is exactly what
   `0027:A2` should have measured and structurally could not.
2. A's C-2 asserts the record's "unchanged" language is simply **false**. That is
   correct and verified: today `interpreterForm` matches
   `shellInterpreters[filepathBase(argv[i])]` at **one** index (argv0 after the
   `env` walk), so `argv[i]` is the thing that gets executed. `0027:§key-discoveries`
   states as **Documented** that "freeing the position adds no new false-refusal
   class beyond the one already accepted." The 12 flips above falsify that
   sentence. It is documented, and it is wrong.

**Verdict: A's reading is the one the evidence supports on the facts; B's is the
one it supports on the remedy.** The class is real (A). But B is right that C1
may legitimately *accept* it — the record simply has to say so in
`0027:§trade-offs`/`0027:§failure-modes` and stop asserting the opposite in
`0027:§key-discoveries` and `0027:§consequences`. The defect is not
"the predicate is wrong"; it is **"the record asserts an untrue thing about the
predicate and stamped an assumption `Verified` that could not see it."**

### CONTESTED-3 — `0027:A6`: does a reviewer-reachable description surface exist today?

**A (C-3):** `table.Categories()` has **zero non-test consumers**, and
`docs/cli-output-contract.md` "defines output entirely as the `CLIError`
envelope with `findings[]`, i.e. refusal-shaped" — so no success-path text
channel exists, and inventing one amends a contract this RDR does not own.

**B (C-1):** softer — "no owner, no deadline, and no fallback if no
reviewer-reachable surface exists". B does not assert the surface cannot exist.

**Verdict: A's fact is right; A's inference is overreached. B's weaker claim is
the safer one, and the reconciled fact sits between them.**

What I checked:

1. **`Categories()` consumers.** Every `.go` reference outside `category.go`'s own
   definition and comments is in a `_test.go` file:
   `emit_witness_0024_test.go:53`, `dump_test.go:874,877,883,895,1028,1044`,
   `roundtrip_test.go:207,1315`, `reserved_key_fixtures_0008_test.go:83`,
   `command_carrier_0025_test.go:853,876,909,1089`. **A's "zero non-test
   consumers" is exactly correct**, and B's ledger does not establish this.
2. **`internal/table/category.go`** — `Category` is `type Category string`, a bare
   string with no description field. `Categories()` returns `[]Category`,
   identifiers only. Confirmed; matches A6's own "absence half is confirmed".
3. **The output contract.** **A overstates here.** `docs/cli-output-contract.md`
   is not refusal-only: it documents a **`flow resolve` success payload**
   carrying `emit` and `dispositions` (§"`flow resolve`, the `emit` answer, and
   its `dispositions`"; §"`flow resolve --plan-only` and the two halves of the
   payload"). A success-path channel demonstrably exists. What does **not** exist
   is a channel for *category descriptions* — nothing in the contract carries
   per-category prose, on either the success or the failure side.

**Reconciled fact — the answer to the question posed:** **No reviewer-reachable
description surface exists today.** The refusal `Detail`
(`load.go::carrierDefect`) fires only on refusal and so cannot carry a promise
about what is *admitted* — which is A6's own correct observation and the reason
the promise needs a new home. `Categories()` is a test-visible artifact with no
CLI consumer. So A6's spike must **create** a surface, not find one.

But A's stronger claim — that creating one is necessarily an amendment to an
output contract this RDR does not own — **does not follow**, because the
contract already accommodates success payloads and a lint-help path (which A6
itself names as a candidate: "a documented lint-help path") need not touch the
`CLIError` envelope at all. A6 is genuinely `Pending`, genuinely un-owned, and
genuinely blocks `C1`'s `promise:` clause under `0027:G-assumptions` — but it is
**not** the "probably cannot exist" that A asserts. B's "no owner or fallback"
is the accurate characterisation of the risk; A's is the accurate
characterisation of the fact.

---

## Disposition recommendation

- `0027:§key-discoveries` + `0027:§consequences` + `0027:D-selection-predicate` ("the already-accepted false-refusal class … unchanged") — **fix now.** Verified false: 12 shell-free vectors flip `old=false → new=true` under C1 where today's predicate examines only argv0-after-`env`. The record asserts as **Documented** a thing the predicate does not do. This is a factual correction to locked-adjacent prose, not a design change.
- `0027:C1` (`predicate:` line, position-freedom itself) — **flag to Stage 6.** The class is real and trivially constructible, but B is right that it may be *accepted*. Stage 6 must choose: bound the predicate (which re-prices the "no wrapper table" economy `0027:§decision-rationale` factor (2) rests on, and collides with `0027:ALT2`'s prohibition), or accept the class and name it in `0027:§trade-offs`/`0027:§failure-modes` as an explicit cost. Do not lock while the record still claims the class is unchanged.
- `0027:A2` (`Verified` stamp) — **fix now.** The stamp is not wrong about what it measured; it is wrong about what it licenses. Its 32-vector corpus (11 green, authored before the phenomenon existed) structurally cannot contain a listed basename in a non-argv0 position. Either re-scope the claim's wording to what the corpus supports, or re-run against an out-of-tree-shaped vector set. It must not stand as evidence that the widened predicate refuses nothing new.
- `0027:A5` (`Verified` stamp + withholding mechanism) — **fix now.** Two independent grounds. (i) The mechanism claim is false on the write path: `cmdbind.go:207` sets `cmd.Stdin` unconditionally and `:723` sends model-authored tag values; a `$(...)` tag value under a declared `["sh","-s"]` write accessor **executes** (reproduced). (ii) `Method: Design Decision` / `Evidence: this RDR's scoping choice` is self-referential, which `0027:G-assumptions` forbids outright. Downgrade the status and correct the mechanism sentence.
- `0027:A6` (`Pending`) + `0027:C1` (`promise:` clause) — **fix now.** Mechanical gate violation: `0027:G-assumptions` forbids settled-fact prose depending on a `Pending` assumption, and `C1`'s `promise:` is present-tense normative prose depending on A6. Established fact: no reviewer-reachable description surface exists today (`Categories()` has zero non-test consumers; `Category` is a bare string; the refusal `Detail` fires only on refusal). Prefer **B's remedy** — split Phase 3 / the `promise:` clause into its own record so the fully-verified predicate half is not held hostage — over A's blanket block.
- `0027:§capability-dependencies` (stdin successor row, `Deferred`, no kata) — **flag to Stage 6.** A dependency with no tracking artifact cannot be cited as an owner. Either open the kata and link it from A5, or restate C1's out-of-scope routing as an unowned admission.
- `0027:C1` (`interpreter set:` line) — **flag to Stage 6.** No version marker, and no test pins `shellInterpreters` content (confirmed: `0027:S6` covers `Categories()` only). Needs a golden test over the map naming `0027:C1`, plus a stated amendment path — both passes converged here from opposite directions (blast-radius growth vs. amendment friction).
- `0027:§references` (`internal/cmdbind/cmdbind.go::resolveArgv0`) — **fix now.** Wrong path; the directory does not exist. Correct is `internal/cli/cmdbind/cmdbind.go`, which `0027:§existing-infrastructure-audit` already uses. One-word edit, and it removes a self-contradiction.
- `0027:§normative-contracts` (`report:` line, detail-string format) — **flag to Stage 6.** B-only and correct: the format is stable in implementation (`load.go:1180`, `argv[i] + " " + argv[j]`) but asserted nowhere; `0027:S5` asserts containment only while `0027:S6` pins `Categories()` order. Cheap to close by tightening S5's Expected.
- `0027:D-selection-predicate` (tie-break attribution) — **flag to Stage 6.** Both passes flagged it; A's coupling argument is right that it degrades from "cosmetic" to "routinely misleading" if the position-free class is accepted. Its disposition follows whatever Stage 6 decides for `0027:C1`, so decide it there, not separately.
- `0027:§approach` / `0027:§prerequisites` (q2q1 test disposition) — **flag to Stage 6.** Real but minor: the record decides the `env` walk's fate both ways and is silent on the tests q2q1 would have added. Merge hygiene, one sentence.
- `0027:MVV` step 3 / `0027:S4` (`env -S` unrunnable) — **dismiss (falsified).** `/usr/bin/env -S "sh -c 'echo HELLO_S'"` returns `HELLO_S`, exit 0, on this dev host. The MVV step runs as written. Secondary note: the A4 spike's portability claim ("BSD env rejects `-S` outright") is stale and should be corrected in the spike file, but nothing in the record depends on it.
- `0027:§performance-expectations` (unbounded argv, O(n²)) — **dismiss (no plausible symptom).** Per-entry, load-time, over hand-authored TOML; at 200 elements this is microseconds. B's own remedy asks only for a stated bound. Not worth a contract line.
- `0027:§decision-rationale` (premortem self-disarming) — **dismiss (subsumed).** Real as a process observation, but A gives it no independent symptom ("n/a") and its consequences are exactly the `0027:C1` / `0027:A2` items above. Folding it there avoids double-counting one defect as two.
- `0027:§metadata` (Profile `large` spanning two surfaces) — **flag to Stage 6.** B-only. Not a defect on its own, but it is the cleanest available remedy for the A6 gate violation, so it should be decided in the same breath as the `promise:`-clause split.
