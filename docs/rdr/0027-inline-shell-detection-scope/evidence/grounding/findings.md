Model: claude-sonnet-5

# Grounding Step-0 codebase claim sweep — RDR 0027

Scope: claims added or edited by the refine + resolve rewrites captured in
`evidence/grounding/post-propose.diff`. The propose-time `Ground-sweep: clean
(17 anchors)` verdict and all 20 `source-anchor` edges are not re-verified
(symbol existence is out of scope per the task). This sweep targets
behavioral claims: call/read counts, counts over fixtures, spike conclusions
vs. spike data, and "no sibling exists" claims.

## Method

Read every diff hunk in `post-propose.diff`, then read the surrounding
element/section through `rdr inspect --select` for each touched claim. Read
`internal/table/load.go` (the `carrierDefect`/`interpreterForm`/
`filepathBase` region, lines ~1031–1192), `internal/table/category.go`
(`Categories()`, lines 50–104), `internal/cli/cmdbind/cmdbind.go`
(`resolveArgv0`, line 437), `internal/table/command_carrier_0025_test.go`
(`TestReq74…`, `TestReq83…`), all 6 `command`-carrying fixtures under
`internal/table/testdata/neg/`, all three spike files under
`evidence/spikes/`, and `evidence/research/stage4-scope-wording.md`. Built
and ran two throwaway Go harnesses (not committed) reproducing
`interpreterForm` and the widened C1 predicate to independently recompute
divergences over the `env`-flag forms `intrastate#q2q1` enumerates.

## Confirmed (no finding)

- **A1** — call count and read count. `interpreterForm(argv)` called exactly
  once at `load.go:1095`; `isInterp` read exactly twice, at `load.go:1114`
  (clause 3 exemption, `!isInterp || …`) and `load.go:1124` (clause 4
  refusal — note: RDR text says `:1124`, the `if isInterp {` line is at
  1124 and the `fail(...)` call spans 1125–1128; the RDR's line number for
  the read itself, 1124, is exact). Repo-wide grep for `interpreterForm`
  and `isInterp` returns only `internal/table/load.go`, no other file.
  CONFIRMED exactly as claimed.
- **Existing Infrastructure Audit, row 1** — "interpreter must sit at argv0
  after an `env` walk that stops on any option flag." Verified by
  reproducing the walk: `env -i sh -c …` and `env -u FOO sh -c …` both stop
  the walk at `-i`/`-u` (non-`NAME=VALUE` tokens) and return
  `isInterp=false` under the CURRENT predicate. CONFIRMED.
- **Existing Infrastructure Audit, row 4** — `resolveArgv0`
  (`internal/cli/cmdbind/cmdbind.go:437`) resolves location (abs / bare via
  `exec.LookPath` / relative via `baseDir` join), never identity. CONFIRMED
  by reading the function body.
- **A2 spike population** — 6 `command`-carrying fixtures exist under
  `internal/table/testdata/neg/*.toml` (grep confirms exactly 6 of 76 total
  `neg-*.toml` files declare `command`); no fixture under
  `internal/table/testdata/pos/` declares `command` (grep empty).
  `TestReq83_TheInterpreterSetIsDenyListedNotClosed` exists at
  `command_carrier_0025_test.go:627` with `command = ["perl","-e","print
  1"]` and asserts `table.Load` returns no error (green). Independently
  recomputed the spike's own `widenedInterpreterForm` against
  `["perl","-e","print 1"]`: stays `(false,"")` under both predicates.
  CONFIRMED — `["perl","-e","print 1"]` stays green, as A2/A3's Evidence
  states.
- **A2 zero-divergence finding** — re-ran the spike's own recorded table by
  hand for a subset (the 11 green rows) and confirm none flips to `true`
  under the widened scan. CONFIRMED against the spike's own pasted output.
- **A3 wrapper argv-word transparency** — re-derived `widenedInterpreterForm`
  independently and ran it over the 9 `env`-option forms
  `intrastate#q2q1` enumerates (`-i`, `-u FOO`, `--unset=FOO`, `-0`, `-C
  DIR`, `--chdir=DIR`, `--`, `A=1`, `-i A=1 --`): all 9 return
  `(true,"sh -c")`, matching the record's Testing Strategy scenario 3 and
  research finding Q3. CONFIRMED.
- **A4 reachable-forms claim** — spike's own harness output for `sh <
  script.sh` (exit 127, `<` treated as a literal filename, script never
  read) and `sh ./gate.sh` (runs the declared file, not inline code) is
  internally consistent with A4's text. CONFIRMED against the spike's own
  data (no independent re-run of the exec harness; the spike's transcript is
  self-consistent and its adversarial round finds no third reachable form).
- **S6 / Testing Strategy scenario 6** — `table.Categories()`
  (`internal/table/category.go:66`) lists `CatCommandShellInterpreter` at
  `category.go:60` (declaration) and appends it at tail position 4 of the 6
  command-carrier categories inside `Categories()` at line 103. C1's change
  is entirely inside `load.go`; `category.go` is untouched by the record's
  proposed edit. CONFIRMED — membership and order match, and the edit-scope
  claim ("the edit touches `load.go` only") is accurate.
- **C1 `base` claim** — "no suffix or alias folding: `python3`, `nodejs`,
  `busybox` stay unlisted spellings." Read `shellInterpreters` map
  (`load.go:1035-1047`): contains exactly `sh, bash, dash, ksh, zsh, csh,
  tcsh, python, ruby, node, php`. None of `python3`, `nodejs`, `busybox` is
  a key. CONFIRMED.
- **Approach section, "the deny-list already carries python, ruby, node and
  php"** — confirmed by the same map read above.
- **Inverse check (task item 7)** — searched for a sibling argv-scanning
  helper, wrapper-aware walk, or interpreter-set membership test elsewhere
  in the repo. `grep -rn shellInterpreters --include=*.go .` returns only
  the declaration (`load.go:1035`) and its one use site (`load.go:1174`,
  inside `interpreterForm`). `grep -rln wrapper --include=*.go` returns 4
  files, none of which implement an argv scanner or interpreter classifier
  (comments about generic "wrapper scripts" in `cmdbind.go`, `clierr.go`,
  `flow_input.go`, `help_all.go`). NOT FOUND — no sibling exists; the
  record's "no wrapper or command-position classifier exists anywhere in
  the audited paths" (A1 Evidence, Existing Infrastructure Audit row 4) is
  accurate.
- **Q5 / C1 `promise:` line** — "No such description exists in the shipped
  code today." Grepped for `command_shell_interpreter` outside test files:
  only `load.go` (refusal detail, `load.go:1125-1128`) and `category.go`
  (constant declaration). No separate user-facing description string exists
  anywhere else (docs, CLI help). CONFIRMED.
- **A5 "no kata tracks the successor"** — corroborated by the record's own
  cited research file `stage4-scope-wording.md` §Q4, which independently
  ran `kata list --status open` and found no stdin/envelope item. Not
  re-run here (kata state is external, time-sensitive, and the record
  already dates the check 2026-09-03, i.e. this Grounding sweep's own run
  date); treated as consistent, not independently re-verified.

## Findings

### Finding 1 — A2's "32 argv vectors" overstates distinct coverage; it is 32 test-case rows, several sharing identical argv

**Claim** (0027:A2 Evidence): "both predicates computed over all 32 `command`
argv vectors in `internal/table/testdata/neg/*.toml` and
`command_carrier_0025_test.go`."

**What I found instead**: The spike's own population breakdown
(`evidence/spikes/a2-widened-predicate-no-new-refusals.md`, "Population
swept") is 6 rows from the neg fixtures + 26 rows from the test file = 32
rows — but these are per-*fixture*/per-*test-case* rows, not 32 distinct
argv vectors. Two of the 6 neg-fixture rows are the literal same vector:
`internal/table/testdata/neg/neg-command-and-path-conflict.toml:96` and
`internal/table/testdata/neg/neg-command-env-conflict.toml:95` both declare
`command = ["rdr-gate", "{artifact}"]`. Independently, `grep -oP 'command =
\[.*\]' internal/table/command_carrier_0025_test.go | sort -u | wc -l`
returns 27 distinct argv literals in the test file (against 61 total
`command = [` occurrences — most test cases reuse `["reader"]`,
`["writer"]`, `["gater"]`, etc.), one more than the spike's 26 rows drawn
from that file.

**Verdict**: REFUTED (as literally stated) / not load-bearing. "32 argv
vectors" is imprecise — it is 32 curated test-case instances (with at least
one exact duplicate and a handful of the file's 27 distinct shapes not
individually re-listed as separate rows), not 32 distinct vectors. The
undercount does not touch the spike's conclusion: rechecking by hand shows
the missing/duplicated rows are unrelated to interpreter forms (`["reader",
...]` variants), and the zero-divergence finding holds regardless. This is
a wording-precision nit in the Evidence field, not a refutation of A2 or of
the "zero new refusals" verdict.

### Finding 2 — "11 green, 21 already-negative" is consistent with the spike's own transcript, but the neg-fixture count folds a duplicate

**Claim** (0027:A2 Evidence): "(11 green, 21 already-negative; no positive
fixture declares `command`)."

**What I found instead**: 32 total − 11 green = 21 already-negative,
arithmetically consistent with the spike's transcript (`total cases: 32,
divergences: 0`, and counting the `green=false` rows in the pasted output
gives 21). This claim is arithmetically self-consistent with the spike file
and is not independently wrong — it inherits Finding 1's imprecision (one
of the 21 "already-negative" rows is a duplicate argv vector counted twice
across two different fixture files) but the arithmetic and the "zero new
refusals" conclusion both still hold.

**Verdict**: CONFIRMED (consistent with spike data) with the same caveat as
Finding 1 — not a separate refutation.

## Spans widened past the listed scope, and why

- Read `internal/table/category.go:50-104` in full (not just the two cited
  line numbers) to confirm S6's membership-and-order claim requires seeing
  the whole `Categories()` body, not just the two anchor lines.
- Read all 6 `command`-carrying `neg/*.toml` fixtures in full (not just
  grepped the `command =` line) to catch the cross-fixture duplicate found
  in Finding 1 — grepping the line alone would have missed that two
  different files declare the identical vector.
- Ran two independent throwaway Go harnesses reproducing `interpreterForm`
  and the widened predicate rather than trusting the spike files' pasted
  transcripts alone, specifically to re-verify the `intrastate#q2q1`
  9-form claim (Testing Strategy scenario 3, Approach section) since that
  claim is asserted in the diff-touched prose but its own worked table
  lives only in `evidence/research/stage4-scope-wording.md`, a file the
  task did not name as spike evidence.
- Spot-checked (not exhaustively verified) the "27 Draft/Final/Implemented
  peers" count in `0027:JC1` (joint-check element) even though `JC1` was
  not touched by `post-propose.diff` — flagging as NOT-VERIFIED below
  rather than silently skipping it, since the task's claim list mentioned
  counts generally. A rough grep (`Status.*: (Draft|Final|Implemented)`
  across `docs/rdr/*.md`) returns 28, not 27, but this is an unreliable
  proxy (front-matter format varies) and JC1 is out of the diff scope, so
  this is reported as NOT-VERIFIED, not REFUTED.

## Not verified

- **0027:JC1** "the 27 Draft/Final/Implemented peers" — out of diff scope
  (JC1 unchanged by `post-propose.diff`); a rough grep returned 28, not 27,
  but the grep is not a reliable proxy for the record's own peer-count
  methodology (front-matter format varies across RDRs, and the record's
  count may exclude records grepped differently, e.g. by directory
  structure vs. flat file, or as of a different date). NOT-VERIFIED —
  flagging for awareness, not asserting REFUTED.
- **gh `pkg/cmd/alias/set/set.go::NewCmdSet`** and its `--shell` help text
  quote — external repo, not vendored locally, and this prose is unchanged
  by `post-propose.diff` (already covered by the propose-time 17-anchor
  ground-sweep). Not re-fetched.
- **Go `os/exec` package doc quote** ("intentionally does not invoke the
  system shell … pipelines, or redirections") — unchanged by the diff, not
  re-verified against the exact stdlib doc text in this pass; the behavior
  it describes (no shell, no redirection) is independently confirmed by
  A4's own harness (`sh < script.sh` unreachable, `<` treated as a literal
  argv word).
