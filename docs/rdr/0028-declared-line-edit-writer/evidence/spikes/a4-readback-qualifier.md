Model: claude-sonnet-5

## Scope

Live spike verifying RDR cli/0028 assumption A4: does the consumer's command
reader report the owned `status` value with the bracketed qualifier
separated, so a planned `Final` written by a qualifier-preserving `replace`
reads back byte-equal under contract 0004:C12?

Full raw transcript: `a4-run.out` (same directory).

## Setup note: no qualified Status line currently exists in the consumer tree

`grep -rn '^- \*\*Status\*\*:' docs/rdr --include='*.md'` over
`intrastate/docs/rdr` (the consumer tree named in the task) turns up only
plain `Draft` / `Implemented` values — no record there currently carries a
bracketed qualifier. So the exact scenario A4 describes (a live `intrastate`
record with `Draft [joint decision → 0016]`) does not exist to test in place.

Two other qualified sources were used instead, both legitimate stand-ins:

1. **`rdr` tool's own testdata** at
   `rdr/tools/rdr/testdata/status/records/` — a controlled fixture set built
   specifically to exercise `rdr status`, including records with `Final
   [joint decision → …]`, `Draft [revised from Final …]`, and `Deferred
   [revisit …]` Status lines. This is the most direct test of the reader's
   parsing behavior since it is what the reader's own test suite targets.
2. **`process/rdr/cli/`** — a sibling consumer repo with real records
   carrying live qualifiers (`Final [joint decision → RFD 0004 §3c …]`,
   `Implemented [joint decision → cli/0133 § A4 …]`, `Superseded [→ cli/0107
   …]`), confirming the qualifier forms in testdata match real usage.

## Qualifier forms found

- `Final [joint decision → 0022-cache-metrics-surface § A1: …]` (testdata 0021)
- `Draft [revised from Final 2026-08-10; re-verify A1,A2 — …]` (testdata 0022)
- `Deferred [revisit when the restart budget is measured]` (testdata 0024)
- `Draft [proposed 2026-09-02 under the RFD 0004 push-down …]` (process cli/0107)
- `Superseded [→ cli/0107 — 2026-09-02 RFD 0004 push-down …]` (process cli/0106)
- `Draft [unblocked — cli/0111 is Implemented, so …]` (process cli/0047)
- `Final [joint decision → RFD 0004 §3c screen-grain: …]` (process cli/0143, cli/0112)
- `Superseded [→ cli/0142 — 2026-09-02 RFD 0004 push-down …]` (process cli/0097)
- `Implemented [joint decision → cli/0133 § A4 (…)]` (process cli/0109)
- `Superseded [→ cli/0140` (process cli/0136)

Distinct qualifier FORM tags observed via the reader itself (`status_form=`):
`joint-decision`, `revised-from`, `revisit-when`, `none` (control, unqualified).

## Command surfaces tested

### 1. `rdr status --help` — flag inventory

`--facts-json` **does not exist**. The actual flags are `-tags`, `-json`,
`-filter`, `-checklist`, `-argv`, `-project`, `-records`, `-repo`,
`-template`, `-facts`. This is a discrepancy from the RDR's own text (0028
lines 349 and 681 both cite `rdr status --facts-json` as the MVV command
reader's argv: `["rdr","status","--facts-json","{artifact}"]`). The real
flag that yields a JSON fact vector is `-json` (optionally narrowed with
`-filter status`).

### 2. `rdr status -tags 0021` (testdata, `Final [joint decision → …]`)

Raw output (relevant lines):
```
status=Final
--tag
status_form=joint-decision
```
`od -c` of the `status=Final` line plus its trailing separator:
```
0000000    s   t   a   t   u   s   =   F   i   n   a   l  \n   -   -   t
0000020    a   g  \n
```
No bracket, no qualifier text — `status=Final\n` exactly. Byte-equal to a
planned `Final`.

### 3. `rdr status -json -filter status 0021`

```json
{
  "facts": [
    {
      "name": "status",
      "kind": "enum",
      "value": "Final"
    }
  ],
  "record": "0021",
  "schema": "2"
}
```
`value` is the bare string `"Final"` — qualifier fully separated, not even
carried as a sibling field at this filtered surface. (The qualifier does
appear as `status_form: "joint-decision"` in the unfiltered `-json` output
alongside it as a second fact entry.)

### 4. `rdr inspect -json -filter metadata 0021`

```json
{
  "label": "Status",
  "canonical": "Status",
  "match": "exact",
  "value": "Final [joint decision → 0022-cache-metrics-surface § A1: whether warm-up counts toward the hit-rate metric]",
  "section": "0021:§metadata",
  "line_start": 9,
  "line_end": 10,
  "status": {
    "value": "Final",
    "qualifier": "joint decision → 0022-cache-metrics-surface § A1: whether warm-up counts toward the hit-rate metric",
    "form": "joint-decision",
    "tier": "canonical",
    "raw": "Final [joint decision → 0022-cache-metrics-surface § A1: whether warm-up counts toward the hit-rate metric]"
  }
}
```
Note the outer `value` field (the metadata line's whole rendered value)
*does* carry the full bracketed text — that is the raw markdown-line value,
not a fact. The nested `status` object separates it cleanly:
`status.value = "Final"`, `status.qualifier = "joint decision → …"`,
`status.tier = "canonical"`, `status.raw` = the full original string. This
answers the `value` / `tier` / `raw` triple the task asked for.

### 5. Artifact-PATH invocation (MVV binds `{artifact}` to a path, not NNNN)

```
rdr status -tags /Users/.../testdata/status/records/0021-cache-warmup-order.md
```
produces identical output to the NNNN form: `status=Final`,
`status_form=joint-decision`, etc. Path-based invocation behaves the same as
record-number lookup.

### 6. Control: unqualified record (testdata 0020, Status: `Draft`)

```
rdr status -tags 0020
```
```
status=Draft
status_form=none
```
`rdr inspect -json -filter metadata 0020` shows `status.qualifier` absent/
empty and `status.raw = "Draft"`. Confirms the separation mechanism is
qualifier-conditional, not always-on structural noise — an unqualified
record round-trips as `Draft` with `status_form=none`.

## Verdict

**A4 holds.** The reader (`rdr status`, either `-tags` or `-json`, filtered
or not) reports the `status` fact as the bare tier value (`Final`), with the
qualifier surfaced separately as `status_form` (a coarse category tag) and,
at the `inspect --json --filter metadata` surface, as a full `status.
qualifier` / `status.raw` pair. A planned write of `Final` (produced by a
qualifier-preserving `replace` that keeps the bracket in the file) reads back
**byte-equal** against `status=Final` — confirmed literally via `od -c`:
`status=Final\n`, no trailing bracket bytes.

**Discrepancy for the author**: the RDR's own MVV fixture text (0028:349,
0028:681) names the command reader as `rdr status --facts-json {artifact}`.
That flag does not exist on the built `rdr status` command. The reader
surface that actually does what A4 needs is `rdr status -json` (optionally
`-filter status`), or `-tags` for the argv-tuple form the rest of 0028's
design already assumes elsewhere. This does not refute A4's substance — the
separation behavior A4 depends on is real and verified — but the MVV
fixture's literal argv is wrong and needs correction before it's runnable.
