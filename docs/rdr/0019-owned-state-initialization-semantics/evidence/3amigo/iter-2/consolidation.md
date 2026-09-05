Model: claude-opus-5[1m]

# 3amigo consolidation — 0019, iteration 2 (delta-scoped re-run)

Delta-scoped to the ids iteration 1 touched, checking whether the rewrites
opened gaps. Three isolated personas, no cross-persona visibility.

## Verdicts

- **PM — PASS, zero findings.** All three iteration-1 concerns resolved. It
  independently corroborated the Priority change: `flow_mvv_0011_test.go`
  hand-transcribes two of `models/rdr.toml`'s `[initial]` values in-tree, which
  is the manual duplication the verb removes.
- **Implementer — six of seven deltas verified sound against source**: the
  kind-dispatched encoder is correct AND complete over all scalar kinds; the
  three carrier binding types are exhaustive over the registry's write branches;
  the refusal ordering is implementable as stated; the arity-1 qualifier is
  exact; the code literals are unpinned everywhere; the mini-checks rows agree
  with the rewritten C1. One Medium: S5.
- **QA — three findings**, two of them new gaps the fixture rewrites opened.

## Hotspot

`0019:S5` — raised independently by Implementer and QA. Two isolated personas
converging on one id is the real signal here.

## Findings and disposition

| id | severity | finding | disposition |
|---|---|---|---|
| `0019:S8` | High (QA) | The different-artifact-path reader fixture is unconstructible: the read-back reader is selected by ROLE (`::Registry.readerFor`) and handed the WRITER's `art`; `flowbind::Reader` reads `load(art.Path)` and uses its own `Path` only for reachability. | **fixed** — the command-backed READ accessor (legal: C1's carrier refusal is scoped to WRITE accessors) is promoted from parenthetical aside to the stated fixture. |
| `0019:S9` | High (QA) | The second-role seal yields a TWO-key store (`{Bkey, seal}`), because `Apply` applies planned tags before sealing and an empty plan never reaches `Apply` (`Executor` short-circuits on `len(plan.Writes)==0`). The scenario would then pass vacuously under the very owned-key-count bug it exists to catch. | **fixed** — restored the clear-through-the-sealing-writer route, which iteration 1 over-rejected: `Apply` deletes the key and sets the seal in ONE atomic `save`, so the mutation lands despite the non-zero exit (the seal only short-circuits the verifying re-read). Setup-step exit 3 stated as expected. |
| `0019:S5` | Medium (Impl + QA) | S5 still demanded a verb-level writer-arity fixture in the writer-routing family, which C1's own load-enforcement clause and the new disposition row make unreachable — a test that passes for the wrong reason. | **fixed** — restated as a LOAD refusal (`flow-model-invalid`) whose purpose is pinning that ordering; C1 and the disposition table updated to name the same code. S6 given standalone assertion terms. |
| `0019:S12` | note (QA, non-blocking) | The torn-state setup mechanism was unstated, and the obvious declarative lever does not work (a `-unreachable` writer applies then seals, so B would be seeded). | **fixed** — the scenario now states the direct-artifact-write setup. |

## Confirmed sound (no edit owed)

`0019:S11` (called "the strongest new scenario" — real oracle, real
discriminating control), `0019:S12`'s oracle, `0019:RT3`/`S4` reconciliation,
`0019:MVV` step 10 (grounded field-by-field against `models/rdr.toml`), the
non-normative refusal-code rule (a real oracle exists: `clierr.Finding.Code` on
the wire, so exit group + pairwise distinctness are assertable without pinning a
literal), `0019:A7`'s statement against C1, the S5/S6 both-arms phrasing as an
assertion form.
