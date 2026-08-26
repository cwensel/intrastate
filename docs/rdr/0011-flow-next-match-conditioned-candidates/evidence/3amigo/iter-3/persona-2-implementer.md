Model: claude-opus-5[1m]

# 3amigo iter-3 Persona 2 — Implementer (RDR cli/0011)

Read via projector: `elements,outline` (all ids), `edges` (all `source-anchor`), and the record body for `0011:C1`/`C2`/`C3`, the six `D-*`, the mini-checks, Approach/Technical Design, A1–A6/A11–A14, MVV, and S1–S7. WIDENED beyond the owned C/D set into the assumption block (A5, A6, A11, A13), the Existing Infrastructure Audit, `0011:JC1`, and `0011:MVV`/S1–S7 — three of the five findings below are about what the contracts do NOT say, and the silence has no line range, so the anchor is the contract that owes the clause plus the assumption or audit row that should have caught it. Grounded against shipped code: `internal/cli/flow_next.go`, `internal/cli/flow_exec.go`, `internal/cli/flow.go`, `internal/cli/respond/respond.go`, `internal/cli/respond/text.go`, `internal/resolve/resolve.go`, `internal/resolve/guard.go`, `internal/table/model.go`, `internal/table/normalize.go`, `models/rdr.toml`, and the five cited test reads.

---

## S1 — `0011:C1` says the probe's match keys come from "the invoked readers", but reader narrowing does not demand match keys, so an owned match key can be structurally absent and C1 silently degrades to `--all`

**Anchor**: `0011:C1` (the phrase defining the assembled view as "owned tags from the invoked readers plus observed --tag tags"), with `0011:A11` and the Existing Infrastructure Audit row `Reader narrowing | internal/cli/flow_exec.go::invokedReaders | union over all rows for next | Reuse | unchanged in both modes`.

**Ground**: `internal/cli/flow_exec.go::invokedReaders` builds its demand set from two sources only — `row.RequiresOwned` and `guardOwnedKeys(m, row)`. And `guardOwnedKeys` (`flow_exec.go:131-145`) opens with

```go
if atom.Block != table.BlockAll && atom.Block != table.BlockUnless {
    continue
}
```

so a `BlockMatch` atom contributes **nothing** to the demand set. `RequiresOwned` is derived at normalization from the write block and clear list alone (`internal/table/normalize.go::renderWrites` returns `slices.Clone(keys)` over `assignments`, which is populated from `rule.Write` and `rule.Clear` only; `0002:C14`). `table.Tag.Required` is loaded (`internal/table/load.go:253`) and read by nobody — the only `.Required` reads in the CLI are `flow_next.go:178-179` on the candidate's own `Required` slice.

Therefore: a model that matches on owned key `k` but never writes/clears `k` and never names `k` in a guard block demands no reader for `k`. No reader is invoked, `k` is absent from `assembledView`, C1's presence test omits every `match.k` atom from the probe, and **every row is a candidate with `{k, absent}`** — the exact "wall of candidates" the Risks section says the fix removes, arrived at structurally rather than through caller error.

Both spikes miss this. `evidence/spikes/a6-rdr-toml-narrowing.md` works only because `models/rdr.toml`'s `stage` is *also* written by every rule (`[rule.write] stage = "prelocked"` at `models/rdr.toml:214`, etc.) and is co-declared on the same reader as the guard keys `status`/`gate_passed` (`keys = ["stage","status","gate_passed"]`, `models/rdr.toml:76/82`) — narrowing succeeds by co-declaration, not by design. A6's own text asserts `stage` is "`required = true` and served by the invoked `rdr-status` reader", which reads as if `required = true` is doing the work; it is not, it is inert. `evidence/spikes/a5-decision-table-no-tag.md` matches on `status`/`size`, which are **observed** tags supplied by `--tag`, so it never exercises an owned match key either.

**Clarification requested**: does C1's predicate require `invokedReaders` to add match-block owned keys to its demand set, or is the degradation intended and merely undocumented?

**Blocks**: the Phase 1 probe-builder edit, and whether Phase 1 is confined to `flow_next.go` at all. If the answer is "add them", `invokedReaders` changes and the Infrastructure Audit's `Reuse | unchanged in both modes` cell is wrong; it also widens the reader set under the **default** mode but not under `--all` (match plays no part there), which contradicts `0011:C1`'s "The invoked reader set and the assembled view are fixed ONCE per invocation … and are the same under `--all`" and contradicts `0011:JC1`'s "Shared modify-anchor: `internal/cli/flow_exec.go::invokedReaders` (Reuse/unchanged in both)" — the joint-check with cli/0010 was cleared on the premise that this anchor is untouched. If the answer is "leave it", C1 owes a stated limit, because a caller reading the help C3 mandates ("a match key the state does not carry") will read `{k, absent}` as "supply the tag" for a key that is `provenance = "owned"` and therefore **cannot** be supplied via `--tag` — `parseTags` refuses it `flow-tag-owned` (A3). The caller is handed a remedy that is refused when attempted.

**Prevents the test**: MVV 6 / S3's "absent ⇒ candidate with `{key, absent}`. Then add `--tag key=<value>` and assert the list narrows" cannot be written for an owned match key at all — the `--tag` half is refused by `flow-tag-owned`. S3's absent case is only writable over an *observed* match key, so the owned-match-key-absent class ships with no oracle in either direction.

---

## S2 — `0011:C2`'s "filter BlockMatch atoms out of `unknown`" collides with C1's owned-key walk on a key that is both a match key and a written key; the mandated oracle fails on `models/rdr.toml`

**Anchor**: `0011:C2` — "the --all branch MUST filter BlockMatch atoms out of `unknown`, and an oracle MUST assert their absence" — read against `0011:C1`'s enumeration of `unknown`'s three sources ("every match atom over an absent key, every guard atom the kernel could not decide, **every owned key no invoked reader established**, and … every gate id").

**Ground**: `unknown` has two independent walks that can produce the same key. `internal/cli/flow_next.go::summarize` today runs the atom walk (`for _, atom := range row.Atoms`) and then, separately, the owned walk (`for _, key := range row.RequiresOwned`), each guarded by its own `slices.Contains`. In `models/rdr.toml` every rule has `stage` in **both** roles: `[rule.match.stage]` at `:211/:221/:231` and `[rule.write] stage = …` at `:214/:224/:234`, and `renderWrites` puts written keys into `RequiresOwned`. So if `stage` is unestablished, the default mode emits `{stage, absent}` from either walk (dedup on the pair collapses them — fine), but under `--all` C2 says the `BlockMatch` atom is filtered while C1 says the owned key is still reported. Both apply to the same `{stage, absent}` pair.

C2 gives one instruction ("filter BlockMatch atoms out") and one oracle ("assert their absence"). Neither says whether "their absence" means *the atom's contribution* is suppressed or *the pair* is absent from `unknown`. Read as the pair, the mandated `--all` oracle **fails on the very model MVV 4 runs** the moment a reader does not establish `stage`, and — worse — it would be *wrong* to make it pass: suppressing `{stage, absent}` under `--all` would delete a fact 0005 reports today from the owned walk, which C2's own "Everything else 0005:C1 requires of flow next … MUST hold identically in both modes" forbids.

**Clarification requested**: is C2's filter scoped to the atom walk's contribution only, leaving the owned walk untouched under `--all`? And is the mandated oracle therefore an assertion over a match key that is **not** in `RequiresOwned` (which C3's discriminating fixtures do not currently promise)?

**Blocks**: the Phase 2 edit ("route it to a probe with the match pattern omitted and the `BlockMatch` atoms filtered out of `unknown`") — the filter's scope is the whole edit, and the two readings produce different code and different payloads.

**Prevents the test**: S2's "under `--all` a row with an absent match key reports NO match entry (the C2 departure)". Over any fixture whose match key is also written, that oracle is unwritable as stated: the row reports the pair regardless, from the other walk. C3's fixture list ("a row whose match key is ABSENT (candidate in both modes, `{key, absent}` by default and no match entry under `--all`)") does not say whether that row writes the key, and `internal/table/normalize.go`'s `malformed_rule_shape` refusal (A5, limit (i): "an ordinary rule with no write block is refused") means the fixture row must write *something* — C3 must say it writes a *different* key, or the fixture cannot discriminate.

---

## S3 — `0011:D-undecided-reporting-shape` mandates a text-mode rendering (`key (reason)`) that `flow next` has no seam to produce, and producing it would break the "modes cannot disagree" invariant the gateway is built on

**Anchor**: `0011:D-undecided-reporting-shape`, final sentence — "Text mode (`--as=text`) renders the same pairs as `key (reason)` on the candidate".

**Ground**: `flow next` has **no verb-specific text renderer**. `runFlowNext` ends at `respond.OK(cmd, respond.Success{Data: payload})`, and `internal/cli/respond/respond.go::OK`'s text branch does exactly two things: emit findings if the payload satisfies `FindingCarrier`, else call `writeTextPayload(cmd.OutOrStdout(), s.Data)`. `flow_next.go`'s only `FindingCarrier` mention is the closing `var _ = clierr.Finding{}` with the comment "The success payloads here carry none, so this file declares no carrier."

`internal/cli/respond/text.go::writeTextPayload` marshals the payload to canonical JSON, unmarshals to `any`, and emits path-qualified leaf lines through `flatten`. Its own comment states the invariant: "Any other payload is rendered field-for-field from the SAME value the JSON branch marshals, so the two modes cannot disagree about what the run reported (RDR 0005 REQ-7/REQ-9/REQ-11/REQ-120)."

So with `unknown` as `[]{key, reason}`, text mode emits, per entry, two lines of the shape `candidates[0].unknown[0].key: stage` and `candidates[0].unknown[0].reason: absent`. It does **not** emit `stage (absent)`, and there is no hook by which `flow next` alone could make it — `flatten` is generic over decoded JSON with no per-verb dispatch.

**Clarification requested**: is `D-undecided-reporting-shape`'s `key (reason)` (a) normative and therefore requiring either a `flow next` text-render seam or a change to `respond/text.go` (neither of which appears in any Phase, in the Infrastructure Audit, or in C1/C2/C3), or (b) descriptive prose meaning "the pair is visible in text mode, in whatever form the generic flattener gives it"?

**Blocks**: Phase 1's scope. Reading (a) puts a change to the shared `respond` gateway in the build — a surface `0005:C1` governs and this record does not claim to override (its Overrides field names only the `flow next` candidate clause and the `unresolved`→`unknown` rename) — and it deliberately makes the two modes disagree field-for-field, which is what the flattener's comment says REQ-7/9/11/120 forbid. Reading (b) leaves a `D-*` sentence that a reviewer will read as a shipped obligation.

**Prevents the test**: MVV 7 / S2's "`unknown` present as `[]`, not omitted … in both modes" is written against the JSON shape; there is no oracle anywhere in S1–S7 for the text rendering this decision mandates, so under reading (a) a load-bearing decision ships unpinned.

---

## S4 — `0011:C1`'s "restricted to keys present" is stated over atoms, but the probe's `Match` is `[]resolve.Tag` and the restriction has to be applied *inside* the `KernelRow()` conversion; C1 names no seam for that

**Anchor**: `0011:C1` — "a one-row probe table that carries the row's match atoms RESTRICTED to keys present in the CLI's assembled view" — and the Existing Infrastructure Audit row `Per-row kernel probe | internal/cli/flow_next.go::excluded | strips every match atom (probe.Match = nil) … | Extend`, plus `0011:A12`, which reasons entirely in terms of "omitting the absent-key match atoms".

**Ground**: the atom→tag conversion is one-way and lossy of block identity. `internal/table/model.go::Row.KernelRow()` walks `r.Atoms`, and for `a.Block == BlockMatch` emits `resolve.Tag{Key: a.Key, Value: seamValue(a.Literal, r.isSet(a.Key))}` — a `resolve.Tag` has only `Key` and `Value`, no `Block`, no `Operator`. The probe today is built as `probe := row.KernelRow(); probe.Match = nil`. To restrict rather than nil out, the CLI must filter `probe.Match` by `Key` against `view` **after** the conversion — which works, since `Tag.Key` survives — but that is a different operation from "omitting atoms", and it interacts with `0002:C13`'s `in`-expansion: A12 states match-block `in` atoms are "already expanded into per-member `eq` rows at normalization", i.e. the expansion is at **row** level (each expanded row is a separate `table.Row` with its own `Suffix`), so a single row's `probe.Match` holds at most one `Tag` per key. The record asserts this ("restricting the set by key presence is per-atom with no cross-atom interaction") but nowhere states the mechanical form — filter `[]resolve.Tag` by key, versus filter `row.Atoms` and re-derive.

The two are not equivalent for `seamValue`. `seamValue(a.Literal, r.isSet(a.Key))` consults `r.isSet(a.Key)` — a per-**row** set-key list, not a per-atom property. Re-deriving the match tags by filtering `row.Atoms` and rebuilding `resolve.Tag`s in `flow_next.go` would require re-implementing `seamValue`'s set-canonicalization (`internal/table/model.go:313-345`, JSON array, `SetEscapeHTML(false)`, sorted, compacted) in a second place — the exact fork the `authority` mini-check's "one test, two uses — C1 forbids a second presence vocabulary" row guards against for presence, but not for value encoding.

**Clarification requested**: is the restriction to be applied as a post-`KernelRow()` filter over `probe.Match` by `Tag.Key`, and is `seamValue` therefore never re-invoked by `flow_next.go`?

**Blocks**: the Phase 1 probe edit's first line. Choosing the atom-side rebuild forks the set-value encoding and would silently mis-compare a set-kinded match key against the kernel.

---

## S5 — `0011:C3`'s `not-evaluated` re-homing rule is unimplementable as stated for `flow_next_0005_test.go:163`, because the gate-id read is a *membership* assertion whose element type changes shape

**Anchor**: `0011:C3` — "The gate-id reads re-home to the `not-evaluated` reason C1 names; a re-homing that has to invent a reason token is not mechanical and is a defect" — and "THREE of the five discard the comma-ok of a helper (`flow_harness_0005_test.go::stringsAt` …) — `flow_next_0005_test.go:404`, `flow_mvv_0005_test.go:66`, and `flow_adversarial_0005_test.go:446`".

**Ground**: the census is exactly right on the shipped build — I confirmed all five reads and that `flow_next_0005_test.go:163` is the one that already asserts `ok` (`unresolved, ok := stringsAt(c, "unresolved")`), while `:404`, `flow_mvv_0005_test.go:66`, and `flow_adversarial_0005_test.go:446` all use `, _`. `stringsAt` (`flow_harness_0005_test.go:426-444`) returns `nil, false` on a non-array or non-string element — so after the type change from `[]string` to `[]{key,reason}`, **every one of the five** returns `nil, false`, and the three `, _` sites go vacuously green. C3 catches this precisely and requires the ok bool. Good.

What C3 does not say is what `stringsAt` becomes. Four of the five reads are membership tests over a flat `[]string` (`slices.Contains(unresolved, "approval")`, `slices.Contains(unresolved, "flag")`, `containsString(unresolved, "flag")`, `slices.Contains(unresolved, "approval")`); one is a presence-of-key test (`_, present := c["unresolved"]`, `:103`). C3 mandates asserting the ok bool but does not say whether `stringsAt` is (a) replaced at these sites by a new `pairsAt`/`unknownKeys` helper in `flow_harness_0005_test.go`, or (b) left alone with each site projecting. `stringsAt` is shared harness used well beyond these five sites, so (a) is not a local edit and (b) means five hand-rolled projections. C3's "mechanical" framing and its exact-five count both read as though the edit is a rename.

**Clarification requested**: does the re-homing add a harness helper (and if so, is that a sixth touched site outside C3's census), or does each of the five project inline?

**Blocks**: Phase 3's "re-home the 5 shipped reads". Minor, but it is the difference between a rename and a harness change, and C3 explicitly forbids counting anything beyond its five against A4's BREAKS-0 — a new harness helper is neither one of the five nor covered by the production-site list C3 enumerates, so the record has no slot for it.

**Prevents the test**: S6's checkable form — "`go test ./internal/cli` green except C3's five named re-homings" — is stated as a *count*, so a harness addition makes the stated pass condition ambiguous.
