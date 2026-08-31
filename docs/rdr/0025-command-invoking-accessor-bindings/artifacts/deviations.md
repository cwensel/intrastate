# Deviations — RDR 0025 command-invoking-accessor-bindings

Recorded, not asked: this is an unattended run. Each entry is a judgement
taken during Phase 1 that a human might otherwise have been consulted on.

## D1 — Phase 1 lands a declaration-only production surface

**Situation.** The launch prompt asks for tests that COMPILE and fail on
behaviour, and prefers not inventing production API. Two facts collide: the
repo's `pre-commit` hook runs `go vet ./...` and refuses a commit that does
not build, and the RDR names production surface (`table.Accessor.Command`,
the six `CatCommand*` constants, `accessor.ExecError`, `Refusal.Detail`,
`Registry`'s new signature) that does not exist on `main`. A test file
citing that surface does not compile, so it cannot be committed at all — and
the commit cadence is load-bearing for the roborev auto-review.

**Taken.** Declare the spec-named surface with NO behaviour, following the
repo's own precedent: `internal/table/category.go` carries a "PHASE 1
DECLARATION ONLY" block for RDR 0008's `Failure.Offending`/`Remedy`/`Rule`
fields, added in that RDR's Phase 1 for exactly this reason. Every
declaration this run adds carries the same marker naming what Phase 2 owes.

**What is deliberately NOT declared**, so the red gate stays honest:

- the six categories are constants but are NOT appended to
  `table.Categories()` — C5 makes registration the contract, so registering
  them would turn REQ-77/REQ-85 green;
- nothing decodes the carrier fields, nothing raises a `command_*` category,
  nothing populates `Refusal.Detail`, and every `cmdbind` method returns a
  not-implemented `*ExecError`;
- `Registry`'s two new parameters are accepted and ignored.

**Consequence.** 79 new tests, 79 failures, all on behaviour.

## D2 — `accessor.ExecError` is fully implemented, not stubbed

C4 specifies the type completely — `struct { Detail string; Err error }`
with `Error()` and `Unwrap()` — and it is a value type with no seam to be
wrong about. Stubbing it would have made `errors.Is(err, exec.ErrNotFound)`
unassertable at all. So the TYPE is real and REQ-58's test was rewritten to
assert the clause's actual claim: that the wrap survives **the trip to the
refusal site**, which is red until the executor threads it.

## D3 — the command binding package is `internal/cli/cmdbind`

C1 states that "the package or file that holds it" is unconstrained
implementation choice; the Technical Design says "a sibling package to
`flowbind`". `internal/cli/cmdbind` is a sibling of `internal/cli/flowbind`.
No test asserts the package path (REQ-15 forbids it); Phase 2 may move it.

## D4 — `goos` is reached through `export_test.go`

C4 requires "an injectable package-level `goos` var … so the refusal is
testable on Unix CI", and `req-list.md`'s ASSUMPTION for REQ-60/61 reads it
as unexported, with tests in that package setting it directly. The suite is
`package cmdbind_test` (external, matching the house style of every other
RDR suite in this repo), so the var is reached through a
`SetGOOSForTest` helper in `export_test.go` — a file compiled only into the
test binary. The seam therefore stays out of the shipped public surface,
which is what the assumption was protecting.

## D5 — the fixture child is a real built binary

The C4 deadline triple, the C4 env allowlist, and C2's byte-for-byte argv
claim are all properties of a REAL spawn. A faked `exec` seam would make
them unassertable, and mocking the unit under test is barred. So
`internal/cli/cmdbind/fixtures_0025_test.go` builds one small helper binary
once per test binary (via `go build` into a temp dir) whose modes cover every
shape the RDR names, and which echoes its own `os.Args` and `os.Environ`
back so the oracles read what the CHILD observed.

Cost, accepted: the cmdbind suite needs a Go toolchain at test time. CI
already runs `go test`, so the toolchain is present by construction.

## D6 — the MVV skips loudly when `git` is absent

`0025:MVV` step 1 binds "an established tool present in CI", and the RDR's
own A3/A7/A8 spikes bound `git config --file`. Where `git` is not on PATH the
MVV runner `t.Skip`s with the reason spelled out rather than degrading to a
weaker tool, because a substituted tool would no longer be the thing the
spike measured. CI runs `ubuntu-latest`, which carries git.

## D7 — the S4 deadline arms are not `testing.Short`-exempt

They spawn real children and sleep to the bound, so
`TestReq49_TheDeadlineBounds…` skips under `-short`. The default `go test
./...` the RDR names as the validate command does not pass `-short`, so the
arms run in the suite of record.

## D8 — REQ-81's forbidden assertion has no code oracle

REQ-81 says no test MAY assert cross-entry reporting order, and REQ-79 says
none may assert a category count. A prohibition on writing a test cannot
itself be a test. Both are covered by asserting the PROPERTY they protect
(append-only, duplicate-free, additions at the tail; a within-entry
precedence ladder). That a future test does not violate them is a review
obligation, recorded here rather than silently dropped.

## D9 — `flow_exec.go`'s single `Registry` call site was updated

C6 changes `Registry`'s signature, which breaks its one production caller.
The caller was updated to `flowbind.Registry(model, "", false)` so the tree
builds. The `""` and `false` are placeholders: Phase 2 owes the real
threading (`filepath.Dir` of the absolutized `--model` path, and the checked
`--allow-commands` lookup), and REQ-91's test is red against exactly that.

## D10 — `--allow-commands` is not yet registered anywhere

REQ-89 and REQ-91 are red because the flag does not exist. Registering it in
Phase 1 would turn the containment walk green while `buildRequest` still
ignored it — a false green on the load-bearing gate. Phase 3 of the RDR's
own plan lands registration and threading in the SAME step, and the tests
are written to that.

---

# Phase 2 — Implementation

## D5 — REQ-14's fixture change lands on RDR 0002's own suite

**Type: TEST-FIXTURE. Status: mechanical translation.**

`0025:C1` relaxes `accessorTable`'s path rule, so the "path is absent or
empty" case moves out of `malformed_accessor_declaration`. REQ-14 states
the consequence directly: "Existing fixtures asserting the old category on
a path-less entry change category, not verdict." Two cases in
`internal/table/accessors_test.go::TestReq19` (absent path, empty path) now
assert `command_and_path_conflict`. Both still refuse; only the category
moved, which is exactly what the clause licenses.

The conflict rule's two halves are keyed differently, and both are forced
by the Phase 1 suite: "both" is keyed on the KEYS being declared (so
`path = ""` beside a `command` is a conflict — REQ-7's third case), while
"neither" is keyed on neither being a USABLE carrier (so `path = ""` alone
is a conflict, which is where the old arm lands).

## D6 — a brace-bearing argv element is a placeholder attempt only when it is a single word

**Type: SPEC-UNDER. Status: derived from the suite's own discriminating pair.**

C5's placeholder clause and its interpreter clause both claim
`["sh", "-c", "cat {artifact}"]`, and the Phase 1 suite requires opposite
answers on the two shapes:

- REQ-74 (`["sh", "-c", "cat {artifact}"]`) ⇒ `command_shell_interpreter`
- REQ-80 (`["sh", "-c", "{role}"]`) ⇒ `command_unknown_placeholder`,
  because "unknown placeholder outranks shell interpreter"

Under a uniform "any element containing a non-whole `{…}` token is a
placeholder defect" reading, REQ-74's first case reports the placeholder
category and the interpreter clause becomes unreachable behind it. The
discriminator the pair forces is WHITESPACE: `{role}`,
`--file={artifact}`, `{artifact}.bak`, and `{artifact` are all single
words and are placeholder attempts; `cat {artifact}` is a command STRING,
which is the shape clause 4 owns.

Grounded in the RDR's own S1 mutant list, which names the placeholder
mutant as an "unknown or partial `{…}` token" — a token, not a command
line — and in C2's "a placeholder is recognized only as a whole argv
element", whose subject is an element that is trying to BE one.

## D7 — the S5 `<clear>` arm's fixture could not produce its own expected outcome

**Type: TEST-FIXTURE. Status: repaired against the spec's own MVV shape.**

`seam_0025_test.go::TestReq127`'s clear arm drove `<clear>` through
`writeThrough`, whose write tool is `stdin-to-file` and whose reader is
`cat {artifact}` in json mode. That tool STORES the literal rather than
removing it, so the read-back observes the key present holding `<clear>`
— which is `read_back_mismatch`, and is asserted as exactly that by
`TestReq46`'s read-back sibling over the byte-identical fixture. The two
arms cannot both hold on one fixture.

The spec resolves it: S5's clear arm expects "the key reads back absent",
and under C3 a command reader establishes absence ONLY through
`exit_absent` — a json-mode omission is UNREADABLE by contract. So the
arm needs the MVV's own shape: a wrapper that performs the removal, and a
reader that exits a listed code with an empty stdout when the key is
gone (`git config --get`'s shape, spike R6/E1g).

Repaired by giving that ONE arm a `clearThrough` helper carrying both.
The other three arms of `TestReq127` keep `writeThrough` unchanged, and
`TestReq46`'s literal-stored mismatch arm is untouched — the two senses
stay separately asserted, which is what the clause requires.

## D8 — `Registry`'s `baseDir` is the model file's directory, absolutized at the CLI

**Type: IMPL-DECISION. Status: recorded; affects future interpretation.**

A12 is `Pending` in the record and C2 states the shape rather than a
verified reachability. It reaches: `buildRequest` already holds the
selected model PATH (`selectModel` returns it as `ref`), so `baseDir` is
`filepath.Dir` of `filepath.Abs(ref)` with no new `table.Model` field and
no loader change. A12's "If wrong" does not apply.

## D9 — cobra resolves a parent's persistent flag into `Flags()` only after a merge

**Type: TEST-FIXTURE. Status: mechanical translation.**

REQ-131 words S7's containment arm as a structural
`Flags().Lookup("allow-commands") != nil` walk. `cobra.Command.Flags()` is
documented as "the complete FlagSet that applies to this command (local
and persistent declared here and by all parents)", but the parent-pflag
merge is LAZY: it runs on parse, or on the first call to
`InheritedFlags()`/`LocalFlags()`. On a freshly built, unexecuted tree
`Flags().Lookup` therefore misses a group-registered persistent flag —
including on the "fifth verb added inside the group" positive control,
which is added after the tree is built and so can never have been primed.

The only implementation that would satisfy the walk cold is a PER-VERB
registration, which C6 forbids in the same clause ("ONE registration, on
the `flow` GROUP") and which the fifth-verb control exists to catch.

Repaired by calling cobra's own `InheritedFlags()` — the accessor that
performs the merge — immediately before each assertion. The assertion
itself stays on `Flags()`, exactly as REQ-131 words it, and still
discriminates: verified that a verb registered OUTSIDE the group resolves
nothing after the identical call.

## D10 — three Phase 1 fixtures misread the shipped CLI surface

**Type: TEST-FIXTURE. Status: mechanical translation.**

Three fixtures asserted against surfaces that do not exist as written; in
each case the RDR fixes nothing about the detail and the repair is to the
shipped shape.

1. **`lint`'s model-invalid code.** `command_gate_0025_test.go` and
   `command_mvv_0025_test.go` asserted `codeModelInvalid`
   (`"flow-model-invalid"`) on `intrastate lint` invocations. Root `lint`
   has always answered `"model-invalid"`
   (`internal/cli/lint.go:183`); the `flow`-prefixed constant is the
   GROUP's. Named `lintModelInvalid` at the suite's own head.
2. **The lint-mutant substitution.** `writeCommandModel` APPENDS
   `"{artifact}"` to the supplied argv, so the entry reads
   `command = ["true", "{artifact}"]` and a substitution keyed on
   `command = ["true"]` never applied.
3. **The spawn sentinel's argv.** `traceArgv` returned
   `["sh", "-c", …]` and was placed INTO a model by
   `writeCommandModel` — which is the C5 `command_shell_interpreter`
   defect REQ-74 asserts, so the model refused at LOAD and never reached
   the gate the sentinel exists to observe. Its own comment already said
   a model must not carry that form. Repaired to a script FILE declared
   as argv0, the shape the MVV's `traceScript` already takes and the one
   C5 names as the remediation.

## D11 — the MVV wrapper could not write through `git config`

**Type: TEST-FIXTURE. Status: repaired against the A3 spike's own shape.**

`writeWrapper`'s body had two defects that made the MVV's step 3
unreachable:

- it passed the tag key VERBATIM to `git config`, which refuses a key
  carrying no section ("key does not contain a section: status"). The
  model's tag key is `status` and the declared reader reads
  `flow.status`, so mapping intrastate's tag name onto the tool's native
  key space is the WRAPPER's job — exactly the translation the Approach
  charts to the wrapper body as "script the model does not carry";
- its envelope split relied on a sed `\n` replacement (GNU-only; a no-op
  under BSD sed) and left no trailing newline, so `while read` dropped
  the final — and, at one key, only — pair. Both halves silently
  produced a zero-pair loop that exited 0 without writing, which the
  read-back then correctly reported as `read_back_mismatch`.

Repaired with `tr` for the split and `read -r pair || [ -n "$pair" ]` for
the last line, plus the `flow.` section prefix.

## D12 — the MVV payload decoders were written against an invented envelope

**Type: TEST-FIXTURE. Status: mechanical translation.**

`writtenTags` decoded `data.written[].name/value` and `reportedTag`
decoded `data.tags[]`; neither field exists. RDR 0005 fixes the CLI
envelope and this RDR changes no payload shape:
`set-state` reports the read-back-CONFIRMED owned tags on
`data.owned` (`flow_state.go::setStatePayload`), and `read-state` reports
per reader on `data.readers[].tags` (`flow_exec.go::readerOutput`).
`data.owned` is the right oracle for the round-trip arm on its own terms:
it is what read-back verified, and a cleared key is absent from it.

## D13 — REQ-91's success control used a write tool that never applies

**Type: TEST-FIXTURE. Status: mechanical translation.**

The discriminating arm ("with the registration intact the identical
invocation succeeds") declared `command = ["true", "{artifact}"]`, which
exits zero and writes nothing, so read-back correctly refused
`read_back_mismatch` — for a reason with nothing to do with the lookup
the arm is about. Repaired to the suite's own applying wrapper.

## D14 — REQ-65's tail arm used a tool that emits no stderr

**Type: TEST-FIXTURE. Status: mechanical translation.**

The "tail only, no applied sense" arm declared `false` as its reader and
then asserted `CLIError.Detail` is non-empty. `false` exits non-zero in
SILENCE, so there is no stderr tail to carry and the assertion could not
hold under any implementation. S3's own wording puts the tail assertion
on FX-exit-codes E1c — `test --bogus`, exit 2 WITH stderr — so the
fixture takes that shape: a script that writes one stderr line and exits
2.

## D15 — `0025:C4`'s deadline triple is necessary but NOT sufficient

**Type: SPEC-DEFECT. Status: implementation closes the gap; the clause's
sufficiency claim is overstated and a successor should say so.**

C4 states the A1 spike's triple — `Setpgid`, a `Cancel` that signals the
group, and `WaitDelay` — as "necessary AND sufficient", and `0025:F4`
bounds the accepted residue to "a child that IGNORES termination at
deadline". Phase 3a and Phase 3b each reproduced, independently, two
liveness failures that the triple as written does not reach and that
`F4`'s residue does not cover, both in children that ignore nothing:

- A child writing more than one pipe buffer (64 KiB) to stderr before
  closing stdout deadlocks against a parent that drains stdout to
  completion first. `WaitDelay` bounds the stdin write and the POST-kill
  drain, not the pre-kill one, so no member of the triple addresses it.
  Measured pre-fix: 5 000 stderr lines burned the full 8s deadline.
- `cmd.Cancel` fires only when the context ends, so the group is signalled
  only on the FAILURE path. A grandchild backgrounded by a tool that
  exited 0 survives holding the inherited pipes. Measured pre-fix: a
  correct tool's read refused after its full 3s deadline.

**Taken.** Both are closed in the implementation without changing the
declared contract: the triple is kept intact and two properties C4 clearly
intends but does not spell out are added — the two output channels drain
CONCURRENTLY, and the group is reaped after the direct child is waited on,
on the success path as much as the failure path. Neither weakens a
declared bound; the stdout cap, the stderr tail, `WaitDelay`, and the
deadline classification from `ctx.Err()` are unchanged.

The reading is grounded in C4's own stated purpose for the group signal —
"dropping the group signal ORPHANS A GRANDCHILD holding the pipe" — which
is the failure observed on the success path, so honouring it there is
what the clause asks for rather than an extension of it.

**Owed to a successor.** C4's "necessary and sufficient" should either
name the concurrent drain and the post-wait group reap as members of the
mechanism, or drop the sufficiency claim. `F4`'s residue statement should
likewise be narrowed to what it now is: a child that survives a SIGKILL to
its process group.

**No new public surface.** The fix adds `reapGroup` and an `invocation`
field, both unexported; no identifier in the RDR's Normative Contracts
changed shape, so this carries no additive SPEC-UNDER.

## D16 — a signalled child's refusal text is new, unspecified prose

**Type: IMPL-DECISION. Status: settled.**

Distinguishing a signalled child from an exited one (D15's sibling, the
`0025:C3` soundness fix) makes the previous refusal text — "produced no
stdout and exited -1" — a statement about an exit code that does not
exist. C3 fixes the CLASSES and the wire, not refusal prose, and the RDR
names no wording for this case because the case was not contemplated.

**Taken.** Read, gate, and write each say the command "was killed by a
signal before it exited", and the read and gate arms add that the entry's
`exit_absent` / `exit_verdicts` map therefore does not apply. The class is
unchanged — `execution_failure` at the binding, which the executor still
reclassifies to `timeout` when `ctx.Err()` is `DeadlineExceeded` — so no
declared class or category moved; only the free-text diagnosis is new.

## D17 — the process-group syscalls are build-tagged; the refusal is not

**Type: IMPL-DECISION. Status: settled.**

`cmdbind.go` shipped as a single untagged file calling `syscall.Kill` and
setting `syscall.SysProcAttr.Setpgid`. Neither identifier exists on
Windows, so `GOOS=windows go build ./...` failed on three lines, and
because `internal/cli/flowbind` imports the package the whole CLI was
unbuildable there — including the `windows/amd64` target `.goreleaser.yaml`
declares. REQ-134 ("the same model passes `intrastate lint` under that
setting, since lint is platform-neutral") is only SATISFIABLE once the
binary builds on Windows, so this makes C4 reachable rather than changing
it.

**Scope of C4's rejection.** C4 rejects a build constraint because "a build
constraint would make the refusal unbuildable-on-Windows rather than
observable, and lint must stay platform-neutral in the same binary". That
reasoning is about the REFUSAL, and it is honoured: `Unsupported()`, the
injectable `goos` var, the pre-spawn platform arm of the refusal ladder,
and all of lint stay untagged in one platform-neutral binary. Only the
platform MECHANISM — the two syscall wrappers — is split, because the
symbols it names do not exist to compile against off Unix. The unedited
`TestReq60_ACommandEntryRefusesOnARefuseListedPlatformBeforeSpawn`,
including its "lint stays platform-neutral in the same binary" arm, is the
evidence the refusal did not move.

**Taken.** `setProcGroup(*exec.Cmd)` and `killGroup(pid int) error` in
`procgroup_unix.go` (`//go:build !windows && !js && !plan9`) and
`procgroup_other.go` (`//go:build windows || js || plan9`). The non-Unix
`killGroup` returns a non-nil error, so `cmd.Cancel`'s existing fallback to
`cmd.Process.Kill()` is taken rather than a termination being falsely
reported. `reapGroup`'s `p == nil || p.Pid <= 1` guard stays UNTAGGED: it
is contract logic (D15/C4), not platform mechanism. `syscall` remains
imported by `cmdbind.go` for `syscall.WaitStatus`, which is portable.

**The one duplication, and its guard.** The split states the refuse-listed
platform set twice — once as a `//go:build` list, once as
`Unsupported()`'s predicate — and a divergence is silent: a platform
compiling the no-op mechanism but NOT refused would spawn children it
could never terminate as a group. `TestProcGroupBuildTagsMatchTheRefuseList`
asserts the two agree in both directions and that the half actually
compiled matches the running platform's verdict.

**No new public surface.** All four added identifiers are unexported; no
identifier in the RDR's Normative Contracts changed shape, so this carries
no additive SPEC-UNDER.
