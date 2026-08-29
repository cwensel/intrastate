model: claude-fable-5 (Fable 5)
variant: full (profile: foundational)

# Repeatability run 3 — RDR 0023 (Resolve-envelope projection opt-out)

Sources read via the projector: elements C1, C2, MVV, S1–S5 (the
`elements` dump also surfaced A1–A9 in full). Widened deliberately past
the contract spans into: §approach, the §technical-design preamble,
§load-bearing-decisions, §disposition, §oracle-discriminability,
§illustrative-code, §phase-1-projection-on-the-verb, and
§phase-2-oracles. Widening was needed to fix the projection site
(before `respond.OK`, success path only), the 14-field struct pin, the
flag's naming rationale, and the wire example.

## 1. Public API

The module's public surface is the CLI verb surface plus the wire
contract; no new Go package API is exported.

### CLI flag

```
flow resolve [--model …] [--outcome …] [--tag k=v]… [--as=json|text] [--plan-only]
```

- `--plan-only`: boolean, default `false`, long form only (no
  shorthand). Registered on `flow resolve` ONLY — never on the shared
  registrars (`registerSelectionFlags`), the flow group's persistent
  set, or the root's persistent set (C1).
- Any other command given `--plan-only` → cobra parse failure: exit 2,
  `command-error` shared bucket. No new refusal code, no new exit
  group is minted (C1, §disposition).
- A pre-RDR binary given the flag → exit 2, `command-error`, so
  version skew is loud: an absent `observed` always means projection,
  never an old binary (§disposition, F5).

### Success output

- Default mode (`--plan-only` absent): byte-identical to today's full
  payload — this RDR changes no default byte.
- Under the flag, the success payload carries exactly the PLAN group:
  `revision`, `rule`, `gates`, `emit`, `next`, `writes`, `clear`,
  `escaped`, plus `escape_class` exactly when the same request's
  default output carries it (producer's own presence rule,
  `flow_resolve.go::escapeClassOf`; `emit` stays present as `{}` per
  0010:C4). It MUST NOT carry the ECHO group: `model`, `observed`,
  `owned`, `readers`, `outcome`. Omitted keys are ABSENT — never
  `null`, never `{}`/`""` placeholders.
- Every carried field renders byte-identically to its default-mode
  rendering (same encoder, HTML escaping disabled); carried keys keep
  default-mode relative declaration order.
- Width: the projected FULL EMITTED LINE (`{"type":"ok","data":{…}}`
  NDJSON record, trailing newline excluded) is STRICTLY SHORTER than
  the default's, in BOTH output modes (text half rests on A8,
  Pending; if refuted the clause narrows to JSON before lock).
- Text mode (`--as=text`): projected lines are a byte-identical
  SUBSET (set membership, not subsequence — text is sorted by
  `respond/text.go::flatten`) of the default run's lines, stable
  across runs, and strictly fewer total rendered bytes.

### Error modes

- REPORT-ONLY, oracle-enforced: the flag never changes what is
  decided or how the run fails. Same request ± flag: identical exit
  codes, byte-identical refusal envelopes, identical OBSERVED
  invoked-reader execution.
- Refusals render `clierr.CLIError` (`code`, `message`, `param`,
  `detail`, `hint`, `findings`) via `respond.Fail`; there is no echo
  member to project, so refusal flag-blindness holds by construction
  (A6) and structurally — the projection site sits after the last
  `respond.Fail` return (Phase 1).

## 2. Three most important internal helpers

1. **The projection function** — GUESS at name/signature:
   `projectPlanOnly(p resolvePayload) resolvePayload` in
   `internal/cli/flow_resolve.go`. Responsibility: given the fully
   assembled default payload, return the same value with the five echo
   fields nilled (A2 mechanism: echo fields are pointer-valued with
   `omitempty`, always populated non-nil in default mode, set to nil
   here so the keys drop from the wire). Implemented FROM the declared
   assignment table (C2): it nils every field the table marks `echo`,
   never from a hard-coded omit-list of its own. Called once, on the
   SUCCESS path only, after payload assembly and before `respond.OK`;
   the flag is read nowhere else (no gate/rule/refusal path reads it).

2. **The partition declaration table** — GUESS at name/shape:
   `resolveFieldGroups`, a standalone `map[string]fieldGroup` (or
   ordered table) keyed by wire field name, one entry per
   `resolvePayload` field, each carrying `echo` or `plan`. C2
   constrains the carrier: it MUST be a standalone table in its own
   declaration site, NOT per-field struct markers on `resolvePayload`
   (a marker there is edited by the same commit that retypes the echo
   fields — the tautology S3 exists to forbid). Responsibility: the
   projection-independent assignment source; the projection derives
   its drop set from it, and the S3 reflective oracle asserts (a)
   every struct field has exactly one entry, (b) entry set == field
   set both directions, (c) a projected run's emitted keys equal the
   table's `plan` side.

3. **The total command-tree walker** — GUESS at name/home:
   `walkAllCommands(root *cobra.Command) []*cobra.Command`, a test
   helper for the S4 oracle (Phase 2; GUESS: lives in the S4 test
   file, not shipped code). Responsibility: descend EVERY child of the
   root — no name-based skip, no `Hidden` gate (explicitly NOT
   `help_all.go::walkCommandTree`, which skips `help`/`completion`) —
   so the oracle can assert the set of commands registering
   `plan-only` (per-command `Lookup("plan-only")`) is exactly
   `{flow resolve}`. Preconditions it must establish and assert:
   materialize both auto-generated commands via cobra's OWN
   initializers (`InitDefaultHelpCmd` is already force-inited;
   `InitDefaultCompletionCmd` must be called since a bare
   `NewRootCmd()` tree lacks `completion`), and assert `completion`
   (the discriminating member) and `help` are IN the walked set;
   vacuity-guarded on the four flow verbs existing. Implementability
   is A5 (Pending) — written for the first time at Phase 2.

## 3. Data model (persisted / crossed boundary)

Nothing new is persisted. The boundary artifact is the resolve success
wire record.

- **`resolvePayload`** (internal/cli): keeps EXACTLY its current 14
  struct fields — `decision_table_0010_test.go` pins
  `NumField() == 14` and the 13-key wire list; A2's conversion changes
  field TYPES only, adds none. Fields by wire key:
  - PLAN: `revision` (loader-produced identity token, presently
    constant-empty — `flow_exec.go::revision` returns `""`; kept as
    the plan-side identity slot), `rule`, `gates`, `emit` (never
    omitted, `{}` when empty — 0010:C4), `next`, `writes`, `clear`,
    `escaped` (unconditionally present), `escape_class` (`omitempty`
    per its producer — the 14th field / conditional 14th key).
  - ECHO (become pointer-valued + `omitempty` under A2): `model`,
    `observed` (the caller's `--tag` set, echoed verbatim except
    set-literal canonical re-encode), `owned`, `readers`, `outcome`.
    GUESS at exact Go types: `model *string` (GUESS — could be a
    small struct pointer), `observed *map[string]string`,
    `owned *map[string]string`, `readers *[]string`,
    `outcome *string`. Default mode populates all five non-nil, so
    empty containers render `{}`/`[]` and never `null` (A9's
    invariant: `tagMap`, `observedTagMap`, `emitMap`, `readerIDs` all
    allocate unconditionally; Phase 1 adds a default-mode empty-shape
    oracle before the conversion lands).
- **Envelope**: NDJSON success `{"type":"ok","data":{…payload…}}`;
  refusal `CLIError{code, message, param, detail, hint, findings}` —
  never projected. Encoder: HTML escaping disabled, struct
  declaration order on the wire.
- **Partition table**: field name → `echo`|`plan`, in-code only (not
  persisted), the normative assignment source for projection and S3.
- **Always-keep core** (C2, forward-binding for any future mode):
  `rule`, `escaped`, `escape_class`, `revision` — no projection mode
  may drop them (projection-invariance, not unconditional presence).
  With one mode, S2's key-set literal is the enforcement; a successor
  adding a second mode owes a per-mode always-keep oracle.

## 4. Top-level pseudo-code of the main operation

```
runFlowResolve(cmd, args):
    # flag registration happened at command construction:
    #   resolveCmd.Flags().Bool("plan-only", false, …)   # resolve only,
    #   long form only, no shorthand, no shared registrar
    model, err   := loadModel(--model)                  # refusal paths…
    if err: return respond.Fail(clierr(...))            # …all flag-blind
    tags, err    := parseTags(--tag…)                   # refuses owned/
    if err: return respond.Fail(...)                    # reserved/malformed
    outcome      := --outcome                           # copied verbatim
    readers      := run reader pass (model, outcome)    # actually EXECUTES;
                                                        # S1 observes this
                                                        # at the invocation
                                                        # site, ± flag equal
    owned        := assemble owned view (prior state overlay)
    result, err  := resolve.Resolve(model, Input{Owned: owned,
                                                 Observed: tags}, outcome)
    if err (unrecognized outcome, no rule, …):
        return respond.Fail(clierr(...))                # last Fail return —
                                                        # projection is
                                                        # unreachable above
    payload := resolvePayload{                          # all 14 fields,
        Revision: revision(model),                      # exactly as today
        Rule: result.RuleID, Gates: gateResults,
        Emit: emitMap(result), Next: …, Writes: …, Clear: …,
        Escaped: …, EscapeClass: escapeClassOf(result), # omitempty
        Model: &modelRef, Observed: &observedTagMap(tags),
        Owned: &tagMap(owned), Readers: &readerIDs(readers),
        Outcome: &outcome,                              # echo ptrs non-nil
    }
    if planOnly := cmd.Flags().GetBool("plan-only"); planOnly:
        # the ONLY read of the flag; success path only, before gateway
        payload = projectPlanOnly(payload)              # nil every field
                                                        # the declaration
                                                        # table marks echo
    return respond.OK(cmd, payload)                     # one renderer for
                                                        # --as=json and
                                                        # --as=text, so both
                                                        # modes agree by
                                                        # construction
```

(GUESS: exact assembly order of readers/owned/kernel-call above — the
RDR fixes only that all of it is upstream of the projection site and
flag-blind; helper names `tagMap`/`observedTagMap`/`emitMap`/
`readerIDs`/`escapeClassOf`/`invokedReaders` are cited by the record,
the rest of the skeleton is a GUESS from those citations.)
