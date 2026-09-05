Model: claude-opus-5

# Critique — RDR 0019, Owned-state initialization semantics

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0019:C1` | The carrier check type-switches on the WRITE binding only, but the emptiness read and the read-back both go through the READ accessor selected by `Registry.readerFor(def.Accessor.Role)`. A model with a file-backed writer and a command-backed reader passes the carrier gate and then cannot answer emptiness. | `init-state` on a model whose read accessor is `command = [...]` either refuses with an unrelated code (`flow-accessor-failed` / allow-commands) after the carrier gate said "file-backed, proceed", or — if the implementer papers over it — seeds on a fabricated emptiness answer. | §1, premortem, AT-1 |
| C-2 | `0019:A7` | `flowbind.Registry`'s write loop is a two-arm switch with a DEFAULT of `&Writer{Path: acc.Path}`, and `commandBacked` returns true for `acc.Path == ""` (the residue). The RDR asserts "three constructed types … mutually exclusive and exhaustive" from a reading of the *branches*, not the default. The default arm is the residue sink only when `commandBacked` is false — i.e. `Path != ""` — so the claimed set is right by accident of one predicate the verb is banned from calling. | An in-memory or future model shape yields `*flowbind.Writer` with `Path: ""`, the carrier gate admits it, and `load("")` returns an empty store: the verb seeds every `[initial]` key into nothing and reports success. | §1, §3, premortem, AT-2 |
| C-3 | `0019:RT3` | Byte-identity between the `init-state` route and the `set-state --write` transcription is asserted for "every value kind BOTH routes admit," naming `int` and `float` among the 9. `valueMembers` renders TOML numerics through `strconv.FormatInt`/`FormatFloat(_, 'g', -1, 64)` while `canonicalValue` stores the argv string VERBATIM after `conformKind`, which has no `float` arm at all and accepts `+5` for `int` via `strconv.Atoi`. | A model with `[initial] threshold = 1.0` seeds `"1"`; the operator's `--write threshold=1.0` writes `"1.0"`; `flow next` guard atoms and any downstream string comparison disagree between the two first-run paths, and the RT3 regression test either fails on day one or is quietly weakened to value-level. | §1, §2, premortem, AT-3 |
| C-4 | `0019:A6` | A6 is `Pending` and is the carrier for the ONLY new data flow the design does not already have. C1 nevertheless fixes the predicate's semantics "regardless of which lands", and the Gate's own Assumption Verification text forbids settled-fact prose depending on a `Pending` assumption — which C1, D-selection-predicate, RT2, RT4, the `disposition` table, S9, S11 and MVV steps 5/7/9 all are. | Implementation stalls at Step 2 with no authorized surface; the implementer invents one (exports `store`, or adds a `Cardinality()` to `ReadBinding`) under time pressure, and the 0004:C3 argument the whole RDR rests on is settled in a commit rather than a record. | §1, §2, premortem, AT-4 |
| C-5 | `0019:C1` | C1 declares the two new refusal code SPELLINGS non-normative "on the precedent of `0028:C1.3`", but the shipped code table's own comment states "The spellings are normative and the group fixes the exit," and `codeWriteEditRefused` is a fixed literal constant, not a floating one. The cited precedent says the opposite of what it is cited for. | Agents and skills that branch on `code` (the documented conformance oracle, REQ-67) get a code they cannot pin; two implementations of the same RDR emit different strings; `docs/cli-output-contract.md` gains an entry no test defends. | §1, premortem, AT-5 |
| C-6 | `0019:S9` | The sealed-artifact scenario's construction is a chain of five incidental implementation facts (`unreachable` is a pure suffix test on the declared Path; `Apply` seals and mutates in one `save`; `Reader.Read` short-circuits; the setup step must exit non-zero "expectedly"). None of these is contract; all are `flowbind` internals. | The one test that pins "count STORE keys, not owned keys" breaks or is deleted the first time `flowbind`'s seal representation changes, and the predicate silently regresses to owned-key counting — seeding beneath an unverifiable state. | §1, §2, premortem, AT-6 |
| C-7 | `0019:C1` | The ALL-quantifier predicate is stated over "EVERY bound artifact", but `runFlowSetState` only checks role bindings for writers it actually needs, and `parseArtifacts` binds whatever the caller names. The RDR never says whether an UNBOUND role's artifact counts as empty, non-empty, or a refusal — and S6 tests only the case where a needed role is unbound at the WRITE side. | Operator binds one of two roles, `init-state` refuses with `flow-artifact-missing` on a first run they thought was complete — or worse, treats the unbound sibling as empty and seeds half the model. | §1, §3, AT-7 |
| C-8 | `0019:MVV` | MVV step 10 runs against the shipped `models/rdr.toml`, whose `[initial]` is `gate_passed = "false"` — a `bool` tag whose initial value is authored as a TOML STRING. It happens to round-trip; nothing in C1 or the MVV says why, and no arm forbids `gate_passed = false` (TOML bool), which normalizes identically today but through a different `valueMembers` arm. | The MVV's "proves the user outcome on the repo's own gate model" claim rests on an authoring coincidence; the moment someone "fixes" the quoting in `models/rdr.toml`, the RDR's headline acceptance step is exercising a different code path than the one it certified. | §1, §2, AT-8 |
| C-9 | `0019:F5` | The empty-scalar asymmetry is disclosed, routed to RDR 0002 as a seed, and then C1's seed path is said to "assume it cannot arrive." But `loadInitial` admits `note = ""` today, `IsDecisionTable` does not exclude it, and no arm of C1 refuses it — so it CAN arrive at a verb whose scalar encoder is `members[0]`, i.e. the empty string, which `flowbind.Writer.Apply` will happily store. | `flow init-state` on a model with an empty-scalar `[initial]` succeeds and persists `{"note":""}` — a value `set-state` refuses to author and no doc describes — creating exactly the "third encoding" C1 claims cannot exist. | §1, §3, premortem, AT-9 |
| C-10 | `0019:ALT1` | Alternative 1 (lint-only) is rejected on "the user outcome is unmet", but the chosen design leaves the outcome unmet for the edit-carried (0028) and command-backed (0025) carriers — both shipped and Implemented — and unmet for every second run, every key added later, and every torn seed. The comparison table scores "First run wall removed: yes" for D without the carrier qualifier that C1 spends a paragraph on. | Operators on `edit`/command flows read the release notes, run `flow init-state`, get `flow-init-carrier-unsupported`, and conclude the feature is broken rather than scoped. | §1, §2, §3, AT-10 |
| C-11 | `0019:D-selection-predicate` | The disclosed residual — clearing the last key returns the artifact to the initializable class — is defended by "the boundary is contract text and a test, not a hidden edge." It is a hidden edge to the operator: nothing in the artifact, the payload, or `read-state` distinguishes "emptied by clears" from "never written," which is precisely the distinction the whole record exists to preserve. | An operator clears keys one at a time to reset, a scheduled `init-state` reseeds, and the operator sees state they explicitly deleted come back — the exact resurrect hazard the RDR rejected Alternative 2 for. | §1, §3, premortem, AT-11 |
| C-12 | `0019:A4` | The A4 consumer list enumerates 8 sites that hard-code "four verbs" and must change, and C1 inherits 0005's per-verb MUSTs "which 0005's test corpus asserts verb-by-verb." Neither the Implementation Plan nor the Testing Strategy has a step for updating 0005's verb-by-verb corpus, `llms.txt`, or the four-verb cardinal in `flowExtendedDesc`/`newFlowCmd`. | `go test ./internal/cli` fails on 0005's surface tests the moment the verb registers; `--help` and `llms.txt` tell the operator there are four verbs while the CLI has five. | §1, §2, premortem, AT-12 |
| C-13 | `0019:F2` | The torn-seed arm makes a re-run a NO-OP that reports missing keys but explicitly MUST NOT repair, and C1 forbids documenting the re-run as recovery. The stated recovery is "explicit `set-state` of the listed keys, or discarding the artifact." Discarding the artifact is not always possible — the artifact may hold non-`[initial]` owned keys another writer committed. | After a partial multi-writer seed, the operator's only correct recovery is to hand-transcribe `[initial]` values into `set-state` — the exact manual transcription this RDR exists to abolish, now on the unhappy path where it is most error-prone. | §1, premortem, AT-13 |
| C-14 | `0019:S8` | S8's read-back-mismatch fixture binds the read accessor to a DIFFERENT artifact path than its writer, pre-seeded with a conflicting value. But `checkAccessorBindings` requires every owned tag be served by exactly one reader, and read-back resolves the reader by ROLE — so the fixture must declare a reader and writer on the same role with different `path` values, a shape whose legality the RDR never establishes and no cited load-time check governs. | The single test defending "present-and-unverified is distinct from read-back-incomplete" is unconstructible as written; the implementer skips it, and a read-back mismatch reports the wrong class to an agent whose documented branch on it is retry. | §1, AT-14 |
| C-15 | `0019:C1` | The ORDERING clause says the carrier refusal "PREEMPTS the `--allow-commands` refusal (`internal/accessor/model.go::Registry.AllowCommands`, also exit 2)". The allow-commands refusal fires inside `cmdbind` at INVOCATION (`cmdbind.go:245`), and C1 separately guarantees no accessor is invoked. The preemption is vacuous, and S10's "pinning the ordering C1 fixes" pins nothing. | A test that claims to guard an ordering guards a tautology; a later change that moves the carrier check after registry invocation passes it. | §1, AT-15 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The carrier gate guards the wrong half of the seam

**Root cause.** The RDR treats "carrier" as a property of the WRITE accessor and builds its entire refusal on a type-switch over `internal/accessor/model.go::Definition.Binding` for `CapWrite` definitions. But the two things the verb actually needs — the emptiness answer and the commit-time read-back — both go through the READ path. `internal/accessor/executor.go` selects the read-back reader with `e.Registry.readerFor(def.Accessor.Role)`: the reader bound to the write accessor's *role*, not the write binding. `internal/cli/flowbind/registry.go`'s read loop applies `commandBacked(acc)` independently of the write loop. Nothing in the loader couples a role's read carrier to its write carrier.

**The passage that enabled it.** `0019:C1`, carrier scope: "The verb type-switches on `internal/accessor/model.go::Definition.Binding` (an exported field holding the exported binding type the registry selected) and admits ONLY `*flowbind.Writer`, refusing `*flowbind.EditWriter` and `*cmdbind.Writer`." And `0019:A7`, which frames the whole question as "the three constructed **binding** types" over "the registry's branches" — the write loop's branches. The RDR knows the hole exists and names it in the place where it can do the least good: `0019:S8` says, parenthetically, "A command-backed READER is the other legal construction — C1's carrier refusal is scoped to WRITE accessors." That sentence is written as a convenience for constructing a mismatch fixture. It is in fact the admission that C1's carrier scope does not cover the surface the predicate reads through.

**Symptom.** An operator with a file-backed writer and a command-backed reader — a legal, load-clean, shipped-supported shape — runs `flow init-state`. The carrier gate says file-backed, admit. The verb then tries to answer emptiness. Whatever A6 lands, if it goes through `ReadBinding`, it hits `cmdbind.Reader`, which spawns a process and refuses without `--allow-commands`, or answers with a key-scoped envelope that has no cardinality. The operator sees a refusal that names commands on a verb they were told is about initialization, or — if the implementer instead reaches around to `flowbind.load` because the write side "is file-backed" — the verb reads a JSON file the READER never uses, gets an empty store, and seeds a model whose actual state lives behind a command. That is silent state duplication under a green exit 0.

### 1.2 The exhaustiveness claim is right for a reason the RDR did not check

**Root cause.** `0019:A7` is `Pending` and asks to "read `::Registry`'s write-construction branches and confirm those three are the complete set a write accessor can yield (no fourth arm, no shared type between two branches)." The actual code is not a three-arm switch. It is:

```go
var binding accessor.Binding = &Writer{Path: acc.Path}
switch {
case len(acc.Edit) != 0:
    binding = NewEditWriter(acc, name)
case commandBacked(acc):
    binding = &cmdbind.Writer{Accessor: acc, Name: name, Config: cfg}
}
```

The DEFAULT is `&Writer{Path: acc.Path}` with no guard. `commandBacked` is `len(acc.Command) != 0 || acc.Path == ""` — its own doc comment states that the `acc.Path == ""` residue "takes the command arm too, where it builds a REFUSING binding rather than a `Path: ""` file binding," and that "C1's load-time exactly-one rule makes the residue unreachable through the loader; a binding built from an in-memory model never passed it." So `*flowbind.Writer` is the file-backed set only because `commandBacked` catches the empty-Path residue first. The RDR's verb is FORBIDDEN from calling `commandBacked` (C1: "Exporting `commandBacked` is not required and is not authorized here") and therefore cannot see the invariant that makes its own type-switch sound. It inherits a safety property from a predicate it is contractually barred from consulting.

**The passage that enabled it.** `0019:A7` Evidence: "Confirmed so far: `::Definition.Binding` is an exported field of an exported struct holding an exported `Binding` interface, and all three constructed types are exported." That is a statement about export status, not about exhaustiveness, and it is stamped as partial progress on an assumption the RDR then treats as settled in C1's normative text.

**Symptom.** The failure is not immediate — it is a landmine. Any future change that adds a write carrier, or that relaxes `commandBacked`'s residue arm (which its own comment flags as load-order-dependent), makes the type-switch admit a `*flowbind.Writer` with `Path: ""`. `flowbind.load("")` returns `os.ErrNotExist` → `store{}`, empty. The verb seeds. `save("")` fails or writes to a stray path. The operator sees "seeded 3 keys" and an artifact that does not exist.

### 1.3 The byte-identity invariant is false for numeric kinds

**Root cause.** `0019:RT3` and `0019:S4` assert that `init-state` into a fresh artifact and `set-state --write` of the same assignments into a second fresh artifact produce byte-identical files, "for every value kind BOTH routes admit," enumerated as nine including `int` and `float`. The two routes do not normalize numerics the same way:

- Loader route: `internal/table/load.go::valueMembers` on TOML `int64` → `strconv.FormatInt(t, 10)`; on `float64` → `strconv.FormatFloat(t, 'g', -1, 64)`.
- Argv route: `internal/cli/flow_input.go::canonicalValue` returns `value` **verbatim** after `table.ConformValue`. `conformKind` has an `int` arm (`strconv.Atoi`) and a `bool` arm and **no `float` arm at all**; `conformDomain` has `enum`, `set`, `int` and no `float`.

So `[initial] threshold = 1.0` becomes `"1"` while `--write threshold=1.0` becomes `"1.0"`. `[initial] n = +5` is not even authorable as TOML, but `--write n=+5` passes `Atoi` and stores `"+5"` where the loader would store `"5"`. `2.50` → `"2.5"` vs `"2.50"`. `1e10` → `"1e+10"` vs `"1e10"`.

**The passage that enabled it.** `0019:A2`, Status **Verified**, Method **Spike**: "RENDERING agrees byte-identically on all 9 kinds both routes admit: enum, scalar string, bool, int, float, set array (sorted), set with duplicates and HTML characters (deduped, escaping off), empty set `[]`, single-member set." The spike compared canonical *forms* for values chosen so the two renderings coincide — an integer like `5`, a float like `0.5`. It did not vary the SPELLING, which is the only axis on which these two routes can differ, because one preserves the author's bytes and the other reformats them. A2's "If wrong" clause is exactly right about the consequence and was never triggered because the table was not adversarial.

**Symptom.** Two operators, two first runs, two different artifacts. One ran `init-state`; one hand-wrote `set-state --write` from the model. Their `flow next` results diverge on any guard atom comparing the value, and `flow read-state` shows `"1"` in one artifact and `"1.0"` in the other for the same declared initial. The RT3 regression test fails the first time someone puts a float in `[initial]` — and the likely resolution under schedule pressure is to weaken RT3 from byte-identity to value-identity, which discards the single strongest property the RDR claims.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`0019:§normative-contracts` — C1's CARRIER SCOPE and EMPTINESS IS A NEW DATA FLOW paragraphs.**

They will be rewritten because they are the only part of C1 that is not a decision. C1 has four genuine decisions in it — `[initial]` is a runtime source; the act is explicit and persisting; the predicate is store-emptiness under ALL; seed values come from the loader-normalized model. Those will hold. Bolted onto them are two paragraphs that say, in normative voice, that the mechanism is undecided:

> "Which carrier serves it — a new `ReadBinding` capability, an exported `flowbind` cardinality probe, or a `read-state`-family surface — is UNDECIDED here and is booked as A6."

and

> "Extending the predicate to those carriers is deliberately out of scope here (see A6, and the successor noted in Consequences); this clause fixes the boundary so the verb cannot silently do the wrong thing at it."

A6 is `Pending`. The Finalization Gate's own Assumption Verification instruction is explicit: "**Status consistency:** no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it." C1's predicate, `0019:D-selection-predicate`, `0019:RT2`, `0019:RT4`, the `disposition` mini-check table, `0019:S9`, `0019:S11`, and MVV steps 5, 7, and 9 are all settled-fact prose depending on A6. This record cannot lock in its present state without the gate lying, and if it locks anyway, the first implementation commit will pick a carrier — and picking a carrier is a contract decision about `0004:C3`, which means C1 gets amended or a successor RDR immediately supersedes these paragraphs.

The carrier-scope paragraph will be rewritten on a second axis too. `0019:§consequences` calls the file-backed-only scope "a live scope gap, not a hypothetical one: those flows keep the manual first-run wall this RDR removes elsewhere." Both excluded carriers are shipped and Implemented (0025, 0028). The first support question after release is "why does `init-state` refuse on my flow," and the answer — "because store emptiness is undefined for your carrier" — is an implementation fact, not a user-comprehensible boundary. The pressure to extend will land inside six weeks, and the carrier-scope paragraph is what has to move.

Third, `0019:S9`'s sealed-artifact construction is load-bearing for C1's "count STORE keys, not owned keys" clause, and it depends on four `flowbind` internals with no contract standing: that `unreachable` is a pure suffix test on the declared Path, that `Apply` seals and mutates in one `save`, that `Reader.Read` short-circuits on the seal, and that the setup step's non-zero exit is "expected." That test is the first thing deleted when `flowbind` changes, and its deletion silently reverts the predicate.

---

## 3. The one assumption that will not survive first contact with a real user

**`0019:A5` / `0019:D-selection-predicate` — that "the store carries no key" is a boundary an operator can reason about.**

Everything downstream of the empty-store predicate is defended on the grounds that the boundary is *disclosed*. `0019:§risks-and-mitigations`: "the no-op payload and the docs must state the boundary in store terms ('while any key remains')." `0019:§consequences`: "the boundary is contract text and a test, not a hidden edge."

It is a hidden edge. The operator's mental model is intent: "I cleared my state to start over." The system's model is cardinality: "your store has 1 key left, so I decline; now it has 0, so I reseed." Nothing in the artifact, nothing in `read-state`, and nothing in `init-state`'s payload lets the operator see which side of the line they are on *before* running the command. `flowbind::load` returns the same zero-key store for an absent file and one emptied by clears — A5 verified this byte-for-byte and treated it as reassurance. It is the problem. The RDR rejected Alternative 2 precisely because "'cleared' and 'unseeded' become indistinguishable" — and then adopted a predicate whose entire safety argument rests on those two states being indistinguishable at the content level, confined (the RDR says) to "the single state where the store carries nothing at all."

That single state is not rare. It is the state of every operator who did what the docs told them: cleared their owned keys to reset. `models/rdr.toml` declares exactly three owned keys. Clear all three — the natural "reset" gesture — and the store is empty. The next scheduled `init-state`, which `0019:§risks-and-mitigations` explicitly contemplates ("a scheduled `init-state` reseeds"), resurrects all three. The RDR's own risk register names this outcome and then dismisses it as "C1's disclosed residual, not a defect."

Worse, the boundary is not just invisible — it is inverted relative to intent. Clearing MORE keys makes the reseed MORE likely. An operator who clears two of three keys is protected; one who finishes the job is not. No amount of documentation makes "clear fewer keys to stay safe" a rule a human internalizes. The `[initial]`-key-added-after-seeding case compounds it: `0019:F4` accepts that a new `[initial]` key requires manual `set-state` forever, so the operator's incentive is to clear everything and reinit — walking straight into the residual.

The assumption that will not survive is not "an emptied store equals an absent store" (that is true of the bytes). It is the assumption that stating this in contract text and pinning it with MVV step 9 constitutes disclosure to a user. It does not. It constitutes disclosure to a reviewer.

---

## 4. Premortem

*Written from twelve weeks after ship.*

`flow init-state` landed in the 0.19 release. The MVV was green: all ten steps, including step 10 against `models/rdr.toml`. Everything below happened anyway.

**Week 1 — the carrier gate.** The first external report came from a team running an `edit`-carried writer over a Markdown status block (0028). `newFlowInitStateCmd` type-switched `def.Binding`, saw `*flowbind.EditWriter`, and returned `flow-init-carrier-unsupported`. Correct per C1. The report was filed as a bug three times before anyone routed it to `0019:§consequences`. The release notes had said "removes the first-run wall for state-machine flows"; the docs constraint in `0019:§activation-step-1-contract-and-reference-docs` — "it states the FILE-BACKED scope … not an unqualified 'every state-machine flow'" — had been written into `docs/cli-reference.md` and into nothing that a user reads first.

**Week 2 — the reader that wasn't checked.** A second team had a file-backed `[write.status]` and a command-backed `[read.status]` on the same role — perfectly legal, `checkAccessorBindings` clean, one reader per owned tag. `runFlowInitState` ran the carrier switch on the WRITE definition, admitted it, and called the emptiness probe. The probe had been implemented as a `Cardinality()` method added to `ReadBinding` — the first of A6's three candidates, chosen in a Tuesday standup because it was the smallest diff. `cmdbind.Reader` got a `Cardinality()` that returned `0, ErrUnsupported`. The verb, having no arm for that, surfaced `flow-accessor-failed`. The user's read: "init-state is broken for command accessors," which contradicted the carrier refusal they'd read about, which was scoped to writers. `0019:S8` had named this shape in a parenthesis — "A command-backed READER is the other legal construction — C1's carrier refusal is scoped to WRITE accessors" — as a fixture-construction convenience. Nobody had read it as a scope hole in C1, because it was filed under a read-back-mismatch scenario.

**Week 3 — the float.** A user added `[initial] confidence = 0.5` and then, on a second machine, ran `flow set-state --write confidence=0.50` to match. `valueMembers` had rendered `0.5` through `FormatFloat(_, 'g', -1, 64)` → `"0.5"`; `canonicalValue` had passed `"0.50"` through verbatim, since `conformKind` has no `float` arm and `conformDomain` has none either. Two artifacts, two byte sequences, one declared initial. `flow next` disagreed between them on a `gte` guard. The RT3 test — S4, promoted from A2's spike — had nine rows, and its float row used `0.5` on both sides, so it had never caught it. The fix shipped as a value-level comparison in RT3, because making the two routes agree meant either reformatting the operator's argv string (changing `set-state`, out of scope) or storing the author's TOML spelling (changing the loader, out of scope). The byte-equality property C1 spends three paragraphs establishing — "which is what makes read-back equality byte equality and leaves no third encoding" — was quietly downgraded to value equality in the test that was supposed to defend it.

**Week 4 — the reset.** An operator on the RDR gate model itself cleared `stage`, `status`, and `gate_passed` to restart a record. The nightly job ran `flow init-state` unconditionally — exactly the automation pattern `0019:A5` was written to survive. The store had zero keys. It reseeded `stage=seeded`, `status=draft`, `gate_passed=false`. The record's state came back from the dead overnight. The operator's ticket quoted the RDR's own Alternative 2 rejection back at the team: "A cleared owned tag silently resurrects — 'cleared' and 'unseeded' become indistinguishable, inverting REQ-107's read-back guarantee." `0019:F4` and `0019:§risks-and-mitigations` had both called this the disclosed residual. MVV step 9 had pinned it in the failing direction, on purpose, so it could not move silently. It had not moved. It had simply been shipped.

**Week 6 — the torn seed.** A two-writer model (two roles, two artifacts) had role A commit and role B fail on a permissions error. Per `0019:F2`, the re-run was a NO-OP with a payload listing B's missing keys, and C1 forbade documenting the re-run as recovery. The operator's actual recovery was to hand-type `set-state --write` for each of B's `[initial]` keys, reading them out of the model file. That is the transcription this RDR exists to abolish, now required on the unhappy path. "Discard the artifact and re-run init" was not available: A's artifact also carried two non-`[initial]` owned keys another writer had committed.

**Week 8 — the verb count.** `--help` still said "The four verbs are the whole skill-integration surface." `flowExtendedDesc` still said "The four verbs divide one job." `llms.txt` had been regenerated in the same commit as `docs.go` but the two cardinals in `flow.go`'s `Long` body and `flowExtendedDesc` were string literals nobody grepped, because `0019:§phase-1-code-implementation` Step 1 says only "shared selection/tag/artifact registration, `ValidateMode`/respond gateway, class refusal." The eight-site consumer list lived in `0019:A4`'s Evidence — verified, thorough, and in a section the implementer read once during review and never opened again, because the Implementation Plan does not reference it.

**Week 10 — the code that wasn't a code.** An agent skill branched on `code == "flow-init-carrier-unsupported"`. A refactor renamed it to `flow-carrier-unsupported`, which was legal: C1 says "the literal spellings are non-normative … so no test may pin the string." The skill broke. `internal/cli/flow_input.go`'s own header comment — "The spellings are normative and the group fixes the exit. They are declared as constants in one place so a producer cannot spell one slightly differently at a second site" — had been true for every other code in the table. C1 had cited `0028:C1.3`'s `codeWriteEditRefused` as precedent for non-normative spellings; `codeWriteEditRefused` is a fixed constant with a fixed literal, cited in `docs/rdr/0028-declared-line-edit-writer/artifacts/deviations.md` as a new surface, not a floating one. The precedent had been read backwards.

**Week 12 — the postmortem's finding.** Nine of the twelve incidents trace to text that was already IN the record: the carrier scope, the sealed-store construction, the residual boundary, the torn-seed recovery, the eight consumer sites, the read-carrier parenthesis in S8. The record had disclosed almost everything that went wrong. It had not been structured so that disclosure reached an implementation step or a user-facing sentence. Disclosure went into Failure Modes, Consequences, and assumption Evidence blocks; the Implementation Plan is nineteen lines across three steps and one activation step, and none of them cite those disclosures. The one thing the record had genuinely not caught — the numeric spelling divergence — was covered by an assumption marked **Verified** by a spike whose table was chosen to agree.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

**AT-1 — the read carrier is checked, not just the write carrier** *(catches C-1)*
```gherkin
Given a model with [write.s] path="s.json" role="s" and [read.s] command=["helper"] role="s"
  And the model loads clean (one reader per owned tag, one writer per [initial] key)
When flow init-state is run against it
Then it MUST refuse with the carrier code naming the READ accessor
  And no accessor process is spawned
  And the refusal is identical whether or not --allow-commands is passed
```
Review-time form: C1's carrier scope must state the predicate over BOTH capabilities for every bound role, and A7 must enumerate `Registry`'s READ loop branches as well as its write loop. As written, A7 asks only about write-construction branches.

**AT-2 — the admitted set is the file-backed set, not a proxy** *(catches C-2)*
```gherkin
Given an accessor.Definition{Capability: CapWrite, Binding: &flowbind.Writer{Path: ""}}
When the verb's carrier type-switch is applied
Then it MUST refuse
  And the refusal MUST NOT depend on flowbind.commandBacked being consulted first
```
Review-time form: A7's Evidence must show `Registry`'s write loop is a switch with an UNGUARDED DEFAULT of `&Writer{Path: acc.Path}`, and must state which invariant makes `Path != ""` hold for that default. It currently says only that three types are exported.

**AT-3 — byte identity across numeric SPELLINGS, not just numeric values** *(catches C-3)*
```gherkin
Given a tag declared kind="float" and a tag declared kind="int"
  And [initial] entries authored as 1.0, 2.50, 1e10, and 5
When init-state seeds a fresh artifact
  And set-state --write is given the same literals verbatim into a second fresh artifact
Then the two artifacts MUST be byte-identical
```
This fails today on every float row and on `+5`. Review-time form: A2's spike table must vary the SPELLING of each numeric, not just its value, and must record that `conformKind`/`conformDomain` have no `float` arm — so the argv route stores the author's bytes while the loader route reformats them. A2's "Verified" stamp should not have been granted on a table where both routes were fed values whose two renderings coincide.

**AT-4 — no normative clause depends on a Pending assumption** *(catches C-4)*
```
Step 1: List every element whose text asserts the store-emptiness predicate as settled.
Step 2: For each, check whether A6 is Verified.
Step 3: If A6 is Pending, the record MUST NOT lock.
Expected at review: C1, D-selection-predicate, RT2, RT4, the disposition table,
  S9, S11, and MVV steps 5/7/9 all fail step 3.
```
This is the Finalization Gate's own Assumption Verification instruction, applied mechanically. Running it before pre-lock would have blocked the record or forced A6 to resolve first.

**AT-5 — refusal-code spellings are normative or they are not** *(catches C-5)*
```gherkin
Given the shipped code table's header comment "The spellings are normative and the group fixes the exit"
  And 0028's codeWriteEditRefused, a fixed constant with a fixed literal
When C1 declares its two new codes' spellings non-normative "on the precedent of 0028:C1.3"
Then the review MUST resolve the contradiction before lock
  And if the codes are normative, C1 MUST fix their literals and S7/S10 MUST pin them
```

**AT-6 — the STORE-key predicate is pinned by contract, not by flowbind internals** *(catches C-6)*
```gherkin
Given a bound artifact whose only key is a non-owned key
When init-state evaluates the emptiness predicate
Then it MUST decline to seed
  And the test constructing that state MUST NOT depend on
      unreachable() being a path-suffix test,
      Apply() sealing and mutating in one save,
      or a setup step exiting non-zero
```
Review-time form: S9 needs a construction that is stable under any `flowbind` seal representation — which likely means the A6 carrier must expose a way to construct or inject a non-owned store key, i.e. A6 must resolve before S9 can be written.

**AT-7 — the ALL quantifier's domain is defined** *(catches C-7)*
```gherkin
Given a model with two write roles, A and B
When init-state is invoked with only role A bound
Then C1 MUST state exactly one of: refuse (artifact-binding family), or treat B as out of the quantifier's domain
  And a test MUST exist for whichever is chosen
```
Neither C1 nor S6 nor S11 covers this. S6 covers a role unbound at the write side of a single-role model; S11 covers two roles both bound.

**AT-8 — the MVV's shipped-model step exercises the declared kind, not an authoring coincidence** *(catches C-8)*
```gherkin
Given models/rdr.toml declares gate_passed as kind="bool" with [initial] gate_passed = "false"
When the same key is authored as the TOML boolean false
Then init-state MUST produce a byte-identical artifact
  And MVV step 10 MUST assert this, or the model MUST be normalized before the step is credited
```

**AT-9 — the empty scalar cannot arrive, or C1 refuses it** *(catches C-9)*
```gherkin
Given a state-machine model with an owned scalar tag `note` and [initial] note = ""
When the model is loaded
Then either the loader MUST refuse it (making C1's "assume it cannot arrive" true)
  Or C1 MUST name an explicit arm for it
Expected today: the loader ADMITS it, C1 has no arm, and the scalar encoder members[0]
  would persist "" — a value set-state refuses to author.
```
`0019:F5` correctly identifies the asymmetry and then routes it out of scope while C1 relies on it not arriving. Those two cannot both stand.

**AT-10 — the outcome claim is scoped everywhere it appears** *(catches C-10)*
```
Step 1: Grep the record for every unqualified claim that the first-run wall is removed.
Step 2: Confirm each carries the file-backed qualifier.
Expected at review: the Decision Rationale scoring table's "First run wall removed: yes"
  row for column D carries no carrier qualifier, and the Problem Statement's
  qualifier appears only in its final sentence.
```

**AT-11 — the residual boundary is observable to the operator before it fires** *(catches C-11)*
```gherkin
Given an artifact whose owned keys have all been cleared
When the operator runs flow read-state
Then the output MUST distinguish "emptied by clears" from "never written"
  Or C1 MUST state that it cannot, and the docs MUST say init-state may reseed a reset artifact
Expected today: read-state cannot distinguish them, and the disclosure lives
  only in RDR prose and MVV step 9.
```

**AT-12 — the eight consumer sites are implementation steps, not assumption prose** *(catches C-12)*
```
Step 1: For each of A4's eight enumerated consumers, confirm the Implementation Plan
  or Testing Strategy names it.
Expected at review: zero of eight appear in Phase 1 or Phase 2; the four-verb
  cardinals in flow.go's Long body and flowExtendedDesc appear nowhere outside A4.
```

**AT-13 — torn-seed recovery does not reintroduce transcription** *(catches C-13)*
```gherkin
Given a torn seed where role B's [initial] keys are un-seeded
  And role A's artifact carries non-[initial] owned keys committed by another writer
When the operator follows C1's stated recovery
Then the recovery MUST NOT require reading [initial] values out of the model by hand
Expected today: it does, because "discard the artifact" loses A's other keys
  and "explicit set-state" is transcription.
```

**AT-14 — S8's fixture is constructible** *(catches C-14)*
```gherkin
Given checkAccessorBindings requires exactly one reader per owned tag
  And read-back resolves the reader by ROLE
When S8 binds "the READ accessor to a different artifact path than its writer"
Then the RDR MUST cite the load-time rule that permits a reader and writer on one role
    with differing `path` values
Expected at review: no such citation exists.
```

**AT-15 — the ordering preemption is not vacuous** *(catches C-15)*
```gherkin
Given the --allow-commands refusal fires inside cmdbind at invocation (cmdbind.go:245)
  And C1 guarantees the carrier refusal precedes any accessor invocation
When S10 asserts the carrier code "rather than the allow-commands refusal"
Then the assertion MUST distinguish the carrier check from any pre-invocation ordering,
    or C1 MUST drop the preemption clause as tautological.
```
