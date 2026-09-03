Model: claude-opus-5

# Persona 3 — QA / Tester

Reviewed the owned set (`S1`–`S27`, `MVV`, `C1`) and widened twice, both times
because a pass/fail criterion could not be settled inside the scenario text:

- **To `A3` and `evidence/spikes/a3-run.out` / `a3-anchor-regex.md`** — sent by
  `0028:S8`, whose Expected ("exactly one line selected per file") is a corpus
  claim whose truth is only checkable against the spike that produced it.
- **To `§failure-modes` (`F4`, `F5`) and `§pre-lock-mini-checks`** — sent by
  `0028:S20` and `0028:S27`, whose witnesses are not named in the scenario and
  which the mini-check tables cite with a different cross-reference than the
  scenarios do.

Every cited code anchor, doc and fixture I checked exists:
`docs/cli-output-contract.md`, `internal/cli/flow.go::registerTagFlag`,
`internal/cli/flow_input.go::parseTags` (registered on `set-state` at
`internal/cli/flow_state.go:199`), `internal/accessor/executor.go:302`
(`readerFor` before `Apply`), `internal/cli/flowbind/flowbind.go:139/162`
(`save`, fixed `0o600`), `internal/accessor/model.go:249` (`Artifact{Role,Path}`,
no context field), `internal/accessor/model.go:411` (`IsClear`),
`internal/cli/cmdbind/cmdbind.go::substitute`, `model.ValueContinues`
(`rdr/tools/rdr/internal/scan/fields.go:146`), and
`rdr/tools/rdr/testdata/status/records/0021-cache-warmup-order.md`. Findings
whose premise did not hold on disk were dropped and are not listed.

The record is unusually testable: all 27 scenarios carry an explicit
**Expected** line, none passes by absence-of-error, and the `disposition` table
enumerates every input class with its category. The findings below are the
residue.

---

### 1. S8's pass criterion is false against the corpus it names — the regression test fails on authoring

**Anchor:** `0028:S8` (widened to `0028:A3` and
`evidence/spikes/a3-run.out`)

`S8` states the Expected as, without qualification, "exactly one line selected
per file" over "all 33 record fixtures". The spike that produced this corpus
reports the opposite for 5 of those 33 files:

```
Total record files scanned: 33
Files where Status-anchor count != 1: 5
  BAD FILE: .../0003-guard-predicate-exhaustiveness-postmortem.md count=0
  BAD FILE: .../0007-guard-predicate-totality-postmortem.md count=0
  BAD FILE: .../0008-recognized-tag-key-ownership-postmortem.md count=0
  BAD FILE: .../0011-flow-next-match-conditioned-candidates-postmortem.md count=0
  BAD FILE: .../BUILD-ORDER.md count=0
```

Verified on disk: those files genuinely carry no `- **Status**:` bullet
(`grep -c` returns 0 for the postmortem; `BUILD-ORDER.md` is not a record). The
spike's own prose demotes them to "informational, not necessarily a defect if
not yet a full record" and its Verdict narrows the claim to "every
*Status-bearing* record scores exactly 1" — but `S8` carries neither the
"Status-bearing" qualifier nor any file-selection predicate. A test author
writing the regression from `S8` alone globs `docs/rdr/*.md`, asserts
`count == 1` per file, and gets 5 red rows with nothing in the record saying
which of the two — corpus or assertion — is wrong.

**Missing pass/fail criterion:** the in-scope predicate for the corpus sweep.
`S8` needs either "every Status-bearing record" plus the rule that identifies
one, or an explicit exclusion list (postmortems, `BUILD-ORDER.md`).

**Blocks:** the C1.3 `select:` corpus regression test — the single test the
record designates as "A3's corpus check as a regression test", and the only
standing defence `F5` names against the decoy risk.

---

### 2. S8's counts are snapshot-pinned to a corpus that has already moved

**Anchor:** `0028:S8`, `0028:MVV` (the `trace` mini-check rows 3a and 3d also
carry "1 of 33" and "1 of 26")

`S8` fixes the corpus cardinality in the Expected: "all 33 record fixtures" and
"the 26-row index". On disk today `docs/rdr/*.md` is **34** files and
`docs/rdr/README.md` carries **28** rows (`| [0026](...)` linked form plus the
two unlinked). The counts drifted between the spike run and now, and they will
drift again on the next RDR seeded — this record itself adds one.

A regression test that asserts `len(files) == 33` or `matches == 26` is a test
that fails for the correct reason (a new record landed) and reports as a
carrier defect. Nothing in `S8` says whether the count is a contract or an
incidental of the snapshot.

**Missing pass/fail criterion:** whether the cardinalities are asserted or
merely descriptive. The invariant that survives corpus growth is per-file
("each Status-bearing record selects exactly one"), never the total.

**Blocks:** a durable form of the same corpus regression test — as written it
is a scheduled false positive.

---

### 3. S20's no-write assertion has no obtainable witness, and its two cross-references point at different, wrong MVV steps

**Anchor:** `0028:S20`, `0028:C1` (C1.3 `write:`), widened to
`§pre-lock-mini-checks` (the `fidelity` and `disposition` tables)

`S20` asserts a post-edit buffer equal to the input is "not written at all — no
staging, no rename", and delegates its witness to "MVV step 5's untouched-files
assertion". Two problems:

1. **The witness cannot distinguish the two outcomes.** "Wrote byte-identical
   content" and "did not write" produce the same file content. The record does
   name an observable that *would* separate them — `A5` establishes "the inode
   always changes" on a real stage-and-rename — but neither `S20`, `C1.3`, nor
   the `disposition` table ("silent by design (S20)") names inode, mtime, or a
   staging-file probe as the assertion. As written the scenario passes by
   absence of a difference that a write would not have produced either.

2. **The cross-reference is wrong, and inconsistently wrong.** `MVV` step 5 is
   the *gate-off refusal* negative ("the same pipeline without
   `--allow-commands` refuses before mutation"), not a no-op-equal-buffer case;
   its untouched files are untouched because the write **refused**, which is a
   different arm of the `disposition` table entirely. Meanwhile the `fidelity`
   table cites the same invariant as "C1.3 `write:`; S20; **MVV 3**", and step 3
   is the happy-path end-state assertion, also not a no-op. `C1.3` `write:`
   itself repeats the step-5 citation. No MVV step exercises an equal-buffer
   edit.

**Missing pass/fail criterion:** the observable that proves absence of a write,
and a correct referent for it.

**Blocks:** the C1.3 `write:` no-op test — the one that holds "re-running an
applied edit is a no-op by construction", which `§performance-expectations`
states as a byte-stability contract.

---

### 4. S12 and the MVV name a wrapped-qualifier fixture that cannot perform the flip they specify

**Anchor:** `0028:S12`, `0028:MVV` (widened to `0028:A3`)

`S12` and `MVV` step 3 both require a wrapped-qualifier record **flipped
`Draft`→`Final`**. The fixture `A3` names for the wrapped shape is
`rdr/tools/rdr/testdata/status/records/0021-cache-warmup-order.md`, whose Status
line reads `- **Status**: Final [joint decision → …` — it is **already
`Final`**. The spike says so explicitly: "word `Final` is already the target
value in this demo, so no value swap was needed to show the truncation."

So the record's one named wrapped fixture cannot exercise the value swap its own
scenario specifies; the byte-preservation half was demonstrated, the swap half
was not. A suitable fixture does exist one directory over —
`.../0022-cache-metrics-surface.md` carries
`- **Status**: Draft [revised from Final 2026-08-10; re-verify A1,A2 — the` with
the qualifier wrapping to line 10 — but the record never names it.

**Missing pass/fail criterion:** a located fixture that is both wrapped and
`Draft`. This is the "fixtures named but not located" case: the name resolves,
the fixture does not fit the scenario.

**Blocks:** `S12` as written, and `MVV` step 3's wrapped-record clause — the
case the record says "the spike was run to settle".

---

### 5. C1.4's `precedence:` clause has no scenario

**Anchor:** `0028:C1` (C1.4 `precedence:`), `0028:S6`

C1.4 states a normative ordering: "within one entry, fail-fast in the order
above, evaluated after 0025:C5's clauses 1–6". `S6` is the only registration
scenario and it covers something else — the wire strings and their append order
in `table.Categories()`, explicitly *not* the evaluation order.

No scenario presents an entry carrying **two** simultaneous defects (say an
`edit` beside `path` *and* an uncompilable anchor, or a `keys` mismatch *and* a
bad `replace`), so nothing states which single category a fail-fast loader must
report. Fail-fast is only observable through a multi-defect input; with
one defect per fixture, any evaluation order passes.

**Missing pass/fail criterion:** the expected reported category for a
multi-defect entry, and whether lint reports one finding or all of them (the
record's `disposition` table says lint "names entry/key/category" but not
cardinality; `S27` fixes the *runtime* refusal's shape via
`docs/cli-output-contract.md` — scalar `param` vs aggregate `findings[]` — and
that doc exists, but `S27` is scoped to refused edits, not to lint).

**Blocks:** the C1.4 precedence test.

---

### 6. C1.5's `one-way:` clause has no scenario

**Anchor:** `0028:C1` (C1.5 `one-way:`), `0028:S17`

C1.5 states a testable consequence: "a deleted line cannot be re-established by
`edit` (no append, C1.3), so after a `clear` the next non-clear write refuses
`edit_anchor_unmatched`". `S17` covers the three `<clear>` arms (delete,
zero-match success, undeclared refusal) but stops at the deletion — no scenario
sequences a `clear` followed by a write against the same artifact.

This is the one clause in C1 whose behaviour is only visible across **two**
invocations; every other clause is single-shot, which is likely why it was
missed.

**Missing pass/fail criterion:** none stated for the second invocation. (The
clause's own text supplies the expected token, so this is a missing scenario
rather than an ambiguous one — but a missing scenario has no line range, which
is why it is anchored to the clause.)

**Blocks:** the C1.5 `one-way:` round-trip test — the test that stops a caller
from assuming `clear` is reversible.

---

### 7. C1.1's `in-memory:` residue arm has no scenario

**Anchor:** `0028:C1` (C1.1 `in-memory:`), `0028:S1`

C1.1 carries a runtime clause beside the load-time one: "the registry's residue
rule (0025:C1 runtime arm) is unchanged — an entry with no carrier builds a
refusing binding, never a file binding". `S1` tests the load-time carrier
conflicts and confirms an `edit`-only entry LOADS, but no scenario constructs
the residue case and asserts a refusing binding.

The clause matters precisely because `edit` changes the discriminator:
`D-selection-predicate` makes `registry.go::commandBacked` three-way, and the
current implementation is `return len(acc.Command) != 0 || acc.Path == ""`
(verified via `A6`) — an expression whose `acc.Path == ""` arm is exactly the
residue path being widened. The regression risk is that adding the `edit` arm
silently reroutes residue to a file binding.

**Missing pass/fail criterion:** the observable that distinguishes a refusing
binding from a file binding at runtime.

**Blocks:** the C1.1 residue test guarding the three-way discriminator change.

---

### 8. F4 and F5's silent-risk claims have no scenario asserting the post-mutation refusal

**Anchor:** `0028:F4`, `0028:F5` (widened from the owned set; sent by `S8`,
which the record designates as the only defence against F5)

Both failure modes turn on a specific, testable claim: a wrong-line or decoy
rewrite is caught because "the read-back through the role's reader refuses
`read_back_mismatch`". That token — `read_back_mismatch` — appears in no
scenario. `S26` tests successful read-back of a planned `Final`; `S25` tests
`read_back_incomplete` being *avoided*; nothing constructs a decoy fixture, lets
the edit apply to the wrong line, and asserts the mismatch fires
post-mutation with the applied sense set.

`F4`/`F5` are the record's stated answer for the one hazard class it cannot
refuse before mutation, so the claim carries weight; it is also the only place
where `Applied()` is **true** on a failure, which `A8` calls out as the sense
distinguishing this from every pre-mutation refusal.

**Missing pass/fail criterion:** no scenario, hence no Expected, for the
post-mutation `read_back_mismatch` arm.

**Blocks:** the F4/F5 silent-corruption test — the test that proves the safety
net the failure-mode analysis leans on actually closes.

---

### 9. S22's negative control asserts behaviour the implementation is contracted never to exhibit

**Anchor:** `0028:S22`, `0028:C1` (C1.3 `target:`)

`S22`'s Expected has two halves. The first is testable: the symlink survives and
points at the new content. The second — "the negative control — renaming onto
the symlink path replaces the symlink with a regular file" — describes what
happens when symlinks are *not* resolved, which C1.3 `target:` forbids
unconditionally ("symlinks resolved").

There is no code path under this contract that renames onto an unresolved
symlink path, so the control cannot be run against the implementation. It is a
property of `os.Rename`, already demonstrated in `A5`'s spike. As a scenario
line it reads as an assertion a test author is expected to write, and they
cannot without stubbing out the clause under test.

**Missing pass/fail criterion:** whether this half is an assertion or a
restatement of spike context. Low severity — the first half is testable and
carries the clause.

**Blocks:** nothing outright; it makes `S22` un-completable as literally
specified.

---

### 10. S23's "no subprocess spawned" names no observation method

**Anchor:** `0028:S23`, `0028:C1` (C1.3 `no subprocess:`)

`S23`'s Expected has three parts. Two are directly checkable —
`Invocations()` counts `Apply` calls as `flowbind.go::Writer` does, per
0004:C14. The third, "no subprocess spawned, `--allow-commands` not consulted by
the write", names no witness: `Invocations()` counts applies, not spawns, and
the record does not say what a test observes to prove a negative about process
creation.

The natural witness is available — running the edit path with the gate off and
asserting success, which distinguishes "not consulted" from "consulted and
permitted" — but `S23` does not name it, and `S25` uses a gate-off fixture for
the opposite assertion (a command-backed *reader* must refuse). A test author
reading both together has no stated way to tell which gate-off outcome is
correct for a pure-`edit` entry whose reader is file-backed.

**Missing pass/fail criterion:** the observable for "no subprocess", and the
expected gate-off outcome for an `edit` entry with a non-command reader.

**Blocks:** the C1.3 `no subprocess:` test.
