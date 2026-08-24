model: claude-opus-5
lens: repeatability-lite (resolve pass)

# Repeatability-lite Resolve — RDR 0005

Origin ledger = `diff.md`'s R-1..R-4 (first pass; the diff findings *are* the
originating concerns).

## Grounding sources

- **Code on `main`**:
  - `internal/resolve::Table.Revision` — "the opaque caller-supplied table
    revision identity. The kernel carries and compares it; it does not parse
    it" (`internal/resolve/resolve.go:205-207`); echoed on outputs at
    `resolve.go:508,516`. Decides R-1.
  - `internal/cli::NewRootCmd` — `cmd.PersistentFlags().String(respond.FlagName,
    "text", "output mode: text | json")` (`internal/cli/root.go:51`); the default
    already exists and is the root's. Decides R-2.
- **`.rdr/resources.md`**: `docs/cli-output-contract.md` is named the
  authoritative output contract — it states the persistent `--as text|json` flag
  but no default, confirming the default lives in root wiring, not the doc. The
  RDR must cite the surface, not restate a default it does not own.
- **RDR's decided text**: none of R-1..R-4 re-litigates an adjudicated option.
  §D8 says `next`/`resolve` "MAY run declared read accessors" without fixing the
  invoked set (R-3 is a genuine silence, not a re-raise); §D9 fixes gate site
  only for `resolve`'s post-selection case, leaving `next`'s scope open (R-4).

## Dispositions

1. **R-1 `revision` undefined — fixed** (pin). §Technical Design *Model
   selection*. Pinned to the loaded model's own revision, carried verbatim:
   RDR 0002's loader produces it, the CLI echoes it, the kernel treats it as
   opaque; the CLI never derives, hashes, or synthesizes it, and a model
   declaring none renders it empty. Collapsed without a tiebreaker because the
   kernel's own doc comment already fixes the ownership — the RDR had simply
   never restated it. Amendment sweep: the two consumer sites (§Load-Bearing
   Decisions *Identity*, §Finalization Gate *Canonical-form / determinism*) both
   read "the same model revision" and are now well-defined by the new clause; no
   stale wording, no competing definition elsewhere.

2. **R-2 `--as` default unstated — fixed** (cut/cite, not restate). §Technical
   Design, end of the *Envelope* prose. States that the flag's persistence,
   domain, and default stay the root contract's
   (`docs/cli-output-contract.md`, `internal/cli::NewRootCmd`) and this RDR
   inherits it unchanged. Deliberately does **not** write "default text" into
   this RDR: duplicating a value another surface owns is exactly the
   single-source drift the lens's `single-source` disposition guards against.

3. **R-3 reader-invocation set unfixed — fixed** (pin). §Normative Contracts
   (`flow next` / `flow resolve` paragraph) + §Technical Design *State input*.
   Pinned to the narrow reading: only readers serving an owned key some
   candidate row requires; a reader no candidate needs must not run and its
   unbound role must not raise `flow-artifact-missing`, which is now explicitly
   scoped to the invoked set's roles. `read-state` is carved out as the
   exception (every declared reader) because a diagnostic read has no candidate
   set to narrow by — that carve-out is what keeps the pin from breaking
   `read-state`'s existing "invoke only declared read accessors … return per
   reader its declared keys" contract. Collapsed rather than escalated: this
   RDR already fixes the identical question for gates one paragraph later
   ("only the selected row's gates run") on an effect-minimization rationale
   (§Decision Rationale: "running effects before knowing the survivor wastes
   calls"). Readers are effects with the same timeout and exit-3 surface, so the
   broad reading would contradict the RDR's own stated principle.

4. **R-4 `next` gate scope under `--evaluate-gates` unfixed — fixed** (pin).
   §Normative Contracts (`flow next` paragraph). Pinned to: gates of every
   reported candidate and no others; a guard-excluded row is not a candidate so
   its gates do not run; each result rides the candidate carrying it; a deny
   constrains only that candidate and `next` still exits 0, because `next`
   enumerates rather than selects and only `resolve` turns a deny into
   `flow-gate-denied`. The exit-0-on-deny half is the load-bearing part — it is
   what stops `next` from silently acquiring `resolve`'s refusal semantics.

## Mini-checks

Cue read re-run because fix R-4 **added** a cue (an outcome/disposition split
across input classes that did not exist in the draft before this pass).

| Mini-check | Fired | Disposition |
| --- | --- | --- |
| disposition | **yes** (added by R-4) | Table written into §Failure Modes: input class × `next --evaluate-gates` × `resolve`, covering allow / deny / indeterminate / timeout / guard-excluded / not-requested, with the silent-vs-loud line. |
| source-authority census | no | R-1 *removed* the only ≥2-source-of-truth candidate by pinning `revision` to one producer; no fallback/derived/propagate arm added. |
| test-discriminability | no | The MVV rows added are positive/negative pairs by construction (unneeded reader: `next`/`resolve` succeed **and** `read-state` refuses), not absence-of-error oracles. |
| round-trip / fidelity | no (pre-existing cue, not added by this pass) | The RDR already carries §Round-Trip / Inverse Invariants with byte-equality stated and the lossy exemption named; none of R-1..R-4 touched it. |
| desk trace | no | No fix put two normative assertions on one output surface; the disposition table above is the joint statement for the one surface R-4 touched. |

## Needs verification (Stage 6)

- **A7 added, Status: Pending** — RDR 0002's normalized model must expose the
  reader→owned-key mapping and each candidate row's gate list without evaluating
  missing facts. Both R-3 and R-4's pins depend on it. Method: Peer RDR. If only
  the gate list is exposed, the reader-narrowing clause falls back to running
  every declared reader and `flow-artifact-missing` widens accordingly. Added to
  the Reconciliation Report table as PENDING (source: round 5).
- No previously Verified assumption was invalidated. A1–A6 stand: R-1/R-2
  resolved *toward* existing verified surfaces rather than changing them, and
  R-3/R-4 add scope clauses inside A3/A4's already-verified verb/accessor
  boundary rather than moving it.

## MVV / Testing coverage added

- MVV: fixture model must declare one reader no candidate needs with its role
  unbound — `next`/`resolve` succeed without invoking it or raising
  `flow-artifact-missing`, `read-state` invokes it and refuses. Plus
  `next --evaluate-gates` over one gated candidate + one guard-excluded row.
- Testing Strategy scenario 8 added mirroring both.

## Tiebreakers

None. All four collapsed on evidence: R-1 and R-2 on code that already fixes the
ownership, R-3 and R-4 on the RDR's own effect-minimization rationale applied
consistently to readers and to `next`.

## Escalation

Variant stays **lite**. No escalation trigger met — the run produced four GUESS
markers (not suspiciously clean), and they localized onto by-design-open areas.
R-1 is load-bearing but was a definition gap on a token this RDR already owns,
resolvable by decision rather than by more draws.

## Convergence

Converged — no open ledger entries. R-1..R-4 all `fixed`; none dismissed, none
charted (no net-new scope surfaced; every fix landed inside an existing contract
clause this RDR already owns).
