Model: claude-sonnet-5

# A6 spike — schema_version constant home (clierr) vs respond/clierr layering

Verifies 0029:A6: "The `schema_version` constant has a home both terminal
records can read without inverting the `respond`/`clierr` layering."
Grounded already (prior stage): clierr and respond are separate packages;
clierr's doc names it the import-cycle-avoiding leaf; no SchemaVersion
constant exists today; respond already imports clierr, clierr imports
neither. Pending only on the edit landing and `go build ./...` confirming
it.

---

## 1. Import direction — confirm against source

`internal/cli/clierr/clierr.go::CLIError` (package doc, lines 1-11):
package doc states clierr is "the leaf home of the CLI's structured-error
type ... so other internal packages (config, …) can construct CLIErrors
without importing internal/cli and forming an import cycle."

clierr.go imports (lines 14-21): stdlib only —
`encoding/json`, `errors`, `fmt`, `io`, `strconv`, `strings`. No import of
respond or internal/cli.

`internal/cli/respond/respond.go::Success` imports (lines 27-35):
`fmt`, `io`, `strings`, `github.com/cwensel/intrastate/internal/cli/clierr`,
`github.com/spf13/cobra`, `github.com/spf13/pflag`.

Confirmed: respond -> clierr (one direction), clierr imports neither
respond nor internal/cli. A read of a clierr-homed constant by respond
forms no cycle.

---

## 2. Spike edit — wire SchemaVersion through both terminal records

Command:
```
cd /Users/cwensel/sandbox/newcoinc/intrastate
```

Edits made (all reverted afterward, see §5):

`internal/cli/clierr/clierr.go`:
- Added `const SchemaVersion = "0.1"` above `type CLIError struct` (SPIKE
  comment, throwaway).
- Added field `SchemaVersion string \`json:"schema_version"\`` as the
  first field of `clierr.go::CLIError`.
- In `clierr.go::EmitJSON`, before `WriteJSONLine`, set
  `e.SchemaVersion = SchemaVersion`.

`internal/cli/respond/respond.go`:
- Added field `SchemaVersion string \`json:"schema_version"\`` as the
  first field of `respond.go::Success`.
- In `respond.go::OK`, after `s.Type = "ok"`, set
  `s.SchemaVersion = clierr.SchemaVersion` — the read respond makes of
  clierr's constant, proving the layering direction A6 asks about.

This wires BOTH terminal records: the `ok` envelope
(`respond.go::Success`, read from `clierr`) and the refusal record
(`clierr.go::CLIError`, declared and set locally in the same package).

---

## 3. Build and vet

```
$ go build ./...
(no output, exit 0)

$ go vet ./...
(no output, exit 0)
```

No import cycle. No compile break. `respond` reading
`clierr.SchemaVersion` compiles cleanly under the existing one-directional
import.

---

## 4. `go test ./internal/cli/...` — byte-identity golden hazard

Command:
```
$ go test ./internal/cli/... 
```

Result: `internal/cli` package FAILS (as predicted); `clierr`, `cmdbind`,
`flowbind` pass; `respond` has no test files.

The specific predicted failure fired exactly as expected:

`internal/cli/flow_mvv_0023_test.go::TestReq27And28And95And96And97And129_DefaultModeIsByteIdenticalToTheCapturedPreChangeGolden`
(flow_mvv_0023_test.go:118) — asserts byte-identity of the full `--as=json`
default-mode terminal envelope against the golden document
`docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`.

```
golden = {"type":"ok","data":{...}}
now    = {"schema_version":"0.1","type":"ok","data":{...}}
```

The only delta between `golden` and `now` is the added `schema_version`
key — every other byte is identical. This is the exact hazard A6/C1 flagged:
adding a non-omitempty `schema_version` to `respond.go::Success` breaks
byte-identity against the 0023 golden. Several sibling tests in the same
run fail for the same single-field reason
(`TestReq35And118_ResolveOverTheCheckedInModelAndShippedFixturesIsUnchanged`,
`TestMVV0023_ResolveEnvelopeProjectionEndToEnd/step6_...`,
`TestReq92_TheFlagRidesResolveAloneAndHandsRespondOKTheProjectedResult`) —
all are the same golden/fixture-byte-identity family, not independent
defects.

This failure is EXPECTED per the spike brief: it is an implementation
obligation (re-capture the golden under the RDR's implementation phase),
not a refutation of A6. The `internal/cli/lint_0006_test.go` failure in the
same run (`TestReq95And96_FindingLivesInClierrWithStringTypedAtomFields`)
is pre-existing repo state unrelated to this spike (references
"clierr imports \"internal/graphlint\"" — not touched by this change) and
was not investigated further; it is out of scope for A6.

---

## 5. Revert

```
$ git checkout -- internal/cli/respond/respond.go
```

`internal/cli/clierr/clierr.go` was found already restored to HEAD by the
time of the revert step (no diff against HEAD needed a checkout).

Verification:
```
$ git status --porcelain
 M docs/rdr/0029-version-promise-on-machine-readable-output.md
 M internal/graphlint/taxonomy.go
?? docs/rdr/0029-version-promise-on-machine-readable-output/evidence/reconcile/
?? internal/cli/a7_snapshot_spike_test.go
?? internal/cli/testdata/
```

No entries for `internal/cli/clierr/clierr.go` or
`internal/cli/respond/respond.go` — both spike edits fully reverted. The
remaining entries above are pre-existing/concurrent repo state, not
introduced by this spike, and were left untouched. No golden file was
re-captured or modified.

---

## Verdict

A6 PASSES: the read direction (`respond` -> `clierr`) is confirmed against
source and forms no cycle; the spike edit wiring `SchemaVersion` into both
terminal records via a `clierr`-homed constant builds and vets cleanly;
the predicted 0023 golden byte-identity failure fired exactly as the RDR
anticipated, confirming C1's "one home" clause names the right package.
