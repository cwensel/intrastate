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
