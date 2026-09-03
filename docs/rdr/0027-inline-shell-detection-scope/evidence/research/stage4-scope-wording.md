Model: claude-opus-5[1m]

# Stage 4 research — the two claim-wording questions

Run 2026-09-03. Prior-art / prior-decision sweep for the two questions the
A4 spike raised. Corpus: intrastate RDRs 0001-0028, JDRs 0001-0003, the kata
tracker, shipped source, and the langref peer-CLI checkouts.

## Q1 — should C1's out-of-scope line name the stdin CHANNEL, or stay shell-only?

### Findings

**The charted successor is channel-scoped, not shell-scoped.** 0025
`evidence/critique/Charted.md`: "Declared stdin appetite (`stdin = "none" |
"envelope"`) on a command entry". Its motivating hazard is the spike-proven
`tee {artifact}` case — and `tee` is not a shell. So the successor A5 routes
to already owns what a command entry CONSUMES generally. Naming non-shell
interpreters describes that remit; it does not overstep it.

**The deny-list is already interpreter-general in code, and its own comment
says so.** `internal/table/load.go:1031`: "`shellInterpreters` is C5's OPEN
deny-list: argv0 forms that take inline code on a flag" — an inline-code
definition — over a map containing `python`, `ruby`, `node`, `php`. The same
comment then says "inline-shell carrier". The shell/interpreter conflation is
ALREADY in shipped code; 0027 is the record that decides the authority
surface, so it is where the conflation is resolved or ratified.

**C1's own predicate lines are already interpreter-general.** The `#` line
says "a listed interpreter word"; `predicate:` says "base(argv[i]) a listed
interpreter". Only the `out of scope` line says "shell". The inconsistency is
localised to one line of one clause.

**Both predicates admit every stdin form — this is a shipped wording gap, not
a 0027 regression.** Computed over old and widened predicates:
`["sh","-s"]`, `["sh"]`, `["sh","-es"]`, `["python","-"]`, `["node","-"]`,
`["ruby","-"]`, `["php","-"]` are all `(false,"")` under BOTH. 0027 is simply
the first record to state the promise precisely.

**House precedent: say the gap out loud.** 0025 `artifacts/verification.md`
records the `python3`/`nodejs` spelling gap under "Observation — not a FAIL",
concluding "the spelling gap is one an amendment would likely want to close" —
rather than omitting an admitted form a reader might assume was caught.

**The seed never contemplated non-shell stdin.** kata zvtg names "`sh -s` /
`sh <script` are a different axis again" and its answer-1 sketch would declare
"Wrappers…, `sh -s`, `sh <script`, and `env -S` … out of scope by name". No
mention of `python -`. The spike finding is new information the seed lacked.

**JDR 0003 explicitly declines this.** It is scoped to the `execution_failure`
sub-reason carrier. Its own text: "0027 amends the same declaration seam
(`cli/0025:C5`) and is bound as a member; no entry lands in it today", and
under *What this does not decide*: "The `carrierDefect`/`Categories()` clause
boundary between 0027 and 0028 — stays with `cli/0025:C5`". Not a joint
decision; 0027 owns it.

**Peer prior art: no convention to borrow.** No peer CLI detects shell-ness
from argv (gh declares via `--shell`/`!`; consul via `-shell`, default true;
roborev/beads are always `sh -c`). No peer maintains a mixed shell/non-shell
interpreter-to-inline-flag list — no ecosystem term exists for a list holding
both `sh` and `python`. No peer enumerates a check's gaps "by name" at all.
Closest analog: consul `command/exec/exec.go:78-82` gates stdin-as-script
(`cmd == "-"`) on the SAME declared shell flag — treating stdin-fed script
delivery as one concern with shell mode, not a separate axis.

### Bearing

The predicate does not change under either answer. What changes is whether a
reviewer who reads the promise literally learns that `["python","-"]` is
admitted. Under shell-only wording they do not.

## Q2 — should C1 or A4 pin `env -S` as GNU-only?

### Findings

**Direct precedent, verbatim — 0025:F9:** "Refusing is what keeps that bound
honest rather than platform-dependent; lint stays platform-neutral, so a model
authored on Windows still validates."

Lint runs at LOAD time over a model that may be authored on one platform and
executed on another. A contract clause qualifying an admitted form by the
author's host would make the static claim host-dependent — what F9 forbids —
and would mislead the review case that matters: a reviewer on darwin reading
a model destined for a Linux runner, where GNU env is the default.

**No portability kata exists** (searched `darwin`, `BSD`, `GNU`, `coreutils`,
`portab`). No record pins a specific tool's platform-dependent flag behaviour
into a normative contract.

**C1 already carries a platform-independent second spelling** of the same
form: "a single `\"sh -c …\"` element handed to a tool that re-splits it".
The form is reachable everywhere; `env -S` is one illustrative spelling.

### Bearing

The spike fact (darwin BSD env rejects `-S`; GNU env 9.11 splits and runs it)
is real and belongs in the evidence record. F9 says it does not belong in the
normative claim.

## Q3 — collateral finding: the q2q1 subsumption claim

0027's Approach and References call intrastate#q2q1 "subsumed". The tracker
deliberately did NOT pre-commit that: zvtg's scope-review comment says the
env option-flag walk "ships independently, and does not foreclose any of the
three answers", and q2q1 says "This kata does not foreclose it."

Verified empirically that subsumption HOLDS on the merits — the widened
predicate returns `(true,"sh -c")` for all nine forms q2q1 enumerates
(`-i`, `-u FOO`, `--unset=FOO`, `-0`, `-C /tmp`, `--chdir=/tmp`, `--`,
`A=1`, and a mixed `-i A=1 --` chain). So the claim is correct, but it is
0027's DECISION to make rather than a fact inherited from the tracker.

## Q4 — collateral finding: the stdin successor has no tracker item

`kata list --status open` carries no stdin/envelope kata; the successor exists
only as a line in 0025's `evidence/critique/Charted.md`. A5's If-wrong ("the
successor never ships or withholds nothing, and `sh -s` stays an unowned
admitted shell") is therefore closer to live than the record assumes.

## Q5 — collateral finding: the promise has no delivery surface yet

C1's `promise:` line says "the lint's user-facing description states the
out-of-scope forms in the words above". No such description exists in code:
the only user-facing text is the refusal detail (`load.go:1127`), which fires
on refusal, never on admission. 0027 Phase 3 plans to create it. So C1's
words ARE the shipped reviewer-facing text — which is what makes their
precision load-bearing rather than editorial.
