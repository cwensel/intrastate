# Verification — RDR 0028 declared-line-edit-writer

Phase 3a Chain-of-Verification. Independent adversarial re-derivation of the
REQ set from `0028-declared-line-edit-writer.md` and `req-list.md` ONLY: for
each clause, a concrete input was named that a correct implementation must
visibly refuse or honour, and each was then RUN against the shipped code
through throwaway probes (unit-level parser calls, `table.Load` over
hand-authored TOML models, `EditWriter.Apply` over real temp files, and
`accessor.Executor.Write` over a two-definition registry). Phase 1's test
files were not read.

Probes covered, by clause:

- **C1.1** — edit-only entry loads; empty `[write.x.edit]` beside `path`;
  `edit` on read and gate entries; the residue/no-carrier runtime arm.
- **C1.2** — both closed vocabularies, adversarially: 14 refused `replace`
  forms (`$1`, `$name`, `{artifact}`, `{tag.nnnn}`, `{other}`, `${}`, `${`,
  `{`, `${9}`, `${nosuch}`, `$`, `a$b`, `${1x}`, `{unclosed}`) and 10 admitted
  ones; 8 admitted / 7 refused `anchor` forms; regexp-quoting of a bound tag
  value carrying `.*|x(`; group-arity stability under both the probe and a
  `(x)(y)` binding; no re-scan of captured text or substituted values.
- **C1.3** — select-exactly-one; re-anchor identity (self de-anchoring AND
  drift onto a sibling's line); the six-step `precedence:` ladder; the
  read-back gate pre-check under both `protected` shapes; byte preservation
  (CRLF, mixed terminators, no final newline, bare `\r`, NUL, UTF-8
  multibyte, `(?m)`/`(?s)` anchors); mode preservation at 0644 and 0755;
  symlink survival; the no-op arm witnessed by an unchanged inode;
  ENOENT/EISDIR refusing before mutation; `Invocations()`; exit group,
  `Applied()` false, and class unchanged.
- **C1.4** — all six wire strings present and in registration order after
  0025:C5's six; `edit_tag_argv0` on read, gate and write entries;
  `{artifact}` at argv0 still admitted; `{tag.<key>}` at a later position
  admitted and an undeclared key still `command_unknown_placeholder`;
  `edit_key_mismatch` both directions; `edit_clear_invalid` over six
  out-of-set spellings; anchor tag keys of owned and recognized provenance.
- **C1.5** — `<clear>` deletes line and terminator (middle line, final line
  with no terminator, CRLF); zero-match `<clear>` succeeds with no write;
  `>=2` still refuses ambiguous; `clear` absent refuses
  `edit_clear_undeclared` ahead of a stale anchor; the literal `<clear>` in a
  `replace` literal segment emits as bytes; one-way (clear, then a non-clear
  write refuses `edit_anchor_unmatched`).
- **C1.6** — whole-element recognition (`x{tag.k}`, `{tag.k}y`, `{tag.}`,
  `{tag.a{b}}` all rejected as placeholders); argv0 position rule; declared
  membership; unbound tag refusing before mutation and naming the
  placeholder.

Two clauses were observed to behave differently from the record. Both are
reproducible; neither was fixed here (Phase 3c owns fixes).

---

## FAIL-1 — a bare `{{` in an `anchor` refuses `edit_anchor_invalid`

**REQ**: REQ-121 (0028 S14, §validation), with REQ-25 and REQ-26 as ground.

**Clause quote** (S14 Expected, verbatim):

> In the anchor, every brace form reaches RE2 untouched — `\d{4}` compiles as
> a quantifier and matches four digits, `{{` is not an escape — and none
> refuses `edit_anchor_invalid`.

S14's Scenario names the three anchor forms by name: "In `anchor`: `\d{4}`,
`a{2,3}`, a literal `{{`." REQ-25 restates the same posture — "every other
`{`, `}` and `$` is passed to RE2 untouched" — and REQ-26 bounds the category
with "never for a brace RE2 itself accepts". `regexp.Compile("{{")` returns
nil error and the compiled pattern matches the literal text `{{`, so `{{` is
a brace RE2 accepts.

**Exact failing input** — a complete model loaded through `table.Load`:

```toml
[read.r]
role = "rec"
path = "a.b"
keys = ["status"]
timeout = "2s"

[write.w]
role = "rec"
keys = ["status"]
timeout = "2s"
read_back = true

[write.w.edit.status]
anchor = "{{"
replace = "L: {status}"
```

**Observed**: the load refuses.

```
edit_anchor_invalid: write w `edit.status` anchor: the `{` at offset 0 opens
neither a `{tag.<key>}` placeholder nor an RE2 repeat quantifier; escape it
as `\{` to anchor a literal brace
```

The two sibling forms S14 names in the same breath, `\d{4}` and `a{2,3}`,
both load clean under the identical fixture, which isolates the divergence to
the `{{` arm alone.

**Specified**: the model loads; the anchor reaches RE2 untouched; no
`edit_anchor_invalid`.

**Site**: `internal/table/edit.go::ParseEditAnchor`, the
`pattern[i] == '{' && !editRepeatSpec(pattern[i:])` arm, and its helper
`editRepeatSpec`, which admits only `{N}`, `{N,}` and `{N,M}`.

**Relation to the recorded deviation D5.** `deviations.md` D5 records the
brace-SHAPE discriminator as a deliberate resolution of a genuine conflict
between C1.4's "carries any other `{…}` form" and C1.2's "never for a brace
RE2 itself accepts", and that reasoning stands. What does NOT hold is D5's
closing claim that the choice "satisfies both clauses exactly on the fixtures
each pins": D5 names the pinned anchor fixtures as `^\d{4}$`, `^a{2,3}$` and
`^\{\{x$` — the third ESCAPED. S14's Expected names a bare, unescaped `{{`,
which the escaped fixture does not exercise and which the shipped
discriminator refuses. So the deviation's scope statement is narrower than
the behaviour it licenses, and the one form S14 names to distinguish the two
templates' escaping dialects is the form that refuses. Recorded as a FAIL
rather than folded into D5 because the divergence is from S14's Expected text
and D5 does not currently disclose it.

---

## FAIL-2 — the read-back gate pre-check is decided BEFORE, not after, C1.6's unbound `{tag.<key>}`

**REQ**: REQ-51 (0028:C1.3 `precedence:`, §normative-contracts); REQ-41 and
REQ-88 as ground.

**Clause quote** (REQ-51, verbatim):

> precedence: apply-time refusals fail-fast within one entry in this order,
> the mirror of C1.4's for load time: (1) the ENTRY-level preconditions —
> C1.6's unbound `{tag.<key>}`, then C1.6's `-`-prefixed bound value (both are
> properties of the binding, decided together before anything is spawned or
> read), then the gate-off command read-back above

The three preconditions of step (1) are ordered, and the gate-off read-back
is named LAST of the three.

**Exact failing input** — an `edit` write entry whose anchor carries an
unbound `{tag.nnnn}`, whose role reader is command-backed, and the
`--allow-commands` gate OFF. Driven through `accessor.Executor.Write`:

- write `w`: role `rec`, keys `["status"]`, `read_back = true`,
  `edit.status = { anchor = "^{tag.nnnn}: ", replace = "L: {status}" }`
- read `r`: role `rec`, keys `["status"]`, `command = ["tool", "{artifact}"]`
- registry `AllowCommands: false`; artifact context map EMPTY, so `nnnn` is
  unbound
- plan: `status = "Final"`

**Observed**: the gate refusal is the one reported.

```
the read-back for the write accessor `w` goes through the command-backed
reader `r`, which requires the allow-commands opt-in (--allow-commands);
refusing before mutation rather than writing and reporting an unverified
read-back
```

**Specified**: the unbound-placeholder refusal is reported, with a Detail
naming `{tag.nnnn}` (REQ-88: "an unbound tag at invocation refuses
`execution_failure` BEFORE spawn, Detail naming the placeholder").

**Site**: the two checks live on opposite sides of the `Apply` boundary. The
gate pre-check is `internal/accessor/executor.go::Executor.Write`, at the
`len(def.Accessor.Edit) != 0 && hasReader && len(reader.Accessor.Command) != 0
&& !e.Registry.AllowCommands` arm, which returns before `binding.Apply` is
called. The unbound-tag check is
`internal/cli/flowbind/edit.go::EditWriter.prepare`, at the
`table.ExpandEditAnchor` `!ok` arm — reachable only once `Apply` has been
entered. So the stated order is inverted by construction, and no ordering of
the code inside `prepare` can correct it.

**Scope of the consequence.** Both arms refuse before mutation and the
artifact is byte-identical either way, so this is a REPORTING-ORDER defect,
not a mutation-safety one: the caller is told to pass `--allow-commands`
when the model would still refuse for an unbound tag after they did. Both
refusals are entry-level and name no rule, so REQ-41's Detail discipline is
satisfied by either; what diverges is only which of the two the caller sees
first, which REQ-51 fixes explicitly.

---

## NOTES — suspicions probed and NOT reproduced

- **NOTE (C1.2, `replace` closure)**: the owned-key placeholder's PRESENCE is
  not required — a `replace` with no `{<key>}` loads and applies, writing a
  constant line. This matches req-list Q1's Reading A, taken for the audit,
  so it is recorded as confirmed behaviour rather than a defect. A reviewer
  expecting Q1's Reading B would read this as a gap.
- **NOTE (C1.3, gate pre-check predicate)**: the pre-check ANDs in
  `len(def.Accessor.Edit) != 0`, so it does not fire for a `path`-carried
  write with a command-backed reader and the gate off. This is deviations.md
  D6, recorded and grounded on C1.6's "the gate [is] untouched". Probed, not
  a FAIL.
- **NOTE (C1.3, re-anchor shift)**: the shift arithmetic was probed with one
  preceding deletion, two preceding deletions, a following deletion (no
  shift), and a zero-match `<clear>` contributing no deletion (A-9). All four
  behaved as REQ-39 and A-9 specify.
- **NOTE (C1.3, byte preservation)**: no case was found where a byte outside
  the selected line changed. Probed with CRLF, mixed CRLF/LF in one file, no
  final terminator, a bare `\r` mid-file, a NUL byte, UTF-8 multibyte text,
  and a regex metacharacter (`Fin.*al$(x)[y]`) as the interpolated value.
- **NOTE (C1.3, no-op arm)**: witnessed by `os.SameFile` across the call on
  both a single-rule and a multi-rule entry; the inode was unchanged in both,
  so nothing was staged or renamed.
- **NOTE (C1.4, registration)**: all six wire strings are registered, in
  clause order, strictly after 0025:C5's six. Asserted as RELATIVE order
  only, never a tail position or a count (REQ-76).
- **NOTE (unreproduced)**: an attempt to make a bound tag value alter the
  anchor's structure — binding `nnnn` to `.*|x(` and to `(x)(y)` — failed in
  both directions: `regexp.QuoteMeta` neutralised the metacharacters and the
  compiled arity was unchanged at 2. No structure-altering input was found.
- **NOTE (unreproduced)**: an attempt to make captured group text or a
  planned value be re-scanned as grammar — a document line carrying
  `${1} {k} $$ {{` and a planned value of `${1} {k} $$` — emitted literally in
  both cases. `parse once:` holds.

## Phase 3b — adversarial failure-mode probes (ADV-N)

Written independently of the Phase 1 suite and of Phase 3a's `FAIL-N`
entries. Anchored in the record's `Trade-offs / Failure Modes` section.
Tests: `internal/accessor/adversarial_0028_test.go`.

All three currently **FAIL** against the implementation at `414cb28`.

---

### ADV-1 — an entry-level precondition refuses in the WRONG exit group

**Failure mode.** `0028:FM`'s third Visible arm covers "the entry-level
preconditions refuse the same class before mutation … an unbound
`{tag.<key>}` … the artifact is untouched". `0028:C1.3` `order:` names the
three entry-level preconditions and then fixes their exit group in the same
sentence — "EXIT GROUP: these refusals are about the REQUEST, not the
environment, so they take the exit-2 group, not `execution_failure`'s
default exit 3 — the executor-facing typed `Err` discriminates at
`flow_exec.go::accessorFailureOf`".

`flowbind/edit.go::prepare` mints the unbound-`{tag.<key>}` refusal as a
bare `&accessor.ExecError{Detail: …}` with no `Err`, so
`executor.go::refusalOf`'s `errors.Is(err, ErrDeclaredRequest)` is false,
`Refusal.DeclaredRequest()` is false, and `accessorFailureOf` routes it to
`codeAccessorFailed` — the environment group, exit 3.

**Why a green suite misses it.** The class is `execution_failure`, the
Detail names the placeholder, and the artifact is byte-identical. Every
assertion the Visible arm invites passes. Only the discriminator differs.

**Why it matters.** Exit 3 promises "repair the environment and re-run the
same request unchanged". A missing `--tag` is the one defect re-running
unchanged can never fix; a retrying pipeline spins forever.

- **Test**: `TestAdv0028_1_UnboundAnchorTagTakesTheRequestExitGroup`
- **Status**: **FAILS**
- **Clauses**: `0028:C1.3` order:/EXIT GROUP (REQ-41, REQ-43),
  `0028:C1.6` binding: (REQ-88)

---

### ADV-2 — the gate-off read-back pre-check refuses in the WRONG exit group

**Failure mode.** The same Visible arm covers "a gate-off command read-back
(Detail names the gate, C1.3) — the artifact is untouched", and
`0028:C1.3` `order:` lists "the read-back gate below" **first** among the
entry-level preconditions its EXIT GROUP: sentence governs.

`internal/accessor/executor.go`'s pre-check builds its refusal with
`refusalOf(def, timeout, ClassExecutionFailure, nil)` — a literal `nil`
error — so `declaredRequest` is false *by construction* and no `Err` can
ever reach it.

This is a **second, distinct site** from ADV-1's: ADV-1's refusal is minted
in the binding and could carry `Err`; this one is minted in the executor,
which has no binding error to wrap. Fixing one leaves the other.

**Why a green suite misses it.** Identical to ADV-1, and MVV step 5 ("the
same pipeline without `--allow-commands` refuses before mutation and both
files are byte-identical") passes unchanged.

**Why it matters.** A forgotten `--allow-commands` is the archetypal
request defect — the record says so itself. Having correctly fenced out the
misleading APPLIED sense, the refusal hands the caller the misleading RETRY
sense instead.

- **Test**: `TestAdv0028_2_GateOffReadBackPreCheckTakesTheRequestExitGroup`
- **Status**: **FAILS**
- **Clauses**: `0028:C1.3` read-back:/order:/EXIT GROUP (REQ-41, REQ-43,
  REQ-59)

---

### ADV-3 — the re-anchor pass skips rules this plan did not name, so one rule silently poisons an unplanned sibling's anchor

**Failure mode.** A third Silent risk, beside the two `0028:FM` files —
and unlike those two, read-back cannot see it. `0028:C1.3` `re-anchor:`
states the invariant with no plan qualifier and then names this exact
consequence as its reason for existing:

> "after the buffer is rewritten in memory, **every rule's** anchor is run
> again over the POST-EDIT buffer and must select exactly its own rewritten
> line … This is what stops a replacement from de-anchoring itself or
> **poisoning a sibling rule's anchor on the next run** (premortem P-4,
> P-13)."

"on the next run" is decisive: a poisoned sibling only matters on a LATER
invocation — which is precisely the invocation that did not plan it now.

`flowbind/edit.go::prepare` `continue`s past any key with no planned value,
so such a rule never enters `plans` and `reAnchor` never sees it. The
fixture uses an **unpinned** sibling anchor, which `0028:C1.1` admits by
contract ("a line is SELECTED when the pattern matches anywhere in it …
authors pin `^…$`" — advice, not a rule).

Witnessed: an entry with keys `status` and `owner`, a plan naming only
`status`, whose replacement value makes the `owner` rule's anchor select
**two** lines of the post-edit buffer. The write applies and reports
SUCCESS.

**Why a green suite misses it.** The entry lints, the planned rule's own
re-anchor passes, the write applies, and read-back verifies `status` — the
only key in the plan. Nothing in the run touches `owner`.

**Why it matters.** The artifact is now in a state no invocation can repair
through this carrier. The next plan naming `owner` refuses
`edit_anchor_ambiguous`, correctly pointing at the `owner` rule — but the
line that broke it was written by the `status` rule an unknown number of
runs earlier, and `0028:FM`'s Recovery arm ("fix the anchor or the
artifact") cannot be followed without that history.

- **Test**: `TestAdv0028_3_ReAnchorCoversRulesThisPlanDidNotName`
- **Status**: **FAILS**
- **Clauses**: `0028:C1.3` re-anchor: (REQ-38, REQ-39), `0028:C1.3` write:
  ("all rules of one entry rewrite ONE buffer")

---

### Phase 3b recorded findings (not edits)

- `internal/table` reports one failure from `zzprobe2_test.go`, an
  **untracked scratch file belonging to a concurrent agent**, not to this
  pass and not to the RDR's surface. It is left in place and uncommitted.
  With it excluded, `internal/table` and `internal/cli` are green and
  `internal/accessor` fails only on the three ADV probes above.
- Not turned into a test, recorded for Phase 3c: `table/edit.go`'s
  `editRepeatSpec` refuses any unescaped `{` that is not a well-formed
  repeat spec, so `[{}]` and `[{]` — patterns RE2 accepts — are
  `edit_anchor_invalid` at lint. `0028:C1.2` says the category fires
  "never for a brace RE2 itself accepts", while `0028:C1.4` says it fires
  on "any other `{…}` form". The implementation documents its own
  reconciliation of the two in-code, so this is a clause conflict for 3c to
  arbitrate rather than an implementation defect to pin.
- Confirmed HANDLED, no test added: `<clear>` deletion preserves the file's
  final-terminator state either way; CRLF survives per line; a bare `\r`,
  NEL and U+2028 stay line content; an unreadable target (ENOENT, EISDIR)
  refuses before mutation with the OS error and does NOT wrap
  `ErrDeclaredRequest` (correct — that one really is the environment);
  cross-entry non-atomicity is stated by the contract, so an entry landing
  beside a sibling entry's refusal is accepted behaviour, not a defect.
- Confirmed ACCEPTED RISK, no test added: `stageAndRename` breaks a hard
  link to the target (the new inode diverges from the old). `0028` names
  this risk and scopes its mitigation to mode and symlinks only.

**Input restriction honoured.** The Phase 1 test files and Phase 3a's
artifacts were not opened. One `grep` over `internal/accessor/*_test.go`
for `flowbind` usage incidentally surfaced matching LINES from
`edit_apply_0028_test.go`, and one `grep` for package-level test
identifiers surfaced helper SIGNATURES from the Phase 1 accessor tests —
both to avoid name collisions. No assertion, fixture or expectation was
read, and every helper in `adversarial_0028_test.go` carries an `adv0028`
prefix.

## Phase 3c — dispositions

Every Phase 3a `FAIL-N` and Phase 3b `ADV-N` above is dispositioned here.
Entries above are unchanged; this section is append-only. Full suite green
and `golangci-lint run` at 0 issues after each commit.

Not weakened, skipped or deleted: no existing test was modified. The three
Phase 3b probes in `internal/accessor/adversarial_0028_test.go` and every
Phase 1 assertion pass unchanged against the fixed implementation.

---

### FAIL-1 — a bare `{{` in an `anchor` refuses `edit_anchor_invalid` — **FIXED**

**What changed.** `internal/table/edit.go`: the anchor's brace
discriminator was `editRepeatSpec` — "a well-formed repeat spec, or refuse".
It is now `editPlaceholderShape` — refuse only a `{` opening a `{word}`
(one or more of letter, digit, `_`, `-`, `.`, closed by `}`, with an
all-digit body excluded as a repeat spec).

**Why that reading.** D5's underlying reconciliation stands and was
re-verified: `regexp.Compile` returns nil for EVERY brace shape probed —
`{{`, `}}`, `{artifact}`, `\d{4}`, `[{}]`, `[{]`, `{`, `{}`, `{,}`,
`{1,2,3}`, `a{2,3` — so C1.2's "never for a brace RE2 itself accepts"
cannot be read as "never for a brace that compiles" without making C1.4's
"carries any other `{…}` form" unreachable. What D5 got wrong was the SHAPE.
C1.4's other two arms are "fails to compile" and "names an undeclared tag
key", both placeholder concerns; "any other `{…}` form" read in that company
means a form an author wrote MEANING a substitution, and the only
substitution an anchor admits is `{tag.<key>}`. The change is strictly wider
on the admit side and identical on the refuse side.

**Constraint set, all satisfied simultaneously:**

| anchor | verdict | source |
| --- | --- | --- |
| `^{{x$` (BARE) | loads | S14 Expected, by name — was refused |
| `[{}]`, `[{]` | load | C1.2; Phase 3b recorded — were refused |
| `^\d{4}$`, `^a{2,3}$`, `^\{\{x$` | load | REQ-25/28/121 — unchanged |
| `^{artifact}$`, `^{status}$` | refuse | REQ-17/26/65 — unchanged |
| `^{tag.absent}$`, `^{tag.nnnn$` | refuse | C1.2 placeholder scan — unchanged |

**Tests.** New: `TestEditBraceShape0028_PlaceholderShapeIsTheDiscriminator`
holds BOTH verdicts in one table, which the shipped Phase 1 pair does not —
`TestReq17And26And65_…` and `TestReq25And28And121_…` are separate, so a
future narrowing could satisfy one by breaking the other silently.
`TestEditBraceShape0028_AdmittedAnchorsAreCarriedVerbatim` holds S14's
operational half: an admitted anchor's bytes reach RE2 unmangled, not merely
un-refused. Eight of the table's arms fail against the pre-correction
discriminator, so the tests are load-bearing.

**Deviation.** D5 UPDATED in place per the brief: status now "CORRECTED
(Phase 3c)", the choice restated, and a "Why the original wording was wrong"
section naming both FAIL-1 and the Phase 3b finding. Type stays SPEC-UNDER.

---

### FAIL-2 — the gate pre-check is decided BEFORE C1.6's unbound `{tag.<key>}` — **NO CODE CHANGE; recorded as D11**

**Assessment: C1.3's stated precedence is satisfiable as-is.** Three
independent grounds, each sufficient; the full reasoning is deviations.md
D11 (SPEC-DEFECT).

1. **The two refusals step (1) names are not the one probed.** Step (1)
   names "C1.6's unbound `{tag.<key>}`" and "C1.6's `-`-prefixed bound
   value"; C1.6 is "`{tag.<key>}` as a COMMAND PLACEHOLDER", both of its
   refusals argv-side in `cmdbind.go::substitute`, and REQ-88's own covering
   test drives them through a `command` argv. FAIL-2's fixture put the
   placeholder in the ANCHOR, whose binding is C1.2's — a site C1.2 gives no
   refusal order and step (1) does not enumerate.
2. **On the path step (1) describes the order is unrealizable in EITHER
   direction.** C1.6's refusals fire inside `executor.go::invokeRead`, the
   read-back reader's spawn, which `Write` reaches only after the gate
   pre-check and after `Apply`. The gate is what prevents that spawn: with
   it off, neither C1.6 refusal is reachable at all. No fixture can witness
   the ordering. The shipped `TestReq51To54And120_…` covers steps (2)–(5)
   only, consistent with this.
3. **The record contradicts itself, and its other statement puts the gate
   first.** REQ-41, one sentence earlier in the same clause, orders the same
   three as "the read-back gate below, C1.6's unbound tag and C1.6's
   `-`-prefixed value". Phase 3b's ADV-2 read that order as governing and
   called the gate "first among the entry-level preconditions". The
   implementation matches REQ-41.

**Cost had one been owed**, recorded so the completion gate can price it:
the gate pre-check must run before `Apply` (C1.3 `read-back:`, A2) and the
anchor's binding runs inside it, so inverting them needs a NEW
`WriteBinding` seam method — new public surface — because
`internal/accessor` cannot import `flowbind` (C1.3 SITE:). Not taken.

Every property C1.3 asserts holds on both arms either way: refusal before
mutation, byte-identical artifact, `execution_failure` class, a Detail
naming the gate or the placeholder, and — after ADV-1/ADV-2 below — the
exit-2 group.

---

### ADV-1 — the unbound anchor tag refuses in the WRONG exit group — **FIXED**

**What changed.** `internal/cli/flowbind/edit.go::prepare`: the
`table.ExpandEditAnchor` `!ok` arm minted a bare `&accessor.ExecError{Detail:
…}`. It now carries `Err: accessor.ErrDeclaredRequest`, so
`refusalOf`'s `errors.Is` discriminator sets `DeclaredRequest()` true and
`accessorFailureOf` routes it to the exit-2 group. One field; the Detail,
class and pre-mutation guarantee are untouched.

**Test.** `TestAdv0028_1_UnboundAnchorTagTakesTheRequestExitGroup`
(Phase 3b's, unmodified) now passes.

---

### ADV-2 — the gate-off read-back pre-check refuses in the WRONG exit group — **FIXED**

**What changed.** `internal/accessor/executor.go`: the pre-check built its
refusal with `refusalOf(def, timeout, ClassExecutionFailure, nil)` — the
literal `nil` made the discriminator false BY CONSTRUCTION, as ADV-2 said, so
this needed a different fix from ADV-1's rather than inheriting it. It now
passes `ErrDeclaredRequest` in the same slot. The sentinel is not an
`*ExecError`, so `refusalOf`'s `errors.As` finds no tail and contributes no
Detail: the gate text assigned immediately below is unchanged and remains the
whole of it.

**Test.** `TestAdv0028_2_GateOffReadBackPreCheckTakesTheRequestExitGroup`
(Phase 3b's, unmodified) now passes.

---

### THIRD SITE (owed by C1.3's exit-group wording, not separately pinned) — **FIXED**

**Verified against C1.3 before fixing.** `order:` enumerates the entry-level
preconditions as a set of three — "the read-back gate below, C1.6's unbound
tag and C1.6's `-`-prefixed value" — and its EXIT GROUP: sentence speaks of
all of them at once ("THESE refusals are about the REQUEST"). ADV-1 and ADV-2
fixed two. C1.6's argv pair, minted through `cmdbind.go::refuse`, omitted the
typed `Err` for the same reason and is owed the same fix.

**What changed.** `internal/cli/cmdbind/cmdbind.go`: a new `refuseRequest`
helper, a SIBLING of `refuse` rather than a widening of it, used by exactly
the two C1.6 arms in `substitute` (the unbound `{tag.<key>}` and the
`-`-prefixed bound value). 0025's own pre-spawn ladder — the gate, the argv0
deny-list, a non-absolute `{artifact}` path — keeps `refuse` and its existing
routing, which this record does not amend.

**Test.** New:
`TestExitGroup0028_CommandTagPreconditionsTakeTheRequestExitGroup`, four
arms — the two refusals wrap `ErrDeclaredRequest` and still name the
placeholder; a bound non-flag value is still admitted (so the fix cannot be
read as "route everything to exit 2"); and a SCOPE CONTROL asserting 0025:C6's
gate refusal does NOT now wrap the sentinel.

---

### ADV-3 — the re-anchor pass skips rules this plan did not name — **FIXED**

**What changed.** `internal/cli/flowbind/edit.go`. `prepare` no longer
`continue`s past a key the plan did not name: it builds the rule for its
ANCHOR alone (`unplanned: true`), and `reAnchor` now runs every declared
rule. C1.3 `re-anchor:` states the invariant with no plan qualifier and names
this exact case as its reason for existing.

**The reading, which the clause under-determines.** Its test is "exactly its
own rewritten line", and an unplanned rule HAS no rewritten line — it writes
nothing. Taken as STABILITY: an unplanned rule must select, after the
rewrite, exactly the lines it selected before it, shifted by any deletions.
The alternative, cardinality, would newly refuse over an anchor already stale
or already ambiguous BEFORE this invocation — a defect this run neither
caused nor touched, and one C1.3 `select:` assigns to the plan that names the
rule. Stability fires when and only when this entry's rewrite moved the
sibling, which is the poisoning. Recorded as **D10 (SPEC-UNDER)**.

**Scope kept narrow, so no other clause moves.** An unplanned rule raises no
refusal of its own at steps (2)–(5): no planned value to scan, `replace`
never parsed or expanded, no deletion, no collision. An anchor it declares
that cannot be parsed, bound or compiled makes it contribute no post-edit
witness rather than condemn a plan that never touched it. Only step (6)
reads it. `shiftHits` reuses the shift arithmetic the planned rules already
use.

**Tests.** `TestAdv0028_3_ReAnchorCoversRulesThisPlanDidNotName` (Phase 3b's,
unmodified) now passes — and passes through its refusal branch, which
requires the artifact byte-identical, not through its accidental-pass branch.
New: `TestReAnchor0028_UnplannedSiblingRulesAreCovered`, the four arms the
ADV-3 fixture does not reach — an untouched sibling still applies; an
already-unmatched sibling is not this plan's refusal; a `<clear>` deletion
that merely SHIFTS the sibling is not poisoning; a rewrite that de-anchors a
sibling refuses before mutation.

**No public surface added:** `unplanned` and `preHits` are unexported fields
of an unexported struct and `shiftHits` is an unexported helper.

---

### Phase 3c new deviations

- **D10** — `re-anchor:` over an UNPLANNED sibling rule: the stability
  reading. **SPEC-UNDER.**
- **D11** — C1.3 `precedence:` step (1): the gate pre-check is decided
  before the anchor's unbound tag. **SPEC-DEFECT**, no code change.
- **D5** — UPDATED in place (not a new id): the brace discriminator
  corrected, with the reason the original scope statement was incomplete.
  Type unchanged (SPEC-UNDER).

No entry carries `Status: needs author decision`: each was settled from the
record's own text. D11 states what a later author would have to build if the
completion gate reads C1.3 step (1) the other way, and prices it as a
contract change rather than a bug fix.
