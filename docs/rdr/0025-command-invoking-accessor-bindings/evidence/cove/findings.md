Model: claude-opus-5[1m] (consolidation) — dual-model pass over claude-opus-5[1m] (findings-a.md) + claude-sonnet-5 (findings-b.md)

# CoVe consolidated findings — RDR 0025

Dual-model cove. Two independent passes ran with no cross-visibility; this file
reconciles them **by passage anchor** (per-pass F-numbers do not correspond).
Every row below was re-verified in the parent against source before landing here
— a pass's verdict alone did not earn a row.

## Convergence

Both passes independently landed on **C6's config subsystem does not exist**.
Independent convergence on one passage is a hotspot, not duplication.

## Origin ledger

| id | Anchor | Class | Finding | Verified-in-parent evidence | Disposition |
|---|---|---|---|---|---|
| CV-1 | 0025:C6 | (a) REFUTED | C6 sites the opt-in in "`internal/cli/config/` discovery" + user-scope `intrastate.toml`. Neither exists: no such dir, zero `intrastate.toml` refs in `*.go`, no user-scope TOML loader, no `allow_commands`/`AllowCommands` anywhere. C6 **specifies a new subsystem while reading as a citation of an existing one**. | `ls internal/cli/` (no `config`); `grep -rn 'intrastate\.toml' --include='*.go'` → none; `grep -rn 'allow_commands\|AllowCommands'` → none | open |
| CV-2 | 0025:C4 (+F1/F2/F3 refs) | (a) REFUTED | C4 promises "`Detail` carries the last 4 KiB of stderr" and C6 a "Detail naming `allow_commands`". `accessor.Refusal` has **no `Detail` field** (`model.go:280-308`). `Detail` exists on `accessor.Finding` (`model.go:138-139`) and `table.Failure` (`category.go`) — different structs on different paths. The command diagnostics C4/C6/F1-F3 promise have **no carrier on the refusal type**. | `model.go:280-308` (Refusal: Class, Accessor, Capability, Role, Timeout, Keys, Expected, Observed, Reason, applied — no Detail) | open |
| CV-3 | 0025:C5 | (d) silence | C5 mints five new `table.Category` constants but is silent on the mandatory `Categories()` registration. `Categories()` (`category.go:52-83`) is a **hand-maintained closed list in declaration order** — a constant absent from it is not in the closed set. RDR 0024 fixed this as "ONE decision". S1 asserts per-mutant category and would pass without registration. | `category.go:52-83` — 28 hand-listed constants | open |
| CV-4 | 0025:A9 | (a) NOT-FOUND | Cites `internal/cli/flowbind/flowbind.go::Reader.run`. No such method; it is `Reader.Read` (`flowbind.go:189`). A9's **substance** (the exhaustive consumer list) re-verified correct by both passes — one-token anchor rename. | `flowbind.go:189` `func (r Reader) Read(...)`; no `Reader) run` in package | open |
| CV-5 | 0025:A4 | (a) NOT-FOUND | Cites `internal/accessor/executor.go::invokeWrite`. Package has `invokeRead` (`:156`) but **no `invokeWrite`**; the real write site is `Executor.Write` (`:235`), which reaches `readerFor` at `:275`. Behavior claim holds; anchor is fictional. | `grep 'func.*invokeWrite'` → none; `executor.go:81,156,235` | open |
| CV-6 | 0025:C3, 0025:D-wire-byte-format | (b) sibling | C3's flat-JSON-string-map read envelope with presence-vs-absence semantics duplicates the decode already normative in `flowbind.go::store` (`:110-136`) for REQ-66/REQ-107. Two decoders of one wire shape, no cross-cite. | `flowbind.go:110-136` | open |
| CV-7 | 0025:C3 `exit_verdicts`, 0025:D-selection-predicate | (b) sibling | The four-way verdict split duplicates `flowbind.go::verdictFor` + `Gate.Gate` (`:78-105`), which D-selection-predicate explicitly **retains unchanged**. Two sites decide allow/deny/indeterminate. | `flowbind.go:78-105` | open |
| CV-8 | 0025:§approach, 0025:§consequences | (c) contradiction | Approach frames the command family as ending the accretion at this seam ("removes the pressure to grow it"), while D-selection-predicate retains the file binding and its magic-suffix vocabulary (`verdictFor`/`unreachable`) unchanged. The accretion is bounded, not ended. | RDR-internal; `flowbind.go::verdictFor`/`::unreachable` retained | open |
| CV-9 | 0025:§prerequisites, 0025:S5 | (d) silence | The one unchecked prerequisite is 0016's fail-closed reader uniqueness. Under today's inherited first-match `readerFor`, "first" is **name-sorted** (`registry.go:26,36,46` build sorted). No scenario pins two-reader read-back behavior, which is exactly the A4 `Deferred` risk surface. | `registry.go` `slices.Sorted(keys(...))` ×3 | open |
| CV-10 | 0025:C5 `command_shell_interpreter` | (d) silence | The interpreter enumeration (`sh -c`, `bash -c`, `python -c`, "env chains") is never stated as **closed or heuristic**, unlike C3's exit maps and 0004's eight-member ValidationCode closure, both of which state their closure explicitly. An implementer cannot tell whether to reject an unlisted interpreter. | RDR-internal; contrast C3 / `accessor.ValidationCodes()` closure | open |
| CV-11 | 0025:A2, 0025:§failure-modes | (d) silence | A2 acknowledges a Windows argv-requoting divergence; no Failure Mode (F1-F7) covers a command entry on a Windows-targeted model. | RDR-internal | open |
| CV-12 | 0025:A1, 0025:§capability-dependencies | (a) nuance | A1's prose can read as though the process-group/deadline triple is already in production `internal/accessor`. It exists in the **A1 spike only**; `executor.go` has no `Setpgid`/`Cancel`/`WaitDelay`. Capability Dependencies correctly books it "Available (Go stdlib)", so this is a prose-precision fix, not a refuted dependency. | `executor.go` — no Setpgid/Cancel/WaitDelay | open |

## Dismissed in consolidation (did not survive the parent's grounding gate)

- **Finalization Gate sub-sections are unfilled template text** (raised by BOTH
  passes; pass-B made it its blocking reason). **Factually true, not a cove
  finding, not blocking at Stage 5.** The gate is authored at Stage 7 by
  `/rdr-finalize`; `gate_written=false` is the correct state for a Draft at
  Stage 5, and the template's own instruction block reads "Complete each item
  … before marking this RDR as **Final**." Dismissed-with-cite: fails ground 3
  (the RDR's own decided text / the flow's stage contract). Cove finds design
  silences, not slots owed by a later stage.
  - One sub-part survives inside CV-*: §risks defers items *to* the
    cross-cutting gate, which is a real Stage-7 obligation — carried forward as
    a note, not a Stage-5 edit.

## Grounding note — why the 32/32 anchor check missed CV-4 and CV-5

All 32 `source-anchor` edges resolved `true` mechanically, yet two cite methods
that do not exist. The resolver evidently grounds at **file/package** granularity,
not **method** granularity. Cove's Step 0 semantic half is what caught them. This
is a durable lesson about the anchor check's resolution, worth carrying beyond
this record.

## Dispositions (resolve half, iteration 1)

| id | Disposition | Edit |
|---|---|---|
| CV-1 | **fixed** | C6 prose + normative block now state the config surface is NEW and this RDR's to build; flag-alone named a complete v1 gate; Capability Dependencies row added (`Build`); **A10 booked Pending** for discovery/precedence/malformed-config |
| CV-2 | **fixed** | C4 normative block names `Detail string` as a NEW field on `accessor.Refusal`; prose corrects "executor unchanged" → classes unchanged + one additive type change, cites why `Reason` is not reused (`0004:C7`) and why the closure argument does not bar it; Capability Dependencies row added (`Extend`) |
| CV-3 | **fixed** | C5 normative block gains a `registration:` line; prose explains an unregistered constant refuses correctly while staying invisible to enumerators; S1 now asserts `Categories()` membership |
| CV-4 | **fixed** | A9 `Reader.run` → `Reader.Read` (substance was correct) |
| CV-5 | **fixed** | A4 `executor.go::invokeWrite` → `::Executor.Write` (behavior claim held) |
| CV-6 | **fixed (reframed)** | Not duplicated logic — `store` decodes intrastate's own artifact file, C3 decodes a child's stdout: two transports, one shared map-of-strings + presence-is-the-answer rule. Kinship now cited in the Wire LBD so the shared rule cannot drift uncited |
| CV-7 | **dismissed-with-cite** | `verdictFor` derives a verdict from a **magic path suffix** — the simulation this RDR displaces, not a rival decoder of the C3 protocol. Disjoint by carrier (`flowbind.go:82-101`); recorded in the authority census rather than "fixed" |
| CV-8 | **fixed** | Consequences "ending the accretion" → "removing the pressure"; Selection LBD now states the suffixes are retained and their retirement is a successor's call |
| CV-9 | **fixed** | S5b added: two-reader read-back pinned to name-sorted first match, asserted **by selected-reader identity** (either reader verifies a correct write, so a success oracle would pass under both) |
| CV-10 | **fixed** | C5 normative block gains an `interpreter set: OPEN` line; prose states the deny-list is explicitly not closed, names what it does and does not guarantee, and makes additions an amendment |
| CV-11 | **fixed** | F7 added: a command entry on Windows refuses `execution_failure` before spawn; lint stays platform-neutral |
| CV-12 | **fixed** | A1 header: "the executor's **existing** context deadline plus … **this RDR adds** in the binding" |

Gate-sub-sections finding: **dismissed-with-cite** (Stage 7's obligation; see above).

## Mini-checks fired (first lens pass owes the cue read)

`authority` · `fidelity` · `disposition` · `trace` — four tables written into the
RDR under *Pre-Lock Mini-Checks*. `test-discriminability` did **not** fire.
The desk trace walked the MVV end to end: **no CONTRADICTION row**.

## Needs (re)verification — carried to Stage 6

- **A10 (new, Pending)** — user-scope config discovery, precedence vs the flag,
  malformed-config behaviour (must not fail open).
- **A1, A4 (still Verified)** — unchanged in substance; edits were precision only.
- No Verified assumption was invalidated by these fixes.

## Convergence

Iteration 1 converged: no open ledger entries. All 12 traced to origin findings;
no net-new scope absorbed (nothing charted — CV-1's config surface is C6's own
dependency, not a successor's).
