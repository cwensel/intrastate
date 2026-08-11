# Critique — dual-model ledger diff (RDR 0009)

Model: claude-fable-5 (diff/resolve session; reconciles the two critique passes)

Passes diffed:
- `critique.md` — Model: claude-fable-5 (session model), 9 rows (A-C-1…A-C-9)
- `critique-modelB.md` — Model: claude-opus-5 (alt base model, per-spawn model
  override on the sub-agent — both passes ran in fresh contexts), 12 rows
  (B-C-1…B-C-12)

Reconciled by passage anchor (IDs are per-file). Health check: both ledgers
name concrete passages, functions (`ExecuteAndEmit`, `planOf`, `mustResolve`),
and user journeys; origins span §1–§3/premortem/AT. Healthy — no rerun needed.

## Anchor-matched merge

| Merged | A (fable) | B (opus) | Anchor | Signal |
|---|---|---|---|---|
| M-1 | A-C-3 | B-C-1 | §Technical Design zero-`Result` / scenario 7b vs `Result` doc "never both and never neither" | AGREE (strong) |
| M-2 | A-C-4 | B-C-5 | Multi-breach clause: collapse + `Count` | AGREE |
| M-3 | A-C-5 | B-C-10 | A4 *Carried forward* / Prerequisites: `NextTags` safety rests on unimplemented RDR 0004 | AGREE |
| M-4 | A-C-8 | B-C-3 | §Trade-offs "Staged, not immediate": table author gets nothing until RDR 0002 | AGREE |
| M-5 | A-C-9 | B-C-8 | Normative error-mechanics pins (verbatim return, one-level `Unwrap`) | AGREE |
| M-6 | A-C-7 | B-C-6 | CLI clause normative on Pending A9 | AGREE |
| M-7 | A-C-1 | — | Phase 1 vs Phase 2 landing order (suite red mid-plan) | A-unique |
| M-8 | A-C-2 | — | A2/Prereqs/scenario 11 "unwrapped exits 1" — false vs `root.go::ExecuteAndEmit` | A-unique (grounded TRUE against source) |
| M-9 | A-C-6 | — | `CheckValid` name promises more than the one property checked | A-unique |
| M-10 | — | B-C-2 | Scenario 6 Mutant A: frozen suite cannot detect check removal | B-unique |
| M-11 | — | B-C-4 | Phase 2 strips ADV-2b/Fixup-1e override values → failure-message legibility | B-unique |
| M-12 | — | B-C-7 | Placement "cheapness preference" vs normative precedence over `unmodeled_outcome` | B-unique |
| M-13 | — | B-C-9 | Proportionality: multiple load-bearing contracts; Gate unfilled | B-unique |
| M-14 | — | B-C-11 | Dormant-row behavior break justified by temporal A5 | B-unique |
| M-15 | — | B-C-12 | Scenario 5 oracle tied to spike artifact at one revision | B-unique |

Dispositions for M-1…M-15: see `resolution.md`.
