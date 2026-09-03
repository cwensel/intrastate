Model: claude-opus-5

Dual-model critique diff for RDR 0028. Pass A = `critique.md` (claude-opus-5, C-1..C-18).
Pass B = `critique-modelB.md` (claude-sonnet-5, C-1..C-10). Ledger IDs are per-file;
diffed by `0028:*` passage anchor, not by ID. All source claims re-grounded against
`internal/accessor/executor.go`, `internal/cli/flowbind/`, `internal/table/load.go`
and the `rdr` projector at HEAD.

---

## 1. Convergence table

Both passes, independently, landed on the same passage in seven places. Convergence
across models is the strongest signal this lens produces — these are the priority
defects.

| 0028 passage | pass-A | pass-B | shared claim | agree on failure mode, or only on the passage? |
|---|---|---|---|---|
| `0028:A10` + `0028:A11` vs `0028:C1` (C1.3 `read-back:` SITE, `order:`) | C-7, C-8 | C-2 | Both assumptions are `Pending` while C1.3 states their conclusions as settled declarative prose inside a `normative` fence — the exact `0028:G-assumptions` "Status consistency" bar. | **Same failure mode, and B states the process rule verbatim.** A goes further: it grounds A11's specific consequence (the deadline arm mints `applied = true`) against source and pairs it with the `disposition` table's "Every refusal … carries `Applied()` false" line, which has no timeout row. A's version is strictly stronger; B's is the cleaner gate citation. Merge as one defect, keep both halves. |
| `0028:C1` (C1.3 `re-anchor:`) | C-6 | C-3 | The re-anchor invariant is under-specified for multi-rule entries. | **Only partly the same failure mode — and they are complementary, not duplicative.** A: the invariant checks *cardinality*, not *identity* — after a sibling `clear` deletes a line, a sibling anchor can select exactly one line that is the WRONG line and still pass. B: the invariant is read as per-rule-own-line only, so a sibling whose anchor now matches zero or two lines is silently unchecked. Reading the clause as written ("every rule's anchor is run again … and must select exactly its own rewritten line"), B's reading is the weaker one — the clause DOES re-run every rule. A's is the live hole: "its own" is asserted but only a count is specified. Fold B's case in as the second half of the same fix. |
| `0028:C1` (C1.5 `one-way:`) | C-14 | C-6 | `clear = "line"` is irreversible and nothing at lint or apply enforces C1.5's prose warning ("never for a shared row"). | **Same failure mode; different symptom emphasis.** A's symptom is the destructive one and is the more serious: a `clear` on a README row deletes four cells of unrelated state and read-back reports the key ABSENT — SUCCESS. B's symptom is the recovery-path one (hand-edit outside intrastate reopens the drift class the RDR exists to close). Both are real; A's is the data-loss case. |
| `0028:C1` (C1.3 `write:` no-op arm) + `0028:S20` | C-9 | C-10 | The no-op write arm is witnessed only by inode identity and is invisible on the CLI surface. | **Same passage, materially different failure modes.** B: caller-visibility — an agent loop cannot distinguish "already applied" from "applied now", which matters for idempotency reasoning. A: the same, PLUS a contract-internal contradiction — C1.3 `no subprocess:` says `Invocations()` counts `Apply` calls, which increments either way, so the no-op is unobservable through the contract's own surfaces while `Invocations()` claims it happened. A's second half is the sharper finding. |
| `0028:§prerequisites` / `0028:F6` (0001, 0006 excluded from the MVV fixture) | C-16 | C-9 | The MVV fixture is curated to exclude exactly the two live README rows the anchor cannot address. | **Same failure mode, near-identical framing.** B adds the constructive remedy (make one of them a NEGATIVE fixture scenario proving the documented refusal, rather than excluding it) and the Scope-Verification gate hook. A adds the day-one arithmetic (the tool works on 32 of 34). Highest-confidence convergent row in the set. |
| `0028:A3` (snapshot verification of a growing corpus) | C-12 | C-1 | A3 is `Verified` against a frozen corpus and asserted as a property of the design; nothing (C1.4 explicitly disclaims it) prevents drift. | **Same passage, DIFFERENT axis — and A's axis is the one that bites first.** A: temporal — the corpus grows by the actions of this feature's own users, and the corpus is documentation *about* the format the anchor matches, so every critique/premortem/template excerpt is a decoy generator. B: cross-consumer — A3 proves nothing for a second adopter with a differently-shaped document family. Both hold. A's is in-fence (this consumer, this corpus, next week); B's is out-of-fence (a hypothetical second adopter, no such adopter exists). |
| `0028:C1` (C1.6) + `0028:§metadata` Profile `large` "one contract" | C-11 | C-7 | C1.6 is a second, independently-scored load-bearing contract riding inside a one-contract Profile claim; the record writes the split signal down ("C1.6 is scored separately, because the matrix above is about the write carrier and C1.6 is not one") and fences past it. | **Same failure mode, and both cite the same self-incriminating sentence.** B's symptom is sharper on the mechanics: a partial revert or partial rollout of C1.6 (a read-path change extending 0025:C2) is architecturally blocked by the Profile undercount. A's is sharper on process: the read-path change ships with no read-path premortem, no read-path alternatives matrix, and one scenario (S7). Merge; keep B's revert-granularity symptom. |

Verified non-overlaps that looked like candidates: A's C-15 (S31 `Applied()` TRUE)
and B's C-4 both touch `0028:S31`, but they are **not** the same claim — see §3,
they are in direct contradiction.

---

## 2. Divergence — rows raised by only one pass

### Unique to pass A

| pass-A | passage | claim | my verdict |
|---|---|---|---|
| C-1 / C-2 | `0028:C1` (C1.3 `read-back:`), `0028:A2` | `Executor.Write` runs a **pre-write baseline read** through `readerFor(role)` (`protectedKeys` → `invokeRead`) BEFORE `binding.Apply`; when the reader is command-backed and the gate is off that read already fails and is already swallowed into `baselineUnread`, and the write already proceeds. A2's `Verified` evidence omits it. | **CONFIRMED — see §4. The single most load-bearing finding in either pass, and pass B missed it entirely.** |
| C-3 | `0028:C1` (C1.3 `write:`) + `0028:§key-discoveries` | Atomicity is scoped to ONE entry; the MVV's journey spans TWO entries over TWO files (`write.record`, `write.readme`), and nothing states what the user has when entry 1 applies and entry 2 refuses. Key Discoveries markets the one-entry boundary as "stronger than Ansible … Puppet … Terraform". | **CONFIRMED.** C1.3 `write:` reads "all rules of one entry rewrite ONE buffer and land in ONE write" — scoped to one entry, verbatim. `0028:MVV` steps 2–3 bind `write.record` and `write.readme` in one pipeline. No section (`§consequences`, `§failure-modes` F1–F8, the `disposition` table) names the split-outcome case. The Key Discoveries bullet is present as quoted. Real, unstated, user-visible gap. |
| C-4 | `0028:MVV` | One `--tag nnnn=NNNN` addresses one record; the fixture declares THREE records; step 3 asserts all three flipped AND that `git diff --stat` shows one line per file. | **CONFIRMED — the MVV is not executable as written.** The fixture paragraph declares three records (`Draft`; `Draft [joint decision → …]`; wrapped-qualifier). Step 2 is a single `resolve \| set-state` with one `--tag nnnn=NNNN`. Three flips need three invocations; the shared README then accumulates three changed rows, so step 3's "`git diff --stat` shows one line per file" is false for the README. The two halves were written independently and never reconciled. |
| C-5 | `0028:C1` (C1.2 `anchor admits:`) | A literal `{tag.` in an anchor is unwritable: `{tag.` is the sole placeholder opener, scanned to its closing `}`, and C1.2 explicitly kills `{{`/`}}` as an escape in `anchor`. | **CONFIRMED as a gap; severity lower than A implies.** C1.2 says verbatim: "In `anchor` ONLY the fixed prefix `{tag.` opens a placeholder, scanned to its closing `}`" and "`{{`/`}}` is NOT an escape there". So `\{tag\.` — RE2's own backslash escape — is very likely the answer, since the scan is for the literal prefix `{tag.` and `\{tag\.` does not contain it. But C1.2 never says so, and `edit_anchor_invalid` fires "when a `{tag.…}` placeholder is malformed" — undefined for the escaped form. The defect is that the clause does not name the expression; it is not necessarily that none exists. **fix-now (one sentence), not blocking.** |
| C-10 | `0028:§problem-statement`, `0028:§capability-dependencies`, `0028:§prerequisites` | The C1.6 reader verb (`rdr index --row-json`) does not exist; the MVV cannot run until another repo ships it; nothing makes it a blocking prerequisite. | **CONFIRMED.** `rdr index` exists in the projector's verb table (`main.go:214`, `:307`, `:465`) but `--row-json` returns zero hits across the whole tool. `0028:§illustrative-code` says "verb illustrative; the consumer's to add." `§capability-dependencies` puts "the consumer adds the verb" in the Spec Impact column, where it is invisible as a blocker. `§prerequisites` lists only README normalization. Phase 4 ("The MVV fixture end to end") cannot run on the day Phase 3 merges. |
| C-13 | `0028:§existing-infrastructure-audit` (`save` row) | The Decision cell reads "Extend or sibling" — an unmade decision in the row whose purpose is to record the decision — and "extend" is not viable. | **CONFIRMED — see §4.** |
| C-15 | `0028:§testing-strategy` (S31) | S31 asserts `read_back_mismatch` with `Applied()` TRUE; the source's mismatch arm leaves `applied` UNSET. | **CONFIRMED — see §4. Directly contradicts pass B's C-4; see §3.** |
| C-17 | `0028:C1` (C1.3 `target:`) | "symlinks resolved" is a two-word obligation with no statement of final-component-vs-whole-path, no TOCTOU discipline, and no bound on where the derived staging directory may land. | **PLAUSIBLE, correctly reasoned, but the escape-the-root half overreaches.** The under-specification is real: C1.3 `target:` says only "symlinks resolved", and S22 asserts only that the symlink survives. But `§consequences` states the disposition A calls insufficient — "`edit` writes exactly where the caller's `--artifact` binding points, as `path` does; a caller who binds a hook file has bound a hook file" — and `path` today has the same property. A symlink is caller-supplied indirection over a caller-supplied path; it is not a privilege boundary this RDR erected. **Real spec gap (which component, TOCTOU); NOT a security hole. chart-as-net-new.** |
| C-18 | `0028:§performance-expectations` | The determinism checklist's "Map order: not reachable" is false, because C1.3 fixes no iteration order over the `edit.<key>` map and application order determines the intermediate buffer the re-anchor pass sees. | **PARTLY REFUTED — see §4. A's reasoning is sound for the general case but the specific mechanism it names is blocked by C1.3 `select:`. Downgrade.** |

### Unique to pass B

| pass-B | passage | claim | my verdict |
|---|---|---|---|
| C-4 | `0028:F4`, `0028:F5`, `0028:S31` | Read-back is the "sole defence" against silent wrong-line corruption, but reader-independence from the anchor is an ambient property of a consumer-authored reader, not a contract this RDR imposes. An author whose reader shares the anchor's blind spot gets corruption reported as success. | **CONFIRMED, and this is pass B's single best unique finding — pass A missed it.** F5 and S31 both name read-back as the defence; nothing in C1.1–C1.6 requires or checks that the reader's selection logic is independent of the anchor's. C1.4 explicitly disclaims lint coverage ("What lint does NOT prove: that an anchor matches exactly one line of a particular file"). B's remedy menu is the right shape: contract the independence, cross-check it, or **document it as an accepted UNMITIGATED risk in `§risks-and-mitigations`** — today it is implied to be mitigated. Option (c) is absorbable at lock. |
| C-5 | `0028:§decision-rationale`, `0028:ALT3` | ALT3 (adapter registry) is rejected because "document knowledge enters the tool", yet the MVV, the Testing Strategy, and the Illustrative Code all encode one consumer's document family into the acceptance bar — the "generic" framing is aspirational. | **PLAUSIBLE but the rejection stands.** Grounded: `§decision-rationale` says "The adapter form (D) is 0025 Alt 2 with a document family attached, and the kata fences document knowledge out"; the MVV and S8/S12/S13/S26 are indeed all RDR-flow-shaped. But the distinction is coherent and B does not engage it — document knowledge in the TOOL (a shipped adapter that knows markdown tables) is a different thing from document knowledge in a FIXTURE (a test corpus that happens to be markdown). ALT3 would put shape-awareness in shipped code; the MVV puts it in testdata. **Not a defect; the fixture-vs-tool distinction is sound.** |
| C-8 | `0028:A6`, `0028:§capability-dependencies` | A6 is `Verified` against nothing: it confirms "no accessor dump/marshal/normalize round-trip exists anywhere in `internal/`", so its own stated failure mode has no live target, and nothing in THIS RDR gates a future one. | **CONFIRMED as accurate reading, REFUTED as a defect of A6.** A6's evidence is exactly as B quotes it and the reasoning is correct. But an assumption whose failure mode has no live target is a **correctly** verified negative — that is what "no such path exists" evidence looks like, and A6's stated conclusion (a third carrier needs one parallel `if a.Edit != nil {…}`, not a rewrite) is independently confirmed against `internal/table/load.go`'s additive per-field copy. Wanting a guard against a future round-trip that does not exist is a request for out-of-fence future-proofing. **chart-as-net-new at most.** |
| B's premortem: A10's "test-only construction site" | `0028:A10` | A10's Pending evidence never audited non-production `accessor.Registry` construction sites; a test helper that does not thread `allowCommands` would silently mean gate-off. | **CONFIRMED as a live risk, and this is worth carrying into A10's verification.** `accessor.Registry{}` is constructed at nine sites in `internal/cli/cmdbind/seam_0025_test.go` (lines 237, 290, 359, 414, 464, 562, 641, 680, 836) plus the one production site `internal/cli/flowbind/registry.go:34`. A10's own "To check" list says "that adding the field breaks no other `accessor.Registry` construction site (tests included)" — so the record ANTICIPATES this. B's contribution is naming the specific hazard: a new bool field's zero value is `false` = gate off, which is fail-safe for spawning but means every one of those nine test sites silently gets gate-off semantics for the new pre-check. **Fold into A10's verification plan.** |

---

## 3. Contradiction — where the two passes disagree about the same fact

### The only direct contradiction: `0028:S31` and the applied sense.

- **Pass A, C-15:** S31's expectation is wrong. `Executor.Write`'s mismatch arm does NOT set `applied`.
- **Pass B, §1 Failure 3 and C-4:** quotes S31's "`Applied()` TRUE — the one arm where the applied sense is set on a failure" **as fact**, and builds its whole third failure narrative on it. B's premortem states outright: "`Applied()` came back true, no error surfaced."

**Pass A is right. Pass B is wrong on this fact.** `internal/accessor/executor.go`, the
mismatch arm, carries an explicit comment and then constructs the refusal without
`r.applied`:

```go
// `applied` stays UNSET: `0004:C14` scopes the applied-but-unverified
// sense to `read_back_incomplete` and a post-mutation `timeout`, and a
// mismatch is not unverified — its verification ran and found the
// artifact wrong (REQ-67).
if mismatch := verifyReadBack(planned, before, observed); mismatch {
    r := refusalOf(def, readTimeout, ClassReadBackMismatch, nil)
    r.Expected = planned
    r.Observed = observedTags(observedValues)
    return WriteResult{Refusal: r}
}
```

Every other refusal arm in the function that intends the applied sense sets
`r.applied = true` explicitly (the deadline arm, the `!hasReader` arm, both
read-back-incomplete arms, the read-back-timeout arm). The mismatch arm alone
does not, deliberately and with a cited contract reason.

**Consequence for the resolve half.** B did not merely miss A's finding — B
*inherited the record's false premise* and reasoned from it. B's C-4 is still a
real and valuable finding (see §2), but its symptom prose ("`Applied()` came back
true") must be corrected before it is used. The two findings compose: S31's
assertion is wrong about the flag AND the defence it tests is unenforced.

### Near-contradiction, resolved as complementary: `0028:C1` C1.3 `re-anchor:`.

A says the pass runs every rule and checks only cardinality; B says the pass does
not re-check siblings at all. **A's reading is the correct one** — the clause says
"every rule's anchor is run again over the POST-EDIT buffer". B's premise is wrong;
B's *scenario* (rewriting key A's line reshapes key B's target) is still a valid
instance of A's cardinality-not-identity hole. Not a contradiction about source, a
contradiction about what the clause says; resolve against A's reading and keep B's
scenario.

---

## 4. Grounding verdicts on the source-anchored claims

### 4.1 Pass A C-1/C-2 — the pre-write baseline read. **CONFIRMED.**

`internal/accessor/executor.go::Executor.Write` (line 262) runs, in order:

1. `e.selects(name, CapWrite)` → def, art, timeout
2. empty-plan short circuit; `WriteBinding` type assertion
3. `nonOwnedPlanKeys` check
4. `reader, hasReader := e.Registry.readerFor(def.Accessor.Role)`
5. `protected := protectedKeys(reader, hasReader, plannedKeys)`; **if `len(protected) != 0`: `raw := e.invokeRead(ctx, reader, art, rt, protected)`** — a full bounded read invocation through that same reader, **before any write**
6. `binding.Apply(applyCtx, art, slices.Clone(planned))`
7. deadline arm / error arm / read-back

Step 5 is exactly what A says it is, and the failure handling is verbatim what A
quotes:

```go
if raw.class != "" {
    // A timed-out, errored, or incomplete pre-write read
    // establishes NO baseline. ...
    baselineUnread = slices.Clone(protected)
} else { ... }
```

So when the reader is command-backed and the gate is off, the failure is **already
detected, already swallowed, and the write already proceeds** — surfacing later as
`read_back_incomplete` with `applied = true`. That is, word for word, the outcome
`0028:A2`'s "If wrong" says must be avoided.

**A2's `Verified` stamp omits a load-bearing fact.** Its Evidence quotes the
`readerFor` line and the `binding.Apply` line and describes the interval between
them as empty ("the order C1.3 requires, within one function"). It never reports
the intervening invocation. Its qualifier — "C1.3's pre-mutation refusal is a NEW
check at an EXISTING site" — is true only in the narrow sense that no gate check
exists there; it is materially misleading, because a *reader invocation* does, with
its own swallowed gate-off failure that C1.3 never names and never dispositions.

**And A's second-order point is confirmed too, and is the sharper one.**
`protectedKeys` (executor.go:526) returns the reader's `RequestedKeys()` **minus**
the planned keys. In the MVV's own fixture the reader declares `keys = ["status"]`
and the writer plans `status`, so `protected` is **empty** and step 5 never runs.
Whether the pre-existing swallowed path fires at all therefore depends on whether
the model author declared a second key on the reader — a variable no clause,
scenario, or MVV step in 0028 mentions. `0028:S25`'s "refuses BEFORE mutation" and
`0028:S23`'s `Invocations()` expectation are both specified against a model of the
executor with no baseline read in it.

C1.3 `read-back:` must say which detector is authoritative and what happens to the
other. Today it names one site as if it were the only one.

### 4.2 Pass A C-15 — S31 and the mismatch arm. **CONFIRMED.** See §3 for the code.

`0028:S31` asserts: "the read-back through the role's reader then refuses
`read_back_mismatch` with `Applied()` TRUE — the one arm where the applied sense is
set on a failure (A8)". The source sets no applied flag on that arm, by an explicit
`0004:C14` scoping comment. The record's only decoy-corruption test asserts a value
that cannot occur. A's predicted repair path — the assertion gets relaxed on first
run, deleting the only check on the applied sense for the corruption path — is the
realistic one.

### 4.3 Pass A C-18 — map iteration order. **PARTLY REFUTED. Downgrade.**

`0028:§performance-expectations` claims: "*Map order*: not reachable — rules resolve
against pre-edit line INDICES held per rule, so iteration order over the `edit.<key>`
tables cannot affect the output buffer (C1.3 'select')."

Checked against C1.3 `select:`, which states: "every rule's anchor is resolved
against the PRE-EDIT content and selections are held as pre-edit line INDICES (a
deletion never shifts a sibling rule's target)", plus `edit_anchor_collision` when
two rules select one line.

The claim holds for the **output buffer**, which is what the checklist row asserts:
selection is against pre-edit content (order-independent), each rule owns a distinct
pre-edit index (collision refuses otherwise), and a rewrite at index *i* is a pure
per-index substitution. Application order does not change the resulting bytes.
**A's specific mechanism — "the ORDER of application into the buffer determines the
intermediate state the re-anchor pass sees" — is blocked**: there is no
order-dependent intermediate state, because the re-anchor pass runs over the single
finished post-edit buffer, not over intermediates.

What survives, and it is worth one sentence: the checklist row scopes its claim to
"the output buffer" while the **refusal path** is not scoped at all. If two rules
carry independent load-time or apply-time defects, which one is reported first is
map-order-dependent unless something fixes it — and C1.4's `precedence:` fixes order
only *within one entry across categories*, not across sibling `edit.<key>` tables.
`0028:S28` ("One entry carrying TWO simultaneous load-time defects") is exactly this
test and its Expected line should pin a deterministic winner.

**Verdict: the determinism-checklist row as written is CORRECT for the output
buffer. A's escalation to flaky refusals is REFUTED for the reason A gives, and
survives only as a narrower point about refusal-reporting order. chart-as-net-new /
one-line clarification, not blocking.**

### 4.4 Both passes on `0028:A11` — the deadline arm. **CONFIRMED, verbatim.**

```go
applyCtx, cancelApply := context.WithTimeout(ctx, timeout)
err := binding.Apply(applyCtx, art, slices.Clone(planned))
appliedDeadline := errors.Is(applyCtx.Err(), context.DeadlineExceeded)
cancelApply()

if appliedDeadline {
    r := refusalOf(def, timeout, ClassTimeout, nil)
    r.applied = true
    ...
}
```

`appliedDeadline` is computed after `Apply` returns, **without consulting `err`**,
and unconditionally sets `applied = true`. A11's own Evidence states this correctly
and A11 is nonetheless `Pending`. The reachability question A11 defers is real: this
is latent for `path` today because `flowbind.Writer.Apply` ignores its context
entirely, but an in-process `edit` that refuses at the anchor stage still returns
through this arm, so a slow machine can mint `ClassTimeout` applied-but-unverified
over an edit that wrote no byte.

Meanwhile `0028:§pre-lock-mini-checks`'s `disposition` table closes with: "Every
refusal is decided BEFORE any byte is written (C1.3 `order:`), carries `Applied()`
false, and adds no new class". **There is no timeout row in that table.** The table
asserts as an invariant precisely what the `Pending` A11 says may not hold, and
C1.3 `order:` states the carve-out as settled ("That arm is 0004's, unchanged here
and out of this RDR's scope (A11)"). This is a `0028:G-assumptions` Status-consistency
violation in the record's most load-bearing clause. Both passes found it; both are
right.

### 4.5 Pass A C-13 — `flowbind.go::save`. **CONFIRMED. "Extend" is not viable.**

`internal/cli/flowbind/flowbind.go:139`:

```go
func save(path string, s store) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o700); err != nil { ... }
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".flow-artifact-*")
	...
	if err := clierr.WriteJSONLine(tmp, s); err != nil { ... }
	...
	if err := os.Chmod(staged, 0o600); err != nil { ... }
	return os.Rename(staged, path)
}
```

Three incompatibilities, each independent:

1. `MkdirAll(dir, 0o700)` — creation-adjacent, and `edit` fences creation out (C1.3 `select:` / `input:`: "a missing target is a stale model").
2. `clierr.WriteJSONLine(tmp, s)` — the signature takes a `store`, not bytes. `edit` writes a raw byte buffer.
3. `os.Chmod(staged, 0o600)` — a fixed mode. C1.3 `write:` and A5 require the ORIGINAL file's mode be preserved. The `fidelity` mini-check table says so explicitly: "unlike `flowbind.go::save`'s fixed 0600".

`save` shares exactly one thing with what `edit` needs: `CreateTemp` beside the
target + `Rename`. Everything else is JSON-store-specific. **The viable dispositions
are "sibling with the stage-and-rename discipline copied" or "extract a
`stageAndRename(dir, write func(io.Writer) error, mode fs.FileMode)` helper both
call". "Extend" is not one of them.** The audit row's "Extend or sibling" defers a
decision the row exists to record, and the record already knows the answer — the
`fidelity` table and A5 both state the mode divergence. This is a one-cell fix.

---

## 5. Merged ledger — the resolve half's worklist

| ID | 0028 passage | claim | raised by | verdict | severity |
|---|---|---|---|---|---|
| D-1 | `0028:A2`, `0028:C1` (C1.3 `read-back:` SITE), `0028:S23`, `0028:S25` | `Executor.Write` runs a pre-write baseline read through `readerFor(role)` (`protectedKeys` → `invokeRead`) BEFORE `Apply`; with a command-backed reader and the gate off it already fails, is already swallowed into `baselineUnread`, and the write already proceeds. A2's Verified evidence omits it and its "new check at an existing site" qualifier is materially misleading. Compounding: `protectedKeys` subtracts the planned keys, so in the MVV's own fixture `protected` is empty and the path does not fire at all — the gate-off behaviour depends on a reader-key-count variable no clause mentions. C1.3 must say which detector is authoritative and what happens to the other. | A | CONFIRMED | blocking-at-lock |
| D-2 | `0028:A10`, `0028:A11` vs `0028:C1` (C1.3 `read-back:` SITE, `order:`) and `0028:§pre-lock-mini-checks` (`disposition`) | Two Critical Assumptions are `Pending` while C1.3 states their conclusions as settled declarative prose inside a `normative` fence, and the `disposition` table asserts "Every refusal … carries `Applied()` false" with no timeout row. Source confirms the deadline arm computes `errors.Is(applyCtx.Err(), context.DeadlineExceeded)` after `Apply` without consulting `err` and sets `applied = true`. Direct `0028:G-assumptions` Status-consistency violation. Add: A10's field is a bool whose zero value means gate-off, and `accessor.Registry{}` is constructed at nine sites in `internal/cli/cmdbind/seam_0025_test.go` besides the one production site. | both | CONFIRMED | blocking-at-lock |
| D-3 | `0028:§testing-strategy` (S31) | S31 asserts `read_back_mismatch` with `Applied()` TRUE, "the one arm where the applied sense is set on a failure". The source's mismatch arm leaves `applied` UNSET by an explicit `0004:C14` scoping comment. The record's only decoy-corruption test asserts a value that cannot occur. **Pass B reasoned from the record's false premise and repeated it as fact — correct B's prose when folding D-9.** | A (B contradicts) | CONFIRMED | blocking-at-lock |
| D-4 | `0028:MVV` (fixture + steps 2–3), `0028:§problem-statement`, `0028:§capability-dependencies`, `0028:§prerequisites` | The MVV is not executable. (a) One `--tag nnnn=NNNN` addresses one record; the fixture declares three; three invocations put three changed rows in the shared README, falsifying step 3's "`git diff --stat` shows one line per file". (b) The C1.6 reader verb `rdr index --row-json` does not exist — `rdr index` is a real verb, `--row-json` returns zero hits in the projector — and "the consumer adds the verb" sits in a Spec Impact cell, not in Prerequisites. Phase 4 cannot run when Phase 3 merges. | A | CONFIRMED | blocking-at-lock |
| D-5 | `0028:C1` (C1.3 `write:`), `0028:§key-discoveries`, `0028:§consequences`, `0028:§failure-modes` | Atomicity is scoped to ONE entry; the MVV's journey spans TWO entries over TWO files. Nothing states what the user has when entry 1 applies and entry 2 refuses — a record reading `Final` and an index reading `Draft`, with the tool reporting "not applied". Key Discoveries markets the one-entry boundary as "stronger than Ansible … Puppet … Terraform" without naming the resulting hole. | A | CONFIRMED | fix-now |
| D-6 | `0028:C1` (C1.3 `re-anchor:`) | The re-anchor invariant says each rule "must select exactly its own rewritten line" but specifies only a CARDINALITY. After a sibling `clear = "line"` deletion, a sibling anchor can select exactly one line that is a DIFFERENT line and still pass. B's variant: a rewrite reshaping a sibling's target line. No scenario in S1–S31 pairs a `clear` with a sibling rule. The invariant needs an identity check (compare against the held pre-edit index, adjusted for deletions), not a count. | both (A's reading of the clause is the correct one) | CONFIRMED | fix-now |
| D-7 | `0028:C1` (C1.5 `one-way:`), `0028:C1` (C1.4 `edit_clear_invalid`) | `clear = "line"` is irreversible and C1.5's "never for a shared row whose other cells carry other state" is prose only — `edit_clear_invalid` checks only that the value is in `{"line"}`, and nothing at apply time checks it. The destructive case (a README row deleted whole) read-backs as SUCCESS, because the key genuinely IS absent. Recovery is by hand, which relocates the drift class the Problem Statement exists to close. | both | CONFIRMED | fix-now |
| D-8 | `0028:§prerequisites`, `0028:F6`, `0028:MVV` (fixture), `0028:G-scope` | The MVV fixture is curated to exclude the two live README rows (0001, 0006) the anchor cannot address — the one class of pre-existing drift the real corpus actually has. The acceptance test proves the composed scenario against data shaped to avoid the known failure. Remedy: include one as a NEGATIVE scenario proving the documented refusal, and have Scope Verification confirm the MVV covers the real corpus's shape distribution. | both | CONFIRMED | fix-now |
| D-9 | `0028:F4`, `0028:F5`, `0028:S31`, `0028:§risks-and-mitigations` | Read-back is named the "sole defence" against silent wrong-line corruption, but it depends on the role reader's selection logic being INDEPENDENT of the anchor's — an ambient property of a consumer-authored reader, imposed by no clause in C1 and disclaimed by C1.4. An author whose reader shares the anchor's blind spot gets corruption reported as success. The record implies this is mitigated; it is not. Minimum absorbable fix: state it as an accepted UNMITIGATED risk in `§risks-and-mitigations` rather than leaving F4/F5 to imply a working defence. Fold with D-3 — S31's applied-sense assertion is separately wrong. | B | CONFIRMED | fix-now |
| D-10 | `0028:C1` (C1.6), `0028:§metadata` (Profile `large`, "one contract"), `0028:§decision-rationale`, `0028:G-proportionality` | C1.6 is a second independently-scored load-bearing contract (it extends 0025:C2, changes the READ and GATE paths, lives in a different binding package) riding inside a one-contract Profile claim. The record writes the split signal down — "C1.6 is scored separately, because the matrix above is about the write carrier and C1.6 is not one" — and fences past it. Consequences: no read-path premortem, no read-path alternatives matrix, one scenario (S7), and C1.6 cannot be reverted or rolled out independently of C1.1–C1.5. Proportionality must either count it and split, or justify two independently-scored decisions under one contract id. | both | CONFIRMED | blocking-at-lock |
| D-11 | `0028:A3`, `0028:S8`, `0028:C1` (C1.4 "what lint does NOT prove") | A3 is `Verified` against a corpus frozen at 33 records / 26 rows; S8 concedes it was already 34/28 when the scenario was written "and this record adds one". Anchor uniqueness is a property of data asserted as a property of design, verified once, with no mechanism (lint mode, literal-anchor dialect, corpus guard) preventing regression — only the after-the-fact refusal. The corpus is documentation ABOUT the format the anchor matches, so every critique, premortem and template excerpt quoting `- **Status**:` is a decoy generator. A3's own If-wrong contemplates the right fix (a literal-anchor dialect) and it is not in v1. | both (A: temporal drift in THIS corpus; B: generality to a second adopter) | CONFIRMED (A's temporal axis); PLAUSIBLE (B's cross-consumer axis) | fix-now (state the mechanism gap honestly in `§risks-and-mitigations`; the literal-anchor dialect is chart-as-net-new) |
| D-12 | `0028:§existing-infrastructure-audit` (`save` row) | The Decision cell reads "Extend or sibling" — an unmade decision in the row whose purpose is to record it. `save` does `MkdirAll(0700)` (creation-adjacent; `edit` fences creation out), `clierr.WriteJSONLine(tmp, s)` (takes a `store`, not bytes), and `Chmod(0600)` (fixed mode; C1.3/A5 require the ORIGINAL mode, and the `fidelity` table says "unlike `flowbind.go::save`'s fixed 0600"). Only `CreateTemp`+`Rename` is shared. "Extend" is not viable; the row must say "sibling" or name an extracted `stageAndRename` helper. | A | CONFIRMED | fix-now |
| D-13 | `0028:C1` (C1.3 `write:` no-op arm), `0028:S20`, `0028:§pre-lock-mini-checks` (`disposition`) | The no-op arm is unobservable through the contract's own surfaces: witnessed only by inode identity (a filesystem property, and `os.Rename` behaviour is platform-variable), while `Invocations()` — which C1.3 `no subprocess:` says counts `Apply` calls "exactly as `flowbind.go::Writer` does" — increments either way. The `disposition` table calls it "silent by design". An agent-loop caller (this record's stated audience) cannot distinguish "already applied" from "applied now" from CLI output alone, which is exactly the idempotency reasoning such a caller needs. | both (A: contract-internal contradiction with `Invocations()`; B: caller-visibility) | CONFIRMED | fix-now |
| D-14 | `0028:C1` (C1.2 `anchor admits:`) | C1.2 closes `anchor`'s placeholder vocabulary on the fixed prefix `{tag.` scanned to its `}` and explicitly kills `{{`/`}}` as an escape there, but never names the expression an author uses for a line containing the literal text `{tag.x}`. RE2's own `\{tag\.` almost certainly works (the scan is for the literal prefix), and `edit_anchor_invalid`'s "malformed placeholder" predicate is undefined for the escaped form. The defect is the silence, not the absence of an escape. One sentence in C1.2 closes it. | A | CONFIRMED (as a gap; A's "unwritable" framing overstates) | fix-now |
| D-15 | `0028:C1` (C1.3 `target:`), `0028:S22` | "symlinks resolved" is a two-word obligation: no statement of final-component vs whole-path resolution, no TOCTOU discipline, no bound on where the derived staging directory may land. S22 asserts only that the symlink survives. A's escape-the-repo framing overreaches — `§consequences` already disposes of caller-bound indirection ("a caller who binds a hook file has bound a hook file") and `path` has the same property today — but the resolution semantics are genuinely unspecified. | A | PLAUSIBLE (spec gap CONFIRMED; security framing REFUTED) | chart-as-net-new |
| D-16 | `0028:§performance-expectations` (determinism checklist, *Map order*), `0028:S28` | A claimed the "Map order: not reachable" row is false because application order determines the intermediate buffer the re-anchor pass sees. **REFUTED for the OUTPUT BUFFER, which is what the row asserts:** C1.3 `select:` resolves every anchor against pre-edit content, holds distinct pre-edit indices per rule, and refuses `edit_anchor_collision` if two rules share a line — so per-index substitution is order-independent and the re-anchor pass runs over the single finished buffer, not over intermediates. What survives is narrower: the row scopes its claim to the output buffer while the REFUSAL-reporting path is unscoped, and C1.4 `precedence:` fixes order only within one entry across categories, not across sibling `edit.<key>` tables. S28's Expected line should pin a deterministic winner. | A | REFUTED (as stated); narrow residue PLAUSIBLE | chart-as-net-new |
| D-17 | `0028:A6` | B read A6 correctly — it is `Verified` on evidence that its stated failure mode "has no live target", since no accessor dump/marshal/normalize round-trip exists in `internal/` — but that is what a correctly verified negative looks like, and A6's positive conclusion (a third carrier needs one parallel `if a.Edit != nil {…}`, not a rewrite) is independently confirmed against `internal/table/load.go`'s additive per-field copy. Gating a future round-trip that does not exist is out-of-fence future-proofing. | B | REFUTED as a defect of A6 | chart-as-net-new |
| D-18 | `0028:§decision-rationale`, `0028:ALT3` | B: ALT3 is rejected on "document knowledge enters the tool" while the MVV, Testing Strategy and Illustrative Code encode one consumer's document family. **REFUTED:** ALT3 would put document-shape awareness in SHIPPED CODE (an adapter registry); the MVV puts it in TESTDATA. The distinction is coherent, and B does not engage it. The generality concern is real but is already carried by D-11's cross-consumer axis; it is not a defect of the ALT3 rejection. | B | REFUTED | chart-as-net-new |
