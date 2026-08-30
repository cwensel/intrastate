Model: claude-opus-5

# Critique — RDR 0025, Command-Invoking Accessor Bindings

Verdict: **BLOCK**. This record specifies a process-execution seam on top of a
mental model of the codebase that is wrong in three load-bearing places. Two of
the three are named as `Pending` assumptions (A11, A12) whose "to verify" text
already contains the false premise, and the third (`intrastate lint` as the
enforcement surface) is asserted as settled fact in the Problem Statement, the
Approach, C5, C6, F1, F2 and the MVV without a single anchor to the function
that would do it. The RDR is 1,500 lines long and its most-repeated claim —
"lint proves the vector is well-formed" — has no verified assumption behind it.

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0025:A12`, `0025:C2` ("the loader records the model file's directory on `table.Model` … set where the file is already opened") | `table.Load` performs no file I/O and receives no path — only `(src []byte, sourceID string)`. `sourceID` is a *display* string (`internal/table/load.go:30`, `:59`, `:69`), used only for line locators, never stored on `Model`. There is no "where the file is already opened" inside the loader; the open is at `internal/cli/flow_input.go:140` and `internal/cli/lint.go:161`, two call sites, one of which the RDR never mentions. And `--model` is stored verbatim (`selectModelPath`, `flow_input.go:112-132`), so the "one absolute-path field" is not absolute. | Separator-bearing `argv0` (`tools/flowstate-write`) — the exact form the RDR's own Illustrative Code and A3's wrapper conclusion make mandatory for every established-tool write — resolves against nothing, or against a relative base, or against cwd. Either every write model refuses `execution_failure` at binding construction, or the RDR's central reversal (`exec.ErrDot`) is silently reintroduced. | §1, §3, premortem, AT-1 |
| C-2 | `0025:A11`, `0025:C4` (`detail:` line: "`executor.go::refusalOf` — the constructor every *post-selection* refusal already passes through") | False on `main`. `refusalOf` is *not* the sole post-selection constructor: `refusalWithKeys` (`executor.go:108-112`) wraps it, and `Executor.Read` routes **every** read refusal through `refusalWithKeys`, not `refusalOf` — including `ClassExecutionFailure` from `invokeRead` (`executor.go:92-93`). A11's own "verify by reading `refusalOf` and every one of its call sites" does not name `refusalWithKeys`. Adding the error parameter to `refusalOf` alone therefore reaches gate and write refusals but silently drops the stderr tail on the entire read path. | The one capability the RDR says binds *directly* to established tools — reads — is the one whose failures carry no diagnosis. `git config` exits 128 with a real stderr message; the user sees `the accessor \`branch_state\` could not be executed` and nothing else. | §1, premortem, AT-2 |
| C-3 | `0025:C4` (`detail:` and the "Rendering is a second hop" paragraph: "the stderr tail renders beneath it") | `clierr.CLIError` has exactly ONE `Detail string` field (`internal/cli/clierr/clierr.go:58`) and `EmitText` renders exactly one `  detail: %s` line (`clierr.go:194-195`). There is no "beneath it". The RDR specifies a precedence rule ("the applied-sense text wins the `CLIError.Detail` slot and the stderr tail renders beneath it") for which no carrier exists, and never opens a Capability Dependency or Existing Infrastructure Audit row for `CLIError` — only for `accessor.Refusal`. | On the failure that matters most — a write that ran, could not be verified, and whose tool printed the reason — the user gets the applied-sense sentence and the tool's stderr is discarded. Or an implementer concatenates a 4 KiB tail into a field `sanitizeLine`'s siblings assume is one line, and the JSON envelope grows a multi-KiB `detail`. | §1, §2, premortem, AT-3 |
| C-4 | `0025:§approach` ("`intrastate lint` validates the declaration"), `0025:C5`, `0025:C6` ("lint (C1/C5) validates regardless"), `0025:F1`, `0025:F2`, `0025:MVV` step 2 | `intrastate lint` is `table.LoadWithAdvisories` + `graphlint.Run` (`internal/cli/lint.go:176-200`). `accessor.Validate` — the eight-arm function the RDR treats as the live validator, including `CodeAmbientArtifactDiscovery` which C4's env policy leans on — has **zero production callers** (verified by sweep). And `table.Load` is fail-fast: one categorized error per *document* (`load.go:20-28`), not per entry. C5's "precedence: WITHIN one entry" clause is written as if `accessorTable` reported per-entry lists; it does not. | An author with a five-entry command model fixes one defect, re-runs lint, gets the next defect, re-runs, gets the next. The "reviewable before it is trusted" promise (F2, C6) degrades to a five-round game of whack-a-mole on a surface the RDR advertises as static review. | §1, §2, premortem, AT-4 |
| C-5 | `0025:C1` ("both or neither is a load-time defect"), `0025:A9` | C1 relaxes `accessorTable`'s `path is absent or empty` arm (`load.go:951-952`) but `Accessor.Path` is a plain `string` (`load.go:986-992`), and `flowbind.Registry` unconditionally constructs `Reader{Path: acc.Path}` / `&Writer{Path: acc.Path}` / `Gate{Path: acc.Path}` for every entry (`registry.go:31-55`). A command entry now yields `Path: ""`. `flowbind.load("")` returns an empty store on `os.ErrNotExist` (`flowbind.go:114-118`) — it does not error. So a selection bug does not crash; it silently binds a command entry to an always-empty artifact. | The most dangerous defect in the record: a mis-dispatched command **read** returns "every key established absent" rather than refusing, and a mis-dispatched command **write** writes to `""` and reports success. The RDR's own "presence is the answer" rule turns the discriminator bug into silent state loss. | §1, premortem, AT-5 |
| C-6 | `0025:S5b`, `0025:A4` ("a role with a path-backed and a command-backed reader, both admissible to today's first-match `readerFor`") | `checkAccessorBindings` (`load.go:1010-1033`) refuses at load: an OWNED tag served by ≠1 reader is `malformed_accessor_binding`; an OBSERVED tag served by >1 is the same. Two readers on one role are authorable **only** with disjoint key sets. So S5b's scenario as written does not load, and the write's read-back reader is then, for the planned owned key, uniquely determined by the loader — not by `readerFor`'s first-match. The scenario asserts a property (`registry` sorts, `readerFor` takes first) that the loader has already made unobservable for the case that matters. | The one test written to pin the pre-0016 behaviour so "0016's landing is a visible change, not a silent one" is unwritable as specified. It will be quietly reshaped into something that passes, and the 0016 landing becomes exactly the silent change the scenario existed to prevent. | §1, §3, premortem, AT-6 |
| C-7 | `0025:C2` ("the substituted value must be an absolute path — the command binding refuses … when the caller-bound artifact path is relative") | `parseArtifacts` (`flow_input.go:239-251`) stores `--artifact role=path` verbatim, and every existing test, fixture and doc example in this repo uses relative paths. C2 accepts this cost in one sentence ("The cost is stated") and never reaches the Consequences list, the MVV, or F1–F8. It is a user-facing behavioural asymmetry between two carrier kinds of the *same* accessor, surfaced only at runtime. | A user converts one read entry from `path` to `command`, changes nothing else, and their entire invocation line stops working with `execution_failure`. Bisecting this against a working path-backed model teaches nothing, because lint passes both. | §1, §3, premortem, AT-7 |
| C-8 | `0025:C6` ("gate site: the command binding's constructor"), `0025:S7` ("the flag omitted on each of the three production `NewExecutor` callers") | The three-`NewExecutor`-callers framing is the wrong axis. `flowbind.Registry(model)` has exactly **one** production call site — `buildRequest` (`flow_exec.go:830`), shared by all three executor verbs. So the flag has one place to be read, and S7's "all three callers refuse identically" tests one code path three times. Meanwhile the real risk C6 never addresses: `buildRequest` takes `cmd`, so the flag must be *registered* on all three verbs or `cmd.Flags().GetBool("allow-commands")` returns `false` for a verb someone forgot — indistinguishable from a user who did not pass it. | A verb gains an executor, nobody registers the flag, and every command invocation through it refuses `execution_failure` naming `allow_commands` while the user is passing `--allow-commands` and cobra is rejecting it as unknown, or silently returning false. The RDR's "gains the flag with it — the coupling is to executor construction" is aspirational; nothing enforces it. | §1, §2, premortem, AT-8 |
| C-9 | `0025:C3` ("the inherited reserved-literal rule composes unchanged … `internal/accessor/executor.go:96`") | Line 96 of `executor.go` is a comment fragment ("A value that reads back as the reserved `<clear>` literal did not"). The rule is at `executor.go:139-143` inside `readOutcome.classify`, gated by a `clearIsUnreadable bool` parameter that `Executor.Write`'s read-back passes as **false** (`executor.go:387`). C3 says the rule "composes unchanged" and S5 asserts it "by class"; it does not compose on the read-back path at all, by deliberate design. | An implementer reads C3, wires the `<clear>`-is-unreadable rule into the command binding's raw-mode decoder, and a legitimate write of a value the tool renders as `<clear>` now refuses on read-back where the contract says it should mismatch. Line-number anchors into a file this RDR expects to be edited are stale on arrival. | §1, premortem, AT-9 |
| C-10 | `0025:C5` (`precedence: WITHIN one entry — accessorTable is fail-fast (one categorized error, never a list)`) and `0025:§testing-strategy` S1 | The fail-fast property is `Load`'s, not `accessorTable`'s (`load.go:20-28`, `run()` at `:73-101` returns on the first step error). `accessorTable` is called three times in sequence for read/write/gate (`load.go:908-919`), each returning on its first bad entry. So "ACROSS entries there is no order" understates it: across *capability tables* the order IS fixed (read, then write, then gate), and a test asserting a gate defect in a model that also has a read defect will deterministically see the read one. S1's "each S1 mutant carries exactly one defect in one entry" is the only way to make the suite pass, and the RDR discovered that constraint without stating why. | Not user-visible directly; it is the reason the C5 test suite will be rewritten. A contributor adds a seventh mutant with two defects, the assertion fails nondeterministically-looking, and the "unspecified which" clause gets cited to delete the assertion rather than fix the model. | §2, AT-10 |
| C-11 | `0025:§metadata` Profile (`foundational`), `0025:G-proportionality` ("the sole author of at most one independent load-bearing contract") | Six normative contracts, six new load categories, a new struct field on `accessor.Refusal`, a new field on `table.Model`, a new CLI flag on three verbs, a new binding package, a change to `refusalOf`'s signature, a change to `flowbind.Registry`'s signature, and an env-composition policy with a four-layer precedence order. The Profile field itself concedes the count ("C1/C2/C5/C6 … with its invocation envelope (C3/C4), inseparable — no split") and then declares no split. C1 (carrier), C4 (env + process safety) and C6 (execution gate) are three independently reversible decisions with three different reversal costs. | Not a user symptom — a delivery symptom. This lands as one branch touching `table`, `accessor`, `cli`, `flowbind` and a new package, and the first production defect in any one of them cannot be reverted without reverting the other five. | §2, premortem |
| C-12 | `0025:A3` ("that wrapper class is small and mechanical, not bespoke per integration"), `0025:§approach` ("Integration honesty") | The spike's own result is that **every** established-tool write needs a wrapper, and the RDR's own Illustrative Code demonstrates this (`tools/flowstate-write`). A3's status is `Verified` on evidence that *refutes* the assumption's headline as originally posed; the assumption text was narrowed until the spike agreed with it. What survives is "reads and gates bind directly" — but every write, the half the Problem Statement says is missing, requires a shell script the model does not carry and lint cannot see. | The user's first real integration is a write. They write a shell script. The RDR's central promise — "the edit becomes a linted, declared command" — is true of the wrapper's *invocation* and false of the wrapper's *body*, which is where every actual failure will live. | §3, premortem, AT-11 |
| C-13 | `0025:F6`, `0025:A3` (hazard P-1, `tee {artifact}` overwrites the artifact) | The spike demonstrated that a stdin-sinking `argv0` destroys the user's artifact and that this surfaces only at read-back. The RDR's response is a Failure Modes bullet and a prose sentence ("the wrapper contract … is the answer"). There is no C5 defect for it, no lint arm, no runtime guard, and the recovery line (`0025:F8`) explicitly says bindings hold no state — no backup, no undo. | A user's state artifact is overwritten with `{"flow.stage":"review"}` by a plausible one-line declaration that passes lint, and the only signal is `read_back_incomplete`. Data loss with a green lint. | §3, premortem, AT-12 |

## 1. The three most likely ways implementation goes wrong

### 1.1 The model file's directory does not exist as a thing the loader can record (C-1)

**Root cause in the RDR.** `0025:A12` and `0025:C2` assert a fact about
`internal/table/load.go` that is false on `main`. The signature is
`Load(src []byte, sourceID string) (*Model, error)`, and the package doc says
so explicitly: *"It takes already-read BYTES plus a source id, never a
filesystem path: this package performs no file I/O and no path resolution
(`0002:EIA`)."* That sentence is a locked peer's contract (RDR 0002), quoted
inside the very function A12 proposes to modify.

**The passage that enabled it.** `0025:C2`: *"the loader records the model
file's directory on `table.Model` (one absolute-path field, set where the file
is already opened)"*. And `0025:A12`'s verification plan: *"Verify by reading
the load site and enumerating in-memory `table.Model` constructions."* The plan
looks for the wrong thing. There is no load site *inside* the loader. There are
two *callers* that open files — `selectModel` (`flow_input.go:140`) and
`runLint` (`lint.go:161`) — and A12 names neither. Worse, both hand `Load` the
user's `--model` string verbatim, which `selectModelPath` never absolutizes. So
even after the field is added, it holds a relative path, and C2's "never
cwd-relative" guarantee is violated by the mechanism C2 chose to honour it.

**Symptom.** The wrapper form is mandatory for every established-tool write
(A3's own conclusion, and the RDR's Illustrative Code writes
`command = ["tools/flowstate-write", "{artifact}"]`). Under the honest reading
of A12, that entry has no base dir to resolve against, so C2's own rule fires:
*"A separator-bearing argv0 in one is refused `execution_failure` at binding
construction."* The user's model lints clean and refuses to run, with a message
about a base dir they never declared. Under the *dishonest* reading — an
implementer wires `filepath.Dir(sourceID)` in and moves on — the resolution is
relative to the process cwd, which is precisely the `exec.ErrDot` reversal C2
cites four paragraphs of prior art to avoid.

### 1.2 The stderr tail reaches gates and writes but never reads (C-2), and has nowhere to render (C-3)

**Root cause.** `0025:A11` and `0025:C4` are built on the claim that
`executor.go::refusalOf` is the chokepoint every post-selection refusal passes
through. It is not. `refusalWithKeys` (`executor.go:108-112`) is a second
constructor, and `Executor.Read` uses it for **every** refusal it mints —
including `ClassExecutionFailure` returned by `invokeRead`
(`executor.go:92-93`). A11's stated method ("reading `refusalOf` and every one
of its call sites") would find `refusalWithKeys` as a *call site* but would not
notice that it is also a *bypass* for the one capability that matters most.

**The passage.** `0025:C4`, `detail:` line, and the "The hook is
**`executor.go::refusalOf`**, not the invocation helpers" paragraph, which
spends six sentences arguing that a single edit "covers every present and future
invocation site." It covers gate and write. It does not cover read.

**Compounding.** Even where the `Detail` does arrive, `0025:C4`'s rendering rule
has no carrier. `clierr.CLIError` holds one `Detail string`
(`clierr/clierr.go:58`); `EmitText` prints one `  detail:` line
(`clierr.go:194-195`). C4's *"the applied-sense text wins the `CLIError.Detail`
slot and the stderr tail renders beneath it"* describes a two-slot renderer
this repo does not have, and no Existing Infrastructure Audit row was opened for
`clierr`. `0025:S3` asserts the tail is present "on an `execution_failure` that
carries no applied sense" — chosen, apparently, because it is the only case the
one-slot renderer can satisfy. The contract's harder half is untested by
construction.

**Symptom.** Reads bind directly to established tools — this is O2's winning
row in the QOC matrix. When `git config --file /abs/x --get flow.stage` exits
128 with `fatal: bad config line 1 in file /abs/x`, the user sees
`error: accessor-failed: the accessor \`branch_state\` could not be executed`
and nothing more. The single most common integration failure produces the least
diagnosable message in the system.

### 1.3 The carrier discriminator fails open into an empty artifact (C-5)

**Root cause.** `0025:C1` relaxes the loader's `path is absent or empty` arm
and `0025:D-selection-predicate` puts the discriminator in
`flowbind/registry.go::Registry`. But `Accessor.Path` is a plain `string`
(`load.go:986-992`, no pointer), so a command entry carries `Path: ""` — an
indistinguishable-from-unset zero value. `flowbind.Registry` today constructs
`Reader{Path: acc.Path}` for every reader with no condition
(`registry.go:31-40`). And `flowbind.load("")` treats `os.ErrNotExist` as an
**empty artifact, not an error** (`flowbind.go:114-118`) — a deliberate 0005
decision so that a first `set-state` can create the file.

**The passage.** `0025:A9` sweeps every reader of `Path` and concludes *"a
locator-less command entry breaks no other partition."* It enumerates the
*readers* correctly and never asks the question that matters: what does the
existing binding **do** with an empty `Path`? The answer — silently succeed —
converts every possible discriminator bug from a crash into data loss.

**Symptom.** A command-backed read whose entry slipped past the discriminator
reports every declared key established-absent. Under RDR 0004's rules that is a
*successful* read of an artifact holding nothing. The resolver then decides a
transition from a state it believes is empty. A command-backed *write* that
slipped past writes JSON to the path `""` and returns success; the read-back
reads `""`, finds the values it just wrote, and confirms. The user's real
artifact is untouched, `set-state` exits 0, and the reported `owned` map is a
fabrication. Nothing in C1–C6, S1–S7, or F1–F8 detects this, because every test
in the record exercises the discriminator on the *happy* branch.

## 2. The one section that will be rewritten within 6 weeks of shipping

**`0025:§normative-contracts` C4** — specifically the `detail:` line and the
"Rendering is a second hop and a **name collision to avoid**" paragraph.

It will be rewritten because it is the only clause in the record that specifies
behaviour across three packages (`accessor`, `cli`, `cli/clierr`) while opening
an infrastructure-audit row for only one of them. Its three components each fail
independently:

- The `refusalOf` hook misses the read path (C-2), so it will need either a
  second hook at `refusalWithKeys` or a change to `invokeRead`'s
  `readOutcome` — which is an interface change to a struct RDR 0004 pins.
- The rendering rule presumes a slot that does not exist (C-3). Someone will
  add `CLIError.Detail2`, or concatenate with a `\n` and break the one-line
  invariant `sanitizeLine` exists to protect, or drop the tail. All three are
  amendments to C4.
- The `env:` line's four-layer precedence with an `LC_` prefix rule, an
  `env_pass` no-glob rule, and an `INTRASTATE_` shadow defect is more policy
  than any other clause in the RDR carries, and its only test
  (`0025:S4b`) is written against "the child's observed environment," which
  requires a test-only introspection seam the RDR never specifies — the same
  gap C4's `platform:` line *did* solve with an injectable `goos` var, and did
  not solve here.

Add C-8's flag-registration gap and C-10's cross-table ordering and C4 becomes
the clause every follow-up kata cites. The rewrite is not a risk; it is
scheduled.

## 3. The one assumption that will not survive first contact with a real user

**`0025:A3`** — *"that wrapper class is small and mechanical, not bespoke per
integration."*

This assumption is stamped `Verified` by a spike whose result contradicted its
original form. The Evidence says so in its own first sentence: *"narrowed from
'a useful class binds directly'."* What the spike actually established is R3:
***every*** established-tool write needs a wrapper, because a fixed argv cannot
carry the planned value and the tools do not read stdin.

The RDR's own Problem Statement says the missing half is the **write**:
*"a model can DECIDE a state change … but no binding invokes a command, so the
change cannot be APPLIED."* A3 verifies that the write half — the half this RDR
exists for — is served by a shell script whose body the model does not carry,
lint cannot see, `command_shell_interpreter` explicitly permits, and C2's
authority bound does not reach. The Approach concedes this in a parenthetical
("accepted, recorded as a Consequence") and then the Decision Rationale scores
O2 a **4** on "Binds established tools without wrappers" by scoring only the
read/gate half — a row-splitting the RDR performs openly (*"The second row is
scored on the *read and gate* half"*) and then does not adjust the total for.
O2's margin over O1 is 4 points; that row alone is 2.

The first real user's first real task is applying a state change. They will
write `tools/flowstate-write`. That script — not the model — will contain the
`git config` invocation, the key iteration, the error handling, and the bug.
The reviewer who reads the model sees `["tools/flowstate-write", "{artifact}"]`
and has learned nothing about what executes. That is the *prose-and-drift
problem the table was built to remove*, relocated by exactly one file — which
is the failure the Problem Statement's second paragraph names as the thing this
RDR must not do.

Corollary that lands the same day: C-7's absolute-path refusal. Every example
in this repo binds `--artifact repo=./state.json`. Converting one entry to
`command` makes that invocation refuse, with lint green.

## 4. Premortem — written from six months after ship

We shipped 0025 in four phases as planned. It has been six months.

**Week 1.** The first internal model converted one read to a command entry:
`["git", "config", "--file", "{artifact}", "--get", "flow.stage"]`, raw mode,
`exit_absent = [1]`. It lints. It refuses. The refusal is
`accessor-failed: the accessor \`branch_state\` could not be executed`. Three
engineers spend an afternoon on it. The cause is C-7: the team's runbook says
`--artifact repo=state/flow.json`, C2 refuses a relative path, and the refusal
`Detail` that was supposed to say *"artifact path must be absolute for a
command entry"* never reaches the terminal, because `accessorFailure`
(`flow_exec.go:384-386`) builds `envErr(codeAccessorFailed, …)` and sets no
`Detail` at all — C4 specified a field on `accessor.Refusal` and a rendering
rule, and the rendering hop was implemented as "wire it later." We add
`ce.Detail = refusal.Detail` in week 2. It is the first amendment to C4.

**Week 3.** The first write lands. `tools/flowstate-write` is 40 lines, not 8 —
the spike's `wrapper-write.sh` did not handle `<clear>`, did not handle a key
whose value contains a newline, and did not `set -euo pipefail`. It resolves
against the model file's directory, except that C-1 bit us in week 2: `Load`
takes no path, so we threaded `filepath.Dir(sourceID)` onto `table.Model` from
`selectModel`, and `runLint` — the other caller — passes a different string. A
model that lints from `docs/` and runs from repo root resolves `argv0` to two
different files. We patch `selectModelPath` to `filepath.Abs` the `--model`
value. That changes every load locator in every findings payload from
`flows/probe.toml` to `/home/ci/work/flows/probe.toml`, and four golden tests in
`internal/cli` break. We special-case the locator. C2 and 0002's "no path
resolution" contract are now both violated, in two places, differently.

**Week 6.** A user runs `set-state --write flow.stage=review --allow-commands`
and it exits 0, reporting `owned: {flow.stage: review}`. The artifact is
unchanged. This is C-5. Their model declared `[write.stage]` with `command` and
no `path`; `flowbind.Registry` had a discriminator, but the model also declared
a `[read.branch_state]` with `command` and the registry's *reader* arm still
read `acc.Path` — one of three near-identical loops in `registry.go:31-55`, and
the reviewer's eye slid over the one that was not changed. `Reader{Path: ""}`
called `flowbind.load("")`, got `os.ErrNotExist`, returned an **empty store**
per `flowbind.go:114-118`, and reported `flow.stage` established-absent. The
write's read-back at `Executor.Write` (`executor.go:352`) called that same
reader through `readerFor`, compared `planned` against an empty artifact, found
`flow.stage` absent, and returned `read_back_mismatch` — except it didn't,
because `verifyReadBack` was comparing against a `before` snapshot that was
*also* empty, and the plan's own key was excluded from `protected` by
`protectedKeys`. It returned success. We find it because a user's state artifact
had been stale for nine days.

**Month 3.** The stderr-tail work is reopened. `Executor.Read` refusals
(`executor.go:92-93` → `refusalWithKeys`) carry no `Detail`, which is C-2, and
reads are the capability that binds directly, so reads are where every tool
error appears. We add the error parameter to `refusalWithKeys` too, then to
`readOutcome`, then discover that `invokeRead` is shared with the read-back path
(`executor.go:352`, `:317`) and a read-back's stderr tail now overwrites the
applied-sense text through the one `CLIError.Detail` slot (C-3). We ship a
truncating concatenation. A tool that writes 4 KiB of stderr now emits a 4 KiB
`detail` string into the JSON envelope, and a downstream consumer's log ingest
starts dropping records.

**Month 4.** Someone declares `command = ["tee", "{artifact}"]` for a write,
because `tee` reads stdin and stdin is where C3 puts the planned tags. It lints
clean — `tee` is not an interpreter, the argv is fixed, the placeholder is
whole-element. It overwrites the state artifact with
`{"flow.stage":"review"}` and exits 0. Read-back through the `git config`
reader fails `bad config line 1` → `read_back_incomplete`. This is C-13,
verbatim: F6 predicted it, named the exact spike case, and shipped no guard.
There is no undo (`0025:F8`). We add a lint advisory for a known stdin-sinking
argv0 list, which is a deny-list over a deny-list.

**Month 5.** RDR 0016 lands. `readerFor` becomes fail-closed. The scenario that
was supposed to make this a visible change — `0025:S5b` — was never written as
specified, because `checkAccessorBindings` (`load.go:1010-1033`) refuses a model
where two readers serve the same owned tag, so the two-reader-one-role fixture
did not load. Whoever implemented S5b reshaped it into a single-reader assertion
that passes before and after 0016. That is C-6. The 0016 landing is silent.

**Month 6.** `--allow-commands` is registered on `flow next` and `set-state` but
not on `flow exec`'s second path, because C6's *"A verb that later gains an
executor gains the flag with it"* is a sentence, not a mechanism, and
`buildRequest` (`flow_exec.go:830`) reads the flag off `cmd` for all three.
`cmd.Flags().GetBool` returns `false, err` for the unregistered verb; we ignore
the error the way every other flag read in `flow_input.go` does. The verb
refuses every command invocation naming a gate the user *did* pass. This is
C-8. S7 tested "all three callers" and passed, because all three callers go
through one code path and the test registered the flag on all three fixtures.

The record was 1,500 lines. Every one of these failures is either asserted-false
in it (C-1, C-2, C-3, C-4, C-6) or named-and-unmitigated (C-5, C-7, C-13). None
of them is a thing we could not have known at review time.

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time tests: each is executable against `main` *before* any
0025 code is written. Each maps to a ledger row.

**AT-1 (C-1) — the loader's path surface**
```gherkin
Given the RDR asserts the loader records the model file's directory
When I read the signature of internal/table/load.go::Load
Then it takes a filesystem path or an *os.File
And table.Model or the loader struct persists that path beyond locator use
And every production caller of table.Load passes an absolute path
```
Result on `main`: FAILS at step 1 (`Load(src []byte, sourceID string)`),
step 2 (`sourceID` is stored on `loader`, never on `Model`), and step 3
(`flow_input.go:149` and `lint.go:176` both pass the raw `--model` value).
A12's Status must be `Refuted`, and C2's argv0 clause must name a mechanism
that exists.

**AT-2 (C-2) — refusalOf is the sole post-selection constructor**
```gherkin
Given C4 hooks the stderr tail into executor.go::refusalOf
When I enumerate every construction of *accessor.Refusal in internal/accessor
Then refusalOf is the only constructor reached after accessor selection
And no capability path constructs a Refusal through any other helper
```
Result on `main`: FAILS. `refusalWithKeys` (`executor.go:108`) is a second
constructor and is the ONLY one `Executor.Read` uses. Fix: C4's `detail:` line
must name both, or the hook must move into `readOutcome`.

**AT-3 (C-3) — the rendering slot exists**
```gherkin
Given C4 says the stderr tail "renders beneath" the applied-sense text
When I read clierr.CLIError and clierr.EmitText
Then there are two distinct detail-bearing fields
Or the RDR opens an Existing Infrastructure Audit row for clierr naming the limit
```
Result on `main`: FAILS on both. One `Detail string` (`clierr.go:58`), one
render line (`clierr.go:194`), and no audit row. Either add the row and a
contract clause for the second slot, or delete the "renders beneath" sentence
and state that the applied sense wins outright and the tail is dropped.

**AT-4 (C-4) — what `intrastate lint` actually validates**
```gherkin
Given the RDR says "intrastate lint validates the declaration" in six places
When I read internal/cli/lint.go::runLint end to end
Then accessor.Validate is invoked on the constructed registry
And a model with three distinct C5-class defects reports three findings
```
Result on `main`: FAILS. `runLint` calls `table.LoadWithAdvisories` and
`graphlint.Run`; `accessor.Validate` has zero production callers. `table.Load`
is fail-fast (`load.go:20-28`), so three defects report one finding. C5's
"precedence: WITHIN one entry" text must be corrected to "one defect per
document," and every "lint validates" claim must be re-anchored to `table.Load`
or the RDR must own wiring `accessor.Validate` in.

**AT-5 (C-5) — the empty-Path failure mode**
```gherkin
Given C1 permits an accessor entry with no path
When flowbind.Reader{Path: ""}.Read is invoked on any artifact
Then it returns an error
```
Result on `main`: FAILS — `flowbind.load("")` returns an empty store on
`os.ErrNotExist` (`flowbind.go:114-118`), so the read succeeds and reports every
key established-absent. This is a two-line check that would have forced a C1
clause: `Accessor.Path` becomes `*string`, or `flowbind`'s constructors refuse
an empty locator. Neither is in the record.

**AT-6 (C-6) — S5b's fixture loads**
```gherkin
Given S5b needs one role served by a path-backed and a command-backed reader
When I author that model and call table.Load
Then it loads
```
Result on `main`: FAILS for the case S5b cares about.
`checkAccessorBindings` (`load.go:1021`) refuses an owned tag served by ≠1
reader. The fixture only loads with disjoint key sets, under which the
read-back reader for the planned key is unique and `readerFor`'s first-match is
unobservable. S5b must be respecified, or dropped with the 0016-visibility
claim withdrawn.

**AT-7 (C-7) — the relative-artifact regression is user-facing**
```gherkin
Given C2 refuses a relative --artifact path for a command entry
When I grep the repo's fixtures, docs and tests for --artifact bindings
Then relative paths are absent
```
Result on `main`: FAILS — relative paths are the norm (`parseArtifacts`,
`flow_input.go:239-251`, stores verbatim). This is a Consequence and a Failure
Mode, and appears as neither. Add both, or absolutize in `parseArtifacts` for
all carriers uniformly.

**AT-8 (C-8) — the gate's single site**
```gherkin
Given C6 gates in the command binding's constructor via the registry
When I count production call sites of flowbind.Registry
Then the count matches the RDR's three-NewExecutor-callers framing
And the flag is registered on every cobra command that reaches buildRequest
```
Result on `main`: `flowbind.Registry` has ONE production call site
(`flow_exec.go:830`), shared. S7's "all three callers refuse identically" is one
path tested thrice. The real hazard — flag registration drift on a verb that
reaches `buildRequest` — is untested. C6 needs a mechanism (a registration
helper the verbs must call, asserted by a test that enumerates commands
reaching `buildRequest`), not a sentence.

**AT-9 (C-9) — line-number anchors resolve**
```gherkin
Given C3 and S5 cite internal/accessor/executor.go:96
When I read that line on main
Then it contains the reserved-literal rule
```
Result on `main`: FAILS — line 96 is a comment fragment. The rule is at
`:139-143` and is parameterized by `clearIsUnreadable`, which the read-back path
passes as `false` (`executor.go:387`). C3's "composes unchanged" is wrong for
the read-back half. Replace every `file:line` anchor with `file::Symbol`.

**AT-10 (C-10) — cross-table defect ordering**
```gherkin
Given C5 says cross-entry reporting order is unspecified
When a model carries a read-entry defect and a gate-entry defect
Then which is reported is genuinely unspecified
```
Result on `main`: FAILS — `loadAccessors` calls `accessorTable` for read, then
write, then gate (`load.go:908-919`), returning on the first error. Across
*tables* the order is fixed and deterministic. C5's precedence clause must say
so, or a test will encode the accidental order as contract.

**AT-11 (C-12) — the wrapper is the integration**
```gherkin
Given A3 claims wrappers are "small and mechanical, not bespoke"
When I read the spike's wrapper-write.sh
Then it handles <clear>, values containing newlines, partial failure, and set -e
And the MVV's write path requires no script the model does not carry
```
Result: FAILS. The spike wrapper is 8 lines and handles none of those; the MVV
explicitly says *"means the declared wrapper for the write."* Either A3 is
downgraded and the QOC row 2 score for O2 is recomputed against O1 (which
collapses the 4-point margin), or the RDR states plainly that v1 delivers a
declared *read* carrier and a wrapper-mediated write.

**AT-12 (C-13) — the stdin-sink hazard has a guard**
```gherkin
Given F6 names tee {artifact} as an artifact-destroying declaration
When that entry is linted
Then some C5 category, advisory, or runtime guard fires
```
Result: FAILS — no clause covers it. C5 has six categories and none applies.
Either add a seventh (an argv0 deny-list for known stdin sinks, matching the
`command_shell_interpreter` precedent), or state in Consequences that a declared
write command can destroy the user's artifact with a green lint and no undo.
