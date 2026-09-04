# REQ list — RDR 0027 inline-shell-detection-scope

Phase 0 spec audit. Source: `docs/rdr/0027-inline-shell-detection-scope.md`
(497 lines; one `normative` fence — `0027:C1` — plus one MVV, 7 Critical
Assumptions, 2 Load-Bearing Decisions, 4 mini-checks, 7 Testing Strategy
scenarios, 3 Alternatives, 1 Cross-Cutting Concerns block).

Element ids carried where a REQ derives from a labelled contract. Quotes are
exact bytes from the record, reflowed only where a clause spans lines in the
source; no wording is changed. The record ships a **single** contract, `C1`,
and it explicitly **OVERRIDES** two lines of `0025:C5` (see REQ-0).

This record has **two halves**, and the REQ set covers both:

- **the predicate half** — what `interpreterForm` computes and what
  `carrierDefect` reports (REQ-1 … REQ-16, REQ-26 … REQ-31);
- **the promise/claim half** — the *wording* verification is held to, shipped
  as a user-facing description on `--help-all` (REQ-17 … REQ-21, REQ-25,
  REQ-38). A predicate test cannot reach this half; S7 is its only oracle.

A tooling note: no RDR projector binary resolves on this host
(`/Users/cwensel/.local/bin/intrastate` is the flow CLI, which has no
`inspect` subcommand; `rdr` is not on `$PATH`). Quotes below were taken by
reading the record whole rather than by `--select`, so they are transcribed
bytes checked against the source lines cited, not tool-copied bytes. Line
cites are given per section for re-verification.

---

## Override / scope frame (`0027:§metadata`)

- [REQ-0] "**Overrides**: 0025:C5's `command_shell_interpreter` line and its
  `interpreter set` line — C1 here is the successor text: the predicate widens
  from argv0-after-an-`env`-walk to any argv position under any prefix, and the
  claim names what the check reads and the two forms it does not" —
  (§metadata, line 20). Consequence for the build: `0025:REQ-74` (the argv0
  wording) and `0025:REQ-83` (the `interpreter set` line) are SUPERSEDED by
  `0027:C1`; every other 0025:C5 line — the clause map, within-entry
  precedence, registration, the other five categories — is unchanged and
  **must not** be restated or edited here.
- [REQ-0a] "The other five C5 categories, C5's registration and precedence
  lines are unchanged and not restated." — (§normative-contracts preamble,
  line 111) — NEGATIVE REQ.
- [REQ-0b] "0025 is not edited; `load.go`'s C5 comments re-cite this clause."
  — (§normative-contracts preamble, line 111). The 0025 record file is not
  touched; the in-source C5 comments are re-pointed at `0027:C1`.

## C1 — the predicate (`0027:C1`, §normative-contracts, lines 115–124)

- [REQ-1] "command_shell_interpreter     # a listed interpreter word followed,
  at ANY later argv position, by one of that interpreter's inline-code flags —
  under any prefix (env and its options, nice, timeout, xargs, doas, …); no
  opt-in in v1" — (0027:C1, §normative-contracts) — the category's wire string
  is unchanged; the comment is the successor text for 0025:REQ-74.
- [REQ-2] "predicate:  refuse iff there exist i < j with base(argv[i]) a listed
  interpreter and argv[j] one of its listed flags" — (0027:C1,
  §normative-contracts) — the quantifier is `exists i < j`; refusal is
  **iff**, so the predicate is total: no other condition refuses and no
  condition exempts.
- [REQ-3] "base = the text after the last `/`, matched exactly (no suffix or
  alias folding: `python3`, `nodejs`, `busybox` stay unlisted spellings under
  the OPEN rule)" — (0027:C1, §normative-contracts).
- [REQ-4] "Nothing before argv[i] is read" — (0027:C1, §normative-contracts) —
  NEGATIVE REQ: no prefix, wrapper, or option-grammar inspection of any kind.
- [REQ-5] "no wrapper table exists and none may be added" — (0027:C1,
  §normative-contracts) — NEGATIVE REQ, normative and forward-binding.
- [REQ-6] "The reported form is argv[i] + \" \" + argv[j] for the lowest i,
  then the lowest j" — (0027:C1, §normative-contracts) — the tie-break is part
  of the contract, not an implementation choice.
- [REQ-7] "reads:      argv WORDS only. The check never splits a word on
  whitespace, never reads stdin, files, PATH, or the resolved binary" —
  (0027:C1, §normative-contracts) — NEGATIVE REQ on four named channels.

## C1 — out of scope, BY NAME (`0027:C1`, line 119)

- [REQ-8] "out of scope, BY NAME (admitted by lint; an interpreter may run): a
  shell string carried in ONE word (`env -S \"sh -c …\"`, or a single `\"sh -c
  …\"` element handed to a tool that re-splits it)" — (0027:C1,
  §normative-contracts) — **form one**: admitted, green, no lint event minted.
- [REQ-9] "an interpreter that reads its script from STDIN (`sh -s`, bare `sh`,
  `sh -es`, `python -`, `node -`) — owned by the charted `stdin = \"none\" |
  \"envelope\"` successor, whose remit is what a command entry CONSUMES and so
  covers the non-shell spellings too" — (0027:C1, §normative-contracts) —
  **form two**: admitted, green.
- [REQ-10] "The channel is the scope: any listed interpreter taking its code on
  stdin rather than as a later argv word is admitted, however spelled." —
  (0027:C1, §normative-contracts) — the out-of-scope line is scoped by CHANNEL,
  **not** by shell spellings; this is what makes `["python","-"]` green
  (`python` IS deny-listed) and is the S4 discriminator.
- [REQ-11] "`sh script.sh` is the sanctioned wrapper-file form and never a
  defect" — (0027:C1, §normative-contracts).
- [REQ-12] "interpreter set: OPEN (deny-listed, not closed), unchanged from
  0025:C5 — an unlisted spelling is admitted; the list grows by amendment of
  THIS clause." — (0027:C1, §normative-contracts) — the amendment site moves
  from 0025:C5 to 0027:C1; the OPEN rule itself is unchanged.
- [REQ-13] "The membership itself is not restated here: it is
  `internal/table/load.go::shellInterpreters`, one map holding each listed name
  with its own inline-code flags (`sh`/`python` → `-c`, `ruby` → `-e`)" —
  (0027:C1, §normative-contracts) — NEGATIVE REQ: the map's membership is not a
  0027 contract and is not edited; the per-interpreter flag list is what S4's
  `python -` row and S5's `ruby -e` row rest on.

## C1 — clause-3 coupling and report (`0027:C1`, lines 121–122)

- [REQ-14] "clause 3 coupling: 0025:C5 clause 3's whitespace exemption keys on
  this predicate, so a brace-bearing string under a wrapper reports the
  interpreter form, never command_unknown_placeholder" — (0027:C1,
  §normative-contracts).
- [REQ-15] "within-entry precedence is 0025:C5's, unchanged" — (0027:C1,
  §normative-contracts) — NEGATIVE REQ.
- [REQ-16] "report:     category string, remediation (\"inline shell is not a
  declared command; put it in a script and declare the script as argv0\") and
  position in `table.Categories()` are unchanged; the detail additionally names
  the two matched words" — (0027:C1, §normative-contracts) — three things
  unchanged, one thing added.

## C1 — the promise clause (`0027:C1`, line 123)

- [REQ-17] "promise:    what a reviewer may rely on is the predicate line and
  the two out-of-scope forms named above, carried in a description Phase 3
  ships on the `--help-all` surface" — (0027:C1, §normative-contracts).
- [REQ-18] "(A6, Verified: `internal/cli/lint.go::lintExtendedDesc` already
  renders a closed taxonomy's text there, outside the refusal path, mirrored
  into `docs/cli-reference.md`)" — (0027:C1, §normative-contracts) — the
  landing site is settled, not open.
- [REQ-19] "THESE words are the text it ships." — (0027:C1,
  §normative-contracts) — the description carries C1's own wording, not a
  paraphrase.
- [REQ-20] "The surface carries no admitted-form disclosure until Phase 3 lands
  it — S7 is that assertion, and until it passes the predicate line alone is
  what holds" — (0027:C1, §normative-contracts).
- [REQ-21] "the lint's description carries it [C1's out-of-scope line]
  verbatim; the docs never say \"closes inline shell\"" —
  (§risks-and-mitigations, line 322) — NEGATIVE REQ on the shipped wording.

## Determinacy (`0027:C1` tail, line 126)

- [REQ-22] "Determinacy: n/a — no new signature, type, or API surface. C1
  replaces `interpreterForm`'s body at its existing name and signature" —
  (§normative-contracts) — NEGATIVE REQ: `interpreterForm`'s name and
  `(string, bool)` signature do not change.
- [REQ-23] "the algorithm is pinned end to end in the clause: the quantifier
  (`exists i < j`), basename matching (text after the last `/`, exact, no
  suffix or alias folding), the tie-break (lowest `i`, then lowest `j`), the
  reported form (`argv[i] + \" \" + argv[j]`), and the read-scope (argv words
  only, never splitting on whitespace)." — (§normative-contracts).

## Load-Bearing Decisions (`0027:D-naming`, `0027:D-selection-predicate`, lines 130–131)

- [REQ-24] "**Naming** — the category stays `command_shell_interpreter`;
  rejected: `command_inline_shell` (a rename re-orders nothing a reviewer sees
  and breaks 0025's tests and `Categories()` wire order for no gain)." —
  (0027:D-naming, §load-bearing-decisions) — NEGATIVE REQ.
- [REQ-25] "when several pairs qualify, the lowest interpreter index and then
  the lowest flag index are reported; order is part of the predicate (a flag
  word before the interpreter word never matches)" —
  (0027:D-selection-predicate, §load-bearing-decisions).
- [REQ-26] "the flag is matched at any later position, as today, so the
  already-accepted false-refusal class (a script argument that is itself `-c` /
  `-e`) is unchanged and is diagnosable from the named words" —
  (0027:D-selection-predicate, §load-bearing-decisions) — the flag-anywhere
  laxity is DELIBERATELY retained.
- [REQ-27] "Two listed interpreters in one argv resolve by the same rule and
  the outer word wins: `[\"python\",\"sh\",\"-c\",\"echo\"]` reports `python
  -c`, not the `sh -c` a reader's eye goes to, because `python` and `sh` both
  carry `-c` (`shellInterpreters`) and `python` holds the lower i. Refusal is
  unaffected — either pair refuses — so this is message attribution only" —
  (0027:D-selection-predicate, §load-bearing-decisions) — the worked
  two-interpreter case is normative for the reported detail.
- [REQ-28] "Position-freedom additionally admits a NEW false-refusal class — a
  listed basename as a data argument under an unlisted argv0 — measured at A7
  … so the class is accepted as a stated cost" —
  (0027:D-selection-predicate, §load-bearing-decisions) — NEGATIVE REQ: the
  class is NOT to be bounded, special-cased, or suppressed.

## Technical Design (§technical-design, lines 103–107)

- [REQ-29] "`interpreterForm` becomes a two-index scan: for the lowest `i`
  whose `filepathBase(argv[i])` is a key of `shellInterpreters`, and the lowest
  `j > i` with `argv[j]` in that key's flags, return `argv[i] + \" \" +
  argv[j]`. Everything before `i` is unread. `shellInterpreters` and
  `filepathBase` are unchanged." — (§technical-design).
- [REQ-30] "Not modelled, on purpose: wrapper option grammars, `--`
  conventions, `env -S` splitting, PATH or binary identity." —
  (§technical-design) — NEGATIVE REQ, five named non-features.
- [REQ-31] "The `env`-chain walk in `internal/table/load.go::interpreterForm`
  is deleted: this record DECIDES that C1 subsumes it" — (§approach, line 93).
  "either way the `env` walk is removed here" — (§prerequisites, line 335).
- [REQ-33] "the edit is confined to `load.go`" / "`category.go`; C1's edit is
  confined to `load.go`" — (§mini-check-authority, line 157) — NEGATIVE REQ:
  `internal/table/category.go` is not edited by the predicate work.
  (Phase 3's description accessor is a separate, additive surface — see
  ASSUMPTION on REQ-38.)

## MVV (`0027:MVV`, §minimum-viable-validation, lines 339–344)

- [REQ-MVV] `0027:MVV` — §minimum-viable-validation, lines 339–344. Five
  numbered steps plus the end-state line, quoted in full:

  1. "Take the 0025 command-carrier model with a green `command` read binding;
     `intrastate lint` passes."
  2. "Swap in one mutant per named wrapper — `[\"env\",\"-i\",\"sh\",\"-c\",\"…\"]`,
     `[\"env\",\"-u\",\"FOO\",\"sh\",\"-c\",\"…\"]`, `[\"nice\",\"sh\",\"-c\",\"…\"]`,
     `[\"timeout\",\"5\",\"sh\",\"-c\",\"…\"]`, `[\"xargs\",\"sh\",\"-c\",\"…\"]`,
     `[\"doas\",\"sh\",\"-c\",\"…\"]` — and one carrying `{artifact}` inside the
     string under `nice`; each refuses with `command_shell_interpreter`, the
     detail naming `sh -c` and the script remediation; the `{artifact}` mutant
     does NOT report `command_unknown_placeholder`."
  3. "Swap in the admitted forms `[\"sh\",\"./gate.sh\"]`,
     `[\"env\",\"-S\",\"sh -c echo\"]`, `[\"sh\",\"-s\"]`, and the non-shell stdin
     spelling `[\"python\",\"-\"]`; each lints green AND loads with its argv
     unchanged (green alone is an absence-of-error oracle a no-op passes;
     `oracle` mini-check). The first two are shown, by running the model under
     `--allow-commands`, to behave as C1 states (the wrapper file runs; the
     `-S` string runs a shell — on a host whose `env` supports `-S`, GNU env,
     since darwin's stock BSD env rejects it; A4). `[\"python\",\"-\"]` is the
     probe that C1's out-of-scope line is stated over the stdin CHANNEL and not
     over shell spellings: it must lint green and be documented as admitted
     (`python` is a `shellInterpreters` member, so a spellings-scoped reading
     would refuse it)."
  4. "The original 0025 REQ-74 probes and the negative fixture still refuse
     with the same category, and `table.Categories()` order is unchanged."
  5. "The description Phase 3 ships is read without provoking a refusal and
     names both out-of-scope forms in C1's words (S7) — the `promise:` clause's
     only test; steps 1–4 all pass with no description at all."

  End-state: "every wrapper mutant red with the right defect, every admitted
  form green and documented, the promise text shipped and asserted, no existing
  probe changed."

## Implementation phases (§implementation-plan, lines 377–387)

- [REQ-34] "Phase 1: Predicate — Replace `interpreterForm`'s body with the
  position-free scan and delete the `env` walk; extend the refusal detail with
  the matched words; leave `shellInterpreters`, `filepathBase`, category
  constants and `Categories()` untouched." — (§implementation-plan).
- [REQ-35] "Phase 2: Probes — Add the wrapper mutants and the admitted-form
  probes beside the REQ-74 test; keep the four existing probes as-is." —
  (§implementation-plan) — NEGATIVE REQ on the four existing probes.
- [REQ-36] "Phase 3: Promise wording — Update the C5 comments in `load.go` to
  cite this record's C1, and **create** the user-facing description of
  `command_shell_interpreter` stating the predicate and the two out-of-scope
  forms in C1's words." — (§implementation-plan).
- [REQ-37] "No description surface exists for `table.Category` today: it is a
  bare string and `Categories()` returns identifiers only … so this phase adds
  the surface as well as the text — the description must be reachable without
  provoking a refusal, which the detail string (`load.go::carrierDefect`) is
  not." — (§implementation-plan).
- [REQ-38] "Where it lands is settled (A6, Verified): the `--help-all`
  extended-help body, following the pattern `internal/cli/lint.go::lintExtendedDesc`
  already uses for `graphlint`'s closed taxonomy — an exported accessor pairing
  the category with its text, consumed by that body, registered via
  `withExtendedHelp` and mirrored into `docs/cli-reference.md` by
  `internal/cli/docs.go::runDocs`." — (§implementation-plan).
- [REQ-39] "`internal/cli/help_all.go` keeps that stream off the respond
  gateway, so it renders with no defect present." — (§implementation-plan) —
  the description must NOT route through the respond gateway.
- [REQ-40] "intrastate#q2q1's disposition recorded (landed first, or closed as
  subsumed) — either way the `env` walk is removed here" — (§prerequisites).

## Testing Strategy scenarios (`0027:S1`–`S7`, §testing-strategy, lines 393–462)

- [REQ-41] "Unit tests beside the existing C5 probes in
  `internal/table/command_carrier_0025_test.go`, driven through `table.Load` so
  each scenario asserts the CATEGORY a model author actually sees. Done = every
  scenario below green, the four REQ-74 probes and
  `neg/neg-command-shell-interpreter.toml` unedited and still green, and `go
  test ./internal/table/` clean." — (§testing-strategy) — NEGATIVE REQ: the
  fixture and the four probes are **unedited**.
- [REQ-42] S1 — "One mutant per wrapper the Problem Statement names — `env -i`,
  `env -u FOO`, `nice`, `timeout 5`, `xargs`, `nohup`, `setsid`, `stdbuf -o0`,
  `chpst`, `doas` — each wrapping `sh -c`. **Expected**: each refuses
  `command_shell_interpreter`, detail naming the matched words `sh -c` and the
  script remediation … the wrapper binary need not exist on the test host,
  since the predicate reads argv words only." — (0027:S1, §testing-strategy) —
  ten named wrappers (the MVV's step 2 names six; S1 is the superset).
- [REQ-43] S2 — "`[\"nice\",\"sh\",\"-c\",\"cat {artifact}\"]` — a brace-bearing
  command string under a wrapper. **Expected**: `command_shell_interpreter`,
  NOT `command_unknown_placeholder`." — (0027:S2, §testing-strategy). Oracle
  row: "assert the category *string*, not merely that load failed"
  (§mini-check-oracle, line 353).
- [REQ-44] S3 — "The `env`-option forms intrastate#q2q1 enumerates — `-i`, `-u
  FOO`, `--unset=FOO`, `-0`, `-C DIR`, `--chdir=DIR`, `--`, and a mixed `-i A=1
  -- sh -c` chain. **Expected**: all refuse `command_shell_interpreter`,
  **detail naming the matched words `sh -c`** as in scenario 1 … Proves the
  deleted `env` walk is subsumed rather than merely removed" — (0027:S3,
  §testing-strategy).
- [REQ-45] S4 — "The admitted forms — `[\"sh\",\"./gate.sh\"]`,
  `[\"env\",\"-S\",\"sh -c echo\"]`, `[\"sh\",\"-s\"]`, `[\"sh\"]`,
  `[\"sh\",\"-es\"]`, `[\"python\",\"-\"]`, `[\"node\",\"-\"]`. **Expected**:
  every one lints GREEN **and the binding loads with its argv unchanged** — not
  merely that no error was returned, which a no-op predicate also satisfies …
  `[\"python\",\"-\"]` and `[\"sh\",\"-es\"]` are the discriminating cases" —
  (0027:S4, §testing-strategy) — seven admitted forms, argv-unchanged oracle.
- [REQ-46] S4 forward-binding clause — "These are C1's out-of-scope line as a
  test: a future predicate change that starts refusing one of them has broken
  the promise, not tightened it." — (0027:S4, §testing-strategy) — NEGATIVE
  REQ.
- [REQ-47] S5 — "The regression set for the accepted false-refusal class —
  `[\"ruby\",\"tool.rb\",\"-e\",\"prod\"]` (a script whose own argument is
  `-e`). **Expected**: refuses, and the detail names `ruby -e` — the two
  matched words — so the author can see which pair collided." — (0027:S5,
  §testing-strategy). Oracle row: "assert the detail *text* contains both
  words; a bare category assertion passes even with today's message"
  (§mini-check-oracle, line 356).
- [REQ-48] S6 — "`table.Categories()` membership and order, and the
  `command_shell_interpreter` wire string. **Expected**: unchanged. The edit
  touches `load.go` only; the category constant and its position live in
  `category.go:60`/`:103` and are not edited" — (0027:S6, §testing-strategy).
  Oracle row: "golden assertion over the full ordered slice, not a `Contains`
  check" (§mini-check-oracle, line 357). See DEVIATION D1 (pre-seeded): the
  golden must be authored as RELATIVE order, since 0025:REQ-79 makes the list's
  size a non-contract and 0028:C1.4 appends five members.
- [REQ-49] S7 — "The description Phase 3 ships for `command_shell_interpreter`,
  read from the `--help-all` extended-help body (A6) with no model loaded and
  no defect present, so nothing routes through the refusal gateway.
  **Expected**: it exists, and its text names both out-of-scope forms — the
  one-word shell string and the stdin-fed interpreter — in C1's channel-scoped
  words." — (0027:S7, §testing-strategy). Oracle row: "assert the description
  text contains BOTH admitted forms (the one-word string and the stdin
  channel), not that a description is non-empty: a placeholder string passes a
  non-empty check" (§mini-check-oracle, line 358).
- [REQ-50] "scenarios 1–5 assert what the predicate does, and none of them
  would fail if the description were missing, stale, or silent about the
  admitted forms." — (0027:S7, §testing-strategy) — S7 is the sole oracle for
  the promise half; it must not be folded into a predicate assertion.

## Mini-check `disposition` (§mini-check-disposition, lines 163–171) — the loud/silent table

- [REQ-51] "interpreter word + its inline-code flag at any later argv position
  (`[\"nice\",\"sh\",\"-c\",\"…\"]`) | refuse — load-time defect |
  `command_shell_interpreter`, detail names the two matched words + the script
  remediation | none (load fails) | **loud**" — (§mini-check-disposition).
- [REQ-52] "shell string inside ONE word (`[\"env\",\"-S\",\"sh -c echo\"]`) |
  green | none | binding loads; a shell may run | **silent by contract** …
  no lint event is minted" — (§mini-check-disposition) — NEGATIVE REQ: no
  advisory, no warning, no event.
- [REQ-53] "interpreter reading its script from stdin (`[\"sh\",\"-s\"]`,
  `[\"sh\"]`, `[\"sh\",\"-es\"]`, `[\"python\",\"-\"]`, `[\"node\",\"-\"]`) |
  green … **silent by contract** — named by CHANNEL, not by shell spelling" —
  (§mini-check-disposition) — bare `["sh"]` is explicitly green.
- [REQ-54] "unlisted spelling (`[\"python3\",\"-c\",\"…\"]`,
  `[\"perl\",\"-e\",\"…\"]`) | green … **silent by contract** — OPEN
  deny-list; the list grows only by amending C1" — (§mini-check-disposition).

## Cross-Cutting Concerns (`0027:G-cross-cutting`, §cross-cutting-concerns, lines 480–487)

- [REQ-55] "`C1` states `no opt-in in v1`, which is deliberate: an opt-in flag
  would make the promise conditional and defeat the visibility the check exists
  for" … "No deprecation window, no compatibility flag, no staged rollout — the
  refusal set moves once, at the version that ships `C1`." —
  (0027:G-cross-cutting) — NEGATIVE REQ: no opt-in flag, no migration flag.
- [REQ-56] "Basename matching is `base = the text after the last /`, matched
  **exactly** — no suffix folding, no alias folding, no case folding … the
  check does no Unicode normalization, no case mapping, and no
  locale-dependent comparison; a listed name is a byte-exact match on the
  post-slash segment." — (0027:G-cross-cutting) — NEGATIVE REQ.
- [REQ-57] "the reported form cannot depend on map iteration order over
  `internal/table/load.go::shellInterpreters`" — (0027:G-cross-cutting) — the
  scan iterates argv positions, never the map.
- [REQ-58] "`interpreterForm` is a pure function over a fixed argv slice with
  no shared state, invoked from the single call site … and this record adds no
  goroutine, no cache and no mutable package-level state." —
  (0027:G-cross-cutting) — NEGATIVE REQ.
- [REQ-59] "This record does not claim byte-identical output,
  content-addressed identity, or replay-stable hashes, so the hash/pre-image
  checklist does not apply." — (0027:G-cross-cutting) — NEGATIVE REQ.
- [REQ-60] "`C1` names the stdin channel out of scope by name and routes it to
  the charted `stdin = \"none\" | \"envelope\"` successor. That routing is an
  **ownership claim, not a closure** … until one does, the forms are admitted
  and documented as such" — (0027:G-cross-cutting) — NEGATIVE REQ: **no code
  lands for the stdin axis in 0027**; no withholding point is added at
  `internal/cli/cmdbind/cmdbind.go:207`.

## Failure Modes (§failure-modes, lines 326–328) — observable behaviour

- [REQ-61] "Visible: a refused wrapper form names `sh -c` (the matched words)
  and the wrapper-file remediation; a false refusal names the two words that
  matched, so the author sees the script argument that collided." —
  (§failure-modes).
- [REQ-62] "Silent: a one-word shell string or a stdin-fed interpreter (`sh
  -s`, bare `sh`, `sh -es`, `python -`, `node -`) lints green — by contract,
  named in C1" — (§failure-modes).
- [REQ-63] "Recovery: none needed for the predicate (no state)" —
  (§failure-modes) — NEGATIVE REQ.

## Consequences (§consequences, lines 311–315) — accepted costs, stated not fixed

- [REQ-64] "Negative: the flag-anywhere laxity stays, so a script argument that
  is itself `-c`/`-e` after an interpreter word is refused (already true
  today); the message now names the words so the fix is visible." —
  (§consequences) — NEGATIVE REQ: do not narrow the flag match.
- [REQ-65] "Negative: `env -S \"sh -c …\"` and the stdin-fed forms (`sh -s`,
  bare `sh`, `sh -es`, `python -`, `node -`) are admitted and documented as
  such" — (§consequences) — "documented as such" is Phase 3 / S7's obligation,
  not a lint event.

## Performance (§performance-expectations, lines 466–470)

- [REQ-66] "the predicate is O(len(argv)²) worst case against O(len(argv))
  today … the scan is not on any hot path and no alternative was rejected for
  cost." — (§performance-expectations) — NEGATIVE REQ: no performance
  assertion or benchmark is owed; the quadratic scan is accepted.

## Rejected alternatives (§alternatives-considered, lines 206–258) — NEGATIVE REQs

- [REQ-67] ALT1 rejected: do not keep the argv0 heuristic and merely reword the
  claim — "it narrows the claim to a predicate that no longer serves the
  visibility purpose" (0027:ALT1, line 220).
- [REQ-68] ALT2 rejected: do not build "a table of wrappers, each with its
  option grammar" — "the position-free scan achieves the class closure with
  zero grammar and zero table" (0027:ALT2, line 236). Reinforces REQ-5.
- [REQ-69] ALT3 rejected: do not demote the refusal to an advisory — "only the
  `sh -s` axis belongs to the successor, and C1 routes exactly that"
  (0027:ALT3, line 251). The category stays a load-time REFUSAL.
- [REQ-70] Briefly rejected, each a NEGATIVE REQ (§briefly-rejected, lines
  255–258): no declared inline-shell opt-in (`shell` field); no exec-time
  identification of the resolved argv0 (`#!` / ELF read); no
  argv0-must-be-a-file-path rule; no spelling folding (`python3` → `python`)
  and no closing of the interpreter set.

---

## ASSUMPTIONS

Implicit choices made where wording was imprecise but a single reading is
defensible.

- ASSUMPTION: REQ-2 — "refuse iff there exist i < j" is evaluated over the
  argv AFTER 0025:C2 placeholder substitution is *not* applied — `carrierDefect`
  runs at LOAD time on the declared vector (`a.Command`), before any
  substitution. The record's examples (`["nice","sh","-c","cat {artifact}"]`)
  carry the raw placeholder, which only makes sense on the declared vector.
- ASSUMPTION: REQ-2 — `i` and `j` range over the whole argv INCLUDING index 0,
  so the pre-existing argv0 forms (`["sh","-c",…]`) still refuse under the
  widened predicate; "position-free" widens the set, it never excludes argv0.
  Confirmed by REQ-41's "the four REQ-74 probes … still green".
- ASSUMPTION: REQ-2/REQ-13 — "one of its listed flags" means the flag list
  keyed by the interpreter found at `i`, not the union over all listed
  interpreters. C1's "one of THAT interpreter's inline-code flags" (REQ-1) and
  A7's "flag drawn from THAT interpreter's own list" pin this; so
  `["ruby","-c"]` does NOT refuse (`ruby` carries only `-e`).
- ASSUMPTION: REQ-6/REQ-25 — the tie-break is "lowest `i` that is a listed
  interpreter AND has some qualifying `j > i`". A listed interpreter at a lower
  `i` with no following flag must not short-circuit the scan; A7's
  `["setup-tests","bash"]` row ("has no following flag so C1 never fires")
  shows a bare trailing interpreter is green, and the illustrative code's
  `continue` on a non-match confirms the outer loop keeps scanning.
- ASSUMPTION: REQ-6 — the reported form's separator is a single ASCII space
  (`argv[i] + " " + argv[j]`), matching today's `interpreterForm` return and
  the record's `sh -c` / `ruby -e` / `python -c` spellings throughout.
- ASSUMPTION: REQ-16 — "the detail additionally names the two matched words"
  is satisfied by the EXISTING detail construction, which already interpolates
  `interpEl` (today `argv[i] + " " + argv[j]`) into the message. The "addition"
  is that the widened predicate makes `interpEl` a mid-argv pair; no new
  message field is required, and S5's oracle is a substring assertion on the
  detail text. Reading taken because REQ-22 forbids a signature change and
  REQ-16 says the remediation text is unchanged.
- ASSUMPTION: REQ-7 — "never reads … PATH" is a statement about the LINT
  predicate only. `internal/cli/cmdbind`'s exec-time argv0 resolution through
  `PATH` (0025:C2/REQ-23) is untouched; the two are different layers and the
  record's `authority` mini-check names `resolveArgv0` as a non-sibling.
- ASSUMPTION: REQ-14 — clause 3's exemption keeps its existing shape
  (`!isInterp || no-whitespace ⇒ defect`); only the WRITER's predicate changes.
  REQ-29's "leave `filepathBase` unchanged" plus the `authority` mini-check's
  "Reuse the read unchanged; the writer's predicate change fixes it"
  (§existing-infrastructure-audit, line 188) pin this.
- ASSUMPTION: REQ-31 — "delete the `env` walk" means the loop skipping `env`
  and `NAME=VALUE` tokens is removed outright, not made conditional. The
  position-free scan subsumes it (an `env`-prefixed argv has the interpreter at
  some `i ≥ 1`), and S3 is the subsumption proof.
- ASSUMPTION: REQ-36/REQ-38 — the "exported accessor pairing the category with
  its text" is a new exported function in `internal/table` beside
  `Categories()` (the record cites `internal/table/category.go::Categories` as
  the precedent shape and `graphlint.BlockingCodes()` as the pattern). Its Go
  name and signature are unconstrained implementation choice: the record pins
  the SURFACE (`--help-all`, mirrored to `docs/cli-reference.md`) and the TEXT
  (C1's words), not the symbol. Consistent with 0025:REQ-15's precedent that
  Go names are not contracts.
- ASSUMPTION: REQ-33/REQ-36 — REQ-33's "edit confined to `load.go`" scopes the
  PREDICATE work (Phase 1) only. Phase 3 necessarily adds a description surface
  which the record itself says does not exist in `category.go` today, so
  touching `category.go` (or a new file in `internal/table`) for the accessor
  is in scope; what REQ-48 forbids is changing `Categories()`' membership,
  order, or the `command_shell_interpreter` wire string.
- ASSUMPTION: REQ-49 — S7 reads the description through the `--help-all`
  rendering path (or the `docs/cli-reference.md` mirror), not by calling the
  accessor directly: the scenario's whole point is that the text is reachable
  by a reviewer "with no model loaded and no defect present". A direct
  accessor-only assertion would pass even if the body never rendered it.
- ASSUMPTION: REQ-45 — "loads with its argv unchanged" is asserted on the
  loaded `table.Accessor`'s `Command` slice being element-wise equal to the
  declared vector, which is the only argv the loader exposes; no exec is
  required for S4 (the MVV's step-3 exec arm is separately scoped to the two
  forms it names, and only on a GNU-`env` host for the `-S` case).
- ASSUMPTION: REQ-42 — S1's mutants are asserted through `table.Load` on model
  fixtures (REQ-41), not by calling `interpreterForm` directly, so the oracle
  is the category a model author sees. Direct unit calls on the unexported
  predicate are permitted as an additional layer but do not satisfy S1.
- ASSUMPTION: REQ-40 — the prerequisite is satisfied by RECORDING q2q1's
  disposition in the deviations artifact; this run does not need q2q1 to have
  landed, because the clause says "either way the `env` walk is removed here".
- ASSUMPTION: REQ-16/REQ-48 — "position in `table.Categories()` … unchanged" is
  read as RELATIVE position (the `command_shell_interpreter` constant's place
  among its neighbours), per the pre-seeded deviation D1 and 0025:REQ-79's "the
  list's total size is not a contract at any point". A tail- or `len`-equality
  golden would contradict 0025:REQ-79.
- ASSUMPTION: REQ-53 — `["sh"]` (bare, single-element) is green under the
  predicate for the trivial reason that no `j > i` exists; no special case is
  written for it. The disposition table lists it as admitted, and the predicate
  already delivers that.
- ASSUMPTION: REQ-19/REQ-21 — "THESE words" and "verbatim" are satisfied by the
  description containing C1's out-of-scope wording substantively (both forms
  named, channel-scoped), not by a byte-for-byte copy of the whole normative
  fence, which contains implementation-facing text (`base(argv[i])`, `i < j`)
  unsuited to user-facing help. S7's oracle is "its text names both
  out-of-scope forms … in C1's channel-scoped words" — a naming assertion, not
  a byte-equality one.

---

## QUESTIONS

Genuinely ambiguous — two readings produce materially different behaviour and
no in-record or predecessor precedent settles them. Recorded, not blocking.

- Q1 — REQ-38 vs pre-seeded deviation D2: **is the Phase 3 description surface
  TOTAL over `Categories()`, or per-category opt-in?** REQ-38 names "an
  exported accessor pairing the category with its text" without saying whether
  every member of `Categories()` must have text. If total, cli/0028's five
  `edit_*` categories (0028:C1.4) fail the accessor's test or ship undescribed,
  and neither record names that obligation; if per-category opt-in, the
  `--help-all` body lists a partial taxonomy, which is weaker than the
  `graphlint.BlockingCodes()` precedent A6 cites (that taxonomy is closed and
  total). **Reading taken for the build**: per-category opt-in — the accessor
  returns text where declared and the extended-help body renders only described
  categories. Grounds: (a) REQ-48 forbids changing `Categories()`, and a
  totality assertion would couple this record to every future append,
  contradicting 0025:REQ-79's append-only, size-is-not-a-contract rule; (b) C1
  ships text for exactly ONE category and S7 asserts exactly that one; (c) the
  cluster-reconcile pre-seeded D2 already flags 0028 breakage under the total
  reading. This is D2's named disposition and must be recorded there.

- Q2 — REQ-19/REQ-21 vs REQ-49: **how literal is "THESE words … verbatim"?**
  A byte-equal copy of C1's `out of scope, BY NAME` line into user-facing help
  would ship contract syntax (`base(argv[i])`, `i < j`, peer-record ids) to a
  CLI reader; a paraphrase risks the premortem's second failure (a reviewer who
  stops at "closes the wrapper class"). No predecessor precedent exists —
  0025 shipped no description surface at all. **Reading taken**: substantive
  fidelity — the description states the predicate in one sentence and names
  both out-of-scope forms with C1's channel scoping, and S7 asserts on those
  named forms rather than on byte equality (see the ASSUMPTION above). Flagged
  because the opposite reading makes S7 a golden-file test whose maintenance
  burden and failure mode differ materially.

## Withdrawn REQs

**REQ-32 (withdrawn in Phase 3c) — not a testable clause.** "Illustrative —
  shape only." (§illustrative-code, line 136) is a CAVEAT ON A CODE BLOCK
  telling the implementer the snippet is non-binding; `rdr inspect` confirms
  line 136 falls inside NO labelled contract element, so it carries no
  normative force. Extracting it as a REQ was a Phase 0 classification error:
  it obliges nobody to do anything, and its "compliance evidence" would be
  the absence of a test, which no assertion can express. Recorded as
  deviation D8. The snippet's actual content is covered behaviourally by the
  C1 predicate REQs (REQ-41..REQ-47).
