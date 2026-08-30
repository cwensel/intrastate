# Recommendation 0025: Command-Invoking Accessor Bindings

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-28
- **Status**: Draft
- **Type**: Feature
- **Profile**: foundational — Seam Lineage ≥2 floors it; locks the
  declared-command carrier and authority bound (C1–C5).
- **Priority**: High
- **Related Issues**: intrastate#v0hb (tracker);
  intrastate#p63c (reader cardinality per role, owned by
  RDR 0016)
- **Seam Lineage**: `internal/accessor/binding.go::WriteBinding`
  — 3rd point-fix; trail: dkcf (read-back seal) + 742a
  (0004 suite oracles). Count ≥2 → Profile floored at
  `foundational`; no accretion disposition.

## Problem Statement

A model author declares accessors so that `intrastate` can determine external
state by delegating a read to an established tool and apply a state change by
delegating a write to one — the product's own thesis. Today the second half is
missing: a model can DECIDE a state change (from a linted, exhaustiveness-proven
table) and can persist the decision into intrastate's own flat-JSON artifact,
but no binding invokes a command, so the change cannot be APPLIED to the
artifact the state actually lives in. The caller discovers this when it takes
the resolved decision and must execute the edit itself — an edit no model
declared and no lint covered, relocating exactly the prose-and-drift problem
the table was built to remove. The read side has the same hole: state only a
command can report (a VCS query, a file probe, a tool's own status verb) cannot
be a declared read.

The system-internal requirement: a binding family for
`internal/accessor/binding.go`'s `ReadBinding`/`GateBinding`/`WriteBinding`
that invokes a DECLARED command under the capability, timeout, read-back and
disposition rules RDR 0004 fixed, with the executed authority bounded and
reviewable statically from the model alone (`intrastate lint` validates the
declaration) — not a shell runner, and not a re-litigation of RDR 0004's
rejected raw shell-out (Alt 2) or global executable allowlist (Alt 3).

## Critical Assumptions

- **A1 [The declared timeout can genuinely bound a command invocation —
  including the child's process *tree*, a non-reading child's stdin pipe, and
  post-kill pipe drain — building on the executor's context deadline plus
  process-group termination and a drain bound in the binding]**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: `internal/accessor/executor.go::invokeRead` wraps every
    invocation in `context.WithTimeout` over the declared bound (the executor,
    not the binding, writes the deadline — `0004:C16`); the spike binds (i) a
    command that sleeps past deadline, (ii) a wrapper whose *grandchild*
    inherits the stdout pipe and sleeps, and (iii) a child that never reads
    stdin, asserting `timeout` refusals, no hang, and no surviving process
    group (premortem P-4).
  - **If wrong**: the CLI hangs or a runaway child outlives the `timeout`
    refusal while still mutating the artifact.
- **A2 [`os/exec` passes the argv vector verbatim to the child: no shell, no
  word-splitting, no glob expansion, on the supported platforms]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: Go `os/exec` package docs and `Cmd.Args` semantics; anchor
    the exact doc sentence at Resolve.
  - **If wrong**: the authority bound of C2 is fiction — a declared argument
    containing shell metacharacters executes something the model never
    declared.
- **A3 [A useful class of established tools binds directly — path-positional
  argv plus either envelope-conformant output, a Resolve-adopted single-value
  raw read mode, or A8's declared exit mapping; everything else is served by
  a thin wrapper conforming to the same C3 envelope, and that wrapper class
  is small and mechanical, not bespoke per integration]**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: bind a real VCS query (e.g. `git config --file {artifact}
    --get <key>`) as a declared read and drive it through `flow` end to end;
    record per tool class (VCS query, file probe, status verb) whether the
    binding was direct, raw-mode, exit-mapped, or wrapped, and what the
    wrapper had to do (premortem P-1/P-2: stdin-as-data corruption and
    envelope-size limits are spike scenarios).
  - **If wrong**: every integration needs a bespoke wrapper script — the
    delegate-to-established-tools thesis fails and the choice tips toward the
    adapter-registry alternative.
- **A4 [Read-back for a command-backed write completes through the role's
  declared reader, which may itself be command-backed, preserving the
  applied-but-unverified sense when it cannot]**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: `internal/accessor/model.go::readerFor` selects the
    read-back reader by artifact role (the model author writes that state by
    declaring `role` on a read entry), and the re-read can fail independently
    of the write (`0004:C13`); `0016:C4` is the normative home for reader
    selection (fail-closed unique `readerFor`) — 0025 conforms as consumer;
    pin the `0016:C4` element id at Resolve.
  - **If wrong**: command writes are unverifiable by construction and every
    such write refuses `read_back_incomplete`.
- **A5 [Kubernetes `ExecAction` documents its probe command as an argv array
  executed without a shell — the demoted prior-art instance of the class
  claim]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: quote the `ExecAction.Command` doc comment from
    `k8s.io/api/core/v1` (not surfaced at Propose; attempts logged in
    evidence/research/prior-art.md).
  - **If wrong**: nothing structural — the class claim already rests on the
    quoted Docker exec-form and Terraform external-program citations.
- **A6 [RDR 0004's closed validation-code set closes 0004's own defect
  vocabulary, not the seam: a successor RDR may add codes for surfaces 0004
  did not define]**
  - **Status**: Pending
  - **Method**: Peer RDR
  - **Evidence**: cite at Resolve the 0004 element that fixes the
    validation-code set's closure scope; `internal/accessor/model.go`
    (`ValidationCodes`, "closed eight-member set") is the code-side reading.
  - **If wrong**: C5's defects must ship as a separate lint family instead of
    extending `accessor.ValidationCode` — carrier shape unchanged, surfacing
    relocated.
- **A7 [A child process's inherited environment does not silently widen the
  authority bound: an explicit env policy (inherit-as-is vs scrubbed) can be
  fixed at Resolve without reshaping the carrier]**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: enumerate what the bound tools actually need from the
    environment (`PATH`, `HOME`, VCS config vars) and fix the policy plus its
    failure surface.
  - **If wrong**: ambient process state reaches the command invisibly — the
    same ambient-authority class `0004:C3` exists to exclude, one seam over.
- **A8 [Established gate/read tools that answer in exit codes (`test -f`,
  `git diff --quiet` conventions) can be admitted via a *declared* per-entry
  exit-code mapping — deny / established-absent as model-visible
  declarations — without laundering execution failure into a verdict]**
  - **Status**: Pending
  - **Method**: Spike
  - **Evidence**: premortem P-5/P-12 show the default non-zero →
    `execution_failure` mapping makes command-gate deny unreachable for
    exit-speaking tools; the spike designs and exercises the mapping shape
    (or concludes such tools always go through wrappers) at Resolve.
  - **If wrong**: exit-speaking gates always need wrappers — narrows A3's
    direct-binding class but changes no contract.
- **A9 [No consumer besides `internal/table/load.go::accessorTable` and the
  registry's constructors requires `path` presence on an accessor entry — a
  locator-less command entry breaks no other partition]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: sweep every reader of `table.Accessor.Path` and every
    validation arm partitioning on locator presence (premortem P-7/P-11: a
    pre-dispatch site that resolves or validates `path` for all accessors
    would reject or bypass command entries).
  - **If wrong**: each such site gains the same carrier discrimination the
    registry gets — mechanical, but it must be enumerated before build.

## Proposed Solution

### Approach

The undecided contract at this seam: **what carries a declared accessor's
real external behaviour**. Today `flowbind` *simulates* command behaviour
through magic `path` suffixes (`internal/cli/flowbind/flowbind.go::verdictFor`,
`::unreachable`) because no declared-command carrier or authority bound
exists — the gap both accreted point-fixes (Seam Lineage) patched around.

Answer: an accessor entry may declare `command`, a **fixed argv vector**
carried on the same TOML entry as its role, keys, and timeout. The vector is
validated at load, executed with **no shell**, and admits exactly one form of
variation: a **closed, tool-defined placeholder vocabulary** (v1: `{artifact}`)
substituted **whole-element only** with the caller-bound artifact path for the
entry's declared role. All other per-invocation data crosses on **stdin as a
JSON object of strings**; read results return on **stdout as a flat JSON
object of strings** — the same shape intrastate's own artifact and
Terraform's external-program protocol already use — and gate results as a
verdict envelope (C3). Command-backed bindings
implement the existing `ReadBinding`/`GateBinding`/`WriteBinding` seam and run
under RDR 0004's executor unchanged: declared timeout, refusal classes,
post-write read-back through the role's reader.

The executed authority is therefore reviewable statically: the reviewer reads
the argv in the model; lint proves the vector is well-formed, the placeholders
are known and whole-element, the entry carries exactly one of
`path`/`command`, and `argv0` is not an inline-code interpreter form (C5).
Nothing a caller supplies at invocation time can change *which words* the
command runs — only which artifact path fills the one declared hole and what
typed data arrives on stdin.

Integration honesty (premortem P-1/P-8/P-9): a tool binds *directly* when it
is path-positional and its semantics fit the envelope or the Resolve-decided
raw/exit-mapping modes (A3, A8); otherwise the integration form is a **thin
declared wrapper that conforms to the same C3 envelope** — its invocation
argv is still in the model and linted, and the envelope bounds what it must
do, but its body is script the model does not carry (accepted, recorded as a
Consequence).

### Technical Design

Components:

- **Carrier** — `internal/table`: `sourceAcc`/`Accessor` gain `command`
  (`[]string`); `internal/table/load.go::accessorTable` today rejects a
  missing `path` ⇒ its rule relaxes to *exactly one of* `path`/`command`
  (C1), with the C5 defects joining the existing load-time rejections.
- **Binding family** — a sibling package to `flowbind` implementing the three
  binding interfaces over `os/exec`. The executor already owns timeout
  enforcement: `internal/accessor/executor.go::invokeRead` wraps invocation in
  `context.WithTimeout` (writer: the executor, from the declared `timeout`
  metadata) ⇒ the binding builds on `exec.CommandContext` and carries no
  timer of its own (A1).
- **Selection** — `internal/cli/flowbind/registry.go::Registry` constructs
  `Reader`/`Writer`/`Gate` from `acc.Path` per capability table ⇒ the carrier
  discriminator slots exactly there: a `command`-carrying entry constructs the
  command binding, a `path`-carrying entry constructs today's file binding.
  Sibling-path check: searched for an existing path-vs-command discriminator —
  none exists (`os/exec` appears nowhere outside tests; `flowbind` is the sole
  non-test binding implementation).
- **Data flow** — per capability, under C3/C4: a read command receives the
  requested key set on stdin and answers with the flat string map (an omitted
  key is not-carried, mirroring `0004:C8`); a gate command answers a verdict
  envelope (three-valued, refusal never laundered into a verdict — the rule
  `internal/cli/flowbind/flowbind.go::Gate` already keeps); a write command
  receives the planned tags on stdin and is verified only by read-back. Field
  grammars, clear-sentinel carriage, and any raw read mode are fixed at
  Resolve (C3; LBD Wire / byte format).

#### Normative Contracts

**C1** — carrier and mutual exclusion.

```normative
[read.<id> | gate.<id> | write.<id>]
command = ["<argv0>", "<arg>", ...]   # TOML array of strings
```

An accessor entry carries **exactly one** of `path` or `command`; both or
neither is a load-time defect. `command` must be non-empty and contain no
empty element. `role`, `keys`, `timeout`, and `read_back` rules are unchanged
from RDR 0002/0004.

**C2**

```normative
{artifact}   # the closed placeholder vocabulary, v1 complete: replaced whole-element by the caller-bound artifact path for the entry's declared role
```

Authority bound. The executed argv is exactly the declared vector after
whole-element placeholder substitution. The placeholder vocabulary is closed
and defined by intrastate, not the model author; v1 is exactly the one token
above. A placeholder is recognized only as a whole argv element. An element that
contains a `{...}` token without being exactly a known placeholder is a
load-time defect (C5) — never silently-literal text. Execution invokes no
shell and performs no other rewriting of the vector.

Substitution contract (premortem P-3): the substituted value must be an
absolute path — the executor refuses the invocation (`execution_failure`)
when the caller-bound artifact path is relative or begins with `-`, so a
path can never be parsed as a flag by the invoked tool.

**C3**

```normative
stdin:         one JSON object, all values strings   # per-invocation data beyond the artifact path
stdout (read): flat JSON object of strings; an omitted key is not-carried
stdout (gate): a verdict envelope — one of RDR 0004's three verdicts, plus a reason
```

Invocation envelope. Per-invocation data beyond the artifact path crosses
only on stdin; read and gate results return only on stdout, in the shapes
above. A write command's success is never taken from its exit status alone —
verification is read-back (`0004:C12`). Field-level grammars (envelope field
names, clear-sentinel carriage, a possible single-value raw-stdout read mode)
are fixed at Resolve by this RDR.

**C4**

```normative
order:    parse the stdout envelope, THEN classify the exit code
deadline: terminate the child's process group; bound the stdin write and the pipe drain
write:    read_back required; verified only through the role's reader, never by exit status
```

Execution-safety inheritance. Command bindings run under RDR 0004's executor
unchanged: the declared timeout
bounds the invocation, and at the deadline the binding terminates the child's
**process group**, not only the direct child, with a bounded stdin write and a
bounded post-kill pipe-drain wait so a non-reading or slow-draining child can
never hang the CLI (mechanism — process-group signal, `WaitDelay`-style drain
bound — fixed at Resolve by A1's spike; premortem P-4). The deadline mints
`ClassTimeout`; spawn failure, a malformed stdout envelope, and — by default —
a non-zero exit are `ClassExecutionFailure`, whose `Detail` carries a bounded
stderr tail (size fixed at Resolve). Ordering is normative (premortem P-12):
the stdout envelope is parsed **before** exit-code classification, so a
well-formed deny envelope with a non-zero exit is a deny, not a failure. By
default a non-zero exit is never a gate deny and never establishes a key
absent; whether an entry may *declare* an exit-code mapping for those two
semantic answers is A8's Resolve question. A write entry carries
`read_back = true` (mandatory for every write since RDR 0004), its success is
never taken from exit status, and it is
verified through the role's declared reader — the re-read runs under its own
bounded timeout (`0004:C15`), never the write's residue — with the
applied-but-unverified sense (`0004:C14`) preserved when read-back cannot
complete.

**C5**

```normative
command_and_path_conflict     # both or neither of path/command declared
command_empty                 # empty vector or empty argv element
command_unknown_placeholder   # unknown or non-whole-element {…} token
command_shell_interpreter     # argv0 + inline-code flag (sh -c, bash -c, python -c, env chains) without explicit opt-in
```

Static validation. New load-time defect classes, surfaced by the same
validation that rejects malformed accessor entries today and reported by
`intrastate lint`. `command_shell_interpreter` is what keeps "not a shell runner" enforced
rather than aspirational (premortem P-3/P-10): a fixed `["bash", "-c", …]`
vector is still a shell runner, so the known interpreter-with-inline-code
forms are load-time defects; the opt-in shape (if any) is fixed at Resolve.
Final literal spellings and their home (table loader vs
`accessor.ValidationCode`) are confirmed at Resolve (A6); the defect classes
are normative.

#### Load-Bearing Decisions

- **Identity** — unchanged: the accessor identity stays RDR 0004's
  `(flow, name, capability)` triple; whether an entry is path- or
  command-backed does not enter identity.
- **Wire / byte format** — stdin/stdout JSON-object-of-strings envelopes
  (C3); exact field grammar deferred, owner: **this RDR at Resolve**, decided
  by the A3 spike against a real tool.
- **Naming** — the field is `command`; the family is "command-backed
  accessors". Rejected: `exec` (implies a shell surface), `run` (verb
  collision with CLI verbs), `argv` (implementation jargon in an
  author-facing carrier).
- **Selection / predicate** — the registry selects the binding constructor by
  carrier field: `path` → file binding (today's `flowbind`), `command` →
  command binding. C1's exactly-one rule makes the selection total; no
  precedence order exists to get wrong.

#### Illustrative Code

Illustrative — intent only; tests must not assert it literally.

```toml
[read.branch_state]
role    = "repo"
command = ["git", "-C", "{artifact}", "config", "--file", ".flowstate", "--get-regexp", "^flow[.]"]
# stdout is not the C3 envelope: binds directly only via the Resolve-decided raw/exit modes (A3/A8), else a thin wrapper
keys    = ["flow.stage"]
timeout = "5s"

[write.stage]
role      = "repo"
command   = ["flowstate-write", "{artifact}"]   # planned tags arrive on stdin
keys      = ["flow.stage"]
timeout   = "10s"
read_back = true
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Bounded child execution under a context deadline | Go stdlib `os/exec` (`CommandContext`) | Available | A1/A2 verify the bound is real |
| Timeout enforcement and refusal classes at the seam | RDR 0004 executor (`internal/accessor`) | Available | C4 inherits, adds no class |
| Reader-per-role read-back resolution | RDR 0016 (in flight, intrastate#p63c) | Deferred | A4 — pin compatibility at Resolve |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Declared-accessor carrier | `internal/table` (`sourceAcc`, `Accessor`, `load.go::accessorTable`) | `path` is required today | Extend | C1 exactly-one rule |
| Execution safety + validation | `internal/accessor` (executor, `Validate`, `ValidationCode`) | validation-code set closed by 0004 | Extend | C5 (home decided via A6) |
| Binding construction | `internal/cli/flowbind/registry.go::Registry` | constructs path-backed bindings only | Extend | Selection LBD |
| Simulated command behaviour | `flowbind::verdictFor` / `::unreachable` path suffixes | magic-string simulation — the accreted seam | Reuse unchanged | command family removes the pressure to grow it |

### Decision Rationale

Scored QOC matrix — options answer the contract question "what carries a
declared command and what exactly may vary at invocation time":
O1 fixed argv + stdin only; **O2 closed-placeholder argv (chosen)**; O3
templated/shell strings; O4 built-in adapter registry. 5 = best.

| Criterion | O1 | O2 | O3 | O4 |
| --- | --- | --- | --- | --- |
| Static reviewability (lint from model alone) | 5 — argv literal | 5 — argv + closed whole-element holes | 1 — executed text is runtime-computed (PA-1/PA-2) | 5 — argv owned by intrastate |
| Binds established tools without wrappers | 2 — artifact path cannot reach argv; `0004:C3` bars ambient discovery as the workaround | 4 — `{artifact}` covers the positional case | 5 — anything spliceable | 2 — only shipped adapters |
| Prior-art alignment | 4 — Terraform external (PA-5) | 5 — Docker exec form (PA-3) + git difftool closed vocabulary (PA-4) + PA-5 envelope | 2 — awf-cli needed an injection linter after the fact (PA-1) | 3 — plugin-registry precedent, none in the peer family |
| Injection / blast radius | 5 — nothing varies | 4 — one path-valued hole, whole-element | 1 — caller data spliced into text | 5 — no author-declared argv |
| Integration cost per new tool | 2 — a protocol wrapper per tool | 4 — direct for path-positional tools, thin wrapper otherwise | 5 — none | 1 — an intrastate release per tool |
| Reversibility | 4 — can add placeholders later | 4 — vocabulary can grow or shrink by validation | 2 — authors embed arbitrary templates immediately | 2 — adapter API becomes public surface |
| **Total** | **22** | **26** | **16** | **18** |

O2 wins on the deciding rows: it is the only option scoring ≥4 on both static
reviewability and binding established tools — the two clauses the Problem
Statement conjoins. O1 concedes the product thesis (every tool needs a
wrapper) to gain no review advantage over O2; O3 concedes review to gain
expressiveness O2 mostly retains; O4 relocates authorship, not authority, and
prices every integration at a release.

Premortem: hardened — the critic's 14-row P-ledger
(evidence/propose-premortem/critic.md) forced into the draft: the
interpreter-argv0 lint defect and dash-guarded absolute-path substitution
(P-3/P-10 → C2/C5), process-group + stdin/drain bounding (P-4 → C4/A1),
envelope-before-exit ordering and the declared exit-mapping question
(P-5/P-12 → C4/A8), the wrapper-honesty restatement and inverted-rejection
repairs (P-1/P-8/P-9 → Approach, Alt 1), locator-presence sweep (P-7/P-11 →
A9), and the read-back staleness/reader-less rows (P-6 → Risks, A4). The
recommendation survives with those folded.
Ground-sweep: clean (18 anchors) — 11 code anchors and 7 peer elements
(0004:C3/C8/C12–C16) all CONFIRMED by a factored checker; one cosmetic note
(role-routing of read-back reads from C12/code, independent-failure from
C13), no load-bearing miss.
Joint-check: fired → 0016 (home: cli/0016 §Normative Contracts C4);
fired → 0020 (home: cli/0025 §Normative Contracts C1 / cli/0020 §Normative
Contracts C1 — mutual tolerance: disjoint rule families at `accessorTable`,
entry shape here, tag admission there; composes) — all three arms ran on the
written proposal. Arm 1 (modify-anchors, `--repo`-resolved): pair with 0020
on `internal/table/load.go::accessorTable` (uncited — 0025's C1 relaxes the
path-required rule there; 0020's tag-admission touch is a disjoint rule
family in the same loader and composes, each record homing its own C1) and
pair with 0016 on `internal/accessor/model.go::readerFor`
(cited — A4 declares the coupling: 0016:C4 is the normative home for
reader selection; 0025's command-write read-back rides it as consumer,
pinning the element id via A4). Arm 2 (contract
literals): no pair involving 0025. Arm 3 (absence, manual): C1 converts the
load-time path-required refusal into acceptance for command entries; grepped
the locked peers (0001–0011) for reliance on that refusal — the line-level
hits are unrelated senses of "path" (resolver disposition, SCXML masking,
output path), and 0004's Alt-2/Alt-3 rejections bound shell and allowlist
shapes, both of which C2/C5 preserve — no locked peer relies on the refusal.

## Alternatives Considered

### Alternative 1: Fixed argv + stdin-only protocol (no placeholders)

**Description**: The Terraform-pure form (PA-5): the declared vector is fully
literal; every per-invocation datum, including the artifact path, crosses on
stdin as JSON.

**Pros**:

- Strictest possible authority bound — nothing varies, ever.
- Exact match to a quoted prior-art protocol.

**Cons**:

- The caller-bound artifact path cannot reach an established tool's argv, so
  no established tool is directly bindable — each needs a protocol-aware
  wrapper, and `0004:C3` (no ambient discovery) rightly bars "the tool
  defaults to cwd" as the workaround.
- Pushes authors toward wrapper scripts, which are exactly the undeclared
  prose surface the product removes.

**Reason for rejection**: it forces the wrapper for *every* tool, including
the path-positional class O2 binds directly (and the exit-speaking class A8
may admit); some tools need wrappers under O2 too (P-8), but O1 pays that cost
universally while O2's one
path-valued, whole-element, dash-guarded hole keeps the review property O1
was buying (P-3 mitigations in C2/C5).

### Alternative 2: Built-in adapter registry

**Description**: The model names a shipped adapter (`adapter = "git-config"`)
plus typed params; intrastate owns every argv shape.

**Pros**:

- Tightest authority: no author-declared argv at all.
- Adapters can be individually hardened and tested.

**Cons**:

- Every new tool integration is an intrastate release — the author cannot
  declare delegation to a tool intrastate has not met.
- Adapter params become a second, ad-hoc command language with its own
  validation burden.

**Reason for rejection**: relocates authorship, not authority; contradicts the
product thesis that the *model author* declares the delegation.

### Briefly Rejected

- **Shell-string / open-templated commands**: re-litigates RDR 0004 Alt 2;
  the peer evidence is decisive — awf-cli carries `command:` shell strings
  with `{{.inputs...}}` interpolation and needed a security-validator plugin
  to warn about "unquoted variable references in commands" after the fact
  (PA-1), and ms-conductor's full Jinja2 templating makes the executed text a
  runtime function (PA-2).
- **Global executable allowlist**: RDR 0004 Alt 3, rejected there —
  "constrains mechanism but not semantic power"; nothing here changes that.
- **Host-language callbacks**: RDR 0004 Alt 5, rejected there — authority
  hidden in code, invisible to lint.

## Context

### Background

Filed as kata intrastate#v0hb after RDR 0004 (Accessor Execution Safety Model)
shipped the capability-bounded execution seam — read, gate, write, each
declared, artifact-role-scoped, timeout-bounded, and (for write) verified by
read-back. No RDR 0001–0024 adjudicates a command-invoking binding (RDR 0014
uses `exec.Command` only in illustrative snippets); RDR 0004's deviations
D1–D18 fence nothing out here, and D17's "write command" language anticipates
process semantics.

The design fork: whether a declared external command is carried as a fixed
argv vector on the accessor's capability-table entry — no shell
interpretation, no interpolation of caller/artifact data — versus a
parameterized/templated command shape that admits bounded substitution. The
lint-validatable-from-the-model-alone authority bound is the load-bearing
fork; carrier placement (generalizing the per-accessor `path` locator vs a
sibling field vs a new accessor shape) is its dependent clause, decided by the
same answer.

### Technical Environment

Go CLI (`intrastate`). The seam: `internal/accessor/binding.go`
(`ReadBinding`/`GateBinding`/`WriteBinding`, `CapRead`/`CapWrite`/`CapGate`);
sole non-test implementation `internal/cli/flowbind` (flat-JSON artifact,
self-described as "deliberately the thinnest thing that satisfies the seam").
Governing peers: RDR 0002 owns the TOML model carrier and the per-accessor
`path` locator; RDR 0004 owns the capability/timeout/read-back/disposition
classes (0004:C12–C14; `ClassTimeout`, `ClassExecutionFailure`) and the
applied-but-unverified sense; RDR 0016 owns read-back comparison-set semantics
and reader-per-role resolution (kata intrastate#p63c). Deviations ledger:
`docs/rdr/0004-accessor-execution-safety-model/artifacts/deviations.md`.

## Research Findings

### Investigation

Prior art was read before enumeration (full record with quotes, queries, and
rejected branches: `evidence/research/prior-art.md`). The peer
state-machine/workflow family splits into shell-string templating (awf-cli,
ms-conductor — PA-1/PA-2) and host-language callbacks (xstate, temporal,
restate — 0004 Alt 5's shape); none carries a statically lint-validatable
command bound, so the bounded forms come from outside the family: Docker's
exec form ("doesn't automatically invoke a command shell … variable
substitution doesn't happen" — PA-3) ⇒ argv-array carriers get no-shell by
construction; git difftool's `$LOCAL`/`$REMOTE` ("evaluated in shell with the
following variables available…" — PA-4) ⇒ when per-invocation data must reach
a declared command, established practice is a closed tool-defined vocabulary
carrying paths, never content interpolation; Terraform's external program
("a list of strings…", "does not execute the program through a shell", JSON
object of strings on stdin/stdout — PA-5) ⇒ the envelope shape for everything
that is not path-positional. Code paths grounding the seam:
`internal/accessor/executor.go::invokeRead` (executor-owned timeout),
`internal/cli/flowbind/registry.go::Registry` (the constructor-selection
slot), `internal/table/load.go::accessorTable` (the `path`-required rule C1
relaxes).

### Key Discoveries

- **Documented** — awf-cli's model carries shell strings with Go-template
  interpolation and ships a lint plugin warning about injection and missing
  timeouts after the fact (PA-1) — the templated-carrier end state.
- **Documented** — Docker exec form and Terraform external-program protocol
  both fix no-shell argv execution; Terraform additionally fixes the
  JSON-object-of-strings stdin/stdout envelope (PA-3, PA-5).
- **Documented** — git difftool demonstrates the closed placeholder
  vocabulary carrying file paths (PA-4).
- **Verified** — the executor already owns timeout enforcement via
  `context.WithTimeout` (`internal/accessor/executor.go::invokeRead`), and no
  path-vs-command discriminator exists anywhere in the tree (`os/exec` absent
  outside tests).
- **Assumed** — established tools can meet the stdout envelope with at most a
  thin wrapper (A3); child env inheritance can be policy-fixed without
  reshaping the carrier (A7).

## Trade-offs

### Consequences

- Positive: a decided transition can finally be APPLIED to the artifact the
  state lives in, by a declared, linted, timeout-bounded command — closing
  the product-thesis gap with authority visible in the model.
- Positive: `flowbind`'s magic-path simulation stops being the only way to
  express command-like behaviour, ending the accretion at this seam.
- Negative: a new process boundary enters the trust surface — child env, PATH
  resolution of `argv0`, and process-listing visibility of argv become
  reviewable concerns (A7, Failure Modes).
- Negative: tools that are not path-positional and not envelope-speaking need
  thin declared wrappers (bounded by A3's spike verdict).

### Risks and Mitigations

- **Risk**: command stdout nondeterminism (ordering, trailing whitespace)
  makes read-back equality flaky.
  **Mitigation**: read-back compares typed tag values from the reader's
  envelope, not raw bytes; comparison-set semantics ride RDR 0016 (A4).
- **Risk**: eventually-consistent external state (a remote-backed query)
  makes a synchronous read-back observe staleness, training operators to
  ignore applied-but-unverified refusals (premortem P-6).
  **Mitigation**: the MVV and A3's spike bind *locally consistent* tools;
  remote-backed accessors are a declared-tool choice the model author makes,
  and the refusal's diagnosis tuple names which reader disagreed — revisit at
  the cross-cutting gate if a retry/settle policy is ever wanted.
- **Risk**: secrets placed in declared argv are visible in process listings
  and in the model.
  **Mitigation**: the carrier is reviewable data by design — document that
  credentials belong in the tool's own config, never in `command`; revisit at
  the cross-cutting gate (secret/credential lifecycle).
- **Risk**: `argv0` resolves via `PATH`, so the same model executes different
  binaries on different hosts.
  **Mitigation**: fold into A7's env policy at Resolve (options: record
  as-declared semantics, or a lint advisory on bare `argv0`).

### Failure Modes

- Visible: load-time C5 defects name the entry and the offending element;
  runtime refusals carry the 0004 diagnosis tuple (accessor, capability,
  role, timeout) with class `timeout` or `execution_failure`.
- Visible: a gate command with a malformed envelope refuses
  `execution_failure` — never a deny the model didn't decide.
- Silent risk: a child that ignores termination at deadline can outlive its
  `timeout` refusal (A1's spike bounds this; diagnosis: refusal class plus a
  lingering process on the artifact).
- Silent risk: a command write whose role has no command-capable reader
  refuses `read_back_incomplete` on every invocation — correct but
  surprising; diagnosis starts at the role's reader declaration. Load-time
  detection of a reader-less write role is reader-cardinality territory —
  RDR 0016's contract, not this one (A4).
- Silent risk: a write tool that treats its stdin as *content* would ingest
  the C3 envelope into the artifact; A3's spike exercises this class, and the
  wrapper contract (envelope in, tool-native invocation out) is the answer
  for tools that read stdin as data (premortem P-2).
- Recovery: fix the model entry, re-run `intrastate lint`, re-invoke;
  bindings hold no state between invocations (`0004:C14` — no retry, no
  undo).

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified (A1–A9)
- [ ] A3 spike verdict on the envelope vs raw single-value read mode

### Minimum Viable Validation

1. Author a model declaring one command-backed read and one command-backed
   write (`read_back = true`) over a single artifact role, delegating to an
   established tool present in CI (e.g. `git config --file {artifact}`).
2. `intrastate lint` accepts it; four mutated copies (path+command conflict,
   empty element, unknown placeholder, `["sh", "-c", …]` interpreter form)
   are each rejected with their C5 defect.
3. Drive a state change end to end: the resolved decision invokes the
   declared write command; the tool — not intrastate — applies the edit to
   the artifact.
4. Read-back runs through the declared command reader and verifies the
   written value; the result reports the written tags.
5. A command that sleeps past its declared timeout refuses `timeout` and
   leaves no orphan process.

End state: a linted model applies and verifies a real state change through
declared commands only, with the executed argv readable in the model.

### Phase 1: Carrier

Extend `internal/table` (`sourceAcc`, `Accessor`, `accessorTable`) with
`command` and the C1/C5 load rules.

### Phase 2: Binding family

Implement command-backed `ReadBinding`/`GateBinding`/`WriteBinding` over
`exec.CommandContext` honoring C2–C4 (a sibling package to `flowbind`).

### Phase 3: Selection

Discriminate the constructor by carrier field at
`internal/cli/flowbind/registry.go::Registry`.

### Phase 4: Surface and proof

Wire lint reporting for the C5 defects and land the MVV scenario plus the
timeout and read-back failure scenarios as tests.

## Validation

### Testing Strategy

Coverage goal: every C1–C5 clause has a test that fails when its rule is
dropped; the MVV scenario is the integration proof.

1. **Scenario**: load-time validation (C1/C5) — a valid command entry plus
   the four MVV mutants (path+command conflict, empty element, unknown or
   partial `{…}` token, `["sh", "-c", …]`).
   **Expected**: the valid entry loads; each mutant is rejected with its own
   named C5 code — asserted by code, never by "validation returned non-empty".
2. **Scenario**: substitution guard (C2) — `{artifact}` bound to an absolute
   path, a relative path, and a `-`-prefixed path; an element embedding the
   token mid-string.
   **Expected**: the absolute path is substituted whole-element and the argv
   the child observes equals the declared vector otherwise byte-for-byte; the
   relative and `-` cases refuse `execution_failure` before spawn; the
   mid-string case is a C5 defect at load.
3. **Scenario**: envelope-before-exit ordering (C4) — a gate command that
   writes a well-formed deny envelope and exits non-zero; one that writes a
   malformed envelope and exits zero; one that exits non-zero with no
   envelope.
   **Expected**: deny, `execution_failure`, `execution_failure` respectively,
   each carrying the 0004 diagnosis tuple and a bounded stderr tail.
4. **Scenario**: deadline (C4, A1) — a child sleeping past `timeout`, a
   wrapper whose grandchild holds the stdout pipe, and a child that never
   reads stdin.
   **Expected**: `timeout` refusal within the bound in all three, no hang,
   and no process from the child's group surviving the refusal.
5. **Scenario**: write read-back (C4) — a command write whose tool applies
   the planned tags; one whose tool exits zero without applying them; one
   whose role reader cannot read a compared key.
   **Expected**: success reporting the written tags; `read_back_mismatch`;
   `read_back_incomplete` (applied-but-unverified) — the write's exit status
   decides none of them.
6. **Scenario**: the MVV end to end against an established tool present in
   CI.
   **Expected**: the declared write command, not intrastate, edits the
   artifact, and the declared command reader verifies it.

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

- Prior-art record with quotes and search budget:
  `docs/rdr/0025-command-invoking-accessor-bindings/evidence/research/prior-art.md`
- Docker Dockerfile reference (exec form vs shell form); git-difftool
  documentation (`difftool.<tool>.cmd`); Terraform `external` data source
  documentation (hashicorp/terraform-provider-external).
- Peer checkouts read: awf-cli (`workflow-syntax.md`, security-validator
  plugin README), ms-conductor (`yaml-schema.md` §Template Syntax) — via the
  StateMachineRes corpus / `state-machines` sibling checkout set.
- Source reviewed: `internal/accessor/binding.go`, `internal/accessor/model.go`,
  `internal/accessor/executor.go`, `internal/cli/flowbind/flowbind.go`,
  `internal/cli/flowbind/registry.go`, `internal/table/source.go`,
  `internal/table/load.go`.
- Related: RDR 0002 (TOML carrier), RDR 0004 (execution safety), RDR 0016
  (reader cardinality, intrastate#p63c); kata intrastate#v0hb.
