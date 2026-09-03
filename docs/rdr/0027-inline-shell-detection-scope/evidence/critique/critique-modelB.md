Model: claude-sonnet-5

# Critique — RDR 0027 (inline-shell detection scope)

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0027:A6` | RDR locks with a `Pending` assumption that gates a normative clause (`C1`'s `promise:` line) and an entire implementation phase (Phase 3); the "spike" is deferred to whoever implements, with no owner, no deadline, and no fallback if no reviewer-reachable surface exists | Phase 3 either ships nothing (promise text never lands anywhere a reviewer reads it) or lands in a place nobody reads (a doc comment), and the RDR's own `promise:` clause is false the day it ships | §1, §2, premortem, AT-1 |
| C-2 | `0027:A5` | The stdin-fed axis (`sh -s`, bare `sh`, `python -`, `node -`) is declared out of scope "by name" and routed to a successor (`stdin = "none"\|"envelope"`) that A5 itself admits no kata tracks — a documented-but-unscheduled deferral | The exact bypass the Background section says was "verified to actually execute the shell" (`env -i sh -c`) is closed, but its sibling stdin bypass (`sh -s`) ships forever-admitted with a paper IOU that nothing forces to be redeemed | §1, §3, premortem, AT-2 |
| C-3 | `0027:D-selection-predicate` | The lowest-i-then-lowest-j tie-break reports the outer wrapper interpreter (`python -c`) instead of the one a human reads as operative (`sh -c`) when two listed interpreters chain in one argv; accepted as "message attribution only" with no test, no spike, and no worked example anywhere in the record | A reviewer or script author sees a refusal citing `python -c` for a binding whose actual hazard is `sh -c`, chases the wrong interpreter, and — worse — this exact shape (`interpreter1 interpreter2 -c ...`) is exactly what a defeated wrapper table would have forced through anyway, so the one class C1 is proudest of closing produces its worst diagnostic on a legitimate closure case | §1, premortem, AT-3 |
| C-4 | `0027:C1` (`interpreter set: OPEN`) | The deny-list stays an OPEN, unfolded, exact-match set (`shellInterpreters`) with no version marker and no forward pointer from this record to where it is amended; C1's own text says "the list grows by amendment of THIS clause" — meaning every future interpreter (`uv run`, `deno`, `bun`, `busybox`, `python3`) triggers a reopen of this exact RDR rather than a lighter-weight table update | First bug report after ship: "why does `uv run python -c ...` lint clean, isn't this what 0027 was for?" — followed by a PR that reopens and re-locks 0027 rather than editing a data table, because the record made itself the amendment surface | §2, premortem, AT-4 |
| C-5 | `0027:§technical-design` (Illustrative Code) / `0027:C1` (predicate line) | The position-free O(n²) scan matches ANY listed interpreter word to ANY later flag word in the same argv, with no adjacency and no "belongs to the same sub-command" boundary — so `["build-tool","--interpreter","ruby","run-step","--flag","-e"]`-shaped tooling wrappers (unrelated to shell invocation) will false-positive the instant a real-world command tool happens to use `-e` or `-c` as its own flag anywhere after a listed interpreter *name* used as a plain argument, not a program | A CI/build binding that legitimately passes the literal string `"ruby"` or `"python"` as a non-executable argument (e.g. `["multi-tool","--lang","ruby","--emit","-e"]`) refuses with `command_shell_interpreter`, and the remediation text ("put it in a script") is nonsensical because there is no shell invocation to extract | §1, premortem, AT-5 |
| C-6 | `0027:A2` | A2's "no new refusals" verification is scoped to *this repo's fixtures and 0025's models* — 32 vectors, 11 green — and the RDR explicitly disclaims coverage of "an out-of-tree user model." Given C-5's false-positive shape, the actual blast radius of the widened predicate on real user models is unverified and unverifiable at lock | The first out-of-tree adopter to upgrade intrastate hits a green-to-red flip A2 was never positioned to catch, and the RDR's own risk mitigation ("A2's spike... before lock") has already run and cannot re-run against models that don't exist yet | §1, §3, premortem |
| C-7 | `0027:§performance-expectations` | Predicate complexity is stated as O(len(argv)²) worst case, dismissed as irrelevant because argv is "single digits of words" in a "hand-authored model entry" — but nothing in C1's contract bounds argv length, and the predicate is now blind to any structural signal (no wrapper table, reads every word against every other word) | A generated or templated model (not hand-authored) with a long argv — plausible once `command` bindings are produced by tooling rather than typed by a human, which is exactly the direction a "reviewer visibility" check is meant to survive — pays quadratic cost with no test or contract line to catch a regression | §2, AT-6 |
| C-8 | `0027:§metadata` (Profile) / `0027:A6` | Profile is declared `large` for "the successor predicate and claim wording," but the record locks two independently-testable surfaces under one Profile: (1) the predicate (`load.go`, fully verified, A1-A4 all `Verified`) and (2) a brand-new user-facing description surface that doesn't exist yet in the codebase (A6 `Pending`, Phase 3, S7) — the Proportionality gate's own split test ("sole author of at most one independent load-bearing contract") is never actually run against this decomposition in the projected record | Stage 8 implementer ships the predicate (fully specified, low risk) and treats Phase 3/A6 as an afterthought bolt-on because the record never forced the two-contract question into the open; the description surface either ships late, ships wrong, or ships silently dropped, and nothing in the gate history shows this was ever weighed | §2, premortem, AT-1 |
| C-9 | `0027:§normative-contracts` (`report:` line) | "the detail additionally names the two matched words" is a normative commitment with no format contract — no specified punctuation, ordering, or stability guarantee for the string (e.g. is it always `"word1 word2"` separated by a space, forever?) — while `table.Categories()` order IS pinned as a golden assertion (S6). The asymmetry means the part of the contract most likely to be read by tooling (a CI grep for the refusal category) is protected, while the part most likely to be read by a *human* (the detail text a reviewer diagnoses from) has no stability contract at all | Six weeks post-ship, someone adjusts the detail message wording for clarity (exactly the kind of copy-edit this record invites — "so a false refusal... is diagnosable from the message") and nothing in the test suite or the RDR notices, because S5's assertion is "the detail text contains both words," not an exact format | §2, AT-3 |

## 1. Three most likely ways implementation goes wrong

**(a) The promise-text surface (A6) ships nowhere reviewer-reachable, or ships late and silently.**

Root cause: A6 is `Pending` at lock time — the RDR's own metadata admits the load-bearing question "where does this text live" is unanswered. The Technical Design's `promise:` clause in `C1` (`0027:C1`) states as fact that "the lint's user-facing description states the out-of-scope forms" and that "Phase 3 creates it and THESE words are the text it ships" — present tense, normative, as if the surface already exists. It does not. The Existing Infrastructure Audit and A6's own evidence agree: `Category` is a bare string, `Categories()` returns identifiers only, and "a doc-comment beside the constant is not reviewer-reachable." Phase 3 (`0027:§phase-3-promise-wording`) punts the actual decision to the implementer: "Where it lands is the implementer's call." A decision this load-bearing — it is the entire second half of the RDR's stated purpose ("what the lint promises a reviewer") — is being handed to whoever happens to pick up the Stage 8 ticket, six weeks or six months from now, with no api/plumbing sketch, no candidate location ranked, and no fallback if no such surface exists in the current CLI output contract.

Symptom: either Phase 3 is dropped silently under implementation-time pressure ("the predicate is the hard part, ship that, describe later" — a extremely common trajectory for a "and also write some docs" phase tacked onto a code-focused RDR), or it lands in a place a reviewer never sees (a `--help` subcommand nobody runs, a doc comment, a CHANGELOG entry), and the RDR is retroactively falsified: the record's central rhetorical move — "what the lint then promises a reviewer is one sentence they can hold" — becomes a sentence that exists in the RDR file and nowhere else.

**(b) The tie-break message names the wrong interpreter for chained-interpreter argv, and nobody notices until a real refusal is misdiagnosed.**

Root cause: `0027:D-selection-predicate` accepts, without a test, spike, or worked scenario anywhere in the record, that `["python","sh","-c","echo"]` reports `python -c` rather than `sh -c`. The stated justification — "a tie-break that preferred the *nearest* pair would make the reported form depend on scanning order rather than on argv position" — is an internal-consistency argument about the *predicate's* determinism, not an argument about what serves the *reader*. The RDR's entire premise is that the deny-list's job is reviewer visibility ("a check whose purpose is visibility," repeated three times across Background, Approach, and Decision Rationale). A diagnostic message that reliably points at the decoy interpreter instead of the operative one directly undermines that stated purpose, and the record accepted it as "message attribution only" — i.e., a known, named, deliberately-shipped defect in the one artifact (the detail string) whose entire job is to make the hazard legible.

Symptom: a reviewer or script author gets a refusal reading `command_shell_interpreter: ... declares the interpreter form python -c ...` for a binding whose actual shell hazard is `sh -c` three tokens later. They rename or edit the `python` reference (wrong fix), re-run, get refused again on `sh -c` (now surfaced as the lowest remaining pair), and burn a cycle discovering the message pointed at the wrong culprit the first time. This is exactly the "reviewer trust in a check whose purpose is visibility" cost the RDR's own Background section names as the thing at stake — inflicted by the RDR's own accepted design choice.

**(c) The position-free scan false-positives on non-shell tooling that happens to use a listed interpreter *name* as a plain string argument, or a listed flag spelling for an unrelated purpose.**

Root cause: `0027:C1`'s predicate is stated as "refuse iff there exist i < j with base(argv[i]) a listed interpreter and argv[j] one of its listed flags... Nothing before argv[i] is read." This is deliberately unstructured — no adjacency requirement, no sub-command boundary, no requirement that argv[i] be an executable position at all. A2's verification (`0027:A2`) checked this is safe against *this repo's* 32 fixture/model vectors and found zero new refusals — but A2 explicitly disclaims generality: "Scope: this repo's fixtures and 0025's models, which is exactly the claim's scope; it says nothing about an out-of-tree user model." The illustrative code (`0027:§illustrative-code`) confirms the scan is a flat double loop over all argv words with no positional semantics beyond ordering. Any `command` binding that legitimately carries the literal token `"ruby"`, `"python"`, `"node"`, `"php"`, `"sh"`, `"bash"`, etc. as a *data argument* (a `--lang ruby` flag to some polyglot build tool, a `--interpreter-name sh` diagnostic flag, a test fixture name) followed anywhere later by a short flag that happens to collide with `-c`/`-e`/`-r`/`--eval` is now refused with a remediation ("put it in a script and declare the script as argv0") that makes no sense because there is no shell invocation to extract.

Symptom: a model author with a legitimate, non-shell `command` binding gets a load-time refusal citing `command_shell_interpreter` and a remediation that doesn't apply to their case. They cannot silence it (no opt-in exists per 0025:C5's "no opt-in in v1," which this RDR reaffirms), so they must restructure their argv to avoid the collision — a false-positive tax paid by users who did nothing wrong, imposed by a predicate this record chose specifically because it required "no wrapper table and no option grammar," i.e., no way to bound what counts as "the interpreter position."

## 2. The section that will be rewritten within 6 weeks of shipping

`0027:C1`'s `interpreter set: OPEN` clause, backed by `shellInterpreters` (`internal/table/load.go:1035-1047`).

This is the highest-confidence reopen candidate, for a reason the record itself names but doesn't act on: C1 states "the list grows by amendment of THIS clause" — meaning the *membership* of the deny-list is normatively owned by this RDR's text, not merely by the Go map. The moment a real user hits `uv run python -c ...`, `python3 -c ...` (explicitly named in A2's evidence as staying green under both predicates — "the widen is over argv POSITION, not over which basenames are listed, so an unlisted spelling is unaffected"), `deno eval`, `bun -e`, or any interpreter not in the eleven-entry map, the fix is a one-line data change to a Go map — but the RDR's own text has made itself the amendment authority for that line ("the list grows by amendment of THIS clause"). That forces a choice at implementation-adjacent time between (a) reopening and re-locking a "large" Architecture RDR to add one map entry, or (b) quietly editing the map and letting the RDR's normative text drift out of sync with `shellInterpreters`, which is precisely the "frozen-at-lock invariant without a version marker" failure this lens is built to catch. Nothing in C1 provides a lighter-weight amendment path (a referenced, separately-versioned data file; an explicit "this list may be extended without RDR amendment, see X" carve-out); the record locks the predicate *shape* well but locks the predicate *content* the same way, and content changes far more often than shape.

The competing candidate is A6/Phase 3 (the promise-text landing site) — but that's better classified as "never actually ships" (see §3/§4) than "rewritten": a `Pending` assumption that never resolves doesn't get revisited, it gets quietly dropped. The deny-list, by contrast, WILL get a concrete forcing function (a real user's real interpreter) within weeks of any meaningful adoption, because C5's own frame (inherited here) is that the list is already known to be incomplete by design ("OPEN, not closed").

## 3. The assumption that will not survive first contact with a real user

`0027:A5`: "The stdin-fed-interpreter hazard... is owned by the charted `stdin = "none" | "envelope"` successor... so routing that axis there leaves no reviewer-visible gap unowned."

This assumption is marked `Verified` by Method "Design Decision" — which is itself a tell: it is not verified against anything external, it is a *scoping choice* dressed as a verified fact. A5's own "Durability caveat" undercuts its status inline: "no kata tracks the successor — it exists only as that `Charted.md` line... so the If-wrong below is nearer than a charted item normally implies." The record is telling the reader, in the assumption's own evidence field, that the assumption is fragile — and then stamping it `Verified` anyway, because the *design decision* to route the hazard there is coherent, even though the thing it's routed to may never arrive.

The Background section states plainly that this exact bypass class was found by roborev triage with a real lint run and that `env -i sh -c 'echo …'` "was verified to actually execute the shell." `sh -s` and bare `sh` are the same severity of hazard (an interpreter admitted by the deny-list, executing attacker/author-supplied code, invisible to the reviewer) and this RDR ships with that hazard *known, named, and left open*, on the strength of a citation to a document (`Charted.md`) that is not a kata, has no owner, and has no schedule. A real user — specifically, the reviewer this entire RDR is written to protect ("A reviewer reading a model under `--allow-commands`...") — will read C1's promise ("if an interpreter and its inline-code flag appear as two separate words... anywhere in argv, the binding is refused") and reasonably build a mental model that inline shell is caught. The very next sentence names the exception, but the RDR's own premortem (`0027:§decision-rationale`) already predicts a reviewer "who read 'closes the wrapper class' stops reading argv" — and A5 is the load-bearing bet that the *promise wording alone* (not a shipped fix) will keep that reviewer reading argv for the stdin cases indefinitely. That bet does not survive contact with a reviewer who trusts the tool, which is the exact reviewer this RDR exists to serve.

## 4. Premortem

It is six weeks after RDR 0027 shipped. `intrastate lint` now refuses `env -i sh -c`, `nice sh -c`, `timeout 5 sh -c`, `xargs sh -c`, and the rest of the wrapper-prefix class, all via the widened `interpreterForm` in `internal/table/load.go`. The predicate change worked exactly as specified: `carrierDefect`'s clause 4 fires on every named wrapper mutant, clause 3's placeholder exemption tracks it correctly, and `table.Categories()` order is untouched. On that narrow slice, the RDR delivered.

Three things did not go as the record implied they would.

First, a support thread opens: a platform team building a polyglot CI orchestration model has a `command` binding like `["polybuild","--target","service-a","--runtime","ruby","--stage","-e"]` — `polybuild` is their own internal build tool, `--runtime ruby` is a plain string argument selecting a language profile, and `--stage -e` selects a pipeline stage abbreviated `-e` for "emit." `intrastate lint` refuses it with `command_shell_interpreter`, citing `ruby -e`, and the remediation says "put it in a script and declare the script as argv0" — but there is no shell invocation anywhere in this command; `polybuild` never touches a shell. The team cannot silence the refusal (no opt-in in v1, reaffirmed by this RDR) and has to rename their own internal flag naming scheme to dodge intrastate's lint, which they experience as the tool dictating their CLI design. `internal/table/load.go::interpreterForm`'s two-index scan (i for a listed interpreter, j > i for a listed flag, "nothing before argv[i] is read") had no adjacency or sub-command boundary by design — this was flagged internally in review as "position-free, no wrapper table," praised as the elegant part of the design, and is now the exact mechanism producing this false positive. Nobody who reviewed 0027 ran a fixture shaped like this because A2's spike was scoped, by its own text, to "this repo's fixtures and 0025's models" and explicitly declined to speak to "an out-of-tree user model."

Second, the promise text never shipped. Phase 3 was scoped as "the implementer's call" for where the description lands (A6 was `Pending` at lock and stayed unresolved through Stage 8 planning). Under delivery pressure, the Stage 8 implementer shipped Phase 1 (the predicate) and Phase 2 (the probes) — both fully specified, both easy to verify against the MVV — and Phase 3 slipped to "follow-up," then fell off the tracker entirely, because no kata was ever opened to track A6's spike the way A5's stdin-successor gap was at least named in `Charted.md`. Six weeks in, `command_shell_interpreter`'s refusal detail is the only user-facing text that exists, exactly as before 0027 — it fires only on refusal, never on admission, so a reviewer reading a *green* `env -S "sh -c echo"` binding still has no in-tool signal that this is a knowingly-admitted shell invocation. The RDR's central deliverable — "one sentence a reviewer can hold" — exists in `docs/rdr/0027-inline-shell-detection-scope.md` and nowhere a reviewer using the CLI will ever see it.

Third, a security-conscious adopter reads 0027's changelog note ("closes the wrapper bypass class"), takes it at face value, and stops manually grep-ing argv for shell tells in their model review process — the exact habit the Problem Statement opens by describing ("They discover it only by reading argv themselves"). Three months later they admit a model with `["python","-"]` reading a script from stdin injected by a compromised CI step. It lints green, per C1's explicit and correctly-documented design (stdin-fed forms are named out-of-scope by channel). The reviewer's trust regression is not a bug in the predicate — C1 does exactly what it says — it is the predicted failure of A5's bet that promise-wording alone, without a shipped enforcement mechanism, keeps a trusting reviewer reading argv for a hazard the tool's own refusal message doesn't fire on.

## 5. Acceptance tests that would have caught each failure at RDR-review time

```
AT-1 (catches C-1, C-8): Promise-text surface has no owner or fallback
  GIVEN RDR 0027 is proposed for Final status
  AND A6 is still `Pending`
  WHEN the Finalization Gate's Assumption Verification is run
  THEN lock is blocked
  AND the record must either:
    (a) resolve A6 to a concrete, named landing site with a spike showing
        it renders without provoking a refusal, or
    (b) split Phase 3 / the `promise:` clause into a separate RDR whose own
        Profile and MVV are scoped to the description surface,
  BEFORE the predicate half (Phase 1/2) is allowed to lock as `large` covering both.

AT-2 (catches C-2): Stdin-fed axis has no forcing function
  GIVEN A5 routes the stdin-fed hazard to the charted `stdin` successor
  AND that successor has no open kata
  WHEN the Finalization Gate's Cross-Cutting or Risk review runs
  THEN require either:
    (a) a kata is opened and linked from A5's evidence, tracked to a
        target milestone, or
    (b) the RDR's shipped promise text (once A6 resolves) explicitly and
        loudly flags the admitted stdin forms as a KNOWN OPEN HAZARD in the
        refusal-adjacent user-facing text, not merely in the RDR file.

AT-3 (catches C-3, C-9): Tie-break message points at the wrong interpreter
  GIVEN a command argv containing two distinct listed interpreters, e.g.
    ["python", "sh", "-c", "echo hi"]
  WHEN `intrastate lint` evaluates the binding
  THEN the refusal detail SHOULD name the pair a reviewer would independently
    identify as the operative shell hazard (`sh -c`), or at minimum the
    Testing Strategy MUST include this exact scenario with an explicit,
    reviewed Expected value (not silently accepted as "message attribution
    only" with zero test coverage)
  AND the detail-string FORMAT (word order, separator, stability across
    versions) is asserted, not merely "contains both words."

AT-4 (catches C-4): Deny-list membership has no amendment path lighter than RDR reopen
  GIVEN `shellInterpreters` needs a new entry (e.g. "uv", "deno", "bun")
  WHEN a maintainer wants to add it
  THEN the RDR MUST state whether this requires re-opening 0027 to Draft,
    or whether a lighter amendment path exists (e.g. "this map may be
    extended by PR without RDR review; C1 governs shape, not membership")
  AND C1's "the list grows by amendment of THIS clause" line must be
    reconciled with whichever answer is chosen before lock.

AT-5 (catches C-5, C-6): Position-free predicate false-positives on non-shell tooling
  GIVEN a `command` binding for a non-shell tool that carries a listed
    interpreter NAME as a plain string argument (not argv0-of-a-wrapper)
    followed later by a short flag colliding with a listed inline-code flag,
    e.g. ["polybuild", "--runtime", "ruby", "--stage", "-e"]
  WHEN `intrastate lint` evaluates the binding
  THEN the binding SHOULD NOT refuse with `command_shell_interpreter`
  AND if the RDR accepts this as a known false-positive class (rather than
    fixing it), that acceptance MUST be written into Trade-offs/Failure Modes
    as an explicit, named cost — not left undiscovered because A2's spike
    was scoped only to in-repo fixtures.

AT-6 (catches C-7): No regression guard on predicate cost
  GIVEN a `command` argv with a large word count (e.g. 200+ elements, as a
    generated/templated binding might produce)
  WHEN `intrastate lint` loads the model
  THEN load time SHOULD remain bounded by a stated contract or test,
    not merely dismissed as irrelevant because today's fixtures are
    "single digits of words" — the RDR names no upper bound on argv length
    anywhere in C1, so nothing prevents this class of input.
```

