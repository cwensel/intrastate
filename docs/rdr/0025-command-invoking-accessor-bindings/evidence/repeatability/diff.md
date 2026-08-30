model: claude-opus-5[1m]
pass: repeatability DIFF (post-barrier, no run authored by this pass)

# Repeatability diff — RDR 0025

Runs compared:

- run-1 — `claude-opus-5[1m]`, variant `full` (profile: foundational)
- run-2 — `claude-sonnet-4-5-20250929`, variant `full`
- run-3 — `claude-fable-5`, variant `full`

run-1's `variant:` reads `full`, so three runs is the intended coverage, not a
`lite` over-run. The model boundary for "alt-model" signal is run-3
(`fable-5`) against run-1/run-2; run-2 (`sonnet`) is also off run-1's model, so
a 2-1 split that isolates run-1 is weaker signal than one that isolates run-3.

---

## 1. Disagreements

### D1 — the `WriteBinding` method name and the `GateBinding` return arity (3-way)

**Element**: `0025:§normative-contracts` (no contract owns it) — with a pull on
`0025:C3`, which fixes the gate *envelope* but not the Go return.

- run-1: `Read(ctx, artifact) (map[string]string, error)`;
  `Gate(ctx, artifact) (accessor.Verdict, string, error)` — a **three**-value
  return, verdict + reason + error; `Apply(ctx, artifact, tags) error`.
  Marked GUESS on the signatures, explicitly: "The RDR fixes only the method
  **names** (`Read`, `Gate`, `Apply`)".
- run-2: `Read(ctx, artifact) (map[string]string, error)`;
  `Gate(ctx, artifact) (Verdict, error)` — **two** values, reason dropped
  entirely; and the write method is named **`Write`**, not `Apply`. Presented
  as "existing interfaces, unchanged signatures", *not* marked GUESS.
- run-3: names `Read`, `Gate`, `Apply`, "each returning a bare `error`", and
  declines to write the parameter lists or return tuples at all.

Three renderings of the same seam, and the two that committed to a shape
disagree on both the gate's arity and the write method's name. run-2 is the
outlier that asserted it without a GUESS marker — an unflagged wrong name
(`Write` for `Apply`) is the worst reading of the three, because a downstream
implementer would take it as fixed.

**RDR passage that let them diverge**: the Approach states command bindings
"implement the existing `ReadBinding`/`GateBinding`/`WriteBinding` seam" and
`0025:A11` says `binding.go`'s "`Read`/`Gate`/`Apply` return a bare `error`" —
that is the *only* place the method names appear, inside an assumption's
Evidence prose rather than a normative block. C3 fixes the gate *wire*
envelope (`verdict` + `reason` on stdout) but never says whether `reason`
crosses back through the Go return. Nothing in C1–C6 carries the interface.

### D2 — the new package / file that holds the command binding (3-way)

**Element**: `0025:§normative-contracts` (unowned; C6 fixes only the
construction *site*, not where the binding type lives).

- run-1: a **new package** `internal/cli/cmdbind`, "a sibling package to
  `flowbind`", with three exported constructors
  `NewReader`/`NewGate`/`NewWriter` taking `(acc, baseDir, allowCommands)`.
  Marked GUESS on the name and on the three-constructor split.
- run-2: **no new package** — "new command-carrier support lives beside the
  existing path-carrier (`internal/cli/flowbind`)", with the binding type
  guessed into `internal/accessor/command_binding.go`, i.e. a *different*
  package (`accessor`, not `flowbind`) from the one its own prose names.
  Internally inconsistent, and marked GUESS only on the file name.
- run-3: names no package or file for the binding at all — describes it as
  "the command binding" and puts the load-time validator arms in
  `internal/table/load.go::accessorTable`.

run-1 asserts a sibling package as if the RDR said so ("the RDR says 'a sibling
package to `flowbind`'"); nothing in C1–C6 or the Approach says that. This is
one contract where the RDR's silence produced three different module layouts,
including one self-contradictory one.

### D3 — the `flowbind.Registry` signature (3-way, all GUESS)

**Element**: `0025:C6` (gate site + "signature change on it plus the
`flowRequest` field"), leaning on `0025:A13` for the current shape.

- run-1: `func Registry(m *table.Model, baseDir string, allowCommands bool) accessor.Registry` — return type carried over from A13 correctly.
- run-2: `func Registry(...existing..., allowCommands bool) *FlowbindRegistry` —
  **invents a return type** `*FlowbindRegistry` and marks it GUESS, and omits
  the base-dir parameter entirely. run-2 never threads the model base dir into
  `Registry` anywhere; it appears in run-2 only as an absent concern.
- run-3: `func Registry(m *table.Model, modelDir string, allowCommands bool) (…bindings…, error)` — adds an **error return** the current signature does
  not have, explicitly marked GUESS on order and spelling.

All three marked some part GUESS, so this is also a GUESS cluster (see G1). The
substantive divergence is run-2 dropping the base dir: `0025:A12` is the only
element that says the base dir threads as a parameter, and it says so in
assumption Evidence, not in C6's normative block. C6 says only that
`buildRequest` "reads the flag and passes it to `flowbind.Registry`" — the
flag, singular. A reader who works from C6 alone loses A12's second parameter,
which is exactly what run-2 did.

### D4 — empty stdout, exit 0, no exit-map hit, on a **json read** (2-1, isolates the alt-model run-3)

**Element**: `0025:C3` / `0025:C4` (the "parse the envelope, THEN classify the
exit code" order; C3's `exit_absent` clause).

- run-1: `execution_failure`. Marked GUESS — "C3 states it explicitly only for
  `output='raw'` … by C4's 'any unlisted exit … is execution_failure' it
  generalizes".
- run-2: `execution_failure` (falls through the `switch` to the terminal
  `refusalWithDetail("execution_failure", ...)`). Not flagged.
- run-3: **`return UNREADABLE per omitted key`** — an explicit
  `if exit == 0 and read(json): return UNREADABLE per omitted key` arm that
  neither other run has. run-3 flags it in its GUESS ledger and says it read
  the gate case as `execution_failure` but the read case as UNREADABLE.

This splits on the model boundary and it is a semantic split, not a naming one:
run-3's arm makes a silent tool (exit 0, no output) yield "every declared key
unreadable" where the other two refuse. Both are defensible from the text,
which is the problem.

**RDR passage**: C3's read/json clause governs "a declared key the object
omits" — it presumes an object was parsed. It says nothing about *no object at
all*. C4's ordering clause routes empty stdout to the exit maps but does not
state the fallthrough when the map does not match on a read. C3 states the
`execution_failure` fallthrough only for `output = "raw"`.

### D5 — where the `--allow-commands` refusal is *decided*: construction vs. invocation (2-1, isolates run-2)

**Element**: `0025:C6` ("gate site: ONE — … `flowbind.Registry` … which builds
refusing command bindings when it is unset").

- run-1: the constructor is the gate site but returns a **refusing binding**,
  never an error — "the refusal is deferred to invocation so that lint and load
  stay ungated". Marked GUESS.
- run-3: same — "gate sits in the constructor", refusal surfaces at invoke.
- run-2: the gate is checked **inside `invoke`** on a `b.allowCommands` field,
  and run-2 separately locates the gate read in
  `flow_exec.go::buildRequest` as "the sole place `--allow-commands` is read".
  Functionally near-equivalent, but run-2 never renders the "refusing binding"
  object C6 names, so C6's fail-closed-by-construction claim has no artifact in
  run-2's reconstruction.

Weaker than D1–D4 (the behaviour converges), but it shows C6's phrase "builds
refusing command bindings" does not by itself force a reader to produce a
refusing-binding type.

### D6 — the carrier-less-entry residue: what object it is and when it refuses (2-1, isolates run-2)

**Element**: `0025:C1` (mutual exclusion) with `0025:A14` (the `Path: ""`
hazard).

- run-1: registry-level `refusingBinding(execution_failure, Detail="malformed
  entry: no carrier")`, refusal deferred to invocation, GUESS-marked on the
  deferral.
- run-3: same shape — "carrier-less entry reaching the constructor (refusing
  binding, never a `Path: ""` file binding)", listed under `execution_failure`.
- run-2: "an entry declaring neither is a REFUSING binding at runtime, never a
  `Path: ""` file binding **(S7)**" — cites the scenario rather than C1/A14,
  and never places the residue in its registry/constructor pseudo-code. run-2's
  §5 pseudo-code has no carrier-selection step at all: it starts inside an
  already-built `commandBinding`.

**RDR passage**: C1's normative block is the TOML grammar only; the
mutual-exclusion *rule* and the residue arm live in C5
(`command_and_path_conflict`, load-time) and in A14's Evidence (runtime,
in-memory models). No normative clause states the runtime residue behaviour —
it is only in an assumption's Evidence and in a scenario.

### D7 — `table.Categories()` total after the six additions (isolated to run-2, but worth naming)

**Element**: `0025:C5` (registration clause).

- run-1 and run-3: state only that the six are appended at the tail in the
  declared order. Neither claims a total.
- run-2: "`accessor.ValidationCodes` count stays at 8 total categories
  (existing 2 + these... GUESS on the prior count's composition)". This is a
  fabricated arithmetic claim — it also renames the list from
  `table.Categories()` to `accessor.ValidationCodes` mid-sentence.

Not a divergence the RDR caused so much as one it failed to foreclose: C5 says
the hand-maintained list "IS the closed set" but never states its current size,
and run-2 filled the hole with a number.

### D8 — non-Unix platform detection predicate (2-1, isolates run-2)

**Element**: `0025:C4` (`platform:` clause).

- run-1 and run-3: "not Unix" / "goos != unix-like", left as the RDR frames it.
- run-2: hardcodes `if goos != "linux" && goos != "darwin"` — an enumeration the
  RDR never gives, and one that refuses on freebsd/openbsd, which do support
  the process-group mechanism C4 relies on.

Minor, but it is a real behavioural difference produced by C4 saying "non-Unix"
without naming the set.

---

## 2. GUESS clusters (candidate RDR rewrites)

Ordered by how many runs marked the same contract GUESS.

### G1 — `flowbind.Registry`'s post-change signature — **3 runs GUESS**

**Element**: `0025:C6`, with `0025:A12` / `0025:A13` supplying the facts.

run-1 GUESSes "the new parameter order"; run-2 GUESSes "the exact `Registry(...)`
parameter list beyond the new `allowCommands` bool" and the return type; run-3
GUESSes "the exact Go signature and parameter order". All three needed a
signature the RDR does not carry in any normative block. This is the strongest
rewrite candidate: C6 should state the post-change `Registry` signature
(parameters and order, including the base dir A12 threads) in its normative
block rather than leaving it split across two Pending assumptions' Evidence
prose.

### G2 — helper/function names and the internal decomposition — **3 runs GUESS**

**Element**: `0025:§normative-contracts` (unowned).

run-1 GUESSes `buildArgv`, `childEnv`, `run` "and that substitution and argv0
resolution live in one helper rather than two"; run-2 GUESSes the
command-binding constructor's file and `readOutcome`/`classify`'s signature;
run-3 GUESSes "the internal helper names (`invoke`, `classifyResult`)". The
three decompositions are genuinely different (run-1 splits argv/env/run into
three; run-3 folds all of it into one `invoke` plus a classifier; run-2 puts
argv+env+spawn in the constructor). If the RDR does not care, that is fine —
but §2 of every run had to invent structure, so the RDR should either name the
seam or say explicitly that the decomposition is unconstrained.

### G3 — `accessor.ExecError`'s shape beyond `Detail string` — **3 runs GUESS**

**Element**: `0025:C4` (`detail:` clause).

run-1: GUESSes "the `Error() string` method body and whether `ExecError` wraps
an inner error (C4 says `errors.As`, which needs no wrapping)". run-3: GUESSes
"it also wraps or stores an underlying cause and implements `Error() string`".
run-2 does not render `ExecError` as a type at all — it lists `Detail` on
`Refusal` and never declares the typed error C4 names, which is arguably worse
than a GUESS. C4 gives the field and the `errors.As` mechanism; it should also
state whether the type wraps a cause, since that decides whether
`errors.Is(err, exec.ErrNotFound)` survives the trip.

### G4 — Go identifier spellings for the six new `Category` constants — **2 runs GUESS**

**Element**: `0025:C5`.

run-3 GUESSes `CatCommandAndPathConflict` and names the `Cat…` family
precedent; run-1 works the same precedent implicitly (it names
`CatUnknownSchemaField` for the strict-decode case) but treats the strings as
the contract. run-2 renders the strings only. Low-stakes if the strings are the
contract — but C5 should say which of the two is normative, because run-3 read
the constants as part of the surface and run-2 did not.

### G5 — `sourceAcc` Go field names and struct tags — **2 runs GUESS**

**Element**: `0025:C1`.

run-1 renders the full struct with tags and marks it GUESS: "C1 gives the TOML
spellings only". run-3 says `sourceAcc` "gains exactly the six fields above with
the Go types shown" without spelling them. run-2 never renders the struct. C1
gives TOML keys plus Go types; the Go field names are absent and one run needed
them.

### G6 — empty-stdout/exit-0 read classification — **2 runs GUESS** (see D4)

**Element**: `0025:C3` / `0025:C4`.

run-1 and run-3 both mark this GUESS, and they resolve it *differently*
(`execution_failure` vs `UNREADABLE`). A GUESS cluster where the guesses
disagree is the highest-value rewrite in this set: the fix is one clause in C3
or C4 stating the read-side fallthrough for empty stdout with no `exit_absent`
match, matching the sentence C3 already has for `output = "raw"`.

---

## 3. Agreement

All three runs rendered identically, with no GUESS and no variance: the six
TOML carrier fields and their Go types (`command []string`, `output *string`
read-only, `exit_absent []int` read-only, `exit_verdicts map[string]string`
gate-only, `env map[string]string`, `env_pass []string`) with `output`'s
pointer-means-omitted rationale; the six C5 category strings and their
conditions in declared tail order; C5's three-level precedence (fail-fast in
clause order within an entry, unspecified across entries in a table,
read→write→gate across tables); C3's complete wire format in both directions
(write stdin = planned tags with `<clear>` as a literal, read/gate stdin = `{}`
always sent, read json = flat map with omission meaning UNREADABLE, read raw =
one trailing newline stripped, gate = verdict/reason envelope, exit maps
consulted only on empty stdout); C4's four-layer env with later-wins ordering,
`LC_` as the one literal prefix rule, `env_pass` whole names only, and the
four `INTRASTATE_*` overlay vars; C4's deadline triple verbatim (`Setpgid`,
`Cancel` = SIGKILL to `-pgid`, `WaitDelay = 500ms` bounding stdin write and
drain) with timeout from `ctx.Err()` and never the wait error; the 1 MiB stdout
cap and 4 KiB stderr tail; the envelope-before-exit-code ordering; C4's
injectable `goos` var with a runtime check rather than a build constraint;
`--allow-commands` as v1's only opt-in, per-invocation, never in the model,
with `lint` ungated; and identity staying `(flow, name, capability)` with
carrier kind outside it. These confirm the RDR is determinate across the whole
wire protocol, the env policy, the deadline mechanism, and the load-time
category set.
