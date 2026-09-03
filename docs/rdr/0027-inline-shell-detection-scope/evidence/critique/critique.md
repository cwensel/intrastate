Model: claude-opus-5

# Critique — RDR 0027 (premortem lens, time-shifted failures)

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0027:A2` | The "zero new refusals" verification is scoped to 32 in-repo fixture vectors and 0025's own models — a population with 11 green `command` bindings, none of which contains a listed interpreter basename in a non-argv0 position. It is structurally incapable of observing the class it is asked to bound, so its Verified stamp certifies nothing about the predicate's blast radius. A five-minute sweep of ordinary tool invocations produces six false refusals the spike could never have seen: `["grep","-r","python","-c","src/"]` (grep's `-c` = count), `["pip","install","python","-c","constraints.txt"]` (pip's `-c` = constraints file), `["rg","python","-c"]`, `["tar","-czf","out.tgz","python","-c"]`, `["npm","run","node","--","-e"]`, `["my-tool","--lang","php","--mode","-r"]`. All are `old=false → new=true`. | A model that lints green on today's binary refuses on upgrade with `command_shell_interpreter` and the remediation "put it in a script and declare the script as argv0" — for a `grep` invocation containing no shell and no inline code. | §1, §3, premortem, AT-1 |
| C-2 | `0027:C1` (`predicate:` line) | The predicate refuses argv in which **no interpreter is executed at all**. `base(argv[i])` is matched against a name list without regard to whether argv[i] is in command position; `["echo","node","--eval"]`, `["cat","php","-r"]`, `["mytool","--interpreter","python","--flag","-c"]` all refuse. The record calls this "the already-accepted false-refusal class, unchanged" (`0027:D-selection-predicate`, Consequences). That is false: today the interpreter must be argv0 (after an `env` walk), so a *real* interpreter is being invoked and the collision is confined to its own operand list. Freeing position converts a bounded collision inside one command's flags into an unbounded collision across every word of every argv. The class is not unchanged; it is a different class with the same name. | An `intrastate lint` refusal claiming a binding "declares the interpreter form `php -r`" for `["cat","php","-r"]`, which spawns `cat`. | §1, §2, premortem, AT-2 |
| C-3 | `0027:A6` | Status `Pending`, and the record proposes to lock with it open. It is not merely unverified — the surface it presupposes probably cannot exist under the governing contract. `table.Categories()` has **zero non-test consumers** (`grep -rn "table.Categories()" internal/cli/ cmd/` returns nothing outside `_test.go`); `docs/cli-output-contract.md` defines output entirely as the `CLIError` envelope with `findings[]`, i.e. refusal-shaped. There is no success-path text channel for a category description, and inventing one is a change to the output contract — a contract this RDR does not own and does not cite as a dependency. | Phase 3 lands as a Go doc comment nobody reads, or as an unreviewed new success-payload key; either way the "promise" the whole record is built on is never delivered to a reviewer. | §1, §2, §3, premortem, AT-3 |
| C-4 | `0027:C1` (`promise:` clause) | C1 normatively asserts "the lint's user-facing description states the out-of-scope forms in the words above" as a present-tense contract while conceding in the same clause that the description does not exist and depends on `A6`, which is `Pending`. This is settled-fact prose depending on an unverified assumption — precisely what `0027:G-assumptions` forbids ("no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it"). The record cannot pass its own gate. | Nothing at first. Then a reviewer cites `0027:C1` `promise:` in a later RDR, discovers it names a surface that never shipped, and the citation resolves to nothing. | §1, §2, premortem, AT-3 |
| C-5 | `0027:A5` | A5's mechanism is empirically wrong for one of three capabilities. It claims the charted `stdin = "none"\|"envelope"` successor closes the stdin axis by "withholding the envelope from a command that never declared it." But `internal/cli/cmdbind/cmdbind.go` sets `cmd.Stdin = bytes.NewReader(stdin)` unconditionally, and a **write** binding sends `stdinObject(planned)` — a JSON object whose *values* are model-authored tag values. A `write` accessor declared `["sh","-s"]` receives author-controlled bytes on stdin. Withholding is not the mechanism for write; the envelope is the payload. Empirically `printf '{}' \| sh -s` exits 127 (`{}: command not found`) — the read/gate case is inert by accident of JSON syntax, not by design — but a planned tag value is not `{}`. A5 is stamped `Verified` by `Method: Design Decision` citing this record's own scoping choice; that is self-referential. | A write accessor declared as a stdin-fed shell executes tag values as shell code. The lint says green, by contract. | §1, §3, premortem, AT-4 |
| C-6 | `0027:C1` (`interpreter set:` line) | The record locks a claim over a deny-list it explicitly does not own and does not version. `shellInterpreters` (`internal/table/load.go:1035`) is 11 entries with per-name flag sets; `node` carries `{"-e","--eval"}`. C1 declares membership "not restated here… the list grows by amendment of THIS clause" — a frozen-at-lock invariant with no version marker. The moment anyone adds `perl` (which the A2 spike explicitly parks) or `python3` (BR4 parks it) or `--eval` variants, C1's blast radius changes silently and the shipped promise text goes stale with no mechanism to notice. `0027:S6` golden-asserts `Categories()` order but nothing asserts `shellInterpreters` content. | Six weeks after ship, someone lands `perl` on the deny-list per BR4's "noted for the successor." Every `["perl","-e",…]`-shaped tool argv anywhere in a model begins refusing, and the promise description still says nothing about it. | §2, premortem, AT-5 |
| C-7 | `0027:D-selection-predicate` | The tie-break is decided on the wrong axis and the record knows it. `["python","sh","-c","echo"]` reports `python -c` because `python` holds the lower `i`; the record calls this "message attribution only." Under C-2's expanded collision class the reported pair is frequently a coincidence and the *operative* pair is elsewhere in argv, so the detail names the wrong two words in exactly the cases where the author most needs help. The remediation text then compounds it. | Refusal detail names `python -c` for an argv whose actual shell is `sh -c` three words later; the author edits the wrong thing. | §1, premortem, AT-2 |
| C-8 | `0027:S5` | The record's sole regression test for the false-refusal class is one hand-picked vector, `["ruby","tool.rb","-e","prod"]`, chosen because it is *sympathetic* — the interpreter really is `ruby`. It is not a regression set; it is a mascot. Nothing in Testing Strategy or MVV exercises a non-interpreter argv0 with an incidental interpreter name mid-argv, which is where the widened predicate's whole new failure surface lives. | Every listed test passes and the six C-1 refusals ship. | §1, §5, AT-1, AT-2 |
| C-9 | `0027:§decision-rationale` (Premortem paragraph) | The record's own premortem names the false-refusal risk, then disarms it in the same sentence: "The first is the already-accepted flag-anywhere class, not new, but the message was the failure." It converts a correctness question into a message-quality question and declares it hardened without a single new vector. This is the exact self-exculpating move that lets C-1/C-2 through — the premortem was performed on the record rather than against it. | n/a (process defect; its symptoms are C-1 and C-2). | §1, §2, premortem |
| C-10 | `0027:§capability-dependencies` | The one row that matters is `Status: Deferred` on a capability with **no kata, no owner, no schedule** — A5's own durability caveat concedes it "exists only as that `Charted.md` line." A dependency with no tracking artifact is not deferred, it is abandoned, and C1's out-of-scope routing ("owned by the charted successor") is a citation to nothing. `0027:F2` calls the resulting admission "silent by contract"; a contract whose counterparty does not exist is not a contract. | Indefinite: `["sh","-s"]` lints green forever, with a record on file saying someone else owns it. | §3, premortem, AT-4 |
| C-11 | `0027:§references` | The References block cites `internal/cmdbind/cmdbind.go::resolveArgv0`; the real path is `internal/cli/cmdbind/cmdbind.go` (the Existing Infrastructure Audit has it right, References does not). A ground-sweep declared "clean (17 anchors) — every `path::Symbol`… confirmed verbatim by a fresh-context read." One of them is wrong, which means the sweep's clean verdict is not load-bearing evidence for the other sixteen. | A reader following the citation finds nothing; more importantly, the record's headline evidence-integrity claim is falsified by inspection. | §1, premortem, AT-6 |
| C-12 | `0027:MVV` step 3 / `0027:S4` | The MVV asserts admitted forms "lint green AND load with argv unchanged," and additionally that `["env","-S","sh -c echo"]` is "shown, by running the model under `--allow-commands`, to behave as C1 states." A4's own spike records that stock darwin BSD `env` **rejects `-S` outright**. So the MVV mandates a behavioural step that cannot execute on the stated development platform, guarded only by a parenthetical. An MVV step that is unrunnable on the dev host is a step that gets skipped. | The `-S` behavioural half is quietly dropped; the one-word-string admission ships asserted only by the predicate returning false, which proves nothing about what runs. | §5, AT-3 |
| C-13 | `0027:§approach` | The Approach unilaterally decides that intrastate#q2q1 (an open, independently-shippable kata) is subsumed and its `env` walk "is deleted." The Prerequisites row then hedges: "landed first, or closed as subsumed — either way the `env` walk is removed here." If q2q1 lands first, it will have added option-flag walking code that Phase 1 immediately deletes, and its own tests will encode the walk's behaviour. The record does not say what happens to those tests. | A merge conflict and an orphaned test suite asserting a walk that no longer exists. | §1, premortem, AT-7 |

---

## 1. The three most likely ways implementation goes wrong

### 1a. The predicate ships and starts refusing bindings that contain no shell

**Root cause in the RDR.** `0027:A2` is the only thing standing between this design and a broad false-refusal regression, and A2's evidence population is 32 argv vectors drawn from `internal/table/testdata/neg/*.toml` (6 vectors, all already-negative, two of them duplicates of `["rdr-gate","{artifact}"]`) plus 26 literals from `command_carrier_0025_test.go`. That leaves **11 green vectors** as the entire universe over which "zero new refusals" was demonstrated. Not one of those 11 was authored to contain a listed interpreter basename in a non-argv0 position, because until this RDR nothing made that position interesting. A2 therefore measures the frequency of a phenomenon in a corpus constructed before the phenomenon existed, and reports zero.

**The passage that enabled it.** `0027:A2`'s own scope disclaimer is the tell: *"Scope: this repo's fixtures and 0025's models, which is exactly the claim's scope; it says nothing about an out-of-tree user model, where the widened refusal is the intended behaviour."* The second clause is an assertion, not a finding. It presumes that any out-of-tree hit is a true positive. That presumption is the entire defect. The record never asks what fraction of out-of-tree hits would be false, and its Consequences line — *"the flag-anywhere laxity stays… (already true today)"* — asserts the class is unchanged when it is not.

I ran the widened predicate exactly as `0027:§illustrative-code` specifies, against ordinary tool invocations a model author would plausibly declare:

```
grep -r python -c src/                    old=false new=true   python -c
pip install python -c constraints.txt     old=false new=true   python -c
rg python -c                              old=false new=true   python -c
tar -czf out.tgz python -c                old=false new=true   python -c
npm run node -- -e                        old=false new=true   node -e
my-tool --lang php --mode -r              old=false new=true   php -r
```

Six new refusals from a list of fourteen guesses. `grep -c` is *count matching lines*. `pip -c` is *constraints file*. `rg -c` is *count*. These are not exotic; `grep` and `pip` are the two most common things a lint-adjacent shell script invokes. The false-positive rate of this predicate against realistic argv is not "the already-accepted class"; it is a new class roughly an order of magnitude larger, and the record contains zero evidence bearing on its size because A2 was pointed at a corpus that cannot contain it.

**Symptom the user sees.** A model that passed `intrastate lint` on Friday fails on Monday's upgrade with:

> declares the interpreter form `python -c`; inline shell is not a declared command; put it in a script and declare the script as argv0

for a binding that reads `command = ["grep", "-r", "python", "-c", "src/"]`. There is no shell. There is no inline code. There is no script to extract. The remediation is not merely unhelpful, it is incoherent — the author's only escape is to wrap a `grep` in a shell script, i.e. to introduce the exact hazard the check exists to prevent. `0027:A2`'s "If wrong" line predicts this word for word ("the refusal's remediation… misdirects, because there is no script to extract") and then stamps the assumption **Verified**.

### 1b. The predicate refuses argv in which no interpreter is ever executed

**Root cause in the RDR.** `0027:C1`'s predicate line reads: *"refuse iff there exist i < j with base(argv[i]) a listed interpreter and argv[j] one of its listed flags. Nothing before argv[i] is read."* The Approach celebrates this: *"nothing before the interpreter word is read at all,"* which is exactly why the design needs "no wrapper table and no option grammar." That is a real economy, and it is bought by discarding the only signal that distinguishes an interpreter from a string that spells one.

Under the old predicate, `argv[i]` is argv0 modulo an `env` walk — it is *the thing that gets executed*. `resolveArgv0` (`internal/cli/cmdbind/cmdbind.go:437`) resolves it and `exec.CommandContext(ctx, argv0, argv[1:]...)` runs it. Position was not incidental; it was the semantics. Delete it and `python` in `["grep","-r","python","-c","src/"]` is a *search pattern*, and the lint reports it as an interpreter form.

**The passage that enabled it.** `0027:D-selection-predicate`: *"the flag is matched at any later position, as today, so the already-accepted false-refusal class (a script argument that is itself `-c` / `-e`) is unchanged."* The word doing the damage is **unchanged**. Today's class is: *the interpreter is genuinely being executed, and one of its script's own arguments collides with an inline-code flag.* Tomorrow's class is: *any argv containing a listed name anywhere, followed anywhere by a matching flag.* Those share a sentence and nothing else. The record's Trade-offs repeats the elision — *"the flag-anywhere laxity stays… (already true today)"* — and `0027:§key-discoveries` states it as a Documented finding: *"Freeing the position adds no new false-refusal class beyond the one already accepted."* It is documented, and it is wrong.

Confirmed by running the predicate: `["echo","node","--eval"]` → refuses `node --eval`. `["cat","php","-r"]` → refuses `php -r`. `["mytool","--interpreter","python","--flag","-c"]` → refuses `python -c`. In all three the executed binary is `echo`, `cat`, `mytool`. No interpreter runs. No shell exists. The lint refuses anyway, and by C1's own text this is not a bug — it is the contract.

**Symptom the user sees.** A refusal that names two words the author can find in their argv and cannot understand the relationship between. Worse: because `0027:D-selection-predicate` fixes lowest-`i`-then-lowest-`j`, when a *real* `sh -c` is present alongside an incidental earlier interpreter name, the message names the incidental pair. `["python","sh","-c","echo"]` reports `python -c`. The record concedes this ("not the `sh -c` a reader's eye goes to") and classifies it as "message attribution only" — which is defensible when both pairs are real, and indefensible under a predicate where the reported pair is routinely a coincidence.

### 1c. Phase 3 has nowhere to land, so the promise never ships and only the predicate does

**Root cause in the RDR.** `0027:A6` is `Pending`. That alone should stop the lock. But the situation is worse than an unverified assumption: the spike A6 proposes is likely to come back negative.

- `table.Categories()` has **no non-test consumer**. `grep -rn "table.Categories()" internal/cli/ cmd/` (excluding `_test.go`) returns nothing. Every reference is in `roundtrip_test.go`, `dump_test.go`, `emit_witness_0024_test.go`, `reserved_key_fixtures_0008_test.go`. The category set is a *test-visible* artifact, not a user-visible one.
- `docs/cli-output-contract.md` — which the RDR itself names as governing ("the output contract governs what a verb may emit") — defines output as the `CLIError` envelope: `code`, `param`, `findings[]` with `message`/`hint`/`locator`. Every field is refusal-shaped. There is no success-path channel for "here is what this category does and does not cover."

So A6's spike must either find a surface that does not exist, or invent one — and inventing one is an amendment to an output contract this RDR does not own, does not cite in Capability Dependencies, and does not run a joint-decision check against.

**The passage that enabled it.** `0027:C1`'s `promise:` clause states normatively: *"the lint's user-facing description states the out-of-scope forms in the words above. No such description exists in the shipped code today… so Phase 3 creates it and THESE words are the text it ships."* This is a contract clause asserting the existence of a thing that does not exist, whose creation depends on a `Pending` assumption. `0027:G-assumptions` prohibits exactly this: *"no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it."* C1 is normative prose. It depends on A6. A6 is Pending. **The record cannot pass its own Finalization Gate as written**, and this is mechanically checkable at lock rather than a matter of judgement.

Compounding it, `0027:§phase-3-promise-wording` punts the decision: *"Where it lands is the implementer's call (A6)."* The record hands its single most important deliverable — the reviewer-facing text that is the entire justification for the predicate widening — to whoever picks up the ticket, with a `Pending` assumption as the brief.

**Symptom the user sees.** Nothing, which is the point. Phase 1 and 2 are mechanical and will land. Phase 3 will land as a doc comment beside the constant (`0027:A6` already concedes "a doc-comment beside the constant is not reviewer-reachable"), and the reviewer will receive a widened refusal predicate with no accompanying statement of what it still admits. The Problem Statement's opening complaint — *"They discover it only by reading argv themselves — the thing the lint exists to spare them"* — is left exactly as it was, plus six new false refusals.

---

## 2. The one section rewritten within 6 weeks of shipping

**`0027:C1`, and specifically its `predicate:` line.**

Not the `promise:` clause — that one simply never ships (§1c). The `predicate:` line is what gets rewritten, and the mechanism is already loaded.

The first support ticket after upgrade will be a false refusal of the C-1 class: `grep -c`, `pip -c`, `tar` with a `python` operand, or the equivalent in someone's private model. The obvious, cheap, locally-correct fix will be proposed within a day: *require argv[i] to be in command position, or at least not to be preceded by a non-wrapper word.* That is `0027:ALT2` — the per-wrapper walk with a bounded wrapper table — which C1 forbids in normative text: *"no wrapper table exists and none may be added."*

So the fix will take one of three shapes, all of which rewrite C1:

1. **A wrapper allow-list smuggled in as a "prefix" rule.** Someone adds "only scan after a known-safe prefix" and calls it not-a-wrapper-table. C1's prohibition is amended or quietly violated. `0027:§decision-rationale` factor (2) — the whole argument for position-freedom — collapses, because a mis-parse becomes a silent admission exactly as ALT2 predicted.
2. **The flag match is narrowed to `j == i+1`.** Cheap, kills most of C-1, and breaks `0027:S5` (`["ruby","tool.rb","-e","prod"]` stops refusing) and the `report:` contract that S5 exercises. It also silently un-refuses `["sh","-x","-c","code"]`. C1's predicate line is rewritten and `0027:S5`, `0027:D-selection-predicate` and the Consequences bullet go stale together.
3. **`shellInterpreters` is pruned.** `python`, `ruby`, `node`, `php` get dropped because they are the false-positive generators, leaving only true shells. This immediately contradicts C1's `out of scope` line, which is *channel*-scoped precisely so that `["python","-"]` is covered — `0027:S4` and the `trace` mini-check row 3c both rest on `python` being a listed member. Pruning the list makes the record's most carefully argued paragraph incoherent.

There is a fourth pressure with a shorter fuse. `0027:C1`'s `interpreter set:` line freezes membership by reference — *"the list grows by amendment of THIS clause"* — with **no version marker**. `0027:BR4` explicitly parks the `python3`/`nodejs` spellings "for the successor that widens the list," and `0027:A2` parks `perl`. Both are queued. When either lands, C1's blast radius changes and the shipped promise text (if it ever shipped) goes stale, and `0027:S6` — which golden-asserts `Categories()` order — does not cover `shellInterpreters` at all. The record locks a deny-list as an invariant and provides no way to detect that the invariant moved.

---

## 3. The one assumption that will not survive first contact with a real user

**`0027:A2`.**

Its claim: *"No binding that lints green today carries a listed interpreter basename as a non-argv0 word followed later by one of that interpreter's flags — widening the predicate refuses nothing currently accepted."*

Its method: compute both predicates over 32 argv vectors from this repository's own test fixtures.

Its verdict: Verified.

The first real user is by definition out-of-tree. `0027:A2` explicitly scopes itself away from them — *"it says nothing about an out-of-tree user model"* — and then explains why that is fine: *"where the widened refusal is the intended behaviour."* That sentence is the assumption that dies. It asserts that every out-of-tree hit is a true positive, and the six vectors in §1a falsify it in the first minute of contact.

The failure is structural, not a matter of insufficient diligence. A2 is being asked *"does this predicate over-refuse?"* and is answering it with a corpus of 11 green `command` bindings authored under the old predicate, in a repository whose own commands are `["rdr-gate","{artifact}"]`-shaped. Asking that corpus about false-positive rate is asking a question the corpus has no information about. The Verified stamp is not wrong about what it measured; it is wrong about what it licenses.

Two competing pressures make the collision certain rather than merely likely:

- The deny-list is deliberately broad — 11 names including four *scripting language* names (`python`, `ruby`, `node`, `php`) whose flags (`-c`, `-e`, `-r`, `--eval`) are among the most heavily reused single-letter flags in Unix. `-c` alone means *count* (grep, rg, wc, sort), *constraints* (pip), *config* (many), *command* (shells), *create* (tar), *check* (various).
- The predicate scans **every** word. So the collision probability is roughly (chance any argv word is one of 11 common names) × (chance any later word is one of ~5 extremely common flags). For argv containing tool names or language names as *data* — which any lint-adjacent model does constantly — this is not a tail event.

`0027:A5` is the runner-up (its "withholding" mechanism is falsified by `cmd.Stdin = bytes.NewReader(stdin)` being unconditional and by the write path carrying author-controlled bytes — see C-5), but A5's blast radius is a known-admitted hazard staying admitted. A2's blast radius is every user's working model breaking on upgrade. A2 is the one.

---

## 4. Premortem — written from eight weeks after ship

It is late October. RDR 0027 landed clean: Phase 1 replaced `interpreterForm`'s body with the eight-line position-free scan, Phase 2 added ten wrapper mutants beside `TestReq74_KnownInterpreterWithInlineCodeIsALoadTimeDefect`, all seven Testing Strategy scenarios went green, `go test ./internal/table/` was clean, and `table.Categories()` order was byte-identical. `0027:S1` through `0027:S6` passed on the first implementation attempt. Nobody had to argue about anything.

`0027:S7` did not pass, because it was never written. `0027:A6`'s spike came back with what anyone could have found in ten minutes: `table.Categories()` has no non-test consumer, and `docs/cli-output-contract.md` describes a `CLIError` envelope with `findings[]` — every field of it produced on refusal. There is no success-path surface on which a category description can be rendered. The implementer, holding `0027:§phase-3-promise-wording`'s instruction that *"where it lands is the implementer's call,"* did the reasonable thing: wrote C1's out-of-scope sentence as a doc comment above `CatCommandShellInterpreter` in `internal/table/category.go`, marked S7 as satisfied-in-spirit, and shipped. `0027:A6`'s own text had already conceded that a doc comment "is not reviewer-reachable." The premortem in `0027:§decision-rationale` had named this failure — *"a reviewer who read 'closes the wrapper class' stops reading argv and an `env -S "sh -c …"` string ships"* — and declared it mitigated on the strength of a description that does not exist.

Week two of the upgrade, the tickets started.

**Journey 1 — the release-gate model.** A platform team runs a gate accessor that greps a manifest: `command = ["grep", "-r", "python", "-c", "manifests/"]`, counting Python service declarations to check a floor. It lints green for eight months. On upgrade, `carrierDefect` clause 4 fires: `interpreterForm` finds `base(argv[2]) == "python"` at `i=2`, scans forward, finds `"-c"` at `j=3`, returns `"python -c"`. The refusal reads *"declares the interpreter form python -c; inline shell is not a declared command; put it in a script and declare the script as argv0."* There is no shell. `grep`'s `-c` is `--count`. The remediation instructs the team to move a `grep` into a shell script — to wrap a shell-free invocation in a shell — in order to satisfy a check whose stated purpose is making shells visible. They do it, because the lint is a hard gate and the message names two words that really are in their argv. Inline shell in the repo goes **up** by one script as a direct consequence of a lint designed to reduce it.

**Journey 2 — the dependency-audit model.** A second team's read accessor is `command = ["pip", "install", "python", "-c", "constraints.txt"]`. Same refusal, `python -c`. Here `-c` is pip's `--constraint`. They file a bug. Triage traces it to `0027:C1` and finds the predicate is behaving exactly as contracted, then finds `0027:A2` stamped **Verified** with the finding *"Zero divergences: no vector moves `old=false → new=true`, so no green binding is newly refused."* The bug is closed as working-as-designed against a record that certifies it cannot happen. Nobody notices for three weeks that A2's population was 32 in-repo fixture vectors with 11 green `command` bindings, none of which could contain the pattern.

**Journey 3 — the one that is not a false positive.** A third model declares a write accessor `command = ["sh", "-s"]`. It lints green — correctly, by contract: `0027:C1`'s out-of-scope line names the stdin channel, and `0027:F2` classifies the admission as *"silent by contract."* `0027:A5` explains that the charted `stdin = "none"|"envelope"` successor owns it, because withholding the envelope "leaves a stdin-reading interpreter nothing to execute." Except `cmdbind.spawn` sets `cmd.Stdin = bytes.NewReader(stdin)` unconditionally, and `(*Writer).Apply` passes `stdinObject(planned)` — a JSON object whose **values are model-authored tag values**. For read and gate bindings the envelope is `stdinObject(nil)` = `{}`, and `printf '{}' | sh -s` exits 127 with `{}: command not found`, so those are inert — by an accident of JSON syntax, not by design. For the write binding the tag values go to `sh -s` as script text. A5's mechanism ("withhold the envelope") is not available for write, because for write the envelope *is* the payload. A5 is stamped **Verified** with `Method: Design Decision` and `Evidence: this RDR's scoping choice` — a self-referential stamp on a mechanism claim that `internal/cli/cmdbind/cmdbind.go:207` and `:723` falsify. The charted successor that was supposed to own this still has no kata; `0027:A5`'s durability caveat said so at lock, and `0027:§capability-dependencies` recorded it as `Deferred` with no owner.

**The retro.** Four things are on the whiteboard.

`0027:A2` verified the wrong population and was believed because it was Verified.

`0027:C1`'s predicate line was reviewed as an elegant simplification — *"no wrapper table and no option grammar"* — and nobody asked what `argv[i]` means once it is no longer the executed binary. `0027:§key-discoveries` had asserted, as Documented, that *"freeing the position adds no new false-refusal class beyond the one already accepted."* Three reviewers read that sentence and none tested it. It took fourteen guessed argv vectors to break.

`0027:S5` — the single regression case for the false-refusal class — was `["ruby","tool.rb","-e","prod"]`, where `ruby` genuinely is the interpreter. The record's only defense against over-refusal was a test case chosen to be sympathetic to the design.

And the record's own premortem paragraph in `0027:§decision-rationale` had named the false-refusal risk, then disposed of it in the same breath: *"The first is the already-accepted flag-anywhere class, not new, but the message was the failure… Neither forces a switch; both hardened the claim, not the predicate."* The premortem was run to be survived, not to find anything. It converted a correctness question into a copy-editing question and moved on, and every downstream lens — grounding, 3amigo — took the disposal at face value; the 3amigo consolidation's six findings are all about test-oracle strength and promise-text delivery, and not one of them is about whether the predicate is correct.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

### AT-1 — A2's population must be able to contain the phenomenon (catches C-1, C-8)

```gherkin
Scenario: the zero-new-refusals claim is tested against argv that could exhibit the pattern
  Given the widened predicate from 0027:§illustrative-code
  And a corpus of at least 50 realistic `command` argv vectors NOT drawn from this
      repository's fixtures — sourced from README examples, CI configs, and Makefiles
      of tools a model author would plausibly bind (grep, rg, pip, npm, tar, docker,
      jq, kubectl, ssh, find, aws)
  When both the current and widened predicates are computed over every vector
  Then the report states the count and the full list of `old=false -> new=true` vectors
  And for each such vector, states whether an interpreter is actually executed
  And 0027:A2 may be stamped Verified only if the false-positive list is EMPTY,
      or if the record adopts a false-positive budget and states it in C1
```

Run against the record as written, this fails immediately: 6 of 14 guessed vectors flip, none executes a shell. The record's actual A2 spike could not have produced this outcome because its corpus (`internal/table/testdata/neg/*.toml` plus `command_carrier_0025_test.go`) contains 11 green vectors and no interpreter names as data.

### AT-2 — the predicate must not refuse argv that executes no interpreter (catches C-2, C-7, C-8)

```gherkin
Scenario Outline: an argv whose executed binary is not a listed interpreter must lint green
  Given command = <argv>
  And base(argv[0]) is not a key of shellInterpreters
  And no listed interpreter appears in a COMMAND position (argv[0], or after a
      pure-prefix wrapper chain)
  When `intrastate lint` runs
  Then the binding lints green

  Examples:
    | argv                                        |
    | ["echo","node","--eval"]                    |
    | ["cat","php","-r"]                          |
    | ["grep","-r","python","-c","src/"]          |
    | ["rg","python","-c"]                        |
    | ["pip","install","python","-c","cons.txt"]  |
    | ["mytool","--interpreter","python","-c"]    |
```

Every row fails under `0027:C1`. Confronting the record with this table at review forces the question `0027:D-selection-predicate` never asks: *is `argv[i]` still the executed binary?* It is not, and once that is on the table the "no wrapper table" economy has to be re-priced against a false-positive rate the record never estimated. Note this test is *incompatible* with C1 as written — that is the finding.

### AT-3 — the promise text must have a demonstrated landing site BEFORE lock (catches C-3, C-4, C-12)

```gherkin
Scenario: A6 resolves to a concrete, reviewer-reachable surface
  Given 0027:A6 is Pending and 0027:C1's `promise:` clause depends on it
  When the A6 spike runs
  Then it names an exact file:symbol where the description is rendered
  And it produces a captured terminal transcript of that text being emitted by a
      real `intrastate` verb on a model with NO defect present
  And it confirms the emission is permitted by docs/cli-output-contract.md, or
      names the contract amendment required and the RDR that owns it
  And until that transcript exists, 0027 does not lock

Scenario: the gate's own status-consistency rule is enforced
  Given 0027:G-assumptions requires that no Pending assumption has settled-fact
        prose depending on it
  When the Finalization Gate is run
  Then 0027:C1's `promise:` clause is flagged, because it asserts normatively that
       "the lint's user-facing description states the out-of-scope forms" while
       depending on Pending A6
```

Run today: `grep -rn "table.Categories()" internal/cli/ cmd/` outside tests returns nothing, and `docs/cli-output-contract.md` documents only refusal-shaped output. The spike returns *no surface exists* and the record is blocked, which is the correct outcome. The second scenario is purely mechanical and would have blocked the lock without any judgement call.

### AT-4 — A5's routing mechanism must be tested against the real stdin path (catches C-5, C-10)

```gherkin
Scenario: the charted stdin successor actually closes the axis it is handed
  Given a WRITE accessor declaring command = ["sh","-s"]
  And a planned tag whose VALUE is `touch /tmp/pwned`
  When the binding is invoked under --allow-commands
  Then the record states what happens

Scenario: the routing target exists as a tracked artifact
  Given 0027:C1 routes the stdin channel to the charted `stdin = "none"|"envelope"`
        successor
  When the tracker is queried
  Then an open kata for that successor exists with an owner
  And if none exists, C1 may not cite it as an owner; the forms are unowned admissions
      and 0027:§capability-dependencies says so
```

Scenario 1 exposes that `cmdbind.spawn` sets `cmd.Stdin` unconditionally (`:207`) and `(*Writer).Apply` sends `stdinObject(planned)` (`:723`), so "withholding the envelope" is not the available mechanism for write — the envelope is the payload. Scenario 2 is a one-command check that turns `0027:A5`'s parenthetical durability caveat into a lock decision instead of a footnote, and forces `0027:§capability-dependencies` to stop recording an owner that does not exist.

### AT-5 — the frozen deny-list needs a version marker and a change detector (catches C-6)

```gherkin
Scenario: shellInterpreters is pinned as an invariant of C1's promise
  Given 0027:C1 declares the interpreter set frozen and amendable only by amending C1
  And 0027:BR4 parks python3/nodejs and 0027:A2 parks perl for a future widening
  Then C1 carries a version marker for the list, and a golden test asserts the exact
       content of shellInterpreters (names AND per-name flags), not merely
       table.Categories() order as 0027:S6 does
  And that test names 0027:C1 so a future widening is forced through the clause
  And the shipped promise text is regenerated from, or asserted against, the list
```

`0027:S6` golden-asserts `Categories()` and nothing asserts `shellInterpreters`. Adding `perl` — which is already queued in two places — silently changes C1's blast radius with every test still green.

### AT-6 — the ground-sweep's clean verdict must be falsifiable (catches C-11)

```gherkin
Scenario: every path::Symbol in the record resolves on main
  Given 0027:§decision-rationale claims "Ground-sweep: clean (17 anchors) — every
        path::Symbol ... confirmed verbatim by a fresh-context read"
  When each cited path is checked against the working tree
  Then all resolve
```

Fails: `0027:§references` cites `internal/cmdbind/cmdbind.go::resolveArgv0`; the file is `internal/cli/cmdbind/cmdbind.go` (the Existing Infrastructure Audit uses the correct path, so the record contradicts itself). One demonstrable miss means the clean verdict was not produced by the process it describes, which withdraws it as evidence for the other sixteen anchors — including the peer-system citations that carry `0027:§decision-rationale`'s prior-art frame.

### AT-7 — the q2q1 subsumption must specify the merge order and the test disposition (catches C-13)

```gherkin
Scenario: intrastate#q2q1 landing first does not orphan tests
  Given 0027:§approach decides q2q1's env walk is subsumed and deleted here
  And 0027:§prerequisites permits q2q1 to land first
  When q2q1 lands first
  Then the record states which of q2q1's tests survive Phase 1's deletion of the walk,
       and which are deleted with it
  And states who closes q2q1 and with what disposition text
```

The record says *"the walk goes either way"* and stops. Either-way is a decision about code; it is silent about the test suite q2q1 will have added to assert the walk it took to write, which Phase 1 then deletes.
