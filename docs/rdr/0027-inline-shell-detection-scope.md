# Recommendation 0027: What the inline-shell check on command bindings actually promises (wrappers, env flags, stdin scripts)

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-31
- **Status**: Draft
- **Type**: Architecture
- **Profile**: large — C1, the successor predicate and claim wording for 0025:C5's `command_shell_interpreter` deny-list; user-facing yes; locks format
- **Priority**: Medium
- **Related Issues**: intrastate#zvtg (defect tracker; stays open until Stage 8); intrastate#v0hb (closed, RDR 0025's tracker); intrastate#q2q1 (open — the contained `env` option-flag sub-fix; subsumed by C1's position-free predicate, may land first without conflict); the charted `stdin = "none"|"envelope"` successor in RDR 0025's `evidence/critique/Charted.md`
- **Predecessors**: 0025-command-invoking-accessor-bindings
- **Overrides**: 0025:C5's `command_shell_interpreter` line and its `interpreter set` line — C1 here is the successor text: the predicate widens from argv0-after-an-`env`-walk to any argv position under any prefix, and the claim names what the check reads and the two forms it does not
- **Seam Lineage**: `internal/table/load.go::interpreterForm` (`area:internal-table`) — no prior accretion. Count provenance: no `kata-scope-review §seam-accretion` emission exists on intrastate#zvtg (routed by `rdr-seed-triage`); taken at seed from `kata list --status closed --label area:internal-table` on 2026-08-31 — the one closed C5 fix at this file, intrastate#b84g, is a different symbol (placeholder-category validation, C5 clause 3), and intrastate#q2q1 is open.

## Problem Statement

A reviewer reading a model under `--allow-commands` (RDR 0025) relies on `intrastate
lint`'s `command_shell_interpreter` refusal to make inline shell visible: a `command`
binding that hands code to `sh -c` must be spelled where the review can see it, not
smuggled. Today that reviewer is told `["env","sh","-c",…]` and `["env","A=1","sh",
"-c",…]` are refused, while `["env","-i","sh","-c",…]`, `["env","-u","FOO","sh",
"-c",…]`, `["nice","sh","-c",…]`, `["timeout","5","sh","-c",…]`, `["xargs","sh",
"-c",…]` and `["sh","-s"]` lint clean and execute a shell. They discover it only by
reading argv themselves — the thing the lint exists to spare them — or by an
`env -i sh -c` that reads as deliberate evasion of a reviewer-visible check.

System-internally, `interpreterForm` walks an `env` chain by skipping `env` and bare
`NAME=VALUE` tokens and stopping at the first non-assignment token, so option flags are
invisible to it. Part of that is a plain conformance gap: 0025:C5 names "env chains" as
in scope and 0025's `verification.md` records the `env`-chain form as a Phase 3a PASS,
yet only assignment-only chains are caught (the contained flag-walk fix is
intrastate#q2q1, not this RDR). The rest is not covered by C5's text at all: `nice`,
`timeout`, `xargs`, `nohup`, `setsid`, `stdbuf`, `chpst`, `doas` each re-open an
argv-position-aware wrapper class, and `sh -s` / `env -S "sh -c …"` are a different
axis — a script carried inside one word or read from stdin, which no argv-level check
can see without modelling the wrapper's own parser.

The decision this RDR owns, exactly once: what is the *authority surface* of C5's
deny-list, and therefore what does it promise? It is a policy/scope decision about
predicate shape plus the claim wording verification is held to, not an implementation
detail — the candidate answers are weighed in §Alternatives Considered.

## Critical Assumptions

- **A1 `interpreterForm` has exactly two consumers, both inside `carrierDefect`: clause 3's whitespace exemption (the `isInterp` read) and clause 4's refusal — so one predicate change reaches both clauses and nothing else.**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/table/load.go::carrierDefect` is the sole call site — `interpreterForm(argv)` is invoked once (`load.go:1095`), and its `isInterp` result is read exactly twice: clause 3's whitespace exemption (`load.go:1114`, `!isInterp || …`) and clause 4's refusal (`load.go:1124`). A repo-wide search for `interpreterForm` returns only that call, the declaration and its doc comment; `isInterp` appears at those three lines and nowhere else. The reuse audit confirms no second argv scanner or command-position classifier exists anywhere in the audited paths.
  - **If wrong**: a second consumer keeps the argv0-only shape and the two disagree on which defect a wrapper form reports.
- **A2 No binding that lints green today carries a listed interpreter basename as a non-argv0 word followed later by one of that interpreter's flags — widening the predicate refuses nothing currently accepted in the repo's fixtures or 0025's verification models.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a2-widened-predicate-no-new-refusals.md` — both predicates computed over 32 `command` argv cases drawn from `internal/table/testdata/neg/*.toml` (6, all already-negative; two of them the same `["rdr-gate", "{artifact}"]` vector) and `command_carrier_0025_test.go` (26 of that file's 27 distinct literals) — 11 green, 21 already-negative; no positive fixture declares `command`. Zero divergences: no vector moves `old=false → new=true`, so no green binding is newly refused. `["perl","-e","print 1"]` stays green under both — the widen is over argv POSITION, not over which basenames are listed, so an unlisted spelling is unaffected. `go test ./internal/table/` green before and after. Scope: this repo's fixtures and 0025's models, which is exactly the claim's scope; it says nothing about an out-of-tree user model, where the widened refusal is the intended behaviour.
  - **If wrong**: a green model turns red on upgrade and the refusal's remediation ("put it in a script") misdirects, because there is no script to extract.
- **A3 Every wrapper form the Problem Statement names (`env -i`, `env -u FOO`, `nice`, `timeout 5`, `xargs`, `nohup`, `setsid`, `stdbuf`, `chpst`, `doas`) places the interpreter name and its inline-code flag as separate argv words in that order, so the position-free predicate refuses each — and a brace-bearing string under a wrapper reports `command_shell_interpreter`, not `command_unknown_placeholder`.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a3-wrapper-argv-transparency.md` — two parts. (1) Each named wrapper was run through a real no-shell `exec.Command` harness; every one presents the interpreter and its inline-code flag as separate, ordered argv words (`["nice","sh","-c",…]`, `["env","-u","FOO","sh","-c",…]`, `["stdbuf","-o0","sh","-c",…]`, …). Four wrappers (`setsid`, `timeout`, `chpst`, `doas`) are absent from stock darwin and report `executable file not found`; their argv layout is still correct, and the layout is all the predicate reads, so host availability does not bear on the claim. (2) Both predicates computed over those exact vectors: all ten flip `false → true` under C1, each reporting `sh -c`, while `["sh","./gate.sh"]`, `["env","-S","sh -c echo"]` and `["sh","-s"]` stay `false`. The second half (category under a wrapper) is structural: `isInterp` gates clause 3's exemption at `load.go:1114` and triggers clause 4 at `:1124` from the single call at `:1095`, so `["nice","sh","-c","cat {artifact}"]` — today `command_unknown_placeholder`, per 0025 deviation D18 — reports the interpreter form once the predicate matches. End-to-end confirmation is the MVV's step-2 `{artifact}` mutant at implementation.
  - **If wrong**: a named wrapper form still lints green and C1's promise is false at lock; or the wrapper form masks the interpreter defect (0025's D6 failure re-opened).
- **A4 From a fixed argv executed without a shell, the only forms that spell or source shell code outside separate argv words are a one-word shell string (`env -S "sh -c …"`, or a single `"sh -c …"` element handed to a tool that re-splits it) and a stdin-fed shell (`sh -s`, bare `sh`); `sh <script` is unreachable (no redirection without a shell) and `sh script.sh` is the sanctioned wrapper-file form.**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a4-reachable-shell-forms-from-fixed-argv.md` — seven probes plus five adversarial candidates, each run through a real no-shell `exec.Command` harness. `sh <script` is empirically unreachable: `<` arrives as a literal filename operand and `sh` exits 127 without touching the script. `sh ./gate.sh` runs the declared FILE, not inline code. The adversarial round found no reachable third form: bundled `sh -cecho …` is rejected by `sh` itself AND fails the predicate's exact flag match; `/bin/sh -c` is still caught via `base()`; a renamed interpreter is the by-design OPEN deny-list case. `sh -es`, `python -` and `node -` do execute caller code and are admitted — but they deliver it over STDIN, the channel C1 already declares out of scope; the spike is what widened that line from "a stdin-fed shell" to the channel, since both the old and widened predicates admit all of them. So the admitted-and-executes set equals C1's named list. Platform note (per 0025:F9, recorded here and not in the contract, which stays platform-neutral): `env -S` requires GNU env — verified against GNU env 9.11, which splits the word and runs the shell — while darwin's stock BSD env rejects `-S` outright; the FORM (a shell string inside one argv word) is reachable on any host, and C1's second spelling, a single `"sh -c …"` element re-split by the receiving tool, carries no such dependency.
  - **If wrong**: a reachable form outside C1's named out-of-scope list executes a shell and the "named by name" claim is incomplete.
- **A5 The stdin-fed-interpreter hazard (`sh -s`, bare `sh`, and the non-shell spellings `python -` / `node -`) is owned by the charted `stdin = "none" | "envelope"` successor: withholding the envelope from a command that never declared it leaves a stdin-reading interpreter nothing to execute, so routing that axis there leaves no reviewer-visible gap unowned.**
  - **Status**: Verified
  - **Method**: Design Decision
  - **Evidence**: this RDR's scoping choice, resting on 0025's `evidence/critique/Charted.md`, whose entry is scoped to "Declared stdin appetite (`stdin = "none" | "envelope"`) **on a command entry**" and whose motivating hazard is `tee {artifact}` — not a shell. The successor's remit is therefore what a command entry CONSUMES, which covers the non-shell spellings on the same mechanism, so C1's channel-scoped wording describes that remit rather than extending it. Corroborating prior art: consul's `command/exec/exec.go` gates stdin-as-script (`cmd == "-"`) on its declared shell flag, treating stdin delivery as one concern with shell mode. **Durability caveat**: no kata tracks the successor — it exists only as that `Charted.md` line (`kata list --status open`, 2026-09-03), so the If-wrong below is nearer than a charted item normally implies; until it ships, C1 names the forms admitted and the reviewer reads argv.
  - **If wrong**: the successor never ships or withholds nothing, and the stdin forms stay unowned admitted interpreters.

## Proposed Solution

### Approach

Decide once what 0025:C5's deny-list is the authority over: **argv words, at any position** — never the contents of a word, stdin, files, or the resolved binary. The predicate is freed from argv0: a binding is refused when a listed interpreter word is followed, at any later argv position, by one of that interpreter's inline-code flags. That closes the whole wrapper class (`env -i`, `env -u FOO`, `nice`, `timeout 5`, `xargs`, `nohup`, `setsid`, `stdbuf`, `chpst`, `doas`, and every wrapper not yet thought of) with **no wrapper table and no option grammar**, because nothing before the interpreter word is read at all. The `env`-chain walk in `internal/table/load.go::interpreterForm` is deleted: this record DECIDES that C1 subsumes it, which the tracker deliberately left open (intrastate#q2q1 "ships independently, and does not foreclose any of the three answers"). Subsumption is verified, not assumed — the position-free predicate refuses all nine forms q2q1 enumerates (`-i`, `-u FOO`, `--unset=FOO`, `-0`, `-C DIR`, `--chdir=`, `--`, `A=1`, and mixed chains), since nothing before the interpreter word is read (`evidence/research/stage4-scope-wording.md` §Q3). q2q1 may still land first; the walk goes either way.

The claim wording is narrowed to exactly that predicate, and the two forms an argv-level check cannot see are declared out of scope **by name**: a shell string carried inside one word (`env -S "sh -c …"`), and an interpreter that reads its script from stdin (`sh -s`, bare `sh`, and the non-shell spellings `python -` / `node -`), the latter routed to the charted `stdin = "none" | "envelope"` successor, which withholds the only script source a fixed argv can reach. That second form is named by its CHANNEL rather than by a list of shell spellings: the deny-list already carries `python`, `ruby`, `node` and `php`, so a promise that said "a stdin-fed shell" would leave a reviewer to discover `["python","-"]` themselves — the one thing this record exists to stop. `sh script.sh` stays the sanctioned wrapper-file form.

What the lint then promises a reviewer is one sentence they can hold: *if an interpreter and its inline-code flag appear as two separate words in that order, anywhere in argv, the binding is refused; a one-word shell string or a stdin-fed interpreter you read yourself.* This keeps C5's own frame — "the interpreter deny-list is defense in depth over that, not the barrier itself" (0025:§normative-contracts, beside C5) ⇒ the deny-list is never re-promoted to a barrier, so no wrapper table is justified; but the cost it raises is no longer defeated by one leading word, which is what a check "whose purpose is visibility" owes the reader.

The prior-art frame supports the shape: no peer CLI detects shell-ness from argv — gh declares it (`pkg/cmd/alias/set/set.go::NewCmdSet`, `--shell` "Declare an alias to be passed through a shell interpreter") ⇒ detection is intrastate's own construction and its promise must be bounded by its own predicate, not borrowed; and Go's `os/exec` package doc ("intentionally does not invoke the system shell … pipelines, or redirections") ⇒ from a fixed argv a shell runs only through an argv word, so words are the complete surface an argv check can own and `sh <script` is not a reachable form.

### Technical Design

`carrierDefect` (`internal/table/load.go::carrierDefect`) computes `interpEl, isInterp := interpreterForm(argv)` once and reads `isInterp` twice: clause 3 exempts a whitespace-bearing brace element from `command_unknown_placeholder` only when the form is an interpreter form, and clause 4 refuses on it. The writer of that state is `interpreterForm` at load time, and it is the only writer (A1) ⇒ changing the predicate in one place moves both clauses together, and the wrong-defect masking a wrapper causes today (`["nice","sh","-c","cat {artifact}"]` reports the placeholder defect) is fixed by the same edit.

`interpreterForm` becomes a two-index scan: for the lowest `i` whose `filepathBase(argv[i])` is a key of `shellInterpreters`, and the lowest `j > i` with `argv[j]` in that key's flags, return `argv[i] + " " + argv[j]`. Everything before `i` is unread. `shellInterpreters` and `filepathBase` are unchanged. The refusal keeps its category, remediation and `table.Categories()` position, and additionally names the two matched words so a false refusal (a script whose own argument happens to be `-c`) is diagnosable from the message.

Not modelled, on purpose: wrapper option grammars, `--` conventions, `env -S` splitting, PATH or binary identity. Each would be a parser the lint has to keep correct forever and a mis-parse is a silent admission; the position-free scan has no such surface.

#### Normative Contracts

The successor text for 0025:C5's `command_shell_interpreter` line and its `interpreter set` line (0025 is not edited; `load.go`'s C5 comments re-cite this clause). The other five C5 categories, C5's registration and precedence lines are unchanged and not restated.

**C1**

```normative
command_shell_interpreter     # a listed interpreter word followed, at ANY later argv position, by one of that interpreter's inline-code flags — under any prefix (env and its options, nice, timeout, xargs, doas, …); no opt-in in v1
predicate:  refuse iff there exist i < j with base(argv[i]) a listed interpreter and argv[j] one of its listed flags; base = the text after the last `/`, matched exactly (no suffix or alias folding: `python3`, `nodejs`, `busybox` stay unlisted spellings under the OPEN rule). Nothing before argv[i] is read; no wrapper table exists and none may be added. The reported form is argv[i] + " " + argv[j] for the lowest i, then the lowest j
reads:      argv WORDS only. The check never splits a word on whitespace, never reads stdin, files, PATH, or the resolved binary
out of scope, BY NAME (admitted by lint; an interpreter may run): a shell string carried in ONE word (`env -S "sh -c …"`, or a single `"sh -c …"` element handed to a tool that re-splits it); an interpreter that reads its script from STDIN (`sh -s`, bare `sh`, `sh -es`, `python -`, `node -`) — owned by the charted `stdin = "none" | "envelope"` successor, whose remit is what a command entry CONSUMES and so covers the non-shell spellings too. The channel is the scope: any listed interpreter taking its code on stdin rather than as a later argv word is admitted, however spelled. `sh script.sh` is the sanctioned wrapper-file form and never a defect
interpreter set: OPEN (deny-listed, not closed), unchanged from 0025:C5 — an unlisted spelling is admitted; the list grows by amendment of THIS clause
clause 3 coupling: 0025:C5 clause 3's whitespace exemption keys on this predicate, so a brace-bearing string under a wrapper reports the interpreter form, never command_unknown_placeholder; within-entry precedence is 0025:C5's, unchanged
report:     category string, remediation ("inline shell is not a declared command; put it in a script and declare the script as argv0") and position in `table.Categories()` are unchanged; the detail additionally names the two matched words
promise:    what a reviewer may rely on is the predicate line and nothing more; the lint's user-facing description states the out-of-scope forms in the words above. No such description exists in the shipped code today — the only user-facing text is the refusal detail, which fires on refusal and never on admission — so Phase 3 creates it and THESE words are the text it ships
```

#### Load-Bearing Decisions

- **Naming** — the category stays `command_shell_interpreter`; rejected: `command_inline_shell` (a rename re-orders nothing a reviewer sees and breaks 0025's tests and `Categories()` wire order for no gain).
- **Selection / predicate** — when several pairs qualify, the lowest interpreter index and then the lowest flag index are reported; order is part of the predicate (a flag word before the interpreter word never matches); the flag is matched at any later position, as today, so the already-accepted false-refusal class (a script argument that is itself `-c` / `-e`) is unchanged and is diagnosable from the named words.

#### Illustrative Code

```go
// Illustrative — shape only.
for i, w := range argv {
    flags, ok := shellInterpreters[filepathBase(w)]
    if !ok { continue }
    for _, later := range argv[i+1:] {
        if slices.Contains(flags, later) { return w + " " + later, true }
    }
}
return "", false
```

Sample refusals (Illustrative): `["timeout","5","sh","-c","…"]` → `sh -c`; `["env","-i","bash","-c","…"]` → `bash -c`; `["xargs","sh","-c","…"]` → `sh -c`. Admitted: `["sh","./gate.sh"]`, `["env","-S","sh -c echo"]`, `["sh","-s"]`.

#### Mini-check: `authority`

Fired by the two-reader `isInterp` decision and the sibling-arm question in the infrastructure audit.

| Input / decision | Writer | Readers | Call sites | Sibling arms | Canonical |
| --- | --- | --- | --- | --- | --- |
| "is this argv an interpreter form?" (`isInterp`, `form`) | `internal/table/load.go::interpreterForm` — the only producer; C1 replaces its body, same name and signature | clause 3's whitespace exemption (`load.go:1114`, `!isInterp \|\| …`); clause 4's refusal (`load.go:1124`) | one — `load.go:1095`, inside `carrierDefect` | none: `cmdbind.go::resolveArgv0` resolves argv0 *location*, never identity, and `shellInterpreters` has one use site (`load.go:1174`). Searched repo-wide; no wrapper or command-position classifier exists | `interpreterForm` is sole source of truth; both readers consume the one call's result, so the predicate change reaches clause 3 and clause 4 together and cannot skew |
| listed interpreter basenames | `shellInterpreters` map (`load.go:1035-1047`) | `interpreterForm` only | one (`load.go:1174`) | none | the map; OPEN deny-list, unchanged by this record |
| category string + position | `internal/table/category.go` (`:60` declaration, `:103` in `Categories()`) | lint output / `table.Categories()` | n/a — not edited | none | `category.go`; C1's edit is confined to `load.go` |

#### Mini-check: `disposition`

Fired by C1's `out of scope, BY NAME` line, which admits input classes while siblings refuse.

| Input class | Exit / lint outcome | Event / error | Artifact / op minted | Silent vs loud |
| --- | --- | --- | --- | --- |
| interpreter word + its inline-code flag at any later argv position (`["nice","sh","-c","…"]`) | refuse — load-time defect | `command_shell_interpreter`, detail names the two matched words + the script remediation | none (load fails) | **loud** |
| the same under a brace-bearing string (`["nice","sh","-c","cat {artifact}"]`) | refuse | `command_shell_interpreter`, **not** `command_unknown_placeholder` — clause 3's exemption keys on the same `isInterp` | none | **loud**, and this is the wrong-defect masking 0025 D18 left open |
| script argument that is itself a listed flag (`["ruby","tool.rb","-e","prod"]`) | refuse — accepted false-refusal class | `command_shell_interpreter`, detail names `ruby -e` so the collision is diagnosable | none | **loud**; already-accepted behaviour, the *message* is what this record adds |
| shell string inside ONE word (`["env","-S","sh -c echo"]`) | green | none | binding loads; a shell may run | **silent by contract** — named in C1's out-of-scope line; no lint event is minted |
| interpreter reading its script from stdin (`["sh","-s"]`, `["sh"]`, `["sh","-es"]`, `["python","-"]`, `["node","-"]`) | green | none | binding loads; a shell may run | **silent by contract** — named by CHANNEL, not by shell spelling; routed to the charted `stdin` successor |
| sanctioned wrapper file (`["sh","./gate.sh"]`) | green | none | binding loads and runs the declared file | **silent, and correct** — never a defect |
| unlisted spelling (`["python3","-c","…"]`, `["perl","-e","…"]`) | green | none | binding loads | **silent by contract** — OPEN deny-list; the list grows only by amending C1 |


### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Withholding stdin from a command that never declared it (so a stdin-fed interpreter has no script) | Future — the charted `stdin = "none" \| "envelope"` successor; no kata tracks it yet (A5) | Deferred | C1 names the stdin channel out of scope and routes it there; until it ships the forms are admitted and named as such |
| Refusal category, remediation text, `Categories()` registration | Predecessor — 0025:C5 | Available | none; C1 changes only the predicate and the claim |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Interpreter-form predicate | `internal/table/load.go::interpreterForm` | interpreter must sit at argv0 after an `env` walk that stops on any option flag | Replace the body (same name, same signature): position-free scan; the `env` walk is deleted as subsumed | C1 predicate line |
| argv word normalisation | `internal/table/load.go::filepathBase` | last-`/` split only, no alias folding | Reuse | C1 `base` definition |
| Clause 3 / clause 4 coupling | `internal/table/load.go::carrierDefect` (`isInterp`) | keyed on the argv0 form, so a wrapper masks the interpreter defect | Reuse the read unchanged; the writer's predicate change fixes it | C1 clause-3 coupling line |
| argv0 reasoning at exec time | `internal/cli/cmdbind/cmdbind.go::resolveArgv0` | resolves location (abs / bare / relative), not identity | Reuse unchanged — the sibling checked at Propose; no wrapper or command-position classifier exists anywhere in the repo | none |
| C5 probes | `internal/table/command_carrier_0025_test.go::TestReq74_KnownInterpreterWithInlineCodeIsALoadTimeDefect` and the `neg/neg-command-shell-interpreter.toml` fixture | four argv0-position probes | Extend with the wrapper mutants | Testing Strategy (Resolve) |

### Decision Rationale

Three factors decide. (1) *What can an argv check own?* Only words — Go's `os/exec` doctrine makes a word the only path to a shell from a fixed argv — so a promise stated over words is complete on its own surface, and anything inside a word or on stdin is out of scope by construction, not by omission. (2) *Cost of the class fix.* Freeing the interpreter's position closes every wrapper with no wrapper table, because nothing before the interpreter word is read; the alternative (a per-wrapper walk) is a parser per wrapper and a silent admission per mis-parse — the same open-set problem C5 already refused to pretend to close, one level up. (3) *What the prose already says.* 0025 scoped the deny-list as "defense in depth … not the barrier", so re-promoting it (wrapper table) is unjustified and demoting it to a hint (Alt 3) throws away the visibility it does give; the honest middle is a predicate that costs a reviewer nothing to remember and that no single leading word defeats.

Rejected: Alt 1 (narrow the claim only) leaves a visibility check that `nice` defeats and asks the reviewer to remember a list of wrappers instead. Alt 2 (per-wrapper walk, bounded set) enumerates the unenumerable and adds mis-parse as a new silent-admission channel. Alt 3 (retire argv as authority; hint only) answers a different hazard — the stdin successor is about what a command *consumes*, not what it *spells* — and a refusal that still fires is load-bearing by behaviour whatever the prose calls it. The briefly-rejected options each contradict a locked decision (opt-in) or move the promise out of lint (exec-time identification).

Premortem: hardened (paragraph). Shipped and failed: the widened scan refuses a team's `["ruby","tool.rb","-e","prod"]`-shaped binding — an interpreter word whose script takes a `-e` of its own — and the remediation ("put it in a script") sends them in circles because there is no inline code; separately, a reviewer who read "closes the wrapper class" stops reading argv and an `env -S "sh -c …"` string ships. The first is the already-accepted flag-anywhere class, not new, but the message was the failure: the refusal now names the two matched words (C1 `report:`), so the fix (rename the flag, or wrap) is obvious from the message. The second is a promise-wording failure: C1's `out of scope, BY NAME` line and the lint's user-facing description carry the two admitted forms in those words, so "closes the wrapper class" is never the whole sentence. Neither forces a switch; both hardened the claim, not the predicate.

Ground-sweep: clean (17 anchors) — every `path::Symbol`, peer element (0025:C5, 0004:ALT2, 0004:ALT3), 0025 artifact/evidence quote and peer-system citation confirmed verbatim by a fresh-context read; `interpreterForm` has no caller outside `internal/table/load.go`.

Joint-check: fired → 0028 (home: cli/0025:C5) — symmetric to 0028's arm-1 fire on `internal/table/load.go::carrierDefect` / `table.Categories()`, disposed cite-don't-restate: C1 here owns only clause 4's predicate, cli/0028:C1 and cli/0028:C4 own the `edit` arm and the tail categories, and both fence deltas against 0025:C5, which owns the clause map and within-entry precedence. Recorded at the batch propose of 2026-08-31 after this record's own check ran clear: all three arms run. Arm 1 (modify-anchors) and arm 2 (contract literals), open-only and repo-resolved: no overlaps. Arm 3 (absence, manual): the 27 Draft/Final/Implemented peers grepped for the admitted-form and wrong-defect tokens (`interpreterForm`, `env chain`, `env -i`, `xargs`, `sh -s`, `env -S`, `command_unknown_placeholder`, `command_shell_interpreter`); hits only in 0025:C5 — the line C1 succeeds, cited not relied on — and 0004's "without raw shell strings", which C1 keeps true. 0026 compared on its written contracts (0026:C1 amends 0025:C4's drain precedence; disjoint from C5); 0028 compared on its problem statement (an in-process `edit` carrier that runs no shell, no subprocess); no shared decision, no bridge surface.

## Alternatives Considered

### Alternative 1: Narrow the claim to the argv0-name heuristic

**Description**: Keep `interpreterForm` as it is (plus the intrastate#q2q1 `env` option walk), reword 0025:C5 to say "argv0 after an `env` chain", and list wrappers (`nice`, `timeout`, `xargs`, …) and shell-string forms as out of scope by name.

**Pros**:

- Smallest change; no predicate risk; fully honest about what runs.
- Matches C5's own "raise the cost, not a barrier" frame.

**Cons**:

- The cost raised is one leading word; the reviewer has to carry the wrapper list in their head, which is the work the lint exists to spare them.
- The wrong-defect masking (`["nice","sh","-c","cat {artifact}"]` → placeholder defect) stays.

**Reason for rejection**: it narrows the claim to a predicate that no longer serves the visibility purpose; the position-free scan costs the same to state and closes the class.

### Alternative 2: Per-wrapper walk with a bounded wrapper table

**Description**: Extend the `env` walk into a table of wrappers, each with its option grammar (which flags take an argument), walk to the command position, then apply the interpreter check; state the wrapper set as the bound.

**Pros**:

- Precise command-position semantics; could later host per-wrapper rules (`env -S` splitting).

**Cons**:

- Every wrapper is a mini-parser (`timeout -k 5 10 sh -c`, `xargs -I{} -n1 sh -c`, `doas -u x`), and a mis-walk is a silent admission — a new failure channel the current code lacks.
- The wrapper set is as open as the interpreter set; the "bound" is a promise to keep enumerating.
- Re-promotes a defense-in-depth check to a barrier, against 0025's own frame and 0004:ALT3's rejection of name lists as authority.

**Reason for rejection**: the position-free scan achieves the class closure with zero grammar and zero table.

### Alternative 3: Retire argv as the authority surface; deny-list becomes a non-load-bearing hint

**Description**: Fold the question into the charted `stdin = "none" | "envelope"` successor, drop the conformance claim from C5, and keep the refusal as an advisory whose text promises nothing.

**Pros**:

- No predicate to maintain; pairs with a successor that is coming anyway.

**Cons**:

- The stdin successor governs what a command *consumes*; inline shell is about what argv *spells* — `stdin = "none"` plus `["sh","-c","…"]` is still inline shell.
- A refusal that still fires is load-bearing by behaviour; calling it a hint only makes the promise vaguer, which is the reviewer-trust failure restated.

**Reason for rejection**: it answers a different hazard and removes the visibility C5 does provide; only the `sh -s` axis belongs to the successor, and C1 routes exactly that.

### Briefly Rejected

- **Declared inline-shell opt-in (a `shell` field, gh's `--shell` / `!` model)**: the peer pattern, but 0025:C5 fixes "no opt-in in v1" and 0004:ALT2 rejected the inline script as the contract; the wrapper-file remediation serves visibility equally.
- **Exec-time identification of the resolved argv0 (`#!` / ELF read)**: a wrapper execs whatever follows it, so the resolved binary is the wrapper; and it moves a lint promise to run time, where the reviewer is absent.
- **Require argv0 to be a file path (no bare names)**: refuses `git`, `rdr`, `kata` and every PATH tool 0025's MVV drives; the interpreter set would still need listing for absolute paths.
- **Fold spellings (`python3` → `python`) or close the interpreter set**: a different amendment to C5's OPEN rule, out of this RDR's scope; noted for the successor that widens the list.

## Context

### Background

Found by the roborev triage of RDR 0025's implementation (`src:roborev`, batch
`rdr-0025`) with real lint runs. Denied, correctly: `env sh -c`, `env A=1 sh -c`,
`bash --login -c`. Admitted: `env -i sh -c`, `env -u FOO sh -c`, `nice sh -c`,
`timeout 5 sh -c`, `xargs sh -c`, `sh -s`; `env -i sh -c 'echo …'` was verified to
actually execute the shell.

Why this is not simply C5 being open: C5's "interpreter set: OPEN (deny-listed, not
closed) — an unlisted interpreter is admitted" covers Phase 3a's `python3` / `nodejs`
non-finding (unlisted argv0 *spellings*). `env -i sh -c` is different — `sh` *is*
listed and C5's own text names `env` chains as in scope — so the option-flag blindness
is a conformance gap against a claim C5 makes, while the `nice`/`timeout`/`xargs`
family is genuinely outside C5's text.

Threat model and impact: low on its own — C5 explicitly declines to make inline shell
impossible and states the deny-list's job as raising the cost, so a bypass voids no
promised guarantee; reachable only under the `--allow-commands` opt-in. The cost is to
reviewer trust in a check whose purpose is visibility.

### Technical Environment

Go. `internal/table/load.go`: `interpreterForm` (the `env`-chain walk) and the
`command_shell_interpreter` refusal class it feeds, evaluated by `intrastate lint`
over `command` bindings admitted under `--allow-commands` (RDR 0025). Governing
records: 0025:C5 (the deny-list contract and its "raise the cost" purpose), 0025's
`verification.md` (Phase 3a's "env chains are caught" PASS text), and the charted
`stdin = "none"|"envelope"` successor in 0025's `evidence/critique/Charted.md`, which
already concedes argv is not a sufficient authority surface.

## Research Findings

### Investigation

Read, not spiked: `internal/table/load.go::interpreterForm`, `::carrierDefect` (clauses 3–4 and the `isInterp` coupling), `::shellInterpreters`, `::filepathBase`; the REQ-74 probes in `internal/table/command_carrier_0025_test.go` and the `neg/neg-command-shell-interpreter.toml` fixture; 0025:C5 with its surrounding prose, 0025's `artifacts/verification.md` Phase 3a C5 row and `evidence/critique/Charted.md`; 0004:ALT2 and 0004:ALT3. Peer systems for the instance question — what does each do for a user-configured command: Go `os/exec` (never a shell), gh (`pkg/cmd/alias/set/set.go::NewCmdSet`, shell mode declared by `!` / `--shell`), roborev (`internal/daemon/hooks.go::(*HookRunner).runHook`, always `sh -c`), beads (`internal/creds/command.go::(CommandSource).Resolve`, always `sh -c`). No peer detects shell-ness from argv. ⚠ no prior-art coverage in the project's reference corpora for the wrapper-re-exec bypass class; the frame rests on the peer instance reads and 0004/0025. Sibling-path check: `internal/cli/cmdbind/cmdbind.go::resolveArgv0` reasons about argv0 location, `internal/cli/flowbind/registry.go` selects by carrier — searched, no wrapper or command-position classifier exists. Accepted citations and the query log are in this record's `evidence/research/prior-art.md`.

### Key Discoveries

- **Documented** — `interpreterForm` already matches the flag at *any* later position; only the interpreter's position is pinned (argv0 after an `env`/assignment walk). Freeing the position adds no new false-refusal class beyond the one already accepted.
- **Documented** — clause 3's placeholder exemption reads `isInterp`, so today a wrapper form carrying `{artifact}` reports `command_unknown_placeholder` — the wrong defect (0025 deviation D6's masking, re-opened under a wrapper). One predicate change fixes both clauses.
- **Documented** — Go `os/exec` never invokes a shell and performs no redirection ⇒ from a fixed argv a shell runs only via an argv word; `sh <script` is unreachable; one-word strings and stdin are the residual axes.
- **Documented** — every peer CLI read either declares shell mode (gh) or is always-shell (roborev hooks, beads credential helpers); none detects. The deny-list's promise cannot be borrowed and must be bounded by its own predicate.
- **Documented** — 0025's prose scopes the deny-list as "defense in depth … not the barrier itself"; `verification.md` asserts only the assignment-only `env` chain; `Charted.md` concedes argv is not the authority for what a command consumes (the `sh -s` axis).
- **Assumed** — the wrapper forms named in the Problem Statement all present the interpreter and its flag as separate words (A3); no green fixture is refused by the widened predicate (A2).

## Trade-offs

### Consequences

- Positive: every wrapper form in the defect tracker, and every wrapper not yet thought of, is refused by one predicate with no table to maintain; intrastate#q2q1's walk is deleted rather than extended.
- Positive: the reviewer's promise is one sentence, with the two admitted forms named in it.
- Negative: the flag-anywhere laxity stays, so a script argument that is itself `-c`/`-e` after an interpreter word is refused (already true today); the message now names the words so the fix is visible.
- Negative: `env -S "sh -c …"` and the stdin-fed forms (`sh -s`, bare `sh`, `sh -es`, `python -`, `node -`) are admitted and documented as such; the second class waits on the stdin successor, which no kata yet tracks (A5).

### Risks and Mitigations

- **Risk**: a green binding somewhere carries an interpreter word mid-argv followed by its flag and turns red on upgrade.
  **Mitigation**: A2's spike over all fixtures and 0025's models before lock; the named-words message makes any survivor self-explanatory.
- **Risk**: "closes the wrapper class" is read as "closes inline shell" and reviewers stop reading argv.
  **Mitigation**: C1's out-of-scope line is normative and the lint's description carries it verbatim; the docs never say "closes inline shell".

### Failure Modes

- Visible: a refused wrapper form names `sh -c` (the matched words) and the wrapper-file remediation; a false refusal names the two words that matched, so the author sees the script argument that collided.
- Silent: a one-word shell string or a stdin-fed interpreter (`sh -s`, bare `sh`, `sh -es`, `python -`, `node -`) lints green — by contract, named in C1; diagnosis is reading argv, which C1's promise tells the reviewer to do for exactly those shapes.
- Recovery: none needed for the predicate (no state); a wrong refusal is fixed by renaming the script's flag or wrapping, as the message says.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] intrastate#q2q1's disposition recorded (landed first, or closed as subsumed) — either way the `env` walk is removed here

### Minimum Viable Validation

1. Take the 0025 command-carrier model with a green `command` read binding; `intrastate lint` passes.
2. Swap in one mutant per named wrapper — `["env","-i","sh","-c","…"]`, `["env","-u","FOO","sh","-c","…"]`, `["nice","sh","-c","…"]`, `["timeout","5","sh","-c","…"]`, `["xargs","sh","-c","…"]`, `["doas","sh","-c","…"]` — and one carrying `{artifact}` inside the string under `nice`; each refuses with `command_shell_interpreter`, the detail naming `sh -c` and the script remediation; the `{artifact}` mutant does NOT report `command_unknown_placeholder`.
3. Swap in the admitted forms `["sh","./gate.sh"]`, `["env","-S","sh -c echo"]`, `["sh","-s"]`, and the non-shell stdin spelling `["python","-"]`; lint passes for each. The first two are shown, by running the model under `--allow-commands`, to behave as C1 states (the wrapper file runs; the `-S` string runs a shell — on a host whose `env` supports `-S`, GNU env, since darwin's stock BSD env rejects it; A4). `["python","-"]` is the probe that C1's out-of-scope line is stated over the stdin CHANNEL and not over shell spellings: it must lint green and be documented as admitted.
4. The original 0025 REQ-74 probes and the negative fixture still refuse with the same category, and `table.Categories()` order is unchanged.
End-state: every wrapper mutant red with the right defect, every admitted form green and documented, no existing probe changed.

#### Mini-check: `oracle`

Fired by MVV step 3 and Testing Strategy scenario 4, whose oracle is "lint passes" — an absence-of-error assertion that a no-op implementation also satisfies.

| MVV / scenario row | Fails if X is wrong, because Y | Negative / failing control |
| --- | --- | --- |
| MVV 2 · S1 — one mutant per wrapper refuses | fails if the predicate is still argv0-anchored, because every mutant puts the interpreter at argv ≥ 1 and would lint green | today's code: all ten mutants green (A3 computed all ten `false → true`) |
| MVV 2 · S2 — `["nice","sh","-c","cat {artifact}"]` reports the interpreter defect | fails if clause 3's exemption reads a different signal than the widened predicate, because the category would come back `command_unknown_placeholder` | assert the category *string*, not merely that load failed; today's code reports the placeholder defect here |
| MVV 3 · S3 — the nine `env`-option forms refuse | fails if the deleted `env` walk was not in fact subsumed, because a form q2q1 enumerates would lint green | today's code: `-i`, `-u FOO`, `--unset=FOO`, `-0`, `-C DIR`, `--chdir=DIR`, `--` all green (the walk stops at the first non-`NAME=VALUE` token) |
| MVV 3 · S4 — admitted forms lint GREEN | **absence-of-error oracle — the weak row.** A no-op predicate passes it. It is discriminating only as a *pair* with S1/S3: S1 forces refusal on the wrapper class, S4 pins the boundary that must not move with it | the discriminating members are `["python","-"]` and `["sh","-es"]`: any implementation reading the out-of-scope line as *shell spellings* refuses `python -` and fails S4. Control = assert green **and** that the binding loads with the argv unchanged, not merely that no error was returned |
| S5 — `["ruby","tool.rb","-e","prod"]` refuses naming `ruby -e` | fails if the detail does not carry the matched words, because the author cannot see which pair collided | assert the detail *text* contains both words; a bare category assertion passes even with today's message |
| S6 — `table.Categories()` membership and order unchanged | fails if the edit strays outside `load.go`, because the wire order a reviewer sees would shift | golden assertion over the full ordered slice, not a `Contains` check |

#### Mini-check: `trace`

Fired by C1's four normative lines (predicate, clause-3 coupling, `report:`, `promise:`) all bearing on one output surface — the category a model author sees. Walked over the MVV.

| Step | Assertions in force | Witness |
| --- | --- | --- |
| 1 — green `command` read binding lints clean | predicate (no `i < j` pair exists); `reads: argv WORDS only` | A2 spike: 11 green vectors, zero move `old=false → new=true`. Consistent |
| 2a — `["nice","sh","-c","…"]` refuses | predicate (`i=1` `sh`, `j=2` `-c`); `report:` names the two matched words | A3 spike: all ten wrappers compute `false → true`, form `sh -c`. Consistent — nothing before argv[i] is read, so `nice` needs no table |
| 2b — `["nice","sh","-c","cat {artifact}"]` reports `command_shell_interpreter`, not the placeholder defect | predicate; clause-3 coupling ("clause 3's whitespace exemption keys on this predicate"); 0025:C5 within-entry precedence, unchanged | A1: clause 3 reads the same `isInterp` the predicate sets (`load.go:1114`), from the one call at `:1095`. Consistent — the two clauses cannot disagree because they read one value |
| 3a — `["sh","./gate.sh"]` green | predicate (`./gate.sh` is not a listed flag); `out of scope` line's sanctioned wrapper-file form | A4: runs the declared FILE, not inline code. Consistent |
| 3b — `["env","-S","sh -c echo"]` green | `reads:` line ("never splits a word on whitespace"); `out of scope` one-word form | A4: the `-S` string runs a shell on GNU env. Consistent — admitted knowingly, and `promise:` requires the shipped description to say so |
| 3c — `["python","-"]` green | `out of scope` line scoped by CHANNEL ("any listed interpreter taking its code on stdin … however spelled"); `predicate` (`-` is not a listed flag) | A4 + A5. Consistent — and this row is why the line is channel-scoped: `python` IS on the deny-list, so a spellings-scoped line would contradict the predicate here |
| 4 — REQ-74 probes and `table.Categories()` unchanged | `report:` line ("category string, remediation and position … are unchanged"); Naming LBD | S6: `category.go` untouched; edit confined to `load.go`. Consistent |

No CONTRADICTION row: every assertion pair that meets on one output surface agrees, and the one coupling that could skew (clause 3 vs clause 4) is structurally prevented by both readers consuming a single `interpreterForm` call.


### Phase 1: Predicate

Replace `interpreterForm`'s body with the position-free scan and delete the `env` walk; extend the refusal detail with the matched words; leave `shellInterpreters`, `filepathBase`, category constants and `Categories()` untouched.

### Phase 2: Probes

Add the wrapper mutants and the admitted-form probes beside the REQ-74 test; keep the four existing probes as-is.

### Phase 3: Promise wording

Update the C5 comments in `load.go` to cite this record's C1, and the lint's user-facing description of `command_shell_interpreter` to state the predicate and the two out-of-scope forms in C1's words.

## Validation

### Testing Strategy

Unit tests beside the existing C5 probes in
`internal/table/command_carrier_0025_test.go`, driven through `table.Load` so
each scenario asserts the CATEGORY a model author actually sees. Done = every
scenario below green, the four REQ-74 probes and
`neg/neg-command-shell-interpreter.toml` unedited and still green, and
`go test ./internal/table/` clean. Scenarios 1–3 are the predicate; 4–5 are
the promise, which is the half a predicate test cannot reach.

1. **Scenario**: One mutant per wrapper the Problem Statement names — `env -i`,
   `env -u FOO`, `nice`, `timeout 5`, `xargs`, `nohup`, `setsid`, `stdbuf -o0`,
   `chpst`, `doas` — each wrapping `sh -c`.
   **Expected**: each refuses `command_shell_interpreter`, detail naming the
   matched words `sh -c` and the script remediation. Backed by A3's spike, which
   computed all ten as `false → true` under the widened predicate
   (`evidence/spikes/a3-wrapper-argv-transparency.md`); the wrapper binary need
   not exist on the test host, since the predicate reads argv words only.

2. **Scenario**: `["nice","sh","-c","cat {artifact}"]` — a brace-bearing command
   string under a wrapper.
   **Expected**: `command_shell_interpreter`, NOT `command_unknown_placeholder`.
   This is the wrong-defect masking 0025 deviation D18 left open under a
   wrapper; it passes only because clause 3's exemption reads the same
   `isInterp` the widened predicate now sets (`load.go:1095/1114/1124`, A1).

3. **Scenario**: The `env`-option forms intrastate#q2q1 enumerates — `-i`,
   `-u FOO`, `--unset=FOO`, `-0`, `-C DIR`, `--chdir=DIR`, `--`, and a mixed
   `-i A=1 -- sh -c` chain.
   **Expected**: all refuse `command_shell_interpreter`, proving the deleted
   `env` walk is subsumed rather than merely removed (Approach; verified in
   `evidence/research/stage4-scope-wording.md` §Q3).

4. **Scenario**: The admitted forms — `["sh","./gate.sh"]`,
   `["env","-S","sh -c echo"]`, `["sh","-s"]`, `["sh"]`, `["sh","-es"]`,
   `["python","-"]`, `["node","-"]`.
   **Expected**: every one lints GREEN. These are C1's out-of-scope line as a
   test: a future predicate change that starts refusing one of them has broken
   the promise, not tightened it. `["python","-"]` and `["sh","-es"]` are the
   discriminating cases — they fail under any implementation that reads the
   out-of-scope line as shell-spellings-only (A4).

5. **Scenario**: The regression set for the accepted false-refusal class —
   `["ruby","tool.rb","-e","prod"]` (a script whose own argument is `-e`).
   **Expected**: refuses, and the detail names `ruby -e` — the two matched
   words — so the author can see which pair collided. The refusal is
   already-accepted behaviour (D-selection-predicate); the message is what this
   record adds, and the premortem names its absence as the failure mode.

6. **Scenario**: `table.Categories()` membership and order, and the
   `command_shell_interpreter` wire string.
   **Expected**: unchanged. The edit touches `load.go` only; the category
   constant and its position live in `category.go:60`/`:103` and are not
   edited (C1 `report:`, Naming LBD).

### Performance Expectations

Omitted: no alternative in this record was weighed on empirical performance
grounds. For the record the predicate is O(len(argv)²) worst case against
O(len(argv)) today, over vectors bounded by a hand-authored model entry
(single digits of words), run once per entry at load time — the scan is not
on any hot path and no alternative was rejected for cost.

## Finalization Gate

> Complete each item with a written response in
> `{ARTIFACT_DIR}/gate.md` before marking this RDR as
> **Final**. Written responses prevent rubber-stamping
> and produce a review record.
>
> First run the mechanical pre-sweep
> (`prompts/gate/tooling-pass.md`): TEMPLATE section
> coverage, Method-label vocabulary, `Source Search`
> self-reference, `Docs Only` on load-bearing claims. It
> catches what the review rounds disturbed; resolve any
> BLOCK before the written responses.
>
> At lock, replace Contradiction Check, Assumption
> Verification, Scope Verification and Proportionality
> with the one-line pointer to gate.md — those four
> judge THIS record at THIS lock and no peer cites
> them. **Cross-Cutting Concerns stays here**, below
> the pointer: it names the project-wide policy other
> RDRs conform to, so it must stay projected and
> citable as `cli/NNNN:G-cross-cutting`. Cite it that
> way, not by section name.

### Contradiction Check

[Gate key: contradiction — a gate response is cited as
`cli/NNNN:G-<key>`, so the key is a stable id and is
not derived from this heading, which may be reworded.]

[State any conflicts between Research Findings and
the Proposed Solution. If none exist, state
"No contradictions found between research findings,
design principles, and proposed solution."]

### Assumption Verification

[Gate key: assumptions]

[Confirm every Critical Assumption Evidence Record
is internally consistent: Status, Method, and
Evidence agree, and "If wrong" is non-empty. List
any record whose Method is `Docs Only` (these block
lock unless paired with a Spike or Source Search
plan) and any that remain `Pending` or `Unverified`
with a plan to verify before implementation begins.
Confirm no `Verified` stamp is self-referential or
proves only an adjacent claim, and that each cited
`path::Symbol` resolves on `main`. **Status
consistency:** no assumption marked `Pending` or
`Unverified` may have settled-fact prose elsewhere in
the RDR depending on it.]

### Scope Verification

[Gate key: scope]

[Confirm the Minimum Viable Validation is in scope
and will be executed during implementation, not
deferred. State the specific test or proof.]

### Cross-Cutting Concerns

[Gate key: cross-cutting]

[Retained at lock — this sub-section stays in the RDR
when the other gate responses move to gate.md, because
peer RDRs cite it as `cli/NNNN:G-cross-cutting` and an
element that is not projected cannot be cited.]

[List only concerns that apply to this RDR. For each,
state either how this RDR addresses it, or which peer
RDR owns the project-wide policy this RDR conforms
to. Omit (rather than N/A-bullet) anything that does
not apply.]

Candidate concerns (include only those that apply):
versioning · build tool compatibility · licensing ·
deployment model · IDE compatibility · incremental
adoption · secret/credential lifecycle · memory
management · concurrency model · character encoding ·
canonical-form / determinism (see note below).

If this RDR claims byte-identical output,
content-addressed identity, or replay-stable hashes,
also confirm: hash function + library, pre-image
byte layout, primitive encodings, map iteration order,
whitespace policy, case folding, empty/null/absent
distinguishability, and a version marker for future
evolution.

### Proportionality

[Gate key: proportionality]

[Is the document right-sized for the change? Flag
any sections that should be trimmed before locking.
The split test is **contract count, not word count**:
confirm this RDR is the sole author of at most one
independent load-bearing contract (per the Normative
Contracts split signal). If it owns more than one
seam, flag it for splitting rather than locking the
seams together.

Re-validate the **Profile** Metadata field against the
contracts you just counted: confirm the value Resolve
wrote still matches (one contract + no user-facing
surface → `small`; etc. per the applicability matrix).
If the lenses that actually ran disagree with the
Profile (e.g. Profile says `small` but the change locks
a contract that warranted `mid`+ lenses, or the lenses
were skipped on a wrong `small`), correct the field and
do not lock until the missing lenses have run. This is
the latch's backstop — a wrong Profile cannot route
past the lens battery undetected. A `Transient`-marked
contract with a named deleting sibling and schedule is a
recorded lifespan disposition, not an under-sized
Profile — do not count it when re-deriving. Also confirm form:
value + one clause naming the contract(s); strip any
matrix/provenance prose left from the template or Seed
(it belongs in the template comment, not the instance).]

## References

- 0025:C5 (the deny-list line this C1 succeeds) and 0025:§normative-contracts prose beside it; 0025 `artifacts/verification.md` Phase 3a C5 row; 0025 `evidence/critique/Charted.md` (stdin appetite)
- 0004:ALT2, 0004:ALT3 (name lists rejected as authority)
- `internal/table/load.go::interpreterForm`, `::carrierDefect`, `::shellInterpreters`, `::filepathBase`; `internal/cli/cmdbind/cmdbind.go::resolveArgv0`; `internal/table/command_carrier_0025_test.go::TestReq74_KnownInterpreterWithInlineCodeIsALoadTimeDefect`
- Go `os/exec` package documentation (no shell, no redirection); gh `pkg/cmd/alias/set/set.go::NewCmdSet`; roborev `internal/daemon/hooks.go::(*HookRunner).runHook`; beads `internal/creds/command.go::(CommandSource).Resolve`
- intrastate#zvtg (tracker), intrastate#q2q1 (the `env` walk this record decides is subsumed; the tracker left that open), intrastate#b84g (clause 3, closed)
- Stage 2 evidence: `evidence/research/prior-art.md`
