model: claude-opus-5

# Tooling Pass — RDR 0021 (lint's normalized-graph export), iter-2

Run: Stage 7 mechanical pre-sweep, immediately before the Finalization Gate's
written responses. Source: `rdr lint --locking 0021` (exit 0, header
`blocking=0 resolution=0 placeholder=5 advisory=9`) at `./lint.txt`, plus
`inspect --json --filter assumptions,edges`. Findings only; no RDR edits made
by this sweep.

## CHECK 1 — Template section coverage

- C1 §finalization-gate (1063-1183) — `gate:inline` plus five
  `placeholder:survived` blocks (1093-1096, 1102-1114, 1120-1122, 1133-1137,
  1158-1183). ALL five sit inside the gate sub-sections, and every one is
  TEMPLATE.md's bracketed guidance for a response not yet authored. This is the
  pre-lock state of a Draft in its locking pass, not a hollowed section: lint's
  own fix text for `gate:inline` is "move these responses to the record's
  artifacts/gate.md and leave the one-line pointer (a cross-file move; no
  patch)" — i.e. the `--outcome lock` step itself. CHECK 1's "surviving template
  bracket is a BLOCK" clause is scoped to a non-Draft (non-locking) RDR and does
  not fire here. NOT a block; cleared by authoring the four responses to gate.md
  and item 4 in place.
- No `template:missing-section` finding. The spine is complete: 45 sections,
  Problem Statement through References, all authored.
- No surviving `this is a seed skeleton` header. No `## Refinement Context
  (cluster re-entry)` block anywhere in the outline (`status_reentry=false`,
  `reentry_target=none`) — no 07.1 re-entry note to clear.
- No `scaffold:row`, no `contract:template-example`.

## CHECK 2 — Method label vocabulary

PASS. Eight assumption rows, every `method.off_vocabulary[]` empty
(`ca_off_vocabulary=0`, `ca_off_vocabulary_ids=[]`). Members in use: `Spike`
(A1, A2, A5, A6), `Source Search` (A3, A7, A8), `Peer RDR` (A4) — all
sanctioned. No Evidence Record lacks a Method field.

## CHECK 3 — Source Search self-reference

PASS. The three `Source Search` rows are A3, A7, A8. Every Evidence anchor on
each resolves into the product tree (`internal/graphlint/reach.go::Reach`,
`::Node`, `::reach`, `internal/guard/declaration.go::AssignmentCount`,
`internal/guard/product.go::Groups`, `internal/graphlint/analysis.go::newAnalysis`,
`internal/graphlint/engine.go::Run`, `internal/graphlint/taxonomy.go::CodeProductTooLarge`).
None resolves to the record file or to anything under this RDR's artifact dir.
No regression.

## CHECK 4 — Docs Only on load-bearing claims

PASS, vacuously. No assumption carries `method.members == ["Docs Only"]`; the
label appears nowhere in the record.

## CHECK 5 — Symbol resolution of Source Search / Spike anchors

PASS. Every `kind == "source-anchor"` edge reports `resolved: true` — 22 of
them, across the CA Evidence fields and the body. None `false`, none ABSENT
(so `--repo` was supplied and the lookups genuinely ran). No blocking
`edge:unresolved` in lint. A5's Evidence is written `clierr.go:174::WriteJSONLine`
— it carries a symbol and resolves, so the stale-line-number component is a
documented NON-finding.

## CHECK 6 — Status consistency

PASS. `metadata[]` Status = `Draft`, no qualifier, `status_form=none`. All eight
assumptions `status.value = Verified` (`ca=all-terminal`, `ca_pending=0`,
`ca_unverified=0`), so there is no Pending/Unverified property relied on as
settled fact. No checklist-vs-gate disagreement is possible: the gate responses
are unwritten template brackets, carrying no status claim to contradict.

## CHECK 9 — Evidence-field budget

PASS. No `evidence:over-budget` finding. (Profile is `large`, not
`foundational`, so this check would be advisory regardless — but there is
nothing to report.)

## CHECK 10 — Linking

PASS. No `label:contracts` / `label:contracts-required` (C1-C5 are labelled),
no `peer-evidence:no-element` (A4's `Method: Peer RDR` cites element-level, not
a bare record), no `edge:unresolved`, no `edge:unresolved-terminal`.

## Advisory — `prose:exactness` (398, 399, 403)

Three hits on "canonical" inside C2: "canonical row order", "canonical atom
order", "canonical node key". The finding's own fix offers discharge by Evidence
Record OR "coverage by the Minimum Viable Validation". Covered by the latter:
MVV step 2 asserts two invocations are byte-identical and carry every C2 field;
RT2 states `export ∘ export` byte identity; S1 pins determinism under map-seed
variation, S2 golden fixtures pin `intrastate.graph/1`, S9 the JSON round-trip.
The determinism checklist the fix text names is answered in C2/C3 — map order
(pre-sorted, no Go map iteration on the wire), encoding (the one shared
non-HTML-escaping encoder `clierr.WriteJSONLine`), empty vs null vs absent
(`[]`/`{}` never `null`; inapplicable members ABSENT), and a version marker
(`schema`, `intrastate.graph/1`). Advisory, discharged, no edit owed.

## Verdict

PASS — no blocking finding; proceed to the Gate's written responses. The five
`placeholder:survived` blocks and `gate:inline` are the gate text the lock moves
out, and clear as part of the lock itself.
