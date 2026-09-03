Model: claude-opus-5

# 3amigo consolidation — cli/0028

Three isolated persona passes (PM, Implementer, QA), each seeing the RDR and its
own persona block only. Hotspots below are a mechanical count over the anchored
ids (`rdr anchors --record 0028` per file, `uniq -c`, `>=2`) — overlap marks a
hotspot passage, not a validated finding, and a single-persona finding is not
thereby weaker.

## Hotspots (ids raised by 2+ personas)

| id | personas | what they independently landed on |
|---|---|---|
| `0028:MVV` | 3 (PM, Impl, QA) | the acceptance pipeline: missing `--tag` on the resolve half (PM-1), the context-map plumb it needs (Impl-2), and two wrong step cross-references (QA-3) |
| `0028:C1` | 3 (PM, Impl, QA) | C1.6's weight (PM-4), C1.3 gate-check site + C1.2 provenance (Impl-1/3), four clauses with no scenario (QA-5/6/7) |
| `0028:S1` | 2 (Impl, QA) | the "neither" arm's wording vs. behaviour (Impl-7); the residue arm has no scenario (QA-7) |
| `0028:S23` | 2 (Impl, QA) | the deadline arm can mint applied-but-unverified (Impl-6); "no subprocess" names no witness (QA-10) |
| `0028:A1` | 2 (PM, Impl) | the seam extension is load-bearing, not incidental (PM-4); it is a two-hop plumb, per-role vs per-invocation unstated (Impl-2) |
| `0028:A7` | 2 (PM, Impl) | verified `--tag` only on `set-state`, never the resolve half (PM-1); `parseTags` refuses two of the three provenances C1.2 admits (Impl-3) |

## Merged findings, severity-ranked

Source key: PM = persona-1-pm.md, IMPL = persona-2-implementer.md, QA = persona-3-qa.md.

### Blocking / high

1. **C1.3's gate pre-check has no reachable site** — IMPL-1. `internal/accessor`
   cannot see `cmdbind.Config.AllowCommands`; `cmdbind` imports `accessor`, so the
   back-import is a cycle. Three candidate sites offered, each a different scope of
   change. Anchors `0028:C1` (C1.3 `read-back:`), `0028:A2`. Blocks MVV step 5, `S25`.
2. **The MVV pipeline cannot bind `{tag.nnnn}` on the resolve half** — PM-1, and
   IMPL-2 from the other side. `--tag` appears only on `set-state`; `flow resolve`
   runs the narrowed readers and parses tags too, so a narrowed C1.6 readme reader
   refuses before `set-state` is reached. Anchors `0028:MVV`, `0028:C1` (C1.6),
   `0028:A7`. Blocks the record's only end-to-end outcome proof. **Hotspot.**
3. **A1's context map is a two-hop plumb, and per-role vs per-invocation is unstated**
   — IMPL-2. `artifactMap()` builds `Artifact{Role,Path}` from `r.artifacts` only and
   never reads `req.observed`. Anchors `0028:A1`, `0028:C1` (C1.2/C1.6), `0028:MVV`.
   **Hotspot.**
4. **C1.2's "any provenance" is refuted by `parseTags`** — IMPL-3. Owned keys and the
   recognized key are refused at parse; only `observed` (and undeclared) bind.
   Anchors `0028:C1` (C1.2/C1.6), `0028:A7`. Blocks `S3`, `S7`, C1.4's
   `edit_anchor_invalid` predicate. **Hotspot.**
5. **S8's pass criterion is false against its own spike** — QA-1. 5 of the 33 named
   fixtures select zero lines (4 postmortems + `BUILD-ORDER.md`); the spike narrows to
   "Status-bearing", `S8` does not. Anchor `0028:S8`. Blocks the corpus regression test
   `F5` names as its only defence.
6. **S12's wrapped fixture is already `Final`** — QA-4. It cannot perform the
   `Draft`→`Final` swap `S12` and MVV step 3 specify; a fitting fixture exists one
   directory over and is not named. Anchor `0028:S12`, `0028:MVV`.

### Major / medium

7. **"Zero caller edits" contradicts the dependency table** — PM-2. The consumer adds
   a verb and types `--tag nnnn=NNNN`. Anchors `0028:§problem-statement`,
   `0028:§capability-dependencies`, `0028:MVV`.
8. **S20's no-write assertion has no obtainable witness, and two wrong MVV refs** —
   QA-3. "wrote byte-identical" and "did not write" are indistinguishable by content;
   `A5`'s inode observable is never named. Anchors `0028:S20`, `0028:C1` (C1.3 `write:`).
9. **S8's cardinalities are snapshot-pinned** — QA-2. 33/26 have already drifted to
   34/28; a count assertion is a scheduled false positive. Anchor `0028:S8`, `0028:MVV`.
10. **`{{`/`}}` escape collides with RE2's repetition quantifier** — IMPL-4. Under
    C1.2 as written `\d{4}` refuses `edit_anchor_invalid`. Anchors `0028:C1` (C1.2),
    `0028:S14`.
11. **C1.3 names no disposition for an unreadable target** — IMPL-5. ENOENT/EISDIR/
    EACCES have no row in the `disposition` table. Anchors `0028:C1` (C1.3),
    `0028:§technical-design`.
12. **The executor's deadline arm can mint applied-but-unverified for a pure-`edit`
    refusal** — IMPL-6, falsifying C1.3 `order:`'s universal claim. Anchors
    `0028:C1` (C1.3 `order:`), `0028:S23`. **Hotspot (S23).**
13. **Four C1 clauses have no scenario** — QA-5/6/7 and QA-8: C1.4 `precedence:`
    (multi-defect entry), C1.5 `one-way:` (two-invocation sequence), C1.1 `in-memory:`
    (residue arm), and F4/F5's post-mutation `read_back_mismatch`. Anchors
    `0028:C1`, `0028:S6`, `0028:S17`, `0028:S1`, `0028:F4`, `0028:F5`. **Hotspot (C1, S1).**
14. **C1.6 is scored nowhere** — PM-4. `§decision-rationale`'s matrix covers the write
    carrier only; no alternative to C1.6 is weighed though `A1`'s `If wrong` names one.
    Anchors `0028:§approach`, `0028:§decision-rationale`, `0028:C1` (C1.6), `0028:A1`.
    **Hotspot (A1).**
15. **The README outcome is proven only against a fixture excluding two live rows** —
    PM-3. No owner, ticket or ordering for normalizing them. Anchors `0028:MVV`,
    `0028:F6`, `0028:§key-discoveries`.

### Minor / low

16. **S23's "no subprocess" names no observation method** — QA-10. **Hotspot (S23).**
17. **S22's negative control cannot be run against the implementation** — QA-9. It is
    a property of `os.Rename` already shown in `A5`'s spike, not an assertion.
18. **C1.1's "the 'neither' arm is unchanged" reads against `S1`** — IMPL-7; also the
    pointer-vs-emptiness discipline for `edit_carrier_conflict`. **Hotspot (S1).**
19. **`§illustrative-code` will not parse as TOML** — IMPL-8. A `;` separator and an
    integer `timeout` against a `*string` field.
20. **`edit` has no analogue for `Writer`'s unreachable-locator seal** — IMPL-9;
    confirm intentional.
21. **C1.2 cites `regexp.Expand` semantics for a hand-rolled parser** — IMPL-10. Bare
    `$1` accepted or refused?
22. **`§decision-rationale`'s process ledger crowds out the rationale** — PM-5. A
    `G-proportionality` answer, not a design defect.

## Isolation and health

Healthy per the lens's Expected signal: every finding anchors a named passage and
names the decision or test it blocks; no persona file references another's output.
Each persona widened and said what sent it — PM into `C1`/dependencies/`A1`, IMPL
three times into `internal/cli` and `executor.go`, QA into `A3`'s spike output and
`§failure-modes`. All three recorded premise-checks that HELD (listed in their own
files) so a later pass does not re-check them.
