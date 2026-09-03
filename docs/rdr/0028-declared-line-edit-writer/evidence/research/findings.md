Model: claude-opus-5

# cli/0028 Stage 4 — accepted citations and rejected branches

Cold-path resolve, 2026-09-03. Nine Critical Assumptions, all Pending at entry.
Method split: five Source Search (own-code), three Spike (live target), one
Peer RDR (JDR 0003 §D1).

## Accepted citations — own-code Source Search (A1, A2, A6, A7, A9)

| A | Anchor | What the source shows |
| --- | --- | --- |
| A1 | `internal/accessor/binding.go::WriteBinding` / `::ReadBinding`, `internal/accessor/model.go::Artifact` | `Apply(ctx, art Artifact, planned []resolve.Tag)` and `Read(ctx, art Artifact, requested []string)` take the SAME `art`; `Artifact` is `{Role, Path}` — no tag/context field. The seam extension is real and unclaimed; leaning a context map on `Artifact` reaches both carriers at one site. |
| A2 | `internal/accessor/executor.go::Executor.Write` | `readerFor(def.Accessor.Role)` precedes `binding.Apply(applyCtx, art, …)` in the same function — the ORDER C3 needs. |
| A6 | `internal/table/load.go::carrierDefect`, `internal/cli/flowbind/registry.go::commandBacked` | Both named exceptions confirmed. `checkAccessorBindings` is carrier-agnostic (arity only). No accessor dump/marshal/normalize round-trip exists, so A6's stated failure mode has no live target. |
| A7 | `internal/cli/flow.go::registerTagFlag`, `internal/cli/flow_input.go::parseTags` | `--tag` registered as "observed tag, as name=value (repeatable)"; `parseTags` refuses an owned key with "owned state is read from the declared read accessors, never supplied by the caller". Observed-context-only is enforced, not merely documented. |
| A9 | `internal/cli/cmdbind/cmdbind.go::substitute` | `if el != ArtifactPlaceholder { out = append(out, el); continue }` — whole-element equality, never substring-scanned. One more vocabulary member is one more branch at the same site. |

### A6 — one site the RDR does not name

`internal/table/load.go` field copy from `sourceAcc` into `Accessor`:
`if a.Path != nil {…}` then `if a.Command != nil {…}`. This is ADDITIVE
per-field copying, not an exclusionary two-carrier branch — a third carrier
needs one parallel `if a.Edit != nil {…}` line, not a rewrite. Recorded on A6's
Evidence line so the implementer does not rediscover it.

### A2 — the nuance that survives verification

The reader BINDING is resolved before `Apply`, and `cmdbind.Config.AllowCommands`
is baked in at construction (`internal/cli/flowbind/registry.go`). But today's
actual gate check lives inside `cmdbind`'s `spawn`, invoked during
Read/Apply execution — not pre-emptively at the resolution point. So C3's
pre-mutation refusal is a NEW check at an EXISTING site (the gate state is
inspectable there), not a re-ordering. A2's "a check at an existing site, not a
new ordering" is correct; the check itself still has to be written.

## Accepted citations — Spike (A3, A4, A5)

See `../spikes/` for commands, programs, and raw output.

| A | Artifact | Verdict |
| --- | --- | --- |
| A3 | `a3-anchor-regex.{go,md}`, `a3-run.out` | RE2 sufficient, zero lookaround (confirmed via `syntax.Parse`). Exactly-one on every Status-bearing record; no decoy scores ≥2. Two gaps below. |
| A4 | `a4-readback-qualifier.md`, `a4-run.out` | Reader separates the qualifier: `status=Final` / `{"value":"Final"}`, byte-equal to a planned `Final`; qualifier rides `status_form` / `status.qualifier`. |
| A5 | `a5-stage-and-rename.{go,md}`, `a5-run.out` | Stage-and-rename is git-clean for content; mode-copy nuance below. |

## Rejected branches / corrections owed to the record

1. **`--facts-json` does not exist** (A4 spike). The real flags are `-tags`,
   `-json`, `-filter`. The MVV fixture's command reader argv
   `["rdr","status","--facts-json","{artifact}"]` cites a nonexistent flag.
   Correction owed at 0028's Approach and MVV.

2. **Unlinked README rows select ZERO lines** (A3 spike). Rows `| 0001 | … |`
   and `| 0006 | … |` carry no `[NNNN](…)` link, so the per-id anchor
   `^\| \[<id>\]\([^)]+\)…` matches nothing → `edit_anchor_unmatched` on real
   README data. Not a regex defect — a corpus-shape mismatch. Author's call.

3. **A wrapped bracket qualifier truncates silently** (A3 spike). A Status
   qualifier that wraps to a continuation line captures only the first line's
   fragment; `${N}` then preserves only part of it. C3 is per-line by contract,
   so this is a SILENT wrong-value write, not a refusal — the one finding that
   escapes both the exactly-one rule and the re-anchor rule. Author's call.
   (The re-anchor rule in C3 does NOT catch it: the rewritten line still
   matches its own anchor.)

4. **A5's git-cleanliness argument is narrower than its prose.** Git tracks the
   executable bit only (100644 vs 100755), not full POSIX perms. For the
   non-executable markdown this RDR targets, `flowbind.go::save`'s fixed 0600 is
   ALREADY git-silent. Mode-copy is still strictly better — correct for an
   executable target, and preserves perms for non-git consumers — but the
   "spurious mode-change diff" motivation only materializes on an executable
   original. A5's claim holds; its stated reason does not, for the named case.

5. **Consequences A5 settles as stated losses, not blockers:** the inode always
   changes and a hard link silently orphans (stale content); xattrs are always
   lost; a symlinked target is correct ONLY if resolved first (renaming onto the
   symlink path replaces the symlink with a regular file — which is exactly what
   C3's "symlinks resolved" buys).

## Reuse audit — no capability already exists

Against the `{RDR_ENV}` reuse-audit paths plus `internal/accessor/`,
`internal/table/`, `internal/cli/flowbind/`, `internal/cli/cmdbind/`.

| Behavior | Verdict |
| --- | --- |
| Line-oriented text editing | ABSENT — no site |
| Atomic staged write | PARTIAL — `flowbind.go::save` is the only one; JSON-encoder-bound (`store = map[string]string` through `clierr.WriteJSONLine`), hardcodes 0600. Discipline reusable, function not. |
| RE2 line selection | ABSENT — all `regexp.` use is in tests; `load.go`'s "anchor" is a rule-id block parser, unrelated |
| Template segmentation | PARTIAL — `commandPlaceholders` is a closed vocabulary but whole-element only; no literal\|group\|placeholder machinery to reuse |
| File mode preservation | ABSENT — no production `os.Stat`/`FileMode` read-and-restore anywhere |
| Validate-then-apply on file CONTENT | ABSENT — tag-level three-phase discipline exists; nothing validates file bytes pre-write |
| Generic text-file writer binding | ABSENT — `flowbind.go::Writer` is JSON-artifact-specific |

Third-party: `go.mod` carries only go-toml/v2, cobra, pflag. No `renameio`,
no atomic-file or line-editing package. Nothing to adopt.

**No reuse finding forces a return to Stage 2/3.** The Existing Infrastructure
Audit's "Reuse pattern; new binding" and "Extend or sibling" dispositions are
confirmed by the audit rather than contradicted.

`docs/cli-output-contract.md` already fixes the refusal envelope C3's refusals
must honor: scalar failures carry `param`, aggregate failures carry `findings[]`
regardless of runtime cardinality, and environment-reachability failures map to
exit 3 vs a request/content defect's exit 2. C3's multi-rule refusals
(unmatched / ambiguous / colliding) must pick a framing under that contract.

## Peer RDR (A8)

`docs/jdr/0003-accessor-binding-seam.md` §D1, `state: settled`,
`cluster: 0004, 0025, 0026, 0027, 0028`. Its Resolved paragraph is "three rules,
no new field, no new class" and it names `0028:A8` explicitly: "`0028:A8` is
confirmed as a constraint of this entry, and its If-wrong 'owned here' is void —
a typed field on the carrier is this registry's, not 0028's."

All three A8 sub-claims map to §D1: applied sense rides `Applied()` (rule 3),
sub-reason rides `Detail` (rule 1), no new refusal class (the option-(c)
analysis). cli/0026:C1 uses the same class for the post-run "ran, output
unproven" case with `Applied()` true; 0028's C3/A8 use it pre-mutation with
`Applied()` false. Same rules, distinguished by `Applied()` exactly as rule 3
intends — aligned, not conflicting. JC1's recorded disposition holds.

A8's Method field reads "Design Review", which is off the stage's closed
vocabulary; the settled peer decision makes it **Peer RDR**.

---

# Author's-round research (2026-09-03, second pass)

Three questions were taken to precedent, prior art, and the live tool before
being put to the author. Two dissolved; one is a data fix.

## Q1 — the "wrapped qualifier" is NOT a defect. Premise refuted.

The A3 spike reported a silent truncation on a wrapped Status qualifier. **This
was a measurement artifact of the spike harness, not a property of the design.**

Verified directly (`/tmp/q1probe`, the real engine fixture
`rdr/tools/rdr/testdata/status/records/0021-cache-warmup-order.md`, whose Status
genuinely wraps):

- Pre-edit line 1: `- **Status**: Draft [joint decision → 0022-cache-metrics-surface §`
- Pre-edit line 2: `  A1: whether warm-up counts toward the hit-rate metric]`
- Rule: anchor `^- \*\*Status\*\*: (\S+)(.*)$`, replace `- **Status**: {status}${2}`
- Anchor selects **1** line (exactly-one rule PASSES)
- Post-edit line 1: `- **Status**: Final [joint decision → 0022-cache-metrics-surface §`
- `diff` is **exactly one line**; the continuation is untouched.
- Reader after the edit: `{"name":"status","value":"Final"}`,
  `status_form=joint-decision`, and `status.qualifier` =
  `"joint decision → 0022-cache-metrics-surface § A1: whether warm-up counts toward the hit-rate metric"` — **fully intact**.

`${2}` captures the trailing remainder OF THE MATCHED LINE and re-emits it
verbatim. The continuation is a DIFFERENT line, never selected, and is therefore
preserved by C3's "preserve every other byte". Nothing is lost.

**Source-level confirmation:** `rdr/tools/rdr/internal/scan/fields.go:146-149`
joins continuation lines into one logical value via `model.ValueContinues`:
```go
val := m[3]
for j := i + 1; j <= len(d.lines) && model.ValueContinues(d.lines[j-1]); j++ {
    val += " " + strings.TrimSpace(d.lines[j-1])
    f.LineEnd = j
}
```
The reader is structurally aware, not line-scoped. This also refutes the
precedent sub-agent's "read-back is vacuously satisfied against a corrupted
baseline" argument, which assumed a line-scoped reader that does not exist.

Adversarial shapes checked, both safe by construction:
- A pinned `^- \*\*Status\*\*: ` anchor never over-selects the continuation line.
- A record whose VALUE sits on the continuation (`- **Status**:` then `  Draft`)
  selects ZERO lines → refuses `edit_anchor_unmatched`, exactly C3's designed
  answer for a stale model.

**Disposition: no contract change. No new reason token. No new Failure Mode.**
C2's existing "out of scope for admission … the reader's read-back is the check"
line already governs, and the case it was worried about does not arise.

## Q2 — unlinked README rows are DRIFT (a pre-existing latent bug)

Three independent lines of evidence:

1. `rdr/models/rdr-write.toml` is the ONLY code that writes or locates index
   rows. `readme-add` emits ``| [NNNN](NNNN-slug.md) | <Title> | <Status> | <Priority> |``
   and `readme-flip-*` locates "the `| [NNNN](` row of the `## Index` table".
   **No rule variant exists for a bare `| NNNN |` row.**
2. Kata `zyh0` — cli/0028's own origin issue — uses
   `anchor = '^\| \[{tag.record_id}\]\('` as its worked example. The design
   assumed the linked form from the start.
3. Both bare rows date to `adc3910` "docs(rdr): add RDR index README", the
   README's founding commit; the other 26 rows were added later and all follow
   the convention. `git log -S` shows neither row was ever touched again.
   Both record files exist under fully linkable slugs.

**Consequence beyond this RDR:** the existing `readme-flip-*` rules already
cannot flip 0001/0006's status cell. That is a live latent bug in the RDR
tooling, independent of cli/0028 — this record merely surfaced it.

**Disposition: normalize the two rows (a data fix, not a contract change).**
`edit_anchor_unmatched` on a bare row is C3 behaving correctly — "an unmatched
anchor is a stale model, not a missing line".

## Q3 — `--facts-json` is a typo, not a planned surface

`rg` over the whole corpus and engine: `--facts-json` occurs ONLY at 0028's
lines 349 and 681. It is nowhere in `rdr`'s source, models, or any other record.
`rdr status -help` shows the real composition: `-json` ("emit the fact vector as
JSON") and `-filter` ("comma-separated fact names to keep").

Verified working against a record by artifact path:
```
$ rdr status -json -filter status <record.md>
{"facts":[{"name":"status","kind":"enum","value":"Final"}],"record":"0021","schema":"2"}
```
`-tags` also works but is a poorer reader: its output interleaves `--tag`
separator lines (it is resolver argv, not a value stream).

**Disposition: correct both sites to `["rdr","status","-json","-filter","status","{artifact}"]`.**

## Prior art — the design is stricter than the field, not weaker

| Tool | Multi-match | Zero-match | Wrapped value |
| --- | --- | --- | --- |
| Ansible `lineinfile` | replaces the LAST match, never refuses | INSERTS (default EOF) | per-line; continuation silently orphaned, no guard |
| Puppet `file_line` | ERRORS unless `multiple => true` | appends by default (`append_on_no_match`) | no guard |
| `sed` | per-line by architecture | n/a | pattern space is one line; `N`/slurp is manual opt-in |
| Ansible `replace` | multiline via opt-in `(?s)` DOTALL | n/a | no truncation detection |
| cli/0028 | REFUSES `edit_anchor_ambiguous` | REFUSES `edit_anchor_unmatched` | value-with-terminator refuses pre-mutation |

**No surveyed tool detects or refuses a wrapped-value truncation.** 0028's
exactly-one-or-refuse and its pre-mutation terminator check are guards ADDED
beyond prior art. Puppet is nearest (errors on ambiguity) but still appends on
no-match; Ansible supports neither half of 0028's posture.

Per-entry atomicity (C3: all rules land in one write or none) is likewise
STRONGER than Ansible (task-by-task, no rollback), Puppet (resource-by-resource
best-effort) and Terraform (explicitly "not transactional"). 0028's boundary is
narrower than a whole run, which is what makes it achievable.

No Go prior art for declarative per-line anchor+replace exists in the langref
checkout set — a negative finding, not a contradiction.

## Incidental defect found (NOT cli/0028's — kata material)

`rdr status` panics with a nil dereference at
`rdr/tools/rdr/internal/scan/fields.go:152` (`d.nodeAt(i)` returns nil, then
`node.ID`) on a file that has NO H1/Metadata heading AND a continuation line.
Minimal repro: a file containing exactly `- **Status**:\n  Draft\n`.
A well-formed record with an empty metadata value is handled fine. Pre-existing
engine bug, unrelated to this RDR's seam.
