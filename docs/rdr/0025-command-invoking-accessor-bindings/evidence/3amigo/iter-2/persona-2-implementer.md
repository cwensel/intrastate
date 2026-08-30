Model: claude-opus-5[1m]

# Persona 2 — Implementer (iter-2, delta)

Scope: 0025:C1–C6 and new assumptions 0025:A10–A12 only. Everything outside
that delta is out of scope for this pass.

Grounding note: every code claim below was checked against
`/Users/cwensel/sandbox/newcoinc/intrastate` on the current tree.

---

## Findings

### F1 — High — `0025:C4` (and `0025:A11`) name two call sites that do not exist

C4's `detail:` line and its prose both assert:

> the executor's `invokeRead`/`invokeGate`/`invokeWrite`, which today discard
> the error and keep only the class, `errors.As` it and copy `Detail` onto the
> `Refusal`

A11 repeats it: "`executor.go::invokeRead`/`invokeGate`/`invokeWrite` today
discard it ... (confirmed on `main`)" and directs the verifier to "read the
three `invoke*` sites".

There is exactly **one** `invoke*` helper in the package. `grep -n "invoke"
internal/accessor/executor.go` yields only `invokeRead` (declared at
`internal/accessor/executor.go:156`, called at lines 91, 293, 370). There is
no `invokeGate` and no `invokeWrite`.

The gate and write error handling is **inline in the exported methods**, not
in helpers:

- gate: `internal/accessor/executor.go` `Executor.Gate` calls
  `binding.Gate(bounded, art)` and then `if err != nil { return
  GateResult{Refusal: refusalOf(def, timeout, ClassExecutionFailure)} }`
  (~line 205–210). There is also a *second* execution-failure arm in the same
  method — the `default:` branch on an unrecognized verdict — plus the
  not-a-`GateBinding` arm above it.
- write: `Executor.Write` calls `binding.Apply(applyCtx, art, ...)` at
  `internal/accessor/executor.go:328` and handles the error at line 340–343.
  `Write` also raises `ClassExecutionFailure` from the non-`WriteBinding` arm
  and from the `nonOwnedPlanKeys` arm — neither of which has an underlying
  `error` to `errors.As`.

So the "three named call sites, one field copy each" plan is not implementable
as written. On Monday I cannot tell whether the intended edit is:

(a) three `errors.As` insertions at the three *error-returning* sites
    (`invokeRead` line ~176, `Gate` line ~209, `Write` line ~340), which is
    what the contract means but not what it names; or
(b) a refactor extracting `invokeGate`/`invokeWrite` helpers to match the
    clause, which is a structural change to 0004's executor that C4's
    "0004 executor unchanged" framing elsewhere is at pains to avoid.

It also leaves undecided what `Detail` is on the execution-failure arms that
have **no** error to inspect (wrong-binding-type, unrecognized verdict,
non-owned plan key). Empty is the only sane answer but the clause does not say
it, and S3 asserts `Detail` non-empty on a specific `execution_failure` case,
so "which execution_failure carries a tail" is now a test-visible question.

**Decision blocked**: where the `errors.As` hook goes, whether 0004's executor
gets refactored, and what `Detail` holds on error-free execution failures.

Note this is *not* the closed-issue "field is inert" concern — the typed-error
mechanism is well specified. The defect is purely that the three named sites
are two-thirds fictional, and A11 instructs the verifier to confirm the claim
by reading them, which will fail.

### F2 — Medium — `0025:C4` routes the `INTRASTATE_*` shadowing defect to a `0025:C5` category whose definition does not cover it

C4's `env:` line ends:

> An `env`/`env_pass` entry naming an `INTRASTATE_*` key is a C5
> `command_output_shape` defect at load

C5's normative block defines `command_output_shape` as exactly three arms:

> `output = "raw"` with keys ≠ 1; `exit_absent` on a non-read; `exit_verdicts`
> on a non-gate or naming a non-verdict

An env-key collision is none of the three, and nothing in C5's prose mentions
`env` or `env_pass` at all. C5 also states "**Five** new load-time defect
classes" and that "S1 asserts five" — and S1's mutant list is the five MVV
mutants, none of which is an env mutant.

This is exactly a rewrite-introduced cross-clause disagreement: C4 was edited
to add the shadowing rule and reached for an existing C5 name; C5 was edited
separately and its category definition was not widened. Implementing literally
means emitting a category spelled `command_output_shape` for a defect that has
nothing to do with output shape, which is a diagnostic the user cannot act on.

**Decision blocked**: whether the env-shadowing defect is a sixth category
(which moves the count C5 pins at five and S1 asserts), an added arm on
`command_output_shape` (which needs C5's definition line amended), or folded
into `command_empty`. Also unstated: is the rule "key has prefix
`INTRASTATE_`" or "key is one of the four overlay names"? The clause says
`INTRASTATE_*`, which is the broader reading, but C4's rationale ("so the
overlay is never shadowed silently") only motivates the four.

### F3 — Medium — `0025:C2` invents a defect category name that `0025:C5` does not carry

C2's closing paragraph:

> a separator-bearing argv0 in one is a load-time `command_empty`-family
> defect rather than a cwd-relative resolve

`command_empty` is defined in C5 as "empty vector or empty argv element". A
separator-bearing argv0 in a model that has no base dir is neither. "family"
is not a thing the implementation has: `internal/table/category.go` holds flat
`Category` string constants (`CatMalformedAccessorDeclaration Category =
"malformed_accessor_declaration"` etc.) and `table.Categories()` returns a flat
literal slice — there is no family/grouping mechanism to hang this on.

Compounding it: C5 pins the closed set at five and normatively fixes their
tail order in `Categories()`, so a sixth cannot be quietly added, and reusing
`command_empty` for "no base dir" produces a refusal message that contradicts
its own category name.

Separately, this arm is only *reachable* from an in-memory model, and every
in-memory model construction is in test code — so the "load-time" framing is
itself questionable: the loader that would report a `Category` is precisely
the path that always *has* a base dir. A12's "If wrong" branch acknowledges
this tension but the clause text does not resolve it.

**Decision blocked**: which category a base-dir-less separator argv0 reports,
and whether the check even lives at load (where it is unreachable) or in the
binding constructor (where it would be an `execution_failure`, matching C2's
other pre-spawn refusals like the relative-artifact one).

### F4 — Medium — `0025:C5`'s "first match in clause order" precedence is unimplementable as a whole-document guarantee

C5:

> `accessorTable` is fail-fast (one categorized error, never a list), so the
> five are judged in the order declared above and the first match is reported

`internal/table/load.go:942` — `func (l *loader) accessorTable(src
map[string]sourceAcc, ...)` — iterates `for id, a := range src`, a Go **map**,
whose iteration order is randomized per run. It returns on the first defect
(line 950 onward: `return nil, bad(...)`). So with two *different* malformed
entries in one table, which entry's category surfaces is nondeterministic
today.

C5's prose narrows this correctly ("`accessorTable` reports one defect per
entry, so a mutant carrying two defects reports the earlier clause, and each
S1 mutant carries exactly one"), but the normative `precedence:` line does not
carry the per-entry scope, and it is the normative line an implementer builds
to. Worse, `accessorTable` is called three times in sequence (read, write,
gate at `load.go:908/912/916`), so there is a *second*, deterministic
precedence axis (capability table order) that neither line mentions.

**Decision blocked**: whether the intra-entry ordering is the only guarantee
(so a two-bad-entry fixture must not assert a specific category), or whether
this RDR also owes a sort over `src` to make whole-table reporting
deterministic. The second reading is a change to shared loader code that
affects every existing category, not just the five.

### F5 — Low — `0025:C6` does not name the verbs it applies to, and the surface it implies is larger than the gate site it fixes

C6's normative line says `--allow-commands` is "a per-invocation flag on every
verb that can execute an accessor", and the prose fixes the gate *site* in the
command binding's constructor, reached via `flowbind.Registry`.

Grounded: `Registry` today takes exactly `func Registry(m *table.Model)
accessor.Registry` (`internal/cli/flowbind/registry.go`), so "the one new
parameter this clause adds" is accurate and the site is well chosen. But the
set of verbs is not enumerated anywhere in the delta. Candidates in the tree
are `flow next` (`internal/cli/flow_next.go:99`), `flow resolve`
(`flow_resolve.go:99`), `flow read-state` (`flow_state.go:51`), `flow
set-state` (`flow_state.go:171`) — and `intrastate lint`, which C6 explicitly
says must validate *without* the flag, so it must NOT get one.

Unstated defaults an implementer must coin-flip: is the flag declared on each
leaf verb or persistent on `flow`? Is it accepted-and-ignored on verbs that
load a model but execute nothing (which would make the "every verb that can
execute" test ambiguous)? Does passing it when the model declares no `command`
entry succeed silently or warn?

**Decision blocked**: flag registration surface (per-verb vs persistent) and
behaviour on verbs that cannot execute.

---

## Closed since iter-1

The following were the delta's targets and the rewrites do address them; I am
not re-raising them.

- `0025:C1` — the carrier is now fully typed (`[]string`, `*string`, `[]int`,
  `map[string]string`) with the `*string` rationale for `output` stated, the
  home named (`internal/table/source.go::sourceAcc`, which today carries
  Role/Path/Keys/Timeout/ReadBack — confirmed), and the strict-decode
  interaction with `CatUnknownSchemaField` spelled out. `sourceAcc`'s
  pointer-per-field convention is honoured. **Closed.**
- `0025:C2` — the relative/`-`-prefixed artifact refusal is now explicit,
  sited in the binding, with the reason grounded in
  `internal/cli/flow_input.go::parseArtifacts` storing paths verbatim
  (confirmed at `flow_input.go:248`, `out[role] = path`). argv0 resolution has
  a stated base dir and a stated reversal rationale. **Closed** apart from F3.
- `0025:C3` — the raw/json split, the `output` default, the omission-is-
  unreadable rule, the exit maps' empty-stdout precondition, and the over-cap-
  stdout interaction are all now decidable without a coin flip. Verdict
  strings are grounded in `internal/accessor/model.go:256-271`
  (`VerdictAllow`/`Deny`/`Indeterminate`). **Closed.**
- `0025:C4` rendering collision — `internal/cli/flow_exec.go:369
  accessorFailure` confirmed: `ClassExecutionFailure` sets no
  `CLIError.Detail`, and `detailMayHaveApplied` is set only on the two
  read-back/timeout-write arms (lines 378, 397). `clierr.CLIError.Detail`
  exists (`internal/cli/clierr/clierr.go:56`). The precedence rule (applied-
  sense wins the slot, tail renders beneath) is stated and the S3 assertion
  picks a case where they do not compete. **Closed.**
- `0025:A10` — grounded. No `internal/cli/config/` directory exists; no
  `allow_commands`/`allow-commands` hit in any non-test `.go` file. The
  "flag alone is a complete gate" argument holds and the Consequences entry
  names the per-invocation cost honestly. **Verified, closed.**
- `0025:A12` — grounded. `table.Model` (`internal/table/model.go:440`) and
  `table.Accessor` (`internal/table/model.go:144`, fields Role/Path/Keys/
  Timeout/ReadBack) carry no source path, and `flowbind.Registry` receives
  only `*table.Model`. The assumption's premise is exactly right; only its
  consequence for the base-dir-less case is under-specified (F3).
