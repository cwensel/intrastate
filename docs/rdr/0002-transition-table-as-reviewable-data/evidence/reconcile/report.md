# Reconcile Report

| item | source | disposition | evidence pointer or plan |
| --- | --- | --- | --- |
| Pre-lock needs-verification lists | 1 | VERIFIED | `evidence/grounding/dispositions.md`, `evidence/3amigo/dispositions.md`, and `evidence/critique/dispositions.md` each report `Needs verification: None`; their findings were folded into the current RDR. |
| Critical Assumptions currently pending or unverified | 2 | VERIFIED | RDR Critical Assumptions A1-A7 are terminal: A1/A3/A6 by spike, A2/A7 by design decision, A4 by peer RDR, and A5 by source search. No `Status: Pending` or `Status: Unverified` remains in the assumption records. |
| Named Resolve spike | 3 | VERIFIED | `cd docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes && GOCACHE=/private/tmp/intrastate-rdr0002-gocache GONOSUMDB='*' GOPROXY=off go run . rdr-fixture.toml kata-fixture.toml`; transcript captured in `docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes/output.txt`. |
| Cluster re-entry escape-edge disturbance to A2/A3 | 1, 3, 4 | VERIFIED | The cluster note in the RDR names A2/A3 as disturbed by missing table-modeled escape-edge semantics. Current A2 is a Design Decision covering non-escape exact-one success plus exactly-one modeled escape fallback; current A3 cites the rerun spike. `rdr-fixture.toml` contains `draft-no-match-escape`, and `output.txt` captures it as `kind=escape` with `escape=no_match`. |
| Post-round exactness delta: deterministic expanded-table ordering | 4 | ACCEPTED | Written into Critical Assumption A7 as a Design Decision. Normative Contracts now define the expanded-table value and deterministic ordering by row identity, then predicate/write keys; MVV scenario 5 covers semantically identical fixtures with different TOML key order. |

Verdict: RECONCILED. All open items have terminal dispositions, no blocker was found, and the RDR is ready for Finalize.
