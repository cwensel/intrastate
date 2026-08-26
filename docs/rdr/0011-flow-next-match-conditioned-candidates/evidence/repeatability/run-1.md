model: claude-sonnet-5
variant: lite (profile: large)

# Reconstruction from RDR 0011: flow next selects by match; --all enumerates the alphabet

## 1. Public API

**CLI surface — `flow next` (modified verb)**

```
intrastate flow next --model <path> [--artifact <role>=<ref>]... [--tag <key>=<value>]...
                      [--all] [--evaluate-gates] [--as=text|json]
```

- New flag: `--all bool` (default `false`), registered on `next` only — NOT on `flow resolve`, `flow read-state`, `flow set-state`, or the `flow` group itself. Non-registration is the enforcement mechanism; an unregistered `--all` on those verbs fails pflag parsing before any `RunE` runs, mapped by the CLI's own error gateway (not cobra's printer, which is silenced) to `{"code":"command-error"}`, `Group: GroupUserEnv`, exit 2. This is the shared "unknown flag" bucket, not a dedicated refusal — no new error code is minted for it.
- Behavior change on the existing default (no `--all`): candidate predicate narrows from "every non-escape row whose guard the kernel does not decide false" (0005's alphabet) to "every non-escape row whose match atoms over PRESENT keys all hold, AND whose guard the kernel does not decide false." Absence of a match key never excludes a row; only a present-and-unequal value does (via the kernel's `no_match`).
- `--all` restores the exact 0005 predicate (guard-only; match plays no part; match keys are never entered into `unknown` under `--all`, regardless of presence).
- `outcomes` remains the model's full declared alphabet in every mode, unconditionally.
- Response payload change: the candidate's field renamed `unresolved []string` → `unknown []{key, reason}`. `unknown` is always present, `[]` rather than omitted when empty, in both modes.
  - `reason` is a closed vocabulary of three tokens: `"absent"` (match/guard/owned key not in view — CLI-derived), `"uncomparable"` (kernel's guard seam could not decide a present key's value — reused from `0007:C8`'s two-member `Reason` enum), and `"not-evaluated"` (a declared gate id not run because `--evaluate-gates` was not passed — minted by this RDR, scoped to gate ids only).
  - Entries deduplicated by the `{key, reason}` pair (not by key alone) and sorted by `(key, reason)`.
- Help text (Short/Long) MUST state: default reports rows the supplied state can take (match+guard both holding or undecided); a match/guard key the state lacks leaves the row a candidate with `unknown` naming it; `--all` ignores match; a "candidate" is a row the state does not exclude, not a row `flow resolve` will select.

**Error modes surfaced (no new `RefusalKind`; `0007:REQ-79` still forbids a sixth):**
- `flow-tag-duplicate` — repeated `--tag <key>` (pre-existing, `internal/cli/flow_input.go::parseTags`).
- `flow-tag-owned` — `--tag` on a key the model declares owned (pre-existing).
- `malformed_rule_shape` / `owned tag %s is served by %d readers; want exactly one` (`checkAccessorBindings`) — model load-time refusals that make CLI presence == kernel comparability exact (A3).
- `flow-no-match`, `flow-artifact-missing` (exit 2), reader's own refusal e.g. `flow-accessor-failed` (exit 3) — `flow resolve`'s three possible dispositions once the demand-set extension (C1) invokes a reader serving a match-only owned key that used to go uninvoked (was uniformly `flow-no-match` before).
- `command-error` / exit 2 — `--all` rejected on the other three verbs (usage-bucket, not a dedicated code).
- Kernel refusal kinds unchanged: `no_match`, `guard_unevaluable`, `owned_state_unavailable`, `unmodeled_outcome` (unreachable on a one-row probe per A1), `ambiguous_match` (vacuous on a one-row probe).

**Internal Go-level API touched (from Illustrative Code / mini-checks / Placement decision):**

```go
// internal/cli/flow_next.go — probe construction, changed
func excluded(row table.Row, owned []accessor.Tag, observed []resolve.Tag) *resolve.Result
// GUESS (RDR states behavior, not the exact signature): today returns a bool;
// C1/A13 require it return the kernel's Result (or equivalent) so summarize can
// read Refusal.Undecided off it instead of collapsing to bool.

// internal/cli/flow_next.go — reporting, changed field type
type unknownEntry struct { Key string; Reason string } // GUESS at field names;
// RDR fixes only the shape {key, reason} and the JSON view, not Go field casing.
func summarize(...) /* candidate.Unknown []unknownEntry, replacing []string Unresolved */

// internal/cli/flow_exec.go — demand-set builder, one added term
func invokedReaders(model *table.Model, outcome string) []Reader
// today: RequiresOwned ∪ guardOwnedKeys(m, row)
// after C1: RequiresOwned ∪ guardOwnedKeys(m, row) ∪ matchOwnedKeys(m, row)  -- GUESS at helper name
// bound at BOTH call sites: flow_next.go:97 (outcome="") and flow_resolve.go:98 (outcome=req.outcome)
```

## 2. Three most important internal helper functions

1. **`excluded` (probe builder), `internal/cli/flow_next.go`** — Per candidate row, builds a one-row kernel probe: `probe := row.KernelRow()`, strips the escape list, binds `Recognized: row.Outcome`, then filters `probe.Match` to keep only `resolve.Tag`s whose `Key` is present in the CLI's assembled view (owned ∪ observed). Filtering happens AFTER `KernelRow()` conversion — never by re-deriving tags from `row.Atoms` — because `KernelRow()` already canonicalizes set-kinded literals (sorted/compacted JSON via `seamValue`) and re-deriving would risk mis-encoding. Responsible for making "presence" a per-key, all-or-nothing decision (a key's several match tags travel or are dropped together) and for returning the kernel's own verdict (plan / `no_match` / `guard_unevaluable` / `owned_state_unavailable`) — never comparing a value itself.

2. **`summarize` (reporting), `internal/cli/flow_next.go`** — Owns the CLI's one presence test (`if _, known := view[atom.Key]; known { continue }` over `assembledView`) and now must (a) walk `row.Atoms` to emit `{key, absent}` for every match/guard atom over an absent key, (b) read `Refusal.Undecided` off a `guard_unevaluable` probe result to emit `{key, uncomparable}` for kernel-side undecided guard atoms (deduping against the walk's `absent` entries on the `{key, reason}` pair), (c) emit `{gate-id, not-evaluated}` for each declared gate id when `--evaluate-gates` was not passed, and (d) under `--all`, suppress ONLY the match-atom-derived entries from this walk while leaving owned-key-absent entries (from the demand-set walk) intact. Produces the final sorted, deduped `unknown` list per candidate.

3. **`invokedReaders` (demand-set builder), `internal/cli/flow_exec.go`** — Computes, once per invocation before the row loop and before any kernel call, the set of readers to invoke: the union of each candidate row's `RequiresOwned` (write/clear-derived), guard-block owned keys, and — the new term this RDR adds — match-block owned keys. Mode-independent (identical under `--all`, since the view itself doesn't change, only whether match participates in the verdict) and shared by both `flow next` and `flow resolve` callers, so the two verbs never assemble different views over the same model. This is the function whose omission of match-owned keys was the defect: without the added term, a match-only owned key never gets read, is never in the view, and every row silently becomes a candidate carrying `{key, absent}` — the default silently degrading to `--all` behavior for that model.

## 3. Data model (persisted / passed across the boundary)

**Candidate payload (JSON response shape, per `flow next` output):**

```
Candidate {
  rule_id:    string
  outcome:    string
  unknown:    []{ key: string, reason: "absent" | "uncomparable" | "not-evaluated" }
              // sorted by (key, reason); deduped by (key, reason); ALWAYS present as [] if empty
  // GUESS: other pre-existing fields (e.g. gate ids, readers list) carried forward unchanged
}

Response {
  outcomes:   []string      // model's full declared alphabet, unconditional, every mode
  candidates: []Candidate
  readers:    []string      // GUESS: field name; S8 confirms this exists and is mode-independent
}
```

**One-row kernel probe (internal, passed to `internal/resolve`):**

```
Probe (== resolve.Table with exactly one row) {
  Row: {
    Match:  []resolve.Tag        // FILTERED: only tags whose Key ∈ assembled view keys
    Guard:  []GuardAtom           // unchanged from row.KernelRow()
    RequiresOwned: []string       // unchanged
    Writes: ...                   // unchanged, escape list stripped (no Escape rows)
  }
  Recognized: row.Outcome         // bound to the row's own outcome, keeps unmodeled_outcome unreachable
}
```

**Assembled view (CLI-side, per invocation):** `map[string]Tag` built from `owned` (from `invokedReaders` → `runReaders` → `OwnedSnapshot`, absent-valued keys omitted) unioned with `observed` (`--tag` pairs, last-write-wins). Same key SET the kernel's `assemble` builds (owned ∪ observed ∪ synthetic `recognized`), per A11.

**Demand set (CLI-side, per invocation):** `set[string]` = `RequiresOwned(rows) ∪ guardOwnedKeys(model, rows) ∪ matchOwnedKeys(model, rows)` — computed once, shared verbatim between `flow next` and `flow resolve`.

**Kernel-side types referenced but NOT modified:** `resolve.Tag{Key, Value}`, `resolve.TagSet`, `RefusalKind` (5-member closed set, unchanged), `Refusal.Undecided []UndecidedRow{RuleID, SourceLocator, Atoms []UndecidedAtom{Key, Block, Operator, Literal, Reason}}` (`0007:C8`'s shape, read but not widened), `Reason` (2-member closed enum: `absent`, `uncomparable`).

GUESS: exact Go struct/field names throughout this section, since the RDR speaks in terms of behavior/shape contracts (JSON keys, dedup/sort rules, closed vocabularies) and cites existing functions by name but does not give full struct definitions for the new `unknown` entry type or the probe-builder's intermediate types.

## 4. Top-level pseudo-code of the main operation (`flow next`, ~35 lines)

```
function FlowNext(model, artifactRefs, observedTags, all bool, evaluateGates bool) Response:
    rows := model.NonEscapeRows()

    # Demand set: union over ALL rows, computed once, mode-independent (C1)
    demand := RequiresOwned(rows) ∪ GuardOwnedKeys(model, rows) ∪ MatchOwnedKeys(model, rows)
    readers := ResolveReaders(model, demand)
    owned, err := RunReaders(readers, artifactRefs)          # aborts verb on any reader refusal
    if err != nil: return AccessorFailure(err)                # exit 3

    view := AssembledView(owned, observedTags)                 # owned ∪ observed, last-write-wins

    candidates := []
    for row in rows:
        probe := row.KernelRow()
        probe.Recognized = row.Outcome
        probe.Escape = nil

        if not all:
            probe.Match = [t for t in probe.Match if t.Key in view]   # presence filter, per-key

        else:
            probe.Match = nil                                  # --all: strip match entirely

        result := Kernel.Resolve(OneRowTable(probe))            # ask the kernel; equality is its call

        switch result.Disposition:
            case Plan, GuardUnevaluable, OwnedStateUnavailable:
                unknown := []
                for atom in row.Atoms:                          # CLI's own presence walk
                    if atom.Key not in view:
                        unknown.append({atom.Key, "absent"})

                if all:
                    unknown := [e for e in unknown if e.key not a MATCH-block atom's key]

                if result.Disposition == GuardUnevaluable:
                    for undecidedAtom in result.Refusal.Undecided:
                        unknown.append_dedup({undecidedAtom.Key, undecidedAtom.Reason})  # absent|uncomparable

                if not evaluateGates:
                    for gateID in row.DeclaredGates:
                        unknown.append({gateID, "not-evaluated"})

                unknown := Dedup(unknown, by={key,reason})
                unknown := SortBy(unknown, key, reason)

                candidates.append(Candidate{row.ID, row.Outcome, unknown})

            case NoMatch:
                continue                                        # excluded; gates never evaluated

            case UnmodeledOutcome, AmbiguousMatch:
                unreachable                                     # A1: cannot arise on a one-row probe

    return Response{Outcomes: model.FullAlphabet(), Candidates: candidates, Readers: readers}
```

## Spans widened past

- Widened past `0011:C1` into the full `§technical-design` → `§approach` and `§load-bearing-decisions` subsections, to recover the exact Go call sites (`invokedReaders`, `excluded`, `summarize`), the demand-set formula, and the placement rationale (CLI-decides-presence / kernel-decides-equality split) that C1's prose references but states more tersely than the surrounding design narrative.
- Widened past `0011:C2`/`0011:C3` into `§illustrative-code` for concrete CLI invocation shape (flags, `--as=json` usage) since the contracts describe behavior but not command-line examples.
- Widened past the `§testing-strategy` (S1–S8) spans into `§mini-checks` (the `disposition` and `authority` tables) to recover a compact enumeration of the input-class → output-shape mapping used for the pseudo-code's `switch` and for the data-model's `reason` vocabulary — the contracts state this discursively across several paragraphs; the tables state it densely and were used to cross-check completeness.
- No Go struct field names, JSON field casing for the new `unknown` entry type, or the exact probe-builder helper function name (`matchOwnedKeys` is invented) are given anywhere in the record; these are marked GUESS in sections 1 and 3 above. The record fixes shapes and vocabularies (closed reason enum, dedup/sort rule, JSON presence-as-`[]`) but deliberately stays at the behavioral/contract level for identifiers not already shipped in 0005's code (e.g., `Row.Atoms`, `KernelRow`, `assembledView`, `invokedReaders`, `excluded`, `summarize` are all pre-existing/cited by name; the new field and any new small helper are not).
