model: Claude Sonnet 5 (claude-sonnet-5)
variant: full (profile: foundational)

# Repeatability reconstruction — RDR 0023 (Resolve-envelope projection opt-out)

Grounded from `elements[]` (kind C, MVV, S) in
`docs/rdr/0023-resolve-envelope-projection.md`, widened into
`§approach`, `§illustrative-code`, `§load-bearing-decisions`,
`§minimum-viable-validation`, `§testing-strategy`, and
`§oracle-discriminability` because C1/C2 leave several signatures,
the exact carrier shape, and step ordering under-determined as bare
contract prose (noted inline as GUESS where the record is genuinely
silent).

## 1. Public API

```go
// internal/cli/flow_resolve.go (existing command, extended)

// registerSelectionFlags is NOT touched — plan-only is resolve-local.
func newFlowResolveCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "resolve",
        Short: "...",
        Long:  flowResolveExtendedDesc, // already documents the "plan" group
        RunE:  runFlowResolve,
    }
    // ... existing flags: --model, --outcome, --tag (repeatable), --as ...
    cmd.Flags().Bool("plan-only", false, "report only the plan group (rule, gates, emit, next, writes, clear, escaped, escape_class, revision); omit the request-echo group")
    // long form only — no shorthand (C1, D-naming)
    return cmd
}

func runFlowResolve(cmd *cobra.Command, args []string) error {
    planOnly, _ := cmd.Flags().GetBool("plan-only")
    // ... existing resolve/plan construction, unchanged ...
    payload := buildResolvePayload(result, req, planOnly)
    return respond.OK(cmd, payload) // projection applied BEFORE respond.OK (C1)
}
```

Types (GUESS on exact field naming/tags beyond what C1/C2/A2 state;
the record fixes behavior, not Go syntax):

```go
// resolvePayload is the success-payload struct rendered by both
// output modes (json/text) via the shared generic renderer.
type resolvePayload struct {
    // ECHO group (JDR 0002 §D1) — omitted (absent, not null) under --plan-only
    Model    *string          `json:"model,omitempty"`
    Observed *observedView    `json:"observed,omitempty"`
    Owned    *map[string]any  `json:"owned,omitempty"`
    Readers  *[]string        `json:"readers,omitempty"`
    Outcome  *string          `json:"outcome,omitempty"`

    // PLAN group — always carried under --plan-only, byte-identical to default
    Revision    string          `json:"revision"`
    Rule        string          `json:"rule"`
    Gates       []gateResult    `json:"gates"`
    Emit        map[string]any  `json:"emit"`               // 0010:C4: never omitted, renders {} when empty
    Next        map[string]any  `json:"next"`
    Writes      map[string]any  `json:"writes"`
    Clear       []string        `json:"clear"`
    Escaped     bool            `json:"escaped"`
    EscapeClass *string         `json:"escape_class,omitempty"` // 0005:A-3 presence rule, unchanged by this RDR
}
```

Error modes:
- `flow next|read-state|set-state --plan-only` → `command-error`
  (existing refusal family), exit code 2 (MVV step 5, F1). No new
  refusal code or exit group is minted (C1).
- Refusal envelopes (`clierr.CLIError` + `findings`) are never
  projected and carry no echo-group member — `--plan-only` cannot
  change a single refusal byte on any verb (A6). Byte-identical
  refusal envelopes and exit codes are asserted ± the flag on
  `flow resolve` itself (S1, MVV step 4).
- No new error is introduced for a malformed `--plan-only` value —
  it is a standard cobra `Bool` flag; a non-boolean literal fails
  cobra's own flag-parsing error, not a resolve-specific one
  (GUESS: the record does not name this path explicitly, but C1
  states "boolean flag" with no bespoke parse-error clause, and no
  other CLI flag in this family is shown to have one).

## 2. Three most important internal helpers

1. **`buildResolvePayload(result, req, planOnly bool) resolvePayload`**
   (name is a GUESS — the record does not name the assembly function,
   only that "the projection is applied to the verb-specific result
   BEFORE `respond.OK`," §approach). Responsibility: construct the
   full payload from the kernel result and request, then — if
   `planOnly` — nil out the pointer-valued echo fields (A2's
   mechanism: pointer-valued echo fields with `omitempty`, always
   populated in default mode, nilled at projection) so the encoder
   drops them as absent keys, never `null`. This is the single seam
   both output modes flow through, which is why JSON and text can
   never disagree (0005:C1 held by construction).

2. **`walkCommandTree(root *cobra.Command) map[string]bool` — the
   NEW total walker** (S4, A5). Responsibility: descend every child
   of the root with no name-based skip and no `Hidden` gate,
   including auto-generated `help` and `completion`, and return the
   set of commands registering `plan-only`. It deliberately does NOT
   reuse `internal/cli/help_all.go::walkCommandTree`, which skips
   children named `help`/`completion` and gates on `Hidden` — a
   strictly weaker negative than C1 requires. Its caller must first
   materialize `help` and `completion` via cobra's own initializers
   (`InitDefaultHelpCmd`, `InitDefaultCompletionCmd` from `ExecuteC`)
   and assert both are present in the walked set as a precondition,
   or the walk can pass vacuously over a tree that never contained
   the command class the clause exists to reach.

3. **Reader-execution recorder (name unspecified — GUESS, e.g.
   `recordReaderInvocation` or a counting seam around
   `flow_exec.go`'s reader pass)**. Responsibility: record, as a
   side effect at the actual reader invocation site, which readers
   ran during a given resolve execution, so S1's invoked-reader
   assertion compares OBSERVED execution across the default and
   projected runs rather than recomputing `invokedReaders(model,
   outcome)` — a pure function of arguments the flag never touches,
   which would make the assertion unfailable/vacuous by construction.
   A projected run that (incorrectly) skips the reader pass must turn
   this assertion red.

## 3. Data model (persisted / boundary-crossing)

Nothing new is persisted; the boundary is the CLI's stdout wire
format (NDJSON `{"type":"ok","data":{…}}` line) and its text
rendering. What crosses the boundary differently under the flag:

- **ECHO group** (input side — supplied by the caller or assembled
  on its behalf): `model` (model reference, string), `observed`
  (`--tag` set, echoed verbatim except a value-preserving
  canonical-set re-encode), `owned` (assembled owned view, prior-owned
  overlaid with `next`), `readers` (invoked reader identities, a pure
  function of model + outcome), `outcome` (`--outcome` string, echoed
  verbatim).
- **PLAN group** (kernel output side + bounded identity
  attestations): `rule` (matched rule id), `gates` (gate results,
  `[]`), `emit` (authored answers, `{}` — never omitted, 0010:C4),
  `next` (`{}`), `writes` (`{}`), `clear` (`[]`), `escaped` (bool),
  `escape_class` (presence-rule field per 0005:A-3 — omitted when
  unescaped or unprobeable), `revision` (loader-produced identity
  token, plan-side by definition, presently always `""` since 0002's
  `[model]` block admits no revision key — contracted but vacant).
- **Assignment carrier** (C2, load-bearing): a standalone table keyed
  by field name, in its own declaration site — NOT a per-field struct
  tag/marker on `resolvePayload` — mapping each field to `echo` or
  `plan`. This independence is structural: A2's mechanism edits the
  echo fields' types/tags on the struct itself, so a marker co-located
  there would be edited in the same commit that changes the
  projection, making the completeness oracle (S3) tautological.
- **Wire format rule**: projection is key deletion only, never
  re-encoding — projected payload = default payload minus echo keys,
  byte-for-byte identical on every surviving field, surviving keys
  keep default-mode declaration order, same non-HTML-escaping encoder.
  Absent means absent (no key at all) — never `null`, `{}`, or `""`.
- **Text mode**: same projected struct flows through the existing
  generic renderer (`respond/text.go::flatten`, which sorts keys from
  the decoded map — deterministic, no map-iteration-order path).
  Projected text lines are a byte-identical subset (per-line) of the
  default-mode text lines.

## 4. Top-level pseudo-code of the main operation

```
function runFlowResolve(cmd, args):
    model      := loadModel(cmd.Flag("model"))
    tags       := parseTags(cmd.Flags("tag"))       # refuses owned/reserved/malformed
    outcome    := cmd.Flag("outcome")
    planOnly   := cmd.Flag("plan-only")              # NEW, default false

    result, err := kernel.Resolve(model, Input{Owned: ownedView, Observed: tags, Outcome: outcome})
    if err != nil:
        # refusal path — CLIError + findings, NO echo-group member,
        # therefore unaffected by planOnly (A6) — respond.Fail(cmd, err)
        return respond.Fail(cmd, err)

    readers := invokeReaders(model, outcome)          # side-effecting; recorded for S1's oracle
    payload := resolvePayload{
        # ECHO group (input side)
        Model:    &model.Reference,
        Observed: &observedView{tags},
        Owned:    &ownedView,
        Readers:  &readerIDs(readers),
        Outcome:  &outcome,

        # PLAN group (output side + identity)
        Revision:    loader.Revision(),                # "" today — 0002 admits no revision key
        Rule:        result.RuleID,
        Gates:       result.Gates,
        Emit:        result.Emit,                      # always {} at minimum, never omitted
        Next:        result.Next,
        Writes:      result.Writes,
        Clear:       result.Clear,
        Escaped:     result.Escaped,
        EscapeClass: escapeClassOf(result),             # presence rule: omitted when unescaped/unprobeable
    }

    if planOnly:
        # Apply BEFORE respond.OK so json/text can never disagree.
        # Nil every pointer-valued ECHO field so `omitempty` drops
        # the key entirely (absent, never null) — mechanism from A2.
        payload.Model    = nil
        payload.Observed = nil
        payload.Owned    = nil
        payload.Readers  = nil
        payload.Outcome  = nil
        # PLAN-group fields (incl. revision, always-keep core) untouched.

    # Both output modes render the SAME payload through one path —
    # JSON: default encoder, HTML-escaping disabled, declaration order.
    # Text: respond/text.go::flatten (sorted keys, deterministic).
    return respond.OK(cmd, payload)
```

Notes:
- `readers`/`invokeReaders` execution happens regardless of
  `planOnly` — the flag is report-only, so the decision path (rule
  selection, gate evaluation, reader invocation) is identical across
  ± flag; only the rendered payload's key set differs (C1 Identity:
  request identity, not decision identity).
- `--plan-only` is registered ONLY on `resolve`'s own command, never
  on shared registrars (`registerSelectionFlags`, the flow group's
  persistent flags, or the root's persistent set) — enforced by the
  S4 total-tree walk, not by convention alone.
