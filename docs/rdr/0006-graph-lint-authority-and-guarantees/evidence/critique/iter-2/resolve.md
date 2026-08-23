Model: claude-opus-5[1m]

# Critique Resolve — RDR 0006, iter-2 (dual-model)

Origin ledger = the reconciled union in `diff.md` (18 A-rows + 2 B-only + 3
divergences, matched by passage anchor). Every entry exits once below.

## Dispositions

- **fixed** — origin: A/C-1 + A/C-11 + B/C-10 (grouping predicate); sections
  touched: `Technical Design` (source-state paragraph), `Contradiction Check`,
  A2. Grounding: RDR 0003 `Final [locked 2026-08-22]` fixes "a row's match
  pattern selects which group the row belongs to"; `internal/resolve/resolve.go::Resolve`
  filters `view.matches(row.Match)` against one concrete view; the
  `0003-0006-0007` cluster gate (F4) scoped A10's discharge to "a one-line
  citation … **not a redefinition**". The draft already carried that correct
  citation paragraph — the offending paragraph immediately after it was the
  defect. Replaced with: membership is the authored match pattern; the
  reachability relation decides only *which* groups are proven.
- **fixed** — origin: B/C-1 (coverage false green) + diff divergence 1; sections:
  `Load-Bearing Decisions` (new *Soundness direction is per invariant*),
  `Normative Contracts` (new exactness clause), desk-trace step 6,
  `Contradiction Check`. Grounding: strong-consult adjudication (`strong-consult.md`)
  — at a merged node the coverage union is a **superset** of each abstracted
  concrete view's, so a real gap can be closed by a row that is not a candidate
  there. Widening the domain of a *universal* claim makes it easier to satisfy.
  The blanket "never a false green" was sound only for existential checks.
  Coverage now certified only at match-key-exact nodes; withheld otherwise.
  **This is the highest-value finding of the pass** — B's framing beat A's.
- **fixed** — origin: A/C-8 + B/C-12 (invariant 2 anti-monotone); section:
  invariant 2. Universal terminal satisfaction is anti-monotone in merging, so
  every converging flow was a structural `graph-dead-end`. Fixed by splitting
  nodes on terminal-participating keys *before* the test — exact, not an
  approximation, so no false green is introduced. The test itself is unchanged.
- **fixed** — origin: A/C-3 (escape rows have no edge); section: reachability
  relation. Grounding: RDR 0002 `Normative Contracts` — an escape row "carries
  neither a write block nor a clear list, normalizes to a row with both empty".
  So an escape rescue is a **self-loop** edge, not a missing one. Modeled as
  such, excluded from invariant 2's progress test and invariant 3's ordinary
  population. This also removes the deadlock A found (0006 prescribed a cure
  RDR 0002 forbids).
- **fixed** — origin: A/C-4 + B/C-3 (envelope depth); sections: `Technical
  Design` envelope paragraph, `Normative Contracts`, `Wire / byte format`,
  desk-trace step 8, Infrastructure Audit. Grounding:
  `clierr.go::EmitJSON` does `json.Marshal(e)` on the `*CLIError`;
  `docs/cli-output-contract.md` states "`{"type":"failed", ...}` is not emitted
  by `Fail`". Findings are a **top-level sibling of `code`** (`.findings`), not
  `.error.findings`. Lint MUST NOT restructure the shipped envelope.
- **fixed** — origin: A/C-16 + B/C-4 (the two omitempty/text halves); same
  sections. Grounding: `respond.go::Success` — `Data any \`json:"data,omitempty"\``,
  so the empty-list receipt needs `Data` assigned a non-nil struct; the `OK`
  text branch renders only `Notes`/`Warnings` and `EmitText` renders no
  findings, so text-mode enumeration is a **booked extension**, not reuse. The
  audit row was corrected from "Reuse" to "Extend".
- **fixed** — origin: A/C-7 + B/C-5(a) (CI job collision); sections: A5 gate
  target, scenario 7, Assumption Verification. Grounding:
  `.github/workflows/ci.yml::jobs.lint` (`name: Lint`) is the shipped
  golangci-lint job. Renamed the new job `graph-lint`; the existing job must
  survive untouched.
- **fixed** — origin: A/C-6 + B/C-5(b) (no checked-in model) → **new A8**
  (`Source Search`, Pending). Grounding: repo's only `.toml` files are
  `.roborev.toml`, `.kata.toml`, and RDR spike fixtures, which A5 excludes by
  name. Booked rather than asserted; A5 cannot flip until A8 lands.
- **fixed** — origin: B/C-5(c) (`--flow` discovery) → **new A9** (`Peer RDR`,
  Pending). Grounding: RDR 0005's audit defers transition-model config
  discovery. Only `--model <path>` is instantiable today.
- **fixed** — origin: A/C-9 (write-replaces semantics) → **new A10** (`Peer
  RDR`, Pending). Grounding: RDR 0002 states only that absence from write/clear
  "MUST NOT imply deletion" — no replace-vs-accumulate rule; RDR 0003 owns the
  marker, not the update rule. Invariant 5 drops the path-accumulation check on
  this semantics, so it is booked, not assumed.
- **fixed** — origin: A/C-10 (root violates always-present); section: invariant
  5's sibling clause. Resolved *toward* the check rather than exempting the
  root: an always-present owned key the `initial` table omits **is** absent at
  the root, which is the defect the code exists for. A6's requested shape must
  cover every always-present owned key; the finding is reported against the
  `initial` declaration so the diagnostic names the authored site.
- **fixed** — origin: A/C-18 (invariant 6 read-set); section: invariant 6.
  Grounding: RDR 0007 (`Final`) narrows `Row.RequiresOwned` to post-guard write
  dependencies ("guard-input decidability is not `RequiresOwned`'s job"); RDR
  0002 fixes it as write-block + clear-list keys. Reading it here would check
  writes-before-writes and never fire. Lint now derives the read-set from the
  row's own match + guard atoms, booked on the minimum input contract.
- **fixed** — origin: A/C-5 (`Finding` type ownership self-contradiction);
  section: `Type ownership`. `Block` is RDR 0002 vocabulary and the failure
  class is RDR 0001's `RefusalKind`, so importing either would create the
  dependency the clause forbids. Resolved by declaring those fields
  **string-typed and meaningless to `clierr`**, with the producing package
  owning the vocabulary — and stating the trade explicitly: the type system does
  not enforce "MUST carry", so MVV scenario 8 does.
- **fixed** — origin: A/C-14 (`graph-unprovable-coverage` overloaded); section:
  finding codes. One code, now four triggers (the exactness gate added a
  fourth), each with a different remedy. Added a closed `reason` discriminator
  (`dimension-not-finite`, `tag-not-single-valued`, `row-can-refuse`,
  `node-not-exact`) with the remedy per row, and required the message to state
  that the last two are **not** model defects.
- **fixed** — origin: B/C-8 (vacuous `exists` report); sections: finding codes,
  advisory-tier clause (both sites), disposition table. RDR 0003 requires the
  report; the closed three-member tier had no code for it. Added
  `graph-vacuous-atom` (info) as the fourth member.
- **fixed** — origin: A/C-12 (per-class coverage has no runtime counterpart);
  sections: escape class-scoping clause, scenario 19. Partly upheld: each class
  *is* separately reachable (`escapeOrRefuse` selects on the kind that
  occurred), so per-class union stands. But the `ambiguous_match` arm is
  unreachable for an overlap-free group, and demanding an escape row there would
  mint one `graph-unreachable-rule` then flags. Now checked only for a group
  carrying `graph-overlap`; vacuously closed otherwise.
- **fixed** — origin: A/C-13 (unenforceable resolver-use clause); section: the
  blocking-gate normative clause. The Approach text was already honest ("at the
  merge boundary, not at resolve time"); the contract overstated it. Aligned:
  the clause now names the merge boundary as the enforcement site and says
  plainly that a locally-edited model can still be resolved against.
- **fixed** — origin: A/C-15 + B/C-9 (unmeasured false-positive rate);
  section: `Consequences`. The three *structural* generators the critiques
  named are now closed by construction (split-node dead-end, exactness-gated
  coverage, `initial`-anchored always-present), and the general shape is stated:
  where the abstraction is imprecise, lint **withholds** rather than accuses.
  The rate claim was replaced with an MVV measurement (scenario 22) and a
  mechanical trigger for the successor RDR. No waiver mechanism added.
- **fixed** — origin: A/C-2 + B/C-2 + B/C-11 (unbounded node count);
  section: `Performance Expectations`, desk-trace step 2. The product bound
  never bounded the node count. Added a published **node ceiling** with
  `graph-product-too-large` on the traversal. B/C-2's circular-merge-trigger
  point is subsumed: with membership now syntactic, the group count is bounded
  by the authored rule count regardless of the lattice, so only the reachability
  filter scales with nodes.
- **fixed** — origin: A/C-17 (Draft closing a locked peer's records); section:
  A2 Evidence, Assumption Verification. A2 claimed A17/A19/A20 "close by
  citation". Withdrawn: this draft supplies their **producer**; each closes on
  the peer as a route-back, tracked with A7.

## Net-new scope

**Charted: none.** Every disposition traces to a ledger entry. Two candidates
were considered and declined: guard-aware pruning (already named as the
successor in `Consequences`, now with a mechanical trigger) and a suppression/
waiver channel (explicitly rejected by the RDR's authority design).

## Needs verification

- **A8** (`Source Search`, Pending) — **new**. No checked-in transition model
  exists; A5's oracle has no subject until one is authored and homed. Blocks
  scenario 7 and 22.
- **A9** (`Peer RDR`, Pending) — **new**. `--flow` depends on RDR 0005 config
  discovery that RDR 0005 defers.
- **A10** (`Peer RDR`, Pending) — **new**. Write-replaces semantics for
  single-valued tags is stated by no peer; invariant 5's per-row reading depends
  on it. RDR 0002 is `Draft`, so a scheduled edit, not a route-back.
- **A5** (Pending, MVV Test) — gate target changed (`graph-lint` job, not
  `lint`); still Pending, now also gated on A8.
- **A6** (Pending, Peer RDR) — unchanged in status, but its requested shape
  **grew**: the `initial` table must now cover every always-present owned key
  (from the C-10 fix). RDR 0002 must carry that.
- **A7** (Pending, Peer RDR) — unchanged; still a route-back on locked RDR 0003.
- **A2** (`Verified`) — **retained**, and now genuinely earned: its A10-closure
  claim rests on the citation paragraph rather than the deleted redefinition.
  Its A17/A19/A20 over-claim was withdrawn rather than the record downgraded.
- **New load-bearing claims added by this pass**: the per-invariant soundness
  split (a derivation, stated inline and asserted by scenario 20); node-splitting
  exactness for invariant 2 (scenario 21); the published node ceiling (an
  implementation constant); escape rows as self-loop edges (grounded in RDR
  0002's normalization clause).

## Tiebreakers

- **One, escalated and resolved by the human.** The grouping fork did not
  collapse on cluster evidence alone, so §strong-consult ran
  (`strong-consult.md`, `NEEDS_DECISION`). Its technical verdict was decisive
  (0006 moves, 0003 stays locked; coverage is a false green), but the remedy
  read as substantive, so it went to the human with a non-normative worked
  example. The human directed applying the reconciliation. The cluster-gate
  prior art (F4: "a one-line citation … not a redefinition") is what reduced it
  from a redesign to a paragraph deletion.

## Review gate

Edits carry no change-history narration; rationale is stated as decisions.
Every fix is grounded in `main` source, locked peer text, or the cluster gate
record, quoted at the point of use. Findings that would have required
re-litigating settled calls were not raised. The needs-verification list is
honest: three new Pending assumptions, one widened requested shape, one
withdrawn over-claim. Mechanical gate: no template brackets or draft
placeholders survive; A1–A10 sequential; 22 scenarios.

**Mini-checks**: not re-run — the grounding pass discharged the cue read and its
two fired tables (`disposition`, `desk trace`) persist in the draft. Both were
**updated** by this pass rather than re-derived: the disposition table gained
five rows and the desk trace's steps 2, 3, 6, and 8 were corrected. No fix
introduced a new cue.
