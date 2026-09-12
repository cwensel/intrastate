# Run plan — rdr-draft-to-lock (cli/0021)

What this run scheduled and decided. Not a position statement — the RDR's
current stage is derived from on-disk evidence (`/rdr-status 0021`).

```
rdr: 0021-lint-normalized-graph-export   profile: large (as read; Draft = provisional)
posture: upfront=false each=false finalize=true   (rule posture-draft-large)
lenses: grounding,3amigo,critique                 (critique runs --auto)
stages: refine -> resolve -> [lenses] -> reconcile -> finalize   stop-after: finalize
```

Entered at Stage 3 (first run; status Draft, form none).

Re-invoked 2026-09-12. Status still Draft, form `none` (not a routed-back
re-entry). Router re-answers `/rdr-refine` as before. The prior run's park was
put to the author, who ruled that its committed refine PASS satisfies the
advance carve-out (evidence/rulings.md, 2026-09-12); this run resumes at
Stage 4. Lens row and posture re-asked on re-entry, both unchanged.

## Ledger

| Stage | Verdict | Blocking | Note |
| --- | --- | --- | --- |
| refine | PASS | no | 5 contradictions resolved, 1 redundancy collapsed; Validation authored (9 scenarios). Commit aed9efe. |
| — | PARKED | yes | Router names `/rdr-refine` again (rule `locate-draft-refine`, guard `ca=all-pending`); refine leaves no done-signal, so the router cannot see commit aed9efe. Advancing to Stage 4 would override `emit.next`. Fork put to user. |
| — | RULED | no | Fork settled by author: accept inherited PASS, enter at Stage 4. See evidence/rulings.md. |
