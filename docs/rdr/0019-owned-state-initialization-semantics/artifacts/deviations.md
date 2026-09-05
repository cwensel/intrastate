# Deviations — RDR 0019 Owned-state initialization semantics

Phase 2 artifact. Every entry carries exactly one Type. A `Status:
mechanical translation` entry was resolved against the record, its
verified assumptions, or source, and the implementation continued past
it. A `Status: needs author decision` entry is a design question the
record's own evidence cannot close; it is recorded and the run
continued.

---

## DEV-1 — `float` is not a declared tag kind

**Type**: TEST-FIXTURE

`0019:S4` names nine value kinds "BOTH routes admit", and lists `float`
among them. Phase 1 read that as a `kind = "float"` declaration and
wrote `initKindModel`'s `k_float` and `initFloatBoundaryModel`'s
`threshold` that way. No such kind exists: RDR 0003's closed five-token
vocabulary in `internal/table/model.go::declaredKinds` is
`{enum, bool, int, set, scalar}`, and the loader refuses anything else
as `malformed_tag_declaration`, so both fixtures failed at LOAD.

The record resolves it against itself. A2's `Status: Verified` evidence
states the fact outright: "(`float` is not a declared kind — the closed
vocabulary in `internal/table/model.go` is `{enum, bool, int, set,
scalar}` — so a TOML float reaches the seed encoder under `scalar`,
which `conformKind` passes through unchecked.)" S4's "float" row is a
float VALUE, not a float KIND.

Both declarations changed to `kind = "scalar"`. The rows keep their
discriminating power: `k_float`'s `[initial] k_float = "1"` still pairs
against `--write k_float=1` for RT3's byte identity, and
`initFloatBoundaryModel`'s `[initial] threshold = 1.0` still normalizes
to `"1"` against `--write threshold=1.0`'s `"1.0"`, which is the whole
of REQ-73 / REQ-82's boundary row.

**Status**: mechanical translation

---

## DEV-2 — reader uniqueness is per-KEY, not per-ROLE

**Type**: TEST-FIXTURE

`0019:C1` fixes FIRST MATCH IN REGISTRY ORDER for the role's reader on
the ground that "role uniqueness on the read side is NOT enforced
today … so two same-role readers are constructible". That is correct.
Phase 1's `initTwoReaderModel` built the two same-role readers over the
SAME KEY, which is a different rule and IS enforced:
`internal/table/load.go::checkAccessorBindings` refuses "owned tag
`stage` is served by 2 readers; want exactly one" at LOAD, so the model
never reached the verb and the test asserted nothing.

The second reader now serves its own owned key (`shadow`, added to the
writer's `keys` so it is writer-served and the model loads). The ROLE
collision — the thing REQ-43 / REQ-46 discriminate — is untouched: two
readers on role `state`, the first file-backed and the second
command-backed, and the gate must admit the model by taking the first.

**Status**: mechanical translation

---

## DEV-3 — the sealing write leaves its owned key beside the seal

**Type**: TEST-FIXTURE

S9 needs "an artifact whose owned keys have ALL been cleared but whose
last write declared an unreachable read-back locator" — a store whose
ONLY key is `::sealedKey`. Phase 1's setup performed the sealing WRITE
and then asserted the post-setup bytes carried no `stage`. They do:
`flowbind::Writer.Apply` applies the mutation AND seals in one atomic
save, deliberately, because "the applied-but-unverified sense
`0004:C14` requires is a claim about an artifact that really was
mutated". The setup's own `t.Fatalf` fired, exactly as REQ-89's
fail-loudly obligation requires.

The setup gains the clear the scenario's own prose names: a `--clear
stage` through the same sealing writer removes the key and re-seals,
leaving exactly `{sealedKey}`. The record anticipates this — S9's own
text says "clearing through a sealing writer both re-seals and exits
non-zero". Both S9 tests take the added step.

**Status**: mechanical translation

---

## DEV-4 — the HTML-escaping predicate was inverted

**Type**: TEST-FIXTURE

`TestReq18And19_0019_SetSeedsTakeTheCanonicalSortedCompactUnescapedForm`
asserted the artifact does NOT contain `<` or `&`, with a message
naming `escapedSetLiteral` as "the defect spelling". Escaping DISABLED
is precisely what leaves those two bytes in the artifact verbatim, so
the predicate asserted the opposite of the clause it cites and could
never pass against a conforming encoder.

The predicate now tests for the ESCAPED spelling (`<`, `&`),
which is what the failure message already described. The assertion is
unchanged.

**Status**: mechanical translation

---

## DEV-5 — the tautology guard could not see an unknown-flag refusal

**Type**: TEST-FIXTURE

`::requireVerbRegistered` fatals when a refusal's CODE is
`command-error`, to exclude the tautology of a missing verb. But
`internal/cli/root.go::cobraErrorToCLIError` maps EVERY cobra/pflag
error to that one code, so an unknown FLAG on a REGISTERED verb carries
it too — and REQ-10 ("no write grammar") and REQ-102 ("no per-role
repair grammar") both assert exactly that unknown-flag refusal, then
call the guard, which rejected it. Seven subtests could never pass.

The guard now also requires cobra's unregistered-command MESSAGE
(`unknown command "<name>" for "<parent>"`), which is what actually
separates the two. It still excludes the tautology it was written for,
and no longer excludes the refusals the contract fixes.

**Status**: mechanical translation

---

## DEV-6 — an inherited persistent flag is not on `Flags()` before Execute

**Type**: TEST-FIXTURE

REQ-11 fixes that `--allow-commands` RESOLVES on `init-state` by
INHERITANCE and is not DECLARED on the verb's own set. Phase 1 probed
the first half with `cmd.Flags().Lookup`. Cobra merges a parent's
persistent set into a child's `Flags()` only during `Execute`, and the
test walks an unexecuted `NewRootCmd()` tree, so that lookup returns
nil for every child of the group — the four shipped verbs included.

The probe is now `cmd.InheritedFlags().Lookup`, which is what observes
inheritance on an unexecuted tree. The second half
(`LocalNonPersistentFlags()` must NOT carry it) is untouched, so the
boundary REQ-11 draws between inheriting and declaring still bites.

**Status**: mechanical translation

---

## DEV-7 — the MVV fixture could not lint clean

**Type**: TEST-FIXTURE

REQ-MVV.1 gates the whole spine on "`flow lint` certifies the model
(0006 arms all green)". Phase 1's `initMVVModel` carried
`[rule.guard.all.note] eq = "hello"` over an OPTIONAL scalar tag with
no finite declared domain, which is a blocking
`graph-unprovable-coverage` under two of 0006's arms at once
(`row-can-refuse` and `dimension-not-finite`). The fixture therefore
failed step 1 by construction, taking the spine and
`TestReq99And100_0019_LintGainsNoArm` with it.

`0019:MVV` requires of the fixture only "one always-present owned key
and one plain owned key, both writer-served". The guard atom is not
part of that, and nothing else in the suite reads `note` through a
guard — it is the `--clear` target of MVV steps 6 and 9. The atom is
dropped from `initMVVModel` and from `initOneKeyInitialModel`, which
Phase 1 wrote as its one-key twin.

**Status**: mechanical translation

---

## DEV-8 — the exit-3 re-run assertion contradicted REQ-67

**Type**: TEST-FIXTURE

`TestReq12_0019_EverySeedIsRoutedThroughTheWriteAccessorWithCommitTimeReadBack`
drove `init-state` twice over ONE artifact and asserted exit 3 both
times. The first run seeds and refuses at exit 3, correctly, and leaves
the store SEALED — non-empty. REQ-67 fixes what the second run then
does: "both leave a non-empty store, so both make a re-run a no-op that
does not repair", which is exit 0. The test asserted exit 3 against the
clause its own sibling
(`TestReq67_0019_AnIncompleteReadBackExitsThreeAndLeavesANonEmptyStore`)
pins as exit 0.

The second assertion now runs over a FRESH artifact, which is the
empty-store condition the first run exercised. Both the exit-3 arm and
REQ-67's no-op re-run stay asserted, by the two tests that own them.

**Status**: mechanical translation

---

## DEV-9 — a read-back MISMATCH is unreachable through the file-backed carrier

**Type**: SPEC-DEFECT

`0019:C1` fixes a PRESENT-AND-UNVERIFIED refusal at exit 2 for "a
read-back that COMPLETED and disagreed", distinct from the
read-back-INCOMPLETE class at exit 3, and S8 states the construction
that is meant to reach it:

> the fixture binds the READ accessor to a different artifact path than
> its writer, pre-seeded with a conflicting value for the same key

and grounds it on `::checkAccessorBindings` constraining "no
relationship between a reader's `Path` and its writer's — nothing
requires the two to name one artifact".

That premise does not hold against the shipped seam, on three
independent counts, each verified in source:

1. **The read-back re-reads the WRITER's artifact, not the reader's.**
   `internal/accessor/executor.go::Executor.Write` resolves `art` from
   the WRITE definition's role (`::selects`, `e.Artifacts[def.Accessor.Role]`)
   and passes that same `art` to `e.invokeRead(ctx, reader, art, …)`. The
   reader's own role never selects an artifact for the read-back.
2. **A reader's declared `Path` is not where it reads.**
   `flowbind::Reader.Read` reads `art.Path` — the caller's
   `--artifact role=path` — and consults its declared `r.Path` only for
   the `::unreachable` short-circuit. Differing declared paths cannot
   make the read-back read a different file.
3. **A reader on a DIFFERENT role is no reader at all for the
   read-back.** `::Registry.readerFor` selects by the writer's role, so
   S8's two-role fixture yields `hasReader == false`, which is
   `ClassReadBackIncomplete` (exit 3) by construction — the very class
   the scenario asks to be distinguished FROM. The fixture's own setup
   `set-state` refuses at exit 3 before the scenario begins.

With the artifact fixed and both reads going through
`flowbind::Reader` over the file the writer just wrote, every remaining
route to `::verifyReadBack`'s true branch is closed: planned keys read
back byte-identically, `classify` turns an unanswered key into
`unread` (which is INCOMPLETE, not mismatch), and `::Writer.Apply`
touches only planned keys so no protected baseline can diverge.

This is not a defect this record introduces. RDR 0005's own suite never
reaches the class either — `internal/cli/flow_setstate_0005_test.go`'s
read-back test asserts the INCOMPLETE/TIMEOUT pair and names the
mismatch code only to exclude it — so `flow-write-readback-mismatch` is
a shipped code with no `flowbind` witness.

**Consequence.** `TestReq63And67And87_0019_AReadBackMismatchRefusesAtExitTwoAndTheReRunDoesNotRepairIt`
is RED and stays red. Its assertion is the contract's, not a fixture
error, so it is not weakened; and the state it asserts is not
constructible, so it is not made to pass.

**Why the evidence cannot close it.** The three routes to a reachable
mismatch are all changes this record forbids or does not carry:

- making the read-back resolve its artifact by the READER's role is a
  change to `internal/accessor`'s read-back seam (`0004:C12`/`C13`),
  which REQ-99 rules out ("no existing verb, flag, artifact shape or
  model schema changes") and which this record carries no override for;
- making `flowbind::Reader` read its DECLARED `Path` rather than the
  caller's artifact contradicts `0004:C3`'s ban on ambient artifact
  discovery, which REQ-3 restates as unchanged;
- retiring the arm — declaring the mismatch class unreachable for the
  file-backed carrier and scoping C1's PRESENT-AND-UNVERIFIED sentence
  and S8 accordingly — is an amendment to a Final record.

Which of the three is right is a design decision, and the record's own
evidence does not pick one: A1–A8 are silent on the read-back's
artifact resolution, and the Testing Strategy asserts S8 as
constructible on a premise source refutes.

**Status**: needs author decision

---

## DEV-10 — the emptiness carrier is one exported store-key probe

**Type**: IMPL-DECISION

A6's `Status: Verified` evidence names candidate (b) — "an exported
`flowbind` cardinality probe" — as the surviving carrier, and
ASSUMPTION-1 records Phase 1's reading of it. The implementation adds
`internal/cli/flowbind::StoreKeys`, which returns the store's KEY SET
rather than only its cardinality.

The widening is not a second surface: it is the same probe answering
the same question, and it is what keeps the NO-OP arm's absent-key
report on the same carrier as the predicate. Routing that report
through the declared readers instead re-enters `::Reader.Read`'s
unreadable short-circuit on a SEALED store, which turns REQ-28's exit-0
no-op into an exit-3 refusal — the exact degradation REQ-29 names and
S9 fails on. One probe, one read of the artifact, no key VALUE crossing
the seam, and the seal REPORTED rather than filtered (C1 fixes the
count as over STORE keys, and the sealed store is the arm where store
keys and owned keys differ).

Recorded because a successor extending the predicate to the
edit-carried or command-backed carriers (REQ-101's named successor)
inherits this seam and its shape.

**Status**: mechanical translation
