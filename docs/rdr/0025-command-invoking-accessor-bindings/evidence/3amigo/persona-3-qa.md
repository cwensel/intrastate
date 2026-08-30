Model: claude-opus-5[1m]

# Persona 3 — QA / Tester

Question: how do I test this? What are the pass/fail criteria?

Owned starting set: `S1`–`S7` (incl. `S4b`, `S5b`), `MVV`, `C1`–`C6`.
**Widened** to `§failure-modes` (F5, F7), `§critical-assumptions` (A4, A10),
`§load-bearing-decisions`, `§prerequisites`, and `§pre-lock-mini-checks`.
What sent me: three of my findings are *silences* — a contract clause with no
scenario, or a scenario whose oracle lives in a section outside the contract
that names it. A silence has no line range inside the owned set, so the owned
set could not answer the question.

Code grounding used: `internal/accessor/model.go`, `internal/accessor/executor.go`,
`internal/accessor/binding.go`, `internal/cli/flowbind/registry.go`,
`internal/table/category.go`, `.github/workflows/ci.yml`, and the two spike
`output.txt` fixtures under `evidence/spikes/`.

---

## High

### H1 — `0025:S3` cites FX-exit-codes as its oracle, but the fixture covers none of the ordering cases the scenario exists to prove

S3 enumerates eight cases and eight expected outcomes and declares "the
verdict/error code split follows normative fixture **FX-exit-codes**
(`evidence/spikes/a3-a7-a8-real-tools/output.txt` E1a–E1h, R6)".

I read that fixture. **Every one of E1a–E1h and R6 has `stdout=""`.** The
fixture is exclusively an empty-stdout exit-code census. But S3's first three
cases — the ones that make the scenario *about* C4's normative ordering — all
have non-empty stdout:

- "a gate command that writes a well-formed deny envelope and exits non-zero"
- "one that writes a malformed envelope and exits zero"
- "a 1 MiB + 1 byte stdout"

C3 states the discriminating rule: "A non-empty stdout is always parsed first
(C4 ordering); the exit maps apply only to an empty stdout." So the cited
fixture is, by construction, a census of the branch that C4's ordering rule
does **not** govern. The eight-element Expected list is also positional
("respectively") against a case list of a different cardinality and different
membership than E1a–E1h, so the reader cannot even align case *n* to fixture
row *n*.

**Test this prevents**: the envelope-before-exit ordering test — the single
highest-value test in the record, the one premortem P-12 exists for. I can
write the empty-stdout half from the fixture. For the non-empty-stdout half I
have no fixture value, and "normative fixture" tells me not to invent one.
Concretely: for the deny-envelope-with-exit-1 case I cannot tell whether the
oracle is `Verdict == deny` with `Refusal == nil`, or a deny that still
carries a populated `Detail` stderr tail; and for the 1 MiB + 1 case I cannot
tell whether the cap is applied to the *captured prefix* (so a well-formed
envelope inside the first 1 MiB still parses, which C4's stated ordering would
require) or to the stream as a whole before any parse (which the "overflow =
execution_failure" wording implies). Those two readings disagree on a real
input, and C4's ordering sentence and C4's bounds sentence each support one.

**Decision blocked**: whether C4's parse-before-classify ordering applies to a
truncated-at-cap stdout, and what evidence supplies the non-empty-stdout
oracle values now that the named fixture does not.

### H2 — `0025:C6` is normative but `0025:A10`, the assumption it rests on, is `Status: Pending`; `0025:S7`'s config arms have no fixture

A10's own **If wrong** line says it: "C6 is unimplementable as written and
S7's config arms have nothing to exercise." A10's Evidence says only the
*absence* is established, and I confirmed it — there is no
`internal/cli/config/` in the repo and no user-scope TOML reader. What is
unfixed is A10's "positive call": search order, precedence against
`--allow-commands`, and absent/malformed behaviour.

S7 is written across that fault line. Its first three arms (gate unset, flag
alone, lint under both) are testable today. Its remaining text — "Once A10
fixes the config surface, the same scenario adds `allow_commands = true` in
config, and the flag-vs-config precedence case A10 decides" — is a scenario
that names a decision that has not been made.

Worse, the last sentence of S7's Expected — "**A malformed config must refuse,
never fail open**" — is a *normative requirement stated only inside a
scenario*. It appears nowhere in C6. C6's normative block covers absent-or-
false only ("absent or false ⇒ every command invocation refuses"); a malformed
file is neither absent nor false. A requirement that lives only in a test's
Expected line is a requirement the implementation is not bound by.

**Test this prevents**: the config-precedence test and the malformed-config
fail-closed test. For precedence I have no rule at all. For malformed I have a
requirement with no contract clause and no defined observable — "refuse" does
not say whether it is a load-time `table.Category` (the C5 family), a
pre-spawn `execution_failure` (the C6 family), or a CLI-level startup error,
and those three surface at different times with different diagnostics.

**Decision blocked**: whether v1 ships the config file at all; if it does, the
precedence rule and the malformed-config refusal *class*. Note S7 is
internally split on this — its first arms assume flag-only v1, its later text
assumes the file exists.

### H3 — `0025:C3`'s "an omitted key is not-carried" is ambiguous across a distinction the code makes normative, and JSON mode — the default — has no scenario at all

C3's read row reads: `stdout (read, output = "json", default): flat JSON
object of strings; an omitted key is not-carried`.

`internal/accessor/binding.go` makes "not-carried" two different things:

    // Read ... returns one KeyValue per key it could read (present or
    // established-absent) and names in unreadable every requested key it
    // could NOT read.

and `KeyValue.Absent` is documented there as "a VALUE ... unreadability is
signalled by the binding refusing the key instead". The two paths diverge in
`executor.go::classify` (lines 133–147): an established-absent key returns a
value; a key the binding did not return at all lands in `unread` and mints
`ClassIncompleteRead`. Same input JSON, two different refusal outcomes,
depending on a mapping C3 does not state.

This matters because C3 *does* state the rule for the other two read shapes.
`exit_absent` is explicit ("establishes every declared key absent" ⇒
`Absent: true`) and raw mode is explicit. Only the default JSON shape is left
to inference.

Compounding it: **no scenario exercises JSON read mode.** S3 covers gate
envelopes plus `exit_absent`/raw absence; S4b is raw mode; S5 is the write
stdin envelope; S6/MVV bind `git config --get` in **raw** mode. The default
`output = "json"` read path — the one every non-established-tool wrapper will
use — has neither a scenario nor a worked fixture.

**Test this prevents**: the JSON-read envelope test. I cannot write a single
assertion for `{"a":"1"}` against `keys = ["a","b"]` — the answer is either
`ReadResult{Values: [a=1], b absent}` or `Refusal{Class: incomplete_read,
Keys: ["b"]}`, and C3's four words support both. This is exactly the
"pass/fail criteria stated as adjectives" failure: "not-carried" is a
description, not an assertion.

**Decision blocked**: whether an omitted key in a JSON read envelope is
established-absent or unreadable.

---

## Medium

### M1 — `0025:F7`'s Windows refusal is a normative behaviour with no scenario and no runner

F7: "a command entry invoked on Windows refuses `execution_failure` naming the
unsupported platform, before spawn ... lint stays platform-neutral, so a model
authored on Windows still validates." Two distinct assertions — a runtime
refusal and a lint pass — on a platform the project does not build for in CI.
`.github/workflows/ci.yml` is `runs-on: ubuntu-latest` at every one of its six
jobs. The repo has exactly one `runtime.GOOS` reference in non-test-adjacent
code (`internal/cli/flow_mvv_0023_test.go:913`).

No S covers F7. C4's `deadline:` row bakes in `Setpgid` and `Cancel` to
`-pgid`, which is why the platform bound exists, but no contract clause states
the refusal — F7 is the only home for a behaviour that gates an entire
platform.

**Test this prevents**: the platform-refusal test. Without a stated seam I
cannot tell whether to test it as a `runtime.GOOS`-switched unit (writable
today on Linux via injection, if the implementation is asked to make the
platform check injectable) or as a build-tagged `_windows_test.go` that CI
never runs (i.e. never actually verified). Those are different implementation
obligations and only one of them is provable on this CI.

**Decision blocked**: whether the platform gate is a testable injected
predicate or a build-constraint, and whether F7's claim gets promoted to a
contract clause so it has a home an S can anchor to.

### M2 — `0025:S5b`'s oracle depends on an invariant split across two files, and the record misattributes where it lives

S5b's Expected is well-formed as an oracle — it correctly bans the weak one
("asserted **by the selected reader's identity**, never by 'read-back
succeeded'"). But its stated mechanism is "`registry.go`'s name-sorted first
match", and A4 grounds it at `internal/accessor/model.go::readerFor`.

I read both. `readerFor` (`internal/accessor/model.go:231`) does a plain
**slice-order** scan:

    for _, d := range reg.Definitions {
        if d.Identity.Capability == CapRead && d.Accessor.Role == role {

Nothing there sorts. The name-sorting happens in a different package, in
`internal/cli/flowbind/registry.go:26`, where readers are appended under
`slices.Sorted(keys(m.Readers))`. So "name-sorted first match" is an
**emergent** property of two files that no single assertion covers, and any
future registry builder (this RDR adds one — the command binding constructor,
per the Selection LBD) that appends without sorting silently changes S5b's
answer with no test failing at the site of the change.

**Test this prevents**: a stable pin on today's inherited first-match
behaviour. I can write S5b as stated and it will pass — but it will keep
passing while the property it pins quietly moves, because the assertion sits
on `readerFor`'s output while the ordering is established in `flowbind`. This
is a time-shifted failure: S5b's stated purpose is "so 0016's landing is a
visible change, not a silent one", and as anchored it cannot deliver that.

**Decision blocked**: whether the name-sort is a contract of registry
construction (and therefore gets its own assertion in the new command-binding
constructor) or an incidental artifact S5b should not lean on.

### M3 — `0025:C4`'s env policy has four composed parts; `0025:S4b` tests one, and the composition order has no oracle

C4's `env:` row is a four-part composition: allowlisted parent vars
(`{PATH, HOME, TMPDIR, LANG, LC_*}`) + `env_pass` + the entry's literal `env`
+ the `INTRASTATE_*` overlay. The `authority` mini-check calls this
"total and ordered". Ordered matters: an entry declaring
`env = { PATH = "/nowhere" }` and the allowlist both name `PATH`; an entry
declaring `env = { INTRASTATE_ROLE = "x" }` collides with the overlay. C4
states the four parts as a sum, never which wins on collision.

S4b tests exactly one negative: `GIT_CONFIG_COUNT` in the parent env does not
reach the child. Nothing tests `env_pass` pass-through (positive), the
`INTRASTATE_*` overlay's presence or values, `LC_*` prefix matching, or any
collision. The illustrative TOML in `§illustrative-code` exercises both
`env` and `env_pass` on `[write.stage]` — but that section is marked
"Illustrative — intent only; tests must not assert it literally", so it
cannot serve as the fixture.

**Test this prevents**: the env-composition test. Three concrete assertions I
cannot write: (a) does entry-literal `env` override an allowlisted parent
`PATH`, or is `PATH` privileged because C4 says "`argv0` still resolves
through the parent's `PATH` before spawn ... which the allowlist passes on
unchanged"? (b) can an entry shadow `INTRASTATE_ROLE`? (c) does `LC_*` mean
prefix-glob over parent vars, when C5 registers `env_pass` as explicitly
"no-glob"? The record uses a glob in the allowlist and bans globs in
`env_pass` without reconciling the two.

**Decision blocked**: the precedence order of C4's four env layers on key
collision, and whether `LC_*` is a literal prefix match.

---

## Low

### L1 — `0025:C5`'s interpreter deny-list is normatively OPEN, so `S1`'s fourth mutant proves less than it appears to

C5 is admirably explicit that this is intended: "interpreter set: OPEN
(deny-listed, not closed) — an unlisted interpreter is admitted", and the prose
says so twice more. No defect in the contract.

The testability note is that S1's mutant list contains exactly one interpreter
form (`["sh", "-c", …]`) and its Expected is "each mutant is rejected with its
own named C5 `table.Category`". A reader of the test suite will read that as
"inline shell is blocked". C5 says the opposite about `perl -e` or a
`busybox sh` alias. There is no scenario asserting the *negative* — that an
unlisted interpreter is admitted and the model loads — which is the assertion
that would keep the deny-list's openness honest in the suite rather than only
in the record.

**Test this prevents**: the "deny-list is open, by design" test. This is
cheap and I can write it from C5's text as it stands; flagging it only because
its absence is how an open deny-list gets mistaken for a closed one by the
next implementer.

**Decision blocked**: none — C5 is decided; this is a coverage addition.

### L2 — `0025:S5`'s `<clear>` arm and `0025:C3` do not address a *read* returning the literal `<clear>`

C3 fixes `<clear>` on the write stdin path ("a planned `<clear>` crosses as the
literal reserved value") and S5's fourth arm asserts "the key reads back
absent". Both correct.

Unaddressed: a command **read** whose child prints `<clear>` as a value. The
existing executor already has a rule for this (`executor.go:96`): a value that
reads back as the reserved literal "did not establish the key's content ... the
key is UNREADABLE", implemented as `classify(requested, true)` at line 102 —
and deliberately *not* applied on the read-back comparison path (line 388),
where the literal's presence is instead the defect. So the behaviour exists and
is inherited, but neither C3 nor any S names it for the command carrier, where
it is newly reachable: a raw-mode read of a `git config` value that literally
contains `<clear>` now mints `incomplete_read` from an otherwise healthy tool.

**Test this prevents**: a raw-read-of-`<clear>` test. I can derive the expected
outcome from the code, but the record's silence means the test would pin
behaviour the RDR never chose — and if the intent is different, the test
enshrines the bug.

**Decision blocked**: whether the inherited clear-literal-is-unreadable rule is
intended to apply to command-backed reads, and whether C3 should say so.

### L3 — `0025:S4`'s ablation arms assert against a build that does not exist

S4's Expected: "the ablations reproduce S3's orphan and S1's blocked `Wait`,
which the test asserts as failures of the ablated build, not of the binding."
The fixture is real — I confirmed `evidence/spikes/a1-deadline/output.txt`
S1 (`survivors=1`, watchdog fired at 5000ms) and S3 (`survivors=1`) against
S2/S4/S5 (`≈1.0s`, `survivors=0`).

The testability gap is *how*. "Ablated build" is not a construct the record
defines: it implies the binding's `Setpgid` / `Cancel` / `WaitDelay` triple is
individually disable-able from a test — a testability requirement on the
implementation that appears nowhere in C4, which states the triple as fixed
behaviour.

**Test this prevents**: the ablation test as written. Without an injection
seam the honest options are to drop the ablations (losing the
necessary-and-sufficient proof that is A1's whole contribution) or to keep the
spike as the only evidence, which makes it a one-time observation rather than a
regression guard — a mechanism that regresses in the binding would still pass
S4's three positive arms if the timeout merely *fires*.

**Decision blocked**: whether C4's deadline triple is required to be
test-injectable, and if not, whether S4's ablation arms are dropped or demoted
to a spike citation.

---

## Notes on what I could *not* fault

- `S1`'s registration assertion is the strongest oracle in the record: it
  correctly identifies that a per-mutant refusal assertion passes without
  `table.Categories()` membership. I confirmed `internal/table/category.go:52`
  returns a hand-maintained literal, so an unregistered constant would indeed
  compile and refuse while staying invisible. Correctly diagnosed.
- `S2`'s "byte-for-byte otherwise" and `S7`'s "asserted by **absence of a
  spawn**, not merely a non-zero exit" are both properly discriminating
  oracles that name the weak alternative and ban it.
- `C4`'s `Detail` field claim checks out: `internal/accessor/model.go:280–301`
  shows `Refusal` with no free-text field, and `Reason`'s own comment
  ("It is never set by this package for a deny") supports reserving it.
- The `§pre-lock-mini-checks` `disposition` table is a better-specified oracle
  set than several of the S entries it backs — the `S3` gap in **H1** is
  partially recoverable from that table's nine rows, though "loud" is an
  adjective, not an assertion, and the table is not cited by any S.
