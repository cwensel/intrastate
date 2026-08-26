# Phase 3a CoVe verification — RDR 0005

Independent spec-derived probing of the shipped `flow` command group. Inputs
were derived from `docs/rdr/0005-skill-integration-cli-contract.md` and
`artifacts/req-list.md` only; the Phase 1 suite (`*_0005_test.go`) was not
read. Every entry below was reproduced by running the built binary.

Probed populations: refusal identity across the whole 25-row stable-code
table, the exit-2/exit-3 boundary, `findings[]` shape and `omitempty`,
byte-fidelity of `<`/`>`/`&` through plan → request → read-back, text/JSON
parity, gate disposition on both verbs, read-accessor narrowing, empty-set
vs clear, and determinism.

---

## FAIL-1 — `flow resolve` drops the results of every non-refusing gate on the selected row

**REQ cited**: REQ-50 (NC, `0005:C1`) — "all gates on the row MUST run, deny
MUST override allow and indeterminate, and **every gate result MUST be
reported**." Reinforced by REQ-47 (FM) — "**no gate result is ever dropped
silently**", stated as binding in both verbs — and fixed in mechanism by
`Technical Design`, *Gates*: "every gate's result is reported, as `gates[]`
on a success or as **one `findings[]` entry per gate** on `flow-gate-denied`
/ `flow-gate-indeterminate`." REQ-51/REQ-97/REQ-98 repeat the carrier as
"`findings[]` one per gate".

**Exact failing input** (self-contained; run from a scratch directory):

```sh
cat > minimal.toml <<'EOF'
outcomes = ["go"]
terminal = ["final-ctx"]

[model]
id = "gatereport"
version = 1
description = "every gate result must be reported"

[initial]
status = "Draft"

[tags.status]
provenance = "owned"
kind = "enum"
domain = ["Draft", "Final"]
single_valued = true
required = true

[tags.recognized]
provenance = "recognized"
kind = "enum"
single_valued = true
required = true

[read.r-main]
role = "main"
path = "m.read"
keys = ["status"]
timeout = "2s"

[write.w-main]
role = "main"
path = "m.write"
keys = ["status"]
timeout = "2s"
read_back = true

[gate.g-allow]
role = "main"
path = "flow.gate.allow"
keys = ["status"]
timeout = "2s"

[gate.g-deny]
role = "main"
path = "flow.gate.deny"
keys = ["status"]
timeout = "2s"

[context.draft.match.status]
eq = "Draft"

[context.final-ctx.match.status]
eq = "Final"

[[rule]]
id = "go-rule"
use = ["draft"]
source = "m:go"
gate = ["g-allow", "g-deny"]
[rule.match.recognized]
eq = "go"
[rule.write]
status = "Final"
EOF
echo '{"status":"Draft"}' > artifact.json

go run ./cmd/intrastate flow resolve \
  --model minimal.toml --artifact main=artifact.json --outcome go --as=json
```

**Observed behaviour** (exit 2):

```json
{"code":"flow-gate-denied","message":"a gate on the selected rule denied the transition","findings":[{"code":"flow-gate-denied","message":"the gate `g-deny` denied the transition: the gate flow.gate.deny denied the transition","param":"g-deny"}]}
```

`g-allow` ran and returned `allow`. Its result appears nowhere in the
envelope — not in `findings[]`, not in `detail`, and not in text mode:

```
error: flow-gate-denied: a gate on the selected rule denied the transition
  flow-gate-denied: the gate `g-deny` denied the transition: ... (param="g-deny")
```

The same drop occurs on the indeterminate arm. With `gate = ["g-allow",
"g-ind"]` on the selected row, `flow-gate-indeterminate` carries one finding
for `g-ind` and none for `g-allow`.

The contrast on the same model and the same row makes the divergence
unambiguous — `flow next --evaluate-gates` reports **both**:

```sh
go run ./cmd/intrastate flow next \
  --model minimal.toml --artifact main=artifact.json --evaluate-gates --as=json
```

```json
"gates":[{"id":"g-allow","result":"allow"},
         {"id":"g-deny","result":"deny","reason":"..."}]
```

**Expected behaviour per the spec quote**: on `flow-gate-denied` and
`flow-gate-indeterminate`, `findings[]` must carry **one entry per gate on
the selected row** — the allowing gates included — so that "every gate result
MUST be reported" (REQ-50) and "no gate result is ever dropped silently"
(REQ-47) hold on the refusal path exactly as they already do on the success
path (`gates[]`) and on `flow next`'s candidate `gates[]`.

**Phase 3c status: FIXED.** `gateVerdictFailure` now builds `findings[]`
from the FULL `[]gateResult`. Precedence is unchanged and selects the
ENVELOPE code only; each finding's own `code` names that gate's own
disposition. Regression:
`TestFail1_EveryGateOnTheSelectedRowIsReportedOnAGateRefusal`
(`internal/cli/flow_adversarial_0005_test.go`), whose non-circular oracle is
`flow next --evaluate-gates` over the same row.

**Where it was**: `internal/cli/flow_exec.go::gateVerdictFailure` partitions
the results and passes only the `denied` slice (or only the
`indeterminate` slice) to `gateRefusal`, which builds one finding per
element of that slice. The full `[]gateResult` from `runGates` is available
at the call site in `internal/cli/flow_resolve.go` and is discarded when
`gateVerdictFailure` returns non-nil.

Not recorded in `artifacts/deviations.md` (DEV-1…DEV-6 do not cover it).

---

## Clean legs — probed and found conforming

No violation was observed on any of the following. Each was run against the
built binary with inputs derived from the spec, not from the test suite.

**Refusal identity — the full input family.** Every code came back with the
exact spelling and the exact carrier the FM table names, at the group's exit:
`flow-tag-owned` (REQ-83), `flow-tag-reserved` (REQ-82),
`flow-tag-duplicate` (REQ-81), `flow-tag-invalid` for an empty value, a
scalar key handed an array, and an absent/empty `--outcome` (REQ-80,
REQ-54), `flow-artifact-invalid` (REQ-88), `flow-artifact-missing` with
`param` = role (REQ-89), `flow-model-not-found` on all four arms — neither
flag, both flags, an unresolvable `--flow` id, and an unreadable `--model`
path (REQ-22, REQ-90), `flow-model-invalid` with a per-category finding
carrying `code` = category slug and `locator` = `file:line` (REQ-24,
REQ-91), `flow-write-invalid` on `<clear>` with the hint "use `--clear …`"
(REQ-62, REQ-84), `flow-write-duplicate` across both flags and twice under
either (REQ-63, REQ-85), `flow-write-unbound` / `flow-clear-unbound`
discriminated by the arriving flag (REQ-86, REQ-87), and every wrong-kind
`--write` value.

**Kernel refusal mapping (REQ-4, REQ-49).** `flow-unmodeled-outcome` (`param`
= the outcome), `flow-no-match`, `flow-ambiguous-match` naming **every**
conflicting row rather than one, `flow-owned-state-unavailable` with one
finding per missing key, and `flow-guard-unevaluable` with the atom's
`key` / `operator` / `literal` / `block` **flat** on the record — never
nested under an `atom` object (REQ-15, REQ-96).

**Exit-code mapping (REQ-19, REQ-20, REQ-21).** The exit-3 population is
exactly the five named classes. Verified at exit 3: read accessor execution
failure during `next`, `resolve`, and `read-state`; gate execution failure
(`flow-accessor-failed`, not a deny — REQ-47, REQ-55); `flow-read-incomplete`;
and `flow-write-readback-incomplete` carrying the "may have been applied"
detail (REQ-104). Verified at exit 2: every input, model, kernel, and gate
refusal, including `GroupInternal`.

**Gate disposition (REQ-43, REQ-46, REQ-50, REQ-51, REQ-116).** Deny
overrides indeterminate on `resolve` (deny + allow + indeterminate on one row
yielded `flow-gate-denied`, not `flow-gate-indeterminate`).
`flow-gate-indeterminate` fired only with no deny present. On `flow next
--evaluate-gates` a deny rides the candidate and the command still exits 0.
A gate could not be laundered into an escape: RDR 0002 refuses an escape rule
carrying a gate list at load (`malformed_escape_declaration`), so the
deny→`no_match`→rescue path is structurally closed.

**REQ-42 — gates of guard-excluded rows.** A row whose `guard.all` is decided
false over the read state, carrying an unreachable gate, was neither reported
nor gated: `flow next --evaluate-gates` exited 0 rather than exit 3.

**REQ-35/36/37 — narrowed read-accessor set.** A model declaring a second
reader (`r-side`, role `side`, serving an owned key no row requires) with its
role left unbound: `flow next` and `flow resolve` both succeeded, invoked only
`r-main`, and raised no `flow-artifact-missing`; `flow read-state` on the same
model invoked `r-side` and refused `flow-artifact-missing` with `param` =
`side`. The narrowing is discriminating, not vacuous.

**REQ-70/71/73 and the MVV encoder leg — byte fidelity.** A set literal
`["a<b","p>q","x&y"]` rendered as itself, never `<` / `>` /
`&`, at every site: the `resolve` plan's `writes`, the `set-state`
request echo, the on-disk artifact, the `read-state` read-back, and text
mode. Plan → request → read-back agreed byte-for-byte. Canonicalization
(sorted, duplicate-free, compact) was applied to `--tag` and `--write` set
values: `["b","a","a"]` → `["a","b"]`.

**REQ-66/107 — empty set vs clear.** `--write 'labels=[]'` reads back as the
present key `[]`; `--clear labels` reads back absent from the reader's `tags`
while remaining in its declared `keys`.

**REQ-77 — read-state payload.** Each reader reports its declared `keys`
beside the tags it returned, so "absent from the artifact" and "not
requested" are distinguishable from the payload alone.

**REQ-34 — `--tag` on `set-state`.** An observed `--tag` was carried as
context and never written to the artifact; an owned key through `--tag` was
still refused `flow-tag-owned`.

**REQ-6/7/11/13/120 — envelope and mode parity.** JSON success emitted
exactly one stdout terminal record and nothing on stderr. Text mode's leaf
set was **exactly** the JSON payload's leaf set on `flow next` — no invented
field, no dropped field. Text-mode failures rendered every finding one per
line, with the same finding-code set as JSON.

**REQ-8/REQ-16 — `findings` `omitempty`.** A single-subject refusal
(`flow-artifact-invalid`) marshalled to `{code, message, param}` with the
`findings` key **absent**, not present-and-empty.

**REQ-52 — escaped plan.** An escape row rescuing `no_match` exited 0 with
`escaped: true` and `escape_class: "no_match"`; `escape_class` was absent
(`omitempty`) on a non-escaped plan.

**REQ-110/111 — determinism.** Three identical invocations of `flow resolve`
(gate-deny refusal) and of `flow next` (success) each produced one distinct
output — same refusal identity, same payload, byte-for-byte.

**REQ-12 — `--as` inherited unchanged.** `--as=yaml` was refused
`flag-invalid-value` from `respond.ValidateMode`, at the root contract's
spelling, before any verb work.

---

## Observations recorded, not raised as failures

- **`revision` is the model id.** `flowRequest.revision()` returns
  `table.Model.KernelTable().Revision`, which RDR 0002 sets to `m.ID`
  (`internal/table/model.go:404`). RDR 0002's layout declares no `revision`
  field, so this is the loader-produced revision identity carried through
  verbatim; the CLI derives, hashes, and synthesizes nothing (REQ-25). The
  clause's "a model that declares none renders `revision` empty" cannot be
  exercised, because 0002 always produces one.

- **`flow-model-invalid` always carries exactly one finding at `file:1`.**
  `table.Failure` is fail-fast and carries no line number
  (`internal/table/category.go`), so "one entry per category hit" is
  satisfiable only with one hit and `locator` is `path:1`. Peer-imposed, not
  CLI drift.

- **`resolve`'s `next` tag-set carries the `<clear>` sentinel** for a cleared
  key (`"next":{"prelock_lens":"<clear>"}`) while `writes` / `clear[]` split
  it correctly. REQ-76 defines `next` as "the next tag-set" and REQ-108's
  copy-through binds only `writes` and `clear[]`, both of which are clean, so
  no transcription path exposes the sentinel to `--write`.

- **`--clear ""` is silently ignored** rather than refused. `parseWrites`
  carries an empty-key arm, but pflag's `StringArray` does not deliver the
  empty element, so the arm is unreachable. No REQ names this input.

---

## Phase 3b — adversarial review

Reviewer worked from `## Trade-offs > ### Failure Modes` (the 25-row code
table plus the gate-disposition mini-check) and the implementation source,
without reading Phase 3a's findings above.

Tests added: `internal/cli/flow_adversarial_0005_test.go`.

Baseline before this file was added: `go test ./...` already fails
`TestReq104And105_ReadBackFailuresExit3AndSayTheWriteMayHaveApplied`
(`internal/cli`) and `TestReq71And72_EmitJSONRendersAngleAndAmpersandAsThemselves`
plus `TestReq71_TheCanonicalLiteralSurvivesTheJSONEmitSiteByteIdentically`
(`internal/cli/clierr`). Those are pre-existing and untouched. No existing
test was modified or weakened.

### ADV-1 — a guard-excluded row is reported by `next` and its gate is invoked

**Failure mode.** `internal/cli::excluded` (`internal/cli/flow_next.go`) is
the CLI's own row-exclusion filter for `flow next`. It decides a row's guard
false only for an `eq` atom under `guard.all` over a present value. Two
shapes RDR 0002's loader admits fall straight through it:

- a `guard.unless` block — the kernel's row verdict is `all ∧ ¬unless`
  (`0007:C5`/`0007:C6`), so an `unless` conjunct decided TRUE excludes the
  row, but `excluded` skips every non-`all` block outright;
- an `all` atom on any operator other than `eq` (`in`, `contains`, `lt`,
  `lte`, `gt`, `gte`).

In both cases the row is reported as a candidate AND, under
`--evaluate-gates`, its gate accessor is INVOKED. Verified against the
kernel over the same state: `flow resolve` refuses `flow-no-match`, so the
row is guard-excluded as a matter of the model, not of the test's opinion.

**Failure Modes row cited.** The gate-disposition mini-check, row
"gate on a guard-excluded row" → `next --evaluate-gates`: "not run, not
reported". Also REQ-42 ("a row whose guards already exclude it is not a
candidate, so its gates MUST NOT run") and REQ-130 (`0005:S8`).

**Why the Phase 1 suite misses it.** `flowGatedNextModel`'s excluded row
uses `guard.all.flag eq "false"` — the one shape `excluded` does decide, so
`TestReq42And46_...` passes over an implementation that decides nothing else.

**Test added.** `TestAdv1_AGuardExcludedRowIsNeitherReportedNorGated`,
sub-tests `unless-block` and `all-block-in-operator`. Each fixture pairs the
excluded row with a gate that DENIES, so a reported `deny` result is direct
evidence the forbidden accessor invocation happened rather than an
inference.

**Currently fails: YES** (both sub-tests; row reported, one `deny` gate
result present on each).

**Phase 3c status: FIXED.** `internal/cli::excluded` now asks the KERNEL for
the row's guard verdict over a one-row probe table, through the same
`guardSeam()` `flow resolve` hands it, so the `all ∧ ¬unless` rule stays RDR
0007's and the typed operator semantics stay RDR 0003's (REQ-113). The probe
strips the row's `Match` pattern (selection is not `next`'s job) and its
escape list. Only `no_match` excludes; `guard_unevaluable` and
`owned_state_unavailable` are undecided, not false.

### ADV-2 — the narrowed reader set ignores owned keys a row reads only in its guard

**Failure mode.** `internal/cli::invokedReaders` (`internal/cli/flow_exec.go`)
computes the invoked read-accessor set from `table.Row.RequiresOwned`. RDR
0002 derives that field from the rule's WRITE BLOCK and CLEAR LIST only —
`internal/resolve::Row`'s own doc comment says so and adds that "the guard's
key set is not required to be a subset of this one".

Consequence: a row that GUARDS on an owned key it does not write demands
that key without naming it in `RequiresOwned`. The declared reader serving
it is never invoked, the kernel receives a view missing the fact, and
`flow resolve` refuses `flow-guard-unevaluable` — for an invocation where
the reader was declared, its artifact role WAS bound, and the artifact held
the value. `flow next` reports the same key as an unresolved fact.

This is the FM section's named silent-failure shape inverted: not a failure
laundered into a success, but a decidable transition turned into a refusal
the caller cannot repair, because nothing in the refusal names the reader
that was skipped.

**Failure Modes row cited.** kernel row `flow-guard-unevaluable`
(`GroupUserEnv` / 2, `findings[]` rows + atoms) — reached here without a
genuine guard-decidability defect. Anchored on REQ-35 ("exactly those
readers serving an owned key some candidate row of the requested model
requires"), REQ-30, REQ-44, and REQ-106.

**Why the Phase 1 suite misses it.** `flowMVVModel`'s `read.orphan` fixture
proves the NEGATIVE half of REQ-35 (a reader no row needs must not run).
Every owned key any fixture row's guard reads is also a key that row writes,
so the positive half — a reader that MUST run — is never tested against a
guard-only demand.

**Test added.** `TestAdv2_AReaderServingOnlyAGuardsOwnedKeyIsStillInvoked`,
sub-tests `resolve` and `next`. Fixture `advGuardOwnedModel` declares
`read.flags` serving exactly one owned key, that key appears only in a
guard, and its role is bound on every invocation — so neither
`flow-artifact-missing` nor REQ-36 narrowing-by-design can explain the
outcome.

**Currently fails: YES.** `resolve` refuses `flow-guard-unevaluable`
(finding: "the atom on `flag` could not be decided (absent)"); `next`
reports `readers = [state]` and `flag` unresolved.

**Phase 3c status: FIXED.** `invokedReaders` now unions `Row.RequiresOwned`
with the row's guard-block (`all` / `unless`) atom keys the model declares
owned. Keys and provenance only — never operators, never literals — so RDR
0003's semantics stay at their seam and both halves remain dumpable normalized
model data (CA A7). No RDR 0002 change. ASSUMPTION A-4 under-read REQ-35 and
is superseded; recorded as **DEV-8**.

### ADV-3 — `revision` renders the model id, which REQ-25 forbids

**Failure mode.** `flowRequest.revision()` returns
`model.KernelTable().Revision`, and `internal/table::Model.KernelTable`
sets `Revision: m.ID`. RDR 0002's `[model]` block admits only `id`,
`version`, `description`, and `metadata` (`internal/table/source.go`) —
there is no revision field, so NO model can declare a revision and REQ-25's
consequent governs every payload this CLI emits:

> "The CLI never derives, hashes, or synthesizes it, so a model that
> declares none renders `revision` empty rather than a CLI-invented value."

Standing the id in there is a CLI-invented value. It also defeats what the
field is named for: REQ-110 makes request identity "(verb, model selection
(`--flow` id or `--model` path) AND model revision, …)", listing the
revision SEPARATELY from the selection — so a revision that is only ever the
id cannot distinguish two revisions of one model, which is the whole
discrimination REQ-110 asks of it.

**Failure Modes row cited.** No FM code row: this is a success-payload
defect, so it produces no refusal at all — which is precisely why it belongs
in an adversarial pass. Anchored on REQ-25 (TD, *Model selection*) and
REQ-110 (`0005:D-identity`), with the payload obligation from
REQ-74..REQ-78.

**Why the Phase 1 suite misses it.**
`TestReq25_RevisionIsNeverCLIDerivedFromPathOrContent`
(`internal/cli/flow_input_0005_test.go:135`) tests only two derivation
shapes — the `--model` PATH and a hex digest — and guards both behind
`if rev != ""`. The model id passes both filters, and the clause the REQ
actually turns on ("declares none → empty") is never asserted.

**Test added.** `TestAdv3_AModelDeclaringNoRevisionRendersRevisionEmpty`.
Runs all four verbs over two documents identical but for `[model].id`, so a
`revision` tracking the id is visibly tracking the wrong identity; asserts
the exact `rev == id` equality rather than mere non-emptiness, so it cannot
be satisfied by an unrelated non-empty value.

**Currently fails: YES** (8/8 sub-tests: `advrevone` and `advrevtwo` × all
four verbs).

**Phase 3c status: FIXED per DEV-7.** `flowRequest.revision()` renders EMPTY.
`internal/table/model.go:404`'s `Revision: m.ID` is UNTOUCHED — that is the
kernel's own table-identity for resolver comparison, a different consumer with
a different meaning — and no `[model].revision` field was added to RDR 0002.

### Attacks the implementation withstood

Recorded so a later pass does not re-run them:

- **Two `[write.<id>]` entries serving one key** (REQ-61's "served by
  exactly one", REQ-86/REQ-87). `internal/cli::writerFor` checks only that
  SOME writer names the key, and `groupByWriter` would silently pick the
  alphabetically-first — but RDR 0002's loader refuses the model first
  (`malformed_accessor_binding`: "written tag X is served by 2 writers;
  want exactly one"), and likewise for two readers over one owned key. The
  CLI's incomplete check is unreachable.
- **Canonical set literal round trip** (REQ-70/71/107/108). An unsorted,
  duplicated caller spelling `["x&y","a<b","a<b"]` is re-canonicalised to
  `["a<b","x&y"]` and survives `set-state` → `read-state` byte-identically
  with `<` and `&` unescaped at every payload site.
- **Empty set vs. cleared key** (REQ-66). `--write labels=[]` reads back as
  the present value `[]`; a `--clear` reads back absent. The distinction
  holds through the artifact format.
- **`escape_class` on a multi-class escape row** (REQ-52, A-3). The
  `return ""` branch in `escapeClassOf` is unreachable: RDR 0002 admits only
  `no_match` and `ambiguous_match` as escapable, and the escape-stripped
  probe necessarily refuses with one of them.
- **A second `recognized`-provenance tag key** (REQ-27, A-2). Refused at
  load — RDR 0008 fences the reserved key in both directions.
- **A gate refusing `gate_indeterminate` under `next --evaluate-gates`**
  (REQ-47). The executor returns `Verdict: indeterminate` alongside the
  refusal, and `runGates` reports it on the candidate rather than dropping
  it; no gate result is lost.


---

## Phase 3c — fixup outcome

All seven defects routed to Phase 3c are closed. `go test ./...`, `go vet
./...`, `gofmt`, and `golangci-lint run` are all clean.

| Item | Kind | Disposition |
| --- | --- | --- |
| FAIL-1 (3a) | implementation | fixed; regression test added |
| ADV-1 (3b) | implementation | fixed; Phase 3b test now green |
| ADV-2 (3b) | implementation | fixed; **DEV-8** records A-4 superseded |
| ADV-3 (3b) | implementation | fixed per **DEV-7** |
| DEV-2 | test oracle | corrected; implementation untouched |
| DEV-3 | test oracle | corrected; implementation untouched |
| DEV-5 | test fixture | corrected; implementation untouched |

Each of the three corrected oracles was verified to still DISCRIMINATE the
behaviour it names, by injecting the defect the REQ forbids and confirming the
oracle fires. The predecessor suites (0001, 0002, 0003, 0004, 0006) are green
and no unrelated test was modified or weakened.
