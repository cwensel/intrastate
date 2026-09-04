# Verification — RDR 0027 inline-shell-detection-scope

## Phase 3a — Chain-of-Verification (independent)

Inputs: `docs/rdr/0027-inline-shell-detection-scope.md` and
`artifacts/req-list.md` only. Phase 1 test files were NOT read — probes were
derived from C1's own wording and run against the shipped binary, so a probe
passing here is independent evidence, not a restatement of an assertion.

**Result: no violations found. 0 FAIL-N entries.**

### Method

The binary was built from this worktree (`go build -o /tmp/rdr0027-probe/intrastate
./cmd/intrastate`). Each probe rewrites the `gate.rdr-lock` `command` vector of
`internal/table/testdata/neg/neg-command-shell-interpreter.toml` into a scratch
copy and runs `intrastate lint --model <scratch>`, discriminating on the
CATEGORY in the refusal detail rather than on exit status.

That discrimination matters and was the first finding of the run: the base
fixture carries an unrelated graph defect (`graph-always-present-owned`), so
every admitted form still exits non-zero. An "any error ⇒ red" oracle would
have scored all seven S4 admitted forms as violations. The reported verdicts
below are `command_shell_interpreter` vs `command_unknown_placeholder` vs
load-green, read out of the detail line.

The reported form was extracted from the detail text
(`declares the interpreter form <FORM>;`) so the tie-break REQs could be
checked on the exact words, not merely on refusal.

### Probes that would make a correct implementation visibly violate the REQ

Each row names the input designed to break the REQ, and what a
plausible-but-wrong implementation would have done with it.

| REQ | Falsifying input | Would break if… | Observed |
| --- | --- | --- | --- |
| REQ-1, REQ-42 (S1) | the ten wrappers, each wrapping `sh -c`: `env -i`, `env -u FOO`, `nice`, `timeout 5`, `xargs`, `nohup`, `setsid`, `stdbuf -o0`, `chpst`, `doas` | predicate still argv0-anchored | all ten refuse `command_shell_interpreter`, detail `sh -c` + script remediation |
| REQ-2, REQ-29 | `["nice","-n","10","timeout","5","env","-i","xargs","sh","-c","x"]` — five stacked wrappers | any wrapper-table or option-grammar walk | refuses, form `sh -c`; nothing before argv[i] read |
| REQ-44 (S3) | the eight `env` option forms: `-i`, `-u FOO`, `--unset=FOO`, `-0`, `-C DIR`, `--chdir=DIR`, `--`, and mixed `-i A=1 -- sh -c` | the deleted `env` walk was removed but not subsumed | all eight refuse, detail names `sh -c` |
| REQ-6, REQ-27 (tie-break) | `["python","sh","-c","echo"]` | tie-break preferred the NEAREST pair, or iterated the `shellInterpreters` map | reports `python -c` — the lower `i`, not the `sh -c` a reader's eye goes to |
| REQ-6 (lowest PAIR, not lowest listed basename) | `["ruby","nice","sh","-c","x"]` — `ruby` is listed at i=0 but carries only `-e`, and no `-e` follows | the scan short-circuited on the first listed basename and returned green | reports `sh -c`; outer loop correctly continued |
| REQ-6 | `["bash","nice","sh","-c","x"]` — two interpreters both carrying `-c` | tie-break on flag proximity | reports `bash -c` (lowest `i`) |
| REQ-6 / A7 | `["setup-tests","bash"]` — listed basename, no following flag | a bare trailing interpreter treated as a match | green |
| REQ-25 (`i < j` ordering) | `["-c","sh"]` and `["mytool","-c","sh"]` — flag BEFORE the interpreter | an unordered `exists i,j` scan, or an off-by-one admitting `j <= i` | both green |
| REQ-2 strict `i<j` | `["sh"]` — single element that is a listed interpreter | a `j >= i` scan would let the word match itself | green |
| REQ-2/REQ-13 (per-interpreter flag keying) | `["ruby","-c","x"]` and `["python","-e","x"]` | flags taken as the UNION over all listed interpreters instead of keyed to the interpreter at `i` | both green — `ruby` carries only `-e`, `python` only `-c` |
| REQ-13 | `["node","--eval","x"]` | only the first flag of a multi-flag entry consulted | refuses, form `node --eval` |
| REQ-3, REQ-56 (basename, exact) | `["/bin/sh","-c"]`, `["/usr/local/bin/bash","-c"]` | no basename split | refuse; forms `/bin/sh -c`, `/usr/local/bin/bash -c` — the full argv WORD is reported per REQ-6, while the LISTED match is on the basename |
| REQ-3, REQ-12, REQ-54 (OPEN, no folding) | `["python3","-c"]`, `["nodejs","-e"]`, `["perl","-e"]`, `["SH","-c"]`, `["sh.exe","-c"]` | suffix, alias, or case folding | all five green |
| REQ-8, REQ-52 (form one: one word) | `["sh -c echo hi"]`, `["env","-S","sh -c echo"]`, `["nice","sh -c echo"]`, `["sh","-cecho"]`, `["sh"," -c"]` | the check split a word on whitespace, or did prefix/substring matching on flags | all five green — no over-claim on the one-word form |
| REQ-9, REQ-10, REQ-53 (form two: stdin CHANNEL) | `["sh","-s"]`, `["sh"]`, `["sh","-es"]`, `["python","-"]`, `["node","-"]`, and under wrappers `["nice","sh","-s"]`, `["env","-i","python","-"]` | the out-of-scope line read as SHELL SPELLINGS rather than as a channel — that reading refuses `python -` and `node -`, which ARE deny-listed names | all seven green. `["python","-"]` and `["sh","-es"]` are the discriminating cases and both pass |
| REQ-11 (sanctioned wrapper file) | `["sh","./gate.sh"]`, `["sh","script.sh"]` | any argv1-is-a-script heuristic | green |
| REQ-7 (never reads the resolved binary) | `["myshell","-c","x"]` — an unlisted name that could be a real shell | binary identity / `#!` read | green |
| REQ-14, REQ-43 (S2, clause-3 coupling) | `["nice","sh","-c","cat {artifact}"]` | clause 3's exemption keyed on a different signal than the widened predicate | reports `command_shell_interpreter`, NOT `command_unknown_placeholder` — the 0025 D6/D18 masking is closed under a wrapper |
| REQ-15 (within-entry precedence unchanged) | control `["mytool","cat {artifact}"]`; and `["-c","sh","cat {artifact}"]` where the flag precedes the interpreter so no interp form exists | the exemption widened to any brace-bearing word | both correctly report `command_unknown_placeholder` — the exemption did not leak |
| REQ-15 | `["nice","sh","-c","{bogus}"]` and `["xargs","-I{}","-n1","sh","-c","x"]` — brace-bearing but whitespace-FREE under a valid interpreter form | the exemption dropped its whitespace condition | both report `command_unknown_placeholder`. Correct: clause 3's shape is `!isInterp \|\| no-whitespace ⇒ defect` and REQ-15 holds 0025:C5's precedence unchanged; REQ-14's exemption is scoped to a brace-bearing *string* (the whitespace-bearing `sh -c "cat {artifact}"` shape), which the whitespace-bearing variant `["xargs","-I {}","-n1","sh","-c","x"]` confirms — it reports `sh -c` |
| REQ-16, REQ-47 (S5), REQ-61 | `["ruby","tool.rb","-e","prod"]` | detail carried the category but not the matched words | refuses, detail names `ruby -e`, so the collided pair is diagnosable |
| REQ-26, REQ-64 (laxity deliberately retained) | `["sh","-c"]` — a "script" literally named `-c` | the flag match narrowed to `j == i+1` or to a script-position rule | refuses. The accepted false-refusal class is intact, not silently fixed |
| REQ-28 (new false-refusal class accepted, not suppressed) | `["busybox","sh","-c","x"]`, `["tini","--","/bin/bash","-c","x"]`, `["sudo","sh","-c","x"]` | the class bounded or special-cased by an argv0 allowlist | all refuse — the class is present as a stated cost, exactly as C1 requires |
| REQ-5, REQ-30, REQ-68 | grep for `wrapperTable` / `knownWrappers` / `wrappers =` across `internal/` | a wrapper table added | no match — none exists |
| REQ-22 | `func interpreterForm(argv []string) (string, bool)` at `internal/table/load.go:1183` | signature or name changed | unchanged |
| REQ-31 | source read of `load.go` for the `env` / `NAME=VALUE` skip loop | walk made conditional rather than deleted | no remnant; S3's eight forms prove subsumption rather than removal |
| REQ-48 (S6) | `git diff main...HEAD -- internal/table/category.go` filtered to `Cat*` constants and `Categories()` | membership, order, or the `command_shell_interpreter` wire string edited | only additions to the description map; `Categories()` body and the tail order are byte-identical |
| REQ-41, REQ-35 | `git diff --stat` on `command_carrier_0025_test.go` and `neg/neg-command-shell-interpreter.toml` | pre-existing probes or the fixture weakened to pass | both empty — unedited |

### The promise / claim half (REQ-17…REQ-21, REQ-25, REQ-36…REQ-39, REQ-49)

The instruction to treat an over-claiming message as a real FAIL-N was probed
directly: every claim the shipped text makes was turned into an input and run.

- **The surface ships and is reachable off the refusal path.** `intrastate lint
  --help-all` with no `--model` and no defect present exits 0 and renders
  `Load category command_shell_interpreter — what it promises:` followed by the
  full body. It is plain text, not a `respond` envelope, confirming REQ-39.
  Mirrored into `docs/cli-reference.md:842` (REQ-18, REQ-38).
- **It names BOTH out-of-scope forms in C1's channel-scoped words** (REQ-49's
  oracle, which a non-empty check would not catch): the one-word shell string
  (`env -S "sh -c …"`, and the single `"sh -c …"` element re-split by a
  receiving tool) and the STDIN channel (`sh -s`, bare `sh`, `sh -es`,
  `python -`, `node -`, with "any listed interpreter taking its code on stdin
  … is admitted, however spelled").
- **No over-claim found.** Each affirmative claim was falsified against the
  binary rather than read: "under any prefix" → the ten wrappers plus a
  five-deep stack all refuse; "never splits a word on whitespace" → the
  one-word forms are green; "never reads … the resolved binary" → `myshell -c`
  is green; "an unlisted spelling (python3, nodejs, busybox) is admitted" →
  green, and `["busybox","sh","-c"]` refusing is not a counter-example since
  `sh` is a separate listed WORD, which is the predicate's own scope; "sh
  script.sh … never a defect" → green. Every promise the text makes is one the
  predicate keeps.
- **REQ-21 (negative) holds.** No occurrence of "closes inline shell", "closes
  the wrapper class", "prevents inline shell", or equivalent in
  `docs/cli-reference.md`, `internal/table/category.go`, or
  `internal/cli/lint.go`. The text opens with "Refused:" and devotes a labelled
  block to "Out of scope, BY NAME — admitted by lint, and an interpreter may
  still run", which is the premortem's second failure closed rather than
  restated.
- **REQ-16 / REQ-20 separation holds.** The refusal detail does NOT carry the
  admitted-form disclosure (grep for "out of scope"/"admitted" over a refusal
  run returns 0); the remediation text is unchanged from 0025:C5. The promise
  lives only on the surface that is reachable without provoking a defect, which
  is the whole point of REQ-37.

### Suite state

`go test ./internal/table/ ./internal/cli/` — both ok.

### Notes

- No Phase 1 test file was opened, grepped, or executed by name; `coverage.md`
  and `deviations.md` were not read. The suite was run whole, which reports
  only pass/fail per package and discloses no assertion text.
- Scratch fixtures were written under `/tmp` and deleted; the worktree carries
  only this artifact.

---

## Phase 3b — Adversarial review (independent agent)

Written from `0027:§failure-modes` plus `artifacts/req-list.md` and the source
tree. Phase 3a's section above was not read before this section was drafted;
`deviations.md` and `coverage.md` were not read at all.

Verdict: **PASS**. Four adversarial tests added; all four pass against the
current implementation, and each was mutation-checked to confirm it is
discriminating rather than vacuous. No conformance defect was found in the
predicate or in the shipped promise text.

Added: `internal/cli/inline_shell_adversarial_0027_test.go` (4 tests, 30 cases).

### Threat model — the three most likely failure modes

The record splits its failure modes by CHANNEL: a **Visible** refusal naming
the two matched words plus the wrapper-file remediation, and a **Silent**
admission "by contract, named in C1; diagnosis is reading argv, which C1's
promise tells the reviewer to do for exactly those shapes". Both halves make
the SHIPPED TEXT the reviewer's contract, so the three modes below are the
three ways that contract can become false.

- **ADV-1 — the promise UNDER-refuses relative to its own text (Silent side).**
  The description names ten spellings as admitted (`env -S "sh -c …"`, `sh -s`,
  bare `sh`, `sh -es`, `python -`, `node -`, `sh script.sh`, plus the unlisted
  `python3`/`nodejs`/`busybox`). Predicate and text are two independently
  maintained sources of truth for one contract and **nothing joins them**: the
  table suite asserts admission against a list transcribed from the *record*,
  the cli suite asserts the *text* mentions "stdin" and "one word". Adding a
  flag to `shellInterpreters` or folding a spelling makes the shipped sentence
  a lie with every existing test green — and `0027:S4` is explicit that a
  change which starts refusing one of them "has broken the promise, not
  tightened it".
  **Caught by** `TestADV1_EveryFormTheDescriptionNamesAsAdmittedActuallyLoads`
  — the admitted forms are transcribed OUT of the rendered `--help-all` body
  (not out of the record) and each is fed to `table.Load`. The anchor check is
  half the oracle: the pair cannot be re-reconciled by deleting the promise.
  *Mutation control*: adding `-s` to `shellInterpreters["sh"]` → FAIL on
  `sh_-s` ("the description tells a reviewer … is ADMITTED, but the loader
  refuses it").

- **ADV-2 — the promise OVER-claims relative to what the predicate detects
  (Visible side).** The text commits to "under any prefix (env and its options,
  nice, timeout, xargs, doas, and **wrappers nobody enumerated**). Nothing
  before the interpreter word is read, so no wrapper table exists and none is
  consulted." That is the record's §Risks line as a live hazard: a reviewer who
  reads it and meets an unenumerated prefix loading green has been told the
  class is closed when it is not. Asserting refusal alone is too weak —
  §failure-modes makes the Visible outcome the two matched words **and** the
  remediation.
  **Caught by**
  `TestADV2_EveryFormTheDescriptionClaimsToRefuseIsRefusedWithTheMatchedWords`
  — every wrapper the text names by name, plus an unenumerated one
  (`runitwrap`), asserted on category + matched words + unchanged remediation.
  *Mutation control*: reintroducing a five-entry wrapper table gate → FAIL on
  `runitwrap_-q_bash_-c_x` ("loads clean").

- **ADV-4 — the promise never reaches the reviewer who is not at a terminal.**
  `0027:C1` promise: makes `docs/cli-reference.md` half the landing site. The
  shipped mirror assertion calls `writeCLIReference` on a live tree and greps
  the RESULT, with a comment saying why ("so the test does not depend on `make
  docs` having been run"). That is a strictly weaker oracle: it asserts the
  generator *would* produce the text, never that the committed file carries it.
  A stale or hand-softened reference leaves every other assertion green while
  the docs reader is told the wrapper class is closed. Drift is gated only in
  `make check`, outside the Go suite.
  **Caught by**
  `TestADV4_TheCommittedCLIReferenceCarriesThePromiseNotJustTheGenerator` —
  reads the checked-in bytes, requires both out-of-scope forms and the
  sanctioned wrapper-file line, and requires the committed text to equal the
  accessor's current text (whitespace-normalized) so a hand-edit that says
  something kinder than the shipped words also fails.
  *Mutation control*: rewording "from STDIN" → "from a pipe" in the committed
  file → FAIL on both the named-form check and the verbatim check.

A fourth test, `TestADV3_TheAdmittedStdinSpellingsAreNotAlsoCoveredByTheRefusedLine`,
holds the description's INTERNAL coherence: `sh -es` sits three lines below a
paragraph promising "one of that interpreter's own inline-code flags" at any
later position, and the only sentence that keeps the two paragraphs from
contradicting each other on the same argv is `C1` reads: — "never splits a word
on whitespace". The test requires that sentence to be present and each admitted
stdin spelling to load with its argv unre-split.

### Attack surfaces probed and found CONFORMANT (no test added, or covered)

Each was exercised against the built binary or through `table.Load`; all
matched `0027:C1` exactly, so no test was added where the shipped suite already
holds the line.

- **Position-free scan, over-refusal.** `["ls","python","-c"]`, `["cat","php","-r"]`
  refuse — the accepted false-positive class, already pinned by the shipped
  suite as ACCEPTED-not-suppressed. Correct per `D-selection-predicate`.
- **Position-free scan, index boundary and `i < j` ordering.** `["sh","--","-c"]`
  and `["sh","./gate.sh","-c"]` refuse (flag at any later position, as written);
  a flag BEFORE the interpreter is green; `["docker","run","-c","alpine","sh"]`
  is green because the flag precedes the interpreter word. `["myprog","--shell=sh","-c"]`
  green — `--shell=sh` is one word and basename matching is byte-exact.
- **Basename matching.** `/bin/sh`, `./sh` refuse; `sh/`, `SH`, `sh\n` green
  (no case, suffix, or alias folding). The reported form is the RAW `argv[i]`
  (`/bin/sh -c`), which is `C1` report: as written, not the basename.
- **Two-interpreter tie-break.** `["python","sh","-c","echo"]` → `python -c`;
  the scan iterates argv positions, never the map, so the form cannot depend on
  map iteration order. Covered by the shipped suite.
- **Clause-3 / clause-4 disagreement after the predicate change.** Probed the
  argv-global reach of `isInterp`: a brace-bearing whitespace element that is
  NOT the interpreter's code string (`["sh","./gate.sh","-c","junk {nope} here"]`,
  `["nice","tool","run {nope} now","sh","-c","x"]`, `["sh","-c","ok","extra {nope} arg"]`)
  is exempted from clause 3 but clause 4 refuses on the same call's `isInterp`,
  so **no green leak exists** — the exemption can only ever change which defect
  is reported, never whether one is. `["nice","sh","-c","{nope}"]` correctly
  reports `command_unknown_placeholder` (no whitespace ⇒ not exempt). The two
  clauses cannot skew because both read one `interpreterForm` call (A1).
- **Capability coverage.** The widened predicate reaches `write` and `gate`
  entries, not only `read`: `carrierDefect` is called once from
  `accessorTable`, which all three capabilities share. The shipped 0027
  predicate suite mutates only `read.state`; the behaviour is nonetheless
  correct, so this is noted rather than tested.
- **`Categories()` ordering / future append (RDR 0028's five `edit_*`).**
  `lintExtendedDesc` renders descriptions by iterating `table.Categories()`, so
  a described-but-unregistered category would silently ship no text — the
  failure `category.go`'s own 0025 comment warns of. Verified: the one
  described category IS registered, and the shipped suite's relative-order
  assertion survives a tail append. No test added: the invariant holds and a
  guard for it would duplicate the rendered-body assertions ADV-1/ADV-2/ADV-4
  already make, which fail on the same drift for the category that matters.
- **`CategoryDescription` opt-in contract (REQ-37).** No shipped test calls the
  accessor directly. Probed: described → `(text, true)`; undescribed and bogus
  → `("", false)`, never an empty string with `true`. Conformant.
- **Over-claiming language elsewhere on the CLI surface.** `grep -i "inline
  shell"` over `internal/cli/`, `internal/table/`, and `docs/cli-reference.md`
  returns exactly one hit — the unchanged 0025 remediation string in the
  refusal detail. Nothing on `flow`'s `--allow-commands` help claims shell
  containment. REQ-21 ("the docs never say 'closes inline shell'") holds.

### Suite state

`gofmt -l internal/` clean; `go vet ./internal/cli/ ./internal/table/` clean;
`go test ./internal/table/ ./internal/cli/` — both ok with the four added tests.

### Notes

- Phase 1 test files were read to avoid duplicating existing coverage, as the
  brief permits; the threat model above was fixed from `§failure-modes` before
  they were opened, and every added test attacks a seam none of them crosses
  (predicate ↔ shipped text, and generator ↔ committed artifact).
- One added-test anchor was initially too strict — it matched a phrase across
  the description's terminal hard-wrap. That was a weakness in the test, not a
  defect in the implementation, and was fixed by normalizing whitespace before
  matching rather than by relaxing the claim.
- Scratch probes were written into the worktree and deleted; nothing outside
  the added test file and this section was changed. `internal/table/load.go`
  and `docs/cli-reference.md` were mutated only to prove the added tests
  discriminate, and restored (`git status` clean but for the new test file).
