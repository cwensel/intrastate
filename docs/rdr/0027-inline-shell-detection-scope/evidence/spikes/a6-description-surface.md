Model: claude-sonnet-5

# Spike: A6 — does a reviewer-reachable, non-refusal landing site exist for a category description?

## Question

A6: "A category description can be surfaced to a reviewer without provoking
a refusal — there is somewhere for Phase 3's promise text to land and be
read." Identify the surface, or determine that none exists and none is
permitted by the output contract.

## Method

Files read in full, commands run as shown. No record file was sed/grepped;
only the RDR projector was used to read 0027 itself.

### 1. `internal/table/category.go` (full file, 204 lines)

- `Category` is `type Category string` (line 11) — a bare string, no
  attached description field.
- `Categories() []Category` (lines 66-107) returns the closed identifier
  set in declaration order, appending the RDR 0025 six at the tail
  (lines 100-105) per `0025:C5`.
- `Failure` (lines 111-150) is the refusal payload: `Category`, `Detail`,
  `Offending`, `Remedy`, `Rule`, `Line`. `Detail` is per-instance refusal
  text, not a per-category static description.
- No consumer of `Categories()` outside `_test.go` files:
  ```
  grep -rn 'table\.Categories\(\)' --include='*.go' internal cmd
  ```
  every hit is in `internal/table/*_test.go` (roundtrip, dump,
  emit_witness, reserved_key_fixtures, command_carrier tests). Confirms
  the record's "absence half": no shipped, non-test code reads the
  category set today.

### 2. CLI verb enumeration

`internal/cli/root.go:190-198` registers exactly three top-level commands:
`version`, `lint` (root, per `0006:C19`), `flow` (group, per `0005:C1`),
plus a hidden `docs` command (`root.go:198`).

Verbs that emit reviewer-facing text with **no defect present**:

- `intrastate lint --help-all` / `intrastate --help-all` — cobra help,
  not routed through the refusal gateway, exits 0, always available.
- `intrastate docs` — hidden, generates `docs/cli-reference.md` and
  `llms.txt` from the live command tree (`internal/cli/docs.go`).
- `intrastate help --all` — same extended-help mechanism, tree-wide.

`internal/cli/lint.go:67-113`, `lintExtendedDesc()`, is the load-bearing
precedent. It builds the `--help-all` body for `lint` by iterating
`graphlint.BlockingCodes()` and `graphlint.AdvisoryCodes()`
(`lint.go:76-78`, `lint.go:85-87`) — i.e., it renders a **closed taxonomy's
own descriptive text live from an exported accessor**, with zero defect
present, registered via `withExtendedHelp(cmd, lintExtendedDesc())` at
`lint.go:58`.

`internal/graphlint/taxonomy.go` is the taxonomy this reads: closed const
code sets (`taxonomy.go:24-44`), an exported accessor per tier
(`BlockingCodes()`/`AdvisoryCodes()`, `taxonomy.go:105-109`), no
per-code description field on the taxonomy itself — the description text
lives beside the accessor call site in `lint.go`, not in `taxonomy.go`.
This is the same shape `table.Category`/`Categories()` would need for a
description surface: not "add a `Description` field to `Category`," but
"add an accessor + a describing loop consumed by an extended-help body,"
mirroring `graphlint`'s pattern exactly.

`internal/cli/help_all.go` (full file read) — the `--help-all` mechanism
is orthogonal to the refusal path: `withExtendedHelp` stores text on
`cmd.Annotations` (help_all.go:70-76); the PreRunE hook prints it and
returns `pflag.ErrHelp`, "deliberately NOT routed through the respond
gateway" (help_all.go:37-41) — i.e., it is architecturally guaranteed to
be reachable without a refusal.

`internal/cli/docs.go` (full file read) — `runDocs` (docs.go:94-140)
walks the same command tree and renders each command's `Long` plus its
extended-help body (`writeCLIReference`, docs.go:170-210) into
`docs/cli-reference.md`, a **committed, generated, non-refusal** file
(staleness-gated by `make check` per docs.go:84-88). This is candidate
(b) named in the spike brief: a documented lint-help path that also
lands in a generated doc file.

### 3. `docs/cli-output-contract.md` (361 lines, read in full)

- Line 5-15: "The binary states the contract itself... `intrastate lint
  --help-all` [carries] the live finding taxonomy." The contract
  explicitly treats `lint --help-all` as the authoritative place a
  taxonomy's description lives — direct precedent for putting
  `command_shell_interpreter`'s promise text in the same place.
- The `findings` table (lines 45-58) documents `code` as "the finding's
  own discriminator — for a model-load failure, the load category slug"
  — confirms `Category` values already ride the wire as `code`, so a
  description keyed off that same closed set is a natural extension,
  not a new envelope shape.
- No clause anywhere in the file forbids adding descriptive help text or
  a generated-docs entry for a category. The contract governs wire
  **payload** shape (JSON envelope, `findings[]`, set-value encoding,
  exit codes) — it says nothing that constrains `--help-all`/`docs`
  content, which is cobra help output explicitly carved out as a
  separate, non-enveloped stream (`help_all.go:37-41`).

### 4. Build sanity

```
go build ./...
```
Repo builds clean (two pre-existing unrelated `main redeclared` errors
surfaced only when compiling `docs/rdr/**/evidence/spikes/*.go` scratch
files under `./...`, not part of `internal/...`; irrelevant here).

## Candidate surfaces and verdicts

| candidate | reachable w/o refusal? | permitted by output contract? | verdict |
| --- | --- | --- | --- |
| (a) description field on `Category`/`Categories()`, read by an existing verb | not yet built; no field exists on `Category` today | yes — contract doesn't govern help content, only wire payload | VIABLE, not yet present |
| (b) documented lint-help path (`lint --help-all`, extended-help body) | YES — architecturally exempt from the refusal gateway (`help_all.go:37-41`) | YES — contract's own line 5-15 names `lint --help-all` as the taxonomy's live home | VERIFIED, precedent exists today (`graphlint` taxonomy) |
| (c) generated docs file (`intrastate docs` → `docs/cli-reference.md`) | YES — same mechanism, committed and staleness-gated | YES — explicitly documented as the file form of `--help-all` (`docs.go:166-170`) | VERIFIED, precedent exists today |

## Conclusion

A6 is **VERIFIED**. A landing site is reachable today, by direct,
already-shipped precedent — not by speculative extension:

- **Site**: `internal/cli/lint.go:67-113` (`lintExtendedDesc`), registered
  at `internal/cli/lint.go:58` via `withExtendedHelp`, rendered by
  `intrastate lint --help-all` (and mirrored into
  `docs/cli-reference.md` by `internal/cli/docs.go:170-210`,
  `runDocs` at `docs.go:94`).
- **Shape**: identical to how `graphlint`'s closed taxonomy already
  surfaces there (`internal/graphlint/taxonomy.go:105-109`,
  `BlockingCodes()`/`AdvisoryCodes()` consumed at `lint.go:76-87`). The
  Phase 3 implementer's job is the same pattern one level down: an
  exported accessor pairing each `table.Category` (or just
  `CatCommandShellInterpreter`) with its description text, consumed by
  `lintExtendedDesc()` or a sibling extended-help body, which already
  renders on `--help-all` — a stream `help_all.go:37-41` guarantees runs
  outside the refusal gateway — with zero defect present.
- `docs/cli-output-contract.md` neither forbids this nor requires a new
  wire shape: it explicitly defers taxonomy description to
  `--help-all` (contract lines 5-15) and constrains only the refusal
  envelope's payload fields, which this surface does not touch.

This is "not built yet but a documented, permitted, local addition"
per the spike brief's distinction — not "no permitted surface exists."
C1's `promise:` clause condition on A6 is satisfied; S7 is writable.
