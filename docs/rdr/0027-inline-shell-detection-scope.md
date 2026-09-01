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
- **Profile**: mid — one contract (the predicate shape and claim wording of 0025:C5's `command_shell_interpreter` deny-list), locking what a reviewer-visible lint promises.
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
argv-position-aware wrapper class, and `sh -s` / `sh <script` / `env -S "sh -c …"` are
a different axis — shell-string forms no argv-level check can see without modelling
the wrapper's own parser.

The decision this RDR owns, exactly once: what is the *authority surface* of C5's
deny-list, and therefore what does it promise? The triage names three defensible
answers — (1) an argv0-name heuristic, narrowed in wording and with the wrapper and
shell-string forms declared out of scope by name; (2) an argv-position-aware wrapper
model that closes the class rather than enumerating binaries, with a stated bound on
the wrapper set; (3) retiring argv as the authority surface in favour of the charted
`stdin = "none"|"envelope"` successor, keeping the deny-list as an explicitly
non-load-bearing hint. It is a policy/scope decision about predicate shape plus the
claim wording verification is held to — not an implementation detail — and is not
decided here.

## Critical Assumptions

- **A1 `interpreterForm` has exactly two consumers, both inside `carrierDefect`: clause 3's whitespace exemption (the `isInterp` read) and clause 4's refusal — so one predicate change reaches both clauses and nothing else.**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: `internal/table/load.go::carrierDefect` — the sole call site; a repo-wide symbol search returns no other consumer (Resolve re-runs it).
  - **If wrong**: a second consumer keeps the argv0-only shape and the two disagree on which defect a wrapper form reports.
- **A2 No binding that lints green today carries a listed interpreter basename as a non-argv0 word followed later by one of that interpreter's flags — widening the predicate refuses nothing currently accepted in the repo's fixtures or 0025's verification models.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: pending — `intrastate lint` over every `command`-carrying fixture under `internal/table/testdata/` and the 0025 MVV models with the C1 predicate; expected: zero new refusals.
  - **If wrong**: a green model turns red on upgrade and the refusal's remediation ("put it in a script") misdirects, because there is no script to extract.
- **A3 Every wrapper form the Problem Statement names (`env -i`, `env -u FOO`, `nice`, `timeout 5`, `xargs`, `nohup`, `setsid`, `stdbuf`, `chpst`, `doas`) places the interpreter name and its inline-code flag as separate argv words in that order, so the position-free predicate refuses each — and a brace-bearing string under a wrapper reports `command_shell_interpreter`, not `command_unknown_placeholder`.**
  - **Status**: Pending
  - **Method**: MVV Test
  - **Evidence**: the MVV's wrapper-mutant set (steps 2–3), one mutant per named wrapper plus one carrying `{artifact}` inside the string.
  - **If wrong**: a named wrapper form still lints green and C1's promise is false at lock; or the wrapper form masks the interpreter defect (0025's D6 failure re-opened).
- **A4 From a fixed argv executed without a shell, the only forms that spell or source shell code outside separate argv words are a one-word shell string (`env -S "sh -c …"`, or a single `"sh -c …"` element handed to a tool that re-splits it) and a stdin-fed shell (`sh -s`, bare `sh`); `sh <script` is unreachable (no redirection without a shell) and `sh script.sh` is the sanctioned wrapper-file form.**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: pending — run each form through the executor's spawn path and record whether a shell executes and whether `lint` admitted it; the out-of-scope list in C1 must equal the admitted-and-executes set.
  - **If wrong**: a reachable form outside C1's named out-of-scope list executes a shell and the "named by name" claim is incomplete.
- **A5 The `sh -s` / bare-`sh` hazard is owned by the charted `stdin = "none" | "envelope"` successor: withholding the envelope from a command that never declared it leaves a stdin-reading shell nothing to execute, so routing that axis there leaves no reviewer-visible gap unowned.**
  - **Status**: Pending
  - **Method**: Design Decision
  - **Evidence**: this RDR's scoping choice, resting on 0025's `evidence/critique/Charted.md` ("not statically detectable from argv, so no C5 arm can catch it"); until the successor ships, C1 names the form admitted and the reviewer reads argv.
  - **If wrong**: the successor never ships or withholds nothing, and `sh -s` stays an unowned admitted shell.

## Proposed Solution

### Approach

Decide once what 0025:C5's deny-list is the authority over: **argv words, at any position** — never the contents of a word, stdin, files, or the resolved binary. The predicate is freed from argv0: a binding is refused when a listed interpreter word is followed, at any later argv position, by one of that interpreter's inline-code flags. That closes the whole wrapper class (`env -i`, `env -u FOO`, `nice`, `timeout 5`, `xargs`, `nohup`, `setsid`, `stdbuf`, `chpst`, `doas`, and every wrapper not yet thought of) with **no wrapper table and no option grammar**, because nothing before the interpreter word is read at all. The `env`-chain walk in `internal/table/load.go::interpreterForm` is deleted as subsumed (intrastate#q2q1 may land first; the walk goes either way).

The claim wording is narrowed to exactly that predicate, and the two forms an argv-level check cannot see are declared out of scope **by name**: a shell string carried inside one word (`env -S "sh -c …"`), and a shell that reads its script from stdin (`sh -s`, bare `sh`), the latter routed to the charted `stdin = "none" | "envelope"` successor, which withholds the only script source a fixed argv can reach. `sh script.sh` stays the sanctioned wrapper-file form.

What the lint then promises a reviewer is one sentence they can hold: *if an interpreter and its inline-code flag appear as two separate words in that order, anywhere in argv, the binding is refused; a one-word shell string or a stdin script you read yourself.* This keeps C5's own frame — "the interpreter deny-list is defense in depth over that, not the barrier itself" (0025:§normative-contracts, beside C5) ⇒ the deny-list is never re-promoted to a barrier, so no wrapper table is justified; but the cost it raises is no longer defeated by one leading word, which is what a check "whose purpose is visibility" owes the reader.

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
out of scope, BY NAME (admitted by lint; a shell may run): a shell string carried in ONE word (`env -S "sh -c …"`, or a single `"sh -c …"` element handed to a tool that re-splits it); a stdin-fed shell (`sh -s`, bare `sh`) — owned by the charted `stdin = "none" | "envelope"` successor. `sh script.sh` is the sanctioned wrapper-file form and never a defect
interpreter set: OPEN (deny-listed, not closed), unchanged from 0025:C5 — an unlisted spelling is admitted; the list grows by amendment of THIS clause
clause 3 coupling: 0025:C5 clause 3's whitespace exemption keys on this predicate, so a brace-bearing string under a wrapper reports the interpreter form, never command_unknown_placeholder; within-entry precedence is 0025:C5's, unchanged
report:     category string, remediation ("inline shell is not a declared command; put it in a script and declare the script as argv0") and position in `table.Categories()` are unchanged; the detail additionally names the two matched words
promise:    what a reviewer may rely on is the predicate line and nothing more; the lint's user-facing description states the out-of-scope forms in the words above
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

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Withholding stdin from a command that never declared it (so a stdin-fed shell has no script) | Future — the charted `stdin = "none" \| "envelope"` successor | Deferred | C1 names `sh -s` out of scope and routes it there; until it ships the form is admitted and named as such |
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

Joint-check: fired → 0028 (home: cli/0025:C5) — symmetric to 0028's arm-1 fire on `internal/table/load.go::carrierDefect` / `table.Categories()`, disposed cite-don't-restate: C1 here owns only clause 4's predicate, 0028:C1/C4 own the `edit` arm and the tail categories, and both fence deltas against 0025:C5, which owns the clause map and within-entry precedence. Recorded at the batch propose of 2026-08-31 after this record's own check ran clear: all three arms run. Arm 1 (modify-anchors) and arm 2 (contract literals), open-only and repo-resolved: no overlaps. Arm 3 (absence, manual): the 27 Draft/Final/Implemented peers grepped for the admitted-form and wrong-defect tokens (`interpreterForm`, `env chain`, `env -i`, `xargs`, `sh -s`, `env -S`, `command_unknown_placeholder`, `command_shell_interpreter`); hits only in 0025:C5 — the line C1 succeeds, cited not relied on — and 0004's "without raw shell strings", which C1 keeps true. 0026 compared on its written contracts (0026:C1 amends 0025:C4's drain precedence; disjoint from C5); 0028 compared on its problem statement (an in-process `edit` carrier that runs no shell, no subprocess); no shared decision, no bridge surface.

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

Constraints carried from the triage: the `env` option-flag walk (`-i`, `-u NAME`,
`-0`, `-C`, `--`) is intrastate#q2q1 and ships independently without foreclosing any
answer; answer 3 may fold into the charted stdin/envelope successor, so that successor
should be paired with, not duplicated by, this RDR; whether `env -S "sh -c …"` is in
scope at all is an open question the chosen answer must settle.

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
- Negative: `env -S "sh -c …"` and `sh -s` are admitted and documented as such; the second waits on the stdin successor.

### Risks and Mitigations

- **Risk**: a green binding somewhere carries an interpreter word mid-argv followed by its flag and turns red on upgrade.
  **Mitigation**: A2's spike over all fixtures and 0025's models before lock; the named-words message makes any survivor self-explanatory.
- **Risk**: "closes the wrapper class" is read as "closes inline shell" and reviewers stop reading argv.
  **Mitigation**: C1's out-of-scope line is normative and the lint's description carries it verbatim; the docs never say "closes inline shell".

### Failure Modes

- Visible: a refused wrapper form names `sh -c` (the matched words) and the wrapper-file remediation; a false refusal names the two words that matched, so the author sees the script argument that collided.
- Silent: a one-word shell string or a stdin-fed shell lints green — by contract, named in C1; diagnosis is reading argv, which C1's promise tells the reviewer to do for exactly those shapes.
- Recovery: none needed for the predicate (no state); a wrong refusal is fixed by renaming the script's flag or wrapping, as the message says.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] intrastate#q2q1's disposition recorded (landed first, or closed as subsumed) — either way the `env` walk is removed here

### Minimum Viable Validation

1. Take the 0025 command-carrier model with a green `command` read binding; `intrastate lint` passes.
2. Swap in one mutant per named wrapper — `["env","-i","sh","-c","…"]`, `["env","-u","FOO","sh","-c","…"]`, `["nice","sh","-c","…"]`, `["timeout","5","sh","-c","…"]`, `["xargs","sh","-c","…"]`, `["doas","sh","-c","…"]` — and one carrying `{artifact}` inside the string under `nice`; each refuses with `command_shell_interpreter`, the detail naming `sh -c` and the script remediation; the `{artifact}` mutant does NOT report `command_unknown_placeholder`.
3. Swap in the admitted forms `["sh","./gate.sh"]`, `["env","-S","sh -c echo"]`, `["sh","-s"]`; lint passes and the first two are shown, by running the model under `--allow-commands`, to behave as C1 states (the wrapper file runs; the `-S` string runs a shell).
4. The original 0025 REQ-74 probes and the negative fixture still refuse with the same category, and `table.Categories()` order is unchanged.
End-state: every wrapper mutant red with the right defect, every admitted form green and documented, no existing probe changed.

### Phase 1: Predicate

Replace `interpreterForm`'s body with the position-free scan and delete the `env` walk; extend the refusal detail with the matched words; leave `shellInterpreters`, `filepathBase`, category constants and `Categories()` untouched.

### Phase 2: Probes

Add the wrapper mutants and the admitted-form probes beside the REQ-74 test; keep the four existing probes as-is.

### Phase 3: Promise wording

Update the C5 comments in `load.go` to cite this record's C1, and the lint's user-facing description of `command_shell_interpreter` to state the predicate and the two out-of-scope forms in C1's words.

## Validation

### Testing Strategy

[Required — never omit. Test scenarios and coverage goals — what to test and
what constitutes "done." For non-functional concerns
(performance, security): state measurement strategy,
not estimates.]

1. **Scenario**: [Description]
   **Expected**: [Result]

### Performance Expectations

[Conditional — omit (don't N/A-bullet) this section unless
comparing alternatives on empirical performance grounds.
Do not include effort estimates or speculative
throughput targets. Rough performance metrics are
appropriate only when comparing alternatives — note
empirical data or obvious gains that support the
chosen approach over a rejected one.]

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
- intrastate#zvtg (tracker), intrastate#q2q1 (subsumed `env` walk), intrastate#b84g (clause 3, closed)
- Stage 2 evidence: `evidence/research/prior-art.md`
