model: claude-sonnet-5
variant: lite (profile: large)

# Reconstruction run — RDR 0029 (version promise on machine-readable output)

Grounded from: 0029:C1-C4, 0029:MVV, 0029:S1-S9 (all read verbatim via
projector), plus widened reads of 0029:§illustrative-code and
0029:§load-bearing-decisions (the contract spans left the exact wire JSON
shape and the naming rationale under-determined — widened past C1/C4 to
avoid guessing a field name or JSON shape the RDR actually pins in prose).
Also widened into current source (`internal/cli/respond/respond.go`,
`internal/cli/clierr/clierr.go`, `internal/version/version.go`, and the
nine existing enumeration seams: `guard.Operators`, `table.Operators`,
`accessor.Verdicts`, `resolve.RefusalKinds`, `graphlint.{Severities,
Reasons,BlockingCodes,AdvisoryCodes}`, `table.Categories`) to reproduce
real Go signatures rather than invent them — C4's census table names these
seams but does not give their signatures, which are load-bearing for a
`respond`/`clierr` reconstruction.

No section was reconstructed from a truncated selection; every widen is
listed above.

---

## 1. Public API

```go
// package respond (internal/cli/respond)

// Success is the terminal-success envelope under --as=json.
// SchemaVersion is new (C1): present on Success, NOT omitempty, begins
// at "0.1" (C1, Load-Bearing Decisions: string "MAJOR.MINOR", never a
// pair of ints or a float).
type Success struct {
	Type          string     `json:"type"`
	SchemaVersion string     `json:"schema_version"`
	Notes         []Advisory `json:"notes,omitempty"`
	Warnings      []Warning  `json:"warnings,omitempty"`
	Data          any        `json:"data,omitempty"`
}

// SchemaVersion is the single exported constant both respond and clierr
// read (C1: "one home"). Lives in clierr — the leaf package — per C1's
// layering argument (respond imports clierr; clierr imports neither).
// GUESS: exact identifier name; C1 fixes the value "0.1" and the
// single-constant requirement but does not spell the Go identifier.
const SchemaVersion = "0.1" // package clierr

func OK(cmd *cobra.Command, s Success) error
func Fail(cmd *cobra.Command, err *clierr.CLIError) error
func ModeOf(cmd *cobra.Command) Mode
func ValidateMode(cmd *cobra.Command) *clierr.CLIError

// package clierr (internal/cli/clierr)

// CLIError is the bare refusal record — no wrapper, no "type" key
// (0005:C1). SchemaVersion is new (C1: present on both terminal
// records, non-omitempty, top level only — never under Findings).
type CLIError struct {
	SchemaVersion string     `json:"schema_version"`
	Code          string     `json:"code"`
	Message       string     `json:"message"`
	Param         string     `json:"param,omitempty"`
	Detail        string     `json:"detail,omitempty"`
	Hint          string     `json:"hint,omitempty"`
	Findings      []Finding  `json:"findings,omitempty"`
	Group         ErrorGroup `json:"-"`
	Cause         error      `json:"-"`
}

func (e *CLIError) Error() string
func ExitCodeFor(err error) int // {0,1,2,3,130} — owed seam: ExitCodes()

// Enumeration seams C4 obliges (frozen unless noted). Existing:
func (guard) Operators() []string          // frozen; also mirrored in table
func (table) Operators() []string          // frozen mirror, same 8 members
func (accessor) Verdicts() []Verdict       // frozen: allow/deny/indeterminate
func (resolve) RefusalKinds() []RefusalKind // frozen: data.escape_class
func (graphlint) Severities() []string     // frozen: blocking/info
func (graphlint) Reasons() []string        // append-only: unprovable-coverage reason
func (graphlint) BlockingCodes() []string  // append-only
func (graphlint) AdvisoryCodes() []string  // growing (amended from closed, A1/S8)
func (table) Categories() []Category       // append-only

// Owed seams this RDR adds (C4 census "Owed" column):
func (respond) Types() []string    // frozen; {"ok"} only — "failed" MUST NOT enter, S6
func (clierr) ExitCodes() []int    // frozen; {0,1,2,3,130}
func (respond) Levels() []string   // frozen; {note,warning}
// GUESS (A5 Pending at time of RDR): signature/location for the CLIError
// code registry (append-only) and the flow-next unknown-reason accessor
// (append-only) — C4 states these two are not mechanical and may land
// as "seam: none (prose-only)" instead; naming them here is a guess of
// the mechanical outcome, not the prose-fallback outcome.
func (clierr) Codes() []string           // GUESS — may not land; A5 Pending
func (graphlint) UnknownReasons() []string // GUESS — may not land; A5 Pending

// package version (internal/version) — unchanged by this RDR, frozen
// field names per C4:
type Info struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}
```

Error modes: `ValidateMode` and verb RunE paths return `*clierr.CLIError`
carrying one of the frozen exit-code classes via `ExitCodeFor`
(0 success/warning, 2 user/internal error, 3 env unavailable, 130
signal). No new error mode is introduced by 0029 itself — it only adds a
field and reclassifies vocabulary tiers; C3's promotion event (`info` →
`blocking`) is an existing exit-2 path, not a new one.

---

## 2. Three most important internal helpers

1. **`clierr` schema-version constant accessor (the "one home").**
   Responsibility: give `respond.OK`/`respond.Fail` (and any future
   terminal-record emitter) the single string both records stamp, so a
   literal can never drift between the two packages. This is the
   mechanism C1's whole "one key meaning two things on one wire" defect
   avoidance rests on. GUESS on exact shape — plain exported const vs. a
   zero-arg func — C1 only requires "a single exported constant."

2. **`internal/cli/flow_input.go::planEnvelope` (and its helper
   `readPlan`) — the tolerant decoder.** Responsibility: decode a
   `--plan` document as either the full `ok` envelope or the bare `data`
   object, discriminating shape by `typePresent := len(bytes.TrimSpace(
   env.Type)) > 0` then `!typePresent && env.Code != ""` (quoted
   verbatim from Illustrative Code). This is the one in-repo consumer
   S3 requires survive the new field unmodified — it must NOT gain
   `DisallowUnknownFields`, by C1's inverse-of-input-rule clause.

3. **`internal/graphlint/taxonomy.go::severityFor`.** Responsibility:
   derive a finding's severity from `!IsBlocking`, so introducing a new
   `info`-severity advisory code (C3) needs no new severity-plumbing
   code path — the severity is structural, not per-code special-cased.
   This is the mechanism that makes C3's "introduction MUST NOT alter
   success disposition" hold by construction rather than by per-site
   discipline.

(A fourth, closely tied but not top-3: the CI seam-diff snapshot check
A7 requires — enumerates each tiered seam's members and each finding
code's severity, fails on tree/snapshot divergence. GUESS on its exact
location/format; C3 requires the mechanism but does not name a file.)

---

## 3. Data model (persisted / cross-boundary)

Wire (JSON) shapes — the only "persisted" artifact this RDR concerns is
the emitted envelope, not on-disk state:

```json
// ok envelope (success)
{
  "type": "ok",
  "schema_version": "0.1",
  "data": { "findings": [] },
  "notes": [ /* omitempty */ ],
  "warnings": [ /* omitempty */ ]
}

// refusal (bare CLIError, no wrapper, no "type" key)
{
  "schema_version": "0.1",
  "code": "graph-lint-failed",
  "message": "the model carries blocking graph-lint findings",
  "findings": [
    {"code": "graph-overlap", "severity": "blocking", "rule": "r3"}
  ]
  // param, detail, hint all omitempty
}
```

Field/vocabulary tiers (the data model's stability contract, C2/C4):

| Field / vocabulary | Tier | Notes |
|---|---|---|
| `schema_version` | n/a (versions the schema itself) | string `MAJOR.MINOR`; non-omitempty on both records; top level only |
| envelope `type` | frozen | `{ok}` only; refusal has no `type` key at all |
| exit-code classes | frozen | `{0,1,2,3,130}` |
| stderr advisory `level` | frozen | `{note,warning}` |
| gate `verdict` | frozen | `{allow,deny,indeterminate}` |
| `data.escape_class` | frozen | `resolve.RefusalKinds` |
| `graph-lint-failed` aggregate code | frozen | single const, no set |
| `findings[].operator` | frozen | pinned in both `guard` and `table` |
| CLIError `code` | append-only | registry owed, A5 Pending |
| `flow next` unknown-`reason` | append-only | accessor owed, A5 Pending |
| `findings[].block` | append-only | grew 2→3 members already (JDR 0001 §D12) |
| graph-lint BLOCKING codes | append-only | |
| `internal/table::Categories()` | append-only | |
| `data.escape_class`'s reason sibling, `graph-unprovable-coverage` reason | append-only | |
| graph-lint ADVISORY codes | growing | amended from `closed`-at-4 (A1, S8) |
| `findings[].class`, `data.dispositions` | no tier | author-authored tokens, not CLI vocabulary |
| `findings[].key/literal/dimension` | no tier | free text |
| `version` payload (`version`,`commit`,`date`) | frozen (names only) | values are build identity, untiered |

---

## 4. Top-level pseudo-code of the main operation

```go
// Terminal-record emission path (respond.OK / respond.Fail), the
// operation C1/C2/C3 constrain end to end.

func OK(cmd *cobra.Command, s Success) error {
    s.Type = "ok"
    s.SchemaVersion = clierr.SchemaVersion // C1: same constant both sides

    if ModeOf(cmd) == ModeJSON {
        return writeJSONLine(stdout, s) // schema_version rides top level
    }
    // text mode: render notes/warnings/findings to stderr/stdout
    return renderText(cmd, s)
}

func Fail(cmd *cobra.Command, err *clierr.CLIError) error {
    err.SchemaVersion = clierr.SchemaVersion // C1: same constant, refusal side

    if ModeOf(cmd) == ModeJSON {
        // bare CLIError, no wrapper, no "type" key (0005:C1)
        return writeJSONLine(stdout, err)
    }
    return renderTextError(cmd, err)
}

// Introducing a new lint finding code (C3), the verdict-preserving path:
func registerAdvisoryCode(code string) {
    // C3: new code enters at severity info
    taxonomy.advisoryCodes = append(taxonomy.advisoryCodes, code) // growing (C4)
    // severityFor derives info from !IsBlocking — no branch added here
}

func lintRun(model Model) (*Success, *clierr.CLIError) {
    findings := graphlint.Evaluate(model)
    if hasBlocking(findings) {
        // verdict-changing path: exit 2, envelope shape changes
        return nil, &clierr.CLIError{
            Code:     "graph-lint-failed", // frozen aggregate
            Findings: findings,
        }
    }
    // info-only or clean: verdict untouched (C3, 0006:C17)
    return &Success{Data: struct{ Findings []Finding }{findings}}, nil
}

// Promotion event (C3): moving a code from info to blocking.
func promote(code string) error {
    if !wasAlreadyRefusedUnderAnotherCode(code) {
        // introduction, not promotion — must land at info (C3 exception clause)
        return errIntroductionNotPromotion
    }
    taxonomy.blockingCodes = append(taxonomy.blockingCodes, code) // append-only
    // MUST disclose in release notes naming the code (C3); MUST NOT
    // ship in a patch release. CI seam-diff (A7) forces the snapshot
    // update that is the reviewable disclosure artifact.
    return nil
}
```

GUESS markers used above: the exact identifier for the shared schema
constant; whether the CLIError code registry and flow-next reason
accessor land as real seams or as `seam: none (prose-only)` rows (A5 is
Pending on both in the record); the CI seam-diff snapshot's file
location/format (A7 requires the mechanism, not a path).
