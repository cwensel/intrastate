# Coverage — RDR 0027 inline-shell-detection-scope

Phase 1 (launch.md — Tests first). One row per REQ in `req-list.md`. Column 2
is the test that would FAIL if a future change broke that clause. An **EMPTY**
column-2 cell is the orphan mark: that REQ has no test, and the empty cell —
not prose — is how it reads as uncovered.

Test files:

- `internal/table/inline_shell_scope_0027_test.go` — the predicate half
  (`0027:C1` predicate/reads/out-of-scope/report lines, S1–S6).
- `internal/cli/inline_shell_promise_0027_test.go` — the promise half
  (`0027:C1` `promise:` line, S7). A predicate test cannot reach this half.
- `internal/cli/inline_shell_mvv_0027_test.go` — `0027:MVV`, one runnable
  end-to-end test carrying all five steps plus the end state.

Test names below are given without the `Test` prefix and without the
`internal/…` path; subtests are named with `/`.

## Red gate

`go test ./internal/table/ ./internal/cli/` at the time of writing: every
failing top-level test is one of the 33 added here and no pre-existing test
regressed. 21 of the 33 are RED; the other 12 are regression guards that
legitimately pass today (table below). The new-behaviour REQs are RED as
required.

**Genuinely red (new behaviour this record ships)** — 21 tests:
Req42, Req5_AnUnenumeratedPrefix, Req44, Req25 (`i<j` refusal arm), Req6_LowestJ,
Req3 (path-spelled arm), Req13 (own-flag-under-wrapper arm), Req6_LowerListed
(later-pair arm), Req43, Req28, Req58, Req66, Req69 (load-time-refusal arm),
Req49_HelpAllSurface, Req49_NamesBothForms, Req10_ChannelScoped,
Req17_PredicateLine, Req39, Req38, Req20, and MVV steps 2 and 5.

**Already-green regression guards** — 12 tests (and 5 subtests of otherwise-red
tests) assert behaviour that already holds today. They are NOT tautological:
each pins a
property the widen could plausibly break, and the record names most of them as
"unchanged" obligations (REQ-41, REQ-46, REQ-64). Listed explicitly so triage
can tell them from the red set:

| Test | Why it is already green | What it would catch |
| --- | --- | --- |
| `Req2_PositionFreedomIncludesArgv0AndDoesNotExcludeIt` | the four 0025 argv0 probes already refuse | a scan written `for i := 1` — passes every wrapper row and silently un-refuses 0025's set |
| `Req45_TheSevenAdmittedFormsLintGreenAndLoadWithArgvUnchanged` | all seven are green under today's argv0 predicate too | a spellings-scoped out-of-scope reading that refuses `python -` / `sh -es` (`0027:S4`'s named discriminators) |
| `Req10_TheStdinChannelIsAdmittedForEveryListedSpelling` | same | a widen over-applied to the stdin channel |
| `Req11_TheSanctionedWrapperFileFormIsNeverADefect` | same | the remediation the refusal prescribes becoming itself a defect |
| `Req12_TheInterpreterSetStaysOpenUnderThePositionFreeScan` | `perl`/`python3`/`nodejs` are unlisted today | position-freedom over-applied to the NAME list (REQ-70 forbids folding) |
| `Req8_AShellStringCarriedInOneWordIsAdmittedAndNeverSplit` | no word-splitting today | a "helpful" implementation that splits a word on whitespace |
| `Req47_TheAcceptedFalseRefusalClassNamesTheTwoWordsThatCollided` | flag-anywhere already fires at argv0 and the detail already interpolates the pair | narrowing the flag match (REQ-26/REQ-64 forbid it) |
| `Req55_NoInlineShellOptInFieldExistsOnAnEntry` | no `shell` field exists | an opt-in escape hatch (REQ-55/REQ-70) |
| `Req48_TheCategoryWireStringAndItsRelativePositionAreUnchanged` | `category.go` is untouched | the edit straying out of `load.go`, or a rename (D-naming) |
| `Req14_TheClause3ExemptionStillRequiresTheInterpreterForm` | clause 3 already keys on `isInterp` | a widen mis-implemented as "any brace-bearing whitespace element is exempt" |
| `Req21_TheShippedWordingNeverClaimsInlineShellIsClosed` | no description ships yet, so no overclaim exists | the premortem's second failure once Phase 3 lands |
| `Req27_TheOuterInterpreterWordWinsTheReportedForm` | `["python","sh","-c",…]` already refuses at argv0 reporting `python -c` | a tie-break that iterated `shellInterpreters` rather than argv positions (REQ-57) |
| `Req25/flag_before_interpreter_is_green`, `Req13/other-interpreter-flag`, `Req3/no-folding`, `Req6_LowerListed/bare-word`, `Req69/no-stdin-code` | the corresponding negative arms already hold | each is the NEGATIVE control that makes its red sibling discriminating |

## Coverage table

| REQ | Test |
| --- | --- |
| REQ-0 | Req48_TheCategoryWireStringAndItsRelativePositionAreUnchanged |
| REQ-0a | Req48_TheCategoryWireStringAndItsRelativePositionAreUnchanged |
| REQ-0b | TestReq0b_C5CommentsInLoadGoReciteThisRecordsSuccessorClause |
| REQ-1 | Req42_EveryNamedWrapperStillRefusesTheInterpreterForm |
| REQ-2 | Req2_PositionFreedomIncludesArgv0AndDoesNotExcludeIt; Req25_TheFlagMustFollowTheInterpreterWordNotPrecedeIt |
| REQ-3 | Req3_BasenameIsThePostSlashSegmentMatchedByteExactly |
| REQ-4 | Req5_AnUnenumeratedPrefixIsRefusedBecauseNothingBeforeTheInterpreterIsRead |
| REQ-5 | Req5_AnUnenumeratedPrefixIsRefusedBecauseNothingBeforeTheInterpreterIsRead |
| REQ-6 | Req6_WithSeveralQualifyingFlagsTheLowestJIsReported; Req6_ALowerListedInterpreterWithNoFollowingFlagDoesNotStopTheScan; Req27_TheOuterInterpreterWordWinsTheReportedForm |
| REQ-7 | Req8_AShellStringCarriedInOneWordIsAdmittedAndNeverSplit; Req42_EveryNamedWrapperStillRefusesTheInterpreterForm |
| REQ-8 | Req8_AShellStringCarriedInOneWordIsAdmittedAndNeverSplit |
| REQ-9 | Req45_TheSevenAdmittedFormsLintGreenAndLoadWithArgvUnchanged |
| REQ-10 | Req10_TheStdinChannelIsAdmittedForEveryListedSpelling; Req10_TheDescriptionScopesTheStdinFormByChannelNotByShellSpelling |
| REQ-11 | Req11_TheSanctionedWrapperFileFormIsNeverADefect |
| REQ-12 | Req12_TheInterpreterSetStaysOpenUnderThePositionFreeScan |
| REQ-13 | Req13_TheFlagListIsKeyedByTheInterpreterFoundNotTheUnion |
| REQ-14 | Req43_ABraceBearingStringUnderAWrapperReportsTheInterpreterDefect; Req14_TheClause3ExemptionStillRequiresTheInterpreterForm |
| REQ-15 | Req14_TheClause3ExemptionStillRequiresTheInterpreterForm |
| REQ-16 | Req47_TheAcceptedFalseRefusalClassNamesTheTwoWordsThatCollided; Req42_EveryNamedWrapperStillRefusesTheInterpreterForm; Req48_TheCategoryWireStringAndItsRelativePositionAreUnchanged |
| REQ-17 | Req17_TheDescriptionAlsoCarriesThePredicateLine |
| REQ-18 | Req38_TheDescriptionReachesTheCLIReferenceMirrorByTheSameDerivation |
| REQ-19 | Req49_TheDescriptionNamesBothOutOfScopeFormsByName; Req17_TheDescriptionAlsoCarriesThePredicateLine |
| REQ-20 | Req20_TheAdmittedFormDisclosureLivesOnTheHelpSurfaceNotInARefusal |
| REQ-21 | Req21_TheShippedWordingNeverClaimsInlineShellIsClosed |
| REQ-22 | Req58_ThePredicateIsPureAndCarriesNoStateAcrossLoads |
| REQ-23 | Req6_WithSeveralQualifyingFlagsTheLowestJIsReported; Req3_BasenameIsThePostSlashSegmentMatchedByteExactly; Req27_TheOuterInterpreterWordWinsTheReportedForm |
| REQ-24 | Req48_TheCategoryWireStringAndItsRelativePositionAreUnchanged |
| REQ-25 | Req25_TheFlagMustFollowTheInterpreterWordNotPrecedeIt; Req27_TheOuterInterpreterWordWinsTheReportedForm |
| REQ-26 | Req47_TheAcceptedFalseRefusalClassNamesTheTwoWordsThatCollided |
| REQ-27 | Req27_TheOuterInterpreterWordWinsTheReportedForm |
| REQ-28 | Req28_TheNewFalseRefusalClassIsAcceptedNotSuppressed |
| REQ-29 | Req42_EveryNamedWrapperStillRefusesTheInterpreterForm; Req6_ALowerListedInterpreterWithNoFollowingFlagDoesNotStopTheScan; Req3_BasenameIsThePostSlashSegmentMatchedByteExactly |
| REQ-30 | Req8_AShellStringCarriedInOneWordIsAdmittedAndNeverSplit; Req5_AnUnenumeratedPrefixIsRefusedBecauseNothingBeforeTheInterpreterIsRead |
| REQ-31 | Req44_TheDeletedEnvWalkIsSubsumedByThePositionFreeScan |
| REQ-33 | Req48_TheCategoryWireStringAndItsRelativePositionAreUnchanged |
| REQ-MVV | ReqMVV0027_TheWrapperClassRefusesTheAdmittedFormsLoadAndThePromiseShips |
| REQ-34 | Req42_EveryNamedWrapperStillRefusesTheInterpreterForm; Req44_TheDeletedEnvWalkIsSubsumedByThePositionFreeScan; Req48_TheCategoryWireStringAndItsRelativePositionAreUnchanged |
| REQ-35 | TestReq35_TheFourExistingReq74ProbesSurviveUnweakened |
| REQ-36 | Req49_TheHelpAllSurfaceDescribesTheShellInterpreterCategory; Req49_TheDescriptionNamesBothOutOfScopeFormsByName |
| REQ-37 | Req49_TheHelpAllSurfaceDescribesTheShellInterpreterCategory; Req20_TheAdmittedFormDisclosureLivesOnTheHelpSurfaceNotInARefusal |
| REQ-38 | Req38_TheDescriptionReachesTheCLIReferenceMirrorByTheSameDerivation |
| REQ-39 | Req39_TheDescriptionRendersOffTheRespondGatewayWithNoDefectPresent |
| REQ-40 | TestReq40_TheQ2q1DispositionIsRecordedInDeviations |
| REQ-41 | Req2_PositionFreedomIncludesArgv0AndDoesNotExcludeIt; ReqMVV0027/step_4 |
| REQ-42 | Req42_EveryNamedWrapperStillRefusesTheInterpreterForm |
| REQ-43 | Req43_ABraceBearingStringUnderAWrapperReportsTheInterpreterDefect |
| REQ-44 | Req44_TheDeletedEnvWalkIsSubsumedByThePositionFreeScan |
| REQ-45 | Req45_TheSevenAdmittedFormsLintGreenAndLoadWithArgvUnchanged |
| REQ-46 | Req45_TheSevenAdmittedFormsLintGreenAndLoadWithArgvUnchanged; Req10_TheStdinChannelIsAdmittedForEveryListedSpelling |
| REQ-47 | Req47_TheAcceptedFalseRefusalClassNamesTheTwoWordsThatCollided |
| REQ-48 | Req48_TheCategoryWireStringAndItsRelativePositionAreUnchanged; ReqMVV0027/step_4 |
| REQ-49 | Req49_TheHelpAllSurfaceDescribesTheShellInterpreterCategory; Req49_TheDescriptionNamesBothOutOfScopeFormsByName; Req10_TheDescriptionScopesTheStdinFormByChannelNotByShellSpelling |
| REQ-50 | Req20_TheAdmittedFormDisclosureLivesOnTheHelpSurfaceNotInARefusal |
| REQ-51 | Req42_EveryNamedWrapperStillRefusesTheInterpreterForm |
| REQ-52 | Req8_AShellStringCarriedInOneWordIsAdmittedAndNeverSplit |
| REQ-53 | Req45_TheSevenAdmittedFormsLintGreenAndLoadWithArgvUnchanged; Req10_TheStdinChannelIsAdmittedForEveryListedSpelling |
| REQ-54 | Req12_TheInterpreterSetStaysOpenUnderThePositionFreeScan |
| REQ-55 | Req55_NoInlineShellOptInFieldExistsOnAnEntry |
| REQ-56 | Req3_BasenameIsThePostSlashSegmentMatchedByteExactly |
| REQ-57 | Req27_TheOuterInterpreterWordWinsTheReportedForm |
| REQ-58 | Req58_ThePredicateIsPureAndCarriesNoStateAcrossLoads |
| REQ-59 | Req58_ThePredicateIsPureAndCarriesNoStateAcrossLoads |
| REQ-60 | Req69_TheCategoryStaysALoadTimeRefusalAndTheStdinAxisShipsNoCode |
| REQ-61 | Req42_EveryNamedWrapperStillRefusesTheInterpreterForm; Req47_TheAcceptedFalseRefusalClassNamesTheTwoWordsThatCollided |
| REQ-62 | Req45_TheSevenAdmittedFormsLintGreenAndLoadWithArgvUnchanged; Req8_AShellStringCarriedInOneWordIsAdmittedAndNeverSplit |
| REQ-63 | Req58_ThePredicateIsPureAndCarriesNoStateAcrossLoads |
| REQ-64 | Req47_TheAcceptedFalseRefusalClassNamesTheTwoWordsThatCollided |
| REQ-65 | Req49_TheDescriptionNamesBothOutOfScopeFormsByName; Req45_TheSevenAdmittedFormsLintGreenAndLoadWithArgvUnchanged |
| REQ-66 | Req66_ALongArgvIsScannedToTheEndWithNoBound |
| REQ-67 | Req42_EveryNamedWrapperStillRefusesTheInterpreterForm |
| REQ-68 | Req5_AnUnenumeratedPrefixIsRefusedBecauseNothingBeforeTheInterpreterIsRead |
| REQ-69 | Req69_TheCategoryStaysALoadTimeRefusalAndTheStdinAxisShipsNoCode |
| REQ-70 | Req55_NoInlineShellOptInFieldExistsOnAnEntry; Req12_TheInterpreterSetStaysOpenUnderThePositionFreeScan; Req3_BasenameIsThePostSlashSegmentMatchedByteExactly |

## Orphan REQs

Phase 1 recorded five empty cells. Phase 3c resolved all of them, leaving
no orphan rows:

- REQ-0b, REQ-35, REQ-40 constrain the SHAPE OF THE DIFF, which is not the
  same as unobservable — each names a fact about the shipped tree a test can
  read directly, so `internal/table/inline_shell_authoring_0027_test.go`
  asserts them (all three mutation-verified).
- REQ-32 was WITHDRAWN: it is a caveat on an illustrative code block, not a
  testable clause, and was extracted in error.
- REQ-33 was never orphaned; it is partially testable and its testable
  consequence is covered.

### REQ-32 — withdrawn, not orphaned

Five REQs carry an empty column-2 cell. Each is a process or authoring
obligation with no runtime behaviour a test can observe; recording them as
uncovered is more honest than pointing a row at a test that does not in fact
constrain them.

- **REQ-0b** — CLOSED in Phase 3c. The clause is a fact about the shipped
  tree, so a test CAN read it:
  `TestReq0b_C5CommentsInLoadGoReciteThisRecordsSuccessorClause` asserts
  `load.go` cites `0027:C1` and that no comment still describes the live
  predicate as argv0-anchored. Mutation-verified (stripping the citation
  fails it).
- **REQ-32** — WITHDRAWN in Phase 3c, so it has no row at all rather than an
  empty one. "Illustrative — shape only." is a caveat on a code block, not a
  normative clause: `rdr inspect` confirms its line falls inside no labelled
  contract element. A clause that obliges nobody to do anything is not a
  testable REQ, and an empty coverage cell would misreport a classification
  error as a coverage gap. See deviation D8.
- **REQ-35** — CLOSED in Phase 3c.
  `TestReq35_TheFourExistingReq74ProbesSurviveUnweakened` asserts the REQ-74
  test still exists and still drives all four probe vectors verbatim, which is
  what "keep the four existing probes as-is" protects: the widening must not
  be paid for by loosening the predecessor's probes. Mutation-verified
  (respelling one probe fails it).
- **REQ-40** — CLOSED in Phase 3c. The prerequisite is discharged by
  RECORDING the disposition, and the record is a readable artifact:
  `TestReq40_TheQ2q1DispositionIsRecordedInDeviations` asserts
  `deviations.md` names q2q1 and names which way it went. The behavioural
  half (the `env` walk is gone) stays covered by
  `Req44_TheDeletedEnvWalkIsSubsumedByThePositionFreeScan`.
  Mutation-verified.
- **REQ-33** is NOT orphaned but is only partially testable: "the edit is
  confined to `load.go`" is a diff fact. What a test can hold is the
  consequence — `Categories()` unchanged — which `Req48_…` asserts.

## Readings taken (unattended override; recorded per the launch prompt)

**Q1 — the Phase 3 description surface is PER-CATEGORY OPT-IN, not total over
`Categories()`.** No test here asserts that every member of `Categories()`
carries description text. Grounds: REQ-48 forbids changing `Categories()`, and
a totality assertion would couple this record to every future append,
contradicting `0025:REQ-79`'s "the list's total size is not a contract at any
point"; C1 ships text for exactly ONE category and `0027:S7` asserts exactly
that one. This is deviation **D2**'s named disposition and is recorded there.

**Q2 — "THESE words … verbatim" is SUBSTANTIVE FIDELITY, not a byte-equal
golden.** `Req49_TheDescriptionNamesBothOutOfScopeFormsByName`,
`Req10_TheDescriptionScopesTheStdinFormByChannelNotByShellSpelling` and
`Req17_TheDescriptionAlsoCarriesThePredicateLine` assert that the text NAMES
both out-of-scope forms, is channel-scoped, and states the position-free
predicate — never byte equality against a golden string. Grounds: a byte-equal
copy of C1's fence would ship contract syntax (`base(argv[i])`, `i < j`,
peer-record ids) to a CLI reader, and `0027:S7`'s own oracle is "its text names
both out-of-scope forms … in C1's channel-scoped words" — a naming assertion.
Each assertion accepts several defensible spellings, so the wording stays an
authoring choice while the two facts a reviewer needs stay pinned.

**D1 (pre-seeded) — `Categories()` order is asserted RELATIVELY.**
`Req48_TheCategoryWireStringAndItsRelativePositionAreUnchanged` and
`MVV/step_4` index each of the six C5 categories in `Categories()` and require
the sequence to ASCEND. There is no tail-slice equality and no `len()`
assertion anywhere in the added tests, so a peer record appending members at
the tail leaves both green while a move, drop, or reorder still fails. This is
D1's named disposition.

**REQ-42 ASSUMPTION — scenarios are driven through `table.Load`,** on the 0025
command-carrier fixture, never by calling the unexported `interpreterForm`
directly, so the oracle is the category a model author sees. The `internal/cli`
MVV additionally drives step 2 and step 3 through `intrastate lint` for the
same reason one level up.

**REQ-38/REQ-36 ASSUMPTION — no Go symbol is pinned.** The promise-half tests
read the RENDERED `--help-all` body and the generated `docs/cli-reference.md`
text; none names the accessor function. The record fixes the SURFACE and the
TEXT, not the symbol, consistent with `0025:REQ-15`.

## Notes on two subtests that pass today for a spec-relevant reason

- `Req44/--unset=FOO` and `Req44/--chdir=DIR` are already green. Today's `env`
  walk skips any token containing `=`, so these two of S3's eight forms already
  refuse — an accident of the walk, not a property of it. They are kept in the
  table because `0027:S3` names all eight and because they must keep refusing
  once the walk is deleted.
- The MVV's steps 1, 3 (both arms) and 4 pass today by design: the record's own
  step 5 says "steps 1–4 all pass with no description at all", and steps 1/3/4
  assert what must NOT move. Steps 2 and 5 are the record's two deltas and both
  are red.

## REQ-MVV output

Recorded from the built binary and the runnable test after Phase 2
implementation landed. All five steps plus the exec arm and the end state PASS.

**Step 1** — the baseline command-carrier model with a green `command` read
binding lints clean, and the declared argv survives onto the loaded binding:

```
$ intrastate lint --model mvv0027-base.toml --as=json
exit 0  {"type":"ok","data":{"findings":[]}}
loaded argv == ["git","config","--file","{artifact}","--get","flow.status"]
```

**Step 2** — every wrapper mutant refuses with `command_shell_interpreter`,
the detail naming the matched words `sh -c` and the unchanged script
remediation. The `{artifact}` mutant does NOT report
`command_unknown_placeholder`:

```
$ intrastate lint --model nice.toml --as=json          # exit 2
{"code":"model-invalid","message":"model does not conform to the transition-model schema",
 "param":"model",
 "detail":"command_shell_interpreter: read state declares the interpreter form sh -c; inline shell is not a declared command; put it in a script and declare the script as argv0",
 "findings":[{"code":"command_shell_interpreter","message":"command_shell_interpreter: read state declares the interpreter form sh -c; inline shell is not a declared command; put it in a script and declare the script as argv0","locator":"nice.toml:1"}]}
```

Identical category, detail and remediation for every named mutant — the wrapper
prefix changes nothing, because nothing before the interpreter word is read:

```
["env","-i","sh","-c","cat flow.status"]      -> exit 2  command_shell_interpreter  form="sh -c"
["env","-u","FOO","sh","-c","cat flow.status"] -> exit 2  command_shell_interpreter  form="sh -c"
["nice","sh","-c","cat flow.status"]          -> exit 2  command_shell_interpreter  form="sh -c"
["timeout","5","sh","-c","cat flow.status"]   -> exit 2  command_shell_interpreter  form="sh -c"
["xargs","sh","-c","cat flow.status"]         -> exit 2  command_shell_interpreter  form="sh -c"
["doas","sh","-c","cat flow.status"]          -> exit 2  command_shell_interpreter  form="sh -c"
["nice","sh","-c","cat {artifact}"]           -> exit 2  command_shell_interpreter  form="sh -c"
```

The last row is the one that moved: it reported `command_unknown_placeholder`
before this build, and clause 3's exemption now keys on the widened predicate.

**Step 3** — every admitted form lints green AND loads with its argv unchanged
(element-wise equal to the declared vector, not merely no error returned):

```
["sh","./gate.sh"]          -> exit 0  {"type":"ok","data":{"findings":[]}}   argv unchanged
["env","-S","sh -c echo"]   -> exit 0  {"type":"ok","data":{"findings":[]}}   argv unchanged
["sh","-s"]                 -> exit 0  {"type":"ok","data":{"findings":[]}}   argv unchanged
["python","-"]              -> exit 0  {"type":"ok","data":{"findings":[]}}   argv unchanged
```

`["python","-"]` is the probe that C1's out-of-scope line is stated over the
stdin CHANNEL and not over shell spellings — `python` IS a `shellInterpreters`
member, so a spellings-scoped reading would refuse it. It is green and is
documented as admitted in the step-5 text below.

**Step 3, exec arm** — BOTH of "the first two" are shown by execution, not by
lint-green. `["sh", "<gate.sh>"]` lints clean and the declared wrapper FILE
actually runs under `--allow-commands`, writing the value the model asks for
(`flow.status = final`) into the artifact. Value-for-value, not a zero exit.

The `["env","-S","<one word>"]` half RUNS on this host — it is not skipped. A4's
platform note ("darwin's stock BSD `env` rejects `-S`") is a non-contractual
observation that does not hold for this machine's `/usr/bin/env`, which accepts
`-S`. The arm is therefore gated on a runtime CAPABILITY PROBE of the resolved
`env` binary rather than on `runtime.GOOS`, so it runs wherever `-S` is present
and skips with a reason only where it genuinely is not. Observed:

```
$ /usr/bin/env -S "A=1"                                 # capability probe: exit 0
$ intrastate lint --model exec-env-s.toml --as=json
exit 0  {"type":"ok","data":{"findings":[]}}
$ intrastate flow set-state --model exec-env-s.toml \
    --artifact state=<tmp>/state.cfg --write status=final --allow-commands --as=json
exit 0
$ cat <tmp>/state.cfg
[flow]
	status = final
```

The declared word is ``sh -c 'git config --file "$0" flow.status `printf
final`'`` with `{artifact}` as a TRAILING argv element, and the oracle is the
recovered value `final`.

The oracle has to separate "a shell ran the word" from "`env -S` split the word
and ran a non-shell binary", because `env -S` is ITSELF a splitter: it performs
quote removal and `${VAR}` expansion on the word before exec. Quote-based
tricks are therefore not discriminating — `env -S` alone concatenates
`"fi""nal"` into `final` with no shell present, so an oracle keyed on that
would pass against a word with the `sh -c` prefix stripped. Two constructs are
discriminating, and both are used: backtick COMMAND SUBSTITUTION (`env -S`
rejects `$(…)` outright and passes backticks through as literal bytes) and the
POSITIONAL PARAMETER `$0` (`env -S` supports only `${VARNAME}` and rejects a
bare `$0`). A shell performs both. Verified by mutation — stripping the `sh -c`
prefix makes `env -S` refuse the word outright and fails the assertion.
Recovering `final` FROM THE PATH BOUND TO `$0` therefore proves the one word
reached a SHELL on two independent counts, rather than merely that some process
ran and exited zero.

The path is passed as a trailing argv element rather than embedded in the word.
`env -S` splits ONLY the `-S` string; every element after it is preserved
verbatim and handed to the utility, so `sh -c <prog> <path>` binds the path to
`$0` inside the program. This keeps `{artifact}` a WHOLE element, which is the
only form `0025:C2` substitutes (a non-whole-element `{…}` token is itself a
refusal, `command_unknown_placeholder`), and it means the path is never text in
either parser — it crosses neither `env -S`'s splitter nor `sh`'s tokenizer. No
quoting of it is required and no `TMPDIR` can corrupt the program: the arm runs
green under paths carrying spaces, tabs, non-ASCII, and quotes, backticks and
`$`. Earlier revisions embedded the path literally and guarded it with a
character precondition; the positional form removes both the quoting and the
skip, so the arm now has no host-path precondition at all.

The `-S` capability probe names NO PATH: it is `env -S "A=1"`, a bare
environment assignment that exercises `-S` parsing and exits 0 without
resolving or executing any utility. The probe must isolate "this `env` has no
`-S`" from every other reason a command can fail, and both obvious
alternatives fail that test — `-S true` conflates it with a missing or
shadowed `true`, and `-S <envPath>` reintroduces the splitter's own whitespace
rule against the resolved path, so an `env` under a directory with a space
would report failure on a host that DOES support `-S`. Either conflation lets
this required execution coverage skip itself. An `env` without `-S` rejects
the flag outright ("illegal option"), a non-zero exit.

Since the form is executable here, the earlier record in this file that the
half "skips on this host" was false and is corrected by this paragraph.

**Step 4** — the four original 0025 REQ-74 probes still refuse with the SAME
category, and the reported form still names the right pair. Position-freedom
widens the set and never excludes argv0:

```
["sh","-c","cat {artifact}"]  -> exit 2  command_shell_interpreter  form="sh -c"
["bash","-c","echo hi"]       -> exit 2  command_shell_interpreter  form="bash -c"
["python","-c","print(1)"]    -> exit 2  command_shell_interpreter  form="python -c"
["env","sh","-c","echo hi"]   -> exit 2  command_shell_interpreter  form="sh -c"
```

`table.Categories()` order holds: the six C5 categories index in ascending
relative order (deviation D1's reading — no tail or `len` coupling). The four
probes and `neg/neg-command-shell-interpreter.toml` are byte-unedited; the whole
shipped 0025 suite is green.

**Step 5** — the description ships and is read with NO model loaded and no
defect provoked, naming both out-of-scope forms in C1's channel-scoped words:

```
$ intrastate lint --help-all

Load category command_shell_interpreter — what it promises:

Refused: a listed interpreter word followed, at ANY later argv
position, by one of that interpreter's own inline-code flags — under any
prefix (env and its options, nice, timeout, xargs, doas, and wrappers
nobody enumerated). Nothing before the interpreter word is read, so no
wrapper table exists and none is consulted.

The check reads argv WORDS only. It never splits a word on whitespace,
and never reads stdin, files, PATH, or the resolved binary. The
interpreter set is an OPEN deny-list, so an unlisted spelling (python3,
nodejs, busybox) is admitted.

Out of scope, BY NAME — admitted by lint, and an interpreter may still
run:

  a shell string carried in ONE word, such as env -S "sh -c …", or a
  single "sh -c …" element handed to a tool that re-splits it;

  an interpreter that reads its script from STDIN — sh -s, bare sh,
  sh -es, python -, node -. The channel is the scope: any listed
  interpreter taking its code on stdin rather than as a later argv word
  is admitted, however spelled.

sh script.sh is the sanctioned wrapper-file form and never a defect.
```

The same text reaches `docs/cli-reference.md` by the same derivation
(`internal/cli/docs.go::runDocs`), and `make check`'s staleness gate is green
against the regenerated file.

**End state** — "every wrapper mutant red with the right defect, every admitted
form green and documented, the promise text shipped and asserted, no existing
probe changed": all four hold. `go test ./...` is clean across the repository.

## REQ-MVV runner

```sh
go test ./internal/cli/ -run TestReqMVV0027_TheWrapperClassRefusesTheAdmittedFormsLoadAndThePromiseShips -v
```
