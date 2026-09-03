Model: claude-opus-5

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0028:C1` (C1.3 `read-back:`) | The gate pre-check is specified against a model of `Executor.Write` that the source contradicts. `Write` already runs a **pre-write baseline read** through `readerFor(role)` (`protectedKeys` → `invokeRead`) BEFORE `binding.Apply`. When the reader is command-backed and the gate is off, that baseline read ALREADY fails today, is ALREADY swallowed into `baselineUnread`, and the write ALREADY proceeds. C1.3's clause is not "a new check at an existing site" — it is a second, redundant check bolted in front of an existing unchecked failure that C1.3 never names. | The gate-off refusal fires (or does not) depending on which of two unsynchronised paths wins; if the implementer wires only C1.3's check, `protected` is non-empty in the consumer fixture and MVV step 5 passes for the wrong reason. If they wire neither, an edit applies and returns `read_back_incomplete` — exactly the outcome A2's If-wrong forbids. | §1, premortem, AT-1 |
| C-2 | `0028:A2` | A2 is stamped **Verified / Source Search** on evidence that reads only the two lines it wanted. It quotes `readerFor` resolving before `Apply` and stops. It never reports that between those two statements the executor performs a full bounded read invocation through that same reader, nor that `flowbind.Reader.Read` for a `path` reader and `cmdbind.Reader` for a command reader take different arms there. The assumption's evidence is a subset selected to confirm the assumption. | Implementation discovers a third invocation of the command reader nobody budgeted; `Invocations()` and gate-off semantics diverge from S23/S25's stated expectations. | §1, §3, AT-1 |
| C-3 | `0028:C1` (C1.3 `order:` / `write:`) | The atomicity claim "all rules of one entry rewrite ONE buffer and land in ONE write" is scoped to ONE ENTRY, but the MVV's user journey spans TWO entries (`write.record`, `write.readme`) over TWO files. Nothing in C1, C1.3, or the Consequences states what happens when entry 1 applies and entry 2 refuses. The record advertises "stronger than Ansible/Puppet/Terraform" atomicity in Research Findings while the delivered boundary is strictly narrower than the delivered user outcome. | The record's Status flips to `Final` and the README row stays `Draft`. The user's corpus is now internally inconsistent, and the tool reports a refusal for the second entry with no indication the first one landed. No MVV step, no scenario, and no failure mode covers this. | §1, §2, premortem, AT-2 |
| C-4 | `0028:MVV` | The MVV binds `--tag nnnn=NNNN` — ONE record identity — on ONE invocation, and its fixture contains THREE records. The scenario is written as if all three flip in one pipeline run, but a single `{tag.nnnn}` binding can address exactly one README row and one record. The MVV as written is not executable. | The implementer either runs the MVV three times (and step 3's "`git diff --stat` shows one line per file" fails on the README, which accumulates three changed rows) or discovers the fixture cannot be built as described and silently redefines the acceptance test. | §1, §5, AT-3 |
| C-5 | `0028:C1` (C1.2 `anchor admits:`) | Placeholder/regex ambiguity is under-specified in the direction that matters. C1.2 says in `anchor` "ONLY the fixed prefix `{tag.` opens a placeholder, scanned to its closing `}`". RE2 patterns legitimately contain `{tag.` followed by a `}` in contexts that are not placeholders (a literal-brace character class, an alternation containing the literal text `{tag.foo}` — e.g. an anchor over a file that itself documents this feature). There is no escape hatch in `anchor` (C1.2 explicitly kills `{{`/`}}` there). A literal `{tag.` is unwritable. | An author anchoring on a line that contains the literal text `{tag.` gets `edit_anchor_invalid` at lint with no way to express what they mean. The consumer's own `rdr-write.toml` documentation lines are exactly such content. | §1, AT-4 |
| C-6 | `0028:C1` (C1.3 `re-anchor:`) | The re-anchor invariant is stated as a total function over rules but is provably wrong for `clear = "line"` under multi-rule entries. C1.3 says re-anchor must select "exactly its own rewritten line (or, for a deleted line, zero lines)". After a deletion, a SIBLING rule's held pre-edit index no longer points at its own line, and the re-anchor pass runs anchors — not indices — over the post-edit buffer. A sibling whose anchor now matches a shifted line satisfies "selects exactly one" while pointing at the wrong content. The invariant checks cardinality, not identity. | A two-key entry with one `clear` silently rewrites the wrong line, passes re-anchor, and lands. Read-back catches it only if the reader reads that particular line. | §1, premortem, AT-5 |
| C-7 | `0028:A11` | A11 is `Pending` while C1.3 `order:` states the deadline hazard as settled prose ("That arm is 0004's, unchanged here and out of this RDR's scope (A11)") and the `disposition` mini-check table asserts "Every refusal … carries `Applied()` false" with no timeout row. The contract states as invariant precisely what the pending assumption says may not hold. This is the Finalization Gate's own "no assumption marked `Pending` may have settled-fact prose elsewhere depending on it" violation, in the record's most load-bearing clause. | On a loaded machine, a refused edit that wrote no byte reports `ClassTimeout` with `applied = true`. The agent loop treats an untouched file as possibly-mutated and stops the flow. Unreproducible, unlogged, blamed on the artifact. | §1, §3, premortem, AT-6 |
| C-8 | `0028:A10` | A10 is `Pending`, and C1.3 `read-back:` states its conclusion as a fixed SITE ("the gate from a field on `accessor.Registry`, set at `flowbind.go::Registry`"). A10's own If-wrong says the fallback moves the refusal minting to the CLI — a different error surface, different Detail composition, and a different phase (`flow_exec.go::phase`). The contract has pinned a site the record has not verified is available. | If A10 fails, C1.3's normative SITE sentence is wrong in a locked contract, and the Detail's shape (S27's assertion) changes with it. Post-lock the contract cannot be amended without a successor. | §1, §3, AT-7 |
| C-9 | `0028:C1` (C1.3 `write:` no-op arm) + `0028:S20` | The no-op arm ("a post-edit buffer equal to the input is not written at all") is unobservable through the contract's own surfaces and is witnessed only by inode identity — a filesystem property, not a contract property. Worse, it is in direct tension with `Invocations()` (C1.3 `no subprocess:` says it counts `Apply` calls "exactly as `flowbind.go::Writer` does"), which increments unconditionally. The disposition table calls the no-op "silent by design", and read-back then reports success on a file the writer declined to touch. | An idempotent re-run reports the write verified when no write occurred. Indistinguishable from a write that was silently dropped by a bug. The only test (S20) asserts an inode, which is a proxy nobody will maintain on a platform where `os.Rename` behaviour differs. | §1, §2, AT-8 |
| C-10 | `0028:§problem-statement` | The problem statement's own scope disclaimer concedes that the consumer must still "add the row-addressed projector verb C1.6 reads with." That verb does not exist. The Illustrative Code marks it "verb illustrative; the consumer's to add." The MVV — this record's acceptance gate — cannot run until a DIFFERENT repo ships a DIFFERENT feature, and no clause, prerequisite, or capability-dependency row makes that a blocking prerequisite. Prerequisites lists only artifact normalization. | The implementation completes, Phase 4's "MVV fixture end to end" cannot be executed, and the record is closed on unit tests while the composed capability — the entire stated reason this RDR exists — remains unproven. This is the exact failure the problem statement blames 0023/0024/0025 for. | §1, §2, §4, AT-9 |
| C-11 | `0028:C1` (C1.6) | C1.6 extends 0025:C2's closed argv vocabulary to a whole FAMILY (`{tag.<key>}`) inside a record whose Proportionality fence claims one seam. C1.6 changes the READ and GATE paths — a different capability, a different binding package (`cmdbind`), a different contract (0025:C2) — while the Profile field and the Proportionality note both assert "one seam — the `edit` write carrier". The Decision Rationale concedes it: "C1.6 is scored separately, because the matrix above is about the write carrier and C1.6 is not one." That sentence is a split signal written down and then ignored. | The read-path change ships under a write-carrier record with no read-path premortem, no read-path alternatives matrix, and one scenario (S7). When `{tag.<key>}` needs to grow (defaults, absence semantics, gate argv), the amendment has no home. | §1, §2, AT-10 |
| C-12 | `0028:A3` | A3 is stamped Verified against 33 records and a 26-row index, then S8 concedes the corpus was already 34 and 28 by the time the scenario was written and "this record adds one." The verification is a snapshot of a monotonically growing corpus whose growth is authored by the very users of this feature. "No decoy scores ≥2" is a property of today's data, asserted as a property of the design. | The first RDR that quotes a `- **Status**:` line in its own prose — a critique file, a premortem, a template excerpt — makes that record unflippable with `edit_anchor_ambiguous`, and the fix is to edit prose the user wrote for a different reason. | §3, premortem, AT-11 |
| C-13 | `0028:§existing-infrastructure-audit` (`save` row) | The audit row for `flowbind.go::save` says "Extend or sibling" — an unmade decision in the row whose whole purpose is to record the decision. `save` today does `MkdirAll(dir, 0o700)`, encodes JSON via `clierr.WriteJSONLine`, and chmods 0600. None of those three behaviours is compatible with `edit`. The record leaves the implementer to discover that "extend" is impossible. | Implementer either forks `save` (two atomic-write disciplines drift apart on the next fix) or parameterises it into a shape neither caller reads well. Cost lands as a review argument, not a design decision. | §2, AT-12 |
| C-14 | `0028:C1` (C1.5 `one-way:`) | `clear = "line"` deletes a line permanently, and C1.5's mitigation is a warning in prose ("declare `clear` only where the LINE is the key"). Nothing at lint or apply time enforces it. C1.4's `edit_clear_invalid` checks only that the value is in `{"line"}`. The record concedes in Consequences that recovery is "by hand". | A model author declares `clear = "line"` on a README row rule (the case C1.5 names as the wrong one), a `<clear>` plan lands, and the whole index row — three cells of unrelated state — is gone. Read-back reports the key absent: SUCCESS. The user finds out from git. | §1, §4, AT-13 |
| C-15 | `0028:§testing-strategy` (S31) | S31 asserts the decoy case produces `read_back_mismatch` with `Applied()` TRUE. `Executor.Write`'s mismatch arm explicitly does NOT set `applied` ("`applied` stays UNSET: `0004:C14` scopes the applied-but-unverified sense to `read_back_incomplete` and a post-mutation `timeout`"). The record's own most-cited safety net is specified against behaviour the source contradicts. | The one test standing between the user and silent wrong-line corruption asserts a flag value that will never be true. It fails on first run and gets "fixed" by relaxing the assertion — deleting the only check on the applied sense for the corruption path. | §1, §4, AT-14 |
| C-16 | `0028:F6` / `0028:§prerequisites` | The record discovers two live README rows (0001, 0006) it cannot address, files it under Failure Modes as "pre-existing data drift", and makes normalizing them a non-blocking prerequisite. The MVV fixture is then constructed to exclude them. The acceptance test was shaped around the data that fails. | Day one on the real corpus: two records cannot be flipped. The user's mental model is "the tool works"; reality is "the tool works on 32 of 34". | §3, AT-15 |
| C-17 | `0028:C1` (C1.3 `target:`) | Symlink resolution is stated as a one-word obligation ("symlinks resolved") with no TOCTOU discipline, no statement of whether resolution is of the final component or the whole path, and no statement of what happens when resolution escapes an expected root. The staging directory is derived from the resolved path, which means an `edit` write can stage and rename into a directory the caller never named. | A caller binds `--artifact record=./link.md` where `link.md` points outside the repo. The tool writes outside the repo, atomically, having stated in Consequences that "a caller who binds a hook file has bound a hook file" — a disclaimer that does not cover the indirection. | §1, AT-16 |
| C-18 | `0028:§performance-expectations` | The determinism checklist claims "Map order: not reachable — rules resolve against pre-edit line INDICES held per rule, so iteration order over the `edit.<key>` tables cannot affect the output buffer." This is false when two rules collide-adjacent or when a `clear` deletes a line another rule's index follows: the ORDER of application into the buffer determines the intermediate state the re-anchor pass sees, and C1.3 fixes no iteration order over the `edit.<key>` map. Go map iteration is randomised. | Flaky refusals: the same model against the same file refuses `edit_anchor_unstable` on some runs and applies on others. Diagnosed as a race, is a map. | §1, premortem, AT-17 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The gate pre-check is designed against an executor that does not exist

**Root cause.** A2 and C1.3 `read-back:` both model `Executor.Write` as: resolve reader → call `Apply` → read back. The actual function, read in full at `internal/accessor/executor.go`, does this:

1. `selects` → definition, artifact, timeout
2. non-owned plan key check
3. `reader, hasReader := e.Registry.readerFor(def.Accessor.Role)`
4. `protected := protectedKeys(reader, hasReader, plannedKeys)`
5. **if `len(protected) != 0`: `raw := e.invokeRead(ctx, reader, art, rt, protected)`** — a full bounded invocation of the role's reader, **before any write**
6. `binding.Apply(...)`
7. read-back

Step 5 is a pre-write baseline read through the same reader C1.3's clause is about. When that reader is command-backed and `--allow-commands` is off, `cmdbind`'s spawn ladder refuses at its first arm (`if !cfg.AllowCommands`), `invokeRead` returns `class: ClassExecutionFailure`, and the executor's handling is:

```go
if raw.class != "" {
    // A timed-out, errored, or incomplete pre-write read
    // establishes NO baseline. Discarding it leaves `before`
    // empty, which makes `0004:C12`'s protected-tag clause
    // vacuously true (REQ-53, REQ-80).
    baselineUnread = slices.Clone(protected)
}
```

The gate-off condition is **already detected, already discarded, and the write already proceeds**. C1.3's pre-mutation refusal is therefore not "a NEW check at an EXISTING site" as A2's qualifier claims — it is a second detector in front of a first detector whose result the executor deliberately throws away. The record does not know this path exists.

Note the further consequence C1.3 has no clause for: `protectedKeys` returns the reader's declared keys **minus the planned keys**. In the MVV's own fixture the reader declares `keys = ["status"]` and the writer plans `status`, so `protected` is **empty** and step 5 does not run. The gate-off failure therefore surfaces or does not surface depending on whether the model author declared a second key on the reader — a variable the RDR never mentions.

**Enabling passage.** `0028:A2`'s Evidence, which quotes exactly the two lines that support the conclusion and reports the intervening thirty as absent. Its stated Method is `Source Search`; a source search that reads a function's two endpoints and not its middle is not a source search. And `0028:C1`'s C1.3 `read-back:` clause, which asserts the SITE with the confidence of verified fact.

**Symptom.** In the fixture the RDR itself ships, MVV step 5 passes because the C1.3 check fires. In the consumer's real model — where a reader legitimately declares more keys than one writer plans — the pre-write baseline read fires first, is swallowed, the edit applies, and the user gets `read_back_incomplete` on an applied write for a write that ran no command. Which is precisely, word for word, the outcome A2's "If wrong" says must be avoided.

### 1.2 The user journey spans two entries and atomicity spans one

**Root cause.** C1.3 `write:` says: "all rules of one entry rewrite ONE buffer and land in ONE write". Research Findings then escalates that into a competitive claim: "per-entry atomicity (C1.3: all rules land in one write or none) is stronger than Ansible … Puppet … and Terraform." The MVV's journey is two entries over two files: `write.record` (the record's Status bullet) and `write.readme` (the index row). Between them there is no transaction, no ordering statement, no compensation, and no failure mode.

The RDR knows the two-file scenario is the point — the Decision Rationale explicitly rejects the narrower "ship `edit` alone and leave the README half unverified" option "because it would make the MVV prove half the stated outcome." So the record commits to the two-file outcome and then defines its atomicity boundary at one file, and never says so in the same place twice.

**Enabling passage.** `0028:C1`, C1.3 `write:`, and the Research Findings bullet that markets the one-entry boundary as a strength ("The boundary is one entry rather than a whole run, which is what makes it achievable"). Achievability is argued; the resulting user-visible hole is never named.

**Symptom.** The user runs a lock. The record file flips to `Final`. The README row's anchor fails — a duplicated row, a drifted shape, one of the two unlinked rows F6 already concedes exist — and refuses `edit_anchor_ambiguous` with `Applied()` false. The tool has reported "not applied". Half the state changed. The user's index and their records now disagree, and the refusal message says nothing happened.

### 1.3 The MVV cannot be executed as written

**Root cause.** Two independent blocks.

First, `{tag.nnnn}` is one binding per invocation. The MVV fixture contains three records ("a record whose Status is `Draft`, one whose Status is `Draft [joint decision → …]` … and one whose bracketed qualifier WRAPS"), and step 2 runs one `resolve | set-state` pipeline with one `--tag nnnn=NNNN`. Step 3 then asserts the end state of all three records AND that "`git diff --stat` shows one line per file". Three records flipped through three invocations produce three changed lines in the one README. The step is self-contradictory.

Second, the C1.6 reader the MVV depends on — `["rdr", "index", "--row-json", "{tag.nnnn}", "{artifact}"]` — does not exist. The Illustrative Code says "verb illustrative; the consumer's to add." Capability Dependencies lists "Row-addressed command reader over a shared artifact" as **Introduced / This RDR (C1.6)** and then adds "the consumer adds the verb" in the Spec Impact column, where it is invisible as a blocker. Prerequisites — the section whose job is exactly this — lists only artifact normalization.

**Enabling passage.** `0028:MVV`, whose fixture paragraph and step 3 were written at different times and never reconciled; and `0028:§problem-statement`'s scope disclaimer, which converts "the acceptance test depends on unshipped work in another repo" into a paragraph about honesty rather than a prerequisite.

**Symptom.** Phase 4 says "The MVV fixture end to end". It cannot run. The implementer ships against unit tests S1–S31, all of which pass, and the composed capability the record exists to deliver goes unproven — the identical failure the Problem Statement diagnoses in 0023/0024/0025 ("each landed what their kata asked and the composed capability still did not exist because it was never a ticket; this RDR is that ticket").

---

## 2. The section rewritten within 6 weeks of shipping

**`0028:C1` — C1.3, the execution-semantics clause.**

Not C1.2 (admission is genuinely well-reasoned and the RE2/quoting asymmetry is the record's strongest work). Not C1.4 (lint categories are mechanical). C1.3, because it is where every unverified assumption, every unmade decision, and every under-scoped invariant has been packed into one prose block that reads as settled.

Specifically, within six weeks:

- The `read-back:` clause's SITE sentence will be wrong, because A10 is `Pending` and A2's evidence is incomplete (C-1, C-2, C-8). Whichever way that resolves, the sentence naming `accessor.Registry` and `flowbind.go::Registry` changes.
- The `order:` clause's timeout carve-out will be revisited, because A11 is `Pending` and the carve-out is the difference between "a refused edit is never reported applied" (what the disposition table asserts) and "a refused edit is sometimes reported applied" (what the source does) (C-7).
- The `write:` clause's no-op arm will be deleted or made observable, because "witnessed by an unchanged inode" is not a maintainable assertion and it contradicts `Invocations()` in the same clause (C-9).
- The `re-anchor:` clause will gain an identity check, because cardinality alone does not distinguish "my line" from "a line" after a sibling deletion (C-6).
- The `write:` clause will gain a multi-entry statement, because the first user to hit a half-applied lock will file it (C-3).
- The `target:` clause's two-word symlink obligation will grow into a paragraph the first time someone binds a link (C-17).

Six separate amendments to one clause. The Proportionality note defends the single fence on the grounds that "they are one design" — but six sub-clauses that each need independent amendment for independent reasons are six decisions wearing one label. The fence is not proportionality; it is a Profile-field optimisation.

The tell is already in the record: C1.3 is 5.6K of the 14.2K contract, it is the only clause whose text argues with itself ("The scoping is exact, not defensive:" followed by a paragraph explaining a case the clause does not cover), and it is the only clause that cites two `Pending` assumptions as settled.

---

## 3. The assumption that will not survive first contact with a real user

**`0028:A3`** — "each anchor selects exactly one line in every existing consumer record and README, template comments and quoted examples included."

A3 is Verified. The spike is real, the method is sound, and the conclusion is true — of a corpus frozen on 2026-08-31 at 33 records and 26 rows. S8 then admits, in the record itself, that the corpus was 34 and 28 by the time the test scenario was written, "and this record adds one."

The assumption is that a regex over prose written by humans, in a corpus that grows by the actions of the very users of this feature, will continue to select exactly one line. Nothing enforces it. It is not a lint check — C1.4 explicitly says "What lint does NOT prove: that an anchor matches exactly one line of a particular file." It is not a runtime guard that helps — the runtime guard's answer is refusal. It is a property of data, asserted as a property of design, verified once.

First contact is not exotic. It is this file. This critique quotes `- **Status**:` — and the moment such a line lands anywhere inside a record file that the Status anchor sweeps, that record refuses `edit_anchor_ambiguous`. Every premortem, every critique, every gate response, every RDR that documents the RDR flow's own metadata format is a decoy generator. The consumer's corpus is a corpus *about* the format the anchor matches. That is the single worst possible domain for a "no decoys exist today" verification.

Worse, the record has already seen the shape of the failure and mislabelled it. F6 reports two live README rows in a form no anchor covers, and files them as "pre-existing data drift, not a defect of this carrier … The remedy is to normalize the data, never to loosen the anchor." That reasoning is correct in isolation and wrong as policy: it establishes that the tool's answer to a corpus that has moved is "the user edits the corpus." The user's corpus will keep moving, because moving it is the user's job.

A3's If-wrong contemplates "the anchor grammar needs a second dialect or a literal-anchor mode." That is the right instinct and it is not in v1. What ships is a design whose correctness is a snapshot.

---

## 4. Premortem (written as if the failure has already happened)

Six weeks after `edit` shipped, the RDR flow's maintainer stopped using it and went back to the `sed` strings in `rdr-write.toml`. Here is how.

**Week 1 — the verb that was never scheduled.** Implementation closed Phases 1–3 clean. All thirty-one scenarios passed. Phase 4 could not start: `rdr index --row-json` did not exist, because it was "the consumer's to add" — a phrase that appeared in the Illustrative Code caption and in a Spec Impact cell, and in neither Prerequisites nor the tracker. The MVV was recorded as deferred to the consumer's kata rdr#yjye. `edit` shipped with its acceptance test unrun. Nobody objected, because thirty-one unit tests is a lot of tests.

**Week 2 — the half-lock.** The maintainer locked RDR 0031. `Executor.Write` ran for `write.record`: `protectedKeys` returned empty (the reader declares `status`, the writer plans `status`), so no baseline read; `Apply` rewrote the Status bullet; read-back through `rdr status -json` reported `Final`; success. Then `Executor.Write` ran for `write.readme`. The index row's anchor matched twice — 0031's slug had been renamed and the old row was still present three lines up, unnoticed. `edit_anchor_ambiguous`, `Applied()` false, exit 3. The message said the write did not occur. The record said `Final`; the index said `Draft`. C1.3's `write:` clause guarantees one entry lands atomically and says nothing about two, so nothing in the tool was wrong. The maintainer fixed it by hand and wrote it off.

**Week 3 — the applied-unverified ghost.** On a loaded laptop during a full-corpus sweep, `Executor.Write` called `binding.Apply` for an `edit` whose anchor matched zero lines. The binding refused instantly, without touching a byte. But the sweep had the machine at load, `applyCtx`'s 5s deadline had elapsed inside the refusal path, and the executor's next three lines are:

```go
appliedDeadline := errors.Is(applyCtx.Err(), context.DeadlineExceeded)
...
if appliedDeadline {
    r := refusalOf(def, timeout, ClassTimeout, nil)
    r.applied = true
```

The CLI reported `flow-write-timeout`, applied-but-unverified, on a file that had not been opened for writing. The agent loop driving the sweep halted the whole batch, as it is supposed to when a write may have half-landed. The maintainer spent an afternoon diffing the corpus to establish nothing had happened. A11 had said this could happen. A11 was `Pending`. C1.3 `order:` had said the arm was "0004's, unchanged here and out of this RDR's scope", and the `disposition` mini-check table had said "Every refusal is decided BEFORE any byte is written, carries `Applied()` false" — a table with no timeout row.

**Week 4 — the gate that fired twice, then not at all.** The maintainer added a second key to the record reader (`status`, `profile`) so a gate could see the profile. Next lock, run without `--allow-commands` by habit. `protectedKeys` now returned `["profile"]`, non-empty, so `Executor.Write` ran the pre-write baseline read through the command reader — which refused at `cmdbind`'s gate arm — and the executor discarded the failure into `baselineUnread` and continued to `Apply`. C1.3's pre-mutation gate check was in the `edit` binding, downstream. Which detector wins depended on where the implementer had put the check; they had put it where C1.3's SITE sentence told them to, at `Executor.Write`'s reader-resolution point, which is *after* the baseline read. The refusal fired correctly. But it fired with a Detail composed for the gate, while the baseline read had already burned a spawn attempt that `Invocations()` does not count and S23 does not model. Nobody noticed for another week.

**Week 5 — `clear` ate a row.** A record was withdrawn. The consumer's model had `clear = "line"` declared on the README row rule — C1.5's prose says "declare `clear` only where the LINE is the key … never for a shared row", but nothing in `edit_clear_invalid` checks it, and the author had copied the record rule. `<clear>` planned, anchor matched one line, line deleted, terminator included. Read-back through the row reader reported the key ABSENT. **Success.** The index row — id, title, date, status, four cells of unrelated state — was gone. C1.5's `one-way:` clause meant `edit` could not put it back. The maintainer restored it from git and deleted `clear` from the model.

**Week 6 — the decoy, and the assertion that had been relaxed.** A new RDR quoted `- **Status**: Draft` in its Problem Statement, discussing this very feature. Its own Status bullet was `Draft`. The anchor matched twice: `edit_anchor_ambiguous`, refused, correct. The maintainer edited their own prose to un-break the tool. Then a different record's real Status bullet drifted into a form the anchor missed while a template comment above it matched — exactly one match, the wrong line. The edit applied to the comment. Read-back through the projector read the real bullet, still `Draft`, and refused `read_back_mismatch`. S31 had been written to assert `Applied()` TRUE on this path; on first run it failed, because `Executor.Write`'s mismatch arm sets no applied flag by contract ("`applied` stays UNSET: `0004:C14` scopes the applied-but-unverified sense to `read_back_incomplete` and a post-mutation `timeout`"). The assertion had been relaxed to "refuses `read_back_mismatch`" during implementation. So the CLI reported a mismatch with no applied sense, the maintainer read "the artifact is wrong" as "the plan was wrong", re-ran with a corrected plan, and the corrupted template comment stayed corrupted in git for eleven days.

**The autopsy.** No single defect was fatal. What was fatal was that C1.3 read as settled on six points where the record held two `Pending` assumptions, one incomplete source search, one unspecified multi-entry boundary, one unobservable no-op, and one cardinality-not-identity invariant — and the only end-to-end proof that would have surfaced any of them was the MVV, which could not run because a verb in another repo was never a prerequisite.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 (C-1, C-2) — the executor's actual call sequence.**
```
Given the reviewer has RDR 0028's A2 and C1.3 `read-back:` clause
When they enumerate EVERY invocation of `readerFor(role)`'s definition inside
     `internal/accessor/executor.go::Executor.Write`, from `selects` to return
Then the enumeration must list the pre-write baseline read
     (`protectedKeys` -> `invokeRead`) as an invocation that occurs BEFORE `Apply`
And A2's Evidence must state what that read does when the reader is
     command-backed and the gate is off
And C1.3 `read-back:` must state which of the two detectors is authoritative
     and what happens to the other
FAILS TODAY: A2's evidence quotes two lines and omits the read; C1.3 names
     one site as if it were the only one.
```

**AT-2 (C-3) — multi-entry atomicity.**
```
Given the MVV's journey writes two entries over two files
When the reviewer asks "entry 1 applies, entry 2 refuses — what does the user have?"
Then C1, Failure Modes, or Consequences must answer in one sentence
And the answer must be reachable from the clause that claims atomicity
FAILS TODAY: C1.3 `write:` scopes atomicity to one entry; no section names
     the two-entry case; Research Findings markets the boundary as a strength.
```

**AT-3 (C-4, C-10) — MVV executability.**
```
Given the MVV fixture declares three records and one `--tag nnnn=NNNN`
When the reviewer traces step 2's single pipeline against step 3's assertions
Then every record in the fixture must be reachable by the bindings step 2 supplies
And `git diff --stat shows one line per file` must be true for the README
And every command in every fixture accessor must resolve to a shipped binary+verb
     or appear in Prerequisites as blocking
FAILS TODAY: one tag cannot address three records; `rdr index --row-json`
     does not exist and is not a prerequisite.
```

**AT-4 (C-5) — literal `{tag.` in an anchor.**
```
Given C1.2 states `{tag.` is the only placeholder opener in `anchor` and that
     `{{`/`}}` is NOT an escape there
When an author needs an anchor matching a line containing the literal text `{tag.x}`
Then C1.2 must name the expression that does it, or state that it is unexpressible
FAILS TODAY: neither. The vocabulary is closed with no escape and no exclusion.
```

**AT-5 (C-6) — re-anchor identity vs cardinality.**
```
Given an entry with two rules, one `clear = "line"` deleting line N
And the sibling rule's anchor, run over the post-edit buffer, matches exactly
     one line — a DIFFERENT line than the one it selected pre-edit
When the reviewer applies C1.3 `re-anchor:` as written
Then the pass must FAIL (identity), not PASS (cardinality)
FAILS TODAY: the clause says "must select exactly its own rewritten line" but
     specifies only a count; no scenario in S1-S31 pairs a `clear` with a sibling.
```

**AT-6 (C-7) — Pending assumption vs settled prose.**
```
Given A11 is Status: Pending
When the reviewer greps C1 and the `disposition` mini-check table for claims
     that depend on the deadline arm's behaviour
Then no clause may state the dependent conclusion as invariant
FAILS TODAY: C1.3 `order:` states the scoping as settled; the disposition table
     asserts "Every refusal ... carries `Applied()` false" with no timeout row.
     The Finalization Gate's own Assumption Verification item forbids exactly this.
```

**AT-7 (C-8) — Pending assumption pinned as a contract SITE.**
```
Given A10 is Status: Pending and its If-wrong relocates refusal minting to the CLI
When the reviewer reads C1.3 `read-back:`'s `SITE:` sentence
Then either A10 must be Verified, or the SITE sentence must be non-normative
FAILS TODAY: the sentence names `accessor.Registry` and `flowbind.go::Registry`
     as fixed, inside a `normative` fence, on a Pending assumption.
```

**AT-8 (C-9) — the no-op arm's observability.**
```
Given C1.3 `write:` states a post-edit-equal buffer is not written at all
When the reviewer asks what contract surface distinguishes it from a written no-op
Then the answer must not be a filesystem property
And it must be reconciled with `Invocations()`, which C1.3 says counts Apply calls
FAILS TODAY: S20's only witness is inode identity; `Invocations()` increments
     either way; the disposition table calls the arm "silent by design".
```

**AT-9 (C-10) — external dependency as prerequisite.**
```
Given Capability Dependencies marks the row-addressed reader "Introduced/This RDR"
     with "the consumer adds the verb" in Spec Impact
When the reviewer asks "can Phase 4 run on the day Phase 3 merges?"
Then the answer must be yes, or the missing verb must be a checked Prerequisite
FAILS TODAY: the answer is no and Prerequisites lists only README normalization.
```

**AT-10 (C-11) — one seam or two.**
```
Given the Profile field and the Proportionality note both assert one seam
When the reviewer reads Decision Rationale's "C1.6 is scored separately, because
     the matrix above is about the write carrier and C1.6 is not one"
Then that sentence must be reconciled with the one-seam claim, or C1.6 splits
FAILS TODAY: the record writes down the split signal and then fences past it.
     C1.6 changes read and gate paths, a different package, a different contract.
```

**AT-11 (C-12) — corpus-frozen verification.**
```
Given A3 verifies anchor uniqueness over a 33-record corpus
When the reviewer asks what prevents record 34 from introducing a decoy
Then the answer must be a mechanism (a lint mode, a literal-anchor dialect,
     a corpus guard), not the refusal that fires after the fact
FAILS TODAY: C1.4 explicitly disclaims it; S8 concedes the corpus already grew.
     The corpus is documentation OF the format the anchor matches.
```

**AT-12 (C-13) — the unmade `save` decision.**
```
Given the Existing Infrastructure Audit's Decision cell reads "Extend or sibling"
When the reviewer checks `flowbind.go::save` for compatibility
Then MkdirAll(0700), WriteJSONLine, and Chmod(0600) must each be dispositioned
FAILS TODAY: none are; the row defers a decision the row exists to record.
```

**AT-13 (C-14) — `clear` on a shared row.**
```
Given C1.5 warns in prose "never for a shared row whose other cells carry
     other state (the README row: deletion would drop the whole row)"
When the reviewer asks what enforces it
Then a lint category or an apply-time refusal must be named
FAILS TODAY: `edit_clear_invalid` checks only the value is in {"line"};
     the destructive case read-backs as SUCCESS.
```

**AT-14 (C-15) — S31's applied-sense assertion.**
```
Given S31 expects `read_back_mismatch` with `Applied()` TRUE
When the reviewer reads `Executor.Write`'s mismatch arm
Then the expectation must match the source
FAILS TODAY: the source sets no applied flag on mismatch, by an explicit
     0004:C14 scoping comment. The record's only decoy-corruption test asserts
     a value that cannot occur.
```

**AT-15 (C-16) — fixture shaped around failing data.**
```
Given F6 names two live README rows the anchor cannot address
And the MVV fixture declares all rows in the linked form
When the reviewer asks what fraction of the real corpus the MVV covers
Then the answer must be stated, and partial adoption must be a shipping decision
FAILS TODAY: it is a non-blocking prerequisite bullet.
```

**AT-16 (C-17) — symlink resolution discipline.**
```
Given C1.3 `target:` says "symlinks resolved" in two words
When the reviewer asks: final component or whole path? resolved when relative to
     the write? what bounds the staging directory the resolved path implies?
Then C1.3 must answer all three
FAILS TODAY: none are answered; S22 asserts only that the symlink survives.
```

**AT-17 (C-18) — map iteration order.**
```
Given Performance Expectations claims map order is "not reachable"
When the reviewer constructs an entry with a `clear` rule and a sibling rule
Then the order of buffer application must not change the re-anchor pass's verdict
FAILS TODAY: pre-edit INDICES fix WHERE each rule writes, not the ORDER of
     application, and the re-anchor pass runs over the resulting buffer.
     C1.3 fixes no iteration order over the `edit.<key>` tables; Go randomises it.
```
