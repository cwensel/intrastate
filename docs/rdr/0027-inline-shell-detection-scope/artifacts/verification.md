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
