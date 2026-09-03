# Recommendation 0028: Write a planned value into a line of a text file without a wrapper script

> Revise during planning; lock at implementation. After lock, content is never
> amended; structure may be migrated to the current template by tooling.
> If wrong, abandon code and iterate RDR.

<!-- Section classes: **Required** (never omit). **Conditional**
(delete the whole section if N/A — do NOT leave it blank or
N/A-bulleted). -->

## Metadata

- **Date**: 2026-08-31
- **Status**: Draft
- **Type**: Feature
- **Profile**: large — C1, the `edit` write carrier: one contract whose six clauses (C1.1–C1.6) lock the `edit` block's declaration grammar and its closed placeholder vocabulary; user-facing yes; locks format
- **Priority**: High
- **Related Issues**: intrastate#zyh0 (feature tracker; stays open until Stage 8); intrastate#c3xz (closed — `set-state --plan`; a prerequisite of this RDR's acceptance scenario, not a joint decision); rdr#yjye (open, blocked by intrastate#zyh0) and rdr#qmkd (closed) — the consumer whose state lives in markdown records; intrastate#v0hb (closed, RDR 0025's tracker)
- **Predecessors**: 0025-command-invoking-accessor-bindings, 0004-accessor-execution-safety-model
- **Overrides**: 0025:C1's exactly-one carrier rule (`path` | `command`), to be extended to admit `edit`; 0025's recorded consequence that a planned value reaches an established tool's text artifact only through a thin declared wrapper whose body the model does not carry; 0025:C2's closed argv placeholder vocabulary (`{artifact}`, "v1 complete"), extended by the `{tag.<key>}` family so a command reader over a shared artifact can be row-addressed (C1.6)
- **Seam Lineage**: `internal/table/load.go` carrier admission (`command_and_path_conflict`) / `internal/accessor/binding.go::WriteBinding` (`area:internal-table`) — no prior accretion. Count provenance: no `kata-scope-review §seam-accretion` emission exists on intrastate#zyh0 (routed by `rdr-seed-triage`); taken at seed from `kata list --status closed --label area:internal-table` on 2026-08-31 — the one closed 0025 fix in this file, intrastate#b84g, is C1.5 placeholder-category validation, a different symbol.

## Problem Statement

A model author whose state lives in text artifacts — the RDR flow's `- **Status**:`
metadata bullet and its README index row are the motivating case (rdr#qmkd, rdr#yjye)
— wants a linted table to decide a state change and `intrastate` to apply it to the
record the state lives in, with zero caller edits and zero caller transcription.
Both words are narrow and the record means them literally: no edit the model does
not declare and lint cannot see, and no planned VALUE typed by the caller. They do
NOT mean zero consumer work — the consumer still binds artifact paths, passes the
record identity as `--tag nnnn=NNNN` (an identity, never a value), and adds the
row-addressed projector verb C1.6 reads with. That consumer-side surface is named in
Capability Dependencies and is the honest scope of "delivered"; what the caller stops
doing is transcribing state and maintaining a wrapper script per field. Their
standing rule is that intrastate is the *only* thing that changes state in a source;
their reader binary never writes. Today they discover the gap the moment a resolved
row needs to land in a markdown field: RDR 0025 gave them the flat-JSON `path` writer
and the `command` writer, and neither can put a planned owned value into a text file
without a wrapper — 0025 says so itself ("a fixed argv cannot carry the planned value
and the tools do not read stdin … a thin declared wrapper … its body is script the
model does not carry"). So the choice is a wrapper script per field, or the caller
running `sed` from an emitted string (the transitional `edit` emits in the consumer's
`rdr-write.toml`) — both relocate exactly the drift the table removed: "an edit no
model declared and no lint covered" (intrastate#v0hb's residue).

System-internally, this asks whether text-artifact writes get a **first-class
declared, in-process write carrier `edit`** — the model carries the edit as data
(an `anchor` regex, a `replace` template, a `mode`), intrastate performs it with no
shell and no subprocess, and the role's reader verifies it by read-back exactly as
today (0004:C12 unchanged) — or whether they stay confined to 0025's accepted
wrapper consequence or 0025's rejected Alt 2 (a shipped adapter registry). If `edit`,
the load-bearing safety clause is the **substitution-admission policy**: what may be
interpolated where (`{<owned key>}` planned values, `{tag.<key>}` declared tags; into
`replace` versus into `anchor`), and why interpolation into *data* is a different
hazard class from the interpolation-into-argv that 0025 rejected. Also to be decided
once: anchor cardinality, `<clear>` semantics, atomicity, the `--allow-commands`
gating of a read-back that runs a command for a write that ran none, and what
`intrastate lint` proves statically. Three prior RDRs (0023/0024/0025) each landed
what their kata asked and the composed capability still did not exist because it was
never a ticket; this RDR is that ticket. Nothing above is decided here.

## Critical Assumptions

- **A1 [The invocation's bound context tags can reach `WriteBinding.Apply` for
  `{tag.<key>}` anchor substitution — today no channel carries them
  (`internal/accessor/binding.go::WriteBinding` takes `(ctx, art, planned)` and
  `internal/accessor/model.go::Artifact` is `{Role, Path}`), so the carriage is a
  seam extension this RDR owns — leaning a context map on `Artifact`, because a
  command READER needs the same tags for C1.6 and `Read` takes the same `art` — and
  Resolve picks the form]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/accessor/binding.go::WriteBinding` declares
    `Apply(ctx, art Artifact, planned []resolve.Tag) error` and
    `internal/accessor/binding.go::ReadBinding` declares
    `Read(ctx, art Artifact, requested []string)` — the SAME `art` value reaches
    both carriers, so one context map on `Artifact` serves the writer's anchor
    tags and C1.6's reader. `internal/accessor/model.go::Artifact` is
    `struct { Role string; Path string }` — no tag or context field exists
    today, confirming the seam extension is unclaimed. Form picked: the context
    map rides `Artifact`, and it is INVOCATION-WIDE, not per-role — the tags come
    from `--tag` on the invocation, so every role's handle carries the same map.
    The carriage is two hops, not one: `internal/cli/flow_exec.go::flowRequest`
    already parses the tags into a `observed []resolveTag` field, but
    `artifactMap()` builds `accessor.Artifact{Role, Path}` from `r.artifacts`
    alone and never reads `observed`, so this RDR also changes `artifactMap()` to
    carry it. That method is the sole non-test `accessor.Artifact` construction
    site and is shared by `resolve`, `set-state` and `read-state`, which is why
    the MVV binds `--tag` on both pipeline halves.
  - **If wrong**: a shared artifact (the README index) cannot be row-addressed
    by writer or reader, and the second half of the MVV is unimplementable
    without a per-record role.
- **A2 [`internal/accessor/executor.go::Executor.Write` resolves the role's
  reader (`readerFor`) BEFORE it calls `Apply`, and the gate state that makes a
  command reader refuse (`cmdbind.Config.AllowCommands`) is reachable at that
  point, so C1.3's pre-mutation refusal of a gate-off read-back is a check at an
  existing site, not a new ordering; failing that, the verb's plan phase in
  `internal/cli/flow_state.go` before any carrier runs is the fallback site]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/accessor/executor.go::Executor.Write` resolves
    `reader, hasReader := e.Registry.readerFor(def.Accessor.Role)` and only
    afterwards calls `binding.Apply(applyCtx, art, slices.Clone(planned))` —
    the order C1.3 requires, within one function. The gate state is baked into
    the reader binding at construction
    (`internal/cli/flowbind/registry.go`: `cfg := cmdbind.Config{BaseDir: baseDir,
    AllowCommands: allowCommands}`), so it is inspectable at the resolution
    point. Qualifier: today's gate check lives inside
    `internal/cli/cmdbind::spawn` and runs only during execution, so C1.3's
    pre-mutation refusal is a NEW check at an EXISTING site — not a
    re-ordering, and not already present.
  - **If wrong**: a caller who forgets `--allow-commands` gets the post-mutation
    `read_back_incomplete` class (applied, unverified) for a write that ran no
    command — C1.3's pre-mutation clause must then be dropped, not weakened.
- **A3 [Go's RE2 `regexp` is expressive enough for the consumer's two anchors —
  the `- **Status**:` bullet with an optional bracketed qualifier as a capture
  group, and the README `| [NNNN](` row with the status cell as a group — with
  no lookaround; AND each anchor selects exactly one line in every existing
  consumer record and README, template comments and quoted examples included
  (the premortem's decoy-line case)]**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a3-anchor-regex.{go,md}`, `a3-run.out` — RE2
    suffices with ZERO lookaround (confirmed via `regexp/syntax.Parse`). The two
    normative fixtures, run over 33 record files and the 26-row index README:
    Status anchor `^- \*\*Status\*\*: (.+)$`; per-record README row anchor
    `^\| \[<id>\]\([^)]+\)[^|]*\| [^|]*\| ([^|]+) \|` with `<id>` regexp-quoted
    from `{tag.nnnn}`. Every Status-bearing record selects exactly one line; no
    decoy (template comment, quoted example) scores ≥2. The qualifier capture
    round-trips: `- **Status**: Draft [joint decision → JDR 0003 §D1]` →
    `- **Status**: Final [joint decision → JDR 0003 §D1]` under `${2}`.
    A WRAPPED qualifier (engine fixture
    `rdr/tools/rdr/testdata/status/records/0021-cache-warmup-order.md`, whose
    qualifier continues onto a second line) is SAFE, not truncating: `${2}`
    captures the trailing remainder OF THE MATCHED LINE and re-emits it
    verbatim, the continuation is a different line that C1.3's "preserve every
    other byte" leaves untouched, and the post-edit `diff` is exactly one line.
    The reader reports the qualifier whole
    (`internal/scan/fields.go` joins continuations via `model.ValueContinues`),
    so 0004:C12 compares a planned `Final` against a reported `Final`.
    A record whose VALUE sits on the continuation line selects ZERO lines and
    refuses `edit_anchor_unmatched` — C1.3's designed answer for a stale model.
  - **If wrong**: the anchor grammar needs a second dialect or a literal-anchor
    mode, and C1.2's "RE2" pin moves; or the consumer's records refuse
    `edit_anchor_ambiguous` on a decoy and the author must tighten the anchor,
    never the tool pick an occurrence.
- **A4 [The consumer's command reader reports the owned `status` value with the
  bracketed qualifier separated (as `rdr status --tags` already does:
  `status=Draft`, `status_form=…`), so a planned `Final` written by a
  qualifier-preserving `replace` reads back byte-equal under 0004:C12]**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a4-readback-qualifier.md`, `a4-run.out` — the
    reader separates the qualifier from the owned value. Normative fixture, on a
    record whose Status is `Final [joint decision → …]`:
    `rdr status -json -filter status <artifact>` reports
    `{"facts":[{"name":"status","kind":"enum","value":"Final"}],…}` — byte-equal
    to a planned `Final` (`od -c` confirms no bracket in the bytes). The
    qualifier rides separate facts: `status_form=joint-decision`, and
    `inspect --json --filter metadata`'s status object carries
    `value`/`qualifier`/`form`/`tier`/`raw` as distinct fields. Control (plain
    `Draft`): `status=Draft`, `status_form=none`. Identical output whether the
    reader is invoked by record number or by artifact path, which is the form
    the MVV binds. Qualifier forms found in the corpus: `joint-decision`,
    `revised-from`, `revisit-when`.
  - **If wrong**: every lock of a joint-decision record refuses
    `read_back_mismatch`, and the qualifier must be dropped by the edit
    (rdr-write.toml's own lock row would then be the wrong precedent).
- **A5 [Staging beside the target and renaming over it — with the ORIGINAL
  file's mode copied, unlike `flowbind.go::save`'s fixed 0600 — is acceptable for
  a git-tracked markdown file: git sees a content change only, and no consumer
  depends on the inode or a hard link]**
  - **Status**: Verified
  - **Method**: Spike
  - **Evidence**: `evidence/spikes/a5-stage-and-rename.{go,md}`, `a5-run.out` —
    stage-and-rename with the original mode copied leaves git reporting a
    content change only: `git status --porcelain` gives ` M target.md` and
    `git diff --summary` is EMPTY. Scope correction the spike forces: git tracks
    the executable bit only (100644 vs 100755), not full POSIX mode, so for the
    non-executable markdown this RDR targets `flowbind.go::save`'s fixed
    `os.Chmod(staged, 0o600)` is ALREADY git-silent. Mode-copy earns its keep on
    an executable original — there the fixed-0600 path emits
    ` mode change 100755 => 100644 target.md` — and for non-git consumers of the
    file's permissions. Settled consequences, accepted as stated losses: the
    inode always changes and a hard link silently orphans (retains pre-edit
    content); xattrs are always lost; a symlinked target is correct ONLY if
    resolved first — renaming onto the symlink path replaces the symlink with a
    regular file, which is exactly what C1.3's "symlinks resolved" buys.
  - **If wrong**: the writer must fall back to in-place truncate-and-write and
    give up torn-write safety, or document the link-breaking consequence. The
    spike also settles symlinked targets (C1.3 resolves them) and xattr loss.
- **A6 [No loader, dump, normalize or `checkAccessorBindings` path assumes the
  carrier set is exactly `{path, command}`, other than
  `internal/table/load.go::carrierDefect` and
  `internal/cli/flowbind/registry.go::commandBacked`, which this RDR extends by
  construction]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: Both named exceptions confirmed —
    `internal/table/load.go::carrierDefect` (the `hasPath`/`hasCommand` conflict
    switch) and `internal/cli/flowbind/registry.go::commandBacked`
    (`return len(acc.Command) != 0 || acc.Path == ""`).
    `internal/table/load.go::checkAccessorBindings` is carrier-agnostic (arity
    only). No accessor dump/marshal/normalize round-trip exists anywhere in
    `internal/`, so this assumption's stated failure mode has no live target.
    One further site the audit adds: `internal/table/load.go`'s loader field
    copy from `sourceAcc` (`if a.Path != nil {…}` then `if a.Command != nil {…}`)
    is ADDITIVE per-field copying, not an exclusionary branch — a third carrier
    needs one parallel `if a.Edit != nil {…}` line, not a rewrite.
  - **If wrong**: a dump/normalize round-trip drops or mangles `edit` silently
    and `intrastate lint` proves a table that is not the one executed.
- **A7 [`flow set-state` accepts `--tag` as context (`registerTagFlag`,
  "context only and is never written"), so the record identity an anchor needs
  (`{tag.nnnn}`) is bindable on the invocation without a reader run and without
  the caller typing a planned VALUE]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/flow.go::registerTagFlag` registers
    `cmd.Flags().StringArray("tag", nil, "observed tag, as name=value (repeatable)…")`,
    wired into `flow set-state` in `internal/cli/flow_state.go`.
    `internal/cli/flow_input.go::parseTags` ENFORCES the context-only semantic
    rather than merely documenting it: a `--tag` naming an owned key is refused
    with "owned state is read from the declared read accessors, never supplied
    by the caller". So `{tag.nnnn}` is bindable on the invocation with no reader
    run and no caller-typed planned value. (The RDR's quoted phrase "context
    only and is never written" is a paraphrase — the source says "observed tag";
    the enforced behaviour matches.)
  - **If wrong**: the MVV's README step needs `set-state` to consume the
    resolve envelope's observed tags — a linkage `--plan` deliberately refuses.
- **A8 [0004's `execution_failure` class plus a Detail carrying the rule id and
  the reason token, with the applied sense riding `Applied()`, is enough for a
  caller (an agent loop) to tell a pre-mutation edit refusal from a command
  carrier's unknown-state failure; no new refusal class or typed reason field
  is needed in v1]**
  - **Status**: Verified
  - **Method**: Peer RDR
  - **Evidence**: JDR 0003 §D1 (`state: settled`;
    `cluster: 0004, 0025, 0026, 0027, 0028`). Its resolution is "three rules, no
    new field, no new class", and it names this assumption by id: "`0028:A8` is
    confirmed as a constraint of this entry, and its If-wrong 'owned here' is
    void — a typed field on the carrier is this registry's, not 0028's." The
    three sub-claims map to §D1's rules: sub-reason text rides `Detail`, the one
    slot (rule 1); `Err` is typed and executor-facing, never caller text
    (rule 2); the applied sense is `Applied()`, set by the executor, never prose
    (rule 3). A typed reason enum was considered and rejected there — "a second
    closed set beside the class set with no caller branch to justify it". The
    alignment holds at cli/0026:C1, which uses `execution_failure` post-run with
    `Applied()` true ("ran, output unproven"), where C1.3 here uses it
    pre-mutation with `Applied()` false — the same class distinguished by
    `Applied()`, exactly as rule 3 intends.
  - **If wrong**: unreachable as written — §D1 settles the question and voids
    this record's claim to own it. A typed reason field would be an amendment
    to JDR 0003, carried here by citation.
- **A9 [A command reader can carry `{tag.<key>}` as a whole argv element under
  the same substitution site as `{artifact}` (`internal/cli/cmdbind`), so C1.6 is
  one more vocabulary member, not a second substitution mechanism]**
  - **Status**: Verified
  - **Method**: Source Search
  - **Evidence**: `internal/cli/cmdbind/cmdbind.go::substitute` compares whole
    elements — `if el != ArtifactPlaceholder { out = append(out, el); continue }`
    — testing each argv element for full equality against the literal
    `{artifact}`, never substring-scanning. A `{tag.<key>}` family is one more
    branch in that same loop at that same site: one vocabulary member, not a
    second substitution mechanism.
  - **If wrong**: C1.6 needs its own substitution pass in `cmdbind`, and the
    whole-element rule must be restated there.
- **A10 [The `--allow-commands` gate state can reach the executor as a field on
  `accessor.Registry`, set at `flowbind.go::Registry` — so C1.3's read-back gate
  pre-check needs no new import, no `Binding` interface method, and no change to
  where the gate is enforced for spawning]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: to verify. What is already established: `internal/accessor`
    imports no carrier package and `internal/cli/cmdbind` imports
    `internal/accessor`, so an `accessor`→`cmdbind` back-import is a cycle;
    `flowbind.go::Registry(m, baseDir, allowCommands)` is the single production
    construction site and already receives the gate for 0025:C6;
    `accessor.Registry` is a plain struct the accessor package owns; and
    "command-backed" is answerable at the executor from `def.Accessor.Command`.
    To check: that adding the field breaks no other `accessor.Registry`
    construction site (tests included) and that no second production path builds
    a registry around the gate.
  - **If wrong**: C1.3's pre-check falls back to A2's named alternative — hoist
    it to `internal/cli/flow_state.go` before `exec.Write` — and the refusal is
    then minted by the CLI rather than the accessor, which moves where the
    gate-naming Detail is composed.
- **A11 [The executor's deadline arm (`internal/accessor/executor.go`) is 0004's
  and is left unchanged by this RDR, so C1.3's no-applied-sense guarantee is
  scoped to refusals the `edit` binding itself mints]**
  - **Status**: Pending
  - **Method**: Source Search
  - **Evidence**: to verify that leaving it is right. Established: the arm
    computes `errors.Is(applyCtx.Err(), context.DeadlineExceeded)` AFTER `Apply`
    returns, without consulting `err`, and sets `applied = true` — so a
    pre-write refusal that ran past the deadline surfaces as `ClassTimeout`
    applied-but-unverified. Latent for `path` today, since
    `flowbind.Writer.Apply` ignores its context entirely. To check: whether an
    in-process `edit` makes this reachable in practice, and whether `timeout` is
    still a required field for a carrier that spawns nothing.
  - **If wrong**: this RDR must amend the deadline arm for non-spawning
    carriers, which widens it into `internal/accessor/executor.go` beyond A10's
    one field.

## Proposed Solution

### Approach

Add a third write carrier, `edit`, beside 0025's `path` and `command`. The model
carries the edit as data: per owned key one **line rule** — an `anchor` (an RE2
regular expression that must select exactly one line) and a `replace` template
(literal text, `${N}` backreferences to the anchor's groups, and `{<key>}` for the
planned value) — plus an optional `clear` disposition. intrastate performs the edit
**in process**: read the whole file, resolve every anchor against the pre-edit
content, refuse before writing a byte if any rule is unmatched, ambiguous, or
colliding or any value carries a line terminator, rewrite the selected lines,
preserve every other byte, and stage-and-rename. No shell, no subprocess, no
`--allow-commands` for the write itself. Read-back is unchanged — 0004:C12 through
the role's reader — and `intrastate lint` proves the declaration statically: one
carrier, compiling anchors, well-formed templates, placeholders closed over the
entry's declared keys.

The load-bearing clause is the **substitution-admission policy** (C1.2): a planned
value is admitted only into `replace`, only as literal bytes, never re-scanned;
a declared context tag is admitted only into `anchor`, only regex-quoted, so a
value can choose *which* line but never *what pattern*. Interpolating into line
data has exactly one structural hazard — an embedded line terminator, refused
before mutation — where argv interpolation has a shell's worth; that asymmetry,
stated as C1.2's `value shape`, is why 0025:C5's rejection does not reach here.

`edit` is the `WriteBinding` `internal/accessor/binding.go` already constrains
to "no shell-out, no host callback": it is a second file binding beside
`flowbind.go::Writer`, selected by carrier exactly as `command` is (0025
Selection LBD), and it ends 0025's recorded wrapper consequence for the one
class it covers — a planned value into one line of a text file. One companion
clause (C1.6) lets the same declared context tag row-address a command READER
(`{tag.<key>}` as a whole argv element), because 0004:C12's read-back of a
shared artifact needs a reader that knows which row — without it the README
half of the consumer scenario is unverifiable by construction.

### Technical Design

The edit lives in the same three places the command carrier does: the loader
admits and lints it (`internal/table/load.go::carrierDefect` gains the `edit`
arm and C1.4's categories, `table.Accessor` gains an `Edit` field), the registry
selects it (`internal/cli/flowbind/registry.go::commandBacked` becomes a
three-way carrier discriminator ⇒ the residue still builds a refusing binding),
and a new file binding implements `WriteBinding.Apply` under 0004's executor
unchanged (timeout, refusal classes, read-back). The one seam extension is
A1's: context tags must reach `Apply` for `{tag.<key>}` anchors, which today's
`Artifact{Role, Path}` cannot carry ⇒ this RDR owns that field, and its form is
sharpened at Resolve.

Data flow at apply: `Executor.Write` → non-owned check (unchanged) → resolve
the role's reader (unchanged; C1.3's gate pre-check hangs here, A2) → `Apply`:
read file → split lines (terminators kept) → for each rule: substitute
`{tag.*}` into the anchor (quoted), compile, select exactly one line → for each
rule: expand the parsed template segments (group text, planned value) → refuse
or rewrite → stage + rename → executor read-back through the role's reader.

#### Normative Contracts

> **Proportionality.** One seam — the `edit` write carrier — expressed as six
> clauses (C1.1 grammar, C1.2 admission, C1.3 execution, C1.4 lint, C1.5 clear,
> C1.6 the `{tag.<key>}` reader placeholder). An implementer holds one contract:
> "a linted line rule applied in process under 0004's executor". The clauses are
> sub-labelled so peers can cite them precisely; they are one design, and the
> single fence is what says so — six separately-fenced contracts would claim
> this record spans six seams, which it does not.

**C1** — the `edit` write carrier.

```normative
--- C1.1 — carrier and grammar ---
[write.<id>]
role, keys, timeout               # unchanged (0002, 0004:C15)
[write.<id>.edit.<key>]           # exactly one table per member of `keys`; a table for a key not in `keys` is a defect (C1.4)
anchor  = "<RE2 pattern>"         # a line is SELECTED when the pattern matches anywhere in it (terminator excluded); authors pin `^…$`; must select exactly one line (C1.3)
replace = "<template>"            # the WHOLE replacement line, terminator excluded (C1.2) — never the matched span; a span-style anchor without `^…$` drops the unmatched prefix/suffix by design
clear   = "line"                  # optional; the disposition of a planned `<clear>` (C1.5); absent ⇒ `<clear>` refuses before mutation

carrier: exactly one of `path` / `command` / `edit` per entry — 0025:C1's exactly-one rule with one member appended, not replaced; `edit` is admissible on WRITE entries only (a read or gate entry carrying `edit` is `edit_carrier_conflict`, C1.4). The "neither" arm keeps its 0025:C5 wire string `command_and_path_conflict`, but its PREDICATE necessarily widens to "none of the three" — `internal/table/load.go::carrierDefect`'s `case !hasPath && !hasCommand:` gains `&& !hasEdit`, or an `edit`-only entry would be refused as carrier-less (S1 asserts it LOADS). Its hardcoded message text, which names two carriers, is updated to name three; the wire string is the contract and the message is not (C1.4). `edit_carrier_conflict` follows the "both" arm's discipline — keyed on the `edit` table being PRESENT, not on it being non-empty, matching `case a.Path != nil && a.Command != nil:` so an empty `edit` beside another carrier is a conflict rather than a silently ignored second carrier
in-memory: the registry's residue rule (0025:C1 runtime arm) is unchanged — an entry with no carrier builds a refusing binding, never a file binding

--- C1.2 — substitution admission ---
replace admits: literal text | ${N} and ${name} — the anchor's capture groups | {<key>} — the planned value of THE key this table is named for, and no other key
anchor  admits: literal RE2 | {tag.<key>} — a tag key the model declares with `observed` provenance, bound on the invocation's context (A1, A7); substituted regexp-quoted, so a bound value is a literal-match fragment and can never alter the pattern's structure. Provenance is NARROW because binding is: `internal/cli/flow_input.go::parseTags` refuses an owned key (`flow-tag-owned`) and the recognized key (`flow-tag-reserved`), so a declared key of either kind is structurally unbindable at invocation, and `{tag.<key>}` naming one is `edit_anchor_invalid` at LINT (C1.4) rather than a guaranteed runtime `execution_failure`. The entry's own planned keys are NOT admissible in `anchor`: the value being written never decides where it is written
parse once: both templates are parsed at LOAD into segments (literal | group | placeholder). At apply each segment emits bytes; captured text and substituted values are never re-scanned for `${…}` or `{…}`. This is 0025:C2's whole-element rule restated for line data: runtime data is data, never grammar
value shape: a planned value or bound tag value containing "\n" or "\r" refuses BEFORE mutation (Detail `edit_value_multiline`) — the single structural hazard of interpolation into line data. No word-splitting, option parsing, PATH resolution or shell exists on this path, which is the difference in hazard class from the argv interpolation 0025:C5 rejects
escapes: `$$` emits a literal `$`; `{{` and `}}` emit literal braces in `replace`; a group that did not participate in the match expands to the empty string — the empty-group rule MATCHES Go `regexp.Expand`'s, but the parse is this clause's own and never a call to `Expand`, which re-scans its template at expansion time and would break `parse once:`
vocabulary: the two templates have DIFFERENT rules, because one of them is a regexp. In `replace` the vocabulary is closed — `{<key>}`, `${N}`/`${name}`, `$$`, `{{`/`}}` — and any other `{…}` or `$…` form is `edit_template_invalid`, including a bare `$1`/`$name` without braces (`Expand` accepts those; this clause does not). In `anchor` ONLY the fixed prefix `{tag.` opens a placeholder, scanned to its closing `}`; every other `{`, `}` and `$` is passed to RE2 untouched, so `\d{4}`, `a{2,3}` and `^…$` are ordinary pattern text and `{{`/`}}` is NOT an escape there. `edit_anchor_invalid` fires when a `{tag.…}` placeholder is malformed or names a key not declared `observed`, or when the pattern fails to compile — never for a brace RE2 itself accepts. Extension is by successor amendment, as 0025:C2's is
out of scope for admission: a value that is well-formed line data but wrong for the DOCUMENT (a `|` inside a markdown table cell) is not a hazard this clause can name without document knowledge — the reader's read-back is the check (0004:C12), post-mutation by contract; a per-rule value guard is a successor's (Briefly Rejected)

--- C1.3 — execution semantics ---
input:    the whole file is read as bytes; lines are split on "\n" and a preceding "\r" stays with the terminator (CRLF preserved per line); a missing final terminator is preserved; the file's mode is preserved on write (A5). A target that cannot be read refuses BEFORE mutation with the OS error in the Detail — it does NOT reuse `edit_anchor_unmatched`, and it does NOT take `flowbind.go::load`'s absent-is-empty precedent: that precedent exists so the FIRST write of a flow artifact can create it, and `edit` fences creation out (`select:`), so for this carrier a missing target is a stale model, reported as such rather than laundered into an anchor result
target:   the caller-bound artifact path for the entry's role (0004:C3), symlinks resolved; the model names no path — an `edit` entry has exactly the authority a `path` entry has over the file the CALLER binds. The reuse of `flowbind.go::save`'s discipline stops at stage-and-rename and does NOT drag along `Writer`'s unreachable-locator seal (`unreachable(path)`/`sealedKey`): that seal keys off a declared-path suffix an `edit` entry does not have, and its purpose is to make the applied-but-unverified sense atomic in a key/value artifact — a markdown target has nowhere to hold it. `edit` therefore has no seal affordance, by construction rather than by omission
select:   every rule's anchor is resolved against the PRE-EDIT content and selections are held as pre-edit line INDICES (a deletion never shifts a sibling rule's target). Each must select exactly one line: 0 ⇒ `edit_anchor_unmatched`; ≥2 ⇒ `edit_anchor_ambiguous`; two rules selecting one line ⇒ `edit_anchor_collision`. Never last-match (Ansible `lineinfile`), never first-match, never insert or append (Puppet `append_on_no_match`, Ansible `insertafter`) — creation is fenced out, and an unmatched anchor is a stale model, not a missing line
re-anchor: after the buffer is rewritten in memory, every rule's anchor is run again over the POST-EDIT buffer and must select exactly its own rewritten line (or, for a deleted line, zero lines); otherwise refuse `edit_anchor_unstable` before any write. This is what stops a replacement from de-anchoring itself or poisoning a sibling rule's anchor on the next run (premortem P-4, P-13)
order:    every refusal in this clause and C1.2's `edit_value_multiline` is decided BEFORE any byte is written. A refused edit is NOT APPLIED and surfaces as 0004's `execution_failure`; a rule-scoped refusal carries a Detail naming the rule (`<id>.edit.<key>`) and the reason token (A8), and the two entry-level preconditions — the read-back gate below and C1.6's unbound tag — name the gate and the placeholder instead, having no rule to name; no new refusal class is introduced, and no pre-write refusal THIS BINDING MINTS ever carries 0004:C14's applied-but-unverified sense (applied sense per JDR 0003 §D1). The scoping is exact, not defensive: `internal/accessor/executor.go`'s deadline arm evaluates `errors.Is(applyCtx.Err(), context.DeadlineExceeded)` AFTER `Apply` returns without consulting its error, and sets `applied = true` — so a slow machine can mint `ClassTimeout` applied-but-unverified over an edit that refused before writing a byte. That arm is 0004's, unchanged here and out of this RDR's scope (A11); the binding's own guarantee is the one stated above. A rename that fails after a good staged write is also NOT APPLIED — the target is untouched by construction
write:    all rules of one entry rewrite ONE buffer and land in ONE write: staged beside the resolved target and renamed over it (`internal/cli/flowbind/flowbind.go::save`'s discipline, mode preserved, A5). A post-edit buffer equal to the input is not written at all (no staging, no rename; S20 asserts it, witnessed by an unchanged inode — no MVV step covers the no-op, step 5's files being untouched by a REFUSAL instead). No lock and no compare-before-rename: a concurrent writer is out of scope, as it is for `path`
terminators: only "\n" — optionally preceded by "\r" — terminates a line for selection and rewriting; a bare "\r", NEL or U+2028 is line content (C1.2 still refuses "\r" in a VALUE). Deleting the final line of a file that had no final terminator also removes the preceding line's terminator, so the file's final-terminator state is preserved either way
no subprocess: `edit` spawns nothing; `--allow-commands` (0025:C6) is not consulted by the write itself. `Invocations()` counts `Apply` calls exactly as `flowbind.go::Writer` does (0004:C14)
read-back: unchanged — 0004:C12/0004:C13 through the role's declared reader. When that reader is command-backed and the gate is off, the write MUST refuse BEFORE mutation (Detail naming the gate), because the reader is resolved before `Apply` (A2) — a forgotten flag must not produce `read_back_incomplete` for a write that ran no command. SITE: both halves of that predicate are answered at the executor without a new import or a `Binding` method — command-backed from `def.Accessor.Command` (already on `Definition`), and the gate from a field on `accessor.Registry`, set at `flowbind.go::Registry`, the single production construction site, which already takes `allowCommands` for 0025:C6. `internal/accessor` does NOT import `cmdbind` (`cmdbind` imports `accessor`; the reverse is a cycle) and `AllowCommands` on `cmdbind.Config` stays where it is — the gate is carried to the accessor as state, never read across the seam (A10)

--- C1.4 — load-time categories (`intrastate lint`) ---
edit_carrier_conflict    # `edit` beside `path` or `command`; or `edit` on a read/gate entry
edit_key_mismatch        # a `keys` member with no `edit.<key>` table, or an `edit.<key>` table for a key not in `keys`
edit_anchor_invalid      # anchor fails to compile as RE2 (checked with every `{tag.<key>}` replaced by a quoted probe), names an undeclared tag key, or carries any other `{…}` form
edit_template_invalid    # replace carries an unknown placeholder, another key's placeholder, a group reference the anchor does not define, or a malformed `${…}`/`{…}` token
edit_clear_invalid       # `clear` outside the closed set {"line"}

registration: appended to `table.Categories()` after 0025:C5's six, in the order above; typed `Cat…` constants beside the others; the wire strings are the contract, the identifiers are not; the list's size is not a contract (0025:C5)
precedence: within one entry, fail-fast in the order above, evaluated after 0025:C5's clauses 1–6 (an `edit` entry never reaches clauses 2–6, which are `command`-only); across entries and tables, 0025:C5's rules apply unchanged
what lint proves: the carrier is one, every anchor compiles, every template parses, every placeholder is closed over the entry's keys and the model's declared tags, every `keys` member has one rule. What lint does NOT prove: that an anchor matches exactly one line of a particular file — that is C1.3's apply-time refusal, by design (a model is linted without its artifacts)

--- C1.5 — `<clear>` on a line rule ---
clear = "line" ⇒ a planned `<clear>` deletes the anchored line, terminator included; read-back asserts the key ABSENT through the role's reader (0004:C11 unchanged). An anchor matching ZERO lines on a `<clear>` plan is SUCCESS with no write — "clearing a key the artifact does not hold MUST succeed" (0004:C11); ≥2 matches still refuses `edit_anchor_ambiguous`
clear absent ⇒ a planned `<clear>` refuses BEFORE mutation, Detail `edit_clear_undeclared`. A Status bullet has no meaningful "cleared" line, and the model author says so by omission
one-way: a deleted line cannot be re-established by `edit` (no append, C1.3), so after a `clear` the next non-clear write refuses `edit_anchor_unmatched` — declare `clear` only where the LINE is the key (a per-record bullet), never for a shared row whose other cells carry other state (the README row: deletion would drop the whole row). A blank-the-cell disposition (`clear = { replace = … }`) is deferred: it is a document-shape decision, and v1's closed set is {"line"}
the literal string `<clear>` is never substituted into `replace` (0004:C11: a re-read holding the literal is a mismatch)

--- C1.6 — `{tag.<key>}` as a command placeholder (0025:C2 extended) ---
[read.<id> | gate.<id> | write.<id>]  command = [..., "{tag.<key>}", ...]
admission: `{tag.<key>}` joins 0025:C2's vocabulary as a family, under 0025:C2's rule unchanged — WHOLE-ELEMENT only, replaced by the bound value of a tag key the model DECLARES; `<key>` undeclared is `command_unknown_placeholder` (0025:C5, unchanged wire string); a `{…}` element that is neither `{artifact}` nor a declared `{tag.<key>}` stays `command_unknown_placeholder`
binding: the value comes from the invocation's context (the same channel as C1.2's anchor tags, A1); an unbound tag at invocation refuses `execution_failure` BEFORE spawn, Detail naming the placeholder — a placeholder is never passed through literally (0025:C2)
why here: 0004:C12 read-back re-reads the SAME role; a shared artifact (an index) has one reader for many records, and without an identity in its argv that reader cannot say which row it read. The consumer's own projector stays the reader (no line-oriented READ carrier is introduced — Briefly Rejected)
no other change to 0025:C1–C1.6: stdin envelopes, exit maps, env overlay, the gate, and the shell-interpreter deny-list are untouched
```

#### Pre-Lock Mini-Checks

Cue-fired at Stage 5 (grounding pass). Three of the five fired; the
source-authority census does not (C1.3's `target:` clause is the single
authority — the caller binds the path, the model names none, no fallback or
derived arm exists) and test-discriminability does not (no MVV step or Expected
line passes by absence-of-error: every one asserts end-state bytes or a named
refusal token, and steps 5–7 plus S21's mode-change row are negative controls).

**`fidelity`** — the byte-preservation invariants C1.3 asserts, each with the
test that holds it.

| Operation | Invariant | Held by |
|---|---|---|
| read → split → rewrite → write | every byte outside the selected line(s) is unchanged; `diff` shows exactly one line per edited file | C1.3 `input:`/`write:`; MVV 3; S12 |
| line terminators | CRLF preserved per line; a preceding `\r` stays with the terminator; bare `\r`, NEL, U+2028 are content, not terminators | C1.3 `terminators:`; S19 |
| final terminator | a missing final terminator is preserved; deleting a final line of an unterminated file also removes the preceding terminator, preserving the state either way | C1.3 `input:`/`terminators:`; S19 |
| file mode | the ORIGINAL file's mode is copied to the staged file before rename (unlike `flowbind.go::save`'s fixed 0600) | C1.3 `write:` (A5); S21 |
| symlinked target | the symlink is resolved before staging, so it survives and points at the new content | C1.3 `target:`; S22 |
| no-op edit | a post-edit buffer equal to the input is not written at all — no staging, no rename; witnessed by an unchanged inode, since content cannot separate it from an identical write | C1.3 `write:`; S20 |
| template round-trip | captured text and substituted values emit as literal bytes, never re-scanned for `${…}`/`{…}` | C1.2 `parse once:`; S15 |
| qualifier round-trip | a wrapped bracketed qualifier: the anchored line is rewritten via `${2}`, the continuation line is byte-identical, and the reader reports the qualifier whole | A3; S12 |

Lossy-exemption sites: none. `clear = "line"` (C1.5) is a declared deletion, not
a fidelity exemption — it is one-way by contract, and the clause says so.

**`disposition`** — every input class, and what it mints.

| Input class | Outcome | Error / category | Artifact minted | Silent or loud |
|---|---|---|---|---|
| well-formed edit, anchor selects 1 | applied | — | rewritten file (1 line) | loud (read-back asserts) |
| post-edit buffer equals input | success, no write | — | none — no staging, no rename | silent by design (S20) |
| `<clear>` plan, `clear = "line"`, 1 match | applied | — | line deleted, key reads ABSENT | loud (read-back) |
| `<clear>` plan, `clear = "line"`, 0 matches | SUCCESS, no write | — | none | silent by design (0004:C11) |
| `<clear>` plan, `clear` absent | refused pre-mutation | `execution_failure` / `edit_clear_undeclared` | none — file untouched | loud |
| anchor selects 0 | refused pre-mutation | `execution_failure` / `edit_anchor_unmatched` | none | loud |
| anchor selects ≥2 | refused pre-mutation | `execution_failure` / `edit_anchor_ambiguous` | none | loud |
| two rules select one line | refused pre-mutation | `execution_failure` / `edit_anchor_collision` | none | loud |
| re-anchor fails post-edit | refused pre-write | `execution_failure` / `edit_anchor_unstable` | none | loud |
| planned or bound value carries `\n`/`\r` | refused pre-mutation | `execution_failure` / `edit_value_multiline` | none | loud |
| command-backed reader, gate off | refused pre-mutation | `execution_failure`, Detail names the gate | none | loud (never `read_back_incomplete`) |
| unbound `{tag.<key>}` at invocation | refused pre-spawn | `execution_failure`, Detail names the placeholder | none | loud |
| target cannot be read (ENOENT/EISDIR/EACCES) | refused pre-mutation | `execution_failure`, Detail carries the OS error — the one refusal in C1 with no `edit_*` reason token, because no rule is at fault | none — nothing opened, nothing staged | loud |
| load-time defects | lint finding | the five `edit_*` categories (C1.4) | lint output names entry/key/category | loud (F1) |
| `none`/`stopped:*` resolve row | applies nothing, exit 0 | — | `dispositions` carried | loud (MVV 4) |

Every refusal is decided BEFORE any byte is written (C1.3 `order:`), carries
`Applied()` false, and adds no new class — `execution_failure` with a Detail,
per JDR 0003 §D1. The two silent arms are the two the contract states as
success; there is no third.

**`trace`** — the MVV walked stepwise, with the assertions in force at each step
and a witness value from the normative fixtures (A3's 33-record corpus and
26-row index; A4's reader output).

| Step | Assertions in force | Witness | Verdict |
|---|---|---|---|
| 1 — `intrastate lint --model` passes | C1.4 all five categories; C1.6 admission; C1.1 exactly-one carrier | fixture declares `edit` alone on both write entries + `{tag.nnnn}` on the readme reader → no category fires; entry carrying only `edit` LOADS (0025:C5 "neither" arm does not fire — C1.1) | consistent |
| 2 — resolve \| set-state exits 0 | C1.6 binding (`--tag nnnn=NNNN` bound on the invocation); A1 one context channel; A7 `--tag` is context, never written | `registerTagFlag` registers `--tag` as observed context; `parseTags` refuses an owned key — so `nnnn` is admissible as context and reaches both the anchor and the argv | consistent |
| 3a — record Status `Draft`→`Final` | C1.2 `{status}` planned value; C1.3 select-exactly-one, re-anchor, write | anchor `^- \*\*Status\*\*: (.+)$` selects 1 of 33; post-edit `- **Status**: Final` re-matches its own anchor → no `edit_anchor_unstable` | consistent |
| 3b — joint-decision record keeps qualifier | C1.2 `${N}` groups; C1.3 re-anchor | `- **Status**: Draft [joint decision → JDR 0003 §D1]` → `Final [...]` under `${2}`; re-anchor selects exactly the rewritten line | consistent |
| 3c — wrapped-qualifier record | C1.3 `input:`/`write:` (byte preservation); A3 | `${2}` re-emits the matched line's remainder verbatim; the continuation is a different line C1.3 leaves untouched; `diff` = 1 line | consistent |
| 3d — README row status cell | C1.6 `{tag.nnnn}` argv identity; C1.2 anchor tag quoted; C1.3 select-one | row anchor `^\| \[<id>\]\([^)]+\)[^\|]*\| [^\|]*\| ([^\|]+) \|` with `<id>` regexp-quoted from `{tag.nnnn}`; selects 1 of 26 | consistent |
| 3e — read-back both files | 0004:C12/C13 through the role's declared reader; A2 reader resolved before `Apply`; A4 | reader reports `status=Final` with the qualifier on `status_form`/`status.qualifier` → byte-equal comparison against planned `Final` succeeds | consistent |
| 3f — every other byte unchanged | C1.3 `write:` one buffer, one write; mode preserved | `git diff --stat` = one line per file; mode unchanged (A5) | consistent |
| 4 — `none`/`stopped:*` row | 0005:C1 plan carriage; C1.3 (nothing selected, nothing written) | applies nothing, exits 0 carrying `dispositions` | consistent |
| 5 — no `--allow-commands` | C1.3 `read-back:` gate pre-check (A2); C1.3 `order:` refuse-before-write | reader is command-backed and gate is off → refuses BEFORE mutation naming the gate; both files byte-identical | consistent |
| 6 — duplicated README row | C1.3 `select:` ≥2 | `edit_anchor_ambiguous`, pre-mutation, files untouched | consistent |
| 7 — self-de-anchoring `replace` | C1.3 `re-anchor:` | `edit_anchor_unstable` from the post-edit pass, before any write | consistent |

No CONTRADICTION row. The trace's one crossing point — step 5's gate pre-check
firing before step 3's mutation — is the ordering C1.3 `read-back:` states and
A2 verified at `internal/accessor/executor.go` (`readerFor` resolved before
`binding.Apply`).

#### Load-Bearing Decisions

- **Identity** — unchanged: the accessor identity is 0004's `(flow, name,
  capability)`; the carrier does not enter it (0025 Identity LBD). A line rule's
  identity is `(entry, key)`.
- **Wire / byte format** — a line is the bytes up to and excluding `\n`, with a
  preceding `\r` treated as part of the terminator; the file is bytes, not a
  decoded string, and RE2 matches over UTF-8 bytes. Rejected: normalising line
  endings on write (would touch bytes no rule selected).
- **Naming** — the carrier is `edit`; the rule tables are `edit.<key>`; the
  fields are `anchor` / `replace` / `clear`. Rejected: `line` (Puppet's type
  name, but it names the artifact not the operation), `sed` (names a tool and
  its dialect), `patch` (implies a diff), `regexp`/`match` for the anchor
  (`match` is already a rule block name in the table grammar).
- **Selection / predicate** — `registry.go::commandBacked` becomes a three-way
  carrier discriminator: `command` → command binding; `edit` → edit binding;
  `path` → file binding; residue → refusing binding, exactly as 0025:C1's runtime
  arm. Load-time exactly-one keeps the selection total; no precedence between
  carriers can matter because two can never coexist past the loader. Anchor
  cardinality is the other predicate: exactly one, refuse otherwise (C1.3) —
  searched the repo for an existing line-selection signal to reuse: none
  exists (no non-test `regexp` import under `internal/`).

#### Illustrative Code

Illustrative — the consumer's two rules as they would read in `rdr-write.toml`
once it is a state-machine (field spellings are C1's; the anchors are A3's spike
subject, not normative):

```toml
[write.record]
role = "record"
keys = ["status"]
timeout = "5s"                      # `Timeout` is *string, parsed with time.ParseDuration — never a bare integer
[write.record.edit.status]
anchor  = '^- \*\*Status\*\*: (?:Draft|Final)(?P<q> \[.*\])?$'
replace = '- **Status**: {status}${q}'

[read.readme]                       # C1.6: the projector reads ONE row, told which by the same context tag
role = "readme"
keys = ["status"]
timeout = "5s"
command = ["rdr", "index", "--row-json", "{tag.nnnn}", "{artifact}"]   # verb illustrative; the consumer's to add

[write.readme]
role = "readme"
keys = ["status"]
timeout = "5s"
[write.readme.edit.status]
anchor  = '^(\| \[{tag.nnnn}\]\([^)]*\) \| [^|]* \| )[^|]*( \| [^|]* \|)$'
replace = '${1}{status}${2}'
```

### Capability Dependencies

| Needed Capability | Source | Status | Spec Impact |
| --- | --- | --- | --- |
| Exactly-one carrier admission at load (`carrierDefect`) | Predecessor 0025:C1/0025:C5 | Available | Extended by one member (C1, C1.4) |
| Carrier-selected binding construction (`commandBacked`) | Predecessor 0025 Selection LBD | Available | Gains the `edit` arm |
| Post-write read-back through the role's reader | Predecessor 0004:C12/0004:C13 | Available | Unchanged; C1.3 adds a pre-mutation gate check (A2) |
| `<clear>` as removal (`accessor.IsClear`) | Predecessor 0004:C11 | Available | Reused (C1.5) |
| Atomic staged write (`flowbind.go::save`) | Existing | Available | Discipline reused; mode handling differs (A5) |
| Context tags crossing the write seam | This RDR | Introduced | A1 — the one seam extension |
| `set-state --plan` (the resolve→apply pipe) | Existing (intrastate#c3xz, landed) | Available | MVV prerequisite only |
| Command reader over the consumer record (`rdr status -json -filter status`) | Existing (0025 command reader) | Available | MVV fixture |
| Row-addressed command reader over a shared artifact (`{tag.<key>}` in argv) | This RDR (C1.6) | Introduced | 0025:C2 vocabulary extended; the consumer adds the verb |

### Existing Infrastructure Audit

| Needed Capability | Existing Surface | Known Limit | Decision | Spec Impact |
| --- | --- | --- | --- | --- |
| Carrier admission + categories | `internal/table/load.go::carrierDefect`, `table.Categories()` | Two carriers | Extend | C1, C1.4 |
| Carrier discriminator | `internal/cli/flowbind/registry.go::commandBacked` | Boolean | Extend to three-way | Selection LBD |
| File write binding | `internal/cli/flowbind/flowbind.go::Writer` | JSON store only | Reuse pattern; new binding | C1.3 |
| Atomic save | `internal/cli/flowbind/flowbind.go::save` | Fixed 0600 mode; JSON encoder | Extend or sibling (mode-preserving, bytes) | C1.3, A5 |
| Placeholder vocabulary | `internal/table/load.go::commandPlaceholders` | argv whole-element only | Sibling: a line-template vocabulary (C1.2) | C1.2 |
| Write seam | `internal/accessor/binding.go::WriteBinding`, `model.go::Artifact` | No context tags cross | Extend (A1) | C1.2 |

### Decision Rationale

The user outcome is "a planned value lands in the line the state lives in, with
nothing the caller types and nothing lint cannot see". The scored matrix
(Alternatives Considered) turns on three rows: **drift removal** — only an
in-process carrier (A) and the self-hosted verb (C) put the edit under lint,
and C puts it behind an opaque argv that 0025:C5 can only see as a program
name; **blast radius** — A adds no process boundary and no `--allow-commands`
dependency for the write, where C inherits both for no safety gain; **prior-art
alignment** — A is the settled class shape (Ansible `lineinfile`'s regexp +
backreference line, Puppet `file_line`'s `multiple => false` error on
ambiguity, both with atomic staging), narrowed where the peers are permissive:
no last-match, no append, no create. The adapter form (D) is 0025 Alt 2 with a
document family attached, and the kata fences document knowledge out. The
wrapper (B) and emit-only (E) forms are the two shapes of the drift the record
exists to remove.

**C1.6 is scored separately, because the matrix above is about the write carrier
and C1.6 is not one.** It is load-bearing rather than incidental: without an
identity in the reader's argv, a shared artifact's one reader cannot say which
row it read, so the README half of the user outcome is unverifiable by
construction — and it extends 0025:C2's *closed* argv vocabulary, which is why
this record declares 0025 as an Overrides target. Its one alternative is the
fallback A1 names: a **per-record role**, binding a distinct reader per record
instead of row-addressing one shared reader. Rejected on the same drift-removal
row as the rest — a role per record makes the model grow with the corpus and
puts record identity in the model where a caller-bound tag belongs, reintroducing
the per-field authoring the record exists to remove. The narrower option — ship
`edit` alone and leave the README half unverified — is rejected because it
would make the MVV prove half the stated outcome.

Rejections, one line each: **B** relocates the edit into script no model
carries and no lint covers (0025's own Consequence). **C** is the same template
semantics behind a process boundary and a gate. **D** makes every artifact
shape an intrastate release. **E** is the status quo: the caller's transcription
is the drift. Insert/append on no match: creation is fenced; an unmatched anchor
is a stale model.

The hardened premortem (evidence: `propose-premortem/critic.md`, P-1…P-22)
hardened rather than switched the choice: the README half of the scenario was
unverifiable as briefed (P-15) — folded as C1.6; anchors could poison themselves
(P-4, P-13) — folded as C1.3's re-anchor invariant; escapes, index stability,
symlink and terminator rules (P-6, P-9, P-10, P-18) — folded into C1.2/C1.3;
`clear` is one-way and shape-blind (P-11, P-20, P-21) — folded into C1.5; the
authority question (P-8, P-19) is answered by C1.3's target clause: the model
names no path, the caller binds it, exactly as for `path`. Two keys on one
line (P-2) is refused in v1 by design (Briefly Rejected).

Premortem: hardened (hardened)
Ground-sweep: clean (35 anchors; one cosmetic citation corrected inline — the emit-payload clause is 0010:C4, not 0010:C3)
Joint-check: fired → 0027, 0026, 0016 (home: JDR 0003 §D1 for cli/0026:C1; cli/0025:C5 for cli/0027:C1; cli/0016:C4) — disposed 2026-08-31: 0026 hoisted (constraint — sub-reason rides `Detail`, applied sense is `Applied()`; C1.3 and A8 aligned, cited not restated); 0027 cite-don't-restate (disjoint deltas, 0025:C5 owns the clause map); 0016 cite-don't-restate (reliance on its fail-closed reader). all three arms run on the written proposal (evidence: `joint-check/arms.md`). Arm 1 (modify-anchors, repo-resolved): 0027 ↔ 0028 on `internal/table/load.go::carrierDefect` — 0027:C1 rewrites clause 4's predicate; C1/C1.4 here add the `edit` arm and five tail categories and never reach clauses 2–6; disjoint clauses of one function, no shared answer, 0027 recorded `clear` before this proposal existed. Arm 2 (contract literals): 0026 ↔ 0028 on `execution_failure` — 0026:C1 uses the class for "ran, output unproven", C1.3/A8 here for "not applied" with `applied: false` in Detail (since aligned to JDR 0003 §D1); the class's meaning for writes has no normative home (0004 is Implemented and silent on it); 0016 ↔ 0028 on `read_back_incomplete` — 0016:C4's fail-closed `readerFor` is what A2/C1.3's pre-mutation check sits on (cite, not restate); 0027 ↔ 0028 on `table.Categories()` — 0027 changes no position, C1.4 appends at the tail; 0014/0021 ↔ 0028 on `intrastate lint` — a tool name, not a decision. Arm 3 (absence, manual): this proposal turns two refusals into acceptances — 0025:C5's "neither" arm for an entry carrying only `edit`, and `command_unknown_placeholder` for a declared `{tag.<key>}` argv element (C1.6). Open peers (Draft/Final at depth 1) grepped for `command_and_path_conflict`, `command_unknown_placeholder`, `{artifact}`, "v1 complete": only 0027 hits, and its reliance (clause-3 whitespace exemption vs the placeholder defect) is preserved by C1.6's "stays `command_unknown_placeholder`" line; 0025 (Implemented) relies on the closed vocabulary and is this record's declared Overrides target — never edited, the coupling rides to 7.1. Dispositions above.

## Alternatives Considered

Questions-Options-Criteria matrix. Options: **A** in-process `edit` carrier
(chosen); **B** thin declared wrapper per field via `command` (0025's accepted
consequence); **C** a self-hosted `intrastate edit` verb invoked through the
`command` carrier; **D** a built-in adapter registry / document-family writer
(0025 Alt 2 revived); **E** emit-only — the table emits the edit, the caller
applies it (today's `rdr-write.toml`).

| Criterion | A `edit` carrier | B wrapper | C self-hosted verb | D adapter | E emit-only |
| --- | --- | --- | --- | --- | --- |
| Correctness fit (value lands, verified by read-back) | Yes; in-process, one atomic write | Yes, if every wrapper is right | Yes | Yes for known families | No — caller applies |
| Drift removal (edit declared + linted) | Full: anchor/template linted | None: body is script | Partial: argv linted, template opaque | Full for shipped families | None |
| Prior-art alignment | lineinfile/file_line shape, narrowed | 0025's accepted form | No peer does this | 0025 Alt 2, rejected | Terraform-external-style emit; a stop-gap by its own header |
| Reversibility | High: a carrier, removable | High | Medium: a public verb | Low: a registry contract | High |
| Blast radius | Loader, registry, one binding, one seam field (A1) | None | New verb + process boundary + gate | Adapter API + per-family code | None |
| Cost | One binding + lint clauses | Per-field script forever | Verb + binding + envelope | Highest | Zero now, paid per run |
| `--allow-commands` needed by the write | No | Yes | Yes | No | n/a |

### Alternative 1: Thin declared wrapper per field (B)

**Description**: keep 0025 as is; each text field gets a 5–8 line script
invoked as a `command` write, reading the planned value from the stdin envelope.

**Pros**:

- Nothing changes in intrastate.
- The wrapper can handle any file shape.

**Cons**:

- The edit is script the model does not carry and lint cannot see —
  intrastate#v0hb's residue, verbatim.
- One script per field, per consumer, forever; every one needs
  `--allow-commands`.

**Reason for rejection**: it is the drift the table was built to remove,
relocated one file over.

### Alternative 2: Self-hosted `intrastate edit` verb via the `command` carrier (C)

**Description**: ship an `edit` verb that reads the C1.3 stdin envelope and an
anchor/template from its argv; models declare
`command = ["intrastate", "edit", "--anchor", "…", "{artifact}"]`.

**Pros**:

- Reuses 0025's carrier and gate unchanged; no loader change.

**Cons**:

- The template semantics (C1.2) still live in intrastate — behind a process
  boundary, a PATH resolution of `intrastate` itself, and an argv 0025:C5
  cannot lint beyond its program name.
- The write is gated by `--allow-commands` for no safety gain.

**Reason for rejection**: same admission policy, worse authority surface.

### Alternative 3: Adapter registry / document-family writer (D)

**Description**: `adapter = "markdown-bullet"` with typed params; intrastate
owns the edit shape per family.

**Pros**:

- Tightest authority; per-family tests.

**Cons**:

- Every artifact shape is an intrastate release; document knowledge enters the
  tool, which the kata fences out.

**Reason for rejection**: 0025 Alt 2's rejection holds — authorship relocated,
not authority.

### Briefly Rejected

- **Emit-only (E)**: today's transitional form; the caller's transcription is the
  drift (rdr-write.toml's own header calls it "a stop-gap, not the design").
- **Insert/append on no match** (Ansible `insertafter`, Puppet
  `append_on_no_match`): creation is fenced out; an unmatched anchor is a stale
  model, not a missing line.
- **Last-match / first-match on ambiguity** (Ansible's documented "only the last
  line found will be replaced"): a silent choice among lines is the failure the
  read-back would then confirm as verified.
- **Planned values admissible in `anchor`**: the value being written would decide
  where it is written; also a regex-injection surface for no use case.
- **A line-oriented text READER carrier** (the mirror of `edit`): out of scope —
  the consumer's reader is its own projector by doctrine, and 0025's command
  reader already covers read-back; a successor may seed it.
- **Multi-line region edits, file creation, rename**: fenced by the kata as filed.
- **Two owned keys on one line** (one rule keyed by anchor carrying several
  placeholders): v1 refuses `edit_anchor_collision`; the consumer's artifacts
  keep one owned key per line, and a successor can key rules by anchor without
  breaking `edit.<key>` tables.
- **A per-rule value guard** (`reject = "<RE2>"` on the planned value, e.g. `|`
  in a table cell): document knowledge in a different coat; v1 relies on
  read-back through the reader that parses the document.
- **A typed pre-mutation refusal class** beside 0004's four: A8 decides whether
  `execution_failure` + Detail suffices; a class change is a CLI-contract change.

## Context

### Background

Filed from the consumer side after RDR 0025 landed (intrastate#zyh0, labelled
`dx`). 0025's A3 spike bound `git config`, `test -f` and `git diff --quiet` —
tools that own their own state — and never wrote a field in a text file; the
wrapper form is recorded only under 0025's Consequences and Approach, not as a
locked clause, so the fork is live. `internal/accessor/binding.go` constrains
the seam to Read/Gate/Write with "no shell-out, no host callback": an
in-process `edit` is a `WriteBinding` consistent with that seam, which is why it
is neither 0004 Alt 5 (host callbacks — the edit is visible to lint) nor 0025
Alt 2 (an adapter registry — authorship stays with the model author). The
charted `stdin = "none"|"envelope"` successor is a different 0025:C1 field with a
different purpose — no overlap, but precedent that C1 is extended by successor
RDRs.

`flow set-state --plan <file|->` carries the resolve→apply pipe this record's
scenario runs over. It needs no plan-carriage clause here: 0005:C1 plus the
no-linkage invariant already decide that fork. **This record owns the composed
consumer acceptance scenario**, below as the Minimum Viable Validation.

Fenced out by the kata as filed: multi-line region edits; file creation or
rename (a consumer's number allocation is artifact creation, not a state
change); any knowledge of a document family — the carrier is line-oriented
text, as generic as the flat-JSON `path` writer, but the MVV fixture is a
consumer record (a markdown file with a `- **Status**:` bullet and a README
index row), not a favourable tool like `git config`.

### Technical Environment

Go. `internal/table/load.go`: the write-binding loader and its exactly-one
carrier admission (`command_and_path_conflict`), the 0025 placeholder
categories and their lint; `internal/accessor/binding.go`: the Read/Gate/Write
seam (`WriteBinding`); `internal/cli/flowbind/`: the `path` file binding and the
carrier-selecting `Registry`; `internal/cli/cmdbind/`: the `command` binding;
`internal/cli/flow_state.go`: `flow set-state` (`--write`/`--clear`/`--plan`).
Governing records: 0004:C10 (payload = planned owned-tag map), 0004:C11
(`<clear>`), 0004:C12/0004:C13 (read-back through the role's reader), 0025:C1
(carrier set), 0025:C2 (closed placeholder vocabulary), 0025:C5 (placeholder
admission — the argv-interpolation rejection this RDR distinguishes itself
from), 0025:C6 (`--allow-commands`), 0010:C4 (`emit` crosses the `flow resolve` payload as a JSON object of
strings — data the caller receives, which is why today's edits are applied by
the caller). The
consumer artifact shapes are the RDR engine's `TEMPLATE.md` Metadata bullets and
the per-project README index table; its reader is the `rdr` projector
(read-only by doctrine).

## Research Findings

### Investigation

Prior art was read before enumerating (evidence: `research/prior-art.md`).
Class: configuration-management line editors — Ansible `lineinfile` ("the
pattern to replace if found. Only the last line found will be replaced";
`backrefs` expand the `line` with the regexp's groups; writes via `mkstemp` +
`atomic_move`) and Puppet `file_line` (`match` replaces; `multiple`: "If set to
false, an exception will be raised if more than one line matches";
`append_on_no_match`). Go in-process precedent: kubebuilder
`pkg/plugin/util/util.go::InsertCode` / `::ReplaceInFile` (literal anchor,
refuse on no match, non-atomic `os.WriteFile`, all occurrences). ⚠ no
prior-art coverage in the state-machine corpus for this problem class (two
queries, no hit); the consumer's own `rdr-write.toml` header is the nearest
instance read — it carries the exact `sed` expression today and names this RDR
as its end state. In-repo: `flowbind.go::save` (atomic staging),
`flowbind.go::Writer.Apply` (`IsClear` ⇒ delete), `registry.go::commandBacked`
(carrier discriminator), `load.go::carrierDefect` (clause-ordered categories),
`executor.go::Executor.Write` (reader resolved before `Apply`). No non-test
`regexp` use exists under `internal/` — the anchor selector is new.

### Key Discoveries

- **Documented** — Ansible and Puppet agree on the carrier shape (regex anchor,
  replacement line with backreferences, atomic staging) and disagree on
  ambiguity (last-match vs error); this RDR takes the error.
- **Documented** — `WriteBinding.Apply(ctx, art, planned)` and
  `Artifact{Role, Path}` carry no context tags; an anchor that needs the record
  identity needs a seam extension (A1).
- **Documented** — `Executor.Write` resolves `readerFor(role)` before `Apply`,
  so a pre-mutation gate check has a site (A2). Resolve adds the qualifier: the
  gate state is baked into the reader binding at construction and so is
  inspectable there, but today's check runs inside `cmdbind::spawn` during
  execution — C1.3's clause is a new check at that existing site.
- **Documented** — `flowbind.go::save` fixes mode 0600 and encodes JSON; the
  edit writer needs the staging discipline with mode preservation and raw bytes
  (A5).
- **Verified** — RE2 covers the consumer's anchors with no lookaround, and each
  selects exactly one line across all 33 records and the 26-row index (A3); the
  consumer's reader splits status from its qualifier, so a planned `Final` reads
  back byte-equal (A4).
- **Verified** — a bracketed qualifier that WRAPS onto a continuation line is
  safe under the per-line model, not truncating: the continuation is a different
  line the anchor never selects, and the reader joins continuations into one
  logical value (`internal/scan/fields.go`, `model.ValueContinues`). No guard is
  owed; C1.2's existing "read-back is the check" disposition governs (A3).
- **Documented** — this RDR is STRICTER than its prior art on both halves of
  selection. Ansible `lineinfile` replaces the last match and silently INSERTS
  on zero matches; Puppet `file_line` errors on ambiguity but appends on
  no-match by default. Neither, nor `sed`, nor Ansible `replace`, nor `perl
  -0777`, detects a wrapped-value truncation — no surveyed tool guards it. C1.3's
  exactly-one-or-refuse and C1.2's pre-mutation terminator check are guards ADDED
  beyond the field, not standard ones recovered from it.
- **Documented** — per-entry atomicity (C1.3: all rules land in one write or none)
  is stronger than Ansible (task-by-task, no rollback), Puppet (resource-by-
  resource best-effort) and Terraform (documented as "not transactional"). The
  boundary is one entry rather than a whole run, which is what makes it
  achievable.
- **Documented** — the RDR index README carries two rows (0001, 0006) in an
  unlinked `| NNNN |` form that no `rdr-write.toml` rule emits or locates; both
  date to the README's founding commit. Pre-existing drift in the consumer's
  data, surfaced by A3's corpus sweep — the existing `readme-flip-*` rules
  already cannot flip those two rows. Normalizing them is a consumer data fix,
  not a change to this contract; `edit_anchor_unmatched` on a bare row is C1.3
  reporting a stale model correctly.

## Trade-offs

### Consequences

- Positive: a decided state change lands in the text artifact it lives in, by a
  declared, linted, in-process edit — the composed capability 0023/0024/0025
  each left one ticket short.
- Positive: the consumer's `rdr-write.toml` sheds its `sed` strings and becomes
  a state-machine over a command reader and `edit` writers, as its header
  planned.
- Negative: a second file-binding and a third carrier arm in loader, registry
  and lint — the seam gains one field for context tags (A1).
- Negative: an `edit` write over a role whose reader is command-backed still
  needs `--allow-commands` on the invocation for its read-back (0025:C6
  unchanged); C1.3 makes a forgotten flag a pre-mutation refusal rather than an
  applied-unverified write.
- Negative: anchors are regexes authored by the model author; a drifted
  artifact shape surfaces at apply time as `edit_anchor_unmatched`, not at lint.
- Negative: `clear = "line"` is one-way — `edit` cannot re-create a line — so a
  cleared per-record bullet is re-established only by hand (C1.5).
- Neutral: `edit` writes exactly where the caller's `--artifact` binding
  points, as `path` does; a caller who binds a hook file has bound a hook file.
  The per-invocation `--allow-commands` consent stays what 0025:C6 made it —
  consent to EXECUTE, not to write — and the read-back's command still needs it.

### Risks and Mitigations

- **Risk**: an anchor that matches the intended line AND a look-alike (a second
  `Status` bullet in a quoted block).
  **Mitigation**: C1.3 refuses ambiguity before writing; the author tightens the
  anchor (`^` and the exact bullet form).
- **Risk**: a replacement that no longer matches its own anchor, so the next
  transition cannot find the line.
  **Mitigation**: C1.3's re-anchor invariant refuses `edit_anchor_unstable`
  before any write (MVV step 7). A lint-time "replace satisfies anchor" probe
  (Puppet's `match`-vs-`line` check) would catch it one stage earlier, but
  cannot be decided without the artifact (C1.4's "what lint does NOT prove") —
  a Testing Strategy candidate, not a second apply-time check.
- **Risk**: a value with a terminator splits the line.
  **Mitigation**: C1.2 `edit_value_multiline`, pre-mutation.
- **Risk**: mode/inode change on rename breaks a consumer's hard link or
  permissions.
  **Mitigation**: A5 spike; mode is copied by contract; symlinks resolved (C1.3).
- **Risk**: a value that is valid line data but breaks the document's own
  grammar (a `|` in a table cell) rewrites the row so the reader parses another
  cell.
  **Mitigation**: read-back mismatch, post-mutation by contract (C1.2); the
  artifact is in git; a value guard is a named successor.
- **Risk**: a reader that DEFAULTS an absent key (reports `Draft` for a missing
  bullet) makes a `clear` read back as a value.
  **Mitigation**: that reader violates 0004:C8 (genuine absence must be
  reported as absence); C1.5's read-back would refuse `read_back_mismatch` — the
  honest outcome — and the fix is the reader's.

### Failure Modes

- Visible: `intrastate lint` names the entry, key and category for every C1.4
  defect; nothing loads.
- Visible: an unmatched, ambiguous or colliding anchor, a multi-line value or an
  undeclared `<clear>` refuses `execution_failure` with a Detail naming
  `<id>.edit.<key>` and the reason token — the artifact is untouched.
- Visible: the two entry-level preconditions refuse the same class before
  mutation, with no rule to name: an unbound `{tag.<key>}` (Detail names the
  placeholder, C1.6) and a gate-off command read-back (Detail names the gate,
  C1.3) — the artifact is untouched.
- Silent risk: an anchor that selects the wrong single line (a look-alike
  elsewhere in the file while the real line has drifted) rewrites it; the
  read-back through the role's reader refuses `read_back_mismatch` only if the
  reader reads the real line — which it does for the consumer's projector, and
  which is the reason read-back is the commit-time check.
- Silent risk: an anchor whose exactly-one match is a decoy (a quoted template
  line) while the real line has drifted rewrites the decoy; the read-back
  through the projector, which reads the real line, refuses
  `read_back_mismatch` post-mutation. A3's corpus check is the lint-time
  defence; the exactly-one rule is why the common case (both present) refuses
  before mutation instead.
- Visible: an index row in a shape the anchor does not cover refuses
  `edit_anchor_unmatched` — the consumer's README carries two such rows
  (0001, 0006) in an unlinked `| NNNN |` form that no `rdr-write.toml` rule
  emits or locates. Pre-existing data drift, not a defect of this carrier: the
  existing `readme-flip-*` rules already cannot flip those rows. The remedy is
  to normalize the data, never to loosen the anchor (a looser anchor buys back
  the decoy risk F5 names).
- Recovery: the artifact is in git; fix the anchor or the artifact, `lint`,
  re-invoke.
- Diagnosis: the Detail carries the rule identity and the match count; the
  staged temp file never survives a refusal.

## Implementation Plan

### Prerequisites

- [ ] All Critical Assumptions verified
- [ ] A1's seam form chosen (Artifact context map vs `Write` argument)
- [ ] A3/A4 spike: the consumer's two anchors and the reader's qualifier split
- [ ] Consumer-side, not blocking this RDR's implementation but blocking its
      adoption on the live README: normalize the two unlinked `| NNNN |` rows
      (0001, 0006) to the linked form (F6). Until then the feature ships
      correct and those two records cannot be flipped through it — partial
      adoption, and the MVV passes without covering them because its fixture
      declares all rows linked.

### Minimum Viable Validation

The consumer's acceptance scenario, carried from intrastate#zyh0. Fixture: a
`rdr-write.toml` re-authored as a state-machine over (a) one 0025 command
reader `["rdr","status","-json","-filter","status","{artifact}"]` on role
`record`, (b) an `edit` writer for the record's Status line on role `record`,
(c) a C1.6 command reader over the README row on role `readme` (`{tag.nnnn}` in
its argv) and an `edit` writer for that row anchored by `{tag.nnnn}`; a record
whose Status is `Draft`, one whose Status is `Draft [joint decision → …]` on one
line, and one whose bracketed qualifier WRAPS onto a continuation line while
still reading `Draft`, so the flip and the preservation are exercised together
(A3's case; `0022-cache-metrics-surface.md`, not the already-`Final` `0021`); a README with all three index rows in the linked
`| [NNNN](NNNN-slug.md) |` form that `rdr-write.toml`'s own `readme-add` rule
emits — the unlinked `| NNNN |` form two live rows carry is out of the fixture
and named in Failure Modes.

1. `intrastate lint --model rdr-write.toml` passes; no wrapper script exists
   anywhere in the fixture.
2. `flow --allow-commands resolve --model rdr-write.toml --artifact record=…
   --artifact readme=… --tag nnnn=NNNN --outcome lock --as json | flow
   --allow-commands set-state --model rdr-write.toml --artifact record=…
   --artifact readme=… --tag nnnn=NNNN --plan -` exits 0. `--tag` is bound on
   BOTH halves and this is load-bearing, not symmetry: `resolve` runs the
   narrowed readers before its kernel call, so the C1.6 readme reader's
   `{tag.nnnn}` must already be bound there or the pipeline refuses at
   `resolve`, before `set-state` is ever reached. Both verbs register the flag
   (`internal/cli/flow_resolve.go`, `internal/cli/flow_state.go`, via
   `registerTagFlag`).
3. End state: the record's Status line reads `Final` (the joint-decision record
   reads `Final [joint decision → …]`, qualifier kept by backreference; the
   wrapped record keeps its continuation line byte-identical and its reader
   reports the qualifier whole); the
   README row's status cell reads `Final`; both confirmed by read-back through
   the command reader; every other byte of both files is unchanged
   (`git diff --stat` shows one line per file).
4. A `none` or `stopped:*` resolve row applies nothing and exits 0 carrying
   `dispositions`.
5. Negative: the same pipeline without `--allow-commands` on `set-state`
   refuses before mutation and both files are byte-identical to before.
6. Negative: a README with the record's row duplicated refuses
   `edit_anchor_ambiguous`; files untouched.
7. Negative: a `replace` whose output no longer matches its own anchor refuses
   `edit_anchor_unstable`; files untouched.

### Phase 1: Grammar and lint

Admit `edit` on write entries — `table.Accessor.Edit`, the `edit.<key>` tables,
C1.4's categories in `carrierDefect`'s order, template and anchor parsed at load
(C1, C1.2, C1.4); dump/normalize round-trip carries it (A6).

### Phase 2: The binding

A mode-preserving, byte-oriented line-edit `WriteBinding` beside
`flowbind.go::Writer`: select, refuse, rewrite, stage-and-rename, `<clear>`
(C1.3, C1.5); `commandBacked` becomes the three-way discriminator.

### Phase 3: Seam carriage and gate pre-check

Context tags cross to `Apply` and `Read` (A1); `{tag.<key>}` joins the command
placeholder vocabulary at `cmdbind`'s substitution site (C1.6, A9); a gate-off
command reader refuses the write before mutation at `Executor.Write`'s
reader-resolution site (A2, C1.3).

### Phase 4: Surface and proof

The MVV fixture end to end; `docs/cli-output-contract.md` names the new
categories and Detail tokens; the consumer's `rdr-write.toml` migration is the
consumer's kata (rdr#yjye), not this plan's.

## Validation

### Testing Strategy

Done = every C1–C1.6 clause has a test, every refusal category fires on a fixture
that earns it, and the MVV's seven steps pass end to end. The normative fixtures
named on A3 and A4 are the expected values; the spike artifacts under
`evidence/spikes/` are what produced them.

**Load-time (C1, C1.4, C1.6) — table-driven over `intrastate lint --model`**

1. **Scenario**: An entry carrying `edit` beside `path`, and one carrying `edit`
   beside `command`; an `edit` on a read entry and on a gate entry.
   **Expected**: `edit_carrier_conflict` each time. An entry carrying only
   `edit` LOADS (0025:C5's "neither" arm no longer fires for it).
2. **Scenario**: A `keys` member with no `edit.<key>` table; an `edit.<key>`
   table for a key not in `keys`.
   **Expected**: `edit_key_mismatch` both ways.
3. **Scenario**: An anchor that fails RE2 compilation; one naming an undeclared
   tag key; one carrying a `{…}` form outside the closed vocabulary.
   **Expected**: `edit_anchor_invalid`. Compilation is checked with every
   `{tag.<key>}` replaced by a quoted probe (C1.4).
4. **Scenario**: A `replace` with an unknown placeholder, another key's
   placeholder, a group reference the anchor does not define, a malformed
   `${…}`/`{…}`.
   **Expected**: `edit_template_invalid` each.
5. **Scenario**: `clear` set to anything outside `{"line"}`.
   **Expected**: `edit_clear_invalid`.
6. **Scenario**: Category registration order and wire strings.
   **Expected**: The five `edit_*` strings append after 0025:C5's six, in C1.4's
   stated order; the wire strings are asserted, the identifiers and list size
   are not (0025:C5).
7. **Scenario**: A `{tag.<key>}` argv element on a read, gate and write entry;
   an undeclared `<key>`; a `{…}` element that is neither `{artifact}` nor a
   declared tag.
   **Expected**: C1.6 admits the declared form whole-element; the other two report
   `command_unknown_placeholder` (0025:C5's wire string, unchanged).

**Apply-time selection and refusal (C1.2, C1.3, C1.5) — over real file fixtures**

8. **Scenario**: The A3 Status anchor over every Status-bearing record in the
   corpus, and the per-record README row anchor over the linked index rows.
   **Expected**: exactly one line selected per in-scope file; no decoy (template
   comment, quoted example) scores ≥2. This is A3's corpus check as a regression
   test. The in-scope predicate is part of the assertion, because the corpus
   contains files that legitimately select zero: a record is Status-bearing when
   it carries a `- **Status**:` bullet, which excludes `*-postmortem.md` and
   `BUILD-ORDER.md` (A3's run reported exactly these as count=0). The invariant
   is PER FILE and the cardinalities are descriptive, never asserted — the
   corpus grows with every RDR seeded (33 records and 26 rows at A3's run; 34
   and 28 when this scenario was written, and this record adds one), so a test
   pinning a total is a scheduled false positive that fails for the correct
   reason.
9. **Scenario**: Anchor matching zero lines; matching ≥2; two rules selecting
   the same line.
   **Expected**: `edit_anchor_unmatched`, `edit_anchor_ambiguous`,
   `edit_anchor_collision` — each BEFORE any byte is written, file
   byte-identical after.
10. **Scenario**: A planned value, and a bound tag value, containing `\n` and
    containing `\r`.
    **Expected**: `edit_value_multiline`, pre-mutation, file untouched.
11. **Scenario**: A `replace` whose output no longer matches its own anchor.
    **Expected**: `edit_anchor_unstable` from the post-edit re-anchor pass,
    before any write (premortem P-4, P-13).
12. **Scenario**: A wrapped-qualifier fixture — a Status whose bracketed
    qualifier continues onto a second line — flipped `Draft`→`Final`. The
    fixture must be wrapped AND `Draft`, which the record A3 used for the
    byte-preservation half is not: `0021-cache-warmup-order.md` is already
    `Final`, so it demonstrated the preservation but never the value swap. The
    one that carries both is
    `rdr/tools/rdr/testdata/status/records/0022-cache-metrics-surface.md`
    (`- **Status**: Draft [revised from Final 2026-08-10; re-verify A1,A2 — the`
    wrapping to the next line).
    **Expected**: the anchored line alone is rewritten, the continuation line is
    byte-identical, `diff` shows exactly one changed line, and the role's reader
    reports the qualifier whole. (A3's normative case; this is what the spike was
    run to settle.)
13. **Scenario**: A record whose Status VALUE sits on the continuation line.
    **Expected**: zero matches → `edit_anchor_unmatched`, not a rewrite.
14. **Scenario**: Escapes and group semantics, per template. In `replace`: `$$`
    → `$`; `{{`/`}}` → literal braces; a group that did not participate expands
    empty; a bare `$1`. In `anchor`: `\d{4}`, `a{2,3}`, a literal `{{`.
    **Expected**: the `replace` arms as C1.2 states, with bare `$1` refusing
    `edit_template_invalid`. In the anchor, every brace form reaches RE2
    untouched — `\d{4}` compiles as a quantifier and matches four digits, `{{`
    is not an escape — and none refuses `edit_anchor_invalid`. The two
    templates are NOT symmetric, which is the point of the scenario.
15. **Scenario**: Captured text and substituted values that themselves contain
    `${…}` or `{…}`.
    **Expected**: emitted as literal bytes, never re-scanned (C1.2 "parse once").
16. **Scenario**: A bound tag value containing RE2 metacharacters, used in an
    anchor.
    **Expected**: regexp-quoted — matches literally, cannot alter the pattern's
    structure (C1.2).
17. **Scenario**: `clear = "line"` with a planned `<clear>`; the same against an
    anchor matching zero lines; `clear` absent with a planned `<clear>`.
    **Expected**: line deleted with terminator and the key read back ABSENT;
    zero-match is SUCCESS with no write (0004:C11); absent `clear` refuses
    `edit_clear_undeclared` pre-mutation.
18. **Scenario**: The literal string `<clear>` reaching a `replace`.
    **Expected**: never substituted (0004:C11).

**Byte preservation and the write (C1.3) — A5's determinism results**

19. **Scenario**: CRLF line endings; a file with no final terminator; a bare
    `\r`, NEL and U+2028 in content; deletion of the final line of a file that
    had no final terminator.
    **Expected**: CRLF preserved per line; final-terminator state preserved
    either way; the three characters treated as line content, not terminators.
20. **Scenario**: A post-edit buffer equal to the input — an already-applied
    edit re-run against its own output.
    **Expected**: not written at all — no staging, no rename. The witness is the
    target's INODE, unchanged across the call: file content cannot separate "did
    not write" from "wrote identical bytes", but stage-and-rename always replaces
    the inode (A5), so an unchanged inode is the observable that proves absence
    of a write. No MVV step covers this — step 5's files are untouched because
    the write REFUSED (a different arm of the disposition table), and step 3 is
    the happy-path end state.
21. **Scenario**: Mode preservation on a 0644 target and on a 0755 target.
    **Expected**: `git diff --summary` empty in both cases. Against the
    fixed-0600 path the 0755 target emits ` mode change 100755 => 100644` —
    A5's fixture, and the reason mode is copied.
22. **Scenario**: A symlinked target.
    **Expected**: the symlink survives and points at the new content (C1.3
    resolves before staging). The contrasting behaviour — renaming onto an
    unresolved symlink path replaces the symlink with a regular file — is a
    property of `os.Rename` already demonstrated in A5's spike, and is stated
    here as the reason the clause resolves first, NOT as an assertion this
    scenario runs: C1.3 `target:` forbids the unresolved path unconditionally, so
    no code path under this contract can exhibit it.
23. **Scenario**: `Invocations()` after an `edit` apply, on an entry whose role's
    reader is FILE-backed, run with the gate OFF.
    **Expected**: counts `Apply` calls exactly as `flowbind.go::Writer` does
    (0004:C14); the write succeeds. Gate-off success is the witness for "no
    subprocess spawned, `--allow-commands` not consulted by the write" —
    `Invocations()` counts applies, not spawns, so it cannot prove that negative,
    and a gate-off run separates "not consulted" from "consulted and permitted".
    The file-backed reader is what distinguishes this from S25, where the gate-off
    run must REFUSE because the reader is command-backed (C1.3 `read-back:`); the
    two gate-off outcomes are opposite by design and the reader's carrier is why.

**Seam and read-back (A1, A2, C1.6)**

24. **Scenario**: Context tags bound on the invocation reaching an anchor's
    `{tag.<key>}`, and the same tags reaching a command reader's argv.
    **Expected**: one channel serves both (A1); an unbound tag refuses
    `execution_failure` before spawn, Detail naming the placeholder (C1.6).
25. **Scenario**: A write whose role's reader is command-backed while the gate
    is off.
    **Expected**: refuses BEFORE mutation, Detail naming the gate — never
    `read_back_incomplete` for a write that ran no command (A2, C1.3). Negative
    control: MVV step 5's whole pipeline without `--allow-commands`, both files
    byte-identical after.
26. **Scenario**: Read-back of a planned `Final` against a record carrying a
    bracketed qualifier.
    **Expected**: byte-equal — the reader reports `status=Final` with the
    qualifier on `status_form`/`status.qualifier` (A4's normative fixture).
27. **Scenario**: A refused edit's error surface.
    **Expected**: 0004's `execution_failure` with `Applied()` false; Detail
    carries the rule id `<id>.edit.<key>` and the reason token; the two
    entry-level preconditions name the gate and the placeholder instead, having
    no rule to name (JDR 0003 §D1). Refusal reporting follows
    `docs/cli-output-contract.md`: a scalar failure carries `param`, an
    aggregate carries `findings[]`, regardless of runtime cardinality.

**Clauses only a second input or a second invocation reaches**

27b. **Scenario**: A bound target that cannot be read — absent (ENOENT), a
    directory (EISDIR), and unreadable (EACCES).
    **Expected**: all three refuse BEFORE mutation with `execution_failure`
    carrying the OS error in the Detail — NOT `edit_anchor_unmatched`, which
    would launder a stale binding into an anchor result, and not
    `flowbind.go::load`'s absent-is-empty treatment, which exists so a flow
    artifact can be created on first write and is wrong for a carrier that
    fences creation out (C1.3 `input:`/`select:`). This is the one refusal in C1
    with no `edit_*` reason token, so the scenario also pins the Detail's shape
    where there is no rule to name.

28. **Scenario**: One entry carrying TWO simultaneous load-time defects — an
    `edit` beside `path` AND an uncompilable anchor; then a `keys` mismatch AND a
    malformed `replace`.
    **Expected**: the earlier category in C1.4's registration order is the one
    reported, and only it (`edit_carrier_conflict` in the first case,
    `edit_key_mismatch` in the second). Fail-fast is observable ONLY through a
    multi-defect input — with one defect per fixture every evaluation order
    passes, which is why C1.4 `precedence:` needs its own scenario.
29. **Scenario**: Two invocations against one artifact — a `clear = "line"` rule
    applying a planned `<clear>`, then a non-clear write through the same rule.
    **Expected**: the first deletes the line and read-back reports the key
    ABSENT; the second refuses `edit_anchor_unmatched`, because `edit` never
    appends (C1.3 `select:`). This is C1.5 `one-way:` — the only clause in C1
    whose behaviour spans two invocations, and the test that stops a caller from
    treating `clear` as reversible.
30. **Scenario**: An entry with NO carrier at all, exercised at runtime rather
    than at load.
    **Expected**: the registry builds a REFUSING binding, never a file binding —
    every declared key reads absent through a file binding, which would confirm
    an unapplied write (C1.1 `in-memory:`, 0025:C1's runtime arm). The guard
    matters because `edit` widens the discriminator: `registry.go::commandBacked`
    is `len(acc.Command) != 0 || acc.Path == ""` today, and its `acc.Path == ""`
    arm IS the residue path, so adding the `edit` arm could silently reroute
    residue to a file binding.
31. **Scenario**: A decoy fixture where the anchor's exactly-one match is the
    WRONG line — a quoted example or template comment that outscores the real
    bullet — so the edit applies to it.
    **Expected**: the write applies, and the read-back through the role's reader
    then refuses `read_back_mismatch` with `Applied()` TRUE — the one arm where
    the applied sense is set on a failure (A8), and the post-mutation net F4 and
    F5 name as the sole defence against silent wrong-line corruption. Pre-mutation
    refusal cannot reach this class by construction, which is why it is tested
    here rather than folded into S27.

**Acceptance** — the MVV's seven steps, run against the `rdr-write.toml`
fixture, are the end-to-end gate; steps 5–7 are the negative half.

### Performance Expectations

Whole-file read, in-memory rewrite, staged write and rename, per entry. The
artifacts in scope are RDR records and an index README — kilobytes — so the cost
is dominated by two syscall round-trips, not by matching. No measurement is
claimed beyond that: nothing here is on a hot path, and the RDR introduces no
loop over files (one entry rewrites ONE buffer in ONE write, C1.3).

Anchor resolution is linear in file lines per rule, with RE2's linear-time
matching guarantee (no backtracking, and C1.2 pins RE2 precisely so a pathological
anchor cannot be written). Rules per entry are bounded by the entry's `keys`.

**Byte-stability is a contract here, so the determinism checklist applies** —
results from A5's spike and C1.3's clause:

- *Hash fn / lib*: none — no hashing on this path.
- *Pre-image byte layout*: the whole file is read as bytes; every byte outside a
  selected line is re-emitted unchanged.
- *Encodings*: no transcoding; bytes pass through. Values are literal bytes,
  never re-scanned (C1.2).
- *Map order*: not reachable — rules resolve against pre-edit line INDICES held
  per rule, so iteration order over the `edit.<key>` tables cannot affect the
  output buffer (C1.3 "select").
- *Whitespace*: preserved; `replace` names the WHOLE replacement line, so
  leading and trailing space is the template author's, not the tool's.
- *Case folding*: none.
- *Empty / null / absent*: a non-participating capture group expands to the
  empty string (C1.2 — the rule matches `regexp.Expand`'s, but the parse is
  C1.2's own, not a call to it); a `<clear>` against a zero-match anchor is
  success with no write; an equal post-edit buffer is not written at all.
- *Version marker*: none in the artifact; the model's `edit` grammar is versioned
  by the closed vocabulary and extended only by successor amendment (C1.2).

Line terminators are the one place byte-stability could silently drift, and C1.3
fixes it: only `\n` (optionally preceded by `\r`) terminates a line, CRLF stays
with its line, and a missing final terminator is preserved. Re-running an
applied edit is a no-op by construction — the re-anchor pass (C1.3) requires each
rule to select exactly its own rewritten line, and an unchanged buffer is not
written.
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

- RDR 0004 (`0004:C10`–`0004:C15`, Alt 2, Alt 5); RDR 0025 (`0025:C1`, `0025:C2`,
  `0025:C5`, `0025:C6`, Alt 2, Consequences, Selection LBD); RDR 0010 (`0010:C4`);
  RDR 0005 (`0005:C1`).
- JDR 0003 §D1 (`docs/jdr/0003-accessor-binding-seam.md`) — the `execution_failure`
  sub-reason carrier; cited, not restated.
- `internal/table/load.go::carrierDefect`, `::commandPlaceholders`;
  `internal/cli/flowbind/registry.go::commandBacked`;
  `internal/cli/flowbind/flowbind.go::save`, `::Writer.Apply`;
  `internal/accessor/binding.go::WriteBinding`; `internal/accessor/model.go::Artifact`;
  `internal/accessor/executor.go::Executor.Write`; `internal/cli/flow_state.go`.
- Ansible `ansible.builtin.lineinfile` (module DOCUMENTATION: `regexp`, `line`,
  `backrefs`, `insertafter`, `create`; `atomic_move`); Puppet `puppetlabs-stdlib`
  `file_line` (`match`, `multiple`, `append_on_no_match`, `replace`);
  kubebuilder `pkg/plugin/util/util.go::InsertCode`, `::ReplaceInFile`.
- The consumer's `models/rdr-write.toml` header ("THE EDIT IS DATA —
  TRANSITIONALLY") and its lock rows.
- intrastate#zyh0 (tracker), intrastate#c3xz (closed, `set-state --plan`),
  intrastate#v0hb (closed, 0025), rdr#yjye, rdr#qmkd.
- Evidence: `evidence/research/prior-art.md`, `evidence/propose-premortem/critic.md`.
